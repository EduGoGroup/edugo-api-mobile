package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Load carga la configuración usando Viper
// Precedencia automática: ENV vars > archivo específico > archivo base > defaults
func Load() (*Config, error) {
	v := viper.New()

	// 1. Configurar ENV vars PRIMERO para que tengan máxima precedencia
	// AutomaticEnv + SetEnvKeyReplacer permite que Viper busque ENV vars automáticamente
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	
	// Bind explicit ENV vars para que Unmarshal() las tome
	// Esto es necesario porque Unmarshal() no usa AutomaticEnv()
	bindEnvVars(v)

	// 2. Configurar defaults
	setDefaults(v)

	// 3. Determinar ambiente
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	// 4. Configurar paths y tipo de archivo
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config") // Por si se ejecuta desde otro directorio

	// 5. Leer archivo base (opcional en cloud mode)
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		// En cloud mode, el archivo puede no existir (se usa solo env vars)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading base config: %w", err)
		}
		// Archivo no encontrado es OK, continuamos con defaults + env vars
	}

	// 6. Merge archivo específico del ambiente
	v.SetConfigName(fmt.Sprintf("config-%s", env))
	if err := v.MergeInConfig(); err != nil {
		// Ignorar si no existe (es opcional)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error merging %s config: %w", env, err)
		}
	}

	// 7. Unmarshal a struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// 8. Validar configuración
	if err := Validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// setDefaults configura los valores por defecto
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")

	v.SetDefault("database.postgres.port", 5432)
	v.SetDefault("database.postgres.max_connections", 25)
	v.SetDefault("database.postgres.ssl_mode", "disable")

	v.SetDefault("database.mongodb.timeout", "10s")

	v.SetDefault("messaging.rabbitmq.prefetch_count", 10)

	v.SetDefault("storage.s3.region", "us-east-1")
	v.SetDefault("storage.s3.endpoint", "")

	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
}

// bindEnvVars vincula explícitamente las variables de entorno
// Esto es necesario para que Unmarshal() tome los valores de ENV vars
func bindEnvVars(v *viper.Viper) {
	// Server
	v.BindEnv("server.port")
	v.BindEnv("server.host")
	v.BindEnv("server.read_timeout")
	v.BindEnv("server.write_timeout")

	// Database - Postgres
	v.BindEnv("database.postgres.host")
	v.BindEnv("database.postgres.port")
	v.BindEnv("database.postgres.database")
	v.BindEnv("database.postgres.user")
	v.BindEnv("database.postgres.password")
	v.BindEnv("database.postgres.max_connections")
	v.BindEnv("database.postgres.ssl_mode")

	// Database - MongoDB
	v.BindEnv("database.mongodb.uri")
	v.BindEnv("database.mongodb.database")
	v.BindEnv("database.mongodb.timeout")

	// Messaging - RabbitMQ
	v.BindEnv("messaging.rabbitmq.url")
	v.BindEnv("messaging.rabbitmq.queues.material_uploaded")
	v.BindEnv("messaging.rabbitmq.queues.assessment_attempt")
	v.BindEnv("messaging.rabbitmq.exchanges.materials")
	v.BindEnv("messaging.rabbitmq.prefetch_count")

	// Storage - S3
	v.BindEnv("storage.s3.region")
	v.BindEnv("storage.s3.bucket_name")
	v.BindEnv("storage.s3.access_key_id")
	v.BindEnv("storage.s3.secret_access_key")
	v.BindEnv("storage.s3.endpoint")

	// Logging
	v.BindEnv("logging.level")
	v.BindEnv("logging.format")
}
