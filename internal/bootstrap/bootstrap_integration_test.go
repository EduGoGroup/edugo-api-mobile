//go:build integration
// +build integration

package bootstrap

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/bootstrap/noop"
	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"go.mongodb.org/mongo-driver/mongo"
)

// TestNormalInitialization tests that all resources initialize successfully
// Requirements: 1.1, 5.1, 5.2
func TestNormalInitialization(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Integration tests disabled. Set RUN_INTEGRATION_TESTS=true to enable")
	}

	ctx := context.Background()

	// Start test containers
	pgContainer, mongoContainer, rabbitContainer := setupTestContainers(t, ctx)
	defer func() {
		pgContainer.Terminate(ctx)
		mongoContainer.Terminate(ctx)
		rabbitContainer.Terminate(ctx)
	}()

	// Get connection details
	pgHost, err := pgContainer.Host(ctx)
	require.NoError(t, err)
	pgPort, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	mongoConnStr, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	rabbitURL, err := rabbitContainer.AmqpURL(ctx)
	require.NoError(t, err)

	// Create config with real connection details
	cfg := createTestConfigWithDetails(pgHost, pgPort.Int(), mongoConnStr, rabbitURL)

	// Create bootstrapper
	b := New(cfg)

	// Initialize infrastructure
	resources, cleanup, err := b.InitializeInfrastructure(ctx)

	// Verify initialization succeeded
	require.NoError(t, err, "initialization should succeed with all resources available")
	require.NotNil(t, resources, "resources should not be nil")
	require.NotNil(t, cleanup, "cleanup function should not be nil")

	// Verify all resources are non-nil
	assert.NotNil(t, resources.Logger, "logger should be initialized")
	assert.NotNil(t, resources.PostgreSQL, "PostgreSQL should be initialized")
	assert.NotNil(t, resources.MongoDB, "MongoDB should be initialized")
	assert.NotNil(t, resources.RabbitMQPublisher, "RabbitMQ publisher should be initialized")
	assert.NotNil(t, resources.S3Client, "S3 client should be initialized (noop)")
	assert.NotEmpty(t, resources.JWTSecret, "JWT secret should be set")

	// Verify PostgreSQL connection works
	err = resources.PostgreSQL.PingContext(ctx)
	assert.NoError(t, err, "PostgreSQL should be pingable")

	// Verify MongoDB connection works
	err = resources.MongoDB.Client().Ping(ctx, nil)
	assert.NoError(t, err, "MongoDB should be pingable")

	// Verify cleanup function works correctly
	err = cleanup()
	assert.NoError(t, err, "cleanup should execute without errors")

	// Verify resources are closed (PostgreSQL should fail to ping after close)
	err = resources.PostgreSQL.PingContext(ctx)
	assert.Error(t, err, "PostgreSQL should be closed after cleanup")
}

// setupTestContainers starts all required test containers
func setupTestContainers(t *testing.T, ctx context.Context) (*postgres.PostgresContainer, *mongodb.MongoDBContainer, *rabbitmq.RabbitMQContainer) {
	t.Helper()

	// Start PostgreSQL
	t.Log("Starting PostgreSQL testcontainer...")
	pgContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("edugo_test"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_pass"),
	)
	require.NoError(t, err, "PostgreSQL container should start")

	// Start MongoDB
	t.Log("Starting MongoDB testcontainer...")
	mongoContainer, err := mongodb.Run(ctx, "mongo:7.0",
		mongodb.WithUsername("test_admin"),
		mongodb.WithPassword("test_pass"),
	)
	if err != nil {
		pgContainer.Terminate(ctx)
		require.NoError(t, err, "MongoDB container should start")
	}

	// Start RabbitMQ
	t.Log("Starting RabbitMQ testcontainer...")
	rabbitContainer, err := rabbitmq.Run(ctx, "rabbitmq:3.12-management-alpine",
		rabbitmq.WithAdminUsername("test_user"),
		rabbitmq.WithAdminPassword("test_pass"),
	)
	if err != nil {
		pgContainer.Terminate(ctx)
		mongoContainer.Terminate(ctx)
		require.NoError(t, err, "RabbitMQ container should start")
	}

	// Wait a bit for containers to be fully ready
	time.Sleep(2 * time.Second)

	return pgContainer, mongoContainer, rabbitContainer
}

