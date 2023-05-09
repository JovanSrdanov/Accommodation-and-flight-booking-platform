package repository

import (
	"accommodation_service/domain/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAccommodationRepository interface {
	Create(accommodation *model.Accommodation) (primitive.ObjectID, error)
	GetById(id primitive.ObjectID) (*model.Accommodation, error)
	GetAll() (model.Accommodations, error)
	Update(accommodation *model.Accommodation) (*model.Accommodation, error)
	Delete(id primitive.ObjectID) error
}
