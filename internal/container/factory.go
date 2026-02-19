package container

import (
	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repositories"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	mockMongo "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/mock/mongodb"
	mockPostgres "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/mock/postgres"
	mongoRepo "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/mongodb/repository"
	postgresRepo "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/postgres/repository"
)

// RepositoryFactory es una fábrica que crea instancias de repositorios,
// pudiendo cambiar entre implementaciones reales (PostgreSQL/MongoDB) y mocks
// según la configuración de desarrollo.
type RepositoryFactory struct {
	config *config.Config
	infra  *InfrastructureContainer
}

// NewRepositoryFactory crea una nueva instancia de RepositoryFactory.
// Esta fábrica permite cambiar entre repositorios reales y mocks basándose
// en la configuración Development.UseMockRepositories.
func NewRepositoryFactory(cfg *config.Config, infra *InfrastructureContainer) *RepositoryFactory {
	return &RepositoryFactory{config: cfg, infra: infra}
}

func (f *RepositoryFactory) CreateUserRepository() repository.UserRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockUserRepository()
	}
	if f.infra.DB == nil {
		panic("PostgreSQL DB connection is nil but mock repositories are disabled")
	}
	return postgresRepo.NewPostgresUserRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateMaterialRepository() repository.MaterialRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockMaterialRepository()
	}
	if f.infra.DB == nil {
		panic("PostgreSQL DB connection is nil but mock repositories are disabled")
	}
	return postgresRepo.NewPostgresMaterialRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateProgressRepository() repository.ProgressRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockProgressRepository()
	}
	if f.infra.DB == nil {
		panic("PostgreSQL DB connection is nil but mock repositories are disabled")
	}
	return postgresRepo.NewPostgresProgressRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateRefreshTokenRepository() repository.RefreshTokenRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockRefreshTokenRepository()
	}
	if f.infra.DB == nil {
		panic("PostgreSQL DB connection is nil but mock repositories are disabled")
	}
	return postgresRepo.NewPostgresRefreshTokenRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateLoginAttemptRepository() repository.LoginAttemptRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockLoginAttemptRepository()
	}
	if f.infra.DB == nil {
		panic("PostgreSQL DB connection is nil but mock repositories are disabled")
	}
	return postgresRepo.NewPostgresLoginAttemptRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateAssessmentRepository() repositories.AssessmentRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockAssessmentRepository()
	}
	if f.infra.DB == nil {
		panic("PostgreSQL DB connection is nil but mock repositories are disabled")
	}
	return postgresRepo.NewPostgresAssessmentRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateAttemptRepository() repositories.AttemptRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockAttemptRepository()
	}
	if f.infra.DB == nil {
		panic("PostgreSQL DB connection is nil but mock repositories are disabled")
	}
	return postgresRepo.NewPostgresAttemptRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateAnswerRepository() repositories.AnswerRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockAnswerRepository()
	}
	if f.infra.DB == nil {
		panic("PostgreSQL DB connection is nil but mock repositories are disabled")
	}
	return postgresRepo.NewPostgresAnswerRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateScreenRepository() repository.ScreenRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockScreenRepository()
	}
	if f.infra.DB == nil {
		panic("PostgreSQL DB connection is nil but mock repositories are disabled")
	}
	return postgresRepo.NewPostgresScreenRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateResourceReader() repository.ResourceReader {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockResourceReader()
	}
	if f.infra.DB == nil {
		panic("PostgreSQL DB connection is nil but mock repositories are disabled")
	}
	return postgresRepo.NewPostgresResourceRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateSummaryRepository() repository.SummaryRepository {
	if f.config.Development.UseMockRepositories {
		return mockMongo.NewMockSummaryRepository()
	}
	// Si MongoDB es opcional y no está disponible, usar mock como fallback
	if f.infra.MongoDB == nil {
		f.infra.Logger.Warn("MongoDB not available, using mock SummaryRepository as fallback")
		return mockMongo.NewMockSummaryRepository()
	}
	return mongoRepo.NewMongoSummaryRepository(f.infra.MongoDB)
}

func (f *RepositoryFactory) CreateAssessmentDocumentRepository() mongoRepo.AssessmentDocumentRepository {
	if f.config.Development.UseMockRepositories {
		return mockMongo.NewMockAssessmentDocumentRepository()
	}
	// Si MongoDB es opcional y no está disponible, usar mock como fallback
	if f.infra.MongoDB == nil {
		f.infra.Logger.Warn("MongoDB not available, using mock AssessmentDocumentRepository as fallback")
		return mockMongo.NewMockAssessmentDocumentRepository()
	}
	return mongoRepo.NewMongoAssessmentDocumentRepository(f.infra.MongoDB)
}
