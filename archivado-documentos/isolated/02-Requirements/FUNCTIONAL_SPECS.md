# Especificaciones Funcionales
# Sistema de Evaluaciones - EduGo

**Versión:** 1.0.0  
**Fecha:** 14 de Noviembre, 2025  
**Proyecto:** edugo-api-mobile - Sistema de Evaluaciones

---

## 1. INTRODUCCIÓN

Este documento detalla las especificaciones funcionales del Sistema de Evaluaciones de EduGo, numeradas con formato RF-XXX (Requerimiento Funcional). Cada especificación incluye:
- Descripción clara
- Prioridad MoSCoW
- Criterios de aceptación específicos
- Dependencias

---

## 2. MÓDULO: OBTENCIÓN DE CUESTIONARIOS

### RF-001: Obtener Cuestionario de un Material

**Prioridad:** MUST  
**Módulo:** Evaluaciones  
**Endpoint:** `GET /v1/materials/:id/assessment`

#### Descripción
El sistema DEBE permitir a un estudiante obtener el cuestionario (quiz) asociado a un material educativo específico, sin revelar las respuestas correctas.

#### Criterios de Aceptación

**AC-001.1:** El endpoint DEBE validar que el usuario está autenticado (JWT válido)

**AC-001.2:** El endpoint DEBE verificar que el material existe y su `processing_status = 'completed'`
```sql
SELECT id, processing_status FROM materials WHERE id = $1
```

**AC-001.3:** El sistema DEBE verificar que el usuario tiene permiso para acceder al material
- MVP: Cualquier usuario autenticado puede acceder
- Post-MVP: Validar membresía en unidad académica

**AC-001.4:** El sistema DEBE consultar MongoDB para obtener las preguntas
```javascript
db.material_assessment.findOne({
  material_id: "uuid-material"
})
```

**AC-001.5:** El sistema DEBE remover campos sensibles antes de enviar al cliente:
- `correct_answer` DEBE ser removido
- `feedback.correct` y `feedback.incorrect` DEBEN ser removidos
- Solo enviar: `id`, `text`, `type`, `options[]`

**AC-001.6:** Response DEBE incluir:
```json
{
  "assessment_id": "uuid",
  "material_id": "uuid",
  "title": "string",
  "total_questions": number,
  "estimated_time_minutes": number,
  "questions": [
    {
      "id": "string",
      "text": "string",
      "type": "multiple_choice",
      "options": [
        {"id": "string", "text": "string"}
      ]
    }
  ]
}
```

**AC-001.7:** Error handling:
- 401 si no autenticado
- 404 si material no existe
- 404 si material no tiene assessment (MongoDB)
- 403 si usuario no tiene permiso (Post-MVP)

#### Dependencias
- Tabla `materials` existente en PostgreSQL
- Colección `material_assessment` existente en MongoDB
- Middleware de autenticación funcional

#### Tests Requeridos
1. Test unitario: Sanitización correcta de respuestas
2. Test de integración: Consulta a MongoDB exitosa
3. Test E2E: Flujo completo con material real
4. Test de seguridad: Verificar que `correct_answer` nunca se envía

---

### RF-002: Validar Existencia de Assessment

**Prioridad:** MUST  
**Módulo:** Evaluaciones  
**Endpoint:** Interno (usado por RF-001)

#### Descripción
El sistema DEBE verificar que un material tiene un assessment disponible antes de permitir intentos.

#### Criterios de Aceptación

**AC-002.1:** Consultar PostgreSQL para metadatos del assessment
```sql
SELECT a.id, a.mongo_document_id, a.total_questions, a.pass_threshold
FROM assessment a
WHERE a.material_id = $1
```

**AC-002.2:** Si no existe registro en PostgreSQL, intentar crear uno automáticamente consultando MongoDB

**AC-002.3:** Si MongoDB tampoco tiene assessment, retornar error 404 con mensaje claro:
```json
{
  "error": "Assessment not available",
  "message": "Este material aún no tiene un cuestionario disponible"
}
```

