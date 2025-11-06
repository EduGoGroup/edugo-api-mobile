# 🧪 Informe 3: Estado de Tests y Plan de Mejora

**Fecha**: 2025-11-06  
**Analista**: Claude Code  
**Scope**: Cobertura actual + Estrategia de tests de integración + Testcontainers

---

## 🎯 Resumen Ejecutivo

**Estado de Tests**: 🟡 Tests unitarios excelentes, integración faltante

**Cobertura Actual**:
- **Tests unitarios**: 89 tests (100% passing) ⭐⭐⭐⭐⭐
- **Cobertura código nuevo**: ≥85% ⭐⭐⭐⭐⭐
- **Cobertura total**: 25.5% ⭐⭐☆☆☆
- **Tests integración**: 0 tests ejecutándose ⭐☆☆☆☆
- **Tests E2E**: 0 tests ⭐☆☆☆☆

**Veredicto**: Proyecto con buena base de tests unitarios, pero **crítica falta de tests de integración**.

---

## 1. Análisis de Cobertura Actual

### 1.1. Tests Unitarios Existentes

**Total**: 89 tests (100% passing ✅)

| Componente | Tests | Coverage | Calidad |
|------------|-------|----------|---------|
| **Scoring Strategies** | 52 | ~95% | ⭐⭐⭐⭐⭐ |
| - MultipleChoice | 7 | 100% | ⭐⭐⭐⭐⭐ |
| - TrueFalse | 24 | 100% | ⭐⭐⭐⭐⭐ |
| - ShortAnswer | 21 | 100% | ⭐⭐⭐⭐⭐ |
| **Services** | 28 | ~85-90% | ⭐⭐⭐⭐⭐ |
| - MaterialService | 5 | 90% | ⭐⭐⭐⭐⭐ |
| - AssessmentService | 7 | 90% | ⭐⭐⭐⭐⭐ |
| - ProgressService | 9 | 95% | ⭐⭐⭐⭐⭐ |
| - StatsService | 6 | 100% | ⭐⭐⭐⭐⭐ |
| **Handlers** | 9 | ~95% | ⭐⭐⭐⭐⭐ |
| - AssessmentHandler | 9 | 95% | ⭐⭐⭐⭐⭐ |

**Hallazgos Positivos**:
- ✅ Cobertura excelente en código nuevo (≥85%)
- ✅ Todos los tests pasando (0 failures)
- ✅ Uso de mocks correcto (testify/mock)
- ✅ Tests bien estructurados (tabla-driven)
- ✅ Edge cases cubiertos

**Hallazgos Negativos**:
- ❌ Cobertura total baja (25.5%) por código legacy
- ❌ No hay tests para repositories
- ❌ No hay tests para handlers antiguos (se eliminarán)

### 1.2. Tests de Integración

**Estado**: ⚠️ Configurados pero NO ejecutándose

**Archivos existentes**:
```
test/integration/
├── mongodb_test.go      (básico, conexión)
├── postgres_test.go     (skipped con t.Skip)
└── rabbitmq_test.go     (básico, conexión)
```

**Problema identificado**:
```go
// postgres_test.go
func TestPostgresConnection(t *testing.T) {
    t.Skip("Requiere testcontainers")  ← SKIPPED
    // ...
}
```

**Razón**: Testcontainers configurado en `go.mod` pero no usado.

**Dependencias disponibles**:
```go
// go.mod
github.com/testcontainers/testcontainers-go v0.39.0
github.com/testcontainers/testcontainers-go/modules/mongodb v0.39.0
github.com/testcontainers/testcontainers-go/modules/postgres v0.39.0
github.com/testcontainers/testcontainers-go/modules/rabbitmq v0.39.0
```

✅ **Testcontainers ya está instalado**, solo falta implementar tests.

### 1.3. Tests E2E

**Estado**: ❌ No existen

**Carpeta**: `test/unit/` está vacía

---

