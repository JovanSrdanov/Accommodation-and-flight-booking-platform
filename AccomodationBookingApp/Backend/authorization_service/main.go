package main

import (
	"authorization_service/communication"
	configuration "authorization_service/configuration"
	"log"
	"os"
)

func main() {
	configuration := configuration.NewConfiguration()
	server := communication.NewServer(configuration)
	log.Printf("Authorization service started, running on %s:%s", os.Getenv("AUTHORIZATION_SERVICE_HOST"), os.Getenv("AUTHORIZATION_SERVICE_PORT"))
	server.Start()
}
