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
		logger     *log.Logger                  = log.New(os.Stdout, "[ticket-repo] ", log.LstdFlags)
		logger2    *log.Logger                  = log.New(os.Stdout, "[flight-repo] ", log.LstdFlags)
		repo       repository.TicketRepositry   = repository.NewTicketRepositry(client, logger)
		flightRepo repository.FlightRepository  = repository.NewFlightRepository(client, logger2)
		serv       service.TicketService        = service.NewTicketService(repo, flightRepo)
		contr      *controller.TicketController = controller.NewTicketController(serv)
	)

	tickets := upperRouterGroup.Group("/ticket")
	{
		tickets.GET("", contr.GetAll)
		tickets.GET(":id", contr.GetById)
		tickets.GET("/getc", contr.GetAllForCustomer)
		tickets.POST("", contr.Create)
		tickets.POST("/buy", contr.BuyTicket)
		tickets.DELETE(":id", contr.Delete)
	}
}
