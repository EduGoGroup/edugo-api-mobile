# Plan Integral: Sistema de Autenticación OAuth2 para Ecosistema EduGo

**Fecha de creación**: 2024-10-31
**Estado**: 📋 EN ANÁLISIS
**Prioridad**: 🔴 ALTA (Seguridad crítica)

---

## 📊 Análisis del Ecosistema Actual

### Proyectos en el Ecosistema

```
EduGo Ecosystem
├── edugo-api-mobile        # API de uso frecuente (endpoints día a día)
│   └── Usuarios: Estudiantes, Docentes, Tutores
├── edugo-api-admin         # API administrativa (endpoints menos frecuentes)
│   └── Usuarios: Administradores, Super admins
├── edugo-worker            # Worker de eventos (escucha RabbitMQ)
│   └── Función: Procesamiento asíncrono
└── edugo-shared            # Librería compartida (modularizada)
    └── Uso: Todos los proyectos dependen de ella
```

---

## 🚨 Problemas Críticos Identificados en Autenticación Actual

### 1. **Hash de Password Inseguro**

**Ubicación**: `edugo-api-mobile/internal/application/service/auth_service.go:116`

```go
// ❌ PROBLEMA: SHA256 NO es seguro para passwords
func hashPassword(password string) string {
    h := sha256.New()
    h.Write([]byte(password))
    return hex.EncodeToString(h.Sum(nil))
}
```

**Por qué es inseguro**:
- SHA256 es extremadamente rápido → vulnerable a ataques de fuerza bruta
- No tiene salt → mismas contraseñas generan mismo hash
- Rainbow tables pueden romper passwords comunes en segundos
- No tiene costo computacional ajustable

**Impacto**: 🔴 CRÍTICO - Todas las contraseñas están en riesgo

---

### 2. **Refresh Token No Implementado Correctamente**

**Código actual** (`auth_service.go:82-92`):
```go
refreshToken, err := s.jwtManager.GenerateToken(
    user.ID().String(),
    user.Email().String(),
    user.Role(),
    7*24*time.Hour, // ❌ Solo es otro JWT con más duración
)
```

**Problemas**:
- ❌ No hay almacenamiento persistente de refresh tokens
- ❌ No hay endpoint para refrescar tokens
- ❌ No hay revocación de tokens
- ❌ No hay rotación de refresh tokens
- ❌ Refresh token es idéntico a access token (solo cambia duración)

**Impacto**: 🔴 CRÍTICO - No hay forma de revocar acceso sin cambiar JWT secret

---

### 3. **No Hay Revocación de Tokens**

**Problemas**:
- Si un token es comprometido, es válido hasta que expire
- No hay lista negra de tokens revocados
- No hay logout real (token sigue siendo válido)
- No hay forma de invalidar sesiones de usuario

**Impacto**: 🟡 ALTO - Sesiones comprometidas no se pueden revocar

---

### 4. **Middleware Duplicado en Cada Proyecto**

**Ubicaciones**:
- `edugo-api-mobile/cmd/main.go:201-236` (35 líneas)
- Probablemente duplicado en `edugo-api-admin`

**Problemas**:
- Código duplicado
- Cambio de seguridad requiere actualizar N proyectos
- Inconsistencias entre proyectos

**Impacto**: 🟡 MODERADO - Mantenibilidad y consistencia

---

### 5. **No Hay Rate Limiting por Usuario**

**Problema**: Cualquiera puede intentar login infinitas veces

**Impacto**: 🟡 ALTO - Vulnerable a ataques de fuerza bruta

---

### 6. **No Hay Auditoría de Autenticación**

**Problemas**:
- No se registran intentos fallidos
- No hay log de sesiones activas
- No hay detección de intentos sospechosos

**Impacto**: 🟡 MODERADO - No hay trazabilidad

---

## 🏗️ Arquitecturas Posibles

### **Opción 1: Servidor de Autenticación Centralizado (Recomendado)**

