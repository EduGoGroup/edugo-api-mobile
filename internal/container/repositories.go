package container

import (
	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repositories"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	mongoRepo "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/mongodb/repository"
)

// RepositoryContainer encapsula todos los repositorios de dominio
// Responsabilidad: Gestionar el acceso a datos (PostgreSQL y MongoDB)
// Implementa el patrón Repository para abstraer la persistencia
type RepositoryContainer struct {
	// PostgreSQL Repositories
	UserRepository         repository.UserRepository
	MaterialRepository     repository.MaterialRepository
	ProgressRepository     repository.ProgressRepository
	RefreshTokenRepository repository.RefreshTokenRepository
	LoginAttemptRepository repository.LoginAttemptRepository

	// Sprint-03: Assessment System Repositories (PostgreSQL)
	AssessmentRepoV2 repositories.AssessmentRepository // Nuevo de Sprint-03
	AttemptRepo      repositories.AttemptRepository
	AnswerRepo       repositories.AnswerRepository

	// MongoDB Repositories
	SummaryRepository      repository.SummaryRepository
	AssessmentRepository   repository.AssessmentRepository // LEGACY: Solo para stats (MongoDB)
	AssessmentDocumentRepo mongoRepo.AssessmentDocumentRepository
}

// NewRepositoryContainer crea y configura todos los repositorios
// Usa RepositoryFactory para crear implementaciones reales o mock según configuración
//
// Parámetros:
//   - infra: Contenedor de infraestructura con conexiones a bases de datos
//   - cfg: Configuración de la aplicación (determina si usar mocks)
//
// Retorna un contenedor con todos los repositorios inicializados
func NewRepositoryContainer(infra *InfrastructureContainer, cfg *config.Config) *RepositoryContainer {
	factory := NewRepositoryFactory(cfg, infra)

	return &RepositoryContainer{
		// PostgreSQL repositories - creados vía factory
		UserRepository:         factory.CreateUserRepository(),
		MaterialRepository:     factory.CreateMaterialRepository(),
		ProgressRepository:     factory.CreateProgressRepository(),
		RefreshTokenRepository: factory.CreateRefreshTokenRepository(),
		LoginAttemptRepository: factory.CreateLoginAttemptRepository(),

		// Sprint-03: Assessment System Repositories (PostgreSQL) - creados vía factory
		AssessmentRepoV2: factory.CreateAssessmentRepository(),
		AttemptRepo:      factory.CreateAttemptRepository(),
		AnswerRepo:       factory.CreateAnswerRepository(),

		// MongoDB repositories - creados vía factory
		SummaryRepository:      factory.CreateSummaryRepository(),
		AssessmentRepository:   factory.CreateLegacyAssessmentRepository(),
		AssessmentDocumentRepo: factory.CreateAssessmentDocumentRepository(),
	}
}
