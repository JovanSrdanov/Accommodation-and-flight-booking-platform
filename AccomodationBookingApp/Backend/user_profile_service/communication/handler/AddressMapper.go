package handler

import (
	user_profile "common/proto/user_profile_service/generated"
	"github.com/google/uuid"
	"user_profile_service/domain/model"
)

type AddressMapper struct{}

func NewAddressMapper() *AddressMapper {
	return &AddressMapper{}
}

type IAddressMapper interface {
	mapToAddressModel(address *user_profile.Address) *model.Address
	mapToProto(address *model.Address) *user_profile.Address
}

func (a AddressMapper) mapToAddressModel(addressProto *user_profile.Address) *model.Address {
	uid, _ := uuid.Parse(addressProto.Id)
	return &model.Address{
		ID:           uid,
		Country:      addressProto.Country,
		City:         addressProto.City,
		Street:       addressProto.Street,
		StreetNumber: addressProto.StreetNumber,
	}
}
func (a AddressMapper) mapToProto(addressModel *model.Address) *user_profile.Address {
	return &user_profile.Address{
		Id:           addressModel.ID.String(),
		Country:      addressModel.Country,
		City:         addressModel.City,
		Street:       addressModel.Street,
		StreetNumber: addressModel.StreetNumber,
	}
}
