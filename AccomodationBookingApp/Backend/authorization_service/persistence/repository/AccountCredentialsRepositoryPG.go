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
	//TODO : probaj da obrises da vidis oce li ti sam postgres izgenerisati ID
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
