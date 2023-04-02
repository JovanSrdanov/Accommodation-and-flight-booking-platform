package token

import (
	JWT "FlightBookingApp/JWT"
	"FlightBookingApp/model"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	accessTokenSecret  = os.Getenv("ACCESS_SECRET_KEY")
	refreshTokenSecret = os.Getenv("REFRESH_SECRET_KEY")
)

func GenerateTokens(account model.Account) (string, string, error) {
	accessTokenString, err := GenerateAccessToken(account)

	if err != nil {
		return "", "", err
	}

	// creating a refresh token
	refreshTokenString, err := GenerateRefreshToken(account)

	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func GenerateAccessToken(account model.Account) (string, error) {
	accessTokenClaims := &JWT.AccessJwtClaims{}

	accessTokenClaims.ID = account.ID
	accessTokenClaims.Roles = []model.Role{account.Role}
	accessTokenClaims.TokenType = "access"

	// session token lasts for 15 minutes
	accessTokenClaims.ExpiresAt = time.Now().UTC().Add(time.Duration(15) * time.Minute).Unix()
	accessTokenClaims.IssuedAt = time.Now().UTC().Unix()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	// Sign and get the complete encoded token as a string using the secret
	accessTokenString, err := accessToken.SignedString([]byte(accessTokenSecret))

	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

func GenerateRefreshToken(account model.Account) (string, error) {
	refreshTokenClaims := &JWT.RefreshJwtClaims{}

	refreshTokenClaims.ID = account.ID
	refreshTokenClaims.Roles = []model.Role{account.Role}
	refreshTokenClaims.TokenType = "refresh"

	// refresh token lasts for a week
	refreshTokenClaims.ExpiresAt = time.Now().UTC().Add(time.Duration(24) * time.Hour * 7).Unix()
	refreshTokenClaims.IssuedAt = time.Now().UTC().Unix()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(refreshTokenSecret))

	if err != nil {
		return "", err
	}

	return refreshTokenString, nil
}

func VerifyToken(tokenString string) (error, *JWT.AccessJwtClaims) {
	claims := &JWT.AccessJwtClaims{}
	token, err := getTokenFromString(tokenString, claims)

	if err != nil {
		return err, claims
	}

	if token.Valid {
		if e := claims.Valid(); e == nil {
			return e, claims
		}
	}

	return nil, claims
}

func getTokenFromString(tokenString string, claims *JWT.AccessJwtClaims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		if claims.TokenType == "access" {
			return []byte(accessTokenSecret), nil
		}
		return []byte(refreshTokenSecret), nil
	})
}
