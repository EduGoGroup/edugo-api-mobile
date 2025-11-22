# REPORTE SPRINT ENTITIES - FASE 1

**Sprint:** Adaptar api-mobile a Entities Centralizadas
**Fecha Ejecución:** 22 de Noviembre, 2025
**Fase:** 1 (Sin ambiente completo - Uso de stubs/mocks)
**Status:** ✅ COMPLETADO (con hallazgos documentados)
**Próxima Fase:** 2 (Con ambiente completo)

---

## 📊 RESUMEN EJECUTIVO

### ✅ Trabajo Completado en Fase 1

| Fase | Tarea | Status | Notas |
|------|-------|--------|-------|
| Fase 0 | Verificar infrastructure entities | ✅ COMPLETADO | Identificado que existe dependencia pero no acceso |
| Fase 1-STUB | Crear stubs locales (7 entities) | ✅ COMPLETADO | Stubs en `internal/infrastructure_stubs/` |
| Fase 2 | Crear 4 Domain Services | ✅ COMPLETADO | MaterialDomain, ProgressDomain, AssessmentDomain, AttemptDomain |
| Fase 8-Parcial | Validación de compilación | ⚠️ BLOQUEADO | Sin Go 1.25 ni internet |
| Fase 9 | Documentación | ✅ COMPLETADO | Este documento |

### ⏳ Trabajo Pendiente para Fase 2

- Fase 3: Actualizar imports en 31 archivos
- Fase 4: Eliminar entities locales
- Fase 5: Crear tests de Domain Services (4 test suites)
- Fase 6: Actualizar tests de repositories (9 archivos)
- Fase 7: Actualizar tests de application services (4 archivos)
- Fase 8-Completa: Validación completa (compilación + tests + coverage)

---

## 🚨 INCONVENIENTES ENCONTRADOS EN FASE 1

### 1. **Sin Conexión a Internet**

**Problema:**
- No se puede descargar Go 1.25.0
- No se puede descargar dependencias de `go mod download`
- No se puede verificar si `github.com/EduGoGroup/edugo-infrastructure/postgres/entities` existe realmente

**Evidencia:**
```bash
$ go version
go: download go1.25.0: Get "https://storage.googleapis.com/...": dial tcp: lookup storage.googleapis.com on [::1]:53: read: connection refused

$ go list github.com/EduGoGroup/edugo-infrastructure/postgres/entities
ENTITIES_NOT_FOUND
```

**Solución Aplicada en Fase 1:**
✅ Crear stubs locales en `internal/infrastructure_stubs/postgres/entities/`

**Acción para Fase 2:**
1. Eliminar `internal/infrastructure_stubs/` completamente
2. Reemplazar imports de stubs por imports reales de infrastructure
3. Validar que infrastructure entities existen y son compatibles

---

### 2. **Imposibilidad de Compilar/Validar**

**Problema:**
- No se puede ejecutar `go build ./...` (requiere Go 1.25.0)
- No se pueden ejecutar tests
- No se puede verificar que el código funciona

**Solución Aplicada en Fase 1:**
✅ Crear código basado en análisis estático y documentación del sprint
✅ Seguir patrones establecidos en entities actuales
✅ Documentar exhaustivamente para validación en Fase 2

**Acción para Fase 2:**
1. Compilar todo el proyecto: `go build ./...`
2. Ejecutar tests: `go test ./...`
3. Corregir errores encontrados
4. Validar coverage ≥ 80%

---

### 3. **No se Actualizaron Imports en Archivos Existentes**

**Problema:**
- 31 archivos aún importan `internal/domain/entity` e `internal/domain/entities`
- Estos imports fallarán en Fase 2 cuando se eliminen las entities locales

**Archivos Afectados (según documento sprint):**

**Importan `domain/entity` (15 archivos):**
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

**Importan `domain/entities` (16 archivos):**
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

**Solución Aplicada en Fase 1:**
⏸️ NO SE HIZO - Requiere compilación para validar
⏸️ Documentado para Fase 2

**Acción para Fase 2:**
Ejecutar script de actualización masiva de imports:
```bash
# Reemplazar imports de entity/
find internal/ -name "*.go" -type f -exec sed -i \
  's|"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"|pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"|g' {} \;

# Reemplazar imports de entities/
find internal/ -name "*.go" -type f -exec sed -i \
  's|"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entities"|pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"|g' {} \;
```

---

## 📁 ARCHIVOS CREADOS EN FASE 1

### Stubs de Infrastructure Entities

