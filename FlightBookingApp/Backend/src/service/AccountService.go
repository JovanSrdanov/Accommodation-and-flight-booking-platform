package service

import (
	utils "FlightBookingApp/Utils"
	"FlightBookingApp/dto"
	"FlightBookingApp/model"
	"FlightBookingApp/repository"
	"FlightBookingApp/token"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type accountService struct {
	accountRepository repository.AccountRepository
}

type AccountService interface {
	Register(account model.Account) (primitive.ObjectID, error)
	Login(loginData dto.LoginRequest) (string, string, error)
	GetAll() (model.Accounts, error)
	GetById(id primitive.ObjectID) (model.Account, error)
	GetByUsername(username string) (model.Account, error)
	Save(model.Account) (model.Account, error)
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

	if !accountToBeLoggedIn.IsActivated {
		return "", "", fmt.Errorf("account not activated")
	}

	return token.GenerateTokens(accountToBeLoggedIn)
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

	// TODO Stefan: move to separate func
	// email activation logic 
	// return time to the nanosecond (1 billionth of a sec)
	rand.Seed(time.Now().UnixNano())
	// create random code for email
	// Go rune data type represent Unicode characters
	var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	emailVerRandRune := make([]rune, 64)
	// create a random slice of runes (characters) to create our emailVerPassword (random string of characters)
	for i := 0; i < 64; i++ {
		emailVerRandRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes) - 1)]
	}

	emailVerPassword := string(emailVerRandRune)
	var emailVerPWhash []byte
	// func GenerateFromPassword(password []byte, const int) ([]byte, error)
	emailVerPWhash, err = bcrypt.GenerateFromPassword([]byte(emailVerPassword), bcrypt.DefaultCost)
	if err != nil {
		return primitive.NewObjectID(), err 
	}
	newAccount.EmailVerificationHash = string(emailVerPWhash)

	// create u.timeout after 48 hours
	timeout := time.Now().Local().AddDate(0, 0, 2)
	newAccount.VerificationTimeout = timeout

	// TODO Stefan: add a transaction, in case the email sending fails abort database modification
	utils.SendConfirmationEmail(newAccount, emailVerPassword)

	return service.accountRepository.Create(&newAccount)
}

func (service *accountService) GetAll() (model.Accounts, error) {
	return service.accountRepository.GetAll()
}

func (service *accountService) GetById(id primitive.ObjectID) (model.Account, error) {
	return service.accountRepository.GetById(id)
}

func (service *accountService) GetByUsername(username string) (model.Account, error) {
	return service.accountRepository.GetByUsername(username)
}

func (service *accountService) Save(account model.Account) (model.Account, error) {
	return service.accountRepository.Save(account)
}

func (service *accountService) Delete(id primitive.ObjectID) error {
	return service.accountRepository.Delete(id)
}