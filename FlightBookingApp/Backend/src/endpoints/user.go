package endpoints

import (
	"FlightBookingApp/controller"
	"FlightBookingApp/dependencyInjection"
	"FlightBookingApp/middleware"
	"FlightBookingApp/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func DefineUserEndpoints(upperRouterGroup *gin.RouterGroup, client *mongo.Client, depContainer *dependencyInjection.DependencyContainer) {
	// Dependencies are instantiated inside account endpoint definition
	userContr := depContainer.GetController("user").(*controller.UserController)
	users := upperRouterGroup.Group("/user")
	users.Use(middleware.ValidateToken())
	{
		users.GET("", middleware.Authorization([]model.Role{model.ADMIN}),
			userContr.GetAll)
		users.GET(":id",
			middleware.Authorization([]model.Role{model.REGULAR_USER, model.ADMIN}),
			userContr.GetById)
		users.GET("logged-in",
			middleware.Authorization([]model.Role{model.REGULAR_USER, model.ADMIN}),
			userContr.GetLoggedInUserInfo)
		users.DELETE(":id",
			middleware.Authorization([]model.Role{model.ADMIN}),
			userContr.Delete)
	}
}
