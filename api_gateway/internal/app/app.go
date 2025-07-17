package app

import (
	"api_gateway/config"
	_ "api_gateway/docs"
	"api_gateway/internal/broker"
	"api_gateway/internal/handler"
	"api_gateway/internal/repository"
	"api_gateway/internal/usecase"
	"api_gateway/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

// @title Order Management API
// @version Microservice
// @description This is a sample server for managing orders, customers, products, and carts.
// @host localhost:8080
// @BasePath /api
func Run() {
	rabbitConn, err := config.InitRabbitMQ()
	if err != nil {
		panic("Error connecting to RabbitMQ")
	}
	db, err := config.ConnectMongoDB()
	if err != nil {
		panic("Error connecting to database")
	}
	defer db.Close()

	messageBroker := broker.NewRabbitMQBroker(rabbitConn)
	authRepo := repository.NewAuthRepository(db.DB)
	authUsecase := usecase.NewAuthUsecase(authRepo, messageBroker)
	authHandler := handler.NewAuthHandler(authUsecase)

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.SetupCORS())
	router.Static("/swagger-ui", "./static/swagger-ui")
	router.GET("swagger/auth_service.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})
	router.GET("/swagger/cart_service.json", middleware.ProxySwaggerDoc(os.Getenv("CART_SERVICE_HOST")+"/swagger/doc.json"))
	router.GET("/swagger/product_service.json", middleware.ProxySwaggerDoc(os.Getenv("PRODUCT_SERVICE_HOST")+"/swagger/doc.json"))
	router.GET("/swagger/customer_service.json", middleware.ProxySwaggerDoc(os.Getenv("CUSTOMER_SERVICE_HOST")+"/swagger/doc.json"))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API Gateway is running",
		})
	})

	api := router.Group("/api")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/registeradmin", authHandler.RegisterAdmin)
		api.GET("/auth/login", authHandler.Login)

	}
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.Any("/product_service/*proxyPath", middleware.ReverseProxyMiddleware(os.Getenv("PRODUCT_SERVICE_HOST"), "/product_service"))
			admin.Any("/customer_service/*proxyPath", middleware.ReverseProxyMiddleware(os.Getenv("CUSTOMER_SERVICE_HOST"), "/customer_service"))
			admin.Any("/cart_service/*proxyPath", middleware.ReverseProxyMiddleware(os.Getenv("CART_SERVICE_HOST"), "/cart_service"))

		}

		user := protected.Group("/user")
		user.Use(middleware.RequireRole("user", "admin"))
		{
			user.Any("/product_service/*proxyPath", middleware.ReverseProxyMiddleware(os.Getenv("PRODUCT_SERVICE_HOST"), "/product_service"))
			user.Any("/customer_service/*proxyPath", middleware.ReverseProxyMiddleware(os.Getenv("CUSTOMER_SERVICE_HOST"), "/customer_service"))
			user.Any("/cart_service/*proxyPath", middleware.ReverseProxyMiddleware(os.Getenv("CART_SERVICE_HOST"), "/cart_service"))
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
