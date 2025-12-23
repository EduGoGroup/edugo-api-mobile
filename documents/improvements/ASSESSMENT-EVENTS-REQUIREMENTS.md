# Requerimientos para Eventos de Assessment

> **Fecha:** 2024-12-23  
> **Estado:** Documentación  
> **Relacionado:** TODO-005 del Sprint de Mejoras

---

## Contexto

El sistema de assessments en `edugo-api-mobile` tiene definido el evento `assessment.generated` pero actualmente **no se está publicando** porque el flujo de generación de assessments ocurre en el worker (`edugo-worker`), no en la API mobile.

Este documento define los eventos necesarios para el sistema de assessments y los requerimientos para integrarlos correctamente.

---

## Eventos Existentes

### `assessment.generated` (Definido, no publicado desde api-mobile)

**Ubicación:** `internal/infrastructure/messaging/rabbitmq/events.go:60-77`

```go
type AssessmentGeneratedPayload struct {
    MaterialID       string `json:"material_id"`
    MongoDocumentID  string `json:"mongo_document_id"`
    QuestionsCount   int    `json:"questions_count"`
    ProcessingTimeMs int    `json:"processing_time_ms,omitempty"`
}
```

**Productor esperado:** `edugo-worker` (después de generar assessment con IA)  
**Consumidor:** `edugo-api-mobile` o servicios de notificación

---

## Eventos Propuestos para api-mobile

### 1. `assessment.attempt.completed`

**Propósito:** Notificar cuando un estudiante completa un intento de evaluación.

```go
type AssessmentAttemptCompletedPayload struct {
    AttemptID      string    `json:"attempt_id"`
    AssessmentID   string    `json:"assessment_id"`
    MaterialID     string    `json:"material_id"`
    StudentID      string    `json:"student_id"`
    SchoolID       string    `json:"school_id"`
    Score          int       `json:"score"`           // 0-100
    Passed         bool      `json:"passed"`
    AttemptNumber  int       `json:"attempt_number"`  // 1, 2, 3...
    TimeSpentSecs  int       `json:"time_spent_seconds"`
    CompletedAt    time.Time `json:"completed_at"`
}
```

**Productor:** `AssessmentAttemptService.CreateAttempt()`  
**Consumidores potenciales:**
- Sistema de gamificación (badges, XP)
- Sistema de notificaciones (notificar al profesor)
- Analytics y reportes
- Sistema de progreso del estudiante

**Ubicación sugerida para publicación:**
```go
// internal/application/service/assessment_attempt_service.go
// Después de persistir el intento exitosamente (línea ~189)
func (s *assessmentAttemptService) CreateAttempt(...) {
    // ... después de guardar answers

    // Publicar evento assessment.attempt.completed
    payload := rabbitmq.AssessmentAttemptCompletedPayload{
        AttemptID:     attempt.ID.String(),
        AssessmentID:  assessment.ID.String(),
        // ...
    }
    // s.publisher.PublishAssessmentAttemptCompleted(ctx, payload)
}
```

---

### 2. `assessment.attempt.started`

**Propósito:** Notificar cuando un estudiante inicia un intento (para tracking de tiempo real).

```go
type AssessmentAttemptStartedPayload struct {
    AttemptID     string    `json:"attempt_id"`
    AssessmentID  string    `json:"assessment_id"`
    MaterialID    string    `json:"material_id"`
    StudentID     string    `json:"student_id"`
    SchoolID      string    `json:"school_id"`
    AttemptNumber int       `json:"attempt_number"`
    StartedAt     time.Time `json:"started_at"`
}
```

**Nota:** Actualmente el servicio no tiene un endpoint separado para "iniciar" un intento. El intento se crea y completa en una sola llamada. Si se requiere tracking de inicio, se necesitaría:
1. Crear endpoint `POST /assessments/:id/attempts/start`
2. Separar flujo de inicio y envío de respuestas

---

### 3. `assessment.first_passed`

**Propósito:** Notificar la primera vez que un estudiante aprueba un assessment (para logros/badges).

