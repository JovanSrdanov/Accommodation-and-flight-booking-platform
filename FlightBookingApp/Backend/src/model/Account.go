package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role int32

const (
	ADMIN Role = iota
	REGULAR_USER
)

// TODO Stefan: separate into account and activation structs
type Account struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Username string             `json:"username" binding:"required,min=3,max=23,alphanum" bson:"username"`
	Password string             `json:"password" binding:"required,min=8,max=24" bson:"password"`
	Role     Role               `json:"role" bson:"role"`
	//Email verification
	Email                 string    `json:"email" binding:"required, email" bson:"email"`
	EmailVerificationHash string    `json:"emailVerificationHash" bson:"emailVerificationHash"`
	VerificationTimeout   time.Time `json:"verificationTimeout" bson:"verificationTimeout"`
	IsActivated           bool      `json:"isActivated" binding:"required" bson:"isActivated"`
	//Tokens
	RefreshToken string             `json:"refreshToken" bson:"refreshToken"`
	UserID       primitive.ObjectID `json:"userId" bson:"userId"`
}

type Accounts []*Account
