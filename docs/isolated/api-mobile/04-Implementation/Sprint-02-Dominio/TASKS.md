# Tareas del Sprint 02 - Capa de Dominio

## Objetivo
Implementar la capa de dominio del Sistema de Evaluaciones con 3 entities principales (Assessment, Attempt, Answer), 5+ value objects, 3 repository interfaces y tests unitarios con >90% coverage, siguiendo principios de Clean Architecture y Domain-Driven Design (DDD).

**Alcance:** Esta capa contiene SOLO lógica de negocio pura, sin dependencias a frameworks, bases de datos, o librerías externas (excepto testing y UUIDs).

---

## Tareas

### TASK-02-001: Crear Entity Assessment
**Tipo:** feature  
**Prioridad:** HIGH  
**Estimación:** 3h  
**Asignado a:** @ai-executor

#### Descripción
Crear la entity Assessment que representa una evaluación asociada a un material educativo. Esta entity encapsula las reglas de negocio relacionadas con evaluaciones: validación de parámetros, límites de intentos, límites de tiempo.

#### Pasos de Implementación

1. Crear archivo en ruta absoluta:
   `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/domain/entities/assessment.go`

2. Implementar struct Assessment con todos los campos del schema PostgreSQL:
   ```go
   package entities
   
   import (
       "errors"
       "time"
       "github.com/google/uuid"
   )
   
   // Assessment representa una evaluación de un material educativo
   // Esta entity corresponde a la tabla `assessment` en PostgreSQL
   type Assessment struct {
       ID                 uuid.UUID
       MaterialID         uuid.UUID
       MongoDocumentID    string  // ObjectId de MongoDB (24 caracteres hex)
       Title              string
       TotalQuestions     int
       PassThreshold      int // Porcentaje 0-100 para aprobar
       MaxAttempts        *int // nil = intentos ilimitados
       TimeLimitMinutes   *int // nil = sin límite de tiempo
       CreatedAt          time.Time
       UpdatedAt          time.Time
   }
   
   // NewAssessment crea una nueva evaluación con validaciones
   // Este constructor aplica fail-fast: si alguna validación falla, retorna error inmediatamente
   func NewAssessment(
       materialID uuid.UUID,
       mongoDocID string,
       title string,
       totalQuestions int,
       passThreshold int,
   ) (*Assessment, error) {
       // Validar material ID
       if materialID == uuid.Nil {
           return nil, ErrInvalidMaterialID
       }
       
       // Validar MongoDB document ID (24 caracteres hexadecimales)
       if len(mongoDocID) != 24 {
           return nil, ErrInvalidMongoDocumentID
       }
       
       // Validar que el título no esté vacío
       if title == "" {
           return nil, ErrEmptyTitle
       }
       
       // Validar total de preguntas (1-100 según schema PostgreSQL)
       if totalQuestions < 1 || totalQuestions > 100 {
           return nil, ErrInvalidTotalQuestions
       }
       
       // Validar umbral de aprobación (0-100)
       if passThreshold < 0 || passThreshold > 100 {
           return nil, ErrInvalidPassThreshold
       }
       
       now := time.Now().UTC()
       return &Assessment{
           ID:                 uuid.New(),
           MaterialID:         materialID,
           MongoDocumentID:    mongoDocID,
           Title:              title,
           TotalQuestions:     totalQuestions,
           PassThreshold:      passThreshold,
           MaxAttempts:        nil, // Default: ilimitado
           TimeLimitMinutes:   nil, // Default: sin límite
           CreatedAt:          now,
           UpdatedAt:          now,
       }, nil
   }
   
   // Validate verifica que la evaluación sea válida en su estado actual
   func (a *Assessment) Validate() error {
       if a.ID == uuid.Nil {
           return ErrInvalidAssessmentID
       }
       if a.MaterialID == uuid.Nil {
           return ErrInvalidMaterialID
       }
       if len(a.MongoDocumentID) != 24 {
           return ErrInvalidMongoDocumentID
       }
       if a.Title == "" {
           return ErrEmptyTitle
       }
       if a.TotalQuestions < 1 || a.TotalQuestions > 100 {
           return ErrInvalidTotalQuestions
       }
       if a.PassThreshold < 0 || a.PassThreshold > 100 {
           return ErrInvalidPassThreshold
       }
       if a.MaxAttempts != nil && *a.MaxAttempts < 1 {
           return ErrInvalidMaxAttempts
       }
       if a.TimeLimitMinutes != nil && (*a.TimeLimitMinutes < 1 || *a.TimeLimitMinutes > 180) {
           return ErrInvalidTimeLimit
       }
       return nil
   }
   
   // CanAttempt verifica si un estudiante puede hacer otro intento
   // Regla de negocio: si MaxAttempts es nil, intentos ilimitados
   func (a *Assessment) CanAttempt(attemptCount int) bool {
       if a.MaxAttempts == nil {
           return true // Ilimitado
       }
       return attemptCount < *a.MaxAttempts
   }
   
   // IsTimeLimited indica si la evaluación tiene límite de tiempo
   func (a *Assessment) IsTimeLimited() bool {
       return a.TimeLimitMinutes != nil && *a.TimeLimitMinutes > 0
   }
   
   // SetMaxAttempts establece el máximo de intentos permitidos
   // Esta es una business rule: mínimo 1 intento debe permitirse
   func (a *Assessment) SetMaxAttempts(max int) error {
       if max < 1 {
           return ErrInvalidMaxAttempts
       }
       a.MaxAttempts = &max
       a.UpdatedAt = time.Now().UTC()
       return nil
   }
   
   // SetTimeLimit establece el límite de tiempo en minutos
   // Business rule: entre 1 y 180 minutos (3 horas)
   func (a *Assessment) SetTimeLimit(minutes int) error {
       if minutes < 1 || minutes > 180 {
           return ErrInvalidTimeLimit
       }
       a.TimeLimitMinutes = &minutes
       a.UpdatedAt = time.Now().UTC()
       return nil
   }
   
   // RemoveMaxAttempts quita el límite de intentos (ilimitados)
   func (a *Assessment) RemoveMaxAttempts() {
       a.MaxAttempts = nil
       a.UpdatedAt = time.Now().UTC()
   }
   
   // RemoveTimeLimit quita el límite de tiempo
   func (a *Assessment) RemoveTimeLimit() {
       a.TimeLimitMinutes = nil
       a.UpdatedAt = time.Now().UTC()
   }
   ```

