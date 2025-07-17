package domain

import "context"

type CustomerUsecase interface {
	GetAll(ctx context.Context) ([]*Customer, error)
	GetByID(ctx context.Context, id string) (*Customer, error)
	Create(ctx context.Context, customer *CustomerRequest) (*Customer, error)
	Update(ctx context.Context, id string, customerReq *CustomerRequest) (*Customer, error)
	Delete(ctx context.Context, id string) (*Customer, error)
}

type CustomerRepository interface {
	GetAll(ctx context.Context) ([]*Customer, error)
	GetByID(ctx context.Context, id string) (*Customer, error)
	Create(ctx context.Context, customer *CustomerRequest) (*Customer, error)
	Update(ctx context.Context, id string, customerReq *CustomerRequest) (*Customer, error)
	Delete(ctx context.Context, id string) (*Customer, error)
}
