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

	//shortened variable names to omit collision with package names
	var (
		logger *log.Logger                  = log.New(os.Stdout, "[flight-repo] ", log.LstdFlags)
		repo   repository.FlightRepository  = repository.NewFlightRepository(client, logger)
		serv   service.FlightService        = service.NewFlightService(repo)
		contr  *controller.FlightController = controller.NewFlightController(serv)
	)

	flights := uppperRouterGroup.Group("/flight")
	{
		//TODO: assgin handlers
		flights.GET("", contr.GetAll)
		flights.GET(":id", contr.GetById)
		flights.POST("", contr.Create)
		flights.DELETE(":id", contr.Delete)
	}
}
