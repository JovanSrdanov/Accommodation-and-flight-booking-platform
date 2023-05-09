package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func GetClient(user string, pass string) (*mongo.Client, error) {
	/*dbUri := "mongodb://mongo:27017"
	log.Println(dbUri + " ALOOOO***")

	clientOptions := options.Client().ApplyURI(dbUri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return client, nil*/

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Println("Gas test1")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongodb:27017"))

	log.Println("Gas test2")

	if err != nil {
		log.Println("database connection error", err)
		return nil, err
	}

	log.Println("Gas test3")

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("err", err)
		return nil, err
	}
	log.Println("Successfully connected and pinged.")
	return client, nil
}
