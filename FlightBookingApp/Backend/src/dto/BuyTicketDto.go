package dto

import "FlightBookingApp/model"

type BuyTicketDto struct {
	Ticket          model.Ticket `json:"ticket" binding:"required" bson:"ticket"`
	NumberOfTickets int32        `json:"numberOfTickets" binding:"required" bson:"numberOfTickets"`
}
