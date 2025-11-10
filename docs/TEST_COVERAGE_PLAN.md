# Plan de Cobertura de Tests - EduGo API Mobile

**Fecha de Creación**: 9 de noviembre de 2025  
**Versión**: 1.0  
**Estado**: En Ejecución

---

## 📋 Resumen Ejecutivo

Este documento define el plan estratégico para incrementar la cobertura de tests del proyecto `edugo-api-mobile` desde el **33.6% actual** hasta alcanzar un **mínimo del 60%**, con metas específicas por módulo según su criticidad.

### Métricas Actuales vs Objetivos

| Categoría | Cobertura Actual | Meta | Prioridad |
|-----------|------------------|------|-----------|
| **Value Objects** | 0% | 100% | 🔴 Crítica |
| **Entities** | 0% | 80% | 🔴 Crítica |
| **Repositories** | 0% | 70% | 🔴 Crítica |
| **Services** | 36.9% | 70% | 🟡 Alta |
| **Handlers** | 41.9% | 60% | 🟡 Alta |
| **Middleware** | 26.5% | 60% | 🟢 Media |
| **Router** | 0% | 50% | 🟢 Media |
| **Total Proyecto** | **33.6%** | **60%+** | - |

### Impacto Esperado

- **Incremento de Cobertura**: +26.4 puntos porcentuales (33.6% → 60%)
- **Módulos Críticos Cubiertos**: 4 módulos (Value Objects, Entities, Repositories, Services)
- **Tests Nuevos Estimados**: ~150-200 tests adicionales
- **Tiempo de Implementación**: 4 semanas

---

## 🎯 Metas de Cobertura por Módulo

### 1. Capa de Dominio (Domain Layer)

#### 1.1 Value Objects - Meta: 100%

**Estado Actual**: 0% (0 tests)  
**Prioridad**: 🔴 **CRÍTICA**  
**Justificación**: Los value objects contienen lógica de validación fundamental que debe estar 100% cubierta.

**Módulos a Cubrir**:

| Archivo | Funciones Críticas | Tests Requeridos | Esfuerzo |
|---------|-------------------|------------------|----------|
| `email.go` | `NewEmail()`, `Validate()` | 8-10 tests | 1h |
| `material_id.go` | `NewMaterialID()`, `Validate()` | 6-8 tests | 1h |
| `user_id.go` | `NewUserID()`, `Validate()` | 6-8 tests | 1h |
| `material_version_id.go` | `NewMaterialVersionID()`, `Validate()` | 6-8 tests | 1h |

**Casos de Prueba Requeridos**:
- ✅ Validación con valores válidos
- ✅ Validación con valores inválidos (vacíos, formato incorrecto)
- ✅ Validación de límites (longitud mínima/máxima)
- ✅ Casos edge (caracteres especiales, Unicode)
- ✅ Serialización/Deserialización (JSON, String)

**Estimación Total**: 4 horas

#### 1.2 Entities - Meta: 80%

**Estado Actual**: 0% (0 tests)  
**Prioridad**: 🔴 **CRÍTICA**  
**Justificación**: Las entities contienen lógica de negocio y reglas de dominio que deben estar bien cubiertas.

**Módulos a Cubrir**:

| Archivo | Funciones Críticas | Tests Requeridos | Esfuerzo |
|---------|-------------------|------------------|----------|
| `material.go` | `NewMaterial()`, `Validate()`, `UpdateStatus()` | 12-15 tests | 2h |
| `material_version.go` | `NewMaterialVersion()`, `Validate()` | 8-10 tests | 1.5h |
| `progress.go` | `NewProgress()`, `UpdateProgress()`, `CalculatePercentage()` | 10-12 tests | 2h |
| `user.go` | `NewUser()`, `Validate()`, `UpdateProfile()` | 10-12 tests | 2h |

**Casos de Prueba Requeridos**:
- ✅ Creación de entities con datos válidos
- ✅ Validación de reglas de negocio
- ✅ Transiciones de estado válidas e inválidas
- ✅ Cálculos y lógica de dominio
- ✅ Relaciones entre entities

**Estimación Total**: 7.5 horas

---

### 2. Capa de Persistencia (Repositories)

