# DEPENDENCIES - API Mobile

## Matriz de Dependencias

```
┌──────────────────────────────────────────────────────────┐
│                   API MOBILE                             │
├──────────────────────────────────────────────────────────┤
│ Dependencias Críticas (sin estas NO funciona)            │
│ ├─ SHARED v1.3.0+                                        │
│ ├─ PostgreSQL 15+                                        │
│ └─ MongoDB 7.0+ (para persistencia de resultados)        │
│                                                          │
│ Dependencias de Integración (para funcionalidad completa)│
│ ├─ RabbitMQ 3.12+                                        │
│ ├─ WORKER (procesamiento async IA)                       │
│ └─ API ADMIN (jerarquía académica)                       │
│                                                          │
│ Dependencias de Desarrollo/Ops                           │
│ ├─ Docker                                                │
│ ├─ Docker Compose                                        │
│ └─ Go 1.21+                                              │
└──────────────────────────────────────────────────────────┘
```

---

## Dependencias Críticas

### 0. edugo-infrastructure v0.1.1 (NUEVO)

**¿Qué es?** Infraestructura compartida (migraciones BD, schemas eventos, Docker Compose)

**⚠️ VERSIÓN CORRECTA:** v0.1.1  
**Estado:** ✅ COMPLETADO (96%)

**Qué usa api-mobile:**

#### Migraciones de Base de Datos
```bash
# infrastructure/database/migrations/
003_create_materials.up.sql      # Tabla materials
004_create_assessment.up.sql     # Tablas assessment, attempt, answer
006_create_progress.up.sql       # Tabla student_progress
008_add_indexes.up.sql           # Índices de performance
```

#### JSON Schemas de Eventos
```bash
# infrastructure/schemas/events/
material.uploaded.json           # Validar eventos que publica api-mobile
evaluation.submitted.json        # Validar eventos que publica api-mobile
```

#### Referencia a docker-compose
```bash
# infrastructure/docker/docker-compose.yml
# Puede usar directamente o copiar a dev-environment
```

**Integración:**
```go
// Validar eventos antes de publicar a RabbitMQ
import "github.com/EduGoGroup/edugo-infrastructure/schemas"

func PublishMaterialEvent(material Material) error {
    event := buildEvent(material)
    
    // Validar contra schema (cuando validator.go esté implementado)
    if err := schemas.Validate("material.uploaded", event); err != nil {
        return fmt.Errorf("invalid event: %w", err)
    }
    
    return publisher.Publish("material-events", "material.uploaded", event)
}
```

---

### 1. edugo-shared v0.7.0

**¿Qué es?** Librería Go compartida con utilidades reutilizables

**⚠️ VERSIÓN CORRECTA:** v0.7.0 (FROZEN hasta post-MVP)  
**❌ NO USAR:** v1.3.0, v1.4.0, v1.5.0 (no existen)

**Módulos de shared v0.7.0 que API Mobile necesita:**

#### NUEVO: evaluation Module (v0.7.0)
```go
// NUEVO en shared v0.7.0
import "github.com/EduGoGroup/edugo-shared/evaluation"

// Modelos de evaluación compartidos
type Assessment struct {
    ID          string
    MaterialID  int64
    Questions   []Question
    TotalPoints int
}

type Question struct {
    ID      string
    Text    string
    Type    QuestionType // MultipleChoice, TrueFalse, ShortAnswer
    Options []Option
    Points  int
}

// Usar en api-mobile
assessment := evaluation.Assessment{
    MaterialID: 42,
    Questions: []evaluation.Question{...},
}
```

**Beneficio:** Consistencia entre api-mobile y worker en modelos de evaluación

---

#### a) Logger Module
```go
// Importar desde SHARED
import "github.com/EduGoGroup/edugo-shared/logger"

// Uso en API Mobile
logger.Info("Evaluación creada", map[string]interface{}{
    "evaluation_id": eval.ID,
    "created_by": eval.CreatedBy,
    "timestamp": time.Now(),
})

logger.Error("Error al guardar evaluación", map[string]interface{}{
    "error": err.Error(),
    "evaluation_id": eval.ID,
})

logger.Debug("Consultando base de datos", map[string]interface{}{
    "table": "evaluations",
    "query": "SELECT * FROM evaluations WHERE id = $1",
})
```

**Configuración esperada:**
```go
// En config/logger.go
type LoggerConfig struct {
    Level    string // "debug", "info", "warn", "error"
    Format   string // "json", "text"
    Output   string // "stdout", "file", "datadog"
}

// Env vars
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout
```

