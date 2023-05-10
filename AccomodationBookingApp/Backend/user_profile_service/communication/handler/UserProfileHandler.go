package handler

import (
	authorization "common/proto/authorization_service/generated"
	user_profile "common/proto/user_profile_service/generated"
	"context"
	"fmt"
	"github.com/google/uuid"
	"user_profile_service/communication"
	"user_profile_service/domain/service"
)

type UserProfileHandler struct {
	user_profile.UnimplementedUserProfileServiceServer
	userProfileService          service.UserProfileService
	authorizationServiceAddress string
}

func NewUserProfileHandler(userProfileService service.UserProfileService, authorizationServerAddress string) *UserProfileHandler {
	return &UserProfileHandler{
		userProfileService:          userProfileService,
		authorizationServiceAddress: authorizationServerAddress,
	}
}

func (handler UserProfileHandler) Create(ctx context.Context, in *user_profile.CreateRequest) (*user_profile.CreateResponse, error) {
	mapper := NewUserProfileMapper()
	id, err := handler.userProfileService.Create(mapper.mapFromCreateRequest(in))

	if err != nil {
		return nil, err
	}
	return &user_profile.CreateResponse{
		Id: id.String(),
	}, nil
}

func (handler UserProfileHandler) Update(ctx context.Context, req *user_profile.UpdateRequest) (*user_profile.UpdateRequest, error) {
	// get account credentials id from logged-in user
	loggedInId, ok := ctx.Value("id").(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("failed to extract id and cast to UUID")
	}

	// get account credentials from acc cred microservice
	accCredClient := communication.NewAccountCredentialsClient(handler.authorizationServiceAddress)
	accCred, err := accCredClient.GetById(ctx, &authorization.GetByIdRequest{Id: loggedInId.String()})
	if err != nil {
		return nil, err
	}

	// get user info
	userInfoId, err := uuid.Parse(accCred.GetAccountCredentials().GetUserProfileId())
	if err != nil {
		return nil, err
	}

	userProfileMapper := NewUserProfileMapper()
	userUpdatedInfo, err := handler.userProfileService.Update(userInfoId, userProfileMapper.mapUpdateRequestToUpdateDto(req))

	return userProfileMapper.mapUpdateDtoToUpdateRequest(userUpdatedInfo), nil
}

func (handler UserProfileHandler) GetById(ctx context.Context, in *user_profile.GetByIdRequest) (*user_profile.GetByIdResponse, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, err
	}
	userProfile, err := handler.userProfileService.GetById(id)
	if err != nil {
		return nil, err
	}

	mapper := NewUserProfileMapper()

	return mapper.mapToGetByIdResponse(userProfile), nil
}

func (handler UserProfileHandler) Delete(ctx context.Context, in *user_profile.DeleteRequest) (*user_profile.DeleteResponse, error) {
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, err
	}
	err = handler.userProfileService.Delete(id)
	if err != nil {
		return nil, err
	}

	return &user_profile.DeleteResponse{Message: "User profile deleted"}, nil
}
func (handler UserProfileHandler) DeleteUser(ctx context.Context, in *user_profile.DeleteUserRequest) (*user_profile.DeleteResponse, error) {
	loggedInId, ok := ctx.Value("id").(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("failed to extract id and cast to UUID")
	}

	// get account credentials from acc cred microservice
	accCredClient := communication.NewAccountCredentialsClient(handler.authorizationServiceAddress)
	accCred, err := accCredClient.GetById(ctx, &authorization.GetByIdRequest{Id: loggedInId.String()})
	if err != nil {
		return nil, err
	}

	// get user info
	userProfileId, err := uuid.Parse(accCred.GetAccountCredentials().GetUserProfileId())
	if err != nil {
		return nil, err
	}

	err = handler.userProfileService.DeleteUser(userProfileId)

	return &user_profile.DeleteResponse{Message: "User deleted"}, err
}
