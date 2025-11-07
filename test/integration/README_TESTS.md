# 🧪 EduGo API Mobile - Integration Tests Suite

**Cobertura E2E**: 17/17 tests ✅ (100%)  
**Tiempo promedio**: ~18.5s/test  
**Testcontainers**: PostgreSQL 15 + MongoDB 7 + RabbitMQ 3.12

---

## 📊 Test Coverage Overview

| Flujo | Tests | Cobertura | Tiempo Promedio | Estado |
|-------|-------|-----------|-----------------|--------|
| **Auth** | 3/3 | ✅ 100% | ~19s | ✅ Pasando |
| **Material** | 4/4 | ✅ 100% | ~18s | ✅ Pasando |
| **Assessment** | 4/4 | ✅ 100% | ~18s | ✅ Pasando |
| **Progress** | 4/4 | ✅ 100% | ~19s | ✅ Pasando |
| **Stats** | 2/2 | ✅ 100% | ~19s | ✅ Pasando |
| **TOTAL** | **17/17** | **✅ 100%** | **~18.5s** | **✅ Pasando** |

---

## 🚀 Quick Start

### Prerequisitos

```bash
# 1. Verificar Docker está corriendo
make docker-check

# 2. Instalar dependencias
go mod download

# 3. Verificar variables de entorno
export RUN_INTEGRATION_TESTS=true
```

### Ejecutar Todos los Tests

```bash
# Opción 1: Makefile
make test-integration

# Opción 2: go test directo
RUN_INTEGRATION_TESTS=true go test -v -tags=integration ./test/integration/ -timeout 15m
```

### Ejecutar Tests por Flujo

```bash
# Auth Flow (3 tests)
RUN_INTEGRATION_TESTS=true go test -v -tags=integration ./test/integration/ -run TestAuthFlow

# Material Flow (4 tests)
RUN_INTEGRATION_TESTS=true go test -v -tags=integration ./test/integration/ -run TestMaterialFlow

# Assessment Flow (4 tests)
RUN_INTEGRATION_TESTS=true go test -v -tags=integration ./test/integration/ -run TestAssessmentFlow

# Progress Flow (4 tests)
RUN_INTEGRATION_TESTS=true go test -v -tags=integration ./test/integration/ -run TestProgressFlow

# Stats Flow (2 tests)
RUN_INTEGRATION_TESTS=true go test -v -tags=integration ./test/integration/ -run TestStatsFlow
```

### Ejecutar Test Individual

```bash
# Ejemplo: Solo test de login success
RUN_INTEGRATION_TESTS=true go test -v -tags=integration ./test/integration/ -run TestAuthFlow_LoginSuccess$
```

---

## 📁 Test Suite Structure

```
test/integration/
├── config.go                         ← Control de ejecución (ENV vars)
├── setup.go                          ← Testcontainers setup (PostgreSQL + MongoDB + RabbitMQ)
├── testhelpers.go                    ← Helpers: schema SQL, seeds, cleanup
├── auth_flow_test.go                 ← 3 tests Auth
├── material_flow_test.go             ← 4 tests Material
├── assessment_flow_test.go           ← 4 tests Assessment
├── progress_stats_flow_test.go       ← 6 tests Progress/Stats
├── example_test.go                   ← Ejemplos de uso
├── README.md                         ← Documentación general
├── README_TESTS.md                   ← Esta documentación (detalles de tests)
└── CONTROL_INTEGRATION_TESTS.md      ← Guía de control de ejecución
```

---

## 🧩 Test Details by Flow

### 1. Auth Flow (3 tests)

#### ✅ TestAuthFlow_LoginSuccess (~19s)
**Objetivo**: Verificar login exitoso con credenciales válidas

**Flujo**:
1. Seed usuario con email `test@edugo.com` y password `Test1234!`
2. POST a `/api/v1/auth/login` con credenciales válidas
3. Verificar status 200
4. Validar presencia de `access_token` y `refresh_token`
5. Validar datos del usuario en respuesta

**Assertions**:
- Status code: 200
- Response contiene: `access_token`, `refresh_token`, `expires_in`, `user`
- User data: `id`, `email`, `first_name`, `last_name`, `role`

#### ✅ TestAuthFlow_LoginInvalidCredentials (~19s)
**Objetivo**: Verificar rechazo de credenciales inválidas

**Flujo**:
1. Seed usuario con password correcto
2. POST con password incorrecto
3. Verificar status 401
4. Validar mensaje de error

**Assertions**:
- Status code: 401
- Error message: "invalid credentials"

#### ✅ TestAuthFlow_LoginNonexistentUser (~18s)
**Objetivo**: Verificar rechazo de usuario inexistente

**Flujo**:
1. POST con email que no existe en BD
2. Verificar status 401
3. Validar mensaje de error

**Assertions**:
- Status code: 401
- Error code: "UNAUTHORIZED"

---

### 2. Material Flow (4 tests)

#### ✅ TestMaterialFlow_CreateMaterial (~18s)
**Objetivo**: Crear material nuevo