**Línea de código en API Mobile:**
```go
package main

import (
    "github.com/EduGoGroup/edugo-shared/logger"
)

func init() {
    logger.Init(LoggerConfig{
        Level: os.Getenv("LOG_LEVEL"),
    })
}
```

---

#### b) Database Module (PostgreSQL)
```go
import "github.com/EduGoGroup/edugo-shared/database"

// Obtener conexión singleton
db := database.GetDB()

// Usar con GORM
var evaluations []Evaluation
db.Find(&evaluations)

// Con transaction
tx := db.BeginTx(ctx, nil)
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()
```

**Configuración:**
```go
// database/config.go
database.Init(database.Config{
    Host:     os.Getenv("DB_HOST"),
    Port:     os.Getenv("DB_PORT"),
    User:     os.Getenv("DB_USER"),
    Password: os.Getenv("DB_PASSWORD"),
    Database: os.Getenv("DB_NAME"),
})
```

**Variables de entorno requeridas:**
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=edugo_user
DB_PASSWORD=secure_password
DB_NAME=edugo_mobile
```

---

#### c) Auth Module (JWT)
```go
import "github.com/EduGoGroup/edugo-shared/auth"

// En middleware
middleware := auth.NewJWTValidator()

// En Gin handlers
func GetEvaluation(c *gin.Context) {
    claims, err := auth.ExtractClaims(c.Request)
    if err != nil {
        c.JSON(401, gin.H{"error": "Unauthorized"})
        return
    }
    
    userID := claims.UserID
    schoolID := claims.SchoolID
}
```

**Validación automática:**
```go
// En router setup
router.Use(auth.ValidateToken())
```

---

#### d) Messaging Module (RabbitMQ)
```go
import "github.com/EduGoGroup/edugo-shared/messaging"

// Publisher: Solicitar quiz automático
publisher := messaging.NewPublisher()
err := publisher.Publish(
    "assessment.requests", // exchange
    "worker.assessment.requests", // routing key
    payload, // []byte
)

// Subscriber: Recibir quiz completado
subscriber := messaging.NewSubscriber()
messages := subscriber.Subscribe(
    "assessment.responses", // exchange
    "api-mobile.assessment.responses", // queue
)

for msg := range messages {
    var response AssessmentResponse
    json.Unmarshal(msg.Body, &response)
    // Procesar respuesta
}
```

**Configuración:**
```bash
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_VHOST=/
```

---

#### e) Models Module (Estructuras compartidas)
```go
import "github.com/EduGoGroup/edugo-shared/models"

// User struct compartido
var user models.User

// School struct compartido
var school models.School

// Permite consistencia entre servicios
```

---

### 2. PostgreSQL 15+

**¿Qué es?** Base de datos relacional para datos transaccionales

**Requisitos:**
- Versión: 15.0 o superior
- Puerto: 5432 (default)
- Máximo 100 conexiones simultáneas (pool)
- Backup diario automático

**Bases de datos a crear:**
```sql
-- Base de datos para API Mobile
CREATE DATABASE edugo_mobile OWNER edugo_user;

-- Crear extensiones
\c edugo_mobile
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
```

**Tablas requeridas en API Mobile:**

```sql
-- 1. Evaluaciones
CREATE TABLE evaluations (
    id BIGSERIAL PRIMARY KEY,
    material_id BIGINT,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50), -- 'manual', 'generated'
    status VARCHAR(50), -- 'draft', 'published', 'closed'
    passing_score INTEGER,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 2. Preguntas
