# ANÁLISIS DE BRECHA: api-mobile vs Infrastructure (VERDAD)

**Fecha:** 22 de Noviembre, 2025  
**Infrastructure Versión:** `postgres/v0.10.0` (main branch - FUENTE DE VERDAD)  
**api-mobile dev:** `postgres/v0.9.0` (desactualizado)  
**Principio:** Infrastructure es la DUEÑA del schema, api-mobile se ADAPTA

---

## 🎯 PRINCIPIO FUNDAMENTAL

> **"Infrastructure tiene la ÚLTIMA PALABRA sobre la estructura de las tablas.  
> api-mobile NO define schemas, solo los CONSUME y se ADAPTA."**

### Regla de Oro

1. ✅ **Infrastructure main = VERDAD ABSOLUTA**
   - Estructura de tablas definida en `postgres/migrations/`
   - Entities reflejan exactamente el schema SQL
   - Tags: `postgres/v0.10.0` es la última versión

2. ✅ **api-mobile se ADAPTA**
   - Importar entities desde infrastructure
   - Ajustar lógica de negocio a campos disponibles
   - Si falta un campo: solicitar a infrastructure para agregarlo

3. ❌ **api-mobile NO puede:**
   - Cambiar estructura de entities
   - Agregar campos no existentes en BD
   - Definir su propio schema

---

## 📊 COMPARACIÓN: Entity Local vs Infrastructure (VERDAD)

### Material Entity

#### 🔴 Entity LOCAL (api-mobile actual - OBSOLETO)

**Ubicación:** `internal/domain/entity/material.go`

```go
type Material struct {
    id               MaterialID              // Value Object
    title            string
    description      string
    authorID         UserID                  // Value Object
    subjectID        string
    s3Key            string                  // ❌ NO EXISTE EN BD
    s3URL            string                  // ❌ NO EXISTE EN BD
    status           enum.MaterialStatus     // Enum custom
    processingStatus enum.ProcessingStatus   // ❌ NO EXISTE EN BD
    createdAt        time.Time
    updatedAt        time.Time
}

// Métodos de negocio embebidos:
func (m *Material) SetS3Info(s3Key, s3URL string) error
func (m *Material) Publish() error
func (m *Material) Archive() error
```

**Problemas:**
- ❌ Usa Value Objects (`MaterialID`, `UserID`) en vez de `uuid.UUID`
- ❌ Campos que NO existen en BD: `s3Key`, `processingStatus`
- ❌ Campos faltantes: `SchoolID`, `UploadedByTeacherID`, `FileType`, `FileSizeBytes`, etc.
- ❌ Lógica de negocio embebida en entity (viola DDD limpio)

---

#### ✅ Entity INFRASTRUCTURE (postgres/v0.10.0 - VERDAD)

**Ubicación:** `github.com/EduGoGroup/edugo-infrastructure/postgres/entities/material.go`

```go
type Material struct {
    ID                     uuid.UUID  `db:"id"`
    SchoolID               uuid.UUID  `db:"school_id"`                // ✅ REQUERIDO
    UploadedByTeacherID    uuid.UUID  `db:"uploaded_by_teacher_id"`   // ✅ REQUERIDO (era AuthorID)
    AcademicUnitID         *uuid.UUID `db:"academic_unit_id"`         // ✅ Nullable
    Title                  string     `db:"title"`
    Description            *string    `db:"description"`              // ✅ Nullable
    Subject                *string    `db:"subject"`                  // ✅ Nullable
    Grade                  *string    `db:"grade"`                    // ✅ Nuevo
    FileURL                string     `db:"file_url"`                 // ✅ (era s3URL)
    FileType               string     `db:"file_type"`                // ✅ Nuevo
    FileSizeBytes          int64      `db:"file_size_bytes"`          // ✅ Nuevo
    Status                 string     `db:"status"`                   // ✅ uploaded, processing, ready, failed
    ProcessingStartedAt    *time.Time `db:"processing_started_at"`    // ✅ Nuevo
    ProcessingCompletedAt  *time.Time `db:"processing_completed_at"`  // ✅ Nuevo
    IsPublic               bool       `db:"is_public"`                // ✅ Nuevo
    CreatedAt              time.Time  `db:"created_at"`
    UpdatedAt              time.Time  `db:"updated_at"`
    DeletedAt              *time.Time `db:"deleted_at"`               // ✅ Soft delete
}

func (Material) TableName() string {
    return "materials"
}
```

