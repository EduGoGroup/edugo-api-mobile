# Sprint: Adaptar api-mobile a Entities Centralizadas

**Proyecto:** edugo-api-mobile  
**Fecha:** 20 de Noviembre, 2025  
**Dependencia:** infrastructure Sprint ENTITIES completado  
**Objetivo:** Reemplazar entities locales por entities centralizadas de infrastructure  
**Prioridad:** ALTA - Facilita mantenimiento y consistencia entre proyectos

---

## 🎯 Contexto

**Problema actual:**
- api-mobile tiene entities duplicadas en dos ubicaciones:
  - `internal/domain/entity/` - 4 entities (Material, User, MaterialVersion, Progress)
  - `internal/domain/entities/` - 3 entities (Assessment, Answer, Attempt)
- Estas entities contienen **lógica de negocio mezclada** con estructura de datos
- Cambios en BD requieren actualizar múltiples proyectos
- Riesgo de inconsistencias entre proyectos

**Solución:**
- Eliminar entities locales de api-mobile
- Importar entities centralizadas desde infrastructure
- **Mover lógica de negocio** de entities a domain services/DTOs
- Entities en infrastructure = solo estructura de datos (sin lógica)

---

## 📊 Análisis de Entities Actuales

### 📁 Entities Encontrados en api-mobile

#### A. Directorio: `internal/domain/entity/`

| # | Archivo | Entity | Líneas | Base de Datos | Tiene Lógica de Negocio |
|---|---------|--------|--------|---------------|-------------------------|
| 1 | `material.go` | `Material` | 177 | PostgreSQL | ✅ Sí (métodos: SetS3Info, Publish, Archive, etc.) |
| 2 | `user.go` | `User` | 67 | PostgreSQL | ❌ No (solo getters) |
| 3 | `material_version.go` | `MaterialVersion` | 109 | PostgreSQL | ❌ No (solo constructor + getters) |
| 4 | `progress.go` | `Progress` | 77 | PostgreSQL | ✅ Sí (método: UpdateProgress con validaciones) |

**Características:**
- Usan `valueobject` para IDs (MaterialID, UserID, etc.)
- Tienen constructores de negocio: `NewMaterial()`, `NewProgress()`
- Tienen reconstructores desde BD: `ReconstructMaterial()`, `ReconstructUser()`
- Fields privados con getters públicos (encapsulación DDD)

#### B. Directorio: `internal/domain/entities/`

| # | Archivo | Entity | Líneas | Base de Datos | Tiene Lógica de Negocio |
|---|---------|--------|--------|---------------|-------------------------|
| 5 | `assessment.go` | `Assessment` | 172 | PostgreSQL | ✅ Sí (métodos: CanAttempt, SetMaxAttempts, Validate, etc.) |
| 6 | `answer.go` | `Answer` | 67 | PostgreSQL | ✅ Sí (método: Validate) |
| 7 | `attempt.go` | `Attempt` | 205 | PostgreSQL | ✅ Sí (métodos: IsPassed, GetCorrectAnswersCount, Validate, etc.) |

**Características:**
- Usan `uuid.UUID` directamente (sin value objects)
- Campos públicos (sin encapsulación)
- Constructores con validaciones: `NewAssessment()`, `NewAnswer()`, `NewAttempt()`
- **Mucha lógica de negocio** embebida

---

### 🗺️ Mapeo a Infrastructure Entities

| Entity Actual (api-mobile) | Entity Infrastructure | Paquete Infrastructure | Complejidad de Adaptación |
|----------------------------|----------------------|------------------------|---------------------------|
| `entity.Material` | `postgres/entities.Material` | `github.com/EduGoGroup/edugo-infrastructure/postgres/entities` | 🟡 MEDIA - Mover lógica a MaterialService |
| `entity.User` | `postgres/entities.User` | `github.com/EduGoGroup/edugo-infrastructure/postgres/entities` | 🟢 BAJA - Sin lógica de negocio |
| `entity.MaterialVersion` | `postgres/entities.MaterialVersion` | `github.com/EduGoGroup/edugo-infrastructure/postgres/entities` | 🟢 BAJA - Solo constructor |
| `entity.Progress` | `postgres/entities.Progress` | `github.com/EduGoGroup/edugo-infrastructure/postgres/entities` | 🟡 MEDIA - Mover UpdateProgress a service |
| `entities.Assessment` | `postgres/entities.Assessment` | `github.com/EduGoGroup/edugo-infrastructure/postgres/entities` | 🔴 ALTA - Mucha lógica de negocio |
| `entities.Answer` | `postgres/entities.AssessmentAnswer` | `github.com/EduGoGroup/edugo-infrastructure/postgres/entities` | 🟡 MEDIA - Validaciones a service |
| `entities.Attempt` | `postgres/entities.AssessmentAttempt` | `github.com/EduGoGroup/edugo-infrastructure/postgres/entities` | 🔴 ALTA - Muchos métodos de negocio |

