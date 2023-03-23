package dto

import (
	"FlightBookingApp/model"

	"github.com/google/uuid"
)

type CreateUserResponse struct {
	ID          uuid.UUID
	Username     string `json:"username" binding:"required,alphanum"`
	Email       string `json:"email" binding:"required, email"`
	Role        model.Role   `json:"role"`
	IsActivated bool   `json:"isActivated" binding:"required"`
}