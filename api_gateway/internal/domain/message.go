package domain

import "context"

type EventType string

const UserRegistered EventType = "user.registered"

type Event struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
}

type UserRegisteredPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type MessageBroker interface {
	Publish(ctx context.Context, routingKey string, event Event) error
	Close() error
}