```
internal/infrastructure_stubs/
├── README.md  ← Documentación de stubs (TEMPORAL)
└── postgres/
    └── entities/
        ├── material.go
        ├── user.go
        ├── material_version.go
        ├── progress.go
        ├── assessment.go
        ├── assessment_answer.go
        └── assessment_attempt.go
```

**Características de los Stubs:**
- Fields públicos con tags GORM
- Tipos nativos de Go (uuid.UUID, string, int, time.Time)
- Método `TableName()` para especificar tabla PostgreSQL
- Sin lógica de negocio (solo estructura de datos)
- Comentarios `TODO FASE 2` para recordar reemplazo

### Domain Services

```
internal/domain/services/
├── material_domain_service.go
├── progress_domain_service.go
├── assessment_domain_service.go
└── attempt_domain_service.go
```

**Características:**
- Extraen lógica de negocio de entities
- Usan imports de stubs (temporal)
- Listos para cambiar a imports reales en Fase 2
- Métodos bien definidos y documentados

---

## 🔄 PLAN COMPLETO PARA FASE 2

### Pre-requisitos de Fase 2

Antes de comenzar Fase 2, asegurar:

- ✅ Conexión a internet estable
- ✅ Go 1.25.0 instalado localmente
- ✅ Acceso a `github.com/EduGoGroup/edugo-infrastructure`
- ✅ Docker corriendo (para tests de integración)
- ✅ golangci-lint instalado

### Paso 1: Verificar Infrastructure Entities (15 min)

```bash
# Navegar a infrastructure (si está clonado localmente)
cd /path/to/edugo-infrastructure

# Verificar que entities existen
ls -la postgres/entities/

# Debe mostrar:
# - material.go
# - user.go
# - material_version.go
# - progress.go
# - assessment.go
# - assessment_answer.go
# - assessment_attempt.go

# Verificar tags/releases
git tag | grep entities

# Debe mostrar algo como:
# postgres/entities/v0.1.0
```

**Si NO existen:**
1. Verificar que se completó el Sprint ENTITIES de infrastructure
2. Si no, completar ese sprint primero antes de continuar aquí

**Si SÍ existen:**
✅ Continuar con Paso 2

---

### Paso 2: Eliminar Stubs y Actualizar Imports (30 min)

```bash
cd /path/to/edugo-api-mobile

# 1. Eliminar stubs temporales
rm -rf internal/infrastructure_stubs/

# 2. Reemplazar imports en Domain Services
find internal/domain/services/ -name "*.go" -type f -exec sed -i \
  's|github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure_stubs/postgres/entities|github.com/EduGoGroup/edugo-infrastructure/postgres/entities|g' {} \;

# 3. Reemplazar imports en archivos existentes (entity/)
find internal/ -name "*.go" -type f -exec sed -i \
  's|"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"|pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"|g' {} \;

# 4. Reemplazar imports en archivos existentes (entities/)
find internal/ -name "*.go" -type f -exec sed -i \
  's|"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entities"|pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"|g' {} \;

# 5. Agregar/actualizar dependencia en go.mod
go get github.com/EduGoGroup/edugo-infrastructure/postgres/entities@latest
go mod tidy

# 6. Verificar que no hay referencias a stubs
grep -r "infrastructure_stubs" internal/
# Debe estar vacío
```

---

### Paso 3: Actualizar Referencias de Tipos (2-3 horas)

Este paso requiere **revisión manual** de cada archivo afectado.

**Cambios típicos necesarios:**

#### Repositories

**Antes:**
```go
import "github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"

func (r *MaterialRepository) Create(material *entity.Material) error {
    // Usar getters
    id := material.ID().String()
    title := material.Title()
}
```

**Después:**
```go
import pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"

func (r *MaterialRepository) Create(material *pgentities.Material) error {
    // Fields públicos directos
    id := material.ID.String()
    title := material.Title
}
```

#### Application Services

**Antes:**
```go
material, err := entity.NewMaterial(title, desc, authorID, subjectID)
if err != nil {
    return nil, err
}

// Método de entity
err = material.SetS3Info(s3Key, s3URL)
```

**Después:**
```go
// Inyectar MaterialDomainService
type MaterialService struct {
    repo        repository.MaterialRepository
    domainSvc   *services.MaterialDomainService  // NUEVO
}

// Crear material manualmente
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

// Usar domain service
err = s.domainSvc.SetS3Info(material, s3Key, s3URL)
```

#### DTOs

**Antes:**
```go
func ToMaterialResponse(material *entity.Material) *MaterialResponse {
    return &MaterialResponse{
        ID:    material.ID().String(),
        Title: material.Title(),
    }
}
```

**Después:**
```go
func ToMaterialResponse(material *pgentities.Material) *MaterialResponse {
    return &MaterialResponse{
        ID:    material.ID.String(),
        Title: material.Title,
    }
}
```

