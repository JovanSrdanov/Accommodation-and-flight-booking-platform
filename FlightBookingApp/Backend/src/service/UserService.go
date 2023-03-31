package service

import (
	"FlightBookingApp/model"
	"FlightBookingApp/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	Create(user model.User) (primitive.ObjectID, error)
	GetAll() (model.Users, error)
	GetById(id primitive.ObjectID) (model.User, error)
	Delete(id primitive.ObjectID) error
}

func NewUserService(userRepository repository.UserRepository) *userService {
	return &userService{
		userRepository: userRepository,
	}
}

func (service *userService) Create(user model.User) (primitive.ObjectID, error) {
	return service.userRepository.Create(&user)
}

func (service *userService) GetAll() (model.Users, error) {
	return service.userRepository.GetAll()
}

func (service *userService) GetById(id primitive.ObjectID) (model.User, error) {
	return service.userRepository.GetById(id)
}
func (service *userService) Delete(id primitive.ObjectID) error {
	return service.userRepository.Delete(id)
}