package handler

import (
	"accommodation_service/domain/model"
	accommodation "common/proto/accommodation_service/generated"
)

type AccommodationMapper struct{}

type IAccommodationMapper interface {
	mapFromCreateRequest(request *accommodation.CreateRequest) *model.Accommodation
}

func NewAccommodationMapper() *AccommodationMapper {
	return &AccommodationMapper{}
}

func (mapper AccommodationMapper) mapFromCreateRequest(request *accommodation.CreateRequest) *model.Accommodation {
	return &model.Accommodation{
		Name:      request.Accommodation.Name,
		Location:  request.Accommodation.Location,
		MinGuests: request.Accommodation.MinGuests,
		MaxGuests: request.Accommodation.MaxGuests,
		Amenities: request.Accommodation.Amenities,
	}
}
