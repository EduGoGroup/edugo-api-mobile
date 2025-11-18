# Sprint 05-B: Testing Avanzado
# Sistema de Evaluaciones - EduGo

**Duración:** 2-3 días  
**Objetivo:** Coverage >80% y tests de performance.

---

## 🎯 Objetivo

Alcanzar coverage >80% y agregar tests de performance:
- 📊 Aumentar coverage de services (61.7% → >80%)
- 📊 Aumentar coverage de handlers (57.1% → >70%)
- 📊 Agregar tests de repositorios (0% → >70%)
- ⚡ Tests de performance (benchmarks)
- 📈 Coverage global >80%

---

## 📋 Tareas del Sprint 05-B

### TASK-05-B-001: Tests de Services (aumentar a >80%)
**Prioridad:** HIGH  
**Estimación:** 4h  
**Coverage Actual:** 61.7%  
**Coverage Objetivo:** >80%

Tests faltantes:
- Edge cases de AssessmentService
- Error handling completo
- Validaciones de negocio

---

### TASK-05-B-002: Tests de Handlers (aumentar a >70%)
**Prioridad:** MEDIUM  
**Estimación:** 3h  
**Coverage Actual:** 57.1%  
**Coverage Objetivo:** >70%

Tests faltantes:
- Todos los códigos de error HTTP
- Validación de request body
- Edge cases de paginación

---

### TASK-05-B-003: Tests de Repositorios
**Prioridad:** HIGH  
**Estimación:** 6h  
**Coverage Actual:** 0%  
**Coverage Objetivo:** >70%

Repositorios a testear:
- `AssessmentRepository` (PostgreSQL)
- `AttemptRepository` (PostgreSQL)
- `AnswerRepository` (PostgreSQL)
- `AssessmentDocumentRepository` (MongoDB)

**Nota:** Usar testcontainers (ya configurados en Sprint 05-A).

---

### TASK-05-B-004: Tests de Performance
**Prioridad:** MEDIUM  
**Estimación:** 3h

Benchmarks a crear:
- `BenchmarkGetAssessment` → Objetivo: <500ms p95
- `BenchmarkCreateAttempt` → Objetivo: <2000ms p95
- `BenchmarkGetAttemptResults` → Objetivo: <300ms p95

Archivo: `test/benchmark/assessment_benchmark_test.go`

```bash
go test ./test/benchmark -bench=. -benchmem
```

---

### TASK-05-B-005: Optimización basada en Benchmarks
**Prioridad:** LOW  
**Estimación:** 4h

Si benchmarks muestran problemas de performance:
- Optimizar queries N+1
- Agregar índices faltantes
- Cachear resultados frecuentes

---

## ✅ Criterios de Validación del Sprint 05-B

Al finalizar el Sprint 05-B:

- [ ] Coverage global >80%
- [ ] Coverage services >80%
- [ ] Coverage handlers >70%
- [ ] Coverage repositorios >70%
- [ ] Benchmarks <2s p95
- [ ] Todos los tests pasando

---

## 🚀 Comandos

```bash
# Coverage completo
go test ./... -cover -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total

# Tests de repositorios
RUN_INTEGRATION_TESTS=true go test ./internal/infrastructure/persistence/... -v -tags=integration -cover

# Benchmarks
go test ./test/benchmark -bench=. -benchmem -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

---

## 📊 Coverage Objetivo Final

| Capa | Sprint 05-A | Sprint 05-B | Incremento |
|------|-------------|-------------|------------|
| **Global** | 36.9% → 60% | 60% → **>80%** | +43.1% |
| Dominio | 94.4% | 94.4% | - |
| Services | 61.7% | **>80%** | +18.3% |
| Handlers | 57.1% | **>70%** | +12.9% |
| Repositories | 0% | **>70%** | +70% |

---

**Sprint:** 05-B/06 (Testing Avanzado)  
**Prerrequisito:** Sprint 05-A completado  
**Fecha Creación:** 2025-11-17
