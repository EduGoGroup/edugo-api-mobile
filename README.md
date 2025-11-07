# EduGo API Mobile

API REST para operaciones frecuentes de docentes y estudiantes en la plataforma EduGo.

## Descripción

Esta API maneja:
- Autenticación de usuarios
- Gestión de materiales educativos (crear, leer, listar)
- Resúmenes generados por IA
- Cuestionarios y evaluaciones
- Seguimiento de progreso de estudiantes
- Estadísticas para docentes

## Tecnologías

- **Lenguaje**: Go 1.21+
- **Framework Web**: Gin
- **Documentación API**: Swagger/OpenAPI (Swaggo)
- **Base de Datos**: PostgreSQL + MongoDB (mock)
- **Autenticación**: JWT (mock)

## Requisitos Previos

- Go 1.21+
- PostgreSQL 12+
- MongoDB 5.0+
- RabbitMQ 3.12+
- Docker (opcional, para desarrollo local)

## Variables de Entorno

La aplicación requiere las siguientes variables de entorno para funcionar:

```bash
# Copiar archivo de ejemplo
cp .env.example .env

# Editar .env con tus valores reales
```

### Variables Requeridas

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `POSTGRES_PASSWORD` | Contraseña de PostgreSQL | `your-secure-password` |
| `MONGODB_URI` | URI de conexión MongoDB | `mongodb://user:pass@host:27017/edugo?authSource=admin` |
| `RABBITMQ_URL` | URL de RabbitMQ | `amqp://user:pass@host:5672/` |
| `JWT_SECRET` | Secret key para JWT | `$(openssl rand -base64 32)` |
| `APP_ENV` | Ambiente (local/dev/qa/prod) | `local` |

### Variables Opcionales

| Variable | Descripción | Default |
|----------|-------------|---------|
| `AWS_ACCESS_KEY_ID` | Access key de AWS S3 | - |
| `AWS_SECRET_ACCESS_KEY` | Secret key de AWS S3 | - |
| `POSTGRES_HOST` | Host de PostgreSQL | `localhost` |
| `POSTGRES_PORT` | Puerto de PostgreSQL | `5432` |

Ver el archivo [`.env.example`](.env.example) para la lista completa de variables y ejemplos.

## Instalación

### Desarrollo Local con Docker Compose

```bash
# 1. Configurar variables de entorno
cp .env.example .env
# Editar .env con valores reales

# 2. Iniciar servicios
docker-compose up -d

# La API estará disponible en http://localhost:9090
```

### Desarrollo Local sin Docker

```bash
# 1. Instalar dependencias
go mod download

# 2. Configurar variables de entorno
cp .env.example .env
# Editar .env con valores reales

# 3. Asegurar que PostgreSQL, MongoDB y RabbitMQ están corriendo

# 4. Generar documentación Swagger
swag init -g cmd/main.go -o docs

# 5. Ejecutar servidor
go run cmd/main.go
```

### Producción con Docker

```bash
# Opción 1: Usar imagen publicada
docker pull ghcr.io/edugogroup/edugo-api-mobile:latest

docker run -d \
  --name edugo-api-mobile \
  -p 8080:8080 \
  -e POSTGRES_PASSWORD=your-password \
  -e MONGODB_URI=mongodb://... \
  -e RABBITMQ_URL=amqp://... \
  -e JWT_SECRET=your-secret \
  -e APP_ENV=prod \
  ghcr.io/edugogroup/edugo-api-mobile:latest

# Opción 2: Build local
docker build -t edugo-api-mobile .
docker run -d -p 8080:8080 --env-file .env edugo-api-mobile
```

### Usando edugo-dev-environment