#### Dependencias
- Tabla `assessment` creada
- Worker que genere assessments y los almacene en MongoDB

---

## 3. MÓDULO: ENVÍO DE RESPUESTAS

### RF-003: Crear Intento de Evaluación

**Prioridad:** MUST  
**Módulo:** Evaluaciones  
**Endpoint:** `POST /v1/materials/:id/assessment/attempts`

#### Descripción
El sistema DEBE permitir a un estudiante enviar sus respuestas a un cuestionario, calcular el puntaje automáticamente, y almacenar el intento de forma permanente.

#### Criterios de Aceptación

**AC-003.1:** Request DEBE incluir:
```json
{
  "answers": [
    {
      "question_id": "string",
      "selected_option": "string"
    }
  ],
  "time_spent_seconds": number
}
```

**AC-003.2:** El sistema DEBE validar:
- Todas las preguntas tienen respuesta
- Los `question_id` existen en el assessment
- `time_spent_seconds` es razonable (>0 y <3600)

**AC-003.3:** El sistema DEBE obtener preguntas CON respuestas correctas desde MongoDB
```javascript
db.material_assessment.findOne({
  material_id: "uuid"
})
```

**AC-003.4:** El sistema DEBE calcular puntaje:
```go
correctCount := 0
for _, answer := range studentAnswers {
    question := findQuestion(questions, answer.QuestionID)
    if answer.SelectedOption == question.CorrectAnswer {
        correctCount++
    }
}
score := (correctCount * 100) / totalQuestions
```

**AC-003.5:** El sistema DEBE persistir en PostgreSQL dentro de una transacción ACID:
```sql
BEGIN;

-- 1. Insertar intento
INSERT INTO assessment_attempt (
    id, assessment_id, student_id, score, max_score,
    time_spent_seconds, started_at, completed_at
) VALUES (...);

-- 2. Insertar respuestas individuales
INSERT INTO assessment_attempt_answer (
    attempt_id, question_id, selected_option, is_correct
) VALUES (...);

COMMIT;
```

**AC-003.6:** El sistema DEBE generar feedback educativo por pregunta:
- Si correcta: mensaje de `feedback.correct` de MongoDB
- Si incorrecta: mensaje de `feedback.incorrect` de MongoDB

**AC-003.7:** Response DEBE incluir:
```json
{
  "attempt_id": "uuid",
  "score": number,
  "max_score": 100,
  "correct_answers": number,
  "total_questions": number,
  "pass_threshold": number,
  "passed": boolean,
  "feedback": [
    {
      "question_id": "string",
      "question_text": "string",
      "selected_option": "string",
      "correct_answer": "string",
      "is_correct": boolean,
      "message": "string"
    }
  ],
  "can_retake": boolean,
  "previous_best_score": number | null
}
```

**AC-003.8:** El sistema DEBE retornar en <2 segundos (p95)

**AC-003.9:** Intentos DEBEN ser inmutables (no editables después de creados)

#### Dependencias
- Tabla `assessment_attempt` creada
- Tabla `assessment_attempt_answer` creada
- Integración con MongoDB funcional

#### Tests Requeridos
1. Test unitario: Cálculo correcto de puntaje
2. Test unitario: Generación de feedback
3. Test de integración: Transacción ACID (rollback en caso de error)
4. Test de integración: Persistencia correcta en PostgreSQL
5. Test E2E: Flujo completo desde envío hasta response
6. Test de performance: Latencia <2 seg con 100 requests concurrentes

---

### RF-004: Validar Respuestas en Servidor

**Prioridad:** MUST  
**Módulo:** Evaluaciones - Seguridad  
**Endpoint:** Interno (usado por RF-003)

#### Descripción
El sistema NUNCA DEBE confiar en validaciones del cliente. TODAS las respuestas DEBEN ser validadas en el servidor contra las preguntas almacenadas en MongoDB.

