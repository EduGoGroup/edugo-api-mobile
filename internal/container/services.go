package container

import (
	"github.com/EduGoGroup/edugo-api-mobile/internal/application/service"
)

// ServiceContainer encapsula todos los servicios de aplicación
// Responsabilidad: Gestionar la lógica de negocio de la aplicación
// Implementa casos de uso y orquesta repositorios
type ServiceContainer struct {
	AuthService       service.AuthService
	MaterialService   service.MaterialService
	ProgressService   service.ProgressService
	SummaryService    service.SummaryService
	AssessmentService service.AssessmentService
	StatsService      service.StatsService
}

// NewServiceContainer crea y configura todos los servicios de aplicación
// Parámetros:
//   - infra: Contenedor de infraestructura (Logger, JWTManager, MessagePublisher)
//   - repos: Contenedor de repositorios para acceso a datos
//
// Retorna un contenedor con todos los servicios inicializados
// Cada servicio recibe sus dependencias específicas según el principio DIP
func NewServiceContainer(infra *InfrastructureContainer, repos *RepositoryContainer) *ServiceContainer {
	return &ServiceContainer{
		// AuthService gestiona autenticación, tokens JWT y sesiones
		AuthService: service.NewAuthService(
			repos.UserRepository,
			repos.RefreshTokenRepository,
			repos.LoginAttemptRepository,
			infra.JWTManager,
			infra.Logger,
		),

		// MaterialService gestiona materiales educativos y versionado
		MaterialService: service.NewMaterialService(
			repos.MaterialRepository,
			infra.MessagePublisher,
			infra.Logger,
		),

		// ProgressService gestiona el progreso de lectura de estudiantes
		ProgressService: service.NewProgressService(
			repos.ProgressRepository,
			infra.Logger,
		),

		// SummaryService gestiona resúmenes de materiales (MongoDB)
		SummaryService: service.NewSummaryService(
			repos.SummaryRepository,
			infra.Logger,
		),

		// AssessmentService gestiona evaluaciones y scoring con Strategy Pattern
		AssessmentService: service.NewAssessmentService(
			repos.AssessmentRepository,
			infra.MessagePublisher,
			infra.Logger,
		),

		// StatsService gestiona estadísticas globales y por material
		// Usa queries paralelas con goroutines para optimización
		StatsService: service.NewStatsService(
			infra.Logger,
			repos.MaterialRepository,
			repos.AssessmentRepository,
			repos.ProgressRepository,
		),
	}
}
