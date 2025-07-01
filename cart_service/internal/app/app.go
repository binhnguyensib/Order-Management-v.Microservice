package app

import (
	"cart_service/config"
	"cart_service/internal/handler"
	"cart_service/internal/repository"
	"cart_service/internal/usecase"
	"cart_service/internal/usecase/helper"
	pb "cart_service/proto/cart"
	"cart_service/proto/product"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func Run() {
	db, err := config.ConnectMongoDB()
	if err != nil {
		panic("Error connecting to database")
	}
	defer db.Close()

	productClient, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	cartRepo := repository.NewCartRepository(db.DB)
	cartUsecase := usecase.NewCartUsecase(cartRepo, productClient)
	cartHandler := handler.NewCartHandler(cartUsecase)
	cartMaper := helper.NewCartMapper(product.NewProductServiceClient(productClient))
	grpcPort := os.Getenv("GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+grpcPort)

	if err != nil {
		print("error when starting gRPC")
	}
	grpcServer := grpc.NewServer()
	cartGRPCHandler := handler.NewCartGRPCHandler(cartUsecase, cartMaper)
	pb.RegisterCartServiceServer(grpcServer, cartGRPCHandler)
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
			"message": "Cart is running",
		})
	})
	api := router.Group("/api")
	{
		api.POST("/carts/:id", cartHandler.AddToCart)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	router.Run(":" + port)

}