#### Criterios de Aceptación

**AC-004.1:** Validación SIEMPRE ocurre en servidor:
```go
// ❌ NUNCA hacer esto (confiar en cliente)
func ProcessAttempt(attemptDTO AttemptDTO) {
    score := attemptDTO.Score // ❌ Cliente puede mentir
}

// ✅ SIEMPRE hacer esto
func ProcessAttempt(answers []Answer) {
    correctAnswers := fetchFromMongoDB(assessmentID)
    score := calculateScore(answers, correctAnswers)
}
```

**AC-004.2:** Respuestas correctas NUNCA expuestas en endpoint GET /assessment

**AC-004.3:** Validación de integridad:
- Verificar que número de respuestas == número de preguntas
- Verificar que todos los `question_id` existen
- Rechazar si hay preguntas duplicadas
- Rechazar si hay `selected_option` inválidas

**AC-004.4:** Logging de intentos sospechosos:
```go
if timeSpent < 5 * len(questions) { // <5 seg por pregunta
    logger.Warn("Suspicious attempt: too fast",
        "attempt_id", attemptID,
        "time_spent", timeSpent,
        "student_id", studentID)
}
```

#### Tests Requeridos
1. Test de seguridad: Intentar enviar `correct_answer` en request (debe ignorarse)
2. Test de seguridad: Intentar enviar `score` calculado en cliente (debe ignorarse)
3. Test de validación: Respuestas faltantes (debe retornar error 400)
4. Test de validación: `question_id` inexistente (debe retornar error 400)

---

## 4. MÓDULO: CONSULTA DE RESULTADOS

### RF-005: Obtener Resultados de un Intento

**Prioridad:** MUST  
**Módulo:** Evaluaciones  
**Endpoint:** `GET /v1/attempts/:id/results`

#### Descripción
El sistema DEBE permitir a un estudiante consultar los resultados detallados de un intento previo, incluyendo puntaje, respuestas correctas/incorrectas, y feedback.

#### Criterios de Aceptación

**AC-005.1:** El sistema DEBE validar que el usuario autenticado es el propietario del intento:
```sql
SELECT student_id FROM assessment_attempt WHERE id = $1
-- Comparar con user_id del JWT
```

**AC-005.2:** El sistema DEBE retornar 403 si el usuario no es el propietario

**AC-005.3:** El sistema DEBE consultar PostgreSQL para datos del intento:
```sql
SELECT
    aa.id, aa.score, aa.max_score, aa.time_spent_seconds,
    aa.completed_at, a.title, m.title as material_title
FROM assessment_attempt aa
INNER JOIN assessment a ON aa.assessment_id = a.id
INNER JOIN materials m ON a.material_id = m.id
WHERE aa.id = $1
```

**AC-005.4:** El sistema DEBE consultar respuestas individuales:
```sql
SELECT
    aaa.question_id, aaa.selected_option, aaa.is_correct
FROM assessment_attempt_answer aaa
WHERE aaa.attempt_id = $1
ORDER BY aaa.question_id
```

**AC-005.5:** El sistema DEBE enriquecer con datos de MongoDB:
- Texto de preguntas
- Texto de opciones
- Feedback educativo

**AC-005.6:** Response DEBE incluir estructura completa (misma que RF-003 response)

#### Dependencias
- Tabla `assessment_attempt` con datos
- Tabla `assessment_attempt_answer` con datos
- Colección `material_assessment` en MongoDB

#### Tests Requeridos
1. Test de integración: Consulta exitosa de intento propio
2. Test de seguridad: 403 al intentar consultar intento ajeno
3. Test E2E: Flujo crear intento → consultar resultados

---

### RF-006: Obtener Historial de Intentos del Usuario

**Prioridad:** MUST  
**Módulo:** Evaluaciones  
**Endpoint:** `GET /v1/users/me/attempts`

