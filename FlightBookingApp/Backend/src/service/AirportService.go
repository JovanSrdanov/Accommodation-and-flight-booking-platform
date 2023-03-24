package service

import (
	"FlightBookingApp/model"
	"FlightBookingApp/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type airportService struct {
	airportRepository repository.AirportRepository
}

type AirportService interface {
	GetAll() (model.Airports, error)
	GetById(id primitive.ObjectID) (model.Airport, error)
}

func NewAirportService(airportRepository repository.AirportRepository) *airportService {
	return &airportService{
		airportRepository: airportRepository,
	}
}
func (service *airportService) GetAll() (model.Airports, error) {
	return service.airportRepository.GetAll()
}

func (service *airportService) GetById(id primitive.ObjectID) (model.Airport, error) {
	return service.airportRepository.GetById(id)
}
