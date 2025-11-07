package config

import (
	"strings"
	"testing"
)

func TestValidate_AllFieldsValid(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Database: DatabaseConfig{
			Postgres: PostgresConfig{
				Password:       "valid-password",
				MaxConnections: 25,
			},
			MongoDB: MongoDBConfig{
				URI: "mongodb://user:pass@localhost:27017/db",
			},
		},
		Messaging: MessagingConfig{
			RabbitMQ: RabbitMQConfig{
				URL: "amqp://user:pass@localhost:5672/",
			},
		},
		Storage: StorageConfig{
			S3: S3Config{
				AccessKeyID:     "test-key-id",
				SecretAccessKey: "test-secret-key",
				BucketName:      "test-bucket",
			},
		},
	}

	err := Validate(cfg)
	if err != nil {
		t.Errorf("Expected no error for valid config, got: %v", err)
	}
}

func TestValidate_MissingPostgresPassword(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Port: 8080},
		Database: DatabaseConfig{
			Postgres: PostgresConfig{
				Password:       "", // Missing
				MaxConnections: 25,
			},
			MongoDB: MongoDBConfig{
				URI: "mongodb://user:pass@localhost:27017/db",
			},
		},
		Messaging: MessagingConfig{
			RabbitMQ: RabbitMQConfig{
				URL: "amqp://user:pass@localhost:5672/",
			},
		},
		Storage: StorageConfig{
			S3: S3Config{
				AccessKeyID:     "test-key",
				SecretAccessKey: "test-secret",
				BucketName:      "test-bucket",
			},
		},
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("Expected error for missing postgres password, got nil")
	}

	if !strings.Contains(err.Error(), "DATABASE_POSTGRES_PASSWORD") {
		t.Errorf("Error should mention DATABASE_POSTGRES_PASSWORD, got: %v", err)
	}
}

func TestValidate_MissingMongoDBURI(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Port: 8080},
		Database: DatabaseConfig{
			Postgres: PostgresConfig{
				Password:       "test-pass",
				MaxConnections: 25,
			},
			MongoDB: MongoDBConfig{
				URI: "", // Missing
			},
		},
		Messaging: MessagingConfig{
			RabbitMQ: RabbitMQConfig{
				URL: "amqp://user:pass@localhost:5672/",
			},
		},
		Storage: StorageConfig{
			S3: S3Config{
				AccessKeyID:     "test-key",
				SecretAccessKey: "test-secret",
				BucketName:      "test-bucket",
			},
		},
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("Expected error for missing MongoDB URI, got nil")
	}

	if !strings.Contains(err.Error(), "DATABASE_MONGODB_URI") {
		t.Errorf("Error should mention DATABASE_MONGODB_URI, got: %v", err)
	}
}

func TestValidate_MissingRabbitMQURL(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Port: 8080},
		Database: DatabaseConfig{
			Postgres: PostgresConfig{
				Password:       "test-pass",
				MaxConnections: 25,
			},
			MongoDB: MongoDBConfig{
				URI: "mongodb://user:pass@localhost:27017/db",
			},
		},
		Messaging: MessagingConfig{
			RabbitMQ: RabbitMQConfig{
				URL: "", // Missing
			},
		},
		Storage: StorageConfig{
			S3: S3Config{
				AccessKeyID:     "test-key",
				SecretAccessKey: "test-secret",
				BucketName:      "test-bucket",
			},
		},
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("Expected error for missing RabbitMQ URL, got nil")
	}

	if !strings.Contains(err.Error(), "MESSAGING_RABBITMQ_URL") {
		t.Errorf("Error should mention MESSAGING_RABBITMQ_URL, got: %v", err)
	}
}

func TestValidate_MissingS3Credentials(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Port: 8080},
		Database: DatabaseConfig{
			Postgres: PostgresConfig{
				Password:       "test-pass",
				MaxConnections: 25,
			},
			MongoDB: MongoDBConfig{
				URI: "mongodb://user:pass@localhost:27017/db",
			},
		},
		Messaging: MessagingConfig{
			RabbitMQ: RabbitMQConfig{
				URL: "amqp://user:pass@localhost:5672/",
			},
		},
		Storage: StorageConfig{
			S3: S3Config{
				AccessKeyID:     "", // Missing
				SecretAccessKey: "", // Missing
				BucketName:      "test-bucket",
			},
		},
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("Expected error for missing S3 credentials, got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "STORAGE_S3_ACCESS_KEY_ID") {
		t.Errorf("Error should mention STORAGE_S3_ACCESS_KEY_ID, got: %v", err)
	}
	if !strings.Contains(errMsg, "STORAGE_S3_SECRET_ACCESS_KEY") {
		t.Errorf("Error should mention STORAGE_S3_SECRET_ACCESS_KEY, got: %v", err)
	}
}

