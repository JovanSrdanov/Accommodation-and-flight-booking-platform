package service

import (
	utils "FlightBookingApp/Utils"
	"FlightBookingApp/dto"
	"FlightBookingApp/errors"
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
	allAccounts, _ := service.accountRepository.GetAll()

	if isAccountValid(newAccount, allAccounts) {
		return service.accountRepository.Create(&newAccount)
	}
	return primitive.NilObjectID, &errors.UsernameOrEmailExistsError{}
}

func  isAccountValid(account model.Account, accounts model.Accounts) bool {
	val, _ := usernameExists(account.Username, accounts)
	if !val && !isEmailTaken(account.Email, accounts) {
		return true
	}
	return false
}

func usernameExists(username string, accounts model.Accounts) (bool, model.Account) {
	for _, value := range accounts{
		if username == value.Username {
			return true, *value
		}
	}
	return false, model.Account{}
}

func isEmailTaken(email string, accounts model.Accounts) bool {
	for _, value := range accounts {
		if email == value.Email {
			return true
		}
	}
	return false
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