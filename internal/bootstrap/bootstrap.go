package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
	"github.com/EduGoGroup/edugo-shared/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

// Bootstrapper orquesta la inicialización de todos los recursos de infraestructura
// REFACTORIZADO: Ahora usa shared/bootstrap internamente vía bridge.go
type Bootstrapper struct {
	config  *config.Config
	options *BootstrapOptions
}

// New crea un nuevo Bootstrapper con la configuración y opciones proporcionadas
// Aplica el patrón de opciones funcionales para permitir configuración flexible
func New(cfg *config.Config, opts ...BootstrapOption) *Bootstrapper {
	// Inicializar opciones con valores por defecto
	options := &BootstrapOptions{
		OptionalResources: make(map[string]bool),
	}

	// Cargar configuración de recursos opcionales desde config
	options.OptionalResources["rabbitmq"] = cfg.Bootstrap.OptionalResources.RabbitMQ
	options.OptionalResources["s3"] = cfg.Bootstrap.OptionalResources.S3

	// Aplicar opciones funcionales
	for _, opt := range opts {
		opt(options)
	}

	return &Bootstrapper{
		config:  cfg,
		options: options,
	}
}

// InitializeInfrastructure inicializa todos los recursos de infraestructura
// Retorna los recursos inicializados, una función de cleanup, y un error si algo falla
//
// MODO MOCK: Si development.use_mock_repositories=true, salta la inicialización de DB
// y retorna recursos mínimos (solo logger). Perfecto para desarrollo frontend sin Docker.
//
// MODO REAL: Delega a shared/bootstrap vía bridgeToSharedBootstrap()
func (b *Bootstrapper) InitializeInfrastructure(ctx context.Context) (*Resources, func() error, error) {
	startTime := time.Now()

	// NUEVO: Verificar si estamos en modo mock
	if b.config.Development.UseMockRepositories {
		return b.initializeMockMode(ctx, startTime)
	}

	// Flujo normal: Llamar al bridge que se conecta con shared/bootstrap
	resources, lifecycleManager, err := bridgeToSharedBootstrap(ctx, b.config, b.options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize infrastructure via shared/bootstrap: %w", err)
	}

	// Log de completado
	duration := time.Since(startTime)
	resources.Logger.Info("infrastructure initialization completed",
		"duration", duration.String(),
	)

	// Crear función de cleanup que usa el lifecycle manager de shared
	cleanup := func() error {
		resources.Logger.Info("starting infrastructure cleanup")

		// shared/lifecycle.Cleanup() no toma contexto
		err := lifecycleManager.Cleanup()
		if err != nil {
			resources.Logger.Error("infrastructure cleanup completed with errors",
				"error", err.Error(),
			)
			return err
		}

		resources.Logger.Info("infrastructure cleanup completed successfully")
		return nil
	}

	return resources, cleanup, nil
}

// initializeMockMode inicializa recursos mínimos para modo mock (sin bases de datos)
// Solo crea logger y retorna recursos nil para DB/Messaging/Storage
func (b *Bootstrapper) initializeMockMode(ctx context.Context, startTime time.Time) (*Resources, func() error, error) {
	// Crear logger básico
	log := logger.NewZapLogger(b.config.Logging.Level, b.config.Logging.Format)

	log.Info("🎭 MODO MOCK ACTIVADO - Iniciando sin conexiones de base de datos",
		"use_mock_repositories", true,
		"ram_saving", "~4GB → ~200MB",
		"startup_improvement", "~30s → ~1.5s",
	)

	// Crear recursos mínimos (sin DB, sin messaging, sin storage)
	resources := &Resources{
		Logger:            log,
		PostgreSQL:        nil, // No inicializar PostgreSQL
		MongoDB:           nil, // No inicializar MongoDB
		RabbitMQPublisher: nil, // No inicializar RabbitMQ
		S3Client:          nil, // No inicializar S3
		JWTSecret:         b.config.Auth.JWT.Secret,
		AuthConfig:        b.config.Auth,
		Config:            b.config, // Configuración completa para factory
	}

	duration := time.Since(startTime)
	log.Info("mock mode initialization completed",
		"duration", duration.String(),
		"mode", "mock",
	)

	// Cleanup no-op para modo mock
	cleanup := func() error {
		log.Info("mock mode cleanup (no-op)")
		return nil
	}

	return resources, cleanup, nil
}

// Mantener compatibilidad con métodos inyectados para tests (deprecated)
// Estos métodos son wrappers sobre las opciones funcionales existentes

// WithInjectedLogger permite inyectar un logger pre-construido (para tests)
// DEPRECATED: Usar WithLogger() directamente
func WithInjectedLogger(log logger.Logger) BootstrapOption {
	return WithLogger(log)
}

// WithInjectedPostgreSQL permite inyectar una conexión PostgreSQL (para tests)
// DEPRECATED: Usar WithPostgreSQL() directamente
func WithInjectedPostgreSQL(db *sql.DB) BootstrapOption {
	return WithPostgreSQL(db)
}

// WithInjectedMongoDB permite inyectar una conexión MongoDB (para tests)
// DEPRECATED: Usar WithMongoDB() directamente
func WithInjectedMongoDB(db *mongo.Database) BootstrapOption {
	return WithMongoDB(db)
}

// WithInjectedRabbitMQ permite inyectar un publisher RabbitMQ (para tests)
// DEPRECATED: Usar WithRabbitMQ() directamente
func WithInjectedRabbitMQ(pub rabbitmq.Publisher) BootstrapOption {
	return WithRabbitMQ(pub)
}

// WithInjectedS3Client permite inyectar un cliente S3 (para tests)
// DEPRECATED: Usar WithS3Client() directamente
func WithInjectedS3Client(client S3Storage) BootstrapOption {
	return WithS3Client(client)
}
