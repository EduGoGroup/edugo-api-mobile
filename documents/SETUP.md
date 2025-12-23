# 🚀 Guía de Instalación y Configuración

## Requisitos Previos

### Software Requerido

| Software | Versión Mínima | Propósito |
|----------|----------------|-----------|
| **Go** | 1.25+ | Compilación y ejecución |
| **Docker** | 20.10+ | Contenedores de servicios |
| **Docker Compose** | 2.0+ | Orquestación local |
| **Make** | 3.81+ | Comandos de desarrollo |
| **Git** | 2.30+ | Control de versiones |

### Herramientas Opcionales

| Herramienta | Propósito |
|-------------|-----------|
| `swag` | Generar documentación Swagger |
| `golangci-lint` | Análisis estático de código |
| `entr` | Watch mode para tests |

---

## 🏃 Quick Start (5 minutos)

```bash
# 1. Clonar repositorio
git clone https://github.com/EduGoGroup/edugo-api-mobile.git
cd edugo-api-mobile

# 2. Copiar configuración
cp .env.example .env

# 3. Levantar servicios (PostgreSQL, MongoDB, RabbitMQ)
docker-compose up -d postgres mongodb rabbitmq

# 4. Esperar a que los servicios estén healthy (~30s)
docker-compose ps

# 5. Ejecutar API
make run

# 6. Verificar
curl http://localhost:8080/health
```

**Swagger UI:** http://localhost:8080/swagger/index.html

---

## 📁 Configuración de Variables de Entorno

### Crear archivo .env

```bash
cp .env.example .env
```

### Variables Requeridas

```bash
# ========================================
# DATABASE SECRETS
# ========================================

# PostgreSQL Password
DATABASE_POSTGRES_PASSWORD=edugo123

# MongoDB Connection URI
DATABASE_MONGODB_URI=mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin

# ========================================
# MESSAGING SECRETS
# ========================================

# RabbitMQ Connection URL
MESSAGING_RABBITMQ_URL=amqp://edugo:edugo123@localhost:5672/

# ========================================
# STORAGE SECRETS (AWS S3)
# ========================================

# Para desarrollo local, puede usar LocalStack o dejar vacío
STORAGE_S3_ACCESS_KEY_ID=test
STORAGE_S3_SECRET_ACCESS_KEY=test
STORAGE_S3_BUCKET_NAME=edugo-materials

# ========================================
# AUTHENTICATION
# ========================================

# JWT Secret (debe coincidir con api-admin)
JWT_SECRET=dev-secret-key-change-in-production

# ========================================
# AMBIENTE
# ========================================

# Opciones: local, dev, qa, prod
APP_ENV=local
```

### Variables por Ambiente

| Variable | Local | Dev | Prod |
|----------|-------|-----|------|
| `APP_ENV` | local | dev | prod |
| `DATABASE_POSTGRES_PASSWORD` | edugo123 | **** | **** |
| `JWT_SECRET` | dev-secret | **** | **** |

---

## 🐳 Docker Compose

### Servicios Disponibles

```yaml
services:
  postgres:      # PostgreSQL 16
  mongodb:       # MongoDB 7.0
  rabbitmq:      # RabbitMQ 3.12 + Management UI
  api-mobile:    # Esta API (opcional)
```

### Comandos Docker

```bash
# Levantar todos los servicios
docker-compose up -d

# Levantar solo infraestructura (sin API)
docker-compose up -d postgres mongodb rabbitmq

# Ver logs
docker-compose logs -f

# Ver estado
docker-compose ps

# Detener servicios
docker-compose down

# Detener y borrar volúmenes (⚠️ borra datos)
docker-compose down -v
```

### Puertos Expuestos

| Servicio | Puerto | URL |
|----------|--------|-----|
| API Mobile | 8080 | http://localhost:8080 |
| PostgreSQL | 5432 | localhost:5432 |
| MongoDB | 27017 | localhost:27017 |
| RabbitMQ AMQP | 5672 | localhost:5672 |
| RabbitMQ UI | 15672 | http://localhost:15672 |

---

## ⚙️ Archivos de Configuración

### Estructura de Configuración

```
config/
├── config.yaml          # Configuración base (todos los ambientes)
├── config-local.yaml    # Override para local
├── config-dev.yaml      # Override para desarrollo
├── config-qa.yaml       # Override para QA
└── config-prod.yaml     # Override para producción
```

### config.yaml (Base)

```yaml
server:
  port: 8080
  host: "localhost"
  read_timeout: 30s
  write_timeout: 30s

database:
  postgres:
    host: "localhost"
    port: 5432
    database: "edugo"
    user: "edugo"
    max_connections: 25
    ssl_mode: "disable"

  mongodb:
    database: "edugo"
    timeout: 10s

messaging:
  rabbitmq:
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

logging:
  level: "info"
  format: "json"

auth:
  jwt:
    issuer: "edugo-central"
  api_admin:
    timeout: 5s
    cache_ttl: 60s
    cache_enabled: true
    remote_enabled: false
    fallback_enabled: false

bootstrap:
  optional_resources:
    rabbitmq: true    # RabbitMQ es opcional
    s3: true          # S3 es opcional

development:
  use_mock_repositories: false
```

