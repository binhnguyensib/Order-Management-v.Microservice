package usecase

import (
	"cart_service/internal/domain"
	"context"
)

var _ domain.CartUsecase = (*cartUsecaseImpl)(nil)

type cartUsecaseImpl struct {
	cartRepo domain.CartRepository
}

func NewCartUsecase(cartRepo domain.CartRepository) domain.CartUsecase {
	return &cartUsecaseImpl{
		cartRepo: cartRepo,
	}
}

func (cu *cartUsecaseImpl) AddToCart(ctx context.Context, customerID string, cart *domain.CartItemRequest) (*domain.Cart, error) {

}
