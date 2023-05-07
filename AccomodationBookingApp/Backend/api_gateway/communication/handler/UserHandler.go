package handler

import (
	"api_gateway/communication"
	"api_gateway/communication/middleware"
	"api_gateway/dto"
	"authorization_service/domain/model"
	"authorization_service/domain/token"
	authorization "common/proto/authorization_service/generated"
	user_profile "common/proto/user_profile_service/generated"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type UserHandler struct {
	authorizationServiceAddress string
	userProfileServiceAddress   string
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
		middleware.Authorization([]model.Role{model.Guest}),
		handler.GetUserInfo)
	userGroup.POST("", handler.CreateUser)
}

func (handler UserHandler) GetUserInfo(ctx *gin.Context) {
	username := ctx.Param("username")

	if username == "" {
		ctx.JSON(http.StatusBadRequest, "Username not provided")
		return
	}

	var userInfo dto.UserInfo

	//TODO errorMessage handling
	handler.addAccountCredentialsInfo(&userInfo, username)
	handler.addUserProfileInfo(&userInfo)

	ctx.JSON(http.StatusOK, userInfo)
}

func (handler UserHandler) addAccountCredentialsInfo(userInfo *dto.UserInfo, username string) error {
	authorizationClient := communication.NewAuthorizationClient(handler.authorizationServiceAddress)
	accountCredentialsInfo, err := authorizationClient.GetByUsername(context.TODO(), &authorization.GetByUsernameRequest{Username: username})

	if err != nil {
		return err
	}

	userProfId, _ := uuid.Parse(accountCredentialsInfo.GetAccountCredentials().GetUserProfileId())
	userInfo.UserProfileID = userProfId
	userInfo.Username = accountCredentialsInfo.GetAccountCredentials().GetUsername()
	return nil
}

func (handler UserHandler) addUserProfileInfo(userInfo *dto.UserInfo) error {
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

	err = handler.CreateAccountCredentials(&user, userProfileId)
	if err != nil {
		deleteErr := handler.DeleteUserProfile(userProfileId)
		if deleteErr != nil {
			ctx.JSON(http.StatusInternalServerError, communication.NewErrorResponse(deleteErr.Error()))
			return
		}
		ctx.JSON(http.StatusBadRequest, communication.NewErrorResponse(err.Error()))
		return
	}

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
func (handler UserHandler) CreateAccountCredentials(user *dto.CreateUser, userProfileId uuid.UUID) error {
	client := communication.NewAuthorizationClient(handler.authorizationServiceAddress)
	_, err := client.Create(context.TODO(), &authorization.CreateRequest{AccountCredentials: &authorization.CreateAccountCredentials{
		Username:      user.Username,
		Password:      user.Password,
		Role:          authorization.Role(user.Role),
		UserProfileId: userProfileId.String(),
	}})

	if err != nil {
		return err
	}

	return nil
}
