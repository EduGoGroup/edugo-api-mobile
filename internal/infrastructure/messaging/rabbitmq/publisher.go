package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-shared/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// Publisher define la interfaz para publicar mensajes
type Publisher interface {
	Publish(ctx context.Context, exchange, routingKey string, body []byte) error
	Close() error
}

// RabbitMQPublisher implementa Publisher usando RabbitMQ
type RabbitMQPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	logger   logger.Logger
}

// NewRabbitMQPublisher crea una nueva instancia de RabbitMQPublisher
func NewRabbitMQPublisher(url, exchange string, log logger.Logger) (*RabbitMQPublisher, error) {
	publisher := &RabbitMQPublisher{
		exchange: exchange,
		logger:   log,
	}

	if err := publisher.Connect(url); err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return publisher, nil
}

// Connect establece la conexión con RabbitMQ y declara el exchange
func (p *RabbitMQPublisher) Connect(url string) error {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to dial RabbitMQ: %w", err)
	}
	p.conn = conn

	// Crear canal
	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close() // Ignorar error en cleanup
		return fmt.Errorf("failed to open channel: %w", err)
	}
	p.channel = channel

	// Declarar exchange de tipo topic
	err = channel.ExchangeDeclare(
		p.exchange, // nombre
		"topic",    // tipo
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // argumentos
	)
	if err != nil {
		_ = channel.Close() // Ignorar error en cleanup
		_ = conn.Close()    // Ignorar error en cleanup
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Habilitar publisher confirms
	err = channel.Confirm(false)
	if err != nil {
		_ = channel.Close() // Ignorar error en cleanup
		_ = conn.Close()    // Ignorar error en cleanup
		return fmt.Errorf("failed to enable publisher confirms: %w", err)
	}

	p.logger.Info("Connected to RabbitMQ successfully",
		zap.String("exchange", p.exchange),
	)

	return nil
}

// Publish publica un mensaje en el exchange especificado
func (p *RabbitMQPublisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	if p.channel == nil {
		return fmt.Errorf("channel is not initialized")
	}

	// Crear confirmaciones
	confirms := p.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	// Publicar mensaje
	err := p.channel.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // mensaje persistente
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		p.logger.Error("Failed to publish message",
			zap.String("exchange", exchange),
			zap.String("routing_key", routingKey),
			zap.Error(err),
		)
		return fmt.Errorf("failed to publish message: %w", err)
	}

	// Esperar confirmación con timeout
	select {
	case confirm := <-confirms:
		if !confirm.Ack {
			p.logger.Warn("Message not acknowledged by broker",
				zap.String("exchange", exchange),
				zap.String("routing_key", routingKey),
			)
			return fmt.Errorf("message not acknowledged by broker")
		}
	case <-time.After(5 * time.Second):
		p.logger.Warn("Timeout waiting for publisher confirmation",
			zap.String("exchange", exchange),
			zap.String("routing_key", routingKey),
		)
		return fmt.Errorf("timeout waiting for publisher confirmation")
	}

	p.logger.Debug("Message published successfully",
		zap.String("exchange", exchange),
		zap.String("routing_key", routingKey),
		zap.Int("body_size", len(body)),
	)

	return nil
}

// Close cierra la conexión y el canal
func (p *RabbitMQPublisher) Close() error {
	var errs []error

	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close channel: %w", err))
		}
	}

	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing RabbitMQ: %v", errs)
	}

	p.logger.Info("RabbitMQ connection closed successfully")
	return nil
}
