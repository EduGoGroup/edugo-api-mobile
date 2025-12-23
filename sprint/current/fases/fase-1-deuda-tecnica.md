# Fase 1: Deuda Técnica Crítica

> **Estado:** ✅ COMPLETADA
> **PR:** #89 (merged)
> **Branch:** `feature/sprint-improvements-plan` (eliminado)
> **Duración real:** ~3 horas

---

## Tareas

### ✅ DEBT-003: Resolver SchoolID hardcodeado
- **Archivo:** `internal/application/service/material_service.go`
- **Problema:** SchoolID hardcodeado como `uuid.Nil`
- **Solución:** Obtener `school_id` del contexto JWT
- **Commit:** `9e83da9`

**Cambios:**
- Agregado `GetSchoolIDFromContext()` en middleware
- Agregado `MustGetSchoolIDFromContext()` en middleware
- MaterialService usa schoolID del contexto

---

### ✅ DEBT-005: Resolver tests unitarios con TODOs
- **Archivos eliminados:**
  - `answer_repository_test.go`
  - `assessment_repository_test.go`
  - `assessment_document_repository_test.go`
- **Razón:** Tests de integración existentes son suficientes
- **Commit:** `3cd38eb`

---

### ✅ DEBT-006: Estandarizar uso de logger
- **Problema:** Uso inconsistente de `zap.Field` vs key-value pairs
- **Solución:** Convertir a formato key-value (compatible con logger.Logger)
- **Commit:** `f920b25`

**Archivos modificados (8):**
- `assessment_attempt_service.go`
- `material_service.go`
- `progress_service.go`
- `stats_service.go`
- `noop/publisher.go`
- `noop/storage.go`
- `rabbitmq/publisher.go`
- `s3/client.go`

---

## Dependencia Externa

- **edugo-infrastructure/postgres v0.13.0** - Columna `school_id` en User
