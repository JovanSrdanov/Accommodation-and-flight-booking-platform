package jwt

import (
	"FlightBookingApp/model"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JwtClaims struct {
	ID primitive.ObjectID `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Roles    []model.Role  `json:"roles,omitempty"`
	jwt.StandardClaims
}

func (claims JwtClaims) Valid() error {
	var now = time.Now().UTC().Unix()
	if claims.VerifyExpiresAt(now, true){
		return nil
	}

	return fmt.Errorf("token is invalid")
}