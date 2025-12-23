package middleware

import (
	"time"

	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
)

// LoggingMiddleware registra información estructurada de cada request HTTP.
// Incluye: request_id, método, path, status, latencia, IP del cliente.
// Útil para debugging, auditoría y monitoreo del sistema.
func LoggingMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capturar tiempo de inicio
		startTime := time.Now()

		// Obtener request ID (generado por RequestIDMiddleware)
		requestID := GetRequestIDFromGin(c)

		// Path y método
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		if rawQuery != "" {
			path = path + "?" + rawQuery
		}
		method := c.Request.Method

		// Procesar request
		c.Next()

		// Calcular latencia
		latency := time.Since(startTime)

		// Obtener status code
		statusCode := c.Writer.Status()

		// Obtener IP del cliente
		clientIP := c.ClientIP()

		// Determinar nivel de log según status code
		logFields := []interface{}{
			"request_id", requestID,
			"method", method,
			"path", path,
			"status", statusCode,
			"latency_ms", latency.Milliseconds(),
			"client_ip", clientIP,
			"user_agent", c.Request.UserAgent(),
		}

		// Agregar user_id si está autenticado
		if userID, exists := c.Get("user_id"); exists {
			logFields = append(logFields, "user_id", userID)
		}

		// Agregar errores si existen
		if len(c.Errors) > 0 {
			logFields = append(logFields, "errors", c.Errors.String())
		}

		// Log según el código de estado
		switch {
		case statusCode >= 500:
			log.Error("Server error", logFields...)
		case statusCode >= 400:
			log.Warn("Client error", logFields...)
		case statusCode >= 300:
			log.Info("Redirect", logFields...)
		default:
			log.Info("Request completed", logFields...)
		}
	}
}

// LoggingConfig permite configurar el middleware de logging
type LoggingConfig struct {
	// SkipPaths son paths que no se loguean (ej: /health para evitar spam)
	SkipPaths []string
	// LogRequestBody si es true, incluye el body del request en el log (cuidado con datos sensibles)
	LogRequestBody bool
	// LogResponseBody si es true, incluye el body de la respuesta en el log
	LogResponseBody bool
}

// DefaultLoggingConfig retorna la configuración por defecto
func DefaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		SkipPaths:       []string{"/health"},
		LogRequestBody:  false,
		LogResponseBody: false,
	}
}

// LoggingMiddlewareWithConfig crea un middleware de logging con configuración personalizada
func LoggingMiddlewareWithConfig(log logger.Logger, config LoggingConfig) gin.HandlerFunc {
	skipPaths := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		// Verificar si el path debe ser omitido
		if skipPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		// Capturar tiempo de inicio
		startTime := time.Now()

		// Obtener request ID (generado por RequestIDMiddleware)
		requestID := GetRequestIDFromGin(c)

		// Path y método
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		if rawQuery != "" {
			path = path + "?" + rawQuery
		}
		method := c.Request.Method

		// Procesar request
		c.Next()

		// Calcular latencia
		latency := time.Since(startTime)

		// Obtener status code
		statusCode := c.Writer.Status()

		// Obtener IP del cliente
		clientIP := c.ClientIP()

		// Construir campos de log
		logFields := []interface{}{
			"request_id", requestID,
			"method", method,
			"path", path,
			"status", statusCode,
			"latency_ms", latency.Milliseconds(),
			"client_ip", clientIP,
		}

		// Agregar user_id si está autenticado
		if userID, exists := c.Get("user_id"); exists {
			logFields = append(logFields, "user_id", userID)
		}

		// Agregar errores si existen
		if len(c.Errors) > 0 {
			logFields = append(logFields, "errors", c.Errors.String())
		}

		// Log según el código de estado
		switch {
		case statusCode >= 500:
			log.Error("Server error", logFields...)
		case statusCode >= 400:
			log.Warn("Client error", logFields...)
		default:
			log.Debug("Request completed", logFields...)
		}
	}
}
