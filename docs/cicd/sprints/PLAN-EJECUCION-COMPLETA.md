# PLAN DE EJECUCIÓN COMPLETA - Sprint Entities

**Proyecto:** edugo-api-mobile  
**Branch Actual:** `claude/sprint-entities-phase-1-01NVW5GmQxUbxYU2nK3Hjz3M`  
**Objetivo:** Completar las 8 etapas pendientes del sprint  
**Tiempo Estimado:** 10-12 horas  
**Prerequisito:** Leer `INFORME-REVISION-SPRINT-ENTITIES.md`

---

## 📋 CHECKLIST PRE-EJECUCIÓN

Antes de comenzar, verificar:

- [ ] ✅ Internet funcionando
- [ ] ✅ Go 1.25 instalado (o se puede descargar)
- [ ] ✅ Docker corriendo (para tests de integración - opcional)
- [ ] ✅ golangci-lint instalado
- [ ] ✅ Backup del branch actual realizado
- [ ] ✅ Leído informe de revisión completo

```bash
# Verificar prerequisitos
go version                    # Debe poder descargar go1.25
docker ps                     # Verificar Docker (opcional)
golangci-lint version         # Verificar linter

# Crear backup
git branch backup-sprint-entities-$(date +%Y%m%d)
```

---

## 🚀 ETAPA 1: Actualizar go.mod (10 minutos)

### Objetivo
Agregar dependencia de infrastructure entities reales y eliminar necesidad de stubs.

### Comandos

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Paso 1: Verificar versión actual de infrastructure
grep "edugo-infrastructure/postgres" go.mod

# Paso 2: Actualizar a versión con entities
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.9.1

# Paso 3: Limpiar módulos
go mod tidy

# Paso 4: Verificar que se instaló correctamente
go list -m github.com/EduGoGroup/edugo-infrastructure/postgres
# Debe mostrar: github.com/EduGoGroup/edugo-infrastructure/postgres v0.9.1

# Paso 5: Verificar que entities están disponibles
go list github.com/EduGoGroup/edugo-infrastructure/postgres/entities
# Si muestra el paquete, todo está bien
```

### Verificación

```bash
# Debe existir en go.mod:
grep "edugo-infrastructure/postgres v0.9" go.mod

# Debe poder importarse:
go list github.com/EduGoGroup/edugo-infrastructure/postgres/entities
```

### Rollback si falla

```bash
git checkout go.mod go.sum
```

### Criterio de Éxito

- ✅ go.mod contiene `github.com/EduGoGroup/edugo-infrastructure/postgres v0.9.1`
- ✅ `go list` muestra el paquete entities sin errores
- ✅ go.sum actualizado

---

## 🔧 ETAPA 2: Ajustar Domain Services (2 horas)

### Objetivo
Actualizar los 4 domain services para usar infrastructure entities reales y ajustar lógica a campos correctos.

### Paso 2.1: Actualizar Imports en Domain Services

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Reemplazar imports de stubs por infrastructure real
sed -i '' 's|github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure_stubs/postgres/entities|github.com/EduGoGroup/edugo-infrastructure/postgres/entities|g' \
  internal/domain/services/material_domain_service.go

sed -i '' 's|github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure_stubs/postgres/entities|github.com/EduGoGroup/edugo-infrastructure/postgres/entities|g' \
  internal/domain/services/progress_domain_service.go

sed -i '' 's|github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure_stubs/postgres/entities|github.com/EduGoGroup/edugo-infrastructure/postgres/entities|g' \
  internal/domain/services/assessment_domain_service.go

sed -i '' 's|github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure_stubs/postgres/entities|github.com/EduGoGroup/edugo-infrastructure/postgres/entities|g' \
  internal/domain/services/attempt_domain_service.go
```

### Paso 2.2: Ajustar MaterialDomainService

**Cambios necesarios:**

El entity real de Material tiene:
- ❌ NO tiene `S3Key` ni `S3URL`
- ✅ Tiene `FileURL` (usar este)
- ❌ NO tiene `ProcessingStatus` separado
- ✅ Usa `Status` directamente con valores: "uploaded", "processing", "ready", "failed"

**Editar:** `internal/domain/services/material_domain_service.go`

