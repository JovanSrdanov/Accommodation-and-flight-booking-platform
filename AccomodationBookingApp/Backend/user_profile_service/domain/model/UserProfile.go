package model

import "github.com/google/uuid"

type UserProfile struct {
	ID        uuid.UUID `json:"id,omitempty" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	AddressID uuid.UUID `json:"address_id"`
	Address   Address   `json:"address" gorm:"foreignKey:AddressID"`
}

// GORM
func (UserProfile) TableName() string {
	return "user_profile"
}
