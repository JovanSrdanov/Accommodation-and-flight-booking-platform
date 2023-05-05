package service

import (
	"authorization_service/domain/model"
	"authorization_service/domain/repository"
	"authorization_service/domain/token"
	"github.com/google/uuid"
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
	//TODO ukloniti ovu inicijalizaciju kad se bude pisala prava logika
	accCred.Salt = "salt"
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

func (service AccountCredentialsService) Login(username, password string) (string, error) {
	accountCredentials, err := service.GetByUsername(username)
	if err != nil {
		return "", status.Errorf(codes.NotFound, "incorrect username")
	}

	if !accountCredentials.IsPasswordCorrect(password) {
		return "", status.Errorf(codes.NotFound, "incorrect password")
	}

	accessToken, err := service.tokenMaker.CreateToken(
		accountCredentials.Username,
		15*time.Minute,
		accountCredentials.Role,
	)
	if err != nil {
		return "", status.Errorf(codes.Internal, "Cannot generate access token")
	}

	return accessToken, nil
}
