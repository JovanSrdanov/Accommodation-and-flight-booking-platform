package repository

import (
	"accommodation_service/domain/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAccommodationRepository interface {
	Create(accommodation *model.Accommodation) (primitive.ObjectID, error)
	GetById(id primitive.ObjectID) (*model.Accommodation, error)
	GetAll() (model.Accommodations, error)
	GetAllMy(hostId string) (model.Accommodations, error)
	Update(accommodation *model.Accommodation) (*model.Accommodation, error)
	Delete(id primitive.ObjectID) error
	DeleteByHostId(id string) error
	GetAmenities() ([]string, error)
	SearchAccommodation(searchDto *model.SearchDto) (model.Accommodations, error)
}