3. Crear errores de dominio en `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/domain/errors/errors.go`:
   ```go
   package errors
   
   import "errors"
   
   // Assessment errors
   var (
       ErrInvalidAssessmentID      = errors.New("domain: invalid assessment ID")
       ErrInvalidMaterialID        = errors.New("domain: invalid material ID")
       ErrInvalidMongoDocumentID   = errors.New("domain: mongo document ID must be exactly 24 characters")
       ErrEmptyTitle               = errors.New("domain: assessment title cannot be empty")
       ErrInvalidTotalQuestions    = errors.New("domain: total questions must be between 1 and 100")
       ErrInvalidPassThreshold     = errors.New("domain: pass threshold must be between 0 and 100")
       ErrInvalidMaxAttempts       = errors.New("domain: max attempts must be at least 1")
       ErrInvalidTimeLimit         = errors.New("domain: time limit must be between 1 and 180 minutes")
   )
   
   // Attempt errors
   var (
       ErrInvalidAttemptID         = errors.New("domain: invalid attempt ID")
       ErrInvalidStudentID         = errors.New("domain: invalid student ID")
       ErrInvalidScore             = errors.New("domain: score must be between 0 and 100")
       ErrInvalidTimeSpent         = errors.New("domain: time spent must be positive and <= 7200 seconds")
       ErrInvalidStartTime         = errors.New("domain: invalid start time")
       ErrInvalidEndTime           = errors.New("domain: end time must be after start time")
       ErrAttemptAlreadyCompleted  = errors.New("domain: attempt already completed, cannot modify")
       ErrNoAnswersProvided        = errors.New("domain: at least one answer must be provided")
   )
   
   // Answer errors
   var (
       ErrInvalidAnswerID          = errors.New("domain: invalid answer ID")
       ErrInvalidQuestionID        = errors.New("domain: invalid question ID")
       ErrInvalidSelectedAnswerID  = errors.New("domain: invalid selected answer ID")
   )
   ```