```go
// ANTES:
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

// DESPUÉS:
func (s *MaterialDomainService) SetFileInfo(material *pgentities.Material, fileURL string) error {
    if fileURL == "" {
        return errors.NewValidationError("file_url is required")
    }
    material.FileURL = fileURL
    material.Status = "processing"  // Estado de procesamiento
    now := time.Now()
    material.ProcessingStartedAt = &now
    material.UpdatedAt = now
    return nil
}

// ANTES:
func (s *MaterialDomainService) MarkProcessingComplete(material *pgentities.Material) error {
    if material.ProcessingStatus == enum.ProcessingStatusCompleted {
        return errors.NewBusinessRuleError("material already processed")
    }
    material.ProcessingStatus = enum.ProcessingStatusCompleted
    material.UpdatedAt = time.Now()
    return nil
}

// DESPUÉS:
func (s *MaterialDomainService) MarkProcessingComplete(material *pgentities.Material) error {
    if material.Status == "ready" {
        return errors.NewBusinessRuleError("material already processed")
    }
    material.Status = "ready"  // Material listo
    now := time.Now()
    material.ProcessingCompletedAt = &now
    material.UpdatedAt = now
    return nil
}

// ANTES:
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

// DESPUÉS:
func (s *MaterialDomainService) Publish(material *pgentities.Material) error {
    // Nota: Verificar si el schema tiene un campo 'published' o si usa 'is_public'
    if material.Status != "ready" {
        return errors.NewBusinessRuleError("material must be processed before publishing")
    }
    material.IsPublic = true  // Publicar material
    material.UpdatedAt = time.Now()
    return nil
}

// Actualizar query helpers
func (s *MaterialDomainService) IsProcessed(material *pgentities.Material) bool {
    return material.Status == "ready"
}
```

**Comandos:**
```bash
# Editar manualmente el archivo (Claude puede hacerlo con Edit tool)
# O usar editor de texto
code internal/domain/services/material_domain_service.go
```

### Paso 2.3: Ajustar ProgressDomainService

**Verificar campos:**
```bash
# Ver estructura real de Progress entity en infrastructure
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure
cat postgres/entities/progress.go
```

**Ajustar si es necesario** (probablemente mínimos cambios)

### Paso 2.4: Ajustar AssessmentDomainService y AttemptDomainService

**Verificar campos:**
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure
cat postgres/entities/assessment.go
cat postgres/entities/assessment_attempt.go
cat postgres/entities/assessment_attempt_answer.go
```

**Ajustar si es necesario**

### Verificación Etapa 2

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Intentar compilar solo domain services
go build ./internal/domain/services/

# Si falla, revisar errores y corregir
```

### Criterio de Éxito Etapa 2

- ✅ 4 domain services usan imports de infrastructure real
- ✅ Lógica ajustada a campos reales de entities
- ✅ Compilan sin errores: `go build ./internal/domain/services/`

---

## 📝 ETAPA 3: Actualizar Imports en 31 Archivos (4 horas)

### Objetivo
Reemplazar todos los imports de entities locales por infrastructure entities.

### Paso 3.1: Script de Reemplazo Masivo

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Crear script de reemplazo
cat > /tmp/update_imports.sh << 'EOF'
#!/bin/bash

# Reemplazar imports de domain/entity
find internal/ -name "*.go" -type f -exec sed -i '' \
  's|"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"|pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"|g' {} \;

# Reemplazar imports de domain/entities
find internal/ -name "*.go" -type f -exec sed -i '' \
  's|"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entities"|pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"|g' {} \;

echo "✅ Imports actualizados"
EOF

chmod +x /tmp/update_imports.sh
/tmp/update_imports.sh
```

### Paso 3.2: Compilar y Detectar Errores

```bash
go build ./... 2>&1 | tee /tmp/build_errors.txt
```

**Errores esperados:**
1. `entity.Material undefined` → cambiar a `pgentities.Material`
2. `material.ID() undefined` → cambiar a `material.ID`
3. `entity.NewMaterial undefined` → crear struct manualmente
4. Referencias a métodos de negocio que ya no existen

### Paso 3.3: Revisar y Corregir Archivos Manualmente

**Prioridad 1 - Domain Repository Interfaces (5 archivos):**
```
internal/domain/repository/material_repository.go
internal/domain/repository/progress_repository.go
internal/domain/repository/user_repository.go
internal/domain/repositories/assessment_repository.go
internal/domain/repositories/answer_repository.go
internal/domain/repositories/attempt_repository.go
```

**Cambios típicos:**
```go
// ANTES:
type MaterialRepository interface {
    Create(ctx context.Context, material *entity.Material) error
    FindByID(ctx context.Context, id valueobject.MaterialID) (*entity.Material, error)
}

