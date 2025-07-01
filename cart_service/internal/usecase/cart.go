package usecase

import (
	"cart_service/internal/domain"
	"cart_service/internal/usecase/helper"
	"cart_service/proto/product"
	"context"

	"google.golang.org/grpc"
)

var _ domain.CartUsecase = (*cartUsecaseImpl)(nil)

type cartUsecaseImpl struct {
	cartRepo      domain.CartRepository
	ProductClient product.ProductServiceClient
	Mapper        helper.CartMapper
}

func NewCartUsecase(cartRepo domain.CartRepository, conn *grpc.ClientConn) domain.CartUsecase {

	return &cartUsecaseImpl{
		cartRepo:      cartRepo,
		ProductClient: product.NewProductServiceClient(conn),
		Mapper:        helper.NewCartMapper(product.NewProductServiceClient(conn)),
	}
}

func (cu *cartUsecaseImpl) AddToCart(ctx context.Context, customerID string, cartItemReq *domain.CartItemRequest) (*domain.Cart, error) {
	productProto, err := cu.ProductClient.GetById(ctx, &product.GetByIdRequest{Id: cartItemReq.ProductID})
	if err != nil {
		return nil, err
	}
	cartItem := &domain.CartItem{
		ProductID:    cartItemReq.ProductID,
		ProductName:  cartItemReq.ProductName,
		Quantity:     cartItemReq.Quantity,
		ProductPrice: productProto.Product.Price,
		Subtotal:     float64(cartItemReq.Quantity) * productProto.Product.Price,
	}

	cart, err := cu.cartRepo.AddToCart(ctx, customerID, cartItem)
	if err != nil {
		return nil, err
	}
	return cart, nil
}
