package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ticket struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Buyer    User               `json:"buyer,omitempty" bson:"buyer"`
	Owner    User               `json:"owner,omitempty" bson:"owner"`
	FlightId primitive.ObjectID `json:"flightId" binding:"required" bson:"flightId"`
}

func NewTicket(Buyer User, Owner User, FlightId primitive.ObjectID) Ticket {
	return Ticket{
		Buyer:    Buyer,
		Owner:    Owner,
		FlightId: FlightId,
	}
}

type Tickets []*Ticket
