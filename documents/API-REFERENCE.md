# 📡 API Reference

## Base URL

| Ambiente | URL |
|----------|-----|
| Local | `http://localhost:8080` |
| Development | `https://api-mobile-dev.edugo.com` |
| Production | `https://api-mobile.edugo.com` |

**Swagger UI:** `{BASE_URL}/swagger/index.html`

---

## 🔐 Autenticación

Todos los endpoints (excepto `/health`) requieren autenticación JWT.

```http
Authorization: Bearer <jwt_token>
```

### Obtener Token
Los tokens se obtienen desde `api-admin`:
```http
POST https://api-admin.edugo.com/v1/auth/login
Content-Type: application/json

{
  "email": "user@school.edu",
  "password": "password123"
}
```

---

## 📋 Resumen de Endpoints

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `GET` | `/v1/materials` | Listar materiales |
| `POST` | `/v1/materials` | Crear material |
| `GET` | `/v1/materials/:id` | Obtener material |
| `GET` | `/v1/materials/:id/versions` | Historial de versiones |
| `POST` | `/v1/materials/:id/upload-url` | URL presignada upload |
| `GET` | `/v1/materials/:id/download-url` | URL presignada download |
| `POST` | `/v1/materials/:id/upload-complete` | Notificar upload completo |
| `GET` | `/v1/materials/:id/summary` | Obtener resumen IA |
| `GET` | `/v1/materials/:id/assessment` | Obtener quiz |
| `POST` | `/v1/materials/:id/assessment/attempts` | Crear intento de quiz |
| `PATCH` | `/v1/materials/:id/progress` | Actualizar progreso (legacy) |
| `GET` | `/v1/materials/:id/stats` | Estadísticas del material |
| `POST` | `/v1/assessments/:id/submit` | Enviar evaluación |
| `GET` | `/v1/attempts/:id/results` | Resultados de intento |
| `GET` | `/v1/users/me/attempts` | Historial de intentos |
| `PUT` | `/v1/progress` | Upsert progreso |
| `GET` | `/v1/stats/global` | Estadísticas globales |

---

## 🏥 Health Check

### GET /health

Verifica el estado del servicio y sus dependencias.

**Autenticación:** No requerida

#### Response 200 - OK
```json
{
  "status": "healthy",
  "postgres": "connected",
  "mongodb": "connected",
  "timestamp": "2024-12-06T10:00:00Z"
}
```

#### Response 503 - Service Unavailable
```json
{
  "status": "unhealthy",
  "postgres": "disconnected",
  "mongodb": "connected",
  "timestamp": "2024-12-06T10:00:00Z"
}
```

---

## 📚 Materials

### GET /v1/materials

Lista todos los materiales disponibles.

**Autenticación:** Requerida

#### Response 200
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "school_id": "660e8400-e29b-41d4-a716-446655440001",
    "uploaded_by_teacher_id": "770e8400-e29b-41d4-a716-446655440002",
    "title": "Introduction to Calculus",
    "description": "A comprehensive guide to differential calculus",
    "subject": "Mathematics",
    "grade": "12th Grade",
    "file_url": "materials/550e8400/calculus.pdf",
    "file_type": "application/pdf",
    "file_size_bytes": 1048576,
    "status": "ready",
    "is_public": false,
    "created_at": "2024-12-06T10:00:00Z",
    "updated_at": "2024-12-06T10:00:00Z"
  }
]
```

---

### POST /v1/materials

Crea un nuevo material educativo.

**Autenticación:** Requerida

#### Request Body
```json
{
  "title": "Introduction to Calculus",
  "description": "A comprehensive guide to differential calculus",
  "subject": "Mathematics",
  "grade": "12th Grade"
}
```

| Campo | Tipo | Requerido | Validación |
|-------|------|-----------|------------|
| `title` | string | ✅ | 3-200 caracteres |
| `description` | string | ❌ | máx 1000 caracteres |
| `subject` | string | ❌ | - |
| `grade` | string | ❌ | - |

#### Response 201 - Created
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "school_id": "660e8400-e29b-41d4-a716-446655440001",
  "uploaded_by_teacher_id": "770e8400-e29b-41d4-a716-446655440002",
  "title": "Introduction to Calculus",
  "description": "A comprehensive guide to differential calculus",
  "subject": "Mathematics",
  "grade": "12th Grade",
  "file_url": "",
  "file_type": "",
  "file_size_bytes": 0,
  "status": "uploaded",
  "is_public": false,
  "created_at": "2024-12-06T10:00:00Z",
  "updated_at": "2024-12-06T10:00:00Z"
}
```

