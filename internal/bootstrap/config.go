package bootstrap

import (
	"database/sql"

	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
	"github.com/EduGoGroup/edugo-shared/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

// ResourceConfig define la configuración de un recurso de infraestructura
type ResourceConfig struct {
	Name     string // Nombre del recurso (e.g., "postgresql", "rabbitmq")
	Optional bool   // Si es true, el recurso puede fallar sin detener la inicialización
	Enabled  bool   // Si es false, el recurso no se inicializa
}

// DefaultResourceConfig retorna la configuración por defecto de recursos
// Por defecto, Logger, PostgreSQL y MongoDB son requeridos
// RabbitMQ y S3 son opcionales para facilitar desarrollo local
func DefaultResourceConfig() map[string]ResourceConfig {
	return map[string]ResourceConfig{
		"logger":     {Name: "logger", Optional: false, Enabled: true},
		"postgresql": {Name: "postgresql", Optional: false, Enabled: true},
		"mongodb":    {Name: "mongodb", Optional: false, Enabled: true},
		"rabbitmq":   {Name: "rabbitmq", Optional: true, Enabled: true},
		"s3":         {Name: "s3", Optional: true, Enabled: true},
	}
}

// BootstrapOption es una función que modifica BootstrapOptions
// Permite configurar el bootstrap usando el patrón de opciones funcionales
type BootstrapOption func(*BootstrapOptions)

// WithLogger inyecta un logger pre-construido
// Útil para testing o cuando se quiere usar un logger específico
func WithLogger(log logger.Logger) BootstrapOption {
	return func(opts *BootstrapOptions) {
		opts.Logger = log
	}
}

// WithPostgreSQL inyecta una conexión PostgreSQL pre-construida
// Útil para testing con bases de datos mock o en memoria
func WithPostgreSQL(db *sql.DB) BootstrapOption {
	return func(opts *BootstrapOptions) {
		opts.PostgreSQL = db
	}
}

// WithMongoDB inyecta una conexión MongoDB pre-construida
// Útil para testing con bases de datos mock o en memoria
func WithMongoDB(db *mongo.Database) BootstrapOption {
	return func(opts *BootstrapOptions) {
		opts.MongoDB = db
	}
}

// WithRabbitMQ inyecta un publisher RabbitMQ pre-construido
// Útil para testing con publishers mock
func WithRabbitMQ(pub rabbitmq.Publisher) BootstrapOption {
	return func(opts *BootstrapOptions) {
		opts.RabbitMQPublisher = pub
	}
}

// WithS3Client inyecta un cliente S3 pre-construido
// Útil para testing con clientes S3 mock
func WithS3Client(client S3Storage) BootstrapOption {
	return func(opts *BootstrapOptions) {
		opts.S3Client = client
	}
}

// WithOptionalResource marca un recurso como opcional
// Si el recurso falla al inicializar, se registra una advertencia
// y se continúa con una implementación noop
func WithOptionalResource(resourceName string) BootstrapOption {
	return func(opts *BootstrapOptions) {
		if opts.OptionalResources == nil {
			opts.OptionalResources = make(map[string]bool)
		}
		opts.OptionalResources[resourceName] = true
	}
}

// WithDisabledResource deshabilita completamente un recurso
// El recurso no se inicializa y se usa una implementación noop
// Recursos soportados: "rabbitmq", "s3"
// Nota: "postgresql", "mongodb" y "logger" no pueden ser deshabilitados
func WithDisabledResource(resourceName string) BootstrapOption {
	return func(opts *BootstrapOptions) {
		if opts.DisabledResources == nil {
			opts.DisabledResources = make(map[string]bool)
		}
		opts.DisabledResources[resourceName] = true
	}
}

// IsResourceDisabled verifica si un recurso está deshabilitado
func (opts *BootstrapOptions) IsResourceDisabled(resourceName string) bool {
	if opts == nil || opts.DisabledResources == nil {
		return false
	}
	return opts.DisabledResources[resourceName]
}

// WithLoggerFactory inyecta una factory personalizada para crear loggers
// Útil para testing avanzado
func WithLoggerFactory(factory LoggerFactory) BootstrapOption {
	return func(opts *BootstrapOptions) {
		opts.LoggerFactory = factory
	}
}

// WithDatabaseFactory inyecta una factory personalizada para crear bases de datos
// Útil para testing avanzado
func WithDatabaseFactory(factory DatabaseFactory) BootstrapOption {
	return func(opts *BootstrapOptions) {
		opts.DatabaseFactory = factory
	}
}

// WithMessagingFactory inyecta una factory personalizada para crear clientes de mensajería
// Útil para testing avanzado
func WithMessagingFactory(factory MessagingFactory) BootstrapOption {
	return func(opts *BootstrapOptions) {
		opts.MessagingFactory = factory
	}
}

// WithStorageFactory inyecta una factory personalizada para crear clientes de almacenamiento
// Útil para testing avanzado
func WithStorageFactory(factory StorageFactory) BootstrapOption {
	return func(opts *BootstrapOptions) {
		opts.StorageFactory = factory
	}
}
