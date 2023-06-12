package handler

import (
	rating "common/proto/rating_service/generated"
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rating_service/domain/service"
)

type RatingHandler struct {
	rating.UnimplementedRatingServiceServer
	ratingService service.RatingService
}

func NewRatingHandler(ratingService service.RatingService) *RatingHandler {
	return &RatingHandler{ratingService: ratingService}
}

func (handler RatingHandler) RateAccommodation(ctx context.Context, in *rating.RateAccommodationRequest) (*rating.RateAccommodationResponse, error) {
	loggedInId := ctx.Value("id")

	mapper := NewRatingMapper()
	err := handler.ratingService.RateAccommodation(loggedInId.(uuid.UUID).String(), mapper.mapFromRateAccommodationRequest(in))
	if err != nil {
		return &rating.RateAccommodationResponse{}, err
	}

	return &rating.RateAccommodationResponse{}, nil
}

func (handler RatingHandler) GetRatingForAccommodation(ctx context.Context, in *rating.RatingForAccommodationRequest) (*rating.RatingForAccommodationResponse, error) {
	mapper := NewRatingMapper()

	accommodationId, _ := primitive.ObjectIDFromHex(in.AccommodationId)
	res, err := handler.ratingService.GetRatingForAccommodation(accommodationId)
	if err != nil {
		return &rating.RatingForAccommodationResponse{}, err
	}
	return mapper.mapToRatingForAccommodationResponse(&res), nil
}

func (handler RatingHandler) GetRecommendedAccommodations(ctx context.Context, in *rating.RecommendedAccommodationsRequest) (*rating.RecommendedAccommodationsResponse, error) {
	mapper := NewRatingMapper()

	res, err := handler.ratingService.GetRecommendedAccommodations(in.GuestId)
	if err != nil {
		return &rating.RecommendedAccommodationsResponse{}, err
	}

	return mapper.mapToRecommendedAccommodationsResponse(&res), nil
}