```go
type AssessmentFirstPassedPayload struct {
    AssessmentID   string    `json:"assessment_id"`
    MaterialID     string    `json:"material_id"`
    StudentID      string    `json:"student_id"`
    SchoolID       string    `json:"school_id"`
    Score          int       `json:"score"`
    AttemptNumber  int       `json:"attempt_number"`  // En qué intento lo logró
    PassedAt       time.Time `json:"passed_at"`
}
```

**Lógica:** Solo publicar si `passed == true` Y es el primer intento aprobado del estudiante para ese assessment.

---

## Cambios Requeridos en Infraestructura

### 1. Inyectar Publisher en AssessmentAttemptService

**Archivo:** `internal/application/service/assessment_attempt_service.go`

```go
type assessmentAttemptService struct {
    // ... repos existentes
    publisher rabbitmq.Publisher  // AGREGAR
    logger    logger.Logger
}

func NewAssessmentAttemptService(
    // ... params existentes
    publisher rabbitmq.Publisher,  // AGREGAR
    logger logger.Logger,
) AssessmentAttemptService {
    // ...
}
```

**Archivo:** `internal/container/services.go`

```go
AssessmentAttemptService: service.NewAssessmentAttemptService(
    repos.AssessmentRepoV2,
    repos.AttemptRepo,
    repos.AnswerRepo,
    repos.AssessmentDocumentRepo,
    infra.MessagePublisher,  // AGREGAR
    infra.Logger,
),
```

### 2. Agregar Payloads y Constructores en events.go

**Archivo:** `internal/infrastructure/messaging/rabbitmq/events.go`

Agregar las estructuras de payload y funciones constructoras para cada evento nuevo.

### 3. Agregar Métodos de Publicación en EventPublisher

**Archivo:** `internal/infrastructure/messaging/rabbitmq/event_publisher.go`

```go
func (p *EventPublisher) PublishAssessmentAttemptCompleted(ctx context.Context, payload AssessmentAttemptCompletedPayload) error {
    event := NewAssessmentAttemptCompletedEvent(payload)
    // validar y publicar
}
```

---

## Prioridad de Implementación

| Evento | Prioridad | Justificación |
|--------|-----------|---------------|
| `assessment.attempt.completed` | 🔴 Alta | Fundamental para tracking y analytics |
| `assessment.first_passed` | 🟡 Media | Útil para gamificación |
| `assessment.attempt.started` | 🟢 Baja | Requiere cambio de arquitectura |

---

## Dependencias Externas

### edugo-worker

El worker actualmente genera assessments y debería publicar `assessment.generated`. Verificar:
- [ ] ¿Está configurado el publisher en worker?
- [ ] ¿Se publica el evento después de guardar en MongoDB?

### edugo-infrastructure

No se requieren cambios en infrastructure para estos eventos. Los schemas de PostgreSQL y MongoDB ya existen.

---

## Plan de Implementación

### Fase 1: Implementar `assessment.attempt.completed` (Estimado: 2-3h)

1. Agregar payload y constructor en `events.go`
2. Agregar método en `event_publisher.go`
3. Modificar `NewAssessmentAttemptService` para recibir publisher
4. Actualizar `container/services.go`
5. Publicar evento en `CreateAttempt()` después de guardar
6. Actualizar tests

### Fase 2: Implementar `assessment.first_passed` (Estimado: 1-2h)

1. Agregar payload y constructor
2. Agregar lógica para detectar "primer aprobado"
3. Publicar solo cuando corresponda

### Fase 3: Evaluar `assessment.attempt.started` (Futuro)

Requiere decisión de producto sobre si separar inicio y envío de intentos.

---

## Notas Adicionales

- Los eventos deben usar el envelope estándar (`Event` struct)
- Todos los eventos deben incluir `school_id` para multi-tenancy
- Considerar agregar `correlation_id` para tracing distribuido
- Los eventos son fire-and-forget (no bloquear flujo principal si falla publicación)

---

## Referencias

- `internal/infrastructure/messaging/rabbitmq/events.go` - Eventos existentes
- `internal/application/service/assessment_attempt_service.go` - Servicio de intentos
- `internal/application/service/progress_service.go` - Ejemplo de publicación de eventos
