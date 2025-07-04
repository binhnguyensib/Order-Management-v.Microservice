package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	Email    string        `json:"email"`
	Password string        `json:"-"`
	User_id  bson.ObjectID `json:"uid" bson:"_id,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	//RefreshToken string `json:"refresh_token"`
	//ExpiresAt    int64  `json:"expires_at"`
	User *User `json:"user"`
}

type TokenClaims struct {
	User_id string `json:"uid"`
	Email   string `json:"email"`
	jwt.RegisteredClaims
}
