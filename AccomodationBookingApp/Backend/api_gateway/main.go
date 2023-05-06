package main

import (
	"api_gateway/startup"
	"log"
)

func main() {
	config := startup.NewConfig()
	server := startup.NewServer(config)
	log.Println("Gateway started...")
	server.Start()
}
