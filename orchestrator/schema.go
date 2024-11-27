package main

import "time"

type Node struct {
	Online      bool
	UID         string
	Ip          string
	Capacity    int    //Number 1-100 representing the capacity of the node as a percentage, the node can control this.
	Version     string //Version of node software
	LastCheckin time.Time
}

type Secret struct {
	UID    string
	Secret string
}

// a change in the changelog
type change struct {
	UID         string
	Date        time.Time
	Description string
}

// A simple changeLog
type changeLog struct {
	Changes []change
}

type indexEntry struct {
	UID           string
	IndexType     string
	IndexFilePath string
	IndexActive   bool //The system will stop when it finds an active index, it will than load the index file and deploy it to nodes.
}

// The indexFile is used to keep tract of detailed information regarding the plans and setups.
// This is the root structure for the file
type indexFile struct {
	changeLog
	indexes []indexEntry
}
