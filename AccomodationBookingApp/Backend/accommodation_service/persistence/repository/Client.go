package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func GetClient(user string, pass string) (*mongo.Client, error) {
	dbUri := "mongodb://mongo:27017"
	log.Println(dbUri + " ALOOOO***")

	clientOptions := options.Client().ApplyURI(dbUri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return client, nil
}
