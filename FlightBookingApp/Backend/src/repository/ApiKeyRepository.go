package repository

import (
	"FlightBookingApp/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type apiKeyRepository struct {
	base Repository
}

type ApiKeyRepository interface {
	Create(key *model.ApiKey) error
	GetByAccountId(id primitive.ObjectID) (model.ApiKey, error)
	GetByValue(value string) (model.ApiKey, error)
}

func NewApiKeyRepository(client *mongo.Client, logger *log.Logger) ApiKeyRepository {
	base := NewRepository(client, logger)
	return &apiKeyRepository{base: base}
}

func (repo *apiKeyRepository) Create(key *model.ApiKey) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	updateOptions := options.Replace().SetUpsert(true) //If it doesn't find it, creates new
	filter := bson.M{"_id": key.AccountId}
	_, err := collection.ReplaceOne(ctx, filter, key, updateOptions)

	if err != nil {
		repo.base.logger.Println(err)
		return err
	}

	return nil
}

func (repo *apiKeyRepository) GetByAccountId(id primitive.ObjectID) (model.ApiKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	result := collection.FindOne(ctx, bson.M{"_id": id})
	if result.Err() != nil {
		return model.ApiKey{}, result.Err()
	}

	var apiKey model.ApiKey
	result.Decode(&apiKey)

	return apiKey, nil
}

func (repo *apiKeyRepository) GetByValue(value string) (model.ApiKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	result := collection.FindOne(ctx, bson.M{"value": value})
	if result.Err() != nil {
		return model.ApiKey{}, result.Err()
	}

	var apiKey model.ApiKey
	result.Decode(&apiKey)

	return apiKey, nil
}

func (repo *apiKeyRepository) getCollection() *mongo.Collection {
	db := repo.base.client.Database("flightDb")
	collection := db.Collection("apiKeys")
	return collection
}
