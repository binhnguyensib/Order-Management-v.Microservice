package grpc

import (
	"context"
	"product_service/internal/domain"
	ps "product_service/proto"
)

type productGRPCHandler struct {
	ps.UnimplementedProductServiceServer
	productUsecase domain.ProductUsecase
}

func NewProductGRPCHandler(usecase domain.ProductUsecase) *productGRPCHandler {
	return &productGRPCHandler{
		productUsecase: usecase,
	}
}

func (ph *productGRPCHandler) GetAll(ctx context.Context, req *ps.GetAllRequest) (*ps.GetAllResponse, error) {
	products, err := ph.productUsecase.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var pbProducts []*ps.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &ps.Product{
			Id:    p.Id.Hex(),
			Name:  p.Name,
			Price: p.Price,
			Stock: int32(p.Stock),
		})
	}

	return &ps.GetAllResponse{Products: pbProducts}, nil
}

func (ph *productGRPCHandler) GetById(ctx context.Context, req *ps.GetByIdRequest) (*ps.GetByIdResponse, error) {
	product, err := ph.productUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &ps.GetByIdResponse{
		Product: &ps.Product{
			Id:    product.Id.Hex(),
			Name:  product.Name,
			Price: product.Price,
			Stock: int32(product.Stock),
		},
	}, nil
}
