package grpc

import (
	"context"
	"customer_service/internal/domain"
	cs "customer_service/proto"
)

type customerGRPCHandler struct {
	cs.UnimplementedCustomerServiceServer
	customerUsecase domain.CustomerUsecase
}

func NewCustomerGRPCHandler(usecase domain.CustomerUsecase) *customerGRPCHandler {
	return &customerGRPCHandler{
		customerUsecase: usecase,
	}
}

func (ch *customerGRPCHandler) GetAll(ctx context.Context, req *cs.GetAllRequest) (*cs.GetAllResponse, error) {
	customers, err := ch.customerUsecase.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var csCustomers []*cs.Customer
	for _, c := range customers {
		csCustomers = append(csCustomers, &cs.Customer{
			Id:    c.Id.Hex(),
			Name:  c.Name,
			Email: c.Email,
			Phone: c.Phone,
		})
	}

	return &cs.GetAllResponse{Customers: csCustomers}, nil
}

func (ch *customerGRPCHandler) GetById(ctx context.Context, req *cs.GetByIdRequest) (*cs.GetByIdResponse, error) {
	customer, err := ch.customerUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &cs.GetByIdResponse{
		Customer: &cs.Customer{
			Id:    customer.Id.Hex(),
			Name:  customer.Name,
			Email: customer.Email,
			Phone: customer.Phone,
		},
	}, nil
}
