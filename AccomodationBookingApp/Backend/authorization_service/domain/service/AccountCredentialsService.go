package service

import (
	"authorization_service/domain/model"
	"authorization_service/domain/repository"
	"authorization_service/domain/token"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type AccountCredentialsService struct {
	accCredRepo repository.IAccountCredentialsRepository
	tokenMaker  token.Maker
}

func NewAccountCredentialsService(accCredRepo repository.IAccountCredentialsRepository,
	tokenMaker token.Maker) *AccountCredentialsService {
	return &AccountCredentialsService{
		accCredRepo: accCredRepo,
		tokenMaker:  tokenMaker,
	}
}

func (service AccountCredentialsService) Create(accCred *model.AccountCredentials) (uuid.UUID, error) {
	id, err := service.accCredRepo.Create(accCred)

	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (service AccountCredentialsService) GetByUsername(username string) (*model.AccountCredentials, error) {
	accountCredentials, err := service.accCredRepo.GetByUsername(username)
	if err != nil {
		return &model.AccountCredentials{}, err
	}

	return accountCredentials, nil
}

func (service AccountCredentialsService) GetById(id uuid.UUID) (*model.AccountCredentials, error) {
	accountCredentials, err := service.accCredRepo.GetById(id)
	if err != nil {
		return &model.AccountCredentials{}, err
	}

	return accountCredentials, nil
}

func (service AccountCredentialsService) Login(username, password string) (string, error) {
	accountCredentials, err := service.GetByUsername(username)
	if err != nil {
		return "", status.Errorf(codes.NotFound, "incorrect username")
	}

	if !accountCredentials.IsPasswordCorrect(password) {
		return "", status.Errorf(codes.NotFound, "incorrect password")
	}

	accessToken, _, err := service.tokenMaker.CreateToken(
		accountCredentials.ID,
		180*time.Minute,
		accountCredentials.Role,
	)
	if err != nil {
		return "", status.Errorf(codes.Internal, "Cannot generate access token")
	}

	return accessToken, nil
}

func (service AccountCredentialsService) ChangeUsername(userId uuid.UUID, username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	oldAccCred, err := service.GetById(userId)
	if err != nil {
		return fmt.Errorf("error while getting logged-in user info")
	}

	if oldAccCred.Username == username {
		return nil
	}

	_, err = service.GetByUsername(username)
	if err == nil {
		return fmt.Errorf("username already exists")
	}

	oldAccCred.Username = username
	err = service.accCredRepo.Update(oldAccCred)
	if err != nil {
		return err
	}

	return nil
}

func (service AccountCredentialsService) ChangePassword(id uuid.UUID, oldPassword, newPassword string) error {
	oldAccCred, err := service.GetById(id)
	if err != nil {
		return err
	}

	if !oldAccCred.IsPasswordCorrect(oldPassword) {
		return status.Errorf(codes.Unauthenticated, "provided old password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	oldAccCred.Password = string(hashedPassword)

	err = service.accCredRepo.Update(oldAccCred)
	if err != nil {
		return err
	}

	return nil
}

func (service AccountCredentialsService) Delete(id uuid.UUID) error {
	return service.accCredRepo.Delete(id)
}
