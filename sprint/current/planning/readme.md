# 📊 Estado Completo del Sprint - EduGo API Mobile

**Fecha de Actualización**: 2025-11-05  
**Sprint Actual**: Fase 2 - Completar TODOs de Servicios  
**Branch**: `fix/debug-sprint-commands`

---

## 🎯 Resumen Ejecutivo

Este documento consolida **todo el estado del sprint actual** en un solo lugar, incluyendo:
- ✅ Trabajo completado (tests y benchmarks)
- ⏳ Trabajo pendiente (optimización de PostgreSQL)
- 📋 Planificación futura (Fase 2 de testing para próximo sprint)

---

# PARTE 1: TRABAJO COMPLETADO ✅

## 📦 Adaptaciones de Corto Plazo - COMPLETADAS

**Fecha de Completitud**: 2025-11-05  
**Estado**: ✅ **100% COMPLETADO**

### 🎯 Objetivos Completados

1. ✅ **Refactorización de MaterialHandler** para mejor inyección de dependencias
2. ✅ **Habilitación de tests S3** previamente skipped
3. ✅ **Implementación de benchmarks** de performance
4. ✅ **Documentación de Fase 2** para siguiente sprint

---

### 📁 Archivos Creados

```
✨ internal/infrastructure/storage/s3/interface.go
✨ internal/infrastructure/http/handler/benchmarks_test.go
✨ sprint/current/planning/fase-2-tests-siguiente-sprint.md
✨ sprint/current/planning/adaptaciones-corto-plazo-completadas.md
```

### 📝 Archivos Modificados

```
📝 internal/infrastructure/http/handler/material_handler.go
📝 internal/infrastructure/http/handler/material_handler_test.go
📝 internal/infrastructure/http/handler/mocks_test.go
```

---

### 1. Refactorización de S3 Client → S3 Storage Interface

**Problema**: MaterialHandler tenía acoplamiento fuerte con implementación concreta de S3Client

**Solución**: Introducir interface S3Storage para mejorar testabilidad

#### Código Implementado:

**`internal/infrastructure/storage/s3/interface.go`** (NUEVO)
```go
package s3

import (
	"context"
	"time"
)

type S3Storage interface {
	GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expires time.Duration) (string, error)
	GeneratePresignedDownloadURL(ctx context.Context, key string, expires time.Duration) (string, error)
}
```

**Cambios en `material_handler.go`**:
- Cambio de `s3Client *s3.S3Client` → `s3Storage s3.S3Storage`
- Actualización del constructor `NewMaterialHandler`
- Todas las llamadas cambiadas de `h.s3Client` → `h.s3Storage`

**Beneficios**:
- ✅ Mejor testabilidad (mock injection)
- ✅ Cumplimiento de SOLID (Dependency Inversion)
- ✅ Preparación para implementaciones alternativas de storage

---

### 2. Tests S3 Habilitados con Mock Completo

**Antes**: Test `TestMaterialHandler_GenerateUploadURL_ValidFileNames` estaba skipped

**Ahora**: ✅ 5 casos de test implementados y pasando

**Casos Testeados**:
1. ✅ Nombre simple válido (`document.pdf`)
2. ✅ Nombre con guiones (`my-document-2024.pdf`)
3. ✅ Nombre con guiones bajos (`my_document_final.pdf`)
4. ✅ Nombre con espacios (`my document.pdf`)
5. ✅ Imagen PNG (`diagram.png`)

**Validaciones del Test**:
- Correcta generación de S3 key (`materials/{id}/{filename}`)
- Propagación de content-type
- Estructura de respuesta (`upload_url`, `s3_key`, `expires_in`)

---

### 3. Suite de Benchmarks de Performance

**Archivo**: `internal/infrastructure/http/handler/benchmarks_test.go`

#### Benchmarks Implementados (11 total):

