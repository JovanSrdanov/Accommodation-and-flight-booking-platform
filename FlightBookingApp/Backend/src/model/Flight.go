package model

import (
	"FlightBookingApp/errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flight struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	//TODO namestiti da smesta UTC
	//TODO Aleksandar (Jovan napisao) , validacija na time, destination i price, ddd na decrease i increase , ne sme da ima negativno dostupnih mesta
	DepartureDateTime time.Time `json:"time" binding:"required" validate:"not-before-current-date" bson:"time"`
	StartPoint        Airport   `json:"startPoint" binding:"required" bson:"startPoint"`
	Destination       Airport   `json:"destination" binding:"required" bson:"destination"`
	Price             float32   `json:"price" binding:"required,min=0" bson:"price"`
	NumberOfSeats     int32     `json:"numberOfSeats" binding:"required,min=0" bson:"numberOfSeats"`
	VacantSeats       int32     `json:"vacantSeats"  bson:"vacantSeats"`
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
