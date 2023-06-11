package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Rating struct {
	AccommodationId primitive.ObjectID `json:"accommodationId,omitempty"`
	GuestId         primitive.ObjectID `json:"guestId,omitempty"`
	Rating          int32              `json:"rating"`
	Date            time.Time          `json:"date"`
}
