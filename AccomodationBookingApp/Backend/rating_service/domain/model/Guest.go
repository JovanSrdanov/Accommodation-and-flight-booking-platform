package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Guest struct {
	GuestId primitive.ObjectID `json:"guestId,omitempty"`
}
