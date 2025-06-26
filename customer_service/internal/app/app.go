package app

import (
	"customer_service/config"
	"customer_service/grpc"
	"customer_service/internal/repository"
	"net"
	"os"

	"customer_service/internal/handler"
	"customer_service/internal/usecase"
	cs "customer_service/proto"

	"github.com/gin-gonic/gin"
	rpc "google.golang.org/grpc"
)

func Run() {
	db, err := config.ConnectMongoDB()
	if err != nil {
		panic("Error connecting to database")
	}
	defer db.Close()

	customerRepo := repository.NewCustomerRepository(db.DB)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo)
	customerHandler := handler.NewCustomerHandler(customerUsecase)

	lis, err := net.Listen("tcp", ":"+string(os.Getenv("GRPC_PORT")))
	if err != nil {
		panic(err)
	}
	grpcServer := rpc.NewServer()
	customerGRPCHandler := grpc.NewCustomerGRPCHandler(customerUsecase)
	cs.RegisterCustomerServiceServer(grpcServer, customerGRPCHandler)
	grpcServer.Serve(lis)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Customer-Service is running",
		})
	})

	api := router.Group("/api")
	{
		api.GET("/customers", customerHandler.GetAll)
		api.GET("/customers/:id", customerHandler.GetByID)
		api.POST("/customers", customerHandler.Create)
		api.PUT("/customers/:id", customerHandler.Update)
		api.DELETE("/customers/:id", customerHandler.Delete)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	router.Run(":" + port)
}
