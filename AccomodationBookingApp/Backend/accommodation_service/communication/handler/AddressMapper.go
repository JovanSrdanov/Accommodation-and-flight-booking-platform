package handler

import (
	"accommodation_service/domain/model"
	accommodation "common/proto/accommodation_service/generated"
)

type AddressMapper struct{}

func NewAddressMapper() *AddressMapper {
	return &AddressMapper{}
}

type IAddressMapper interface {
	mapToAddressModel(address *accommodation.Address) *model.Address
	mapToProto(address *model.Address) *accommodation.Address
}

func (a AddressMapper) mapToAddressModel(addressProto *accommodation.Address) model.Address {
	return model.Address{
		Country:      addressProto.Country,
		City:         addressProto.City,
		Street:       addressProto.Street,
		StreetNumber: addressProto.StreetNumber,
	}
}
func (a AddressMapper) mapToProto(addressModel *model.Address) *accommodation.Address {
	return &accommodation.Address{
		Country:      addressModel.Country,
		City:         addressModel.City,
		Street:       addressModel.Street,
		StreetNumber: addressModel.StreetNumber,
	}
}