## 2. Estrategia de Tests de Integración

### 2.1. Objetivos

1. **Cobertura de flujos críticos** end-to-end
2. **Validación de integraciones** reales (PostgreSQL, MongoDB, RabbitMQ)
3. **Aislamiento completo** (contenedores efímeros)
4. **Reproducibilidad** (sin dependencias de servicios externos)
5. **CI/CD friendly** (ejecutables en GitHub Actions)

### 2.2. Alcance de Tests Propuestos

#### 🔴 Críticos (Prioridad 1)

1. **Auth Flow**
   - Login → Genera tokens → Acceso a recursos protegidos
   - Refresh token → Nuevo access token
   - Logout → Tokens revocados

2. **Material Flow**
   - Crear material → Guardar en PostgreSQL
   - Subir a S3 (mock) → Publicar evento RabbitMQ
   - Consultar material con versiones → LEFT JOIN correcto

3. **Assessment Flow**
   - Obtener assessment de MongoDB
   - Enviar respuestas → Calcular puntaje con Strategy Pattern
   - Validar feedback generado correctamente
   - Verificar persistencia en MongoDB

#### 🟡 Importantes (Prioridad 2)

4. **Progress Flow**
   - Actualizar progreso → UPSERT sin duplicados
   - Múltiples updates → Idempotencia verificada
   - Completar material (100%) → Flag is_completed

5. **Stats Flow**
   - Consultar estadísticas globales
   - Validar queries paralelas (5 simultáneas)
   - Verificar agregaciones correctas

#### 🟢 Opcionales (Prioridad 3)

6. **Summary Flow** (si aplica)
7. **Error Handling Flow**
8. **Concurrency Flow** (race conditions)

---

## 3. Plan de Implementación con Testcontainers

### 3.1. Arquitectura de Tests Propuesta

```
test/integration/
├── testcontainers/
│   ├── setup.go           ← Setup compartido
│   ├── postgres.go        ← Contenedor PostgreSQL
│   ├── mongodb.go         ← Contenedor MongoDB
│   ├── rabbitmq.go        ← Contenedor RabbitMQ
│   └── s3mock.go          ← Mock de S3 (opcional)
│
├── auth_flow_test.go      ← Tests de autenticación
├── material_flow_test.go  ← Tests de materiales
├── assessment_flow_test.go← Tests de evaluaciones
├── progress_flow_test.go  ← Tests de progreso
└── stats_flow_test.go     ← Tests de estadísticas
```

### 3.2. Setup Compartido (testcontainers/setup.go)

**Estrategia**: Contenedores compartidos para toda la suite de tests

```go
package testcontainers

import (
    "context"
    "testing"
    "github.com/testcontainers/testcontainers-go"
    postgresContainer "github.com/testcontainers/testcontainers-go/modules/postgres"
    mongoContainer "github.com/testcontainers/testcontainers-go/modules/mongodb"
    rabbitmqContainer "github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

type TestContainers struct {
    Postgres *postgresContainer.PostgresContainer
    MongoDB  *mongoContainer.MongoDBContainer
    RabbitMQ *rabbitmqContainer.RabbitMQContainer
    
    PostgresURI string
    MongoURI    string
    RabbitURI   string
}

// SetupContainers inicia todos los contenedores necesarios
// Se ejecuta UNA VEZ por suite de tests
func SetupContainers(ctx context.Context) (*TestContainers, error) {
    tc := &TestContainers{}
    
    // PostgreSQL
    pgContainer, err := postgresContainer.Run(ctx,
        "postgres:16-alpine",
        postgresContainer.WithDatabase("edugo_test"),
        postgresContainer.WithUsername("test_user"),
        postgresContainer.WithPassword("test_pass"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(30*time.Second),
        ),
    )
    if err != nil {
        return nil, err
    }
    tc.Postgres = pgContainer
    
    // Obtener URI de conexión
    tc.PostgresURI, err = pgContainer.ConnectionString(ctx, "sslmode=disable")
    if err != nil {
        return nil, err
    }
    
    // MongoDB
    mongoContainer, err := mongoContainer.Run(ctx,
        "mongo:7",
        mongoContainer.WithUsername("test_admin"),
        mongoContainer.WithPassword("test_pass"),
    )
    if err != nil {
        return nil, err
    }
    tc.MongoDB = mongoContainer
    tc.MongoURI, err = mongoContainer.ConnectionString(ctx)
    
    // RabbitMQ
    rabbitContainer, err := rabbitmqContainer.Run(ctx,
        "rabbitmq:3.12-alpine",
    )
    if err != nil {
        return nil, err
    }
    tc.RabbitMQ = rabbitContainer
    tc.RabbitURI, err = rabbitContainer.AmqpURL(ctx)
    
    return tc, nil
}

// TeardownContainers detiene todos los contenedores
func (tc *TestContainers) TeardownContainers(ctx context.Context) error {
    if err := tc.Postgres.Terminate(ctx); err != nil {
        return err
    }
    if err := tc.MongoDB.Terminate(ctx); err != nil {
        return err
    }
    if err := tc.RabbitMQ.Terminate(ctx); err != nil {
        return err
    }
    return nil
}
```

