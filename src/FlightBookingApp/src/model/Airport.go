package model

import (
	"github.com/google/uuid"
)

type Airport struct {
	ID      uuid.UUID `json:"iD"`
	Name    string    `json:"name"`
	Address Address   `json:"address"`
}
