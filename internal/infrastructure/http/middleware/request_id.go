package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDKey es la key usada para almacenar/recuperar el request ID del contexto
const RequestIDKey contextKey = "request_id"

// RequestIDHeader es el nombre del header HTTP para el request ID
const RequestIDHeader = "X-Request-ID"

// RequestIDMiddleware genera o propaga un request ID único para cada request.
// Si el cliente envía un X-Request-ID header, lo usa; si no, genera uno nuevo.
// El request ID se propaga en:
// - El contexto de Go (para logs y servicios internos)
// - El header de respuesta X-Request-ID
// - El contexto de Gin (c.Get("request_id"))
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener request ID del header si existe
		requestID := c.GetHeader(RequestIDHeader)

		// Si no existe, generar uno nuevo
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Almacenar en el contexto de Gin para acceso fácil en handlers
		c.Set(string(RequestIDKey), requestID)

		// Agregar al contexto de Go para propagación a servicios
		ctx := context.WithValue(c.Request.Context(), RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		// Agregar el request ID al header de respuesta
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID extrae el request ID del contexto.
// Retorna el request ID si existe, o string vacío si no.
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}

	return ""
}

// GetRequestIDFromGin extrae el request ID del contexto de Gin.
// Útil cuando se tiene acceso directo al contexto de Gin.
func GetRequestIDFromGin(c *gin.Context) string {
	if c == nil {
		return ""
	}

	if requestID, exists := c.Get(string(RequestIDKey)); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}

	return ""
}

// MustGetRequestID extrae el request ID del contexto.
// Si no existe, genera uno nuevo (útil para casos edge).
func MustGetRequestID(ctx context.Context) string {
	requestID := GetRequestID(ctx)
	if requestID == "" {
		return uuid.New().String()
	}
	return requestID
}
