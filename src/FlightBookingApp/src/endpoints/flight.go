package endpoints

import (
	"FlightBookingApp/controller"
	"FlightBookingApp/repository"
	"FlightBookingApp/service"
	"github.com/gin-gonic/gin"
)

func DefineFlightEndpoints(uppperRouterGroup *gin.RouterGroup) *gin.RouterGroup {

	var (
		repository repository.FlightRepository = repository.NewFlightRepository()
		service    service.FlightService       = service.NewFlightService(repository)
		controller controller.FlightController = controller.NewFlightController(service)
	)

	flights := uppperRouterGroup.Group("/flight")
	{
		//TODO: assgin handlers
		flights.GET("", controller.GetAll)
		flights.GET(":id", controller.GetById)
		flights.POST("", controller.Create)
		flights.DELETE(":id", nil)
	}

	return flights
}
