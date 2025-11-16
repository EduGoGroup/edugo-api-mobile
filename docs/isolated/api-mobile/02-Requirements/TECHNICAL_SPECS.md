# Especificaciones Técnicas
# Sistema de Evaluaciones - EduGo

**Versión:** 1.0.0  
**Fecha:** 14 de Noviembre, 2025  
**Proyecto:** edugo-api-mobile - Sistema de Evaluaciones

---

## 1. STACK TECNOLÓGICO

### 1.1 Backend (API Mobile)

| Componente | Tecnología | Versión | Justificación |
|------------|-----------|---------|---------------|
| **Lenguaje** | Go | 1.21+ | Performance, concurrencia nativa, consistencia con proyecto |
| **Framework Web** | Gin | 1.9+ | Más rápido, middleware robusto, usado en api-mobile |
| **Documentación API** | Swaggo | 1.16+ | Genera OpenAPI 3.0 desde anotaciones Go |
| **ORM** | GORM | 1.25+ | Usado actualmente en api-mobile |
| **Validación** | go-playground/validator | 10.15+ | Validación struct-based |
| **Testing** | Testify | 1.8+ | Assertions y mocks |
| **Testcontainers** | shared/testing | v0.6.2 | Reutilizar módulo compartido |

### 1.2 Bases de Datos

| Base de Datos | Versión | Uso | Driver Go |
|---------------|---------|-----|-----------|
| **PostgreSQL** | 15+ | Datos relacionales (intentos, respuestas) | github.com/lib/pq |
| **MongoDB** | 7.0+ | Documentos (preguntas, feedback) | go.mongodb.org/mongo-driver |

### 1.3 Infraestructura

| Servicio | Versión | Propósito |
|----------|---------|-----------|
| **Docker** | 24+ | Contenedores |
| **Docker Compose** | 2.20+ | Orquestación local |
| **RabbitMQ** | 3.12+ | Mensajería asíncrona (Post-MVP) |

### 1.4 CI/CD

| Herramienta | Propósito |
|-------------|-----------|
| **GitHub Actions** | Pipeline CI/CD |
| **golangci-lint** | Linting estático |
| **go test -cover** | Coverage de tests |
| **Swagger UI** | Documentación interactiva |

---

## 2. ARQUITECTURA DE SOFTWARE

### 2.1 Clean Architecture (Hexagonal)

```
api-mobile/
├── cmd/
│   └── api/
│       └── main.go                    # Entry point
│
├── internal/
│   ├── domain/                        # Capa de Dominio (Entities, Value Objects)
│   │   ├── entities/
│   │   │   ├── assessment.go          # Entity: Assessment
│   │   │   ├── attempt.go             # Entity: Attempt
│   │   │   └── answer.go              # Entity: Answer
│   │   │
│   │   ├── value_objects/
│   │   │   ├── assessment_id.go       # VO: AssessmentID
│   │   │   ├── attempt_id.go          # VO: AttemptID
│   │   │   ├── question_id.go         # VO: QuestionID
│   │   │   └── score.go               # VO: Score (0-100)
│   │   │
│   │   └── repositories/              # Interfaces de Repositorios
│   │       ├── assessment_repository.go
│   │       ├── attempt_repository.go
│   │       └── question_repository.go # MongoDB
│   │
│   ├── application/                   # Capa de Aplicación (Use Cases, Services)
│   │   ├── services/
│   │   │   ├── assessment_service.go  # Business logic
│   │   │   └── scoring_service.go     # Cálculo de puntajes
│   │   │
│   │   └── dto/
│   │       ├── assessment_dto.go      # DTOs de request/response
│   │       └── attempt_dto.go
│   │
│   └── infrastructure/                # Capa de Infraestructura (Implementaciones)
│       ├── persistence/
│       │   ├── postgres/
│       │   │   ├── assessment_repo.go # Implementación PostgreSQL
│       │   │   └── attempt_repo.go
│       │   │
│       │   └── mongodb/
│       │       └── question_repo.go   # Implementación MongoDB
│       │
│       ├── http/
│       │   ├── handlers/
│       │   │   └── assessment_handler.go # HTTP handlers
│       │   │
│       │   ├── middleware/
│       │   │   └── auth_middleware.go
│       │   │
│       │   └── routes/
│       │       └── assessment_routes.go
│       │
│       └── messaging/                 # RabbitMQ (Post-MVP)
│           └── event_publisher.go
│
├── scripts/
│   └── postgresql/
│       └── 06_assessments.sql         # Schema BD
│
└── tests/
    ├── unit/
    ├── integration/
    └── e2e/
```

