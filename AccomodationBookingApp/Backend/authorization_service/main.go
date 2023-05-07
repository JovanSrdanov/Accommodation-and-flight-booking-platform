package main

import (
	"authorization_service/startup"
	"log"
	"os"
)

func main() {
	configuration := startup.NewConfiguration()
	server := startup.NewServer(configuration)
	log.Printf("Authorization service started, running on %s:%s", os.Getenv("AUTHORIZATION_SERVICE_HOST"), os.Getenv("AUTHORIZATION_SERVICE_PORT"))
	server.Start()
}
