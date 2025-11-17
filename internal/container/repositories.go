package container

import (
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repositories"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	mongoRepo "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/mongodb/repository"
	postgresRepo "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/postgres/repository"
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
	SummaryRepository         repository.SummaryRepository
	AssessmentRepository      repository.AssessmentRepository // Legacy
	AssessmentDocumentRepo    mongoRepo.AssessmentDocumentRepository
}

// NewRepositoryContainer crea y configura todos los repositorios
// Parámetros:
//   - infra: Contenedor de infraestructura con conexiones a bases de datos
//
// Retorna un contenedor con todos los repositorios inicializados
// Los repositorios de PostgreSQL usan infra.DB
// Los repositorios de MongoDB usan infra.MongoDB
func NewRepositoryContainer(infra *InfrastructureContainer) *RepositoryContainer {
	return &RepositoryContainer{
		// PostgreSQL repositories
		UserRepository:         postgresRepo.NewPostgresUserRepository(infra.DB),
		MaterialRepository:     postgresRepo.NewPostgresMaterialRepository(infra.DB),
		ProgressRepository:     postgresRepo.NewPostgresProgressRepository(infra.DB),
		RefreshTokenRepository: postgresRepo.NewPostgresRefreshTokenRepository(infra.DB),
		LoginAttemptRepository: postgresRepo.NewPostgresLoginAttemptRepository(infra.DB),

		// Sprint-03: Assessment System Repositories (PostgreSQL)
		AssessmentRepoV2: postgresRepo.NewPostgresAssessmentRepository(infra.DB),
		AttemptRepo:      postgresRepo.NewPostgresAttemptRepository(infra.DB),
		AnswerRepo:       postgresRepo.NewPostgresAnswerRepository(infra.DB),

		// MongoDB repositories
		SummaryRepository:     mongoRepo.NewMongoSummaryRepository(infra.MongoDB),
		AssessmentRepository:  mongoRepo.NewMongoAssessmentRepository(infra.MongoDB), // Legacy
		AssessmentDocumentRepo: mongoRepo.NewMongoAssessmentDocumentRepository(infra.MongoDB),
	}
}