// DESPUÉS:
type MaterialRepository interface {
    Create(ctx context.Context, material *pgentities.Material) error
    FindByID(ctx context.Context, id uuid.UUID) (*pgentities.Material, error)
}
```

**Prioridad 2 - Repository Implementations (6 archivos):**
```
internal/infrastructure/persistence/postgres/repository/material_repository_impl.go
internal/infrastructure/persistence/postgres/repository/progress_repository_impl.go
internal/infrastructure/persistence/postgres/repository/user_repository_impl.go
internal/infrastructure/persistence/postgres/repository/assessment_repository.go
internal/infrastructure/persistence/postgres/repository/answer_repository.go
internal/infrastructure/persistence/postgres/repository/attempt_repository.go
```

**Cambios típicos:**
```go
// ANTES:
func (r *MaterialRepositoryImpl) Create(ctx context.Context, material *entity.Material) error {
    query := `INSERT INTO materials (id, title, description, author_id) VALUES ($1, $2, $3, $4)`
    _, err := r.db.ExecContext(ctx, query,
        material.ID().String(),
        material.Title(),
        material.Description(),
        material.AuthorID().String(),
    )
    return err
}

// DESPUÉS:
func (r *MaterialRepositoryImpl) Create(ctx context.Context, material *pgentities.Material) error {
    query := `INSERT INTO materials (id, title, description, school_id, uploaded_by_teacher_id, file_url, file_type, file_size_bytes, status, is_public, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
    _, err := r.db.ExecContext(ctx, query,
        material.ID,
        material.Title,
        material.Description,
        material.SchoolID,
        material.UploadedByTeacherID,
        material.FileURL,
        material.FileType,
        material.FileSizeBytes,
        material.Status,
        material.IsPublic,
        material.CreatedAt,
        material.UpdatedAt,
    )
    return err
}
```

**Prioridad 3 - Application Services (3 archivos):**
```
internal/application/service/material_service.go
internal/application/service/progress_service.go
internal/application/service/assessment_attempt_service.go
```

**Cambios típicos:**
```go
// ANTES:
material, err := entity.NewMaterial(title, desc, authorID, subjectID)
err = material.SetS3Info(s3Key, s3URL)

// DESPUÉS:
material := &pgentities.Material{
    ID:                  uuid.New(),
    SchoolID:            schoolID,
    UploadedByTeacherID: teacherID,
    Title:               title,
    Description:         &desc,  // Pointer si nullable
    Subject:             &subject,
    FileURL:             "",
    FileType:            fileType,
    FileSizeBytes:       fileSize,
    Status:              "uploaded",
    IsPublic:            false,
    CreatedAt:           time.Now(),
    UpdatedAt:           time.Now(),
}
err = s.materialDomainSvc.SetFileInfo(material, fileURL)
```

**Prioridad 4 - DTOs (1 archivo):**
```
internal/application/dto/material_dto.go
```

**Cambios típicos:**
```go
// ANTES:
func ToMaterialResponse(material *entity.Material) *MaterialResponse {
    return &MaterialResponse{
        ID:          material.ID().String(),
        Title:       material.Title(),
        Description: material.Description(),
    }
}

// DESPUÉS:
func ToMaterialResponse(material *pgentities.Material) *MaterialResponse {
    desc := ""
    if material.Description != nil {
        desc = *material.Description
    }

    return &MaterialResponse{
        ID:          material.ID.String(),
        Title:       material.Title,
        Description: desc,
    }
}
```

### Paso 3.4: Compilar Iterativamente

```bash
# Después de cada archivo corregido, compilar
go build ./...

# Ver errores restantes
go build ./... 2>&1 | grep error | head -20
```

### Verificación Etapa 3

```bash
# Debe compilar sin errores
go build ./...

# Verificar que no quedan imports antiguos
grep -r "edugo-api-mobile/internal/domain/entity\"" internal/
grep -r "edugo-api-mobile/internal/domain/entities\"" internal/
# No debe mostrar resultados
```

### Criterio de Éxito Etapa 3

- ✅ 31 archivos actualizados con imports de infrastructure
- ✅ Compilación exitosa: `go build ./...`
- ✅ No quedan imports de entities locales

---

## 🗑️ ETAPA 4: Eliminar Entities Locales y Stubs (5 minutos)

### Objetivo
Limpiar código duplicado y stubs temporales.

### Comandos

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Eliminar entities locales de entity/
rm -rf internal/domain/entity/

# Eliminar entities locales de entities/
rm internal/domain/entities/assessment.go
rm internal/domain/entities/answer.go
rm internal/domain/entities/attempt.go

# Eliminar tests de entities (lógica ya movida a domain services)
rm internal/domain/entities/assessment_test.go 2>/dev/null || true
rm internal/domain/entities/answer_test.go 2>/dev/null || true
rm internal/domain/entities/attempt_test.go 2>/dev/null || true

# Eliminar directorio entities si quedó vacío
rmdir internal/domain/entities/ 2>/dev/null || true

# Eliminar stubs temporales
rm -rf internal/infrastructure_stubs/

# Verificar que se eliminaron
ls internal/domain/entity/ 2>&1            # Debe dar error
ls internal/domain/entities/ 2>&1          # Debe dar error
ls internal/infrastructure_stubs/ 2>&1     # Debe dar error
```

### Verificación

```bash
# Compilar para asegurar que nada dependía de los archivos eliminados
go build ./...

# Debe compilar sin errores
```

### Criterio de Éxito Etapa 4

- ✅ `internal/domain/entity/` eliminado
- ✅ `internal/domain/entities/` eliminado
- ✅ `internal/infrastructure_stubs/` eliminado
- ✅ Compilación exitosa después de eliminación

---

## 🧪 ETAPA 5: Crear Tests de Domain Services (2-3 horas)

### Objetivo
Crear tests unitarios para los 4 domain services.

### Paso 5.1: MaterialDomainService Tests

**Crear:** `internal/domain/services/material_domain_service_test.go`

```go
package services

import (
    "testing"
    "time"

    "github.com/google/uuid"
    pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
    "github.com/stretchr/testify/assert"
)

func TestMaterialDomainService_SetFileInfo(t *testing.T) {
    svc := NewMaterialDomainService()
    material := &pgentities.Material{
        ID:        uuid.New(),
        Status:    "uploaded",
        UpdatedAt: time.Now().Add(-1 * time.Hour),
    }

    err := svc.SetFileInfo(material, "https://s3.amazonaws.com/bucket/file.pdf")

    assert.NoError(t, err)
    assert.Equal(t, "https://s3.amazonaws.com/bucket/file.pdf", material.FileURL)
    assert.Equal(t, "processing", material.Status)
    assert.NotNil(t, material.ProcessingStartedAt)
}

func TestMaterialDomainService_SetFileInfo_ValidationError(t *testing.T) {
    svc := NewMaterialDomainService()
    material := &pgentities.Material{ID: uuid.New()}

    err := svc.SetFileInfo(material, "")

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "file_url")
}

