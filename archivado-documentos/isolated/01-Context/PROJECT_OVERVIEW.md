# PROJECT OVERVIEW - API Mobile

## Información General

**Proyecto:** EduGo API Mobile  
**Tipo:** API REST Microservicio  
**Puerto:** 8080  
**Lenguaje:** Go 1.21+  
**Marco de Trabajo:** Gin Framework  
**Especificación de Origen:** spec-01-evaluaciones  
**Estado:** En Desarrollo (Sprint 1/6)

---

## Propósito del Proyecto

API REST especializada en la gestión del **Sistema de Evaluaciones** de EduGo. Proporciona endpoints para que la aplicación móvil administre, realice y controle evaluaciones académicas.

### Responsabilidades Principales
- Gestión completa de evaluaciones (creación, edición, eliminación)
- Sistema de preguntas (múltiple choice, verdadero/falso, desarrollo)
- Asignación de evaluaciones a estudiantes
- Seguimiento de respuestas y progreso
- Generación de reportes de evaluación
- Integración con Worker para generación automática de quizzes

---

## Arquitectura del Proyecto

### Stack Tecnológico
```
┌─────────────────────────────────────────────────────┐
│ API Mobile (Gin Framework - Go)                     │
├─────────────────────────────────────────────────────┤
│ Layers:                                             │
│ ├─ HTTP Handlers (routes + middleware)             │
│ ├─ Service Layer (lógica de negocio)               │
│ ├─ Repository Layer (acceso a datos)               │
│ └─ Domain Models (estructuras de datos)             │
├─────────────────────────────────────────────────────┤
│ Dependencias Externas:                              │
│ ├─ PostgreSQL 15+ (datos relacionales)             │
│ ├─ MongoDB 7.0+ (almacenamiento de resultados)     │
│ ├─ RabbitMQ 3.12+ (comunicación con Worker)        │
│ └─ shared v1.3.0+ (librerías compartidas)          │
└─────────────────────────────────────────────────────┘
```

### Estructura de Carpetas
```
api-mobile/
├── cmd/
│   └── api-mobile/
│       └── main.go              # Punto de entrada
├── internal/
│   ├── handlers/                # HTTP handlers
│   ├── services/                # Lógica de negocio
│   ├── repositories/            # Acceso a datos
│   ├── models/                  # Estructuras de datos
│   ├── middleware/              # Middleware de Gin
│   └── config/                  # Configuración
├── pkg/
│   └── [estructuras compartidas]
├── migrations/                  # Migraciones de BD
├── docker/
│   └── Dockerfile
├── go.mod
├── go.sum
└── [configuración]
```

---

## Responsabilidades por Módulo

### 1. Evaluaciones (Core)
- **Crear evaluación** → POST /api/v1/evaluaciones
- **Listar evaluaciones** → GET /api/v1/evaluaciones
- **Obtener detalle** → GET /api/v1/evaluaciones/:id
- **Editar evaluación** → PUT /api/v1/evaluaciones/:id
- **Eliminar evaluación** → DELETE /api/v1/evaluaciones/:id
- **Publicar/Cerrar** → POST /api/v1/evaluaciones/:id/publish

### 2. Preguntas
- **Crear pregunta** → POST /api/v1/evaluaciones/:id/preguntas
- **Listar preguntas** → GET /api/v1/evaluaciones/:id/preguntas
- **Editar pregunta** → PUT /api/v1/preguntas/:id
- **Eliminar pregunta** → DELETE /api/v1/preguntas/:id
- **Reordenar preguntas** → POST /api/v1/evaluaciones/:id/reorder

### 3. Asignaciones
- **Asignar a estudiantes** → POST /api/v1/evaluaciones/:id/assign
- **Listar asignaciones** → GET /api/v1/evaluaciones/:id/assignments
- **Cambiar fecha límite** → PUT /api/v1/assignments/:id
- **Marcar como completada** → POST /api/v1/assignments/:id/complete

### 4. Respuestas (Capturas)
- **Enviar respuestas** → POST /api/v1/evaluaciones/:id/submit
- **Obtener respuestas guardadas** → GET /api/v1/evaluaciones/:id/draft
- **Validar respuestas** → Lógica interna en Service

