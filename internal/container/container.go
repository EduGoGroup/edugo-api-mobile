package container

import (
	"database/sql"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/service"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/handler"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
	mongoRepo "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/mongodb/repository"
	postgresRepo "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/postgres/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/storage/s3"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

// Container es el contenedor de dependencias de API Mobile
// Implementa el patrón Dependency Injection
type Container struct {
	// Infrastructure
	DB               *sql.DB
	MongoDB          *mongo.Database
	Logger           logger.Logger
	JWTManager       *auth.JWTManager
	MessagePublisher rabbitmq.Publisher
	S3Client         *s3.S3Client

	// Repositories
	UserRepository         repository.UserRepository
	MaterialRepository     repository.MaterialRepository
	ProgressRepository     repository.ProgressRepository
	SummaryRepository      repository.SummaryRepository
	AssessmentRepository   repository.AssessmentRepository
	RefreshTokenRepository repository.RefreshTokenRepository
	LoginAttemptRepository repository.LoginAttemptRepository

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
func NewContainer(db *sql.DB, mongoDB *mongo.Database, publisher rabbitmq.Publisher, s3Client *s3.S3Client, jwtSecret string, logger logger.Logger) *Container {
	c := &Container{
		DB:               db,
		MongoDB:          mongoDB,
		MessagePublisher: publisher,
		S3Client:         s3Client,
		Logger:           logger,
		JWTManager:       auth.NewJWTManager(jwtSecret, "edugo-mobile"),
	}

	// Inicializar repositories (capa de infraestructura)
	c.UserRepository = postgresRepo.NewPostgresUserRepository(db)
	c.MaterialRepository = postgresRepo.NewPostgresMaterialRepository(db)
	c.ProgressRepository = postgresRepo.NewPostgresProgressRepository(db)
	c.RefreshTokenRepository = postgresRepo.NewPostgresRefreshTokenRepository(db)
	c.LoginAttemptRepository = postgresRepo.NewPostgresLoginAttemptRepository(db)
	c.SummaryRepository = mongoRepo.NewMongoSummaryRepository(mongoDB)
	c.AssessmentRepository = mongoRepo.NewMongoAssessmentRepository(mongoDB)

	// Inicializar services (capa de aplicación)
	c.AuthService = service.NewAuthService(
		c.UserRepository,
		c.RefreshTokenRepository,
		c.LoginAttemptRepository,
		c.JWTManager,
		logger,
	)
	c.MaterialService = service.NewMaterialService(
		c.MaterialRepository,
		c.MessagePublisher,
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
		c.MessagePublisher,
		logger,
	)
	c.StatsService = service.NewStatsService(
		logger,
		c.MaterialRepository,
		c.AssessmentRepository,
		c.ProgressRepository,
	)

	// Inicializar handlers (capa de infraestructura HTTP)
	c.AuthHandler = handler.NewAuthHandler(
		c.AuthService,
		logger,
	)
	c.MaterialHandler = handler.NewMaterialHandler(
		c.MaterialService,
		c.S3Client,
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
