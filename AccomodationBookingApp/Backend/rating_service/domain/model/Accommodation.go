package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Accommodation struct {
	AccommodationId primitive.ObjectID `json:"accommodationId,omitempty"`
}