#### 2.1 PostgreSQL Repositories - Meta: 70%

**Estado Actual**: 0% (0 tests)  
**Prioridad**: 🔴 **CRÍTICA**  
**Justificación**: Los repositories son la interfaz con la base de datos y deben garantizar persistencia correcta.

**Módulos a Cubrir**:

| Repository | Operaciones Críticas | Tests Requeridos | Esfuerzo |
|------------|---------------------|------------------|----------|
| `user_repository_impl.go` | FindByID, FindByEmail, Update | 8-10 tests | 3h |
| `material_repository_impl.go` | Create, FindByID, Update, List, FindByAuthor | 12-15 tests | 4h |
| `progress_repository_impl.go` | Upsert, FindByMaterialAndUser, CountActiveUsers | 10-12 tests | 3h |
| `refresh_token_repository_impl.go` | Store, FindByTokenHash, Revoke, DeleteExpired | 8-10 tests | 2h |
| `login_attempt_repository_impl.go` | RecordAttempt, GetRecentAttempts | 6-8 tests | 2h |

**Casos de Prueba Requeridos**:
- ✅ CRUD básico (Create, Read, Update, Delete)
- ✅ Búsquedas con resultados existentes
- ✅ Búsquedas sin resultados (not found)
- ✅ Operaciones con datos inválidos
- ✅ Constraints de base de datos (unique, foreign keys)
- ✅ Transacciones y rollback

**Estrategia de Testing**:
- Usar **testcontainers** con PostgreSQL real
- Limpiar datos entre tests (TRUNCATE CASCADE)
- Verificar estado de BD después de operaciones
- Probar casos de concurrencia cuando aplique

**Estimación Total**: 14 horas

#### 2.2 MongoDB Repositories - Meta: 70%

**Estado Actual**: 0% (0 tests)  
**Prioridad**: 🔴 **CRÍTICA**

**Módulos a Cubrir**:

| Repository | Operaciones Críticas | Tests Requeridos | Esfuerzo |
|------------|---------------------|------------------|----------|
| `assessment_repository_impl.go` | SaveAssessment, FindAssessmentByMaterialID, SaveResult | 8-10 tests | 3h |
| `summary_repository_impl.go` | SaveSummary, FindByMaterialID | 6-8 tests | 2h |

**Casos de Prueba Requeridos**:
- ✅ Guardar documentos nuevos
- ✅ Actualizar documentos existentes
- ✅ Búsquedas por ID
- ✅ Índices únicos (duplicados)
- ✅ Documentos anidados

**Estrategia de Testing**:
- Usar **testcontainers** con MongoDB real
- Limpiar colecciones entre tests
- Verificar índices y constraints

**Estimación Total**: 5 horas

---

### 3. Capa de Aplicación (Services)

#### 3.1 Services - Meta: 70%

**Estado Actual**: 36.9% (30 tests existentes)  
**Prioridad**: 🟡 **ALTA**  
**Justificación**: Los services contienen lógica de negocio compleja que requiere cobertura robusta.

**Módulos a Mejorar**:

| Service | Cobertura Actual | Meta | Tests Adicionales | Esfuerzo |
|---------|------------------|------|-------------------|----------|
| `material_service.go` | ~40% | 70% | 8-10 tests | 2h |
| `progress_service.go` | ~35% | 70% | 10-12 tests | 2.5h |
| `stats_service.go` | ~30% | 70% | 8-10 tests | 2h |
| `assessment_service.go` | ~45% | 70% | 6-8 tests | 1.5h |
| `auth_service.go` | ~40% | 70% | 8-10 tests | 2h |
| `summary_service.go` | ~20% | 70% | 10-12 tests | 2.5h |

**Casos de Prueba Adicionales Requeridos**:
- ✅ Casos edge no cubiertos
- ✅ Manejo de errores de repositories
- ✅ Validaciones de datos de entrada
- ✅ Lógica de negocio compleja
- ✅ Interacciones entre múltiples repositories
- ✅ Casos de concurrencia

**Estrategia de Testing**:
- Usar **mocks** para repositories
- Verificar llamadas a dependencies
- Probar todos los paths de ejecución
- Validar transformaciones de datos

