# Configuration Guide - EduGo API Mobile

## Overview

This document describes the configuration system for the EduGo API Mobile application. The system uses a combination of YAML files for public configuration and environment variables for secrets.

## Quick Start

### Local Development Setup

1. **Copy the example environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your local values:**
   ```bash
   # Database credentials
   DATABASE_POSTGRES_PASSWORD=your-local-password
   DATABASE_MONGODB_URI=mongodb://user:pass@localhost:27017/edugo?authSource=admin

   # Messaging
   MESSAGING_RABBITMQ_URL=amqp://user:pass@localhost:5672/

   # Storage
   STORAGE_S3_ACCESS_KEY_ID=your-aws-key
   STORAGE_S3_SECRET_ACCESS_KEY=your-aws-secret
   ```

3. **Run the application:**
   ```bash
   make run
   # or
   go run cmd/main.go
   # or
   docker-compose up
   ```

## Configuration Architecture

### Precedence Order

Configuration values are loaded in the following order (highest to lowest precedence):

1. **Environment Variables** (highest)
2. **Environment-specific YAML** (`config-{env}.yaml`)
3. **Base YAML** (`config.yaml`)
4. **Defaults** (in code)

### File Structure

```
config/
├── config.yaml              # Base configuration (all environments)
├── config-local.yaml        # Local development overrides
├── config-dev.yaml          # Development server overrides
├── config-qa.yaml           # QA/Staging overrides
└── config-prod.yaml         # Production overrides

.env.example                 # Template with all required variables
.env                         # Your local secrets (gitignored)
```

## Configuration Variables

### Server Configuration

| Variable | Type | Default | Description | Source |
|----------|------|---------|-------------|--------|
| `server.port` | int | 8080 | HTTP server port | YAML/ENV |
| `server.host` | string | "0.0.0.0" | HTTP server host | YAML/ENV |
| `server.read_timeout` | duration | 30s | HTTP read timeout | YAML/ENV |
| `server.write_timeout` | duration | 30s | HTTP write timeout | YAML/ENV |

**Environment Variable Mapping:**
- `SERVER_PORT` → `server.port`
- `SERVER_HOST` → `server.host`
- `SERVER_READ_TIMEOUT` → `server.read_timeout`
- `SERVER_WRITE_TIMEOUT` → `server.write_timeout`

### Database Configuration

#### PostgreSQL

| Variable | Type | Default | Description | Source |
|----------|------|---------|-------------|--------|
| `database.postgres.host` | string | "localhost" | PostgreSQL host | YAML/ENV |
| `database.postgres.port` | int | 5432 | PostgreSQL port | YAML/ENV |
| `database.postgres.database` | string | "edugo" | Database name | YAML/ENV |
| `database.postgres.user` | string | "edugo" | Database user | YAML/ENV |
| `database.postgres.password` | string | - | Database password | **ENV ONLY** ⚠️ |
| `database.postgres.max_connections` | int | 25 | Max connections | YAML/ENV |
| `database.postgres.ssl_mode` | string | "disable" | SSL mode | YAML/ENV |

**Environment Variable Mapping:**
- `DATABASE_POSTGRES_HOST` → `database.postgres.host`
- `DATABASE_POSTGRES_PORT` → `database.postgres.port`
- `DATABASE_POSTGRES_DATABASE` → `database.postgres.database`
- `DATABASE_POSTGRES_USER` → `database.postgres.user`
- `DATABASE_POSTGRES_PASSWORD` → `database.postgres.password` ⚠️ **Required**
- `DATABASE_POSTGRES_MAX_CONNECTIONS` → `database.postgres.max_connections`
- `DATABASE_POSTGRES_SSL_MODE` → `database.postgres.ssl_mode`

#### MongoDB

| Variable | Type | Default | Description | Source |
|----------|------|---------|-------------|--------|
| `database.mongodb.uri` | string | - | MongoDB connection URI | **ENV ONLY** ⚠️ |
| `database.mongodb.database` | string | "edugo" | Database name | YAML/ENV |
| `database.mongodb.timeout` | duration | 10s | Connection timeout | YAML/ENV |

**Environment Variable Mapping:**
- `DATABASE_MONGODB_URI` → `database.mongodb.uri` ⚠️ **Required**
- `DATABASE_MONGODB_DATABASE` → `database.mongodb.database`
- `DATABASE_MONGODB_TIMEOUT` → `database.mongodb.timeout`

