# Plan de Trabajo - Migración de Mocks a Implementación Real

Este documento contiene el plan de trabajo para completar la migración del proyecto desde handlers mock a la implementación real con Container DI.

---

## 📋 Estado General del Proyecto

**Objetivo**: Conectar toda la implementación real, eliminar mocks, y completar funcionalidades pendientes.

**Branch actual**: `feature/conectar`

---

## ✅ FASE 1: Conectar Implementación Real con Container DI

**Estado**: ✅ COMPLETADA

**Commit**: `3332c05` - "feat: conectar implementación real con Container DI"

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

### Commits Estimados

| Fase | Commits | Estado |
|------|---------|--------|
| Fase 1 | 1 | ✅ Completado |
| Fase 2 | 3 | ⏳ Pendiente |
| Fase 3 | 1 | ⏳ Pendiente |
| Fase 4 | 1 | ⏳ Pendiente |
| **TOTAL** | **6** | **1/6 completados** |

### Archivos Principales Modificados en Fase 1

- [x] `cmd/main.go` - Refactorizado completamente (+192 líneas)
- [x] `internal/application/service/material_service.go` - Formateo menor
- [x] Eliminados: `internal/domain/.gitkeep`, `internal/infrastructure/http/.gitkeep`

---

## 🎯 Próximos Pasos Inmediatos

**Cuando reanudes el trabajo**:

1. Revisar este documento (`sprint/README.md`)
2. Verificar el estado del branch `feature/conectar`
3. Continuar con **FASE 2.1**: Implementar funcionalidad S3
4. Marcar casillas completadas según avances
5. Actualizar este documento con hallazgos o cambios al plan

---

## 📝 Notas Importantes

### Decisiones Tomadas

- ✅ Los handlers viejos (`internal/handlers/`) NO fueron eliminados en Fase 1 para mantener atomicidad del commit
- ✅ Se decidió usar `logger.NewZapLogger()` en lugar de `logger.NewLogger()`
- ✅ CORS configurado como wildcard (*) por ahora, puede ajustarse en producción
- ✅ Health check mejorado con validación real de conexiones a DBs

### Puntos de Atención

- ⚠️ **RabbitMQ**: Aún no está conectado. La aplicación fallará si intenta publicar eventos.
- ⚠️ **S3**: No configurado. Las subidas de materiales no generarán URLs aún.
- ⚠️ **Queries complejas**: Algunos servicios tienen implementaciones básicas que necesitan refinamiento.
- ⚠️ **Variables de entorno**: Asegurarse de tenerlas todas configuradas antes de ejecutar.

### Referencias Útiles

- **Container DI**: `internal/container/container.go`
- **Handlers Reales**: `internal/infrastructure/http/handler/`
- **Servicios**: `internal/application/service/`
- **Repositorios**: `internal/infrastructure/persistence/{postgres,mongodb}/repository/`
- **Tests Existentes**: `test/integration/`

---

**Última actualización**: 2025-10-31
**Responsable**: Claude Code
**Branch**: `feature/conectar`
