package token

import (
	JWT "FlightBookingApp/JWT"
	"FlightBookingApp/model"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// struct for storing metadata for access and refresh tokens
type TokenDetails struct {
	AccessToken string
	RefreshToken string
	AccessUuid string
	RefreshUuid string
	AccessTokenExpires int64
	RefreshTokenExpires int64
}

var (
	accessTokenSecret = os.Getenv("ACCESS_SECRET_KEY")
	refreshTokenSecret = os.Getenv("REFRESH_SECRET_KEY")
)

func GenerateToken(account model.Account) (string, string, error) {
	tokenDetails := &TokenDetails{}

	tokenDetails.AccessTokenExpires = time.Now().UTC().Add(time.Duration(15) * time.Minute).Unix()
	tokenDetails.AccessUuid = uuid.New().String()

	tokenDetails.RefreshTokenExpires = time.Now().UTC().Add(time.Duration(24) * time.Hour * 7).Unix()
	tokenDetails.RefreshUuid = uuid.New().String()

	accessTokenString, err := generateAccessToken(account, tokenDetails)

	if err != nil {
		return "", "", err
	}

	// creating a refresh token
	refreshTokenString, err := generateRefreshToken(account, tokenDetails)

	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString,  nil
}

func generateAccessToken(account model.Account, tokenDetails *TokenDetails) (string, error){
	accessTokenClaims := &JWT.AccessJwtClaims{}

	accessTokenClaims.ID = account.ID
	accessTokenClaims.Roles = []model.Role{account.Role}
	accessTokenClaims.AccessID = tokenDetails.AccessUuid

	// session token lasts for 15 minutes
	accessTokenClaims.ExpiresAt = tokenDetails.AccessTokenExpires
	accessTokenClaims.IssuedAt = time.Now().UTC().Unix()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	// Sign and get the complete encoded token as a string using the secret
	accessTokenString, err := accessToken.SignedString([]byte(accessTokenSecret))

	tokenDetails.AccessToken = accessTokenString

	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

func generateRefreshToken(account model.Account, tokenDetails *TokenDetails) (string, error) {
	refreshTokenClaims := &JWT.RefreshJwtClaims{}
	
	refreshTokenClaims.ID = account.ID
	refreshTokenClaims.Roles = []model.Role{account.Role}
	refreshTokenClaims.RefreshID = tokenDetails.RefreshUuid

	refreshTokenClaims.ExpiresAt = tokenDetails.RefreshTokenExpires
	refreshTokenClaims.IssuedAt = time.Now().UTC().Unix()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(refreshTokenSecret))

	tokenDetails.RefreshToken = refreshTokenString

	if err != nil {
		return "",  err
	}

	return refreshTokenString, nil
}

func VerifyToken(tokenString string) (bool, *JWT.AccessJwtClaims){
	claims := &JWT.AccessJwtClaims{}
	token, _ := getTokenFromString(tokenString, claims)

	if token.Valid {
		if e := claims.Valid() ; e == nil {
			return true, claims
		}
	}

	return false, claims
}

func getTokenFromString(tokenString string, claims *JWT.AccessJwtClaims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(accessTokenSecret), nil
	})
}