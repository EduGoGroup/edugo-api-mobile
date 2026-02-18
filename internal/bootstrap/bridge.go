package bootstrap

import (
	"context"
	"fmt"

	"github.com/EduGoGroup/edugo-api-mobile/internal/bootstrap/adapter"
	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
	sharedBootstrap "github.com/EduGoGroup/edugo-shared/bootstrap"
	"github.com/EduGoGroup/edugo-shared/lifecycle"
	sharedLogger "github.com/EduGoGroup/edugo-shared/logger"
	"github.com/sony/gobreaker"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm/logger"
)

// bridgeToSharedBootstrap es el puente entre shared/bootstrap y la API de api-mobile
// Convierte la configuración de api-mobile al formato de shared, llama shared/bootstrap,
// y adapta los recursos retornados usando los adapters
func bridgeToSharedBootstrap(
	ctx context.Context,
	cfg *config.Config,
	opts *BootstrapOptions,
) (*Resources, *lifecycle.Manager, error) {
	// Si hay recursos inyectados (mocks), retornarlos directamente sin llamar a shared/bootstrap
	if opts != nil && opts.Logger != nil && opts.PostgreSQL != nil && opts.MongoDB != nil {
		resources := &Resources{
			Logger:            opts.Logger,
			PostgreSQL:        opts.PostgreSQL,
			MongoDB:           opts.MongoDB,
			RabbitMQPublisher: opts.RabbitMQPublisher,
			S3Client:          opts.S3Client,
			JWTSecret:         cfg.Auth.JWT.Secret,
			AuthConfig:        cfg.Auth,
			Config:            cfg,
		}
		// Crear lifecycle manager vacío para mocks
		lifecycleManager := lifecycle.NewManager(opts.Logger)
		return resources, lifecycleManager, nil
	}

	// 1. Crear configuraciones para cada componente de shared/bootstrap
	postgresConfig := sharedBootstrap.PostgreSQLConfig{
		Host:     cfg.Database.Postgres.Host,
		Port:     cfg.Database.Postgres.Port,
		User:     cfg.Database.Postgres.User,
		Password: cfg.Database.Postgres.Password,
		Database: cfg.Database.Postgres.Database,
		SSLMode:  cfg.Database.Postgres.SSLMode,
	}

	mongoConfig := sharedBootstrap.MongoDBConfig{
		URI:      cfg.Database.MongoDB.URI,
		Database: cfg.Database.MongoDB.Database,
	}

	rabbitMQConfig := sharedBootstrap.RabbitMQConfig{
		URL: cfg.Messaging.RabbitMQ.URL,
	}

	s3Config := sharedBootstrap.S3Config{
		Bucket:          cfg.Storage.S3.BucketName,
		Region:          cfg.Storage.S3.Region,
		AccessKeyID:     cfg.Storage.S3.AccessKeyID,
		SecretAccessKey: cfg.Storage.S3.SecretAccessKey,
	}

	// 2. Crear struct con todas las configs (será pasado como interface{})
	bootstrapConfig := struct {
		Environment string
		PostgreSQL  sharedBootstrap.PostgreSQLConfig
		MongoDB     sharedBootstrap.MongoDBConfig
		RabbitMQ    sharedBootstrap.RabbitMQConfig
		S3          sharedBootstrap.S3Config
	}{
		Environment: cfg.Environment,
		PostgreSQL:  postgresConfig,
		MongoDB:     mongoConfig,
		RabbitMQ:    rabbitMQConfig,
		S3:          s3Config,
	}

	// 3. Crear wrapper de factories personalizado que retiene tipos concretos
	// IMPORTANTE: PostgreSQLFactory requiere un logger de GORM
	// Configuramos el nivel de logging de GORM según el entorno
	gormLogLevel := logger.Silent
	if cfg.Environment == "development" || cfg.Environment == "local" {
		gormLogLevel = logger.Info
	}
	gormLogger := logger.Default.LogMode(gormLogLevel)

	sharedFactories := &sharedBootstrap.Factories{
		Logger:     sharedBootstrap.NewDefaultLoggerFactory(),
		PostgreSQL: sharedBootstrap.NewDefaultPostgreSQLFactory(gormLogger),
		MongoDB:    sharedBootstrap.NewDefaultMongoDBFactory(),
		RabbitMQ:   sharedBootstrap.NewDefaultRabbitMQFactory(),
		S3:         sharedBootstrap.NewDefaultS3Factory(),
	}

	wrapper := newCustomFactoriesWrapper(sharedFactories)
	customFactories := createCustomFactories(wrapper)

	// 4. Crear lifecycle manager de shared (sin logger por ahora, lo configuraremos después)
	lifecycleManager := lifecycle.NewManager(opts.Logger)

	// 5. Configurar opciones de bootstrap de shared
	var bootstrapOpts []sharedBootstrap.BootstrapOption

	// Recursos opcionales
	if cfg.Bootstrap.OptionalResources.RabbitMQ {
		bootstrapOpts = append(bootstrapOpts, sharedBootstrap.WithOptionalResources("rabbitmq"))
	}
	if cfg.Bootstrap.OptionalResources.S3 {
		bootstrapOpts = append(bootstrapOpts, sharedBootstrap.WithOptionalResources("s3"))
	}
	if cfg.Bootstrap.OptionalResources.MongoDB {
		bootstrapOpts = append(bootstrapOpts, sharedBootstrap.WithOptionalResources("mongodb"))
	}

	// 6. Llamar shared/bootstrap
	_, err := sharedBootstrap.Bootstrap(
		ctx,
		bootstrapConfig,
		customFactories,
		lifecycleManager,
		bootstrapOpts...,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("shared bootstrap failed: %w", err)
	}

	// 7. Adaptar recursos de shared a tipos de api-mobile usando los tipos retenidos
	// 7. Crear loggerAdapter primero para poder usarlo en lifecycle
	if wrapper.logrusLogger == nil {
		return nil, nil, fmt.Errorf("logger not initialized")
	}
	loggerAdapter := adapter.NewLoggerAdapter(wrapper.logrusLogger)

	// 8. Crear lifecycle manager con logger configurado
	lifecycleManagerWithLogger := lifecycle.NewManager(loggerAdapter)

	// 9. Adaptar recursos de shared a tipos de api-mobile usando los tipos retenidos
	resources, err := adaptSharedResources(wrapper, loggerAdapter, cfg, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to adapt shared resources: %w", err)
	}

	return resources, lifecycleManagerWithLogger, nil
}

