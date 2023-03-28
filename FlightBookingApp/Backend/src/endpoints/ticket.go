package endpoints

import (
	"FlightBookingApp/controller"
	"FlightBookingApp/dependencyInjection"
	"FlightBookingApp/repository"
	"FlightBookingApp/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

func DefineTicketEndpoints(upperRouterGroup *gin.RouterGroup, client *mongo.Client, depContainer *dependencyInjection.DependencyContainer) {

	//shortened variable names to omit collision with package names
	var (
		logger     *log.Logger                  = log.New(os.Stdout, "[ticket-repo] ", log.LstdFlags)
		repo       repository.TicketRepositry   = repository.NewTicketRepositry(client, logger)
		flightRepo repository.FlightRepository  = depContainer.GetRepository("flight").(repository.FlightRepository)
		serv       service.TicketService        = service.NewTicketService(repo, flightRepo)
		contr      *controller.TicketController = controller.NewTicketController(serv)
	)

	tickets := upperRouterGroup.Group("/ticket")
	{
		tickets.GET("", contr.GetAll)
		tickets.GET(":id", contr.GetById)
		tickets.GET("/myTickets", contr.GetAllForCustomer)
		tickets.POST("", contr.Create)
		tickets.POST("/buy", contr.BuyTicket)
		tickets.DELETE(":id", contr.Delete)
	}
}
