# Criterios de Aceptación
# Sistema de Evaluaciones - EduGo

**Versión:** 1.0.0  
**Fecha:** 14 de Noviembre, 2025  
**Proyecto:** edugo-api-mobile - Sistema de Evaluaciones

---

## 1. DEFINICIÓN DE "DONE"

Un requisito funcional se considera **DONE** cuando cumple TODOS estos criterios:

✅ **Código Implementado**
- Funcionalidad completa según especificación
- Código revisado (code review)
- Sin warnings de linting (`golangci-lint run`)

✅ **Tests Pasando**
- Tests unitarios con >85% coverage
- Tests de integración exitosos
- Tests E2E del flujo completo

✅ **Documentación**
- Swagger annotations completas
- Comentarios en funciones públicas (godoc)
- README actualizado si aplica

✅ **CI/CD**
- Pipeline de GitHub Actions verde
- Build exitoso
- Tests automáticos pasando

✅ **Validación Manual**
- Probado manualmente en entorno local
- Casos edge testeados
- Performance aceptable (<2 seg para operaciones críticas)

---

## 2. CRITERIOS POR MÓDULO FUNCIONAL

### 2.1 MÓDULO: Obtención de Cuestionarios

#### AC-MOD-001: GET /v1/materials/:id/assessment

**Criterio de Aceptación 1: Autenticación Requerida**

✅ **GIVEN** un usuario sin token JWT  
✅ **WHEN** intenta acceder a GET /v1/materials/{id}/assessment  
✅ **THEN** el sistema retorna 401 Unauthorized con mensaje de error claro

**Validación:**
```bash
curl -X GET http://localhost:8080/v1/materials/uuid-1/assessment
# Esperado: 401
# {
#   "error": "unauthorized",
#   "message": "missing or invalid token"
# }
```

---

**Criterio de Aceptación 2: Material Existe**

✅ **GIVEN** un usuario autenticado  
✅ **AND** un material_id que NO existe en la base de datos  
✅ **WHEN** solicita GET /v1/materials/{id}/assessment  
✅ **THEN** el sistema retorna 404 Not Found

**Validación:**
```bash
curl -X GET http://localhost:8080/v1/materials/99999999-9999-9999-9999-999999999999/assessment \
  -H "Authorization: Bearer valid-token"
# Esperado: 404
# {
#   "error": "not_found",
#   "message": "Material not found"
# }
```

---

**Criterio de Aceptación 3: Assessment Disponible**

✅ **GIVEN** un material que existe pero su processing_status != 'completed'  
✅ **WHEN** solicita GET /v1/materials/{id}/assessment  
✅ **THEN** el sistema retorna 404 con mensaje "Assessment not available yet"

**Validación SQL:**
```sql
-- Setup: Crear material sin assessment
INSERT INTO materials (id, title, processing_status) 
VALUES ('uuid-test', 'Material Sin Assessment', 'pending');

-- Test: Debe retornar 404
```

---

**Criterio de Aceptación 4: Respuestas Correctas Nunca Expuestas**

✅ **GIVEN** un material con assessment disponible  
✅ **WHEN** solicita GET /v1/materials/{id}/assessment  
✅ **THEN** el response JSON NO contiene campos `correct_answer` ni `feedback`

**Validación:**
```go
// Test de seguridad
func TestAssessment_NeverExposeCorrectAnswers(t *testing.T) {
    resp := getAssessment(t, validMaterialID, validToken)
    
    for _, question := range resp.Questions {
        // ⚠️ CRÍTICO: Estos campos NUNCA deben estar presentes
        assert.Nil(t, question.CorrectAnswer, "correct_answer must be nil")
        assert.Nil(t, question.Feedback, "feedback must be nil")
    }
}
```

---

**Criterio de Aceptación 5: Estructura de Response Correcta**

✅ **GIVEN** un material con assessment disponible  
✅ **WHEN** solicita GET /v1/materials/{id}/assessment  
✅ **THEN** el response contiene:
- `assessment_id` (UUID)
- `material_id` (UUID)
- `title` (string no vacío)
- `total_questions` (entero >0)
- `estimated_time_minutes` (entero >0)
- `questions` (array de objetos con: id, text, type, options)

