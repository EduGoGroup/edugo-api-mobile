# TECH STACK - API Mobile

## Resumen Ejecutivo

| Layer | Tecnología | Versión | Propósito |
|-------|-----------|---------|----------|
| **Language** | Go | 1.21+ | Backend, compilado a binario |
| **Framework** | Gin | v1.9+ | Web framework HTTP |
| **ORM** | GORM | v1.25+ | Abstracción PostgreSQL |
| **Primary DB** | PostgreSQL | 15+ | Datos relacionales (evaluaciones) |
| **Secondary DB** | MongoDB | 7.0+ | Documentos (resultados) |
| **Message Broker** | RabbitMQ | 3.12+ | Comunicación async con Worker |
| **Config Management** | Viper | Latest | Multi-environment config |
| **Authentication** | JWT (shared) | Custom | Token-based auth |
| **Logging** | Structured (shared) | Custom | Centralized logging |
| **Containerization** | Docker | 20.10+ | Container runtime |
| **Orchestration** | Docker Compose | 2.0+ | Local development |

---

## Stack Detallado por Capa

### Capa 1: Aplicación (Go)

```
┌─────────────────────────────────────┐
│   API Mobile Application (Go)       │
├─────────────────────────────────────┤
│ • Binary ejecutable compilado       │
│ • Goroutines para concurrencia      │
│ • Channels para sincronización      │
│ • Context para timeouts/cancellation│
│ • Interfaces para inyección deps    │
└─────────────────────────────────────┘
```

**Características de Go:**
- Compilado: Eficiente en CPU y memoria
- Concurrencia nativa: Maneja miles de goroutines
- Garbage collection: Automático
- Tipado estático: Errores en tiempo de compilación
- Deploye simple: Single binary, sin dependencias runtime

**Estructura de proyecto Go:**
```
api-mobile/
├── cmd/
│   └── api-mobile/
│       └── main.go         # Entry point
├── internal/
│   ├── handlers/           # HTTP handlers (Gin)
│   ├── services/           # Business logic
│   ├── repositories/       # Data access (GORM)
│   ├── models/             # Domain models
│   ├── middleware/         # Gin middleware
│   ├── config/             # Configuration
│   └── errors/             # Custom errors
├── migrations/             # DB migrations (GORM)
├── docker/
│   └── Dockerfile
├── go.mod                  # Module definition
├── go.sum                  # Dependency hashes
└── Makefile               # Build commands
```

---

### Capa 2: Web Framework (Gin)

**¿Qué es Gin?** Framework HTTP ligero y rápido para Go

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    router := gin.Default()
    
    // Middleware
    router.Use(AuthMiddleware())
    router.Use(LoggingMiddleware())
    
    // Rutas
    v1 := router.Group("/api/v1")
    {
        v1.POST("/evaluations", CreateEvaluation)
        v1.GET("/evaluations/:id", GetEvaluation)
        v1.PUT("/evaluations/:id", UpdateEvaluation)
        v1.DELETE("/evaluations/:id", DeleteEvaluation)
        
        v1.POST("/evaluations/:id/preguntas", CreateQuestion)
        v1.GET("/evaluations/:id/preguntas", ListQuestions)
    }
    
    // Correr servidor
    router.Run(":8080")
}

func CreateEvaluation(c *gin.Context) {
    var req CreateEvaluationRequest
    
    // Bind JSON request
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Procesar
    evaluation := &Evaluation{
        Title: req.Title,
        Type:  req.Type,
    }
    
    // Guardar
    if err := db.Create(evaluation).Error; err != nil {
        c.JSON(500, gin.H{"error": "Internal error"})
        return
    }
    
    c.JSON(201, evaluation)
}
```

**Ventajas de Gin:**
- Rapidísimo (benchmarks: 40x más rápido que algunos competitors)
- Middleware system similar a Express.js
- Validación automática de JSON
- Routing eficiente
- Error handling elegante
- Supporta gzip, logging, custom recovery

**Estructura de handlers:**
```go
package handlers

import "github.com/gin-gonic/gin"

type EvaluationHandler struct {
    service services.EvaluationService
}

func NewEvaluationHandler(service services.EvaluationService) *EvaluationHandler {
    return &EvaluationHandler{service: service}
}

func (h *EvaluationHandler) Create(c *gin.Context) {
    // Logic
}

func (h *EvaluationHandler) List(c *gin.Context) {
    // Logic
}

// Register en router
func RegisterRoutes(router *gin.Engine, handler *EvaluationHandler) {
    g := router.Group("/api/v1")
    {
        g.POST("/evaluations", handler.Create)
        g.GET("/evaluations", handler.List)
    }
}
```

---

### Capa 3: ORM (GORM)

**¿Qué es GORM?** Object-Relational Mapping para Go

```go
package models

