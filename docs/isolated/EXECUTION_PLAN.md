# EXECUTION PLAN - API Mobile

## Información del Proyecto

**Proyecto:** EduGo API Mobile  
**Objetivo:** API REST para gestión de evaluaciones académicas  
**Duración:** 6 Sprints (12 semanas)  
**Equipo:** Backend engineers + DevOps  
**Repositorio:** https://github.com/EduGoGroup/edugo-api-mobile

---

## Fase 1: Setup e Inicialización (Sprint 1)

### 1.1 Configuración del Proyecto
- [ ] Clonar repo desde GitHub
- [ ] Instalar Go 1.21+
- [ ] Ejecutar `go mod download && go mod tidy`
- [ ] Configurar `.env.local` para desarrollo
- [ ] Verificar compilación: `go build ./cmd/api-mobile`

### 1.2 Configuración de Base de Datos
- [ ] Crear base de datos PostgreSQL `edugo_mobile`
- [ ] Ejecutar migraciones GORM: `go run ./cmd/migrate`
- [ ] Verificar tablas creadas:
  ```bash
  psql -U edugo_user -d edugo_mobile -c "\dt"
  ```
- [ ] Crear índices de performance
- [ ] Insertar datos de prueba

### 1.3 Configuración de MongoDB
- [ ] Crear base de datos MongoDB `edugo_assessments`
- [ ] Crear colecciones (`evaluation_results`, `evaluation_audit`)
- [ ] Crear índices en MongoDB
- [ ] Verificar conexión desde aplicación

### 1.4 Configuración de RabbitMQ
- [ ] Levantar RabbitMQ (local o Docker)
- [ ] Acceder a Management UI (localhost:15672)
- [ ] Crear exchange `assessment.requests`
- [ ] Crear queue `worker.assessment.requests`
- [ ] Crear exchange `assessment.responses`
- [ ] Crear queue `api-mobile.assessment.responses`
- [ ] Bindear exchanges a queues

### 1.5 Configuración de Docker
- [ ] Crear `docker/Dockerfile`
- [ ] Crear `docker-compose.yml`
- [ ] Verificar build: `docker build -t edugo/api-mobile:latest .`
- [ ] Levantar stack: `docker-compose up`
- [ ] Verificar health check: `curl http://localhost:8080/api/v1/health`

### 1.6 Integración con SHARED
- [ ] Agregar a `go.mod`: `require github.com/EduGoGroup/edugo-shared v1.3.0`
- [ ] Descargar: `go mod download`
- [ ] Importar módulos (logger, database, auth, messaging)
- [ ] Probar logger: logs en stdout
- [ ] Probar database: conexión a PostgreSQL
- [ ] Probar auth: validar tokens JWT
- [ ] Probar messaging: conexión a RabbitMQ

### 1.7 Estructura Inicial de Código
```
api-mobile/
├── cmd/
│   └── api-mobile/
│       └── main.go              # Entry point
├── internal/
│   ├── handlers/
│   │   └── evaluation_handler.go
│   ├── services/
│   │   └── evaluation_service.go
│   ├── repositories/
│   │   └── evaluation_repository.go
│   ├── models/
│   │   └── evaluation.go
│   ├── middleware/
│   │   └── auth.go
│   └── config/
│       └── config.go
├── migrations/
│   └── init.sql
├── docker/
│   └── Dockerfile
├── go.mod
├── go.sum
└── docker-compose.yml
```

### 1.8 Pipeline CI/CD Básico
- [ ] Crear GitHub Actions workflow para tests
- [ ] Crear workflow para build Docker
- [ ] Crear workflow para push a registry
- [ ] Configurar linting (golangci-lint)
- [ ] Configurar tests automáticos (go test)

**Checklist de Completación Sprint 1:**
- [ ] Proyecto compila sin errores
- [ ] Base de datos lista
- [ ] Docker funciona
- [ ] SHARED integrado
- [ ] Health check responde 200 OK
- [ ] Logger genera logs JSON
- [ ] Tests pasan (go test ./...)

