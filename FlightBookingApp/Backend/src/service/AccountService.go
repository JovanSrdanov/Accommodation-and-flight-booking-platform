package service

import (
	"FlightBookingApp/dto"
	"FlightBookingApp/model"
	"FlightBookingApp/repository"
	"FlightBookingApp/token"
	"FlightBookingApp/utils"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type accountService struct {
	accountRepository repository.AccountRepository
	userRepository    repository.UserRepository
}

type AccountService interface {
	Register(registrationInfo dto.AccountRegistration) (dto.CreateUserResponse, error)
	Login(loginData dto.LoginRequest) (string, string, error)
	GetAll() (model.Accounts, error)
	GetById(id primitive.ObjectID) (model.Account, error)
	GetByUsername(username string) (model.Account, error)
	GetByRefreshToken(token string) (model.Account, error)
	Save(model.Account) (model.Account, error)
	Delete(id primitive.ObjectID) error
}

func NewAccountService(accountRepository repository.AccountRepository, userRepository repository.UserRepository) *accountService {
	return &accountService{
		accountRepository: accountRepository,
		userRepository:    userRepository,
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

	accessTokenString, refreshTokenString, err := token.GenerateTokens(accountToBeLoggedIn)
	if err != nil {
		return "", "", fmt.Errorf("error while generating tokens")
	}

	accountToBeLoggedIn.RefreshToken = refreshTokenString
	_, err = service.accountRepository.Save(accountToBeLoggedIn)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (service *accountService) Register(registrationInfo dto.AccountRegistration) (dto.CreateUserResponse, error) {
	_, err := service.accountRepository.GetByUsername(registrationInfo.Username)
	if err == nil {
		return dto.CreateUserResponse{}, fmt.Errorf("username already exists")
	}

	_, err = service.accountRepository.GetByEmail(registrationInfo.Email)
	if err == nil {
		return dto.CreateUserResponse{}, fmt.Errorf("email already exists")
	}

	hashedPassword, err := utils.HashPassword(registrationInfo.Password)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	newAccount := model.Account{
		Username:    registrationInfo.Username,
		Password:    hashedPassword,
		Email:       registrationInfo.Email,
		Role:        model.REGULAR_USER,
		IsActivated: false,
	}

	//Pravljenje user-a
	newUser := model.User{
		Name:    registrationInfo.Name,
		Surname: registrationInfo.Surname,
		Address: registrationInfo.Address,
	}

	userId, err := service.userRepository.Create(&newUser)
	if err != nil {
		return dto.CreateUserResponse{}, fmt.Errorf("an error has occured while trying to create the user")
	}

	newAccount.UserID = userId

	err = PrepareAndSendConfirmationEmail(&newAccount)
	if err != nil {
		return dto.CreateUserResponse{}, fmt.Errorf("an error has occured while sending the verification email")
	}

	id, err1 := service.accountRepository.Create(&newAccount)
	if err1 != nil {
		return dto.CreateUserResponse{}, fmt.Errorf("an error has occured while trying to create your account")
	}

	newAccount.ID = id

	response := dto.CreateUserResponse{
		ID:          newAccount.ID,
		Role:        newAccount.Role,
		IsActivated: newAccount.IsActivated,
	}

	return response, nil
}

// email activation logic
func PrepareAndSendConfirmationEmail(account *model.Account) error {

	// return time to the nanosecond (1 billionth of a sec)
	rand.Seed(time.Now().UnixNano())
	// create random code for email
	// Go rune data type represent Unicode characters
	var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	emailVerRandRune := make([]rune, 64)
	// create a random slice of runes (characters) to create our emailVerPassword (random string of characters)
	for i := 0; i < 64; i++ {
		emailVerRandRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}

	emailVerPassword := string(emailVerRandRune)
	var emailVerPWhash []byte
	// func GenerateFromPassword(password []byte, const int) ([]byte, error)
	emailVerPWhash, err := bcrypt.GenerateFromPassword([]byte(emailVerPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	account.EmailVerificationHash = string(emailVerPWhash)

	// create u.timeout after 48 hours
	timeout := time.Now().Local().AddDate(0, 0, 2)
	account.VerificationTimeout = timeout

	// TODO Stefan: add a transaction, in case the email sending fails abort database modification
	err = utils.SendConfirmationEmail(*account, emailVerPassword)
	if err != nil {
		return err
	}

	return nil
}

func (service *accountService) GetAll() (model.Accounts, error) {
	return service.accountRepository.GetAll()
}

func (service *accountService) GetById(id primitive.ObjectID) (model.Account, error) {
	return service.accountRepository.GetById(id)
}

func (service *accountService) GetByRefreshToken(token string) (model.Account, error) {
	return service.accountRepository.GetByRefreshToken(token)
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
