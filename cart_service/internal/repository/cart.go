package repository

import (
	rdis "cart_service/config"
	"cart_service/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	redisKey := fmt.Sprintf("cart:%s", CustomerId)
	collection := cr.conn.Collection("carts")
	cached, err := rdis.Get(ctx, redisKey)
	if err == nil && cached != "" {
		var c domain.Cart
		if ummarshalErr := json.Unmarshal([]byte(cached), &c); ummarshalErr != nil {
			Logger.Info("Cache hit for", redisKey)
		}
		found := false
		for i, existingItem := range c.Items {
			if existingItem.ProductID == cartItem.ProductID {
				c.Items[i].Quantity += cartItem.Quantity
				c.TotalItems += cartItem.Quantity
				c.TotalPrice += cartItem.ProductPrice * float64(cartItem.Quantity)
				found = true
				break
			}
		}
		if !found {
			c.Items = append(c.Items, cartItem)
		}
		cr.recalCartTotals(&c)
		_, err = collection.UpdateOne(ctx, bson.M{"_id": c.Id}, bson.M{"$set": c})
		if err != nil {
			return nil, err
		}
		Logger.WithFields(
			logrus.Fields{
				"method":     "AddToCart",
				"customerId": CustomerId}).Info("Add item to cart")
		rdis.Del(ctx, redisKey)
		return &c, nil
	} else {

		var existingCart domain.Cart
		err = collection.FindOne(ctx, bson.M{"customer_id": CustomerId}).Decode(&existingCart)
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
		bytes, marshalerr := json.Marshal(existingCart)
		if marshalerr == nil {
			go func() {
				ctxBg := context.Background()
				rdis.Set(ctxBg, redisKey, bytes, time.Hour)
			}()
		}
		Logger.WithFields(
			logrus.Fields{
				"method":     "AddToCart",
				"customerId": CustomerId}).Info("Add item to cart")
		return &existingCart, nil
	}

}

func (cr *cartRepositoryImpl) GetCartByCustomerId(ctx context.Context, customerId string) (*domain.Cart, error) {
	redisKey := fmt.Sprintf("cart:%s", customerId)
	cached, err := rdis.Get(ctx, redisKey)
	if err == nil && cached != "" {
		var c domain.Cart
		if ummarshalErr := json.Unmarshal([]byte(cached), &c); ummarshalErr != nil {
			Logger.Info("Cache hit for", redisKey)
			return &c, nil
		}
	}
	collection := cr.conn.Collection("carts")
	var cart domain.Cart
	err = collection.FindOne(ctx, bson.M{"customer_id": customerId}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			Logger.Errorf("Cart of %v is empty", customerId)
			return nil, err
		}
		Logger.Errorf("Failed to find cart for customer %v", customerId)
		return nil, err
	}

	bytes, marshalerr := json.Marshal(cart)
	if marshalerr == nil {
		go func() {
			ctxBg := context.Background()
			rdis.Set(ctxBg, redisKey, bytes, time.Hour)
		}()
	}
	Logger.WithFields(
		logrus.Fields{
			"method":     "GetCartByCustomerId",
			"customerId": customerId}).Info("Fetch cart form DB")
	return &cart, nil
}

func (cr *cartRepositoryImpl) UpdateCartItem(ctx context.Context, customerID string, cartItem *domain.CartItem) (*domain.Cart, error) {
	redisKey := fmt.Sprintf("cart:%s", customerID)
	collection := cr.conn.Collection("carts")
	cached, err := rdis.Get(ctx, redisKey)
	if err == nil && cached != "" {
		var c domain.Cart
		if ummarshalErr := json.Unmarshal([]byte(cached), &c); ummarshalErr != nil {
			Logger.Info("Cache hit for", redisKey)
		}
		found := false
		for i, existingItem := range c.Items {
			if existingItem.ProductID == cartItem.ProductID {
				c.Items[i].Quantity = cartItem.Quantity
				c.Items[i].Subtotal = cartItem.ProductPrice * float64(cartItem.Quantity)
				found = true
				break
			}
		}
		if !found {
			Logger.Errorf("Item not found in cart")
			return nil, err
		}

		cr.recalCartTotals(&c)
		_, err = collection.UpdateOne(ctx, bson.M{"_id": c.Id}, bson.M{"$set": c})
		if err != nil {
			Logger.Errorf("Failed to update cart")
			return nil, err
		}
		Logger.WithFields(
			logrus.Fields{
				"method":     "UpdateCartItem",
				"customerId": customerID}).Info("Update cart item successfully")
		rdis.Del(ctx, redisKey)
		return &c, nil
	} else {

		var existingCart domain.Cart
		err = collection.FindOne(ctx, bson.M{"customer_id": customerID}).Decode(&existingCart)
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
		bytes, marshalerr := json.Marshal(existingCart)
		if marshalerr == nil {
			go func() {
				ctxBg := context.Background()
				rdis.Set(ctxBg, redisKey, bytes, time.Hour)
			}()
		}
		Logger.WithFields(
			logrus.Fields{
				"method":     "UpdateCartItem",
				"customerId": customerID}).Info("Update cart item successfully")
		return &existingCart, nil
	}

}

