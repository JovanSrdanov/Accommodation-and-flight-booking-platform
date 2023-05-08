package handler

import (
	user_profile "common/proto/user_profile_service/generated"
	"user_profile_service/domain/model"
)

type UserProfileMapper struct{}

type IUserProfileMapper interface {
	mapFromCreateRequest(request *user_profile.CreateRequest) *model.UserProfile
	mapToGetByIdResponse(userProfile *model.UserProfile) *model.UserProfile
	mapUpdateRequestToUpdateDto(request *user_profile.UpdateRequest) *model.UpdateProfileDto
	mapUpdateDtoToUpdateRequest(updateProfile *model.UpdateProfileDto) *user_profile.UpdateRequest
}

func NewUserProfileMapper() *UserProfileMapper {
	return &UserProfileMapper{}
}

func (mapper UserProfileMapper) mapFromCreateRequest(request *user_profile.CreateRequest) *model.UserProfile {
	addressMapper := NewAddressMapper()
	return &model.UserProfile{
		Name:    request.UserProfile.Name,
		Surname: request.UserProfile.Surname,
		Email:   request.UserProfile.Email,
		Address: *addressMapper.mapToAddressModel(request.UserProfile.Address),
	}
}

func (mapper UserProfileMapper) mapToGetByIdResponse(userProfile *model.UserProfile) *user_profile.GetByIdResponse {
	addressMapper := NewAddressMapper()
	return &user_profile.GetByIdResponse{UserProfile: &user_profile.UserProfile{
		Id:      userProfile.ID.String(),
		Name:    userProfile.Name,
		Surname: userProfile.Surname,
		Email:   userProfile.Email,
		Address: addressMapper.mapToProto(&userProfile.Address),
	}}
}

func (mapper UserProfileMapper) mapUpdateRequestToUpdateDto(request *user_profile.UpdateRequest) *model.UpdateProfileDto {
	addressMapper := NewAddressMapper()
	return &model.UpdateProfileDto{
		Name:    request.Name,
		Surname: request.Name,
		Email:   request.Email,
		Address: *addressMapper.mapToAddressModel(request.Address),
	}
}

func (mapper UserProfileMapper) mapUpdateDtoToUpdateRequest(updateProfile *model.UpdateProfileDto) *user_profile.UpdateRequest {
	addressMapper := NewAddressMapper()

	return &user_profile.UpdateRequest{
		Name:    updateProfile.Name,
		Surname: updateProfile.Surname,
		Email:   updateProfile.Email,
		Address: addressMapper.mapToProto(&updateProfile.Address),
	}
}
