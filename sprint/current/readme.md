# Sprint: Mejoras y Refactorizaciones - edugo-api-mobile

> **Fecha de inicio:** 2024-12-23  
> **Estado:** En Ejecución - Fase 1  
> **Branch base:** `dev`  
> **Branch activo:** `feature/unify-config-standard` (reutilizado para Fase 1)

---

## Resumen Ejecutivo

Este sprint aborda las mejoras documentadas en `documents/improvements/`, organizadas en **6 fases** incrementales. Cada fase:
- Crea una rama desde `dev`
- Implementa los cambios
- Ejecuta tests, lint y compilación
- Crea PR a `dev`

---

## Fases del Sprint

| Fase | Nombre | Prioridad | Esfuerzo | Commits Est. |
|------|--------|-----------|----------|--------------|
| 1 | Deuda Técnica Crítica | 🔴 Alta | 4-6h | 3-4 |
| 2 | TODOs de Autorización | 🔴 Alta | 4-6h | 2-3 |
| 3 | TODOs de Funcionalidad | 🟡 Media | 6-8h | 4-5 |
| 4 | Refactorizaciones de Infraestructura | 🟡 Media | 4-6h | 3-4 |
| 5 | Limpieza de Código Legacy | 🟢 Baja | 2-4h | 2-3 |
| 6 | Mejoras de Observabilidad | 🟢 Baja | 4-6h | 3-4 |

**Total estimado:** 24-36 horas de desarrollo

---

## Fase 1: Deuda Técnica Crítica

**Branch:** `feature/debt-critical`  
**Prioridad:** 🔴 Alta  
**Duración estimada:** 4-6 horas

### Objetivo
Resolver la deuda técnica más crítica que afecta la funcionalidad del sistema.

### Tareas

- [x] **DEBT-003**: Resolver SchoolID hardcodeado ✅ (23 Dic 2024)
  - Archivo: `internal/application/service/material_service.go:63-64`
  - ✅ Agregado `SchoolID` a JWT Claims en `edugo-shared` (release auth/v0.10.0)
  - ✅ Creado helper `GetSchoolIDFromContext()` en middleware
  - ✅ Creado helper `MustGetSchoolIDFromContext()` en middleware
  - ✅ Actualizado `MaterialService.CreateMaterial` para usar schoolID del contexto
  - ✅ Actualizados tests y mocks con nuevo parámetro
  - **Commit:** `fix(material): obtener schoolID del contexto de autenticación (DEBT-003)`

- [ ] **DEBT-005**: Resolver tests unitarios con TODOs
  - Archivos afectados:
    - `answer_repository_test.go`
    - `assessment_repository_test.go`
    - `assessment_document_repository_test.go`
  - Decisión: Eliminar tests unitarios comentados (los de integración son suficientes)
  - **Commit:** `test: eliminar tests unitarios redundantes con TODOs`

- [ ] **DEBT-006**: Estandarizar uso de logger
  - Auditar archivos con logging inconsistente
  - Estandarizar en formato key-value pairs
  - **Commit:** `refactor(logger): estandarizar formato de logging`

### Validación
```bash
go build ./...
go test ./...
golangci-lint run
```

### PR
- **Título:** `fix: resolver deuda técnica crítica (DEBT-003, DEBT-005, DEBT-006)`
- **Labels:** `debt`, `priority-high`

---

## Fase 2: TODOs de Autorización

**Branch:** `feature/auth-todos`  
**Prioridad:** 🔴 Alta  
**Duración estimada:** 4-6 horas  
**Dependencia:** Fase 1 (usa helper de contexto)

### Objetivo
Completar funcionalidades de autorización pendientes.

### Tareas

- [ ] **TODO-003**: Verificación de rol admin en Progress Handler
  - Archivo: `internal/infrastructure/http/handler/progress_handler.go:109-110`
  - Agregar bypass para roles `admin` y `super_admin`
  - **Commit:** `feat(progress): agregar bypass de admin para actualizar progreso de otros usuarios`

- [ ] **Crear middleware genérico de autorización por rol**
  - Archivo: `internal/infrastructure/http/middleware/authorization.go`
  - Crear `RequireRole(roles ...string)` middleware
  - Crear `RequireAnyRole(roles ...string)` middleware
  - **Commit:** `feat(middleware): crear middleware genérico de autorización por rol`