**Características clave**:
- ✅ Contenedores compartidos (no uno por test)
- ✅ Setup UNA VEZ (rápido)
- ✅ Teardown automático
- ✅ Puertos aleatorios (no colisiones)
- ✅ Timeout configurado (no espera infinita)

### 3.3. Ejemplo: Auth Flow Test

```go
package integration

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "your-project/test/integration/testcontainers"
)

var (
    testContainers *testcontainers.TestContainers
    container      *container.Container
)

// TestMain se ejecuta UNA VEZ antes de todos los tests
func TestMain(m *testing.M) {
    ctx := context.Background()
    
    // Setup contenedores
    tc, err := testcontainers.SetupContainers(ctx)
    if err != nil {
        log.Fatalf("Failed to setup containers: %v", err)
    }
    testContainers = tc
    
    // Ejecutar migraciones
    if err := runMigrations(tc.PostgresURI); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }
    
    // Inicializar container DI con URIs de test
    container = initializeTestContainer(
        tc.PostgresURI,
        tc.MongoURI,
        tc.RabbitURI,
    )
    
    // Ejecutar tests
    code := m.Run()
    
    // Cleanup
    tc.TeardownContainers(ctx)
    os.Exit(code)
}

func TestAuthFlow_CompleteLogin(t *testing.T) {
    ctx := context.Background()
    
    // 1. Crear usuario de prueba
    user := createTestUser(t, container.UserRepository)
    
    // 2. Login
    loginReq := dto.LoginRequest{
        Email:    user.Email,
        Password: "test_password",
    }
    loginResp, err := container.AuthService.Login(ctx, loginReq)
    require.NoError(t, err)
    require.NotEmpty(t, loginResp.AccessToken)
    require.NotEmpty(t, loginResp.RefreshToken)
    
    // 3. Verificar que access token es válido
    claims, err := container.JWTManager.ValidateToken(loginResp.AccessToken)
    require.NoError(t, err)
    assert.Equal(t, user.ID.String(), claims.UserID)
    
    // 4. Verificar refresh token está en BD
    tokenHash := auth.HashToken(loginResp.RefreshToken)
    token, err := container.RefreshTokenRepository.FindByTokenHash(ctx, tokenHash)
    require.NoError(t, err)
    assert.NotNil(t, token)
    assert.Equal(t, user.ID, token.UserID)
    
    // 5. Usar access token para acceder a recurso protegido
    material, err := container.MaterialService.GetMaterial(ctx, "some-id")
    // Validar que funciona...
}

func TestAuthFlow_RefreshToken(t *testing.T) {
    // Test de refresh token rotation...
}

func TestAuthFlow_Logout(t *testing.T) {
    ctx := context.Background()
    
    // 1. Login
    user := createTestUser(t, container.UserRepository)
    loginResp, _ := container.AuthService.Login(ctx, dto.LoginRequest{...})
    
    // 2. Logout
    err := container.AuthService.Logout(ctx, user.ID.String(), loginResp.RefreshToken)
    require.NoError(t, err)
    
    // 3. Verificar refresh token revocado
    tokenHash := auth.HashToken(loginResp.RefreshToken)
    token, _ := container.RefreshTokenRepository.FindByTokenHash(ctx, tokenHash)
    assert.NotNil(t, token.RevokedAt)  ← Debe estar revocado
}
```

