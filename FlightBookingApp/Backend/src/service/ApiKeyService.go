package service

import (
	"FlightBookingApp/model"
	"FlightBookingApp/repository"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApiKeyService struct {
	apiKeyRepo repository.ApiKeyRepository
}

func NewApiKeyService(apiKeyRepo repository.ApiKeyRepository) *ApiKeyService {
	return &ApiKeyService{
		apiKeyRepo: apiKeyRepo,
	}
}
func (service *ApiKeyService) Create(key *model.ApiKey) error {
	return service.apiKeyRepo.Create(key)
}

func (service *ApiKeyService) GetByAccountId(id primitive.ObjectID) (model.ApiKey, error) {
	key, err := service.apiKeyRepo.GetByAccountId(id)
	if !key.IsValid() {
		return model.ApiKey{}, errors.New("Api key expired")
	}
	return key, err
}
