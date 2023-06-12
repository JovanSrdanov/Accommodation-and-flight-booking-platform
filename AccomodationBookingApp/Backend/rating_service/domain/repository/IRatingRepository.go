package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rating_service/domain/model"
)

type IRatingRepository interface {
	RateAccommodation(rating *model.Rating) error
	GetRatingForAccommodation(id primitive.ObjectID) (model.RatingResponse, error)
}
