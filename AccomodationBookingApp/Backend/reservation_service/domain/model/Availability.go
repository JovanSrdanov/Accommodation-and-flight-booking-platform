package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Availability struct {
	ID                     primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	AccommodationId        primitive.ObjectID `json:"accommodationId" binding:"required" bson:"accommodationId"`
	HostId                 string             `json:"hostId" binding:"required" bson:"hostId"`
	IsAutomaticReservation bool               `json:"isAutomaticReservation" binding:"required" bson:"isAutomaticReservation"`
	AvailableDates         []*PriceWithDate   `json:"availableDates,omitempty" bson:"availableDates"`
}

type Availabilities []*Availability