CREATE TABLE questions (
    id BIGSERIAL PRIMARY KEY,
    evaluation_id BIGINT NOT NULL REFERENCES evaluations(id),
    type VARCHAR(50), -- 'multiple_choice', 'true_false', 'short_answer'
    text TEXT NOT NULL,
    position INTEGER,
    points INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Opciones de preguntas
CREATE TABLE question_options (
    id BIGSERIAL PRIMARY KEY,
    question_id BIGINT NOT NULL REFERENCES questions(id),
    text VARCHAR(500),
    is_correct BOOLEAN DEFAULT FALSE,
    position INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Asignaciones de evaluaciones
CREATE TABLE evaluation_assignments (
    id BIGSERIAL PRIMARY KEY,
    evaluation_id BIGINT NOT NULL REFERENCES evaluations(id),
    student_id BIGINT NOT NULL,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    due_date TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50), -- 'pending', 'in_progress', 'submitted', 'graded'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 5. Respuestas en borrador
CREATE TABLE answer_drafts (
    id BIGSERIAL PRIMARY KEY,
    assignment_id BIGINT NOT NULL REFERENCES evaluation_assignments(id),
    question_id BIGINT NOT NULL REFERENCES questions(id),
    answer TEXT,
    saved_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índices para performance
CREATE INDEX idx_evaluations_status ON evaluations(status);
CREATE INDEX idx_evaluations_created_by ON evaluations(created_by);
CREATE INDEX idx_questions_evaluation_id ON questions(evaluation_id);
CREATE INDEX idx_assignments_evaluation_id ON evaluation_assignments(evaluation_id);
CREATE INDEX idx_assignments_student_id ON evaluation_assignments(student_id);
CREATE INDEX idx_answer_drafts_assignment_id ON answer_drafts(assignment_id);
```

**Conexión desde API Mobile:**
```go
package database

import (
    "fmt"
    "github.com/EduGoGroup/edugo-shared/database"
    "os"
)

func InitPostgres() error {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )
    
    return database.Init(dsn)
}
```

---

### 3. MongoDB 7.0+

**¿Qué es?** Base de datos documental para almacenamiento flexible

**Requisitos:**
- Versión: 7.0 o superior
- Puerto: 27017 (default)
- Máximo 100 conexiones simultáneas
- Backup diario automático

**Colecciones requeridas:**

```javascript
// 1. Resultados de evaluaciones
db.createCollection("evaluation_results", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["evaluation_id", "assignment_id", "student_id"],
            properties: {
                _id: { bsonType: "objectId" },
                evaluation_id: { bsonType: "int" },
                assignment_id: { bsonType: "int" },
                student_id: { bsonType: "int" },
                answers: {
                    bsonType: "array",
                    items: {
                        bsonType: "object",
                        properties: {
                            question_id: { bsonType: "int" },
                            answer: { bsonType: "string" },
                            is_correct: { bsonType: "bool" },
                            points_earned: { bsonType: "int" }
                        }
                    }
                },
                total_score: { bsonType: "int" },
                max_score: { bsonType: "int" },
                percentage: { bsonType: "double" },
                submitted_at: { bsonType: "date" },
                feedback: { bsonType: "string" }
            }
        }
    }
});

// Crear índices
db.evaluation_results.createIndex({ "evaluation_id": 1 });
db.evaluation_results.createIndex({ "student_id": 1 });
db.evaluation_results.createIndex({ "submitted_at": -1 });

// 2. Auditoría
db.createCollection("evaluation_audit");
db.evaluation_audit.createIndex({ "evaluation_id": 1, "timestamp": -1 });
```

**Conexión desde API Mobile:**
```go
package database

import (
    "context"
    "github.com/EduGoGroup/edugo-shared/database"
    "os"
)

func InitMongoDB() error {
    mongoURI := os.Getenv("MONGO_URI")
    mongoDB := os.Getenv("MONGO_DB_NAME")
    
    return database.InitMongo(mongoURI, mongoDB)
}

// Uso en services
func SaveEvaluationResult(ctx context.Context, result *EvaluationResult) error {
    collection := database.GetMongoCollection("evaluation_results")
    _, err := collection.InsertOne(ctx, result)
    return err
}
```

**Variables de entorno:**
```bash
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=edugo_assessments
MONGO_USER=admin
MONGO_PASSWORD=secure_password
```

---

## Dependencias de Integración

### 4. RabbitMQ 3.12+

**¿Qué es?** Message broker para comunicación asíncrona entre servicios

**Requisitos:**
- Versión: 3.12 o superior
- Puerto: 5672 (AMQP)
- Puerto 15672 (Management UI)
- Max queues: 10000

**Exchanges y Queues a crear:**

```bash
# Exchange para requests de generación de quizzes
rabbitmqctl declare_exchange assessment.requests type:direct durable:true

# Queue para worker (consumer)
rabbitmqctl declare_queue worker.assessment.requests durable:true
rabbitmqctl bind_queue assessment.requests worker.assessment.requests worker.assessment.requests

# Exchange para responses de quizzes generados
rabbitmqctl declare_exchange assessment.responses type:direct durable:true

# Queue para API Mobile (consumer)
rabbitmqctl declare_queue api-mobile.assessment.responses durable:true
rabbitmqctl bind_queue assessment.responses api-mobile.assessment.responses api-mobile.assessment.responses
```

**Publicar mensaje (desde API Mobile):**
```go
package messaging

import (
    "encoding/json"
    "github.com/EduGoGroup/edugo-shared/messaging"
)

type GenerateQuizRequest struct {
    RequestID   string                 `json:"request_id"`
    Type        string                 `json:"type"` // "generate_assessment"
    MaterialID  int                    `json:"material_id"`
    MaterialURL string                 `json:"material_url"`
    Config      map[string]interface{} `json:"config"`
}