#### Response 400 - Bad Request
```json
{
  "error": "title is required",
  "code": "VALIDATION_ERROR"
}
```

---

### GET /v1/materials/:id

Obtiene un material por su ID.

**Autenticación:** Requerida

#### Path Parameters
| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `id` | UUID | ID del material |

#### Response 200
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "school_id": "660e8400-e29b-41d4-a716-446655440001",
  "uploaded_by_teacher_id": "770e8400-e29b-41d4-a716-446655440002",
  "title": "Introduction to Calculus",
  "description": "A comprehensive guide to differential calculus",
  "subject": "Mathematics",
  "grade": "12th Grade",
  "file_url": "materials/550e8400/calculus.pdf",
  "file_type": "application/pdf",
  "file_size_bytes": 1048576,
  "status": "ready",
  "processing_started_at": "2024-12-06T10:01:00Z",
  "processing_completed_at": "2024-12-06T10:05:00Z",
  "is_public": false,
  "created_at": "2024-12-06T10:00:00Z",
  "updated_at": "2024-12-06T10:05:00Z"
}
```

#### Response 404 - Not Found
```json
{
  "error": "material not found",
  "code": "NOT_FOUND"
}
```

---

### GET /v1/materials/:id/versions

Obtiene el material con su historial de versiones.

**Autenticación:** Requerida

#### Response 200
```json
{
  "material": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Introduction to Calculus v2",
    "...": "..."
  },
  "versions": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440003",
      "material_id": "550e8400-e29b-41d4-a716-446655440000",
      "version_number": 2,
      "title": "Introduction to Calculus v2",
      "content_url": "materials/550e8400/calculus-v2.pdf",
      "changed_by": "770e8400-e29b-41d4-a716-446655440002",
      "created_at": "2024-12-07T10:00:00Z"
    },
    {
      "id": "990e8400-e29b-41d4-a716-446655440004",
      "material_id": "550e8400-e29b-41d4-a716-446655440000",
      "version_number": 1,
      "title": "Introduction to Calculus",
      "content_url": "materials/550e8400/calculus.pdf",
      "changed_by": "770e8400-e29b-41d4-a716-446655440002",
      "created_at": "2024-12-06T10:00:00Z"
    }
  ]
}
```

---

### POST /v1/materials/:id/upload-url

Genera una URL presignada para subir un archivo a S3.

**Autenticación:** Requerida

#### Request Body
```json
{
  "file_name": "calculus.pdf",
  "content_type": "application/pdf"
}
```

#### Response 200
```json
{
  "upload_url": "https://s3.amazonaws.com/edugo-materials/materials/550e8400/calculus.pdf?X-Amz-Algorithm=AWS4-HMAC-SHA256&...",
  "file_url": "materials/550e8400/calculus.pdf",
  "expires_in": 900
}
```

| Campo | Descripción |
|-------|-------------|
| `upload_url` | URL presignada para PUT del archivo |
| `file_url` | Key de S3 para referencia |
| `expires_in` | Segundos hasta expiración (15 min) |

---

### GET /v1/materials/:id/download-url

Genera una URL presignada para descargar un archivo de S3.

**Autenticación:** Requerida

#### Response 200
```json
{
  "download_url": "https://s3.amazonaws.com/edugo-materials/materials/550e8400/calculus.pdf?X-Amz-Algorithm=AWS4-HMAC-SHA256&...",
  "expires_in": 3600
}
```

#### Response 404 - File Not Uploaded
```json
{
  "error": "material file not uploaded yet",
  "code": "FILE_NOT_FOUND"
}
```

---

### POST /v1/materials/:id/upload-complete

Notifica que el archivo fue subido exitosamente a S3.

**Autenticación:** Requerida

#### Request Body
```json
{
  "file_url": "materials/550e8400/calculus.pdf",
  "file_type": "application/pdf",
  "file_size_bytes": 1048576
}
```

#### Response 204 - No Content
Sin cuerpo. El material pasa a estado `processing`.

---

### GET /v1/materials/:id/summary

Obtiene el resumen generado por IA del material.

**Autenticación:** Requerida

#### Response 200
```json
{
  "material_id": "550e8400-e29b-41d4-a716-446655440000",
  "summary_text": "Este documento presenta los fundamentos del cálculo diferencial, comenzando con el concepto de límites...",
  "key_points": [
    "Definición de límite",
    "Concepto de derivada",
    "Reglas de derivación",
    "Aplicaciones prácticas"
  ],
  "word_count": 250,
  "reading_time_minutes": 2,
  "generated_at": "2024-12-06T10:05:00Z"
}
```

#### Response 404 - Summary Not Ready
```json
{
  "error": "summary not found for this material",
  "code": "NOT_FOUND"
}
```

---

## 📝 Assessments (Evaluaciones)

### GET /v1/materials/:id/assessment

Obtiene el cuestionario de un material **SIN respuestas correctas**.

**Autenticación:** Requerida

#### Response 200
```json
{
  "assessment_id": "aa0e8400-e29b-41d4-a716-446655440000",
  "material_id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Quiz: Introduction to Calculus",
  "total_questions": 10,
  "pass_threshold": 70,
  "max_attempts": null,
  "time_limit_minutes": 30,
  "estimated_time_minutes": 15,
  "questions": [
    {
      "id": "q1",
      "text": "¿Cuál es la derivada de x²?",
      "type": "multiple_choice",
      "options": [
        { "id": "a", "text": "x" },
        { "id": "b", "text": "2x" },
        { "id": "c", "text": "x²" },
        { "id": "d", "text": "2" }
      ]
    },
    {
      "id": "q2",
      "text": "La integral es la operación inversa de la derivada",
      "type": "true_false",
      "options": [
        { "id": "true", "text": "Verdadero" },
        { "id": "false", "text": "Falso" }
      ]
    }
  ]
}
```

> ⚠️ **Nota:** Las respuestas correctas y explicaciones NO se incluyen en esta respuesta para evitar trampas.

---

### POST /v1/materials/:id/assessment/attempts

Crea un intento de evaluación, calcula score y retorna feedback.

**Autenticación:** Requerida

#### Request Body
```json
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

