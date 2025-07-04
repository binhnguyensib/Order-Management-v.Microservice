package utils

import (
	"api_gateway/internal/domain"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(email, uid string) (string, error) {
	expirationTime := time.Now().Add(time.Hour)
	tokenClaims := &domain.TokenClaims{
		Email:   email,
		User_id: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Authentication service",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseJWT(tokenString string) (*domain.TokenClaims, error) {
	tokenClaims := &domain.TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, tokenClaims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return tokenClaims, nil
}

func ValidateJWT(tokenString string) (string, error) {
	tokenClaims, err := ParseJWT(tokenString)
	if err != nil {
		return "", err
	}
	if tokenClaims.ExpiresAt.Time.Before(time.Now()) {
		return "", jwt.ErrTokenExpired
	}
	return tokenClaims.Email, nil
}
