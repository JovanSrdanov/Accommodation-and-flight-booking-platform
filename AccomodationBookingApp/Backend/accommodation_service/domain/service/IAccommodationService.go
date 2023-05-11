package service

import (
	"accommodation_service/domain/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUserProfileService interface {
	Create(userProf *model.Accommodation) (primitive.ObjectID, error)
	GetById(id primitive.ObjectID) (*model.Accommodation, error)
	GetAll() (model.Accommodations, error)
	GetAllMy(hostId string) (model.Accommodations, error)
	Update(id primitive.ObjectID, dto *model.Accommodation) (*model.Accommodation, error)
	Delete(id primitive.ObjectID) error
	GetAmenities() ([]string, error)
}
