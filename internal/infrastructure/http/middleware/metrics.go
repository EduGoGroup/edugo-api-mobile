package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// httpRequestsTotal cuenta el número total de requests HTTP
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestDuration mide la duración de los requests HTTP
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestsInFlight mide el número de requests en progreso
	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)

	// httpResponseSize mide el tamaño de las respuestas HTTP
	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "path", "status"},
	)

	// httpErrorsTotal cuenta errores HTTP (4xx y 5xx)
	httpErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors (4xx and 5xx)",
		},
		[]string{"method", "path", "status", "error_type"},
	)
)

// MetricsMiddleware registra métricas de Prometheus para cada request HTTP.
// Métricas registradas:
// - http_requests_total: Contador de requests por method, path, status
// - http_request_duration_seconds: Histograma de latencias
// - http_requests_in_flight: Gauge de requests en progreso
// - http_response_size_bytes: Histograma de tamaño de respuestas
// - http_errors_total: Contador de errores (4xx, 5xx)
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Incrementar requests en progreso
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		// Capturar tiempo de inicio
		startTime := time.Now()

		// Obtener path normalizado (sin IDs dinámicos para evitar cardinality explosion)
		path := normalizePath(c.FullPath())
		if path == "" {
			path = "unknown"
		}
		method := c.Request.Method

		// Procesar request
		c.Next()

		// Calcular duración
		duration := time.Since(startTime).Seconds()

		// Obtener status code
		statusCode := c.Writer.Status()
		status := strconv.Itoa(statusCode)

		// Registrar métricas
		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path, status).Observe(duration)
		httpResponseSize.WithLabelValues(method, path, status).Observe(float64(c.Writer.Size()))

		// Registrar errores
		if statusCode >= 400 {
			errorType := categorizeError(statusCode)
			httpErrorsTotal.WithLabelValues(method, path, status, errorType).Inc()
		}
	}
}

// normalizePath normaliza el path para evitar cardinality explosion.
// Reemplaza IDs dinámicos con placeholders.
// Ejemplo: /v1/materials/123/summary -> /v1/materials/:id/summary
func normalizePath(path string) string {
	if path == "" {
		return ""
	}
	// Gin ya usa FullPath() que retorna el patrón con :id
	// Por ejemplo: /v1/materials/:id/summary
	return path
}

// categorizeError categoriza el error según el status code
func categorizeError(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "server_error"
	case statusCode == 404:
		return "not_found"
	case statusCode == 401:
		return "unauthorized"
	case statusCode == 403:
		return "forbidden"
	case statusCode == 400:
		return "bad_request"
	case statusCode == 429:
		return "rate_limited"
	default:
		return "client_error"
	}
}

// MetricsConfig permite configurar el middleware de métricas
type MetricsConfig struct {
	// SkipPaths son paths que no registran métricas (ej: /metrics, /health)
	SkipPaths []string
	// Namespace es el prefijo para todas las métricas (ej: "edugo_mobile")
	Namespace string
}

// DefaultMetricsConfig retorna la configuración por defecto
func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		SkipPaths: []string{"/health", "/metrics"},
		Namespace: "edugo_mobile",
	}
}

// MetricsMiddlewareWithConfig crea un middleware de métricas con configuración personalizada
func MetricsMiddlewareWithConfig(config MetricsConfig) gin.HandlerFunc {
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

		// Incrementar requests en progreso
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		// Capturar tiempo de inicio
		startTime := time.Now()

		// Obtener path normalizado
		path := normalizePath(c.FullPath())
		if path == "" {
			path = "unknown"
		}
		method := c.Request.Method

		// Procesar request
		c.Next()

		// Calcular duración
		duration := time.Since(startTime).Seconds()

		// Obtener status code
		statusCode := c.Writer.Status()
		status := strconv.Itoa(statusCode)

		// Registrar métricas
		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path, status).Observe(duration)
		httpResponseSize.WithLabelValues(method, path, status).Observe(float64(c.Writer.Size()))

		// Registrar errores
		if statusCode >= 400 {
			errorType := categorizeError(statusCode)
			httpErrorsTotal.WithLabelValues(method, path, status, errorType).Inc()
		}
	}
}
