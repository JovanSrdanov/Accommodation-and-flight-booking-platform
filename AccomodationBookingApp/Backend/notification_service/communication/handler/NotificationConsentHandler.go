package handler

import (
	"common/NotificationMessaging"
	notification "common/proto/notification_service/generated"
	"common/saga/messaging"
	"context"
	"fmt"
	"notification_service/domain/service"
	"notification_service/utils"
)

type NotificationConsentHandler struct {
	notification.UnimplementedNotificationServiceServer
	notificationConsentService service.NotificationConsentService
	subscriber                 messaging.Subscriber
	publisher                  messaging.Publisher
}

func NewNotificationConsentHandler(notificationConsentService service.NotificationConsentService, subscriber messaging.Subscriber, publisher messaging.Publisher) *NotificationConsentHandler {

	handler := &NotificationConsentHandler{
		notificationConsentService: notificationConsentService,
		publisher:                  publisher,
		subscriber:                 subscriber,
	}
	handler.subscriber.Subscribe(handler.HandleMessages)

	return handler
}

func (handler *NotificationConsentHandler) HandleMessages(message *NotificationMessaging.NotificationMessage) {
	notificationConsent, err := handler.notificationConsentService.GetById(message.AccountID)
	if err != nil {
		return
	}
	switch message.MessageType {
	case "RequestMade":
		if notificationConsent.HostResponded {
			handler.publisher.Publish(message)
		}
	case "ReservationCanceled":
		if notificationConsent.ReservationCanceled {
			handler.publisher.Publish(message)
		}
	case "HostRatingGiven":
		if notificationConsent.HostRatingGiven {
			handler.publisher.Publish(message)
		}
	case "AccommodationRatingGiven":
		if notificationConsent.AccommodationRatingGiven {
			handler.publisher.Publish(message)
		}
	case "ProminentHost":
		if notificationConsent.ProminentHost {
			handler.publisher.Publish(message)
		}
	case "HostResponded":
		if notificationConsent.HostResponded {
			handler.publisher.Publish(message)
		}

	}

}

func (handler *NotificationConsentHandler) Create(ctx context.Context, req *notification.CreateRequest) (*notification.CreateResponse, error) {
	notificationConsent, err := mapFromCreateRequest(req)
	if err != nil {
		return nil, err
	}

	return &notification.CreateResponse{}, handler.notificationConsentService.Create(notificationConsent)
}

func (handler *NotificationConsentHandler) UpdateMyNotificationConsent(ctx context.Context, req *notification.UpdateMyNotificationConsentRequest) (*notification.UpdateMyNotificationConsentResponse, error) {
	loggedInId, err := utils.GetTokenInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract id")
	}

	UpdateMyNotificationConsent := mapFromUpdateMyNotificationConsentRequest(req, loggedInId)
	_, err = handler.notificationConsentService.UpdateMyNotificationConsent(UpdateMyNotificationConsent)

	if err != nil {
		return nil, err
	}
	return &notification.UpdateMyNotificationConsentResponse{Message: "Updated"}, nil

}

func (handler *NotificationConsentHandler) GetMyNotificationSettings(ctx context.Context, req *notification.EmptyRequest) (*notification.GetMyNotificationSettingsResponse, error) {
	loggedInId, err := utils.GetTokenInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract id")
	}

	notificationConsent, err := handler.notificationConsentService.GetById(loggedInId)
	if err != nil {
		return nil, err
	}

	return mapToRequest(notificationConsent), nil
}