// createTestConfigWithDetails creates a test configuration with provided connection details
func createTestConfigWithDetails(pgHost string, pgPort int, mongoConnStr, rabbitURL string) *config.Config {
	return &config.Config{
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Auth: config.AuthConfig{
			JWT: config.JWTConfig{
				Secret: "test-jwt-secret-key",
			},
		},
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{
				Host:     pgHost,
				Port:     pgPort,
				Database: "edugo_test",
				User:     "test_user",
				Password: "test_pass",
				SSLMode:  "disable",
			},
			MongoDB: config.MongoDBConfig{
				URI:      mongoConnStr,
				Database: "edugo_test",
			},
		},
		Messaging: config.MessagingConfig{
			RabbitMQ: config.RabbitMQConfig{
				URL: rabbitURL,
				Exchanges: config.ExchangeConfig{
					Materials: "materials",
				},
			},
		},
		Storage: config.StorageConfig{
			S3: config.S3Config{
				Region:          "us-east-1",
				BucketName:      "test-bucket",
				AccessKeyID:     "test-key",
				SecretAccessKey: "test-secret",
				Endpoint:        "",
			},
		},
		Bootstrap: config.BootstrapConfig{
			OptionalResources: config.OptionalResourcesConfig{
				RabbitMQ: false, // Required for this test
				S3:       true,  // Optional (will use noop)
			},
		},
	}
}

// createTestConfig creates a test configuration with provided connection strings (for simpler tests)
func createTestConfig(pgConnStr, mongoConnStr, rabbitURL string) *config.Config {
	return &config.Config{
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Auth: config.AuthConfig{
			JWT: config.JWTConfig{
				Secret: "test-jwt-secret-key",
			},
		},
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{
				Host:     "localhost",
				Port:     5432,
				Database: "edugo_test",
				User:     "test_user",
				Password: "test_pass",
				SSLMode:  "disable",
			},
			MongoDB: config.MongoDBConfig{
				URI:      mongoConnStr,
				Database: "edugo_test",
			},
		},
		Messaging: config.MessagingConfig{
			RabbitMQ: config.RabbitMQConfig{
				URL: rabbitURL,
				Exchanges: config.ExchangeConfig{
					Materials: "materials",
				},
			},
		},
		Storage: config.StorageConfig{
			S3: config.S3Config{
				Region:          "us-east-1",
				BucketName:      "test-bucket",
				AccessKeyID:     "test-key",
				SecretAccessKey: "test-secret",
				Endpoint:        "",
			},
		},
		Bootstrap: config.BootstrapConfig{
			OptionalResources: config.OptionalResourcesConfig{
				RabbitMQ: false, // Required for this test
				S3:       true,  // Optional (will use noop)
			},
		},
	}
}

