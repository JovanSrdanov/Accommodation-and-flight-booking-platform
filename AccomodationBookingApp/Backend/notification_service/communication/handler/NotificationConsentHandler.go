package handler

import (
	notification "common/proto/notification_service/generated"
	"context"
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
