package main

import (
	"context"
	"log"
	"notification_service/config"
	"notification_service/handler"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Can't load config: %v", err)
	}

	// Connect to RabbitMQ
	rabbitConn, err := config.ConnectRabbitMQ(appConfig.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	// Initialize email consumer
	emailConsumer := handler.NewEmailConsumer(rabbitConn, appConfig.SMTP)
	if err := emailConsumer.Setup(); err != nil {
		log.Fatalf("Failed to setup email consumer: %v", err)
	}

	// Start consumer in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := emailConsumer.Start(ctx); err != nil {
			config.Logger.Errorf("Email consumer failed: %v", err)
		}
	}()

	// Setup HTTP server
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "notification_service",
		})
	})

	// Keep the existing HTTP endpoint for testing
	router.POST("/send", handler.SendEmail(appConfig.SMTP))

	// Start HTTP server in background
	go func() {
		config.Logger.Infof("HTTP server starting on port %s", appConfig.ServerPort)
		if err := router.Run(":" + appConfig.ServerPort); err != nil {
			log.Fatalf("Can't start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	config.Logger.Infof("Received signal %s, shutting down...", sig)
	cancel()

	config.Logger.Info("Notification service stopped")
}
