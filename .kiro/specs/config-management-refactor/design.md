# Design Document

## Overview

Este diseño refactoriza el sistema de configuración de EduGo API Mobile para simplificar el manejo de variables de entorno y configuración, eliminando duplicación y lógica innecesaria. El nuevo sistema se basa en tres principios:

1. **Separación de responsabilidades**: YAML para configuración pública, ENV para secretos
2. **Convención sobre configuración**: Usar las capacidades nativas de Viper sin código manual
3. **Developer Experience**: Una sola forma de configurar que funciona en todos los contextos

## Architecture

### Estructura de Archivos de Configuración

```
config/
├── config.yaml              # Base: valores comunes a todos los ambientes
├── config-local.yaml        # Local: desarrollo en máquina del dev
├── config-dev.yaml          # Dev: servidor de desarrollo
├── config-qa.yaml           # QA: ambiente de pruebas
└── config-prod.yaml         # Prod: producción

.env.example                 # Template con todas las variables requeridas
.env                         # Local: valores reales (gitignored)

internal/config/
├── config.go                # Structs de configuración
├── loader.go                # Loader simplificado con Viper
├── validator.go             # Validación de configuración
└── loader_test.go           # Tests del loader

tools/
└── configctl/
    ├── main.go              # CLI para gestionar configuración
    ├── add.go               # Comando para agregar variables
    ├── validate.go          # Comando para validar configuración
    └── generate.go          # Comando para generar documentación
```

### Principio de Separación: Público vs Secreto

**Configuración Pública (YAML)**:
- Puertos, hosts, timeouts
- Nombres de colas, exchanges
- Nombres de bases de datos
- Configuración de logging
- Límites y thresholds

**Secretos (ENV)**:
- Passwords de bases de datos
- API keys y tokens
- Connection strings con credenciales
- Claves de cifrado

### Convención de Nombres

**Variables de Entorno**:
- Formato: `SECTION_SUBSECTION_KEY`
- Ejemplo: `DATABASE_POSTGRES_PASSWORD`
- Viper automáticamente mapea `database.postgres.password` → `DATABASE_POSTGRES_PASSWORD`

**Campos en YAML**:
- Formato: snake_case
- Jerarquía con indentación
- Comentarios para indicar secretos

## Components and Interfaces

### 1. Config Loader (Simplificado)

```go
// loader.go
package config

import (

"fmt"
"os"
"strings"

"github.com/spf13/viper"
)

func Load() (*Config, error) {
    v := viper.New()

    // 1. Configurar Viper
    v.SetConfigType("yaml")
    v.AddConfigPath("./config")
    v.AddConfigPath("../config")

    // 2. Cargar archivo base
    v.SetConfigName("config")
    if err := v.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, fmt.Errorf("error reading base config: %w", err)
        }
    }

    // 3. Merge archivo específico del ambiente
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "local"
    }

    v.SetConfigName(fmt.Sprintf("config-%s", env))
    if err := v.MergeInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, fmt.Errorf("error merging %s config: %w", env, err)
        }
    }


    // 4. ENV vars (precedencia automática)
    v.AutomaticEnv()
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // 5. Unmarshal
    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("error unmarshaling config: %w", err)
    }

    // 6. Validar
    if err := Validate(&cfg); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &cfg, nil
}
```

**Cambios clave**:
- ❌ Eliminado: `BindEnv()` manual para cada variable
- ❌ Eliminado: `Set()` manual para forzar precedencia
- ✅ Agregado: `AutomaticEnv()` con `SetEnvKeyReplacer()`
- ✅ Simplificado: Viper maneja la precedencia automáticamente

### 2. Config Structs (Actualizado)

```go
// config.go
package config

import "time"

type Config struct {
    Server    ServerConfig    `mapstructure:"server"`
    Database  DatabaseConfig  `mapstructure:"database"`
    Messaging MessagingConfig `mapstructure:"messaging"`
    Storage   StorageConfig   `mapstructure:"storage"`
    Logging   LoggingConfig   `mapstructure:"logging"`
}

type ServerConfig struct {
    Port         int           `mapstructure:"port"`
    Host         string        `mapstructure:"host"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
    Postgres PostgresConfig `mapstructure:"postgres"`
    MongoDB  MongoDBConfig  `mapstructure:"mongodb"`
}