**Características:**
- ✅ Reflejo EXACTO de `postgres/migrations/005_create_materials.up.sql`
- ✅ Usa tipos nativos Go (`uuid.UUID`, `string`, `int64`)
- ✅ Tags `db:` para sqlx/database/sql
- ✅ Sin lógica de negocio (solo estructura de datos)
- ✅ Campos nullable usan pointers (`*string`, `*uuid.UUID`, `*time.Time`)

---

### 📋 Mapeo de Campos: Local → Infrastructure

| Campo Local | Campo Infrastructure | Tipo Local | Tipo Infrastructure | Acción Requerida |
|-------------|---------------------|------------|---------------------|------------------|
| `id` (MaterialID) | `ID` | ValueObject | `uuid.UUID` | ✅ Convertir VO → UUID |
| `authorID` (UserID) | `UploadedByTeacherID` | ValueObject | `uuid.UUID` | ✅ Renombrar + Convertir |
| `title` | `Title` | `string` | `string` | ✅ OK |
| `description` | `Description` | `string` | `*string` | ⚠️ Nullable |
| `subjectID` | `Subject` | `string` | `*string` | ⚠️ Nullable + Renombrar |
| `s3Key` | ❌ **NO EXISTE** | `string` | - | ❌ **ELIMINAR** |
| `s3URL` | `FileURL` | `string` | `string` | ✅ Renombrar |
| `status` | `Status` | Enum | `string` | ✅ Cambiar a string |
| `processingStatus` | ❌ **NO EXISTE** | Enum | - | ❌ **ELIMINAR** (usar `Status`) |
| ❌ Falta | `SchoolID` | - | `uuid.UUID` | ✅ **AGREGAR** |
| ❌ Falta | `AcademicUnitID` | - | `*uuid.UUID` | ✅ **AGREGAR** |
| ❌ Falta | `Grade` | - | `*string` | ✅ **AGREGAR** |
| ❌ Falta | `FileType` | - | `string` | ✅ **AGREGAR** |
| ❌ Falta | `FileSizeBytes` | - | `int64` | ✅ **AGREGAR** |
| ❌ Falta | `ProcessingStartedAt` | - | `*time.Time` | ✅ **AGREGAR** |
| ❌ Falta | `ProcessingCompletedAt` | - | `*time.Time` | ✅ **AGREGAR** |
| ❌ Falta | `IsPublic` | - | `bool` | ✅ **AGREGAR** |
| ❌ Falta | `DeletedAt` | - | `*time.Time` | ✅ **AGREGAR** |

---

### Progress Entity

#### Entity LOCAL (api-mobile - OBSOLETO)

```go
type Progress struct {
    id             ProgressID    // Composite VO (MaterialID + UserID)
    materialID     MaterialID    // Value Object
    userID         UserID        // Value Object
    percentage     int
    lastPage       int
    status         enum.ProgressStatus  // Enum custom
    lastAccessedAt time.Time
    createdAt      time.Time
    updatedAt      time.Time
}

func (p *Progress) UpdateProgress(percentage, lastPage int) error  // Lógica embebida
```

#### Entity INFRASTRUCTURE (VERDAD)

```go
type Progress struct {
    MaterialID     uuid.UUID `db:"material_id"`  // PK compuesta
    UserID         uuid.UUID `db:"user_id"`      // PK compuesta
    Percentage     int       `db:"percentage"`
    LastPage       int       `db:"last_page"`
    Status         string    `db:"status"`       // not_started, in_progress, completed
    LastAccessedAt time.Time `db:"last_accessed_at"`
    CreatedAt      time.Time `db:"created_at"`
    UpdatedAt      time.Time `db:"updated_at"`
}
```

**Diferencias:**
- ❌ No usa ProgressID composite
- ❌ MaterialID y UserID son `uuid.UUID` directo
- ❌ Status es `string` no enum
- ✅ Sin lógica de negocio embebida

---

### Assessment Entity

#### Entity LOCAL (api-mobile - OBSOLETO)

```go
type Assessment struct {
    ID              uuid.UUID  // ✅ Ya usa UUID (entities/ no entity/)
    MaterialID      uuid.UUID
    MongoDocumentID string
    Title           string              // ❌ Requerido
    TotalQuestions  int
    PassThreshold   int                 // ❌ Requerido
    MaxAttempts     *int
    TimeLimitMinutes *int
    Status          string
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

// Lógica embebida:
func (a *Assessment) Validate() error
func (a *Assessment) CanAttempt(attemptCount int) bool
func (a *Assessment) SetMaxAttempts(max int) error
```

#### Entity INFRASTRUCTURE (VERDAD)