func PublishGenerateQuiz(req GenerateQuizRequest) error {
    publisher := messaging.NewPublisher()
    
    payload, _ := json.Marshal(req)
    
    return publisher.Publish(
        "assessment.requests",        // exchange
        "worker.assessment.requests", // routing key
        payload,
    )
}
```

**Consumir mensaje (desde API Mobile):**
```go
func ConsumeAssessmentResponses(ctx context.Context) {
    subscriber := messaging.NewSubscriber()
    
    messages := subscriber.Subscribe(
        "assessment.responses",
        "api-mobile.assessment.responses",
    )
    
    for {
        select {
        case msg := <-messages:
            var response AssessmentResponse
            json.Unmarshal(msg.Body, &response)
            
            // Procesar respuesta del Worker
            SaveGeneratedAssessment(response)
            msg.Ack(false)
            
        case <-ctx.Done():
            return
        }
    }
}
```

**Configuración:**
```bash
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_VHOST=/
RABBITMQ_MAX_RETRIES=3
RABBITMQ_RETRY_DELAY=5s
```

---

### 5. WORKER Microservicio

**¿Qué es?** Servicio que genera automáticamente preguntas de evaluación usando IA

**Versión:** Latest (sigue semántico versionamiento)

**Comunicación:**
- RabbitMQ: Recibe requests, envía responses
- No HTTP directo desde API Mobile

**Contrato de mensaje - Request:**
```json
{
  "request_id": "uuid-12345",
  "type": "generate_assessment",
  "material_id": 42,
  "material_path": "s3://edugo-bucket/materials/42/content.pdf",
  "config": {
    "num_questions": 10,
    "difficulty": "medium",
    "language": "es",
    "question_types": ["multiple_choice", "true_false"]
  },
  "timestamp": "2025-11-15T10:30:00Z"
}
```

**Contrato de mensaje - Response:**
```json
{
  "request_id": "uuid-12345",
  "status": "success",
  "evaluation_id": 1001,
  "questions_generated": 10,
  "questions": [
    {
      "id": "q-1",
      "text": "¿Cuál es la capital de España?",
      "type": "multiple_choice",
      "options": [
        {"text": "Madrid", "correct": true, "order": 1},
        {"text": "Barcelona", "correct": false, "order": 2}
      ],
      "points": 5,
      "explanation": "Madrid es la capital de España"
    }
  ],
  "timestamp": "2025-11-15T10:45:00Z"
}
```

**Caso de error:**
```json
{
  "request_id": "uuid-12345",
  "status": "error",
  "error_message": "Material format not supported",
  "error_code": "INVALID_FORMAT",
  "timestamp": "2025-11-15T10:45:00Z"
}
```

**Cómo API Mobile debería manejar:**
```go
func HandleAssessmentResponse(msg []byte) error {
    var response WorkerResponse
    json.Unmarshal(msg, &response)
    
    if response.Status == "error" {
        // Log error, notificar usuario
        logger.Error("Worker error", map[string]interface{}{
            "request_id": response.RequestID,
            "error": response.ErrorMessage,
        })
        
        // Actualizar request status en BD
        UpdateGenerationRequest(response.RequestID, "failed")
        return nil
    }
    
    // Guardar preguntas en PostgreSQL
    for _, q := range response.Questions {
        SaveQuestion(response.EvaluationID, q)
    }
    
    // Marcar como completado
    UpdateGenerationRequest(response.RequestID, "completed")
    return nil
}
```

---

### 6. API ADMIN Microservicio

**¿Qué es?** API que maneja jerarquía académica (escuelas, unidades, docentes)

**Versión:** Latest

**Comunicación:**
- HTTP REST (direct calls)
- Shared PostgreSQL (algunas tablas)

**Endpoints que API Mobile podría llamar:**
```go
// Obtener información de escuela
GET /api/v1/schools/{school_id}

// Obtener información de docente
GET /api/v1/teachers/{teacher_id}

// Obtener información de estudiante
GET /api/v1/students/{student_id}

// Obtener unidades académicas (jerarquía)
GET /api/v1/schools/{school_id}/academic-units
```

**Ejemplo de integración:**
```go
package adapters

import (
    "github.com/go-resty/resty/v2"
    "os"
)

var adminAPIClient = resty.New().
    SetBaseURL(os.Getenv("API_ADMIN_URL")). // http://api-admin:8081
    SetHeader("Authorization", "Bearer "+getToken())