type PostgresConfig struct {
    Host           string `mapstructure:"host"`
    Port           int    `mapstructure:"port"`
    Database       string `mapstructure:"database"`
    User           string `mapstructure:"user"`
    Password       string `mapstructure:"password"` // ENV: DATABASE_POSTGRES_PASSWORD
    MaxConnections int    `mapstructure:"max_connections"`
    SSLMode        string `mapstructure:"ssl_mode"`
}

type MongoDBConfig struct {
    URI      string        `mapstructure:"uri"` // ENV: DATABASE_MONGODB_URI
    Database string        `mapstructure:"database"`
    Timeout  time.Duration `mapstructure:"timeout"`
}

type MessagingConfig struct {
    RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
}

type RabbitMQConfig struct {
    URL           string         `mapstructure:"url"` // ENV: MESSAGING_RABBITMQ_URL
    Queues        QueuesConfig   `mapstructure:"queues"`
    Exchanges     ExchangeConfig `mapstructure:"exchanges"`
    PrefetchCount int            `mapstructure:"prefetch_count"`
}

type QueuesConfig struct {
    MaterialUploaded  string `mapstructure:"material_uploaded"`
    AssessmentAttempt string `mapstructure:"assessment_attempt"`
}

type ExchangeConfig struct {
    Materials string `mapstructure:"materials"`
}

type StorageConfig struct {
    S3 S3Config `mapstructure:"s3"`
}

type S3Config struct {
    Region          string `mapstructure:"region"`
    BucketName      string `mapstructure:"bucket_name"`
    AccessKeyID     string `mapstructure:"access_key_id"`     // ENV: STORAGE_S3_ACCESS_KEY_ID
    SecretAccessKey string `mapstructure:"secret_access_key"` // ENV: STORAGE_S3_SECRET_ACCESS_KEY
    Endpoint        string `mapstructure:"endpoint"`
}

type LoggingConfig struct {
    Level  string `mapstructure:"level"`
    Format string `mapstructure:"format"`
}
```

### 3. Validator (Nuevo)

```go
// validator.go
package config

import (
    "fmt"
    "strings"
)

func Validate(cfg *Config) error {
    var errors []string

    // Validar secretos requeridos
    if cfg.Database.Postgres.Password == "" {
        errors = append(errors, "DATABASE_POSTGRES_PASSWORD is required")
    }
    if cfg.Database.MongoDB.URI == "" {
        errors = append(errors, "DATABASE_MONGODB_URI is required")
    }
    if cfg.Messaging.RabbitMQ.URL == "" {
        errors = append(errors, "MESSAGING_RABBITMQ_URL is required")
    }
    if cfg.Storage.S3.AccessKeyID == "" {
        errors = append(errors, "STORAGE_S3_ACCESS_KEY_ID is required")
    }
    if cfg.Storage.S3.SecretAccessKey == "" {
        errors = append(errors, "STORAGE_S3_SECRET_ACCESS_KEY is required")
    }

    // Validar valores públicos
    if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
        errors = append(errors, "server.port must be between 1 and 65535")
    }
    if cfg.Database.Postgres.MaxConnections <= 0 {
        errors = append(errors, "database.postgres.max_connections must be positive")
    }

    if len(errors) > 0 {
        return fmt.Errorf("configuration validation failed:\n  - %s",
            strings.Join(errors, "\n  - "))
    }

    return nil
}
```

### 4. ConfigCTL CLI Tool

```go
// tools/configctl/main.go
package main

import (

"fmt"
"os"

"github.com/spf13/cobra"
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "configctl",
        Short: "Configuration management tool for EduGo API",
    }

    rootCmd.AddCommand(addCmd())
    rootCmd.AddCommand(validateCmd())
    rootCmd.AddCommand(generateDocsCmd())

    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