```
                        ┌─────────────────────────┐
                        │   Auth Service          │
                        │  (edugo-auth-service)   │
                        │                         │
                        │  - Login                │
                        │  - Refresh Token        │
                        │  - Revocación           │
                        │  - Gestión de usuarios  │
                        └───────────┬─────────────┘
                                    │
                    ┌───────────────┼───────────────┐
                    │               │               │
                    ▼               ▼               ▼
         ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
         │ api-mobile   │  │ api-admin    │  │ worker       │
         │              │  │              │  │              │
         │ Solo valida  │  │ Solo valida  │  │ Solo valida  │
         │ tokens       │  │ tokens       │  │ tokens       │
         └──────────────┘  └──────────────┘  └──────────────┘
```

**Ventajas**:
- ✅ Single source of truth para autenticación
- ✅ Gestión centralizada de usuarios y permisos
- ✅ Auditoría centralizada
- ✅ Revocación centralizada de tokens
- ✅ Más fácil implementar features avanzadas (MFA, OAuth2, SAML)
- ✅ Las APIs solo validan tokens (más livianas)

**Desventajas**:
- ❌ Punto único de falla (requiere alta disponibilidad)
- ❌ Proyecto adicional a mantener
- ❌ Latencia adicional en autenticación

**Cuándo usar**:
- Ecosistemas con 3+ servicios ✅ (tenemos 3)
- Necesitas SSO (Single Sign-On) ✅
- Planeas agregar más servicios en el futuro ✅

---

### **Opción 2: Autenticación Descentralizada con Librería Compartida**

```
         ┌──────────────────┐  ┌──────────────────┐
         │   api-mobile     │  │   api-admin      │
         │                  │  │                  │
         │ ┌──────────────┐ │  │ ┌──────────────┐ │
         │ │ Auth Module  │ │  │ │ Auth Module  │ │
         │ │ (from shared)│ │  │ │ (from shared)│ │
         │ └──────────────┘ │  │ └──────────────┘ │
         └──────────────────┘  └──────────────────┘
                   │                      │
                   └──────────┬───────────┘
                              ▼
                    ┌─────────────────┐
                    │  Shared DB      │
                    │  (users, tokens)│
                    └─────────────────┘
```

**Ventajas**:
- ✅ No necesita servicio adicional
- ✅ Menos latencia (sin hop adicional)
- ✅ Menos infraestructura

**Desventajas**:
- ❌ Código de autenticación duplicado en cada API
- ❌ Necesitan acceso a BD de usuarios (acoplamiento)
- ❌ Difícil hacer cambios (actualizar todas las APIs)
- ❌ No hay SSO real

**Cuándo usar**:
- Ecosistema muy pequeño (1-2 servicios)
- No planeas crecer mucho

---

### **Opción 3: Híbrida (Inicio Recomendado)**

```
Phase 1 (Ahora):
- Implementar autenticación en api-mobile
- Extraer a shared lo máximo posible
- api-admin consume mismo código de shared

Phase 2 (Después):
- Migrar a Auth Service cuando el ecosistema crezca
- Mover lógica de shared a servicio dedicado
```

**Ventajas**:
- ✅ Rápido de implementar
- ✅ Mejora inmediata de seguridad
- ✅ Preparado para migrar a centralizado

---

## 🎯 Decisión Arquitectónica Recomendada

### **Recomendación: Opción 3 (Híbrida) → Migrar a Opción 1**

**Fase Inmediata (Sprint actual)**:
1. Implementar OAuth2 en `edugo-api-mobile`
2. Extraer componentes a `edugo-shared/auth`
3. `edugo-api-admin` consume mismo código

**Fase Futura (Q1 2026)**:
1. Crear `edugo-auth-service` cuando crezca el ecosistema
2. Migrar lógica de autenticación al servicio dedicado

---

## 📦 Componentes OAuth2 Necesarios

### **1. Componentes en `edugo-shared/auth`**

