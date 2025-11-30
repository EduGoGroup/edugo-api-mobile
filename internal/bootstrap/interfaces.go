package bootstrap

import (
	"context"
	"database/sql"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/storage/s3"
	"github.com/EduGoGroup/edugo-shared/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

// LoggerFactory crea instancias de logger
type LoggerFactory interface {
	Create(level, format string) (logger.Logger, error)
}

// DatabaseFactory crea conexiones a bases de datos
type DatabaseFactory interface {
	CreatePostgreSQL(ctx context.Context, cfg *config.Config, log logger.Logger) (*sql.DB, error)
	CreateMongoDB(ctx context.Context, cfg *config.Config, log logger.Logger) (*mongo.Database, error)
}

// MessagingFactory crea clientes de mensajería
type MessagingFactory interface {
	CreatePublisher(url, exchange string, log logger.Logger) (rabbitmq.Publisher, error)
}

// StorageFactory crea clientes de almacenamiento
type StorageFactory interface {
	CreateS3Client(ctx context.Context, cfg s3.S3Config, log logger.Logger) (S3Storage, error)
}

// S3Storage define las operaciones de almacenamiento en S3
// Esta interfaz es reutilizada desde s3/interface.go para mantener consistencia
type S3Storage interface {
	// GeneratePresignedUploadURL genera una URL presignada para subir archivos
	GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expires time.Duration) (string, error)

	// GeneratePresignedDownloadURL genera una URL presignada para descargar archivos
	GeneratePresignedDownloadURL(ctx context.Context, key string, expires time.Duration) (string, error)
}

// Resources encapsula todos los recursos de infraestructura inicializados
// Esta estructura se pasa al container para inyección de dependencias
type Resources struct {
	Logger            logger.Logger
	PostgreSQL        *sql.DB
	MongoDB           *mongo.Database
	RabbitMQPublisher rabbitmq.Publisher
	S3Client          S3Storage
	JWTSecret         string
	AuthConfig        config.AuthConfig // Configuración de autenticación (api-admin)
	Config            *config.Config    // Configuración completa de la aplicación
}

// BootstrapOptions permite configurar el comportamiento del bootstrap
// Soporta inyección de recursos pre-construidos (para testing) y
// configuración de recursos opcionales
type BootstrapOptions struct {
	// Recursos pre-construidos (para testing)
	Logger            logger.Logger
	PostgreSQL        *sql.DB
	MongoDB           *mongo.Database
	RabbitMQPublisher rabbitmq.Publisher
	S3Client          S3Storage

	// Configuración de recursos opcionales
	OptionalResources map[string]bool

	// Factories personalizadas (para testing avanzado)
	LoggerFactory    LoggerFactory
	DatabaseFactory  DatabaseFactory
	MessagingFactory MessagingFactory
	StorageFactory   StorageFactory
}
