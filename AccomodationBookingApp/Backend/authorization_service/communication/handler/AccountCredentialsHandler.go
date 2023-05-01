package handler

import (
	"authorization_service/domain/service"
	authorizationProto "common/proto/authorization_service/generated"
	"context"
)

type AccountCredentialsHandler struct {
	authorizationProto.UnimplementedAuthorizationServiceServer
	accCredService service.IAccountCredentialsService
}

func NewAccountCredentialsHandler(accCredService service.IAccountCredentialsService) *AccountCredentialsHandler {
	return &AccountCredentialsHandler{accCredService: accCredService}
}

func (handler AccountCredentialsHandler) Create(ctx context.Context, request *authorizationProto.CreateRequest) (*authorizationProto.CreateResponse, error) {
	mapper := NewAccountCredentialsMapper()
	accCred := mapper.mapFromCreateRequest(request)
	result, err := handler.accCredService.Create(accCred)
	if err != nil {
		return nil, err
	}

	return &authorizationProto.CreateResponse{
		Id: result.String(),
	}, nil
}
func (handler AccountCredentialsHandler) GetByEmail(ctx context.Context, request *authorizationProto.GetByEmailRequest) (*authorizationProto.GetByEmailResponse, error) {
	result, err := handler.accCredService.GetByEmail(request.Email)
	if err != nil {
		return nil, err
	}
	mapper := NewAccountCredentialsMapper()
	return mapper.mapToGetByEmailResponse(result), nil
}
