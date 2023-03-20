package service

import (
	"FlightBookingApp/model"
	"FlightBookingApp/repository"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type flightService struct {
	flightRepository repository.FlightRepository
}

type FlightService interface {
	Create(flight model.Flight) model.Flight
	GetAll() []model.Flight
	GetById(id uuid.UUID) (model.Flight, error)
	Delete(id primitive.ObjectID) error
}

func NewFlightService(flightRepository repository.FlightRepository) FlightService {
	return &flightService{
		flightRepository: flightRepository,
	}
}

func (service *flightService) Create(flight model.Flight) model.Flight {
	return service.flightRepository.Create(flight)
}

func (service *flightService) GetAll() []model.Flight {
	return service.flightRepository.GetAll()
}

func (service *flightService) GetById(id uuid.UUID) (model.Flight, error) {
	return service.flightRepository.GetById(id)
}
func (service *flightService) Delete(id primitive.ObjectID) error {
	return service.flightRepository.Delete(id)
}
