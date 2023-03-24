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

func DefineAirportEndpoints(upperRouterGroup *gin.RouterGroup, client *mongo.Client) {

	//shortened variable names to omit collision with package names
	var (
		logger *log.Logger                   = log.New(os.Stdout, "[airport-repo] ", log.LstdFlags)
		repo   repository.AirportRepository  = repository.NewAirportRepository(client, logger)
		serv   service.AirportService        = service.NewAirportService(repo)
		contr  *controller.AirportController = controller.NewAirportController(serv)
	)

	airports := upperRouterGroup.Group("/airport")
	{
		airports.GET("", contr.GetAll)
		airports.GET(":id", contr.GetById)
	}
}
