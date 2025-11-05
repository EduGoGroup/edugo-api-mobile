# Fase 2: Plan de Testing para Siguiente Sprint

**Fecha de Creación**: 2025-11-05  
**Estado**: Pendiente para próximo sprint  
**Prioridad**: Media-Alta

---

## 📊 Estado Actual de Cobertura

### ✅ Handlers Completados (100% Crítico)

#### AuthHandler
- **Tests**: 19 tests pasando
- **Cobertura**: ~85% (crítico y normal)
- **Seguridad**: ✅ Validada
- **Performance**: ✅ Benchmarks incluidos

**Tests implementados**:
- ✅ Login (success, invalid credentials, invalid request, service error)
- ✅ Refresh (success, invalid token, expired token, invalid request)
- ✅ Logout (success, unauthenticated, invalid request, service error)
- ✅ RevokeAll (success, unauthenticated, service error)

#### MaterialHandler
- **Tests**: 10 tests pasando
- **Cobertura**: ~80% (crítico y normal)
- **Seguridad**: ✅ Path traversal prevention validado
- **Performance**: ✅ Benchmarks incluidos

**Tests implementados**:
- ✅ Constructor validation
- ✅ Path traversal prevention (6 casos críticos de seguridad)
- ✅ Valid file names (5 casos con diferentes formatos)
- ✅ CreateMaterial (success, invalid request)
- ✅ GetMaterial (success)
- ✅ GenerateUploadURL (material not found, invalid request)
- ✅ GenerateDownloadURL (file not uploaded)

### ⏳ Handlers Pendientes

#### HealthHandler
- **Tests actuales**: 4 tests pasando, 7 skipped
- **Cobertura**: ~30% (solo tests estructurales)
- **Requiere**: Testcontainers para PostgreSQL y MongoDB

#### Handlers Sin Tests
- `AssessmentHandler` - 0% cobertura
- `ProgressHandler` - 0% cobertura
- `StatsHandler` - 0% cobertura
- `SummaryHandler` - 0% cobertura

---

## 🎯 Objetivos de Fase 2

### Objetivo Principal
Alcanzar **80%+ de cobertura global** en todos los handlers HTTP del proyecto.

### Métricas de Éxito
- [ ] HealthHandler: 80%+ cobertura con testcontainers
- [ ] AssessmentHandler: 75%+ cobertura
- [ ] ProgressHandler: 75%+ cobertura
- [ ] StatsHandler: 75%+ cobertura
- [ ] SummaryHandler: 75%+ cobertura
- [ ] Cobertura global de handlers: 80%+

---

## 📋 Plan de Implementación

### 1. HealthHandler con Testcontainers (Prioridad: Alta)

**Objetivo**: Implementar tests de integración reales para health checks

**Tareas**:
```
- [ ] Configurar testcontainers para PostgreSQL
- [ ] Configurar testcontainers para MongoDB
- [ ] Implementar test: Check_AllHealthy (DB real)
- [ ] Implementar test: Check_PostgreSQL_Degraded
- [ ] Implementar test: Check_MongoDB_Degraded
- [ ] Implementar test: Check_BothDatabases_Down
- [ ] Implementar test: Check_ResponseTime_Acceptable
- [ ] Benchmark: HealthCheck con DBs reales
```

**Archivos a modificar**:
- `internal/infrastructure/http/handler/health_handler_test.go`

**Estimación**: 4-6 horas

**Ejemplo de implementación**:
```go
func TestHealthHandler_Check_WithTestContainers(t *testing.T) {
    // Setup PostgreSQL testcontainer
    pgContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15-alpine"),
        // ... configuración
    )
    require.NoError(t, err)
    defer pgContainer.Terminate(ctx)

    // Setup MongoDB testcontainer
    mongoContainer, err := mongodb.RunContainer(ctx,
        testcontainers.WithImage("mongo:7"),
        // ... configuración
    )
    require.NoError(t, err)
    defer mongoContainer.Terminate(ctx)

    // Conectar a las bases de datos
    db := setupPostgreSQLConnection(pgContainer)
    mongoDB := setupMongoDBConnection(mongoContainer)

    // Crear handler con conexiones reales
    handler := NewHealthHandler(db, mongoDB)

    // Ejecutar test...
}
```

---

### 2. AssessmentHandler Tests (Prioridad: Media-Alta)

**Objetivo**: Implementar suite completa de tests para evaluaciones

**Funcionalidades a testear**:
```
- [ ] CreateAssessment (success, invalid request, unauthorized)
- [ ] GetAssessment (success, not found, unauthorized)
- [ ] ListAssessments (success, pagination, filters)
- [ ] UpdateAssessment (success, not found, unauthorized)
- [ ] DeleteAssessment (success, not found, unauthorized)
- [ ] SubmitAnswer (success, invalid format, time expired)
- [ ] GetResults (success, not completed, unauthorized)
```

