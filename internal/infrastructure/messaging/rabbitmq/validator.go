package rabbitmq

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/EduGoGroup/edugo-infrastructure/schemas"
)

var (
	validator     *schemas.EventValidator
	validatorOnce sync.Once
	validatorErr  error
)

// InitValidator inicializa el validador de eventos una sola vez (singleton)
// Debe llamarse al inicio de la aplicación
func InitValidator() error {
	validatorOnce.Do(func() {
		validator, validatorErr = schemas.NewEventValidator()
	})
	return validatorErr
}

// GetValidator retorna la instancia del validador (debe haberse inicializado antes)
func GetValidator() (*schemas.EventValidator, error) {
	if validator == nil {
		return nil, fmt.Errorf("validator not initialized, call InitValidator() first")
	}
	return validator, nil
}

// ValidateEvent valida un evento contra su schema correspondiente
// El evento debe implementar métodos EventType() y EventVersion()
func ValidateEvent(event interface{}) error {
	v, err := GetValidator()
	if err != nil {
		return err
	}

	// Serializar evento a JSON para validación
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshaling event: %w", err)
	}

	// Determinar tipo y versión del evento
	eventType, eventVersion, err := getEventTypeAndVersion(event)
	if err != nil {
		return err
	}

	// Validar contra schema
	if err := v.ValidateJSON(eventJSON, eventType, eventVersion); err != nil {
		return fmt.Errorf("event validation failed for %s v%s: %w", eventType, eventVersion, err)
	}

	return nil
}

// getEventTypeAndVersion extrae el tipo y versión del evento
func getEventTypeAndVersion(event interface{}) (string, string, error) {
	// Si es un Event con envelope estándar
	if e, ok := event.(Event); ok {
		return e.EventType, e.EventVersion, nil
	}

	// Fallback para tipos legacy (si existen)
	return "", "", fmt.Errorf("event must be of type Event with standard envelope, got: %T", event)
}
