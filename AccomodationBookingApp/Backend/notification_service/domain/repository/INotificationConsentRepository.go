package repository

import (
	"github.com/google/uuid"
	"notification_service/domain/model"
)

type INotificationConsentRepository interface {
	Create(notificationConsent *model.NotificationConsent) error
	Update(notificationConsent *model.NotificationConsent) (*model.NotificationConsent, error)
	GetById(id uuid.UUID) (*model.NotificationConsent, error)
	Delete(id uuid.UUID) error
}
