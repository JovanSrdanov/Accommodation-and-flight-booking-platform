package repository

import (
	"accommodation_service/domain/model"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccommodationRepositoryMongo struct {
	dbClient *mongo.Client
}

func NewAccommodationRepositoryMongo(dbClient *mongo.Client) (*AccommodationRepositoryMongo, error) {
	return &AccommodationRepositoryMongo{dbClient: dbClient}, nil
}

func (repo AccommodationRepositoryMongo) Create(accommodation *model.Accommodation) (primitive.ObjectID, error) {
	fmt.Println(accommodation.Name + " TEST*********")
	accommodation.ID = primitive.NewObjectID()
	return accommodation.ID, nil
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
