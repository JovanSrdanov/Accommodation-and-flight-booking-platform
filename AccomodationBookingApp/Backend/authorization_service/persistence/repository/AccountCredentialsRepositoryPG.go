package repository

import (
	"authorization_service/domain/model"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Postgres repository
type AccountCredentialsRepositoryPG struct {
	dbClient *gorm.DB
}

func NewAccountCredentialsRepositoryPG(dbClient *gorm.DB) (*AccountCredentialsRepositoryPG, error) {
	//Syncs database schema with current application model
	err := dbClient.AutoMigrate(&model.AccountCredentials{})
	if err != nil {
		return nil, err
	}
	return &AccountCredentialsRepositoryPG{dbClient: dbClient}, nil
}

func (repo AccountCredentialsRepositoryPG) Create(accCred *model.AccountCredentials) (uuid.UUID, error) {
	_, err := repo.GetByUsername(accCred.Username)
	if err == nil {
		return uuid.Nil, fmt.Errorf("username taken")
	}

	accCred.ID, _ = uuid.NewUUID()

	result := repo.dbClient.Create(accCred)
	if result.Error != nil {
		return uuid.UUID{}, result.Error
	}

	return accCred.ID, nil
}

func (repo AccountCredentialsRepositoryPG) GetByUsername(username string) (*model.AccountCredentials, error) {
	var accCred model.AccountCredentials

	result := repo.dbClient.Where("username = ?", username).First(&accCred)
	if result.Error != nil {
		return &model.AccountCredentials{}, result.Error
	}

	return &accCred, nil
}

func (repo AccountCredentialsRepositoryPG) GetById(id uuid.UUID) (*model.AccountCredentials, error) {
	var accCred model.AccountCredentials

	result := repo.dbClient.Where("id = ?", id).First(&accCred)
	if result.Error != nil {
		return &model.AccountCredentials{}, result.Error
	}

	return &accCred, nil
}

func (repo AccountCredentialsRepositoryPG) Update(accCred *model.AccountCredentials) error {
	if err := repo.dbClient.Save(&accCred).Error; err != nil {
		return err
	}

	return nil
}
func (repo AccountCredentialsRepositoryPG) Delete(id uuid.UUID) error {

	result := repo.dbClient.Where("user_profile_id = ?", id).Delete(&model.AccountCredentials{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("account credentials with given user profile id not found")
	}

	return nil
}
