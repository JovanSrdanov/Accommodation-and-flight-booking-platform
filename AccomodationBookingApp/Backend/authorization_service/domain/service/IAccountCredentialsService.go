package service

import (
	"authorization_service/domain/model"
	"github.com/google/uuid"
	"time"
)

type IAccountCredentialsService interface {
	Create(accCred *model.AccountCredentials) (uuid.UUID, error)
	GetByUsername(username string) (*model.AccountCredentials, error)
	GetById(id uuid.UUID) (*model.AccountCredentials, error)
	Login(username, password string) (string, model.Role, time.Time, error)
	ChangePassword(id uuid.UUID, oldPassword, newPassword string) error
	ChangeUsername(id uuid.UUID, username string) error
}
