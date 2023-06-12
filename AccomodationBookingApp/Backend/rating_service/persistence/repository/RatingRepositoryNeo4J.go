package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rating_service/domain/model"
)

type RatingRepositoryNeo4J struct {
	dbClient *neo4j.DriverWithContext
}

func NewRatingRepositoryNeo4J(dbClient *neo4j.DriverWithContext) (*RatingRepositoryNeo4J, error) {
	return &RatingRepositoryNeo4J{dbClient: dbClient}, nil
}

func (repo RatingRepositoryNeo4J) RateAccommodation(rating *model.Rating) error {
	return nil
}

func (repo RatingRepositoryNeo4J) GetRatingForAccommodation(id primitive.ObjectID) (model.RatingResponse, error) {
	return model.RatingResponse{}, nil
}