**MongoDB URI Format:**
```
mongodb://user:password@host:port/database?authSource=admin
```

### Messaging Configuration (RabbitMQ)

| Variable | Type | Default | Description | Source |
|----------|------|---------|-------------|--------|
| `messaging.rabbitmq.url` | string | - | RabbitMQ AMQP URL | **ENV ONLY** ⚠️ |
| `messaging.rabbitmq.queues.material_uploaded` | string | "edugo.material.uploaded" | Queue name | YAML/ENV |
| `messaging.rabbitmq.queues.assessment_attempt` | string | "edugo.assessment.attempt" | Queue name | YAML/ENV |
| `messaging.rabbitmq.exchanges.materials` | string | "edugo.materials" | Exchange name | YAML/ENV |
| `messaging.rabbitmq.prefetch_count` | int | 10 | Prefetch count | YAML/ENV |

**Environment Variable Mapping:**
- `MESSAGING_RABBITMQ_URL` → `messaging.rabbitmq.url` ⚠️ **Required**
- `MESSAGING_RABBITMQ_QUEUES_MATERIAL_UPLOADED` → `messaging.rabbitmq.queues.material_uploaded`
- `MESSAGING_RABBITMQ_QUEUES_ASSESSMENT_ATTEMPT` → `messaging.rabbitmq.queues.assessment_attempt`
- `MESSAGING_RABBITMQ_EXCHANGES_MATERIALS` → `messaging.rabbitmq.exchanges.materials`
- `MESSAGING_RABBITMQ_PREFETCH_COUNT` → `messaging.rabbitmq.prefetch_count`

**RabbitMQ URL Format:**
```
amqp://user:password@host:port/
```

### Storage Configuration (AWS S3)

| Variable | Type | Default | Description | Source |
|----------|------|---------|-------------|--------|
| `storage.s3.region` | string | "us-east-1" | AWS region | YAML/ENV |
| `storage.s3.bucket_name` | string | "edugo-materials" | S3 bucket name | YAML/ENV |
| `storage.s3.access_key_id` | string | - | AWS access key | **ENV ONLY** ⚠️ |
| `storage.s3.secret_access_key` | string | - | AWS secret key | **ENV ONLY** ⚠️ |
| `storage.s3.endpoint` | string | "" | Custom endpoint (for Localstack) | YAML/ENV |

**Environment Variable Mapping:**
- `STORAGE_S3_REGION` → `storage.s3.region`
- `STORAGE_S3_BUCKET_NAME` → `storage.s3.bucket_name`
- `STORAGE_S3_ACCESS_KEY_ID` → `storage.s3.access_key_id` ⚠️ **Required**
- `STORAGE_S3_SECRET_ACCESS_KEY` → `storage.s3.secret_access_key` ⚠️ **Required**
- `STORAGE_S3_ENDPOINT` → `storage.s3.endpoint`

### Logging Configuration

| Variable | Type | Default | Description | Source |
|----------|------|---------|-------------|--------|
| `logging.level` | string | "info" | Log level (debug, info, warn, error) | YAML/ENV |
| `logging.format` | string | "json" | Log format (json, text) | YAML/ENV |

**Environment Variable Mapping:**
- `LOGGING_LEVEL` → `logging.level`
- `LOGGING_FORMAT` → `logging.format`

## Environment-Specific Configuration

### Local Development (`APP_ENV=local`)

- Uses `config-local.yaml`
- Typically connects to local services (localhost)
- Debug logging enabled
- Text format logs for readability

### Development Server (`APP_ENV=dev`)

- Uses `config-dev.yaml`
- Connects to development infrastructure
- Debug logging enabled
- JSON format logs

### QA/Staging (`APP_ENV=qa`)

- Uses `config-qa.yaml`
- Connects to QA infrastructure
- Info logging
- JSON format logs

### Production (`APP_ENV=prod`)

- Uses `config-prod.yaml`
- Connects to production infrastructure
- Warn logging (minimal)
- JSON format logs
- SSL required for databases

## Development Environments

### Running with IDE

#### IntelliJ IDEA / GoLand

1. Install the **EnvFile** plugin
2. Edit Run Configuration
3. Add `.env` file in EnvFile tab
4. Run the application

