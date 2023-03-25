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

type accountRepository struct {
	base     Repository
	Accounts []model.Account
}

type AccountRepository interface {
	Create(newAccount *model.Account) (primitive.ObjectID, error)
	GetAll() (model.Accounts, error)
	GetById(id primitive.ObjectID) (model.Account, error)
	Delete(id primitive.ObjectID) error
}

func NewAccountRepository(client *mongo.Client, logger *log.Logger) *accountRepository {
	base := NewRepository(client, logger)
	return &accountRepository{base: base}
}

func (repo *accountRepository) Create(account *model.Account) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	account.ID = primitive.NewObjectID()
	result, err := collection.InsertOne(ctx, &account)

	if err != nil {
		repo.base.logger.Println(err)
		return primitive.ObjectID{}, err
	}

	id := result.InsertedID.(primitive.ObjectID)
	repo.base.logger.Printf("Inserted entity, id = '%s'\n", id)
	return id, nil
}

func (repo *accountRepository) GetAll() (model.Accounts, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	var accounts model.Accounts
	fligtsCursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		repo.base.logger.Println(err)
		return nil, err
	}
	err = fligtsCursor.All(ctx, &accounts)
	if err != nil {
		repo.base.logger.Println(err)
		return nil, err
	}

	return accounts, nil
}

func (repo *accountRepository) GetById(id primitive.ObjectID) (model.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	result := collection.FindOne(ctx, bson.M{"_id": id})
	if result.Err() != nil {
		return model.Account{}, result.Err()
	}

	var account model.Account
	result.Decode(&account)

	return account, nil
}

func (repo *accountRepository) Delete(id primitive.ObjectID) error {
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

func (repo *accountRepository) getCollection() *mongo.Collection {
	db := repo.base.client.Database("flightDb")
	collection := db.Collection("accounts")
	return collection
}