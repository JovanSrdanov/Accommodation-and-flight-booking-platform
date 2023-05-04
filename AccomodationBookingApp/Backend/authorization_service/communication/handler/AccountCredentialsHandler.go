package handler

import (
	"authorization_service/domain/service"
	authorizationProto "common/proto/authorization_service/generated"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountCredentialsHandler struct {
	authorizationProto.UnimplementedAuthorizationServiceServer
	accCredService service.IAccountCredentialsService
}

func NewAccountCredentialsHandler(accCredService service.IAccountCredentialsService) *AccountCredentialsHandler {
	return &AccountCredentialsHandler{
		accCredService: accCredService,
	}
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
func (handler AccountCredentialsHandler) GetByUsername(ctx context.Context, request *authorizationProto.GetByUsernameRequest) (*authorizationProto.GetByUsernameResponse, error) {
	result, err := handler.accCredService.GetByUsername(request.Username)
	if err != nil {
		return nil, err
	}
	mapper := NewAccountCredentialsMapper()
	return mapper.mapToGetByUsernameResponse(result), nil
}

// Login is an unary rpc
func (handler AccountCredentialsHandler) Login(ctx context.Context, req *authorizationProto.LoginRequest) (*authorizationProto.LoginResponse, error) {
	accCred, err := handler.accCredService.GetByUsername(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}

	accessToken, err := handler.accCredService.Login(accCred)
	if err != nil {
		return nil, err
	}

	res := &authorizationProto.LoginResponse{AccessToken: accessToken}

	return res, nil
}
