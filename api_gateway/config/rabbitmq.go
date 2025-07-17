package config

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

var (
	Logger = logrus.New()
)

func InitRabbitMQ() (*amqp.Connection, error) {
	url := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(url)
	if err != nil {
		Logger.Errorf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	Logger.Info("Connected to RabbitMQ successfully")
	return conn, nil
}