4. Crear tests unitarios en `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/domain/entities/assessment_test.go`:
   ```go
   package entities_test
   
   import (
       "testing"
       "time"
       
       "github.com/google/uuid"
       "github.com/stretchr/testify/assert"
       "github.com/stretchr/testify/require"
       
       "edugo-api-mobile/internal/domain/entities"
       domainErrors "edugo-api-mobile/internal/domain/errors"
   )
   
   func TestNewAssessment_Success(t *testing.T) {
       materialID := uuid.New()
       mongoDocID := "507f1f77bcf86cd799439011"
       title := "Cuestionario: Introducción a Pascal"
       totalQuestions := 5
       passThreshold := 70
       
       assessment, err := entities.NewAssessment(
           materialID,
           mongoDocID,
           title,
           totalQuestions,
           passThreshold,
       )
       
       require.NoError(t, err)
       require.NotNil(t, assessment)
       
       assert.NotEqual(t, uuid.Nil, assessment.ID)
       assert.Equal(t, materialID, assessment.MaterialID)
       assert.Equal(t, mongoDocID, assessment.MongoDocumentID)
       assert.Equal(t, title, assessment.Title)
       assert.Equal(t, totalQuestions, assessment.TotalQuestions)
       assert.Equal(t, passThreshold, assessment.PassThreshold)
       assert.Nil(t, assessment.MaxAttempts, "MaxAttempts should be nil by default")
       assert.Nil(t, assessment.TimeLimitMinutes, "TimeLimitMinutes should be nil by default")
       assert.False(t, assessment.CreatedAt.IsZero())
       assert.False(t, assessment.UpdatedAt.IsZero())
   }
   
   func TestNewAssessment_InvalidMaterialID(t *testing.T) {
       _, err := entities.NewAssessment(
           uuid.Nil,
           "507f1f77bcf86cd799439011",
           "Title",
           5,
           70,
       )
       
       assert.ErrorIs(t, err, domainErrors.ErrInvalidMaterialID)
   }
   
   func TestNewAssessment_InvalidMongoDocumentID(t *testing.T) {
       testCases := []struct {
           name       string
           mongoDocID string
       }{
           {"empty", ""},
           {"too short", "123"},
           {"too long", "507f1f77bcf86cd799439011EXTRA"},
           {"wrong length", "507f1f77bcf86cd79943901"},
       }
       
       for _, tc := range testCases {
           t.Run(tc.name, func(t *testing.T) {
               _, err := entities.NewAssessment(
                   uuid.New(),
                   tc.mongoDocID,
                   "Title",
                   5,
                   70,
               )
               
               assert.ErrorIs(t, err, domainErrors.ErrInvalidMongoDocumentID)
           })
       }
   }
   
   func TestNewAssessment_EmptyTitle(t *testing.T) {
       _, err := entities.NewAssessment(
           uuid.New(),
           "507f1f77bcf86cd799439011",
           "",
           5,
           70,
       )
       
       assert.ErrorIs(t, err, domainErrors.ErrEmptyTitle)
   }
   
   func TestNewAssessment_InvalidTotalQuestions(t *testing.T) {
       testCases := []struct {
           name           string
           totalQuestions int
       }{
           {"zero questions", 0},
           {"negative questions", -1},
           {"too many questions", 101},
           {"way too many", 1000},
       }
       
       for _, tc := range testCases {
           t.Run(tc.name, func(t *testing.T) {
               _, err := entities.NewAssessment(
                   uuid.New(),
                   "507f1f77bcf86cd799439011",
                   "Title",
                   tc.totalQuestions,
                   70,
               )
               
               assert.ErrorIs(t, err, domainErrors.ErrInvalidTotalQuestions)
           })
       }
   }
   
   func TestNewAssessment_InvalidPassThreshold(t *testing.T) {
       testCases := []struct {
           name          string
           passThreshold int
       }{
           {"negative threshold", -1},
           {"above 100", 101},
           {"way above 100", 150},
       }
       
       for _, tc := range testCases {
           t.Run(tc.name, func(t *testing.T) {
               _, err := entities.NewAssessment(
                   uuid.New(),
                   "507f1f77bcf86cd799439011",
                   "Title",
                   5,
                   tc.passThreshold,
               )
               
               assert.ErrorIs(t, err, domainErrors.ErrInvalidPassThreshold)
           })
       }
   }
   
   func TestAssessment_Validate_Success(t *testing.T) {
       assessment, _ := entities.NewAssessment(
           uuid.New(),
           "507f1f77bcf86cd799439011",
           "Title",
           5,
           70,
       )
       
       err := assessment.Validate()
       assert.NoError(t, err)
   }
   
   func TestAssessment_CanAttempt_Unlimited(t *testing.T) {
       assessment, _ := entities.NewAssessment(
           uuid.New(),
           "507f1f77bcf86cd799439011",
           "Title",
           5,
           70,
       )
       
       // Sin límite de intentos (MaxAttempts = nil)
       assert.True(t, assessment.CanAttempt(0))
       assert.True(t, assessment.CanAttempt(10))
       assert.True(t, assessment.CanAttempt(100))
       assert.True(t, assessment.CanAttempt(1000))
   }
   
   func TestAssessment_CanAttempt_WithLimit(t *testing.T) {
       assessment, _ := entities.NewAssessment(
           uuid.New(),
           "507f1f77bcf86cd799439011",
           "Title",
           5,
           70,
       )
       
       // Establecer límite de 3 intentos
       err := assessment.SetMaxAttempts(3)
       require.NoError(t, err)
       
       assert.True(t, assessment.CanAttempt(0), "should allow attempt 1")
       assert.True(t, assessment.CanAttempt(1), "should allow attempt 2")
       assert.True(t, assessment.CanAttempt(2), "should allow attempt 3")
       assert.False(t, assessment.CanAttempt(3), "should NOT allow attempt 4")
       assert.False(t, assessment.CanAttempt(4), "should NOT allow attempt 5")
   }
   
   func TestAssessment_SetMaxAttempts_Success(t *testing.T) {
       assessment, _ := entities.NewAssessment(
           uuid.New(),
           "507f1f77bcf86cd799439011",
           "Title",
           5,
           70,
       )
       
       oldUpdatedAt := assessment.UpdatedAt
       time.Sleep(10 * time.Millisecond)
       
       err := assessment.SetMaxAttempts(5)
       
       assert.NoError(t, err)
       assert.NotNil(t, assessment.MaxAttempts)
       assert.Equal(t, 5, *assessment.MaxAttempts)
       assert.True(t, assessment.UpdatedAt.After(oldUpdatedAt), "UpdatedAt should be updated")
   }
   
   func TestAssessment_SetMaxAttempts_Invalid(t *testing.T) {
       assessment, _ := entities.NewAssessment(
           uuid.New(),
           "507f1f77bcf86cd799439011",
           "Title",
           5,
           70,
       )
       
       testCases := []int{0, -1, -10}
       
       for _, maxAttempts := range testCases {
           err := assessment.SetMaxAttempts(maxAttempts)
           assert.ErrorIs(t, err, domainErrors.ErrInvalidMaxAttempts)
       }
   }
   
   func TestAssessment_RemoveMaxAttempts(t *testing.T) {
       assessment, _ := entities.NewAssessment(
           uuid.New(),
           "507f1f77bcf86cd799439011",
           "Title",
           5,
           70,
       )
       
       // Primero establecer límite
       assessment.SetMaxAttempts(3)
       assert.NotNil(t, assessment.MaxAttempts)
       
       // Luego quitar límite
       assessment.RemoveMaxAttempts()
       assert.Nil(t, assessment.MaxAttempts)
       
       // Ahora debería permitir intentos ilimitados
       assert.True(t, assessment.CanAttempt(100))
   }
   
   func TestAssessment_IsTimeLimited(t *testing.T) {
       assessment, _ := entities.NewAssessment(
           uuid.New(),
           "507f1f77bcf86cd799439011",
           "Title",
           5,
           70,
       )
       
       // Sin límite de tiempo
       assert.False(t, assessment.IsTimeLimited())
       
       // Con límite de tiempo
       assessment.SetTimeLimit(30)
       assert.True(t, assessment.IsTimeLimited())
       
       // Quitar límite
       assessment.RemoveTimeLimit()
       assert.False(t, assessment.IsTimeLimited())
   }
   
   func TestAssessment_SetTimeLimit_Success(t *testing.T) {
       assessment, _ := entities.NewAssessment(
           uuid.New(),
           "507f1f77bcf86cd799439011",
           "Title",
           5,
           70,
       )
       
       testCases := []int{1, 15, 30, 60, 120, 180}
       
       for _, minutes := range testCases {
           err := assessment.SetTimeLimit(minutes)
           assert.NoError(t, err)
           assert.NotNil(t, assessment.TimeLimitMinutes)
           assert.Equal(t, minutes, *assessment.TimeLimitMinutes)
       }
   }
   
   func TestAssessment_SetTimeLimit_Invalid(t *testing.T) {
       assessment, _ := entities.NewAssessment(
           uuid.New(),
           "507f1f77bcf86cd799439011",
           "Title",
           5,
           70,
       )
       
       testCases := []int{0, -1, 181, 200, 1000}
       
       for _, minutes := range testCases {
           err := assessment.SetTimeLimit(minutes)
           assert.ErrorIs(t, err, domainErrors.ErrInvalidTimeLimit)
       }
   }
   ```

#### Criterios de Aceptación
- [ ] Archivo `assessment.go` creado en ruta absoluta especificada
- [ ] Struct Assessment con 10 campos (ID, MaterialID, MongoDocumentID, Title, TotalQuestions, PassThreshold, MaxAttempts, TimeLimitMinutes, CreatedAt, UpdatedAt)
- [ ] Constructor NewAssessment() con validaciones fail-fast
- [ ] Método Validate() verifica estado completo de la entity
- [ ] Método CanAttempt() implementa regla de negocio de intentos
- [ ] Métodos SetMaxAttempts(), RemoveMaxAttempts(), SetTimeLimit(), RemoveTimeLimit()
- [ ] Archivo `errors.go` con errores de dominio
- [ ] Tests unitarios con >90% coverage
- [ ] Tests cubren casos exitosos y todos los casos de error
- [ ] Tests verifican reglas de negocio (intentos ilimitados, límites de tiempo, etc.)

