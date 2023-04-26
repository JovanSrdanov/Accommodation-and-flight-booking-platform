package repository

import (
	"authorization_service/domain/model"
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
	accCred.ID, _ = uuid.NewUUID()

	result := repo.dbClient.Create(accCred)
	if result.Error != nil {
		return uuid.UUID{}, result.Error
	}

	return accCred.ID, nil
}

func (repo AccountCredentialsRepositoryPG) GetByEmail(email string) (*model.AccountCredentials, error) {
	var accCred model.AccountCredentials

	result := repo.dbClient.Where("email = ?", email).First(&accCred)
	if result.Error != nil {
		return &model.AccountCredentials{}, result.Error
	}

	return &accCred, nil
}