**Estimación Total**: 12.5 horas

#### 3.2 Scoring Strategies - Estado: ✅ COMPLETO

**Cobertura Actual**: 95.7%  
**Acción**: Mantener cobertura actual, no requiere trabajo adicional.

---

### 4. Capa de Infraestructura HTTP

#### 4.1 Handlers - Meta: 60%

**Estado Actual**: 41.9% (47 tests existentes)  
**Prioridad**: 🟡 **ALTA**

**Módulos a Mejorar**:

| Handler | Cobertura Actual | Meta | Tests Adicionales | Esfuerzo |
|---------|------------------|------|-------------------|----------|
| `progress_handler.go` | 0% | 60% | 8-10 tests | 2.5h |
| `stats_handler.go` | 0% | 60% | 6-8 tests | 2h |
| `summary_handler.go` | 0% | 60% | 6-8 tests | 2h |
| `material_handler.go` | ~50% | 60% | 4-6 tests | 1.5h |
| `assessment_handler.go` | ~50% | 60% | 4-6 tests | 1.5h |

**Casos de Prueba Requeridos**:
- ✅ Request con datos válidos (200 OK)
- ✅ Request con datos inválidos (400 Bad Request)
- ✅ Request sin autenticación (401 Unauthorized)
- ✅ Request sin permisos (403 Forbidden)
- ✅ Recurso no encontrado (404 Not Found)
- ✅ Errores del service (500 Internal Server Error)

**Estrategia de Testing**:
- Usar **mocks** para services
- Usar `httptest.NewRecorder()` para capturar responses
- Verificar status codes y response bodies
- Probar headers (Content-Type, Authorization)

**Estimación Total**: 9.5 horas

#### 4.2 Middleware - Meta: 60%

**Estado Actual**: 26.5% (1 test existente)  
**Prioridad**: 🟢 **MEDIA**

**Módulos a Mejorar**:

| Middleware | Cobertura Actual | Meta | Tests Adicionales | Esfuerzo |
|------------|------------------|------|-------------------|----------|
| `auth.go` | ~20% | 60% | 8-10 tests | 2.5h |
| `cors.go` | ~40% | 60% | 4-6 tests | 1.5h |

**Casos de Prueba Requeridos**:
- ✅ Token JWT válido
- ✅ Token JWT inválido/expirado
- ✅ Token ausente
- ✅ CORS headers correctos
- ✅ Preflight requests (OPTIONS)

**Estimación Total**: 4 horas

#### 4.3 Router - Meta: 50%

**Estado Actual**: 0%  
**Prioridad**: 🟢 **MEDIA**

**Módulos a Cubrir**:

| Archivo | Tests Requeridos | Esfuerzo |
|---------|------------------|----------|
| `router.go` | 6-8 tests | 2h |

**Casos de Prueba Requeridos**:
- ✅ Rutas registradas correctamente
- ✅ Middleware aplicado en orden correcto
- ✅ Rutas protegidas requieren autenticación
- ✅ Rutas públicas accesibles sin auth

**Estimación Total**: 2 horas

---

### 5. Otros Módulos

#### 5.1 Config - Estado: ✅ COMPLETO

**Cobertura Actual**: 95.9%  
**Acción**: Mantener cobertura actual.

#### 5.2 Bootstrap - Estado: ✅ ACEPTABLE

**Cobertura Actual**: 56.7%  
**Acción**: Mantener cobertura actual, no es crítico mejorar.

#### 5.3 Database Clients - Meta: 50%

**Estado Actual**: ~30%  
**Prioridad**: 🟢 **BAJA**

**Módulos**:
- `mongodb.go`: Tests de conexión y configuración
- `postgres.go`: Tests de conexión y configuración

**Estimación Total**: 2 horas

---

## 📅 Timeline de Implementación

### Semana 1: Capa de Dominio (Crítica)

**Objetivo**: Cubrir Value Objects y Entities al 80%+

| Día | Tareas | Responsable | Horas |
|-----|--------|-------------|-------|
| Lunes | Tests para `email.go`, `material_id.go` | TBD | 2h |
| Martes | Tests para `user_id.go`, `material_version_id.go` | TBD | 2h |
| Miércoles | Tests para `material.go`, `user.go` | TBD | 4h |
| Jueves | Tests para `progress.go`, `material_version.go` | TBD | 3.5h |
| Viernes | Revisión y ajustes | TBD | 1h |

