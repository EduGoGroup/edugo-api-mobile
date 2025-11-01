# Plan de Trabajo - Migración de Mocks a Implementación Real

**📅 Última actualización**: 2024-10-31 23:30
**🎯 Progreso**: 6/11 commits (55%)
**⏱️ Tiempo invertido**: ~9 horas
**👉 Próxima tarea**: [FASE 2.1: RabbitMQ](#21-implementar-messaging-rabbitmq)

---

## 📊 Vista Rápida de Progreso

```
✅ FASE 0: Autenticación OAuth2      COMPLETADA (5 commits)
✅ FASE 1: Container DI              COMPLETADA (1 commit)
⏳ FASE 2: TODOs de Servicios        0/3 commits ← EMPEZAR AQUÍ
⏳ FASE 3: Limpieza                  0/1 commits
⏳ FASE 4: Testing                   0/1 commits
```

---

## 📋 Estado General del Proyecto

**Objetivo**: Conectar toda la implementación real, eliminar mocks, y completar funcionalidades pendientes.

**Branch actual**: `feature/conectar`

---

## ✅ FASE 0: Mejorar Autenticación OAuth2 - **COMPLETADA**

**Estado**: ✅ **COMPLETADA 2025-10-31**
**Commits**: 5 (3 en shared + 2 en api-mobile)
**Tiempo real**: 9 horas

### Pasos Completados

- [x] **PASO 0.1**: bcrypt seguro (Commits: `8d7005a`, `e8a177c`)
- [x] **PASO 0.3**: Refresh tokens (Commits: `8fed9d7`, `24b10f6`)
- [x] **PASO 0.4**: Middleware compartido (Commits: `4330be1`, `c09e347`)
- [x] **PASO 0.5**: Rate limiting (Commit: `204aeea`)

### Tags Publicados en edugo-shared

- [x] `auth/v0.0.1` - bcrypt implementation
- [x] `auth/v0.0.2` - refresh token generator
- [x] `middleware/gin/v0.0.1` - JWT middleware reutilizable

### Mejoras de Seguridad

- [x] bcrypt cost 12 (vs SHA256 inseguro)
- [x] Refresh tokens con revocación en BD
- [x] Logout funcional
- [x] Revocación de todas las sesiones
- [x] Rate limiting (5 intentos/15 min)
- [x] Middleware JWT compartido
- [x] Type-safe helpers (GetUserID, etc.)
- [x] Access tokens 15 min (vs 24 horas antes)

---

## ✅ FASE 1: Conectar Implementación Real con Container DI - **COMPLETADA**

**Estado**: ✅ **COMPLETADA 2025-10-31**
**Commit**: `3332c05`

### Tareas Completadas

- [x] Refactorizar cmd/main.go para inicializar PostgreSQL y MongoDB
- [x] Instanciar Container de dependencias con todas las capas
- [x] Reemplazar handlers mock por handlers reales del Container
- [x] Implementar funciones auxiliares de inicialización (DB, logger, middleware)
- [x] Agregar health check que valida estado de PostgreSQL y MongoDB
- [x] Implementar JWT middleware para autenticación en rutas protegidas
- [x] Conectar handlers reales de auth, material, progress, assessment, summary y stats

### Detalles Técnicos Implementados

- Conexión PostgreSQL con pool configurado y validación de ping
- Conexión MongoDB con timeout y validación de conexión
- Logger Zap inicializado desde configuración
- Container DI inicializa repositorios → servicios → handlers
- CORS middleware configurado
- Eliminación de archivos .gitkeep de carpetas con contenido

### Variables de Entorno Requeridas

La aplicación ahora requiere las siguientes variables de entorno:
- `POSTGRES_PASSWORD` - Contraseña de PostgreSQL
- `MONGODB_URI` - URI de conexión a MongoDB
- `RABBITMQ_URL` - URL de RabbitMQ
- `JWT_SECRET` - Secret para firmar tokens JWT
- `APP_ENV` - Ambiente (local, dev, qa, prod)

---

## 🚧 FASE 2: Completar TODOs de Servicios

**Estado**: ⏳ PENDIENTE

**Estimación**: 3 commits separados por funcionalidad

### Tareas Pendientes

#### 2.1. Implementar Funcionalidad S3

- [ ] Configurar cliente AWS S3 desde configuración
- [ ] Implementar generación de URLs firmadas para subida de materiales
- [ ] Agregar método en MaterialService para generar presigned URLs
- [ ] Integrar con handler CreateMaterial
- [ ] Crear commit: "feat: implementar generación de URLs firmadas S3"

**Archivos a modificar**:
- `internal/application/service/material_service.go`
- `internal/config/config.go` (agregar config de S3)
- `config/config.yaml` (agregar configuración S3)

**TODOs relacionados en código**:
- `internal/handlers/materials.go:line 46` - TODO: Generar URL firmada de S3

---

#### 2.2. Implementar Messaging RabbitMQ

- [ ] Configurar conexión a RabbitMQ en main.go
- [ ] Crear publisher/producer para eventos
- [ ] Implementar publicación de evento `material_uploaded`
- [ ] Implementar publicación de evento `assessment_attempt_recorded`
- [ ] Agregar publisher al Container de dependencias
- [ ] Integrar eventos con servicios correspondientes
- [ ] Crear commit: "feat: implementar messaging RabbitMQ para eventos"

**Archivos a crear**:
- `internal/infrastructure/messaging/rabbitmq/publisher.go`
- `internal/infrastructure/messaging/events.go` (definir eventos)

**Archivos a modificar**:
- `cmd/main.go` (inicializar RabbitMQ)
- `internal/container/container.go` (agregar publisher)
- `internal/application/service/material_service.go`
- `internal/application/service/assessment_service.go`

**TODOs relacionados en código**:
- `internal/handlers/materials.go:line 66` - TODO: Publicar evento material_uploaded a RabbitMQ
- `internal/handlers/materials.go:line 153` - TODO: Publicar evento assessment_attempt_recorded

---

#### 2.3. Implementar Consultas Complejas en Servicios

- [ ] Implementar queries de materiales con versiones
- [ ] Implementar cálculo de puntajes en AssessmentService
- [ ] Implementar generación de feedback detallado
- [ ] Implementar actualización de progreso de lectura (UPSERT)
- [ ] Implementar query complejo de estadísticas
- [ ] Crear commit: "feat: implementar consultas complejas en servicios"

**Archivos a modificar**:
- `internal/application/service/material_service.go`
- `internal/application/service/assessment_service.go`
- `internal/application/service/progress_service.go`
- `internal/application/service/stats_service.go`
- `internal/infrastructure/persistence/postgres/repository/material_repository_impl.go`
- `internal/infrastructure/persistence/postgres/repository/progress_repository_impl.go`
- `internal/infrastructure/persistence/mongodb/repository/assessment_repository_impl.go`

**TODOs relacionados en código**:
- `internal/handlers/materials.go:line 63` - TODO: Registrar versión en material_version
- `internal/handlers/materials.go:line 64` - TODO: Calcular file_hash
- `internal/handlers/materials.go:line 65` - TODO: Verificar deduplicación
- `internal/handlers/materials.go:line 126` - TODO: Validar cada respuesta comparando con correct_answer
- `internal/handlers/materials.go:line 127` - TODO: Generar DetailedFeedback
- `internal/handlers/materials.go:line 128` - TODO: Calcular puntaje
- `internal/handlers/materials.go:line 129` - TODO: Persistir en quiz_attempt
- `internal/handlers/materials.go:line 166` - TODO: Upsert en reading_log con GREATEST
- `internal/handlers/materials.go:line 197` - TODO: Query complejo PostgreSQL

---

## 🧹 FASE 3: Limpieza y Consolidación

**Estado**: ⏳ PENDIENTE

**Estimación**: 1 commit consolidado

### Tareas Pendientes

#### 3.1. Eliminar Código Duplicado

- [ ] Eliminar carpeta `internal/handlers/` (handlers viejos con mocks)
- [ ] Eliminar archivo `internal/middleware/auth.go` (middleware viejo)
- [ ] Verificar que no hay referencias a código eliminado
- [ ] Actualizar imports si es necesario

**Archivos a eliminar**:
- `internal/handlers/auth.go`
- `internal/handlers/materials.go`
- `internal/middleware/auth.go`

---

#### 3.2. Consolidar Modelos

- [ ] Analizar modelos duplicados en `internal/models/`
- [ ] Migrar modelos necesarios a `internal/application/dto/`
- [ ] Actualizar referencias en handlers y servicios
- [ ] Eliminar carpeta `internal/models/` si queda vacía

**Archivos a revisar**:
- `internal/models/request/` vs `internal/application/dto/`
- `internal/models/response/` vs `internal/application/dto/`
- `internal/models/mongodb/` (verificar uso real)

**Decisión pendiente**: Determinar si `internal/models/enum/` debe moverse a `internal/domain/valueobject/` o mantenerse como está.

---

- [ ] Crear commit: "refactor: eliminar handlers mock y consolidar modelos"

---

## 🧪 FASE 4: Testing

**Estado**: ⏳ PENDIENTE

**Estimación**: 1 commit

### Tareas Pendientes

#### 4.1. Tests de Integración

- [ ] Test completo de flujo de autenticación (login → JWT → acceso a recursos)
- [ ] Test de creación y consulta de materiales
- [ ] Test de evaluaciones (obtener assessment → registrar intento → validar puntaje)
- [ ] Test de actualización de progreso de lectura
- [ ] Test de obtención de estadísticas
- [ ] Verificar que health check funciona con DBs reales

**Archivos a crear**:
- `test/integration/auth_flow_test.go`
- `test/integration/material_flow_test.go`
- `test/integration/assessment_flow_test.go`
- `test/integration/progress_flow_test.go`

**Archivos existentes a completar**:
- `test/integration/postgres_test.go` (ya existe, agregar más tests)

---

- [ ] Crear commit: "test: agregar tests de integración para flujo completo"

---

## 📊 Resumen de Progreso

### Commits por Fase

| Fase | Commits | Estado |
|------|---------|--------|
| Fase 0 | 5/5 | ✅ Completada |
| Fase 1 | 1/1 | ✅ Completada |
| Fase 2 | 0/3 | ⏳ Pendiente |
| Fase 3 | 0/1 | ⏳ Pendiente |
| Fase 4 | 0/1 | ⏳ Pendiente |
| **TOTAL** | **6/11** | **55% completado** |

---

## 🎯 Cómo Retomar el Trabajo

### **Inicio de Sesión - 3 Pasos**:

1. **Ver estado del proyecto**:
   ```bash
   git status
   git log -5 --oneline
   ```

2. **Leer vista rápida**:
   ```bash
   cat sprint/README.md | head -20
   # O abrir: sprint/MASTER_PLAN_VISUAL.md
   ```

3. **Buscar próxima tarea sin marcar**:
   - Buscar el primer `- [ ]` en este documento
   - Esa es la siguiente tarea a realizar

### **Durante el Trabajo**:

1. Marcar `- [ ]` como `- [x]` al completar cada tarea
2. Actualizar sección "📊 Vista Rápida de Progreso" arriba
3. Hacer commits atómicos (código que compila)

---

## 📝 Notas Importantes

### Decisiones Tomadas

- ✅ Implementación propia de OAuth2 (vs Firebase/Auth0)
- ✅ bcrypt cost 12 para passwords
- ✅ Refresh tokens con revocación en BD
- ✅ Rate limiting: 5 intentos en 15 minutos
- ✅ Access tokens válidos 15 minutos (renovables)
- ✅ Middleware compartido en edugo-shared
- ✅ Health check mejorado con validación de DBs

### Puntos de Atención para FASE 2

- ⚠️ **RabbitMQ**: Configurar antes de publicar eventos
- ⚠️ **S3**: Configurar cliente AWS antes de generar URLs
- ⚠️ **Queries complejas**: Implementaciones básicas necesitan refinamiento

### Referencias Útiles

- 📄 **Plan detallado**: [MASTER_PLAN.md](MASTER_PLAN.md) (código completo)
- 📄 **Plan visual**: [MASTER_PLAN_VISUAL.md](MASTER_PLAN_VISUAL.md) (checkboxes)
- 📄 **Análisis OAuth2**: [AUTH_PROVIDERS_COMPARISON.md](AUTH_PROVIDERS_COMPARISON.md)
- 📁 **Container DI**: `internal/container/container.go`
- 📁 **Handlers**: `internal/infrastructure/http/handler/`
- 📁 **Servicios**: `internal/application/service/`

---

**Última actualización**: 2025-10-31 23:30
**Responsable**: Claude Code + Jhoan Medina
**Branch**: `feature/conectar`
**Estado**: ✅ 55% completado | ⏳ 3-4 días restantes
