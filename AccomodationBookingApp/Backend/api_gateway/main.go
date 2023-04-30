package main

import (
	"api_gateway/communication"
	"api_gateway/configuration"
)

func main() {
	config := configuration.NewConfig()
	server := communication.NewServer(config)
	server.Start()
}
