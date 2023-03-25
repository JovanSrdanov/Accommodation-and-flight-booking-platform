package service

import (
	utils "FlightBookingApp/Utils"
	"FlightBookingApp/dto"
	"FlightBookingApp/model"
	"FlightBookingApp/repository"
	"FlightBookingApp/token"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type accountService struct {
	accountRepository repository.AccountRepository
}

type AccountService interface {
	Register(account model.Account) (primitive.ObjectID, error)
	Login(loginData dto.LoginRequest) (string, string, error)
	GetAll() (model.Accounts, error)
	GetById(id primitive.ObjectID) (model.Account, error)
	Delete(id primitive.ObjectID) error
}

func NewAccountService(accountRepository repository.AccountRepository) *accountService {
	return &accountService {
		accountRepository:  accountRepository,
	}
}

func (service *accountService) Login(loginData dto.LoginRequest) (string, string, error) {
	accountToBeLoggedIn, err := service.accountRepository.GetByUsername(loginData.Username)
	if err != nil {
		return "", "", fmt.Errorf("username of password invalid")
	}

	err = utils.CheckPassword(loginData.Password, accountToBeLoggedIn.Password)
	if err != nil {
		return "", "", fmt.Errorf("username of password invalid")
	}

	return token.GenerateToken(accountToBeLoggedIn)
}

func (service *accountService) Register(newAccount model.Account) (primitive.ObjectID, error) {
	_, err := service.accountRepository.GetByUsername(newAccount.Username)
	if err == nil {
		return primitive.NewObjectID(), fmt.Errorf("username already exists") 
	}

	_, err = service.accountRepository.GetByEmail(newAccount.Email)
	if err == nil {
		return primitive.NewObjectID(), fmt.Errorf("email already exists") 
	}

	return service.accountRepository.Create(&newAccount)
}

func (service *accountService) GetAll() (model.Accounts, error) {
	return service.accountRepository.GetAll()
}

func (service *accountService) GetById(id primitive.ObjectID) (model.Account, error) {
	return service.accountRepository.GetById(id)
}

func (service *accountService) Delete(id primitive.ObjectID) error {
	return service.accountRepository.Delete(id)
}