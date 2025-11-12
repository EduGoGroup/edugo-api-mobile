package adapter

import (
	"context"
	"fmt"

	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
	"github.com/EduGoGroup/edugo-shared/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

// MessagePublisherAdapter adapta *amqp.Channel (de shared/bootstrap) a rabbitmq.Publisher (interfaz de api-mobile)
// Este adapter permite que el código de api-mobile siga usando la interfaz rabbitmq.Publisher
// mientras internamente usamos *amqp.Channel retornado por shared/bootstrap
type MessagePublisherAdapter struct {
	channel  *amqp.Channel
	exchange string
	logger   logger.Logger
}

// NewMessagePublisherAdapter crea un nuevo adapter de message publisher
// channel: canal de RabbitMQ retornado por shared/bootstrap
// exchange: nombre del exchange a usar para publicar mensajes
// logger: logger para registrar errores y eventos
func NewMessagePublisherAdapter(
	channel *amqp.Channel,
	exchange string,
	log logger.Logger,
) rabbitmq.Publisher {
	return &MessagePublisherAdapter{
		channel:  channel,
		exchange: exchange,
		logger:   log,
	}
}

// Publish publica un mensaje al exchange especificado con el routing key dado
// Implementa la interfaz rabbitmq.Publisher
func (a *MessagePublisherAdapter) Publish(
	ctx context.Context,
	exchange string,
	routingKey string,
	body []byte,
) error {
	// Si no se especifica exchange, usar el default del adapter
	if exchange == "" {
		exchange = a.exchange
	}

	// Publicar mensaje con contexto
	err := a.channel.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Mensajes persistentes
		},
	)

	if err != nil {
		a.logger.Error("failed to publish message",
			"exchange", exchange,
			"routing_key", routingKey,
			"error", err,
		)
		return fmt.Errorf("failed to publish message to exchange %s: %w", exchange, err)
	}

	a.logger.Debug("message published successfully",
		"exchange", exchange,
		"routing_key", routingKey,
		"body_size", len(body),
	)

	return nil
}

// Close cierra el canal de RabbitMQ
// Implementa la interfaz rabbitmq.Publisher
func (a *MessagePublisherAdapter) Close() error {
	if a.channel == nil {
		return nil
	}

	a.logger.Info("closing RabbitMQ channel")
	
	err := a.channel.Close()
	if err != nil {
		a.logger.Error("failed to close RabbitMQ channel", "error", err)
		return fmt.Errorf("failed to close RabbitMQ channel: %w", err)
	}

	a.logger.Info("RabbitMQ channel closed successfully")
	return nil
}

// Verificar en compile-time que MessagePublisherAdapter implementa rabbitmq.Publisher
var _ rabbitmq.Publisher = (*MessagePublisherAdapter)(nil)
