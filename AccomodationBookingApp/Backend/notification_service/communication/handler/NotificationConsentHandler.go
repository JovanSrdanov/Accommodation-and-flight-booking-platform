package handler

import (
	notification "common/proto/notification_service/generated"
	"context"
	"fmt"
	"github.com/google/uuid"
	"notification_service/domain/service"
)

type NotificationConsentHandler struct {
	notification.UnimplementedNotificationServiceServer
	notificationConsentService service.NotificationConsentService
}

func NewNotificationConsentHandler(notificationConsentService service.NotificationConsentService) *NotificationConsentHandler {
	return &NotificationConsentHandler{notificationConsentService: notificationConsentService}
}

func (handler *NotificationConsentHandler) Create(ctx context.Context, req *notification.CreateRequest) (*notification.CreateResponse, error) {
	notificationConsent, err := mapFromCreateRequest(req)
	if err != nil {
		return nil, err
	}

	return &notification.CreateResponse{}, handler.notificationConsentService.Create(notificationConsent)
}

func (handler *NotificationConsentHandler) UpdateMyNotificationConsent(ctx context.Context, req *notification.UpdateMyNotificationConsentRequest) (*notification.UpdateMyNotificationConsentResponse, error) {

	loggedInId, ok := ctx.Value("id").(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("failed to extract id and cast to UUID")
	}
	UpdateMyNotificationConsent := mapFromUpdateMyNotificationConsentRequest(req, loggedInId)
	_, err := handler.notificationConsentService.UpdateMyNotificationConsent(UpdateMyNotificationConsent)

	if err != nil {
		return nil, err
	}
	return &notification.UpdateMyNotificationConsentResponse{Message: "Updated"}, nil

}

func (handler *NotificationConsentHandler) GetMyNotificationSettings(ctx context.Context, req *notification.EmptyRequest) (*notification.GetMyNotificationSettingsResponse, error) {
	loggedInId, ok := ctx.Value("id").(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("failed to extract id and cast to UUID")
	}

	notificationConsent, err := handler.notificationConsentService.GetById(loggedInId)
	if err != nil {
		return nil, err
	}

	return mapToRequest(notificationConsent), nil
}
