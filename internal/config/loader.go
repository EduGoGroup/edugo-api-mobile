package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Load carga la configuración usando Viper
// Precedencia: ENV vars > archivo específico > archivo base > defaults
func Load() (*Config, error) {
	v := viper.New()

	// 1. Defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("database.postgres.max_connections", 25)
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")

	// 2. Determinar ambiente
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	// 3. Configuración de archivos
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")  // Por si se ejecuta desde otro directorio

	// Leer archivo base
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading base config: %w", err)
	}

	// Merge archivo específico del ambiente
	v.SetConfigName(fmt.Sprintf("config-%s", env))
	if err := v.MergeInConfig(); err != nil {
		// Ignorar si no existe (es opcional)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error merging %s config: %w", env, err)
		}
	}

	// 4. Environment variables (highest precedence)
	v.AutomaticEnv()
	v.SetEnvPrefix("EDUGO_MOBILE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind variables de ambiente específicas (sin prefijo)
	v.BindEnv("database.postgres.password", "POSTGRES_PASSWORD")
	v.BindEnv("database.mongodb.uri", "MONGODB_URI")
	v.BindEnv("messaging.rabbitmq.url", "RABBITMQ_URL")

	// 5. Unmarshal a struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// 6. Validar configuración
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}
