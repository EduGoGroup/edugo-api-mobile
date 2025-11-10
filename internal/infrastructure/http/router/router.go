package router

import (
	"github.com/EduGoGroup/edugo-api-mobile/internal/container"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/handler"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/middleware"
	ginmiddleware "github.com/EduGoGroup/edugo-shared/middleware/gin"
	"github.com/gin-gonic/gin"
)

// SetupRouter configura todas las rutas de la aplicación con sus respectivos handlers y middleware.
// Separa las rutas públicas de las protegidas y organiza los endpoints por recursos.
func SetupRouter(c *container.Container, healthHandler *handler.HealthHandler) *gin.Engine {
	r := gin.Default()

	// Middleware global
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.ClientInfoMiddleware()) // Extraer IP y User-Agent del cliente

	// Health check (público, sin versión)
	r.GET("/health", healthHandler.Check)

	// Swagger UI (público) con detección dinámica de host
	SetupSwaggerUI(r)

	// Grupo de rutas API v1
	v1 := r.Group("/v1")
	{
		// Rutas públicas de autenticación
		setupAuthPublicRoutes(v1, c)

		// Rutas protegidas (requieren JWT)
		setupProtectedRoutes(v1, c)
	}

	return r
}

// setupAuthPublicRoutes configura las rutas públicas de autenticación.
func setupAuthPublicRoutes(rg *gin.RouterGroup, c *container.Container) {
	rg.POST("/auth/login", c.Handlers.AuthHandler.Login)
	rg.POST("/auth/refresh", c.Handlers.AuthHandler.Refresh)
}

// setupProtectedRoutes configura todas las rutas que requieren autenticación JWT.
func setupProtectedRoutes(rg *gin.RouterGroup, c *container.Container) {
	protected := rg.Group("")
	protected.Use(ginmiddleware.JWTAuthMiddleware(c.Infrastructure.JWTManager))
	{
		// Rutas de autenticación protegidas
		setupAuthProtectedRoutes(protected, c)

		// Rutas de materiales
		setupMaterialRoutes(protected, c)

		// Rutas de evaluaciones (assessments)
		setupAssessmentRoutes(protected, c)

		// Rutas de progreso (progress)
		setupProgressRoutes(protected, c)

		// Rutas de estadísticas globales
		setupStatsRoutes(protected, c)
	}
}

// setupAuthProtectedRoutes configura las rutas de autenticación que requieren JWT.
func setupAuthProtectedRoutes(rg *gin.RouterGroup, c *container.Container) {
	rg.POST("/auth/logout", c.Handlers.AuthHandler.Logout)
	rg.POST("/auth/revoke-all", c.Handlers.AuthHandler.RevokeAll)
}

// setupMaterialRoutes configura todas las rutas relacionadas con materiales educativos.
func setupMaterialRoutes(rg *gin.RouterGroup, c *container.Container) {
	materials := rg.Group("/materials")
	{
		// CRUD básico de materiales
		materials.GET("", c.Handlers.MaterialHandler.ListMaterials)
		materials.POST("", c.Handlers.MaterialHandler.CreateMaterial)
		materials.GET("/:id", c.Handlers.MaterialHandler.GetMaterial)
		materials.POST("/:id/upload-complete", c.Handlers.MaterialHandler.NotifyUploadComplete)

		// Historial de versiones de materiales
		materials.GET("/:id/versions", c.Handlers.MaterialHandler.GetMaterialWithVersions)

		// URLs presignadas para S3
		materials.POST("/:id/upload-url", c.Handlers.MaterialHandler.GenerateUploadURL)
		materials.GET("/:id/download-url", c.Handlers.MaterialHandler.GenerateDownloadURL)

		// Resúmenes de materiales
		materials.GET("/:id/summary", c.Handlers.SummaryHandler.GetSummary)

		// Evaluaciones (assessments)
		materials.GET("/:id/assessment", c.Handlers.AssessmentHandler.GetAssessment)
		materials.POST("/:id/assessment/attempts", c.Handlers.AssessmentHandler.RecordAttempt)

		// Progreso del estudiante
		materials.PATCH("/:id/progress", c.Handlers.ProgressHandler.UpdateProgress)

		// Estadísticas de materiales
		materials.GET("/:id/stats", c.Handlers.StatsHandler.GetMaterialStats)
	}
}

// setupAssessmentRoutes configura todas las rutas relacionadas con evaluaciones.
func setupAssessmentRoutes(rg *gin.RouterGroup, c *container.Container) {
	assessments := rg.Group("/assessments")
	{
		// Submit de evaluación con cálculo automático de score y feedback detallado
		assessments.POST("/:id/submit", c.Handlers.AssessmentHandler.SubmitAssessment)
	}
}

// setupProgressRoutes configura todas las rutas relacionadas con progreso de usuarios.
func setupProgressRoutes(rg *gin.RouterGroup, c *container.Container) {
	progress := rg.Group("/progress")
	{
		// Endpoint UPSERT idempotente de progreso (Fase 5)
		// PUT (no POST) porque la operación es idempotente
		progress.PUT("", c.Handlers.ProgressHandler.UpsertProgress)
	}
}

// setupStatsRoutes configura todas las rutas relacionadas con estadísticas globales del sistema.
func setupStatsRoutes(rg *gin.RouterGroup, c *container.Container) {
	stats := rg.Group("/stats")
	{
		// Estadísticas globales del sistema (Fase 6)
		// TODO: Agregar middleware de autorización para solo admins
		stats.GET("/global", c.Handlers.StatsHandler.GetGlobalStats)
	}
}
