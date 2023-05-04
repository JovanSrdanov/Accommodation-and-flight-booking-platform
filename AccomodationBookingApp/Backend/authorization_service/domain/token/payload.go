package token

import (
	"authorization_service/domain/model"
	"errors"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload of the token
type Payload struct {
	// token ID
	//ID        uuid.UUID  `json:"id"`
	Username string `json:"username"`
	//Role      model.Role `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration, role model.Role) (*Payload, error) {
	//tokenID, err := uuid.NewRandom()
	//if err != nil {
	//	return nil, err
	//}

	payload := &Payload{
		//ID:        tokenID,
		Username: username,
		//Role:      role,
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