### 2.2 Principios de Diseño

#### Inversión de Dependencias
```go
// ✅ CORRECTO: Domain no depende de Infrastructure
// domain/repositories/attempt_repository.go
package repositories

type AttemptRepository interface {
    Create(ctx context.Context, attempt *entities.Attempt) error
    FindByID(ctx context.Context, id uuid.UUID) (*entities.Attempt, error)
    FindByStudentID(ctx context.Context, studentID uuid.UUID) ([]*entities.Attempt, error)
}

// infrastructure/persistence/postgres/attempt_repo.go
package postgres

type PostgresAttemptRepository struct {
    db *gorm.DB
}

func (r *PostgresAttemptRepository) Create(ctx context.Context, attempt *entities.Attempt) error {
    // Implementación específica de PostgreSQL
}
```

#### Separación de Concerns
- **Domain:** Lógica de negocio pura, sin dependencias externas
- **Application:** Orquestación de use cases, DTOs
- **Infrastructure:** Implementaciones concretas (DB, HTTP, etc.)

---

## 3. MODELO DE DATOS

### 3.1 PostgreSQL - Schema Completo

#### Tabla: `assessment`
```sql
CREATE TABLE assessment (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    material_id UUID NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    mongo_document_id VARCHAR(24) NOT NULL, -- ObjectId de MongoDB
    title VARCHAR(255) NOT NULL,
    total_questions INTEGER NOT NULL CHECK (total_questions > 0),
    pass_threshold INTEGER NOT NULL DEFAULT 70 CHECK (pass_threshold >= 0 AND pass_threshold <= 100),
    max_attempts INTEGER DEFAULT NULL, -- NULL = ilimitado (Post-MVP)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT unique_material_assessment UNIQUE (material_id)
);

CREATE INDEX idx_assessment_material_id ON assessment(material_id);
CREATE INDEX idx_assessment_mongo_document_id ON assessment(mongo_document_id);
```

#### Tabla: `assessment_attempt`
```sql
CREATE TABLE assessment_attempt (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    assessment_id UUID NOT NULL REFERENCES assessment(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    score INTEGER NOT NULL CHECK (score >= 0 AND score <= 100),
    max_score INTEGER NOT NULL DEFAULT 100,
    time_spent_seconds INTEGER NOT NULL CHECK (time_spent_seconds > 0),
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CHECK (completed_at >= started_at)
);

CREATE INDEX idx_attempt_assessment_id ON assessment_attempt(assessment_id);
CREATE INDEX idx_attempt_student_id ON assessment_attempt(student_id);
CREATE INDEX idx_attempt_completed_at ON assessment_attempt(completed_at DESC);
CREATE INDEX idx_attempt_student_assessment ON assessment_attempt(student_id, assessment_id);
```

#### Tabla: `assessment_attempt_answer`
```sql
CREATE TABLE assessment_attempt_answer (
    attempt_id UUID NOT NULL REFERENCES assessment_attempt(id) ON DELETE CASCADE,
    question_id VARCHAR(100) NOT NULL, -- ID de pregunta en MongoDB
    selected_option VARCHAR(10) NOT NULL, -- 'a', 'b', 'c', 'd'
    is_correct BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (attempt_id, question_id)
);

CREATE INDEX idx_answer_attempt_id ON assessment_attempt_answer(attempt_id);
CREATE INDEX idx_answer_question_id ON assessment_attempt_answer(question_id);
```

#### Tabla: `material_summary_link` (Opcional - Mejor práctica)
```sql
CREATE TABLE material_summary_link (
    material_id UUID PRIMARY KEY REFERENCES materials(id) ON DELETE CASCADE,
    mongo_document_id VARCHAR(24) NOT NULL, -- ObjectId de MongoDB
    summary_version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_summary_link_mongo_id ON material_summary_link(mongo_document_id);
```