**Archivos a Revisar Manualmente (31 archivos):**
- Ver lista completa en sección "INCONVENIENTE #3" arriba

---

### Paso 4: Eliminar Entities Locales (5 min)

```bash
cd /path/to/edugo-api-mobile

# Eliminar entities de entity/
rm -rf internal/domain/entity/

# Eliminar entities de entities/
rm internal/domain/entities/assessment.go
rm internal/domain/entities/answer.go
rm internal/domain/entities/attempt.go

# Eliminar tests de entities (lógica ya no está ahí)
rm internal/domain/entities/assessment_test.go
rm internal/domain/entities/answer_test.go
rm internal/domain/entities/attempt_test.go

# Si directorio entities/ quedó vacío, eliminarlo
rmdir internal/domain/entities/ 2>/dev/null || true
```

---

### Paso 5: Crear Tests de Domain Services (2-3 horas)

Crear 4 test suites:

1. `internal/domain/services/material_domain_service_test.go`
2. `internal/domain/services/progress_domain_service_test.go`
3. `internal/domain/services/assessment_domain_service_test.go`
4. `internal/domain/services/attempt_domain_service_test.go`

**Template para tests:**

```go
package services_test

import (
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"

    pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
    "github.com/EduGoGroup/edugo-api-mobile/internal/domain/services"
    "github.com/EduGoGroup/edugo-shared/common/types/enum"
)

func TestMaterialDomainService_SetS3Info(t *testing.T) {
    svc := services.NewMaterialDomainService()

    material := &pgentities.Material{
        ID:               uuid.New(),
        Title:            "Test Material",
        ProcessingStatus: enum.ProcessingStatusPending,
    }

    // Test válido
    err := svc.SetS3Info(material, "test-key", "https://test-url.com")
    assert.NoError(t, err)
    assert.Equal(t, "test-key", material.S3Key)
    assert.Equal(t, "https://test-url.com", material.S3URL)
    assert.Equal(t, enum.ProcessingStatusProcessing, material.ProcessingStatus)

    // Test inválido (parámetros vacíos)
    err = svc.SetS3Info(material, "", "")
    assert.Error(t, err)
}

// ... más tests ...
```

**Migrar tests de entities:**
- Lógica de `internal/domain/entities/assessment_test.go` → `assessment_domain_service_test.go`
- Lógica de `internal/domain/entities/answer_test.go` → `attempt_domain_service_test.go`
- Lógica de `internal/domain/entities/attempt_test.go` → `attempt_domain_service_test.go`

---

### Paso 6: Actualizar Tests de Repositories (2 horas)

Archivos a actualizar (9):
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

1. Actualizar imports
2. Cambiar de getters a acceso directo de fields
3. Ajustar constructores (entities de infrastructure no tienen `NewMaterial`, etc.)

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

---

### Paso 7: Actualizar Tests de Application Services (1 hora)

Archivos a actualizar (4):
```
internal/application/service/auth_service_test.go
internal/application/service/material_service_test.go
internal/application/service/progress_service_test.go
internal/application/service/assessment_attempt_service.go (si tiene tests)
```

**Cambios:**
1. Actualizar imports
2. Actualizar mocks para trabajar con nuevos tipos
3. Ajustar assertions (de `material.Title()` a `material.Title`)

---

### Paso 8: Validación Final (60 min)

#### 8.1 Compilación

```bash
cd /path/to/edugo-api-mobile

# Compilar todo
go build ./...

# Debe compilar sin errores
```

**Si hay errores:**
- Revisar imports
- Verificar tipos
- Corregir referencias a métodos que ya no existen

#### 8.2 Tests Unitarios

```bash
# Ejecutar tests unitarios
go test -short ./...

# Debe pasar sin errores
```

#### 8.3 Tests de Integración

```bash
# Asegurar Docker corriendo
docker ps

# Ejecutar todos los tests
go test ./...

# Debe pasar sin errores
```

#### 8.4 Coverage

```bash
# Generar coverage
go test -coverprofile=coverage.out ./...

# Ver coverage total
go tool cover -func=coverage.out | grep total

# Debe ser ≥ 80% (igual o mejor que antes)
```

#### 8.5 Linter

```bash
# Ejecutar linter
golangci-lint run

# No debe haber nuevos errores críticos
```

---

### Paso 9: Documentación (60 min)

Crear/actualizar documentación:

1. **README.md:**
   - Mencionar uso de entities centralizadas
   - Documentar nuevos domain services
   - Actualizar arquitectura si es necesario

