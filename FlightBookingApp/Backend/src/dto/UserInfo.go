package dto

import (
	"FlightBookingApp/model"
)

type UserInfo struct {
	Name     string  `json:"name" binding:"required" bson:"name"`
	Surname  string  `json:"surname" binding:"required" bson:"surname"`
	Address  model.Address `json:"address" binding:"required" bson:"address"`
}