func TestValidate_InvalidServerPort(t *testing.T) {
	testCases := []struct {
		name string
		port int
	}{
		{"Zero port", 0},
		{"Negative port", -1},
		{"Port too high", 70000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &Config{
				Server: ServerConfig{Port: tc.port},
				Database: DatabaseConfig{
					Postgres: PostgresConfig{
						Password:       "test-pass",
						MaxConnections: 25,
					},
					MongoDB: MongoDBConfig{
						URI: "mongodb://user:pass@localhost:27017/db",
					},
				},
				Messaging: MessagingConfig{
					RabbitMQ: RabbitMQConfig{
						URL: "amqp://user:pass@localhost:5672/",
					},
				},
				Storage: StorageConfig{
					S3: S3Config{
						AccessKeyID:     "test-key",
						SecretAccessKey: "test-secret",
						BucketName:      "test-bucket",
					},
				},
			}

			err := Validate(cfg)
			if err == nil {
				t.Fatalf("Expected error for invalid port %d, got nil", tc.port)
			}

			if !strings.Contains(err.Error(), "server.port") {
				t.Errorf("Error should mention server.port, got: %v", err)
			}
		})
	}
}

func TestValidate_InvalidMaxConnections(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Port: 8080},
		Database: DatabaseConfig{
			Postgres: PostgresConfig{
				Password:       "test-pass",
				MaxConnections: 0, // Invalid
			},
			MongoDB: MongoDBConfig{
				URI: "mongodb://user:pass@localhost:27017/db",
			},
		},
		Messaging: MessagingConfig{
			RabbitMQ: RabbitMQConfig{
				URL: "amqp://user:pass@localhost:5672/",
			},
		},
		Storage: StorageConfig{
			S3: S3Config{
				AccessKeyID:     "test-key",
				SecretAccessKey: "test-secret",
				BucketName:      "test-bucket",
			},
		},
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("Expected error for invalid max_connections, got nil")
	}

	if !strings.Contains(err.Error(), "max_connections") {
		t.Errorf("Error should mention max_connections, got: %v", err)
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Port: 0}, // Invalid
		Database: DatabaseConfig{
			Postgres: PostgresConfig{
				Password:       "", // Missing
				MaxConnections: -1, // Invalid
			},
			MongoDB: MongoDBConfig{
				URI: "", // Missing
			},
		},
		Messaging: MessagingConfig{
			RabbitMQ: RabbitMQConfig{
				URL: "", // Missing
			},
		},
		Storage: StorageConfig{
			S3: S3Config{
				AccessKeyID:     "", // Missing
				SecretAccessKey: "", // Missing
				BucketName:      "", // Missing
			},
		},
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("Expected error for multiple validation failures, got nil")
	}

	errMsg := err.Error()
	
	// Should contain multiple error messages
	expectedErrors := []string{
		"DATABASE_POSTGRES_PASSWORD",
		"DATABASE_MONGODB_URI",
		"MESSAGING_RABBITMQ_URL",
		"STORAGE_S3_ACCESS_KEY_ID",
		"STORAGE_S3_SECRET_ACCESS_KEY",
		"server.port",
		"max_connections",
		"bucket_name",
	}

	for _, expected := range expectedErrors {
		if !strings.Contains(errMsg, expected) {
			t.Errorf("Error message should contain '%s', got: %v", expected, errMsg)
		}
	}
}

func TestValidate_ErrorMessageFormat(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Port: 8080},
		Database: DatabaseConfig{
			Postgres: PostgresConfig{
				Password:       "", // Missing
				MaxConnections: 25,
			},
			MongoDB: MongoDBConfig{
				URI: "mongodb://user:pass@localhost:27017/db",
			},
		},
		Messaging: MessagingConfig{
			RabbitMQ: RabbitMQConfig{
				URL: "amqp://user:pass@localhost:5672/",
			},
		},
		Storage: StorageConfig{
			S3: S3Config{
				AccessKeyID:     "test-key",
				SecretAccessKey: "test-secret",
				BucketName:      "test-bucket",
			},
		},
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	errMsg := err.Error()
	
	// Should have helpful message format
	if !strings.Contains(errMsg, "Configuration validation failed") {
		t.Error("Error should start with 'Configuration validation failed'")
	}
	
	if !strings.Contains(errMsg, ".env") {
		t.Error("Error should mention .env file")
	}
}