**Características**:
- ✅ `TestMain` setup UNA VEZ
- ✅ Contenedores compartidos entre tests
- ✅ Base de datos real (no mocks)
- ✅ Validación completa de flujo
- ✅ Tests independientes (cleanup entre tests si necesario)

### 3.4. Estrategia de Datos de Prueba

**Problema**: Tests necesitan datos iniciales

**Solución 1: Fixtures SQL**:
```go
func runMigrations(dbURI string) error {
    // 1. Ejecutar migraciones normales
    // 2. Ejecutar fixtures de prueba
    db, _ := sql.Open("postgres", dbURI)
    defer db.Close()
    
    fixtures := []string{
        "INSERT INTO users ...",
        "INSERT INTO materials ...",
        // ...
    }
    
    for _, query := range fixtures {
        _, err := db.Exec(query)
        if err != nil {
            return err
        }
    }
    return nil
}
```

**Solución 2: Factories de Test**:
```go
func createTestUser(t *testing.T, repo repository.UserRepository) *entity.User {
    user := &entity.User{
        Email:    "test@example.com",
        Password: "hashed_password",
        // ...
    }
    err := repo.Create(context.Background(), user)
    require.NoError(t, err)
    return user
}
```

**Recomendación**: Usar **Factories** (más flexible).

### 3.5. Gestión de Base de Datos entre Tests

**Estrategia Recomendada**: **Transacciones por test**

```go
func TestWithTransaction(t *testing.T) {
    ctx := context.Background()
    
    // Iniciar transacción
    tx, err := container.DB.Begin()
    require.NoError(t, err)
    defer tx.Rollback()  ← Rollback al final (cleanup automático)
    
    // Usar repositorios con la transacción
    repo := postgresRepo.NewPostgresUserRepository(tx)
    
    // Test...
    user := createTestUser(t, repo)
    
    // Al terminar, rollback automático
    // (siguiente test tendrá BD limpia)
}
```

**Beneficios**:
- ✅ BD limpia entre tests
- ✅ Aislamiento perfecto
- ✅ Rápido (no recrear contenedores)

**Alternativa**: Truncar tablas entre tests (más lento).

---

## 4. Plan de Trabajo Detallado

### Fase 1: Setup Base (4 horas)

**Tareas**:
1. Crear `test/integration/testcontainers/setup.go` (2h)
   - SetupContainers()
   - Contenedores PostgreSQL, MongoDB, RabbitMQ
   - Gestión de puertos y URIs
   
2. Configurar TestMain en cada archivo de test (1h)
   - Reutilizar contenedores
   - Ejecutar migraciones
   - Cleanup al final

3. Crear factories de datos de prueba (1h)
   - `createTestUser()`
   - `createTestMaterial()`
   - `createTestAssessment()`

**Entregable**: Infraestructura base funcionando

### Fase 2: Tests Críticos (6 horas)

**Tareas**:
1. `auth_flow_test.go` (2h)
   - TestAuthFlow_CompleteLogin
   - TestAuthFlow_RefreshToken
   - TestAuthFlow_Logout
   - TestAuthFlow_RevokeAllSessions

