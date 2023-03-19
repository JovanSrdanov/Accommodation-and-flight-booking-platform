package repository

import (
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type flightRepository struct {
	Flights []model.Flight
}

type FlightRepository interface {
	Create(flight model.Flight) model.Flight
	GetAll() []model.Flight
	GetById(id uuid.UUID) (model.Flight, error)
	Delete(id primitive.ObjectID) error
}

// TODO srediti kad se bude radilo sa pravom bazom
func NewFlightRepository() FlightRepository {
	return &flightRepository{}
}

func (repository *flightRepository) Create(flight model.Flight) model.Flight {
	repository.Flights = append(repository.Flights, flight)
	return flight
}

func (repository *flightRepository) GetAll() []model.Flight {
	return repository.Flights
}

func (repository *flightRepository) GetById(id uuid.UUID) (model.Flight, error) {
	for _, flight := range repository.Flights {
		if flight.ID == id {
			return flight, nil
		}
	}
	return model.Flight{}, &errors.NotFoundError{}
}
func (repository *flightRepository) Delete(id primitive.ObjectID) error {
	return nil
}