```go
// tools/configctl/add.go
package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func addCmd() *cobra.Command {
	var (
		varType     string
		isSecret    bool
		defaultVal  string
		description string
	)

	cmd := &cobra.Command{
		Use:   "add [hierarchy.path] [name]",
		Short: "Add a new configuration variable",
		Long: `Add a new configuration variable to the system.

Examples:
  # Add a public config variable
  configctl add database.postgres.pool_size --type int --default 10 --desc "Connection pool size"

  # Add a secret variable
  configctl add auth.jwt.secret --type string --secret --desc "JWT signing secret"
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			return addVariable(path, varType, isSecret, defaultVal, description)
		},
	}

	cmd.Flags().StringVar(&varType, "type", "string", "Variable type (string, int, bool, duration)")
	cmd.Flags().BoolVar(&isSecret, "secret", false, "Mark as secret (ENV only)")
	cmd.Flags().StringVar(&defaultVal, "default", "", "Default value")
	cmd.Flags().StringVar(&description, "desc", "", "Description")
	cmd.MarkFlagRequired("desc")

	return cmd
}

func addVariable(path, varType string, isSecret bool, defaultVal, description string) error {
	// 1. Validar path
	if err := validatePath(path); err != nil {
		return err
	}

	// 2. Actualizar config.go (agregar campo al struct)
	if err := updateConfigStruct(path, varType, description); err != nil {
		return fmt.Errorf("failed to update config.go: %w", err)
	}

	// 3. Si es secreto, actualizar .env.example
	if isSecret {
		if err := updateEnvExample(path, description); err != nil {
			return fmt.Errorf("failed to update .env.example: %w", err)
		}
		fmt.Printf("✓ Added secret variable to .env.example\n")
		fmt.Printf("  ENV var: %s\n", pathToEnvVar(path))
	} else {
		// 4. Si es público, actualizar YAMLs
		if err := updateYAMLFiles(path, defaultVal); err != nil {
			return fmt.Errorf("failed to update YAML files: %w", err)
		}
		fmt.Printf("✓ Added public variable to YAML files\n")
	}

	// 5. Actualizar validator.go si es requerido
	if isSecret {
		if err := updateValidator(path); err != nil {
			return fmt.Errorf("failed to update validator: %w", err)
		}
	}

	fmt.Println("\n✓ Configuration variable added successfully!")
	fmt.Println("\nNext steps:")
	if isSecret {
		fmt.Printf("  1. Add %s to your .env file\n", pathToEnvVar(path))
		fmt.Println("  2. Update deployment secrets in your cloud provider")
	} else {
		fmt.Println("  1. Review the default values in config-*.yaml files")
		fmt.Println("  2. Adjust per-environment values as needed")
	}

	return nil
}
```

## Data Models

### Archivo YAML (Público)

```yaml
# config/config.yaml
server:
  port: 8080
  host: "0.0.0.0"
  read_timeout: 30s
  write_timeout: 30s

database:
  postgres:
    host: "localhost"
    port: 5432
    database: "edugo"
    user: "edugo"
    # password: set via DATABASE_POSTGRES_PASSWORD env var
    max_connections: 25
    ssl_mode: "disable"

  mongodb:
    # uri: set via DATABASE_MONGODB_URI env var
    database: "edugo"
    timeout: 10s

messaging:
  rabbitmq:
    # url: set via MESSAGING_RABBITMQ_URL env var
    queues:
      material_uploaded: "edugo.material.uploaded"
      assessment_attempt: "edugo.assessment.attempt"
    exchanges:
      materials: "edugo.materials"
    prefetch_count: 10

storage:
  s3:
    region: "us-east-1"
    bucket_name: "edugo-materials"
    # access_key_id: set via STORAGE_S3_ACCESS_KEY_ID env var
    # secret_access_key: set via STORAGE_S3_SECRET_ACCESS_KEY env var
    endpoint: "" # Optional, for Localstack

logging:
  level: "info"
  format: "json"
```

### Archivo ENV (Secretos)

```bash
# .env.example

# ========================================
# DATABASE SECRETS
# ========================================

# PostgreSQL password
# Used by: database.postgres.password
DATABASE_POSTGRES_PASSWORD=change-me