#### Descripción
El sistema DEBE permitir a un estudiante consultar su historial completo de intentos de evaluaciones.

#### Criterios de Aceptación

**AC-006.1:** El sistema DEBE consultar todos los intentos del usuario autenticado:
```sql
SELECT
    aa.id, aa.score, aa.max_score, aa.completed_at,
    m.title as material_title, m.id as material_id,
    a.pass_threshold
FROM assessment_attempt aa
INNER JOIN assessment a ON aa.assessment_id = a.id
INNER JOIN materials m ON a.material_id = m.id
WHERE aa.student_id = $1
ORDER BY aa.completed_at DESC
LIMIT 50
```

**AC-006.2:** El sistema DEBE soportar paginación (query params: `limit`, `offset`)

**AC-006.3:** El sistema DEBE incluir indicador de aprobado/reprobado:
```go
passed := attempt.Score >= assessment.PassThreshold
```

**AC-006.4:** Response DEBE incluir:
```json
{
  "attempts": [
    {
      "attempt_id": "uuid",
      "material_id": "uuid",
      "material_title": "string",
      "score": number,
      "max_score": 100,
      "passed": boolean,
      "completed_at": "ISO8601"
    }
  ],
  "total_count": number,
  "page": number,
  "limit": number
}
```

#### Tests Requeridos
1. Test de integración: Consulta con múltiples intentos
2. Test de integración: Paginación correcta
3. Test E2E: Crear 3 intentos → consultar historial → verificar 3 resultados

---

## 5. MÓDULO: GESTIÓN DE ASSESSMENTS (Admin - Post-MVP)

### RF-007: Crear Assessment Manualmente

**Prioridad:** COULD  
**Módulo:** Administración de Assessments  
**Endpoint:** `POST /v1/materials/:id/assessment` (Admin)

#### Descripción
El sistema DEBERÍA permitir a un administrador o profesor crear un assessment manualmente si el worker no lo generó automáticamente.

#### Criterios de Aceptación

**AC-007.1:** Solo usuarios con rol `teacher` o `admin` pueden acceder

**AC-007.2:** Request DEBE incluir:
```json
{
  "title": "string",
  "pass_threshold": number, // 0-100
  "questions": [
    {
      "text": "string",
      "type": "multiple_choice",
      "options": [
        {"id": "a", "text": "string"}
      ],
      "correct_answer": "string",
      "feedback": {
        "correct": "string",
        "incorrect": "string"
      }
    }
  ]
}
```

**AC-007.3:** El sistema DEBE validar:
- `pass_threshold` entre 0-100
- Cada pregunta tiene al menos 2 opciones
- `correct_answer` existe en `options`
- No hay opciones duplicadas

**AC-007.4:** El sistema DEBE almacenar en MongoDB y crear registro en PostgreSQL

**AC-007.5:** Post-MVP: Integrar con sistema de jerarquía (solo profesores de la unidad)

#### Dependencias
- Sistema de jerarquía académica implementado
- Middleware de autorización por rol

---

## 6. MÓDULO: ANALYTICS Y REPORTES (Post-MVP)

### RF-008: Obtener Estadísticas de un Material

**Prioridad:** SHOULD  
**Módulo:** Analytics  
**Endpoint:** `GET /v1/materials/:id/stats`

#### Descripción
El sistema DEBERÍA proporcionar estadísticas agregadas de rendimiento estudiantil en un material específico.

#### Criterios de Aceptación

**AC-008.1:** Solo accesible por profesores y administradores

**AC-008.2:** El sistema DEBE calcular:
```sql
SELECT
    COUNT(DISTINCT student_id) as total_students,
    AVG(score) as average_score,
    MIN(score) as min_score,
    MAX(score) as max_score,
    COUNT(*) as total_attempts
FROM assessment_attempt aa
INNER JOIN assessment a ON aa.assessment_id = a.id
WHERE a.material_id = $1
```

