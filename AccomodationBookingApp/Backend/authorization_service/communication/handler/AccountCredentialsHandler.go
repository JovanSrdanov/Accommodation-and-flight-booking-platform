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
	accessToken, err := handler.accCredService.Login(req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	res := &authorizationProto.LoginResponse{AccessToken: accessToken}
	return res, nil
}

func (handler AccountCredentialsHandler) ChangeUsername(ctx context.Context, req *authorizationProto.ChangeUsernameRequest) (*emptypb.Empty, error) {
	loggedInId, ok := ctx.Value("id").(uuid.UUID)
	if !ok {
		return &emptypb.Empty{}, fmt.Errorf("failed to extract id and cast to UUID")
	}

	err := handler.accCredService.ChangeUsername(loggedInId, req.GetUsername())
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (handler AccountCredentialsHandler) ChangePassword(ctx context.Context, req *authorizationProto.ChangePasswordRequest) (*emptypb.Empty, error) {
	loggedInId, ok := ctx.Value("id").(uuid.UUID)
	if !ok {
		return &emptypb.Empty{}, fmt.Errorf("failed to extract id and cast to UUID")
	}

	err := handler.accCredService.ChangePassword(loggedInId, req.GetOldPassword(), req.GetNewPassword())
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (handler AccountCredentialsHandler) CheckIfDeleted(ctx context.Context, req *authorizationProto.CheckIfDeletedRequest) (*authorizationProto.CheckIfDeletedResponse, error) {
	loggedInId, ok := ctx.Value("id").(uuid.UUID)
	if !ok {
		return &authorizationProto.CheckIfDeletedResponse{
			Response: false,
		}, fmt.Errorf("failed to extract id and cast to UUID")
	}

	_, err := handler.accCredService.GetById(loggedInId)
	if err != nil {
		return &authorizationProto.CheckIfDeletedResponse{
			Response: true,
		}, nil
	} else {
		return &authorizationProto.CheckIfDeletedResponse{
			Response: false,
		}, nil
	}
}
