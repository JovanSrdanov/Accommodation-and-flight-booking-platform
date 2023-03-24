package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Airport struct {
	ID      primitive.ObjectID `json:"id, omitempty" bson:"_id"`
	Name    string             `json:"name" binding:"required" bson:"name"`
	Address Address            `json:"address" binding:"required" bson:"address"`
}
