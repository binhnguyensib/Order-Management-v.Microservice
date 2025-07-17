package app

import (
	"context"
	"customer_service/config"
	_ "customer_service/docs"
	"customer_service/grpc"
	"customer_service/internal/consumer"
	"customer_service/internal/handler"
	"customer_service/internal/repository"
	"customer_service/internal/usecase"
	cs "customer_service/proto"
	"net"
	"os"
	"os/signal"
	"syscall"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	rpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// @title Order Management API
// @version Microservice
// @description This is a sample server for managing orders, customers, products, and carts.
// @host localhost:8081
// @BasePath /api

var Logger = logrus.New()

func Run() {
	Logger.SetFormatter(&logrus.TextFormatter{})
	db, err := config.ConnectMongoDB()
	if err != nil {
		Logger.Error("Error when connecting to MongoDB")
	}
	defer db.Close()

	config.InitRedis()

	rabbitConn, err := config.InitRabbitMQ()
	if err != nil {
		Logger.Errorf("Failed to connect to RabbitMQ: %v", err)
	} else {
		defer rabbitConn.Close()
	}

	customerRepo := repository.NewCustomerRepository(db.DB)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	if rabbitConn != nil {
		consumer := consumer.NewCustomerConsumer(rabbitConn, customerUsecase)
		if err := consumer.Setup(); err != nil {
			Logger.Errorf("Failed to setup consumer: %v", err)
		} else {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Start consumer in background
			go func() {
				if err := consumer.Start(ctx); err != nil {
					Logger.Errorf("Consumer failed: %v", err)
				}
			}()

			// Handle graceful shutdown
			go func() {
				quit := make(chan os.Signal, 1)
				signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
				<-quit
				Logger.Info("Shutting down consumer...")
				cancel()
			}()
		}
	}
	grpcPort := os.Getenv("GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		Logger.Error("Error when starting gRPC server")
	}
	grpcServer := rpc.NewServer()
	customerGRPCHandler := grpc.NewCustomerGRPCHandler(customerUsecase)
	cs.RegisterCustomerServiceServer(grpcServer, customerGRPCHandler)
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
			"message": "Customer-Service is running",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	admin := router.Group("/api/admin")
	{
		admin.GET("/customers", customerHandler.GetAll)
		admin.POST("/customers", customerHandler.Create)
		admin.DELETE("/customers/:id", customerHandler.Delete)

	}

	user := router.Group("/api/user")
	{
		user.GET("/customers/:id", customerHandler.GetByID)
		user.PUT("/customers/:id", customerHandler.Update)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	router.Run(":" + port)
}