// TestOptionalResourceFailure tests that initialization succeeds with noop when optional resource fails
// Requirements: 3.1, 3.2, 3.3
func TestOptionalResourceFailure(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Integration tests disabled. Set RUN_INTEGRATION_TESTS=true to enable")
	}

	ctx := context.Background()

	// Start only PostgreSQL and MongoDB (no RabbitMQ)
	pgContainer, mongoContainer := setupMinimalContainers(t, ctx)
	defer func() {
		pgContainer.Terminate(ctx)
		mongoContainer.Terminate(ctx)
	}()

	// Get connection details
	pgHost, err := pgContainer.Host(ctx)
	require.NoError(t, err)
	pgPort, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	mongoConnStr, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	// Create config with INVALID RabbitMQ URL
	cfg := createTestConfigWithDetails(pgHost, pgPort.Int(), mongoConnStr, "amqp://invalid:5672")

	// Mark RabbitMQ as optional
	cfg.Bootstrap.OptionalResources.RabbitMQ = true

	// Create bootstrapper
	b := New(cfg, WithOptionalResource("rabbitmq"))

	// Initialize infrastructure
	resources, cleanup, err := b.InitializeInfrastructure(ctx)

	// Verify initialization succeeded despite RabbitMQ failure
	require.NoError(t, err, "initialization should succeed with optional resource failure")
	require.NotNil(t, resources, "resources should not be nil")
	require.NotNil(t, cleanup, "cleanup function should not be nil")

	// Verify required resources are initialized
	assert.NotNil(t, resources.Logger, "logger should be initialized")
	assert.NotNil(t, resources.PostgreSQL, "PostgreSQL should be initialized")
	assert.NotNil(t, resources.MongoDB, "MongoDB should be initialized")

	// Verify RabbitMQ is replaced with noop implementation
	assert.NotNil(t, resources.RabbitMQPublisher, "RabbitMQ publisher should not be nil")
	_, isNoop := resources.RabbitMQPublisher.(*noop.NoopPublisher)
	assert.True(t, isNoop, "RabbitMQ publisher should be noop implementation")

	// Verify noop publisher works without errors
	err = resources.RabbitMQPublisher.Publish(ctx, "test-exchange", "test-key", []byte("test"))
	assert.NoError(t, err, "noop publisher should not return errors")

	// Cleanup
	err = cleanup()
	assert.NoError(t, err, "cleanup should execute without errors")
}

// setupMinimalContainers starts only PostgreSQL and MongoDB
func setupMinimalContainers(t *testing.T, ctx context.Context) (*postgres.PostgresContainer, *mongodb.MongoDBContainer) {
	t.Helper()

	// Start PostgreSQL
	t.Log("Starting PostgreSQL testcontainer...")
	pgContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("edugo_test"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_pass"),
	)
	require.NoError(t, err, "PostgreSQL container should start")

	// Start MongoDB
	t.Log("Starting MongoDB testcontainer...")
	mongoContainer, err := mongodb.Run(ctx, "mongo:7.0",
		mongodb.WithUsername("test_admin"),
		mongodb.WithPassword("test_pass"),
	)
	if err != nil {
		pgContainer.Terminate(ctx)
		require.NoError(t, err, "MongoDB container should start")
	}

	// Wait a bit for containers to be fully ready
	time.Sleep(2 * time.Second)

	return pgContainer, mongoContainer
}

// TestRequiredResourceFailure tests that initialization fails when a required resource fails
// Requirements: 3.5, 5.2, 5.3
func TestRequiredResourceFailure(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Integration tests disabled. Set RUN_INTEGRATION_TESTS=true to enable")
	}

	ctx := context.Background()

	// Start only MongoDB (no PostgreSQL)
	mongoContainer := setupMongoOnly(t, ctx)
	defer mongoContainer.Terminate(ctx)

	mongoConnStr, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	// Create config with INVALID PostgreSQL connection
	cfg := createTestConfigWithDetails("invalid-host", 5432, mongoConnStr, "amqp://invalid:5672")

	// PostgreSQL is required by default
	cfg.Bootstrap.OptionalResources.RabbitMQ = true
	cfg.Bootstrap.OptionalResources.S3 = true

	// Create bootstrapper
	b := New(cfg)

	// Initialize infrastructure
	resources, cleanup, err := b.InitializeInfrastructure(ctx)

	// Verify initialization failed with descriptive error
	require.Error(t, err, "initialization should fail when required resource fails")
	assert.Contains(t, err.Error(), "PostgreSQL", "error should mention PostgreSQL")
	assert.Nil(t, resources, "resources should be nil on failure")
	assert.Nil(t, cleanup, "cleanup function should be nil on failure")
}

