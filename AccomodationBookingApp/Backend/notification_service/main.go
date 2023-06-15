package main

import "notification_service/startup"

func main() {
	config := startup.NewConfiguration()
	server := startup.NewServer(config)
	server.Start()
}
