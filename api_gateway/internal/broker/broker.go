package broker

import (
	"api_gateway/internal/domain"
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

const (
	exchangeName = "customer_events"
)

type rabbitMQBroker struct {
	conn   *amqp.Connection
	logger *logrus.Logger
}

func NewRabbitMQBroker(conn *amqp.Connection) domain.MessageBroker {
	return &rabbitMQBroker{
		conn:   conn,
		logger: logrus.New(),
	}
}

func (r *rabbitMQBroker) Publish(ctx context.Context, routingKey string, event domain.Event) error {
	ch, err := r.conn.Channel()
	if err != nil {
		r.logger.Errorf("Failed to open channel: %v", err)
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.logger.Errorf("Failed to declare exchange: %v", err)
		return err
	}
	body, err := json.Marshal(event)
	if err != nil {
		r.logger.Errorf("Failed to marshal event: %v", err)
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		})
	if err != nil {
		r.logger.Errorf("Failed to publish message: %v", err)
		return err
	}

	r.logger.Infof("Published event %s with routing key %s", event.Type, routingKey)
	return nil
}

func (r *rabbitMQBroker) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