// TestRequiredResourceFailureWithPartialCleanup tests that partial cleanup is executed
// Requirements: 3.5, 5.2, 5.3
func TestRequiredResourceFailureWithPartialCleanup(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Integration tests disabled. Set RUN_INTEGRATION_TESTS=true to enable")
	}

	ctx := context.Background()

	// Start PostgreSQL but not MongoDB
	pgContainer := setupPostgresOnly(t, ctx)
	defer pgContainer.Terminate(ctx)

	pgHost, err := pgContainer.Host(ctx)
	require.NoError(t, err)
	pgPort, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Create config with valid PostgreSQL but INVALID MongoDB
	cfg := createTestConfigWithDetails(pgHost, pgPort.Int(), "mongodb://invalid:27017/test", "amqp://invalid:5672")

	// MongoDB is required by default
	cfg.Bootstrap.OptionalResources.RabbitMQ = true
	cfg.Bootstrap.OptionalResources.S3 = true

	// Create bootstrapper
	b := New(cfg)

	// Initialize infrastructure
	resources, cleanup, err := b.InitializeInfrastructure(ctx)

	// Verify initialization failed
	require.Error(t, err, "initialization should fail when MongoDB fails")
	assert.Contains(t, err.Error(), "MongoDB", "error should mention MongoDB")
	assert.Nil(t, resources, "resources should be nil on failure")
	assert.Nil(t, cleanup, "cleanup function should be nil on failure")

	// Note: Partial cleanup is executed internally before returning error
	// PostgreSQL connection should have been closed by the bootstrap system
	// We can't verify this directly since resources is nil, but the test
	// verifies that the error is returned correctly and no panic occurs
}

// setupMongoOnly starts only MongoDB container
func setupMongoOnly(t *testing.T, ctx context.Context) *mongodb.MongoDBContainer {
	t.Helper()

	t.Log("Starting MongoDB testcontainer...")
	mongoContainer, err := mongodb.Run(ctx, "mongo:7.0",
		mongodb.WithUsername("test_admin"),
		mongodb.WithPassword("test_pass"),
	)
	require.NoError(t, err, "MongoDB container should start")

	time.Sleep(2 * time.Second)
	return mongoContainer
}

// setupPostgresOnly starts only PostgreSQL container
func setupPostgresOnly(t *testing.T, ctx context.Context) *postgres.PostgresContainer {
	t.Helper()

	t.Log("Starting PostgreSQL testcontainer...")
	pgContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("edugo_test"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_pass"),
	)
	require.NoError(t, err, "PostgreSQL container should start")

	time.Sleep(2 * time.Second)
	return pgContainer
}

