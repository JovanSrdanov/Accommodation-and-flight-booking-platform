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
	if primitive.ObjectID.IsZero(accommodation.ID) {
		accommodation.ID = primitive.NewObjectID()
	}

	_, err := collection.InsertOne(ctx, &accommodation)
	if err != nil {
		log.Println(err)
		return primitive.ObjectID{}, err
	}

	return accommodation.ID, nil
}

func (repo AccommodationRepositoryMongo) Delete(id primitive.ObjectID) error {
	return nil
}

func (repo AccommodationRepositoryMongo) DeleteByHostId(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	filter := bson.D{{"hostId", id}}

	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (repo AccommodationRepositoryMongo) Update(accommodation *model.Accommodation) (*model.Accommodation, error) {
	return accommodation, nil
}

func (repo AccommodationRepositoryMongo) GetById(id primitive.ObjectID) (*model.Accommodation, error) {
	collection := repo.getCollection()

	filter := bson.M{"_id": id}

	// Find the document by ID
	var result model.Accommodation
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
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

func (repo AccommodationRepositoryMongo) GetAllMy(hostId string) (model.Accommodations, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	filter := bson.D{{"hostId", hostId}}
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

func (repo AccommodationRepositoryMongo) SearchAccommodation(searchDto *model.SearchDto) (model.Accommodations, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollection()

	//TODO: Prosiriti da gleda i hostId gde je host promoted?
	filter := bson.M{
		"$or": []bson.M{
			{"address.country": bson.M{"$regex": primitive.Regex{Pattern: searchDto.Location, Options: "i"}}},
			{"address.city": bson.M{"$regex": primitive.Regex{Pattern: searchDto.Location, Options: "i"}}},
			{"address.street": bson.M{"$regex": primitive.Regex{Pattern: searchDto.Location, Options: "i"}}},
		},
		//"amenities": bson.M{"$all": searchDto.Amenities},
		"minGuests": bson.M{"$lte": searchDto.MinGuests},
		"maxGuests": bson.M{"$gte": searchDto.MinGuests},
	}

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
