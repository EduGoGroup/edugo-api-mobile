# 🎯 PLAN MAESTRO - Vista Rápida con Checkboxes

**📅 Última actualización**: 2025-11-05
**🌿 Branch**: `fix/debug-sprint-commands`
**📊 Progreso**: 10/11 commits (91%) | ~15 horas invertidas

---

## 🚀 PRÓXIMA TAREA

**👉 FASE 2.3: Completar Queries Complejas** (1 día restante)

Ver detalles completos en sección [FASE 2](#fase-2-completar-todos-de-servicios) abajo ⬇️

---

## 📋 ÍNDICE DE FASES

- [✅ FASE 0: Autenticación OAuth2](#fase-0-autenticación-oauth2) - **COMPLETADA**
- [✅ FASE 1: Container DI](#fase-1-container-di) - **COMPLETADA**
- [🔵 FASE 2: TODOs de Servicios](#fase-2-completar-todos-de-servicios) - **EN PROGRESO** (67% - 2/3 pasos)
- [⏳ FASE 3: Limpieza](#fase-3-limpieza-y-consolidación) - **SIGUIENTE**
- [⏳ FASE 4: Testing](#fase-4-testing-de-integración) - PENDIENTE

---

## ✅ FASE 0: Autenticación OAuth2

**Estado**: ✅ **COMPLETADA 2024-10-31**
**Commits**: 5/5 ✅
**Tiempo**: 9 horas

### Pasos Completados

- [x] **PASO 0.1**: bcrypt en edugo-shared ✅
  - [x] Crear `auth/password.go` con bcrypt
  - [x] Crear tests (8 tests, 100% passing)
  - [x] Commit + tag `auth/v0.0.1`
  - [x] Actualizar api-mobile
  - [x] Eliminar SHA256 inseguro
  - **Commits**: `8d7005a` (shared), `e8a177c` (api-mobile)

- [x] **PASO 0.3**: Refresh Tokens ✅
  - [x] SUB-PASO 0.3.1: Crear tabla `refresh_tokens`
  - [x] SUB-PASO 0.3.2: Crear `RefreshToken` en shared (tag `auth/v0.0.2`)
  - [x] SUB-PASO 0.3.3: Crear RefreshTokenRepository
  - [x] SUB-PASO 0.3.4: Modificar AuthService (3 métodos nuevos)
  - [x] SUB-PASO 0.3.5: Crear endpoints `/refresh`, `/logout`, `/revoke-all`
  - **Commits**: `8fed9d7` (shared), `24b10f6` (api-mobile)

- [x] **PASO 0.4**: Middleware JWT Compartido ✅
  - [x] Crear `middleware/gin/jwt_auth.go`
  - [x] Crear `middleware/gin/context.go` (helpers tipados)
  - [x] Tests (17 tests, 100% passing)
  - [x] Commit + tag `middleware/gin/v0.0.1`
  - [x] Migrar api-mobile al middleware
  - [x] Eliminar middleware local (−35 líneas)
  - **Commits**: `4330be1` (shared), `c09e347` (api-mobile)

- [x] **PASO 0.5**: Rate Limiting ✅
  - [x] Crear tabla `login_attempts`
  - [x] Crear LoginAttemptRepository
  - [x] Implementar rate limiting en AuthService
  - [x] Max 5 intentos en 15 minutos
  - [x] Tracking de IP + User-Agent
  - **Commits**: `204aeea` (api-mobile)

---

## ✅ FASE 1: Container DI

**Estado**: ✅ **COMPLETADA 2024-10-31**
**Commits**: 1/1 ✅

- [x] Refactorizar cmd/main.go para inicializar PostgreSQL y MongoDB
- [x] Instanciar Container de dependencias
- [x] Reemplazar handlers mock por handlers reales
- [x] Implementar funciones auxiliares (DB, logger, middleware)
- [x] Health check con validación de DBs
- [x] JWT middleware para rutas protegidas
- **Commit**: `3332c05`

---

## 🔵 FASE 2: Completar TODOs de Servicios

**Estado**: 🔵 **EN PROGRESO** (67% completado)
**Commits**: 2/3 ✅
**Esfuerzo estimado**: 1 día restante

### 📍 PRÓXIMA TAREA: PASO 2.3 (completar queries restantes)

- [x] **PASO 2.1**: Implementar RabbitMQ Messaging ✅ (PR #15 merged - commit ce03298)
  - [x] Configurar conexión a RabbitMQ en main.go
  - [x] Crear publisher/producer para eventos
  - [x] Implementar publicación de evento `material_uploaded`
  - [x] Implementar publicación de evento `assessment_attempt_recorded`
  - [x] Agregar publisher al Container DI
  - [x] Integrar eventos con MaterialService y AssessmentService
  - [x] **Archivos creados**:
    - `internal/infrastructure/messaging/rabbitmq/publisher.go` ✅
    - `internal/infrastructure/messaging/rabbitmq/publisher_test.go` ✅
  - [x] **Archivos modificados**:
    - `cmd/main.go` (inicializar RabbitMQ) ✅
    - `internal/container/container.go` ✅
    - `internal/application/service/material_service.go` ✅
    - `internal/application/service/assessment_service.go` ✅
  - [x] Commit: "feat: implementar messaging RabbitMQ para eventos de dominio" ✅

- [x] **PASO 2.2**: Implementar S3 URLs Firmadas ✅ (PR #16 merged - commit 2f2a8af)
  - [x] Configurar cliente AWS S3 desde configuración
  - [x] Implementar generación de presigned URLs (upload + download)
  - [x] Agregar método en MaterialService
  - [x] Integrar con handler CreateMaterial
  - [x] Crear interface S3Storage para testabilidad (commit 5cb1af4)
  - [x] **Archivos creados**:
    - `internal/infrastructure/storage/s3/client.go` ✅
    - `internal/infrastructure/storage/s3/client_test.go` ✅
    - `internal/infrastructure/storage/s3/interface.go` ✅
  - [x] **Archivos modificados**:
    - `internal/application/service/material_service.go` ✅
    - `internal/infrastructure/http/handler/material_handler.go` ✅
    - `internal/container/container.go` ✅
    - `internal/config/config.go` (agregar S3 config) ✅
    - `config/config.yaml` ✅
  - [x] Commit: "feat: implementar URLs presignadas de AWS S3 para materiales" ✅

- [🔵] **PASO 2.3**: Implementar Queries Complejas ⏳ (EN PROGRESO - ~20% completado)
  - [x] **Optimización PostgreSQL**: Índice en `materials.updated_at` ✅ (commits 896ca73 + 59062dd)
    - [x] Crear script SQL `scripts/postgresql/05_indexes_materials.sql`
    - [x] Índice descendente para queries `ORDER BY updated_at DESC`
    - [x] Validar con EXPLAIN ANALYZE (46% mejora local, 5-10x estimado producción)
    - [x] Tests y compilación pasando
  - [ ] Queries de materiales con versiones (PENDIENTE)
  - [ ] Cálculo de puntajes en AssessmentService (PENDIENTE)
  - [ ] Generación de feedback detallado (PENDIENTE)
  - [ ] Actualización de progreso (UPSERT) (PENDIENTE)
  - [ ] Query complejo de estadísticas (PENDIENTE)
  - [ ] **Archivos pendientes a modificar**:
    - `internal/application/service/material_service.go`
    - `internal/application/service/assessment_service.go`
    - `internal/application/service/progress_service.go`
    - `internal/application/service/stats_service.go`
    - `internal/infrastructure/persistence/postgres/repository/material_repository_impl.go`
    - `internal/infrastructure/persistence/mongodb/repository/assessment_repository_impl.go`
  - [ ] Commit pendiente: "feat: implementar consultas complejas en servicios"

---

## ⏳ FASE 3: Limpieza y Consolidación

**Estado**: ⏳ **PENDIENTE**
**Commits**: 0/1
**Esfuerzo estimado**: 0.5-1 día

- [ ] **PASO 3.1**: Eliminar Código Duplicado
  - [ ] Eliminar carpeta `internal/handlers/` (handlers viejos con mocks)
  - [ ] Eliminar archivo `internal/middleware/auth.go` (middleware viejo)
  - [ ] Verificar que no hay referencias
  - [ ] **Archivos a eliminar**:
    - `internal/handlers/auth.go`
    - `internal/handlers/materials.go`
    - `internal/middleware/auth.go`

- [ ] **PASO 3.2**: Consolidar Modelos
  - [ ] Analizar modelos duplicados en `internal/models/`
  - [ ] Migrar a `internal/application/dto/`
  - [ ] Actualizar referencias
  - [ ] Eliminar carpeta `internal/models/` si queda vacía

- [ ] Commit: "refactor: eliminar handlers mock y consolidar modelos"

---

## ⏳ FASE 4: Testing de Integración

**Estado**: ⏳ **PENDIENTE**
**Commits**: 0/1
**Esfuerzo estimado**: 1-2 días

- [ ] **PASO 4.1**: Crear Tests de Integración
  - [ ] Test de flujo completo de autenticación
  - [ ] Test de creación y consulta de materiales
  - [ ] Test de evaluaciones (assessment → intento → puntaje)
  - [ ] Test de actualización de progreso
  - [ ] Test de estadísticas
  - [ ] Verificar health check con DBs reales
  - [ ] **Archivos a crear**:
    - `test/integration/auth_flow_test.go`
    - `test/integration/material_flow_test.go`
    - `test/integration/assessment_flow_test.go`
    - `test/integration/progress_flow_test.go`
  - [ ] Commit: "test: agregar tests de integración para flujo completo"

---

## 📊 Tracking de Progreso - Última Sesión

### **Completado Recientemente (2025-11-05)**:

```
✅ FASE 0: Autenticación OAuth2 (100%)
   ├── ✅ PASO 0.1: bcrypt
   ├── ✅ PASO 0.3: Refresh tokens
   ├── ✅ PASO 0.4: Middleware compartido
   └── ✅ PASO 0.5: Rate limiting

✅ FASE 1: Container DI (100%)

🔵 FASE 2: TODOs de Servicios (67% - 2/3 pasos completados)
   ├── ✅ PASO 2.1: RabbitMQ (PR #15 merged)
   ├── ✅ PASO 2.2: S3 URLs (PR #16 merged)
   └── 🔵 PASO 2.3: Queries complejas (20% - índice materials completado)

Commits: 10 en api-mobile + 3 en shared = 13 total
Tags: auth/v0.0.1, auth/v0.0.2, middleware/gin/v0.0.1
PRs merged: #15 (RabbitMQ), #16 (S3)
```

### **Próxima Sesión - Empezar Aquí** 👇

```
🔵 FASE 2.3: Completar Queries Complejas (80% restante) ← CONTINUAR AQUÍ
   ⏳ Queries de materiales con versiones
   ⏳ Cálculo de puntajes en AssessmentService
   ⏳ Generación de feedback detallado
   ⏳ Actualización de progreso (UPSERT)
   ⏳ Query complejo de estadísticas
```

---

## 🎯 Resumen de Archivos Modificados/Creados

### En edugo-shared (3 tags publicados):
- [x] `auth/password.go` + tests (tag: auth/v0.0.1)
- [x] `auth/refresh_token.go` + tests (tag: auth/v0.0.2)
- [x] `middleware/gin/*.go` + tests (tag: middleware/gin/v0.0.1)

### En edugo-api-mobile (FASE 0 + 1):
- [x] `scripts/postgresql/03_refresh_tokens.sql`
- [x] `scripts/postgresql/04_login_attempts.sql`
- [x] `internal/domain/repository/refresh_token_repository.go`
- [x] `internal/domain/repository/login_attempt_repository.go`
- [x] `internal/infrastructure/persistence/postgres/repository/refresh_token_repository_impl.go`
- [x] `internal/infrastructure/persistence/postgres/repository/login_attempt_repository_impl.go`
- [x] `internal/application/service/auth_service.go` (actualizado)
- [x] `internal/application/dto/auth_dto.go` (actualizado)
- [x] `internal/infrastructure/http/handler/auth_handler.go` (actualizado)
- [x] `internal/infrastructure/http/handler/material_handler.go` (actualizado)
- [x] `internal/infrastructure/http/handler/progress_handler.go` (actualizado)
- [x] `internal/infrastructure/http/handler/assessment_handler.go` (actualizado)
- [x] `internal/container/container.go` (actualizado)
- [x] `cmd/main.go` (actualizado)

### Archivos Nuevos (FASE 2):
- [x] `internal/infrastructure/messaging/rabbitmq/publisher.go` ✅ (PASO 2.1)
- [x] `internal/infrastructure/messaging/rabbitmq/publisher_test.go` ✅ (PASO 2.1)
- [x] `internal/infrastructure/storage/s3/client.go` ✅ (PASO 2.2)
- [x] `internal/infrastructure/storage/s3/client_test.go` ✅ (PASO 2.2)
- [x] `internal/infrastructure/storage/s3/interface.go` ✅ (PASO 2.2)
- [x] `scripts/postgresql/05_indexes_materials.sql` ✅ (PASO 2.3)
- [x] `internal/infrastructure/http/handler/benchmarks_test.go` ✅ (Testing)

---

## 📚 Documentación de Referencia

Para detalles completos de implementación de cada paso, ver:
- **[MASTER_PLAN.md](MASTER_PLAN.md)** - Código completo de cada paso (1,300+ líneas)
- **[AUTH_PROVIDERS_COMPARISON.md](AUTH_PROVIDERS_COMPARISON.md)** - Por qué implementación propia
- **[OAUTH2_ARCHITECTURE_PLAN.md](OAUTH2_ARCHITECTURE_PLAN.md)** - Arquitectura técnica
- **[SOCIAL_LOGIN_ROADMAP.md](SOCIAL_LOGIN_ROADMAP.md)** - Futuro: Google/Apple/Facebook

---

## 🔥 Comandos Rápidos de Retomo

```bash
# Ver estado actual
git status
git log -5 --oneline

# Ver plan maestro
cat sprint/MASTER_PLAN_VISUAL.md

# Buscar próxima tarea
grep "⏳ PASO" sprint/MASTER_PLAN_VISUAL.md | head -1

# Ver documentación detallada del paso
cat sprint/MASTER_PLAN.md | grep -A 50 "PASO 2.1"
```

---

**Última actualización**: 2025-11-05
**Responsable**: Claude Code + Jhoan Medina
**Próxima tarea**: FASE 2.3 (Completar Queries Complejas restantes)
