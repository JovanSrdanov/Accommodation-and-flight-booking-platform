package token

import (
	JWT "FlightBookingApp/JWT"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jWTPrivateToken = os.Getenv("SECRET_KEY")
)

func GenerateToken(claims *JWT.JwtClaims, expirationTime time.Time) (string, error) {
	claims.ExpiresAt = expirationTime.Unix()
	claims.IssuedAt = time.Now().UTC().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(jWTPrivateToken))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (bool, *JWT.JwtClaims){
	claims := &JWT.JwtClaims{}
	token, _ := getTokenFromString(tokenString, claims)

	if token.Valid {
		if e := claims.Valid() ; e == nil {
			return true, claims
		}
	}

	return false, claims
}

func getTokenFromString(tokenString string, claims *JWT.JwtClaims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(jWTPrivateToken), nil
	})
}