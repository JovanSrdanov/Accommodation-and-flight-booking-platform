package repository

import (
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type flightRepository struct {
	base Repository
}

type FlightRepository interface {
	Create(flight *model.Flight) (primitive.ObjectID, error)
	GetAll() (model.Flights, error)
	GetById(id primitive.ObjectID) (model.Flight, error)
	Delete(id primitive.ObjectID) error
}

// NoSQL: Constructor which reads db configuration from environment
func NewFlightRepository(client *mongo.Client, logger *log.Logger) *flightRepository {
	base := NewRepository(client, logger)
	return &flightRepository{base: base}
}

func (repo *flightRepository) Create(flight *model.Flight) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	flight.ID = primitive.NewObjectID()
	result, err := collection.InsertOne(ctx, &flight)
	if err != nil {
		repo.base.logger.Println(err)
		return primitive.ObjectID{}, err
	}
	id := result.InsertedID.(primitive.ObjectID)
	repo.base.logger.Printf("Inserted entity, id = '%s'\n", id)
	return id, nil
}

func (repo *flightRepository) GetAll() (model.Flights, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		repo.base.logger.Println(err)
		return nil, err
	}

	var flights model.Flights
	err = cursor.All(ctx, &flights)
	if err != nil {
		repo.base.logger.Println(err)
		return nil, err
	}

	return flights, nil
}

func (repo *flightRepository) GetById(id primitive.ObjectID) (model.Flight, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	result := collection.FindOne(ctx, bson.M{"_id": id})
	if result.Err() != nil {
		return model.Flight{}, result.Err()
	}

	var flight model.Flight
	result.Decode(&flight)

	return flight, nil
}
func (repo *flightRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return &errors.NotFoundError{}
	}
	repo.base.logger.Printf("Deleted entity, id: %s", id.String())
	return nil
}

func (repo *flightRepository) getCollection() *mongo.Collection {
	db := repo.base.client.Database("flightDb")
	collection := db.Collection("flights")
	return collection
}
