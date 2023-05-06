package handler

import (
	"api_gateway/communication"
	"api_gateway/domain/model"
	authorization "common/proto/authorization_service/generated"
	user_profile "common/proto/user_profile_service/generated"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type UserHandler struct {
	authorizationServiceAddress string
	userProfileServiceAddress   string
}

func NewUserHandler(authorizationServiceAddress string, userProfileServiceAddress string) *UserHandler {
	return &UserHandler{authorizationServiceAddress: authorizationServiceAddress,
		userProfileServiceAddress: userProfileServiceAddress}
}

func (handler UserHandler) Init(router *gin.RouterGroup) {
	userGroup := router.Group("/user")
	userGroup.GET("/:username/info", handler.GetUserInfo)
}

func (handler UserHandler) GetUserInfo(ctx *gin.Context) {
	username := ctx.Param("username")

	if username == "" {
		ctx.JSON(http.StatusBadRequest, "Username not provided")
		return
	}

	var userInfo model.UserInfo

	//TODO Error handling
	handler.addAccountCredentialsInfo(&userInfo, username)
	handler.addUserProfileInfo(&userInfo)

	ctx.JSON(http.StatusOK, userInfo)
}

func (handler UserHandler) addAccountCredentialsInfo(userInfo *model.UserInfo, username string) error {
	authorizationClient := communication.NewAuthorizationClient(handler.authorizationServiceAddress)
	accountCredentialsInfo, err := authorizationClient.GetByUsername(context.TODO(), &authorization.GetByUsernameRequest{Username: username})
	log.Println(accountCredentialsInfo)

	if err != nil {
		return err
	}

	userProfId, _ := uuid.Parse(accountCredentialsInfo.GetAccountCredentials().GetUserProfileId())
	userInfo.UserProfileID = userProfId
	userInfo.Username = accountCredentialsInfo.GetAccountCredentials().GetUsername()
	return nil
}

func (handler UserHandler) addUserProfileInfo(userInfo *model.UserInfo) error {
	userProfileClient := communication.NewUserProfileClient(handler.userProfileServiceAddress)

	userProfileInfo, err := userProfileClient.GetById(context.TODO(), &user_profile.GetByIdRequest{Id: userInfo.UserProfileID.String()})

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