```go
type Assessment struct {
    ID                uuid.UUID  `db:"id"`
    MaterialID        uuid.UUID  `db:"material_id"`
    MongoDocumentID   string     `db:"mongo_document_id"`
    QuestionsCount    int        `db:"questions_count"`
    TotalQuestions    *int       `db:"total_questions"`   // ✅ Nullable
    Title             *string    `db:"title"`             // ✅ Nullable
    PassThreshold     *int       `db:"pass_threshold"`    // ✅ Nullable
    MaxAttempts       *int       `db:"max_attempts"`
    TimeLimitMinutes  *int       `db:"time_limit_minutes"`
    Status            string     `db:"status"`
    CreatedAt         time.Time  `db:"created_at"`
    UpdatedAt         time.Time  `db:"updated_at"`
    DeletedAt         *time.Time `db:"deleted_at"`        // ✅ Nuevo
}
```

**Diferencias:**
- ⚠️ `Title` es nullable en BD (`*string`) no requerido
- ⚠️ `PassThreshold` es nullable (`*int`)
- ✅ Tiene `QuestionsCount` y `TotalQuestions` (sincronizados)
- ✅ Soft delete con `DeletedAt`

---

## 🚨 PROBLEMAS EN DOMAIN SERVICES CREADOS (Fase 1)

El programador anterior creó 4 domain services basándose en **stubs incorrectos**.

### MaterialDomainService - REQUIERE CORRECCIONES

**Código Actual (INCORRECTO):**

```go
func (s *MaterialDomainService) SetS3Info(material *pgentities.Material, s3Key, s3URL string) error {
    material.S3Key = s3Key      // ❌ Campo NO EXISTE en infrastructure
    material.S3URL = s3URL      // ❌ Campo NO EXISTE en infrastructure
    material.ProcessingStatus = enum.ProcessingStatusProcessing  // ❌ Campo NO EXISTE
    return nil
}
```

**Código Correcto (ADAPTADO A INFRASTRUCTURE):**

```go
func (s *MaterialDomainService) SetFileInfo(material *pgentities.Material, fileURL string, fileType string, fileSizeBytes int64) error {
    if fileURL == "" {
        return errors.NewValidationError("file_url is required")
    }

    material.FileURL = fileURL           // ✅ Campo real
    material.FileType = fileType         // ✅ Campo real
    material.FileSizeBytes = fileSizeBytes  // ✅ Campo real
    material.Status = "processing"       // ✅ Usa Status en vez de ProcessingStatus
    now := time.Now()
    material.ProcessingStartedAt = &now  // ✅ Campo real
    material.UpdatedAt = now

    return nil
}

func (s *MaterialDomainService) MarkProcessingComplete(material *pgentities.Material) error {
    if material.Status == "ready" {
        return errors.NewBusinessRuleError("material already processed")
    }

    material.Status = "ready"  // ✅ Estado correcto según migration
    now := time.Now()
    material.ProcessingCompletedAt = &now  // ✅ Campo real
    material.UpdatedAt = now

    return nil
}

func (s *MaterialDomainService) Publish(material *pgentities.Material) error {
    if material.Status != "ready" {
        return errors.NewBusinessRuleError("material must be processed before publishing")
    }

    material.IsPublic = true  // ✅ Campo real para publicación
    material.UpdatedAt = time.Now()

    return nil
}

func (s *MaterialDomainService) Archive(material *pgentities.Material) error {
    // Soft delete
    now := time.Now()
    material.DeletedAt = &now  // ✅ Campo real
    material.UpdatedAt = now

    return nil
}
```

---

### ProgressDomainService - REQUIERE CORRECCIONES MENORES

**Código Actual:**

```go
func (s *ProgressDomainService) UpdateProgress(progress *pgentities.Progress, percentage, lastPage int) error {
    progress.Status = enum.ProgressStatusInProgress  // ❌ Usa enum
    // ...
}
```

**Código Correcto:**

```go
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
        progress.Status = "not_started"  // ✅ String según migration
    } else if percentage >= 100 {
        progress.Status = "completed"
    } else {
        progress.Status = "in_progress"
    }

    return nil
}
```

---

### AssessmentDomainService - REQUIERE VALIDACIONES AJUSTADAS

**Ajustar validaciones para campos nullable:**

