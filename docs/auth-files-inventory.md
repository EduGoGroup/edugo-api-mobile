# Inventario de Archivos de Auth a Eliminar
## Sprint 3 - Limpieza de Código Duplicado

**Fecha**: 2025-11-24
**Proyecto**: edugo-api-mobile

---

## Archivos a ELIMINAR (código duplicado)

| Archivo | Líneas | Descripción |
|---------|--------|-------------|
| `internal/application/service/auth_service.go` | 367 | Login, Refresh, Logout, RevokeAll |
| `internal/application/dto/auth_dto.go` | 51 | DTOs de auth |
| `internal/infrastructure/http/handler/auth_handler.go` | 213 | Handlers HTTP de auth |

**Subtotal código**: 631 líneas

## Tests relacionados a ELIMINAR

| Archivo | Líneas | Descripción |
|---------|--------|-------------|
| `internal/application/service/auth_service_test.go` | 492 | Tests de auth service |
| `internal/infrastructure/http/handler/auth_handler_test.go` | 576 | Tests de auth handler |
| `test/integration/auth_flow_test.go` | 158 | Tests de integración auth |

**Subtotal tests**: 1,226 líneas

## Archivos a MODIFICAR (no eliminar)

| Archivo | Líneas | Acción |
|---------|--------|--------|
| `internal/infrastructure/http/middleware/auth.go` | 54 | MODIFICAR: Cambiar validación local por remota |
| `internal/infrastructure/http/middleware/auth_test.go` | 259 | MODIFICAR: Actualizar tests |
| `internal/infrastructure/http/router/router.go` | - | MODIFICAR: Eliminar rutas de auth |
| `internal/container/container.go` | - | MODIFICAR: Eliminar dependencias auth |

## Rutas a Eliminar del Router

```
POST /v1/auth/login      -> Mover a api-admin (ya existe)
POST /v1/auth/refresh    -> Mover a api-admin (ya existe)
POST /v1/auth/logout     -> Mover a api-admin (ya existe)
POST /v1/auth/revoke-all -> Mover a api-admin (implementar si necesario)
GET  /v1/auth/me         -> Mover a api-admin (ya existe)
```

## Dependencias Externas

- `github.com/EduGoGroup/edugo-shared/auth` - JWTManager (seguir usando para validación)
- `repository.UserReader` - Eliminar dependencia
- `repository.RefreshTokenRepository` - Eliminar dependencia  
- `repository.LoginAttemptRepository` - Eliminar dependencia

## Resumen

| Categoría | Líneas |
|-----------|--------|
| Código a eliminar | 631 |
| Tests a eliminar | 1,226 |
| **Total** | **1,857 líneas** |

## Estrategia

1. El middleware `auth.go` se **modifica** para validar tokens remotamente con api-admin
2. Los endpoints de auth (`/login`, `/refresh`, `/logout`) se **eliminan** de api-mobile
3. Los clientes (apple-app) ya apuntan a api-admin para auth (Sprint 2 completado)
4. api-mobile solo necesita **validar** tokens, no emitirlos