#### Comandos de Validación
```bash
# Cambiar a directorio del proyecto
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Compilar el paquete (verificar que no hay errores de sintaxis)
go build ./internal/domain/entities

# Ejecutar tests del entity Assessment
go test ./internal/domain/entities -v -run TestAssessment

# Verificar coverage
go test ./internal/domain/entities -cover -coverprofile=coverage.out
go tool cover -func=coverage.out | grep assessment.go
# Esperado: >90%

# Generar reporte HTML de coverage
go tool cover -html=coverage.out -o coverage.html
open coverage.html  # Verificar visualmente las líneas cubiertas
```

#### Dependencias
- Requiere: Go 1.21+ instalado
- Usa: `github.com/google/uuid` para generación de IDs
- Usa: `github.com/stretchr/testify` v1.8.4 para assertions en tests

```bash
# Instalar dependencias si no existen
go get github.com/google/uuid@latest
go get github.com/stretchr/testify@v1.8.4
```

#### Tiempo Estimado
3 horas

---

### TASK-02-002: Crear Entity Attempt
**Tipo:** feature  
**Prioridad:** HIGH  
**Estimación:** 3h  
**Asignado a:** @ai-executor

#### Descripción
Crear la entity Attempt que representa un intento de un estudiante en una evaluación. Esta entity es **INMUTABLE** después de ser creada (siguiendo evento de dominio "intento completado"). Encapsula las reglas de cálculo de score y validación de respuestas.

#### Pasos de Implementación

1. Crear archivo:
   `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/domain/entities/attempt.go`

2. Implementar struct Attempt INMUTABLE:
   ```go
   package entities
   
   import (
       "time"
       "github.com/google/uuid"
   )
   
   // Attempt representa un intento de un estudiante en una evaluación
   // Esta entity es INMUTABLE: una vez creada, no se puede modificar
   // Corresponde a la tabla `assessment_attempt` en PostgreSQL
   type Attempt struct {
       ID                  uuid.UUID
       AssessmentID        uuid.UUID
       StudentID           uuid.UUID
       Score               int // 0-100
       MaxScore            int // Siempre 100
       TimeSpentSeconds    int // Tiempo total en segundos
       StartedAt           time.Time
       CompletedAt         time.Time
       CreatedAt           time.Time
       Answers             []*Answer // Respuestas del intento
       IdempotencyKey      *string // Para prevenir duplicados (Post-MVP)
   }
   
   // NewAttempt crea un nuevo intento COMPLETO
   // Business rule: un intento se crea YA COMPLETADO, no hay estado "en progreso"
   func NewAttempt(
       assessmentID uuid.UUID,
       studentID uuid.UUID,
       answers []*Answer,
       startedAt time.Time,
       completedAt time.Time,
   ) (*Attempt, error) {
       // Validaciones básicas
       if assessmentID == uuid.Nil {
           return nil, ErrInvalidAssessmentID
       }
       
       if studentID == uuid.Nil {
           return nil, ErrInvalidStudentID
       }
       
       if len(answers) == 0 {
           return nil, ErrNoAnswersProvided
       }
       
       if startedAt.IsZero() {
           return nil, ErrInvalidStartTime
       }
       
       if completedAt.IsZero() || !completedAt.After(startedAt) {
           return nil, ErrInvalidEndTime
       }
       
       // Calcular tiempo gastado
       timeSpent := int(completedAt.Sub(startedAt).Seconds())
       if timeSpent <= 0 || timeSpent > 7200 { // Máximo 2 horas
           return nil, ErrInvalidTimeSpent
       }
       
       // Calcular score basándose en respuestas correctas
       correctAnswers := 0
       for _, answer := range answers {
           if answer.IsCorrect {
               correctAnswers++
           }
       }
       
       totalQuestions := len(answers)
       score := (correctAnswers * 100) / totalQuestions
       
       return &Attempt{
           ID:               uuid.New(),
           AssessmentID:     assessmentID,
           StudentID:        studentID,
           Score:            score,
           MaxScore:         100,
           TimeSpentSeconds: timeSpent,
           StartedAt:        startedAt.UTC(),
           CompletedAt:      completedAt.UTC(),
           CreatedAt:        time.Now().UTC(),
           Answers:          answers,
           IdempotencyKey:   nil,
       }, nil
   }
   
   // NewAttemptWithIdempotency crea intento con clave de idempotencia (Post-MVP)
   func NewAttemptWithIdempotency(
       assessmentID uuid.UUID,
       studentID uuid.UUID,
       answers []*Answer,
       startedAt time.Time,
       completedAt time.Time,
       idempotencyKey string,
   ) (*Attempt, error) {
       attempt, err := NewAttempt(assessmentID, studentID, answers, startedAt, completedAt)
       if err != nil {
           return nil, err
       }
       
       attempt.IdempotencyKey = &idempotencyKey
       return attempt, nil
   }
   
   // IsPassed indica si el intento aprobó la evaluación
   func (a *Attempt) IsPassed(passThreshold int) bool {
       return a.Score >= passThreshold
   }
   
   // GetCorrectAnswersCount retorna la cantidad de respuestas correctas
   func (a *Attempt) GetCorrectAnswersCount() int {
       count := 0
       for _, answer := range a.Answers {
           if answer.IsCorrect {
               count++
           }
       }
       return count
   }
   
   // GetIncorrectAnswersCount retorna la cantidad de respuestas incorrectas
   func (a *Attempt) GetIncorrectAnswersCount() int {
       return len(a.Answers) - a.GetCorrectAnswersCount()
   }
   
   // GetTotalQuestions retorna el total de preguntas respondidas
   func (a *Attempt) GetTotalQuestions() int {
       return len(a.Answers)
   }
   
   // GetAccuracyPercentage retorna el porcentaje de precisión (alias de Score)
   func (a *Attempt) GetAccuracyPercentage() int {
       return a.Score
   }
   
   // GetAverageTimePerQuestion retorna el tiempo promedio por pregunta en segundos
   func (a *Attempt) GetAverageTimePerQuestion() int {
       if len(a.Answers) == 0 {
           return 0
       }
       return a.TimeSpentSeconds / len(a.Answers)
   }
   
   // Validate verifica la integridad del intento
   func (a *Attempt) Validate() error {
       if a.ID == uuid.Nil {
           return ErrInvalidAttemptID
       }
       if a.AssessmentID == uuid.Nil {
           return ErrInvalidAssessmentID
       }
       if a.StudentID == uuid.Nil {
           return ErrInvalidStudentID
       }
       if a.Score < 0 || a.Score > 100 {
           return ErrInvalidScore
       }
       if a.TimeSpentSeconds <= 0 || a.TimeSpentSeconds > 7200 {
           return ErrInvalidTimeSpent
       }
       if a.StartedAt.IsZero() {
           return ErrInvalidStartTime
       }
       if a.CompletedAt.IsZero() || !a.CompletedAt.After(a.StartedAt) {
           return ErrInvalidEndTime
       }
       if len(a.Answers) == 0 {
           return ErrNoAnswersProvided
       }
       
       // Verificar que el score calculado coincide
       correctCount := a.GetCorrectAnswersCount()
       expectedScore := (correctCount * 100) / len(a.Answers)
       if a.Score != expectedScore {
           return errors.New("domain: score mismatch with answers")
       }
       
       return nil
   }
   ```

