package handler

import (
	"cart_service/internal/domain"
	"cart_service/internal/usecase/helper"
	pb "cart_service/proto/cart"
	"context"
)

type cartGRPCHandler struct {
	pb.UnimplementedCartServiceServer
	cartUsecase domain.CartUsecase
	cartMapper  helper.CartMapper
}

func NewCartGRPCHandler(cartUsecase domain.CartUsecase, cartMapper helper.CartMapper) *cartGRPCHandler {
	return &cartGRPCHandler{
		cartUsecase: cartUsecase,
		cartMapper:  cartMapper,
	}
}

func (css *cartGRPCHandler) AddToCart(ctx context.Context, req *pb.AddToCartRequest) (*pb.AddToCartResponse, error) {
	cartItemRequest := css.cartMapper.ToCartItemRequest(req.CartItemRequest)

	cart, err := css.cartUsecase.AddToCart(ctx, req.CustomerId, cartItemRequest)
	if err != nil {
		return nil, err
	}

	return &pb.AddToCartResponse{Cart: css.cartMapper.ToProtoCart(cart)}, nil
}
