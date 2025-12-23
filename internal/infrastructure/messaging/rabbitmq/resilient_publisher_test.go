package rabbitmq

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPublisher es un mock del Publisher para testing
type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	args := m.Called(ctx, exchange, routingKey, body)
	return args.Error(0)
}

func (m *MockPublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewResilientPublisher(t *testing.T) {
	mockPub := new(MockPublisher)
	config := DefaultResilientPublisherConfig()

	rp := NewResilientPublisher(mockPub, config)

	assert.NotNil(t, rp)
	assert.Equal(t, gobreaker.StateClosed, rp.State())
}

func TestResilientPublisher_Publish_Success(t *testing.T) {
	mockPub := new(MockPublisher)
	config := DefaultResilientPublisherConfig()
	rp := NewResilientPublisher(mockPub, config)

	ctx := context.Background()
	payload := []byte(`{"test":"data"}`)

	mockPub.On("Publish", ctx, "exchange", "routing.key", payload).Return(nil)

	err := rp.Publish(ctx, "exchange", "routing.key", payload)

	assert.NoError(t, err)
	mockPub.AssertExpectations(t)
}

func TestResilientPublisher_Publish_Error(t *testing.T) {
	mockPub := new(MockPublisher)
	config := DefaultResilientPublisherConfig()
	rp := NewResilientPublisher(mockPub, config)

	ctx := context.Background()
	payload := []byte(`{"test":"data"}`)
	expectedErr := errors.New("connection failed")

	mockPub.On("Publish", ctx, "exchange", "routing.key", payload).Return(expectedErr)

	err := rp.Publish(ctx, "exchange", "routing.key", payload)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockPub.AssertExpectations(t)
}

func TestResilientPublisher_CircuitBreaker_Opens(t *testing.T) {
	mockPub := new(MockPublisher)

	// Configuración para abrir el circuito rápidamente
	config := ResilientPublisherConfig{
		Name:             "test-cb",
		MaxRequests:      1,
		Interval:         10 * time.Second,
		Timeout:          1 * time.Second,
		FailureThreshold: 3, // Abrir después de 3 fallos consecutivos
	}

	var stateChanges []gobreaker.State
	var mu sync.Mutex
	config.OnStateChange = func(name string, from, to gobreaker.State) {
		mu.Lock()
		stateChanges = append(stateChanges, to)
		mu.Unlock()
	}

	rp := NewResilientPublisher(mockPub, config)

	ctx := context.Background()
	payload := []byte(`{"test":"data"}`)
	expectedErr := errors.New("connection failed")

	// Configurar el mock para fallar siempre
	mockPub.On("Publish", ctx, "exchange", "routing.key", payload).Return(expectedErr)

	// Hacer 3 llamadas fallidas para abrir el circuito
	for i := 0; i < 3; i++ {
		_ = rp.Publish(ctx, "exchange", "routing.key", payload)
	}

	// El circuito debería estar abierto ahora
	assert.Equal(t, gobreaker.StateOpen, rp.State())

	// La siguiente llamada debería fallar inmediatamente con circuit breaker open
	err := rp.Publish(ctx, "exchange", "routing.key", payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit breaker open")

	// Verificar que hubo cambio de estado
	mu.Lock()
	assert.Contains(t, stateChanges, gobreaker.StateOpen)
	mu.Unlock()
}

func TestResilientPublisher_CircuitBreaker_HalfOpen(t *testing.T) {
	mockPub := new(MockPublisher)

	// Configuración con timeout corto para pasar a half-open rápido
	config := ResilientPublisherConfig{
		Name:             "test-cb-halfopen",
		MaxRequests:      1,
		Interval:         10 * time.Second,
		Timeout:          100 * time.Millisecond, // Timeout muy corto para test
		FailureThreshold: 2,
	}

	rp := NewResilientPublisher(mockPub, config)

	ctx := context.Background()
	payload := []byte(`{"test":"data"}`)
	connErr := errors.New("connection failed")

	// Configurar mock para fallar primero, luego tener éxito
	call := mockPub.On("Publish", ctx, "exchange", "routing.key", payload)
	callCount := int32(0)

	call.RunFn = func(args mock.Arguments) {
		count := atomic.AddInt32(&callCount, 1)
		if count <= 2 {
			call.ReturnArguments = mock.Arguments{connErr}
		} else {
			call.ReturnArguments = mock.Arguments{nil}
		}
	}

	// Abrir el circuito con 2 fallos
	_ = rp.Publish(ctx, "exchange", "routing.key", payload)
	_ = rp.Publish(ctx, "exchange", "routing.key", payload)

	assert.Equal(t, gobreaker.StateOpen, rp.State())

	// Esperar a que pase a half-open
	time.Sleep(150 * time.Millisecond)

	assert.Equal(t, gobreaker.StateHalfOpen, rp.State())

	// La siguiente llamada exitosa debería cerrar el circuito
	err := rp.Publish(ctx, "exchange", "routing.key", payload)
	assert.NoError(t, err)
	assert.Equal(t, gobreaker.StateClosed, rp.State())
}

func TestResilientPublisher_Close(t *testing.T) {
	mockPub := new(MockPublisher)
	config := DefaultResilientPublisherConfig()
	rp := NewResilientPublisher(mockPub, config)

	mockPub.On("Close").Return(nil)

	err := rp.Close()

	assert.NoError(t, err)
	mockPub.AssertExpectations(t)
}

func TestResilientPublisher_Counts(t *testing.T) {
	mockPub := new(MockPublisher)
	config := DefaultResilientPublisherConfig()
	rp := NewResilientPublisher(mockPub, config)

	ctx := context.Background()
	payload := []byte(`{"test":"data"}`)

	mockPub.On("Publish", ctx, "exchange", "routing.key", payload).Return(nil)

	// Hacer algunas llamadas exitosas
	_ = rp.Publish(ctx, "exchange", "routing.key", payload)
	_ = rp.Publish(ctx, "exchange", "routing.key", payload)

	counts := rp.Counts()
	assert.Equal(t, uint32(2), counts.TotalSuccesses)
	assert.Equal(t, uint32(0), counts.TotalFailures)
}

func TestDefaultResilientPublisherConfig(t *testing.T) {
	config := DefaultResilientPublisherConfig()

	assert.Equal(t, "rabbitmq-publisher", config.Name)
	assert.Equal(t, uint32(3), config.MaxRequests)
	assert.Equal(t, 60*time.Second, config.Interval)
	assert.Equal(t, 30*time.Second, config.Timeout)
	assert.Equal(t, uint32(5), config.FailureThreshold)
	assert.Nil(t, config.OnStateChange)
}
