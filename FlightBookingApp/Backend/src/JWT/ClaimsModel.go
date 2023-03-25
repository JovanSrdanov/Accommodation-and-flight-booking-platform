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
	Password string `json:"password,omitempty"`
	Roles    []model.Role  `json:"roles,omitempty"`
	jwt.StandardClaims
}

//TODO Stefan: add an env file
const ip = "192.168.0.107"	//issuer

func (claims JwtClaims) Valid() error {
	var now = time.Now().UTC().Unix()
	if claims.VerifyExpiresAt(now, true) && claims.VerifyIssuer(ip, true) {
		return nil
	}

	return fmt.Errorf("token is invalid")
}