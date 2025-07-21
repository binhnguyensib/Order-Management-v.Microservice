package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"notification_service/config"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/mail.v2"
)

type EmailRequest struct {
	To      string `json:"to" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

type WelcomeEmailData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

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

// Consumer để nhận events từ RabbitMQ
type EmailConsumer struct {
	conn       *amqp.Connection
	smtpConfig *config.SMTPConfig
}

func NewEmailConsumer(conn *amqp.Connection, smtpConfig *config.SMTPConfig) *EmailConsumer {
	return &EmailConsumer{
		conn:       conn,
		smtpConfig: smtpConfig,
	}
}

func (ec *EmailConsumer) Setup() error {
	ch, err := ec.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"customer_events", // name (giống với API Gateway)
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"notification_user_queue", // name
		true,                      // durable
		false,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,
		"customer.create", // routing key (giống với API Gateway)
		"customer_events", // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	config.Logger.Info("Email consumer setup completed")
	return nil
}

func (ec *EmailConsumer) Start(ctx context.Context) error {
	ch, err := ec.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Qos(1, 0, false)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		"notification_user_queue",
		"",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	config.Logger.Info("Email consumer started, waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				config.Logger.Warn("Message channel closed")
				return nil
			}

			ec.processMessage(msg)
		}
	}
}

func (ec *EmailConsumer) processMessage(msg amqp.Delivery) {
	var event Event
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		config.Logger.Errorf("Failed to unmarshal event: %v", err)
		msg.Nack(false, false)
		return
	}

	config.Logger.Infof("Received event: %s", event.Type)

	switch event.Type {
	case UserRegistered:
		if err := ec.handleUserRegistered(event.Payload); err != nil {
			config.Logger.Errorf("Failed to handle user registered event: %v", err)
			msg.Nack(false, true) // Requeue for retry
			return
		}
	default:
		config.Logger.Warnf("Unknown event type: %s", event.Type)
	}

	msg.Ack(false)
	config.Logger.Infof("Successfully processed event: %s", event.Type)
}

func (ec *EmailConsumer) handleUserRegistered(payload json.RawMessage) error {
	var data UserRegisteredPayload
	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}

	// Send welcome email
	return ec.sendWelcomeEmail(data.Email, data.Name)
}

func (ec *EmailConsumer) sendWelcomeEmail(to, name string) error {
	// Load welcome template
	tmpl, err := template.ParseFiles("templates/welcome.html")
	if err != nil {
		config.Logger.Errorf("Failed to load welcome template: %v", err)
		// Fallback to simple text email
		return ec.sendSimpleWelcomeEmail(to, name)
	}

	// Render template
	var bodyBuffer bytes.Buffer
	templateData := WelcomeEmailData{
		Name:  name,
		Email: to,
	}

	if err := tmpl.Execute(&bodyBuffer, templateData); err != nil {
		config.Logger.Errorf("Failed to render template: %v", err)
		return ec.sendSimpleWelcomeEmail(to, name)
	}

	// Create and send email
	msg := mail.NewMessage()
	msg.SetHeader("From", ec.smtpConfig.Sender)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "🎉 Chào mừng bạn đến với hệ thống của chúng tôi!")
	msg.SetBody("text/html", bodyBuffer.String())

	dialer := mail.NewDialer(
		ec.smtpConfig.Host,
		ec.smtpConfig.Port,
		ec.smtpConfig.Username,
		ec.smtpConfig.Password,
	)

	if err := dialer.DialAndSend(msg); err != nil {
		config.Logger.Errorf("Failed to send welcome email to %s: %v", to, err)
		return err
	}

	config.Logger.Infof("Welcome email sent successfully to %s", to)
	return nil
}

func (ec *EmailConsumer) sendSimpleWelcomeEmail(to, name string) error {
	body := `
Xin chào ` + name + `!

Cảm ơn bạn đã đăng ký tài khoản với chúng tôi. 
Chúng tôi rất vui mừng có bạn tham gia!

Bây giờ bạn có thể:
- Duyệt và mua sắm các sản phẩm
- Thêm sản phẩm vào giỏ hàng  
- Đặt hàng và theo dõi đơn hàng
- Theo dõi tình trạng giao hàng

Chúc bạn có trải nghiệm tuyệt vời!

Trân trọng,
Đội ngũ hỗ trợ
`

	msg := mail.NewMessage()
	msg.SetHeader("From", ec.smtpConfig.Sender)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Chào mừng bạn đến với hệ thống của chúng tôi!")
	msg.SetBody("text/plain", body)

	dialer := mail.NewDialer(
		ec.smtpConfig.Host,
		ec.smtpConfig.Port,
		ec.smtpConfig.Username,
		ec.smtpConfig.Password,
	)

	return dialer.DialAndSend(msg)
}

// HTTP Handler (giữ nguyên để test)
func SendEmail(smtpConfig *config.SMTPConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req EmailRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
			return
		}

		msg := mail.NewMessage()
		msg.SetHeader("From", smtpConfig.Sender)
		msg.SetHeader("To", req.To)
		msg.SetHeader("Subject", req.Subject)
		msg.SetBody("text/plain", req.Body)

		dialer := mail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.Password)

		if err := dialer.DialAndSend(msg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi gửi email: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email đã được gửi thành công!"})
	}
}