| # | Benchmark | ns/op | B/op | allocs/op | Descripción |
|---|-----------|-------|------|-----------|-------------|
| 1 | `BenchmarkAuthHandler_Login` | 11,306,302 | 4,208 | 36 | Login secuencial |
| 2 | `BenchmarkAuthHandler_Login_Parallel` | 1,426,971 | 3,792 | 36 | Login paralelo (7.9x más rápido) ⚡ |
| 3 | `BenchmarkAuthHandler_Refresh` | 13,081,361 | 3,144 | 27 | Token refresh |
| 4 | `BenchmarkMaterialHandler_CreateMaterial` | 16,479,838 | 4,103 | 34 | Crear material |
| 5 | `BenchmarkMaterialHandler_GenerateUploadURL` | 15,133,235 | 3,730 | 34 | Generar URL upload |
| 6 | `BenchmarkMaterialHandler_GenerateUploadURL_Parallel` | 1,920,740 | 3,694 | 34 | URL upload paralelo (7.8x más rápido) ⚡ |
| 7 | `BenchmarkMaterialHandler_ListMaterials` | 21,587,397 | 27,055 | 114 | Listar 50 materiales |
| 8 | `BenchmarkMaterialHandler_GetMaterial` | 5,974,420 | 2,154 | 16 | Obtener material |
| 9 | `BenchmarkJSONSerialization` | 472 | 352 | 1 | Serialización JSON |
| 10 | `BenchmarkPathTraversalValidation` | **12** | **0** | **0** | Validación seguridad ✅ |
| 11 | `BenchmarkErrorHandling` | 254 | 480 | 6 | Manejo de errores |

#### Análisis de Performance:

**Excelente** ✅:
- PathTraversalValidation: **12ns** sin allocaciones (óptimo)
- JSONSerialization: **472ns** excelente
- Paralelización: **7-8x speedup** en operaciones I/O

**Áreas de Mejora** ⚠️:
- ErrorHandling: 480 bytes/op (considerar object pooling)
- ListMaterials: 27KB/op con 50 items (optimizar serialización)

**Comando de ejecución**:
```bash
go test -bench=. -benchmem -benchtime=1s ./internal/infrastructure/http/handler/...
```

---

### 4. Estado de Tests Actual

| Handler | Tests Pasando | Tests Skipped | Cobertura Estimada |
|---------|---------------|---------------|-------------------|
| **AuthHandler** | 19 ✅ | 0 | ~85% |
| **MaterialHandler** | 10 ✅ | 0 | ~80% |
| **HealthHandler** | 4 ✅ | 7 ⏭️ | ~30% |
| AssessmentHandler | 0 | - | 0% |
| ProgressHandler | 0 | - | 0% |
| StatsHandler | 0 | - | 0% |
| SummaryHandler | 0 | - | 0% |

**Total**: 
- ✅ **33 tests pasando**
- ⏭️ **7 tests skipped** (requieren testcontainers)
- ❌ **0 tests fallando**
- 🎯 **11 benchmarks** funcionando

---

### 🚀 Comandos Útiles

```bash
# Ejecutar todos los tests
go test ./internal/infrastructure/http/handler/...

# Ejecutar tests con verbose
go test -v ./internal/infrastructure/http/handler/...

# Ejecutar solo tests de material
go test ./internal/infrastructure/http/handler/... -run TestMaterialHandler

# Ejecutar benchmarks
go test -bench=. -benchmem ./internal/infrastructure/http/handler/...

# Ejecutar benchmarks de auth
go test -bench=BenchmarkAuth.* -benchmem ./internal/infrastructure/http/handler/...

# Ver cobertura
go test -coverprofile=coverage.out ./internal/infrastructure/http/handler/...
go tool cover -html=coverage.out
```

---

---

# PARTE 2: TRABAJO PENDIENTE ⏳

## 🔧 Optimización de PostgreSQL - Índice en Materials

**Estado**: ⏳ **0% COMPLETADO** (0/24 tareas)  
**Estimación**: 10-15 minutos  
**Objetivo**: Crear índice descendente en `materials.updated_at` para optimizar queries con `ORDER BY updated_at DESC`

**Mejora Esperada**: 5-10x más rápido (de 50-200ms a 5-20ms)

---

## 📋 Plan de Ejecución Detallado

### Fase 1: Preparación y Validación ⏳ (0/4)

- [ ] **1.1** - Verificar conexión a PostgreSQL local
  ```bash
  psql -d edugo_db_local -c "SELECT current_database(), version();"
  ```

- [ ] **1.2** - Verificar existencia de tabla materials
  ```bash
  psql -d edugo_db_local -c "SELECT COUNT(*) FROM materials;"
  ```

- [ ] **1.3** - Verificar índices existentes
  ```bash
  psql -d edugo_db_local -c "SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'materials';"
  ```

