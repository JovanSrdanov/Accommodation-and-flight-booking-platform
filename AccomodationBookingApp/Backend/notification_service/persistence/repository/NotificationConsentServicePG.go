package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"notification_service/domain/model"
)

type NotificationConsentServicePG struct {
	dbClient *gorm.DB
}

func NewNotificationConsentServicePG(dbClient *gorm.DB) (*NotificationConsentServicePG, error) {
	err := dbClient.AutoMigrate(&model.NotificationConsent{})
	if err != nil {
		return nil, err
	}

	return &NotificationConsentServicePG{dbClient: dbClient}, nil
}

func (repo NotificationConsentServicePG) Create(notificationConsent *model.NotificationConsent) error {
	result := repo.dbClient.Create(notificationConsent)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo NotificationConsentServicePG) Update(accCred *model.NotificationConsent) error {
	//TODO implement me
	panic("implement me")
}

func (repo NotificationConsentServicePG) Delete(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
