package repository

import (
	"github.com/google/uuid"
	"user_profile_service/domain/model"
)

type IUserProfileRepository interface {
	Create(userProf *model.UserProfile) (uuid.UUID, error)
	GetById(id uuid.UUID) (*model.UserProfile, error)
	Update(userProf *model.UserProfile) (*model.UserProfile, error)
	Delete(id uuid.UUID) error
}
