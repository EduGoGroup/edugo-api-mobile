# 📋 Requerimientos para API-Admin

> **Documento para:** Equipo de api-admin  
> **Creado:** Diciembre 2024  
> **Actualizado:** Diciembre 2024  
> **Prioridad:** 🔴 Alta  
> **Relacionado con:** TODO-001 en edugo-api-mobile

---

## ⚠️ Estado de Implementación

> **Última revisión de código:** 23 Diciembre 2024

| Componente | Estado | Notas |
|------------|--------|-------|
| `school_id` en JWT Claims | ❌ **No implementado** | Falta agregar campo |
| `school_id` en Login Response | ❌ **No implementado** | Falta incluir en respuesta |
| Relación User-School en BD | ⚠️ **Por verificar** | Revisar si existe columna |

### Código Actual en api-admin

**Archivo:** `internal/shared/crypto/jwt_manager.go`
```go
// Estado ACTUAL - Sin school_id
type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}
```

**Archivo:** `internal/auth/dto/auth_dto.go`
```go
// Estado ACTUAL - Sin school_id
type TokenClaims struct {
    UserID    string    `json:"user_id"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    TokenID   string    `json:"jti"`
    IssuedAt  time.Time `json:"iat"`
    ExpiresAt time.Time `json:"exp"`
}

type UserInfo struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    FullName  string `json:"full_name"`
    Role      string `json:"role"`
}
```

---

## 🎯 Objetivo

Incluir `school_id` en el JWT para que los microservicios downstream (como edugo-api-mobile) puedan asociar correctamente los recursos a la escuela del usuario.

---

## 📌 Contexto

### Problema Actual

En `edugo-api-mobile`, cuando un docente crea un material educativo, el sistema no sabe a qué escuela pertenece porque el JWT no incluye esa información:

```go
// internal/application/service/material_service.go:64
schoolID := uuid.New() // ← Se genera UUID aleatorio (INCORRECTO)
```

### Impacto

| Problema | Severidad |
|----------|-----------|
| Materiales no se asocian a la escuela correcta | 🔴 Alta |
| Queries de filtrado por escuela no funcionan | 🔴 Alta |
| Violación de aislamiento multi-tenant | 🔴 Alta |

---

## ✅ Requerimiento: Agregar `school_id` al JWT

### 1. Modificar Claims del JWT

**Archivo:** `internal/shared/crypto/jwt_manager.go`

```go
// ANTES (estado actual)
type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

// DESPUÉS (requerido)
type Claims struct {
    UserID   string `json:"user_id"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    SchoolID string `json:"school_id"` // ← AGREGAR
    jwt.RegisteredClaims
}
```

### 2. Modificar DTOs de Auth

**Archivo:** `internal/auth/dto/auth_dto.go`

```go
// Agregar SchoolID a TokenClaims
type TokenClaims struct {
    UserID    string    `json:"user_id"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    SchoolID  string    `json:"school_id"` // ← AGREGAR
    TokenID   string    `json:"jti"`
    IssuedAt  time.Time `json:"iat"`
    ExpiresAt time.Time `json:"exp"`
}

// Agregar SchoolID a UserInfo
type UserInfo struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    FullName  string `json:"full_name"`
    Role      string `json:"role"`
    SchoolID  string `json:"school_id"` // ← AGREGAR
}
```

### 3. Incluir `school_id` al Generar Token

**Archivo:** `internal/auth/handler/auth_handler.go` o servicio relacionado

```go
func (h *AuthHandler) Login(c *gin.Context) {
    // ... validar credenciales ...

    user, err := h.userRepo.FindByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }

    // Obtener school_id del usuario
    schoolID := user.SchoolID // ← Obtener de la BD

    claims := &crypto.Claims{
        UserID:   user.ID.String(),
        Email:    user.Email,
        Role:     user.Role,
        SchoolID: schoolID.String(), // ← INCLUIR
        RegisteredClaims: jwt.RegisteredClaims{
            // ...
        },
    }

    token, err := h.jwtManager.GenerateToken(claims)
    // ...
}
```

### 4. Actualizar Endpoint de Validación (si aplica)

Si api-admin tiene un endpoint `/auth/validate` que devuelve claims, incluir `school_id`:

```json
{
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "docente@escuela.edu",
    "role": "teacher",
    "school_id": "987fcdeb-51a2-3c4d-e5f6-789012345678"
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

## ❓ Preguntas para el Equipo de API-Admin

1. **¿La tabla `users` tiene columna `school_id`?**
   - Si no existe, ¿de dónde se obtiene la relación usuario-escuela?

2. **¿Un usuario puede pertenecer a múltiples escuelas?**
   - Si es así, ¿cuál usar? ¿La escuela activa/principal?

3. **¿Qué pasa con usuarios `super_admin` que no pertenecen a una escuela específica?**
   - Sugerencia: `school_id` puede ser `null` o string vacío para ellos

4. **¿Hay un endpoint `/auth/validate` que deba actualizarse?**
   - Si es así, incluir `school_id` en la respuesta

---

## 🔄 Pasos de Implementación

### En API-Admin (PENDIENTE)

- [ ] Verificar que tabla `users` tenga columna `school_id`
- [ ] Agregar `SchoolID` a struct `Claims` en `jwt_manager.go`
- [ ] Agregar `SchoolID` a struct `TokenClaims` en `auth_dto.go`
- [ ] Agregar `SchoolID` a struct `UserInfo` en `auth_dto.go`
- [ ] Modificar generación de token en login para incluir `school_id`
- [ ] Actualizar respuesta de `/auth/validate` (si existe)
- [ ] Agregar tests unitarios
- [ ] Documentar en Swagger/OpenAPI

### En edugo-api-mobile (posterior, depende de API-Admin)

- [ ] Actualizar `RemoteAuthMiddleware` para extraer `school_id`
- [ ] Crear helper `MustGetSchoolID(c *gin.Context)`
- [ ] Usar en `material_service.go` en lugar de `uuid.New()`
- [ ] Agregar tests

---

## 📅 Timeline Sugerido

| Fase | Descripción | Estimado | Estado |
|------|-------------|----------|--------|
| 1 | Implementar en api-admin | 2-4 horas | ❌ Pendiente |
| 2 | Actualizar edugo-api-mobile | 1-2 horas | ⏳ Bloqueado |
| 3 | Testing E2E | 1 hora | ⏳ Bloqueado |

---

## 📞 Contacto

Para dudas sobre este requerimiento, contactar al equipo de edugo-api-mobile.
