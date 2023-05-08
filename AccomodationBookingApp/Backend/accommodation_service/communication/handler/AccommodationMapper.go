package handler

import (
	"accommodation_service/domain/model"
	accommodation "common/proto/accommodation_service/generated"
)

type UserProfileMapper struct{}

type IUserProfileMapper interface {
	mapFromCreateRequest(request *accommodation.CreateRequest) *model.Accommodation
}

func NewAccommodationMapper() *UserProfileMapper {
	return &UserProfileMapper{}
}

func (mapper UserProfileMapper) mapFromCreateRequest(request *accommodation.CreateRequest) *model.Accommodation {
	return &model.Accommodation{
		Name: request.UserProfile.Name,
	}
}
