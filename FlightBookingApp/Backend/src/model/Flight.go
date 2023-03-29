package model

import (
	"FlightBookingApp/errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flight struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id" example:"641c4e542ce3f0c9dc339384"`
	//TODO namestiti da smesta UTC
	DepartureDateTime time.Time `json:"departureDateTime" binding:"required" validate:"not-before-current-date" bson:"departureDateTime" example:"2023-07-18T21:54:42.123Z"`
	StartPoint        Airport   `json:"startPoint" binding:"required" bson:"startPoint"`
	Destination       Airport   `json:"destination" binding:"required" bson:"destination"`
	Price             float32   `json:"price" binding:"required,min=0" bson:"price" example:"3000"`
	NumberOfSeats     int32     `json:"numberOfSeats" binding:"required,min=0" bson:"numberOfSeats" example:"350"`
	VacantSeats       int32     `json:"vacantSeats"  bson:"vacantSeats" example:"32"`
	Canceled          bool      `json:"canceled"  bson:"canceled"`
}
type Flights []*Flight

func (flight *Flight) SelfValidate() error {
	if flight.StartPoint.ID == flight.Destination.ID {
		return &errors.SameStartPointAndDestinationError{}
	}
	return nil
}

func (flight *Flight) DecreaseVacantSeats(number int32) {
	if flight.VacantSeats > 0 {
		flight.VacantSeats -= number
	}
}
func (flight *Flight) IncreaseVacantSeats(number int32) {
	if flight.VacantSeats < flight.NumberOfSeats {
		flight.VacantSeats += number
	}
}

func (flight *Flight) HasPassed() bool {
	return flight.DepartureDateTime.Before(time.Now())
}
