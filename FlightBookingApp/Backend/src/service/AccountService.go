package service

import (
	JWT "FlightBookingApp/JWT"
	utils "FlightBookingApp/Utils"
	"FlightBookingApp/dto"
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"FlightBookingApp/repository"
	"FlightBookingApp/token"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type accountService struct {
	accountRepository repository.AccountRepository
}

type AccountService interface {
	Register(account model.Account) (primitive.ObjectID, error)
	Login(loginData dto.LoginRequest) (string, error)
	GetAll() (model.Accounts, error)
	GetById(id primitive.ObjectID) (model.Account, error)
	Delete(id primitive.ObjectID) error
}

func NewAccountService(accountRepository repository.AccountRepository) *accountService {
	return &accountService {
		accountRepository:  accountRepository,
	}
}

func (service *accountService) Login(loginData dto.LoginRequest) (string, error) {
	allAccounts, _ := service.accountRepository.GetAll()

	val, accountToBeLoggedIn := isLoginDataValid(loginData, allAccounts)
	if !val {
		return "", fmt.Errorf("username of password invalid")
	}

	// after verifying login credentials, generates claims for the jwt session token and generates it
	var claims = &JWT.JwtClaims{}
	claims.ID = accountToBeLoggedIn.ID
	claims.Username = loginData.Username
	claims.Roles = []model.Role{accountToBeLoggedIn.Role}

	var tokenCreationTime = time.Now().UTC()

	// session token lasts for 30 minutes
	var expirationTime = tokenCreationTime.Add(time.Duration(30) * time.Minute)
	return token.GenerateToken(claims, expirationTime)
}

// if login data is valid returns both true and the id of the logged in account
func isLoginDataValid(loginData dto.LoginRequest, accounts model.Accounts) (bool, model.Account){
	val, account := usernameExists(loginData.Username, accounts)
	if val && isPasswordValid(loginData.Password, accounts) {
		return true, account
	}
	return false, model.Account{}
}

func isPasswordValid(password string, accounts model.Accounts) bool{
	for _, value := range accounts{
		if utils.CheckPassword(password, value.Password) == nil{
			return true
		}
	}
	return false
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