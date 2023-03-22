package repository

import (
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type flightRepository struct {
	base    Repository
	Flights []model.Flight
}

type FlightRepository interface {
	Create(flight model.Flight) model.Flight
	GetAll() []model.Flight
	GetById(id uuid.UUID) (model.Flight, error)
	Delete(id primitive.ObjectID) error
}

// NoSQL: Constructor which reads db configuration from environment
func NewFlightRepository(ctx context.Context, logger *log.Logger) (*flightRepository, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	return &flightRepository{
		base: Repository{
			client: client,
			logger: logger,
		},
	}, nil
}

// TODO Aleksandar: napraviti connection pooling da ne bi morale konstantno da se otvaraju i zatvaraju konekcije sa bazom
func (repo *flightRepository) Create(flight model.Flight) model.Flight {
	repo.base.Connect(context.Background())
	defer repo.base.Disconnect(context.Background())
	repo.Flights = append(repo.Flights, flight)
	return flight
}

func (repo *flightRepository) GetAll() []model.Flight {
	repo.base.Connect(context.Background())
	defer repo.base.Disconnect(context.Background())
	return repo.Flights
}

func (repo *flightRepository) GetById(id uuid.UUID) (model.Flight, error) {
	repo.base.Connect(context.Background())
	defer repo.base.Disconnect(context.Background())
	for _, flight := range repo.Flights {
		if flight.ID == id {
			return flight, nil
		}
	}
	return model.Flight{}, &errors.NotFoundError{}
}
func (repo *flightRepository) Delete(id primitive.ObjectID) error {
	return nil
}
