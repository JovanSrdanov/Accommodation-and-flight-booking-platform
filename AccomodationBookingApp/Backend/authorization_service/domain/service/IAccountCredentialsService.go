package service

import (
	"authorization_service/domain/model"
	"github.com/google/uuid"
)

type IAccountCredentialsService interface {
	Create(accCred *model.AccountCredentials) (uuid.UUID, error)
	GetByUsername(username string) (*model.AccountCredentials, error)
	Login(username, password string) (string, error)
}
