package token

import (
	"authorization_service/domain/model"
	"github.com/google/uuid"
	"time"
)

// Maker is an interface for managing tokens
type Maker interface {
	CreateToken(id uuid.UUID, duration time.Duration, role model.Role) (string, model.Role, error)
	VerifyToken(token string) (*Payload, error)
}
