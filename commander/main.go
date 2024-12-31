package main

import (
	"database/sql"
	"errors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"time"
)

const (
	FastscaleVersonMajor  = 0
	FastscaleVersionMinor = 1
	FastscaleVersionPatch = 0
)

func main() {
	//Love a cool banner
	log.Println(`
        ______           __  _____            __   
	   / ____/___ ______/ /_/ ___/_________ _/ /__ 
	  / /_  / __ ./ ___/ __/\__ \/ ___/ __ ./ / _ \
	 / __/ / /_/ (__  ) /_ ___/ / /__/ /_/ / /  __/
	/_/    \__,_/____/\__//____/\___/\__,_/_/\___/

	Commander v1.0.0
	`)

	//create the conf folder if it is not already there
	_, err := os.ReadDir("./conf")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			//create the folder
			err = os.Mkdir("./conf", 7777)
			if err != nil {
				log.Fatal("Unable to create configuration folder!", err.Error())
			}
		} else {
			log.Fatal("Unknown filesystem error", err)
		}
	}

	//try to read the config data from env
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//map the content to an object
	var configData ConfigFile
	configData.load()

	//connect to the local commander database

	//make sure the database file specified exists
	_, err = os.Stat(configData.DatabaseFile)

	//create the database file
	if errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(configData.DatabaseFile)
		if err != nil {
			log.Fatal("Was not able to create database file at the following path", configData.DatabaseFile, err.Error())
		}
	} else if err != nil {
		log.Fatal("An Unknown error occured", err.Error())
	}

	db, err := sql.Open("sqlite3", "file:"+configData.DatabaseFile+"?cache=shared")
	if err != nil {
		log.Fatal("Could not open database, this is bad!!", err.Error())
	}

	defer db.Close()

	//Sync the tables for the local DB
	createTables(db)

	log.Println("Sqlite3 status changed to up")

	//Create a connection to the rabbitMQ instance
	conn, err := amqp.Dial(configData.RbmqAddr)
	if err != nil {
		log.Fatalln("Could not connect to Rabbit MQ. If you are not expecting this, this is BAD, it will break the entire network.", err.Error())
		return
	}
	defer conn.Close()

	_, err = conn.Channel()
	if err != nil {
		log.Fatalln("Was not able to create channel with rabbitMQ, BAD BAD BAD!", err.Error())
		return
	}

	log.Println("RabbitMQ status changed to up")

	now := time.Now()
	err = setDatabaseVariable(VarLastStarted, now.String(), db)
	if err != nil {
		log.Fatalln("cannot access database variables! ", err)
	}

	//Try to connect to the main postgres application server

	appDb, err := sql.Open("postgres", configData.PostgresAddress)
	if err != nil {
		log.Fatalln("Unable to connect to application database! This is bad!!", err.Error())
	}

	err = appDb.Ping()
	if err != nil {
		log.Fatalln("Unable to ping application database, check to make sure connection is not being blocked by a firewall!", err.Error())
	}

	log.Println("Postgres status changed to up")

	//Load the model into system

	var mm ModelManager

	err = mm.loadModel("./conf/model.json", appDb)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("No model detected, creating the default model")
			err = mm.createDefaultModel()
			if err != nil {
				log.Fatal("Unable to create the default mode, ", err)
			}
			err = mm.loadModel("./conf/model.json", appDb)
			if err != nil {
				log.Fatal("Failed to load the model for the second time after creating the default model!", err.Error())
			}
		} else {
			//todo: check for backups and run the last known backup
			log.Fatal("Error when loading model", err)
		}
	}

	//Start the services in different threads

	//start the api to listen for registration requests
	go func() {
		err := startRegistrationServer()
		if err != nil {
			log.Fatal("Failed to start the registration server", err.Error())
		}
	}()

}