**Casos de seguridad críticos**:
- Prevención de acceso a evaluaciones de otros usuarios
- Validación de tiempos de expiración
- Validación de respuestas (XSS, injection)

**Estimación**: 6-8 horas

---

### 3. ProgressHandler Tests (Prioridad: Media)

**Objetivo**: Validar tracking de progreso de estudiantes

**Funcionalidades a testear**:
```
- [ ] GetUserProgress (success, unauthorized, different user)
- [ ] UpdateProgress (success, invalid percentage, unauthorized)
- [ ] GetMaterialProgress (success, material not found)
- [ ] ListProgressBySubject (success, pagination)
- [ ] GetCompletionStats (success, empty data)
```

**Estimación**: 4-5 horas

---

### 4. StatsHandler Tests (Prioridad: Media)

**Objetivo**: Validar generación de estadísticas

**Funcionalidades a testear**:
```
- [ ] GetGlobalStats (success, admin only)
- [ ] GetUserStats (success, own user, unauthorized)
- [ ] GetMaterialStats (success, material not found)
- [ ] GetSubjectStats (success, date filters)
- [ ] ExportStats (success, format validation)
```

**Estimación**: 4-5 horas

---

### 5. SummaryHandler Tests (Prioridad: Baja-Media)

**Objetivo**: Validar generación de resúmenes

**Funcionalidades a testear**:
```
- [ ] GenerateSummary (success, material not found, invalid params)
- [ ] GetSummary (success, not generated, unauthorized)
- [ ] ListSummaries (success, pagination)
- [ ] RegenerateSummary (success, already processing)
```

**Estimación**: 3-4 horas

---

## 🔧 Infraestructura de Testing Necesaria

### Testcontainers Setup

**Dependencias a agregar** (si no están):
```go
import (
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/modules/mongodb"
)
```

### Helper Functions Globales

Crear archivo `internal/infrastructure/http/handler/testcontainers_helpers.go`:
```go
package handler

import (
    "context"
    "database/sql"
    
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/modules/mongodb"
    "go.mongodb.org/mongo-driver/mongo"
)

func SetupPostgreSQLTestContainer(ctx context.Context) (*sql.DB, func(), error) {
    // Implementar setup de PostgreSQL
}

func SetupMongoDBTestContainer(ctx context.Context) (*mongo.Database, func(), error) {
    // Implementar setup de MongoDB
}
```

---

## 📊 Benchmarks Adicionales

### Benchmarks Completados ✅
- AuthHandler: Login, Login_Parallel, Refresh
- MaterialHandler: CreateMaterial, GenerateUploadURL, GenerateUploadURL_Parallel, ListMaterials, GetMaterial
- Utils: JSONSerialization, PathTraversalValidation, ErrorHandling

### Benchmarks Pendientes para Fase 2
```
- [ ] BenchmarkAssessmentHandler_SubmitAnswer
- [ ] BenchmarkAssessmentHandler_GetResults
- [ ] BenchmarkProgressHandler_UpdateProgress
- [ ] BenchmarkStatsHandler_GetGlobalStats
- [ ] BenchmarkHealthHandler_Check (con testcontainers)
```

---

## 🎯 Resultados de Performance Actuales

### Benchmarks Base (Apple M1 Pro)

| Benchmark | ns/op | B/op | allocs/op | Notas |
|-----------|-------|------|-----------|-------|
| **Auth** ||||
| Login | 11,306,302 | 4,208 | 36 | Con mock de DB (10ms) |
| Login_Parallel | 1,426,971 | 3,792 | 36 | 7.9x más rápido con concurrencia |
| Refresh | 13,081,361 | 3,144 | 27 | Con mock de DB (12ms) |
| **Material** ||||
| CreateMaterial | 16,479,838 | 4,103 | 34 | Con mock de DB (15ms) |
| GenerateUploadURL | 15,133,235 | 3,730 | 34 | DB (5ms) + S3 SDK (8ms) |
| GenerateUploadURL_Parallel | 1,920,740 | 3,694 | 34 | 7.8x más rápido |
| ListMaterials | 21,587,397 | 27,055 | 114 | 50 items, mock (20ms) |
| GetMaterial | 5,974,420 | 2,154 | 16 | Query simple (5ms) |
| **Utils** ||||
| JSONSerialization | 472 | 352 | 1 | LoginResponse completo |
| PathTraversalValidation | 12 | 0 | 0 | Sin allocaciones ✅ |
| ErrorHandling | 254 | 480 | 6 | AppError wrapping |

### Análisis de Performance

**Fortalezas** ✅:
- PathTraversalValidation: Extremadamente rápido (12ns), 0 allocaciones
- JSONSerialization: Óptimo (472ns) para estructuras complejas
- Paralelización: 7-8x speedup en operaciones I/O bound

