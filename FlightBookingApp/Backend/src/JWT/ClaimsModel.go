package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtClaims struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Roles    []int  `json:"roles,omitempty"`
	jwt.StandardClaims
}

//TODO Stefan temp, remove later
const ip = "192.168.0.107"

func (claims JwtClaims) Valid() error {
	var now = time.Now().UTC().Unix()
	if claims.VerifyExpiresAt(now, true) && claims.VerifyIssuer(ip, true) {
		return nil
	}

	return fmt.Errorf("token is invalid")
}