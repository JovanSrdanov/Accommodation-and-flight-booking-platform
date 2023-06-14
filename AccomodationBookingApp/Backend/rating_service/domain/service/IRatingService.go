package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rating_service/domain/model"
)

type IRatingService interface {
	RateAccommodation(guestId string, rating *model.Rating) error
	GetRatingForAccommodation(id primitive.ObjectID) (model.RatingResponse, error)
	GetRatingGuestGaveHost(hostID, guestID string) (float32, error)
	GetRatingGuestGaveAccommodation(accommodationID, guestID string) (float32, error)
}
