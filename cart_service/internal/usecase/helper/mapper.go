package helper

import (
	"cart_service/internal/domain"
	cs "cart_service/proto/cart"
	ps "cart_service/proto/product"
)

type CartMapper interface {
	ToCartItemRequest(*cs.CartItemRequest) *domain.CartItemRequest
	//ToProtoCartItem(*domain.CartItem) *cs.CartItem
	ToProtoCart(*domain.Cart) *cs.Cart
	//ToProduct(*ps.GetByIdRequest)
}

type cartMapper struct {
	productClient ps.ProductServiceClient
}

func NewCartMapper(productClient ps.ProductServiceClient) CartMapper {
	return &cartMapper{
		productClient: productClient,
	}
}

func (m *cartMapper) ToCartItemRequest(protoReq *cs.CartItemRequest) *domain.CartItemRequest {
	if protoReq == nil {
		return nil
	}
	return &domain.CartItemRequest{
		ProductID:   protoReq.ProductId,
		ProductName: protoReq.ProductName,
		Quantity:    int(protoReq.Quantity),
	}
}

func (m *cartMapper) ToProtoCart(cart *domain.Cart) *cs.Cart {
	if cart == nil {
		return nil
	}
	var cartItems []*cs.CartItem
	for _, c := range cart.Items {
		cartItems = append(cartItems, &cs.CartItem{
			ProductId:    c.ProductID,
			ProductName:  c.ProductName,
			ProductPrice: float32(c.ProductPrice),
			Quantity:     int32(c.Quantity),
			Subtotal:     float32(c.Subtotal),
		})
	}
	return &cs.Cart{
		Id:         cart.Id.Hex(),
		CustomerId: cart.CustomerID,
		CartItem:   cartItems,
		TotalItem:  int32(cart.TotalItems),
		TotalPrice: float32(cart.TotalPrice),
	}
}
