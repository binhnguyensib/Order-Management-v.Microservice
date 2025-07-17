package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

var Logger = logrus.New()

type User struct {
	User_id   bson.ObjectID `json:"uid" bson:"_id,omitempty"`
	Name      string
	Email     string    `json:"email"`
	Password  string    `json:"-" bson:"password"`
	Role      string    `json:"role" bson:"role"`
	Status    string    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdateAt  time.Time `json:"updated_at" bson:"updated_at"`
}

// Hash password before store in database
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		Logger.Error("Fail when hash password")
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Check password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	User        *User  `json:"user"`
}

type TokenClaims struct {
	User_id string `json:"uid"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}
