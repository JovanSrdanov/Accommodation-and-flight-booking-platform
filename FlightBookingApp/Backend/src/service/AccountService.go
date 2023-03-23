package service

import (
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"FlightBookingApp/repository"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type accountService struct {
	accountRepository repository.AccountRepository
}

type AccountService interface {
	Register(account model.Account) (model.Account, error)
	GetAll() []model.Account
	GetById(id uuid.UUID) (model.Account, error)
	Delete(id primitive.ObjectID) error
}

func NewAccountService(accountRepository repository.AccountRepository) AccountService {
	return &accountService {
		accountRepository:  accountRepository,
	}
}

func (service *accountService) Register(account model.Account) (model.Account, error) {
	if isAccountValid(account, service) {
		service.accountRepository.Create(account)
		return account, nil
	}

	return account, &errors.UsernameOrEmailExistsError{}
}

func  isAccountValid(account model.Account, service *accountService) bool {
	return isUsernameTaken(account.Username, service) && isEmailTaken(account.Email, service)
}

func isUsernameTaken(username string, service *accountService) bool {
	for _, value := range service.accountRepository.GetAll() {
		if username == value.Username {
			return false
		}
	}
	return true
}

func isEmailTaken(email string, service *accountService) bool {
	for _, value := range service.accountRepository.GetAll() {
		if email == value.Email {
			return false
		}
	}
	return true
}

func (service *accountService) GetAll() []model.Account {
	return service.accountRepository.GetAll()
}

func (service *accountService) GetById(id uuid.UUID) (model.Account, error) {
	return service.accountRepository.GetById(id)
}

func (service *accountService) Delete(id primitive.ObjectID) error {
	return service.accountRepository.Delete(id)
}