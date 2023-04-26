package handler

import (
	"authorization_service/domain/model"
	authorization "common/proto/authorization_service/generated"
)

type AccountCredentialsMapper struct{}

func NewAccountCredentialsMapper() *AccountCredentialsMapper {
	return &AccountCredentialsMapper{}
}

type IAccountCredentialsMapper interface {
	mapFromCreateRequest(request *authorization.CreateRequest) *model.AccountCredentials
	mapToGetByEmailResponse(accCred *model.AccountCredentials) *authorization.GetByEmailResponse
}

func (a AccountCredentialsMapper) mapFromCreateRequest(request *authorization.CreateRequest) *model.AccountCredentials {
	return &model.AccountCredentials{
		Email:    request.GetAccountCredentials().Email,
		Password: request.GetAccountCredentials().Password,
		Role:     model.Role(request.GetAccountCredentials().Role),
	}
}

func (a AccountCredentialsMapper) mapToGetByEmailResponse(accCred *model.AccountCredentials) *authorization.GetByEmailResponse {
	return &authorization.GetByEmailResponse{
		AccountCredentials: &authorization.AccountCredentials{
			Id:       accCred.ID.String(),
			Email:    accCred.Email,
			Password: accCred.Password,
			Salt:     accCred.Salt,
			Role:     authorization.Role(accCred.Role),
		},
	}
}