3. Crear tests en `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/domain/entities/attempt_test.go`:
   ```go
   package entities_test
   
   import (
       "testing"
       "time"
       
       "github.com/google/uuid"
       "github.com/stretchr/testify/assert"
       "github.com/stretchr/testify/require"
       
       "edugo-api-mobile/internal/domain/entities"
       domainErrors "edugo-api-mobile/internal/domain/errors"
   )
   
   func TestNewAttempt_Success(t *testing.T) {
       assessmentID := uuid.New()
       studentID := uuid.New()
       
       // Crear 5 respuestas: 3 correctas, 2 incorrectas (60%)
       answers := []*entities.Answer{
           {ID: uuid.New(), QuestionID: "q1", SelectedAnswerID: "a1", IsCorrect: true, TimeSpentSeconds: 30},
           {ID: uuid.New(), QuestionID: "q2", SelectedAnswerID: "a2", IsCorrect: false, TimeSpentSeconds: 45},
           {ID: uuid.New(), QuestionID: "q3", SelectedAnswerID: "a3", IsCorrect: true, TimeSpentSeconds: 60},
           {ID: uuid.New(), QuestionID: "q4", SelectedAnswerID: "a4", IsCorrect: false, TimeSpentSeconds: 50},
           {ID: uuid.New(), QuestionID: "q5", SelectedAnswerID: "a5", IsCorrect: true, TimeSpentSeconds: 40},
       }
       
       startedAt := time.Now().Add(-5 * time.Minute)
       completedAt := time.Now()
       
       attempt, err := entities.NewAttempt(assessmentID, studentID, answers, startedAt, completedAt)
       
       require.NoError(t, err)
       require.NotNil(t, attempt)
       
       assert.NotEqual(t, uuid.Nil, attempt.ID)
       assert.Equal(t, assessmentID, attempt.AssessmentID)
       assert.Equal(t, studentID, attempt.StudentID)
       assert.Equal(t, 60, attempt.Score, "Score should be 60% (3/5 correct)")
       assert.Equal(t, 100, attempt.MaxScore)
       assert.Equal(t, 5, len(attempt.Answers))
       assert.True(t, attempt.TimeSpentSeconds > 0)
   }
   
   func TestNewAttempt_ScoreCalculation(t *testing.T) {
       testCases := []struct {
           name          string
           correctCount  int
           totalQuestions int
           expectedScore int
       }{
           {"0% - all wrong", 0, 5, 0},
           {"20% - 1 of 5", 1, 5, 20},
           {"40% - 2 of 5", 2, 5, 40},
           {"60% - 3 of 5", 3, 5, 60},
           {"80% - 4 of 5", 4, 5, 80},
           {"100% - all correct", 5, 5, 100},
           {"100% - single question", 1, 1, 100},
           {"75% - 3 of 4", 3, 4, 75},
       }
       
       for _, tc := range testCases {
           t.Run(tc.name, func(t *testing.T) {
               answers := make([]*entities.Answer, tc.totalQuestions)
               for i := 0; i < tc.totalQuestions; i++ {
                   isCorrect := i < tc.correctCount
                   answers[i] = &entities.Answer{
                       ID:                uuid.New(),
                       QuestionID:        "q" + string(rune(i)),
                       SelectedAnswerID:  "a1",
                       IsCorrect:         isCorrect,
                       TimeSpentSeconds:  30,
                   }
               }
               
               startedAt := time.Now().Add(-3 * time.Minute)
               completedAt := time.Now()
               
               attempt, err := entities.NewAttempt(uuid.New(), uuid.New(), answers, startedAt, completedAt)
               
               require.NoError(t, err)
               assert.Equal(t, tc.expectedScore, attempt.Score)
           })
       }
   }
   
   func TestNewAttempt_InvalidAssessmentID(t *testing.T) {
       answers := []*entities.Answer{{ID: uuid.New(), QuestionID: "q1", IsCorrect: true}}
       
       _, err := entities.NewAttempt(
           uuid.Nil,
           uuid.New(),
           answers,
           time.Now().Add(-1*time.Minute),
           time.Now(),
       )
       
       assert.ErrorIs(t, err, domainErrors.ErrInvalidAssessmentID)
   }
   
   func TestNewAttempt_InvalidStudentID(t *testing.T) {
       answers := []*entities.Answer{{ID: uuid.New(), QuestionID: "q1", IsCorrect: true}}
       
       _, err := entities.NewAttempt(
           uuid.New(),
           uuid.Nil,
           answers,
           time.Now().Add(-1*time.Minute),
           time.Now(),
       )
       
       assert.ErrorIs(t, err, domainErrors.ErrInvalidStudentID)
   }
   
   func TestNewAttempt_NoAnswers(t *testing.T) {
       _, err := entities.NewAttempt(
           uuid.New(),
           uuid.New(),
           []*entities.Answer{},
           time.Now().Add(-1*time.Minute),
           time.Now(),
       )
       
       assert.ErrorIs(t, err, domainErrors.ErrNoAnswersProvided)
   }
   
   func TestNewAttempt_InvalidEndTime(t *testing.T) {
       answers := []*entities.Answer{{ID: uuid.New(), QuestionID: "q1", IsCorrect: true}}
       now := time.Now()
       
       // CompletedAt antes de StartedAt
       _, err := entities.NewAttempt(
           uuid.New(),
           uuid.New(),
           answers,
           now,
           now.Add(-1*time.Minute), // Completado ANTES de empezar
       )
       
       assert.ErrorIs(t, err, domainErrors.ErrInvalidEndTime)
   }
   
   func TestAttempt_IsPassed(t *testing.T) {
       testCases := []struct {
           name          string
           score         int
           passThreshold int
           shouldPass    bool
       }{
           {"60% with 70% threshold - fail", 60, 70, false},
           {"70% with 70% threshold - pass", 70, 70, true},
           {"80% with 70% threshold - pass", 80, 70, true},
           {"100% with 70% threshold - pass", 100, 70, true},
           {"0% - fail", 0, 70, false},
       }
       
       for _, tc := range testCases {
           t.Run(tc.name, func(t *testing.T) {
               // Crear intento con score específico
               totalQuestions := 10
               correctCount := (tc.score * totalQuestions) / 100
               
               answers := make([]*entities.Answer, totalQuestions)
               for i := 0; i < totalQuestions; i++ {
                   answers[i] = &entities.Answer{
                       ID:                uuid.New(),
                       QuestionID:        "q",
                       IsCorrect:         i < correctCount,
                       TimeSpentSeconds:  30,
                   }
               }
               
               attempt, _ := entities.NewAttempt(
                   uuid.New(),
                   uuid.New(),
                   answers,
                   time.Now().Add(-5*time.Minute),
                   time.Now(),
               )
               
               assert.Equal(t, tc.shouldPass, attempt.IsPassed(tc.passThreshold))
           })
       }
   }
   
   func TestAttempt_GetCorrectAnswersCount(t *testing.T) {
       answers := []*entities.Answer{
           {ID: uuid.New(), QuestionID: "q1", IsCorrect: true},
           {ID: uuid.New(), QuestionID: "q2", IsCorrect: false},
           {ID: uuid.New(), QuestionID: "q3", IsCorrect: true},
           {ID: uuid.New(), QuestionID: "q4", IsCorrect: true},
       }
       
       attempt, _ := entities.NewAttempt(
           uuid.New(),
           uuid.New(),
           answers,
           time.Now().Add(-2*time.Minute),
           time.Now(),
       )
       
       assert.Equal(t, 3, attempt.GetCorrectAnswersCount())
       assert.Equal(t, 1, attempt.GetIncorrectAnswersCount())
       assert.Equal(t, 4, attempt.GetTotalQuestions())
   }
   
   func TestAttempt_GetAverageTimePerQuestion(t *testing.T) {
       answers := []*entities.Answer{
           {ID: uuid.New(), QuestionID: "q1", IsCorrect: true, TimeSpentSeconds: 30},
           {ID: uuid.New(), QuestionID: "q2", IsCorrect: true, TimeSpentSeconds: 60},
           {ID: uuid.New(), QuestionID: "q3", IsCorrect: true, TimeSpentSeconds: 90},
       }
       
       startedAt := time.Now().Add(-3 * time.Minute)
       completedAt := time.Now()
       
       attempt, _ := entities.NewAttempt(
           uuid.New(),
           uuid.New(),
           answers,
           startedAt,
           completedAt,
       )
       
       avgTime := attempt.GetAverageTimePerQuestion()
       assert.True(t, avgTime > 0)
       assert.True(t, avgTime <= 180) // ~180 segundos / 3 preguntas = ~60 seg promedio
   }
   
   func TestAttempt_Validate_Success(t *testing.T) {
       answers := []*entities.Answer{
           {ID: uuid.New(), QuestionID: "q1", IsCorrect: true, TimeSpentSeconds: 60},
       }
       
       attempt, _ := entities.NewAttempt(
           uuid.New(),
           uuid.New(),
           answers,
           time.Now().Add(-2*time.Minute),
           time.Now(),
       )
       
       err := attempt.Validate()
       assert.NoError(t, err)
   }
   ```

