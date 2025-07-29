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

type PriceUpdateResult struct {
	ProductName string  `json:"product_name"`
	NewPrice    float64 `json:"new_price"`
	OldPrice    float64 `json:"old_price"`
	Status      string  `json:"status"`
	Message     string  `json:"message"`
}

type BulkPriceUpdateResult struct {
	TotalProcessed int                  `json:"total_processed"`
	TotalSuccess   int                  `json:"total_success"`
	TotalFailed    int                  `json:"total_failed"`
	Results        []*PriceUpdateResult `json:"results"`
}
