package endpoints

import (
	"FlightBookingApp/controller"
	"FlightBookingApp/dependencyInjection"
	"FlightBookingApp/middleware"
	"FlightBookingApp/model"
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
		jwtServ    service.JwtService           = depContainer.GetService("jwt").(service.JwtService)
		contr      *controller.TicketController = controller.NewTicketController(serv, jwtServ)
	)
	depContainer.RegisterEntityDependencyBundle("ticket", repo, serv, contr)

	tickets := upperRouterGroup.Group("/ticket")
	tickets.Use(middleware.ValidateToken())
	{
		tickets.GET("",
			middleware.Authorization([]model.Role{model.REGULAR_USER}),
			contr.GetAll)
		tickets.GET(":id",
			middleware.Authorization([]model.Role{model.REGULAR_USER}),
			contr.GetById)
		tickets.GET("/myTickets",
			middleware.Authorization([]model.Role{model.REGULAR_USER}),
			contr.GetAllForCustomer)
		tickets.POST("",
			middleware.Authorization([]model.Role{model.ADMIN}),
			contr.Create)

		tickets.DELETE(":id", contr.Delete)

		tickets.POST("/buy",
			middleware.Authorization([]model.Role{model.REGULAR_USER}),
			contr.BuyTicket)

	}
}
