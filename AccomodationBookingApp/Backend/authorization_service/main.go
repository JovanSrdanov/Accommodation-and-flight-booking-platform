package main

import (
	"authorization_service/startup"
	"authorization_service/startup/configuration"
	"log"
	"os"
)

func main() {
	configuration := configuration.NewConfiguration()
	server := startup.NewServer(configuration)
	log.Printf("Authorization service started, running on %s:%s", os.Getenv("AUTHORIZATION_SERVICE_HOST"), os.Getenv("AUTHORIZATION_SERVICE_PORT"))
	server.Start()
}
