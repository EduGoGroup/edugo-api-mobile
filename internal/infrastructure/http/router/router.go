package router

import (
	"github.com/EduGoGroup/edugo-api-mobile/internal/container"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/handler"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/middleware"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// SetupRouter configura todas las rutas de la aplicación con sus respectivos handlers y middleware.
// Separa las rutas públicas de las protegidas y organiza los endpoints por recursos.
//
// NOTA: Las rutas de autenticación (login, refresh, logout, etc.) han sido migradas
// a api-admin como parte de la centralización de autenticación (Sprint 3).
// Este servicio valida tokens contra api-admin usando RemoteAuthMiddleware.
func SetupRouter(c *container.Container, healthHandler *handler.HealthHandler) *gin.Engine {
	r := gin.New() // Usar gin.New() en lugar de Default() para control total de middleware

	// Middleware global (orden importante)
	r.Use(gin.Recovery())                   // 1. Recuperar de panics
	r.Use(middleware.RequestIDMiddleware()) // 2. Generar/propagar request ID para tracing
	r.Use(middleware.MetricsMiddlewareWithConfig(middleware.MetricsConfig{
		SkipPaths: []string{"/health", "/metrics"}, // No registrar métricas de endpoints de infraestructura
	})) // 3. Métricas Prometheus
	r.Use(middleware.LoggingMiddlewareWithConfig(c.Infrastructure.Logger, middleware.LoggingConfig{
		SkipPaths: []string{"/health", "/metrics"}, // No loguear endpoints de infraestructura
	})) // 4. Logging estructurado con request_id
	r.Use(middleware.CORS())                 // 5. CORS headers
	r.Use(middleware.ClientInfoMiddleware()) // 6. Extraer IP y User-Agent del cliente

	// Endpoints de infraestructura (públicos, sin versión)
	r.GET("/health", healthHandler.Check)
	r.GET("/metrics", gin.WrapH(promhttp.Handler())) // Prometheus metrics

	// Swagger UI (público) con detección dinámica de host
	SetupSwaggerUI(r)

	// Grupo de rutas API v1
	v1 := r.Group("/v1")
	{
		// Rutas protegidas (requieren JWT validado contra api-admin)
		setupProtectedRoutes(v1, c)
	}

	return r
}

// setupProtectedRoutes configura todas las rutas que requieren autenticación JWT.
// Los tokens son validados remotamente contra api-admin usando AuthClient.
func setupProtectedRoutes(rg *gin.RouterGroup, c *container.Container) {
	protected := rg.Group("")

	// Middleware de autenticación remota con api-admin
	// Reemplaza la validación local de JWT por validación centralizada
	protected.Use(middleware.RemoteAuthMiddleware(middleware.RemoteAuthConfig{
		AuthClient: c.Infrastructure.AuthClient,
		Logger:     c.Infrastructure.Logger,
		SkipPaths:  []string{}, // Todas las rutas en este grupo requieren auth
	}))
	{
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

// setupMaterialRoutes configura todas las rutas relacionadas con materiales educativos.
func setupMaterialRoutes(rg *gin.RouterGroup, c *container.Container) {
	materials := rg.Group("/materials")
	{
		// Lectura de materiales (requiere permiso materials:read)
		materials.GET("",
			middleware.RequirePermission(enum.PermissionMaterialsRead),
			c.Handlers.MaterialHandler.ListMaterials,
		)
		materials.GET("/:id",
			middleware.RequirePermission(enum.PermissionMaterialsRead),
			c.Handlers.MaterialHandler.GetMaterial,
		)
		materials.GET("/:id/versions",
			middleware.RequirePermission(enum.PermissionMaterialsRead),
			c.Handlers.MaterialHandler.GetMaterialWithVersions,
		)
		materials.GET("/:id/download-url",
			middleware.RequirePermission(enum.PermissionMaterialsDownload),
			c.Handlers.MaterialHandler.GenerateDownloadURL,
		)
		materials.GET("/:id/summary",
			middleware.RequirePermission(enum.PermissionMaterialsRead),
			c.Handlers.SummaryHandler.GetSummary,
		)
		materials.GET("/:id/assessment",
			middleware.RequirePermission(enum.PermissionAssessmentsRead),
			c.Handlers.AssessmentHandler.GetMaterialAssessment,
		)
		materials.GET("/:id/stats",
			middleware.RequirePermission(enum.PermissionStatsUnit),
			c.Handlers.StatsHandler.GetMaterialStats,
		)

		// Creación/modificación de materiales (requiere permisos específicos)
		materials.POST("",
			middleware.RequirePermission(enum.PermissionMaterialsCreate),
			c.Handlers.MaterialHandler.CreateMaterial,
		)
		materials.POST("/:id/upload-complete",
			middleware.RequirePermission(enum.PermissionMaterialsCreate),
			c.Handlers.MaterialHandler.NotifyUploadComplete,
		)
		materials.POST("/:id/upload-url",
			middleware.RequirePermission(enum.PermissionMaterialsCreate),
			c.Handlers.MaterialHandler.GenerateUploadURL,
		)
		materials.PUT("/:id",
			middleware.RequirePermission(enum.PermissionMaterialsUpdate),
			c.Handlers.MaterialHandler.UpdateMaterial,
		)

		// Intentos de evaluación (requiere permiso assessments:attempt)
		materials.POST("/:id/assessment/attempts",
			middleware.RequirePermission(enum.PermissionAssessmentsAttempt),
			c.Handlers.AssessmentHandler.CreateMaterialAttempt,
		)
	}
}

// setupAssessmentRoutes configura todas las rutas relacionadas con evaluaciones.
func setupAssessmentRoutes(rg *gin.RouterGroup, c *container.Container) {
	// Rutas de intentos (attempts) - Sprint-04
	attempts := rg.Group("/attempts")
	{
		attempts.GET("/:id/results",
			middleware.RequirePermission(enum.PermissionAssessmentsViewResults),
			c.Handlers.AssessmentHandler.GetAttemptResults,
		)
	}

	// Rutas de historial de usuario - Sprint-04
	users := rg.Group("/users")
	{
		users.GET("/me/attempts",
			middleware.RequirePermission(enum.PermissionAssessmentsViewResults),
			c.Handlers.AssessmentHandler.GetUserAttemptHistory,
		)
	}
}

// setupProgressRoutes configura todas las rutas relacionadas con progreso de usuarios.
func setupProgressRoutes(rg *gin.RouterGroup, c *container.Container) {
	progress := rg.Group("/progress")
	{
		// Endpoint UPSERT idempotente de progreso (Fase 5)
		// PUT (no POST) porque la operación es idempotente
		progress.PUT("",
			middleware.RequirePermission(enum.PermissionProgressUpdate),
			c.Handlers.ProgressHandler.UpsertProgress,
		)
	}
}

// setupStatsRoutes configura todas las rutas relacionadas con estadísticas globales del sistema.
func setupStatsRoutes(rg *gin.RouterGroup, c *container.Container) {
	stats := rg.Group("/stats")
	{
		// Estadísticas globales del sistema (requiere permiso stats:global)
		stats.GET("/global",
			middleware.RequirePermission(enum.PermissionStatsGlobal),
			c.Handlers.StatsHandler.GetGlobalStats,
		)
	}
}
