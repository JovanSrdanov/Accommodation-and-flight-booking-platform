package model

import "github.com/google/uuid"

type User struct {
	ID uuid.UUID `json:"id,omitempty" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Surname     string `json:"surname" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}