**AC-008.3:** El sistema DEBE calcular distribución de puntajes:
```go
buckets := []int{0-20, 21-40, 41-60, 61-80, 81-100}
// Histograma de distribución
```

**AC-008.4:** Response DEBE incluir:
```json
{
  "material_id": "uuid",
  "total_students": number,
  "total_attempts": number,
  "average_score": number,
  "min_score": number,
  "max_score": number,
  "pass_rate": number, // % de intentos que pasaron threshold
  "score_distribution": {
    "0-20": number,
    "21-40": number,
    "41-60": number,
    "61-80": number,
    "81-100": number
  }
}
```

#### Tests Requeridos
1. Test de integración: Cálculos correctos con datos de prueba
2. Test de autorización: 403 para estudiantes

---

### RF-009: Identificar Preguntas Problemáticas

**Prioridad:** COULD  
**Módulo:** Analytics  
**Endpoint:** `GET /v1/materials/:id/question-stats`

#### Descripción
El sistema PODRÍA identificar preguntas con alta tasa de error para que profesores las revisen.

#### Criterios de Aceptación

**AC-009.1:** El sistema DEBE calcular tasa de error por pregunta:
```sql
SELECT
    question_id,
    COUNT(*) as total_answers,
    SUM(CASE WHEN is_correct THEN 1 ELSE 0 END) as correct_count,
    ROUND(100.0 * SUM(CASE WHEN NOT is_correct THEN 1 ELSE 0 END) / COUNT(*), 2) as error_rate
FROM assessment_attempt_answer aaa
INNER JOIN assessment_attempt aa ON aaa.attempt_id = aa.id
INNER JOIN assessment a ON aa.assessment_id = a.id
WHERE a.material_id = $1
GROUP BY question_id
ORDER BY error_rate DESC
```

**AC-009.2:** Response DEBE incluir:
```json
{
  "questions": [
    {
      "question_id": "string",
      "question_text": "string",
      "total_answers": number,
      "correct_count": number,
      "error_rate": number, // %
      "is_problematic": boolean // error_rate > 70%
    }
  ]
}
```

**AC-009.3:** Marcar como problemática si error_rate > 70%

---

## 7. MÓDULO: NOTIFICACIONES (Post-MVP)

### RF-010: Notificar Docentes de Intentos Completados

**Prioridad:** SHOULD  
**Módulo:** Notificaciones Asíncronas  
**Componente:** Worker

#### Descripción
El sistema DEBERÍA notificar a los docentes cuando un estudiante completa un cuestionario, especialmente si el puntaje es bajo (<60%).

#### Criterios de Aceptación

**AC-010.1:** API Mobile DEBE publicar evento a RabbitMQ tras crear intento:
```json
{
  "event_type": "assessment_attempt_recorded",
  "attempt_id": "uuid",
  "material_id": "uuid",
  "student_id": "uuid",
  "score": number,
  "timestamp": "ISO8601"
}
```

**AC-010.2:** Worker DEBE consumir evento y:
1. Identificar docentes de la unidad académica
2. Generar notificación (email/push)
3. Enviar solo si score < 60% (configurable)

**AC-010.3:** Notificación DEBE incluir:
- Nombre del estudiante
- Título del material
- Puntaje obtenido
- Link al detalle del intento

#### Dependencias
- Worker funcionando
- RabbitMQ configurado
- Sistema de jerarquía académica (para identificar docentes)
- Servicio de email/push configurado

---

## 8. MÓDULO: REINTENTOS (Post-MVP)

### RF-011: Permitir Múltiples Intentos

**Prioridad:** SHOULD  
**Módulo:** Evaluaciones  
**Endpoint:** Modificación de RF-003

#### Descripción
El sistema DEBERÍA permitir a un estudiante realizar múltiples intentos de un cuestionario, registrando todos los intentos pero mostrando el mejor puntaje.

