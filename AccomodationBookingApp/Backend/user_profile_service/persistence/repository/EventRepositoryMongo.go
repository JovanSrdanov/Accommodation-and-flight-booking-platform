package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"user_profile_service/event_sourcing"
)

type EventRepositoryMongo struct {
	events *mongo.Collection
}

func NewEventRepositoryMongo(client *mongo.Client, dbName, collectionName string) (*EventRepositoryMongo, error) {
	events := client.Database(dbName).Collection(collectionName)
	return &EventRepositoryMongo{
		events: events,
	}, nil
}

func (repo *EventRepositoryMongo) Save(event *event_sourcing.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.events.InsertOne(ctx, event)
	if err != nil {
		return err
	}
	return nil
}

func (repo *EventRepositoryMongo) Read(sagaId primitive.ObjectID, action string) (*event_sourcing.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"sagaId": sagaId,
		"action": action,
	}
	var result event_sourcing.Event

	err := repo.events.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, err
}

func (repo *EventRepositoryMongo) Delete(sagaId primitive.ObjectID, action string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"sagaId": sagaId,
		"action": action,
	}

	_, err := repo.events.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	return nil
}
