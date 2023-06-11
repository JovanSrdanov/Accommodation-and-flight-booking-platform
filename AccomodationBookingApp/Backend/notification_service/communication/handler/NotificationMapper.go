package handler

import (
	notification "common/proto/notification_service/generated"
	"github.com/google/uuid"
	"notification_service/domain/model"
)

func mapFromCreateRequest(request *notification.CreateRequest) *model.NotificationConsent {

	return &model.NotificationConsent{
		UserProfileID:            uuid.UUID{},
		RequestMade:              request.RequestMade,
		ReservationCanceled:      request.ReservationCanceled,
		HostRatingGiven:          request.HostRatingGiven,
		AccommodationRatingGiven: request.AccommodationRatingGiven,
		ProminentHost:            request.ProminentHost,
		HostResponded:            request.HostResponded,
	}
}
