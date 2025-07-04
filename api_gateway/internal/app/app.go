package app

import (
	"api_gateway/config"
	"api_gateway/internal/handler"
	"api_gateway/internal/repository"
	"api_gateway/internal/usecase"
	"api_gateway/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func Run() {
	db, err := config.ConnectMongoDB()
	if err != nil {
		panic("Error connecting to database")
	}
	defer db.Close()

	authRepo := repository.NewAuthRepository(db.DB)
	authUsecase := usecase.NewAuthUsecase(authRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.SetupCORS())
	router.Static("/swagger-ui", "./static/swagger-ui")
	router.GET("/swagger/product_service.json", middleware.ProxySwaggerDoc("http://localhost:8082/swagger/doc.json"))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API Gateway is running",
		})
	})

	api := router.Group("/api")
	{
		api.POST("/users", authHandler.Register)
		api.GET("/users", authHandler.Login)

	}
	api.Any("/product_service/*proxyPath", middleware.ReverseProxyMiddleware(os.Getenv("PRODUCT_SERVICE_HOST"), "/api/product_service"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
