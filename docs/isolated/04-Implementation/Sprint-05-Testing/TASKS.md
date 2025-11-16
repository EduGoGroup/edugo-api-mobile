# Tareas del Sprint 05 - Testing Completo

## Objetivo
Crear suite completa de tests que garantice calidad y correctitud del Sistema de Evaluaciones.

---

## Tareas

### TASK-05-001: Tests Unitarios Dominio (>90%)
**Prioridad:** HIGH  
**Estimación:** 3h

#### Implementación
Completar tests de entities, value objects y asegurar >90% coverage.

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Ejecutar tests de dominio
go test ./internal/domain/... -v -cover

# Generar reporte
go test ./internal/domain/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

#### Criterios
- [ ] Coverage >90% en domain/entities
- [ ] Coverage >90% en domain/valueobjects
- [ ] Tests de todas las business rules
- [ ] Tests de validaciones

---

### TASK-05-002: Tests Integración con Testcontainers
**Prioridad:** HIGH  
**Estimación:** 4h

#### Implementación
Tests de repositorios con PostgreSQL y MongoDB reales.

```bash
go test ./tests/integration -v -tags=integration
```

#### Criterios
- [ ] Tests con PostgreSQL container
- [ ] Tests con MongoDB container
- [ ] Tests de transacciones ACID
- [ ] Contenedores limpios después

---

### TASK-05-003: Tests E2E Flujos Completos
**Prioridad:** HIGH  
**Estimación:** 4h

#### Implementación
Tests de flujos completos usando httptest.

Archivo: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/tests/e2e/assessment_flow_test.go`

```go
//go:build e2e

func TestE2E_CompleteAssessmentFlow(t *testing.T) {
    // 1. Obtener assessment
    // 2. Responder preguntas
    // 3. Obtener calificación
    // 4. Verificar historial
}

func TestE2E_MultipleAttempts(t *testing.T) {
    // Verificar que se pueden hacer múltiples intentos
}

func TestE2E_AttemptsLimit(t *testing.T) {
    // Verificar límite de intentos
}
```

#### Criterios
- [ ] 5+ flujos E2E cubiertos
- [ ] Tests usan API HTTP real

---

### TASK-05-004: Tests de Seguridad
**Prioridad:** HIGH  
**Estimación:** 2h

#### Implementación
```go
func TestSecurity_CorrectAnswersNeverExposed(t *testing.T) {
    // Verificar TODOS los endpoints
    // Ninguno debe retornar correct_answer
}

func TestSecurity_ScoreValidatedServerSide(t *testing.T) {
    // Cliente no puede mentir sobre score
}
```

#### Criterios
- [ ] Respuestas correctas NUNCA expuestas
- [ ] Score validado en servidor
- [ ] JWT requerido en todos los endpoints

---

### TASK-05-005: Tests de Performance
**Prioridad:** MEDIUM  
**Estimación:** 2h

#### Implementación
```go
func BenchmarkGetAssessment(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // GET /v1/materials/:id/assessment
    }
}

func BenchmarkCreateAttempt(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // POST /v1/materials/:id/assessment/attempts
    }
}
```

#### Criterios
- [ ] GET assessment <500ms p95
- [ ] POST attempt <2000ms p95

```bash
go test ./tests/benchmark -bench=. -benchmem
```

---

**Sprint:** 05/06
