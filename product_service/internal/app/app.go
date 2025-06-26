package app

import (
	"log"
	"net"
	"os"
	"product_service/config"
	"product_service/grpc"
	"product_service/internal/handler"
	"product_service/internal/repository"
	"product_service/internal/usecase"
	pb "product_service/proto"

	"github.com/gin-gonic/gin"
	rpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run() {
	db, err := config.ConnectMongoDB()
	if err != nil {
		panic("Error connecting to database")
	}
	defer db.Close()

	productRepo := repository.NewProductRepository(db.DB)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUsecase)

	grpcPort := os.Getenv("GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+grpcPort)

	if err != nil {
		print("error when starting gRPC")
	}
	grpcServer := rpc.NewServer()
	productGRPCHandler := grpc.NewProductGRPCHandler(productUsecase)
	pb.RegisterProductServiceServer(grpcServer, productGRPCHandler)
	reflection.Register(grpcServer)
	go func() {
		log.Printf("gRPC Server listening on :%s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Product-Service is running",
		})
	})
	api := router.Group("/api")
	{
		api.GET("/products", productHandler.GetAll)
		api.GET("/products/:id", productHandler.GetByID)
		api.POST("/products", productHandler.Create)
		api.PUT("/products/:id", productHandler.Update)
		api.DELETE("/products/:id", productHandler.Delete)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	router.Run(":" + port)
}
