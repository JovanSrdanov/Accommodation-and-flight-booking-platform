package repository

import (
	"github.com/google/uuid"
	"notification_service/domain/model"
)

type INotificationConsentRepository interface {
	Create(accCred *model.NotificationConsent) error
	Update(accCred *model.NotificationConsent) error
	GetById(id uuid.UUID) (*model.NotificationConsent, error)
	Delete(id uuid.UUID) error
}
