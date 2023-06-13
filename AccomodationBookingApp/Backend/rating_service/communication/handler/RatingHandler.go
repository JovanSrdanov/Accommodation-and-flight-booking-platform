package handler

import (
	rating "common/proto/rating_service/generated"
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rating_service/domain/model"
	"rating_service/domain/service"
	"time"
)

type RatingHandler struct {
	rating.UnimplementedRatingServiceServer
	ratingService service.RatingService
}

func NewRatingHandler(ratingService service.RatingService) *RatingHandler {
	return &RatingHandler{ratingService: ratingService}
}

func (handler RatingHandler) RateAccommodation(ctx context.Context, in *rating.RateAccommodationRequest) (*rating.EmptyResponse, error) {
	loggedInId := ctx.Value("id")

	mapper := NewRatingMapper()
	err := handler.ratingService.RateAccommodation(loggedInId.(uuid.UUID).String(), mapper.mapFromRateAccommodationRequest(in))
	if err != nil {
		return &rating.EmptyResponse{}, err
	}

	return &rating.EmptyResponse{}, nil
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

	return mapper.mapToRecommendedAccommodationsResponse(res), nil
}

func (handler RatingHandler) DeleteRatingForAccommodation(ctx context.Context, in *rating.RatingForAccommodationRequest) (*rating.SimpleResponse, error) {
	loggedInId := ctx.Value("id").(uuid.UUID).String()
	message, err := handler.ratingService.DeleteRatingForAccommodation(in.AccommodationId, loggedInId)
	if err != nil {
		return nil, err
	}

	return &rating.SimpleResponse{Message: message}, nil
}

func (handler RatingHandler) RateHost(ctx context.Context, in *rating.RateHostRequest) (*rating.EmptyResponse, error) {
	loggedInId := ctx.Value("id").(uuid.UUID).String()
	err := handler.ratingService.RateHost(&model.RateHostDto{
		HostId:  in.Rating.HostId,
		GuestId: loggedInId,
		Rating:  in.Rating.Rating,
		Date:    time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &rating.EmptyResponse{}, nil
}

func (handler RatingHandler) GetRatingForHost(ctx context.Context, in *rating.RatingForHostRequest) (*rating.RatingForHostResponse, error) {
	mapper := NewRatingMapper()

	res, err := handler.ratingService.GetRatingForHost(in.HostId)
	if err != nil {
		return nil, err
	}

	return mapper.mapToRatingForHostResponse(&res), nil
}

func (handler RatingHandler) DeleteRatingForHost(ctx context.Context, in *rating.RatingForHostRequest) (*rating.SimpleResponse, error) {
	loggedInId := ctx.Value("id").(uuid.UUID).String()
	message, err := handler.ratingService.DeleteRatingForHost(in.HostId, loggedInId)
	if err != nil {
		return nil, err
	}

	return &rating.SimpleResponse{Message: message}, nil
}

func (handler RatingHandler) CalculateRatingForHost(ctx context.Context, in *rating.RatingForHostRequest) (*rating.CalculateRatingForHostResponse, error) {
	hostRating, err := handler.ratingService.CalculateRatingForHost(in.HostId)
	if err != nil {
		return nil, err
	}

	return &rating.CalculateRatingForHostResponse{Rating: &rating.SimpleHostRating{
		AvgRating: hostRating.AvgRating,
		HostId:    hostRating.HostId,
	}}, nil
}

func (handler RatingHandler) CalculateRatingForAccommodation(ctx context.Context, in *rating.RatingForAccommodationRequest) (*rating.CalculateRatingForAccommodationResponse, error) {
	accommodationRating, err := handler.ratingService.CalculateRatingForAccommodation(in.AccommodationId)
	if err != nil {
		return nil, err
	}

	return &rating.CalculateRatingForAccommodationResponse{Rating: &rating.SimpleAccommodationRating{
		AvgRating:       accommodationRating.AvgRating,
		AccommodationId: accommodationRating.AccommodationId,
	}}, nil
}