### 3.2 MongoDB - Colecciones

#### Colección: `material_assessment`
```javascript
{
  "_id": ObjectId("..."),
  "material_id": "uuid-material",
  "title": "Cuestionario: Introducción a Pascal",
  "questions": [
    {
      "id": "q1",
      "text": "¿Qué es un compilador?",
      "type": "multiple_choice",
      "options": [
        {"id": "a", "text": "Un programa que traduce código fuente a código máquina"},
        {"id": "b", "text": "Un tipo de variable en Pascal"},
        {"id": "c", "text": "Una estructura de control"},
        {"id": "d", "text": "Un editor de texto"}
      ],
      "correct_answer": "a",
      "feedback": {
        "correct": "¡Correcto! Un compilador traduce código fuente a código máquina ejecutable.",
        "incorrect": "Incorrecto. Revisa la sección 'Herramientas de Desarrollo' en el resumen."
      }
    }
  ],
  "version": 1,
  "created_at": ISODate("2025-01-15T12:30:00Z"),
  "updated_at": ISODate("2025-01-15T12:30:00Z")
}
```

**Validación de Schema:**
```javascript
db.createCollection("material_assessment", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["material_id", "title", "questions", "version"],
      properties: {
        material_id: { bsonType: "string" },
        title: { bsonType: "string" },
        questions: {
          bsonType: "array",
          minItems: 1,
          items: {
            bsonType: "object",
            required: ["id", "text", "type", "options", "correct_answer"],
            properties: {
              id: { bsonType: "string" },
              text: { bsonType: "string" },
              type: { enum: ["multiple_choice", "true_false", "short_answer"] },
              options: { bsonType: "array" },
              correct_answer: { bsonType: "string" },
              feedback: { bsonType: "object" }
            }
          }
        },
        version: { bsonType: "int" }
      }
    }
  }
});
```

**Índices:**
```javascript
db.material_assessment.createIndex({ "material_id": 1 }, { unique: true });
db.material_assessment.createIndex({ "created_at": -1 });
```

---

## 4. API REST - CONTRATOS

### 4.1 Convenciones

- **Versionado:** `/v1/` en todas las rutas
- **Formato:** JSON (Content-Type: application/json)
- **Autenticación:** JWT en header `Authorization: Bearer {token}`
- **Códigos HTTP:**
  - 200: Success
  - 201: Created
  - 400: Bad Request (validación)
  - 401: Unauthorized (sin token o inválido)
  - 403: Forbidden (sin permisos)
  - 404: Not Found
  - 500: Internal Server Error

### 4.2 Endpoints Detallados

#### GET /v1/materials/:id/assessment
```go
// Handler signature
func (h *AssessmentHandler) GetAssessment(c *gin.Context) {
    materialID := c.Param("id")
    userID := c.GetString("user_id") // De JWT middleware
    
    assessment, err := h.service.GetAssessmentForStudent(c.Request.Context(), materialID, userID)
    if err != nil {
        // Handle error
    }
    
    c.JSON(200, assessment)
}

// Response DTO
type AssessmentDTO struct {
    AssessmentID         string        `json:"assessment_id"`
    MaterialID           string        `json:"material_id"`
    Title                string        `json:"title"`
    TotalQuestions       int           `json:"total_questions"`
    EstimatedTimeMinutes int           `json:"estimated_time_minutes"`
    Questions            []QuestionDTO `json:"questions"`
}

type QuestionDTO struct {
    ID      string      `json:"id"`
    Text    string      `json:"text"`
    Type    string      `json:"type"`
    Options []OptionDTO `json:"options"`
    // ⚠️ NUNCA incluir: CorrectAnswer, Feedback
}
```

**Swagger Annotation:**
```go
// GetAssessment godoc
// @Summary      Obtener cuestionario de un material
// @Description  Retorna las preguntas de evaluación sin respuestas correctas
// @Tags         assessments
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Material ID (UUID)"
// @Success      200  {object}  AssessmentDTO
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /v1/materials/{id}/assessment [get]
```

