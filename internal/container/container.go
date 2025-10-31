package container

import (
	"database/sql"

	"github.com/edugo/api-mobile/internal/application/service"
	"github.com/edugo/api-mobile/internal/domain/repository"
	"github.com/edugo/api-mobile/internal/infrastructure/http/handler"
	mongoRepo "github.com/edugo/api-mobile/internal/infrastructure/persistence/mongodb/repository"
	postgresRepo "github.com/edugo/api-mobile/internal/infrastructure/persistence/postgres/repository"
	"github.com/edugo/shared/pkg/auth"
	"github.com/edugo/shared/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

// Container es el contenedor de dependencias de API Mobile
// Implementa el patrón Dependency Injection
type Container struct {
	// Infrastructure
	DB         *sql.DB
	MongoDB    *mongo.Database
	Logger     logger.Logger
	JWTManager *auth.JWTManager

	// Repositories
	UserRepository       repository.UserRepository
	MaterialRepository   repository.MaterialRepository
	ProgressRepository   repository.ProgressRepository
	SummaryRepository    repository.SummaryRepository
	AssessmentRepository repository.AssessmentRepository

	// Services
	AuthService       service.AuthService
	MaterialService   service.MaterialService
	ProgressService   service.ProgressService
	SummaryService    service.SummaryService
	AssessmentService service.AssessmentService
	StatsService      service.StatsService

	// Handlers
	AuthHandler       *handler.AuthHandler
	MaterialHandler   *handler.MaterialHandler
	ProgressHandler   *handler.ProgressHandler
	SummaryHandler    *handler.SummaryHandler
	AssessmentHandler *handler.AssessmentHandler
	StatsHandler      *handler.StatsHandler
}

// NewContainer crea un nuevo contenedor e inicializa todas las dependencias
func NewContainer(db *sql.DB, mongoDB *mongo.Database, jwtSecret string, logger logger.Logger) *Container {
	c := &Container{
		DB:         db,
		MongoDB:    mongoDB,
		Logger:     logger,
		JWTManager: auth.NewJWTManager(jwtSecret, "edugo-mobile"),
	}

	// Inicializar repositories (capa de infraestructura)
	c.UserRepository = postgresRepo.NewPostgresUserRepository(db)
	c.MaterialRepository = postgresRepo.NewPostgresMaterialRepository(db)
	c.ProgressRepository = postgresRepo.NewPostgresProgressRepository(db)
	c.SummaryRepository = mongoRepo.NewMongoSummaryRepository(mongoDB)
	c.AssessmentRepository = mongoRepo.NewMongoAssessmentRepository(mongoDB)

	// Inicializar services (capa de aplicación)
	c.AuthService = service.NewAuthService(
		c.UserRepository,
		c.JWTManager,
		logger,
	)
	c.MaterialService = service.NewMaterialService(
		c.MaterialRepository,
		logger,
	)
	c.ProgressService = service.NewProgressService(
		c.ProgressRepository,
		logger,
	)
	c.SummaryService = service.NewSummaryService(
		c.SummaryRepository,
		logger,
	)
	c.AssessmentService = service.NewAssessmentService(
		c.AssessmentRepository,
		logger,
	)
	c.StatsService = service.NewStatsService(logger)

	// Inicializar handlers (capa de infraestructura HTTP)
	c.AuthHandler = handler.NewAuthHandler(
		c.AuthService,
		logger,
	)
	c.MaterialHandler = handler.NewMaterialHandler(
		c.MaterialService,
		logger,
	)
	c.ProgressHandler = handler.NewProgressHandler(
		c.ProgressService,
		logger,
	)
	c.SummaryHandler = handler.NewSummaryHandler(
		c.SummaryService,
		logger,
	)
	c.AssessmentHandler = handler.NewAssessmentHandler(
		c.AssessmentService,
		logger,
	)
	c.StatsHandler = handler.NewStatsHandler(
		c.StatsService,
		logger,
	)

	return c
}

// Close cierra los recursos del contenedor
func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
