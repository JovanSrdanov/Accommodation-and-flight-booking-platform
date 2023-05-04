package token

import (
	"authorization_service/domain/model"
	"time"
)

// Maker is an interface for managing tokens
type Maker interface {
	CreateToken(username string, duration time.Duration, role model.Role) (string, error)
	VerifyToken(token string) (*Payload, error)
}