import "gorm.io/gorm"

type Evaluation struct {
    ID          int64     `gorm:"primaryKey"`
    MaterialID  int64
    Title       string    `gorm:"column:title"`
    Description string
    Type        string    // 'manual', 'generated'
    Status      string    // 'draft', 'published', 'closed'
    PassingScore int
    CreatedBy   int64
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"` // Soft delete
    
    // Relations
    Questions []Question `gorm:"foreignKey:EvaluationID"`
}

// Definir table name
func (Evaluation) TableName() string {
    return "evaluations"
}
```

**Operaciones CRUD:**
```go
// CREATE
evaluation := &Evaluation{
    Title: "Quiz Matemáticas",
    Type: "manual",
}
db.Create(evaluation)

// READ
var eval Evaluation
db.First(&eval, id) // Por ID
db.Where("status = ?", "published").Find(&evals) // Con filtros

// UPDATE
db.Model(&evaluation).Update("status", "published")
db.Model(&evaluation).Updates(&Evaluation{
    Title: "Nuevo título",
    Status: "published",
})

// DELETE (soft delete)
db.Delete(&evaluation)
db.Unscoped().Delete(&evaluation) // Hard delete

// Transacciones
tx := db.BeginTx(ctx, nil)
tx.Create(&evaluation)
tx.Create(&question)
tx.Commit()
```

**Relaciones:**
```go
type Evaluation struct {
    ID BIGSERIAL PRIMARY KEY
    // ...
    Questions []Question `gorm:"foreignKey:EvaluationID;constraint:OnDelete:CASCADE"`
}

type Question struct {
    ID            BIGSERIAL PRIMARY KEY
    EvaluationID  int64
    Evaluation    Evaluation `gorm:"foreignKey:EvaluationID"`
    // ...
    Options []QuestionOption `gorm:"foreignKey:QuestionID"`
}

// Queries con relaciones
var eval Evaluation
db.Preload("Questions").Preload("Questions.Options").First(&eval, id)

// Lazy load
db.Model(&eval).Association("Questions").Find(&questions)
```

**Migraciones:**
```go
package migrations

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
    // AutoMigrate crea/actualiza schema basado en structs
    return db.AutoMigrate(
        &Evaluation{},
        &Question{},
        &QuestionOption{},
        &EvaluationAssignment{},
        &AnswerDraft{},
    )
}

// Migraciones personalizadas
func CreateIndexes(db *gorm.DB) error {
    return db.Migrator().CreateIndex(&Evaluation{}, "status")
}
```

---

### Capa 4: Bases de Datos

#### PostgreSQL 15+

**¿Qué es?** Sistema relacional para datos ACID, transacciones, integridad referencial

**Características usadas en API Mobile:**

1. **Tipos de datos:**
   ```sql
   -- INT para IDs
   id BIGSERIAL PRIMARY KEY
   
   -- VARCHAR/TEXT para strings
   title VARCHAR(255)
   description TEXT
   
   -- ENUM para estados
   type VARCHAR(50) CHECK (type IN ('manual', 'generated'))
   
   -- TIMESTAMP para auditoría
   created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
   ```

2. **Constraints:**
   ```sql
   -- Primary keys
   PRIMARY KEY (id)
   
   -- Foreign keys (relaciones)
   FOREIGN KEY (evaluation_id) REFERENCES evaluations(id)
   
   -- Unique constraints
   UNIQUE (email)
   
   -- NOT NULL
   title VARCHAR(255) NOT NULL
   ```

3. **Índices (Performance):**
   ```sql
   CREATE INDEX idx_evaluations_status ON evaluations(status);
   CREATE INDEX idx_evaluations_created_by ON evaluations(created_by);
   CREATE INDEX idx_questions_evaluation_id ON questions(evaluation_id);
   ```

4. **Transacciones (Data integrity):**
   ```sql
   BEGIN;
   INSERT INTO evaluations VALUES (...);
   INSERT INTO questions VALUES (...);
   COMMIT;
   
   -- O ROLLBACK si error
   ```

5. **Connection pooling:**
   ```go
   // GORM maneja automáticamente
   // Default: max 25 conexiones
   db.DB().SetMaxOpenConns(100)
   db.DB().SetMaxIdleConns(10)
   db.DB().SetConnMaxLifetime(time.Hour)
   ```

---

#### MongoDB 7.0+

**¿Qué es?** Base de datos documental NoSQL, flexible schema

```javascript
// Estructura de documento en MongoDB
{
  "_id": ObjectId("507f1f77bcf86cd799439011"),
  "evaluation_id": 1,
  "assignment_id": 1,
  "student_id": 42,
  "answers": [
    {
      "question_id": 101,
      "answer": "Option B",
      "is_correct": true,
      "points_earned": 5,
      "timestamp": ISODate("2025-11-15T10:30:00Z")
    }
  ],
  "total_score": 45,
  "max_score": 50,
  "percentage": 90.0,
  "status": "graded",
  "submitted_at": ISODate("2025-11-15T10:45:00Z"),
  "feedback": "Excelente desempeño",
  "metadata": {
    "ip_address": "192.168.1.1",
    "user_agent": "Mobile App v1.0",
    "device_id": "uuid-12345"
  }
}
```

**Operaciones en Go (usando mongo-go-driver):**
```go
import "go.mongodb.org/mongo-driver/mongo"