---

## Fase 2: Evaluaciones CRUD (Sprint 2)

### 2.1 Modelo de Evaluación
```go
type Evaluation struct {
    ID          int64
    MaterialID  *int64
    Title       string
    Description string
    Type        string    // 'manual', 'generated'
    Status      string    // 'draft', 'published', 'closed'
    PassingScore int
    CreatedBy   int64
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   *time.Time
}
```

### 2.2 Endpoints a Implementar

#### POST /api/v1/evaluations
```
Request:
{
  "title": "Quiz Matemáticas",
  "description": "Evaluación de cálculo",
  "type": "manual",
  "passing_score": 60
}

Response: 201 Created
{
  "id": 1,
  "title": "Quiz Matemáticas",
  "status": "draft",
  "created_at": "2025-11-15T10:30:00Z",
  "created_by": 42
}
```

#### GET /api/v1/evaluations
```
Query params:
- page=1&limit=10
- status=draft|published|closed
- created_by=42

Response: 200 OK
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 50
  }
}
```

#### GET /api/v1/evaluations/:id
```
Response: 200 OK
{
  "id": 1,
  "title": "Quiz Matemáticas",
  "description": "...",
  "type": "manual",
  "status": "draft",
  "passing_score": 60,
  "created_by": 42,
  "created_at": "2025-11-15T10:30:00Z",
  "updated_at": "2025-11-15T10:30:00Z"
}

No encontrado: 404 Not Found
{
  "error": "Evaluation not found",
  "code": "EVALUATION_NOT_FOUND"
}
```

#### PUT /api/v1/evaluations/:id
```
Request:
{
  "title": "Quiz Matemáticas Avanzadas",
  "passing_score": 70
}

Response: 200 OK
{
  "id": 1,
  "title": "Quiz Matemáticas Avanzadas",
  "passing_score": 70,
  "updated_at": "2025-11-15T10:45:00Z"
}
```

#### DELETE /api/v1/evaluations/:id
```
Response: 204 No Content
(soft delete)
```

### 2.3 Service Layer
```go
type EvaluationService interface {
    Create(ctx context.Context, req CreateEvaluationRequest) (*Evaluation, error)
    GetByID(ctx context.Context, id int64) (*Evaluation, error)
    List(ctx context.Context, filters EvaluationFilters) ([]*Evaluation, error)
    Update(ctx context.Context, id int64, req UpdateEvaluationRequest) (*Evaluation, error)
    Delete(ctx context.Context, id int64) error
}
```

### 2.4 Repository Layer (GORM)
```go
type EvaluationRepository interface {
    Create(ctx context.Context, evaluation *Evaluation) error
    GetByID(ctx context.Context, id int64) (*Evaluation, error)
    List(ctx context.Context, filters EvaluationFilters) ([]*Evaluation, int64, error)
    Update(ctx context.Context, evaluation *Evaluation) error
    Delete(ctx context.Context, id int64) error
}
```

### 2.5 Tests Unitarios
- [ ] Tests para service (mocks de repository)
- [ ] Tests para repository (base de datos real)
- [ ] Tests para handlers (mocks de service)
- [ ] Cobertura mínima: 80%

### 2.6 Tests de Integración
- [ ] Test: Crear evaluación y verificar en BD
- [ ] Test: Listar evaluaciones con filtros
- [ ] Test: Actualizar evaluación existente
- [ ] Test: Borrar evaluación (soft delete)
- [ ] Test: Validaciones de request

**Checklist de Completación Sprint 2:**
- [ ] Endpoints CRUD funcionales
- [ ] Tests unitarios pasan
- [ ] Tests de integración pasan
- [ ] Swagger documentado
- [ ] Cobertura >= 80%

---

## Fase 3: Preguntas y Opciones (Sprint 3)

### 3.1 Modelos

