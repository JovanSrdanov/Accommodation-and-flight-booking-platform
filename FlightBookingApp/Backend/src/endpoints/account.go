package endpoints

import (
	"FlightBookingApp/controller"
	"FlightBookingApp/middleware"
	"FlightBookingApp/model"
	"FlightBookingApp/repository"
	"FlightBookingApp/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func DefineAccountEndpoints(upperRouterGroup *gin.RouterGroup, client *mongo.Client) {
	var (
		accLogger *log.Logger                  = log.New(os.Stdout, "[account-repo] ", log.LstdFlags)
		accRepo   repository.AccountRepository  = repository.NewAccountRepository(client, accLogger)
		accServ   service.AccountService        = service.NewAccountService(accRepo)
		accContr  *controller.AccountController = controller.NewAccountController(accServ)
	)

	// temp, only testing authorization
	accounts := upperRouterGroup.Group("/account")
	accounts.Use(middleware.ValidateToken())
	{
		accounts.GET("", middleware.Authrorization([]model.Role{model.ADMIN}),
										 accContr.GetAll)
		accounts.GET(":id", middleware.Authrorization([]model.Role{model.REGULAR_USER}),
		 										accContr.GetById)
		//accounts.POST("/register", accContr.Register)
		//accounts.POST("/login", accContr.Login)
		accounts.DELETE(":id", accContr.Delete)
		accounts.GET("/refresh-token/:token", accContr.RefreshAccessToken)
	}

	//temp, should be in accounts group
	test := upperRouterGroup.Group("/account")
	{
		test.POST("/login", accContr.Login)
		test.POST("/register", accContr.Register)
		test.GET("/emailver/:username/:verPass", accContr.VerifyEmail)
	}
}