- [ ] **1.4** - Medir performance baseline (ANTES del índice)
  ```bash
  psql -d edugo_db_local -c "EXPLAIN ANALYZE SELECT * FROM materials ORDER BY updated_at DESC LIMIT 20;"
  ```

**Resultado esperado**: Baseline documentado, `idx_materials_updated_at` NO existe aún

---

### Fase 2: Creación del Script ⏳ (0/4)

- [ ] **2.1** - Verificar carpeta de scripts SQL
  ```bash
  ls -la scripts/postgresql/
  ```

- [ ] **2.2** - Identificar número secuencial para el nuevo script
  ```bash
  ls scripts/postgresql/ | grep -E '^[0-9]+_' | sort -V | tail -1
  ```

- [ ] **2.3** - Crear archivo `scripts/postgresql/06_indexes_materials.sql`
  
  **Contenido del archivo**:
  ```sql
  -- ============================================================
  -- Migration: 06_indexes_materials.sql
  -- Description: Agregar índice descendente en materials.updated_at
  --              para optimizar queries de listado cronológico
  -- Author: Claude Code / EduGo Team
  -- Date: 2025-11-05
  -- ============================================================

  -- Objetivo:
  -- Mejorar performance de queries que ordenan materiales por fecha
  -- de actualización más reciente (patrón común en la aplicación).
  --
  -- Queries beneficiadas:
  -- 1. SELECT * FROM materials ORDER BY updated_at DESC LIMIT N;
  -- 2. SELECT * FROM materials WHERE course_id = X ORDER BY updated_at DESC;
  -- 3. SELECT * FROM materials WHERE type = 'Y' ORDER BY updated_at DESC;
  --
  -- Mejora esperada: 5-10x más rápido (de 50-200ms a 5-20ms)

  -- Crear índice descendente de forma idempotente
  CREATE INDEX IF NOT EXISTS idx_materials_updated_at
  ON materials(updated_at DESC);

  -- Verificación:
  -- Después de ejecutar este script, verificar con:
  -- SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'materials';
  --
  -- Validar uso del índice con:
  -- EXPLAIN ANALYZE SELECT * FROM materials ORDER BY updated_at DESC LIMIT 20;
  -- Debe mostrar: "Index Scan using idx_materials_updated_at"

  -- Rollback (si es necesario):
  -- DROP INDEX IF EXISTS idx_materials_updated_at;
  ```

- [ ] **2.4** - Validar sintaxis SQL
  ```bash
  psql -d edugo_db_local -c "BEGIN; \i scripts/postgresql/06_indexes_materials.sql; ROLLBACK;"
  ```

**Resultado esperado**: Script SQL creado y validado sin errores de sintaxis

---

### Fase 3: Ejecución Local ⏳ (0/4)

- [ ] **3.1** - Ejecutar script de migración
  ```bash
  psql -d edugo_db_local -f scripts/postgresql/06_indexes_materials.sql
  ```

- [ ] **3.2** - Verificar creación del índice
  ```bash
  psql -d edugo_db_local -c "SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'materials' AND indexname = 'idx_materials_updated_at';"
  ```

- [ ] **3.3** - Validar que el índice es utilizado
  ```bash
  psql -d edugo_db_local -c "EXPLAIN ANALYZE SELECT * FROM materials ORDER BY updated_at DESC LIMIT 20;"
  ```
  - Debe mostrar: `Index Scan using idx_materials_updated_at`

- [ ] **3.4** - Probar idempotencia del script
  ```bash
  psql -d edugo_db_local -f scripts/postgresql/06_indexes_materials.sql
  ```
  - Debe mostrar: `NOTICE: relation "idx_materials_updated_at" already exists, skipping`

**Resultado esperado**: Índice creado, verificado y funcionando

---

### Fase 4: Validación de Aplicación ⏳ (0/4)

- [ ] **4.1** - Verificar que la aplicación compila
  ```bash
  go build ./...
  ```

- [ ] **4.2** - Ejecutar suite de tests unitarios
  ```bash
  go test ./... -v
  ```

- [ ] **4.3** - Ejecutar tests de integración (si existen)
  ```bash
  go test ./... -tags=integration -v
  ```

- [ ] **4.4** - Probar manualmente endpoint (opcional)
  ```bash
  # Levantar servidor
  go run cmd/main.go
  
  # En otra terminal
  curl -X GET "http://localhost:8080/api/materials?sort=updated_at&order=desc&limit=20"
  ```

**Resultado esperado**: Aplicación funciona correctamente, tests pasan

