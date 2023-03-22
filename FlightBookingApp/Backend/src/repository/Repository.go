package repository

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

type Repository struct {
	client *mongo.Client
	logger *log.Logger
}

// Nothing to do with repo but contains logic that is used in all repositories
func GetClient() (*mongo.Client, error) {
	dbUri := os.Getenv("MONGO_DB_URI")
	return mongo.NewClient(options.Client().ApplyURI(dbUri))
}

func (repo *Repository) Connect(ctx context.Context) error {
	err := repo.client.Connect(ctx)
	if err != nil {
		return err
	}
	repo.logger.Println("Connected")
	return nil
}

func (repo *Repository) Disconnect(ctx context.Context) error {
	err := repo.client.Disconnect(ctx)
	if err != nil {
		return err
	}
	repo.logger.Println("Disconnected")
	return nil
}

func (repo *Repository) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := repo.client.Ping(ctx, readpref.Primary())
	if err != nil {
		repo.logger.Println(err)
	}

	// Print available databases
	databases, err := repo.client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		repo.logger.Println(err)
	}
	fmt.Println(databases)
}
