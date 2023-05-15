package service

import (
	authorization "common/proto/authorization_service/generated"
	"github.com/google/uuid"
	"log"
	"user_profile_service/communication/orchestrator"
	"user_profile_service/domain/model"
	"user_profile_service/domain/repository"
)

type UserProfileService struct {
	userProfRepo       repository.IUserProfileRepository
	deleteOrchestrator *orchestrator.DeleteUserOrchestrator
}

func NewUserProfileService(userProfRepo repository.IUserProfileRepository, deleteOrchestrator *orchestrator.DeleteUserOrchestrator) *UserProfileService {
	return &UserProfileService{
		userProfRepo:       userProfRepo,
		deleteOrchestrator: deleteOrchestrator}
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
	log.Println("new address: ", dto.Address)
	userInfo.Address.Street = dto.Address.Street
	userInfo.Address.City = dto.Address.City
	userInfo.Address.Country = dto.Address.Country
	userInfo.Address.StreetNumber = dto.Address.StreetNumber

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
func (service UserProfileService) DeleteUser(accCredId string, userProfileID uuid.UUID, role authorization.Role) error {
	return service.deleteOrchestrator.Start(accCredId, userProfileID, role)
}

func (service UserProfileService) Delete(id uuid.UUID) error {
	return service.userProfRepo.Delete(id)
}
