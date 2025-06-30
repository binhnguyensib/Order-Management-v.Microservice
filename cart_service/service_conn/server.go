package serviceconn

import (
	"cart_service/internal/domain"
	pb "cart_service/proto/cart"
	"cart_service/service_conn/helper"
	"context"
)

type cartServiceServer struct {
	pb.UnimplementedCartServiceServer
	cartUsecase domain.CartUsecase
	cartMapper  helper.CartMapper
}

func (css *cartServiceServer) AddToCart(ctx context.Context, req *pb.AddToCartRequest) (*pb.AddToCartResponse, error) {
	cartItemRequest := css.cartMapper.ToCartItemRequest(req.CartItemRequest)

	cart, err := css.cartUsecase.AddToCart(ctx, req.CustomerId, cartItemRequest)
	if err != nil {
		return nil, err
	}

	return &pb.AddToCartResponse{Cart: css.cartMapper.ToProtoCart(cart)}, nil
}
