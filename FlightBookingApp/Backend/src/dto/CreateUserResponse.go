package dto

import (
	"FlightBookingApp/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserResponse struct {
	ID          primitive.ObjectID
	Role        model.Role   `json:"role"`
	IsActivated bool   `json:"isActivated" binding:"required"`
}