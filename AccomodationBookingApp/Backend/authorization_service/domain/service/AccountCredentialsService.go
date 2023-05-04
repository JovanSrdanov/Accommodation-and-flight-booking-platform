package service

import (
	"authorization_service/domain/model"
	"authorization_service/domain/repository"
	"authorization_service/token"
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

func (service AccountCredentialsService) Login(accCred *model.AccountCredentials) (string, error) {
	accountCredentials, err := service.GetByUsername(accCred.Username)
	if err != nil || !accountCredentials.IsPasswordCorrect(accCred.Password) {
		return "", status.Errorf(codes.NotFound, "incorrect username/password")
	}

	accessToken, err := service.tokenMaker.CreateToken(
		accCred.Username,
		time.Duration(15*time.Minute),
		accCred.Role,
	)
	if err != nil {
		return "", status.Errorf(codes.Internal, "Cannot generate access token")
	}

	return accessToken, nil
}