# MongoDB connection URI (includes credentials)
# Used by: database.mongodb.uri
# Format: mongodb://user:password@host:port/database?authSource=admin
DATABASE_MONGODB_URI=mongodb://edugo:change-me@localhost:27017/edugo?authSource=admin

# ========================================
# MESSAGING SECRETS
# ========================================

# RabbitMQ connection URL (includes credentials)
# Used by: messaging.rabbitmq.url
# Format: amqp://user:password@host:port/
MESSAGING_RABBITMQ_URL=amqp://edugo:change-me@localhost:5672/

# ========================================
# STORAGE SECRETS
# ========================================

# AWS S3 Access Key ID
# Used by: storage.s3.access_key_id
STORAGE_S3_ACCESS_KEY_ID=your-access-key-id

# AWS S3 Secret Access Key
# Used by: storage.s3.secret_access_key
STORAGE_S3_SECRET_ACCESS_KEY=your-secret-access-key

# ========================================
# APPLICATION
# ========================================

# Environment name (local, dev, qa, prod)
# Determines which config-{env}.yaml file to load
APP_ENV=local
```

## Error Handling

### Errores de Configuración

1. **Archivo no encontrado**: Continuar con defaults (permitir solo ENV vars en cloud)
2. **Variable requerida faltante**: Fallar rápido con mensaje claro
3. **Valor inválido**: Fallar con mensaje indicando formato esperado
4. **Conflicto de precedencia**: No puede ocurrir (Viper lo maneja automáticamente)

### Mensajes de Error Claros

```
Configuration validation failed:
  - DATABASE_POSTGRES_PASSWORD is required
  - STORAGE_S3_ACCESS_KEY_ID is required
  - server.port must be between 1 and 65535

Please check your .env file or environment variables.
For local development, copy .env.example to .env and fill in the values.
```

## Testing Strategy

### Unit Tests

1. **Loader Tests**:
   - Test precedencia: ENV > YAML específico > YAML base > defaults
   - Test carga con archivo faltante
   - Test carga solo con ENV vars (cloud mode)

2. **Validator Tests**:
   - Test validación de campos requeridos
   - Test validación de rangos y formatos
   - Test mensajes de error claros

3. **ConfigCTL Tests**:
   - Test agregar variable pública
   - Test agregar variable secreta
   - Test validación de paths
   - Test dry-run mode

### Integration Tests

1. Test carga de configuración en diferentes ambientes
2. Test que la aplicación inicia correctamente con configuración válida
3. Test que la aplicación falla rápido con configuración inválida

### Test Fixtures

```
test/fixtures/config/
├── valid/
│   ├── config.yaml
│   ├── config-test.yaml
│   └── .env
├── invalid/
│   ├── missing-secrets/
│   └── invalid-values/
└── cloud-mode/
    └── .env (solo ENV vars, sin YAML)
```


## Configuración por Entorno de Desarrollo

### 1. IDE (IntelliJ IDEA / GoLand)

**Run Configuration**:
```
Environment Variables: (load from .env file)
Working Directory: $PROJECT_DIR$
```

El IDE puede cargar automáticamente el archivo `.env` usando plugins como EnvFile.

### 2. Editor de Texto (Zed / VSCode)

**`.zed/tasks.json`**:
```json
{
  "tasks": [
    {
      "label": "Run API",
      "command": "go run cmd/main.go",
      "env_file": ".env"
    }
  ]
}
```

**`.vscode/launch.json`**:
```json
{
  "configurations": [
    {
      "name": "Launch API",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/main.go",
      "envFile": "${workspaceFolder}/.env"
    }
  ]
}
```

### 3. Make

**Makefile** (actualizado):
```makefile
# Cargar .env si existe
ifneq (,$(wildcard .env))
    include .env
    export
endif

run: ## Ejecutar en modo desarrollo
	@echo "🚀 Ejecutando $(APP_NAME) (ambiente: $(APP_ENV))..."
	@go run $(MAIN_PATH)
