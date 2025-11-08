package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/bootstrap/noop"
	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/storage/s3"
	"github.com/EduGoGroup/edugo-shared/logger"
	"go.uber.org/zap"
)

// Bootstrapper orquesta la inicialización de todos los recursos de infraestructura
// Aplica el patrón Factory con soporte para recursos opcionales, inyección de mocks,
// y gestión del ciclo de vida
type Bootstrapper struct {
	config    *config.Config
	options   *BootstrapOptions
	resources *Resources
	lifecycle *LifecycleManager
	logger    logger.Logger
}

// New crea un nuevo Bootstrapper con la configuración y opciones proporcionadas
// Aplica el patrón de opciones funcionales para permitir configuración flexible
func New(cfg *config.Config, opts ...BootstrapOption) *Bootstrapper {
	// Inicializar opciones con valores por defecto
	options := &BootstrapOptions{
		OptionalResources: make(map[string]bool),
	}

	// Cargar configuración de recursos opcionales desde config
	// Esto permite configurar recursos opcionales via archivo YAML o ENV vars
	options.OptionalResources["rabbitmq"] = cfg.Bootstrap.OptionalResources.RabbitMQ
	options.OptionalResources["s3"] = cfg.Bootstrap.OptionalResources.S3

	// Aplicar opciones funcionales (pueden sobrescribir la configuración)
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
// Los recursos opcionales que fallan se reemplazan con implementaciones noop
func (b *Bootstrapper) InitializeInfrastructure(ctx context.Context) (*Resources, func() error, error) {
	startTime := time.Now()

	// Inicializar logger primero (necesario para logging de otros recursos)
	if err := b.initializeLogger(); err != nil {
		return nil, nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	b.logger.Info("starting infrastructure initialization")

	// Inicializar lifecycle manager
	b.lifecycle = NewLifecycleManager(b.logger)

	// Inicializar recursos struct
	b.resources = &Resources{
		Logger:    b.logger,
		JWTSecret: b.config.Auth.JWT.Secret,
	}

	// Inicializar PostgreSQL (recurso requerido)
	if err := b.initializePostgreSQL(ctx); err != nil {
		b.logger.Error("failed to initialize PostgreSQL", zap.Error(err))
		// Ejecutar cleanup parcial antes de retornar error
		_ = b.lifecycle.Cleanup()
		return nil, nil, fmt.Errorf("failed to initialize required resource PostgreSQL: %w", err)
	}

	// Inicializar MongoDB (recurso requerido)
	if err := b.initializeMongoDB(ctx); err != nil {
		b.logger.Error("failed to initialize MongoDB", zap.Error(err))
		// Ejecutar cleanup parcial antes de retornar error
		_ = b.lifecycle.Cleanup()
		return nil, nil, fmt.Errorf("failed to initialize required resource MongoDB: %w", err)
	}

	// Inicializar RabbitMQ (recurso opcional)
	if err := b.initializeRabbitMQ(); err != nil {
		if b.isOptional("rabbitmq") {
			b.logger.Warn("RabbitMQ initialization failed, using noop implementation",
				zap.Error(err),
			)
			b.resources.RabbitMQPublisher = noop.NewNoopPublisher(b.logger)
		} else {
			b.logger.Error("failed to initialize RabbitMQ", zap.Error(err))
			_ = b.lifecycle.Cleanup()
			return nil, nil, fmt.Errorf("failed to initialize required resource RabbitMQ: %w", err)
		}
	}

	// Inicializar S3 (recurso opcional)
	if err := b.initializeS3(ctx); err != nil {
		if b.isOptional("s3") {
			b.logger.Warn("S3 initialization failed, using noop implementation",
				zap.Error(err),
			)
			b.resources.S3Client = noop.NewNoopS3Storage(b.logger)
		} else {
			b.logger.Error("failed to initialize S3", zap.Error(err))
			_ = b.lifecycle.Cleanup()
			return nil, nil, fmt.Errorf("failed to initialize required resource S3: %w", err)
		}
	}

	duration := time.Since(startTime)
	b.logger.Info("infrastructure initialization completed",
		zap.Duration("duration", duration),
	)

	// Retornar recursos y función de cleanup
	cleanup := func() error {
		return b.lifecycle.Cleanup()
	}

	return b.resources, cleanup, nil
}

// initializeLogger inicializa el logger
// Usa el logger inyectado si está disponible, de lo contrario crea uno nuevo
func (b *Bootstrapper) initializeLogger() error {
	// Si ya hay un logger inyectado, usarlo
	if b.options.Logger != nil {
		b.logger = b.options.Logger
		b.logger.Info("using injected logger")
		return nil
	}

	// Crear factory (usar inyectada o default)
	factory := b.options.LoggerFactory
	if factory == nil {
		factory = NewDefaultFactories()
	}

	// Crear logger
	log, err := factory.Create(b.config.Logging.Level, b.config.Logging.Format)
	if err != nil {
		return fmt.Errorf("factory failed to create logger: %w", err)
	}

	b.logger = log
	b.logger.Info("logger initialized successfully",
		zap.String("level", b.config.Logging.Level),
		zap.String("format", b.config.Logging.Format),
	)

	return nil
}

// initializePostgreSQL inicializa la conexión a PostgreSQL
// Usa la conexión inyectada si está disponible, de lo contrario crea una nueva
func (b *Bootstrapper) initializePostgreSQL(ctx context.Context) error {
	b.logger.Info("initializing PostgreSQL",
		zap.String("host", b.config.Database.Postgres.Host),
		zap.Int("port", b.config.Database.Postgres.Port),
		zap.String("database", b.config.Database.Postgres.Database),
	)

	// Si ya hay una conexión inyectada, usarla
	if b.options.PostgreSQL != nil {
		b.resources.PostgreSQL = b.options.PostgreSQL
		b.logger.Info("using injected PostgreSQL connection")
		return nil
	}

	// Crear factory (usar inyectada o default)
	factory := b.options.DatabaseFactory
	if factory == nil {
		factory = NewDefaultFactories()
	}

	// Crear conexión
	db, err := factory.CreatePostgreSQL(ctx, b.config, b.logger)
	if err != nil {
		return fmt.Errorf("factory failed to create PostgreSQL connection: %w", err)
	}

	b.resources.PostgreSQL = db

	// Registrar cleanup
	b.lifecycle.Register("postgresql", func() error {
		b.logger.Info("closing PostgreSQL connection")
		return db.Close()
	})

	b.logger.Info("PostgreSQL initialized successfully",
		zap.String("host", b.config.Database.Postgres.Host),
		zap.Int("port", b.config.Database.Postgres.Port),
	)

	return nil
}

// initializeMongoDB inicializa la conexión a MongoDB
// Usa la conexión inyectada si está disponible, de lo contrario crea una nueva
func (b *Bootstrapper) initializeMongoDB(ctx context.Context) error {
	b.logger.Info("initializing MongoDB",
		zap.String("database", b.config.Database.MongoDB.Database),
	)

	// Si ya hay una conexión inyectada, usarla
	if b.options.MongoDB != nil {
		b.resources.MongoDB = b.options.MongoDB
		b.logger.Info("using injected MongoDB connection")
		return nil
	}

	// Crear factory (usar inyectada o default)
	factory := b.options.DatabaseFactory
	if factory == nil {
		factory = NewDefaultFactories()
	}

	// Crear conexión
	db, err := factory.CreateMongoDB(ctx, b.config, b.logger)
	if err != nil {
		return fmt.Errorf("factory failed to create MongoDB connection: %w", err)
	}

	b.resources.MongoDB = db

	// Registrar cleanup
	b.lifecycle.Register("mongodb", func() error {
		b.logger.Info("closing MongoDB connection")
		return db.Client().Disconnect(context.Background())
	})

	b.logger.Info("MongoDB initialized successfully",
		zap.String("database", b.config.Database.MongoDB.Database),
	)

	return nil
}

// initializeRabbitMQ inicializa el publisher de RabbitMQ
// Usa el publisher inyectado si está disponible, de lo contrario crea uno nuevo
func (b *Bootstrapper) initializeRabbitMQ() error {
	b.logger.Info("initializing RabbitMQ")

	// Si ya hay un publisher inyectado, usarlo
	if b.options.RabbitMQPublisher != nil {
		b.resources.RabbitMQPublisher = b.options.RabbitMQPublisher
		b.logger.Info("using injected RabbitMQ publisher")
		return nil
	}

	// Crear factory (usar inyectada o default)
	factory := b.options.MessagingFactory
	if factory == nil {
		factory = NewDefaultFactories()
	}

	// Crear publisher
	pub, err := factory.CreatePublisher(
		b.config.Messaging.RabbitMQ.URL,
		b.config.Messaging.RabbitMQ.Exchanges.Materials,
		b.logger,
	)
	if err != nil {
		return fmt.Errorf("factory failed to create RabbitMQ publisher: %w", err)
	}

	b.resources.RabbitMQPublisher = pub

	// Registrar cleanup
	b.lifecycle.Register("rabbitmq", func() error {
		b.logger.Info("closing RabbitMQ connection")
		return pub.Close()
	})

	b.logger.Info("RabbitMQ initialized successfully")

	return nil
}

// initializeS3 inicializa el cliente S3
// Usa el cliente inyectado si está disponible, de lo contrario crea uno nuevo
func (b *Bootstrapper) initializeS3(ctx context.Context) error {
	b.logger.Info("initializing S3",
		zap.String("region", b.config.Storage.S3.Region),
		zap.String("bucket", b.config.Storage.S3.BucketName),
	)

	// Si ya hay un cliente inyectado, usarlo
	if b.options.S3Client != nil {
		b.resources.S3Client = b.options.S3Client
		b.logger.Info("using injected S3 client")
		return nil
	}

	// Crear factory (usar inyectada o default)
	factory := b.options.StorageFactory
	if factory == nil {
		factory = NewDefaultFactories()
	}

	// Crear configuración S3
	s3Config := s3.S3Config{
		Region:          b.config.Storage.S3.Region,
		BucketName:      b.config.Storage.S3.BucketName,
		AccessKeyID:     b.config.Storage.S3.AccessKeyID,
		SecretAccessKey: b.config.Storage.S3.SecretAccessKey,
		Endpoint:        b.config.Storage.S3.Endpoint,
	}

	// Crear cliente
	client, err := factory.CreateS3Client(ctx, s3Config, b.logger)
	if err != nil {
		return fmt.Errorf("factory failed to create S3 client: %w", err)
	}

	b.resources.S3Client = client

	// S3 client no requiere cleanup explícito
	b.logger.Info("S3 initialized successfully",
		zap.String("region", b.config.Storage.S3.Region),
		zap.String("bucket", b.config.Storage.S3.BucketName),
	)

	return nil
}

// isOptional verifica si un recurso está marcado como opcional
func (b *Bootstrapper) isOptional(resourceName string) bool {
	// Verificar en opciones personalizadas
	if optional, exists := b.options.OptionalResources[resourceName]; exists {
		return optional
	}

	// Verificar en configuración por defecto
	defaultConfig := DefaultResourceConfig()
	if rc, exists := defaultConfig[resourceName]; exists {
		return rc.Optional
	}

	// Por defecto, los recursos son requeridos
	return false
}
