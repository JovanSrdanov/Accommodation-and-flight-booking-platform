package handler

import (
	"api_gateway/communication"
	"api_gateway/communication/middleware"
	"api_gateway/dto"
	"authorization_service/domain/model"
	"authorization_service/domain/token"
	authorization "common/proto/authorization_service/generated"
	notification "common/proto/notification_service/generated"
	user_profile "common/proto/user_profile_service/generated"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
)

type UserHandler struct {
	authorizationServiceAddress string
	userProfileServiceAddress   string
	notificationServiceAddress  string
	tokenMaker                  token.Maker
}

func NewUserHandler(authorizationServiceAddress string, userProfileServiceAddress string,
	tokenMaker token.Maker) *UserHandler {
	return &UserHandler{authorizationServiceAddress: authorizationServiceAddress,
		userProfileServiceAddress: userProfileServiceAddress,
		tokenMaker:                tokenMaker,
	}
}

func (handler UserHandler) Init(router *gin.RouterGroup) {
	userGroup := router.Group("/user")
	userGroup.GET("/:username/info",
		middleware.ValidateToken(handler.tokenMaker),
		middleware.Authorization([]model.Role{model.Guest, model.Host}),
		handler.GetUserInfo)
	userGroup.GET("/logged-in-info",
		middleware.ValidateToken(handler.tokenMaker),
		middleware.Authorization([]model.Role{model.Guest, model.Host}),
		handler.GetLoggedInUserInfo,
	)
	userGroup.POST("", handler.CreateUser)
}

func (handler UserHandler) GetUserInfo(ctx *gin.Context) {
	username := ctx.Param("username")

	if username == "" {
		ctx.JSON(http.StatusBadRequest, "Username not provided")
		return
	}

	log.Println("username for get all: ", username)
	log.Println("GIN KONTEEEEEEEKS: ", ctx.GetHeader("Authorization"))
	ctxGrpc := createGrpcContextFromGinContext(ctx)

	var userInfo dto.UserInfo

	//TODO errorMessage handling
	// it's important to pass the ctx to all handlers that need to call a grpc method with a client,
	// because it has the auth header embedded in it
	handler.addAccountCredentialsInfo(&userInfo, username, ctxGrpc)
	handler.addUserProfileInfo(&userInfo, ctxGrpc)

	ctx.JSON(http.StatusOK, userInfo)
}

func (handler UserHandler) addAccountCredentialsInfo(userInfo *dto.UserInfo, username string, ctx context.Context) error {
	authorizationClient := communication.NewAuthorizationClient(handler.authorizationServiceAddress)

	accountCredentialsInfo, err := authorizationClient.GetByUsername(ctx, &authorization.GetByUsernameRequest{Username: username})
	log.Println("accCredInfo from authClient GetByUsername: ", accountCredentialsInfo)
	log.Println("ERRRRRRRRRRRRRRROR: ", err)

	if err != nil {
		return err
	}

	userProfId, _ := uuid.Parse(accountCredentialsInfo.GetAccountCredentials().GetUserProfileId())
	userInfo.UserProfileID = userProfId
	userInfo.Username = accountCredentialsInfo.GetAccountCredentials().GetUsername()
	return nil
}

func (handler UserHandler) addUserProfileInfo(userInfo *dto.UserInfo, ctx context.Context) error {
	userProfileClient := communication.NewUserProfileClient(handler.userProfileServiceAddress)
	log.Println("userInfoId for userProfileInfo: ", userInfo.UserProfileID.String())

	userProfileInfo, err := userProfileClient.GetById(ctx, &user_profile.GetByIdRequest{Id: userInfo.UserProfileID.String()})
	log.Println("userProfileInfo from userProfileClient GetById: ", userProfileInfo)

	if err != nil {
		return err
	}

	userInfo.Name = userProfileInfo.UserProfile.Name
	userInfo.Surname = userProfileInfo.UserProfile.Surname
	userInfo.Email = userProfileInfo.UserProfile.Email

	userInfo.Address.Country = userProfileInfo.UserProfile.Address.Country
	userInfo.Address.City = userProfileInfo.UserProfile.Address.City
	userInfo.Address.Street = userProfileInfo.UserProfile.Address.Street
	userInfo.Address.StreetNumber = userProfileInfo.UserProfile.Address.StreetNumber

	return nil
}

