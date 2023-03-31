package service

import (
	"FlightBookingApp/model"
	"FlightBookingApp/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type jwtService struct {
	accountRepository repository.AccountRepository
	userRepository    repository.UserRepository
}

type JwtService interface {
	GetUser(userAccountID any) (model.User, error)
}

func NewJwtService(accountRepository repository.AccountRepository, userRepository repository.UserRepository) *jwtService {
	return &jwtService{
		accountRepository: accountRepository,
		userRepository:    userRepository,
	}
}

func (service *jwtService) GetUser(userAccountID any) (model.User, error) {

	userAccount, err := service.accountRepository.GetById(userAccountID.(primitive.ObjectID))
	if err != nil {
		return model.User{}, err
	}

	user, err := service.userRepository.GetById(userAccount.UserID)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
