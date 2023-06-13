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

func (service RatingService) GetRecommendedAccommodations(guestId string) (model.RecommendedAccommodations, error) {
	return service.ratingRepo.GetRecommendedAccommodations(guestId)
}

func (service RatingService) DeleteRatingForAccommodation(accommodationId string, guestId string) (string, error) {
	return service.ratingRepo.DeleteRatingForAccommodation(accommodationId, guestId)
}

func (service RatingService) RateHost(rating *model.RateHostDto) error {
	return service.ratingRepo.RateHost(rating)
}

func (service RatingService) GetRatingForHost(hostId string) (model.HostRatingResponse, error) {
	return service.ratingRepo.GetRatingForHost(hostId)
}

func (service RatingService) DeleteRatingForHost(hostId string, guestId string) (string, error) {
	return service.ratingRepo.DeleteRatingForHost(hostId, guestId)
}
