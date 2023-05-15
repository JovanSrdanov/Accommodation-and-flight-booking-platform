package token

import (
	"authorization_service/domain/model"
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"strconv"
	"time"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil,
			fmt.Errorf("invalid key size: must be exactly %d characters, provided key size is: %d, key is: %s", chacha20poly1305.KeySize, len(symmetricKey), symmetricKey)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(id uuid.UUID, duration time.Duration, role model.Role) (string, Payload, error) {
	payload, err := NewPayload(id, role, duration)
	if err != nil {
		return "", Payload{}, nil
	}

	footer := map[string]interface{}{
		"RoleAndExp": "role:" + strconv.Itoa(int(role)) + ", expiration date: " + payload.ExpiredAt.String(),
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, footer)
	return token, *payload, err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