See `.idea/runConfigurations/README.md` for detailed instructions.

#### VSCode

The `.vscode/launch.json` is already configured to load `.env`:

```json
{
  "name": "Launch API",
  "type": "go",
  "request": "launch",
  "program": "${workspaceFolder}/cmd/main.go",
  "envFile": "${workspaceFolder}/.env"
}
```

#### Zed

The `.zed/debug.json` is already configured to load `.env`:

```json
{
  "label": "Go: Debug main (Delve)",
  "envFile": "${workspaceFolder}/.env"
}
```

### Running with Make

The Makefile automatically loads `.env`:

```bash
make run
make build
make test
```

### Running with Docker Compose

Docker Compose automatically loads `.env`:

```bash
docker-compose up
```

## Cloud Deployment

### AWS Secrets Manager

For production deployments, use AWS Secrets Manager:

```bash
# Store secrets
aws secretsmanager create-secret \
  --name edugo-api-mobile/prod \
  --secret-string '{
    "DATABASE_POSTGRES_PASSWORD": "...",
    "DATABASE_MONGODB_URI": "...",
    "MESSAGING_RABBITMQ_URL": "...",
    "STORAGE_S3_ACCESS_KEY_ID": "...",
    "STORAGE_S3_SECRET_ACCESS_KEY": "..."
  }'
```

Then configure your deployment to inject these as environment variables.

### Kubernetes Secrets

For Kubernetes deployments:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: edugo-secrets
type: Opaque
stringData:
  DATABASE_POSTGRES_PASSWORD: "..."
  DATABASE_MONGODB_URI: "..."
  MESSAGING_RABBITMQ_URL: "..."
  STORAGE_S3_ACCESS_KEY_ID: "..."
  STORAGE_S3_SECRET_ACCESS_KEY: "..."
```

Reference in deployment:

```yaml
env:
  - name: DATABASE_POSTGRES_PASSWORD
    valueFrom:
      secretKeyRef:
        name: edugo-secrets
        key: DATABASE_POSTGRES_PASSWORD
```

## Configuration Management Tools

### Validate Configuration

```bash
make config-validate
# or
./bin/configctl validate
```

### Generate Documentation

```bash
make config-docs
# or
./bin/configctl generate-docs
```

## Troubleshooting

### Application fails to start with "Configuration validation failed"

Check that all required environment variables are set:

```bash
# Required variables:
DATABASE_POSTGRES_PASSWORD
DATABASE_MONGODB_URI
MESSAGING_RABBITMQ_URL
STORAGE_S3_ACCESS_KEY_ID
STORAGE_S3_SECRET_ACCESS_KEY
```

### Environment variables not being loaded

1. Verify `.env` file exists
2. Check file format (no spaces around `=`)
3. Restart your IDE/terminal
4. For Docker: `docker-compose down && docker-compose up`

### Configuration not updating

1. Restart the application
2. Check `APP_ENV` is set correctly
3. Verify the correct config file is being loaded
4. Check environment variable precedence

## Security Best Practices

1. **Never commit `.env` files** - They are in `.gitignore`
2. **Never put secrets in YAML files** - Use environment variables
3. **Use strong passwords** - Generate with `openssl rand -base64 32`
4. **Rotate secrets regularly** - Especially in production
5. **Use cloud secret managers** - For production deployments
6. **Limit access** - Only give access to those who need it

## Adding New Configuration Variables

To add a new configuration variable, use the configctl tool:

```bash
# Public variable (goes in YAML)
./bin/configctl add database.redis.host \
  --type string \
  --default localhost \
  --desc "Redis host"

# Secret variable (goes in ENV)
./bin/configctl add auth.jwt.secret \
  --type string \
  --secret \
  --desc "JWT signing secret"
```

Or manually:

1. Add field to struct in `internal/config/config.go`
2. Add to `bindEnvVars()` in `internal/config/loader.go`
3. Add to `setDefaults()` if needed
4. Add validation in `internal/config/validator.go` if required
5. Update YAML files or `.env.example`
6. Update this documentation

## Support

For questions or issues with configuration:

1. Check this documentation
2. Review `.env.example` for required variables
3. Run `make config-validate` to check configuration files
4. Check application logs for specific error messages