**Flujo**:
1. Seed usuario autenticado
2. POST a `/api/v1/materials` con título y descripción
3. Verificar status 201
4. Validar estructura del material creado

**Assertions**:
- Status code: 201
- Response contiene: `id`, `title`, `description`, `author_id`, `status`
- Material existe en BD

#### ✅ TestMaterialFlow_GetMaterial (~18s)
**Objetivo**: Obtener material existente

**Flujo**:
1. Seed material en BD
2. GET a `/api/v1/materials/:id`
3. Verificar status 200
4. Validar datos del material

**Assertions**:
- Status code: 200
- Material ID coincide
- Todos los campos presentes

#### ✅ TestMaterialFlow_GetMaterialNotFound (~18s)
**Objetivo**: Manejar material inexistente (404)

**Flujo**:
1. GET a `/api/v1/materials/:id` con ID inexistente
2. Verificar status 404

**Assertions**:
- Status code: 404

#### ✅ TestMaterialFlow_ListMaterials (~18s)
**Objetivo**: Listar materiales del usuario

**Flujo**:
1. Seed 2 materiales para el usuario
2. GET a `/api/v1/materials`
3. Verificar status 200
4. Validar que retorna 2 materiales

**Assertions**:
- Status code: 200
- Count de materiales: 2
- Cada material tiene estructura completa

---

### 3. Assessment Flow (4 tests)

#### ✅ TestAssessmentFlow_GetAssessment (~18s)
**Objetivo**: Obtener assessment de un material

**Flujo**:
1. Seed material y assessment en MongoDB
2. GET a `/api/v1/materials/:id/assessment`
3. Verificar status 200
4. Validar 2 preguntas

**Assertions**:
- Status code: 200
- Assessment tiene `material_id` y `questions`
- Número de preguntas: 2

#### ✅ TestAssessmentFlow_GetAssessmentNotFound (~18s)
**Objetivo**: Manejar assessment inexistente (404)

**Flujo**:
1. GET a `/api/v1/materials/:id/assessment` con ID inexistente
2. Verificar status 404

**Assertions**:
- Status code: 404

#### ✅ TestAssessmentFlow_SubmitAssessment (~18s)
**Objetivo**: Enviar respuestas y recibir score + feedback

**Flujo**:
1. Seed assessment con 2 preguntas (respuestas: A, B)
2. POST a `/api/v1/assessments/:id/submit` con respuestas correctas
3. Verificar status 200
4. Validar score 100% (2/2 correctas)
5. Verificar feedback detallado

**Assertions**:
- Status code: 200
- Score: 100%
- Total questions: 2
- Correct answers: 2
- Feedback: 2 items con `is_correct: true`

#### ✅ TestAssessmentFlow_SubmitAssessmentDuplicate (~17s)
**Objetivo**: Validar manejo de submissions duplicadas

**Flujo**:
1. Primer submit: exitoso (200 OK)
2. Segundo submit: mismo assessment, mismo usuario
3. Verificar comportamiento (200 o 409 según config MongoDB)

**Assertions**:
- Primer submit: 200 OK
- Segundo submit: 200 o 409 (flexible para tests)

---

### 4. Progress Flow (4 tests)

#### ✅ TestProgressFlow_UpsertProgress (~18.4s)
**Objetivo**: Crear progreso inicial (upsert)

**Flujo**:
1. Seed usuario y material
2. PUT a `/api/v1/progress` con 50%, page 25
3. Verificar status 200
4. Validar datos del progreso

**Assertions**:
- Status code: 200
- Progress percentage: 50
- Last page: 25
- Message: "progress updated successfully"

#### ✅ TestProgressFlow_UpsertProgressUpdate (~20.3s)
**Objetivo**: Actualizar progreso existente (idempotencia)

**Flujo**:
1. Primer upsert: 30%, page 15 (200 OK)
2. Segundo upsert: 75%, page 38 (200 OK)
3. Verificar progreso actualizado a 75%

**Assertions**:
- Ambos upserts: 200 OK
- Progreso final: 75%, page 38
- Idempotencia funcionando ✅

#### ✅ TestProgressFlow_UpsertProgressUnauthorized (~19.4s)
**Objetivo**: Validar autorización (solo propio progreso)

**Flujo**:
1. Seed 2 usuarios: user1, user2
2. User1 autenticado intenta actualizar progreso de user2
3. Verificar status 403

**Assertions**:
- Status code: 403 Forbidden
- Error message: "you can only update your own progress"

#### ✅ TestProgressFlow_UpsertProgressInvalidData (~18.6s)
**Objetivo**: Validar datos de entrada

**Flujo**:
1. Test 1: Porcentaje inválido (150) → 400
2. Test 2: user_id faltante → 400

**Assertions**:
- Ambos casos: 400 Bad Request
- Validaciones de binding funcionando ✅

---

### 5. Stats Flow (2 tests)

#### ✅ TestStatsFlow_GetMaterialStats (~19.2s)
**Objetivo**: Obtener estadísticas de un material

