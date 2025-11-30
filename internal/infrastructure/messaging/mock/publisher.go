package mock

import (
	"context"

	"github.com/EduGoGroup/edugo-shared/logger"
)

// MockPublisher es un mock de rabbitmq.Publisher que simplemente loguea
// los mensajes sin enviarlos realmente a RabbitMQ.
// Útil para desarrollo y testing sin dependencia de RabbitMQ.
type MockPublisher struct {
	log logger.Logger
}

// NewMockPublisher crea un nuevo MockPublisher
func NewMockPublisher(log logger.Logger) *MockPublisher {
	return &MockPublisher{
		log: log,
	}
}

// Publish simula la publicación de un mensaje logueando la operación
// En lugar de enviar a RabbitMQ, simplemente registra el evento en los logs
func (p *MockPublisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	p.log.Info("📨 [MOCK] Mensaje RabbitMQ (no enviado)",
		"exchange", exchange,
		"routing_key", routingKey,
		"body_size", len(body),
		"body_preview", truncate(string(body), 100),
	)
	return nil
}

// Close simula el cierre de la conexión (no-op para mock)
func (p *MockPublisher) Close() error {
	p.log.Info("🔌 [MOCK] RabbitMQ publisher cerrado (no-op)")
	return nil
}

// truncate trunca un string a una longitud máxima
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Verificar que MockPublisher implementa la interfaz Publisher en tiempo de compilación
// Si no implementa todos los métodos, esto causará un error de compilación
var _ interface {
	Publish(ctx context.Context, exchange, routingKey string, body []byte) error
	Close() error
} = (*MockPublisher)(nil)
