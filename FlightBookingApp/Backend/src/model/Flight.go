package model

import (
	"github.com/google/uuid"
	"time"
)

type Flight struct {
	//TODO: namestiti da automatski generise uuid
	ID uuid.UUID `json:"id, omitempty"`
	//TODO namestiti da smesta UTC
	Time        time.Time `json:"time" binding:"required"`
	StartPoint  Airport   `json:"startPoint" binding:"required"`
	Destination Airport   `json:"destination" binding:"required"`
	Price       float32   `json:"price" binding:"required"`
	VacantSeats int32     `json:"vacantSeats" binding:"required"`
}

func (flight *Flight) decreaseVacantSeats(number int32) {
	flight.VacantSeats -= number
}
func (flight *Flight) increaseVacantSeats(number int32) {
	flight.VacantSeats += number
}