**Nota importante:** Infrastructure entities tendrán nombres exactos de tablas PostgreSQL:
- `assessment` → `Assessment`
- `assessment_answers` → `AssessmentAnswer`
- `assessment_attempts` → `AssessmentAttempt`

---

## 🔍 Análisis de Dependencias

### Archivos que Importan `domain/entity` (15 archivos)

```
internal/application/dto/material_dto.go
internal/application/service/auth_service_test.go
internal/application/service/material_service.go
internal/application/service/progress_service_test.go
internal/application/service/material_service_test.go
internal/application/service/progress_service.go
internal/infrastructure/persistence/postgres/repository/progress_repository_impl.go
internal/infrastructure/persistence/postgres/repository/material_repository_impl.go
internal/infrastructure/persistence/postgres/repository/material_repository_impl_test.go
internal/infrastructure/persistence/postgres/repository/user_repository_impl_test.go
internal/infrastructure/persistence/postgres/repository/user_repository_impl.go
internal/infrastructure/persistence/postgres/repository/progress_repository_impl_test.go
internal/domain/repository/user_repository.go
internal/domain/repository/progress_repository.go
internal/domain/repository/material_repository.go
```

### Archivos que Importan `domain/entities` (16 archivos)

```
internal/application/service/assessment_attempt_service.go
internal/infrastructure/persistence/postgres/repository/attempt_repository.go
internal/infrastructure/persistence/postgres/repository/attempt_repository_integration_test.go
internal/infrastructure/persistence/postgres/repository/attempt_repository_test.go
internal/infrastructure/persistence/postgres/repository/answer_repository_integration_test.go
internal/infrastructure/persistence/postgres/repository/answer_repository_test.go
internal/infrastructure/persistence/postgres/repository/assessment_repository_test.go
internal/infrastructure/persistence/postgres/repository/answer_repository.go
internal/infrastructure/persistence/postgres/repository/assessment_repository_integration_test.go
internal/infrastructure/persistence/postgres/repository/assessment_repository.go
internal/domain/repositories/attempt_repository.go
internal/domain/repositories/answer_repository.go
internal/domain/repositories/assessment_repository.go
internal/domain/entities/assessment_test.go
internal/domain/entities/answer_test.go
internal/domain/entities/attempt_test.go
```

**Total:** 31 archivos afectados

---

## 🚨 Complejidades Detectadas

### 1. Lógica de Negocio Embebida en Entities

**Entity: Material** (`entity.Material`)
```go
// Métodos de negocio que deben moverse:
- SetS3Info(s3Key, s3URL string) error
- MarkProcessingComplete() error
- Publish() error
- Archive() error
- IsDraft() bool
- IsPublished() bool
- IsProcessed() bool
```
**Solución:** Mover a `MaterialService` (ya existe)

---

**Entity: Progress** (`entity.Progress`)
```go
// Método de negocio:
- UpdateProgress(percentage, lastPage int) error
```
**Solución:** Mover a `ProgressService` (ya existe)

---

**Entity: Assessment** (`entities.Assessment`)
```go
// Métodos de negocio que deben moverse:
- Validate() error
- CanAttempt(attemptCount int) bool
- IsTimeLimited() bool
- SetMaxAttempts(max int) error
- SetTimeLimit(minutes int) error
- RemoveMaxAttempts()
- RemoveTimeLimit()
```
**Solución:** Crear `AssessmentDomainService` o incluir en `AssessmentAttemptService`

---

**Entity: Answer** (`entities.Answer`)
```go
// Método de negocio:
- Validate() error
```
**Solución:** Validación en `AssessmentAttemptService` antes de crear

---

**Entity: Attempt** (`entities.Attempt`)
```go
// Métodos de negocio que deben moverse:
- IsPassed(passThreshold int) bool
- GetCorrectAnswersCount() int
- GetIncorrectAnswersCount() int
- GetTotalQuestions() int
- GetAccuracyPercentage() int
- GetAverageTimePerQuestion() int
- Validate() error
```
**Solución:** Crear helper/calculator en `AssessmentAttemptService` o DTOs

---

### 2. Uso de Value Objects vs UUIDs

