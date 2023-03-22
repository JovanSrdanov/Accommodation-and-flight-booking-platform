package main

import (
	"FlightBookingApp/endpoints"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	server := gin.Default()
	logger := log.New(os.Stdout, "[flight-app-api] ", log.LstdFlags)

	apiRoutes := server.Group("/api")
	{
		_, err := endpoints.DefineFlightEndpoints(apiRoutes)
		if err != nil {

			logger.Fatal("VARVARIN")
			return
		}
	}

	port := os.Getenv("PORT")
	//port := "4200"
	server.Run(":" + port)
}
