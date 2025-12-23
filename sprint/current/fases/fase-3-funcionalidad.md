# Fase 3: TODOs de Funcionalidad

> **Estado:** ✅ COMPLETADA
> **PR:** Pendiente crear
> **Branch:** `feature/functionality-todos`
> **Estimado:** 6-8 horas
> **Real:** ~4 horas

---

## Tareas

### ✅ TODO-004: URL real de S3 en Material Service
- **Archivo:** `internal/application/service/material_service.go`
- **Problema:** Evento MaterialUploaded usaba URL placeholder
- **Solución:** Mover publicación a NotifyUploadComplete con datos reales
- **Commit:** `ae1d9b1`

**Cambios:**
- Eliminar publicación de evento en CreateMaterial
- Agregar publicación en NotifyUploadComplete con FileURL, FileSizeBytes, FileType reales

---

### ✅ TODO-006: Implementar FindByIDWithVersions completo
- **Archivo:** `internal/infrastructure/persistence/postgres/repository/material_repository_impl.go:367`
- **Problema:** Método solo retorna material sin versiones
- **Solución:** Implementar query real a `material_versions`
- **Commit:** `457cf27`

**Cambios realizados:**
- Query separado para obtener versiones ordenadas por `version_number DESC`
- Mapeo completo de todos los campos de `MaterialVersion`
- Retorna material con slice de versiones

---

### ✅ TODO-007: Publicar evento material_completed
- **Archivo:** `internal/application/service/progress_service.go`
- **Problema:** No se notifica cuando un estudiante completa un material
- **Solución:** Publicar evento cuando progress = 100%
- **Commit:** `0ebef2c`

**Cambios realizados:**
- Agregar `MaterialCompletedPayload` y `NewMaterialCompletedEvent` en events.go
- Inyectar `rabbitmq.Publisher` en `ProgressService`
- Publicar evento `material.completed` cuando `percentage == 100`
- Fire-and-forget (no afecta flujo principal si falla publicación)
- Tests actualizados con mock del publisher

---

### ✅ TODO-005: Preparar restauración de eventos de Assessment
- **Tipo:** Documentación
- **Problema:** Eventos de assessment comentados, necesitan schema en infrastructure
- **Solución:** Documentar requerimientos para edugo-infrastructure
- **Commit:** Pendiente

**Entregable:**
- ✅ Documento `documents/improvements/ASSESSMENT-EVENTS-REQUIREMENTS.md`
- Eventos propuestos: `assessment.attempt.completed`, `assessment.first_passed`, `assessment.attempt.started`
- Plan de implementación en 3 fases
- Cambios requeridos en servicio y container documentados
