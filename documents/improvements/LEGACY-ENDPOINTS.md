# 🔌 Endpoints Legacy

> **Última revisión:** Diciembre 2024  
> **Endpoints legacy identificados:** 2

Este documento detalla los endpoints legacy que deben migrarse o eliminarse, incluyendo guías de migración para clientes.

---

## LEG-001: PATCH /v1/materials/:id/progress

### Información General

| Atributo | Valor |
|----------|-------|
| **Endpoint** | `PATCH /v1/materials/:id/progress` |
| **Handler** | `ProgressHandler.UpdateProgress` |
| **Archivo** | `internal/infrastructure/http/handler/progress_handler.go:40-65` |
| **Router** | `internal/infrastructure/http/router/router.go:92` |
| **Estado** | ⚠️ Legacy - Usar `PUT /v1/progress` |
| **Fecha de Sunset** | Por definir |

### Request Actual

```http
PATCH /v1/materials/550e8400-e29b-41d4-a716-446655440000/progress
Authorization: Bearer <token>
Content-Type: application/json

{
  "percentage": 75,
  "last_page": 45
}
```

### Response Actual

```http
HTTP/1.1 204 No Content
```

### Problemas del Endpoint Legacy

1. **No es idempotente:** Múltiples llamadas pueden causar comportamiento inesperado
2. **No retorna datos:** El cliente no sabe el estado final del progreso
3. **Material ID en path:** Inconsistente con el nuevo diseño REST
4. **No valida user_id explícitamente:** Solo usa el del JWT

---

### Nuevo Endpoint Recomendado

```http
PUT /v1/progress
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": "660e8400-e29b-41d4-a716-446655440001",
  "material_id": "550e8400-e29b-41d4-a716-446655440000",
  "progress_percentage": 75,
  "last_page": 45
}
```

### Response Nuevo

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "user_id": "660e8400-e29b-41d4-a716-446655440001",
  "material_id": "550e8400-e29b-41d4-a716-446655440000",
  "progress_percentage": 75,
  "last_page": 45,
  "message": "progress updated successfully"
}
```

### Ventajas del Nuevo Endpoint

| Característica | Legacy | Nuevo |
|----------------|--------|-------|
| Idempotencia | ❌ No | ✅ Sí (UPSERT) |
| Retorna datos | ❌ No (204) | ✅ Sí (200 con body) |
| Validación explícita | ❌ Solo JWT | ✅ Body + JWT |
| Múltiples llamadas seguras | ❌ | ✅ |

### Guía de Migración para Clientes

#### Paso 1: Actualizar el código del cliente

**Antes (JavaScript):**
```javascript
async function updateProgress(materialId, percentage, lastPage) {
  await fetch(`/v1/materials/${materialId}/progress`, {
    method: 'PATCH',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ percentage, last_page: lastPage })
  });
}
```

**Después (JavaScript):**
```javascript
async function updateProgress(userId, materialId, percentage, lastPage) {
  const response = await fetch('/v1/progress', {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      user_id: userId,
      material_id: materialId,
      progress_percentage: percentage,
      last_page: lastPage
    })
  });
  return response.json(); // Ahora retorna datos
}
```

#### Paso 2: Actualizar manejo de respuesta

```javascript
// El nuevo endpoint retorna el estado actualizado
const result = await updateProgress(userId, materialId, 75, 45);
console.log(`Progress: ${result.progress_percentage}%`);
```

#### Paso 3: Obtener user_id del token

```javascript
// El user_id debe obtenerse del token JWT decodificado
function getUserIdFromToken(token) {
  const payload = JSON.parse(atob(token.split('.')[1]));
  return payload.user_id;
}
```

---

## LEG-002: POST /v1/assessments/:id/submit

### Información General

| Atributo | Valor |
|----------|-------|
| **Endpoint** | `POST /v1/assessments/:id/submit` |
| **Handler** | `AssessmentHandler.SubmitAssessment` |
| **Archivo** | `internal/infrastructure/http/handler/assessment_handler.go:103-193` |
| **Router** | `internal/infrastructure/http/router/router.go:116` |
| **Estado** | ⚠️ Legacy - Usar sistema Sprint-04 |
| **Fecha de Sunset** | Por definir |

### Request Actual

```http
POST /v1/assessments/550e8400-e29b-41d4-a716-446655440000/submit
Authorization: Bearer <token>
Content-Type: application/json

