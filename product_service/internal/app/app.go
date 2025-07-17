package app

import (
	"net"
	"os"
	"product_service/config"
	_ "product_service/docs"
	"product_service/grpc"
	"product_service/internal/handler"
	"product_service/internal/repository"
	"product_service/internal/usecase"
	pb "product_service/proto"

	"github.com/gin-gonic/gin"
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
var Logger = logrus.New()

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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	admin := router.Group("/api/admin")
	{
		admin.POST("/products", productHandler.Create)
		admin.PUT("/products/:id", productHandler.Update)
		admin.DELETE("/products/:id", productHandler.Delete)
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
