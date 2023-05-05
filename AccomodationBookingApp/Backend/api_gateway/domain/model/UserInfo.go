package model

import "github.com/google/uuid"

type UserInfo struct {
	Username string    `json:"username"`
	UserID   uuid.UUID `json:"userId"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Email    string    `json:"email"`
	Address  Address   `json:"address"`
}
