package domain

import "context"

type ProductUsecase interface {
	GetAll(ctx context.Context) ([]*Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, product *ProductRequest) (*Product, error)
	Update(ctx context.Context, id string, productReq *ProductRequest) (*Product, error)
	Delete(ctx context.Context, id string) (*Product, error)
}

type ProductRepository interface {
	GetAll(ctx context.Context) ([]*Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, product *ProductRequest) (*Product, error)
	Update(ctx context.Context, id string, productReq *ProductRequest) (*Product, error)
	Delete(ctx context.Context, id string) (*Product, error)
}
