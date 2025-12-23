package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockHealthChecker implementa HealthChecker para testing
type MockHealthChecker struct {
	healthy bool
	err     error
}

func (m *MockHealthChecker) CheckHealth(ctx context.Context) error {
	if !m.healthy {
		return m.err
	}
	return nil
}

func TestHealthHandler_Check_Basic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns healthy when all services are mocked", func(t *testing.T) {
		handler := NewHealthHandler(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

		handler.Check(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response HealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "healthy", response.Status)
		assert.Equal(t, "mock", response.Postgres)
		assert.Equal(t, "mock", response.MongoDB)
		assert.Equal(t, "edugo-api-mobile", response.Service)
	})
}

func TestHealthHandler_CheckDetailed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns detailed response when detail=1", func(t *testing.T) {
		handler := NewHealthHandler(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/health?detail=1", nil)

		handler.Check(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response DetailedHealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "healthy", response.Status)
		assert.Equal(t, "edugo-api-mobile", response.Service)
		assert.NotEmpty(t, response.TotalTime)
		assert.NotEmpty(t, response.Components)

		// Verificar componentes
		assert.Equal(t, "mock", response.Components["postgres"].Status)
		assert.Equal(t, "mock", response.Components["mongodb"].Status)
		assert.Equal(t, "not_configured", response.Components["rabbitmq"].Status)
		assert.True(t, response.Components["rabbitmq"].Optional)
		assert.Equal(t, "not_configured", response.Components["s3"].Status)
		assert.True(t, response.Components["s3"].Optional)
	})

	t.Run("returns healthy RabbitMQ when checker is healthy", func(t *testing.T) {
		rabbitChecker := &MockHealthChecker{healthy: true}
		handler := NewHealthHandlerWithCheckers(nil, nil, rabbitChecker, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/health?detail=1", nil)

		handler.Check(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response DetailedHealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "healthy", response.Components["rabbitmq"].Status)
		assert.True(t, response.Components["rabbitmq"].Optional)
	})

	t.Run("returns unhealthy RabbitMQ when checker fails", func(t *testing.T) {
		rabbitChecker := &MockHealthChecker{healthy: false, err: errors.New("connection refused")}
		handler := NewHealthHandlerWithCheckers(nil, nil, rabbitChecker, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/health?detail=1", nil)

		handler.Check(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response DetailedHealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		// Estado general sigue siendo healthy porque RabbitMQ es opcional
		assert.Equal(t, "healthy", response.Status)
		assert.Equal(t, "unhealthy", response.Components["rabbitmq"].Status)
		assert.Equal(t, "connection refused", response.Components["rabbitmq"].Error)
		assert.True(t, response.Components["rabbitmq"].Optional)
	})

	t.Run("returns healthy S3 when checker is healthy", func(t *testing.T) {
		s3Checker := &MockHealthChecker{healthy: true}
		handler := NewHealthHandlerWithCheckers(nil, nil, nil, s3Checker)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/health?detail=1", nil)

		handler.Check(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response DetailedHealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "healthy", response.Components["s3"].Status)
		assert.True(t, response.Components["s3"].Optional)
	})

	t.Run("returns unhealthy S3 when checker fails", func(t *testing.T) {
		s3Checker := &MockHealthChecker{healthy: false, err: errors.New("bucket not found")}
		handler := NewHealthHandlerWithCheckers(nil, nil, nil, s3Checker)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/health?detail=1", nil)

		handler.Check(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response DetailedHealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		// Estado general sigue siendo healthy porque S3 es opcional
		assert.Equal(t, "healthy", response.Status)
		assert.Equal(t, "unhealthy", response.Components["s3"].Status)
		assert.Equal(t, "bucket not found", response.Components["s3"].Error)
		assert.True(t, response.Components["s3"].Optional)
	})

	t.Run("includes latency for all components", func(t *testing.T) {
		rabbitChecker := &MockHealthChecker{healthy: true}
		s3Checker := &MockHealthChecker{healthy: true}
		handler := NewHealthHandlerWithCheckers(nil, nil, rabbitChecker, s3Checker)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/health?detail=1", nil)

		handler.Check(c)

		var response DetailedHealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verificar que todos los componentes tienen latencia
		for name, component := range response.Components {
			assert.NotEmpty(t, component.Latency, "component %s should have latency", name)
		}
	})
}

func TestNewHealthHandlerWithCheckers(t *testing.T) {
	rabbitChecker := &MockHealthChecker{healthy: true}
	s3Checker := &MockHealthChecker{healthy: true}

	handler := NewHealthHandlerWithCheckers(nil, nil, rabbitChecker, s3Checker)

	assert.NotNil(t, handler)
	assert.Equal(t, rabbitChecker, handler.rabbitMQChecker)
	assert.Equal(t, s3Checker, handler.s3Checker)
}
