package service

import (
	"github.com/google/uuid"
	"user_profile_service/domain/model"
)

type IUserProfileService interface {
	Create(userProf *model.UserProfile) (uuid.UUID, error)
	GetById(id uuid.UUID) (*model.UserProfile, error)
	Update(id uuid.UUID, dto *model.UpdateProfileDto) (*model.UpdateProfileDto, error)
	Delete(id uuid.UUID) error
}
