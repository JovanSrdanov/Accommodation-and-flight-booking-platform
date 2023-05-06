package handler

import (
	"authorization_service/domain/model"
	"authorization_service/domain/service"
	authorizationProto "common/proto/authorization_service/generated"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
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
	userProfileId, err := uuid.Parse(request.GetAccountCredentials().GetUserProfileId())
	if err != nil {
		return nil, err
	}

	accCred, err := model.NewAccountCredentials(request.GetAccountCredentials().Username,
		request.GetAccountCredentials().Password, "", model.Role(request.GetAccountCredentials().Role), userProfileId)
	if err != nil {
		return nil, err
	}

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
	log.Println("METADATAAAAAAAA: ", metaData)
	values := metaData["authorization"]
	token := strings.TrimPrefix(values[0], "Bearer ")
	var footerData map[string]interface{}
	if err := paseto.ParseFooter(token, &footerData); err != nil {
		return nil, fmt.Errorf("failed to parse token footer")
	}

	log.Println("Logged in user username: ", footerData["Username"])
	/////////////

	result, err := handler.accCredService.GetByUsername(request.Username)
	if err != nil {
		return nil, err
	}
	mapper := NewAccountCredentialsMapper()
	return mapper.mapToGetByUsernameResponse(result), nil
}

// Login is an unary rpc
func (handler AccountCredentialsHandler) Login(ctx context.Context, req *authorizationProto.LoginRequest) (*authorizationProto.LoginResponse, error) {
	accessToken, err := handler.accCredService.Login(req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	res := &authorizationProto.LoginResponse{AccessToken: accessToken}
	return res, nil
}
