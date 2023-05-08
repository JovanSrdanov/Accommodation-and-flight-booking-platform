package service

import (
	"authorization_service/domain/model"
	"authorization_service/domain/repository"
	"authorization_service/domain/token"
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

func (service AccountCredentialsService) Login(username, password string) (string, model.Role, error) {
	accountCredentials, err := service.GetByUsername(username)
	if err != nil {
		return "", -1, status.Errorf(codes.NotFound, "incorrect username")
	}

	if !accountCredentials.IsPasswordCorrect(password) {
		return "", -1, status.Errorf(codes.NotFound, "incorrect password")
	}

	accessToken, role, err := service.tokenMaker.CreateToken(
		accountCredentials.ID,
		15*time.Minute,
		accountCredentials.Role,
	)
	if err != nil {
		return "", -1, status.Errorf(codes.Internal, "Cannot generate access token")
	}

	return accessToken, role, nil
}

func (service AccountCredentialsService) Update(id uuid.UUID, newUsername, newPassword string) error {
	accCred, err := service.GetById(id)
	if err != nil {
		return err
	}

	accCred.Username = newUsername

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	accCred.Password = string(hashedPassword)

	err = service.accCredRepo.Update(accCred)
	if err != nil {
		return err
	}

	return nil
}