```go
edugo-shared/auth/
├── jwt_manager.go              // ✅ YA EXISTE
├── password.go                 // 🆕 NUEVO - bcrypt hash
├── token_store.go              // 🆕 NUEVO - almacenamiento de tokens
├── refresh_token.go            // 🆕 NUEVO - gestión de refresh tokens
├── claims.go                   // 🆕 NUEVO - estructura de claims extendida
└── oauth2/
    ├── authorization_code.go   // 🆕 OPCIONAL - flujo OAuth2
    ├── client_credentials.go   // 🆕 OPCIONAL - para M2M
    └── pkce.go                 // 🆕 OPCIONAL - PKCE para móviles
```

### **2. Componentes en `edugo-shared/middleware/gin`**

```go
edugo-shared/middleware/gin/
├── jwt_auth.go                 // 🆕 NUEVO - middleware JWT
├── context_helpers.go          // 🆕 NUEVO - extraer claims
├── rate_limiter.go             // 🆕 NUEVO - rate limiting
└── cors.go                     // 🆕 NUEVO - CORS reutilizable
```

### **3. Componentes en `edugo-api-mobile`**

```go
internal/
├── application/
│   ├── service/
│   │   ├── auth_service.go           // ✏️ MODIFICAR - usar nuevos componentes
│   │   └── token_service.go          // 🆕 NUEVO - gestión de tokens
│   └── dto/
│       ├── auth_dto.go               // ✏️ MODIFICAR - agregar refresh
│       └── token_dto.go              // 🆕 NUEVO - DTOs de tokens
├── domain/
│   └── repository/
│       ├── user_repository.go        // ✅ YA EXISTE
│       └── refresh_token_repository.go // 🆕 NUEVO - tokens persistidos
└── infrastructure/
    ├── http/
    │   └── handler/
    │       ├── auth_handler.go       // ✏️ MODIFICAR - endpoints nuevos
    │       └── token_handler.go      // 🆕 NUEVO - refresh endpoint
    └── persistence/
        └── postgres/
            └── repository/
                └── refresh_token_repository_impl.go // 🆕 NUEVO
```

### **4. Base de Datos**

```sql
-- 🆕 Nueva tabla: refresh_tokens
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash VARCHAR(64) NOT NULL UNIQUE,  -- SHA256 del token
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    client_info JSONB,                       -- Info del cliente (IP, User-Agent)
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    revoked_at TIMESTAMP,
    replaced_by UUID REFERENCES refresh_tokens(id),  -- Token rotation
    INDEX idx_user_id (user_id),
    INDEX idx_token_hash (token_hash),
    INDEX idx_expires_at (expires_at)
);

-- 🆕 Nueva tabla: login_attempts (rate limiting)
CREATE TABLE login_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    identifier VARCHAR(255) NOT NULL,        -- Email o IP
    attempt_type VARCHAR(20) NOT NULL,       -- 'email' o 'ip'
    attempted_at TIMESTAMP DEFAULT NOW(),
    successful BOOLEAN,
    INDEX idx_identifier_time (identifier, attempted_at)
);

-- 🆕 Nueva tabla: user_sessions (auditoría)
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    refresh_token_id UUID REFERENCES refresh_tokens(id),
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    last_activity TIMESTAMP DEFAULT NOW(),
    ended_at TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at)
);
```

---

## 🔐 Flujos OAuth2 a Implementar

### **Flujo 1: Login (Password Grant - Simplificado)**

```
┌─────────┐                                          ┌──────────────┐
│ Cliente │                                          │  api-mobile  │
│ (App)   │                                          │              │
└────┬────┘                                          └──────┬───────┘
     │                                                       │
     │  POST /v1/auth/login                                 │
     │  { email, password }                                 │
     ├──────────────────────────────────────────────────────>│
     │                                                       │
     │                                   1. Validar password│
     │                                   2. Generar access  │
     │                                   3. Generar refresh │
     │                                   4. Guardar refresh │
     │                                                       │
     │  200 OK                                               │
     │  { access_token, refresh_token, expires_in }         │
     │<──────────────────────────────────────────────────────┤
     │                                                       │
```

### **Flujo 2: Refresh Token**

