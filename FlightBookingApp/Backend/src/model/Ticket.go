package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ticket struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Buyer    string             `json:"buyer,omitempty" bson:"buyer"`
	Owner    string             `json:"owner,omitempty" bson:"owner"`
	FlightId primitive.ObjectID `json:"flightId" binding:"required" bson:"flightId"`
}
type Tickets []*Ticket