### Precedencia de Configuración

```
1. Variables de entorno (.env)     ← Mayor prioridad
2. config-{APP_ENV}.yaml
3. config.yaml                     ← Menor prioridad
```

---

## 🛠️ Comandos Make

### Desarrollo

```bash
# Ejecutar en modo desarrollo
make run

# Desarrollo completo (deps + swagger + run)
make dev

# Build binario
make build

# Build para debugging
make build-debug
```

### Testing

```bash
# Todos los tests
make test

# Solo tests unitarios (rápido)
make test-unit

# Solo tests de integración (con testcontainers)
make test-integration

# Reporte de cobertura
make coverage-report

# Verificar cobertura mínima (60%)
make coverage-check

# Benchmarks
make benchmark
```

### Calidad de Código

```bash
# Formatear código
make fmt

# Análisis estático
make vet

# Linter completo
make lint

# Auditoría completa
make audit
```

### Docker

```bash
# Build imagen Docker
make docker-build

# Levantar con compose
make docker-up

# Detener servicios
make docker-down

# Ver logs
make docker-logs
```

### Swagger

```bash
# Regenerar documentación Swagger
make swagger
```

### Utilidades

```bash
# Limpiar archivos generados
make clean

# Información del proyecto
make info

# Pipeline CI completo
make ci

# Pre-commit (rápido)
make pre-commit
```

---

## 🔧 Modos de Ejecución

### 1. Modo Normal (con bases de datos reales)

```bash
# 1. Levantar infraestructura
docker-compose up -d postgres mongodb rabbitmq

# 2. Ejecutar API
make run
```

**Requiere:** PostgreSQL, MongoDB, RabbitMQ running

### 2. Modo Mock (sin bases de datos)

```yaml
# config/config-local.yaml
development:
  use_mock_repositories: true
```

```bash
# Ejecutar sin Docker
make run
```

**Ideal para:** Desarrollo frontend, pruebas rápidas

### 3. Modo Docker Completo

```bash
# Levantar todo incluyendo la API
docker-compose up -d

# API disponible en puerto 9090
curl http://localhost:9090/health
```

---

## 🔌 Configuración de Servicios Externos

### AWS S3 (Producción)

```bash
# .env
STORAGE_S3_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
STORAGE_S3_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
STORAGE_S3_BUCKET_NAME=edugo-materials-prod
```

### LocalStack (Desarrollo Local)

```yaml
# config/config-local.yaml
storage:
  s3:
    endpoint: "http://localhost:4566"  # LocalStack
```

```bash
# Levantar LocalStack
docker run -d -p 4566:4566 localstack/localstack
```

### RabbitMQ Management UI

- **URL:** http://localhost:15672
- **Usuario:** edugo
- **Password:** edugo123

---

## 🧪 Ejecutar Tests

### Tests Unitarios

```bash
# Rápido, sin dependencias externas
make test-unit
```

### Tests de Integración

```bash
# Usa testcontainers (levanta containers automáticamente)
make test-integration
```

### Cobertura

```bash
# Generar reporte HTML
make coverage-report

# Abrir reporte
open coverage/coverage.html
```

---

## 🐛 Troubleshooting

### Error: "connection refused" a PostgreSQL

```bash
# Verificar que el contenedor está running
docker-compose ps postgres

# Ver logs
docker-compose logs postgres

# Reiniciar servicio
docker-compose restart postgres
```

### Error: "MongoDB authentication failed"

```bash
# Verificar URI en .env
DATABASE_MONGODB_URI=mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin

# Verificar con mongosh
docker exec -it edugo-mongodb mongosh -u edugo -p edugo123 --authenticationDatabase admin
```

### Error: "swag: command not found"

```bash
# Instalar swag
go install github.com/swaggo/swag/cmd/swag@latest

# Verificar PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### Puerto 8080 en uso

```bash
# Encontrar proceso
lsof -i :8080

# Matar proceso
kill -9 <PID>

# O cambiar puerto en config
# config/config-local.yaml
server:
  port: 8081
```

### RabbitMQ no conecta

RabbitMQ es **opcional** por defecto. Si necesitas verificar:

```bash
# Ver logs
docker-compose logs rabbitmq

# Verificar health
curl -u edugo:edugo123 http://localhost:15672/api/healthchecks/node
```

---

## 📋 Checklist de Setup

- [ ] Go 1.25+ instalado
- [ ] Docker y Docker Compose instalados
- [ ] Repositorio clonado
- [ ] `.env` creado desde `.env.example`
- [ ] `docker-compose up -d` ejecutado
- [ ] Servicios healthy (`docker-compose ps`)
- [ ] `make run` ejecuta sin errores
- [ ] `curl localhost:8080/health` retorna `healthy`
- [ ] Swagger UI accesible en `/swagger/index.html`

---

## 🔗 Referencias

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [GORM](https://gorm.io/docs/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [Docker Compose](https://docs.docker.com/compose/)
