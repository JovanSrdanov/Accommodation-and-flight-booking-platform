package model

import "github.com/google/uuid"

type Role int32

const (
	ADMIN Role = iota
	REGULAR_USER
)

type Account struct {
	ID uuid.UUID `json:"id,omitempty" binding:"required"`
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Email string `json:"email" binding:"required, email"`
	Role Role `json:"role"`
	IsActivated bool `json:"isActivated" binding:"required"`
}