```

### 4. Docker Compose

**docker-compose.yml** (simplificado):
```yaml
services:
  api-mobile:
    build: .
    env_file:
      - .env
    ports:
      - "9090:8080"
    depends_on:
      - postgres
      - mongodb
      - rabbitmq
```

**Ventaja**: Un solo archivo `.env` para todo.

## Migración desde Sistema Actual

### Paso 1: Limpiar YAMLs

Remover todos los secretos de los archivos YAML:
- ❌ `config-local.yaml`: Remover `password`, `uri`, `url`
- ✅ Mantener solo configuración pública

### Paso 2: Actualizar .env.example

Agregar todas las variables de secretos con documentación clara.

### Paso 3: Simplificar loader.go

Remover:
- `BindEnv()` manual
- `Set()` manual
- Lógica de precedencia manual

Agregar:
- `AutomaticEnv()`
- `SetEnvKeyReplacer()`

### Paso 4: Crear validator.go

Mover toda la lógica de validación desde `config.go` a `validator.go`.

### Paso 5: Actualizar Makefile

Agregar carga automática de `.env`.

### Paso 6: Actualizar docker-compose.yml

Usar `env_file: .env` en lugar de variables individuales.

### Paso 7: Crear ConfigCTL

Implementar la herramienta CLI para gestión de configuración.

### Paso 8: Documentación

Crear `CONFIG.md` con toda la documentación de variables.

## Compatibilidad con Cloud Secrets

### AWS Secrets Manager

```go
// Ejemplo de integración (opcional, para futuro)
func LoadFromAWS(secretName string) (*Config, error) {
    // 1. Obtener secreto de AWS
    secret, err := getSecretFromAWS(secretName)
    if err != nil {
        return nil, err
    }

    // 2. Setear como ENV vars
    for key, value := range secret {
        os.Setenv(key, value)
    }

    // 3. Cargar configuración normalmente
    return Load()
}
```

### Kubernetes Secrets

```yaml
# deployment.yaml
env:
  - name: DATABASE_POSTGRES_PASSWORD
    valueFrom:
      secretKeyRef:
        name: edugo-secrets
        key: postgres-password
  - name: DATABASE_MONGODB_URI
    valueFrom:
      secretKeyRef:
        name: edugo-secrets
        key: mongodb-uri
```

El sistema de configuración no necesita cambios, solo recibe las ENV vars.

## Ventajas del Nuevo Diseño

1. **Simplicidad**: Menos código, menos bugs
2. **Claridad**: Obvio qué es público y qué es secreto
3. **Consistencia**: Una sola forma de configurar en todos los contextos
4. **Mantenibilidad**: Fácil agregar nuevas variables con ConfigCTL
5. **Seguridad**: Secretos nunca en archivos versionados
6. **Flexibilidad**: Compatible con desarrollo local y despliegue en cloud
7. **Developer Experience**: Setup rápido con `.env.example`

## Desventajas y Trade-offs

1. **Requiere migración**: Hay que actualizar archivos existentes
2. **Convención estricta**: Nombres de ENV vars deben seguir el patrón
3. **Dependencia de Viper**: Pero ya la tenemos, solo la usamos mejor

## Decisiones de Diseño

### ¿Por qué AutomaticEnv() en lugar de BindEnv()?

- **AutomaticEnv()**: Viper automáticamente busca ENV vars para cualquier clave
- **BindEnv()**: Requiere binding manual para cada variable
- **Decisión**: AutomaticEnv() es más simple y escalable

### ¿Por qué separar validator.go?

- **Separación de responsabilidades**: Loader carga, Validator valida
- **Testabilidad**: Más fácil testear validación por separado
- **Claridad**: Validaciones en un solo lugar

### ¿Por qué ConfigCTL en lugar de scripts?

- **Type-safe**: Go en lugar de bash
- **Reutilizable**: Puede usarse en CI/CD
- **Mantenible**: Más fácil de extender

### ¿Por qué mantener YAML para públicos?

- **Legibilidad**: YAML es más legible que ENV vars para configuración compleja
- **Estructura**: YAML soporta jerarquía y tipos nativamente
- **Defaults**: Fácil definir valores por defecto por ambiente
