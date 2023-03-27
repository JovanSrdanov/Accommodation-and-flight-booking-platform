package repository

import (
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type ticketRepository struct {
	base Repository
}

type TicketRepositry interface {
	Create(ticket *model.Ticket) (primitive.ObjectID, error)
	GetAll() (model.Tickets, error)
	GetById(id primitive.ObjectID) (model.Ticket, error)
	Delete(id primitive.ObjectID) error
}

func NewTicketRepositry(client *mongo.Client, logger *log.Logger) *ticketRepository {
	base := NewRepository(client, logger)
	return &ticketRepository{base: base}
}

func (repo *ticketRepository) Create(ticket *model.Ticket) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	ticket.ID = primitive.NewObjectID()
	result, err := collection.InsertOne(ctx, &ticket)
	if err != nil {
		repo.base.logger.Println(err)
		return primitive.ObjectID{}, err
	}
	id := result.InsertedID.(primitive.ObjectID)
	repo.base.logger.Printf("Inserted entity, id = '%s'\n", id)
	return id, nil
}

func (repo *ticketRepository) GetAll() (model.Tickets, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		repo.base.logger.Println(err)
		return nil, err
	}

	var tickets model.Tickets
	err = cursor.All(ctx, &tickets)
	if err != nil {
		repo.base.logger.Println(err)
		return nil, err
	}

	return tickets, nil
}

func (repo *ticketRepository) GetById(id primitive.ObjectID) (model.Ticket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	result := collection.FindOne(ctx, bson.M{"_id": id})
	if result.Err() != nil {
		return model.Ticket{}, result.Err()
	}

	var ticket model.Ticket
	result.Decode(&ticket)

	return ticket, nil
}
func (repo *ticketRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return &errors.NotFoundError{}
	}
	repo.base.logger.Printf("Deleted entity, id: %s", id.String())
	return nil
}

func (repo *ticketRepository) getCollection() *mongo.Collection {
	db := repo.base.client.Database("flightDb")
	collection := db.Collection("tickets")
	return collection
}
