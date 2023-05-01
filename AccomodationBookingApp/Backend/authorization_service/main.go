package main

import (
	"authorization_service/communication"
	configuration "authorization_service/configuration"
)

func main() {
	configuration := configuration.NewConfiguration()
	server := communication.NewServer(configuration)
	server.Start()
}