**Validación:**
```bash
curl -X GET http://localhost:8080/v1/materials/valid-uuid/assessment \
  -H "Authorization: Bearer valid-token"

# Esperado: 200 OK
# {
#   "assessment_id": "uuid",
#   "material_id": "uuid",
#   "title": "Cuestionario: ...",
#   "total_questions": 5,
#   "estimated_time_minutes": 10,
#   "questions": [
#     {
#       "id": "q1",
#       "text": "¿Pregunta?",
#       "type": "multiple_choice",
#       "options": [
#         {"id": "a", "text": "Opción A"},
#         {"id": "b", "text": "Opción B"}
#       ]
#     }
#   ]
# }
```

---

**Criterio de Aceptación 6: Performance <500ms**

✅ **GIVEN** 100 requests concurrentes  
✅ **WHEN** solicitan GET /v1/materials/{id}/assessment  
✅ **THEN** p95 de latencia es <500ms

**Validación:**
```bash
# Load test con Apache Bench
ab -n 1000 -c 100 \
  -H "Authorization: Bearer valid-token" \
  http://localhost:8080/v1/materials/uuid-1/assessment

# Verificar: Time per request (p95) < 500ms
```

---

### 2.2 MÓDULO: Envío de Respuestas

#### AC-MOD-002: POST /v1/materials/:id/assessment/attempts

**Criterio de Aceptación 1: Request Válido**

✅ **GIVEN** un usuario autenticado  
✅ **AND** un material con assessment disponible  
✅ **WHEN** envía POST con todas las respuestas correctamente formateadas  
✅ **THEN** el sistema acepta el request y procesa el intento

**Request Válido:**
```json
{
  "answers": [
    {"question_id": "q1", "selected_option": "a"},
    {"question_id": "q2", "selected_option": "b"},
    {"question_id": "q3", "selected_option": "c"},
    {"question_id": "q4", "selected_option": "d"},
    {"question_id": "q5", "selected_option": "a"}
  ],
  "time_spent_seconds": 420
}
```

---

**Criterio de Aceptación 2: Validación de Respuestas Completas**

✅ **GIVEN** un assessment con 5 preguntas  
✅ **WHEN** envía POST con solo 3 respuestas  
✅ **THEN** el sistema retorna 400 Bad Request con mensaje "incomplete answers"

**Validación:**
```bash
curl -X POST http://localhost:8080/v1/materials/uuid-1/assessment/attempts \
  -H "Authorization: Bearer valid-token" \
  -H "Content-Type: application/json" \
  -d '{
    "answers": [
      {"question_id": "q1", "selected_option": "a"}
    ],
    "time_spent_seconds": 60
  }'

# Esperado: 400
# {
#   "error": "validation_error",
#   "message": "incomplete answers: expected 5, got 1"
# }
```

---

**Criterio de Aceptación 3: Cálculo Correcto de Puntaje**

✅ **GIVEN** un assessment con 5 preguntas  
✅ **AND** respuestas correctas: q1=a, q2=b, q3=c, q4=d, q5=a  
✅ **WHEN** estudiante envía: q1=a (✓), q2=x (✗), q3=c (✓), q4=d (✓), q5=x (✗)  
✅ **THEN** el score calculado es 60/100 (3 de 5 correctas)

**Validación:**
```go
func TestScoringService_CalculateScore(t *testing.T) {
    service := NewScoringService()
    
    correctAnswers := []Question{
        {ID: "q1", CorrectAnswer: "a"},
        {ID: "q2", CorrectAnswer: "b"},
        {ID: "q3", CorrectAnswer: "c"},
        {ID: "q4", CorrectAnswer: "d"},
        {ID: "q5", CorrectAnswer: "a"},
    }
    
    studentAnswers := []Answer{
        {QuestionID: "q1", SelectedOption: "a"},
        {QuestionID: "q2", SelectedOption: "x"},
        {QuestionID: "q3", SelectedOption: "c"},
        {QuestionID: "q4", SelectedOption: "d"},
        {QuestionID: "q5", SelectedOption: "x"},
    }
    
    score := service.CalculateScore(studentAnswers, correctAnswers)
    assert.Equal(t, 60, score)
}
```

---

**Criterio de Aceptación 4: Persistencia Atómica (ACID)**

✅ **GIVEN** un intento válido  
✅ **WHEN** ocurre error al guardar respuestas individuales (después de guardar intento)  
✅ **THEN** el sistema hace ROLLBACK completo (no se guarda nada)

