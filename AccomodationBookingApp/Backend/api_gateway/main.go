package main

import (
	"api_gateway/startup"
	"api_gateway/startup/configuration"
	"log"
)

func main() {
	config := configuration.NewConfig()
	server := startup.NewServer(config)
	log.Println("Gateway started...")
	server.Start()
}
