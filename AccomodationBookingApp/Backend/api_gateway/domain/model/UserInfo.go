package model

import "github.com/google/uuid"

type UserInfo struct {
	Username      string    `json:"username"`
	UserProfileID uuid.UUID `json:"userProfileID"`
	Name          string    `json:"name"`
	Surname       string    `json:"surname"`
	Email         string    `json:"email"`
	Address       Address   `json:"address"`
}