```go
type Question struct {
    ID           int64
    EvaluationID int64
    Type         string // 'multiple_choice', 'true_false', 'short_answer'
    Text         string
    Position     int
    Points       int
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type QuestionOption struct {
    ID         int64
    QuestionID int64
    Text       string
    IsCorrect  bool
    Position   int
}
```

### 3.2 Endpoints

#### POST /api/v1/evaluations/:id/preguntas
```
Request:
{
  "type": "multiple_choice",
  "text": "¿Cuál es la capital de España?",
  "points": 5,
  "options": [
    {"text": "Madrid", "is_correct": true},
    {"text": "Barcelona", "is_correct": false},
    {"text": "Valencia", "is_correct": false}
  ]
}

Response: 201 Created
{
  "id": 101,
  "evaluation_id": 1,
  "type": "multiple_choice",
  "text": "¿Cuál es la capital de España?",
  "position": 1,
  "points": 5,
  "options": [...]
}
```

#### GET /api/v1/evaluaciones/:id/preguntas

#### PUT /api/v1/preguntas/:id

#### DELETE /api/v1/preguntas/:id

#### POST /api/v1/evaluaciones/:id/reorder
```
Request:
{
  "questions": [
    {"id": 101, "position": 1},
    {"id": 102, "position": 2},
    {"id": 103, "position": 3}
  ]
}

Response: 200 OK
{
  "updated": 3
}
```

### 3.3 Validaciones
- Pregunta no puede estar sin evaluación
- Evaluación debe estar en estado 'draft' para agregar preguntas
- Mínimo 2 opciones para multiple choice
- Al menos 1 opción correcta
- Puntos >= 1
- Text no vacío

### 3.4 Tests
- [ ] Tests CRUD de preguntas
- [ ] Tests de validaciones
- [ ] Tests de reordenamiento
- [ ] Tests de relación evaluación-preguntas

**Checklist de Completación Sprint 3:**
- [ ] Preguntas CRUD funcional
- [ ] Validaciones implementadas
- [ ] Tests >= 80%
- [ ] Swagger actualizado

---

## Fase 4: Asignaciones a Estudiantes (Sprint 4)

### 4.1 Modelo

```go
type EvaluationAssignment struct {
    ID           int64
    EvaluationID int64
    StudentID    int64
    AssignedAt   time.Time
    DueDate      *time.Time
    Status       string // 'pending', 'in_progress', 'submitted', 'graded'
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

### 4.2 Endpoints

#### POST /api/v1/evaluaciones/:id/assign
```
Request:
{
  "student_ids": [1, 2, 3, 4, 5],
  "due_date": "2025-12-01T23:59:59Z"
}

Response: 201 Created
{
  "assigned": 5,
  "assignments": [
    {
      "id": 1001,
      "evaluation_id": 1,
      "student_id": 1,
      "status": "pending",
      "due_date": "2025-12-01T23:59:59Z"
    }
  ]
}
```

#### GET /api/v1/evaluaciones/:id/assignments
```
Response: 200 OK
{
  "data": [
    {
      "id": 1001,
      "student_id": 1,
      "student_name": "Juan García",
      "status": "pending",
      "due_date": "2025-12-01T23:59:59Z",
      "submission_date": null
    }
  ],
  "summary": {
    "total": 10,
    "pending": 5,
    "in_progress": 2,
    "submitted": 2,
    "graded": 1
  }
}
```

#### PUT /api/v1/assignments/:id
```
Request:
{
  "due_date": "2025-12-15T23:59:59Z"
}

Response: 200 OK
```

#### POST /api/v1/assignments/:id/start
```
Request: (sin body)

Response: 200 OK
{
  "id": 1001,
  "status": "in_progress",
  "started_at": "2025-11-15T14:30:00Z"
}
```

### 4.3 Flujo de Evaluación
```
1. Evaluación en estado 'draft'
2. Docente publica: POST /evaluaciones/:id/publish
   → Status cambia a 'published'
