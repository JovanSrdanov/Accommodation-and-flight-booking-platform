package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Accommodation struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name      string             `json:"name" binding:"required" bson:"name"`
	Address   Address            `json:"address" binding:"required" bson:"address"`
	MinGuests int32              `json:"minGuests" binding:"required" bson:"minGuests"`
	MaxGuests int32              `json:"maxGuests" binding:"required" bson:"maxGuests"`
	Amenities []string           `json:"amenities" binding:"required" bson:"amenities"`
	Images    []string           `json:"images" binding:"required" bson:"images"`
	HostId    string             `json:"hostId,omitempty" bson:"hostId"`
}

type Accommodations []*Accommodation