// Insertar
collection := client.Database("edugo").Collection("evaluation_results")
result, err := collection.InsertOne(ctx, result)

// Buscar
var result EvaluationResult
err = collection.FindOne(ctx, bson.M{
    "evaluation_id": 1,
    "student_id": 42,
}).Decode(&result)

// Actualizar
_, err = collection.UpdateOne(ctx, 
    bson.M{"_id": id},
    bson.M{"$set": bson.M{"status": "graded"}},
)

// Borrar
_, err = collection.DeleteOne(ctx, bson.M{"_id": id})

// Agregaciones
pipeline := mongo.Pipeline{
    bson.D{{Key: "$match", Value: bson.D{{Key: "evaluation_id", Value: 1}}}},
    bson.D{{Key: "$group", Value: bson.D{
        {Key: "_id", Value: "$student_id"},
        {Key: "avg_score", Value: bson.D{{Key: "$avg", Value: "$total_score"}}},
    }}},
}
cursor, err := collection.Aggregate(ctx, pipeline)
```

**Índices (Performance):**
```javascript
// Crear índices para queries comunes
db.evaluation_results.createIndex({ "evaluation_id": 1 })
db.evaluation_results.createIndex({ "student_id": 1 })
db.evaluation_results.createIndex({ "submitted_at": -1 })
db.evaluation_results.createIndex({ "evaluation_id": 1, "student_id": 1 })
```

**Ventaja sobre PostgreSQL:**
- Schema flexible (sin migración de tabla necesaria)
- Datos anidados naturales (respuestas con detalles)
- Queries analíticas potentes (aggregation framework)
- Mejor performance para lecturas de documentos complejos

---

### Capa 5: Message Broker (RabbitMQ)

**¿Qué es?** Sistema de colas para comunicación asíncrona confiable

```
┌────────────────┐
│  API Mobile    │
└────────┬───────┘
         │ Publica mensaje
         ▼
    ┌─────────────┐
    │ RabbitMQ    │
    │ Exchange    │
    └──────┬──────┘
           │
    ┌──────┴──────┐
    │             │
    ▼             ▼
┌────────────┐ ┌────────────┐
│ Queue 1    │ │ Queue 2    │
│ (Worker)   │ │ (Other)    │
└────────────┘ └────────────┘
    │
    ▼
┌────────────┐
│ Worker     │
│ consume    │
└────────────┘
```

**Publish desde API Mobile:**
```go
package messaging

import (
    "encoding/json"
    "github.com/streadway/amqp"
)

type Publisher struct {
    conn *amqp.Connection
    ch   *amqp.Channel
}

