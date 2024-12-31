package main

import (
	"github.com/TboOffical/fastscale/commander/utils"
	"log"
	"os"
)

/*
The following are constants for database variables, to make sure that the name of the variable can be changed without having to change the name throughout the program.
And that the names cant be misspelled and new keys created without a need for them
*/
const (
	VarLastStarted         = "last_started"
	RegistrationServerPort = 8881
)

/*
Controller is a struct that represents a running controller that has registered with this commander in order to receive the running configuration
*/
type Controller struct {
	id             string
	ip             string
	nodesConnected int
}

/*
ConfigFile defines the structure of the json config file that is required for the commander to start, or even create the necessary files.
*/
type ConfigFile struct {
	RbmqAddr          string `json:"rbmq_addr"`
	RbmqManagementUrl string
	DatabaseFile      string `json:"database_file"`
	PostgresAddress   string
}

func (c *ConfigFile) load() {
	RbmqAdr, err := utils.NullErr(os.Getenv("rbmq_address"))
	if err != nil {
		log.Fatalln("Was not able to find a rabbitMq address in .env, use 'rbmq_address'")
	}
	DatabFile, err := utils.NullErr(os.Getenv("database_file"))
	if err != nil {
		log.Fatalln("Was not able to find a database file in .env, use 'database_file'")
	}
	RbmqManagement, err := utils.NullErr(os.Getenv("rbmq_managment"))
	if err != nil {
		log.Fatalln("Was not able to find a database file in .env, use 'database_file'")
	}
	PostgressAddress, err := utils.NullErr(os.Getenv("postgres_address"))
	if err != nil {
		log.Fatalln("Was not able to find a postgres conn string in .env, use 'postgres_address'")
	}

	c.RbmqAddr = RbmqAdr
	c.DatabaseFile = DatabFile
	c.RbmqManagementUrl = RbmqManagement
	c.PostgresAddress = PostgressAddress
}
