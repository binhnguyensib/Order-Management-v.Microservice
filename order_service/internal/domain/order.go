package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	Id           bson.ObjectID `bson:"_id,omitempty"`
	CustomerID   string
	Items        []*OrderItem
	TotalItems   int
	TotalPrice   float64
	Status       OrderStatus
	ShippingInfo ShippingInfo
	CreateAt     time.Time
	UpdateAt     time.Time
}

type OrderItem struct {
	ProductID    string
	ProductName  string
	ProductPrice float64
	Quantity     int
	Subtotal     float64
}

type ShippingInfo struct {
	Name        string
	Phone       string
	City        string
	ShippingFee float64
}
