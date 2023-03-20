package main

import (
	"FlightBookingApp/endpoints"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	//Has default logging and recovery middleware
	server := gin.Default()

	apiRoutes := server.Group("/api")
	{
		endpoints.DefineFlightEndpoints(apiRoutes)
	}

	port := os.Getenv("PORT")
	//port := "4200"
	server.Run(":" + port)
}
