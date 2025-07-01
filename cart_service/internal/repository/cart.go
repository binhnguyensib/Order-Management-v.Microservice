package repository

import (
	"cart_service/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type cartRepositoryImpl struct {
	conn *mongo.Database
}

func NewCartRepository(db *mongo.Database) domain.CartRepository {
	return &cartRepositoryImpl{
		conn: db,
	}
}

func (cr *cartRepositoryImpl) AddToCart(ctx context.Context, CustomerId string, cartItem *domain.CartItem) (*domain.Cart, error) {
	collection := cr.conn.Collection("cart_service")
	var existingCart domain.Cart
	err := collection.FindOne(ctx, bson.M{"customer_id": CustomerId}).Decode(&existingCart)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			newCart := &domain.Cart{
				CustomerID: CustomerId,
				Items:      []*domain.CartItem{cartItem},
				TotalItems: cartItem.Quantity,
				TotalPrice: cartItem.ProductPrice * float64(cartItem.Quantity),
			}
			result, err := collection.InsertOne(ctx, newCart)
			if err != nil {
				return nil, err
			}
			if insertedId, ok := result.InsertedID.(bson.ObjectID); ok {
				newCart.Id = insertedId
			}
			return newCart, nil
		}
	}
	found := false
	for i, existingItem := range existingCart.Items {
		if existingItem.ProductID == cartItem.ProductID {
			existingCart.Items[i].Quantity += cartItem.Quantity
			existingCart.TotalItems += cartItem.Quantity
			existingCart.TotalPrice += cartItem.ProductPrice * float64(cartItem.Quantity)
			found = true
			break
		}
	}
	if !found {
		existingCart.Items = append(existingCart.Items, cartItem)
	}
	cr.recalCartTotals(&existingCart)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": existingCart.Id}, bson.M{"$set": existingCart})
	if err != nil {

		return nil, err
	}
	return &existingCart, nil
}

func (cr *cartRepositoryImpl) recalCartTotals(cart *domain.Cart) {
	totalItems := 0
	totalPrice := 0.0
	for _, item := range cart.Items {
		totalItems += item.Quantity
		totalPrice += item.ProductPrice * float64(item.Quantity)
	}
	cart.TotalItems = totalItems
	cart.TotalPrice = totalPrice
}
