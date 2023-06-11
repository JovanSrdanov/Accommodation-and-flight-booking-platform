package service

import (
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
