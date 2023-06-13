package handler

import (
	rating "common/proto/rating_service/generated"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rating_service/domain/model"
	"time"
)

type RatingMapper struct {
}

func NewRatingMapper() *RatingMapper {
	return &RatingMapper{}
}

func (mapper RatingMapper) mapFromRateAccommodationRequest(request *rating.RateAccommodationRequest) *model.Rating {
	accommodationId, _ := primitive.ObjectIDFromHex(request.Rating.AccommodationId)

	return &model.Rating{
		AccommodationId: accommodationId,
		GuestId:         "iz jwt-a",
		Rating:          request.Rating.Rating,
		Date:            time.Now(),
	}
}

func (mapper RatingMapper) mapToRatingForAccommodationResponse(in *model.RatingResponse) *rating.RatingForAccommodationResponse {
	ratingsProto := make([]*rating.Rating, 0)

	for _, val := range in.Ratings {
		ratingsProto = append(ratingsProto, &rating.Rating{
			GuestId: val.GuestId,
			Date:    val.Date.Format("2006-01-02"),
			Rating:  val.Rating,
		})
	}

	return &rating.RatingForAccommodationResponse{Rating: &rating.AccommodationRating{
		AvgRating:       in.AvgRating,
		AccommodationId: in.AccommodationId,
		Ratings:         ratingsProto,
	}}
}

func (mapper RatingMapper) mapToRatingForHostResponse(in *model.HostRatingResponse) *rating.RatingForHostResponse {
	ratingsProto := make([]*rating.Rating, 0)

	for _, val := range in.Ratings {
		ratingsProto = append(ratingsProto, &rating.Rating{
			GuestId: val.GuestId,
			Date:    val.Date.Format("2006-01-02"),
			Rating:  val.Rating,
		})
	}

	return &rating.RatingForHostResponse{Rating: &rating.HostRating{
		AvgRating: in.AvgRating,
		HostId:    in.HostId,
		Ratings:   ratingsProto,
	}}
}

func (mapper RatingMapper) mapToRecommendedAccommodationsResponse(in *model.RecommendedAccommodations) *rating.RecommendedAccommodationsResponse {
	slice := make([]string, 0)

	for _, val := range in.AccommodationsIds {
		slice = append(slice, val.Hex())
	}

	return &rating.RecommendedAccommodationsResponse{AccommodationId: slice}
}
