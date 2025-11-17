# Tareas del Sprint 04 - Services y API REST

## Objetivo
Implementar la capa de aplicación (services) que orquesta la lógica de negocio y la capa de presentación (handlers/controllers) que expone 4 endpoints REST para el Sistema de Evaluaciones.

---

## Tareas

### TASK-04-001: Implementar AssessmentService
**Tipo:** feature  
**Prioridad:** HIGH  
**Estimación:** 4h

#### Descripción
Service que orquesta operaciones de evaluaciones (obtener assessment, verificar permisos, etc.)

#### Implementación

Archivo: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/application/services/assessment_service.go`

```go
package services

import (
    "context"
    "github.com/google/uuid"
    
    "edugo-api-mobile/internal/domain/entities"
    "edugo-api-mobile/internal/domain/repositories"
)

type AssessmentService struct {
    assessmentRepo repositories.AssessmentRepository
    mongoRepo      repositories.QuestionRepository
}

func NewAssessmentService(
    assessmentRepo repositories.AssessmentRepository,
    mongoRepo repositories.QuestionRepository,
) *AssessmentService {
    return &AssessmentService{
        assessmentRepo: assessmentRepo,
        mongoRepo:      mongoRepo,
    }
}

// GetAssessmentByMaterialID obtiene assessment SIN respuestas correctas
func (s *AssessmentService) GetAssessmentByMaterialID(
    ctx context.Context,
    materialID uuid.UUID,
) (*AssessmentDTO, error) {
    // 1. Buscar assessment en PostgreSQL
    assessment, err := s.assessmentRepo.FindByMaterialID(ctx, materialID)
    if err != nil {
        return nil, err
    }
    if assessment == nil {
        return nil, ErrAssessmentNotFound
    }
    
    // 2. Buscar preguntas en MongoDB (sin respuestas correctas)
    questions, err := s.mongoRepo.FindQuestionsByMaterialID(ctx, materialID.String())
    if err != nil {
        return nil, err
    }
    
    // 3. Convertir a DTO (sanitizado)
    return &AssessmentDTO{
        AssessmentID:      assessment.ID,
        MaterialID:        assessment.MaterialID,
        Title:             assessment.Title,
        TotalQuestions:    assessment.TotalQuestions,
        EstimatedMinutes:  assessment.TimeLimitMinutes,
        Questions:         sanitizeQuestions(questions), // Remover respuestas
    }, nil
}

// VerifyAttemptAllowed verifica si estudiante puede hacer otro intento
func (s *AssessmentService) VerifyAttemptAllowed(
    ctx context.Context,
    studentID, assessmentID uuid.UUID,
) (bool, error) {
    assessment, err := s.assessmentRepo.FindByID(ctx, assessmentID)
    if err != nil {
        return false, err
    }
    if assessment == nil {
        return false, ErrAssessmentNotFound
    }
    
    attemptCount, err := s.attemptRepo.CountByStudentAndAssessment(ctx, studentID, assessmentID)
    if err != nil {
        return false, err
    }
    
    return assessment.CanAttempt(attemptCount), nil
}

// Helper: sanitizar preguntas (remover respuestas correctas)
func sanitizeQuestions(questions []Question) []QuestionDTO {
    result := make([]QuestionDTO, len(questions))
    for i, q := range questions {
        result[i] = QuestionDTO{
            ID:      q.ID,
            Text:    q.Text,
            Type:    q.Type,
            Options: q.Options,
            // ❌ NO incluir correct_answer ni feedback
        }
    }
    return result
}
```

#### Criterios de Aceptación
- [ ] Service orquesta repositorios (no tiene lógica de negocio)
- [ ] NUNCA expone respuestas correctas al cliente
- [ ] Método `sanitizeQuestions()` remueve campos sensibles
- [ ] Usa DTOs para responses (no entities directamente)
- [ ] Context.Context en todos los métodos

#### Comandos de Validación
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
go test ./internal/application/services -v -run TestAssessmentService
```

#### Tiempo Estimado
4 horas

---

### TASK-04-002: Implementar ScoringService
**Tipo:** feature  
**Prioridad:** HIGH  
**Estimación:** 3h

#### Descripción
Service que valida respuestas del cliente contra MongoDB y calcula scores en servidor (NUNCA confiar en cliente).

#### Implementación

Archivo: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/application/services/scoring_service.go`

```go
package services

import (
    "context"
    "edugo-api-mobile/internal/domain/entities"
)

type ScoringService struct {
    mongoRepo repositories.QuestionRepository
}

