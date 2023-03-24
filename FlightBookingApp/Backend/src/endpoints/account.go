package endpoints

import (
	"FlightBookingApp/controller"
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

	accounts := upperRouterGroup.Group("/account")
	{
		accounts.GET("", accContr.GetAll)
		accounts.GET(":id", accContr.GetById)
		accounts.POST("", accContr.Register)
		accounts.DELETE(":id", accContr.Delete)
	}
}