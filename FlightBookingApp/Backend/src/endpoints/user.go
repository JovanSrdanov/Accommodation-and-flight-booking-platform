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

func DefineUserEndpoints(upperRouterGroup *gin.RouterGroup, client *mongo.Client) {
	var (
		accLogger  *log.Logger                   = log.New(os.Stdout, "[account-repo] ", log.LstdFlags)
		userLogger *log.Logger                   = log.New(os.Stdout, "[user-repo] ", log.LstdFlags)
		accRepo    repository.AccountRepository  = repository.NewAccountRepository(client, accLogger)
		userRepo   repository.UserRepository     = repository.NewUserRepository(client, userLogger)
		userServ   service.UserService           = service.NewUserService(userRepo)
		userContr  *controller.UserController     = controller.NewUserController(userServ, accRepo)
	)

	users := upperRouterGroup.Group("/user")
	users.Use(middleware.ValidateToken())
	{
		users.GET("", middleware.Authrorization([]model.Role{model.ADMIN}),
			userContr.GetAll)
		users.GET(":id",
			middleware.Authrorization([]model.Role{model.REGULAR_USER, model.ADMIN}),
			userContr.GetById)
		users.GET("logged-in",
			middleware.Authrorization([]model.Role{model.REGULAR_USER, model.ADMIN}),
			userContr.GetLoggedInUserInfo)
		users.DELETE(":id",
			middleware.Authrorization([]model.Role{model.ADMIN}),
			userContr.Delete)
	}
}