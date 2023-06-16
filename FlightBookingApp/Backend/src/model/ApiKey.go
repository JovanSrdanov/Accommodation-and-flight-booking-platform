package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ApiKey struct {
	AccountId      primitive.ObjectID `json:"accountId" bson:"_id"`
	Value          string             `json:"value,omitempty"  bson:"value"`
	ExpirationDate time.Time          `json:"expirationDate"  bson:"expiration_date"`
}

func (key *ApiKey) IsValid() bool {
	return key.ExpirationDate.After(time.Now())
}
