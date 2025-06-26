package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type Customer struct {
	Id       bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string        `bson:"name" json:"name"`
	Email    string        `bson:"email" json:"email"`
	Phone    string        `bson:"phone" json:"phone"`
	Password string        `bson:"password" json:"-"`
}

type CustomerRegister struct {
	Name     string `bson:"name" json:"name"`
	Email    string `bson:"email" json:"email"`
	Phone    string `bson:"phone" json:"phone"`
	Password string `bson:"password" json:"-"`
}

type CustomerLogin struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"-"`
}

func (c *Customer) HashPassword() bool {
	hash, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return false
	}
	c.Password = string(hash)
	return true
}

func (c *Customer) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password))
	return err == nil
}
