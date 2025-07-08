package repository

import (
	"cart_service/internal/domain"
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var Logger = logrus.New()

type cartRepositoryImpl struct {
	conn *mongo.Database
}

func NewCartRepository(db *mongo.Database) domain.CartRepository {
	return &cartRepositoryImpl{
		conn: db,
	}
}

func (cr *cartRepositoryImpl) AddToCart(ctx context.Context, CustomerId string, cartItem *domain.CartItem) (*domain.Cart, error) {
	collection := cr.conn.Collection("carts")
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

func (cr *cartRepositoryImpl) GetCartByCustomerId(ctx context.Context, customerId string) (*domain.Cart, error) {
	collection := cr.conn.Collection("carts")
	var cart domain.Cart
	err := collection.FindOne(ctx, bson.M{"customer_id": customerId}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			Logger.Errorf("Cart of %v is empty", customerId)
			return nil, err
		}
		Logger.Errorf("Failed to find cart for customer %v", customerId)
		return nil, err
	}
	return &cart, nil
}

func (cr *cartRepositoryImpl) UpdateCartItem(ctx context.Context, customerID string, cartItem *domain.CartItem) (*domain.Cart, error) {
	collection := cr.conn.Collection("carts")
	var existingCart domain.Cart
	err := collection.FindOne(ctx, bson.M{"customer_id": customerID}).Decode(&existingCart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			Logger.Errorf("Cart of %v is empty", customerID)
			return nil, err
		}
		Logger.Errorf("Failed to find cart for customer %v", customerID)
		return nil, err
	}

	found := false
	for i, existingItem := range existingCart.Items {
		if existingItem.ProductID == cartItem.ProductID {
			existingCart.Items[i].Quantity = cartItem.Quantity
			existingCart.Items[i].Subtotal = cartItem.ProductPrice * float64(cartItem.Quantity)
			found = true
			break
		}
	}
	if !found {
		Logger.Errorf("Item not found in cart")
		return nil, err
	}

	cr.recalCartTotals(&existingCart)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": existingCart.Id}, bson.M{"$set": existingCart})
	if err != nil {
		Logger.Errorf("Failed to update cart")
		return nil, err
	}
	return &existingCart, nil
}

func (cr *cartRepositoryImpl) RemoveCartItem(ctx context.Context, customerID string, productID string) (*domain.Cart, error) {
	collection := cr.conn.Collection("carts")
	var existingCart domain.Cart
	err := collection.FindOne(ctx, bson.M{"customer_id": customerID}).Decode(&existingCart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			Logger.Errorf("Cart of %v is empty", customerID)
			return nil, err
		}
		Logger.Errorf("Failed to find cart for customer %v", customerID)
		return nil, err
	}

	var updatedItems []*domain.CartItem
	for _, item := range existingCart.Items {
		if item.ProductID != productID {
			updatedItems = append(updatedItems, item)
		}
	}

	if len(updatedItems) == len(existingCart.Items) {
		Logger.Error("Item not found in cart")
		return nil, err
	}

	existingCart.Items = updatedItems
	cr.recalCartTotals(&existingCart)

	_, err = collection.UpdateOne(ctx, bson.M{"_id": existingCart.Id}, bson.M{"$set": existingCart})
	if err != nil {
		Logger.Error("Failed to update cart after removing item")
		return nil, err
	}
	return &existingCart, nil
}

func (cr *cartRepositoryImpl) ClearCart(ctx context.Context, customerID string) error {
	collection := cr.conn.Collection("carts")
	result, err := collection.DeleteOne(ctx, bson.M{"customer_id": customerID})
	if err != nil {
		Logger.Error("Failed to clear cart for customer", "customer_id", customerID, "error", err)
		return err
	}
	if result.DeletedCount == 0 {
		Logger.Warn("Customer's cart is already empty", "customer_id", customerID)
	}
	return nil
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
