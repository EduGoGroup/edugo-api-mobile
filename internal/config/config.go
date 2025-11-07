package config

import (
	"fmt"
	"time"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Messaging MessagingConfig `mapstructure:"messaging"`
	Storage   StorageConfig   `mapstructure:"storage"`
	Logging   LoggingConfig   `mapstructure:"logging"`
}

// ServerConfig configuración del servidor HTTP
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig configuración de bases de datos
type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
	MongoDB  MongoDBConfig  `mapstructure:"mongodb"`
}

// PostgresConfig configuración de PostgreSQL
type PostgresConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	Database       string `mapstructure:"database"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"` // ENV: DATABASE_POSTGRES_PASSWORD (required)
	MaxConnections int    `mapstructure:"max_connections"`
	SSLMode        string `mapstructure:"ssl_mode"`
}

// MongoDBConfig configuración de MongoDB
type MongoDBConfig struct {
	URI      string        `mapstructure:"uri"` // ENV: DATABASE_MONGODB_URI (required, format: mongodb://user:pass@host:port/db?authSource=admin)
	Database string        `mapstructure:"database"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

// MessagingConfig configuración de RabbitMQ
type MessagingConfig struct {
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
}

// RabbitMQConfig configuración de RabbitMQ
type RabbitMQConfig struct {
	URL           string         `mapstructure:"url"` // ENV: MESSAGING_RABBITMQ_URL (required, format: amqp://user:pass@host:port/)
	Queues        QueuesConfig   `mapstructure:"queues"`
	Exchanges     ExchangeConfig `mapstructure:"exchanges"`
	PrefetchCount int            `mapstructure:"prefetch_count"`
}

// QueuesConfig nombres de colas
type QueuesConfig struct {
	MaterialUploaded  string `mapstructure:"material_uploaded"`
	AssessmentAttempt string `mapstructure:"assessment_attempt"`
}

// ExchangeConfig nombres de exchanges
type ExchangeConfig struct {
	Materials string `mapstructure:"materials"`
}

// StorageConfig configuración de almacenamiento
type StorageConfig struct {
	S3 S3Config `mapstructure:"s3"`
}

// S3Config configuración de AWS S3
type S3Config struct {
	Region          string `mapstructure:"region"`
	BucketName      string `mapstructure:"bucket_name"`
	AccessKeyID     string `mapstructure:"access_key_id"`     // ENV: STORAGE_S3_ACCESS_KEY_ID (required)
	SecretAccessKey string `mapstructure:"secret_access_key"` // ENV: STORAGE_S3_SECRET_ACCESS_KEY (required)
	Endpoint        string `mapstructure:"endpoint"`          // Optional, for Localstack in development
}

// LoggingConfig configuración de logging
type LoggingConfig struct {
	Level  string `mapstructure:"level"`  // debug, info, warn, error
	Format string `mapstructure:"format"` // json, text
}

// GetPostgresConnectionString construye la cadena de conexión PostgreSQL
func (c *PostgresConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}
