package service

import (
	"common/NotificationMessaging"
	accommodation "common/proto/accommodation_service/generated"
	"common/saga/messaging"
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	client2 "rating_service/communication/client"
	"rating_service/domain/model"
	"rating_service/domain/repository"
)

type RatingService struct {
	ratingRepo repository.IRatingRepository
	publisher  messaging.Publisher
}

func NewRatingService(ratingRepo repository.IRatingRepository, publisher messaging.Publisher) *RatingService {
	return &RatingService{ratingRepo: ratingRepo,
		publisher: publisher}
}

func (service RatingService) RateAccommodation(guestId string, rating *model.Rating) error {
	res := service.ratingRepo.RateAccommodation(guestId, rating)
	if res != nil {
		return res
	}

	hostId, err2 := getHostIdForAccommodationId(rating.AccommodationId.Hex())
	if err2 != nil {
		return err2
	}
	log.Println(hostId)

	accountID, err := uuid.Parse(hostId)
	if err != nil {
		log.Fatal(err)
		return err
	}

	message := NotificationMessaging.NotificationMessage{
		MessageType:            "AccommodationRatingGiven",
		MessageForNotification: "Your accommodation has been rated",
		AccountID:              accountID,
	}
	service.publisher.Publish(message)

	return res
}

func getHostIdForAccommodationId(accommodationId string) (string, error) {
	accommodationHost := os.Getenv("ACCOMMODATION_SERVICE_HOST")
	accommodationPort := os.Getenv("ACCOMMODATION_SERVICE_PORT")
	client := client2.NewAccommodationClient(accommodationHost + ":" + accommodationPort)

	fullAccInfo, err := client.GetById(context.TODO(), &accommodation.GetByIdRequest{Id: accommodationId})
	if err != nil {
		return "", err
	}

	hostId := fullAccInfo.Accommodation.HostId
	return hostId, nil
}

func (service RatingService) GetRatingForAccommodation(id primitive.ObjectID) (model.RatingResponse, error) {
	return service.ratingRepo.GetRatingForAccommodation(id)
}

func (service RatingService) GetRecommendedAccommodations(guestId string) ([]model.RecommendedAccommodation, error) {
	return service.ratingRepo.GetRecommendedAccommodations(guestId)
}

func (service RatingService) DeleteRatingForAccommodation(accommodationId string, guestId string) (string, error) {
	errorMess, err := service.ratingRepo.DeleteRatingForAccommodation(accommodationId, guestId)
	if err != nil {
		return "Big puc kod brisanja accommodation rating " + errorMess, err
	}

	hostId, err := getHostIdForAccommodationId(accommodationId)
	if err != nil {
		return "Big puc kod brisanja accommodation rating", err
	}

	log.Println(hostId)
	accountID, err := uuid.Parse(hostId)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	message := NotificationMessaging.NotificationMessage{
		MessageType:            "AccommodationRatingGiven",
		MessageForNotification: "A rating for your accommodation has been deleted",
		AccountID:              accountID,
	}
	service.publisher.Publish(message)

	return errorMess, nil
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

func (service RatingService) CalculateRatingForAccommodation(accommodationId string) (model.SimpleRatingResponse, error) {
	return service.ratingRepo.CalculateRatingForAccommodation(accommodationId)
}

func (service RatingService) CalculateRatingForHost(hostId string) (model.SimpleHostRatingResponse, error) {
	return service.ratingRepo.CalculateRatingForHost(hostId)
}

func (service RatingService) GetRatingGuestGaveHost(hostID, guestID string) (float32, error) {
	return service.ratingRepo.GetRatingGuestGaveHost(hostID, guestID)
}

func (service RatingService) GetRatingGuestGaveAccommodation(accommodationID, guestID string) (float32, error) {
	return service.ratingRepo.GetRatingGuestGaveAccommodation(accommodationID, guestID)
}
