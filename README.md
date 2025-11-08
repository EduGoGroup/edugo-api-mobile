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

- **Lenguaje**: Go 1.25.3
- **Framework Web**: Gin
- **Documentación API**: Swagger/OpenAPI (Swaggo)
- **Base de Datos**: PostgreSQL + MongoDB (mock)
- **Autenticación**: JWT (mock)

## Requisitos Previos

- Go 1.25.3
- PostgreSQL 12+
- MongoDB 5.0+
- RabbitMQ 3.12+
- Docker (opcional, para desarrollo local)

## Configuración

### Setup Rápido

1. **Copiar el archivo de ejemplo:**
   ```bash
   cp .env.example .env
   ```

2. **Editar `.env` con tus valores:**
   ```bash
   # Database
   DATABASE_POSTGRES_PASSWORD=your-password
   DATABASE_MONGODB_URI=mongodb://user:pass@localhost:27017/edugo?authSource=admin
   
   # Messaging
   MESSAGING_RABBITMQ_URL=amqp://user:pass@localhost:5672/
   
   # Storage
   STORAGE_S3_ACCESS_KEY_ID=your-aws-key
   STORAGE_S3_SECRET_ACCESS_KEY=your-aws-secret
   
   # Application
   APP_ENV=local
   ```

3. **Ejecutar la aplicación:**
   ```bash
   make run
   # o
   go run cmd/main.go
   # o
   docker-compose up
   ```

### Variables Requeridas

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `DATABASE_POSTGRES_PASSWORD` | Contraseña de PostgreSQL | `your-secure-password` |
| `DATABASE_MONGODB_URI` | URI de conexión MongoDB | `mongodb://user:pass@host:27017/edugo?authSource=admin` |
| `MESSAGING_RABBITMQ_URL` | URL de RabbitMQ | `amqp://user:pass@host:5672/` |
| `STORAGE_S3_ACCESS_KEY_ID` | AWS S3 Access Key | `AKIAIOSFODNN7EXAMPLE` |
| `STORAGE_S3_SECRET_ACCESS_KEY` | AWS S3 Secret Key | `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY` |
| `APP_ENV` | Ambiente (local/dev/qa/prod) | `local` |

**📖 Para documentación completa de configuración, ver [CONFIG.md](CONFIG.md)**

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

# 2. Instalar swag CLI para documentación Swagger
go install github.com/swaggo/swag/cmd/swag@latest

# 3. Configurar variables de entorno
cp .env.example .env
# Editar .env con valores reales

# 4. Asegurar que PostgreSQL, MongoDB y RabbitMQ están corriendo

# 5. Ejecutar servidor (Swagger se regenera automáticamente)
go run cmd/main.go
```

**Nota**: Ya no es necesario ejecutar `swag init` manualmente. La aplicación regenera automáticamente la documentación Swagger al iniciar.

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

**Características de Swagger:**
- ✅ **Regeneración automática**: La documentación se actualiza automáticamente al iniciar la aplicación
- ✅ **Detección dinámica de puerto**: Swagger UI detecta automáticamente el puerto en el que corre la aplicación
- ✅ **Pruebas directas**: Puedes probar todos los endpoints directamente desde la interfaz

**Nota**: La primera vez que ejecutes la aplicación, asegúrate de tener instalado `swag` CLI (ver sección de Instalación)

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

## Documentación Swagger

### Regeneración Automática

La documentación Swagger se regenera automáticamente cada vez que inicias la aplicación. Esto garantiza que la documentación siempre esté actualizada con los últimos cambios en el código.

**Requisitos**:
- Tener instalado `swag` CLI: `go install github.com/swaggo/swag/cmd/swag@latest`
- Asegurarse de que `swag` esté en tu PATH

### Regeneración Manual (Opcional)

Si necesitas regenerar la documentación manualmente:

```bash
swag init -g cmd/main.go -o docs
```

### Acceso a Swagger UI

La interfaz de Swagger UI está disponible en:
```
http://localhost:{PORT}/swagger/index.html
```

Donde `{PORT}` es el puerto configurado en tu archivo de configuración (por defecto 8080).

**Características**:
- Detección automática del puerto y host
- Prueba de endpoints directamente desde la interfaz
- Documentación completa de todos los endpoints, parámetros y respuestas
- Soporte para autenticación Bearer token

### Troubleshooting

#### Error: "swag: command not found"

**Problema**: La aplicación no puede encontrar el comando `swag`.

**Solución**:
```bash
# Instalar swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Verificar que está en el PATH
which swag

# Si no está en el PATH, agregar $GOPATH/bin a tu PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

#### Swagger UI no carga o muestra errores

**Problema**: La interfaz de Swagger no se carga correctamente.

**Soluciones**:
1. Verificar que la aplicación esté corriendo: `curl http://localhost:8080/health`
2. Verificar que los archivos de documentación existan: `ls -la docs/`
3. Revisar los logs de la aplicación para errores de regeneración
4. Regenerar manualmente: `swag init -g cmd/main.go -o docs`

#### Los endpoints no funcionan desde Swagger UI

**Problema**: Al hacer clic en "Try it out", las peticiones fallan.

**Soluciones**:
1. Verificar que el puerto en la URL coincida con el puerto de la aplicación
2. Para endpoints protegidos, hacer clic en "Authorize" e ingresar el token Bearer
3. Verificar que CORS esté configurado correctamente
4. Revisar la consola del navegador para errores de red

#### La documentación no refleja cambios recientes

**Problema**: Los cambios en las anotaciones Swagger no aparecen en la UI.

**Soluciones**:
1. Reiniciar la aplicación (la regeneración es automática)
2. Limpiar caché del navegador y recargar Swagger UI
3. Verificar que las anotaciones Swagger estén correctamente formateadas
4. Regenerar manualmente: `swag init -g cmd/main.go -o docs`

#### Advertencia: "no se pudo regenerar Swagger"

**Problema**: La aplicación muestra una advertencia al iniciar.

**Causa**: `swag` no está instalado o no está en el PATH.

**Impacto**: La aplicación continúa funcionando con la documentación existente.

**Solución**: Instalar `swag` CLI como se indica arriba.

## Puerto

Por defecto: `8080`

Para cambiar, editar en `cmd/main.go`:
```go
port := ":8080"
```

## Licencia

MIT