2. **docs/MIGRATION_ENTITIES_TO_INFRASTRUCTURE.md** (NUEVO):
   ```markdown
   # Migración a Entities Centralizadas

   ## Qué Cambió

   - Eliminadas entities locales de `internal/domain/entity/` y `internal/domain/entities/`
   - Importadas entities centralizadas de `github.com/EduGoGroup/edugo-infrastructure/postgres/entities`
   - Creados Domain Services para lógica de negocio

   ## Cómo Trabajar con las Nuevas Entities

   ### Crear Material

   **Antes:**
   ```go
   material, err := entity.NewMaterial(...)
   ```

   **Ahora:**
   ```go
   material := &pgentities.Material{
       ID:    uuid.New(),
       Title: title,
       // ...
   }

   // Validar con domain service
   materialSvc := services.NewMaterialDomainService()
   err := materialSvc.Validate(material)
   ```

   ## Agregar Lógica de Negocio

   ❌ **NO agregar lógica en entities**
   ✅ **Agregar lógica en Domain Services**

   Las entities en infrastructure son solo estructura de datos.
   ```

3. **CHANGELOG.md:**
   - Agregar entrada de migración a entities centralizadas

---

## 🎯 CRITERIOS DE ÉXITO PARA FASE 2

Sprint completado cuando:

- [ ] Stubs eliminados (`internal/infrastructure_stubs/` no existe)
- [ ] Imports actualizados (todos usan infrastructure entities)
- [ ] 31 archivos actualizados sin errores
- [ ] 7 entities locales eliminados
- [ ] 4 domain services creados y testeados
- [ ] 4 test suites de domain services pasando
- [ ] 9 test suites de repositories pasando (actualizados)
- [ ] 4 test suites de services pasando (actualizados)
- [ ] Compilación exitosa: `go build ./...` ✅
- [ ] Tests pasando: `go test ./...` ✅
- [ ] Coverage ≥ 80%
- [ ] Linter sin nuevos errores críticos
- [ ] Documentación actualizada
- [ ] PR creado y mergeado

---

## 📝 NOTAS IMPORTANTES

### Value Objects vs UUIDs

Infrastructure entities usan `uuid.UUID` directamente, no value objects.

**Conversión necesaria en repositories:**

```go
// De BD → Value Object
materialID := valueobject.MaterialID{UUID: material.ID}

// De Value Object → BD
material.ID = materialID.UUID()
```

### Encapsulación

Infrastructure entities tienen **fields públicos** (para GORM/JSON).

Si necesitas encapsulación en domain layer:
- Opción A: Wrappear entity en domain object con métodos
- Opción B: Usar domain services (recomendado) ✅

### Lógica de Negocio

**REGLA ESTRICTA:** Entities en infrastructure NO tienen lógica de negocio.

**Dónde poner lógica:**
- ✅ Domain Services (validaciones complejas, reglas de negocio)
- ✅ Application Services (orquestación, casos de uso)
- ✅ Value Objects (validaciones de formato)
- ❌ Entities (solo estructura de datos)

### Tests

**Estrategia cambió:**
- ❌ Tests de entities (antes validaban lógica embebida)
- ✅ Tests de domain services (validan lógica de negocio)
- ✅ Tests de repositories (mapeo entity ↔ BD)
- ✅ Tests de application services (orquestación)

---

## 🚀 SIGUIENTES PASOS (DESPUÉS DE FASE 2)

Una vez completado Sprint ENTITIES en api-mobile:

1. **api-administracion**: Ejecutar su propio sprint de adaptación (similar a este)
2. **worker**: Ejecutar su sprint (más simple, solo MongoDB entities)
3. **Validar consistencia**: Verificar que todos los proyectos usan mismo patrón
4. **Documentar aprendizajes**: Crear guía maestra de uso de entities centralizadas

---

## 📚 REFERENCIAS

- Sprint original: `docs/cicd/sprints/SPRINT-ENTITIES-ADAPTATION.md`
- Documentación DDD: https://martinfowler.com/bliki/EvansClassification.html
- Clean Architecture: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

---

## 👥 CONTACTO Y SOPORTE

Si tienes dudas durante Fase 2, revisar:

1. Este documento (`SPRINT-ENTITIES-PHASE1-REPORT.md`)
2. Sprint original (`SPRINT-ENTITIES-ADAPTATION.md`)
3. Stubs de referencia (antes de eliminarlos)
4. Domain Services creados (para ver patrones)

---

**Generado por:** Claude Code (Sprint Entities Adaptation - Fase 1)
**Fecha:** 22 de Noviembre, 2025
**Versión:** 1.0
**Siguiente Paso:** Ejecutar Fase 2 con ambiente completo
