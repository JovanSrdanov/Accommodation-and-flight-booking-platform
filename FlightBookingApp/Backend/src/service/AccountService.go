package service

import (
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"FlightBookingApp/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type accountService struct {
	accountRepository repository.AccountRepository
}

type AccountService interface {
	Register(account model.Account) (primitive.ObjectID, error)
	GetAll() (model.Accounts, error)
	GetById(id primitive.ObjectID) (model.Account, error)
	Delete(id primitive.ObjectID) error
}

func NewAccountService(accountRepository repository.AccountRepository) *accountService {
	return &accountService {
		accountRepository:  accountRepository,
	}
}

func (service *accountService) Register(newAccount model.Account) (primitive.ObjectID, error) {
	allAccounts, _ := service.accountRepository.GetAll()

	if isAccountValid(newAccount, allAccounts) {
		return service.accountRepository.Create(&newAccount)
	}
	return primitive.NilObjectID, &errors.UsernameOrEmailExistsError{}
}

func  isAccountValid(account model.Account, accounts model.Accounts) bool {
	return isUsernameTaken(account.Username, accounts) && isEmailTaken(account.Email, accounts)
}

func isUsernameTaken(username string, accounts model.Accounts) bool {
	for _, value := range accounts{
		if username == value.Username {
			return false
		}
	}
	return true
}

func isEmailTaken(email string, accounts model.Accounts) bool {
	for _, value := range accounts {
		if email == value.Email {
			return false
		}
	}
	return true
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