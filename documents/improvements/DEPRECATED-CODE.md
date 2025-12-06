# 🗑️ Código Deprecado

> **Última revisión:** Diciembre 2024  
> **Estado:** Pendiente de eliminación

Este documento lista todo el código marcado como DEPRECATED que debe ser eliminado en futuras versiones.

---

## DEP-001: Funciones WithInjected* en Bootstrap

**Archivo:** `internal/bootstrap/bootstrap.go`  
**Líneas:** 138-166  
**Prioridad:** 🔴 Alta  
**Esfuerzo:** 30 minutos

### Descripción

Existen 5 funciones wrapper deprecadas que simplemente llaman a las funciones equivalentes sin el prefijo `Injected`. Estas fueron mantenidas temporalmente para compatibilidad pero ya no son necesarias.

### Código a Eliminar

```go
// internal/bootstrap/bootstrap.go:138-166

// WithInjectedLogger permite inyectar un logger pre-construido (para tests)
// DEPRECATED: Usar WithLogger() directamente
func WithInjectedLogger(log logger.Logger) BootstrapOption {
	return WithLogger(log)
}

// WithInjectedPostgreSQL permite inyectar una conexión PostgreSQL (para tests)
// DEPRECATED: Usar WithPostgreSQL() directamente
func WithInjectedPostgreSQL(db *sql.DB) BootstrapOption {
	return WithPostgreSQL(db)
}

// WithInjectedMongoDB permite inyectar una conexión MongoDB (para tests)
// DEPRECATED: Usar WithMongoDB() directamente
func WithInjectedMongoDB(db *mongo.Database) BootstrapOption {
	return WithMongoDB(db)
}

// WithInjectedRabbitMQ permite inyectar un publisher RabbitMQ (para tests)
// DEPRECATED: Usar WithRabbitMQ() directamente
func WithInjectedRabbitMQ(pub rabbitmq.Publisher) BootstrapOption {
	return WithRabbitMQ(pub)
}

// WithInjectedS3Client permite inyectar un cliente S3 (para tests)
// DEPRECATED: Usar WithS3Client() directamente
func WithInjectedS3Client(client S3Storage) BootstrapOption {
	return WithS3Client(client)
}
```

### Pasos de Migración

1. **Buscar usos actuales:**
   ```bash
   grep -r "WithInjected" --include="*.go" .
   ```

2. **Reemplazar usos:**
   - `WithInjectedLogger` → `WithLogger`
   - `WithInjectedPostgreSQL` → `WithPostgreSQL`
   - `WithInjectedMongoDB` → `WithMongoDB`
   - `WithInjectedRabbitMQ` → `WithRabbitMQ`
   - `WithInjectedS3Client` → `WithS3Client`

3. **Eliminar las funciones deprecadas**

4. **Ejecutar tests:**
   ```bash
   make test
   ```

### Impacto
- Ningún impacto en runtime
- Posible impacto en tests si usan las funciones antiguas
- Mejora la claridad del API

---

## DEP-002: Repositorio Legacy de Assessments (MongoDB)

**Archivo:** `internal/container/repositories.go`  
**Líneas:** 28  
**Prioridad:** 🟡 Media  
**Esfuerzo:** 2-4 horas

### Descripción

Coexisten dos sistemas de repositorios para assessments:

1. **Legacy (MongoDB):** `AssessmentRepository` - Sistema original
2. **Nuevo (PostgreSQL):** `AssessmentRepoV2` - Sistema Sprint-03

El sistema legacy debería eliminarse una vez que todos los clientes migren al nuevo.

### Código Actual

```go
// internal/container/repositories.go

type RepositoryContainer struct {
    // ...

    // Sprint-03: Assessment System Repositories (PostgreSQL)
    AssessmentRepoV2 repositories.AssessmentRepository // Nuevo de Sprint-03
    AttemptRepo      repositories.AttemptRepository
    AnswerRepo       repositories.AnswerRepository

    // MongoDB Repositories
    SummaryRepository      repository.SummaryRepository
    AssessmentRepository   repository.AssessmentRepository // Legacy  ← ELIMINAR
    AssessmentDocumentRepo mongoRepo.AssessmentDocumentRepository
}
```

### Archivos Relacionados a Revisar

```
internal/
├── container/
│   ├── repositories.go           # Declaración del repositorio
│   └── factory.go                # CreateLegacyAssessmentRepository()
├── domain/
│   └── repository/
│       └── assessment_repository.go  # Interface legacy
├── infrastructure/
│   └── persistence/
│       ├── mongodb/
│       │   └── repository/
│       │       └── assessment_repository.go  # Implementación MongoDB
│       └── mock/
│           └── mongodb/
│               └── assessment_repository.go  # Mock implementation
└── application/
    └── service/
        └── assessment_service.go  # Usa el repositorio legacy
```

### Pasos de Migración

1. **Identificar consumidores:**
   ```bash
   grep -r "AssessmentRepository" --include="*.go" . | grep -v "V2"
   ```

2. **Migrar `AssessmentService` al nuevo sistema**

3. **Eliminar archivos:**
   - `internal/domain/repository/assessment_repository.go`
   - `internal/infrastructure/persistence/mongodb/repository/assessment_repository.go`
   - `internal/infrastructure/persistence/mock/mongodb/assessment_repository.go`

