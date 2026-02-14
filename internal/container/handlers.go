package container

import (
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/handler"
)

// HandlerContainer encapsula todos los handlers HTTP
// Responsabilidad: Gestionar la capa de presentación HTTP (REST API)
// Implementa el patrón Adapter entre HTTP y servicios de aplicación
type HandlerContainer struct {
	MaterialHandler   *handler.MaterialHandler
	ProgressHandler   *handler.ProgressHandler
	SummaryHandler    *handler.SummaryHandler
	AssessmentHandler *handler.AssessmentHandler
	StatsHandler      *handler.StatsHandler
	ScreenHandler     *handler.ScreenHandler // Dynamic UI - Phase 1
}

// NewHandlerContainer crea y configura todos los handlers HTTP
// Parámetros:
//   - infra: Contenedor de infraestructura (Logger, S3Client)
//   - services: Contenedor de servicios de aplicación
//
// Retorna un contenedor con todos los handlers inicializados
// Cada handler actúa como adaptador entre Gin y los servicios
func NewHandlerContainer(infra *InfrastructureContainer, services *ServiceContainer) *HandlerContainer {
	return &HandlerContainer{
		// MaterialHandler gestiona CRUD de materiales y URLs presignadas de S3
		MaterialHandler: handler.NewMaterialHandler(
			services.MaterialService,
			infra.S3Client,
			infra.Logger,
		),

		// ProgressHandler gestiona el progreso de lectura
		ProgressHandler: handler.NewProgressHandler(
			services.ProgressService,
			infra.Logger,
		),

		// SummaryHandler gestiona la obtención de resúmenes
		SummaryHandler: handler.NewSummaryHandler(
			services.SummaryService,
			infra.Logger,
		),

		// AssessmentHandler gestiona evaluaciones y attempts (Sprint-04)
		AssessmentHandler: handler.NewAssessmentHandler(
			services.AssessmentAttemptService,
			infra.Logger,
		),

		// StatsHandler gestiona estadísticas globales y por material
		StatsHandler: handler.NewStatsHandler(
			services.StatsService,
			infra.Logger,
		),

		// ScreenHandler gestiona definiciones de pantalla dinámicas (Dynamic UI - Phase 1)
		ScreenHandler: handler.NewScreenHandler(
			services.ScreenService,
			infra.Logger,
		),
	}
}
