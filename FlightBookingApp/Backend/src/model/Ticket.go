package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ticket struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Buyer    Account            `json:"buyer" binding:"required" bson:"buyer"`
	Owner    Account            `json:"owner" binding:"required" bson:"owner"`
	Flight   Flight             `json:"-" bson:"-"`
	FlightId primitive.ObjectID `json:"flightId" binding:"required" bson:"flightId"`
}
type Tickets []*Ticket