### 5. Resultados
- **Obtener resultados** → GET /api/v1/evaluaciones/:id/results
- **Detalles de respuesta** → GET /api/v1/assignments/:id/answer-detail
- **Exportar resultados** → GET /api/v1/evaluaciones/:id/results/export

### 6. Integración con Worker
- **Solicitar quiz automático** → POST /api/v1/evaluaciones/material/:id/generate-quiz
- **Obtener estado** → GET /api/v1/requests/:id/status
- **Resultados generados** → Consumir mensajes de RabbitMQ

---

## Flujos Principales

### Flujo 1: Crear y Ejecutar Evaluación Manual
```
1. Docente crea evaluación (POST /evaluaciones)
2. Docente agrega preguntas (POST /preguntas)
3. Docente asigna a estudiantes (POST /assign)
4. Estudiante abre evaluación en app móvil
5. Estudiante responde preguntas
6. Estudiante envía evaluación (POST /submit)
7. API valida respuestas
8. API calcula resultados
9. Docente revisa resultados en panel administrativo
```

### Flujo 2: Generar Quiz Automático con IA
```
1. Docente sube material educativo
2. Docente solicita generar quiz (POST /generate-quiz)
3. API publica mensaje a RabbitMQ (Worker)
4. Worker procesa material con OpenAI
5. Worker genera preguntas
6. Worker publica respuesta a RabbitMQ
7. API recibe y guarda preguntas en MongoDB
8. API retorna ID de evaluación creada
9. Flujo normal de evaluación (pasos 3-9 del Flujo 1)
```

### Flujo 3: Seguimiento de Progreso
```
1. Estudiante responde evaluación (POST /submit)
2. API calcula puntuación
3. API actualiza progreso en PostgreSQL
4. API publica evento a RabbitMQ (auditoría)
5. App móvil solicita progreso (GET /progress)
6. API retorna datos actualizados
```

---

## Entidades de Base de Datos

### PostgreSQL (Datos Relacionales)

```sql
-- Tabla de Evaluaciones
CREATE TABLE evaluations (
  id SERIAL PRIMARY KEY,
  material_id INT,
  title VARCHAR(255),
  description TEXT,
  type ENUM('manual', 'generated'),
  status ENUM('draft', 'published', 'closed'),
  passing_score INT,
  created_by INT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- Tabla de Preguntas
CREATE TABLE questions (
  id SERIAL PRIMARY KEY,
  evaluation_id INT REFERENCES evaluations(id),
  type ENUM('multiple_choice', 'true_false', 'short_answer'),
  text TEXT,
  position INT,
  points INT,
  created_at TIMESTAMP
);

-- Tabla de Opciones (para multiple choice)
CREATE TABLE question_options (
  id SERIAL PRIMARY KEY,
  question_id INT REFERENCES questions(id),
  text VARCHAR(500),
  is_correct BOOLEAN,
  position INT
);

-- Tabla de Asignaciones
CREATE TABLE evaluation_assignments (
  id SERIAL PRIMARY KEY,
  evaluation_id INT REFERENCES evaluations(id),
  student_id INT,
  assigned_at TIMESTAMP,
  due_date TIMESTAMP,
  status ENUM('pending', 'in_progress', 'submitted', 'graded')
);

-- Tabla de Respuestas (draft)
CREATE TABLE answer_drafts (
  id SERIAL PRIMARY KEY,
  assignment_id INT REFERENCES evaluation_assignments(id),
  question_id INT REFERENCES questions(id),
  answer TEXT,
  saved_at TIMESTAMP
);
```

### MongoDB (Resultados Evaluaciones)

```javascript
// Colección: evaluation_results
{
  _id: ObjectId,
  evaluation_id: 1,
  assignment_id: 1,
  student_id: 42,
  answers: [
    {
      question_id: 101,
      answer: "Option B",
      is_correct: true,
      points_earned: 5,
      timestamp: "2025-11-15T10:30:00Z"
    }
  ],
  total_score: 45,
  max_score: 50,
  percentage: 90,
  status: "graded",
  submitted_at: "2025-11-15T10:45:00Z",
  feedback: "Excelente desempeño"
}
```

---

## Configuración Requerida

