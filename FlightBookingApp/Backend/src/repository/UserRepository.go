package repository

import (
	"FlightBookingApp/errors"
	"FlightBookingApp/model"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	base  Repository
	Users []model.User
}

type UserRepository interface {
	Create(user *model.User) (primitive.ObjectID, error)
	GetAll() (model.Users, error)
	GetById(id primitive.ObjectID) (model.User, error)
	Delete(id primitive.ObjectID) error
}

func NewUserRepository(client *mongo.Client, logger *log.Logger) *userRepository {
	base := NewRepository(client, logger)
	return &userRepository{base: base}
}

func (repo *userRepository) Create(user *model.User) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	user.ID = primitive.NewObjectID()
	result, err := collection.InsertOne(ctx, &user)

	if err != nil {
		repo.base.logger.Println(err)
		return primitive.ObjectID{}, err
	}
	
	id := result.InsertedID.(primitive.ObjectID)
	repo.base.logger.Printf("Inserted entity, id = '%s'\n", id)
	return id, nil
}

func (repo *userRepository) GetAll() (model.Users, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	var users model.Users
	fligtsCursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		repo.base.logger.Println(err)
		return nil, err
	}
	err = fligtsCursor.All(ctx, &users)
	if err != nil {
		repo.base.logger.Println(err)
		return nil, err
	}

	return users, nil
}

func (repo *userRepository) GetById(id primitive.ObjectID) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	result := collection.FindOne(ctx, bson.M{"_id": id})
	if result.Err() != nil {
		return model.User{}, result.Err()
	}

	var user model.User
	result.Decode(&user)

	return user, nil
}

func (repo *userRepository) getCollection() *mongo.Collection {
	db := repo.base.client.Database("flightDb")
	collection := db.Collection("users")
	return collection
}

func (repo *userRepository) Delete(id primitive.ObjectID) error {
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