**Flujo**:
1. Seed material
2. GET a `/api/v1/materials/:id/stats`
3. Verificar status 200
4. Validar estructura de stats

**Assertions**:
- Status code: 200
- Stats object no nulo

#### ✅ TestStatsFlow_GetGlobalStats (~18.0s)
**Objetivo**: Obtener estadísticas globales del sistema

**Flujo**:
1. Seed múltiples materiales
2. GET a `/api/v1/stats/global`
3. Verificar status 200
4. Validar estructura de global stats

**Assertions**:
- Status code: 200
- Global stats object no nulo

---

## 🔧 Helpers Disponibles

### Cleanup

```go
// Limpiar PostgreSQL
CleanDatabase(t, app.DB)

// Limpiar MongoDB
CleanMongoCollections(t, app.MongoDB)
```

### Seeds

```go
// Usuarios
userID, email := SeedTestUser(t, app.DB)
userID, email := SeedTestUserWithEmail(t, app.DB, "custom@email.com")

// Materiales
materialID := SeedTestMaterial(t, app.DB, authorID)
materialID := SeedTestMaterialWithTitle(t, app.DB, authorID, "Custom Title")

// Assessments
assessmentID := SeedTestAssessment(t, app.MongoDB, materialID)
```

### Setup

```go
// Setup completo de app de test
app := SetupTestApp(t)
defer app.Cleanup()

// Acceder a handlers
app.Container.Handlers.AuthHandler
app.Container.Handlers.MaterialHandler
app.Container.Handlers.AssessmentHandler
app.Container.Handlers.ProgressHandler
app.Container.Handlers.StatsHandler
```

---

## 🐳 Testcontainers

### Servicios Levantados

```yaml
PostgreSQL 15-alpine:
  Puerto: Random (asignado dinámicamente)
  Database: testdb
  User: testuser
  Password: testpass

MongoDB 7.0:
  Puerto: Random
  Database: edugo_test
  Collections: material_assessments, assessment_attempts, assessment_results

RabbitMQ 3.12-management-alpine:
  Puerto: Random
  Usuario: guest
  Password: guest
```

### Schema SQL Auto-Creado

```sql
- users (8 columnas)
- refresh_tokens (6 columnas)
- subjects (4 columnas)
- materials (11 columnas)
- material_versions (9 columnas)
- material_progress (9 columnas)
- assessments (4 columnas)
- student_classes (4 columnas)
- guardian_student_relation (5 columnas)
```

---

## 🎯 Best Practices

### Escribir Tests

1. **Arrange-Act-Assert pattern**
```go
// Arrange
app := SetupTestApp(t)
defer app.Cleanup()
CleanDatabase(t, app.DB)

// Act
req := httptest.NewRequest(...)
router.ServeHTTP(w, req)

// Assert
assert.Equal(t, http.StatusOK, w.Code)
```

2. **Usar Helpers**
```go
userID, _ := SeedTestUser(t, app.DB)
materialID := SeedTestMaterial(t, app.DB, userID)
```

3. **Cleanup Siempre**
```go
defer app.Cleanup()
```

4. **Logging para Debugging**
```go
t.Logf("Response status: %d", w.Code)
t.Logf("Response body: %s", w.Body.String())
```

---

## 🚨 Troubleshooting

### Tests Fallan con "Docker not available"

```bash
# Verificar Docker
docker ps

# Iniciar Docker Desktop
open /Applications/Docker.app
```

### Tests Timeout

```bash
# Aumentar timeout
go test -timeout 20m ...
```

### Tests Fallan Aleatoriamente

```bash
# Limpiar testcontainers
docker system prune -a

# Re-ejecutar tests
make test-integration
```

### Error "database error"

```bash
# Verificar schema SQL está sincronizado
# Ver: test/integration/testhelpers.go - initializeSchema()
```

---

## 📈 Performance

### Tiempos Típicos

```
Setup testcontainers: ~5-7s
Ejecución test individual: ~18-20s
Suite completa (17 tests): ~5-6 minutos
```

### Optimización

```bash
# Ejecutar en paralelo (experimental)
go test -v -tags=integration -parallel 2 ./test/integration/

# Ejecutar solo tests rápidos
go test -v -tags=integration -short ./test/integration/
```

---

## 🔗 Referencias

- [Testcontainers Go Documentation](https://golang.testcontainers.org/)
- [Gin Testing Guide](https://gin-gonic.com/docs/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)

---

## 📝 Changelog

### v1.0.0 (2025-11-06)
- ✅ 17 tests de integración E2E
- ✅ 100% cobertura de flujos críticos
- ✅ Testcontainers funcionando (PostgreSQL + MongoDB + RabbitMQ)
- ✅ Schema SQL robusto (9 tablas)
- ✅ Helpers reutilizables
- ✅ Documentación completa

---

**Mantenido por**: EduGo Team  
**Última actualización**: 2025-11-06  
**Estado**: ✅ Production Ready
