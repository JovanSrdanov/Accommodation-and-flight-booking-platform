package service

import (
	"github.com/google/uuid"
	"user_profile_service/domain/model"
	"user_profile_service/domain/repository"
)

type UserProfileService struct {
	userProfRepo repository.IUserProfileRepository
}

func NewUserProfileService(userProfRepo repository.IUserProfileRepository) *UserProfileService {
	return &UserProfileService{userProfRepo: userProfRepo}
}

func (service UserProfileService) Create(userProf *model.UserProfile) (uuid.UUID, error) {
	return service.userProfRepo.Create(userProf)
}

func (service UserProfileService) GetById(id uuid.UUID) (*model.UserProfile, error) {
	return service.userProfRepo.GetById(id)
}
