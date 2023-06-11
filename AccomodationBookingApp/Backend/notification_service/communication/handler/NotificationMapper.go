package handler

import (
	notification "common/proto/notification_service/generated"
	"github.com/google/uuid"
	"notification_service/domain/model"
)

func mapFromCreateRequest(request *notification.CreateRequest) (*model.NotificationConsent, error) {

	id, err := uuid.Parse(request.UserProfileID)
	if err != nil {
		return nil, err
	}
	return &model.NotificationConsent{
		UserProfileID:            id,
		RequestMade:              request.RequestMade,
		ReservationCanceled:      request.ReservationCanceled,
		HostRatingGiven:          request.HostRatingGiven,
		AccommodationRatingGiven: request.AccommodationRatingGiven,
		ProminentHost:            request.ProminentHost,
		HostResponded:            request.HostResponded,
	}, nil
}
