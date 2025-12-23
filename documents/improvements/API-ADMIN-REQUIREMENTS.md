# 📋 Requerimientos para API-Admin

> **Documento para:** Equipo de api-admin  
> **Creado:** Diciembre 2024  
> **Actualizado:** 23 Diciembre 2024  
> **Prioridad:** ✅ Completado  
> **Relacionado con:** TODO-001 en edugo-api-mobile

---

## ✅ Estado de Implementación - COMPLETADO

> **Última revisión de código:** 23 Diciembre 2024  
> **PR de implementación:** [PR #64](https://github.com/EduGoGroup/edugo-api-administracion/pull/64) - Merged  
> **Branch:** `feature/add-school-id-to-jwt` → `dev` → `main`

| Componente | Estado | PR/Commit |
|------------|--------|-----------|
| `school_id` en JWT Claims | ✅ **Implementado** | PR #64 |
| `school_id` en Login Response | ✅ **Implementado** | PR #64 |
| `school_id` en User Entity | ✅ **Implementado** | PR #49 (infrastructure) |
| `SwitchContext` endpoint | ✅ **Implementado** | PR #64 |
| Relación User-School en BD | ✅ **Implementado** | Columna `school_id` en users |

### Resumen de Cambios Implementados

**En edugo-infrastructure (postgres/v0.13.0):**
- Agregada columna `school_id` a entidad `User`
- GitHub Release: `postgres/v0.13.0`

**En edugo-api-administracion:**
- `internal/shared/crypto/jwt_manager.go`: Agregado `SchoolID` a Claims
- `internal/auth/dto/auth_dto.go`: Agregados DTOs `SwitchContextRequest`, `SwitchContextResponse`, `ContextInfo`
- `internal/auth/service/auth_service.go`: Implementado método `SwitchContext`
- `internal/auth/handler/auth_handler.go`: Agregado endpoint `POST /v1/auth/switch-context`
- `internal/domain/repository/unit_membership_repository.go`: Agregado método `FindByUserAndSchool`
- `internal/container/container.go`: Inyección de `membershipRepo` a `AuthService`

---

## 🎯 Funcionalidad Implementada

### 1. JWT con school_id

El JWT ahora incluye `school_id` del usuario:

```go
type Claims struct {
    UserID   string `json:"user_id"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    SchoolID string `json:"school_id"` // ✅ Implementado
    jwt.RegisteredClaims
}
```

### 2. Endpoint de Cambio de Contexto

**POST /v1/auth/switch-context**

Permite a un usuario cambiar su contexto activo a otra escuela donde tenga membresía:

```json
// Request
{
    "school_id": "987fcdeb-51a2-3c4d-e5f6-789012345678"
}

// Response
{
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG...",
    "expires_in": 3600,
    "token_type": "Bearer",
    "context": {
        "school_id": "987fcdeb-51a2-3c4d-e5f6-789012345678",
        "school_name": "Escuela Ejemplo",
        "role": "teacher",
        "user_id": "123e4567-e89b-12d3-a456-426614174000",
        "email": "docente@escuela.edu"
    }
}
```

### 3. Arquitectura Multi-tenant

El sistema soporta:
- **1:1** - `users.school_id`: Escuela principal/default del usuario
- **N:N** - `memberships`: Múltiples escuelas con diferentes roles por usuario
- **Cambio de contexto**: El usuario puede cambiar su escuela activa obteniendo un nuevo JWT con el rol correspondiente

---

## 🔄 Siguientes Pasos en edugo-api-mobile

### Ya Desbloqueado

Con la implementación completada en api-admin, ahora se puede:

- [x] ~~Esperar implementación en api-admin~~ ✅
- [ ] Actualizar `RemoteAuthMiddleware` para extraer `school_id`
- [ ] Crear helper `MustGetSchoolID(c *gin.Context)`
- [ ] Usar en `material_service.go` en lugar de `uuid.New()`
- [ ] Agregar tests

### Código a Implementar en api-mobile

```go
// internal/infrastructure/http/middleware/auth.go
func MustGetSchoolID(c *gin.Context) uuid.UUID {
    schoolIDStr, exists := c.Get("school_id")
    if !exists {
        // Fallback: obtener de usuario si no está en JWT
        panic("school_id not found in context")
    }
    schoolID, err := uuid.Parse(schoolIDStr.(string))
    if err != nil {
        panic("invalid school_id format")
    }
    return schoolID
}

// internal/application/service/material_service.go
func (s *MaterialService) CreateMaterial(...) {
    // ANTES (hardcodeado):
    // schoolID := uuid.New()

    // DESPUÉS (del contexto):
    schoolID := middleware.MustGetSchoolID(c)
    // ...
}
```

---

## 📊 Especificación del Campo

| Atributo | Valor |
|----------|-------|
| **Nombre del claim** | `school_id` |
| **Tipo** | `string` (UUID format) |
| **Requerido** | Sí, para usuarios con rol `teacher` y `admin` |
| **Formato** | UUID v4 como string |
| **Ejemplo** | `"987fcdeb-51a2-3c4d-e5f6-789012345678"` |

---

## 📅 Timeline - COMPLETADO

| Fase | Descripción | Estimado | Estado |
|------|-------------|----------|--------|
| 1 | Implementar en api-admin | 2-4 horas | ✅ Completado (23 Dic 2024) |
| 2 | Actualizar edugo-api-mobile | 1-2 horas | ⏳ Pendiente (Fase 1 Sprint) |
| 3 | Testing E2E | 1 hora | ⏳ Pendiente |

---

## 📞 Contacto

Para dudas sobre este requerimiento, contactar al equipo de edugo-api-mobile.

---

## 📝 Historial de Cambios

| Fecha | Cambio | Autor |
|-------|--------|-------|
| Dic 2024 | Documento creado | Claude Code |
| 23 Dic 2024 | Implementación completada en api-admin | Claude Code |