#### Criterios de Aceptación
- [ ] Entity Attempt es INMUTABLE (no hay setters)
- [ ] Constructor NewAttempt() calcula score automáticamente
- [ ] Score calculado correctamente: (correctas / total) * 100
- [ ] Validación de tiempos: CompletedAt > StartedAt
- [ ] Validación de TimeSpent: >0 y <=7200 segundos
- [ ] Método IsPassed() compara con threshold
- [ ] Métodos de consulta: GetCorrectAnswersCount(), GetIncorrectAnswersCount(), GetAverageTimePerQuestion()
- [ ] Tests cubren cálculo de score (0%, 20%, 40%, 60%, 80%, 100%)
- [ ] Tests cubren validaciones de tiempos
- [ ] Coverage >90%

#### Comandos de Validación
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

go test ./internal/domain/entities -v -run TestAttempt
go test ./internal/domain/entities -cover | grep attempt.go
```

#### Dependencias
- Requiere: TASK-02-003 (Entity Answer) debe existir
- Usa: github.com/google/uuid
- Usa: github.com/stretchr/testify v1.8.4

#### Tiempo Estimado
3 horas

---

### TASK-02-003: Crear Entity Answer
**Tipo:** feature  
**Prioridad:** HIGH  
**Estimación:** 2h  
**Asignado a:** @ai-executor

#### Descripción
Crear la entity Answer que representa una respuesta individual a una pregunta dentro de un intento. Entity simple con validaciones básicas.

#### Pasos de Implementación

1. Crear archivo:
   `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/domain/entities/answer.go`

2. Implementar struct Answer:
   ```go
   package entities
   
   import (
       "time"
       "github.com/google/uuid"
   )
   
   // Answer representa una respuesta individual a una pregunta en un intento
   // Corresponde a la tabla `assessment_attempt_answer` en PostgreSQL
   type Answer struct {
       ID                  uuid.UUID
       AttemptID           uuid.UUID
       QuestionID          string  // ID de la pregunta en MongoDB
       SelectedAnswerID    string  // ID de la opción seleccionada en MongoDB
       IsCorrect           bool
       TimeSpentSeconds    int
       CreatedAt           time.Time
   }
   
   // NewAnswer crea una nueva respuesta
   func NewAnswer(
       attemptID uuid.UUID,
       questionID string,
       selectedAnswerID string,
       isCorrect bool,
       timeSpent int,
   ) (*Answer, error) {
       if attemptID == uuid.Nil {
           return nil, ErrInvalidAttemptID
       }
       
       if questionID == "" {
           return nil, ErrInvalidQuestionID
       }
       
       if selectedAnswerID == "" {
           return nil, ErrInvalidSelectedAnswerID
       }
       
       if timeSpent < 0 {
           return nil, ErrInvalidTimeSpent
       }
       
       return &Answer{
           ID:                uuid.New(),
           AttemptID:         attemptID,
           QuestionID:        questionID,
           SelectedAnswerID:  selectedAnswerID,
           IsCorrect:         isCorrect,
           TimeSpentSeconds:  timeSpent,
           CreatedAt:         time.Now().UTC(),
       }, nil
   }
   
   // Validate verifica la validez de la respuesta
   func (a *Answer) Validate() error {
       if a.ID == uuid.Nil {
           return ErrInvalidAnswerID
       }
       if a.AttemptID == uuid.Nil {
           return ErrInvalidAttemptID
       }
       if a.QuestionID == "" {
           return ErrInvalidQuestionID
       }
       if a.SelectedAnswerID == "" {
           return ErrInvalidSelectedAnswerID
       }
       if a.TimeSpentSeconds < 0 {
           return ErrInvalidTimeSpent
       }
       return nil
   }
   ```

3. Tests en `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/domain/entities/answer_test.go`:
   ```go
   package entities_test
   
   import (
       "testing"
       
       "github.com/google/uuid"
       "github.com/stretchr/testify/assert"
       "github.com/stretchr/testify/require"
       
       "edugo-api-mobile/internal/domain/entities"
       domainErrors "edugo-api-mobile/internal/domain/errors"
   )
   
   func TestNewAnswer_Success(t *testing.T) {
       attemptID := uuid.New()
       questionID := "q1"
       selectedAnswerID := "a2"
       isCorrect := true
       timeSpent := 45
       
       answer, err := entities.NewAnswer(attemptID, questionID, selectedAnswerID, isCorrect, timeSpent)
       
       require.NoError(t, err)
       require.NotNil(t, answer)
       
       assert.NotEqual(t, uuid.Nil, answer.ID)
       assert.Equal(t, attemptID, answer.AttemptID)
       assert.Equal(t, questionID, answer.QuestionID)
       assert.Equal(t, selectedAnswerID, answer.SelectedAnswerID)
       assert.Equal(t, isCorrect, answer.IsCorrect)
       assert.Equal(t, timeSpent, answer.TimeSpentSeconds)
   }
   
   func TestNewAnswer_InvalidAttemptID(t *testing.T) {
       _, err := entities.NewAnswer(uuid.Nil, "q1", "a1", true, 30)
       assert.ErrorIs(t, err, domainErrors.ErrInvalidAttemptID)
   }
   
   func TestNewAnswer_InvalidQuestionID(t *testing.T) {
       _, err := entities.NewAnswer(uuid.New(), "", "a1", true, 30)
       assert.ErrorIs(t, err, domainErrors.ErrInvalidQuestionID)
   }
   
   func TestNewAnswer_InvalidSelectedAnswerID(t *testing.T) {
       _, err := entities.NewAnswer(uuid.New(), "q1", "", true, 30)
       assert.ErrorIs(t, err, domainErrors.ErrInvalidSelectedAnswerID)
   }
   
   func TestNewAnswer_NegativeTimeSpent(t *testing.T) {
       _, err := entities.NewAnswer(uuid.New(), "q1", "a1", true, -10)
       assert.ErrorIs(t, err, domainErrors.ErrInvalidTimeSpent)
   }
   
   func TestAnswer_Validate_Success(t *testing.T) {
       answer, _ := entities.NewAnswer(uuid.New(), "q1", "a1", true, 30)
       err := answer.Validate()
       assert.NoError(t, err)
   }
   ```

#### Criterios de Aceptación
- [ ] Struct Answer con 7 campos
- [ ] Constructor NewAnswer() con validaciones
- [ ] Método Validate()
- [ ] Tests de casos válidos e inválidos
- [ ] Coverage >90%

#### Comandos de Validación
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
go test ./internal/domain/entities -v -run TestAnswer
go test ./internal/domain/entities -cover | grep answer.go
```