{
  "responses": {
    "q1": "b",
    "q2": true,
    "q3": "Paris"
  }
}
```

### Response Actual

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "id": "result-uuid",
  "score": 85.5,
  "total_questions": 10,
  "correct_answers": 8,
  "feedback": [
    {
      "question_id": "q1",
      "is_correct": true,
      "user_answer": "b",
      "correct_answer": "b",
      "explanation": "Correcto!"
    }
  ]
}
```

### Problemas del Endpoint Legacy

1. **Usa MongoDB:** El nuevo sistema usa PostgreSQL para mejor consistencia
2. **Sin tracking de tiempo:** No registra cuánto tiempo tomó cada respuesta
3. **Sin historial robusto:** Difícil consultar intentos anteriores
4. **ID es de assessment:** Debería ser material_id para consistencia
5. **Formato de respuestas:** Map genérico en lugar de estructura tipada

---

### Nuevo Sistema (Sprint-04)

#### Obtener Quiz

```http
GET /v1/materials/550e8400-e29b-41d4-a716-446655440000/assessment
Authorization: Bearer <token>
```

#### Crear Intento

```http
POST /v1/materials/550e8400-e29b-41d4-a716-446655440000/assessment/attempts
Authorization: Bearer <token>
Content-Type: application/json

{
  "answers": [
    {
      "question_id": "q1",
      "selected_answer_id": "b",
      "time_spent_seconds": 45
    },
    {
      "question_id": "q2",
      "selected_answer_id": "true",
      "time_spent_seconds": 20
    }
  ],
  "time_spent_seconds": 180
}
```

### Response Nuevo

```http
HTTP/1.1 201 Created
Content-Type: application/json

{
  "attempt_id": "bb0e8400-e29b-41d4-a716-446655440000",
  "assessment_id": "aa0e8400-e29b-41d4-a716-446655440000",
  "score": 80,
  "max_score": 100,
  "correct_answers": 8,
  "total_questions": 10,
  "pass_threshold": 70,
  "passed": true,
  "time_spent_seconds": 180,
  "started_at": "2024-12-06T14:30:00Z",
  "completed_at": "2024-12-06T14:33:00Z",
  "feedback": [
    {
      "question_id": "q1",
      "question_text": "¿Cuál es la derivada de x²?",
      "selected_option": "2x",
      "correct_answer": "2x",
      "is_correct": true,
      "message": "¡Correcto!"
    }
  ],
  "can_retake": true,
  "previous_best_score": 60
}
```

### Ventajas del Nuevo Sistema

| Característica | Legacy | Nuevo (Sprint-04) |
|----------------|--------|-------------------|
| Base de datos | MongoDB | PostgreSQL |
| Tracking de tiempo | ❌ | ✅ Por pregunta |
| Historial de intentos | Básico | ✅ Completo con paginación |
| Mejor score anterior | ❌ | ✅ `previous_best_score` |
| Puede reintentar | ❌ No sabe | ✅ `can_retake` |
| Formato tipado | Map genérico | DTOs estructurados |
| Consistencia ACID | ❌ | ✅ |

### Nuevos Endpoints Disponibles

```http
# Obtener quiz (sin respuestas correctas)
GET /v1/materials/:id/assessment

# Crear intento y obtener resultado
POST /v1/materials/:id/assessment/attempts

# Ver resultados de un intento específico
GET /v1/attempts/:id/results

# Ver historial de todos mis intentos
GET /v1/users/me/attempts?limit=10&offset=0
```

### Guía de Migración para Clientes

