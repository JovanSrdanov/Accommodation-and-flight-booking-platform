package repository

import (
	"api_gateway/domain/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type UniqueVisitorRepositoryMongo struct {
	dbClient *mongo.Client
}

func NewUniqueVisitorRepositoryMongo(dbClient *mongo.Client) (*UniqueVisitorRepositoryMongo, error) {
	return &UniqueVisitorRepositoryMongo{
		dbClient: dbClient,
	}, nil
}

func (repo UniqueVisitorRepositoryMongo) CreateUniqueVisitor(newUniqueVisitor *model.UniqueVisitor) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := repo.getCollectionUniqueVisitor()

	newUniqueVisitor.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, &newUniqueVisitor)
	if err != nil {
		log.Println(err)
		return primitive.NilObjectID, err
	}

	return newUniqueVisitor.ID, nil
}

func (repo UniqueVisitorRepositoryMongo) GetVisitorByIpAndBrowser(ipAddress, browser string) (*model.UniqueVisitor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println(ipAddress + " " + browser)

	collection := repo.getCollectionUniqueVisitor()
	filter := bson.M{
		"ipAddress": ipAddress,
		"browser":   browser,
	}

	var visitor model.UniqueVisitor
	err := collection.FindOne(ctx, filter).Decode(&visitor)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Fatal(err)
		return nil, err
	}

	return &visitor, nil
}

func (repo UniqueVisitorRepositoryMongo) getCollectionUniqueVisitor() *mongo.Collection {
	db := repo.dbClient.Database("uniqueSiteVisitorsDb")
	collection := db.Collection("visitors")

	return collection
}
