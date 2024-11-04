package comms

import (
	"errors"
	"fmt"
	"github.com/TboOffical/fastscale/orchestrator/utils"
	"net"
	"time"
)

// PORT is the default port to use for communications, change at your own risk
const PORT = 7777

/*
All the communications between the orchestrator DB nodes, and regular nodes are implemented below according to the protocol
*/

// Connection a single connection between two nodes
type Connection struct {
	Online    bool
	TargetIP  string
	SourceKey string
	TargetKey string
	Pipe      chan []byte
}

func NewConnection(targetIP string) Connection {
	return Connection{
		TargetIP: targetIP,
	}
}

func (c *Connection) Connect() error {
	if c.Online {
		//Already online. Do nothing
	}

	//Validate IP Address
	ip := net.ParseIP(c.TargetIP)
	if ip == nil {
		//Invalid IP
		return errors.New("invalid IP Address")
	}

	//Create an error channel to receive errors from the handler
	var err error
	errChan := make(chan error)

	//Start the handler(on a different thread)
	go Handler(c, errChan)

	//Wait for the handler to finish, hopefully never.
	err = <-errChan

	return err
}

func Handler(c *Connection, ec chan error) {
	s, err := net.ResolveUDPAddr("udp4", c.TargetIP+":"+fmt.Sprint(PORT))
	if err != nil {
		ec <- err
		return
	}
	uc, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		ec <- err
		return
	}

	//Send ready for client key message
	_, err = uc.Write(utils.HexStringToBytes("FF"))
	if err != nil {
		ec <- err
		return
	}

	//Variable to hold the byes of the key
	key := make([]byte, 768)

	//If the key is not received in 5 seconds, close the connection
	timer := time.NewTimer(5 * time.Second)

	for {
		select {
		case <-timer.C:
			ec <- errors.New("key timeout; no key received from client")
			break
		}

		pkSize, _, err := uc.ReadFromUDP(key)
		if err != nil {
			ec <- err
			break
		}

		if pkSize != 768 {
			err = uc.Close()
			if err != nil {
				ec <- err
				return
			}
			ec <- errors.New("invalid key size")
			return
		}

	}

}
