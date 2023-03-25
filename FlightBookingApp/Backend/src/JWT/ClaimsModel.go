package jwt

import (
	"FlightBookingApp/model"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccessJwtClaims struct {
	ID primitive.ObjectID `json:"id,omitempty"`
	AccessID string `json:"accessId,omitempty"`
	Roles    []model.Role  `json:"roles,omitempty"`
	jwt.StandardClaims
}

type RefreshJwtClaims struct {
	ID primitive.ObjectID `json:"id,omitempty"`
	RefreshID string `json:"refreshId,omitempty"`
	Roles    []model.Role  `json:"roles,omitempty"`
	jwt.StandardClaims
}

func (claims AccessJwtClaims) Valid() error {
	var now = time.Now().UTC().Unix()
	if claims.VerifyExpiresAt(now, true){
		return nil
	}

	return fmt.Errorf("token is invalid")
}