func (cr *cartRepositoryImpl) RemoveCartItem(ctx context.Context, customerID string, productID string) (*domain.Cart, error) {
	redisKey := fmt.Sprintf("cart:%s", customerID)
	cached, err := rdis.Get(ctx, redisKey)
	if err == nil && cached != "" {
		var c domain.Cart
		if ummarshalErr := json.Unmarshal([]byte(cached), &c); ummarshalErr != nil {
			Logger.Info("Cache hit for", redisKey)
		}
		var updatedItems []*domain.CartItem
		for _, item := range c.Items {
			if item.ProductID != productID {
				updatedItems = append(updatedItems, item)
			}
		}

		if len(updatedItems) == len(c.Items) {
			Logger.Error("Item not found in cart")
			return nil, err
		}

		c.Items = updatedItems
		cr.recalCartTotals(&c)

		_, err = cr.conn.Collection("carts").UpdateOne(ctx, bson.M{"_id": c.Id}, bson.M{"$set": c})
		if err != nil {
			Logger.Error("Failed to update cart after removing item")
			return nil, err
		}
		Logger.WithFields(
			logrus.Fields{
				"method":     "RemoveCartItem",
				"customerId": customerID}).Info("Remove cart item successfully")
		rdis.Del(ctx, redisKey)
		return &c, nil
	} else {
		collection := cr.conn.Collection("carts")
		var existingCart domain.Cart
		err = collection.FindOne(ctx, bson.M{"customer_id": customerID}).Decode(&existingCart)
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
		bytes, marshalerr := json.Marshal(existingCart)
		if marshalerr == nil {
			go func() {
				ctxBg := context.Background()
				rdis.Set(ctxBg, redisKey, bytes, time.Hour)
			}()
		}
		Logger.WithFields(
			logrus.Fields{
				"method":     "RemoveCartItem",
				"customerId": customerID}).Info("Remove cart item successfully")
		return &existingCart, nil
	}

}

func (cr *cartRepositoryImpl) ClearCart(ctx context.Context, customerID string) error {
	redisKey := fmt.Sprintf("cart:%s", customerID)
	cached, err := rdis.Get(ctx, redisKey)
	if err == nil && cached != "" {
		var c domain.Cart
		if ummarshalErr := json.Unmarshal([]byte(cached), &c); ummarshalErr != nil {
			Logger.Info("Cache hit for", redisKey)
		}
		result, err := cr.conn.Collection("carts").DeleteOne(ctx, bson.M{"customer_id": customerID})
		if err != nil {
			Logger.Error("Failed to clear cart for customer", "customer_id", customerID, "error", err)
			return err
		}
		if result.DeletedCount == 0 {
			Logger.Warn("Customer's cart is already empty", "customer_id", customerID)
		}
		Logger.WithFields(
			logrus.Fields{
				"method":     "ClearCart",
				"customerId": customerID}).Info("Clear cart successfully")
		return nil
	} else {
		collection := cr.conn.Collection("carts")
		result, err := collection.DeleteOne(ctx, bson.M{"customer_id": customerID})
		if err != nil {
			Logger.Error("Failed to clear cart for customer", "customer_id", customerID, "error", err)
			return err
		}
		if result.DeletedCount == 0 {
			Logger.Warn("Customer's cart is already empty", "customer_id", customerID)
		}
		Logger.WithFields(
			logrus.Fields{
				"method":     "ClearCart",
				"customerId": customerID}).Info("Clear cart successfully")
		return nil
	}

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
