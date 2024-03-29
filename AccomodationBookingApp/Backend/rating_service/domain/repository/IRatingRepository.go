package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rating_service/domain/model"
)

type IRatingRepository interface {
	GetRecommendedAccommodations(guestId string) ([]model.RecommendedAccommodation, error)

	RateAccommodation(guestId string, rating *model.Rating) error
	GetRatingForAccommodation(id primitive.ObjectID) (model.RatingResponse, error)
	DeleteRatingForAccommodation(accommodationId string, guestId string) (string, error)

	RateHost(rating *model.RateHostDto) error
	GetRatingForHost(hostId string) (model.HostRatingResponse, error)
	DeleteRatingForHost(hostId string, guestId string) (string, error)

	CalculateRatingForAccommodation(accommodationId string) (model.SimpleRatingResponse, error)
	CalculateRatingForHost(hostId string) (model.SimpleHostRatingResponse, error)

	GetRatingGuestGaveHost(hostID, guestID string) (float32, error)
	GetRatingGuestGaveAccommodation(accommodationID, guestID string) (float32, error)
}
