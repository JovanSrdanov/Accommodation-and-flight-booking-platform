package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

func GetClient() (*mongo.Client, error) {
	dbUri := os.Getenv("MONGO_DB_URI")
	//dbUri := "mongodb://root:pass@localhost:27017"
	return mongo.NewClient(options.Client().ApplyURI(dbUri))
}

func Connect(ctx context.Context, client *mongo.Client, logger *log.Logger) error {
	err := client.Connect(ctx)
	if err != nil {
		return err
	}
	logger.Println("Connected")
	return nil
}

func Disconnect(ctx context.Context, client *mongo.Client, logger *log.Logger) error {
	err := client.Disconnect(ctx)
	if err != nil {
		return err
	}
	logger.Println("Disconnected")
	return nil
}

func Ping(client *mongo.Client, logger log.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Println(err)
	}

	// Print available databases
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		logger.Println(err)
	}
	fmt.Println(databases)
}
