package config

import (
	"os"
	"testing"
)

func TestLoad_WithEnvVars(t *testing.T) {
	// Setup: Set required environment variables
	os.Setenv("APP_ENV", "local")
	os.Setenv("DATABASE_POSTGRES_PASSWORD", "test-password")
	os.Setenv("DATABASE_MONGODB_URI", "mongodb://test:test@localhost:27017/test")
	os.Setenv("MESSAGING_RABBITMQ_URL", "amqp://test:test@localhost:5672/")
	os.Setenv("STORAGE_S3_ACCESS_KEY_ID", "test-key-id")
	os.Setenv("STORAGE_S3_SECRET_ACCESS_KEY", "test-secret-key")
	os.Setenv("STORAGE_S3_BUCKET_NAME", "test-bucket")
	// JWT secret necesario para validación
	os.Setenv("AUTH_JWT_SECRET", "test-jwt-secret")

	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("DATABASE_POSTGRES_PASSWORD")
		os.Unsetenv("DATABASE_MONGODB_URI")
		os.Unsetenv("MESSAGING_RABBITMQ_URL")
		os.Unsetenv("STORAGE_S3_ACCESS_KEY_ID")
		os.Unsetenv("STORAGE_S3_SECRET_ACCESS_KEY")
		os.Unsetenv("STORAGE_S3_BUCKET_NAME")
		os.Unsetenv("AUTH_JWT_SECRET")
	}()

	// Execute
	cfg, err := Load()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if cfg.Database.Postgres.Password != "test-password" {
		t.Errorf("Expected password 'test-password', got '%s'", cfg.Database.Postgres.Password)
	}

	if cfg.Database.MongoDB.URI != "mongodb://test:test@localhost:27017/test" {
		t.Errorf("Expected MongoDB URI, got '%s'", cfg.Database.MongoDB.URI)
	}

	if cfg.Messaging.RabbitMQ.URL != "amqp://test:test@localhost:5672/" {
		t.Errorf("Expected RabbitMQ URL, got '%s'", cfg.Messaging.RabbitMQ.URL)
	}

	if cfg.Storage.S3.AccessKeyID != "test-key-id" {
		t.Errorf("Expected S3 access key 'test-key-id', got '%s'", cfg.Storage.S3.AccessKeyID)
	}

	if cfg.Storage.S3.SecretAccessKey != "test-secret-key" {
		t.Errorf("Expected S3 secret key 'test-secret-key', got '%s'", cfg.Storage.S3.SecretAccessKey)
	}
}

func TestLoad_EnvVarsOverrideYAML(t *testing.T) {
	// Setup: Set environment variables that should override YAML values
	os.Setenv("APP_ENV", "local")
	os.Setenv("SERVER_PORT", "9999") // Override default 8080
	os.Setenv("DATABASE_POSTGRES_PASSWORD", "env-password")
	os.Setenv("DATABASE_MONGODB_URI", "mongodb://env:env@localhost:27017/env")
	os.Setenv("MESSAGING_RABBITMQ_URL", "amqp://env:env@localhost:5672/")
	os.Setenv("STORAGE_S3_ACCESS_KEY_ID", "env-key")
	os.Setenv("STORAGE_S3_SECRET_ACCESS_KEY", "env-secret")
	os.Setenv("STORAGE_S3_BUCKET_NAME", "env-bucket")
	// JWT secret necesario para validación
	os.Setenv("AUTH_JWT_SECRET", "test-jwt-secret")

	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DATABASE_POSTGRES_PASSWORD")
		os.Unsetenv("DATABASE_MONGODB_URI")
		os.Unsetenv("MESSAGING_RABBITMQ_URL")
		os.Unsetenv("STORAGE_S3_ACCESS_KEY_ID")
		os.Unsetenv("STORAGE_S3_SECRET_ACCESS_KEY")
		os.Unsetenv("STORAGE_S3_BUCKET_NAME")
		os.Unsetenv("AUTH_JWT_SECRET")
	}()

	// Execute
	cfg, err := Load()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// ENV var should override YAML/default
	if cfg.Server.Port != 9999 {
		t.Errorf("Expected port 9999 from ENV, got %d", cfg.Server.Port)
	}

	if cfg.Database.Postgres.Password != "env-password" {
		t.Errorf("Expected password from ENV, got '%s'", cfg.Database.Postgres.Password)
	}
}

func TestLoad_MissingRequiredEnvVars(t *testing.T) {
	// Setup: Clear all environment variables
	os.Unsetenv("DATABASE_POSTGRES_PASSWORD")
	os.Unsetenv("DATABASE_MONGODB_URI")
	os.Unsetenv("MESSAGING_RABBITMQ_URL")
	os.Unsetenv("STORAGE_S3_ACCESS_KEY_ID")
	os.Unsetenv("STORAGE_S3_SECRET_ACCESS_KEY")

	// Execute
	_, err := Load()

	// Assert: Should fail validation
	if err == nil {
		t.Fatal("Expected error for missing required env vars, got nil")
	}

	// Error message should mention missing variables
	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Expected non-empty error message")
	}
}

func TestLoad_Defaults(t *testing.T) {
	// Setup: Set only required env vars
	os.Setenv("DATABASE_POSTGRES_PASSWORD", "test")
	os.Setenv("DATABASE_MONGODB_URI", "mongodb://test:test@localhost:27017/test")
	os.Setenv("MESSAGING_RABBITMQ_URL", "amqp://test:test@localhost:5672/")
	os.Setenv("STORAGE_S3_ACCESS_KEY_ID", "test")
	os.Setenv("STORAGE_S3_SECRET_ACCESS_KEY", "test")
	os.Setenv("STORAGE_S3_BUCKET_NAME", "test-bucket")
	// JWT secret necesario para validación
	os.Setenv("AUTH_JWT_SECRET", "test-jwt-secret")

	defer func() {
		os.Unsetenv("DATABASE_POSTGRES_PASSWORD")
		os.Unsetenv("DATABASE_MONGODB_URI")
		os.Unsetenv("MESSAGING_RABBITMQ_URL")
		os.Unsetenv("STORAGE_S3_ACCESS_KEY_ID")
		os.Unsetenv("STORAGE_S3_SECRET_ACCESS_KEY")
		os.Unsetenv("STORAGE_S3_BUCKET_NAME")
		os.Unsetenv("AUTH_JWT_SECRET")
	}()

	// Execute
	cfg, err := Load()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check defaults are applied
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", cfg.Server.Port)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected default host '0.0.0.0', got '%s'", cfg.Server.Host)
	}

	if cfg.Database.Postgres.MaxConnections != 25 {
		t.Errorf("Expected default max_connections 25, got %d", cfg.Database.Postgres.MaxConnections)
	}

	if cfg.Logging.Level != "info" {
		t.Errorf("Expected default log level 'info', got '%s'", cfg.Logging.Level)
	}
}