```
┌─────────┐                                          ┌──────────────┐
│ Cliente │                                          │  api-mobile  │
└────┬────┘                                          └──────┬───────┘
     │                                                       │
     │  POST /v1/auth/refresh                               │
     │  { refresh_token }                                   │
     ├──────────────────────────────────────────────────────>│
     │                                                       │
     │                                   1. Validar refresh │
     │                                   2. Verificar no    │
     │                                      revocado        │
     │                                   3. Generar nuevo   │
     │                                      access          │
     │                                   4. Rotar refresh   │
     │                                      (opcional)      │
     │                                                       │
     │  200 OK                                               │
     │  { access_token, refresh_token, expires_in }         │
     │<──────────────────────────────────────────────────────┤
     │                                                       │
```

### **Flujo 3: Logout**

```
┌─────────┐                                          ┌──────────────┐
│ Cliente │                                          │  api-mobile  │
└────┬────┘                                          └──────┬───────┘
     │                                                       │
     │  POST /v1/auth/logout                                │
     │  Authorization: Bearer {access_token}                │
     ├──────────────────────────────────────────────────────>│
     │                                                       │
     │                                   1. Extraer user_id │
     │                                   2. Revocar refresh │
     │                                   3. Cerrar sesión   │
     │                                                       │
     │  204 No Content                                       │
     │<──────────────────────────────────────────────────────┤
     │                                                       │
```

### **Flujo 4: Revocación de Todas las Sesiones**

```
POST /v1/auth/revoke-all
Authorization: Bearer {access_token}

Uso: Cuando usuario cambia password o detecta actividad sospechosa
```

---

## 📝 Plan de Implementación Detallado

### **FASE 1: Mejorar Seguridad Básica (1-2 días)**

#### **Tarea 1.1: Implementar bcrypt en edugo-shared**

```bash
# Archivos a crear/modificar en edugo-shared:
/Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared/auth/password.go
```

**Código**:
```go
package auth

import "golang.org/x/crypto/bcrypt"

const bcryptCost = 12

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
    return string(bytes), err
}

func VerifyPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword(
        []byte(hashedPassword),
        []byte(password),
    )
}
```

**Checklist**:
- [ ] Crear archivo `password.go` en edugo-shared/auth
- [ ] Implementar `HashPassword()` con bcrypt cost 12
- [ ] Implementar `VerifyPassword()`
- [ ] Agregar tests unitarios
- [ ] Commit en shared
- [ ] Crear tag (ej: v0.1.0 → v0.2.0)
- [ ] Push tag a GitHub

---

#### **Tarea 1.2: Migrar api-mobile a bcrypt**

**Archivos a modificar**:
- `internal/application/service/auth_service.go`

**Cambios**:
```go
// ANTES
passwordHash := hashPassword(req.Password)
if user.PasswordHash() != passwordHash {
    return nil, errors.NewUnauthorizedError("invalid credentials")
}

// DESPUÉS
err := auth.VerifyPassword(user.PasswordHash(), req.Password)
if err != nil {
    return nil, errors.NewUnauthorizedError("invalid credentials")
}
```

**Checklist**:
- [ ] Actualizar edugo-shared en go.mod
- [ ] Eliminar función `hashPassword()` local
- [ ] Usar `auth.VerifyPassword()` en login
- [ ] Compilar y verificar
- [ ] Commit

---

### **FASE 2: Implementar Refresh Tokens (2-3 días)**

#### **Tarea 2.1: Crear tabla refresh_tokens**

**Archivo a crear**:
```sql
scripts/postgresql/03_refresh_tokens.sql
```

**Checklist**:
- [ ] Crear script SQL
- [ ] Ejecutar en entorno local
- [ ] Verificar índices

---

#### **Tarea 2.2: Crear RefreshToken en shared**

**Archivo a crear en shared**:
```go
edugo-shared/auth/refresh_token.go
```