**Áreas de Mejora** ⚠️:
- ErrorHandling: 480 bytes/op, considerar pooling de objetos
- ListMaterials: 27KB/op con 50 items, optimizar serialización
- Operaciones con mock I/O: Añadir caching donde aplique

**Recomendaciones**:
1. Implementar object pooling para errores frecuentes
2. Considerar streaming para listas grandes (>100 items)
3. Añadir caching de respuestas para endpoints de lectura
4. Implementar rate limiting basado en estos benchmarks

---

## ⚠️ Consideraciones de Seguridad

### Tests de Seguridad Críticos (Prioritarios)

Para cada handler, implementar validaciones de:

1. **Autorización**:
   - [ ] Usuario autenticado puede acceder solo a sus recursos
   - [ ] Admin puede acceder a recursos globales
   - [ ] Profesor puede acceder a recursos de sus estudiantes

2. **Validación de Input**:
   - [ ] XSS prevention en campos de texto
   - [ ] SQL injection prevention (debe ser manejado por ORM)
   - [ ] Path traversal prevention (ya implementado en MaterialHandler ✅)

3. **Rate Limiting** (considerar para Fase 2):
   - [ ] Limitar intentos de login
   - [ ] Limitar creación de recursos
   - [ ] Limitar consultas de estadísticas

---

## 📈 Métricas de Calidad Esperadas

### Cobertura de Código
- **Actual**: ~60% en handlers implementados
- **Objetivo Fase 2**: 80%+ global
- **Crítico**: 100% en validaciones de seguridad

### Tiempo de Ejecución
- **Actual**: ~15s para suite completa
- **Objetivo Fase 2**: <30s con testcontainers
- **CI/CD**: <2min en pipeline completo

### Mantenibilidad
- Reutilizar mocks existentes en `mocks_test.go`
- Mantener patrón AAA (Arrange-Act-Assert)
- Documentar tests complejos con comentarios claros

---

## 🚀 Orden de Implementación Recomendado

### Sprint 1 (HealthHandler + Base)
1. Setup de testcontainers infrastructure
2. HealthHandler tests completos
3. Helpers reutilizables para otros handlers

**Entregables**:
- ✅ Testcontainers configurados
- ✅ HealthHandler 80%+ coverage
- ✅ Helpers documentados

---

### Sprint 2 (AssessmentHandler)
1. AssessmentHandler CRUD tests
2. Tests de seguridad (autorización)
3. Tests de validación (respuestas, tiempo)

**Entregables**:
- ✅ AssessmentHandler 75%+ coverage
- ✅ Security tests documentados

---

### Sprint 3 (Progress + Stats + Summary)
1. ProgressHandler tests
2. StatsHandler tests
3. SummaryHandler tests
4. Benchmarks adicionales

**Entregables**:
- ✅ Cobertura global 80%+
- ✅ Suite completa de benchmarks
- ✅ Documentación de performance

---

## 📚 Referencias

### Archivos Completados (Ejemplos)
- `internal/infrastructure/http/handler/auth_handler_test.go` - Patrón de tests de auth
- `internal/infrastructure/http/handler/material_handler_test.go` - Tests de seguridad (path traversal)
- `internal/infrastructure/http/handler/mocks_test.go` - Mocks reutilizables
- `internal/infrastructure/http/handler/testing_helpers.go` - Helpers comunes
- `internal/infrastructure/http/handler/benchmarks_test.go` - Suite de benchmarks

### Documentación Útil
- Testcontainers Go: https://golang.testcontainers.org/
- Go testing best practices: https://go.dev/doc/tutorial/add-a-test
- Gin testing guide: https://gin-gonic.com/docs/testing/

---

## ✅ Criterios de Aceptación (Fase 2)

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
- [ ] README con instrucciones de ejecución de tests

---

## 🎯 Resumen Ejecutivo

**Estado Actual**:
- ✅ 29 tests pasando (AuthHandler: 19, MaterialHandler: 10)
- ✅ 11 benchmarks implementados
- ✅ Seguridad crítica validada (path traversal)
- ✅ Infrastructure de testing establecida

**Fase 2 - Objetivos**:
- 🎯 +40-50 tests adicionales
- 🎯 HealthHandler con testcontainers
- 🎯 4 handlers adicionales testeados
- 🎯 80%+ cobertura global

**Estimación Total Fase 2**: 21-28 horas de desarrollo

**Entregables Finales**:
- Suite completa de tests unitarios y de integración
- Benchmarks de performance documentados
- Infrastructure de testcontainers reutilizable
- Documentación de cobertura y métricas

---

**Última actualización**: 2025-11-05  
**Autor**: Claude Code + Jhoan Medina  
**Versión**: 1.0
