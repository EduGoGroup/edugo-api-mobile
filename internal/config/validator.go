package config

import (
	"fmt"
	"strings"

	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// Validate valida que la configuración tenga los campos obligatorios y valores válidos
func Validate(cfg *Config) error {
	var validationErrors []string

	// Validar secretos requeridos
	if cfg.Database.Postgres.Password == "" {
		validationErrors = append(validationErrors, "DATABASE_POSTGRES_PASSWORD is required")
	}
	if cfg.Database.MongoDB.URI == "" {
		validationErrors = append(validationErrors, "DATABASE_MONGODB_URI is required")
	}
	if cfg.Messaging.RabbitMQ.URL == "" {
		validationErrors = append(validationErrors, "MESSAGING_RABBITMQ_URL is required")
	}
	if cfg.Storage.S3.AccessKeyID == "" {
		validationErrors = append(validationErrors, "STORAGE_S3_ACCESS_KEY_ID is required")
	}
	if cfg.Storage.S3.SecretAccessKey == "" {
		validationErrors = append(validationErrors, "STORAGE_S3_SECRET_ACCESS_KEY is required")
	}
	// Secret para firmar JWT (vía auth.jwt.secret -> ENV AUTH_JWT_SECRET)
	if cfg.Auth.JWT.Secret == "" {
		validationErrors = append(validationErrors, "AUTH_JWT_SECRET is required")
	}

	// Validar valores públicos
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		validationErrors = append(validationErrors, "server.port must be between 1 and 65535")
	}
	if cfg.Database.Postgres.MaxConnections <= 0 {
		validationErrors = append(validationErrors, "database.postgres.max_connections must be positive")
	}
	if cfg.Storage.S3.BucketName == "" {
		validationErrors = append(validationErrors, "storage.s3.bucket_name is required")
	}

	// Si hay errores, retornar un error compuesto con mensaje claro
	if len(validationErrors) > 0 {
		errorMsg := fmt.Sprintf("Configuration validation failed:\n  - %s\n\nPlease check your .env file or environment variables.\nFor local development, copy .env.example to .env and fill in the values.",
			strings.Join(validationErrors, "\n  - "))
		return errors.NewValidationError(errorMsg)
	}

	return nil
}
