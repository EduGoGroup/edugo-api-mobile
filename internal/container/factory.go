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

type RepositoryFactory struct {
	config *config.Config
	infra  *InfrastructureContainer
}

func NewRepositoryFactory(cfg *config.Config, infra *InfrastructureContainer) *RepositoryFactory {
	return &RepositoryFactory{config: cfg, infra: infra}
}

func (f *RepositoryFactory) CreateUserRepository() repository.UserRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockUserRepository()
	}
	return postgresRepo.NewPostgresUserRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateMaterialRepository() repository.MaterialRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockMaterialRepository()
	}
	return postgresRepo.NewPostgresMaterialRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateProgressRepository() repository.ProgressRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockProgressRepository()
	}
	return postgresRepo.NewPostgresProgressRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateRefreshTokenRepository() repository.RefreshTokenRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockRefreshTokenRepository()
	}
	return postgresRepo.NewPostgresRefreshTokenRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateLoginAttemptRepository() repository.LoginAttemptRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockLoginAttemptRepository()
	}
	return postgresRepo.NewPostgresLoginAttemptRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateAssessmentRepository() repositories.AssessmentRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockAssessmentRepository()
	}
	return postgresRepo.NewPostgresAssessmentRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateAttemptRepository() repositories.AttemptRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockAttemptRepository()
	}
	return postgresRepo.NewPostgresAttemptRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateAnswerRepository() repositories.AnswerRepository {
	if f.config.Development.UseMockRepositories {
		return mockPostgres.NewMockAnswerRepository()
	}
	return postgresRepo.NewPostgresAnswerRepository(f.infra.DB)
}

func (f *RepositoryFactory) CreateSummaryRepository() repository.SummaryRepository {
	if f.config.Development.UseMockRepositories {
		return mockMongo.NewMockSummaryRepository()
	}
	return mongoRepo.NewMongoSummaryRepository(f.infra.MongoDB)
}

func (f *RepositoryFactory) CreateLegacyAssessmentRepository() repository.AssessmentRepository {
	if f.config.Development.UseMockRepositories {
		return mockMongo.NewMockLegacyAssessmentRepository()
	}
	return mongoRepo.NewMongoAssessmentRepository(f.infra.MongoDB)
}

func (f *RepositoryFactory) CreateAssessmentDocumentRepository() mongoRepo.AssessmentDocumentRepository {
	if f.config.Development.UseMockRepositories {
		return mockMongo.NewMockAssessmentDocumentRepository()
	}
	return mongoRepo.NewMongoAssessmentDocumentRepository(f.infra.MongoDB)
}