---

### Fase 5: Control de Versiones ⏳ (0/5)

- [ ] **5.1** - Verificar estado de Git
  ```bash
  git status
  ```

- [ ] **5.2** - Agregar script al staging
  ```bash
  git add scripts/postgresql/06_indexes_materials.sql
  ```

- [ ] **5.3** - Crear commit con mensaje descriptivo
  ```bash
  git commit -m "perf(db): agregar índice en materials.updated_at para optimizar ordenamiento

  - Crear script 06_indexes_materials.sql
  - Índice descendente (DESC) para queries con ORDER BY updated_at DESC
  - Script idempotente con IF NOT EXISTS
  - Mejora esperada: 5-10x más rápido (50-200ms → 5-20ms)
  - Sin cambios en código Go (optimización transparente)

  Queries beneficiadas:
  - Listado de materiales recientes
  - Filtros por curso/tipo + ordenamiento cronológico

  Validado con EXPLAIN ANALYZE en ambiente local.

  🤖 Generated with [Claude Code](https://claude.com/claude-code)

  Co-Authored-By: Claude <noreply@anthropic.com>"
  ```

- [ ] **5.4** - Actualizar plan de sprint con checkboxes completados
  - Marcar todas las casillas en este documento
  - Actualizar `sprint/current/readme.md`

- [ ] **5.5** - Crear commit de documentación
  ```bash
  git add sprint/current/planning/ESTADO-COMPLETO-SPRINT.md sprint/current/readme.md
  git commit -m "docs(sprint): marcar optimización de índice como completada

  - Actualizar ESTADO-COMPLETO-SPRINT.md
  - Documentar resultado de validación
  - Sprint completado exitosamente

  🤖 Generated with [Claude Code](https://claude.com/claude-code)

  Co-Authored-By: Claude <noreply@anthropic.com>"
  ```

**Resultado esperado**: Cambios committeados con mensajes apropiados

---

### Fase 6: Preparación para Deployment [OPCIONAL] ⏳ (0/3)

- [ ] **6.1** - Documentar instrucciones para QA
- [ ] **6.2** - Documentar consideraciones para producción
- [ ] **6.3** - Notificar al equipo sobre cambio pendiente

**Resultado esperado**: Documentación lista para DevOps

---

## 🎯 Resumen de Pendientes

| Fase | Tareas | Completadas | Estado |
|------|--------|-------------|--------|
| Fase 1: Preparación | 4 | 0/4 | ⏳ Pendiente |
| Fase 2: Script | 4 | 0/4 | ⏳ Pendiente |
| Fase 3: Ejecución | 4 | 0/4 | ⏳ Pendiente |
| Fase 4: Validación | 4 | 0/4 | ⏳ Pendiente |
| Fase 5: Git | 5 | 0/5 | ⏳ Pendiente |
| Fase 6: Deployment (opcional) | 3 | 0/3 | ⏳ Pendiente |
| **TOTAL** | **24** | **0/24** | **0% Completado** |

**Estimación Total**: 10-15 minutos

**Próximo Comando**: 
```bash
/03-execution
```

---

## 🚨 Manejo de Errores

### Error: PostgreSQL no conecta
**Solución**:
```bash
# macOS
brew services start postgresql

# Linux
sudo systemctl start postgresql
```

### Error: Tabla materials no existe
**Solución**:
```bash
# Ejecutar migraciones previas
ls scripts/postgresql/*.sql | sort -V | xargs -I {} psql -d edugo_db_local -f {}
```

### Error: Índice no se usa en EXPLAIN ANALYZE
**Causa**: Tabla muy pequeña (<100 registros), optimizador elige Seq Scan
**Solución**: Es comportamiento esperado, documentar y validar en QA con datos reales

### Rollback en caso de problemas
```sql
DROP INDEX IF EXISTS idx_materials_updated_at;
```

---

---

# PARTE 3: PLANIFICACIÓN FUTURA 📋

## 🔮 Fase 2 de Testing - Próximo Sprint

**Estado**: 📋 **DOCUMENTACIÓN** (no es trabajo inmediato)  
**Objetivo**: Alcanzar **80%+ de cobertura global** en todos los handlers  
**Estimación Total**: 21-28 horas de desarrollo

---

## 📊 Estado Actual vs Objetivo

### Cobertura Actual de Handlers

