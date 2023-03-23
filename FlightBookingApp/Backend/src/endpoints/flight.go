package endpoints

import (
	"FlightBookingApp/controller"
	"FlightBookingApp/repository"
	"FlightBookingApp/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

func DefineFlightEndpoints(uppperRouterGroup *gin.RouterGroup, client *mongo.Client) {

	logger := log.New(os.Stdout, "[flight-repo] ", log.LstdFlags)
	repository := repository.NewFlightRepository(client, logger)

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
}
