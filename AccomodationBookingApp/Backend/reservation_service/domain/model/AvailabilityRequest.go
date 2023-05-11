package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type AvailabilityRequest struct {
	PriceWithDate   PriceWithDate      `json:"priceWithDate" binding:"required" bson:"priceWithDate"`
	AccommodationId primitive.ObjectID `json:"accommodationId" binding:"required" bson:"accommodationId"`
}
