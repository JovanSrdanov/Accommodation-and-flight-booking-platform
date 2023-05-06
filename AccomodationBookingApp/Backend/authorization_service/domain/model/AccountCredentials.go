package model

import (
	"fmt"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Role int8

const (
	Host Role = iota
	Guest
)

type AccountCredentials struct {
	ID            uuid.UUID `json:"id,omitempty" gorm:"primaryKey"`
	Username      string    `json:"username" gorm:"unique" `
	Password      string    `json:"password"`
	Salt          string    `json:"salt,omitempty"`
	Role          Role      `json:"role"`
	UserProfileID uuid.UUID `json:"userProfileID"`
}

func NewAccountCredentials(username string, password string, salt string, role Role, userProfileId uuid.UUID) (*AccountCredentials, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	accountCredentials := &AccountCredentials{
		Username:      username,
		Password:      string(hashedPassword),
		Salt:          salt,
		Role:          role,
		UserProfileID: userProfileId,
	}

	return accountCredentials, nil
}

func (accountCredentials *AccountCredentials) IsPasswordCorrect(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(accountCredentials.Password), []byte(password))
	return err == nil
}

// GORM needs this for its fluent api queries
func (AccountCredentials) TableName() string {
	return "account_credentials"
}