#### POST /v1/materials/:id/assessment/attempts
```go
// Request DTO
type CreateAttemptRequest struct {
    Answers          []AnswerDTO `json:"answers" binding:"required,min=1"`
    TimeSpentSeconds int         `json:"time_spent_seconds" binding:"required,min=1"`
}

type AnswerDTO struct {
    QuestionID     string `json:"question_id" binding:"required"`
    SelectedOption string `json:"selected_option" binding:"required"`
}

// Response DTO
type AttemptResultDTO struct {
    AttemptID        string           `json:"attempt_id"`
    Score            int              `json:"score"`
    MaxScore         int              `json:"max_score"`
    CorrectAnswers   int              `json:"correct_answers"`
    TotalQuestions   int              `json:"total_questions"`
    PassThreshold    int              `json:"pass_threshold"`
    Passed           bool             `json:"passed"`
    Feedback         []FeedbackDTO    `json:"feedback"`
    CanRetake        bool             `json:"can_retake"`
    PreviousBestScore *int            `json:"previous_best_score"`
}

type FeedbackDTO struct {
    QuestionID     string `json:"question_id"`
    QuestionText   string `json:"question_text"`
    SelectedOption string `json:"selected_option"`
    CorrectAnswer  string `json:"correct_answer"`
    IsCorrect      bool   `json:"is_correct"`
    Message        string `json:"message"`
}
```

**Validaciones:**
```go
// Validar que todas las preguntas tienen respuesta
if len(req.Answers) != assessment.TotalQuestions {
    return errors.New("incomplete answers")
}

// Validar que no hay respuestas duplicadas
seen := make(map[string]bool)
for _, answer := range req.Answers {
    if seen[answer.QuestionID] {
        return errors.New("duplicate question_id")
    }
    seen[answer.QuestionID] = true
}
```

---

## 5. REQUISITOS DE PERFORMANCE

### 5.1 Latencia

| Endpoint | p50 | p95 | p99 | Max Aceptable |
|----------|-----|-----|-----|---------------|
| GET /assessment | <100ms | <200ms | <300ms | 500ms |
| POST /attempts | <500ms | <1500ms | <2000ms | 3000ms |
| GET /attempts/:id | <100ms | <200ms | <300ms | 500ms |
| GET /users/me/attempts | <200ms | <400ms | <600ms | 1000ms |

### 5.2 Throughput

| Endpoint | Requests/seg | Usuarios Concurrentes |
|----------|--------------|----------------------|
| GET /assessment | 100 req/s | 500 usuarios |
| POST /attempts | 50 req/s | 250 usuarios |

### 5.3 Escalabilidad

- **Base de Datos:**
  - PostgreSQL: Hasta 100,000 intentos/día
  - MongoDB: Hasta 10,000 assessments

- **API:**
  - Stateless: Escalado horizontal con load balancer
  - Pool de conexiones: Max 100 conexiones PostgreSQL

---

## 6. REQUISITOS DE SEGURIDAD

### 6.1 Autenticación y Autorización

```go
// Middleware de autenticación JWT
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "missing token"})
            c.Abort()
            return
        }
        
        claims, err := ValidateJWT(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}
```

### 6.2 Validación de Inputs

```go
// SIEMPRE validar con struct tags
type CreateAttemptRequest struct {
    Answers          []AnswerDTO `json:"answers" binding:"required,min=1,max=50,dive"`
    TimeSpentSeconds int         `json:"time_spent_seconds" binding:"required,min=1,max=7200"`
}

// SIEMPRE sanitizar IDs
func ValidateUUID(id string) error {
    _, err := uuid.Parse(id)
    return err
}
```

### 6.3 Protección Contra Trampas

```go
// NUNCA exponer respuestas correctas antes de submit
func (s *AssessmentService) GetAssessment(ctx context.Context, materialID string) (*AssessmentDTO, error) {
    questions, err := s.questionRepo.FindByMaterialID(ctx, materialID)
    if err != nil {
        return nil, err
    }
    
    // ⚠️ CRÍTICO: Sanitizar
    for i := range questions {
        questions[i].CorrectAnswer = ""      // ❌ Remover
        questions[i].Feedback = nil          // ❌ Remover
    }
    
    return toDTO(questions), nil
}

// SIEMPRE validar en servidor
func (s *ScoringService) CalculateScore(answers []Answer, correctAnswers []Question) int {
    // NUNCA confiar en score enviado por cliente
    correctCount := 0
    for _, answer := range answers {
        question := findQuestion(correctAnswers, answer.QuestionID)
        if answer.SelectedOption == question.CorrectAnswer {
            correctCount++
        }
    }
    return (correctCount * 100) / len(answers)
}
```