**Implementación**:
```go
package auth

import (
    "crypto/rand"
    "encoding/base64"
    "time"
)

type RefreshToken struct {
    Token     string
    ExpiresAt time.Time
}

func GenerateRefreshToken(ttl time.Duration) (*RefreshToken, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return nil, err
    }

    return &RefreshToken{
        Token:     base64.URLEncoding.EncodeToString(bytes),
        ExpiresAt: time.Now().Add(ttl),
    }, nil
}
```

**Checklist**:
- [ ] Implementar generación de refresh tokens
- [ ] Tests unitarios
- [ ] Commit + tag en shared

---

#### **Tarea 2.3: Crear repositorio de refresh tokens**

**Archivos a crear en api-mobile**:
- `internal/domain/repository/refresh_token_repository.go` (interfaz)
- `internal/infrastructure/persistence/postgres/repository/refresh_token_repository_impl.go`

**Checklist**:
- [ ] Definir interfaz
- [ ] Implementar repositorio PostgreSQL
- [ ] Agregar al Container DI
- [ ] Tests de integración

---

#### **Tarea 2.4: Modificar AuthService para usar refresh tokens**

**Cambios**:
1. Login: Generar y guardar refresh token
2. Crear método `RefreshAccessToken()`
3. Crear método `RevokeRefreshToken()`

**Checklist**:
- [ ] Modificar `Login()` para retornar refresh token
- [ ] Implementar `RefreshAccessToken()`
- [ ] Implementar `RevokeRefreshToken()`
- [ ] Tests unitarios

---

#### **Tarea 2.5: Crear endpoints de refresh y logout**

**Archivos a modificar/crear**:
- `internal/infrastructure/http/handler/auth_handler.go`
- `cmd/main.go` (agregar rutas)

**Nuevas rutas**:
```go
POST /v1/auth/refresh       // Refrescar token
POST /v1/auth/logout        // Cerrar sesión (requiere auth)
POST /v1/auth/revoke-all    // Revocar todas las sesiones
```

**Checklist**:
- [ ] Implementar handler de refresh
- [ ] Implementar handler de logout
- [ ] Implementar handler de revoke-all
- [ ] Actualizar Swagger
- [ ] Tests de integración

---

### **FASE 3: Middleware Reutilizable (1 día)**

#### **Tarea 3.1: Crear middleware en shared**

**Archivo a crear en shared**:
```go
edugo-shared/middleware/gin/jwt_auth.go
edugo-shared/middleware/gin/context.go
```

**Checklist**:
- [ ] Implementar `JWTAuthMiddleware()`
- [ ] Implementar helpers de contexto
- [ ] Tests
- [ ] Commit + tag en shared

---

#### **Tarea 3.2: Migrar api-mobile al middleware compartido**

**Cambios en api-mobile**:
- Eliminar `jwtAuthMiddleware()` de `main.go`
- Usar `gin.JWTAuthMiddleware()` de shared

**Checklist**:
- [ ] Actualizar shared
- [ ] Reemplazar middleware local
- [ ] Verificar todos los handlers
- [ ] Commit

---

### **FASE 4: Rate Limiting y Auditoría (1-2 días)**

#### **Tarea 4.1: Implementar rate limiting en login**

**Estrategia**: Redis o PostgreSQL para almacenar intentos

**Checklist**:
- [ ] Crear tabla login_attempts
- [ ] Implementar rate limiter
- [ ] Aplicar en endpoint de login
- [ ] Tests

---

#### **Tarea 4.2: Implementar auditoría de sesiones**

**Checklist**:
- [ ] Crear tabla user_sessions
- [ ] Registrar sesiones en login
- [ ] Actualizar en refresh
- [ ] Cerrar en logout
- [ ] Endpoint para ver sesiones activas

---

### **FASE 5: Aplicar a api-admin (1 día)**

**Checklist**:
- [ ] Actualizar edugo-shared en api-admin
- [ ] Usar mismos servicios de autenticación
- [ ] Usar mismo middleware
- [ ] Verificar que compila
- [ ] Tests de integración

---

## 🔍 Decisiones Pendientes a Tomar

