package handler

import (
	rating "common/proto/rating_service/generated"
	"context"
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
	mapper := NewRatingMapper()
	err := handler.ratingService.RateAccommodation(mapper.mapFromRateAccommodationRequest(in))
	if err != nil {
		return &rating.RateAccommodationResponse{}, err
	}

	return &rating.RateAccommodationResponse{}, nil
}

func (handler RatingHandler) GetRatingForAccommodation(ctx context.Context, in *rating.RatingForAccommodationRequest) (*rating.RatingForAccommodationResponse, error) {
	return &rating.RatingForAccommodationResponse{}, nil
}