#### Dependencias
- Requiere: Go 1.21+
- Usa: github.com/google/uuid
- Usa: github.com/stretchr/testify v1.8.4

#### Tiempo Estimado
2 horas

---

### TASK-02-004: Crear Value Objects
**Tipo:** feature  
**Prioridad:** MEDIUM  
**Estimación:** 3h  
**Asignado a:** @ai-executor

#### Descripción
Crear value objects para encapsular valores primitivos con validaciones. En DDD, value objects son objetos inmutables que representan conceptos del dominio.

#### Pasos de Implementación

1. Crear `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/domain/valueobjects/score.go`:
   ```go
   package valueobjects
   
   import "fmt"
   
   // Score representa un puntaje de 0 a 100
   type Score struct {
       value int
   }
   
   // NewScore crea un Score válido
   func NewScore(value int) (Score, error) {
       if value < 0 || value > 100 {
           return Score{}, fmt.Errorf("score must be between 0 and 100, got %d", value)
       }
       return Score{value: value}, nil
   }
   
   // Value retorna el valor del score
   func (s Score) Value() int {
       return s.value
   }
   
   // IsPassing verifica si el score aprueba con un threshold
   func (s Score) IsPassing(threshold int) bool {
       return s.value >= threshold
   }
   
   // IsFailin

g verifica si el score reprueba
   func (s Score) IsFailing(threshold int) bool {
       return !s.IsPassing(threshold)
   }
   
   // String implementa Stringer
   func (s Score) String() string {
       return fmt.Sprintf("%d%%", s.value)
   }
   
   // Equals compara dos scores
   func (s Score) Equals(other Score) bool {
       return s.value == other.value
   }
   ```

