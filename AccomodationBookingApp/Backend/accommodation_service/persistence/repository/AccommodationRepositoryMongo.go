package repository

import (
	"accommodation_service/domain/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
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
	accommodation.ID = primitive.NewObjectID()

	//TODO Strahinja: ove iz pasetoa izvuci id

	accommodation.HostId = "paseto"

	log.Println(accommodation)

	result, err := collection.InsertOne(ctx, &accommodation)
	if err != nil {
		log.Println(err)
		return primitive.ObjectID{}, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
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

func (repo AccommodationRepositoryMongo) GetAll() (model.Accommodations, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var accommodations model.Accommodations
	err = cursor.All(ctx, &accommodations)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return accommodations, nil
}

func (repo AccommodationRepositoryMongo) GetAllMy() (model.Accommodations, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	//TODO Strahinja: Iz pasetoa izvuci id ulogovanog hosta

	filter := bson.D{{"hostId", "paseto"}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var accommodations model.Accommodations
	err = cursor.All(ctx, &accommodations)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return accommodations, nil
}

func (repo AccommodationRepositoryMongo) GetAmenities() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := repo.dbClient.Database("accommodationDb")
	collection := db.Collection("amenities")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var amenities model.Amenities
	err = cursor.All(ctx, &amenities)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var amenitiesString []string
	for _, value := range amenities {
		amenitiesString = append(amenitiesString, value.Name)
	}

	return amenitiesString, nil
}

func (repo AccommodationRepositoryMongo) getCollection() *mongo.Collection {
	db := repo.dbClient.Database("accommodationDb")
	collection := db.Collection("accommodations")
	return collection
}
