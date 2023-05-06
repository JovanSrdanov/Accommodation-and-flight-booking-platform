package handler

import (
	"api_gateway/communication"
	"api_gateway/domain/model"
	authorization "common/proto/authorization_service/generated"
	user_profile "common/proto/user_profile_service/generated"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"log"
	"net/http"
)

type UserInfoHandler struct {
	authorizationServiceAddress string
	userProfileServiceAddress   string
}

func NewUserInfoHandler(authorizationServiceAddress string, userProfileServiceAddress string) *UserInfoHandler {
	return &UserInfoHandler{authorizationServiceAddress: authorizationServiceAddress,
		userProfileServiceAddress: userProfileServiceAddress}
}

func (handler UserInfoHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/user/{username}/info", handler.GetUserInfo)
	if err != nil {
		panic(err)
	}
}

func (handler UserInfoHandler) GetUserInfo(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	username := pathParams["username"]

	var userInfo model.UserInfo

	//TODO Error handling
	handler.addAccountCredentialsInfo(&userInfo, username)
	handler.addUserProfileInfo(&userInfo)

	response, err := json.Marshal(userInfo)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response, _ = json.Marshal(err.Error())
		w.Write(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler UserInfoHandler) addAccountCredentialsInfo(userInfo *model.UserInfo, username string) error {
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

func (handler UserInfoHandler) addUserProfileInfo(userInfo *model.UserInfo) error {
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