2. `material_flow_test.go` (2h)
   - TestMaterialFlow_CreateAndRetrieve
   - TestMaterialFlow_WithVersions
   - TestMaterialFlow_UploadComplete

3. `assessment_flow_test.go` (2h)
   - TestAssessmentFlow_SubmitAndScore
   - TestAssessmentFlow_DetailedFeedback
   - TestAssessmentFlow_DuplicatePrevention

**Entregable**: Flujos críticos cubiertos

### Fase 3: Tests Importantes (4 horas)

**Tareas**:
1. `progress_flow_test.go` (2h)
   - TestProgressFlow_UpsertIdempotency
   - TestProgressFlow_CompleteMarking
   - TestProgressFlow_MultipleUpdates

2. `stats_flow_test.go` (2h)
   - TestStatsFlow_GlobalStats
   - TestStatsFlow_ParallelQueries
   - TestStatsFlow_EmptySystem

**Entregable**: Flujos importantes cubiertos

### Fase 4: CI/CD Integration (2 horas)

**Tareas**:
1. Configurar GitHub Actions workflow
2. Makefile targets para tests de integración
3. Documentación de cómo ejecutar tests

**Entregable**: Tests ejecutables en CI/CD

### Resumen del Plan

| Fase | Esfuerzo | Prioridad | Bloqueante |
|------|----------|-----------|------------|
| Fase 1: Setup | 4h | 🔴 | Sí (para todo) |
| Fase 2: Críticos | 6h | 🔴 | Sí (para deploy) |
| Fase 3: Importantes | 4h | 🟡 | No |
| Fase 4: CI/CD | 2h | 🟡 | No |
| **TOTAL** | **16h** | - | - |

---

## 5. Comandos de Makefile Propuestos

```makefile
# Agregar al Makefile existente:

test-integration-setup: ## Setup testcontainers (primera vez)
	@echo "🐳 Verificando Docker..."
	@docker ps > /dev/null || (echo "❌ Docker no está corriendo" && exit 1)
	@echo "✅ Docker listo"

test-integration-run: ## Ejecutar tests de integración
	@echo "🧪 Ejecutando tests de integración..."
	@go test -v -tags=integration ./test/integration/... -timeout 10m
	@echo "✅ Tests completados"

test-integration-coverage: ## Tests integración con coverage
	@echo "📊 Tests integración con coverage..."
	@mkdir -p coverage
	@go test -tags=integration -coverprofile=coverage/integration.out \
		-covermode=atomic ./test/integration/... -timeout 10m
	@go tool cover -html=coverage/integration.out -o coverage/integration.html
	@echo "✅ Reporte: coverage/integration.html"

test-integration-watch: ## Watch mode para tests integración
	@echo "👀 Watch mode activado..."
	@find test/integration -name "*.go" | entr -c go test -v -tags=integration ./test/integration/...

test-all: test test-integration-run ## Ejecutar TODOS los tests (unit + integration)
	@echo "✅ Todos los tests completados"

docker-check: ## Verificar Docker disponible
	@docker ps > /dev/null || (echo "❌ Iniciar Docker Desktop" && exit 1)
	@echo "✅ Docker disponible"
```

**Uso**:
```bash
# Primera vez
make test-integration-setup

# Ejecutar tests
make test-integration-run

# Con coverage
make test-integration-coverage

# Todos los tests
make test-all
```

---

## 6. Consideraciones Importantes

### 6.1. Performance

**Problema**: Tests de integración son lentos

**Solución**:
- ✅ Contenedores compartidos (no uno por test)
- ✅ Setup UNA VEZ en TestMain
- ✅ Transacciones para aislamiento (rápido)
- ✅ Fixtures mínimos

**Tiempo estimado**:
- Setup contenedores: ~10s
- Cada test: ~100-500ms
- Suite completa: ~2-3 min

### 6.2. CI/CD

**GitHub Actions** (ejemplo):
```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  integration-tests:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Run integration tests
      run: make test-integration-run
```