### Variables de Entorno (.env)
```bash
# Base de datos PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=edugo_user
DB_PASSWORD=secure_password
DB_NAME=edugo_mobile

# Base de datos MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=edugo_assessments

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest

# Shared Library
SHARED_LOG_LEVEL=info
SHARED_CONTEXT_TIMEOUT=30s

# API Configuration
API_PORT=8080
API_ENV=development
API_TIMEOUT=30s

# Authentication
JWT_SECRET=your_secret_key
JWT_EXPIRY=24h
```

### Docker Compose (desarrollo)
```yaml
services:
  api-mobile:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - MONGO_URI=mongodb://mongo:27017
      - RABBITMQ_HOST=rabbitmq
    depends_on:
      - postgres
      - mongo
      - rabbitmq
```

---

## Dependencias del Proyecto

### Dependencias Go Principales
```
go get github.com/gin-gonic/gin@latest              # Web framework
go get gorm.io/gorm@latest                          # ORM
go get gorm.io/driver/postgres@latest               # PostgreSQL driver
go get go.mongodb.org/mongo-driver@latest           # MongoDB driver
go get github.com/streadway/amqp@latest             # RabbitMQ client
go get github.com/spf13/viper@latest                # Configuration management
go get github.com/golang-jwt/jwt/v5@latest          # JWT
```

### Dependencias Internas
- **shared v1.3.0+**
  - `logger` - Logging centralizado
  - `database` - Conexión a PostgreSQL
  - `auth` - Validación JWT
  - `messaging` - Cliente RabbitMQ
  - `models` - Estructuras compartidas

### Dependencias Externas
- **PostgreSQL 15+** - Base de datos relacional
- **MongoDB 7.0+** - Almacenamiento documentos
- **RabbitMQ 3.12+** - Message broker
- **Worker microservicio** - Generación de quizzes con IA

---

## Compilación y Despliegue

### Compilación Local
```bash
# Descargar dependencias
go mod download
go mod tidy

# Compilar ejecutable
go build -o api-mobile ./cmd/api-mobile

# Ejecutar
./api-mobile
```

### Compilación Docker
```bash
# Build imagen
docker build -t edugo/api-mobile:latest -f docker/Dockerfile .

# Run contenedor
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e MONGO_URI=mongodb://mongo:27017 \
  edugo/api-mobile:latest
```

### Migraciones de Base de Datos
```bash
# Crear/actualizar esquema
go run ./cmd/migrate

# Ver estado migraciones
go run ./cmd/migrate status
```

---

## Testing

### Pruebas Unitarias
```bash
# Ejecutar todas
go test ./...

# Con cobertura
go test -cover ./...

# Verbose
go test -v ./...
```

### Pruebas de Integración
```bash
# Requiere servicios levantados (PostgreSQL, MongoDB, RabbitMQ)
go test -tags=integration ./...
```

### Documentación API (Swagger)
```bash
# Generar documentación
swag init -g cmd/api-mobile/main.go

# Acceder en http://localhost:8080/swagger/index.html
```

---

## Métricas y Monitoreo

### Healthcheck
```
GET /api/v1/health
Respuesta: {"status": "ok", "timestamp": "2025-11-15T10:30:00Z"}
```

### Readiness
```
GET /api/v1/ready
Respuesta: {"ready": true, "services": {"postgres": "ok", "mongo": "ok", "rabbitmq": "ok"}}
```

### Métricas Prometheus
```
GET /metrics
- http_requests_total
- http_request_duration_seconds
- db_query_duration_seconds
```

---

## Sprint Planning (6 Sprints)

| Sprint | Funcionalidad | Duración |
|--------|---------------|----------|
| 1 | Setup + Evaluaciones CRUD | 2 semanas |
| 2 | Preguntas + Opciones | 2 semanas |
| 3 | Asignaciones + Flujo básico | 2 semanas |
| 4 | Respuestas + Validación | 2 semanas |
| 5 | Resultados + Reportes | 2 semanas |
| 6 | Integración Worker + Optimizaciones | 2 semanas |

---

## Contacto y Referencias

- **Repositorio GitHub:** https://github.com/EduGoGroup/edugo-api-mobile
- **Especificación Completa:** docs/ESTADO_PROYECTO.md (repo análisis)
- **Documentación Técnica:** Este directorio (01-Context/)
