package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func doesUserExist(user string, file ConfigFile) bool {
	client := http.Client{
		Timeout: time.Second * 30,
	}

	parsedString := strings.Replace(file.RbmqManagementUrl, "amqp://", "", 1)
	parsedString = strings.Split(parsedString, "@")[0]
	finalParsed := strings.Split(parsedString, ":")

	splitUsername := finalParsed[0]
	splitPassword := finalParsed[1]

	URL, err := url.Parse(file.RbmqManagementUrl)
	if err != nil {
		log.Fatalln("RabbitMQ management URL is not valid", err.Error())
	}

	req, err := http.NewRequest("GET", URL.Scheme+URL.Host+"/api/users/"+user, nil)
	req.SetBasicAuth(splitUsername, splitPassword)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("RabbitMQ not responding!", err.Error())
	}

	if resp.StatusCode != 200 {
		return false
	}
	return true
}

func startupRabbitCheck(file ConfigFile) error {
	//check to make sure that users exist for all of the controllers
	return nil
}