**Entregables**:
- ✅ 4 archivos de test para value objects
- ✅ 4 archivos de test para entities
- ✅ Cobertura de dominio >= 80%

**Criterios de Éxito**:
- Todos los tests pasan
- Cobertura de value objects = 100%
- Cobertura de entities >= 80%

---

### Semana 2: Repositories (Crítica)

**Objetivo**: Cubrir Repositories al 70%+

| Día | Tareas | Responsable | Horas |
|-----|--------|-------------|-------|
| Lunes | Tests para `user_repository_impl.go` | TBD | 3h |
| Martes | Tests para `material_repository_impl.go` | TBD | 4h |
| Miércoles | Tests para `progress_repository_impl.go` | TBD | 3h |
| Jueves | Tests para `assessment_repository_impl.go` (MongoDB) | TBD | 3h |
| Viernes | Tests para `refresh_token_repository_impl.go`, `login_attempt_repository_impl.go` | TBD | 4h |

**Entregables**:
- ✅ 7 archivos de test para repositories
- ✅ Cobertura de repositories >= 70%
- ✅ Tests de integración con testcontainers funcionando

**Criterios de Éxito**:
- Todos los tests pasan
- Testcontainers se levantan correctamente
- Cobertura de repositories >= 70%

---

### Semana 3: Services y Handlers (Alta Prioridad)

**Objetivo**: Mejorar cobertura de Services y Handlers

| Día | Tareas | Responsable | Horas |
|-----|--------|-------------|-------|
| Lunes | Mejorar tests de `material_service.go`, `auth_service.go` | TBD | 4h |
| Martes | Mejorar tests de `progress_service.go`, `stats_service.go` | TBD | 4.5h |
| Miércoles | Mejorar tests de `assessment_service.go`, `summary_service.go` | TBD | 4h |
| Jueves | Tests para `progress_handler.go`, `stats_handler.go` | TBD | 4.5h |
| Viernes | Tests para `summary_handler.go`, mejorar handlers existentes | TBD | 3.5h |

**Entregables**:
- ✅ Tests adicionales para 6 services
- ✅ Tests nuevos para 3 handlers
- ✅ Cobertura de services >= 70%
- ✅ Cobertura de handlers >= 60%

**Criterios de Éxito**:
- Todos los tests pasan
- Cobertura de services >= 70%
- Cobertura de handlers >= 60%

---

### Semana 4: Middleware, Router y Validación Final

**Objetivo**: Completar cobertura restante y validar metas

| Día | Tareas | Responsable | Horas |
|-----|--------|-------------|-------|
| Lunes | Mejorar tests de `auth.go`, `cors.go` middleware | TBD | 4h |
| Martes | Tests para `router.go` | TBD | 2h |
| Miércoles | Tests para database clients (opcional) | TBD | 2h |
| Jueves | Ejecutar suite completa, verificar cobertura | TBD | 3h |
| Viernes | Ajustes finales, documentación | TBD | 3h |

**Entregables**:
- ✅ Tests para middleware
- ✅ Tests para router
- ✅ Cobertura total >= 60%
- ✅ Reporte final de cobertura
- ✅ Documentación actualizada

**Criterios de Éxito**:
- Todos los tests pasan (100%)
- Cobertura total >= 60%
- Todas las metas por módulo alcanzadas
- CI/CD ejecutando tests automáticamente

---

## 👥 Asignación de Responsables

### Opción 1: Equipo Completo

| Responsable | Área Asignada | Semanas |
|-------------|---------------|---------|
| **Desarrollador A** | Domain Layer (Value Objects + Entities) | Semana 1 |
| **Desarrollador B** | Repositories (PostgreSQL + MongoDB) | Semana 2 |
| **Desarrollador C** | Services | Semana 3 |
| **Desarrollador D** | Handlers + Middleware | Semanas 3-4 |
| **Tech Lead** | Revisión, validación final | Semana 4 |

### Opción 2: Desarrollador Individual

