package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Airport struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id" example:"641d71cedd5e51a214a40c14"`
	Name    string             `json:"name" binding:"required" bson:"name" example:"Nikola Tesla"`
	Address Address            `json:"address" binding:"required" bson:"address"`
}

type Airports []*Airport
