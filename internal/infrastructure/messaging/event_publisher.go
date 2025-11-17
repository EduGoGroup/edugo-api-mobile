package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
)

// EventPublisher es un wrapper sobre RabbitMQPublisher que valida eventos antes de publicarlos
type EventPublisher struct {
	publisher rabbitmq.Publisher
	exchange  string
}

// NewEventPublisher crea una nueva instancia de EventPublisher
func NewEventPublisher(publisher rabbitmq.Publisher, exchange string) *EventPublisher {
	return &EventPublisher{
		publisher: publisher,
		exchange:  exchange,
	}
}

// PublishMaterialUploaded publica un evento material.uploaded validado
func (p *EventPublisher) PublishMaterialUploaded(ctx context.Context, payload MaterialUploadedPayload) error {
	// Crear evento con envelope estándar
	event := NewMaterialUploadedEvent(payload)

	// Validar evento antes de publicar
	if err := ValidateEvent(event); err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}

	// Serializar evento
	body, err := event.ToJSON()
	if err != nil {
		return fmt.Errorf("error serializing event: %w", err)
	}

	// Publicar en RabbitMQ
	routingKey := "material.uploaded"
	if err := p.publisher.Publish(ctx, p.exchange, routingKey, body); err != nil {
		return fmt.Errorf("error publishing event: %w", err)
	}

	return nil
}

// PublishAssessmentGenerated publica un evento assessment.generated validado
func (p *EventPublisher) PublishAssessmentGenerated(ctx context.Context, payload AssessmentGeneratedPayload) error {
	// Crear evento con envelope estándar
	event := NewAssessmentGeneratedEvent(payload)

	// Validar evento antes de publicar
	if err := ValidateEvent(event); err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}

	// Serializar evento
	body, err := event.ToJSON()
	if err != nil {
		return fmt.Errorf("error serializing event: %w", err)
	}

	// Publicar en RabbitMQ
	routingKey := "assessment.generated"
	if err := p.publisher.Publish(ctx, p.exchange, routingKey, body); err != nil {
		return fmt.Errorf("error publishing event: %w", err)
	}

	return nil
}

// PublishEvent es un método genérico que valida y publica cualquier evento
// El evento debe estar registrado en getEventTypeAndVersion()
func (p *EventPublisher) PublishEvent(ctx context.Context, event interface{}) error {
	// Validar evento
	if err := ValidateEvent(event); err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}

	// Serializar evento
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error serializing event: %w", err)
	}

	// Obtener routing key basado en tipo de evento
	routingKey, err := getRoutingKey(event)
	if err != nil {
		return err
	}

	// Publicar en RabbitMQ
	if err := p.publisher.Publish(ctx, p.exchange, routingKey, body); err != nil {
		return fmt.Errorf("error publishing event: %w", err)
	}

	return nil
}

// getRoutingKey retorna la routing key basada en el tipo de evento
func getRoutingKey(event interface{}) (string, error) {
	// Si es un Event con envelope estándar, usar EventType
	if e, ok := event.(Event); ok {
		return e.EventType, nil
	}

	return "", fmt.Errorf("event must be of type Event with standard envelope, got: %T", event)
}

// Close cierra la conexión del publisher subyacente
func (p *EventPublisher) Close() error {
	return p.publisher.Close()
}
