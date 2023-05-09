package repository

import (
	"accommodation_service/domain/model"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type AccommodationRepositoryMongo struct {
	dbClient *mongo.Client
}

func NewAccommodationRepositoryMongo(dbClient *mongo.Client) (*AccommodationRepositoryMongo, error) {
	return &AccommodationRepositoryMongo{dbClient: dbClient}, nil
}

func (repo AccommodationRepositoryMongo) Create(accommodation *model.Accommodation) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	log.Println("UZEO COLLECTION")

	accommodation.ID = primitive.NewObjectID()

	log.Println(accommodation)

	result, err := collection.InsertOne(ctx, &accommodation)
	if err != nil {
		log.Println("GAS: " + err.Error())
		return primitive.ObjectID{}, err
	}
	id := result.InsertedID.(primitive.ObjectID)
	log.Println("Inserted entity, id = '%s'\n", id)
	return id, nil
}

func (repo AccommodationRepositoryMongo) Delete(id primitive.ObjectID) error {
	return nil
}

func (repo AccommodationRepositoryMongo) Update(accommodation *model.Accommodation) (*model.Accommodation, error) {
	return accommodation, nil
}

func (repo AccommodationRepositoryMongo) GetById(id primitive.ObjectID) (*model.Accommodation, error) {
	return &model.Accommodation{}, nil
}

func (repo AccommodationRepositoryMongo) getCollection() *mongo.Collection {
	db := repo.dbClient.Database("accommodationDb")
	collection := db.Collection("accommodations")
	return collection
}