- [ ] **Aplicar middleware a endpoints sensibles**
  - Revisar router y aplicar donde sea necesario
  - **Commit:** `feat(router): aplicar middleware de autorización a endpoints sensibles`

### Validación
```bash
go build ./...
go test ./...
golangci-lint run
```

### PR
- **Título:** `feat(auth): completar TODOs de autorización`
- **Labels:** `feature`, `security`, `priority-high`

---

## Fase 3: TODOs de Funcionalidad

**Branch:** `feature/functionality-todos`  
**Prioridad:** 🟡 Media  
**Duración estimada:** 6-8 horas

### Objetivo
Completar funcionalidades pendientes relacionadas con eventos y persistencia.

### Tareas

- [ ] **TODO-004**: URL real de S3 en Material Service
  - Archivo: `internal/application/service/material_service.go:116-117`
  - Mover publicación de evento a `NotifyUploadComplete`
  - Usar datos reales de S3 en el payload
  - **Commit:** `fix(material): usar URL real de S3 en evento MaterialUploaded`

- [ ] **TODO-006**: Implementar FindByIDWithVersions completo
  - Archivo: `internal/infrastructure/persistence/postgres/repository/material_repository_impl.go:369`
  - Implementar join con `material_versions`
  - Agregar tests de integración
  - **Commit:** `feat(material): implementar FindByIDWithVersions con join a versiones`

- [ ] **TODO-007**: Publicar evento material_completed
  - Archivo: `internal/application/service/progress_service.go:110-118`
  - Definir estructura del evento
  - Implementar publicación cuando progress = 100%
  - **Commit:** `feat(progress): publicar evento material_completed cuando progreso llega a 100%`

- [ ] **TODO-005**: Preparar restauración de eventos de Assessment
  - Documentar qué schema se necesita en edugo-infrastructure
  - Crear issue/ticket para definir schema
  - **Commit:** `docs: documentar requerimientos para eventos de assessment`

### Validación
```bash
go build ./...
go test ./...
golangci-lint run
```

### PR
- **Título:** `feat: completar TODOs de funcionalidad (eventos y persistencia)`
- **Labels:** `feature`, `priority-medium`

---

## Fase 4: Refactorizaciones de Infraestructura

**Branch:** `feature/infra-refactor`  
**Prioridad:** 🟡 Media  
**Duración estimada:** 4-6 horas

### Objetivo
Mejorar la resiliencia y robustez de la infraestructura.

### Tareas

- [ ] **REF-004**: Implementar Circuit Breaker para servicios externos
  - Crear `internal/infrastructure/messaging/resilient_publisher.go`
  - Usar `sony/gobreaker` (ya en go.mod)
  - Integrar en bootstrap
  - **Commit:** `feat(infra): implementar circuit breaker para RabbitMQ publisher`

- [ ] **REF-006**: Implementar Healthcheck detallado
  - Mejorar `HealthHandler` con checks individuales
  - Agregar parámetro `?detail=1` para info detallada
  - Incluir latencias y estados de cada servicio
  - **Commit:** `feat(health): implementar healthcheck detallado con checks individuales`

- [ ] **TODO-008**: Implementar lógica de deshabilitación de recursos
  - Archivo: `internal/bootstrap/config.go:96-97`
  - Completar `WithDisabledResource` para deshabilitar recursos
  - **Commit:** `feat(bootstrap): implementar deshabilitación completa de recursos`

### Validación
```bash
go build ./...
go test ./...
golangci-lint run
```

### PR
- **Título:** `feat(infra): refactorizaciones de infraestructura (circuit breaker, healthcheck)`
- **Labels:** `infrastructure`, `refactor`, `priority-medium`

---

## Fase 5: Limpieza de Código Legacy

**Branch:** `feature/legacy-cleanup`  
**Prioridad:** 🟢 Baja  
**Duración estimada:** 2-4 horas

### Objetivo
Eliminar código legacy y deprecado que ya no se usa.

### Tareas

- [ ] **DEP-002**: Limpiar repositorio legacy de Assessments
  - Verificar que no hay referencias activas
  - Eliminar o marcar claramente como legacy
  - Actualizar documentación
  - **Commit:** `refactor: limpiar referencias a repositorio legacy de assessments`

- [ ] **DEBT-004**: Documentar plan de consolidación de sistemas Assessment
  - Crear documento de migración
  - Definir timeline para eliminación completa
  - **Commit:** `docs: crear plan de consolidación de sistemas de assessment`

