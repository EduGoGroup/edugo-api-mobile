# Estrategia de Testing
# Sistema de Evaluaciones - EduGo

**Versión:** 1.0.0  
**Fecha:** 14 de Noviembre, 2025

---

## 1. PIRÁMIDE DE TESTING

```
           /\
          /E2E\         10% - Tests End-to-End (flujos completos HTTP)
         /------\
        /  INT   \      20% - Tests de Integración (Testcontainers)
       /----------\
      /   UNIT     \    70% - Tests Unitarios (dominio puro)
     /--------------\
```

### Distribución Objetivo

| Tipo | Porcentaje | Cantidad Estimada | Tiempo Ejecución |
|------|------------|-------------------|------------------|
| **Unitarios** | 70% | ~100 tests | <5s |
| **Integración** | 20% | ~30 tests | <30s |
| **E2E** | 10% | ~15 tests | <2min |
| **TOTAL** | 100% | ~145 tests | <3min |

---

## 2. ESTRATEGIA DE COVERAGE

### Objetivos de Coverage por Capa

| Capa | Coverage Objetivo | Justificación |
|------|------------------|---------------|
| **Domain (entities, value objects)** | >90% | Lógica de negocio crítica |
| **Application (services)** | >85% | Orquestación importante |
| **Infrastructure (repos)** | >70% | Código más simple, tests caros |
| **Handlers (HTTP)** | >80% | API pública, tests E2E |
| **GLOBAL** | >80% | Estándar de industria |

### Comandos de Medición
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Coverage global
go test ./... -cover -coverprofile=coverage.out

# Coverage por paquete
go tool cover -func=coverage.out

# Coverage específico por capa
go test ./internal/domain/... -cover
go test ./internal/application/... -cover
go test ./internal/infrastructure/... -cover -tags=integration
```

---

## 3. HERRAMIENTAS DE TESTING

### 3.1 Framework de Testing

**Testify** (github.com/stretchr/testify v1.8.4)
- `assert` - Assertions que no detienen test
- `require` - Assertions que detienen test inmediatamente
- `suite` - Test suites con setup/teardown
- `mock` - Mocking (para casos específicos)

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
    result, err := SomeFunction()
    
    require.NoError(t, err) // Detiene test si falla
    assert.Equal(t, expected, result) // Continúa si falla
}
```

### 3.2 Testcontainers

**Testcontainers-go** (v0.27.0)
- Levanta PostgreSQL y MongoDB en Docker
- Tests de integración con BDs reales
- Cleanup automático

```go
import (
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestIntegration(t *testing.T) {
    ctx := context.Background()
    
    pgContainer, _ := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15-alpine"),
    )
    defer pgContainer.Terminate(ctx)
    
    // Tests contra PostgreSQL real
}
```

### 3.3 Coverage Reporting

**go test -cover**
```bash
# Generar coverage
go test ./... -coverprofile=coverage.out

# Reporte HTML
go tool cover -html=coverage.out -o coverage.html

# CI/CD: Codecov
# Upload automático desde GitHub Actions
```

---

## 4. TIPOS DE TESTS POR CAPA

### 4.1 Domain Layer (Unitarios Puros)

**Qué testear:**
- Validaciones de entities
- Business rules (CanAttempt, IsPassed, etc.)
- Cálculo de scores
- Value objects

**Características:**
- ✅ Sin dependencias externas (no DB, no HTTP)
- ✅ Fast (<5s para todos)
- ✅ Determinísticos
- ✅ Coverage >90%

**Ejemplo:**
```go
func TestAssessment_CanAttempt(t *testing.T) {
    assessment := &entities.Assessment{MaxAttempts: &3}
    
    assert.True(t, assessment.CanAttempt(0))
    assert.True(t, assessment.CanAttempt(2))
    assert.False(t, assessment.CanAttempt(3))
}
```

### 4.2 Application Layer (Mocks de Repositorios)

**Qué testear:**
- Orquestación de servicios
- Lógica de sanitización (remover correct_answers)
- Error handling

**Características:**
- ⚠️ Usa mocks de repositorios
- ✅ Fast (<10s)
- ✅ Coverage >85%

**Ejemplo:**
```go
func TestAssessmentService_GetAssessment(t *testing.T) {
    mockRepo := new(MockAssessmentRepo)
    mockRepo.On("FindByMaterialID", mock.Anything).Return(assessment, nil)
    
    service := NewAssessmentService(mockRepo, mockMongoRepo)
    
    result, err := service.GetAssessmentByMaterialID(ctx, materialID)
    
    require.NoError(t, err)
    // Verificar que correct_answer fue removido
    assert.NotContains(t, result, "correct_answer")
}
```

### 4.3 Infrastructure Layer (Testcontainers)

**Qué testear:**
- Queries SQL/MongoDB reales
- Transacciones ACID
- Connection pooling
- Error handling de BD

**Características:**
- ⚠️ Requiere Docker (Testcontainers)
- ⚠️ Slower (~30s)
- ✅ Tests realistas
- ✅ Coverage >70%

**Tag:**
```go
//go:build integration
```

### 4.4 HTTP Handlers (Tests E2E)