func GetTeacher(teacherID int) (*models.Teacher, error) {
    resp, err := adminAPIClient.
        R().
        SetResult(&models.Teacher{}).
        Get("/api/v1/teachers/" + string(teacherID))
    
    if err != nil {
        return nil, err
    }
    
    return resp.Result().(*models.Teacher), nil
}
```

**Variables de entorno:**
```bash
API_ADMIN_URL=http://api-admin:8081
API_ADMIN_TIMEOUT=10s
API_ADMIN_RETRY_TIMES=3
```

---

## Dependencias de Desarrollo

### 7. Docker

**¿Qué es?** Containerización de aplicación

**Requisitos:**
- Docker versión 20.10+
- Docker Compose versión 2.0+

**Dockerfile:**
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api-mobile ./cmd/api-mobile

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/api-mobile .
EXPOSE 8080

CMD ["./api-mobile"]
```

**Docker Compose (en dev-environment):**
```yaml
services:
  api-mobile:
    build:
      context: .
      dockerfile: docker/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=edugo_user
      - DB_PASSWORD=edugo_pass
      - DB_NAME=edugo_mobile
      - MONGO_URI=mongodb://mongo:27017
      - RABBITMQ_HOST=rabbitmq
      - API_ADMIN_URL=http://api-admin:8081
    depends_on:
      - postgres
      - mongo
      - rabbitmq
      - api-admin
    networks:
      - edugo-network
```

---

### 8. Go 1.21+

**¿Qué es?** Lenguaje de programación

**Requisitos:**
- Go versión 1.21 o superior
- GOPATH correctamente configurado

**Setup local:**
```bash
# Descargar dependencias
go mod download

# Verificar dependencias
go mod tidy

# Compilar
go build -o api-mobile ./cmd/api-mobile

# Ejecutar
./api-mobile

# Tests
go test ./...

# Cobertura
go test -cover ./...
```

---

## Matriz de Compatibilidad

| Componente | Versión Requerida | Versión Máxima | Notas |
|-----------|------------------|----------------|-------|
| Go | 1.21+ | Latest | 1.22+ soportado |
| PostgreSQL | 15+ | 16+ | Compatible |
| MongoDB | 7.0+ | Latest | 8.0+ soportado |
| RabbitMQ | 3.12+ | Latest | 3.13+ soportado |
| SHARED | v1.3.0+ | v1.x | Breaking changes en v2.0 |
| Gin | v1.9+ | Latest | Compatible |
| GORM | v1.25+ | Latest | Compatible |

---

## Checklist de Instalación de Dependencias

```markdown
## Dependencias Críticas
- [ ] PostgreSQL 15 levantado en localhost:5432
- [ ] Base de datos "edugo_mobile" creada
- [ ] Tablas de evaluaciones creadas
- [ ] MongoDB 7.0 levantado en localhost:27017
- [ ] Colecciones de resultados creadas
- [ ] RabbitMQ 3.12 levantado en localhost:5672
- [ ] Exchanges y queues declaradas

## Dependencias de Código
- [ ] Go 1.21+ instalado
- [ ] go mod download ejecutado
- [ ] go mod tidy ejecutado
- [ ] SHARED v1.3.0+ en go.mod
- [ ] go build compila sin errores

## Dependencias de Integración
- [ ] WORKER disponible en RabbitMQ
- [ ] API ADMIN disponible en http://localhost:8081
- [ ] Credenciales JWT configuradas

## Dependencias de Desarrollo
- [ ] Docker instalado
- [ ] Docker Compose instalado
- [ ] Dockerfile presente
- [ ] docker-compose.yml presente
```

---

## Resolución de Problemas de Dependencias

### PostgreSQL no conecta
```bash
# Verificar conexión
psql -h localhost -U edugo_user -d edugo_mobile

# Ver logs
docker logs postgres

# Recrear si es necesario
docker compose down postgres
docker compose up postgres
```

### MongoDB no conecta
```bash
# Verificar conexión
mongosh --host localhost:27017

# Ver logs
docker logs mongo

# Verificar colecciones
use edugo_assessments
show collections
```

### RabbitMQ no conecta
```bash
# Acceder a management UI
http://localhost:15672
# user: guest, password: guest

# Ver queues
rabbitmqctl list_queues

# Listar exchanges
rabbitmqctl list_exchanges
```

### SHARED incompatible
```bash
# Verificar versión importada
cat go.mod | grep shared

# Actualizar SHARED
go get -u github.com/EduGoGroup/edugo-shared@v1.3.0

# Verificar compatibilidad
go mod tidy
go build ./cmd/api-mobile
```
