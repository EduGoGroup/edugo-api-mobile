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
	Password       string `mapstructure:"password"`        // Desde ENV
	MaxConnections int    `mapstructure:"max_connections"`
	SSLMode        string `mapstructure:"ssl_mode"`
}

// MongoDBConfig configuración de MongoDB
type MongoDBConfig struct {
	URI      string        `mapstructure:"uri"`      // Desde ENV
	Database string        `mapstructure:"database"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

// MessagingConfig configuración de RabbitMQ
type MessagingConfig struct {
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
}

// RabbitMQConfig configuración de RabbitMQ
type RabbitMQConfig struct {
	URL           string        `mapstructure:"url"` // Desde ENV
	Queues        QueuesConfig  `mapstructure:"queues"`
	Exchanges     ExchangeConfig `mapstructure:"exchanges"`
	PrefetchCount int           `mapstructure:"prefetch_count"`
}

// QueuesConfig nombres de colas
type QueuesConfig struct {
	MaterialUploaded   string `mapstructure:"material_uploaded"`
	AssessmentAttempt  string `mapstructure:"assessment_attempt"`
}

// ExchangeConfig nombres de exchanges
type ExchangeConfig struct {
	Materials string `mapstructure:"materials"`
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

// Validate valida que la configuración tenga los campos obligatorios
func (c *Config) Validate() error {
	if c.Database.Postgres.Password == "" {
		return fmt.Errorf("POSTGRES_PASSWORD is required")
	}
	if c.Database.MongoDB.URI == "" {
		return fmt.Errorf("MONGODB_URI is required")
	}
	if c.Messaging.RabbitMQ.URL == "" {
		return fmt.Errorf("RABBITMQ_URL is required")
	}
	return nil
}
