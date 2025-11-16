# Sprint 02: Capa de Dominio
# Sistema de Evaluaciones - EduGo

**Duración:** 3 días  
**Objetivo:** Implementar capa de dominio (entities, value objects, interfaces) siguiendo Clean Architecture.

---

## 🎯 Objetivo del Sprint

Crear la capa de dominio del Sistema de Evaluaciones, incluyendo:
- 3 entities principales (Assessment, Attempt, Answer)
- 5+ value objects
- 3 repository interfaces
- Business rules del dominio
- Tests unitarios con >90% coverage

---

## 📋 Tareas del Sprint

Ver archivo [TASKS.md](./TASKS.md) para lista detallada.

**Resumen:**
- 3 entities con validaciones
- 5 value objects (Score, AssessmentID, etc.)
- 3 repository interfaces
- Tests unitarios completos

---

## 🔗 Dependencias

Ver archivo [DEPENDENCIES.md](./DEPENDENCIES.md).

**Críticas:**
- Sprint 01 completado (schema PostgreSQL creado)
- Go 1.21+ instalado
- Clean Architecture structure creada

---

## ❓ Decisiones y Preguntas

Ver archivo [QUESTIONS.md](./QUESTIONS.md).

**Decisiones clave:**
- Usar Value Objects para Score, AssessmentID
- Business rules en entities (no en services)
- Repository interfaces en dominio (implementación en infrastructure)

---

## ✅ Validación

Ver archivo [VALIDATION.md](./VALIDATION.md) para checklist completo.

**Criterios de éxito:**
- [ ] 3 entities creadas con business logic
- [ ] 5+ value objects con validaciones
- [ ] 3 repository interfaces definidas
- [ ] Tests unitarios >90% coverage
- [ ] Sin dependencias externas en dominio

---

## 📊 Entregables

1. `internal/domain/entities/assessment.go`
2. `internal/domain/entities/attempt.go`
3. `internal/domain/entities/answer.go`
4. `internal/domain/valueobjects/score.go`
5. `internal/domain/repositories/assessment_repository.go`
6. Tests unitarios en `internal/domain/entities/*_test.go`

---

## 🚀 Comandos Rápidos

```bash
# Ejecutar tests de dominio
go test ./internal/domain/... -v

# Ver coverage
go test ./internal/domain/... -cover

# Generar reporte HTML de coverage
go test ./internal/domain/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Verificar que no hay dependencias externas en dominio
go list -f '{{.Imports}}' ./internal/domain/... | grep -v "internal/domain"
```

---

**Generado con:** Claude Code  
**Sprint:** 02/06  
**Última actualización:** 2025-11-14
