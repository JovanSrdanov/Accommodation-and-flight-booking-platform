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
	addressMapper := NewAddressMapper()
	return &model.Accommodation{
		Name:      request.Accommodation.Name,
		Address:   addressMapper.mapToAddressModel(request.Accommodation.Address),
		MinGuests: request.Accommodation.MinGuests,
		MaxGuests: request.Accommodation.MaxGuests,
		Amenities: request.Accommodation.Amenities,
		Images:    request.Accommodation.Images,
	}
}

func (mapper AccommodationMapper) mapToGetByIdResponse(model *model.Accommodation) *accommodation.GetByIdResponse {
	addressMapper := NewAddressMapper()
	return &accommodation.GetByIdResponse{
		Accommodation: &accommodation.Accommodation{
			Name:      model.Name,
			Address:   addressMapper.mapToProto(&model.Address),
			MinGuests: model.MinGuests,
			MaxGuests: model.MaxGuests,
			Amenities: model.Amenities,
			Images:    model.Images,
			HostId:    model.HostId,
		},
	}
}

func (mapper AccommodationMapper) mapToGetAllResponse(model model.Accommodations) *accommodation.GetAllResponse {
	addressMapper := NewAddressMapper()
	accommodationsProto := make([]*accommodation.AccommodationFull, 0)

	for _, value := range model {
		accommodationsProto = append(accommodationsProto, &accommodation.AccommodationFull{
			Id:        value.ID.String(),
			Name:      value.Name,
			Address:   addressMapper.mapToProto(&value.Address),
			MinGuests: value.MinGuests,
			MaxGuests: value.MaxGuests,
			Amenities: value.Amenities,
			Images:    value.Images,
			HostId:    value.HostId,
		})
	}

	return &accommodation.GetAllResponse{
		Accommodation: accommodationsProto,
	}
}
