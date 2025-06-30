package domain

import (
	"context"
)

type CartUsecase interface {
	AddToCart(ctx context.Context, customerID string, cart *CartItemRequest) (*Cart, error)
	//GetCartByCustomerId(ctx context.Context, customerID string) (*Cart, error)
	//UpdateCartItem(ctx context.Context, customerID string, cartItem *CartItemRequest) (*Cart, error)
	//RemoveCartItem(ctx context.Context, customerID string, productID string) (*Cart, error)
	//ClearCart(ctx context.Context, customerID string) error
}

type CartRepository interface {
	AddToCart(ctx context.Context, customerID string, cartItem *CartItem) (*Cart, error)
	//GetCartByCustomerId(ctx context.Context, customerID string) (*Cart, error)
	//UpdateCartItem(ctx context.Context, customerID string, cartItem *CartItem) (*Cart, error)
	//RemoveCartItem(ctx context.Context, customerID string, productID string) (*Cart, error)
	//ClearCart(ctx context.Context, customerID string) error
}