**Validación:**
```go
func TestAttemptRepository_TransactionRollback(t *testing.T) {
    db, mock := setupMockDB(t)
    repo := NewAttemptRepository(db)
    
    // Simular error en segunda inserción
    mock.ExpectBegin()
    mock.ExpectExec("INSERT INTO assessment_attempt").WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectExec("INSERT INTO assessment_attempt_answer").WillReturnError(errors.New("db error"))
    mock.ExpectRollback()
    
    err := repo.CreateAttempt(ctx, attempt)
    assert.Error(t, err)
    
    // Verificar que no quedó basura en BD
    count := db.Model(&Attempt{}).Where("id = ?", attempt.ID).Count()
    assert.Equal(t, 0, count)
}
```

---

**Criterio de Aceptación 5: Feedback Educativo Presente**

✅ **GIVEN** un intento completado  
✅ **WHEN** el sistema calcula resultados  
✅ **THEN** cada respuesta incorrecta tiene un mensaje de feedback educativo

**Validación:**
```json
// Response esperado
{
  "attempt_id": "uuid",
  "score": 60,
  "feedback": [
    {
      "question_id": "q1",
      "is_correct": true,
      "message": "¡Correcto! Un compilador traduce código fuente..."
    },
    {
      "question_id": "q2",
      "is_correct": false,
      "correct_answer": "b",
      "message": "Incorrecto. Revisa la sección 'Tipos de Datos' en el resumen."
    }
  ]
}
```

---

**Criterio de Aceptación 6: Tiempo de Respuesta <2 segundos**

✅ **GIVEN** un intento con 5 preguntas  
✅ **WHEN** se envía POST /attempts  
✅ **THEN** el sistema responde en <2 segundos (p95)

**Validación:**
```bash
# Performance test
for i in {1..100}; do
  time curl -X POST http://localhost:8080/v1/materials/uuid-1/assessment/attempts \
    -H "Authorization: Bearer token-$i" \
    -H "Content-Type: application/json" \
    -d @valid-attempt.json
done | awk '{print $2}' | sort -n | tail -5

# Verificar: p95 < 2.0 segundos
```

---

**Criterio de Aceptación 7: Inmutabilidad de Intentos**

✅ **GIVEN** un intento ya creado  
✅ **WHEN** intenta modificar el intento (PUT /attempts/:id)  
✅ **THEN** el sistema retorna 405 Method Not Allowed

**Validación:**
```bash
curl -X PUT http://localhost:8080/v1/attempts/uuid-1 \
  -H "Authorization: Bearer valid-token" \
  -d '{"score": 100}'

# Esperado: 405 Method Not Allowed
```

**Validación DB:**
```sql
-- No debe existir UPDATE en assessment_attempt
-- Solo INSERT permitido
```

---

### 2.3 MÓDULO: Consulta de Resultados

#### AC-MOD-003: GET /v1/attempts/:id/results

**Criterio de Aceptación 1: Solo Propietario Puede Acceder**

✅ **GIVEN** un intento creado por estudiante A  
✅ **WHEN** estudiante B intenta acceder a GET /attempts/{id}  
✅ **THEN** el sistema retorna 403 Forbidden

**Validación:**
```go
func TestAttemptHandler_OnlyOwnerCanAccess(t *testing.T) {
    // Crear intento de studentA
    attemptID := createAttempt(t, studentAToken, materialID)
    
    // Intentar acceder con studentB
    resp := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/v1/attempts/"+attemptID+"/results", nil)
    req.Header.Set("Authorization", "Bearer "+studentBToken)
    
    router.ServeHTTP(resp, req)
    
    assert.Equal(t, 403, resp.Code)
    assert.Contains(t, resp.Body.String(), "forbidden")
}
```

---

**Criterio de Aceptación 2: Datos Completos en Response**

✅ **GIVEN** un intento existente  
✅ **WHEN** el propietario solicita GET /attempts/{id}/results  
✅ **THEN** el response contiene:
- Datos del intento (score, time_spent, completed_at)
- Todas las respuestas individuales
- Feedback educativo por pregunta
- Datos del material asociado

