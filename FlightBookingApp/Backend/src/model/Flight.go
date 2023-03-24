package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Flight struct {
	ID primitive.ObjectID `json:"id, omitempty" bson:"_id"`
	//TODO namestiti da smesta UTC
	Time        time.Time `json:"time" binding:"required" bson:"time"`
	StartPoint  Airport   `json:"startPoint" binding:"required" bson:"startPoint"`
	Destination Airport   `json:"destination" binding:"required" bson:"destination"`
	Price       float32   `json:"price" binding:"required" bson:"price"`
	VacantSeats int32     `json:"vacantSeats" binding:"required" bson:"vacantSeats"`
}
type Flights []*Flight

func (flight *Flight) decreaseVacantSeats(number int32) {
	flight.VacantSeats -= number
}
func (flight *Flight) increaseVacantSeats(number int32) {
	flight.VacantSeats += number
}
