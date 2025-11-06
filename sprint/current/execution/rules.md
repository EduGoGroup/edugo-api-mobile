# Reglas de Ejecución del Sprint

## 📋 Estado de Fases

### Fases del Sprint

- [x] **Fase 1**: Preparación de Infraestructura de Base de Datos (7 tareas)
- [x] **Fase 2**: Implementar Queries de Materiales con Versionado (5 tareas)
- [ ] **Fase 3**: Implementar Cálculo de Puntajes en Evaluaciones (8 tareas)
- [ ] **Fase 4**: Implementar Generación de Feedback Detallado (6 tareas)
- [ ] **Fase 5**: Implementar UPSERT de Progreso (6 tareas)
- [ ] **Fase 6**: Implementar Estadísticas Globales (9 tareas)
- [ ] **Fase 7**: Validación Integral y Refinamiento (7 tareas)
- [ ] **Fase 8**: Commit Atómico y Preparación para PR (4 tareas)

**Total**: 8 fases, 52 tareas granulares

---

## 🎯 Instrucciones para el Agente de Ejecución

### Regla 1: Ejecución Incremental
- **Solo ejecutar UNA fase a la vez**
- Completar todas las tareas de la fase antes de continuar
- NO saltar fases ni ejecutar parcialmente

### Regla 2: Selección de Fase
- **Identificar la próxima fase disponible** según la lista de arriba
- La próxima fase es la primera con `[ ]` (casilla sin marcar)
- Leer el plan detallado en `sprint/current/planning/readme.md` para obtener las tareas específicas de esa fase

### Regla 3: Dependencias
- **Respetar dependencias entre fases** según el plan
- Si una fase depende de otra no completada, DETENER y reportar
- Consultar el grafo de dependencias en `sprint/current/planning/readme.md`

### Regla 4: Validaciones Obligatorias
- **Código debe compilar**: Ejecutar `go build ./...` después de cada tarea
- **Tests deben pasar**: Ejecutar `go test ./...` si hay tests nuevos
- **Solo marcar tarea como completada** si todas las validaciones pasan

### Regla 5: Actualización de Estado
Cuando una fase se complete exitosamente:

1. **Generar reporte de ejecución**: `sprint/current/execution/fase-N-[timestamp].md`
2. **Hacer commit atómico** según mensaje sugerido en el plan
3. **Actualizar este archivo (`rules.md`)**:
   - Marcar la fase como completada: `- [x]`
   - Agregar resumen de la fase al final de este documento (sección "Resúmenes de Fases Completadas")
4. **Actualizar plan**: Marcar casillas en `sprint/current/planning/readme.md`

### Regla 6: Manejo de Errores
Si encuentras un error que no puedes resolver:
- **DETENER inmediatamente**
- **NO continuar** con tareas dependientes
- **Generar reporte** con el error y contexto
- **NO actualizar** este archivo (dejar fase sin marcar)
- **Presentar opciones** al usuario

### Regla 7: Contexto de Fases Previas
- **Leer los resúmenes** al final de este documento antes de comenzar
- Los resúmenes proveen contexto de lo que se implementó en fases anteriores
- Usa este contexto para mantener consistencia arquitectónica

---

## 📊 Resúmenes de Fases Completadas

### ✅ Fase 1: Preparación de Infraestructura de Base de Datos

**Fecha de completitud**: 2025-11-05
**Commit**: Incluido en fase 2 (no se hizo commit separado)

**Resumen**:
- **Scripts creados**:
  - `scripts/postgresql/04_material_versions.sql` - Tabla y índices para versionado de materiales
  - `scripts/postgresql/05_user_progress_upsert.sql` - Constraint UNIQUE e índices para UPSERT de progreso
  - `scripts/mongodb/02_assessment_results.js` - Colección e índices para resultados de evaluaciones

- **Estructuras de BD creadas**:
  - Tabla `material_versions` con constraint UNIQUE(material_id, version_number)
  - Índices: `idx_material_versions_material_id`, `idx_material_versions_created_at`
  - Constraint `unique_user_material UNIQUE(user_id, material_id)` en `user_progress`
  - Índices: `idx_user_progress_user_material`, `idx_user_progress_updated_at`
  - Colección MongoDB `assessment_results` con índice UNIQUE en {assessment_id, user_id}

