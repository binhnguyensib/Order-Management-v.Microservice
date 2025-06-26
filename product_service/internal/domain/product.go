package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type Product struct {
	Id    bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  string        `json:"name"`
	Price float64       `json:"price"`
	Stock int           `json:"stock"`
}

type ProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}