| Handler | Tests | Cobertura | Estado |
|---------|-------|-----------|--------|
| AuthHandler | 19 ✅ | ~85% | ✅ Completo |
| MaterialHandler | 10 ✅ | ~80% | ✅ Completo |
| HealthHandler | 4 ✅ / 7 ⏭️ | ~30% | ⚠️ Parcial |
| AssessmentHandler | 0 | 0% | ❌ Pendiente |
| ProgressHandler | 0 | 0% | ❌ Pendiente |
| StatsHandler | 0 | 0% | ❌ Pendiente |
| SummaryHandler | 0 | 0% | ❌ Pendiente |

**Total Actual**: 33 tests, ~50% cobertura global  
**Objetivo Fase 2**: 70-80 tests, 80%+ cobertura global

---

## 📋 Plan de Implementación por Sprints

### Sprint 1: HealthHandler + Testcontainers (6 horas)

**Objetivo**: Completar HealthHandler con tests de integración reales

**Tareas**:
- [ ] Setup de testcontainers para PostgreSQL
- [ ] Setup de testcontainers para MongoDB
- [ ] Implementar test: `Check_AllHealthy` (DB real)
- [ ] Implementar test: `Check_PostgreSQL_Degraded`
- [ ] Implementar test: `Check_MongoDB_Degraded`
- [ ] Implementar test: `Check_BothDatabases_Down`
- [ ] Implementar test: `Check_ResponseTime_Acceptable`
- [ ] Benchmark: `HealthCheck` con DBs reales

**Entregables**:
- ✅ Testcontainers configurados y reutilizables
- ✅ HealthHandler 80%+ coverage
- ✅ Helper functions documentadas

**Ejemplo de código**:
```go
func TestHealthHandler_Check_WithTestContainers(t *testing.T) {
    // Setup PostgreSQL testcontainer
    pgContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15-alpine"),
    )
    require.NoError(t, err)
    defer pgContainer.Terminate(ctx)

    // Setup MongoDB testcontainer
    mongoContainer, err := mongodb.RunContainer(ctx,
        testcontainers.WithImage("mongo:7"),
    )
    require.NoError(t, err)
    defer mongoContainer.Terminate(ctx)

    // Conectar y testear...
}
```

---

### Sprint 2: AssessmentHandler (8 horas)

**Objetivo**: Suite completa de tests para evaluaciones

**Funcionalidades a testear**:
- [ ] CreateAssessment (success, invalid request, unauthorized)
- [ ] GetAssessment (success, not found, unauthorized)
- [ ] ListAssessments (success, pagination, filters)
- [ ] UpdateAssessment (success, not found, unauthorized)
- [ ] DeleteAssessment (success, not found, unauthorized)
- [ ] SubmitAnswer (success, invalid format, time expired)
- [ ] GetResults (success, not completed, unauthorized)

**Tests de seguridad críticos**:
- Prevención de acceso a evaluaciones de otros usuarios
- Validación de tiempos de expiración
- Validación de respuestas (XSS, injection)

**Entregables**:
- ✅ AssessmentHandler 75%+ coverage
- ✅ Tests de seguridad documentados
- ✅ Benchmarks de operaciones críticas

---

### Sprint 3: Progress + Stats + Summary (7 horas)

#### ProgressHandler (4 horas)
- [ ] GetUserProgress (success, unauthorized, different user)
- [ ] UpdateProgress (success, invalid percentage, unauthorized)
- [ ] GetMaterialProgress (success, material not found)
- [ ] ListProgressBySubject (success, pagination)
- [ ] GetCompletionStats (success, empty data)

#### StatsHandler (4 horas)
- [ ] GetGlobalStats (success, admin only)
- [ ] GetUserStats (success, own user, unauthorized)
- [ ] GetMaterialStats (success, material not found)
- [ ] GetSubjectStats (success, date filters)
- [ ] ExportStats (success, format validation)

#### SummaryHandler (3 horas)
- [ ] GenerateSummary (success, material not found)
- [ ] GetSummary (success, not generated, unauthorized)
- [ ] ListSummaries (success, pagination)
- [ ] RegenerateSummary (success, already processing)

**Entregables**:
- ✅ ProgressHandler 75%+ coverage
- ✅ StatsHandler 75%+ coverage
- ✅ SummaryHandler 75%+ coverage
- ✅ Cobertura global 80%+

---

## 🔧 Infraestructura Necesaria

### Testcontainers Setup

