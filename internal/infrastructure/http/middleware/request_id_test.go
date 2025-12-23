package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestIDMiddleware_GeneratesNewID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestIDMiddleware())

	var capturedRequestID string
	router.GET("/test", func(c *gin.Context) {
		capturedRequestID = GetRequestIDFromGin(c)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verificar que se generó un request ID
	assert.NotEmpty(t, capturedRequestID)

	// Verificar que es un UUID válido
	_, err := uuid.Parse(capturedRequestID)
	assert.NoError(t, err, "El request ID debe ser un UUID válido")

	// Verificar que el header de respuesta contiene el request ID
	responseRequestID := w.Header().Get(RequestIDHeader)
	assert.Equal(t, capturedRequestID, responseRequestID)
}

func TestRequestIDMiddleware_UsesExistingID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestIDMiddleware())

	existingID := "test-request-id-12345"
	var capturedRequestID string

	router.GET("/test", func(c *gin.Context) {
		capturedRequestID = GetRequestIDFromGin(c)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(RequestIDHeader, existingID)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verificar que usa el ID existente
	assert.Equal(t, existingID, capturedRequestID)

	// Verificar que el header de respuesta contiene el mismo ID
	responseRequestID := w.Header().Get(RequestIDHeader)
	assert.Equal(t, existingID, responseRequestID)
}

func TestGetRequestID_FromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestIDMiddleware())

	var capturedFromContext string
	router.GET("/test", func(c *gin.Context) {
		// Obtener del contexto de Go
		capturedFromContext = GetRequestID(c.Request.Context())
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.NotEmpty(t, capturedFromContext)
	_, err := uuid.Parse(capturedFromContext)
	assert.NoError(t, err)
}

func TestGetRequestID_NilContext(t *testing.T) {
	//nolint:staticcheck // SA1012: Intencionalmente probamos el comportamiento con nil
	result := GetRequestID(nil)
	assert.Empty(t, result)
}

func TestGetRequestID_EmptyContext(t *testing.T) {
	ctx := context.Background()
	result := GetRequestID(ctx)
	assert.Empty(t, result)
}

func TestGetRequestIDFromGin_NilContext(t *testing.T) {
	result := GetRequestIDFromGin(nil)
	assert.Empty(t, result)
}

func TestMustGetRequestID_ReturnsExisting(t *testing.T) {
	expectedID := "existing-request-id"
	ctx := context.WithValue(context.Background(), RequestIDKey, expectedID)

	result := MustGetRequestID(ctx)
	assert.Equal(t, expectedID, result)
}

func TestMustGetRequestID_GeneratesNew(t *testing.T) {
	ctx := context.Background()

	result := MustGetRequestID(ctx)

	require.NotEmpty(t, result)
	_, err := uuid.Parse(result)
	assert.NoError(t, err, "Debe generar un UUID válido cuando no existe")
}

func TestRequestIDMiddleware_UniquePerRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestIDMiddleware())

	var ids []string
	router.GET("/test", func(c *gin.Context) {
		ids = append(ids, GetRequestIDFromGin(c))
		c.Status(http.StatusOK)
	})

	// Hacer múltiples requests
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// Verificar que todos los IDs son únicos
	idSet := make(map[string]bool)
	for _, id := range ids {
		assert.False(t, idSet[id], "Los request IDs deben ser únicos")
		idSet[id] = true
	}
}
