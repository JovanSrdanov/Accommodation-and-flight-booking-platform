package main

import (
	"FlightBookingApp/endpoints"
	"github.com/gin-gonic/gin"
)

func main() {
	//Has default logging and recovery middleware
	server := gin.Default()

	apiRoutes := server.Group("/api")
	{
		endpoints.DefineFlightEndpoints(apiRoutes)
	}

	server.Run(":4200")
}
