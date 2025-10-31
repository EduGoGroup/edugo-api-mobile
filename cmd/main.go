package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-api-mobile/internal/handlers"
	"github.com/EduGoGroup/edugo-api-mobile/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/EduGoGroup/edugo-api-mobile/docs" // Swagger docs generados por swag init
)

// @title EduGo API Mobile
// @version 1.0
// @description API para operaciones frecuentes de docentes y estudiantes en EduGo
// @termsOfService http://edugo.com/terms/

// @contact.name Equipo EduGo
// @contact.email soporte@edugo.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token con formato: Bearer {token}

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Error loading configuration: %v", err)
	}

	// Mostrar ambiente activo
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}
	log.Printf("🌍 Environment: %s", env)
	log.Printf("📊 Log Level: %s, Format: %s", cfg.Logging.Level, cfg.Logging.Format)

	// Configurar Gin
	r := gin.Default()

	// Middleware global
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.RateLimiter())

	// Health check
	r.GET("/health", HealthCheck)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rutas públicas
	v1 := r.Group("/v1")
	{
		// Autenticación
		v1.POST("/auth/login", handlers.Login)
	}

	// Rutas protegidas (requieren JWT)
	protected := v1.Group("")
	protected.Use(middleware.AuthRequired())
	{
		// Materials
		materials := protected.Group("/materials")
		{
			materials.GET("", handlers.GetMaterials)
			materials.POST("", handlers.CreateMaterial)
			materials.GET("/:id", handlers.GetMaterialDetail)
			materials.POST("/:id/upload-complete", handlers.UploadComplete)
			materials.GET("/:id/summary", handlers.GetMaterialSummary)
			materials.GET("/:id/assessment", handlers.GetAssessment)
			materials.POST("/:id/assessment/attempts", handlers.RecordAttempt)
			materials.PATCH("/:id/progress", handlers.UpdateProgress)
			materials.GET("/:id/stats", handlers.GetMaterialStats)
		}
	}

	// Iniciar servidor usando configuración
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🚀 API Mobile running on http://localhost:%d", cfg.Server.Port)
	log.Printf("📚 Swagger UI: http://localhost:%d/swagger/index.html", cfg.Server.Port)
	log.Printf("🗄️  PostgreSQL: %s:%d/%s", cfg.Database.Postgres.Host, cfg.Database.Postgres.Port, cfg.Database.Postgres.Database)
	log.Printf("🍃 MongoDB: %s", cfg.Database.MongoDB.Database)
	log.Printf("🐰 RabbitMQ: Connected")

	if err := r.Run(addr); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}
}

// HealthCheck godoc
// @Summary Health check
// @Description Verifica que la API está funcionando
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "edugo-api-mobile",
		"version": "1.0.0",
	})
}