| Semana | Área | Horas Estimadas |
|--------|------|-----------------|
| Semana 1 | Domain Layer | 12.5h |
| Semana 2 | Repositories | 19h |
| Semana 3 | Services + Handlers | 20.5h |
| Semana 4 | Middleware + Validación | 14h |
| **Total** | - | **66 horas** |

**Nota**: Si se trabaja solo, el timeline se puede extender a 6-8 semanas trabajando 8-10 horas por semana.

---

## 🎯 Priorización de Tests Faltantes

### Prioridad 1: CRÍTICA (Semanas 1-2)

**Impacto**: Alto - Afecta lógica de negocio fundamental

1. **Value Objects** (4h)
   - `email.go` - Validación de emails
   - `material_id.go` - Validación de IDs
   - `user_id.go` - Validación de IDs de usuario
   - `material_version_id.go` - Validación de versiones

2. **Entities** (7.5h)
   - `material.go` - Lógica de materiales
   - `progress.go` - Cálculo de progreso
   - `user.go` - Validación de usuarios
   - `material_version.go` - Versionado

3. **Repositories** (19h)
   - `user_repository_impl.go` - Persistencia de usuarios
   - `material_repository_impl.go` - Persistencia de materiales
   - `progress_repository_impl.go` - Persistencia de progreso
   - `assessment_repository_impl.go` - Persistencia de assessments

**Total Prioridad 1**: 30.5 horas

### Prioridad 2: ALTA (Semana 3)

**Impacto**: Medio-Alto - Mejora confiabilidad de servicios

4. **Services** (12.5h)
   - Mejorar cobertura de 6 services existentes
   - Agregar tests para casos edge
   - Validar manejo de errores

5. **Handlers Críticos** (6.5h)
   - `progress_handler.go` - Nuevo
   - `stats_handler.go` - Nuevo
   - `summary_handler.go` - Nuevo

**Total Prioridad 2**: 19 horas

### Prioridad 3: MEDIA (Semana 4)

**Impacto**: Medio - Completa cobertura general

6. **Handlers Existentes** (3h)
   - Mejorar `material_handler.go`
   - Mejorar `assessment_handler.go`

7. **Middleware** (4h)
   - Mejorar `auth.go`
   - Mejorar `cors.go`

8. **Router** (2h)
   - Tests para `router.go`

**Total Prioridad 3**: 9 horas

### Prioridad 4: BAJA (Opcional)

**Impacto**: Bajo - Nice to have

9. **Database Clients** (2h)
10. **Otros módulos** (variable)

---

## 📊 Métricas de Seguimiento

### Métricas Semanales

Ejecutar al final de cada semana:

```bash
make coverage-report
```

**Métricas a Rastrear**:

| Métrica | Semana 1 | Semana 2 | Semana 3 | Semana 4 | Meta Final |
|---------|----------|----------|----------|----------|------------|
| **Cobertura Total** | 40% | 48% | 56% | 60%+ | 60%+ |
| **Value Objects** | 100% | 100% | 100% | 100% | 100% |
| **Entities** | 80% | 80% | 80% | 80% | 80% |
| **Repositories** | 0% | 70% | 70% | 70% | 70% |
| **Services** | 37% | 37% | 70% | 70% | 70% |
| **Handlers** | 42% | 42% | 60% | 60% | 60% |
| **Tests Totales** | +30 | +50 | +60 | +10 | +150 |

### Dashboard de Progreso

Crear un dashboard simple en el README o en GitHub Actions:

```markdown
## 📊 Test Coverage Progress

| Module | Current | Target | Status |
|--------|---------|--------|--------|
| Value Objects | 100% | 100% | ✅ |
| Entities | 80% | 80% | ✅ |
| Repositories | 45% | 70% | 🔄 |
| Services | 55% | 70% | 🔄 |
| Handlers | 50% | 60% | 🔄 |
| **Total** | **52%** | **60%** | 🔄 |
```

---

## ✅ Criterios de Éxito

### Criterios Obligatorios