**Problema:**
- Entities en `entity/` usan value objects: `valueobject.MaterialID`, `valueobject.UserID`
- Entities en `entities/` usan `uuid.UUID` directamente
- Infrastructure entities usarán **tipos nativos de Go** (`uuid.UUID`)

**Solución:**
- Adaptar en repositories: convertir entre value objects y UUIDs
- Ejemplo: `materialID.UUID()` → `uuid.UUID`

---

### 3. Encapsulación (Fields Privados vs Públicos)

**Problema:**
- `entity.Material` tiene fields privados + getters (DDD estricto)
- `entities.Assessment` tiene fields públicos (anémico)
- Infrastructure entities tendrán **fields públicos con tags** (para GORM/serialización)

**Solución:**
- Repositories deben mapear explícitamente
- Domain layer puede wrappear entities en domain objects si necesita encapsulación

---

### 4. Tests de Entities

**Problema:**
- Existen tests unitarios de entities:
  - `internal/domain/entities/assessment_test.go`
  - `internal/domain/entities/answer_test.go`
  - `internal/domain/entities/attempt_test.go`
- Estos tests verifican **lógica de negocio**

**Solución:**
- ❌ Eliminar tests de entities (la lógica ya no estará ahí)
- ✅ Crear tests de domain services que contendrán la lógica movida
- ✅ Mantener tests de repositories (cambiarán imports pero lógica igual)

---

## 📋 Tareas del Sprint

### Etapa 0: Preparación

#### Tarea 0.1: Verificar que infrastructure entities existen

```bash
# Clonar o actualizar edugo-infrastructure
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Verificar que entities fueron creadas en Sprint anterior
ls -la postgres/entities/
ls -la mongodb/entities/

# Verificar que hay releases con tags
git tag | grep entities
```

**Criterio de éxito:**
- Existen 14 archivos en `postgres/entities/`
- Existen tags `postgres/entities/v0.1.0` y `mongodb/entities/v0.1.0`

---

### Etapa 1: Actualizar go.mod de api-mobile

#### Tarea 1.1: Agregar dependencia de infrastructure

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Agregar dependency (ajustar versión según tag real)
go get github.com/EduGoGroup/edugo-infrastructure/postgres/entities@postgres/entities/v0.1.0
go get github.com/EduGoGroup/edugo-infrastructure/mongodb/entities@mongodb/entities/v0.1.0

go mod tidy
```

**Criterio de éxito:** `go.mod` contiene las nuevas dependencias

---

### Etapa 2: Crear Domain Services para Lógica de Negocio

**IMPORTANTE:** Antes de eliminar entities, debemos extraer la lógica de negocio.

#### Tarea 2.1: Crear MaterialDomainService (para lógica de Material)

**Ubicación:** `internal/domain/services/material_domain_service.go`

```go
package services

