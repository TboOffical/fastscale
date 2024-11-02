package main

import (
	"net"
	"time"
)

/*
PipesNode a single DB node
*/
type PipesNode struct {
	Online      bool
	UID         string
	Ip          net.IP
	Capacity    int    //Number 1-100 representing the capacity of the node as a percentage, the node can control this.
	Version     string //Pipes version of node
	LastCheckin time.Time
}