Si usas el [repositorio unificado de desarrollo](https://github.com/EduGoGroup/edugo-dev-environment):

```bash
# 1. Clonar el repositorio
git clone https://github.com/EduGoGroup/edugo-dev-environment.git
cd edugo-dev-environment

# 2. Configurar .env (usa su propio .env)
cp .env.example .env
# Editar .env según instrucciones del repositorio

# 3. Iniciar todos los servicios
docker-compose up -d

# La API Mobile estará en http://localhost:8081
```

## Uso

### Iniciar Servidor

```bash
go run cmd/main.go
```

El servidor estará disponible en `http://localhost:8080`

### Swagger UI

Accede a la documentación interactiva en:
```
http://localhost:8080/swagger/index.html
```

### Health Check

```bash
curl http://localhost:8080/health
```

## Endpoints Principales

### Autenticación

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/v1/auth/login` | Login y obtención de JWT |

### Materiales (requieren autenticación)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/v1/materials` | Listar materiales con filtros |
| POST | `/v1/materials` | Crear nuevo material |
| GET | `/v1/materials/:id` | Detalle de material + URL PDF |
| POST | `/v1/materials/:id/upload-complete` | Notificar upload completado |
| GET | `/v1/materials/:id/summary` | Obtener resumen generado |
| GET | `/v1/materials/:id/assessment` | Obtener quiz |
| POST | `/v1/materials/:id/assessment/attempts` | Enviar respuestas de quiz |
| PATCH | `/v1/materials/:id/progress` | Actualizar progreso |
| GET | `/v1/materials/:id/stats` | Estadísticas (solo docentes) |

## Autenticación

Incluir header en requests protegidos:

```
Authorization: Bearer {jwt_token}
```

Ejemplo:
```bash
curl -H "Authorization: Bearer eyJhbGci..." \
     http://localhost:8080/v1/materials
```

## Ejemplo de Uso

### 1. Login

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "docente@example.com",
    "password": "password123"
  }'
```

Respuesta:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "...",
  "expires_at": "2025-01-29T12:00:00Z",
  "user": {
    "id": "user-uuid-123",
    "name": "María González",
    "email": "docente@example.com",
    "role": "teacher"
  }
}
```

### 2. Listar Materiales

```bash
curl -H "Authorization: Bearer {token}" \
     "http://localhost:8080/v1/materials?unit_id=uuid-5a&status=new"
```

### 3. Crear Material

```bash
curl -X POST http://localhost:8080/v1/materials \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Introducción a Pascal",
    "description": "Material base sobre Pascal",
    "subject_id": "uuid-prog",
    "unit_ids": ["uuid-5a", "uuid-5b"]
  }'
```

## Estado Actual

**Implementación**: Código base con datos MOCK

Este código provee la estructura completa con:
- ✅ Rutas definidas
- ✅ Handlers con firmas correctas
- ✅ Modelos de request/response
- ✅ Documentación Swagger completa
- ✅ Middleware de autenticación (mock)
- ⏳ Datos MOCK (retornan datos estáticos)

### Próximos Pasos

Para convertir en código producción:

1. **Configuración**:
   - Agregar archivo `.env` con variables de entorno
   - Configurar conexiones a PostgreSQL, MongoDB, S3, RabbitMQ

2. **Servicios Reales**:
   - Implementar capa de servicios con lógica real
   - Implementar repositorios para PostgreSQL y MongoDB
   - Implementar cliente S3 para URLs firmadas
   - Implementar publicador de eventos RabbitMQ

3. **Autenticación Real**:
   - Generar y validar JWT real (ej: con `github.com/golang-jwt/jwt`)
   - Hash de contraseñas con bcrypt
   - Refresh tokens con Redis

4. **Validaciones**:
   - Validaciones de negocio
   - Manejo de errores robusto
   - Logging estructurado

5. **Testing**:
   - Unit tests para handlers y servicios
   - Integration tests con bases de datos de prueba
   - E2E tests

## Estructura del Proyecto

```
api-mobile/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── handlers/            # HTTP handlers
│   │   ├── auth.go
│   │   └── materials.go
│   ├── models/
│   │   ├── enum/            # Enums
│   │   ├── request/         # DTOs de request
│   │   └── response/        # DTOs de response
│   ├── services/            # Lógica de negocio (TODO)
│   └── middleware/          # Middleware HTTP
├── docs/                    # Swagger docs generados
├── go.mod
└── README.md
```

## Generar Swagger Docs

Después de modificar anotaciones Swagger:

```bash
swag init -g cmd/main.go -o docs
```

## Puerto

Por defecto: `8080`

Para cambiar, editar en `cmd/main.go`:
```go
port := ":8080"
```

## Licencia

MIT
