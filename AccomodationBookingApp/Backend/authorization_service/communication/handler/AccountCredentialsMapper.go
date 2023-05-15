package handler

import (
	"authorization_service/domain/model"
	authorization "common/proto/authorization_service/generated"
	"github.com/google/uuid"
)

type AccountCredentialsMapper struct{}

func NewAccountCredentialsMapper() *AccountCredentialsMapper {
	return &AccountCredentialsMapper{}
}

type IAccountCredentialsMapper interface {
	mapFromCreateRequest(request *authorization.CreateRequest) *model.AccountCredentials
	mapToGetByUsernameResponse(accCred *model.AccountCredentials) *authorization.GetByUsernameResponse
	mapToModelAccountCredentials(accCred *authorization.AccountCredentials) (*model.AccountCredentials, error)
}

func (a AccountCredentialsMapper) mapFromCreateRequest(request *authorization.CreateRequest) *model.AccountCredentials {
	//TODO handle error
	userProfileID, _ := uuid.Parse(request.GetAccountCredentials().UserProfileId)
	return &model.AccountCredentials{
		Username:      request.GetAccountCredentials().Username,
		Password:      request.GetAccountCredentials().Password,
		Role:          model.Role(request.GetAccountCredentials().Role),
		UserProfileID: userProfileID,
	}
}

func (a AccountCredentialsMapper) mapToGetByUsernameResponse(accCred *model.AccountCredentials) *authorization.GetByUsernameResponse {
	return &authorization.GetByUsernameResponse{
		AccountCredentials: &authorization.AccountCredentials{
			Id:            accCred.ID.String(),
			Username:      accCred.Username,
			Password:      accCred.Password,
			Role:          authorization.Role(accCred.Role),
			UserProfileId: accCred.UserProfileID.String(),
		},
	}
}

func (a AccountCredentialsMapper) mapToModelAccountCredentials(accCred *authorization.AccountCredentials) (*model.AccountCredentials, error) {
	accID, err := uuid.Parse(accCred.GetId())
	if err != nil {
		return nil, err
	}

	userProfileId, err := uuid.Parse(accCred.GetUserProfileId())
	if err != nil {
		return nil, err
	}

	return &model.AccountCredentials{
		ID:            accID,
		Username:      accCred.GetUsername(),
		Password:      accCred.GetPassword(),
		Role:          model.Role(accCred.GetRole()),
		UserProfileID: userProfileId,
	}, nil
}