// TestMockInjection tests that mocks are used instead of real implementations
// Requirements: 4.1, 4.2, 4.4, 4.5
func TestMockInjection(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Integration tests disabled. Set RUN_INTEGRATION_TESTS=true to enable")
	}

	ctx := context.Background()

	// Create mock resources
	mockLogger := logger.NewZapLogger("info", "json")
	mockDB := &sql.DB{}
	mockMongoDB := &mongo.Database{}
	mockPublisher := noop.NewNoopPublisher(mockLogger)
	mockS3 := noop.NewNoopS3Storage(mockLogger)

	// Create minimal config (won't be used for connections since we inject mocks)
	cfg := &config.Config{
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Auth: config.AuthConfig{
			JWT: config.JWTConfig{
				Secret: "mock-jwt-secret",
			},
		},
		Bootstrap: config.BootstrapConfig{
			OptionalResources: config.OptionalResourcesConfig{
				RabbitMQ: true,
				S3:       true,
			},
		},
	}

	// Create bootstrapper with all mocked resources
	b := New(cfg,
		WithLogger(mockLogger),
		WithPostgreSQL(mockDB),
		WithMongoDB(mockMongoDB),
		WithRabbitMQ(mockPublisher),
		WithS3Client(mockS3),
	)

	// Initialize infrastructure
	resources, cleanup, err := b.InitializeInfrastructure(ctx)

	// Verify initialization succeeded
	require.NoError(t, err, "initialization should succeed with mocked resources")
	require.NotNil(t, resources, "resources should not be nil")
	require.NotNil(t, cleanup, "cleanup function should not be nil")

	// Verify mocks are used instead of real implementations
	assert.Equal(t, mockLogger, resources.Logger, "should use injected logger")
	assert.Equal(t, mockDB, resources.PostgreSQL, "should use injected PostgreSQL")
	assert.Equal(t, mockMongoDB, resources.MongoDB, "should use injected MongoDB")
	assert.Equal(t, mockPublisher, resources.RabbitMQPublisher, "should use injected RabbitMQ publisher")
	assert.Equal(t, mockS3, resources.S3Client, "should use injected S3 client")

	// Verify application works with mocks (no panics or errors)
	assert.NotPanics(t, func() {
		// Simulate using the resources
		_ = resources.Logger
		_ = resources.PostgreSQL
		_ = resources.MongoDB
		_ = resources.RabbitMQPublisher.Publish(ctx, "test", "test", []byte("test"))
	}, "application should work with mocked resources")

	// Cleanup should work with mocks
	err = cleanup()
	assert.NoError(t, err, "cleanup should work with mocked resources")
}

// TestMockInjectionWithCustomFactories tests using custom factories for advanced testing
// Requirements: 4.1, 4.2, 4.4, 4.5
func TestMockInjectionWithCustomFactories(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Integration tests disabled. Set RUN_INTEGRATION_TESTS=true to enable")
	}

	ctx := context.Background()

	// Create mock factory that tracks calls
	mockFactory := &MockDatabaseFactory{
		createPostgreSQLCalled: false,
		createMongoDBCalled:    false,
	}

	// Create config
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
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{
				Host:     "mock-host",
				Port:     5432,
				Database: "mock-db",
			},
			MongoDB: config.MongoDBConfig{
				Database: "mock-db",
			},
		},
		Bootstrap: config.BootstrapConfig{
			OptionalResources: config.OptionalResourcesConfig{
				RabbitMQ: true,
				S3:       true,
			},
		},
	}

	// Create bootstrapper with custom factory
	b := New(cfg, WithDatabaseFactory(mockFactory))

	// Initialize infrastructure
	resources, cleanup, err := b.InitializeInfrastructure(ctx)

	// Verify initialization succeeded
	require.NoError(t, err, "initialization should succeed with custom factory")
	require.NotNil(t, resources, "resources should not be nil")

	// Verify custom factory was used
	assert.True(t, mockFactory.createPostgreSQLCalled, "custom factory should be called for PostgreSQL")
	assert.True(t, mockFactory.createMongoDBCalled, "custom factory should be called for MongoDB")

	// Cleanup
	err = cleanup()
	assert.NoError(t, err)
}

// MockDatabaseFactory is a mock implementation of DatabaseFactory for testing
type MockDatabaseFactory struct {
	createPostgreSQLCalled bool
	createMongoDBCalled    bool
}

func (m *MockDatabaseFactory) CreatePostgreSQL(ctx context.Context, cfg *config.Config, log logger.Logger) (*sql.DB, error) {
	m.createPostgreSQLCalled = true
	// Return a mock DB (note: this won't be functional, but that's ok for this test)
	return &sql.DB{}, nil
}

func (m *MockDatabaseFactory) CreateMongoDB(ctx context.Context, cfg *config.Config, log logger.Logger) (*mongo.Database, error) {
	m.createMongoDBCalled = true
	// Return a mock MongoDB (note: this won't be functional, but that's ok for this test)
	return &mongo.Database{}, nil
}
