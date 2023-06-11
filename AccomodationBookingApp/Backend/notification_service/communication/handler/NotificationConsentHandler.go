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
	return &notification.CreateResponse{}, handler.notificationConsentService.Create(mapFromCreateRequest(req))
}
