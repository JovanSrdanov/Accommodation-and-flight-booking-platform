package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role int32

const (
	ADMIN Role = iota
	REGULAR_USER
)

//TODO Stefan: implement account email activation
type Account struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Username string `json:"username" binding:"required,alphanum" bson:"username"`
	Password string `json:"password" binding:"required,min=6" bson:"password"`
	Email string `json:"email" binding:"required, email" bson:"email"`
	Role Role `json:"role" bson:"role"`
	IsActivated bool `json:"isActivated" binding:"required" bson:"isActivated"`
}

type Accounts []*Account