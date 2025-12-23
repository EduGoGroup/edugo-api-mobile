# Fase 2: TODOs de Autorización

> **Estado:** ✅ COMPLETADA
> **PR:** #90 (merged)
> **Branch:** `feature/auth-todos` (eliminado)
> **Duración real:** ~45 minutos

---

## Tareas

### ✅ TODO-003: Bypass admin en Progress Handler
- **Archivo:** `internal/infrastructure/http/handler/progress_handler.go`
- **Problema:** Admin no podía actualizar progreso de otros usuarios
- **Solución:** Verificar rol admin/super_admin antes de rechazar
- **Commit:** `14f949b`

**Cambios:**
- Agregado `IsAdminRole()` helper en middleware/auth.go
- Agregado `HasRole()` helper genérico
- Progress handler permite admin actualizar progreso de cualquier usuario

---

### ✅ Crear middleware genérico de autorización por rol
- **Archivo:** `internal/infrastructure/http/middleware/remote_auth.go`
- **Commit:** `1de2caf`

**Shortcuts agregados:**
- `RequireAdmin()` - admin o super_admin
- `RequireSuperAdmin()` - solo super_admin
- `RequireTeacher()` - teacher, admin o super_admin
- `RequireStudentOrAbove()` - cualquier rol autenticado

---

### ✅ Aplicar middleware a endpoints sensibles
- **Archivo:** `internal/infrastructure/http/router/router.go`
- **Commit:** `072c7fe`

**Endpoints protegidos:**
| Endpoint | Middleware |
|----------|------------|
| `POST /materials` | RequireTeacher() |
| `POST /materials/:id/upload-complete` | RequireTeacher() |
| `POST /materials/:id/upload-url` | RequireTeacher() |
| `GET /stats/global` | RequireAdmin() |