**Validación:**
```json
// Response esperado
{
  "attempt_id": "uuid",
  "material_id": "uuid",
  "material_title": "Introducción a Pascal",
  "score": 80,
  "max_score": 100,
  "correct_answers": 4,
  "total_questions": 5,
  "time_spent_seconds": 420,
  "completed_at": "2025-11-14T10:30:00Z",
  "passed": true,
  "feedback": [
    {
      "question_id": "q1",
      "question_text": "¿Qué es un compilador?",
      "selected_option": "a",
      "correct_answer": "a",
      "is_correct": true,
      "message": "¡Correcto! ..."
    }
  ]
}
```

---

#### AC-MOD-004: GET /v1/users/me/attempts

**Criterio de Aceptación 1: Historial Completo del Usuario**

✅ **GIVEN** un estudiante con 3 intentos previos  
✅ **WHEN** solicita GET /v1/users/me/attempts  
✅ **THEN** el sistema retorna los 3 intentos ordenados por fecha descendente

**Validación:**
```go
func TestAttemptHandler_HistoryOrderedByDate(t *testing.T) {
    // Crear 3 intentos en orden temporal
    attempt1 := createAttempt(t, token, material1, time.Now().Add(-2*time.Hour))
    attempt2 := createAttempt(t, token, material2, time.Now().Add(-1*time.Hour))
    attempt3 := createAttempt(t, token, material1, time.Now())
    
    resp := getAttemptHistory(t, token)
    
    assert.Len(t, resp.Attempts, 3)
    assert.Equal(t, attempt3.ID, resp.Attempts[0].ID) // Más reciente primero
    assert.Equal(t, attempt2.ID, resp.Attempts[1].ID)
    assert.Equal(t, attempt1.ID, resp.Attempts[2].ID)
}
```

---

**Criterio de Aceptación 2: Paginación Funcional**

✅ **GIVEN** un estudiante con 100 intentos  
✅ **WHEN** solicita GET /v1/users/me/attempts?limit=10&offset=20  
✅ **THEN** el sistema retorna intentos 21-30

**Validación:**
```bash
curl -X GET "http://localhost:8080/v1/users/me/attempts?limit=10&offset=20" \
  -H "Authorization: Bearer valid-token"

# Esperado: 200 OK
# {
#   "attempts": [ ... ], // 10 items
#   "total_count": 100,
#   "page": 3,
#   "limit": 10
# }
```

---

## 3. CRITERIOS DE INTEGRACIÓN

### 3.1 Integración con MongoDB

**Criterio de Aceptación: Lectura Exitosa de Preguntas**

✅ **GIVEN** MongoDB está disponible  
✅ **AND** existe colección `material_assessment` con documento válido  
✅ **WHEN** el sistema consulta preguntas  
✅ **THEN** obtiene estructura correcta sin errores

**Validación:**
```go
func TestQuestionRepository_FindByMaterialID(t *testing.T) {
    // Setup MongoDB testcontainer
    ctx := context.Background()
    mongoContainer, db := setupMongoTestContainer(t)
    defer mongoContainer.Terminate(ctx)
    
    // Seed data
    seedMongoAssessment(t, db, "material-uuid-1")
    
    // Test
    repo := NewMongoQuestionRepository(db)
    questions, err := repo.FindByMaterialID(ctx, "material-uuid-1")
    
    assert.NoError(t, err)
    assert.Len(t, questions, 5)
    assert.NotEmpty(t, questions[0].Text)
    assert.NotEmpty(t, questions[0].CorrectAnswer)
}
```

---

### 3.2 Integración con PostgreSQL

**Criterio de Aceptación: Constraints de Integridad**

✅ **GIVEN** un intento con attempt_id válido  
✅ **WHEN** intenta insertar respuesta con question_id duplicado  
✅ **THEN** PostgreSQL rechaza con error de PRIMARY KEY

**Validación:**
```sql
-- Setup
INSERT INTO assessment_attempt (id, assessment_id, student_id, score, ...) 
VALUES ('attempt-1', 'assessment-1', 'student-1', 80, ...);

-- Test: Intentar duplicar
INSERT INTO assessment_attempt_answer (attempt_id, question_id, selected_option, is_correct)
VALUES ('attempt-1', 'q1', 'a', true);

-- Esto debe fallar
INSERT INTO assessment_attempt_answer (attempt_id, question_id, selected_option, is_correct)
VALUES ('attempt-1', 'q1', 'b', false);

-- Esperado: ERROR: duplicate key value violates unique constraint
```

---

## 4. CRITERIOS DE PERFORMANCE

