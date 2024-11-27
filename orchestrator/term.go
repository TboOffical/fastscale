package main

/*
Deals with TCP communications between the orchestrator and anything that wants to connect to it.
*/

type Connection struct {
	IpAddr string
	Chan   chan string
}

type CommandHandler struct {
	CommandAlies string              //The command that the handler will respond to
	HandleFunc   func(string) string //The function that will process the command
}

type Router struct {
	ActiveConnections []Connection
	Commands          []CommandHandler
}
