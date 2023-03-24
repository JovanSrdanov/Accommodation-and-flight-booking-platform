package main

import (
	"context"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strconv"
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

	clientOptions := options.Client().ApplyURI(dbUri).
		SetMinPoolSize(uint64(minPoolSize)).
		SetMaxPoolSize(uint64(maxPoolSize))

	logger.Printf("Min pool size: %d", minPoolSize)
	logger.Printf("Max pool size: %d", maxPoolSize)

	if setPoolMonitor {
		clientOptions.SetPoolMonitor(initPoolMonitor())
		logger.Println("Pool monitoring on")
	}

	return mongo.NewClient(clientOptions)
}

func Connect(ctx context.Context, client *mongo.Client, logger *log.Logger) {
	err := client.Connect(ctx)
	if err != nil {
		logger.Println(err.Error())
	}
	logger.Println("Connected")
}

func Disconnect(ctx context.Context, client *mongo.Client, logger *log.Logger) {
	err := client.Disconnect(ctx)
	if err != nil {
		logger.Println(err.Error())
	}
	logger.Println("Disconnected")
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