### 4.1 Latencia de Endpoints

| Endpoint | p50 | p95 | p99 | Validación |
|----------|-----|-----|-----|------------|
| GET /assessment | <100ms | <200ms | <300ms | Load test con 1000 requests |
| POST /attempts | <500ms | <1500ms | <2000ms | Load test con 500 requests |
| GET /attempts/:id | <100ms | <200ms | <300ms | Load test con 1000 requests |

**Validación:**
```bash
# Usar Apache Bench o k6
ab -n 1000 -c 50 \
  -H "Authorization: Bearer token" \
  http://localhost:8080/v1/materials/uuid-1/assessment

# Verificar percentiles en output
```

---

### 4.2 Throughput

**Criterio de Aceptación: 100 req/s sostenido**

✅ **GIVEN** API corriendo en producción  
✅ **WHEN** recibe 100 requests por segundo durante 1 minuto  
✅ **THEN** tasa de error <1% y latencia p95 <500ms

**Validación:**
```bash
# k6 load test
k6 run --vus 100 --duration 60s load-test.js

# Verificar:
# - http_req_failed < 1%
# - http_req_duration(p95) < 500ms
```

---

## 5. CRITERIOS DE SEGURIDAD

### 5.1 Autenticación JWT

**Criterio de Aceptación: Token Inválido Rechazado**

✅ **GIVEN** un token JWT expirado  
✅ **WHEN** intenta acceder a cualquier endpoint protegido  
✅ **THEN** el sistema retorna 401 Unauthorized

**Validación:**
```bash
# Token expirado (exp = 2020-01-01)
curl -X GET http://localhost:8080/v1/materials/uuid-1/assessment \
  -H "Authorization: Bearer expired.token.here"

# Esperado: 401
```

---

### 5.2 Sanitización de Inputs

**Criterio de Aceptación: SQL Injection Prevención**

✅ **GIVEN** un atacante intenta SQL injection  
✅ **WHEN** envía material_id malicioso  
✅ **THEN** el sistema rechaza con error 400

**Validación:**
```bash
curl -X GET "http://localhost:8080/v1/materials/'; DROP TABLE users; --/assessment" \
  -H "Authorization: Bearer valid-token"

# Esperado: 400 Bad Request (UUID inválido)
# ⚠️ NUNCA: Error SQL o tabla eliminada
```

---

## 6. CRITERIOS DE OBSERVABILIDAD

### 6.1 Logging

**Criterio de Aceptación: Logs Estructurados**

✅ **GIVEN** cualquier operación en el sistema  
✅ **WHEN** se ejecuta  
✅ **THEN** se genera log estructurado JSON con campos: timestamp, level, message, context

**Validación:**
```json
// Log esperado
{
  "timestamp": "2025-11-14T10:30:00Z",
  "level": "info",
  "message": "Assessment attempt created",
  "student_id": "uuid",
  "attempt_id": "uuid",
  "score": 80,
  "duration_ms": 1234
}
```

---

### 6.2 Health Check

**Criterio de Aceptación: Health Endpoint Funcional**

✅ **GIVEN** todos los servicios funcionando  
✅ **WHEN** consulta GET /health  
✅ **THEN** retorna 200 con status "healthy"

✅ **GIVEN** PostgreSQL caído  
✅ **WHEN** consulta GET /health  
✅ **THEN** retorna 503 con status "unhealthy"

**Validación:**
```bash
# Caso exitoso
curl http://localhost:8080/health
# {
#   "status": "healthy",
#   "postgres": true,
#   "mongodb": true
# }

# Caso con error
docker stop postgres-container
curl http://localhost:8080/health
# {
#   "status": "unhealthy",
#   "postgres": false,
#   "mongodb": true
# }
```

---

## 7. CRITERIOS DE TESTING

### 7.1 Coverage Mínimo

**Criterio de Aceptación: >80% Global**

✅ **GIVEN** todos los tests ejecutados  
✅ **WHEN** se calcula coverage  
✅ **THEN** coverage global es >80%

**Validación:**
```bash
go test ./... -cover -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total

# Esperado: total: (statements)    82.5%
```

---

### 7.2 Tests de Integración con Testcontainers

**Criterio de Aceptación: Tests Autónomos**

✅ **GIVEN** máquina limpia sin PostgreSQL/MongoDB instalados  
✅ **WHEN** ejecuta tests de integración  
✅ **THEN** tests pasan usando testcontainers

