package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"strconv"
	"time"
)

func GetClient(logger *log.Logger) (*mongo.Client, error) {
	dbUri := os.Getenv("MONGO_DB_URI")
	minPoolSize, _ := strconv.ParseInt(os.Getenv("CON_POOL_MIN_POOL_SIZE"), 10, 32)
	maxPoolSize, _ := strconv.ParseInt(os.Getenv("CON_POOL_MAX_POOL_SIZE"), 10, 32)
	setPoolMonitor, _ := strconv.ParseBool(os.Getenv("CONNECTION_POOL_MONITORING"))

	//dbUri := "mongodb://root:pass@localhost:27017"
	//minPoolSize, _ := 4
	//maxPoolSize, _ := 10
	//setPoolMonitor, _ := true

	options := options.Client().ApplyURI(dbUri).
		SetMinPoolSize(uint64(minPoolSize)).
		SetMaxPoolSize(uint64(maxPoolSize))

	logger.Printf("Min pool size: %d", minPoolSize)
	logger.Printf("Max pool size: %d", maxPoolSize)

	if setPoolMonitor {
		options.SetPoolMonitor(initPoolMonitor())
		logger.Println("Pool monitoring on")
	}

	return mongo.NewClient(options)
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

func initPoolMonitor() *event.PoolMonitor {
	logger := log.New(os.Stdout, "[connection-pool-monitor] ", log.LstdFlags)
	return &event.PoolMonitor{
		Event: func(e *event.PoolEvent) {
			switch e.Type {
			case event.PoolCreated:
				// mongo-connection pool creation events.
				logger.Println("Connection pool created")
			case event.ConnectionCreated:
				logger.Println("Connection created")
			case event.ConnectionClosed:
				// mongo-connection gets closed
				logger.Println("Connection closed")
			case event.GetSucceeded:
				// succesful mongo connection established by the pool
				logger.Println("New connection established")
			case event.ConnectionReturned:
				// connection instance returns to the pool after completing a mongo transaction.
				logger.Println("Connection returned to pool")
			case event.PoolCleared:
				// PoolCleared event will remove all the connections from the pool
				logger.Println("Pool cleared")
			}
		},
	}
}
