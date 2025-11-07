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

	// 1. Configurar defaults
	setDefaults(v)

	// 2. Determinar ambiente
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	// 3. Configurar paths y tipo de archivo
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config") // Por si se ejecuta desde otro directorio

	// 4. Leer archivo base (opcional en cloud mode)
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		// En cloud mode, el archivo puede no existir (se usa solo env vars)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading base config: %w", err)
		}
		// Archivo no encontrado es OK, continuamos con defaults + env vars
	}

	// 5. Merge archivo específico del ambiente
	v.SetConfigName(fmt.Sprintf("config-%s", env))
	if err := v.MergeInConfig(); err != nil {
		// Ignorar si no existe (es opcional)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error merging %s config: %w", env, err)
		}
	}

	// 6. Configurar ENV vars (precedencia automática)
	// AutomaticEnv hace que Viper busque automáticamente variables de entorno
	// que coincidan con las claves de configuración
	v.AutomaticEnv()
	// SetEnvKeyReplacer convierte "database.postgres.password" → "DATABASE_POSTGRES_PASSWORD"
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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
