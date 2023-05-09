package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Accommodation struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name string             `json:"name" binding:"required" bson:"name"`
}

type Accommodations []*Accommodation