- **Validación**: Scripts ejecutables sin errores (validado mediante tests de siguiente fase)

**Archivos clave**:
- `scripts/postgresql/04_material_versions.sql`
- `scripts/postgresql/05_user_progress_upsert.sql`
- `scripts/mongodb/02_assessment_results.js`

**Impacto**: Todas las tablas, colecciones e índices necesarios para las siguientes fases están creados y listos para uso.

---

### ✅ Fase 2: Implementar Queries de Materiales con Versionado

**Fecha de completitud**: 2025-11-05 21:49
**Commit**: `4d6e5a2` - "feat(materials): agregar endpoint para consultar materiales con versionado histórico"
**Reporte completo**: `sprint/current/execution/fase-2-2025-11-05-2149.md`

**Resumen**:
- **Funcionalidad implementada**: Endpoint `GET /api/v1/materials/{id}/versions` que retorna material con historial completo de versiones

- **Implementación técnica**:
  - Query SQL con LEFT JOIN a `material_versions` ordenado por `version_number DESC`
  - Método `FindByIDWithVersions()` en MaterialRepositoryImpl
  - Método `GetMaterialWithVersions()` en MaterialService con logging y validación
  - Handler HTTP con validación de UUID y códigos apropiados (200, 400, 404, 500)
  - DTOs: `MaterialVersionResponse`, `MaterialWithVersionsResponse`

- **Tests creados**: 5 tests unitarios (100% coverage)
  - Material con versiones (happy path)
  - Material sin versiones (array vacío)
  - Material no encontrado (404)
  - UUID inválido (400)
  - Error de base de datos (500)

- **Decisiones arquitectónicas**:
  - LEFT JOIN vs INNER JOIN: Se usó LEFT JOIN para incluir materiales sin versiones
  - Array vacío vs null: Siempre retornar array vacío para evitar null checks en frontend
  - Logging de tiempo de ejecución: Medición de performance del query complejo

**Archivos modificados**:
- `internal/application/dto/material_dto.go` (+50 líneas)
- `internal/application/service/material_service.go` (+80 líneas)
- `internal/application/service/material_service_test.go` (+350 líneas, nuevo)
- `internal/infrastructure/http/handler/material_handler.go` (+60 líneas)
- `internal/infrastructure/http/router/router.go` (+2 líneas)
- `internal/infrastructure/persistence/postgres/repository/material_repository_impl.go` (+120 líneas)

**Validaciones**:
- ✅ Compilación exitosa
- ✅ 5/5 tests pasando
- ✅ Código protegido con JWT
- ✅ Manejo de errores con error types de edugo-shared

**Problemas resueltos durante ejecución**:
1. Missing import "time" → Agregado
2. Mock incompleto de Logger → Agregados métodos Sync() y With()
3. Error codes incorrectos → Corregidos a ErrorCodeNotFound, etc.
4. mock.Anything en time.Time → Cambiado a time.Now()

**Impacto**: Funcionalidad de versionado de materiales lista para producción. Frontend puede consultar historial completo de cambios en materiales educativos.

---

## 🎯 Próxima Fase a Ejecutar

**Fase 3**: Implementar Cálculo de Puntajes en Evaluaciones (8 tareas)

**Tareas a ejecutar** (ver detalle en `sprint/current/planning/readme.md`):
- 3.1 - Definir interfaces de Strategy Pattern
- 3.2 - MultipleChoiceStrategy
- 3.3 - TrueFalseStrategy
- 3.4 - ShortAnswerStrategy
- 3.5 - SaveResult en AssessmentRepositoryImpl
- 3.6 - CalculateScore en AssessmentService
- 3.7 - Tests para ScoringStrategy
- 3.8 - Tests para AssessmentService

**Commit esperado**: `feat(assessments): implementar cálculo automático de puntajes con Strategy Pattern`

---

_Este archivo es actualizado automáticamente por el agente de ejecución después de completar cada fase._

_Última actualización: 2025-11-05 21:49 - Fase 2 completada_
