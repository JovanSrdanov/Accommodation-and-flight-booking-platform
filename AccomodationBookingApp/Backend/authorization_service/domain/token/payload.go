package token

import (
	"authorization_service/domain/model"
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload of the token
type Payload struct {
	// token ID
	ID        uuid.UUID  `json:"id"`
	Role      model.Role `json:"role"`
	IssuedAt  time.Time  `json:"issued_at"`
	ExpiredAt time.Time  `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(id uuid.UUID, role model.Role, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		ID:        id,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
