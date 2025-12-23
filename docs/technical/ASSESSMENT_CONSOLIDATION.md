# Plan de Consolidacion de Sistemas de Assessment

> **Fecha de creacion:** 2024-12-23  
> **Estado:** Planificado  
> **Prioridad:** Media  
> **Esfuerzo estimado:** 8-12 horas

---

## Resumen Ejecutivo

El proyecto tiene **dos sistemas de assessments** que necesitan consolidarse:

| Sistema | Base de Datos | Ubicacion | Estado |
|---------|---------------|-----------|--------|
| **Nuevo** | PostgreSQL | `internal/domain/repositories/` | ACTIVO |
| **Legacy** | MongoDB | `internal/domain/repository/` | PARCIALMENTE ACTIVO |

### Uso Actual

**Sistema Nuevo (PostgreSQL):**
- `AssessmentRepository` - CRUD de assessments
- `AttemptRepository` - Intentos de evaluacion
- `AnswerRepository` - Respuestas de usuarios
- Usado por: `AssessmentAttemptService`

**Sistema Legacy (MongoDB):**
- `AssessmentStats` - **ACTIVO** (usado por `StatsService`)
- `AssessmentReader` - DEPRECATED (no se usa)
- `AssessmentWriter` - DEPRECATED (no se usa)

---

## Objetivo de la Consolidacion

1. **Migrar estadisticas a PostgreSQL** - Mover `CountCompletedAssessments` y `CalculateAverageScore` al nuevo sistema
2. **Eliminar codigo legacy** - Remover interfaces y repositorios MongoDB no usados
3. **Simplificar arquitectura** - Un solo sistema de assessments

---

## Plan de Migracion

### Fase 1: Preparacion (2-3 horas)

- [ ] Crear queries de estadisticas en PostgreSQL
- [ ] Agregar metodos a `PostgresAssessmentRepository`:
  - `CountCompletedAssessments(ctx) (int64, error)`
  - `CalculateAverageScore(ctx) (float64, error)`
- [ ] Crear tests de integracion para nuevos metodos

### Fase 2: Migracion de Datos (2-3 horas)

- [ ] Evaluar si hay datos historicos en MongoDB que migrar
- [ ] Crear script de migracion si es necesario
- [ ] Validar integridad de datos migrados

### Fase 3: Actualizacion de StatsService (2-3 horas)

- [ ] Modificar `StatsService` para usar PostgreSQL
- [ ] Actualizar inyeccion de dependencias en container
- [ ] Actualizar tests de `StatsService`

### Fase 4: Limpieza (2-3 horas)

- [ ] Eliminar `internal/domain/repository/assessment_repository.go`
- [ ] Eliminar `internal/infrastructure/persistence/mongodb/repository/assessment_repository_impl.go`
- [ ] Eliminar mocks legacy en `mock/mongodb/stubs.go`
- [ ] Eliminar `CreateLegacyAssessmentRepository` de factory
- [ ] Actualizar documentacion

---

## Queries de Estadisticas (PostgreSQL)

### CountCompletedAssessments

```sql
-- Contar intentos completados (score > 0 o status = 'completed')
SELECT COUNT(DISTINCT id)
FROM assessment_attempts
WHERE status = 'completed';
```

### CalculateAverageScore

```sql
-- Promedio de scores de intentos completados
SELECT COALESCE(AVG(score), 0.0)
FROM assessment_attempts
WHERE status = 'completed'
AND score IS NOT NULL;
```

---

## Dependencias

### Antes de iniciar:

1. Confirmar que no hay otros servicios usando el legacy repo
2. Verificar si hay datos historicos valiosos en MongoDB
3. Coordinar con equipo si hay consumidores externos

### Archivos afectados:

| Archivo | Accion |
|---------|--------|
| `internal/domain/repository/assessment_repository.go` | ELIMINAR |
| `internal/domain/repositories/assessment_repository.go` | MODIFICAR (agregar stats) |
| `internal/infrastructure/persistence/mongodb/repository/assessment_repository_impl.go` | ELIMINAR |
| `internal/infrastructure/persistence/postgres/repository/assessment_repository.go` | MODIFICAR |
| `internal/infrastructure/persistence/mock/mongodb/stubs.go` | MODIFICAR |
| `internal/container/factory.go` | MODIFICAR |
| `internal/container/repositories.go` | MODIFICAR |
| `internal/container/services.go` | MODIFICAR |
| `internal/application/service/stats_service.go` | MODIFICAR |

---

## Riesgos y Mitigaciones

| Riesgo | Probabilidad | Impacto | Mitigacion |
|--------|--------------|---------|------------|
| Perdida de datos historicos | Baja | Alto | Backup antes de migrar |
| Diferencias en calculos | Media | Medio | Tests comparativos |
| Downtime de estadisticas | Baja | Bajo | Feature flag para rollback |

---

## Criterios de Exito

- [ ] `StatsService` funciona con PostgreSQL
- [ ] Tests de integracion pasan
- [ ] No hay referencias al sistema legacy
- [ ] Estadisticas globales retornan mismos valores (o mejores)
- [ ] Codigo legacy eliminado

---

## Timeline Sugerido

| Semana | Actividad |
|--------|-----------|
| 1 | Fase 1: Preparacion |
| 2 | Fase 2: Migracion de datos |
| 3 | Fase 3: Actualizacion de StatsService |
| 4 | Fase 4: Limpieza y documentacion |

**Nota:** Este timeline asume dedicacion parcial. Puede comprimirse a 1-2 semanas con dedicacion completa.

---

## Referencias

- Sistema nuevo: `internal/domain/repositories/assessment_repository.go`
- Sistema legacy: `internal/domain/repository/assessment_repository.go`
- StatsService: `internal/application/service/stats_service.go`
- Container: `internal/container/services.go`

---

**Ultima actualizacion:** 2024-12-23
