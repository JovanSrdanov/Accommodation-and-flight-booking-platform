package model

import (
	"github.com/google/uuid"
)

type Airport struct {
	//TODO: namestiti da automatski generise uuid
	ID      uuid.UUID `json:"id,omitempty"`
	Name    string    `json:"name" binding:"required"`
	Address Address   `json:"address" binding:"required"`
}
