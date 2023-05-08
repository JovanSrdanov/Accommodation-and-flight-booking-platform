package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetClient(user string, pass string) (*mongo.Client, error) {
	dbUri := "mongodb://" + user + ":" + pass + "@mongo:27017"

	clientOptions := options.Client().ApplyURI(dbUri)
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		return nil, err
	}

	return client, nil
}