// ValidateAndScore valida respuestas contra MongoDB y calcula score
// CRÍTICO: Validación servidor-side, NUNCA confiar en score del cliente
func (s *ScoringService) ValidateAndScore(
    ctx context.Context,
    mongoDocID string,
    userAnswers []UserAnswerDTO,
) ([]*entities.Answer, int, error) {
    // 1. Obtener preguntas con respuestas correctas de MongoDB
    assessment, err := s.mongoRepo.FindByDocumentIDWithCorrectAnswers(ctx, mongoDocID)
    if err != nil {
        return nil, 0, err
    }
    
    // 2. Validar cada respuesta contra respuesta correcta
    answers := make([]*entities.Answer, len(userAnswers))
    correctCount := 0
    
    for i, userAnswer := range userAnswers {
        // Buscar pregunta correspondiente
        question := findQuestion(assessment.Questions, userAnswer.QuestionID)
        if question == nil {
            return nil, 0, ErrInvalidQuestionID
        }
        
        // Comparar respuesta del usuario con respuesta correcta
        isCorrect := question.CorrectAnswer == userAnswer.SelectedAnswerID
        
        if isCorrect {
            correctCount++
        }
        
        // Crear entity Answer
        answer, err := entities.NewAnswer(
            uuid.Nil, // Attempt ID se asigna después
            userAnswer.QuestionID,
            userAnswer.SelectedAnswerID,
            isCorrect,
            userAnswer.TimeSpentSeconds,
        )
        if err != nil {
            return nil, 0, err
        }
        
        answers[i] = answer
    }
    
    // 3. Calcular score (servidor-side)
    score := (correctCount * 100) / len(userAnswers)
    
    return answers, score, nil
}
```

#### Criterios de Aceptación
- [ ] Score SIEMPRE calculado en servidor (no acepta score del cliente)
- [ ] Validación contra MongoDB con `correct_answer`
- [ ] Retorna entities.Answer con IsCorrect correctamente calculado
- [ ] Maneja caso de questionID inválido

#### Tiempo Estimado
3 horas

---

### TASK-04-003: Implementar AssessmentHandler (4 endpoints REST)
**Tipo:** feature  
**Prioridad:** HIGH  
**Estimación:** 5h

#### Descripción
Handler que expone 4 endpoints REST usando Gin framework.

#### Implementación

Archivo: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/interfaces/http/handlers/assessment_handler.go`

```go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    
    "edugo-api-mobile/internal/application/services"
)

type AssessmentHandler struct {
    assessmentService *services.AssessmentService
    scoringService    *services.ScoringService
    attemptService    *services.AttemptService
}

// GetAssessment godoc
// @Summary Obtener cuestionario de un material
// @Tags Evaluaciones
// @Security BearerAuth
// @Param id path string true "Material ID (UUID)"
// @Success 200 {object} AssessmentResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/materials/{id}/assessment [get]
func (h *AssessmentHandler) GetAssessment(c *gin.Context) {
    materialID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid material ID"})
        return
    }
    
    assessment, err := h.assessmentService.GetAssessmentByMaterialID(c.Request.Context(), materialID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
        return
    }
    if assessment == nil {
        c.JSON(http.StatusNotFound, ErrorResponse{Error: "Assessment not found"})
        return
    }
    
    c.JSON(http.StatusOK, assessment)
}

// CreateAttempt godoc
// @Summary Crear intento y obtener calificación
// @Tags Evaluaciones
// @Security BearerAuth
// @Param id path string true "Material ID"
// @Param request body CreateAttemptRequest true "Respuestas del estudiante"
// @Success 201 {object} AttemptResponse
// @Router /v1/materials/{id}/assessment/attempts [post]
func (h *AssessmentHandler) CreateAttempt(c *gin.Context) {
    var req CreateAttemptRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
        return
    }
    
    // Obtener student ID del JWT
    studentID := getStudentIDFromContext(c)
    
    // Crear intento
    attempt, err := h.attemptService.CreateAttempt(c.Request.Context(), studentID, req)
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, attempt)
}

// GetAttemptResults godoc
// @Summary Obtener resultados de un intento
// @Tags Evaluaciones
// @Security BearerAuth
// @Param id path string true "Attempt ID"
// @Success 200 {object} AttemptResultsResponse
// @Router /v1/attempts/{id}/results [get]
func (h *AssessmentHandler) GetAttemptResults(c *gin.Context) {
    attemptID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid attempt ID"})
        return
    }
    
    studentID := getStudentIDFromContext(c)
    
    results, err := h.attemptService.GetAttemptResults(c.Request.Context(), attemptID, studentID)
    if err != nil {
        c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, results)
}

// GetUserAttempts godoc
// @Summary Listar intentos del usuario autenticado
// @Tags Evaluaciones
// @Security BearerAuth
// @Param limit query int false "Límite de resultados" default(10)
// @Param offset query int false "Offset para paginación" default(0)
// @Success 200 {object} AttemptsListResponse
// @Router /v1/users/me/attempts [get]
func (h *AssessmentHandler) GetUserAttempts(c *gin.Context) {
    studentID := getStudentIDFromContext(c)
    
    limit := c.DefaultQuery("limit", "10")
    offset := c.DefaultQuery("offset", "0")
    
    attempts, err := h.attemptService.GetUserAttempts(c.Request.Context(), studentID, limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, attempts)
}
```

