package repository

import (
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

type flightRepository struct {
	Flights []model.Flight
	cli     *mongo.Client
	logger  *log.Logger
}

type FlightRepository interface {
	Create(flight model.Flight) model.Flight
	GetAll() []model.Flight
	GetById(id uuid.UUID) (model.Flight, error)
	Delete(id primitive.ObjectID) error
}

// NoSQL: Constructor which reads db configuration from environment
func NewFlightRepository(ctx context.Context, logger *log.Logger) (*flightRepository, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	logger.Println("Successful connection")

	return &flightRepository{
		cli:    client,
		logger: logger,
	}, nil

}

// Disconnect from database
func (pr *flightRepository) Disconnect(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (repo *flightRepository) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := repo.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		repo.logger.Println(err)
	}

	// Print available databases
	databases, err := repo.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		repo.logger.Println(err)
	}
	fmt.Println(databases)
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