func (handler UserHandler) CreateUser(ctx *gin.Context) {
	var user dto.CreateUser

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
		return
	}

	userProfileId, err := handler.CreateUserProfile(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
		return
	}

	_, err = handler.CreateAccountCredentials(&user, userProfileId)
	//	accountID, err := handler.CreateAccountCredentials(&user, userProfileId)
	if err != nil {
		deleteErr := handler.DeleteUserProfile(userProfileId)
		if deleteErr != nil {
			ctx.JSON(http.StatusInternalServerError, communication.NewErrorResponse(deleteErr.Error()))
			return
		}
		ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
		return
	}

	//err = handler.CreateNotificationConsent(accountID)
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
	//	return
	//}

	ctx.JSON(http.StatusCreated, communication.NewCreatedUserResponse(user.Username, userProfileId))
}

func (handler UserHandler) CreateUserProfile(user *dto.CreateUser) (uuid.UUID, error) {
	client := communication.NewUserProfileClient(handler.userProfileServiceAddress)
	response, err := client.Create(context.TODO(), &user_profile.CreateRequest{UserProfile: &user_profile.UserProfile{
		Name:    user.Name,
		Surname: user.Surname,
		Email:   user.Email,
		Address: &user_profile.Address{
			Country:      user.Address.Country,
			City:         user.Address.City,
			Street:       user.Address.Street,
			StreetNumber: user.Address.StreetNumber,
		},
	}})

	if err != nil {
		return uuid.UUID{}, err
	}

	userProfileId, _ := uuid.Parse(response.Id)
	return userProfileId, nil
}

func (handler UserHandler) DeleteUserProfile(id uuid.UUID) error {
	client := communication.NewUserProfileClient(handler.userProfileServiceAddress)
	_, err := client.Delete(context.TODO(), &user_profile.DeleteRequest{Id: id.String()})

	return err
}
func (handler UserHandler) CreateAccountCredentials(user *dto.CreateUser, userProfileId uuid.UUID) (uuid.UUID, error) {
	client := communication.NewAuthorizationClient(handler.authorizationServiceAddress)
	response, err := client.Create(context.TODO(), &authorization.CreateRequest{AccountCredentials: &authorization.CreateAccountCredentials{
		Username:      user.Username,
		Password:      user.Password,
		Role:          authorization.Role(user.Role),
		UserProfileId: userProfileId.String(),
	}})

	if err != nil {
		return uuid.UUID{}, err
	}
	id, err := uuid.Parse(response.Id)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, err
}

func (handler UserHandler) GetLoggedInUserInfo(ctx *gin.Context) {
	loggedInAccCredIdFromCtx := ctx.Keys["id"]
	if loggedInAccCredIdFromCtx == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "logged-in account credentials not provided"})
		return
	}
	loggedInAccCredId := fmt.Sprintf("%v", loggedInAccCredIdFromCtx)

	grpcCtx := createGrpcContextFromGinContext(ctx)

	authorizationClient := communication.NewAuthorizationClient(handler.authorizationServiceAddress)
	loggedInAccCred, err := authorizationClient.GetById(grpcCtx, &authorization.GetByIdRequest{
		Id: loggedInAccCredId,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not get logged-in account credentials"})
		return
	}

	var userInfo dto.UserInfo

	handler.addAccountCredentialsInfo(&userInfo, loggedInAccCred.GetAccountCredentials().GetUsername(), grpcCtx)
	handler.addUserProfileInfo(&userInfo, grpcCtx)

	ctx.JSON(http.StatusOK, userInfo)
}

func (handler UserHandler) CreateNotificationConsent(id uuid.UUID) error {
	client := communication.NewNotificationClient(handler.notificationServiceAddress)
	_, err := client.Create(context.TODO(), &notification.CreateRequest{
		UserProfileID:            id.String(),
		RequestMade:              true,
		ReservationCanceled:      true,
		HostRatingGiven:          true,
		AccommodationRatingGiven: true,
		ProminentHost:            true,
		HostResponded:            true,
	})

	if err != nil {
		return err
	}
	return nil

}

func createGrpcContextFromGinContext(ctx *gin.Context) context.Context {
	authHeader := ctx.GetHeader("Authorization")
	accessToken := authHeader[len("Bearer "):]
	md := metadata.New(map[string]string{"Authorization": accessToken})
	ctxGrpc := metadata.NewOutgoingContext(context.TODO(), md)
	return ctxGrpc
}