### 6.4 Rate Limiting (Post-MVP)

```go
// Límite por rol
var rateLimits = map[string]int{
    "student": 100, // 100 req/min
    "teacher": 200,
    "admin":   500,
}
```

---

## 7. REQUISITOS DE OBSERVABILIDAD

### 7.1 Logging Estructurado

```go
import "github.com/edugogroup/edugo-shared/pkg/logger"

func (s *AssessmentService) CreateAttempt(ctx context.Context, req CreateAttemptRequest) error {
    logger.Info("Creating assessment attempt",
        "student_id", req.StudentID,
        "assessment_id", req.AssessmentID,
        "total_answers", len(req.Answers))
    
    // Business logic
    
    logger.Info("Assessment attempt created",
        "attempt_id", attemptID,
        "score", score,
        "duration_ms", elapsed.Milliseconds())
    
    return nil
}
```

### 7.2 Métricas

```go
// Prometheus metrics
var (
    attemptDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "assessment_attempt_duration_seconds",
            Help: "Time to process assessment attempt",
        },
        []string{"status"},
    )
    
    attemptScore = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name: "assessment_attempt_score",
            Help: "Distribution of assessment scores",
        },
    )
)
```

### 7.3 Health Checks

```go
// GET /health
func HealthCheck(c *gin.Context) {
    dbStatus := checkPostgres()
    mongoStatus := checkMongoDB()
    
    healthy := dbStatus && mongoStatus
    statusCode := 200
    if !healthy {
        statusCode = 503
    }
    
    c.JSON(statusCode, gin.H{
        "status": map[bool]string{true: "healthy", false: "unhealthy"}[healthy],
        "postgres": dbStatus,
        "mongodb": mongoStatus,
    })
}
```

---

## 8. REQUISITOS DE TESTING

### 8.1 Coverage Mínimo

| Capa | Coverage Objetivo | Herramienta |
|------|-------------------|-------------|
| Domain | >90% | go test -cover |
| Application | >85% | go test -cover |
| Infrastructure | >70% | go test + testcontainers |
| **Global** | **>80%** | CI/CD enforced |

### 8.2 Tipos de Tests

#### Tests Unitarios
```go
// domain/entities/attempt_test.go
func TestAttempt_CalculateScore(t *testing.T) {
    attempt := &Attempt{
        Answers: []Answer{
            {QuestionID: "q1", SelectedOption: "a", IsCorrect: true},
            {QuestionID: "q2", SelectedOption: "b", IsCorrect: false},
        },
    }
    
    score := attempt.CalculateScore()
    assert.Equal(t, 50, score)
}
```

#### Tests de Integración
```go
// tests/integration/assessment_repo_test.go
func TestAssessmentRepository_Create(t *testing.T) {
    // Setup testcontainer PostgreSQL
    ctx := context.Background()
    container, db := setupTestDB(t)
    defer container.Terminate(ctx)
    
    repo := postgres.NewAssessmentRepository(db)
    
    // Test
    assessment := &entities.Assessment{
        MaterialID: uuid.New(),
        Title: "Test Assessment",
    }
    
    err := repo.Create(ctx, assessment)
    assert.NoError(t, err)
    assert.NotNil(t, assessment.ID)
}
```

#### Tests End-to-End
```go
// tests/e2e/assessment_flow_test.go
func TestFullAssessmentFlow(t *testing.T) {
    // 1. Setup test server
    router := setupTestRouter()
    
    // 2. Obtener assessment
    resp := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/v1/materials/uuid-1/assessment", nil)
    router.ServeHTTP(resp, req)
    assert.Equal(t, 200, resp.Code)
    
    // 3. Crear intento
    body := `{"answers": [...], "time_spent_seconds": 120}`
    req, _ = http.NewRequest("POST", "/v1/materials/uuid-1/assessment/attempts", strings.NewReader(body))
    router.ServeHTTP(resp, req)
    assert.Equal(t, 201, resp.Code)
    
    // 4. Consultar resultados
    // ...
}
```