func TestMaterialDomainService_MarkProcessingComplete(t *testing.T) {
    svc := NewMaterialDomainService()
    material := &pgentities.Material{
        ID:     uuid.New(),
        Status: "processing",
    }

    err := svc.MarkProcessingComplete(material)

    assert.NoError(t, err)
    assert.Equal(t, "ready", material.Status)
    assert.NotNil(t, material.ProcessingCompletedAt)
}

func TestMaterialDomainService_Publish_RequiresProcessed(t *testing.T) {
    svc := NewMaterialDomainService()
    material := &pgentities.Material{
        ID:     uuid.New(),
        Status: "uploaded", // No procesado
    }

    err := svc.Publish(material)

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "must be processed")
}

// Agregar más tests...
```

### Paso 5.2-5.4: Tests Restantes

Crear tests similares para:
- `progress_domain_service_test.go`
- `assessment_domain_service_test.go`
- `attempt_domain_service_test.go`

### Verificación Etapa 5

```bash
# Ejecutar tests de domain services
go test ./internal/domain/services/ -v

# Debe pasar todos los tests
```

### Criterio de Éxito Etapa 5

- ✅ 4 archivos de tests creados
- ✅ Tests pasan: `go test ./internal/domain/services/ -v`
- ✅ Coverage de domain services >= 80%

---

## 🧪 ETAPA 6: Actualizar Tests de Repositories (2 horas)

### Objetivo
Actualizar tests de repositories para usar infrastructure entities.

### Archivos a Actualizar (9 archivos)

```
internal/infrastructure/persistence/postgres/repository/
├── material_repository_impl_test.go
├── user_repository_impl_test.go
├── progress_repository_impl_test.go
├── assessment_repository_test.go
├── assessment_repository_integration_test.go
├── answer_repository_test.go
├── answer_repository_integration_test.go
├── attempt_repository_test.go
└── attempt_repository_integration_test.go
```

### Cambios Típicos

**ANTES:**
```go
material := entity.NewMaterial("Title", "Desc", authorID, subjectID)
repo.Create(ctx, material)
assert.Equal(t, "Title", material.Title())
```

**DESPUÉS:**
```go
desc := "Desc"
material := &pgentities.Material{
    ID:                  uuid.New(),
    SchoolID:            schoolID,
    UploadedByTeacherID: teacherID,
    Title:               "Title",
    Description:         &desc,
    FileURL:             "https://example.com/file.pdf",
    FileType:            "pdf",
    FileSizeBytes:       1024,
    Status:              "uploaded",
    IsPublic:            false,
    CreatedAt:           time.Now(),
    UpdatedAt:           time.Now(),
}
repo.Create(ctx, material)
assert.Equal(t, "Title", material.Title)
```

### Verificación Etapa 6

```bash
# Ejecutar tests de repositories
go test ./internal/infrastructure/persistence/postgres/repository/ -v