import (
    "time"
    pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
    "github.com/EduGoGroup/edugo-shared/common/errors"
    "github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// MaterialDomainService contiene reglas de negocio de Material
type MaterialDomainService struct{}

func NewMaterialDomainService() *MaterialDomainService {
    return &MaterialDomainService{}
}

// SetS3Info establece información de S3 (extraído de entity.Material)
func (s *MaterialDomainService) SetS3Info(material *pgentities.Material, s3Key, s3URL string) error {
    if s3Key == "" || s3URL == "" {
        return errors.NewValidationError("s3_key and s3_url are required")
    }
    material.S3Key = s3Key
    material.S3URL = s3URL
    material.ProcessingStatus = enum.ProcessingStatusProcessing
    material.UpdatedAt = time.Now()
    return nil
}

// MarkProcessingComplete marca procesamiento completado
func (s *MaterialDomainService) MarkProcessingComplete(material *pgentities.Material) error {
    if material.ProcessingStatus == enum.ProcessingStatusCompleted {
        return errors.NewBusinessRuleError("material already processed")
    }
    material.ProcessingStatus = enum.ProcessingStatusCompleted
    material.UpdatedAt = time.Now()
    return nil
}

// Publish publica el material
func (s *MaterialDomainService) Publish(material *pgentities.Material) error {
    if material.Status == enum.MaterialStatusPublished {
        return errors.NewBusinessRuleError("material is already published")
    }
    if material.ProcessingStatus != enum.ProcessingStatusCompleted {
        return errors.NewBusinessRuleError("material must be processed before publishing")
    }
    material.Status = enum.MaterialStatusPublished
    material.UpdatedAt = time.Now()
    return nil
}

// Archive archiva el material
func (s *MaterialDomainService) Archive(material *pgentities.Material) error {
    if material.Status == enum.MaterialStatusArchived {
        return errors.NewBusinessRuleError("material is already archived")
    }
    material.Status = enum.MaterialStatusArchived
    material.UpdatedAt = time.Now()
    return nil
}

// Query helpers
func (s *MaterialDomainService) IsDraft(material *pgentities.Material) bool {
    return material.Status == enum.MaterialStatusDraft
}

func (s *MaterialDomainService) IsPublished(material *pgentities.Material) bool {
    return material.Status == enum.MaterialStatusPublished
}

func (s *MaterialDomainService) IsProcessed(material *pgentities.Material) bool {
    return material.ProcessingStatus == enum.ProcessingStatusCompleted
}
```

**Criterio de éxito:** Service creado con toda la lógica de `entity.Material`

---

#### Tarea 2.2: Crear ProgressDomainService (para lógica de Progress)

**Ubicación:** `internal/domain/services/progress_domain_service.go`

```go
package services

import (
    "time"
    pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
    "github.com/EduGoGroup/edugo-shared/common/errors"
    "github.com/EduGoGroup/edugo-shared/common/types/enum"
)

type ProgressDomainService struct{}

func NewProgressDomainService() *ProgressDomainService {
    return &ProgressDomainService{}
}

// UpdateProgress actualiza el progreso con validaciones
func (s *ProgressDomainService) UpdateProgress(progress *pgentities.Progress, percentage, lastPage int) error {
    if percentage < 0 || percentage > 100 {
        return errors.NewValidationError("percentage must be between 0 and 100")
    }

    progress.Percentage = percentage
    progress.LastPage = lastPage
    progress.LastAccessedAt = time.Now()
    progress.UpdatedAt = time.Now()

    // Business rule: determinar status según percentage
    if percentage == 0 {
        progress.Status = enum.ProgressStatusNotStarted
    } else if percentage >= 100 {
        progress.Status = enum.ProgressStatusCompleted
    } else {
        progress.Status = enum.ProgressStatusInProgress
    }

    return nil
}
```

**Criterio de éxito:** Service creado con lógica de `entity.Progress`

---

#### Tarea 2.3: Crear AssessmentDomainService (para lógica de Assessment)

**Ubicación:** `internal/domain/services/assessment_domain_service.go`

```go
package services

import (
    "time"
    "github.com/google/uuid"
    pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
    domainErrors "github.com/EduGoGroup/edugo-api-mobile/internal/domain/errors"
)

type AssessmentDomainService struct{}

func NewAssessmentDomainService() *AssessmentDomainService {
    return &AssessmentDomainService{}
}

// ValidateAssessment valida entity
func (s *AssessmentDomainService) ValidateAssessment(a *pgentities.Assessment) error {
    if a.ID == uuid.Nil {
        return domainErrors.ErrInvalidAssessmentID
    }
    if a.MaterialID == uuid.Nil {
        return domainErrors.ErrInvalidMaterialID
    }
    if len(a.MongoDocumentID) != 24 {
        return domainErrors.ErrInvalidMongoDocumentID
    }
    if a.Title == "" {
        return domainErrors.ErrEmptyTitle
    }
    if a.TotalQuestions < 1 || a.TotalQuestions > 100 {
        return domainErrors.ErrInvalidTotalQuestions
    }
    if a.PassThreshold < 0 || a.PassThreshold > 100 {
        return domainErrors.ErrInvalidPassThreshold
    }
    if a.MaxAttempts != nil && *a.MaxAttempts < 1 {
        return domainErrors.ErrInvalidMaxAttempts
    }
    if a.TimeLimitMinutes != nil && (*a.TimeLimitMinutes < 1 || *a.TimeLimitMinutes > 180) {
        return domainErrors.ErrInvalidTimeLimit
    }
    return nil
}

// CanAttempt verifica si puede hacer otro intento
func (s *AssessmentDomainService) CanAttempt(assessment *pgentities.Assessment, attemptCount int) bool {
    if assessment.MaxAttempts == nil {
        return true // Ilimitado
    }
    return attemptCount < *assessment.MaxAttempts
}

// IsTimeLimited indica si tiene límite de tiempo
func (s *AssessmentDomainService) IsTimeLimited(assessment *pgentities.Assessment) bool {
    return assessment.TimeLimitMinutes != nil && *assessment.TimeLimitMinutes > 0
}

// SetMaxAttempts establece máximo de intentos
func (s *AssessmentDomainService) SetMaxAttempts(assessment *pgentities.Assessment, max int) error {
    if max < 1 {
        return domainErrors.ErrInvalidMaxAttempts
    }
    assessment.MaxAttempts = &max
    assessment.UpdatedAt = time.Now().UTC()
    return nil
}

// SetTimeLimit establece límite de tiempo
func (s *AssessmentDomainService) SetTimeLimit(assessment *pgentities.Assessment, minutes int) error {
    if minutes < 1 || minutes > 180 {
        return domainErrors.ErrInvalidTimeLimit
    }
    assessment.TimeLimitMinutes = &minutes
    assessment.UpdatedAt = time.Now().UTC()
    return nil
}

// RemoveMaxAttempts quita límite de intentos
func (s *AssessmentDomainService) RemoveMaxAttempts(assessment *pgentities.Assessment) {
    assessment.MaxAttempts = nil
    assessment.UpdatedAt = time.Now().UTC()
}

// RemoveTimeLimit quita límite de tiempo
func (s *AssessmentDomainService) RemoveTimeLimit(assessment *pgentities.Assessment) {
    assessment.TimeLimitMinutes = nil
    assessment.UpdatedAt = time.Now().UTC()
}
```

**Criterio de éxito:** Service creado con lógica de `entities.Assessment`

---

#### Tarea 2.4: Crear AttemptDomainService (para lógica de Attempt)

**Ubicación:** `internal/domain/services/attempt_domain_service.go`

```go
package services

import (
    "errors"
    "time"
    "github.com/google/uuid"
    pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
    domainErrors "github.com/EduGoGroup/edugo-api-mobile/internal/domain/errors"
)

type AttemptDomainService struct{}

func NewAttemptDomainService() *AttemptDomainService {
    return &AttemptDomainService{}
}

// IsPassed indica si aprobó
func (s *AttemptDomainService) IsPassed(attempt *pgentities.AssessmentAttempt, passThreshold int) bool {
    return attempt.Score >= passThreshold
}

// GetCorrectAnswersCount cuenta respuestas correctas
func (s *AttemptDomainService) GetCorrectAnswersCount(answers []*pgentities.AssessmentAnswer) int {
    count := 0
    for _, answer := range answers {
        if answer.IsCorrect {
            count++
        }
    }
    return count
}

// GetIncorrectAnswersCount cuenta respuestas incorrectas
func (s *AttemptDomainService) GetIncorrectAnswersCount(answers []*pgentities.AssessmentAnswer) int {
    return len(answers) - s.GetCorrectAnswersCount(answers)
}

// GetAverageTimePerQuestion calcula tiempo promedio
func (s *AttemptDomainService) GetAverageTimePerQuestion(attempt *pgentities.AssessmentAttempt, totalQuestions int) int {
    if totalQuestions == 0 {
        return 0
    }
    return attempt.TimeSpentSeconds / totalQuestions
}

// ValidateAttempt valida intento
func (s *AttemptDomainService) ValidateAttempt(attempt *pgentities.AssessmentAttempt, answers []*pgentities.AssessmentAnswer) error {
    if attempt.ID == uuid.Nil {
        return domainErrors.ErrInvalidAttemptID
    }
    if attempt.AssessmentID == uuid.Nil {
        return domainErrors.ErrInvalidAssessmentID
    }
    if attempt.StudentID == uuid.Nil {
        return domainErrors.ErrInvalidStudentID
    }
    if attempt.Score < 0 || attempt.Score > 100 {
        return domainErrors.ErrInvalidScore
    }
    if attempt.TimeSpentSeconds <= 0 || attempt.TimeSpentSeconds > 7200 {
        return domainErrors.ErrInvalidTimeSpent
    }
    if attempt.StartedAt.IsZero() {
        return domainErrors.ErrInvalidStartTime
    }
    if attempt.CompletedAt.IsZero() || !attempt.CompletedAt.After(attempt.StartedAt) {
        return domainErrors.ErrInvalidEndTime
    }
    if len(answers) == 0 {
        return domainErrors.ErrNoAnswersProvided
    }

    // Verificar que el score calculado coincide
    correctCount := s.GetCorrectAnswersCount(answers)
    expectedScore := (correctCount * 100) / len(answers)
    if attempt.Score != expectedScore {
        return errors.New("domain: score mismatch with answers")
    }

    return nil
}

// ValidateAnswer valida respuesta individual
func (s *AttemptDomainService) ValidateAnswer(answer *pgentities.AssessmentAnswer) error {
    if answer.ID == uuid.Nil {
        return domainErrors.ErrInvalidAnswerID
    }
    if answer.AttemptID == uuid.Nil {
        return domainErrors.ErrInvalidAttemptID
    }
    if answer.QuestionID == "" {
        return domainErrors.ErrInvalidQuestionID
    }
    if answer.SelectedAnswerID == "" {
        return domainErrors.ErrInvalidSelectedAnswerID
    }
    if answer.TimeSpentSeconds < 0 {
        return domainErrors.ErrInvalidTimeSpent
    }
    return nil
}
```

**Criterio de éxito:** Service creado con lógica de `entities.Attempt` y `entities.Answer`

---

### Etapa 3: Actualizar Imports en Archivos Existentes

#### Tarea 3.1: Actualizar imports en repositories (PostgreSQL entities)

**Buscar y reemplazar en todos los archivos:**

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Reemplazar import de domain/entity
find internal/ -name "*.go" -type f -exec sed -i '' \
  's|"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"|pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"|g' {} \;

# Reemplazar import de domain/entities
find internal/ -name "*.go" -type f -exec sed -i '' \
  's|"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entities"|pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"|g' {} \;
```

**Nota:** Esto cambiará los imports pero necesitarás ajustar el código manualmente:
- `entity.Material` → `pgentities.Material`
- `entities.Assessment` → `pgentities.Assessment`

---

#### Tarea 3.2: Actualizar referencias de tipos en el código

**Archivos críticos a revisar manualmente:**

1. **Repositories:**
   - `internal/infrastructure/persistence/postgres/repository/*.go`
   - Cambiar `entity.Material` → `pgentities.Material`
   - Cambiar `entities.Assessment` → `pgentities.Assessment`
   - Eliminar llamadas a métodos de negocio (ya no existen en entities)

2. **Domain Repository Interfaces:**
   - `internal/domain/repository/*.go`
   - `internal/domain/repositories/*.go`
   - Actualizar signatures de métodos

3. **Application Services:**
   - `internal/application/service/*.go`
   - Inyectar nuevos domain services creados en Etapa 2
   - Reemplazar llamadas a métodos de entity por llamadas a domain service

**Ejemplo de cambio en MaterialService:**

**Antes:**
```go
// En MaterialService
material, err := entity.NewMaterial(title, desc, authorID, subjectID)
if err != nil {
    return nil, err
}

// ...más tarde
err = material.SetS3Info(s3Key, s3URL)
if err != nil {
    return nil, err
}
```

**Después:**
```go
// En MaterialService (inyectar MaterialDomainService)
type MaterialService struct {
    repo        repository.MaterialRepository
    domainSvc   *services.MaterialDomainService  // NUEVO
}

// En método
material := &pgentities.Material{
    ID:          uuid.New(),
    Title:       title,
    Description: desc,
    AuthorID:    authorID.UUID(), // Convertir value object a UUID
    SubjectID:   subjectID,
    Status:      enum.MaterialStatusDraft,
    CreatedAt:   time.Now(),
    UpdatedAt:   time.Now(),
}

// ...más tarde
err = s.domainSvc.SetS3Info(material, s3Key, s3URL)
if err != nil {
    return nil, err
}
```

---

#### Tarea 3.3: Actualizar DTOs

**Archivo:** `internal/application/dto/material_dto.go`

**Antes:**
```go
func ToMaterialResponse(material *entity.Material) *MaterialResponse {
    return &MaterialResponse{
        ID:          material.ID().String(),
        Title:       material.Title(),
        Description: material.Description(),
        // ...
    }
}
```

**Después:**
```go
func ToMaterialResponse(material *pgentities.Material) *MaterialResponse {
    return &MaterialResponse{
        ID:          material.ID.String(),
        Title:       material.Title,
        Description: material.Description,
        // ...
    }
}
```

**Nota:** Cambiar de getters a acceso directo de fields públicos.

---

### Etapa 4: Eliminar Entities Locales

#### Tarea 4.1: Eliminar entities antiguos

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Eliminar entities de entity/
rm -rf internal/domain/entity/

# Eliminar entities de entities/ (pero conservar tests temporalmente)
rm internal/domain/entities/assessment.go
rm internal/domain/entities/answer.go
rm internal/domain/entities/attempt.go
```

**Criterio de éxito:** Entities eliminados, carpetas vacías o removidas

---

#### Tarea 4.2: Eliminar tests de entities

```bash
# Eliminar tests de entities (la lógica ya no está en entities)
rm internal/domain/entities/assessment_test.go
rm internal/domain/entities/answer_test.go
rm internal/domain/entities/attempt_test.go

# Si el directorio entities/ quedó vacío, eliminarlo
rmdir internal/domain/entities/ 2>/dev/null || true
```

**Criterio de éxito:** Tests de entities eliminados

---

### Etapa 5: Crear Tests de Domain Services

#### Tarea 5.1: Tests de MaterialDomainService

**Ubicación:** `internal/domain/services/material_domain_service_test.go`

Crear tests unitarios que verifiquen:
- ✅ `SetS3Info` valida parámetros
- ✅ `MarkProcessingComplete` cambia status correctamente
- ✅ `Publish` solo funciona si está procesado
- ✅ `Archive` funciona correctamente

**Criterio de éxito:** Tests pasan

---

#### Tarea 5.2: Tests de ProgressDomainService

**Ubicación:** `internal/domain/services/progress_domain_service_test.go`

Crear tests que verifiquen:
- ✅ `UpdateProgress` valida rangos (0-100)
- ✅ Status cambia según percentage

**Criterio de éxito:** Tests pasan

---

#### Tarea 5.3: Tests de AssessmentDomainService

**Ubicación:** `internal/domain/services/assessment_domain_service_test.go`

Migrar lógica de `internal/domain/entities/assessment_test.go` a estos nuevos tests.

**Criterio de éxito:** Tests migrados y pasan

---

#### Tarea 5.4: Tests de AttemptDomainService

**Ubicación:** `internal/domain/services/attempt_domain_service_test.go`

Migrar lógica de `internal/domain/entities/attempt_test.go` y `answer_test.go`.

**Criterio de éxito:** Tests migrados y pasan

---

### Etapa 6: Actualizar Tests de Repositories

#### Tarea 6.1: Actualizar repository tests

**Archivos afectados:**
```
internal/infrastructure/persistence/postgres/repository/material_repository_impl_test.go
internal/infrastructure/persistence/postgres/repository/user_repository_impl_test.go
internal/infrastructure/persistence/postgres/repository/progress_repository_impl_test.go
internal/infrastructure/persistence/postgres/repository/assessment_repository_test.go
internal/infrastructure/persistence/postgres/repository/assessment_repository_integration_test.go
internal/infrastructure/persistence/postgres/repository/answer_repository_test.go
internal/infrastructure/persistence/postgres/repository/answer_repository_integration_test.go
internal/infrastructure/persistence/postgres/repository/attempt_repository_test.go
internal/infrastructure/persistence/postgres/repository/attempt_repository_integration_test.go
```

**Cambios necesarios:**
1. Actualizar imports a `pgentities`
2. Cambiar referencias de tipos
3. Cambiar de getters a acceso directo de fields
4. Ajustar constructores (entities de infrastructure no tienen NewMaterial, etc.)

**Ejemplo:**

**Antes:**
```go
material := entity.NewMaterial(title, desc, authorID, subjectID)
```

**Después:**
```go
material := &pgentities.Material{
    ID:          uuid.New(),
    Title:       title,
    Description: desc,
    AuthorID:    authorID.UUID(),
    SubjectID:   subjectID,
    Status:      enum.MaterialStatusDraft,
    CreatedAt:   time.Now(),
    UpdatedAt:   time.Now(),
}
```

**Criterio de éxito:** Todos los tests de repositories compilan y pasan

---

### Etapa 7: Actualizar Tests de Application Services

#### Tarea 7.1: Actualizar service tests

**Archivos afectados:**
```
internal/application/service/auth_service_test.go
internal/application/service/material_service_test.go
internal/application/service/progress_service_test.go
internal/application/service/assessment_attempt_service.go (si tiene tests)
```

**Cambios:**
1. Actualizar imports
2. Actualizar mocks para trabajar con nuevos types
3. Ajustar assertions (cambiar de `material.Title()` a `material.Title`)

**Criterio de éxito:** Tests de services pasan

---

### Etapa 8: Validación Final

#### Tarea 8.1: Compilación

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Compilar todo el proyecto
go build ./...
```

**Criterio de éxito:** ✅ Compila sin errores

---

#### Tarea 8.2: Ejecutar tests completos

```bash
# Ejecutar todos los tests
go test ./... -v

# Si tienes tests de integración (requieren Docker):
make test-integration
```

**Criterio de éxito:** ✅ Todos los tests pasan

---

#### Tarea 8.3: Verificar coverage

```bash
go test ./... -cover
```

**Criterio de éxito:** Coverage >= 80% (igual o mejor que antes)

---

#### Tarea 8.4: Ejecutar linter

```bash
golangci-lint run
```

**Criterio de éxito:** Sin warnings críticos

---

### Etapa 9: Documentación

#### Tarea 9.1: Actualizar README

Actualizar `README.md` para mencionar:
- Dependencies de infrastructure
- Uso de entities centralizadas
- Nuevos domain services

---

#### Tarea 9.2: Crear documento de migración

**Ubicación:** `docs/MIGRATION_ENTITIES_TO_INFRASTRUCTURE.md`

Documentar:
- Qué cambió
- Por qué se hizo
- Cómo trabajar con los nuevos entities
- Cómo agregar lógica de negocio (en domain services, no en entities)

**Criterio de éxito:** Documento creado

---

## 📊 Estimación de Esfuerzo

| Etapa | Tareas | Tiempo Estimado |
|------|--------|-----------------|
| Etapa 0: Preparación | Verificar infrastructure | 10 min |
| Etapa 1: go.mod | Actualizar dependencies | 10 min |
| Etapa 2: Domain Services | Crear 4 services | 3-4 horas |
| Etapa 3: Actualizar Imports | 31 archivos afectados | 2-3 horas |
| Etapa 4: Eliminar Entities | rm archivos | 5 min |
| Etapa 5: Tests Domain Services | 4 test suites | 2-3 horas |
| Etapa 6: Tests Repositories | 9 archivos | 2 horas |
| Etapa 7: Tests Services | 4 archivos | 1 hora |
| Etapa 8: Validación | Build + tests | 30 min |
| Etapa 9: Documentación | README + doc | 1 hora |
| **TOTAL** | | **12-15 horas** |

---

## 🔗 Dependencias

**Antes de este sprint:**
- ✅ infrastructure Sprint ENTITIES completado
- ✅ Entities de PostgreSQL disponibles en infrastructure
- ✅ Tags de release creados en infrastructure

**Después de este sprint:**
- ➡️ api-mobile usa entities centralizadas
- ➡️ Lógica de negocio separada en domain services
- ➡️ Preparado para futuros cambios en schema de BD (solo actualizar infrastructure)

---

## ⚠️ Notas Importantes

### 1. Value Objects vs UUIDs

**Problema:** api-mobile usa value objects (`valueobject.MaterialID`), infrastructure entities usan `uuid.UUID`.

**Solución:**
```go
// En repositories, convertir:
materialID := valueobject.MaterialID{UUID: material.ID}  // BD → VO
material.ID = materialID.UUID()                          // VO → BD
```

---

### 2. Encapsulación

Infrastructure entities tienen **fields públicos** para facilitar serialización (GORM, JSON).

Si necesitas encapsulación en domain layer:
- Opción A: Wrappear entity en un domain object con métodos
- Opción B: Usar domain services (recomendado)

---

### 3. Lógica de Negocio

**REGLA:** Entities en infrastructure NO tienen lógica de negocio.

**Dónde poner lógica:**
- ✅ **Domain Services** (validaciones complejas, reglas de negocio)
- ✅ **Application Services** (orquestación, casos de uso)
- ✅ **Value Objects** (validaciones de formato)
- ❌ **Entities** (solo estructura de datos)

---

### 4. Constructores

Infrastructure entities NO tienen constructores `New*()`.

**En su lugar:**
```go
// Crear entity manualmente
material := &pgentities.Material{
    ID:          uuid.New(),
    Title:       title,
    Description: desc,
    // ...campos requeridos
    CreatedAt:   time.Now(),
    UpdatedAt:   time.Now(),
}

// Validar con domain service
if err := materialDomainSvc.Validate(material); err != nil {
    return nil, err
}
```

---

### 5. Tests

**Cambio de estrategia:**
- ❌ Tests de entities (antes validaban lógica embebida)
- ✅ Tests de domain services (validan lógica de negocio)
- ✅ Tests de repositories (mapeo entity ↔ BD)
- ✅ Tests de application services (orquestación)

---

## 📈 Criterios de Éxito del Sprint

- [ ] 7 entities locales eliminados
- [ ] 4 domain services creados con lógica de negocio
- [ ] 31 archivos actualizados con imports de infrastructure
- [ ] Todos los tests pasan (unit + integration)
- [ ] Compilación exitosa sin errores
- [ ] Coverage >= 80%
- [ ] Documentación actualizada
- [ ] go.mod contiene dependencies de infrastructure

---

## 🚀 Siguiente Paso

Una vez completado este sprint:
1. **api-administracion** puede ejecutar su propio sprint de adaptación (similar a este)
2. **worker** puede ejecutar su sprint (más simple, solo MongoDB entities)

---

## 📚 Referencias

- [Sprint ENTITIES de Infrastructure](../02-infrastructure/SPRINT-ENTITIES.md)
- [Documentación de DDD - Domain Services](https://martinfowler.com/bliki/EvansClassification.html)
- [Clean Architecture - Use Cases](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

**Generado por:** Claude Code  
**Fecha:** 20 de Noviembre, 2025
