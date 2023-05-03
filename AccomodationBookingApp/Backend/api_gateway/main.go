package main

import (
	"api_gateway/communication"
	"api_gateway/configuration"
	"log"
)

func main() {
	config := configuration.NewConfig()
	server := communication.NewServer(config)
	log.Println("Gateway started...")
	server.Start()
}