**Beneficios**:
- ✅ Tests ejecutan en cada PR
- ✅ No requiere servicios externos
- ✅ Testcontainers funciona en GitHub Actions

### 6.3. Docker Requirement

**⚠️ Importante**: Testcontainers requiere Docker

**Solución para desarrolladores**:
- Docker Desktop debe estar corriendo
- Agregar `make docker-check` antes de tests
- Documentar en README

### 6.4. Recursos de Sistema

**Contenedores necesitan**:
- CPU: ~1-2 cores
- RAM: ~2-4 GB
- Disk: ~1 GB

**Recomendación**: Máquinas con ≥8GB RAM.

---

## 7. Métricas de Éxito

### Objetivos Post-Implementación

| Métrica | Actual | Objetivo | Delta |
|---------|--------|----------|-------|
| **Tests unitarios** | 89 | 100+ | +11 |
| **Tests integración** | 0 | 15+ | +15 |
| **Cobertura código nuevo** | 85% | 85%+ | - |
| **Cobertura total** | 25.5% | 40%+ | +14.5% |
| **Tiempo de ejecución** | 10s | 3 min | - |

### KPIs de Calidad

- ✅ **100% de flujos críticos** con tests integración
- ✅ **0 tests skipped** en suite integración
- ✅ **<5 min** tiempo total de ejecución
- ✅ **Ejecutable en CI/CD** sin servicios externos

---

## 8. Priorización de Tests

### 🔴 Implementar YA (Bloqueantes para Producción)

1. **Auth Flow**: Login, refresh, logout
2. **Assessment Flow**: Submit, scoring, feedback
3. **Material Flow**: CRUD básico

**Razón**: Funcionalidades core que deben funcionar perfectamente.

### 🟡 Implementar Próximo Sprint

4. **Progress Flow**: UPSERT, idempotencia
5. **Stats Flow**: Queries paralelas
6. **RabbitMQ Integration**: Eventos publicados correctamente

**Razón**: Importantes pero no críticas para MVP.

### 🟢 Backlog

7. **Error Handling**: Manejo de errores edge cases
8. **Performance**: Tests de carga
9. **Security**: Tests de seguridad

**Razón**: Nice to have, no bloqueantes.

---

## 9. Recomendaciones Finales

### Para el Equipo

1. **Empezar con Fase 1 (Setup)** completa antes de tests
2. **Un desarrollador** dedicado a infraestructura de tests
3. **Code review** estricto de tests (calidad igual que código)
4. **Documentar** ejemplos de cómo escribir tests de integración

### Para el Proyecto

1. **No mergear** código sin tests de integración en próximos PRs
2. **Ejecutar tests** en CI/CD obligatoriamente
3. **Aumentar cobertura** gradualmente (objetivo 40% → 60% → 80%)
4. **Refactorizar** código legacy agregando tests

---

## 10. Conclusión

### Estado Actual

**Fortalezas**:
- ✅ Tests unitarios excelentes (89 tests)
- ✅ Cobertura alta en código nuevo (≥85%)
- ✅ Testcontainers ya configurado

**Debilidades**:
- ❌ Sin tests de integración ejecutándose
- ❌ Cobertura total baja (25.5%)
- ❌ Código legacy sin tests

### Plan de Acción

```
Fase 1 (4h)  →  Fase 2 (6h)  →  Fase 3 (4h)  →  Fase 4 (2h)
   Setup     →   Críticos    →  Importantes  →   CI/CD
                                                 
Total: 16 horas de trabajo
```

### Veredicto

**Estado de Tests**: 🟡 Bueno pero **incompleto**

**Acción Requerida**: Implementar tests de integración es **crítico** antes de producción.

**Prioridad**: 🔴 **ALTA** (después de FASE 3 del Plan Maestro)

---

**Siguiente Paso**: Ver `04-resumen-ejecutivo.md` para consolidación y plan final.
