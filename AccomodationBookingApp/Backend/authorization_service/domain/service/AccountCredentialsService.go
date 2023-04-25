package service

import (
	"authorization_service/domain/model"
	"authorization_service/domain/repository"
	"github.com/google/uuid"
)

type AccountCredentialsService struct {
	accCredRepo repository.IAccountCredentialsRepository
}

func NewAccountCredentialsService(accCredRepo repository.IAccountCredentialsRepository) *AccountCredentialsService {
	return &AccountCredentialsService{accCredRepo: accCredRepo}
}

func (service AccountCredentialsService) Create(accCred model.AccountCredentials) (uuid.UUID, error) {
	id, err := service.accCredRepo.Create(accCred)

	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (service AccountCredentialsService) GetByEmail(email string) (*model.AccountCredentials, error) {
	accountCredentials, err := service.accCredRepo.GetByEmail(email)

	if err != nil {
		return &model.AccountCredentials{}, err
	}
	return accountCredentials, nil
}
