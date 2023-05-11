package main

import (
	"reservation_service/startup"
)

func main() {
	config := startup.NewConfiguration()
	server := startup.NewServer(config)
	server.Start()
}
