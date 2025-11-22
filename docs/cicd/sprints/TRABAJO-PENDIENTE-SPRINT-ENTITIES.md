# TRABAJO PENDIENTE - Sprint Entities (75% Completado)

**Fecha:** 22 de Noviembre, 2025  
**Estado:** 75% completado - Listo para finalizar  
**Tiempo restante estimado:** 30-60 minutos

---

## ✅ TRABAJO COMPLETADO (Excelente Progreso)

### 1. Infrastructure Actualizado ✅
- go.mod actualizado a `postgres/v0.10.0` (VERDAD de infrastructure)
- Dependencia funcionando correctamente

### 2. Domain Services Corregidos ✅
- **MaterialDomainService**: FileURL, Status (uploaded/processing/ready/failed), soft delete
- **ProgressDomainService**: Status como strings (not_started/in_progress/completed)
- **AssessmentDomainService**: Campos nullable manejados correctamente
- **AttemptDomainService**: Score/TimeSpent como *float64/*int
- Todos compilan sin errores ✅

### 3. Limpieza Completada ✅
- Stubs temporales eliminados (`internal/infrastructure_stubs/`)
- Entities locales eliminados (`internal/domain/entity/`, `internal/domain/entities/`)

### 4. DTOs Actualizados ✅
- `material_dto.go`: Completamente refactorizado a estructura REAL
- Campos adaptados: SchoolID, UploadedByTeacherID, FileURL, FileType, etc.
- MaterialVersion adaptado: Title, ContentURL, ChangedBy
- Todos compilan sin errores ✅

### 5. Interfaces de Repositorios Actualizados ✅
- `internal/domain/repository/*.go`: 3 archivos corregidos
- `internal/domain/repositories/*.go`: 3 archivos corregidos
- Nombres de entities actualizados (Answer → AssessmentAttemptAnswer, etc.)
- Todos compilan sin errores ✅

### 6. Documentación Creada ✅
- `INFORME-REVISION-SPRINT-ENTITIES.md`: Análisis completo del sprint
- `PLAN-EJECUCION-COMPLETA.md`: Scripts paso a paso
- `ANALISIS-BRECHA-INFRASTRUCTURE.md`: Comparación Local vs Infrastructure (VERDAD)

---

## ⚠️ TRABAJO PENDIENTE (2 archivos)

### Archivo 1: `answer_repository.go`

**Ubicación:** `internal/infrastructure/persistence/postgres/repository/answer_repository.go`

**Errores a corregir:**

1. **Campos renombrados:**
   ```go
   // ❌ INCORRECTO (campos antiguos)
   QuestionID        string
   SelectedAnswerID  string

   // ✅ CORRECTO (según infrastructure REAL)
   QuestionIndex     int       // 0-based index
   StudentAnswer     *string   // Nullable, puede ser JSON o texto
   ```

2. **Campos nullable - Necesitan punteros:**
   ```go
   // ❌ INCORRECTO
   IsCorrect:        isCorrect,         // bool
   TimeSpentSeconds: timeSpentSecs,     // int

   // ✅ CORRECTO
   IsCorrect:        &isCorrect,        // *bool
   TimeSpentSeconds: &timeSpentSecs,    // *int
   ```

3. **Métodos de validación:**
   ```go
   // ❌ answer.Validate() ya NO existe en entity
   // ✅ Usar AttemptDomainService.ValidateAnswer(answer)
   ```

4. **Eliminación de `*` innecesaria:**
   ```go
   // Línea 65:
   // ❌ &*pgentities.AssessmentAttemptAnswer{...}
   // ✅ &pgentities.AssessmentAttemptAnswer{...}
   ```

**Script de corrección:**
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Reemplazar nombres de campos
sed -i '' 's/QuestionID/QuestionIndex/g' internal/infrastructure/persistence/postgres/repository/answer_repository.go
sed -i '' 's/SelectedAnswerID/StudentAnswer/g' internal/infrastructure/persistence/postgres/repository/answer_repository.go

# Luego revisar manualmente para:
# - Cambiar string → int en QuestionIndex
# - Cambiar string → *string en StudentAnswer
# - Agregar & a IsCorrect y TimeSpentSeconds
# - Eliminar llamadas a .Validate()
# - Quitar * en &*pgentities
```

---

### Archivo 2: `assessment_attempt_service.go`

**Ubicación:** `internal/application/service/assessment_attempt_service.go`

**Errores a corregir:**

1. **Campos nullable - Verificar nil antes de usar:**
   ```go
   // ❌ INCORRECTO
   Title: assessment.Title,           // *string
   TotalQuestions: assessment.TotalQuestions,  // *int

   // ✅ CORRECTO
   Title: getStringValue(assessment.Title),           // "" si nil
   TotalQuestions: getIntValue(assessment.TotalQuestions),  // 0 si nil

   // O con helpers:
   func getStringValue(s *string) string {
       if s == nil { return "" }
       return *s
   }
   ```

2. **Métodos movidos a Domain Services:**
   ```go
   // ❌ assessment.CanAttempt(attemptCount) ya NO existe
   // ✅ Usar AssessmentDomainService
   assessmentSvc := services.NewAssessmentDomainService()
   canAttempt := assessmentSvc.CanAttempt(assessment, attemptCount)
   ```

3. **Constructores eliminados:**
   ```go
   // ❌ pgentities.NewAttempt() NO existe
   // ✅ Crear struct manualmente
   attempt := &pgentities.AssessmentAttempt{
       ID:           uuid.New(),
       AssessmentID: assessmentID,
       StudentID:    studentID,
       StartedAt:    time.Now(),
       Status:       "in_progress",
       CreatedAt:    time.Now(),
       UpdatedAt:    time.Now(),
   }
   ```

4. **Comparación de punteros:**
   ```go
   // ❌ prev.Score > best  (prev.Score es *float64)
   // ✅
   if prev.Score != nil && *prev.Score > float64(best) {
       best = int(*prev.Score)
   }
   ```

5. **Campos eliminados:**
   ```go
   // ❌ attempt.Answers NO existe
   // ✅ Las respuestas son una tabla separada, cargar con:
   answers, err := s.answerRepo.FindByAttemptID(ctx, attemptID)
   ```

---

## 🔧 SCRIPTS DE CORRECCIÓN RÁPIDA

### Script 1: Corregir answer_repository.go

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Paso 1: Reemplazos automáticos
sed -i '' 's/QuestionID:/QuestionIndex:/g' internal/infrastructure/persistence/postgres/repository/answer_repository.go
sed -i '' 's/SelectedAnswerID:/StudentAnswer:/g' internal/infrastructure/persistence/postgres/repository/answer_repository.go
sed -i '' 's/\.QuestionID/.QuestionIndex/g' internal/infrastructure/persistence/postgres/repository/answer_repository.go
sed -i '' 's/\.SelectedAnswerID/.StudentAnswer/g' internal/infrastructure/persistence/postgres/repository.go
sed -i '' 's/&\*pgentities/\&pgentities/g' internal/infrastructure/persistence/postgres/repository/answer_repository.go

# Paso 2: Inyectar AttemptDomainService
# Agregar al struct:
# attemptDomainSvc *services.AttemptDomainService

# Paso 3: Reemplazar answer.Validate() por:
# err := r.attemptDomainSvc.ValidateAnswer(answer)
```

### Script 2: Helpers para assessment_attempt_service.go

Agregar al inicio del archivo:

```go
// Helpers para manejar campos nullable
func getStringValue(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}

func getIntValue(i *int) int {
    if i == nil {
        return 0
    }
    return *i
}

func getFloat64Value(f *float64) float64 {
    if f == nil {
        return 0.0
    }
    return *f
}
```

Luego reemplazar usos directos por estos helpers.

---

## 📊 ESTADO DE COMPILACIÓN ACTUAL

```bash
# Paquetes que compilan: ~28 de 30 (93%)
# Errores restantes: ~25 errores en 2 archivos
# Tiempo estimado: 30-60 minutos

# Verificar:
go build ./...
```

---

## ✅ CRITERIOS DE ÉXITO FINAL

Sprint completado cuando:

- [ ] `answer_repository.go` compila sin errores
- [ ] `assessment_attempt_service.go` compila sin errores
- [ ] `go build ./...` pasa completamente
- [ ] Tests básicos funcionan: `go test ./internal/domain/services/`
- [ ] Commit final creado con mensaje descriptivo

---

## 🎯 PRÓXIMOS PASOS INMEDIATOS

1. **Corregir answer_repository.go** (15-20 min)
   - Usar scripts de corrección rápida
   - Revisar manualmente tipos de campos
   - Inyectar AttemptDomainService para validaciones

2. **Corregir assessment_attempt_service.go** (20-30 min)
   - Agregar helpers para campos nullable
   - Inyectar AssessmentDomainService
   - Reemplazar métodos de entity por domain service
   - Crear structs manualmente en vez de constructores

3. **Validación Final** (10 min)
   - `go build ./...`
   - `go test ./internal/domain/services/`
   - Verificar que no hay imports de entities locales

4. **Commit Final** (5 min)
   - Commit con mensaje descriptivo
   - Push a branch (si autorizado)

---

## 📚 REFERENCIAS RÁPIDAS

### Estructura REAL de AssessmentAttemptAnswer

```go
type AssessmentAttemptAnswer struct {
    ID               uuid.UUID  `db:"id"`
    AttemptID        uuid.UUID  `db:"attempt_id"`
    QuestionIndex    int        `db:"question_index"`     // 0-based
    StudentAnswer    *string    `db:"student_answer"`     // Nullable
    IsCorrect        *bool      `db:"is_correct"`         // Nullable
    PointsEarned     *float64   `db:"points_earned"`      // Nullable
    MaxPoints        *float64   `db:"max_points"`         // Nullable
    TimeSpentSeconds *int       `db:"time_spent_seconds"` // Nullable
    AnsweredAt       time.Time  `db:"answered_at"`
    CreatedAt        time.Time  `db:"created_at"`
    UpdatedAt        time.Time  `db:"updated_at"`
}
```

### Estructura REAL de AssessmentAttempt

```go
type AssessmentAttempt struct {
    ID               uuid.UUID  `db:"id"`
    AssessmentID     uuid.UUID  `db:"assessment_id"`
    StudentID        uuid.UUID  `db:"student_id"`
    StartedAt        time.Time  `db:"started_at"`
    CompletedAt      *time.Time `db:"completed_at"`       // Nullable
    Score            *float64   `db:"score"`              // Nullable
    MaxScore         *float64   `db:"max_score"`          // Nullable
    Percentage       *float64   `db:"percentage"`         // Nullable
    TimeSpentSeconds *int       `db:"time_spent_seconds"` // Nullable
    IdempotencyKey   *string    `db:"idempotency_key"`    // Nullable
    Status           string     `db:"status"`             // in_progress, completed, abandoned
    CreatedAt        time.Time  `db:"created_at"`
    UpdatedAt        time.Time  `db:"updated_at"`
}
```

---

**Generado por:** Claude Code  
**Progreso:** 75% → 100% (falta 25%)  
**Siguiente sesión:** Completar 2 archivos restantes  
**Tiempo estimado:** 30-60 minutos
