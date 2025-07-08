package app

import (
	"cart_service/config"
	_ "cart_service/docs"
	"cart_service/internal/handler"
	"cart_service/internal/repository"
	"cart_service/internal/usecase"
	"cart_service/internal/usecase/helper"
	pb "cart_service/proto/cart"
	"cart_service/proto/product"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// @title Order Management API
// @version Microservice
// @description This is a sample server for managing orders, customers, products, and carts.
// @host localhost:8083
// @BasePath /api

var Logger = logrus.New()

func Run() {
	db, err := config.ConnectMongoDB()
	if err != nil {
		Logger.Error("Error connecting to database")
	}
	defer db.Close()

	config.InitRedis()

	productClient, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		Logger.Error("Error when create new gRPC client to Product service", err)
	}

	cartRepo := repository.NewCartRepository(db.DB)
	cartUsecase := usecase.NewCartUsecase(cartRepo, productClient)
	cartHandler := handler.NewCartHandler(cartUsecase)
	cartMaper := helper.NewCartMapper(product.NewProductServiceClient(productClient))
	grpcPort := os.Getenv("GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		Logger.Error("Error when starting gRPC server")
	}
	grpcServer := grpc.NewServer()
	cartGRPCHandler := handler.NewCartGRPCHandler(cartUsecase, cartMaper)
	pb.RegisterCartServiceServer(grpcServer, cartGRPCHandler)
	reflection.Register(grpcServer)
	go func() {
		Logger.Infof("gRPC Server listening on :%s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			Logger.Errorf("Failed to serve gRPC: %v", err)
		}
	}()

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Cart is running",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api")
	{
		api.POST("/carts/:id", cartHandler.AddToCart)
		api.GET("/carts/:id", cartHandler.GetCartByCustomerId)
		api.DELETE("/carts/:id", cartHandler.ClearCart)
		api.PUT("/carts/:id", cartHandler.UpdateCartItem)
		api.DELETE("/carts/:id/item/:product_id", cartHandler.RemoveCartItem)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	router.Run(":" + port)

}
