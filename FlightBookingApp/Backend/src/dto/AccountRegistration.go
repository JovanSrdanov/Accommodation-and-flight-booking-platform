package dto

import (
	"FlightBookingApp/model"
)

type AccountRegistration struct {
	Username string  `json:"username" binding:"required" bson:"username"`
	Password string  `json:"password" binding:"required" bson:"password"`
	Email    string  `json:"email" bindong:"required, email" bson:"email"`
	Name     string  `json:"name" binding:"required" bson:"name"`
	Surname  string  `json:"surname" binding:"required" bson:"surname"`
	Address  model.Address `json:"address" binding:"required" bson:"address"`
}