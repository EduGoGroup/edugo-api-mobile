package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

// ResilientPublisherConfig configura el comportamiento del circuit breaker para el publisher
type ResilientPublisherConfig struct {
	// Name identifica el circuit breaker en logs y métricas
	Name string
	// MaxRequests es el número máximo de requests permitidos en estado half-open
	MaxRequests uint32
	// Interval es el período de tiempo para resetear contadores en estado closed
	Interval time.Duration
	// Timeout es el tiempo que el circuit breaker permanece en estado open
	Timeout time.Duration
	// FailureThreshold es el número de fallos consecutivos para abrir el circuito
	FailureThreshold uint32
	// OnStateChange callback opcional cuando cambia el estado
	OnStateChange func(name string, from, to gobreaker.State)
}

// DefaultResilientPublisherConfig retorna la configuración por defecto
func DefaultResilientPublisherConfig() ResilientPublisherConfig {
	return ResilientPublisherConfig{
		Name:             "rabbitmq-publisher",
		MaxRequests:      3,
		Interval:         60 * time.Second,
		Timeout:          30 * time.Second,
		FailureThreshold: 5,
	}
}

// ResilientPublisher envuelve un Publisher con Circuit Breaker
type ResilientPublisher struct {
	publisher Publisher
	cb        *gobreaker.CircuitBreaker
	config    ResilientPublisherConfig
}

// NewResilientPublisher crea un nuevo publisher con circuit breaker
func NewResilientPublisher(publisher Publisher, config ResilientPublisherConfig) *ResilientPublisher {
	settings := gobreaker.Settings{
		Name:        config.Name,
		MaxRequests: config.MaxRequests,
		Interval:    config.Interval,
		Timeout:     config.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Abrir el circuito si hay más de FailureThreshold fallos consecutivos
			return counts.ConsecutiveFailures >= config.FailureThreshold
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			if config.OnStateChange != nil {
				config.OnStateChange(name, from, to)
			}
		},
	}

	return &ResilientPublisher{
		publisher: publisher,
		cb:        gobreaker.NewCircuitBreaker(settings),
		config:    config,
	}
}

// Publish publica un mensaje usando el circuit breaker
func (rp *ResilientPublisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	_, err := rp.cb.Execute(func() (interface{}, error) {
		err := rp.publisher.Publish(ctx, exchange, routingKey, body)
		return nil, err
	})

	if err != nil {
		// Si el circuit breaker está abierto, retornar error específico
		if err == gobreaker.ErrOpenState {
			return fmt.Errorf("circuit breaker open for RabbitMQ publisher: %w", err)
		}
		if err == gobreaker.ErrTooManyRequests {
			return fmt.Errorf("circuit breaker half-open, too many requests: %w", err)
		}
		return err
	}

	return nil
}

// State retorna el estado actual del circuit breaker
func (rp *ResilientPublisher) State() gobreaker.State {
	return rp.cb.State()
}

// Counts retorna las estadísticas del circuit breaker
func (rp *ResilientPublisher) Counts() gobreaker.Counts {
	return rp.cb.Counts()
}

// Close cierra el publisher subyacente
func (rp *ResilientPublisher) Close() error {
	return rp.publisher.Close()
}
