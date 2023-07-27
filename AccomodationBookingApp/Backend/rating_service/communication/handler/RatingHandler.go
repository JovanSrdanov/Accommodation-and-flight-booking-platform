package handler

import (
	"common/NotificationMessaging"
	rating "common/proto/rating_service/generated"
	"common/saga/messaging"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"rating_service/domain/model"
	"rating_service/domain/service"
	"rating_service/utils"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RatingHandler struct {
	rating.UnimplementedRatingServiceServer
	ratingService service.RatingService
	publisher     messaging.Publisher
}

func NewRatingHandler(ratingService service.RatingService, publisher messaging.Publisher) *RatingHandler {
	return &RatingHandler{ratingService: ratingService,
		publisher: publisher}
}

func (handler RatingHandler) RateAccommodation(ctx context.Context, in *rating.RateAccommodationRequest) (*rating.EmptyResponse, error) {
	loggedInId, err := utils.GetTokenInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract id")
	}

	mapper := NewRatingMapper()
	err = handler.ratingService.RateAccommodation(loggedInId.String(), mapper.mapFromRateAccommodationRequest(in))
	if err != nil {
		return &rating.EmptyResponse{}, err
	}

	//GRPC da vidis ciji je smestaj

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
	loggedInId, err := utils.GetTokenInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract id")
	}

	message, err := handler.ratingService.DeleteRatingForAccommodation(in.AccommodationId, loggedInId.String())
	if err != nil {
		return nil, err
	}

	return &rating.SimpleResponse{Message: message}, nil
}

func (handler RatingHandler) RateHost(ctx context.Context, in *rating.RateHostRequest) (*rating.EmptyResponse, error) {
	oldProminenet, err := prominentHostHttp(in.Rating.HostId)
	if err != nil {
		return nil, err
	}

	loggedInId, err := utils.GetTokenInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract id")
	}

	err = handler.ratingService.RateHost(&model.RateHostDto{
		HostId:  in.Rating.HostId,
		GuestId: loggedInId.String(),
		Rating:  in.Rating.Rating,
		Date:    time.Now(),
	})
	if err != nil {
		return nil, err
	}

	newProminent, err := prominentHostHttp(in.Rating.HostId)
	if err != nil {
		return nil, err
	}

	if oldProminenet != newProminent {
		accountID, err := uuid.Parse(in.Rating.HostId)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		MFN := ""
		if newProminent {
			MFN = "You are now a prominent host"
		} else {
			MFN = "You are no longer prominent host"
		}

		message := NotificationMessaging.NotificationMessage{
			MessageType:            "ProminentHost",
			MessageForNotification: MFN,
			AccountID:              accountID,
		}
		handler.publisher.Publish(message)
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
	oldProminenet, err := prominentHostHttp(in.HostId)
	if err != nil {
		return nil, err
	}

	loggedInId, err := utils.GetTokenInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract id")
	}

	message, err := handler.ratingService.DeleteRatingForHost(in.HostId, loggedInId.String())
	if err != nil {
		return nil, err
	}
	newProminent, err := prominentHostHttp(in.HostId)
	if err != nil {
		return nil, err
	}
	if oldProminenet != newProminent {
		accountID, err := uuid.Parse(in.HostId)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		MFN := ""
		if newProminent {
			MFN = "You are now a prominent host"
		} else {
			MFN = "You are no longer prominent host"
		}

		message := NotificationMessaging.NotificationMessage{
			MessageType:            "ProminentHost",
			MessageForNotification: MFN,
			AccountID:              accountID,
		}
		handler.publisher.Publish(message)
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

func (handler RatingHandler) GetRatingGuestGaveHost(ctx context.Context, in *rating.GetRatingGuestGaveHostRequest) (*rating.GetRatingGuestGaveHostResponse, error) {
	hostRating, err := handler.ratingService.GetRatingGuestGaveHost(in.HostId, in.GuestId)
	if err != nil {
		return nil, err
	}

	return &rating.GetRatingGuestGaveHostResponse{Rating: hostRating}, nil
}

func (handler RatingHandler) GetRatingGuestGaveAccommodation(ctx context.Context, in *rating.GetRatingGuestGaveAccommodationRequest) (*rating.GetRatingGuestGaveAccommodationResponse, error) {
	hostRating, err := handler.ratingService.GetRatingGuestGaveAccommodation(in.AccommodationId, in.GuestId)
	if err != nil {
		return nil, err
	}

	return &rating.GetRatingGuestGaveAccommodationResponse{Rating: hostRating}, nil
}

func prominentHostHttp(hostId string) (bool, error) {
	client := &http.Client{}

	apiHost := os.Getenv("API_GATEWAY_HOST")
	req, err := http.NewRequest("GET", "http://"+apiHost+":8000/api-2/accommodation/prominent-host/"+hostId, nil)
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if string(body) == "true" {
		return true, nil
	} else {
		return false, nil
	}
}
