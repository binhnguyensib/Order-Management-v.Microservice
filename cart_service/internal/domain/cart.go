package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Cart struct {
	Id         bson.ObjectID `json:"id" bson:"_id,omitempty"`
	CustomerID string        `json:"customer_id" bson:"customer_id"`
	Items      []*CartItem   `json:"items" bson:"items"`
	TotalItems int           `json:"total_items" bson:"total_items"`
	TotalPrice float64       `json:"total_price" bson:"total_price"`
}

type CartItem struct {
	ProductID    string  `json:"product_id" bson:"product_id"`
	ProductName  string  `json:"product_name" bson:"product_name"`
	ProductPrice float64 `json:"product_price" bson:"product_price"`
	Quantity     int     `json:"quantity" bson:"quantity"`
	Subtotal     float64 `json:"subtotal" bson:"subtotal"`
}

type CartItemRequest struct {
	ProductID   string `json:"product_id" bson:"product_id"`
	ProductName string `json:"product_name" bson:"product_name"`
	Quantity    int    `json:"quantity" bson:"quantity"`
}