- ✅ **Cobertura total >= 60%**
- ✅ **Value Objects = 100%**
- ✅ **Entities >= 80%**
- ✅ **Repositories >= 70%**
- ✅ **Services >= 70%**
- ✅ **Handlers >= 60%**
- ✅ **Todos los tests pasan (100%)**
- ✅ **CI/CD ejecuta tests automáticamente**
- ✅ **Build falla si cobertura cae por debajo del umbral**

### Criterios Deseables

- ✅ Middleware >= 60%
- ✅ Router >= 50%
- ✅ Documentación de testing completa
- ✅ Badges de cobertura en README
- ✅ Reportes de cobertura publicados automáticamente

---

## 🚧 Riesgos y Mitigaciones

### Riesgo 1: Testcontainers Lentos

**Descripción**: Tests de integración con testcontainers pueden ser lentos (15-20s cada uno)

**Impacto**: Alto - Puede ralentizar desarrollo

**Mitigación**:
- Reutilizar contenedores entre tests cuando sea posible
- Ejecutar tests de integración solo en CI, no localmente
- Usar `make test-unit` para desarrollo rápido
- Paralelizar tests de integración

### Riesgo 2: Falta de Tiempo

**Descripción**: 66 horas de trabajo pueden no ser suficientes

**Impacto**: Alto - No alcanzar metas

**Mitigación**:
- Priorizar tests críticos (Prioridad 1 y 2)
- Extender timeline si es necesario
- Asignar más desarrolladores
- Reducir meta de cobertura a 55% si es necesario

### Riesgo 3: Tests Frágiles

**Descripción**: Tests pueden romperse con cambios pequeños en el código

**Impacto**: Medio - Mantenimiento costoso

**Mitigación**:
- Usar mocks apropiadamente
- No testear detalles de implementación
- Testear comportamiento, no estructura
- Revisar tests en code reviews

### Riesgo 4: Cobertura sin Calidad

**Descripción**: Alcanzar cobertura alta pero con tests de baja calidad

**Impacto**: Alto - Falsa sensación de seguridad

**Mitigación**:
- Revisar calidad de tests en code reviews
- Seguir patrón AAA (Arrange-Act-Assert)
- Testear casos edge y errores
- Usar mutation testing (opcional)

---

## 📚 Recursos y Referencias

### Documentación Interna

- **Guía de Testing**: `docs/TESTING_GUIDE.md`
- **Guía de Tests Unitarios**: `docs/TESTING_UNIT_GUIDE.md`
- **Guía de Tests de Integración**: `docs/TESTING_INTEGRATION_GUIDE.md`
- **Reporte de Análisis**: `docs/TEST_ANALYSIS_REPORT.md`

### Comandos Útiles

```bash
# Tests unitarios rápidos
make test-unit

# Tests de integración
make test-integration

# Cobertura completa
make coverage-report

# Verificar umbral mínimo
make coverage-check

# Ver cobertura en navegador
open coverage/coverage.html
```

### Herramientas

- **testify**: Assertions y mocks
- **testcontainers-go**: Contenedores para tests
- **go tool cover**: Análisis de cobertura
- **golangci-lint**: Linting de código

---

## 📝 Notas Finales

### Principios de Testing

1. **Tests Rápidos**: Tests unitarios deben ser < 100ms
2. **Tests Independientes**: Cada test debe poder ejecutarse solo
3. **Tests Legibles**: Nombres descriptivos, patrón AAA
4. **Tests Mantenibles**: No testear detalles de implementación
5. **Tests Valiosos**: Testear comportamiento, no cobertura por cobertura

### Mantenimiento Continuo

Una vez alcanzadas las metas:

- **Mantener cobertura**: No permitir que baje del 60%
- **Tests en PRs**: Requerir tests para nuevo código
- **Revisión de tests**: Incluir tests en code reviews
- **Refactoring**: Mejorar tests existentes continuamente
- **Documentación**: Mantener guías actualizadas

### Próximos Pasos

1. ✅ Revisar y aprobar este plan
2. ✅ Asignar responsables
3. ✅ Comenzar Semana 1: Domain Layer
4. ✅ Ejecutar métricas semanales
5. ✅ Ajustar plan según progreso

---

**Última actualización**: 9 de noviembre de 2025  
**Próxima revisión**: 16 de noviembre de 2025 (fin de Semana 1)