func (p *Publisher) PublishGenerateQuiz(req GenerateQuizRequest) error {
    body, _ := json.Marshal(req)
    
    return p.ch.Publish(
        "assessment.requests",        // exchange
        "worker.assessment.requests", // routing key
        false,                        // mandatory
        false,                        // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
}
```

**Subscribe desde API Mobile:**
```go
func (s *Subscriber) ConsumeAssessmentResponses(handler func([]byte) error) error {
    queue, err := s.ch.QueueDeclare(
        "api-mobile.assessment.responses", // queue name
        true,  // durable
        false, // exclusive
        false, // no-wait
        nil,   // arguments
    )
    if err != nil {
        return err
    }
    
    msgs, err := s.ch.Consume(
        queue.Name,
        "",    // consumer tag
        false, // auto-ack (false = manual ack)
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil,   // args
    )
    
    go func() {
        for d := range msgs {
            if err := handler(d.Body); err == nil {
                d.Ack(false) // Manual acknowledge
            } else {
                d.Nack(false, true) // Requeue
            }
        }
    }()
    
    return nil
}
```

**Dead Letter Queue (para errores):**
```go
// Declarar exchange DLX
s.ch.ExchangeDeclare(
    "assessment.dlx",
    "direct",
    true, false, false, false,
    nil,
)

// Declarar queue con DLX
s.ch.QueueDeclare(
    "api-mobile.assessment.responses",
    true, false, false, false,
    amqp.Table{
        "x-dead-letter-exchange": "assessment.dlx",
        "x-message-ttl": 3600000, // 1 hora
    },
)
```

---

### Capa 6: Configuración (Viper)

**¿Qué es Viper?** Gestor de configuración multi-source para Go

```go
package config

import "github.com/spf13/viper"

func InitConfig() error {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml") // o "json", "env"
    viper.AddConfigPath(".")
    viper.AddConfigPath("/etc/edugo/")
    
    // Leer variables de entorno
    viper.AutomaticEnv()
    viper.BindEnv("db.host", "DB_HOST")
    viper.BindEnv("db.port", "DB_PORT")
    
    // Valores por defecto
    viper.SetDefault("api.port", 8080)
    viper.SetDefault("api.timeout", 30*time.Second)
    
    return viper.ReadInConfig()
}

// Usar en código
func GetDBHost() string {
    return viper.GetString("db.host")
}

func GetAPIPort() int {
    return viper.GetInt("api.port")
}

// Escuchar cambios (hot reload)
viper.OnConfigChange(func(e fsnotify.Event) {
    logger.Info("Config file changed")
})
viper.WatchConfig()
```

**Archivo config.yaml:**
```yaml
api:
  port: 8080
  environment: development
  timeout: 30s

database:
  postgres:
    host: localhost
    port: 5432
    user: edugo_user
    password: ${DB_PASSWORD}  # From env var
    name: edugo_mobile
    
  mongo:
    uri: mongodb://localhost:27017
    database: edugo_assessments

rabbitmq:
  host: localhost
  port: 5672
  user: guest
  password: guest

shared:
  log_level: info
  context_timeout: 30s
```

---

### Capa 7: Autenticación (JWT vía SHARED)

**¿Qué es JWT?** Token basado en claims para autenticación sin sesión

```go
package auth

type Claims struct {
    UserID    int64
    SchoolID  int64
    Email     string
    Roles     []string
    ExpiresAt int64
}

// Validar token en middleware
func ValidateToken() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Missing token"})
            c.Abort()
            return
        }
        
        claims, err := ParseToken(token[7:]) // Remove "Bearer "
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // Inyectar en contexto
        c.Set("user_id", claims.UserID)
        c.Set("school_id", claims.SchoolID)
        c.Next()
    }
}

// Usar en handler
func CreateEvaluation(c *gin.Context) {
    userID := c.GetInt64("user_id")
    schoolID := c.GetInt64("school_id")
    
    eval := &Evaluation{
        CreatedBy: userID,
        SchoolID: schoolID,
    }
    // ...
}
```

---

### Capa 8: Logging (vía SHARED)

**¿Qué es?** Sistema centralizado de logs estructurados

```go
import "github.com/EduGoGroup/edugo-shared/logger"

// Inicializar
logger.Init(logger.Config{
    Level:  "info",
    Format: "json",
})

// Uso
logger.Info("Evaluación creada", map[string]interface{}{
    "evaluation_id": eval.ID,
    "created_by": eval.CreatedBy,
    "timestamp": time.Now(),
})

logger.Error("Error guardando evaluación", map[string]interface{}{
    "error": err.Error(),
    "evaluation_id": eval.ID,
    "stack": debug.Stack(),
})

logger.Warn("Timeout en consulta", map[string]interface{}{
    "query_time": 5.5,
    "table": "evaluations",
})

logger.Debug("Query ejecutada", map[string]interface{}{
    "sql": "SELECT * FROM evaluations WHERE id = $1",
    "params": []interface{}{1},
})
```

**Salida JSON:**
```json
{
  "timestamp": "2025-11-15T10:30:00.123Z",
  "level": "INFO",
  "message": "Evaluación creada",
  "evaluation_id": 1,
  "created_by": 42,
  "request_id": "uuid-12345",
  "user_id": 42,
  "service": "api-mobile"
}
```

---

### Capa 9: Containerización (Docker)

**¿Qué es Docker?** Empaqueta aplicación + dependencias en contenedor aislado

```dockerfile
# Multi-stage build: Compilar + Runtime pequeño
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Descargar dependencias (cacheable layer)
COPY go.mod go.sum ./
RUN go mod download

# Copiar código
COPY . .

