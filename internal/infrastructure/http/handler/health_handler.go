package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// HealthChecker define la interfaz para verificar el estado de un servicio
type HealthChecker interface {
	// CheckHealth verifica el estado del servicio y retorna un error si no está saludable
	CheckHealth(ctx context.Context) error
}

// HealthHandler maneja el endpoint de health check del sistema.
type HealthHandler struct {
	db              *sql.DB
	mongoDB         *mongo.Database
	rabbitMQChecker HealthChecker
	s3Checker       HealthChecker
}

// NewHealthHandler crea una nueva instancia de HealthHandler con las dependencias necesarias.
func NewHealthHandler(db *sql.DB, mongoDB *mongo.Database) *HealthHandler {
	return &HealthHandler{
		db:      db,
		mongoDB: mongoDB,
	}
}

// NewHealthHandlerWithCheckers crea un HealthHandler con checkers adicionales para RabbitMQ y S3
func NewHealthHandlerWithCheckers(
	db *sql.DB,
	mongoDB *mongo.Database,
	rabbitMQChecker HealthChecker,
	s3Checker HealthChecker,
) *HealthHandler {
	return &HealthHandler{
		db:              db,
		mongoDB:         mongoDB,
		rabbitMQChecker: rabbitMQChecker,
		s3Checker:       s3Checker,
	}
}

// HealthResponse representa la respuesta del endpoint de health check.
type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Postgres  string `json:"postgres"`
	MongoDB   string `json:"mongodb"`
	Timestamp string `json:"timestamp"`
}

// DetailedHealthResponse representa la respuesta detallada del endpoint de health check.
type DetailedHealthResponse struct {
	Status     string                     `json:"status"`
	Service    string                     `json:"service"`
	Version    string                     `json:"version"`
	Timestamp  string                     `json:"timestamp"`
	Components map[string]ComponentHealth `json:"components"`
	TotalTime  string                     `json:"total_time"`
}

// ComponentHealth representa el estado de salud de un componente individual.
type ComponentHealth struct {
	Status   string `json:"status"`
	Latency  string `json:"latency"`
	Error    string `json:"error,omitempty"`
	Optional bool   `json:"optional,omitempty"`
}

// Check godoc
// @Summary Health check
// @Description Verifica que la API y sus dependencias (PostgreSQL, MongoDB) estén funcionando correctamente
// @Tags health
// @Produce json
// @Param detail query string false "Incluir detalles de cada componente (1=detallado)"
// @Success 200 {object} HealthResponse "System is healthy or degraded with status details"
// @Failure 500 {object} ErrorResponse "System is unhealthy"
// @Router /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	// Si se solicita detalle, retornar respuesta detallada
	if c.Query("detail") == "1" {
		h.checkDetailed(c)
		return
	}

	// Verificar PostgreSQL
	pgStatus := "healthy"
	if h.db == nil {
		pgStatus = "mock"
	} else if err := h.db.Ping(); err != nil {
		pgStatus = "unhealthy"
	}

	// Verificar MongoDB
	mongoStatus := "healthy"
	if h.mongoDB == nil {
		mongoStatus = "mock"
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := h.mongoDB.Client().Ping(ctx, nil); err != nil {
			mongoStatus = "unhealthy"
		}
	}

	// Determinar estado general del sistema
	status := "healthy"
	if pgStatus == "unhealthy" || mongoStatus == "unhealthy" {
		status = "degraded"
	}

	response := HealthResponse{
		Status:    status,
		Service:   "edugo-api-mobile",
		Version:   "1.0.0",
		Postgres:  pgStatus,
		MongoDB:   mongoStatus,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}

// checkDetailed realiza un health check detallado con latencias de cada componente
func (h *HealthHandler) checkDetailed(c *gin.Context) {
	startTotal := time.Now()
	components := make(map[string]ComponentHealth)
	overallStatus := "healthy"
	ctx := c.Request.Context()

	// Check PostgreSQL
	pgHealth := h.checkPostgres(ctx)
	components["postgres"] = pgHealth
	if pgHealth.Status == "unhealthy" {
		overallStatus = "degraded"
	}

	// Check MongoDB
	mongoHealth := h.checkMongoDB(ctx)
	components["mongodb"] = mongoHealth
	if mongoHealth.Status == "unhealthy" {
		overallStatus = "degraded"
	}

	// Check RabbitMQ (opcional)
	rabbitHealth := h.checkRabbitMQ(ctx)
	components["rabbitmq"] = rabbitHealth
	// RabbitMQ es opcional, no afecta el estado general

	// Check S3 (opcional)
	s3Health := h.checkS3(ctx)
	components["s3"] = s3Health
	// S3 es opcional, no afecta el estado general

	totalTime := time.Since(startTotal)

	response := DetailedHealthResponse{
		Status:     overallStatus,
		Service:    "edugo-api-mobile",
		Version:    "1.0.0",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Components: components,
		TotalTime:  totalTime.String(),
	}

	c.JSON(http.StatusOK, response)
}

// checkPostgres verifica el estado de PostgreSQL
func (h *HealthHandler) checkPostgres(ctx context.Context) ComponentHealth {
	if h.db == nil {
		return ComponentHealth{
			Status:  "mock",
			Latency: "0ms",
		}
	}

	start := time.Now()
	err := h.db.PingContext(ctx)
	latency := time.Since(start)

	if err != nil {
		return ComponentHealth{
			Status:  "unhealthy",
			Latency: latency.String(),
			Error:   err.Error(),
		}
	}

	return ComponentHealth{
		Status:  "healthy",
		Latency: latency.String(),
	}
}

// checkMongoDB verifica el estado de MongoDB
func (h *HealthHandler) checkMongoDB(ctx context.Context) ComponentHealth {
	if h.mongoDB == nil {
		return ComponentHealth{
			Status:  "mock",
			Latency: "0ms",
		}
	}

	start := time.Now()
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err := h.mongoDB.Client().Ping(pingCtx, nil)
	latency := time.Since(start)

	if err != nil {
		return ComponentHealth{
			Status:  "unhealthy",
			Latency: latency.String(),
			Error:   err.Error(),
		}
	}

	return ComponentHealth{
		Status:  "healthy",
		Latency: latency.String(),
	}
}

// checkRabbitMQ verifica el estado de RabbitMQ
func (h *HealthHandler) checkRabbitMQ(ctx context.Context) ComponentHealth {
	if h.rabbitMQChecker == nil {
		return ComponentHealth{
			Status:   "not_configured",
			Latency:  "0ms",
			Optional: true,
		}
	}

	start := time.Now()
	checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err := h.rabbitMQChecker.CheckHealth(checkCtx)
	latency := time.Since(start)

	if err != nil {
		return ComponentHealth{
			Status:   "unhealthy",
			Latency:  latency.String(),
			Error:    err.Error(),
			Optional: true,
		}
	}

	return ComponentHealth{
		Status:   "healthy",
		Latency:  latency.String(),
		Optional: true,
	}
}

// checkS3 verifica el estado de S3
func (h *HealthHandler) checkS3(ctx context.Context) ComponentHealth {
	if h.s3Checker == nil {
		return ComponentHealth{
			Status:   "not_configured",
			Latency:  "0ms",
			Optional: true,
		}
	}

	start := time.Now()
	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	err := h.s3Checker.CheckHealth(checkCtx)
	latency := time.Since(start)

	if err != nil {
		return ComponentHealth{
			Status:   "unhealthy",
			Latency:  latency.String(),
			Error:    err.Error(),
			Optional: true,
		}
	}

	return ComponentHealth{
		Status:   "healthy",
		Latency:  latency.String(),
		Optional: true,
	}
}
