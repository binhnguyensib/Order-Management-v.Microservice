package cronservice

import (
	"context"
	"product_service/internal/domain"
	"time"

	"path/filepath"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type productCronService struct {
	productUsecase domain.ProductUsecase
	cron           *cron.Cron
	supplierDir    string
	logger         *logrus.Logger
}

func NewProductCronService(productUsecase domain.ProductUsecase) *productCronService {
	return &productCronService{
		productUsecase: productUsecase,
		cron:           cron.New(),
		supplierDir:    "./supplier",
		logger:         logrus.New(),
	}
}

func (pcs *productCronService) Start() {
	// Schedule the cron job to update new products price form supplier every day at 21:00
	pcs.cron.AddFunc("0 21 * * *", func() {
		now := time.Now()
		formattedDate := now.Format("20060102")
		result := "price_update_" + formattedDate + ".xlsx"
		var report *domain.BulkPriceUpdateResult
		var ctx = context.Background()
		report, err := pcs.productUsecase.UpdatePricesFromExcel(ctx, result)
		if err != nil {
			pcs.logger.WithError(err).Error("Failed to update prices from Excel")
			return
		}
		pcs.logger.Infof("Successfully updated prices from Excel: %+v", report)
	})

	//test 30s
	pcs.cron.AddFunc("@every 30s", func() {
		now := time.Now()
		formattedDate := now.Format("20060102")
		filename := "price_update_" + formattedDate + ".xlsx"
		filepath := filepath.Join(pcs.supplierDir, filename)
		var report *domain.BulkPriceUpdateResult
		var ctx = context.Background()
		report, err := pcs.productUsecase.UpdatePricesFromExcel(ctx, filepath)
		if err != nil {
			pcs.logger.WithError(err).Error("Failed to update prices from Excel")
			return
		}
		pcs.logger.Infof("Successfully updated prices from Excel: %+v", report)
	})

	pcs.cron.Start()
	pcs.logger.Info("Product cron service started")
}
