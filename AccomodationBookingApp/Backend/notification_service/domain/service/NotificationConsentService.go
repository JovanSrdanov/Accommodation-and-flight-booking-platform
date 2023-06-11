package service

import (
	"github.com/google/uuid"
	"notification_service/domain/model"
	"notification_service/domain/repository"
)

type NotificationConsentService struct {
	notificationConsentRepository repository.INotificationConsentRepository
}

func NewNotificationConsentService(notificationConsentRepository repository.INotificationConsentRepository) *NotificationConsentService {
	return &NotificationConsentService{notificationConsentRepository: notificationConsentRepository}
}

func (service NotificationConsentService) Create(notificationConsent *model.NotificationConsent) error {
	err := service.notificationConsentRepository.Create(notificationConsent)
	if err != nil {
		return err
	}
	return nil

}

func (service NotificationConsentService) GetById(id uuid.UUID) (*model.NotificationConsent, error) {
	return service.notificationConsentRepository.GetById(id)
}