#### Criterios de Aceptación

**AC-011.1:** Assessment DEBE tener campo `max_attempts` (nullable, default ilimitado)
```sql
ALTER TABLE assessment ADD COLUMN max_attempts INTEGER DEFAULT NULL;
```

**AC-011.2:** El sistema DEBE validar antes de permitir nuevo intento:
```sql
SELECT COUNT(*) as attempt_count
FROM assessment_attempt
WHERE assessment_id = $1 AND student_id = $2
```

**AC-011.3:** Si `attempt_count >= max_attempts`, retornar error 403

**AC-011.4:** El sistema DEBE calcular mejor puntaje:
```sql
SELECT MAX(score) as best_score
FROM assessment_attempt
WHERE assessment_id = $1 AND student_id = $2
```

**AC-011.5:** Response de RF-003 DEBE incluir:
```json
{
  "can_retake": boolean,
  "attempts_used": number,
  "attempts_remaining": number | null,
  "previous_best_score": number | null
}
```

---

## 9. MÓDULO: BANCO ALEATORIO (Post-MVP)

### RF-012: Seleccionar Preguntas Aleatorias

**Prioridad:** WON'T HAVE (Fase 2)  
**Módulo:** Evaluaciones  
**Endpoint:** Modificación de RF-001

#### Descripción
El sistema PODRÍA seleccionar aleatoriamente N preguntas de un banco más grande (ej: 5 de 20) para prevenir memorización de respuestas.

#### Criterios de Aceptación

**AC-012.1:** MongoDB DEBE almacenar banco completo (ej: 20 preguntas)

**AC-012.2:** Assessment DEBE tener campo `questions_per_attempt` (ej: 5)

**AC-012.3:** El sistema DEBE seleccionar aleatoriamente:
```go
selectedQuestions := selectRandom(allQuestions, questionsPerAttempt)
```

**AC-012.4:** El sistema DEBE registrar qué preguntas se mostraron en cada intento

**AC-012.5:** Validación DEBE comparar solo contra preguntas del intento específico

---

## 10. PRIORIZACIÓN MOSCOW

### MUST HAVE (MVP - 2 semanas)
- ✅ RF-001: Obtener Cuestionario
- ✅ RF-002: Validar Existencia de Assessment
- ✅ RF-003: Crear Intento de Evaluación
- ✅ RF-004: Validar Respuestas en Servidor
- ✅ RF-005: Obtener Resultados de un Intento
- ✅ RF-006: Obtener Historial de Intentos

### SHOULD HAVE (Post-MVP - Fase 1)
- 🟡 RF-008: Estadísticas de Material
- 🟡 RF-010: Notificar Docentes
- 🟡 RF-011: Múltiples Intentos

### COULD HAVE (Post-MVP - Fase 2)
- 🟢 RF-007: Crear Assessment Manualmente
- 🟢 RF-009: Preguntas Problemáticas

### WON'T HAVE (Futuro)
- ⚪ RF-012: Banco Aleatorio
- ⚪ Tipos de preguntas avanzadas
- ⚪ Respuestas cortas con NLP

---

## 11. MATRIZ DE TRAZABILIDAD

| RF | Objetivo de Negocio | Criterio de Éxito | Sprint |
|----|---------------------|-------------------|--------|
| RF-001 | OB-01 | CF-01 | Sprint-02 |
| RF-003 | OB-02, OB-03 | CF-02, CT-02 | Sprint-04 |
| RF-004 | Seguridad | CF-05 | Sprint-04 |
| RF-005 | OB-03 | CF-04 | Sprint-04 |
| RF-006 | OB-03 | CF-04 | Sprint-04 |
| RF-008 | OB-04 | Post-MVP | Sprint-06 |

---

**Generado con:** Claude Code  
**Total Especificaciones:** 12 (6 MUST, 3 SHOULD, 2 COULD, 1 WON'T)  
**Última actualización:** 2025-11-14
