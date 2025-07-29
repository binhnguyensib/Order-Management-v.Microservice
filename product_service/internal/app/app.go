package app

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"product_service/config"
	_ "product_service/docs"
	"product_service/grpc"
	"product_service/internal/handler"
	"product_service/internal/repository"
	"product_service/internal/usecase"
	cronservice "product_service/internal/worker"
	pb "product_service/proto"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	rpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// @title Order Management API
// @version Microservice
// @description This is a sample server for managing orders, customers, products, and carts.
// @host localhost:8082
// @BasePath /api
var (
	Logger          = logrus.New()
	CronHealthCheck = cron.New()
)

func simpleHealthCheck() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	url := fmt.Sprintf("http://localhost:%s/health", port)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	Logger.Info("Performing health check...")
	start := time.Now()

	resp, err := client.Get(url)
	duration := time.Since(start)

	if err != nil {
		Logger.WithFields(logrus.Fields{
			"url":      url,
			"error":    err.Error(),
			"duration": duration,
		}).Error("Health check FAILED - Server not responding")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		Logger.WithFields(logrus.Fields{
			"url":         url,
			"status_code": resp.StatusCode,
			"duration":    duration,
		}).Info("Health check PASSED - Server is alive")
	} else {
		Logger.WithFields(logrus.Fields{
			"url":         url,
			"status_code": resp.StatusCode,
			"duration":    duration,
		}).Warn("Health check WARNING - Server responded with non-200 status")
	}
}

func Run() {
	Logger.SetFormatter(&logrus.TextFormatter{})
	db, err := config.ConnectMongoDB()
	if err != nil {
		Logger.Error("Error when connecting to MongoDB")
	}
	defer db.Close()

	config.InitRedis()

	productRepo := repository.NewProductRepository(db.DB)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUsecase)
	productCronService := cronservice.NewProductCronService(productUsecase)
	productCronService.Start()
	grpcPort := os.Getenv("GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+grpcPort)

	if err != nil {
		Logger.Error("Error when starting gRPC server")
	}
	grpcServer := rpc.NewServer()
	productGRPCHandler := grpc.NewProductGRPCHandler(productUsecase)
	pb.RegisterProductServiceServer(grpcServer, productGRPCHandler)
	reflection.Register(grpcServer)
	go func() {
		Logger.Infof("gRPC server start at port :%s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			Logger.Errorf("Failed to serve gRPC: %v", err)
		}
	}()

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Product-Service is running",
		})
	})

	CronHealthCheck.AddFunc("@every 15m", simpleHealthCheck)
	CronHealthCheck.Start()
	defer CronHealthCheck.Stop()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	admin := router.Group("/api/admin")
	{
		admin.POST("/products", productHandler.Create)
		admin.PUT("/products/:id", productHandler.Update)
		admin.DELETE("/products/:id", productHandler.Delete)
		admin.POST("/products/supplier-upload", productHandler.SupplierUpload)
	}
	user := router.Group("/api/user")
	{
		user.GET("/products", productHandler.GetAll)
		user.GET("/products/:id", productHandler.GetByID)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	router.Run(":" + port)
}