**Qué testear:**
- Flujos completos (GET assessment → POST attempt → GET results)
- Validación de inputs HTTP
- Status codes correctos
- Middleware (auth, CORS)

**Características:**
- ⚠️ Requiere API corriendo
- ⚠️ Slower (~1-2min)
- ✅ Valida integración completa

**Tag:**
```go
//go:build e2e
```

**Ejemplo:**
```go
func TestE2E_AssessmentFlow(t *testing.T) {
    router := setupTestRouter()
    
    // 1. GET assessment
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/v1/materials/"+id+"/assessment", nil)
    router.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)
    
    // 2. POST attempt
    // 3. Verify score
}
```

---

## 5. CI/CD INTEGRATION

### 5.1 GitHub Actions

**Workflow de Tests:**
```yaml
# .github/workflows/ci.yml
jobs:
  test:
    steps:
      - name: Unit Tests
        run: go test ./internal/domain/... -v -cover
      
      - name: Integration Tests
        run: go test ./tests/integration -v -tags=integration
        
      - name: E2E Tests
        run: go test ./tests/e2e -v -tags=e2e
      
      - name: Upload Coverage
        uses: codecov/codecov-action@v3
```

### 5.2 Coverage Reporting

**Codecov Integration:**
- Coverage badge en README
- Trend de coverage en PRs
- Alertas si coverage baja

---

## 6. ESTRATEGIA POR FUNCIONALIDAD

### RF-001: GET /v1/materials/:id/assessment

**Tests Unitarios (Domain):**
- ✅ Assessment.Validate()
- ✅ Sanitización de respuestas correctas

**Tests Integración:**
- ✅ Query PostgreSQL por material_id
- ✅ Query MongoDB con proyección (sin correct_answer)

**Tests E2E:**
- ✅ Flujo HTTP completo con JWT
- ✅ Verificar que response NO tiene correct_answer
- ✅ 404 si material no existe

### RF-002: POST /v1/materials/:id/assessment/attempts

**Tests Unitarios:**
- ✅ Attempt.NewAttempt() calcula score correcto
- ✅ Validaciones de answers

**Tests Integración:**
- ✅ Transacción ACID (attempt + answers atómicas)
- ✅ Score guardado correctamente

**Tests E2E:**
- ✅ POST con respuestas, obtener score
- ✅ Verificar límite de intentos
- ✅ Score calculado en servidor (no en cliente)

### RF-003: GET /v1/attempts/:id/results

**Tests E2E:**
- ✅ Usuario solo ve sus propios intentos
- ✅ 403 si intenta ver intento de otro usuario

### RF-004: GET /v1/users/me/attempts

**Tests E2E:**
- ✅ Paginación funciona (limit/offset)
- ✅ Ordenamiento por fecha DESC

---

## 7. TESTS DE SEGURIDAD

### Casos Críticos

**SEC-001: Respuestas Correctas NUNCA Expuestas**
```go
func TestSecurity_CorrectAnswersNeverExposed(t *testing.T) {
    // Probar TODOS los endpoints
    endpoints := []string{
        "GET /v1/materials/:id/assessment",
        "POST /v1/materials/:id/assessment/attempts",
        "GET /v1/attempts/:id/results",
        "GET /v1/users/me/attempts",
    }
    
    for _, endpoint := range endpoints {
        response := callEndpoint(endpoint)
        assertNotContains(t, response, "correct_answer")
        assertNotContains(t, response, "feedback")
    }
}
```

**SEC-002: Score Validado en Servidor**
```go
func TestSecurity_ScoreServerSideValidation(t *testing.T) {
    // Cliente envía respuestas
    // Servidor calcula score comparando con MongoDB
    // Ignorar cualquier "score" que cliente envíe
}
```

**SEC-003: JWT Requerido**
```go
func TestSecurity_JWTRequired(t *testing.T) {
    // Request sin JWT → 401 Unauthorized
}
```

---

## 8. TESTS DE PERFORMANCE

### Benchmarks

```go
func BenchmarkGetAssessment(b *testing.B) {
    for i := 0; i < b.N; i++ {
        service.GetAssessmentByMaterialID(ctx, materialID)
    }
}

func BenchmarkCreateAttempt(b *testing.B) {
    for i := 0; i < b.N; i++ {
        service.CreateAttempt(ctx, studentID, answers)
    }
}
```

**Ejecutar:**
```bash
go test ./... -bench=. -benchmem
```

**Objetivos:**
- GET assessment: <500ms p95
- POST attempt: <2000ms p95 (incluye validación MongoDB + transacción PostgreSQL)

---

## 9. COMANDOS ÚTILES

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Todos los tests
go test ./... -v

# Solo unitarios (fast)
go test ./internal/domain/... -v

# Solo integración (requiere Docker)
go test ./tests/integration -v -tags=integration

# Solo E2E
go test ./tests/e2e -v -tags=e2e

# Coverage global
go test ./... -cover

# Coverage detallado
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out

# Tests en paralelo
go test ./... -v -parallel 4

# Tests con timeout
go test ./... -v -timeout 5m

# Benchmarks
go test ./... -bench=. -benchmem
```

---

**Generado con:** Claude Code  
**Tokens:** ~3K  
**Estado:** Estrategia de Testing Completa
