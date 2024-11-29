package main

import (
	"errors"
	"fmt"
	"github.com/TboOffical/fastscale/orchestrator/utils"
	"github.com/charmbracelet/log"
	"net"
	"strings"
	"time"
)

/*
Formatting stuff, when we send responses to the terminal, we want the format to be somewhat standardized
*/

func CommandResponseToString(response CommandResponse) string {
	workingOutput := ""

	//start with the date
	now := time.Now()
	dateString := now.Format("2006-01-02T15:04:05")
	workingOutput += dateString + "|"

	//was the command successful
	workingOutput += utils.If(response.Success, "OK", "FAIL").(string) + "|"

	//what was the commands output
	output := response.Output
	output = strings.ReplaceAll(output, "|", "") //make sure the "|" is not in the output so the whole thing can be split up easily

	workingOutput += output

	return workingOutput
}

func CustomResponseToString(Text string, Pass bool) string {
	workingOutput := ""

	//start with the date
	now := time.Now()
	dateString := now.Format("2006-01-02T15:04:05")
	workingOutput += dateString + "|"

	//was the command successful
	workingOutput += utils.If(Pass, "OK", "FAIL").(string) + "|"

	//what was the commands output
	output := Text
	output = strings.ReplaceAll(output, "|", "") //make sure the "|" is not in the output so the whole thing can be split up easily

	workingOutput += output

	return workingOutput
}

//Thought the best place to put this would be here

/*
Deals with the commands themselvs
*/

type Command struct {
	Alias string
	Args  []string
}

type CommandResponse struct {
	Success bool
	Output  string
}

// ParseInputString takes a string and returns the command and arguments
func ParseInputString(command string) ([]Command, error) {
	commandsRawBytes := strings.Split(command, ";")

	if len(commandsRawBytes) == 1 {
		return nil, errors.New("no commands found. did you forget a semicolon? (;)")
	}

	//the last peice of the command split is leftover
	commandsRawBytes = utils.Remove(commandsRawBytes, len(commandsRawBytes)-1)

	var workingOutput []Command

	for _, commandRaw := range commandsRawBytes {
		commandParts := strings.Split(commandRaw, " ")
		if len(commandParts) == 0 {
			return nil, errors.New("no command alias found")
		}

		workingOutput = append(workingOutput, Command{
			Alias: commandParts[0],
			Args:  commandParts[1:],
		})
	}

	return workingOutput, nil

}

func Verify(text string) bool {
	if len(text) > MAX_SENDBACK {
		return false
	}
	return true
}

/*
Deals with TCP communications between the orchestrator and anything that wants to connect to it.
*/

// Connection Ya sends things in, and ya receive things out. Simple
type Connection struct {
	IpAddr string
	Conn   net.Conn
	Router *Application
}

func (c *Connection) Handle() {
	//todo: go through handshake process

	resp := CustomResponseToString("Welcome to FastScale. \nThis is "+c.Router.Name, true)
	_, err := c.Conn.Write([]byte(resp))
	if err != nil {
		log.Error("error before handling started: ", err)
		return
	}

	//Map channel inputs to the connection
	for {
		_, err = c.Conn.Write([]byte("\nFS? "))
		if err != nil {
			log.Error("error before handling started: ", err)
			return
		}

		buffer := make([]byte, 1024) //Create a buffer to store incoming data
		_, err = c.Conn.Read(buffer)
		if err != nil {
			log.Error("Error reading from connection", err)
			err := c.Conn.Close()
			if err != nil {
				return
			}
			return
		}

		//save command to variable
		command := string(buffer)

		//parse command
		commands, err := ParseInputString(command)
		if err != nil {
			log.Error("Error parsing command", err)
			resp := CustomResponseToString(err.Error(), false)
			_, err := c.Conn.Write([]byte(resp))
			if err != nil {
				log.Error("error before handling started: ", err)
				return
			}
			continue
		}

		//record the number of matched handlers to make sure the command got a match later on
		matches := 0

		//create a place to store the command responses
		var commandResponses []CommandResponse

		//look through each command and check if there is a matching handler in the router
		for _, cmd := range commands {
			for _, handler := range c.Router.Commands {
				if handler.CommandAlies == cmd.Alias {
					//Make sure to indicate to the system that we found a valid handler for the ommand
					matches++
					//We have found an appropriate handler, now we can execute the command and append a new CommandResponse to the commandResponses object
					response := handler.HandleFunc(cmd)

					//Add it to the responses
					commandResponses = append(commandResponses, response)
				}
			}
		}

		if matches == 0 {
			//No handlers were found, send an error to the client and wait for next input.
			resp := CustomResponseToString("Command Unmatched!", false)
			verified := Verify(resp)
			if !verified {
				return
			}
			_, err := c.Conn.Write([]byte(resp))
			if err != nil {
				log.Error("Could not write response to conn", err)
				return
			}
			continue
		}

		for _, resp := range commandResponses {
			resp := CommandResponseToString(resp)
			verified := Verify(resp)
			if !verified {
				return
			}
			_, err := c.Conn.Write([]byte(resp))
			if err != nil {
				log.Error("Could not write response to conn", err)
				return
			}
		}

	}
}

type CommandHandler struct {
	CommandAlies string                        //The command that the handler will respond to
	HandleFunc   func(Command) CommandResponse //The function that will process the command
}

type Application struct {
	Name     string
	Commands []CommandHandler
}

// Listen listens for incoming connections
func (r *Application) Listen() error {
	//Listen for incoming connections
	listen, err := net.Listen("tcp", fmt.Sprint(":", FSP_PORT))

	log.Print("Application " + r.Name + " is now online")

	if err != nil {
		log.Error("Error starting TCP listener", err)
		return err
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Warn("Was not able to accept a connection", err)
			continue
		}

		//Create a new connection
		connection := Connection{
			IpAddr: conn.RemoteAddr().String(),
			Conn:   conn,
			Router: r,
		}

		go connection.Handle()

	}

}
