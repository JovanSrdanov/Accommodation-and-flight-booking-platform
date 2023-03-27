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
		accLogger *log.Logger                   = log.New(os.Stdout, "[account-repo] ", log.LstdFlags)
		userLogger *log.Logger									= log.New(os.Stdout, "[user-repo] ", log.LstdFlags)
		accRepo   repository.AccountRepository  = repository.NewAccountRepository(client, accLogger)
		userRepo  repository.UserRepository 		= repository.NewUserRepository(client, userLogger)
		accServ   service.AccountService        = service.NewAccountService(accRepo, userRepo)
		accContr  *controller.AccountController = controller.NewAccountController(accServ)
	)

	authenticatedAccounts := upperRouterGroup.Group("/account")
	authenticatedAccounts.Use(middleware.ValidateToken())
	{
		authenticatedAccounts.GET("", middleware.Authrorization([]model.Role{model.ADMIN}),
										 accContr.GetAll)
		authenticatedAccounts.GET(":id", 
															middleware.Authrorization([]model.Role{model.REGULAR_USER, model.ADMIN}),
		 													accContr.GetById)
		authenticatedAccounts.DELETE(":id",
																	middleware.Authrorization([]model.Role{model.ADMIN}),
																	accContr.Delete)
		authenticatedAccounts.GET("/refresh-token/:token", 
									middleware.Authrorization([]model.Role{model.REGULAR_USER, model.ADMIN}),
									accContr.RefreshAccessToken)
	}

	// anyone can use these 
	unauthenticated := upperRouterGroup.Group("/account")
	{
		unauthenticated.POST("/login", accContr.Login)
		unauthenticated.POST("/register", accContr.Register)
		unauthenticated.GET("/emailver/:username/:verPass", accContr.VerifyEmail)
	}
}