package model

import (
	"github.com/google/uuid"
	"time"
)

type Flight struct {
	ID          uuid.UUID `json:"id"`
	Time        time.Time `json:"time"`
	StartPoint  Airport   `json:"startPoint"`
	Destination Airport   `json:"destination"`
	//TODO u kojoj valuti?
	Price       float32 `json:"price"`
	VacantSeats int32   `json:"vacantSeats"`
}

func (flight *Flight) decreaseVacantSeats(number int32) {
	flight.VacantSeats -= number
}
func (flight *Flight) increaseVacantSeats(number int32) {
	flight.VacantSeats += number
}
