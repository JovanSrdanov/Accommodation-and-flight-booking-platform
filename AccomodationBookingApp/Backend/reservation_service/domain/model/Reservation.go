package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Reservation struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	DateRange       DateRange          `json:"dateRange" binding:"required" bson:"dateRange"`
	Price           int32              `json:"price,omitempty" bson:"price"`
	NumberOfGuests  int32              `json:"numberOfGuests" binding:"required" bson:"numberOfGuests"`
	Status          string             `json:"status,omitempty" bson:"status"`
	AccommodationId primitive.ObjectID `json:"accommodationId" binding:"required" bson:"accommodationId"`
	GuestId         string             `json:"guestId,omitempty" bson:"guestId"`
}

type Reservations []*Reservation
