package usecase

import (
	"context"
	"customer_service/internal/domain"
)

var _ domain.CustomerUsecase = (*customerUsecaseImpl)(nil)

type customerUsecaseImpl struct {
	customerRepo domain.CustomerRepository
}

func NewCustomerUsecase(customerRepo domain.CustomerRepository) domain.CustomerUsecase {
	return &customerUsecaseImpl{
		customerRepo: customerRepo,
	}
}
func (cu *customerUsecaseImpl) GetAll(ctx context.Context) ([]*domain.Customer, error) {
	customers, err := cu.customerRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (cu *customerUsecaseImpl) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	customer, err := cu.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (cu *customerUsecaseImpl) Create(ctx context.Context, customer *domain.CustomerRequest) (*domain.Customer, error) {
	cus, err := cu.customerRepo.Create(ctx, customer)
	if err != nil {
		return nil, err
	}
	return cus, nil
}

func (cu *customerUsecaseImpl) Update(ctx context.Context, id string, customerReq *domain.CustomerRequest) (*domain.Customer, error) {
	cus, err := cu.customerRepo.Update(ctx, id, customerReq)
	if err != nil {
		return nil, err
	}
	return cus, nil
}

func (cu *customerUsecaseImpl) Delete(ctx context.Context, id string) (*domain.Customer, error) {
	cus, err := cu.customerRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return cus, nil
}
