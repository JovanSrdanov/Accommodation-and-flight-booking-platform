package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func GetClient(dbPort string, dbName string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("DB NAME: " + dbName + ", DB PORT: " + dbPort)

	dbUri := "mongodb://" + dbName + ":" + dbPort

	log.Println("DB_URI: " + dbUri)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))

	if err != nil {
		log.Println("database connection error", err)
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("err", err)
		return nil, err
	}
	log.Println("Successfully connected and pinged.")
	return client, nil
}