# Compilar (CGO=0 para compatibilidad)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -o api-mobile \
    -ldflags="-X main.Version=$(git describe --tags)" \
    ./cmd/api-mobile

# Runtime: Alpine (muy pequeño)
FROM alpine:latest

# Instalar CA certs para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar solo binario compilado
COPY --from=builder /app/api-mobile .

# Configuración
EXPOSE 8080
ENV PORT=8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Ejecutar
CMD ["./api-mobile"]
```

**Tamaño resultante:** ~50MB (pequeño, rápido)

---

### Capa 10: Orchestración Local (Docker Compose)

**¿Qué es?** Define multi-contenedor stack para desarrollo local

```yaml
version: '3.8'

services:
  # API Mobile
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
      - LOG_LEVEL=debug
    depends_on:
      postgres:
        condition: service_healthy
      mongo:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy
    networks:
      - edugo-network
    volumes:
      - .:/app  # Live reload en desarrollo

  # PostgreSQL
  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=edugo_user
      - POSTGRES_PASSWORD=edugo_pass
      - POSTGRES_DB=edugo_mobile
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations/init.sql:/docker-entrypoint-initdb.d/01-init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U edugo_user"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - edugo-network

  # MongoDB
  mongo:
    image: mongo:7.0
    environment:
      - MONGO_INITDB_ROOT_USERNAME=admin
      - MONGO_INITDB_ROOT_PASSWORD=admin
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db
    healthcheck:
      test: echo 'db.runCommand("ping").ok' | mongosh localhost:27017/test --quiet
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - edugo-network

  # RabbitMQ
  rabbitmq:
    image: rabbitmq:3.12-management-alpine
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      - RABBITMQ_DEFAULT_USER=guest
      - RABBITMQ_DEFAULT_PASS=guest
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    healthcheck:
      test: rabbitmq-diagnostics ping
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - edugo-network

volumes:
  postgres_data:
  mongo_data:
  rabbitmq_data:

networks:
  edugo-network:
    driver: bridge
```

---

## Flujo de Datos Técnico

```
Cliente HTTP
    │
    ▼
┌─ Gin Router ─┐
│ (Handlers)   │
└──────┬───────┘
       │
       ▼
┌──────────────────────┐
│ Middleware Stack:    │
│ 1. Auth (JWT)        │
│ 2. Logging           │
│ 3. CORS              │
│ 4. Recovery (panic)  │
└──────┬───────────────┘
       │
       ▼
┌──────────────────────┐
│ Service Layer        │
│ (Business logic)     │
└──────┬───────────────┘
       │
       ├─────────────────┬──────────────┬──────────────┐
       │                 │              │              │
       ▼                 ▼              ▼              ▼
    GORM            RabbitMQ      MongoDB Driver   External APIs
  (PostgreSQL)      (Publish)     (InsertOne)     (API Admin)
       │                 │              │              │
       └─────────────────┴──────────────┴──────────────┘
                         │
                         ▼
                   HTTP Response (JSON)
```

---

## Comparativa: Gin vs Alternativas

| Aspecto | Gin | Echo | Chi |
|--------|-----|------|-----|
| Performance | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| Middleware | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| Validación | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| Comunidad | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| Documentación | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |

**Por qué Gin en EduGo:**
- Performance crítica para evaluaciones
- Validación JSON automática
- Middleware simples pero potentes
- Adoptado en WORKER también (consistencia)

---

## Benchmarks Esperados

**Hardware típico (MacBook Pro):**
- Requests/segundo: ~10,000 (sin BD)
- Latencia promedio: 5-50ms (con BD)
- Uso de memoria: ~50-100MB
- CPU por request: <1ms

**Con carga (100 concurrent users):**
- Throughput: ~5,000 req/sec
- P95 latency: ~100ms
- P99 latency: ~200ms
- Error rate: <0.1%

---

## Actualizaciones de Seguridad

| Componente | Política | Cadencia |
|-----------|----------|----------|
| Go | Latest minor version | Cada 1-2 meses |
| Gin | v1.9+ | Mensual |
| PostgreSQL driver | Latest | Mensual |
| MongoDB driver | Latest | Mensual |
| GORM | v1.25+ | Mensual |

---

## Recomendaciones

1. **Go:** Actualizar a 1.22+ cuando sea estable
2. **PostgreSQL:** Usar 15+ en producción, 13+ en desarrollo
3. **MongoDB:** 7.0+ recomendado, 6.0+ mínimo
4. **RabbitMQ:** 3.12+ para mejor performance
5. **SHARED:** Mantener en sincronía, evitar saltarse versiones
