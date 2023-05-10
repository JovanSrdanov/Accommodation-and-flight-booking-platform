package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UpdatePriceAndDate struct {
	PriceWithDate   PriceWithDate      `json:"priceWithDate" binding:"required" bson:"priceWithDate"`
	AccommodationId primitive.ObjectID `json:"availabilityId" binding:"required" bson:"availabilityId"`
}
