package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

type Config struct {
	SMTP        *SMTPConfig
	RabbitMQURL string
	ServerPort  string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		Logger.Warn("No .env file found")
	}

	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	return &Config{
		SMTP: &SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     port,
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			Sender:   os.Getenv("FROM_EMAIL"),
		},
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		ServerPort:  getEnv("SERVER_PORT", "8085"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func ConnectRabbitMQ(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		Logger.Errorf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	Logger.Info("Connected to RabbitMQ successfully")
	return conn, nil
}