---

## 9. DEPENDENCIAS EXTERNAS

### 9.1 Bibliotecas Go

```go
// go.mod
module github.com/edugogroup/edugo-api-mobile

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/google/uuid v1.4.0
    github.com/lib/pq v1.10.9
    go.mongodb.org/mongo-driver v1.13.0
    gorm.io/gorm v1.25.5
    gorm.io/driver/postgres v1.5.4
    
    // Testing
    github.com/stretchr/testify v1.8.4
    github.com/testcontainers/testcontainers-go v0.26.0
    
    // Shared
    github.com/edugogroup/edugo-shared v0.6.2
)
```

### 9.2 Servicios Externos

| Servicio | Dependencia | Criticidad | Fallback |
|----------|-------------|------------|----------|
| PostgreSQL | Datos relacionales | 🔴 Crítica | Retry + Circuit breaker |
| MongoDB | Preguntas | 🔴 Crítica | Retry + Circuit breaker |
| RabbitMQ | Notificaciones | 🟡 Media | Logs (sin notificación) |

---

## 10. COMPATIBILIDAD

### 10.1 Versiones de Go

- **Mínimo:** Go 1.21
- **Recomendado:** Go 1.21+ (latest patch)
- **Incompatible:** Go <1.21

### 10.2 Bases de Datos

- **PostgreSQL:** 15.0 - 16.x
- **MongoDB:** 7.0 - 7.x

### 10.3 Clients

- **API Mobile:** Cualquier cliente HTTP que soporte JWT
- **Swagger UI:** Navegadores modernos (Chrome, Firefox, Safari)

---

## 11. DECISIONES ARQUITECTÓNICAS (ADRs)

### ADR-001: Usar Clean Architecture

**Decisión:** Implementar Clean Architecture (Hexagonal) con separación en capas Domain, Application, Infrastructure.

**Razones:**
- Consistencia con api-mobile existente
- Testabilidad (mock de interfaces)
- Independencia de frameworks y BDs

**Alternativas Consideradas:**
- MVC tradicional (rechazado: acoplamiento)
- Arquitectura en capas simple (rechazado: menos testeable)

---

### ADR-002: PostgreSQL para Intentos, MongoDB para Preguntas

**Decisión:** Almacenar intentos en PostgreSQL y preguntas en MongoDB.

**Razones:**
- PostgreSQL: Transacciones ACID para intentos (crítico)
- MongoDB: Flexibilidad de schema para preguntas (worker genera estructura variable)
- Arquitectura híbrida ya existente

**Alternativas Consideradas:**
- Todo en PostgreSQL (rechazado: preguntas en JSONB menos flexible)
- Todo en MongoDB (rechazado: transacciones ACID complejas)

---

### ADR-003: Validación Siempre en Servidor

**Decisión:** Nunca confiar en validaciones del cliente, siempre re-validar en servidor.

**Razones:**
- Seguridad: Cliente puede ser modificado
- Integridad: Respuestas correctas solo en servidor
- Auditoría: Registro completo de validaciones

**Alternativas Consideradas:**
- Validación solo en cliente (rechazado: inseguro)
- Validación híbrida (rechazado: redundante)

---

## 12. LIMITACIONES CONOCIDAS

| Limitación | Impacto | Plan de Mitigación |
|------------|---------|-------------------|
| Worker debe existir y funcionar | 🔴 Bloqueante | Verificar antes de Sprint 1, Plan B: generar manualmente |
| Esquema MongoDB puede variar | 🟡 Alto | Adapter pattern, validación de schema |
| Sin sistema de jerarquía académica (MVP) | 🟢 Bajo | Cualquier usuario puede acceder (Post-MVP: filtrar por unidad) |

---

**Generado con:** Claude Code  
**Última actualización:** 2025-11-14  
**Próxima revisión:** Tras completar Sprint 1
