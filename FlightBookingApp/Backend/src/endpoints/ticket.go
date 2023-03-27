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

func DefineTicketEndpoints(upperRouterGroup *gin.RouterGroup, client *mongo.Client) {

	//shortened variable names to omit collision with package names
	var (
		logger *log.Logger                  = log.New(os.Stdout, "[ticket-repo] ", log.LstdFlags)
		repo   repository.TicketRepositry   = repository.NewTicketRepositry(client, logger)
		serv   service.TicketService        = service.NewTicketService(repo)
		contr  *controller.TicketController = controller.NewTicketController(serv)
	)

	flights := upperRouterGroup.Group("/ticket")
	{
		flights.GET("", contr.GetAll)
		flights.GET(":id", contr.GetById)
		flights.POST("", contr.Create)
		flights.DELETE(":id", contr.Delete)
	}
}