- [ ] **Eliminar código comentado restante**
  - Buscar bloques de código comentado
  - Eliminar o crear issues para funcionalidad faltante
  - **Commit:** `refactor: eliminar código comentado residual`

### Validación
```bash
go build ./...
go test ./...
golangci-lint run
```

### PR
- **Título:** `refactor: limpieza de código legacy y deprecado`
- **Labels:** `cleanup`, `refactor`, `priority-low`

---

## Fase 6: Mejoras de Observabilidad

**Branch:** `feature/observability`  
**Prioridad:** 🟢 Baja  
**Duración estimada:** 4-6 horas

### Objetivo
Mejorar la capacidad de debugging y monitoreo del sistema.

### Tareas

- [ ] **REF-005**: Agregar Request ID y Tracing
  - Crear middleware `RequestID()` en `internal/infrastructure/http/middleware/`
  - Propagar request_id en logs
  - Propagar en headers de RabbitMQ
  - **Commit:** `feat(observability): agregar middleware de Request ID`

- [ ] **Mejorar logging estructurado**
  - Agregar request_id a todos los logs de handlers
  - Agregar contexto adicional (endpoint, method, duration)
  - **Commit:** `feat(logging): mejorar logging estructurado con contexto`

- [ ] **Agregar métricas básicas**
  - Contador de requests por endpoint
  - Histograma de latencias
  - Contador de errores por tipo
  - **Commit:** `feat(metrics): agregar métricas básicas de endpoints`

### Validación
```bash
go build ./...
go test ./...
golangci-lint run
```

### PR
- **Título:** `feat(observability): mejoras de observabilidad (request ID, logging, métricas)`
- **Labels:** `observability`, `feature`, `priority-low`

---

## Dependencias Externas

### ✅ Cambios en api-admin - COMPLETADO (23 Dic 2024)

La dependencia de api-admin para resolver **DEBT-003** ya fue implementada:

1. ✅ **`school_id` agregado al JWT** en api-admin (PR #64 - merged)
2. ✅ **Columna `school_id` en User** en infrastructure (postgres/v0.13.0)
3. ✅ **Endpoint `POST /auth/switch-context`** para cambio de escuela
4. ✅ Ver detalles actualizados en `documents/improvements/API-ADMIN-REQUIREMENTS.md`

**Estado:** Ya no se requiere workaround temporal. Se puede usar directamente el `school_id` del JWT.

---

## Checklist de Validación por Fase

Antes de crear cada PR, verificar:

```bash
# 1. Compilación
go build ./...

# 2. Tests
go test ./... -v

# 3. Linting
golangci-lint run

# 4. Formateo
go fmt ./...

# 5. Verificar imports
goimports -w .

# 6. Pre-commit hooks
pre-commit run --all-files
```

---

## Orden de Ejecución

```
Fase 1 (Deuda Técnica) ──┐
                         ├──► Fase 3 (Funcionalidad)
Fase 2 (Autorización) ───┘
         │
         ▼
Fase 4 (Infraestructura)
         │
         ▼
Fase 5 (Legacy Cleanup)
         │
         ▼
Fase 6 (Observabilidad)
```

**Nota:** Fases 1 y 2 pueden ejecutarse en paralelo si diferentes personas las trabajan.

---

## Resumen de Archivos Clave

| Archivo | Fases que lo modifican |
|---------|----------------------|
| `material_service.go` | 1, 3 |
| `progress_handler.go` | 2 |
| `progress_service.go` | 3 |
| `material_repository_impl.go` | 3 |
| `router.go` | 2 |
| `bootstrap.go` / `config.go` | 4 |
| `health_handler.go` | 4 |
| Middleware (nuevos) | 2, 6 |

---

## Historial de Cambios

| Fecha | Cambio | Autor |
|-------|--------|-------|
| 2024-12-23 | Creación del plan | Claude Code |
| 2024-12-23 | Dependencia api-admin completada (school_id en JWT) | Claude Code |
| 2024-12-23 | Inicio de Fase 1 | Claude Code |
| 2024-12-23 | **DEBT-003 completado** - SchoolID del contexto JWT | Claude Code |

---

**Próximo paso:** Comenzar con Fase 1 - crear rama `feature/debt-critical` desde `dev`