2. Tests en `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/domain/valueobjects/score_test.go`

3. Crear otros value objects similares:
   - `assessment_id.go` - Wrapper de UUID con validación
   - `question_id.go` - String validado para IDs de preguntas MongoDB
   - `mongo_document_id.go` - String de 24 caracteres hex
   - `time_spent.go` - Validación de tiempo (0-7200 segundos)

#### Criterios de Aceptación
- [ ] 5 value objects creados
- [ ] Cada uno inmutable (sin setters)
- [ ] Método Equals() para comparación
- [ ] String() para representación
- [ ] Tests con coverage >90%

#### Comandos de Validación
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
go test ./internal/domain/valueobjects -v
go test ./internal/domain/valueobjects -cover
```

#### Tiempo Estimado
3 horas

---

### TASK-02-005: Crear Repository Interfaces
**Tipo:** feature  
**Prioridad:** HIGH  
**Estimación:** 2h  
**Asignado a:** @ai-executor

#### Descripción
Crear interfaces de repositorios en la capa de dominio. Estas interfaces definen el contrato, pero NO la implementación (eso va en infrastructure).

#### Pasos de Implementación

1. Crear `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/domain/repositories/assessment_repository.go`:
   ```go
   package repositories
   
   import (
       "context"
       "github.com/google/uuid"
       "edugo-api-mobile/internal/domain/entities"
   )
   
   // AssessmentRepository define el contrato para persistencia de evaluaciones
   type AssessmentRepository interface {
       // FindByID busca una evaluación por ID
       FindByID(ctx context.Context, id uuid.UUID) (*entities.Assessment, error)
       
       // FindByMaterialID busca una evaluación por material ID
       FindByMaterialID(ctx context.Context, materialID uuid.UUID) (*entities.Assessment, error)
       
       // Save guarda una evaluación (INSERT o UPDATE)
       Save(ctx context.Context, assessment *entities.Assessment) error
       
       // Delete elimina una evaluación
       Delete(ctx context.Context, id uuid.UUID) error
   }
   ```

2. Crear `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/domain/repositories/attempt_repository.go`:
   ```go
   package repositories
   
   import (
       "context"
       "github.com/google/uuid"
       "edugo-api-mobile/internal/domain/entities"
   )
   
   // AttemptRepository define el contrato para persistencia de intentos
   type AttemptRepository interface {
       // FindByID busca un intento por ID
       FindByID(ctx context.Context, id uuid.UUID) (*entities.Attempt, error)
       
       // FindByStudentAndAssessment busca intentos de un estudiante en una evaluación
       FindByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) ([]*entities.Attempt, error)
       
       // Save guarda un intento (solo INSERT, no UPDATE - inmutable)
       Save(ctx context.Context, attempt *entities.Attempt) error
       
       // CountByStudentAndAssessment cuenta intentos de un estudiante
       CountByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) (int, error)
       
       // FindByStudent busca todos los intentos de un estudiante (historial)
       FindByStudent(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]*entities.Attempt, error)
   }
   ```

3. Crear `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/domain/repositories/answer_repository.go`

#### Criterios de Aceptación
- [ ] 3 repository interfaces creadas
- [ ] Métodos con context.Context como primer parámetro
- [ ] Métodos retornan (*Entity, error) o ([]Entity, error)
- [ ] Solo interfaces, SIN implementación

#### Comandos de Validación
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
go build ./internal/domain/repositories
```

#### Tiempo Estimado
2 horas

---

### TASK-02-006: Tests Unitarios Completos de Dominio
**Tipo:** test  
**Prioridad:** HIGH  
**Estimación:** 2h  
**Asignado a:** @ai-executor

#### Descripción
Asegurar que la capa de dominio tiene >90% coverage y todos los edge cases están cubiertos.

#### Pasos de Implementación

1. Ejecutar coverage:
   ```bash
   cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
   go test ./internal/domain/... -coverprofile=coverage.out
   go tool cover -html=coverage.out -o coverage.html
   ```

2. Identificar gaps de coverage y agregar tests faltantes

3. Agregar tests de edge cases:
   - Nil pointers
   - Zero values
   - Boundary conditions
   - Concurrent access (si aplica)

#### Criterios de Aceptación
- [ ] Coverage >90% en `internal/domain/entities`
- [ ] Coverage >90% en `internal/domain/valueobjects`
- [ ] Tests de todos los métodos públicos
- [ ] Tests de todas las validaciones
- [ ] Tests de todas las business rules

#### Comandos de Validación
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Coverage por paquete
go test ./internal/domain/entities -cover
go test ./internal/domain/valueobjects -cover

# Coverage detallado
go test ./internal/domain/... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Verificar que cada paquete tiene >90%
go tool cover -func=coverage.out | grep -E "(entities|valueobjects)" | awk '{print $3}' | sed 's/%//' | awk '{if ($1 < 90) exit 1}'
```

#### Tiempo Estimado
2 horas

---

## Resumen del Sprint

**Total de Tareas:** 6  
**Estimación Total:** 17 horas  
**Archivos de Código a Crear:** ~15 archivos Go + tests

**Entregables:**
1. `internal/domain/entities/assessment.go` + tests
2. `internal/domain/entities/attempt.go` + tests
3. `internal/domain/entities/answer.go` + tests
4. `internal/domain/valueobjects/*.go` (5 archivos) + tests
5. `internal/domain/repositories/*.go` (3 interfaces)
6. `internal/domain/errors/errors.go`
7. Coverage >90% en toda la capa de dominio

**Criterios de Éxito Globales:**
- [ ] 3 entities implementadas con business logic
- [ ] 5+ value objects implementados
- [ ] 3 repository interfaces definidas
- [ ] Tests unitarios >90% coverage
- [ ] Sin dependencias a frameworks externos (solo UUID y testing)
- [ ] Todos los tests pasan
- [ ] golangci-lint sin errores

---

**Generado con:** Claude Code  
**Sprint:** 02/06  
**Última actualización:** 2025-11-14