**Dependencias a agregar**:
```go
import (
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/modules/mongodb"
)
```

**Helper Functions** (crear archivo `testcontainers_helpers.go`):
```go
func SetupPostgreSQLTestContainer(ctx context.Context) (*sql.DB, func(), error) {
    // Implementar setup de PostgreSQL
}

func SetupMongoDBTestContainer(ctx context.Context) (*mongo.Database, func(), error) {
    // Implementar setup de MongoDB
}
```

---

## 📈 Métricas de Calidad Esperadas

### Cobertura de Código
- **Actual**: ~50% en handlers implementados
- **Objetivo Fase 2**: 80%+ global
- **Crítico**: 100% en validaciones de seguridad

### Tiempo de Ejecución
- **Actual**: ~15s para suite completa
- **Objetivo Fase 2**: <30s con testcontainers
- **CI/CD**: <2min en pipeline completo

### Benchmarks Adicionales
```
- [ ] BenchmarkAssessmentHandler_SubmitAnswer
- [ ] BenchmarkProgressHandler_UpdateProgress
- [ ] BenchmarkStatsHandler_GetGlobalStats
- [ ] BenchmarkHealthHandler_Check (con testcontainers)
```

---

## 🎯 Criterios de Aceptación - Fase 2

### Funcionales
- [ ] Todos los handlers tienen tests de CRUD completos
- [ ] HealthHandler con testcontainers funcionando
- [ ] Cobertura global ≥80%

### Seguridad
- [ ] Tests de autorización en todos los endpoints protegidos
- [ ] Tests de validación de input en todos los endpoints
- [ ] Tests de prevención de ataques comunes (XSS, injection)

### Performance
- [ ] Benchmarks para todos los endpoints críticos
- [ ] Suite completa ejecuta en <30s
- [ ] Documentación de métricas de performance

### Calidad
- [ ] Tests siguen patrón AAA consistente
- [ ] Mocks reutilizables y bien documentados
- [ ] README con instrucciones de ejecución

---

## 📚 Referencias y Recursos

### Archivos de Referencia (Completados)
- `auth_handler_test.go` - Patrón de tests de autenticación
- `material_handler_test.go` - Tests de seguridad (path traversal)
- `mocks_test.go` - Mocks reutilizables
- `testing_helpers.go` - Helpers comunes
- `benchmarks_test.go` - Suite de benchmarks

### Documentación Útil
- Testcontainers Go: https://golang.testcontainers.org/
- Go testing best practices: https://go.dev/doc/tutorial/add-a-test
- Gin testing guide: https://gin-gonic.com/docs/testing/

---

---

# RESUMEN GLOBAL 🎯

## Estado Consolidado del Sprint

| Categoría | Estado | Detalle |
|-----------|--------|---------|
| **Tests Completados** | ✅ 100% | 33 tests, 11 benchmarks |
| **Refactorización S3** | ✅ 100% | Interface implementada |
| **Optimización PostgreSQL** | ⏳ 0% | 0/24 tareas pendientes |
| **Planificación Fase 2** | 📋 Documentado | Listo para próximo sprint |

---

## Próximos Pasos Inmediatos

### 1. Completar Optimización de PostgreSQL (10-15 min)
```bash
/03-execution
```
Esto ejecutará las 24 tareas pendientes del plan de optimización

### 2. Push y PR (5 min)
```bash
git push origin fix/debug-sprint-commands
# Crear Pull Request en GitHub
```

### 3. Planificar Próximo Sprint (Fase 2 Testing)
- Revisar documento de Fase 2
- Estimar 21-28 horas de desarrollo
- Priorizar HealthHandler con testcontainers

---

## Comandos Rápidos

```bash
# Ver estado actual
git status
git log -5 --oneline

# Ejecutar tests
go test ./internal/infrastructure/http/handler/...

# Ejecutar benchmarks
go test -bench=. -benchmem ./internal/infrastructure/http/handler/...

# Ejecutar plan de PostgreSQL
/03-execution

# Ver cobertura
go test -coverprofile=coverage.out ./internal/infrastructure/http/handler/...
go tool cover -html=coverage.out
```

---

**Última actualización**: 2025-11-05  
**Autor**: Claude Code + Jhoan Medina  
**Documento Maestro**: Consolida readme.md + adaptaciones-corto-plazo-completadas.md + fase-2-tests-siguiente-sprint.md
