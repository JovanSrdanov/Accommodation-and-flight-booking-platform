package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rating_service/domain/model"
	"rating_service/domain/repository"
)

type RatingService struct {
	ratingRepo repository.IRatingRepository
}

func NewRatingService(ratingRepo repository.IRatingRepository) *RatingService {
	return &RatingService{ratingRepo: ratingRepo}
}

func (service RatingService) RateAccommodation(guestId string, rating *model.Rating) error {
	//TODO Strahinja: Ovde proveriti da li guest sme da uradi rate
	return service.ratingRepo.RateAccommodation(guestId, rating)
}

func (service RatingService) GetRatingForAccommodation(id primitive.ObjectID) (model.RatingResponse, error) {
	return service.ratingRepo.GetRatingForAccommodation(id)
}