// adaptSharedResources adapta los tipos concretos retenidos a los tipos esperados por api-mobile
func adaptSharedResources(
	wrapper *customFactoriesWrapper,
	loggerAdapter sharedLogger.Logger,
	cfg *config.Config,
	opts *BootstrapOptions,
) (*Resources, error) {

	// 2. PostgreSQL: ya tenemos *sql.DB retenido
	if wrapper.sqlDB == nil {
		return nil, fmt.Errorf("PostgreSQL connection not initialized")
	}

	// 3. MongoDB: obtener *mongo.Database del cliente retenido (opcional)
	var mongoDatabase *mongo.Database
	if wrapper.mongoClient != nil {
		mongoDatabase = wrapper.mongoClient.Database(cfg.Database.MongoDB.Database)
	}

	// 4. RabbitMQ: crear adapter con el channel retenido y envolver con circuit breaker
	// Si está deshabilitado, usar noop publisher
	var rabbitMQPublisher rabbitmq.Publisher
	if opts.IsResourceDisabled("rabbitmq") {
		loggerAdapter.Info("RabbitMQ está deshabilitado, usando noop publisher")
		// rabbitMQPublisher permanece nil, el container usará noop
	} else if wrapper.rabbitChannel != nil {
		basePublisher := adapter.NewMessagePublisherAdapter(
			wrapper.rabbitChannel,
			cfg.Messaging.RabbitMQ.Exchanges.Materials,
			loggerAdapter,
		)

		// Envolver con Circuit Breaker si está habilitado
		cbConfig := cfg.Messaging.RabbitMQ.CircuitBreaker
		if cbConfig.Enabled {
			// Usar configuración por defecto y sobrescribir con valores del config si están definidos
			resilientConfig := rabbitmq.DefaultResilientPublisherConfig()

			// Sobrescribir solo valores configurados (diferentes de zero value)
			if cbConfig.MaxRequests > 0 {
				resilientConfig.MaxRequests = cbConfig.MaxRequests
			}
			if cbConfig.Interval > 0 {
				resilientConfig.Interval = cbConfig.Interval
			}
			if cbConfig.Timeout > 0 {
				resilientConfig.Timeout = cbConfig.Timeout
			}
			if cbConfig.FailureThreshold > 0 {
				resilientConfig.FailureThreshold = cbConfig.FailureThreshold
			}

			// Configurar callback de cambio de estado
			resilientConfig.OnStateChange = func(name string, from, to gobreaker.State) {
				loggerAdapter.Warn("circuit breaker state changed",
					"name", name,
					"from", from.String(),
					"to", to.String(),
				)
			}

			rabbitMQPublisher = rabbitmq.NewResilientPublisher(basePublisher, resilientConfig)
			loggerAdapter.Info("RabbitMQ publisher wrapped with circuit breaker",
				"max_requests", resilientConfig.MaxRequests,
				"timeout", resilientConfig.Timeout.String(),
				"failure_threshold", resilientConfig.FailureThreshold,
			)
		} else {
			rabbitMQPublisher = basePublisher
		}
	}

	// 5. S3: crear adapter con el cliente retenido
	// Si está deshabilitado, usar noop storage
	var s3Storage S3Storage
	if opts.IsResourceDisabled("s3") {
		loggerAdapter.Info("S3 está deshabilitado, usando noop storage")
		// s3Storage permanece nil, el container usará noop
	} else if wrapper.s3Client != nil {
		s3Storage = adapter.NewStorageClientAdapter(
			wrapper.s3Client,
			cfg.Storage.S3.BucketName,
			loggerAdapter,
		)
	}

	// 6. Construir Resources de api-mobile
	resources := &Resources{
		Logger:            loggerAdapter,
		PostgreSQL:        wrapper.sqlDB,
		MongoDB:           mongoDatabase,
		RabbitMQPublisher: rabbitMQPublisher,
		S3Client:          s3Storage,
		JWTSecret:         cfg.Auth.JWT.Secret,
		AuthConfig:        cfg.Auth,
		Config:            cfg,
	}

	return resources, nil
}
