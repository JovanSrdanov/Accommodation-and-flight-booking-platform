package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Repository struct {
	client *mongo.Client
	logger *log.Logger
}

func NewRepository(client *mongo.Client, logger *log.Logger) Repository {
	repository := Repository{
		client: client,
		logger: logger,
	}
	return repository
}
