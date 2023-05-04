package main

import (
	"user_profile_service/communication"
	"user_profile_service/configuration"
)

func main() {
	config := configuration.NewConfiguration()
	server := communication.NewServer(config)
	server.Start()
}
