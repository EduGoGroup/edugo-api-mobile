package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// HealthHandler maneja el endpoint de health check del sistema.
type HealthHandler struct {
	db      *sql.DB
	mongoDB *mongo.Database
}

// NewHealthHandler crea una nueva instancia de HealthHandler con las dependencias necesarias.
func NewHealthHandler(db *sql.DB, mongoDB *mongo.Database) *HealthHandler {
	return &HealthHandler{
		db:      db,
		mongoDB: mongoDB,
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

// Check godoc
// @Summary Health check
// @Description Verifica que la API y sus dependencias (PostgreSQL, MongoDB) estén funcionando correctamente
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse "System is healthy or degraded with status details"
// @Failure 500 {object} ErrorResponse "System is unhealthy"
// @Router /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
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
