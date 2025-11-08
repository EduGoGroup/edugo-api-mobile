package bootstrap

import (
	"context"
	"database/sql"

	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/database"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/storage/s3"
	"github.com/EduGoGroup/edugo-shared/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

// DefaultFactories implementa todas las factories con implementaciones reales
// Delega la creación de recursos a las implementaciones existentes en internal/infrastructure
type DefaultFactories struct{}

// NewDefaultFactories crea una nueva instancia de DefaultFactories
func NewDefaultFactories() *DefaultFactories {
	return &DefaultFactories{}
}

// Create crea una instancia de logger usando la implementación de edugo-shared
func (f *DefaultFactories) Create(level, format string) (logger.Logger, error) {
	// logger.NewZapLogger no retorna error, siempre tiene éxito
	log := logger.NewZapLogger(level, format)
	return log, nil
}

// CreatePostgreSQL crea una conexión a PostgreSQL usando la implementación existente
func (f *DefaultFactories) CreatePostgreSQL(ctx context.Context, cfg *config.Config, log logger.Logger) (*sql.DB, error) {
	return database.InitPostgreSQL(ctx, cfg, log)
}

// CreateMongoDB crea una conexión a MongoDB usando la implementación existente
func (f *DefaultFactories) CreateMongoDB(ctx context.Context, cfg *config.Config, log logger.Logger) (*mongo.Database, error) {
	return database.InitMongoDB(ctx, cfg, log)
}

// CreatePublisher crea un publisher de RabbitMQ usando la implementación existente
func (f *DefaultFactories) CreatePublisher(url, exchange string, log logger.Logger) (rabbitmq.Publisher, error) {
	return rabbitmq.NewRabbitMQPublisher(url, exchange, log)
}

// CreateS3Client crea un cliente S3 usando la implementación existente
func (f *DefaultFactories) CreateS3Client(ctx context.Context, cfg s3.S3Config, log logger.Logger) (S3Storage, error) {
	return s3.NewS3Client(ctx, cfg, log)
}

// Verificar que DefaultFactories implementa todas las interfaces de factory
var (
	_ LoggerFactory    = (*DefaultFactories)(nil)
	_ DatabaseFactory  = (*DefaultFactories)(nil)
	_ MessagingFactory = (*DefaultFactories)(nil)
	_ StorageFactory   = (*DefaultFactories)(nil)
)
