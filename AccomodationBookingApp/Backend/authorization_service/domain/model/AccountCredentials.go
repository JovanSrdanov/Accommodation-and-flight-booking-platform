package model

import (
	"github.com/google/uuid"
	_ "github.com/google/uuid"
)

type Role int8

const (
	Host Role = iota
	Guest
)

type AccountCredentials struct {
	ID       uuid.UUID `json:"id,omitempty" gorm:"primaryKey"`
	Email    string    `json:"email" gorm:"unique" `
	Password string    `json:"password"`
	Salt     string    `json:"salt,omitempty"`
	Role     Role      `json:"role"`
}

// GORM needs this for its fluent api queries
func (AccountCredentials) TableName() string {
	return "account_credentials"
}
