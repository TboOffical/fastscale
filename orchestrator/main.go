package main

import (
	"fmt"
	"github.com/charmbracelet/log"
)

func main() {
	//Love a cool banner
	log.Print(`
		______           __  _____            __   
	   / ____/___ ______/ /_/ ___/_________ _/ /__ 
	  / /_  / __ ./ ___/ __/\__ \/ ___/ __ ./ / _ \
	 / __/ / /_/ (__  ) /_ ___/ / /__/ /_/ / /  __/
	/_/    \__,_/____/\__//____/\___/\__,_/_/\___/      
	`)

	//just for testing

	testCommand := CommandHandler{
		CommandAlies: "test",
		HandleFunc: func(c Command) CommandResponse {
			return CommandResponse{
				Success: true,
				Output:  "hello from test command" + fmt.Sprint(c.Args),
			}
		},
	}

	r := Application{
		Name: "TestApplication",
		Commands: []CommandHandler{
			testCommand,
		},
	}

	err := r.Listen()
	if err != nil {
		log.Error("error occured", err)
		return
	}

	////Run the initialization sequence
	//
	////todo: Check for brain data backup
	//
	////Check for the existence of the data directory
	//err := checkDirectories()
	//if err != nil {
	//	SetupComplete = false
	//}
	//
	//if !SetupComplete {
	//
	//}

}
