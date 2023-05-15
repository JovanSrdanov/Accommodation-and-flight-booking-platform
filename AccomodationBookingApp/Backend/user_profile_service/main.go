package main

import (
	"user_profile_service/startup"
)

func main() {
	config := startup.NewConfiguration()
	server := startup.NewServer(config)
	server.Start()
}
