package service

import (
	"authorization_service/domain/model"
	"github.com/google/uuid"
)

type IAccountCredentialsService interface {
	Create(accCred *model.AccountCredentials) (uuid.UUID, error)
	GetByUsername(username string) (*model.AccountCredentials, error)
	GetById(id uuid.UUID) (*model.AccountCredentials, error)
	Login(username, password string) (string, error)
	Update(id uuid.UUID, newUsername, newPassword string) error
}
