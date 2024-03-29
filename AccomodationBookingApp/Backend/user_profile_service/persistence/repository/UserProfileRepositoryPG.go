package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"user_profile_service/domain/model"
)

type UserProfileRepositoryPG struct {
	dbClient *gorm.DB
}

func NewUserProfileRepositoryPG(dbClient *gorm.DB) (*UserProfileRepositoryPG, error) {
	//Syncs database schema with current application model
	//TODO izmestiti u poseban address repo?
	err := dbClient.AutoMigrate(&model.Address{}, &model.UserProfile{})
	if err != nil {
		return nil, err
	}
	return &UserProfileRepositoryPG{dbClient: dbClient}, err
}

func (repo UserProfileRepositoryPG) Create(userProf *model.UserProfile) (uuid.UUID, error) {
	if userProf.ID == uuid.Nil {
		userProf.ID, _ = uuid.NewUUID()
	}

	if userProf.Address.ID == uuid.Nil {
		userProf.Address.ID, _ = uuid.NewUUID()
	}

	result := repo.dbClient.Create(userProf)
	if result.Error != nil {
		return uuid.UUID{}, result.Error
	}

	return userProf.ID, nil
}

func (repo UserProfileRepositoryPG) GetById(id uuid.UUID) (*model.UserProfile, error) {
	var userProf model.UserProfile

	//With eager loading
	result := repo.dbClient.Where("id = ?", id).Preload("Address").First(&userProf)
	if result.Error != nil {
		return &model.UserProfile{}, result.Error
	}

	return &userProf, nil
}

func (repo UserProfileRepositoryPG) Update(userProf *model.UserProfile) (*model.UserProfile, error) {
	if err := repo.dbClient.Save(&userProf).Error; err != nil {
		return &model.UserProfile{}, err
	}
	repo.dbClient.Save(&userProf.Address)

	return userProf, nil
}
func (repo UserProfileRepositoryPG) Delete(id uuid.UUID) error {

	transaction := repo.dbClient.Begin()
	//For panic recovery
	defer func() {
		r := recover()
		if r != nil {
			transaction.Rollback()
		}
	}()

	var userProfile model.UserProfile
	//Because of cascade deletion
	result := transaction.First(&userProfile, id)
	if result.Error != nil {
		transaction.Rollback()
		return result.Error
	}

	// DeleteUserProfile the user profile and its associated address
	result = transaction.Delete(&userProfile)
	if result.Error != nil {
		transaction.Rollback()
		return result.Error
	}

	result = transaction.Delete(&model.Address{ID: userProfile.AddressID})
	if result.Error != nil {
		transaction.Rollback()
		return result.Error
	}

	err := transaction.Commit().Error
	if err != nil {
		transaction.Rollback()
		return err
	}

	return nil
}