```go
func (s *AssessmentDomainService) ValidateAssessment(a *pgentities.Assessment) error {
    if a.ID == uuid.Nil {
        return errors.NewValidationError("invalid assessment id")
    }
    if a.MaterialID == uuid.Nil {
        return errors.NewValidationError("invalid material id")
    }
    if a.MongoDocumentID == "" {
        return errors.NewValidationError("invalid mongo document id")
    }

    // Title es opcional (nullable)
    // PassThreshold es opcional (nullable)

    if a.QuestionsCount < 1 {
        return errors.NewValidationError("questions count must be >= 1")
    }

    if a.PassThreshold != nil && (*a.PassThreshold < 0 || *a.PassThreshold > 100) {
        return errors.NewValidationError("pass threshold must be between 0 and 100")
    }

    if a.MaxAttempts != nil && *a.MaxAttempts < 1 {
        return errors.NewValidationError("max attempts must be >= 1")
    }

    return nil
}
```

---

## 📋 CAMPOS FALTANTES EN api-mobile

Si api-mobile necesita campos que **NO EXISTEN** en infrastructure, debe:

1. ❌ **NO agregarlos localmente** (viola principio)
2. ✅ **Solicitar a infrastructure** crear migration
3. ✅ **Esperar nueva versión** de infrastructure
4. ✅ **Actualizar go.mod** cuando esté disponible

### Ejemplo: Si necesitáramos campo `OriginalFileName`

```bash
# 1. Crear issue en infrastructure repo
gh issue create --repo EduGoGroup/edugo-infrastructure \
  --title "feat(postgres): agregar campo original_filename a materials" \
  --body "Necesitamos almacenar el nombre original del archivo subido..."

# 2. Esperar que se cree migration y nueva versión

# 3. Actualizar api-mobile cuando esté disponible
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.11.0
```

---

## ✅ VERSIÓN CORRECTA A USAR

### Versión Actual en api-mobile dev

```
github.com/EduGoGroup/edugo-infrastructure/postgres v0.9.0
```

### Última Versión Disponible en Infrastructure

```
github.com/EduGoGroup/edugo-infrastructure/postgres v0.10.0  (main)
```

### Recomendación

✅ **Actualizar a `postgres/v0.10.0`** (última versión estable en main)

```bash
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.10.0
go mod tidy
```

---

## 🎯 PLAN DE ADAPTACIÓN CORRECTO

### Fase 1: Actualizar go.mod a VERDAD

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Actualizar a última versión de infrastructure (VERDAD)
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.10.0
go mod tidy
```

### Fase 2: Corregir Domain Services

Ajustar los 4 domain services para usar campos REALES de infrastructure:

1. `material_domain_service.go` - Cambios mayores (SetFileInfo, Status, etc.)
2. `progress_domain_service.go` - Cambios menores (string en vez de enum)
3. `assessment_domain_service.go` - Ajustar validaciones nullable
4. `attempt_domain_service.go` - Verificar campos

### Fase 3: Eliminar Stubs y Entities Locales

```bash
rm -rf internal/infrastructure_stubs/
rm -rf internal/domain/entity/
rm -rf internal/domain/entities/
```

### Fase 4: Actualizar 31 Archivos

Seguir plan de ejecución pero con **VERDAD de infrastructure**, no stubs.

---

## 📊 RESUMEN DE CAMBIOS CRÍTICOS

| Aspecto | Entity Local | Infrastructure (VERDAD) | Impacto |
|---------|--------------|------------------------|---------|
| **Material.AuthorID** | ✅ Existe | ❌ NO existe (es UploadedByTeacherID) | 🔴 ALTO |
| **Material.S3Key** | ✅ Existe | ❌ NO existe | 🔴 ALTO |
| **Material.ProcessingStatus** | ✅ Existe (enum) | ❌ NO existe (usa Status) | 🔴 ALTO |
| **Material.SchoolID** | ❌ NO existe | ✅ Existe (requerido) | 🔴 ALTO |
| **Material.FileType** | ❌ NO existe | ✅ Existe (requerido) | 🟡 MEDIO |
| **Progress PK** | Composite VO | 2 UUIDs separados | 🟡 MEDIO |
| **Assessment.Title** | string | *string (nullable) | 🟡 MEDIO |
| **Todos los enums** | Custom enums | strings | 🟢 BAJO |

---

## ⚠️ ADVERTENCIA FINAL

> **El programador anterior NO tenía acceso a infrastructure actualizado.**  
> Los stubs creados son **aproximaciones incorrectas**.  
> **TODO el código de Fase 1 debe revisarse contra infrastructure REAL.**

**Siguiente paso:** Seguir `PLAN-EJECUCION-COMPLETA.md` pero usando **entities REALES de infrastructure v0.10.0**, no los stubs.

---

**Generado por:** Claude Code  
**Fecha:** 22 de Noviembre, 2025  
**Fuente de Verdad:** `github.com/EduGoGroup/edugo-infrastructure/postgres@v0.10.0`