3. Docente asigna: POST /evaluaciones/:id/assign
   → Crea assignments con status 'pending'
4. Estudiante inicia: POST /assignments/:id/start
   → Status cambia a 'in_progress'
5. Estudiante responde: POST /evaluaciones/:id/submit
   → Status cambia a 'submitted'
6. API calcula: Resultados persistidos
   → Status cambia a 'graded'
```

### 4.4 Validaciones
- Solo se puede asignar si evaluación está publicada
- No se puede reasignar a mismo estudiante
- Due date debe ser en el futuro (opcional)
- Solo se puede iniciar si status es 'pending'

### 4.5 Tests
- [ ] Tests de asignación
- [ ] Tests de cambio de estado
- [ ] Tests de validaciones

**Checklist de Completación Sprint 4:**
- [ ] Asignaciones funcionales
- [ ] Estados correctos
- [ ] Tests >= 80%

---

## Fase 5: Respuestas y Validación (Sprint 5)

### 5.1 Modelos

```go
type AnswerDraft struct {
    ID           int64
    AssignmentID int64
    QuestionID   int64
    Answer       string
    SavedAt      time.Time
}

type SubmissionRequest struct {
    AssignmentID int64
    Answers      map[int64]string // question_id -> answer
}
```

### 5.2 Endpoints

#### GET /api/v1/evaluaciones/:id/draft
```
Obtener respuestas guardadas en borrador

Response: 200 OK
{
  "assignment_id": 1001,
  "evaluation_id": 1,
  "answers": {
    "101": "Option A",
    "102": "True",
    "103": "El proceso es..."
  }
}
```

#### POST /api/v1/evaluaciones/:id/submit
```
Request:
{
  "assignment_id": 1001,
  "answers": {
    "101": "Option B",
    "102": "False",
    "103": "La respuesta es..."
  }
}

Response: 200 OK
{
  "submission_id": "uuid-12345",
  "status": "submitted",
  "submitted_at": "2025-11-15T15:00:00Z",
  "will_be_graded_at": "2025-11-15T15:01:00Z"
}
```

#### POST /api/v1/evaluaciones/:id/save-draft
```
Guardar respuestas sin enviar (respuestas parciales)

Request:
{
  "assignment_id": 1001,
  "answers": {
    "101": "Option A"
  }
}

Response: 200 OK
{
  "saved": true,
  "timestamp": "2025-11-15T14:50:00Z"
}
```

### 5.3 Lógica de Validación

```go
type ValidationEngine struct {
    // Validar que respuesta corresponde a pregunta
    // Validar tipos de respuesta (string para short_answer, option para MC)
    // Validar que opción existe
    // Validar que no hay respuestas duplicadas
}

func (ve *ValidationEngine) ValidateAnswers(eval *Evaluation, answers map[int64]string) error {
    // 1. Verificar que assignment existe
    // 2. Verificar que assignment está en_progress
    // 3. Para cada respuesta:
    //    - Verificar que pregunta existe en evaluación
    //    - Verificar que respuesta es válida para tipo de pregunta
    // 4. Retornar error si alguna validación falla
}
```

### 5.4 Cálculo de Puntuación

```go
type ScoringEngine struct {
    // Calcular puntos por pregunta
    // Calcular total score
    // Calcular porcentaje
    // Determinar si aprobó (>= passing_score)
}

