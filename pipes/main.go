package main

import "github.com/TboOffical/fastscale/orchestrator/comms"

func main() {
	conn := comms.NewConnection("192.168.1.5")
	err := conn.Connect()
	if err != nil {
		panic(err)
		return
	}
}
