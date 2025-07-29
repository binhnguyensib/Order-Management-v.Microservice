package domain

import "context"

type ProductUsecase interface {
	GetAll(ctx context.Context) ([]*Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, product *ProductRequest) (*Product, error)
	Update(ctx context.Context, id string, productReq *ProductRequest) (*Product, error)
	Delete(ctx context.Context, id string) (*Product, error)

	UpdatePricesFromExcel(ctx context.Context, filePath string) (*BulkPriceUpdateResult, error)
}

type ProductRepository interface {
	GetAll(ctx context.Context) ([]*Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)

	GetAndUpdateByName(ctx context.Context, name string, newPrice float64) (*PriceUpdateResult, error)

	Create(ctx context.Context, product *ProductRequest) (*Product, error)

	Update(ctx context.Context, id string, productReq *ProductRequest) (*Product, error)

	Delete(ctx context.Context, id string) (*Product, error)
}