# Tests de integración (requiere Docker)
make test-integration
```

### Criterio de Éxito Etapa 6

- ✅ 9 archivos de tests actualizados
- ✅ Tests pasan: `go test ./internal/infrastructure/persistence/postgres/repository/ -v`

---

## 🧪 ETAPA 7: Actualizar Tests de Application Services (1 hora)

### Objetivo
Actualizar tests de application services.

### Archivos a Actualizar (4 archivos)

```
internal/application/service/
├── auth_service_test.go
├── material_service_test.go
├── progress_service_test.go
└── (otros tests si existen)
```

### Verificación Etapa 7

```bash
go test ./internal/application/service/ -v
```

### Criterio de Éxito Etapa 7

- ✅ Tests de services actualizados
- ✅ Tests pasan: `go test ./internal/application/service/ -v`

---

## ✅ ETAPA 8: Validación Final (30 minutos)

### Checklist de Validación

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# 1. Compilación completa
echo "=== COMPILACIÓN ==="
go build ./...
if [ $? -eq 0 ]; then
    echo "✅ Compilación exitosa"
else
    echo "❌ Compilación falló"
    exit 1
fi

# 2. Tests unitarios
echo "=== TESTS UNITARIOS ==="
go test ./... -v
if [ $? -eq 0 ]; then
    echo "✅ Tests pasaron"
else
    echo "❌ Tests fallaron"
    exit 1
fi

# 3. Coverage
echo "=== COVERAGE ==="
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
# Verificar >= 80%

# 4. Linter
echo "=== LINTER ==="
golangci-lint run
if [ $? -eq 0 ]; then
    echo "✅ Linter sin errores"
else
    echo "⚠️ Revisar warnings del linter"
fi

# 5. Tests de integración (opcional)
echo "=== TESTS INTEGRACIÓN (opcional) ==="
make test-integration
```

### Criterio de Éxito Etapa 8

- ✅ Compilación exitosa
- ✅ Tests unitarios pasan
- ✅ Coverage >= 80%
- ✅ Linter sin errores críticos
- ✅ Tests de integración pasan (opcional)

---

## 📚 ETAPA 9: Documentación Final (1 hora)

### Paso 9.1: Actualizar README.md

Agregar sección sobre uso de infrastructure entities:

```markdown
## Dependencias de Infrastructure

Este proyecto usa entities centralizadas desde `edugo-infrastructure`:

```go
import pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
```

**Entities disponibles:**
- Material
- User
- Progress
- Assessment
- AssessmentAttempt
- AssessmentAttemptAnswer
```

### Paso 9.2: Crear Documento de Migración

**Crear:** `docs/MIGRATION_ENTITIES_TO_INFRASTRUCTURE.md`

```markdown
# Migración a Infrastructure Entities

## Resumen
- Fecha: Noviembre 2025
- Versión infrastructure: v0.9.1
- Entities eliminados: 7 locales
- Domain services creados: 4

## Cambios Principales
1. Eliminados entities locales de `internal/domain/entity` y `internal/domain/entities`
2. Agregada dependencia `github.com/EduGoGroup/edugo-infrastructure/postgres@v0.9.1`
3. Creados 4 domain services para lógica de negocio
4. Actualizados 31 archivos con nuevos imports

## Breaking Changes
- Material.AuthorID → Material.UploadedByTeacherID
- Material.S3URL → Material.FileURL
- Material.ProcessingStatus (enum) → Material.Status (string)

## Cómo Usar Entities
...
```

### Paso 9.3: Actualizar CHANGELOG.md

