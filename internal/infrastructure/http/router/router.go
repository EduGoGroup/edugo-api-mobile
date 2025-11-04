package router

import (
	"github.com/EduGoGroup/edugo-api-mobile/internal/container"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/handler"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/middleware"
	ginmiddleware "github.com/EduGoGroup/edugo-shared/middleware/gin"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configura todas las rutas de la aplicación con sus respectivos handlers y middleware.
// Separa las rutas públicas de las protegidas y organiza los endpoints por recursos.
func SetupRouter(c *container.Container, healthHandler *handler.HealthHandler) *gin.Engine {
	r := gin.Default()

	// Middleware global
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// Health check (público, sin versión)
	r.GET("/health", healthHandler.Check)

	// Swagger UI (público)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
	rg.POST("/auth/login", c.AuthHandler.Login)
	rg.POST("/auth/refresh", c.AuthHandler.Refresh)
}

// setupProtectedRoutes configura todas las rutas que requieren autenticación JWT.
func setupProtectedRoutes(rg *gin.RouterGroup, c *container.Container) {
	protected := rg.Group("")
	protected.Use(ginmiddleware.JWTAuthMiddleware(c.JWTManager))
	{
		// Rutas de autenticación protegidas
		setupAuthProtectedRoutes(protected, c)

		// Rutas de materiales
		setupMaterialRoutes(protected, c)
	}
}

// setupAuthProtectedRoutes configura las rutas de autenticación que requieren JWT.
func setupAuthProtectedRoutes(rg *gin.RouterGroup, c *container.Container) {
	rg.POST("/auth/logout", c.AuthHandler.Logout)
	rg.POST("/auth/revoke-all", c.AuthHandler.RevokeAll)
}

// setupMaterialRoutes configura todas las rutas relacionadas con materiales educativos.
func setupMaterialRoutes(rg *gin.RouterGroup, c *container.Container) {
	materials := rg.Group("/materials")
	{
		// CRUD básico de materiales
		materials.GET("", c.MaterialHandler.ListMaterials)
		materials.POST("", c.MaterialHandler.CreateMaterial)
		materials.GET("/:id", c.MaterialHandler.GetMaterial)
		materials.POST("/:id/upload-complete", c.MaterialHandler.NotifyUploadComplete)

		// Resúmenes de materiales
		materials.GET("/:id/summary", c.SummaryHandler.GetSummary)

		// Evaluaciones (assessments)
		materials.GET("/:id/assessment", c.AssessmentHandler.GetAssessment)
		materials.POST("/:id/assessment/attempts", c.AssessmentHandler.RecordAttempt)

		// Progreso del estudiante
		materials.PATCH("/:id/progress", c.ProgressHandler.UpdateProgress)

		// Estadísticas de materiales
		materials.GET("/:id/stats", c.StatsHandler.GetMaterialStats)
	}
}
