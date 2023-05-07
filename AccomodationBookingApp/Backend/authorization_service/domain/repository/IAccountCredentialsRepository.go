package repository

import (
	"authorization_service/domain/model"
	"github.com/google/uuid"
)

type IAccountCredentialsRepository interface {
	Create(accCred *model.AccountCredentials) (uuid.UUID, error)
	GetByUsername(username string) (*model.AccountCredentials, error)
	GetById(id uuid.UUID) (*model.AccountCredentials, error)
	Update(accCred *model.AccountCredentials) error
}
