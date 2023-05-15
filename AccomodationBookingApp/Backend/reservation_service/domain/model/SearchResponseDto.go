package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type SearchResponseDto struct {
	AccommodationId *primitive.ObjectID `json:"accommodationId" binding:"required"  bson:"accommodationId"`
	Price           int32               `json:"price" binding:"required" bson:"price"`
}
