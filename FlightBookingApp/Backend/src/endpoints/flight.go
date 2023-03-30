package endpoints

import (
	"FlightBookingApp/controller"
	"FlightBookingApp/dependencyInjection"
	"FlightBookingApp/repository"
	"FlightBookingApp/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func DefineFlightEndpoints(upperRouterGroup *gin.RouterGroup, client *mongo.Client, depContainer *dependencyInjection.DependencyContainer) {

	//shortened variable names to omit collision with package names
	var (
		logger *log.Logger                  = log.New(os.Stdout, "[flight-repo] ", log.LstdFlags)
		repo   repository.FlightRepository  = repository.NewFlightRepository(client, logger)
		serv   service.FlightService        = service.NewFlightService(repo)
		contr  *controller.FlightController = controller.NewFlightController(serv)
	)
	depContainer.RegisterEntityDependencyBundle("flight", repo, serv, contr)

	flights := upperRouterGroup.Group("/flight")
	{
		flights.GET("", contr.GetAll)
		flights.GET(":id", contr.GetById)
		flights.GET("search", contr.Search)
		flights.POST("", contr.Create)
		flights.PATCH(":id", contr.Cancel)
	}
	searchFlights := upperRouterGroup.Group("/search-flights")
	{

		searchFlights.GET("", contr.Search)

	}
}
