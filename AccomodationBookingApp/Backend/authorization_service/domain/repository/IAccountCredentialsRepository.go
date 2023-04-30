package repository

import (
	"authorization_service/domain/model"
	"github.com/google/uuid"
)

type IAccountCredentialsRepository interface {
	Create(accCred *model.AccountCredentials) (uuid.UUID, error)
	GetByEmail(email string) (*model.AccountCredentials, error)
}