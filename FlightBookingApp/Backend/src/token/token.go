package token

import (
	JWT "FlightBookingApp/JWT"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//TODO Stefan: move to a separate file
const (
	jWTPrivateTOken = "SecretTokenSecretToken"
	ip = "192.168.0.107"
)

func GenerateToken(claims *JWT.JwtClaims, expirationTime time.Time) (string, error) {
	claims.ExpiresAt = expirationTime.Unix()
	claims.IssuedAt = time.Now().UTC().Unix()
	claims.Issuer = ip

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(jWTPrivateTOken))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}