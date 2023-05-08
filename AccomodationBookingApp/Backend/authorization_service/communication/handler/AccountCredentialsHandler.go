package handler

import (
	"authorization_service/domain/model"
	"authorization_service/domain/service"
	authorizationProto "common/proto/authorization_service/generated"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
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
		request.GetAccountCredentials().Password, model.Role(request.GetAccountCredentials().Role), userProfileId)
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
	//loggedInUserId, ok := ctx.Value("id").(uuid.UUID)
	//if !ok {
	//	return nil, fmt.Errorf("failed to extract id and cast to UUID")
	//}
	//
	//loggedInUserRole, err := token.ExtractTokenInfoFromContext(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//
	//log.Printf("Logged in user username: %s, role: %s", loggedInUserId, loggedInUserRole)
	/////////////

	result, err := handler.accCredService.GetByUsername(request.Username)
	if err != nil {
		return nil, err
	}
	mapper := NewAccountCredentialsMapper()
	return mapper.mapToGetByUsernameResponse(result), nil
}

func (handler AccountCredentialsHandler) GetById(ctx context.Context, req *authorizationProto.GetByIdRequest) (*authorizationProto.GetByUsernameResponse, error) {
	accId, err := uuid.Parse(req.GetId())

	result, err := handler.accCredService.GetById(accId)
	if err != nil {
		return nil, err
	}

	mapper := NewAccountCredentialsMapper()
	return mapper.mapToGetByUsernameResponse(result), nil
}

// Login is an unary rpc
func (handler AccountCredentialsHandler) Login(ctx context.Context, req *authorizationProto.LoginRequest) (*authorizationProto.LoginResponse, error) {
	accessToken, role, err := handler.accCredService.Login(req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	res := &authorizationProto.LoginResponse{AccessToken: accessToken, Role: authorizationProto.Role(role)}
	return res, nil
}

func (handler AccountCredentialsHandler) Update(ctx context.Context, req *authorizationProto.UpdateRequest) (*emptypb.Empty, error) {
	loggedInId, ok := ctx.Value("id").(uuid.UUID)
	if !ok {
		return &emptypb.Empty{}, fmt.Errorf("failed to extract id and cast to UUID")
	}

	//loggedInIdInfoAsString := loggedInId.(string)
	//loggedInId, err := uuid.Parse(loggedInIdInfoAsString)
	//if err != nil {
	//	return &emptypb.Empty{}, err
	//}

	err := handler.accCredService.Update(loggedInId, req.GetUsername(), req.GetPassword())
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
