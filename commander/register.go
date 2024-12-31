package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type RegisterHandlerInput struct {
	RegisterKey      string
	IP               string
	Mac              string
	ReturnRabbitCred bool //True if the controller lost, and/or needs new rabbit cred. Regenerates password.
}

func registerHandler(c *gin.Context) {

}

func startRegistrationServer() error {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	log.Println("Registration server started on port", RegistrationServerPort)

	r.POST("/register", registerHandler)

	err := r.Run(fmt.Sprint(":", RegistrationServerPort))
	if err != nil {
		return err
	}

	return nil
}
