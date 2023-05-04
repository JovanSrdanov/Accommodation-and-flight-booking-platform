package handler

import (
	"authorization_service/domain/service"
	authorizationProto "common/proto/authorization_service/generated"
	"context"
	"fmt"
	"github.com/o1egl/paseto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
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
	// TODO Stefan: only for testing purposes, remove later
	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no metadata provided")
	}
	token := metaData["authorization"][0]
	var footerData map[string]interface{}
	if err := paseto.ParseFooter(token, &footerData); err != nil {
		return nil, fmt.Errorf("failed to parse token footer")
	}

	log.Println("Logged in user id: ", footerData["ID"])
	///////////////

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

	// TODO Stefan not sure if this is a good solution
	metadata.AppendToOutgoingContext(ctx, "authorization", accessToken)

	return res, nil
}