#### Criterios de Aceptación
- [ ] 4 endpoints implementados
- [ ] Swagger annotations en todos los endpoints
- [ ] Validación de UUIDs en params
- [ ] Error handling consistente
- [ ] DTOs para request/response (no entities)

#### Tiempo Estimado
5 horas

---

### TASK-04-004: Configurar Rutas y Middleware
**Tipo:** feature  
**Prioridad:** HIGH  
**Estimación:** 2h

#### Implementación

Archivo: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/interfaces/http/routes/routes.go`

```go
package routes

import (
    "github.com/gin-gonic/gin"
    
    "edugo-api-mobile/internal/interfaces/http/handlers"
    "edugo-api-mobile/internal/interfaces/http/middleware"
)

func SetupRoutes(r *gin.Engine, h *handlers.AssessmentHandler) {
    v1 := r.Group("/v1")
    v1.Use(middleware.AuthMiddleware()) // JWT auth requerido
    
    {
        // GET /v1/materials/:id/assessment
        v1.GET("/materials/:id/assessment", h.GetAssessment)
        
        // POST /v1/materials/:id/assessment/attempts
        v1.POST("/materials/:id/assessment/attempts", h.CreateAttempt)
        
        // GET /v1/attempts/:id/results
        v1.GET("/attempts/:id/results", h.GetAttemptResults)
        
        // GET /v1/users/me/attempts
        v1.GET("/users/me/attempts", h.GetUserAttempts)
    }
}
```

#### Tiempo Estimado
2 horas

---

### TASK-04-005: Swagger Annotations
**Tipo:** docs  
**Prioridad:** MEDIUM  
**Estimación:** 2h

#### Comandos
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Instalar swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generar docs
swag init -g cmd/api/main.go -o docs/swagger

# Verificar
ls docs/swagger/swagger.json
```

#### Tiempo Estimado
2 horas

---

### TASK-04-006: Tests E2E
**Tipo:** test  
**Prioridad:** HIGH  
**Estimación:** 4h

#### Implementación

Archivo: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/tests/e2e/assessment_e2e_test.go`

```go
//go:build e2e

package e2e_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func TestE2E_AssessmentFlow(t *testing.T) {
    // 1. Setup: Iniciar servidor de test
    router := setupTestRouter()
    
    // 2. GET /v1/materials/:id/assessment
    req, _ := http.NewRequest("GET", "/v1/materials/"+testMaterialID+"/assessment", nil)
    req.Header.Set("Authorization", "Bearer "+testJWT)
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    // Verificar que NO incluye correct_answer
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.NotContains(t, response, "correct_answer")
    
    // 3. POST /v1/materials/:id/assessment/attempts
    attemptReq := map[string]interface{}{
        "answers": []map[string]interface{}{
            {"question_id": "q1", "selected_answer_id": "a1"},
            {"question_id": "q2", "selected_answer_id": "a2"},
        },
    }
    body, _ := json.Marshal(attemptReq)
    
    req, _ = http.NewRequest("POST", "/v1/materials/"+testMaterialID+"/assessment/attempts", bytes.NewBuffer(body))
    req.Header.Set("Authorization", "Bearer "+testJWT)
    req.Header.Set("Content-Type", "application/json")
    
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
    
    // 4. Verificar que score fue calculado en servidor
    var attemptResp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &attemptResp)
    assert.Contains(t, attemptResp, "score")
    assert.IsType(t, float64(0), attemptResp["score"])
}
```

#### Tiempo Estimado
4 horas

---

## Resumen

**Total Tareas:** 6  
**Estimación:** 20 horas  
**Entregables:** 4 endpoints REST + Swagger + Tests E2E

---

**Sprint:** 04/06
