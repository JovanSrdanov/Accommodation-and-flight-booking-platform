package main

import (
	"FlightBookingApp/controller"
	"FlightBookingApp/repository"
	"FlightBookingApp/service"
	"github.com/gin-gonic/gin"
)

// TODO delete
func testHandler(ctx *gin.Context) {
	ctx.JSON(200, "Radim sve ti jebem")
}

var (
	flightRepository repository.FlightRepository = repository.NewFlightRepository()
	flightService    service.FlightService       = service.NewFlightService(flightRepository)
	flightController controller.FlightController = controller.NewFlightController(flightService)
)

func main() {
	//Has default logging and recovery middleware
	server := gin.Default()

	apiRoutes := server.Group("/api")
	{
		flights := apiRoutes.Group("/flight")
		{
			//TODO: assgin handlers
			flights.GET("", flightController.GetAll)
			flights.GET("/:id", nil)
			flights.POST("", flightController.Create)
			flights.DELETE("/:id", nil)
		}
	}

	server.Run(":4200")
}
