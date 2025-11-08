# Configuración por Ambientes - API Mobile

## Archivos

- `config.yaml` - Configuración base (común a todos los ambientes)
- `config-local.yaml` - Desarrollo local
- `config-dev.yaml` - Servidor de desarrollo
- `config-qa.yaml` - QA/Staging
- `config-prod.yaml` - Producción

## Cómo Funciona

La aplicación carga configuración con esta precedencia (mayor a menor):

1. **Variables de ambiente** (ej: `EDUGO_MOBILE_SERVER_PORT=9090`)
2. **Archivo específico** (ej: `config-dev.yaml`)
3. **Archivo base** (`config.yaml`)
4. **Defaults** (valores por defecto en código)

## Uso

### Cambiar Ambiente

```bash
# Local (default)
go run cmd/main.go

# Development
APP_ENV=dev go run cmd/main.go

# QA
APP_ENV=qa go run cmd/main.go

# Production
APP_ENV=prod go run cmd/main.go
```

### Sobrescribir con ENV Variables

```bash
# Cambiar puerto
APP_ENV=dev EDUGO_MOBILE_SERVER_PORT=9090 go run cmd/main.go

# Cambiar log level
EDUGO_MOBILE_LOGGING_LEVEL=debug go run cmd/main.go
```

## Variables de Ambiente

Prefijo: `EDUGO_MOBILE_`

Formato: Reemplazar `.` con `_` y convertir a MAYÚSCULAS

Ejemplos:
- `server.port` → `EDUGO_MOBILE_SERVER_PORT`
- `database.postgres.host` → `EDUGO_MOBILE_DATABASE_POSTGRES_HOST`
- `logging.level` → `EDUGO_MOBILE_LOGGING_LEVEL`
- `bootstrap.optional_resources.rabbitmq` → `EDUGO_MOBILE_BOOTSTRAP_OPTIONAL_RESOURCES_RABBITMQ`
- `bootstrap.optional_resources.s3` → `EDUGO_MOBILE_BOOTSTRAP_OPTIONAL_RESOURCES_S3`

## Recursos Opcionales

La aplicación soporta recursos de infraestructura opcionales que pueden fallar sin detener el arranque.
Por defecto, RabbitMQ y S3 son opcionales para facilitar el desarrollo local.

### Configurar Recursos Opcionales

**Via archivo YAML:**

```yaml
bootstrap:
  optional_resources:
    rabbitmq: true  # true = opcional, false = requerido
    s3: true        # true = opcional, false = requerido
```

**Via variables de ambiente:**

```bash
# Hacer RabbitMQ requerido (falla si no está disponible)
EDUGO_MOBILE_BOOTSTRAP_OPTIONAL_RESOURCES_RABBITMQ=false go run cmd/main.go

# Hacer S3 requerido (falla si no está disponible)
EDUGO_MOBILE_BOOTSTRAP_OPTIONAL_RESOURCES_S3=false go run cmd/main.go
```

### Comportamiento

- **Recurso Opcional (true)**: Si falla al inicializar, la aplicación registra una advertencia y continúa con una implementación noop
- **Recurso Requerido (false)**: Si falla al inicializar, la aplicación detiene el arranque con un error

### Recursos Soportados

- `rabbitmq` - Sistema de mensajería RabbitMQ (opcional por defecto)
- `s3` - Almacenamiento AWS S3 (opcional por defecto)
- `postgresql` - Base de datos PostgreSQL (siempre requerido)
- `mongodb` - Base de datos MongoDB (siempre requerido)

## Secretos

**NUNCA** commitear secretos en archivos YAML.

Usar variables de ambiente:
- `POSTGRES_PASSWORD` - Password de PostgreSQL
- `MONGODB_URI` - URI completa de MongoDB
- `RABBITMQ_URL` - URL completa de RabbitMQ

## Agregar Nueva Configuración

1. Agregar campo en `internal/config/config.go`
2. Agregar valor en `config.yaml` (base)
3. Sobrescribir en archivos específicos si es necesario
4. Usar en código: `cfg.NuevoCampo`
