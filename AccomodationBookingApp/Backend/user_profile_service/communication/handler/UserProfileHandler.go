package handler

import (
	user_profile "common/proto/user_profile_service/generated"
	"context"
	"github.com/google/uuid"
	"user_profile_service/domain/service"
)

type UserProfileHandler struct {
	user_profile.UnimplementedUserProfileServiceServer
	userProfileService service.UserProfileService
}

func NewUserProfileHandler(userProfileService service.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{userProfileService: userProfileService}
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
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, err
	}
	err = handler.userProfileService.Delete(id)

	return &user_profile.DeleteResponse{}, err
}