func (se *ScoringEngine) CalculateScore(eval *Evaluation, answers map[int64]string) (*Score, error) {
    score := &Score{
        TotalScore:  0,
        MaxScore:    0,
        Percentage:  0,
        AnswerDetails: make([]AnswerDetail, 0),
    }
    
    // 1. Iterar cada pregunta
    // 2. Obtener respuesta del estudiante
    // 3. Comparar con respuesta correcta
    // 4. Sumar puntos si es correcta
    // 5. Calcular total y porcentaje
    // 6. Determinar si aprobó
    
    return score, nil
}
```

### 5.5 Tests
- [ ] Tests de validación de respuestas
- [ ] Tests de cálculo de puntuación
- [ ] Tests de casos límite (sin respuestas, todas correctas, todas incorrectas)
- [ ] Tests de guardado de draft
- [ ] Tests de sumisión completa

**Checklist de Completación Sprint 5:**
- [ ] Endpoints de respuestas funcionales
- [ ] Validación robusta
- [ ] Cálculo de puntuación correcto
- [ ] Tests >= 80%

---

## Fase 6: Resultados, Reportes e Integración con Worker (Sprint 6)

### 6.1 Persistencia de Resultados en MongoDB

```go
type EvaluationResult struct {
    ID           primitive.ObjectID
    EvaluationID int64
    AssignmentID int64
    StudentID    int64
    Answers      []AnswerDetail
    TotalScore   int
    MaxScore     int
    Percentage   float64
    Status       string // "graded"
    SubmittedAt  time.Time
    Feedback     string
    Metadata     map[string]interface{}
}

func (es *EvaluationService) SaveResults(ctx context.Context, result *EvaluationResult) error {
    collection := db.GetMongoCollection("evaluation_results")
    _, err := collection.InsertOne(ctx, result)
    return err
}
```

### 6.2 Endpoints de Resultados

#### GET /api/v1/evaluaciones/:id/results
```
Obtener resultados consolidados de una evaluación

Response: 200 OK
{
  "evaluation_id": 1,
  "total_students": 10,
  "submitted": 8,
  "pending": 2,
  "average_score": 75.5,
  "highest_score": 95,
  "lowest_score": 45,
  "results": [
    {
      "student_id": 1,
      "student_name": "Juan García",
      "score": 80,
      "percentage": 80,
      "status": "graded",
      "submitted_at": "2025-11-15T15:00:00Z"
    }
  ]
}
```

#### GET /api/v1/assignments/:id/answer-detail
```
Detalle completo de respuesta de un estudiante

Response: 200 OK
{
  "assignment_id": 1001,
  "student_id": 1,
  "evaluation_title": "Quiz Matemáticas",
  "answers": [
    {
      "question_id": 101,
      "question_text": "¿Cuál es la capital de España?",
      "question_type": "multiple_choice",
      "student_answer": "Option B",
      "correct_answer": "Option A",
      "is_correct": false,
      "points_earned": 0,
      "points_available": 5
    }
  ],
  "total_score": 80,
  "max_score": 100,
  "percentage": 80
}
```

#### GET /api/v1/evaluaciones/:id/results/export
```
Exportar resultados en CSV

Response: 200 OK (CSV file)
student_id,student_name,score,percentage,status,submitted_at
1,Juan García,80,80,graded,2025-11-15T15:00:00Z
2,María López,85,85,graded,2025-11-15T15:05:00Z
...
```

### 6.3 Integración con Worker (Generación Automática de Quizzes)

#### POST /api/v1/evaluaciones/material/:id/generate-quiz
```
Request:
{
  "num_questions": 10,
  "difficulty": "medium",
  "language": "es"
}

Response: 202 Accepted
{
  "request_id": "uuid-12345",
  "status": "processing",
  "estimated_time": "30s"
}
```

#### GET /api/v1/requests/:id/status
```
Consultar estado de generación

