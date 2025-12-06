# 🔄 Oportunidades de Refactorización

> **Última revisión:** Diciembre 2024  
> **Propósito:** Mejoras que no son bugs ni deuda técnica, pero mejorarían la calidad del código

---

## REF-001: Extraer Scoring Engine a Paquete Separado

### Ubicación Actual

El motor de puntuación está embebido en `assessment_attempt_service.go`:

```
internal/application/service/
├── assessment_attempt_service.go  (incluye lógica de scoring)
└── scoring/                       (paquete separado existe pero incompleto)
    ├── calculator.go
    └── ...
```

### Propuesta

Extraer toda la lógica de scoring a `internal/application/service/scoring/`:

```
internal/application/service/scoring/
├── engine.go           # Interface principal
├── calculator.go       # Cálculo de puntuación
├── feedback.go         # Generación de feedback
├── validator.go        # Validación de respuestas
└── types.go           # Tipos específicos de scoring
```

### Beneficios

1. **Testabilidad:** Scoring puede testearse independientemente
2. **Reutilización:** Otros servicios pueden usar el scoring
3. **Claridad:** Separación de responsabilidades
4. **Extensibilidad:** Fácil agregar nuevos tipos de preguntas

### Ejemplo de Interface

```go
// scoring/engine.go
package scoring

type Engine interface {
    // CalculateScore calcula el puntaje de un conjunto de respuestas
    CalculateScore(assessment *Assessment, answers []Answer) (*Result, error)

    // GenerateFeedback genera feedback detallado por pregunta
    GenerateFeedback(assessment *Assessment, answers []Answer) ([]Feedback, error)

    // ValidateAnswers valida que las respuestas sean válidas
    ValidateAnswers(assessment *Assessment, answers []Answer) error
}

type Result struct {
    Score           int
    MaxScore        int
    CorrectAnswers  int
    TotalQuestions  int
    Passed          bool
    Feedback        []Feedback
}

type Feedback struct {
    QuestionID    string
    QuestionText  string
    UserAnswer    string
    CorrectAnswer string
    IsCorrect     bool
    Explanation   string
}
```

---

## REF-002: Implementar Repository Pattern Completo para Progress

### Situación Actual

El `ProgressRepository` implementa Interface Segregation pero podría beneficiarse de un patrón más robusto.

```go
// Actual
type ProgressRepository interface {
    ProgressReader
    ProgressWriter
    ProgressStats
}
```

### Propuesta

Agregar métodos de batch y transaccionales:

```go
type ProgressRepository interface {
    ProgressReader
    ProgressWriter
    ProgressStats
    ProgressBatch      // Nuevo
    ProgressTransaction // Nuevo
}

type ProgressBatch interface {
    // BatchUpsert actualiza múltiples progresos en una transacción
    BatchUpsert(ctx context.Context, progresses []*Progress) error

    // BatchFindByUser obtiene todos los progresos de un usuario
    BatchFindByUser(ctx context.Context, userID UserID) ([]*Progress, error)
}

type ProgressTransaction interface {
    // WithTransaction ejecuta operaciones en una transacción
    WithTransaction(ctx context.Context, fn func(repo ProgressRepository) error) error
}
```

### Caso de Uso

Sincronizar progreso offline de múltiples materiales:

```go
func (s *progressService) SyncOfflineProgress(ctx context.Context, batch []ProgressUpdate) error {
    return s.repo.WithTransaction(ctx, func(repo repository.ProgressRepository) error {
        for _, update := range batch {
            if _, err := repo.Upsert(ctx, &update.Progress); err != nil {
                return err // Rollback automático
            }
        }
        return nil
    })
}
```

---

## REF-003: Mejorar Error Handling con Error Types

### Situación Actual

Los errores se manejan con funciones helper de `edugo-shared/common/errors`:

```go
return nil, errors.NewDatabaseError("create material", err)
return nil, errors.NewValidationError("invalid author_id format")
return nil, errors.NewNotFoundError("material", id)
```

### Propuesta

Crear tipos de error específicos del dominio:

```go
// internal/domain/errors/errors.go
package errors

import "fmt"

// MaterialError errores relacionados con materiales
type MaterialError struct {
    Op      string // Operación que falló
    ID      string // ID del material (si aplica)
    Err     error  // Error subyacente
    Code    string // Código para el cliente
    Message string // Mensaje para el usuario
}

func (e *MaterialError) Error() string {
    if e.ID != "" {
        return fmt.Sprintf("material %s: %s: %v", e.ID, e.Op, e.Err)
    }
    return fmt.Sprintf("material: %s: %v", e.Op, e.Err)
}

func (e *MaterialError) Unwrap() error {
    return e.Err
}

// Constructores específicos
func MaterialNotFound(id string) *MaterialError {
    return &MaterialError{
        Op:      "find",
        ID:      id,
        Code:    "MATERIAL_NOT_FOUND",
        Message: fmt.Sprintf("material %s not found", id),
    }
}

func MaterialCreationFailed(err error) *MaterialError {
    return &MaterialError{
        Op:      "create",
        Err:     err,
        Code:    "MATERIAL_CREATION_FAILED",
        Message: "failed to create material",
    }
}
```

### Beneficios

1. **Type-safe:** Permite `errors.As()` para manejo específico
2. **Contexto rico:** Incluye operación, ID, error original
3. **Consistencia:** Todos los errores de material tienen la misma estructura
4. **Testing:** Fácil verificar tipos de error específicos

---

## REF-004: Implementar Circuit Breaker para Servicios Externos