| Campo | Tipo | Requerido | Validación |
|-------|------|-----------|------------|
| `answers` | array | ✅ | mín 1 elemento |
| `answers[].question_id` | string | ✅ | - |
| `answers[].selected_answer_id` | string | ✅ | - |
| `answers[].time_spent_seconds` | int | ✅ | ≥ 0 |
| `time_spent_seconds` | int | ✅ | 1-7200 (2h max) |

#### Response 201 - Created
```json
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
      "message": "¡Correcto! La derivada de x^n es n*x^(n-1)"
    },
    {
      "question_id": "q2",
      "question_text": "La integral es la operación inversa de la derivada",
      "selected_option": "Verdadero",
      "correct_answer": "Verdadero",
      "is_correct": true,
      "message": "¡Correcto!"
    }
  ],
  "can_retake": true,
  "previous_best_score": 60
}
```

---

### POST /v1/assessments/:id/submit

**Legacy endpoint** - Envía evaluación con scoring automático.

**Autenticación:** Requerida

#### Request Body
```json
{
  "responses": {
    "q1": "b",
    "q2": true,
    "q3": "Paris"
  }
}
```

#### Response 200
```json
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
      "explanation": "La derivada de x² es 2x"
    }
  ]
}
```

#### Response 409 - Conflict
```json
{
  "error": "assessment already completed by this user",
  "code": "ASSESSMENT_ALREADY_COMPLETED"
}
```

---

### GET /v1/attempts/:id/results

Obtiene los resultados detallados de un intento específico.

**Autenticación:** Requerida (solo el dueño del intento)

#### Response 200
```json
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
  "feedback": [...],
  "can_retake": true
}
```

#### Response 403 - Forbidden
```json
{
  "error": "attempt does not belong to user",
  "code": "FORBIDDEN"
}
```

---

### GET /v1/users/me/attempts

Obtiene el historial de intentos del usuario autenticado.

**Autenticación:** Requerida

#### Query Parameters
| Parámetro | Tipo | Default | Descripción |
|-----------|------|---------|-------------|
| `limit` | int | 10 | Máx resultados (1-100) |
| `offset` | int | 0 | Resultados a saltar |

#### Response 200
```json
{
  "attempts": [
    {
      "attempt_id": "bb0e8400-e29b-41d4-a716-446655440000",
      "assessment_id": "aa0e8400-e29b-41d4-a716-446655440000",
      "material_id": "550e8400-e29b-41d4-a716-446655440000",
      "material_title": "Introduction to Calculus",
      "score": 80,
      "max_score": 100,
      "passed": true,
      "time_spent_seconds": 180,
      "completed_at": "2024-12-06T14:33:00Z"
    }
  ],
  "total_count": 15,
  "page": 1,
  "limit": 10
}
```

