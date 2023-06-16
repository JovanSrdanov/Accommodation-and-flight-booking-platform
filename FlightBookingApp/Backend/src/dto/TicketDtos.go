package dto

import (
	"FlightBookingApp/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TicketFullInfo struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Buyer      model.User         `json:"buyer,omitempty" binding:"required" bson:"buyer"`
	Owner      model.User         `json:"owner,omitempty" binding:"required" bson:"owner"`
	FlightId   primitive.ObjectID `json:"flightId" binding:"required" bson:"flightId"`
	FlightInfo model.Flight       `json:"flightInfo" bson:"flightInfo"`
}

type BuyTicketDto struct {
	FlightId        primitive.ObjectID `json:"flightId" binding:"required" bson:"flightId"`
	NumberOfTickets int32              `json:"numberOfTickets" binding:"required,min=1" bson:"numberOfTickets"`
}

type BuyTicketApiKeyDto struct {
	ApiKey          string             `json:"apiKey" binding:"required"`
	FlightId        primitive.ObjectID `json:"flightId" binding:"required" bson:"flightId"`
	NumberOfTickets int32              `json:"numberOfTickets" binding:"required,min=1" bson:"numberOfTickets"`
}
