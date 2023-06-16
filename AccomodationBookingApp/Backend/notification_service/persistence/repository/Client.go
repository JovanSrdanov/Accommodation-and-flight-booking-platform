package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func GetPostgresClient(host, user, password, dbname, port string) (*gorm.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	return gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
}

func GetMongoClient(dbName, dbPort string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