### **Decisión 1: ¿Crear servicio de autenticación dedicado?**

| Opción | Cuándo | Esfuerzo |
|--------|--------|----------|
| **Ahora** | Si planeas agregar muchos servicios pronto | 🔴 ALTO |
| **Después** | Implementar en las APIs, migrar cuando crezca | 🟢 BAJO |

**Recomendación**: Implementar en APIs primero, migrar a servicio en 6 meses.

---

### **Decisión 2: ¿Token rotation en refresh?**

**Token Rotation**: Cada vez que se refresca, se genera nuevo refresh token y se invalida el anterior.

| Opción | Seguridad | Complejidad |
|--------|-----------|-------------|
| **Con rotation** | 🟢 ALTA | 🟡 MEDIA |
| **Sin rotation** | 🟡 MEDIA | 🟢 BAJA |

**Recomendación**: Implementar rotation (mejor seguridad, complejidad manejable).

---

### **Decisión 3: ¿Implementar OAuth2 completo o simplificado?**

**OAuth2 Completo**: Authorization Code, Client Credentials, PKCE, etc.

| Opción | Cuándo usar | Esfuerzo |
|--------|-------------|----------|
| **Completo** | Si tendrás clientes de terceros (public API) | 🔴 ALTO |
| **Simplificado** | Solo para apps internas de EduGo | 🟢 BAJO |

**Recomendación**: Simplificado por ahora (Password Grant + Refresh Token).

---

### **Decisión 4: ¿Redis para tokens o PostgreSQL?**

| Opción | Pros | Contras |
|--------|------|---------|
| **Redis** | Muy rápido, TTL automático | Infraestructura adicional |
| **PostgreSQL** | Ya lo tienes, más simple | Más lento, requiere limpieza manual |

**Recomendación**: PostgreSQL para MVP, migrar a Redis si hay problemas de performance.

---

## 📊 Resumen de Esfuerzo Estimado

| Fase | Tareas | Esfuerzo | Commits |
|------|--------|----------|---------|
| **Fase 1** | bcrypt | 1-2 días | 2 (shared + api-mobile) |
| **Fase 2** | Refresh tokens | 2-3 días | 4-5 |
| **Fase 3** | Middleware compartido | 1 día | 2 (shared + api-mobile) |
| **Fase 4** | Rate limiting + auditoría | 1-2 días | 2-3 |
| **Fase 5** | Aplicar a api-admin | 1 día | 1 |
| **TOTAL** | | **6-9 días** | **11-14 commits** |

---

## 🎯 Entregables

### **En edugo-shared**
- [ ] `auth/password.go` - Hash con bcrypt
- [ ] `auth/refresh_token.go` - Generación de refresh tokens
- [ ] `middleware/gin/jwt_auth.go` - Middleware JWT
- [ ] `middleware/gin/context.go` - Helpers de contexto
- [ ] `middleware/gin/rate_limiter.go` - Rate limiting
- [ ] Tests completos con >80% coverage

### **En edugo-api-mobile**
- [ ] Repositorio de refresh tokens
- [ ] AuthService actualizado
- [ ] Endpoints: /refresh, /logout, /revoke-all
- [ ] Migraciones SQL
- [ ] Tests de integración
- [ ] Documentación Swagger actualizada

### **En edugo-api-admin**
- [ ] Mismo código que api-mobile
- [ ] Tests de integración

---

## 🚀 Próximos Pasos Inmediatos

**¿Quieres proceder?**

1. **Opción A - Implementar ahora**: Comenzar con Fase 1 (bcrypt)
2. **Opción B - Analizar más**: Discutir decisiones pendientes
3. **Opción C - Otro enfoque**: Proponer arquitectura diferente

**Mi recomendación**: Empezar con Fase 1 (bcrypt) hoy mismo, es:
- Rápido (1-2 horas)
- Mejora crítica de seguridad
- No rompe nada existente
- Prepara terreno para Fase 2

---

**Última actualización**: 2025-10-31
**Próxima revisión**: Después de implementar Fase 1
