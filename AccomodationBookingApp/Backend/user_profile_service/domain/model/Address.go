package model

import "github.com/google/uuid"

type Address struct {
	ID           uuid.UUID `json:"id,omitempty" gorm:"primaryKey"`
	Country      string    `json:"country" `
	City         string    `json:"city"`
	Street       string    `json:"street"`
	StreetNumber string    `json:"street_number"`
}

// GORM
func (Address) TableName() string {
	return "address"
}
