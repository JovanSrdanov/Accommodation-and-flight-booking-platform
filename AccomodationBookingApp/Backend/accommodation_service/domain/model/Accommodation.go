package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Accommodation struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name      string             `json:"name" binding:"required" bson:"name"`
	Location  string             `json:"location" binding:"required" bson:"location"`
	MinGuests int32              `json:"minGuests" binding:"required" bson:"minGuests"`
	MaxGuests int32              `json:"maxGuests" binding:"required" bson:"maxGuests"`
	Amenities []string           `json:"amenities" binding:"required" bson:"amenities"`
}

type Accommodations []*Accommodation
