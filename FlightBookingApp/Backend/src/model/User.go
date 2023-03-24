package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `json:"id,omitempty" binding:"required"  bson:"_id"`
	Name        string `json:"name" binding:"required" bson:"name"`
	Surname     string `json:"surname" binding:"required" bson:"surname"`
	PhoneNumber string `json:"phoneNumber" binding:"required" bson:"phoneNumber"`
}

type Users []*User