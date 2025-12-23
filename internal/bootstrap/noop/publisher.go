package noop

import (
	"context"

	"github.com/EduGoGroup/edugo-shared/logger"
)

// NoopPublisher es una implementación noop de rabbitmq.Publisher
// Se utiliza cuando RabbitMQ no está disponible o está configurado como opcional
type NoopPublisher struct {
	logger logger.Logger
}

// NewNoopPublisher crea una nueva instancia de NoopPublisher
func NewNoopPublisher(log logger.Logger) *NoopPublisher {
	return &NoopPublisher{logger: log}
}

// Publish simula la publicación de un mensaje sin realizar ninguna operación real
// Registra un mensaje de debug para facilitar el debugging
func (p *NoopPublisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	p.logger.Debug("noop publisher: message not published (RabbitMQ not available)",
		"exchange", exchange,
		"routing_key", routingKey,
		"body_size", len(body),
	)
	return nil
}

// Close simula el cierre de la conexión sin realizar ninguna operación real
func (p *NoopPublisher) Close() error {
	p.logger.Debug("noop publisher: close called (no-op)")
	return nil
}
