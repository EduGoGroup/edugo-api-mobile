# 📋 Requerimientos para API-Admin

> **Documento para:** Equipo de api-admin  
> **Creado:** Diciembre 2024  
> **Prioridad:** 🔴 Alta  
> **Relacionado con:** TODO-001 en edugo-api-mobile

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

**Archivo probable:** `internal/auth/claims.go` o donde esté definido `CustomClaims`

```go
// ANTES
type CustomClaims struct {
    jwt.RegisteredClaims
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
}

// DESPUÉS
type CustomClaims struct {
    jwt.RegisteredClaims
    UserID   string `json:"user_id"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    SchoolID string `json:"school_id"` // ← AGREGAR
}
```

### 2. Incluir `school_id` al Generar Token

**Archivo probable:** `internal/service/auth_service.go` o donde se maneje el login

```go
func (s *authService) Login(ctx context.Context, email, password string) (*TokenResponse, error) {
    // ... validar credenciales ...

    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil {
        return nil, err
    }

    claims := &CustomClaims{
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   user.ID.String(),
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
        UserID:   user.ID.String(),
        Email:    user.Email,
        Role:     user.Role,
        SchoolID: user.SchoolID.String(), // ← AGREGAR
    }

    token, err := s.jwtManager.GenerateToken(claims)
    // ...
}
```

### 3. Actualizar Endpoint de Validación (si aplica)

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

### En API-Admin

- [ ] Agregar `SchoolID` a struct `CustomClaims`
- [ ] Modificar generación de token en login para incluir `school_id`
- [ ] Actualizar respuesta de `/auth/validate` (si existe)
- [ ] Agregar tests unitarios
- [ ] Documentar en Swagger/OpenAPI

### En edugo-api-mobile (posterior)

- [ ] Actualizar `RemoteAuthMiddleware` para extraer `school_id`
- [ ] Crear helper `MustGetSchoolID(c *gin.Context)`
- [ ] Usar en `material_service.go` en lugar de `uuid.New()`
- [ ] Agregar tests

---

## 📅 Timeline Sugerido

| Fase | Descripción | Estimado |
|------|-------------|----------|
| 1 | Implementar en api-admin | 2-4 horas |
| 2 | Actualizar edugo-api-mobile | 1-2 horas |
| 3 | Testing E2E | 1 hora |

---

## 📞 Contacto

Para dudas sobre este requerimiento, contactar al equipo de edugo-api-mobile.
