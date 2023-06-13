package dto

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type UserInfo struct {
	Username      string    `json:"username"`
	UserProfileID uuid.UUID `json:"userProfileID"`
	Name          string    `json:"name"`
	Surname       string    `json:"surname"`
	Email         string    `json:"email"`
	Address       Address   `json:"address"`
}

type BasicUserInfo struct {
	HostId   string `json:"hostId"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
}

type Role int8

const (
	Host Role = iota
	Guest
)

type CreateUser struct {
	//AccountCredentials
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Role     Role   `json:"role,omitempty"`
	//UserProfile
	Name    string  `json:"name,omitempty"`
	Surname string  `json:"surname,omitempty"`
	Email   string  `json:"email,omitempty"`
	Address Address `json:"address"`
}

func (role *Role) UnmarshalJSON(data []byte) error {
	var roleString string
	err := json.Unmarshal(data, &roleString)
	if err != nil {
		return err
	}
	switch strings.ToLower(roleString) {
	case "guest":
		*role = Guest
	case "host":
		*role = Host
	default:
		return fmt.Errorf("invalid Role value: %s", roleString)
	}
	return nil
}
