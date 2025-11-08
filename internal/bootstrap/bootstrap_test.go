package bootstrap

import (
	"context"
	"database/sql"
	"testing"

	"github.com/EduGoGroup/edugo-api-mobile/internal/bootstrap/noop"
	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

// TestNew verifica que el constructor New crea un Bootstrapper correctamente
func TestNew(t *testing.T) {
	cfg := &config.Config{
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Bootstrap: config.BootstrapConfig{
			OptionalResources: config.OptionalResourcesConfig{
				RabbitMQ: true,
				S3:       true,
			},
		},
	}

	b := New(cfg)

	assert.NotNil(t, b)
	assert.Equal(t, cfg, b.config)
	assert.NotNil(t, b.options)
	assert.NotNil(t, b.options.OptionalResources)
}

// TestNewWithOptions verifica que las opciones funcionales se aplican correctamente
func TestNewWithOptions(t *testing.T) {
	cfg := &config.Config{
		Bootstrap: config.BootstrapConfig{
			OptionalResources: config.OptionalResourcesConfig{
				RabbitMQ: true,
				S3:       true,
			},
		},
	}
	mockLogger := logger.NewZapLogger("info", "json")

	b := New(cfg,
		WithLogger(mockLogger),
		WithOptionalResource("rabbitmq"),
		WithOptionalResource("s3"),
	)

	assert.NotNil(t, b)
	assert.Equal(t, mockLogger, b.options.Logger)
	assert.True(t, b.options.OptionalResources["rabbitmq"])
	assert.True(t, b.options.OptionalResources["s3"])
}

// TestInitializeInfrastructureWithMocks verifica la inicialización con todos los recursos mockeados
func TestInitializeInfrastructureWithMocks(t *testing.T) {
	ctx := context.Background()

	// Crear mocks
	mockLogger := logger.NewZapLogger("info", "json")
	mockDB := &sql.DB{}
	mockMongoDB := &mongo.Database{}
	mockPublisher := noop.NewNoopPublisher(mockLogger)
	mockS3 := noop.NewNoopS3Storage(mockLogger)

	cfg := &config.Config{
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Auth: config.AuthConfig{
			JWT: config.JWTConfig{
				Secret: "test-secret",
			},
		},
		Bootstrap: config.BootstrapConfig{
			OptionalResources: config.OptionalResourcesConfig{
				RabbitMQ: true,
				S3:       true,
			},
		},
	}

	// Crear bootstrapper con todos los recursos inyectados
	b := New(cfg,
		WithLogger(mockLogger),
		WithPostgreSQL(mockDB),
		WithMongoDB(mockMongoDB),
		WithRabbitMQ(mockPublisher),
		WithS3Client(mockS3),
	)

	// Inicializar infraestructura
	resources, cleanup, err := b.InitializeInfrastructure(ctx)

	require.NoError(t, err)
	require.NotNil(t, resources)
	require.NotNil(t, cleanup)

	// Verificar que los recursos son los inyectados
	assert.Equal(t, mockLogger, resources.Logger)
	assert.Equal(t, mockDB, resources.PostgreSQL)
	assert.Equal(t, mockMongoDB, resources.MongoDB)
	assert.Equal(t, mockPublisher, resources.RabbitMQPublisher)
	assert.Equal(t, mockS3, resources.S3Client)
	assert.Equal(t, "test-secret", resources.JWTSecret)

	// Ejecutar cleanup (no debería fallar con mocks)
	err = cleanup()
	assert.NoError(t, err)
}

// TestIsOptional verifica la lógica de recursos opcionales
func TestIsOptional(t *testing.T) {
	cfg := &config.Config{
		Bootstrap: config.BootstrapConfig{
			OptionalResources: config.OptionalResourcesConfig{
				RabbitMQ: true,
				S3:       true,
			},
		},
	}

	tests := []struct {
		name         string
		resourceName string
		options      []BootstrapOption
		expected     bool
	}{
		{
			name:         "rabbitmq is optional by default",
			resourceName: "rabbitmq",
			options:      nil,
			expected:     true,
		},
		{
			name:         "s3 is optional by default",
			resourceName: "s3",
			options:      nil,
			expected:     true,
		},
		{
			name:         "postgresql is required by default",
			resourceName: "postgresql",
			options:      nil,
			expected:     false,
		},
		{
			name:         "mongodb is required by default",
			resourceName: "mongodb",
			options:      nil,
			expected:     false,
		},
		{
			name:         "custom optional resource",
			resourceName: "postgresql",
			options:      []BootstrapOption{WithOptionalResource("postgresql")},
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(cfg, tt.options...)
			result := b.isOptional(tt.resourceName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestInitializeLoggerWithInjected verifica que se usa el logger inyectado
func TestInitializeLoggerWithInjected(t *testing.T) {
	mockLogger := logger.NewZapLogger("debug", "text")

	cfg := &config.Config{
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Bootstrap: config.BootstrapConfig{
			OptionalResources: config.OptionalResourcesConfig{
				RabbitMQ: true,
				S3:       true,
			},
		},
	}

	b := New(cfg, WithLogger(mockLogger))
	err := b.initializeLogger()

	require.NoError(t, err)
	assert.Equal(t, mockLogger, b.logger)
}

// TestInitializeLoggerWithFactory verifica que se crea un logger nuevo cuando no hay inyectado
func TestInitializeLoggerWithFactory(t *testing.T) {
	cfg := &config.Config{
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Bootstrap: config.BootstrapConfig{
			OptionalResources: config.OptionalResourcesConfig{
				RabbitMQ: true,
				S3:       true,
			},
		},
	}

	b := New(cfg)
	err := b.initializeLogger()

	require.NoError(t, err)
	assert.NotNil(t, b.logger)
}
