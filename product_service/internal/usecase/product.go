package usecase

import (
	"context"
	"fmt"
	"product_service/internal/domain"
	"strconv"

	"github.com/xuri/excelize/v2"
)

var _ domain.ProductUsecase = (*productUsecaseImpl)(nil)

type productUsecaseImpl struct {
	productRepo domain.ProductRepository
}

func NewProductUsecase(productRepo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecaseImpl{
		productRepo: productRepo,
	}
}

func (pu *productUsecaseImpl) GetAll(ctx context.Context) ([]*domain.Product, error) {
	products, err := pu.productRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil

}

func (pu *productUsecaseImpl) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	product, err := pu.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (pu *productUsecaseImpl) Create(ctx context.Context, product *domain.ProductRequest) (*domain.Product, error) {
	productCreated, err := pu.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}
	return productCreated, nil
}
func (pu *productUsecaseImpl) Update(ctx context.Context, id string, productReq *domain.ProductRequest) (*domain.Product, error) {
	productUpdated, err := pu.productRepo.Update(ctx, id, productReq)
	if err != nil {
		return nil, err
	}
	return productUpdated, nil
}

func (pu *productUsecaseImpl) Delete(ctx context.Context, id string) (*domain.Product, error) {
	productDeleted, err := pu.productRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return productDeleted, nil
}

func (pu *productUsecaseImpl) UpdatePricesFromExcel(ctx context.Context, filePath string) (*domain.BulkPriceUpdateResult, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("no sheet found in excel file")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from sheet: %w", err)
	}

	response := &domain.BulkPriceUpdateResult{
		Results: make([]*domain.PriceUpdateResult, 0),
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			continue
		}
		productName := row[0]
		newPrice, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			continue
		}

		result, err := pu.productRepo.GetAndUpdateByName(ctx, productName, newPrice)
		if err != nil {
			response.Results = append(response.Results, &domain.PriceUpdateResult{
				ProductName: productName,
				NewPrice:    newPrice,
				Status:      "error",
				Message:     fmt.Sprintf("failed to update price: %v", err),
			})
			continue
		}

		response.Results = append(response.Results, result)
		response.TotalProcessed++
		if result.Status == "success" {
			response.TotalSuccess++
		} else {
			response.TotalFailed++
		}
	}

	return response, nil
}
