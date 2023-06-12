package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type RecommendedAccommodations struct {
	AccommodationsIds []primitive.ObjectID `json:"accommodationsIds"`
}