Response: 200 OK
{
  "request_id": "uuid-12345",
  "status": "processing|completed|error",
  "progress": 50,
  "evaluation_id": null, // Si está completado
  "error_message": null  // Si hubo error
}
```

#### Consumidor de Respuestas del Worker
```go
func (es *EvaluationService) ConsumeWorkerResponses(ctx context.Context) error {
    subscriber := messaging.NewSubscriber()
    
    messages := subscriber.Subscribe(
        "assessment.responses",
        "api-mobile.assessment.responses",
    )
    
    for {
        select {
        case msg := <-messages:
            var response WorkerResponse
            json.Unmarshal(msg.Body, &response)
            
            if response.Status == "error" {
                es.logger.Error("Worker error", map[string]interface{}{
                    "request_id": response.RequestID,
                    "error": response.ErrorMessage,
                })
                es.UpdateGenerationRequest(response.RequestID, "failed")
            } else {
                // Guardar preguntas
                for _, q := range response.Questions {
                    es.repo.CreateQuestion(ctx, q)
                }
                
                // Actualizar evaluation status
                es.repo.Update(ctx, &Evaluation{
                    ID: response.EvaluationID,
                    Status: "published",
                })
                
                // Marcar request como completado
                es.UpdateGenerationRequest(response.RequestID, "completed")
            }
            
            msg.Ack(false)
            
        case <-ctx.Done():
            return nil
        }
    }
}
```

### 6.4 Publicador de Requests a Worker
```go
func (es *EvaluationService) RequestQuizGeneration(ctx context.Context, req GenerateQuizRequest) error {
    // 1. Crear evaluation en estado 'generating'
    eval := &Evaluation{
        MaterialID: req.MaterialID,
        Type: "generated",
        Status: "generating",
    }
    es.repo.Create(ctx, eval)
    
    // 2. Guardar request de generación
    genReq := &GenerationRequest{
        ID: uuid.New().String(),
        EvaluationID: eval.ID,
        Status: "pending",
    }
    es.repo.SaveGenerationRequest(genReq)
    
    // 3. Publicar a RabbitMQ
    publisher := messaging.NewPublisher()
    payload, _ := json.Marshal(map[string]interface{}{
        "request_id": genReq.ID,
        "evaluation_id": eval.ID,
        "material_id": req.MaterialID,
        "config": req,
    })
    
    publisher.Publish(
        "assessment.requests",
        "worker.assessment.requests",
        payload,
    )
    
    // 4. Retornar request_id para polling
    return nil
}
```

### 6.5 Event Publishing
```go
func (es *EvaluationService) PublishEvaluationEvent(ctx context.Context, event *EvaluationEvent) error {
    publisher := messaging.NewPublisher()
    
    payload, _ := json.Marshal(event)
    
    return publisher.Publish(
        "evaluation.events",
        "evaluation."+event.Type,
        payload,
    )
}

// Eventos a publicar:
// - evaluation.created
// - evaluation.published
// - question.created
// - assignment.created
// - submission.received
// - results.calculated
```

### 6.6 Reportes Avanzados

#### GET /api/v1/evaluaciones/:id/analytics
```
Response: 200 OK
{
  "evaluation_id": 1,
  "title": "Quiz Matemáticas",
  "statistics": {
    "total_assigned": 10,
    "total_submitted": 8,
    "submission_rate": 0.8,
    "average_score": 75.5,
    "median_score": 78,
    "std_deviation": 12.3,
    "pass_rate": 0.75,
    "score_distribution": {
      "0-20": 0,
      "21-40": 1,
      "41-60": 2,
      "61-80": 3,
      "81-100": 2
    }
  },
  "question_analytics": [
    {
      "question_id": 101,
      "question_text": "¿Cuál es la capital de España?",
      "correct_count": 6,
      "incorrect_count": 2,
      "accuracy_rate": 0.75,
      "avg_time_seconds": 15
    }
  ]
}
```

### 6.7 Optimizaciones
- [ ] Caché de evaluaciones frecuentes (Redis)
- [ ] Índices de MongoDB para queries de resultados
- [ ] Batching de cálculos de resultados
- [ ] Async processing de reportes complejos

### 6.8 Tests
- [ ] Tests de integración con Worker
- [ ] Tests de endpoints de resultados
- [ ] Tests de analytics
- [ ] Tests de export

**Checklist de Completación Sprint 6:**
- [ ] Resultados persistidos en MongoDB
- [ ] Endpoints funcionales
- [ ] Integración Worker completa
- [ ] Reportes funcionando
- [ ] Tests >= 80%
- [ ] Documentación Swagger 100%

---

## Post-Sprints: Producción y Mantenimiento

### Fase 7: Preparación para Producción
- [ ] Load testing (10,000 req/sec)
- [ ] Security scanning (OWASP)
- [ ] Performance profiling
- [ ] Logging en producción
- [ ] Monitoreo con Prometheus
- [ ] Alertas configuradas
- [ ] Runbooks documentados
- [ ] Disaster recovery testing

### Fase 8: Deployment
- [ ] Build imagen Docker multistage
- [ ] Push a registry privado
- [ ] Configurar Kubernetes manifests
- [ ] Deploy a staging
- [ ] Smoke tests
- [ ] Deploy a producción (blue-green)
- [ ] Monitoreo post-deployment
- [ ] Rollback plan

### Fase 9: Mantenimiento Continuo
- [ ] Monitoreo de errores (Sentry)
- [ ] Analítica de performance
- [ ] Updates de dependencias
- [ ] Security patches
- [ ] Escalado automático
- [ ] Backups regulares
- [ ] Auditoría de acceso

---

## Criterios de Aceptación Globales

- [ ] API responde < 100ms en P95
- [ ] Uptime >= 99.9%
- [ ] Error rate < 0.1%
- [ ] Cobertura de tests >= 80%
- [ ] Swagger documentación 100%
- [ ] Security: OWASP Top 10 mitigado
- [ ] Logging centralizado funcionando
- [ ] Monitoreo y alertas configurados
- [ ] CI/CD pipeline automatizado
- [ ] Documentación completa

---

## Dependencias Críticas

```mermaid
Sprint 1 (Setup)
     ↓