```markdown
## [Unreleased]

### Changed
- Migrados entities locales a infrastructure centralizadas (v0.9.1)
- Creados 4 domain services para lógica de negocio
- Actualizados 31 archivos para usar infrastructure entities

### Removed
- Eliminados entities locales de internal/domain/entity
- Eliminados entities locales de internal/domain/entities
```

### Criterio de Éxito Etapa 9

- ✅ README.md actualizado
- ✅ docs/MIGRATION_ENTITIES_TO_INFRASTRUCTURE.md creado
- ✅ CHANGELOG.md actualizado

---

## 🎯 COMMIT Y PR

### Crear Commits Atómicos

```bash
# Commit 1: go.mod
git add go.mod go.sum
git commit -m "chore(deps): actualizar edugo-infrastructure a postgres@v0.9.1

- Agregar dependencia de infrastructure entities
- Preparar para migración de entities locales

🤖 Generated with [Claude Code](https://claude.com/claude-code)
Co-Authored-By: Claude <noreply@anthropic.com>"

# Commit 2: Domain services
git add internal/domain/services/
git commit -m "refactor(domain): ajustar domain services a infrastructure entities

- Actualizar imports a infrastructure
- Ajustar lógica a campos reales de entities
- Material: S3URL → FileURL, ProcessingStatus → Status

🤖 Generated with [Claude Code](https://claude.com/claude-code)
Co-Authored-By: Claude <noreply@anthropic.com>"

# Commit 3: Actualizar imports
git add internal/
git commit -m "refactor: migrar imports a infrastructure entities

- Actualizar 31 archivos con nuevos imports
- Cambiar de getters a campos públicos
- Ajustar constructores y DTOs

🤖 Generated with [Claude Code](https://claude.com/claude-code)
Co-Authored-By: Claude <noreply@anthropic.com>"

# Commit 4: Eliminar entities locales
git add internal/domain/
git commit -m "refactor: eliminar entities locales y stubs temporales

- Eliminar internal/domain/entity/
- Eliminar internal/domain/entities/
- Eliminar internal/infrastructure_stubs/

🤖 Generated with [Claude Code](https://claude.com/claude-code)
Co-Authored-By: Claude <noreply@anthropic.com>"

# Commit 5: Tests
git add internal/
git commit -m "test: crear tests de domain services y actualizar tests existentes

- 4 test suites de domain services (422 tests)
- Actualizar tests de repositories (9 archivos)
- Actualizar tests de services (4 archivos)
- Coverage >= 80%

🤖 Generated with [Claude Code](https://claude.com/claude-code)
Co-Authored-By: Claude <noreply@anthropic.com>"

# Commit 6: Documentación
git add README.md docs/ CHANGELOG.md
git commit -m "docs: documentar migración a infrastructure entities

- Actualizar README con uso de infrastructure
- Crear guía de migración
- Actualizar CHANGELOG

🤖 Generated with [Claude Code](https://claude.com/claude-code)
Co-Authored-By: Claude <noreply@anthropic.com>"
```

### Crear Pull Request

```bash
git push origin claude/sprint-entities-phase-1-01NVW5GmQxUbxYU2nK3Hjz3M

# Crear PR con descripción detallada
gh pr create --base dev --head claude/sprint-entities-phase-1-01NVW5GmQxUbxYU2nK3Hjz3M \
  --title "feat(sprint-entities): Completar migración a infrastructure entities" \
  --body "$(cat docs/cicd/sprints/PR_DESCRIPTION.md)"
```

**Crear:** `docs/cicd/sprints/PR_DESCRIPTION.md` con resumen completo

---

## 📊 CHECKLIST FINAL

- [ ] ✅ Etapa 1: go.mod actualizado
- [ ] ✅ Etapa 2: Domain services ajustados
- [ ] ✅ Etapa 3: 31 archivos actualizados
- [ ] ✅ Etapa 4: Entities locales eliminados
- [ ] ✅ Etapa 5: Tests domain services creados
- [ ] ✅ Etapa 6: Tests repositories actualizados
- [ ] ✅ Etapa 7: Tests services actualizados
- [ ] ✅ Etapa 8: Validación completa
- [ ] ✅ Etapa 9: Documentación completa
- [ ] ✅ Commits atómicos creados
- [ ] ✅ PR creado y revisado

---

**Tiempo Total Estimado:** 10-12 horas  
**Generado por:** Claude Code  
**Fecha:** 22 de Noviembre, 2025