#### Paso 1: Actualizar obtención de quiz

**Antes:**
```javascript
// El legacy usaba el mismo endpoint para todo
const quiz = await fetch(`/v1/assessments/${assessmentId}`);
```

**Después:**
```javascript
// Nuevo: obtener quiz por material_id
const quiz = await fetch(`/v1/materials/${materialId}/assessment`);
```

#### Paso 2: Actualizar envío de respuestas

**Antes:**
```javascript
const result = await fetch(`/v1/assessments/${id}/submit`, {
  method: 'POST',
  body: JSON.stringify({
    responses: { q1: 'b', q2: true }
  })
});
```

**Después:**
```javascript
const result = await fetch(`/v1/materials/${materialId}/assessment/attempts`, {
  method: 'POST',
  body: JSON.stringify({
    answers: [
      { question_id: 'q1', selected_answer_id: 'b', time_spent_seconds: 45 },
      { question_id: 'q2', selected_answer_id: 'true', time_spent_seconds: 20 }
    ],
    time_spent_seconds: 180
  })
});
```

#### Paso 3: Actualizar UI para nuevos campos

```javascript
const data = await result.json();

// Nuevos campos disponibles
if (data.passed) {
  showSuccessMessage(`¡Aprobaste con ${data.score}/${data.max_score}!`);
} else {
  showRetryMessage(`Necesitas ${data.pass_threshold}% para aprobar`);
}

if (data.can_retake) {
  showRetryButton();
}

if (data.previous_best_score) {
  showImprovement(data.score - data.previous_best_score);
}
```

---

## 📅 Plan de Deprecación

### Timeline Propuesto

```
┌─────────────────────────────────────────────────────────────────┐
│                    TIMELINE DE DEPRECACIÓN                       │
└─────────────────────────────────────────────────────────────────┘

    Enero 2025          Febrero 2025         Marzo 2025
        │                    │                    │
        ▼                    ▼                    ▼
   ┌─────────┐         ┌─────────┐         ┌─────────┐
   │ Fase 1  │         │ Fase 2  │         │ Fase 3  │
   │ Warning │         │ Sunset  │         │ Remove  │
   └─────────┘         └─────────┘         └─────────┘
        │                    │                    │
        │                    │                    │
   Agregar headers      Desactivar          Eliminar
   de deprecación       endpoints           código
```

### Fase 1: Deprecation Warning (Enero 2025)

Agregar headers HTTP de deprecación:

```go
func deprecationMiddleware(sunset time.Time, successor string) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Deprecation", "true")
        c.Header("Sunset", sunset.Format(time.RFC1123))
        c.Header("Link", fmt.Sprintf("<%s>; rel=\"successor-version\"", successor))
        c.Next()
    }
}

// Aplicar al router
materials.PATCH("/:id/progress",
    deprecationMiddleware(time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC), "/v1/progress"),
    c.Handlers.ProgressHandler.UpdateProgress,
)
```

### Fase 2: Sunset Period (Febrero 2025)

- Endpoints siguen funcionando
- Logs de uso para identificar clientes
- Comunicación directa con equipos

### Fase 3: Removal (Marzo 2025)

- Eliminar handlers legacy
- Eliminar rutas del router
- Actualizar documentación
- Actualizar tests

---

## 📊 Monitoreo de Uso

### Queries para Identificar Uso

```sql
-- Si tienes logs en base de datos
SELECT
    endpoint,
    COUNT(*) as calls,
    COUNT(DISTINCT user_id) as unique_users,
    MAX(created_at) as last_used
FROM api_logs
WHERE endpoint IN (
    '/v1/materials/%/progress',
    '/v1/assessments/%/submit'
)
AND created_at >= NOW() - INTERVAL '30 days'
GROUP BY endpoint;
```

### Métricas a Monitorear

- Llamadas por día al endpoint legacy
- Usuarios únicos usando endpoint legacy
- Última fecha de uso
- Errores en migración al nuevo endpoint
