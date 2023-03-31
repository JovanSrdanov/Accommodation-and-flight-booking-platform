package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name string `json:"name" binding:"required" bson:"name"`
	Surname string `json:"surname" binding:"required" bson:"surname"`
	Address Address `json:"address" binding:"required" bson:"address"` 
}

type Users []*User