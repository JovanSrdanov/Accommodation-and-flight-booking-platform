package dto

import "time"

type CreateApiKeyDto struct {
	ExpirationDate time.Time `json:"expirationDate"`
}
