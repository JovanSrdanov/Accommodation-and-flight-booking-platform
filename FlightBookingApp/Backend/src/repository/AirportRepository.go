package repository

import (
	"FlightBookingApp/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type airportRepository struct {
	base Repository
}

type AirportRepository interface {
	GetAll() (model.Airports, error)
	GetById(id primitive.ObjectID) (model.Airport, error)
}

func NewAirportRepository(client *mongo.Client, logger *log.Logger) *airportRepository {
	base := NewRepository(client, logger)
	return &airportRepository{base: base}
}
func (repo *airportRepository) GetAll() (model.Airports, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		repo.base.logger.Println(err)
		return nil, err
	}

	var airports model.Airports
	err = cursor.All(ctx, &airports)
	if err != nil {
		repo.base.logger.Println(err)
		return nil, err
	}

	return airports, nil
}

func (repo *airportRepository) GetById(id primitive.ObjectID) (model.Airport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	result := collection.FindOne(ctx, bson.M{"_id": id})
	if result.Err() != nil {
		return model.Airport{}, result.Err()
	}

	var airport model.Airport
	result.Decode(&airport)

	return airport, nil
}

func (repo *airportRepository) getCollection() *mongo.Collection {
	db := repo.base.client.Database("flightDb")
	collection := db.Collection("airports")
	return collection
}