### Situación Actual

Las llamadas a servicios externos (S3, RabbitMQ) no tienen circuit breaker:

```go
// Actual - sin protección
if err := s.messagePublisher.Publish(ctx, exchange, routingKey, body); err != nil {
    s.logger.Warn("failed to publish event", "error", err)
    // Continúa sin reintentos ni circuit breaker
}
```

### Propuesta

Usar `sony/gobreaker` (ya en go.mod):

```go
// internal/infrastructure/messaging/resilient_publisher.go
package messaging

import (
    "context"
    "github.com/sony/gobreaker"
)

type ResilientPublisher struct {
    publisher Publisher
    cb        *gobreaker.CircuitBreaker
}

func NewResilientPublisher(pub Publisher) *ResilientPublisher {
    cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
        Name:        "rabbitmq-publisher",
        MaxRequests: 3,                // Máx requests en half-open
        Interval:    10 * time.Second, // Intervalo de conteo
        Timeout:     30 * time.Second, // Tiempo en open antes de half-open
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            // Abrir si >50% de requests fallan
            failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
            return counts.Requests >= 3 && failureRatio >= 0.5
        },
        OnStateChange: func(name string, from, to gobreaker.State) {
            log.Printf("Circuit breaker %s: %s -> %s", name, from, to)
        },
    })

    return &ResilientPublisher{publisher: pub, cb: cb}
}

func (r *ResilientPublisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
    _, err := r.cb.Execute(func() (interface{}, error) {
        return nil, r.publisher.Publish(ctx, exchange, routingKey, body)
    })
    return err
}
```

### Uso

```go
// En bootstrap
publisher := messaging.NewResilientPublisher(rawPublisher)
```

---

## REF-005: Agregar Request ID y Tracing

### Situación Actual

No hay request ID ni distributed tracing:

```go
// Logs actuales
s.logger.Info("material created", "material_id", id)
// No hay correlación entre logs de diferentes servicios
```

### Propuesta

Agregar middleware de request ID y propagarlo:

```go
// internal/infrastructure/http/middleware/request_id.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader(RequestIDHeader)
        if requestID == "" {
            requestID = uuid.New().String()
        }

        c.Set("request_id", requestID)
        c.Header(RequestIDHeader, requestID)

        c.Next()
    }
}

// Helper para obtener request ID
func GetRequestID(c *gin.Context) string {
    if id, exists := c.Get("request_id"); exists {
        return id.(string)
    }
    return ""
}
```

### Logger con Request ID

```go
// En handlers
requestID := middleware.GetRequestID(c)
h.logger.Info("processing request",
    "request_id", requestID,
    "endpoint", c.FullPath(),
    "method", c.Request.Method,
)
```

### Propagación a Servicios

```go
// En llamadas a RabbitMQ
headers := map[string]interface{}{
    "x-request-id": requestID,
}
publisher.PublishWithHeaders(ctx, exchange, routingKey, body, headers)
```

---

## REF-006: Implementar Healthcheck Detallado

### Situación Actual

El healthcheck es básico:

```go
func (h *HealthHandler) Check(c *gin.Context) {
    // Solo verifica conexiones básicas
}
```

### Propuesta

Implementar healthcheck con niveles de detalle:

```go
// GET /health           → básico (para load balancer)
// GET /health?detail=1  → detallado (para debugging)

type HealthResponse struct {
    Status    string                 `json:"status"`
    Timestamp time.Time              `json:"timestamp"`
    Version   string                 `json:"version,omitempty"`
    Checks    map[string]CheckResult `json:"checks,omitempty"`
}

type CheckResult struct {
    Status      string        `json:"status"`
    Latency     time.Duration `json:"latency_ms"`
    Message     string        `json:"message,omitempty"`
    LastSuccess time.Time     `json:"last_success,omitempty"`
}

func (h *HealthHandler) Check(c *gin.Context) {
    detailed := c.Query("detail") == "1"

    response := HealthResponse{
        Timestamp: time.Now(),
        Version:   version.Current,
    }

    if detailed {
        response.Checks = map[string]CheckResult{
            "postgres": h.checkPostgres(),
            "mongodb":  h.checkMongoDB(),
            "rabbitmq": h.checkRabbitMQ(),
            "s3":       h.checkS3(),
        }
    }

    // Determinar status general
    allHealthy := true
    for _, check := range response.Checks {
        if check.Status != "healthy" {
            allHealthy = false
            break
        }
    }

    if allHealthy {
        response.Status = "healthy"
        c.JSON(http.StatusOK, response)
    } else {
        response.Status = "degraded"
        c.JSON(http.StatusServiceUnavailable, response)
    }
}
```

---

## 📊 Priorización de Refactorizaciones

| ID | Título | Impacto | Esfuerzo | Prioridad |
|----|--------|---------|----------|-----------|
| REF-005 | Request ID y Tracing | Alto | Medio | 🔴 1 |
| REF-004 | Circuit Breaker | Alto | Bajo | 🔴 2 |
| REF-006 | Healthcheck Detallado | Medio | Bajo | 🟡 3 |
| REF-001 | Scoring Engine | Medio | Medio | 🟡 4 |
| REF-003 | Error Types | Medio | Medio | 🟢 5 |
| REF-002 | Progress Repository | Bajo | Alto | 🟢 6 |

---

## 🗓️ Historial de Refactorizaciones

| Fecha | ID | PR | Descripción |
|-------|-----|-----|-------------|
| - | - | - | Ninguna refactorización completada aún |
