package service

import (
	"FlightBookingApp/model"
	"FlightBookingApp/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type flightService struct {
	flightRepository repository.FlightRepository
}

type FlightService interface {
	Create(flight model.Flight) (primitive.ObjectID, error)
	GetAll() (model.Flights, error)
	GetById(id primitive.ObjectID) (model.Flight, error)
	Delete(id primitive.ObjectID) error
}

func NewFlightService(flightRepository repository.FlightRepository) FlightService {
	return &flightService{
		flightRepository: flightRepository,
	}
}

func (service *flightService) Create(flight model.Flight) (primitive.ObjectID, error) {
	return service.flightRepository.Create(&flight)
}

func (service *flightService) GetAll() (model.Flights, error) {
	return service.flightRepository.GetAll()
}

func (service *flightService) GetById(id primitive.ObjectID) (model.Flight, error) {
	return service.flightRepository.GetById(id)
}
func (service *flightService) Delete(id primitive.ObjectID) error {
	return service.flightRepository.Delete(id)
}
