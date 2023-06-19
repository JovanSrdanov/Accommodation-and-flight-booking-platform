package endpoints

import (
	"FlightBookingApp/controller"
	"FlightBookingApp/dependencyInjection"
	"FlightBookingApp/middleware"
	"FlightBookingApp/model"
	"FlightBookingApp/repository"
	"FlightBookingApp/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func DefineAccountEndpoints(upperRouterGroup *gin.RouterGroup, client *mongo.Client, depContainer *dependencyInjection.DependencyContainer) {
	var (
		accLogger    *log.Logger                   = log.New(os.Stdout, "[account-repo] ", log.LstdFlags)
		userLogger   *log.Logger                   = log.New(os.Stdout, "[user-repo] ", log.LstdFlags)
		apiKeyLogger *log.Logger                   = log.New(os.Stdout, "[api-key-repo] ", log.LstdFlags)
		accRepo      repository.AccountRepository  = repository.NewAccountRepository(client, accLogger)
		userRepo     repository.UserRepository     = repository.NewUserRepository(client, userLogger)
		apiKeyRepo   repository.ApiKeyRepository   = repository.NewApiKeyRepository(client, apiKeyLogger)
		accServ      service.AccountService        = service.NewAccountService(accRepo, userRepo)
		apiKeyServ   *service.ApiKeyService        = service.NewApiKeyService(apiKeyRepo)
		accContr     *controller.AccountController = controller.NewAccountController(accServ, apiKeyServ)
		// For easier dependency chain following
		//TODO stefan : vidi je l ovde treba da se bas vraca struct ili mozda interfejs
		userServ service.UserService = service.NewUserService(userRepo)
		//TODO stefan : ubaci da se u controller stavlja accService, a ne repo
		userContr *controller.UserController = controller.NewUserController(userServ, accRepo)
	)
	depContainer.RegisterEntityDependencyBundle("user", userRepo, userServ, userContr)
	depContainer.RegisterEntityDependencyBundle("account", accRepo, accServ, accContr)
	depContainer.RegisterEntityDependencyBundle("apiKey", apiKeyRepo, apiKeyServ, nil)

	authenticatedAccounts := upperRouterGroup.Group("/account")
	authenticatedAccounts.Use(middleware.ValidateToken())
	{
		authenticatedAccounts.GET("", middleware.Authorization([]model.Role{model.ADMIN}),
			accContr.GetAll)
		authenticatedAccounts.GET(":id",
			middleware.Authorization([]model.Role{model.REGULAR_USER, model.ADMIN}),
			accContr.GetById)
		authenticatedAccounts.GET("logged/info",
			middleware.Authorization([]model.Role{model.REGULAR_USER, model.ADMIN}),
			accContr.GetLoggedInfo)
		authenticatedAccounts.DELETE(":id",
			middleware.Authorization([]model.Role{model.ADMIN}),
			accContr.Delete)
		authenticatedAccounts.POST("api-key",
			middleware.Authorization([]model.Role{model.REGULAR_USER}),
			accContr.CreateApiKey)
		authenticatedAccounts.GET("api-key",
			middleware.Authorization([]model.Role{model.REGULAR_USER}),
			accContr.GetApiKey)
	}

	// anyone can use these
	unauthenticated := upperRouterGroup.Group("/account")
	{
		unauthenticated.POST("/login", accContr.Login)
		unauthenticated.POST("/register", accContr.Register)
		unauthenticated.GET("/refresh-token", accContr.RefreshAccessToken)
		unauthenticated.GET("/emailver/:username/:verPass", accContr.VerifyEmail)
		unauthenticated.GET("/logout", accContr.Logout)
	}
}
