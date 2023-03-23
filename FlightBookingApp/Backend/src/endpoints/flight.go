package endpoints

import (
	"FlightBookingApp/controller"
	"FlightBookingApp/repository"
	"FlightBookingApp/service"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func DefineFlightEndpoints(uppperRouterGroup *gin.RouterGroup) (*gin.RouterGroup, error) {

	logger := log.New(os.Stdout, "[flight-repo] ", log.LstdFlags)
	repository, err := repository.NewFlightRepository(context.Background(), logger)
	if err != nil {
		return nil, err
	}

	var (
		service    service.FlightService       = service.NewFlightService(repository)
		controller controller.FlightController = controller.NewFlightController(service)
	)

	flights := uppperRouterGroup.Group("/flight")
	{
		//TODO: assgin handlers
		flights.GET("", controller.GetAll)
		flights.GET(":id", controller.GetById)
		flights.POST("", controller.Create)
		flights.DELETE(":id", controller.Delete)
	}

	return flights, nil
}