Sprint 2 (CRUD Evaluaciones)
     ↓
Sprint 3 (Preguntas) ← Depende Sprint 2
     ↓
Sprint 4 (Asignaciones) ← Depende Sprint 2
     ↓
Sprint 5 (Respuestas) ← Depende Sprint 4
     ↓
Sprint 6 (Resultados + Worker) ← Depende Sprint 5
     ↓
Post-Sprints (Prod + Monitoring)
```

---

## Riesgos y Mitigaciones

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|-------------|--------|-----------|
| Compatibilidad SHARED | Media | Alto | Tests de integración tempranos |
| Performance BD | Media | Alto | Índices desde sprint 1 |
| Integración Worker | Media | Alto | Mocks de RabbitMQ en tests |
| Cambios de esquema | Baja | Alto | Migraciones versionadas |
| Datos de prueba | Baja | Medio | Seeders desde el inicio |

---

## Recursos Requeridos

- 2 Backend engineers (Go)
- 1 DevOps (Docker, Kubernetes)
- 1 QA (Testing)
- Acceso a repos de EduGoGroup
- Infraestructura: PostgreSQL, MongoDB, RabbitMQ
- Registry Docker privado
- Herramientas: Git, GitHub, Docker, Go IDE

---

## Comunicación y Reportes

- **Daily standup:** 15 min (problemas bloqueantes)
- **Sprint review:** Al final de cada sprint (demo funcionalidades)
- **Sprint retrospective:** Lecciones aprendidas
- **Weekly status:** Email con KPIs
  - Tareas completadas
  - Bloqueantes
  - Velocity vs. planned
  - Bugs/Deuda técnica

---

## Diccionario de Estados

| Estado | Significado |
|--------|-----------|
| `draft` | Evaluación en edición, no visible a estudiantes |
| `published` | Evaluación visible, puede ser asignada |
| `closed` | Evaluación cerrada, no se pueden enviar respuestas |
| `pending` | Asignación pendiente de iniciar |
| `in_progress` | Estudiante inició evaluación |
| `submitted` | Estudiante envió respuestas |
| `graded` | Respuestas calificadas, resultados disponibles |

---

**Próxima revisión:** Después de Sprint 1 (semana 2)  
**Última actualización:** 15 de Noviembre, 2025