**Validación:**
```bash
# En máquina limpia (solo Docker instalado)
go test ./tests/integration -v

# Esperado:
# - Testcontainers levanta PostgreSQL/MongoDB
# - Tests pasan
# - Contenedores se limpian automáticamente
```

---

## 8. CRITERIOS DE DEPLOYMENT

### 8.1 Migraciones de Base de Datos

**Criterio de Aceptación: Migraciones Idempotentes**

✅ **GIVEN** schema ya aplicado  
✅ **WHEN** ejecuta migraciones nuevamente  
✅ **THEN** no hay errores (operaciones IF NOT EXISTS)

**Validación:**
```bash
# Primera ejecución
psql -U postgres -d edugo < scripts/postgresql/06_assessments.sql
# Esperado: Success

# Segunda ejecución (idempotente)
psql -U postgres -d edugo < scripts/postgresql/06_assessments.sql
# Esperado: Success (sin errores de "ya existe")
```

---

### 8.2 Health Check en CI/CD

**Criterio de Aceptación: Pipeline Valida Health**

✅ **GIVEN** aplicación desplegada en staging  
✅ **WHEN** pipeline de CI/CD ejecuta  
✅ **THEN** verifica health check antes de promover a producción

**Validación (GitHub Actions):**
```yaml
# .github/workflows/deploy.yml
- name: Health Check
  run: |
    response=$(curl -s -o /dev/null -w "%{http_code}" http://staging/health)
    if [ $response -ne 200 ]; then
      echo "Health check failed"
      exit 1
    fi
```

---

## 9. RESUMEN DE CRITERIOS SMART

### 9.1 Específicos (Specific)
- ✅ Cada criterio define QUÉ debe cumplirse exactamente
- ✅ Sin ambigüedades: "Response <2 seg" vs "Response rápido"

### 9.2 Medibles (Measurable)
- ✅ Métricas cuantificables: >80% coverage, <2 seg latencia
- ✅ Tests automatizables para validación

### 9.3 Alcanzables (Achievable)
- ✅ Basados en arquitectura existente de api-mobile
- ✅ Tecnologías ya probadas (Go, PostgreSQL, MongoDB)

### 9.4 Relevantes (Relevant)
- ✅ Alineados con objetivos de negocio (OB-01 a OB-04)
- ✅ Priorizados según impacto (MUST > SHOULD > COULD)

### 9.5 Temporales (Time-bound)
- ✅ MVP completo en 2 semanas (14 días)
- ✅ Post-MVP en fases posteriores

---

## 10. CHECKLIST FINAL DE ACEPTACIÓN

Antes de considerar el proyecto DONE, verificar:

### Funcionalidad
- [ ] Todos los endpoints REST implementados
- [ ] Validaciones de seguridad en servidor
- [ ] Respuestas correctas nunca expuestas al cliente
- [ ] Cálculo de puntaje correcto
- [ ] Feedback educativo presente

### Base de Datos
- [ ] Schema PostgreSQL ejecutable
- [ ] Constraints de integridad funcionando
- [ ] Índices optimizados creados
- [ ] Migraciones idempotentes

### Tests
- [ ] Coverage global >80%
- [ ] Tests unitarios pasando
- [ ] Tests de integración con testcontainers pasando
- [ ] Tests E2E del flujo completo pasando

### Performance
- [ ] Latencia p95 <2 seg para POST /attempts
- [ ] Latencia p95 <500ms para GET endpoints
- [ ] Throughput >50 req/s sostenido

### Seguridad
- [ ] Autenticación JWT funcionando
- [ ] Autorización (solo propietario accede a intentos)
- [ ] SQL injection prevenido
- [ ] Sanitización de inputs

### Observabilidad
- [ ] Logs estructurados JSON
- [ ] Health check endpoint funcional
- [ ] Métricas Prometheus (Post-MVP)

### Documentación
- [ ] Swagger annotations completas
- [ ] README actualizado
- [ ] Comentarios godoc en funciones públicas

### CI/CD
- [ ] Pipeline GitHub Actions verde
- [ ] Linting pasando (golangci-lint)
- [ ] Build exitoso
- [ ] Tests automáticos en CI

---

**Generado con:** Claude Code  
**Total Criterios:** 47 criterios medibles  
**Última actualización:** 2025-11-14