---

## 📊 Progress (Progreso)

### PATCH /v1/materials/:id/progress

**Legacy endpoint** - Actualiza progreso de lectura.

**Autenticación:** Requerida

#### Request Body
```json
{
  "percentage": 75,
  "last_page": 45
}
```

#### Response 204 - No Content
Sin cuerpo.

---

### PUT /v1/progress

**Nuevo endpoint** - Upsert idempotente de progreso.

**Autenticación:** Requerida

#### Request Body
```json
{
  "user_id": "660e8400-e29b-41d4-a716-446655440001",
  "material_id": "550e8400-e29b-41d4-a716-446655440000",
  "progress_percentage": 75,
  "last_page": 45
}
```

| Campo | Tipo | Requerido | Validación |
|-------|------|-----------|------------|
| `user_id` | UUID | ✅ | Debe coincidir con JWT |
| `material_id` | UUID | ✅ | Debe existir |
| `progress_percentage` | int | ✅ | 0-100 |
| `last_page` | int | ❌ | ≥ 0 |

#### Response 200
```json
{
  "user_id": "660e8400-e29b-41d4-a716-446655440001",
  "material_id": "550e8400-e29b-41d4-a716-446655440000",
  "progress_percentage": 75,
  "last_page": 45,
  "message": "progress updated successfully"
}
```

#### Response 403 - Forbidden
```json
{
  "error": "you can only update your own progress",
  "code": "FORBIDDEN"
}
```

---

## 📈 Stats (Estadísticas)

### GET /v1/materials/:id/stats

Obtiene estadísticas de un material específico.

**Autenticación:** Requerida

#### Response 200
```json
{
  "material_id": "550e8400-e29b-41d4-a716-446655440000",
  "total_views": 150,
  "unique_users": 45,
  "average_progress": 68.5,
  "completion_rate": 42.0,
  "average_assessment_score": 75.3,
  "total_assessment_attempts": 38
}
```

---

### GET /v1/stats/global

Obtiene estadísticas globales del sistema.

**Autenticación:** Requerida (solo admins - TODO)

#### Response 200
```json
{
  "total_materials": 250,
  "active_users_last_30_days": 1200,
  "completed_assessments": 5400,
  "average_assessment_score": 72.5,
  "average_progress": 58.3
}
```

---

## ❌ Códigos de Error

### Estructura de Error
```json
{
  "error": "descripción del error",
  "code": "ERROR_CODE"
}
```

### Códigos Comunes

| Código HTTP | Error Code | Descripción |
|-------------|------------|-------------|
| 400 | `INVALID_REQUEST` | Request body inválido |
| 400 | `VALIDATION_ERROR` | Error de validación de campos |
| 400 | `INVALID_MATERIAL_ID` | UUID de material inválido |
| 400 | `INVALID_ATTEMPT_ID` | UUID de intento inválido |
| 400 | `INVALID_FILENAME` | Nombre de archivo con caracteres inválidos |
| 401 | `UNAUTHORIZED` | Token JWT faltante o inválido |
| 403 | `FORBIDDEN` | Sin permisos para la acción |
| 404 | `NOT_FOUND` | Recurso no encontrado |
| 404 | `FILE_NOT_FOUND` | Archivo no subido aún |
| 409 | `ASSESSMENT_ALREADY_COMPLETED` | Evaluación ya completada |
| 500 | `INTERNAL_ERROR` | Error interno del servidor |
| 500 | `S3_ERROR` | Error con AWS S3 |
| 500 | `DATABASE_ERROR` | Error de base de datos |

---

## 🔄 Rate Limiting

> ⚠️ **Nota:** Rate limiting no está implementado actualmente en api-mobile. Se recomienda implementar en el API Gateway o load balancer.

### Límites Sugeridos
| Endpoint | Límite |
|----------|--------|
| Auth endpoints | 5 req/min por IP |
| GET endpoints | 100 req/min por usuario |
| POST/PUT endpoints | 30 req/min por usuario |
| Upload URLs | 10 req/min por usuario |

---

## 📝 Notas de Versionamiento

- **Versión actual:** v1
- **Base path:** `/v1/...`
- **Compatibilidad:** Backward compatible dentro de v1
- **Deprecation:** Endpoints legacy marcados en documentación
