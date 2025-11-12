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
// REFACTORIZADO: Ahora delega a shared/bootstrap vía bridgeToSharedBootstrap()
func (b *Bootstrapper) InitializeInfrastructure(ctx context.Context) (*Resources, func() error, error) {
	startTime := time.Now()

	// Llamar al bridge que se conecta con shared/bootstrap
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
