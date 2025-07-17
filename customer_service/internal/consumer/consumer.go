package consumer

import (
	"context"
	"customer_service/internal/domain"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

const (
	exchangeName = "customer_events"
	queueName    = "customer_creation_queue"
	routingKey   = "customer.create"
)

type EventType string

const (
	UserRegistered EventType = "user.registered"
)

type Event struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type UserRegisteredPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type CustomerConsumer struct {
	conn            *amqp.Connection
	customerUsecase domain.CustomerUsecase
	logger          *logrus.Logger
}

func NewCustomerConsumer(conn *amqp.Connection, customerUsecase domain.CustomerUsecase) *CustomerConsumer {
	return &CustomerConsumer{
		conn:            conn,
		customerUsecase: customerUsecase,
		logger:          logrus.New(),
	}
}

func (c *CustomerConsumer) Setup() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare exchange
	err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	// Declare queue
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *CustomerConsumer) Start(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Qos(1, 0, false)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	c.logger.Info("Customer consumer started, waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				c.logger.Warn("Message channel closed")
				return nil
			}

			c.processMessage(msg)
		}
	}
}

func (c *CustomerConsumer) processMessage(msg amqp.Delivery) {
	var event Event
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		c.logger.Errorf("Failed to unmarshal event: %v", err)
		msg.Nack(false, false)
		return
	}

	c.logger.Infof("Received event: %s", event.Type)

	switch event.Type {
	case UserRegistered:
		var payload UserRegisteredPayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			c.logger.Errorf("Failed to unmarshal payload: %v", err)
			msg.Nack(false, false)
			return
		}

		c.handleUserRegistered(msg, payload)
	default:
		c.logger.Warnf("Unknown event type: %s", event.Type)
		msg.Ack(false)
	}
}

func (c *CustomerConsumer) handleUserRegistered(msg amqp.Delivery, payload UserRegisteredPayload) {
	// Create customer request
	customerReq := &domain.CustomerRequest{
		Name:  payload.Name,
		Email: payload.Email,
		Phone: "", // Will be updated later by user
	}

	// Create customer
	customer, err := c.customerUsecase.Create(context.Background(), customerReq)
	if err != nil {
		c.logger.Errorf("Failed to create customer: %v", err)
		msg.Nack(false, true) // Requeue for retry
		return
	}

	c.logger.Infof("Customer created successfully: ID=%s, Email=%s",
		customer.Id.Hex(), customer.Email)

	msg.Ack(false)
}
