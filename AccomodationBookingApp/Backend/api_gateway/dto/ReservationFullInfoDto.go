package dto

import (
	"time"
)

type ReservationFullInfo struct {
	ID              string    `json:"id" bson:"_id"`
	DateRange       DateRange `json:"dateRange" binding:"required" bson:"dateRange"`
	Price           int32     `json:"price" bson:"price"`
	NumberOfGuests  int32     `json:"numberOfGuests" binding:"required" bson:"numberOfGuests"`
	Status          string    `json:"status" bson:"status"`
	AccommodationId string    `json:"accommodationId" binding:"required" bson:"accommodationId"`
	GuestId         string    `json:"guestId" bson:"guestId"`

	AccommodationName string  `json:"accommodationName"`
	Address           Address `json:"address"`
}

type DateRange struct {
	From time.Time `json:"from" binding:"required" bson:"from"`
	To   time.Time `json:"to" binding:"required" bson:"to"`
}
