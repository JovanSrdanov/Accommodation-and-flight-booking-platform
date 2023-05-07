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

func (service UserProfileService) Update(id uuid.UUID, dto *model.UpdateProfileDto) (*model.UpdateProfileDto, error) {
	userInfo, err := service.GetById(id)
	if err != nil {
		return &model.UpdateProfileDto{}, err
	}

	userInfo.Name = dto.Name
	userInfo.Surname = dto.Surname
	userInfo.Email = dto.Email
	userInfo.Address = dto.Address

	userInfo, err = service.userProfRepo.Update(userInfo)
	if err != nil {
		return nil, err
	}

	return &model.UpdateProfileDto{
		Name:    userInfo.Name,
		Surname: userInfo.Surname,
		Email:   userInfo.Email,
		Address: userInfo.Address,
	}, nil
}