4. **Actualizar container:**
   - Eliminar campo `AssessmentRepository`
   - Eliminar método `CreateLegacyAssessmentRepository()`

5. **Ejecutar tests de integración**

### Riesgos
- Los datos existentes en MongoDB deben migrarse o mantenerse accesibles
- Clientes externos usando el endpoint `/assessments/:id/submit` deben migrarse

---

## DEP-003: Endpoint Legacy UpdateProgress (PATCH)

**Archivo:** `internal/infrastructure/http/handler/progress_handler.go`  
**Líneas:** 26-65  
**Router:** `internal/infrastructure/http/router/router.go:92`  
**Prioridad:** 🟡 Media  
**Esfuerzo:** 1-2 horas

### Descripción

El endpoint `PATCH /v1/materials/:id/progress` es legacy. El nuevo endpoint `PUT /v1/progress` es idempotente y más RESTful.

### Código Legacy

```go
// Handler
// @Summary Update reading progress (legacy endpoint)
// @Description Updates user's reading progress for a material (percentage and last page read).
//              Legacy endpoint - consider using PUT /progress instead.
// @Router /v1/materials/{id}/progress [patch]
func (h *ProgressHandler) UpdateProgress(c *gin.Context) {
    // ...implementación...
}

// Router
materials.PATCH("/:id/progress", c.Handlers.ProgressHandler.UpdateProgress)
```

### Nuevo Endpoint (Recomendado)

```go
// @Summary Upsert progress idempotently (new UPSERT endpoint)
// @Description Updates user progress using idempotent UPSERT operation.
//              Multiple calls with same data are safe.
// @Router /v1/progress [put]
func (h *ProgressHandler) UpsertProgress(c *gin.Context) {
    // ...implementación...
}
```

### Diferencias

| Aspecto | Legacy (PATCH) | Nuevo (PUT) |
|---------|----------------|-------------|
| **Ruta** | `/v1/materials/:id/progress` | `/v1/progress` |
| **Material ID** | Path parameter | Body parameter |
| **User ID** | Solo del JWT | Body + validación contra JWT |
| **Idempotencia** | No garantizada | Garantizada con UPSERT |
| **Respuesta** | 204 No Content | 200 con datos actualizados |

### Plan de Deprecación

1. **Fase 1:** Agregar header `Deprecation` en respuesta
   ```go
   c.Header("Deprecation", "true")
   c.Header("Sunset", "2025-03-01")
   c.Header("Link", "</v1/progress>; rel=\"successor-version\"")
   ```

2. **Fase 2:** Documentar migración en changelog

3. **Fase 3:** Eliminar endpoint después de fecha de sunset

---

## DEP-004: Endpoint Legacy SubmitAssessment

**Archivo:** `internal/infrastructure/http/handler/assessment_handler.go`  
**Líneas:** 103-193  
**Router:** `internal/infrastructure/http/router/router.go:116`  
**Prioridad:** 🟡 Media  
**Esfuerzo:** 2-3 horas

### Descripción

El endpoint `POST /v1/assessments/:id/submit` es del sistema legacy. El nuevo sistema usa `POST /v1/materials/:id/assessment/attempts`.

### Código Legacy

```go
// @Summary Submit assessment with automatic scoring and detailed feedback
// @Description Calcula automáticamente el puntaje de una evaluación
// @Router /v1/assessments/{id}/submit [post]
func (h *AssessmentHandler) SubmitAssessment(c *gin.Context) {
    // Usa assessmentService (legacy)
    result, err := h.assessmentService.CalculateScore(...)
}
```

### Nuevo Endpoint (Sprint-04)

```go
// @Summary Crear intento de evaluación y obtener calificación
// @Description Crea un intento, valida respuestas en servidor, calcula score
// @Router /v1/materials/{id}/assessment/attempts [post]
func (h *AssessmentHandler) CreateMaterialAttempt(c *gin.Context) {
    // Usa assessmentAttemptService (nuevo)
    result, err := h.assessmentAttemptService.CreateAttempt(...)
}
```

### Diferencias

| Aspecto | Legacy | Nuevo (Sprint-04) |
|---------|--------|-------------------|
| **Ruta** | `/v1/assessments/:id/submit` | `/v1/materials/:id/assessment/attempts` |
| **ID en path** | Assessment ID | Material ID |
| **Formato respuestas** | `map[string]interface{}` | `[]UserAnswerDTO` con tiempo |
| **Persistencia** | MongoDB | PostgreSQL |
| **Tracking de tiempo** | No | Sí, por pregunta |
| **Historial** | No | Sí, con `/users/me/attempts` |

---

## 📋 Checklist de Eliminación

- [ ] DEP-001: Funciones WithInjected* (bootstrap.go)
- [ ] DEP-002: Repositorio Legacy Assessments (MongoDB)
- [ ] DEP-003: Endpoint PATCH /materials/:id/progress
- [ ] DEP-004: Endpoint POST /assessments/:id/submit

---

## 🗓️ Historial de Eliminaciones

| Fecha | Código | PR | Descripción |
|-------|--------|-----|-------------|
| - | - | - | Ninguna eliminación registrada aún |
