package repository

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"notification_service/domain/model"
)

type NotificationConsentRepositoryPG struct {
	dbClient *gorm.DB
}

func (repo NotificationConsentRepositoryPG) Delete(id uuid.UUID) error {
	result := repo.dbClient.Where("user_profile_id = ?", id).Delete(&model.NotificationConsent{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notification consent with given user profile id not found")
	}

	return nil
}

func NewNotificationConsentServicePG(dbClient *gorm.DB) (*NotificationConsentRepositoryPG, error) {
	err := dbClient.AutoMigrate(&model.NotificationConsent{})
	if err != nil {
		return nil, err
	}

	return &NotificationConsentRepositoryPG{dbClient: dbClient}, nil
}

func (repo NotificationConsentRepositoryPG) Create(notificationConsent *model.NotificationConsent) error {
	result := repo.dbClient.Create(notificationConsent)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (repo NotificationConsentRepositoryPG) GetById(id uuid.UUID) (*model.NotificationConsent, error) {
	notificationConsent := &model.NotificationConsent{}
	result := repo.dbClient.First(notificationConsent, "user_profile_id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return notificationConsent, nil
}

func (repo NotificationConsentRepositoryPG) Update(notificationConsent *model.NotificationConsent) (*model.NotificationConsent, error) {
	repo.dbClient.Save(&notificationConsent)
	return nil, nil
}
