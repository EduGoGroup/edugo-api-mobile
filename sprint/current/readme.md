# Sprint: Completar Queries Complejas - FASE 2.3

## Descripción

Completar la implementación de queries complejas en los servicios de la aplicación. Este sprint se enfoca en finalizar el PASO 2.3 del plan maestro, que incluye consultas avanzadas para materiales, evaluaciones, progreso y estadísticas.

## Contexto

Este sprint es la continuación de la FASE 2 (TODOs de Servicios). Ya se completaron:
- ✅ PASO 2.1: RabbitMQ Messaging (PR #15 merged)
- ✅ PASO 2.2: S3 URLs Firmadas (PR #16 merged)
- 🔵 PASO 2.3: Queries Complejas (20% completado - solo optimización de índice PostgreSQL)

Falta completar el 80% restante del PASO 2.3.

## Requisitos Funcionales

### RF-1: Queries de Materiales con Versiones
- [ ] Implementar consulta de materiales que incluya información de versiones
- [ ] Soportar filtrado por versión específica
- [ ] Optimizar consulta con joins eficientes

### RF-2: Cálculo de Puntajes en AssessmentService
- [ ] Implementar lógica de cálculo de puntajes basado en respuestas
- [ ] Soportar diferentes tipos de evaluación (multiple choice, verdadero/falso, etc.)
- [ ] Almacenar resultados en MongoDB

### RF-3: Generación de Feedback Detallado
- [ ] Generar feedback por pregunta en evaluaciones
- [ ] Incluir explicaciones de respuestas correctas/incorrectas
- [ ] Formatear feedback para consumo del frontend

### RF-4: Actualización de Progreso (UPSERT)
- [ ] Implementar UPSERT para actualización de progreso de usuario
- [ ] Evitar duplicados en la tabla de progreso
- [ ] Actualizar timestamp de última actualización

### RF-5: Query Complejo de Estadísticas
- [ ] Implementar query de estadísticas globales
- [ ] Incluir métricas de materiales, evaluaciones y progreso
- [ ] Optimizar con agregaciones eficientes

## Requisitos Técnicos

### RT-1: Seguir Clean Architecture
- Mantener separación de capas (domain, application, infrastructure)
- Usar DTOs para transferencia de datos
- Implementar interfaces en domain, implementaciones en infrastructure

### RT-2: Tests Unitarios
- Crear tests para cada método nuevo implementado
- Alcanzar mínimo 80% de cobertura en código nuevo
- Incluir casos edge (datos vacíos, valores nulos, etc.)

### RT-3: Performance
- Queries deben ejecutar en <100ms para datasets pequeños (<1000 registros)
- Usar índices apropiados en PostgreSQL
- Optimizar queries N+1 en MongoDB

### RT-4: Manejo de Errores
- Usar error types de `edugo-shared/common/errors`
- Logging apropiado con contexto
- Retornar errores de aplicación en handlers

## Entregables Esperados

### 1. Código Implementado

**Archivos a Modificar**:
- `internal/application/service/material_service.go`
- `internal/application/service/assessment_service.go`
- `internal/application/service/progress_service.go`
- `internal/application/service/stats_service.go`
- `internal/infrastructure/persistence/postgres/repository/material_repository_impl.go`
- `internal/infrastructure/persistence/postgres/repository/progress_repository_impl.go`
- `internal/infrastructure/persistence/mongodb/repository/assessment_repository_impl.go`

**Archivos de Tests**:
- Tests unitarios para cada servicio modificado
- Tests de repositorio con mocks

### 2. Documentación

- [ ] Comentarios en código explicando queries complejas
- [ ] Ejemplos de uso en comentarios
- [ ] Actualizar README si es necesario

### 3. Validación

- [ ] `go build ./...` pasa sin errores
- [ ] `go test ./...` todos los tests pasan
- [ ] Verificación manual de endpoints (opcional pero recomendado)

### 4. Commit Atómico

**Mensaje sugerido**:
```
feat: implementar consultas complejas en servicios

- Agregar queries de materiales con versiones
- Implementar cálculo de puntajes en AssessmentService
- Generar feedback detallado por pregunta
- Implementar UPSERT para actualización de progreso
- Agregar query de estadísticas globales

Incluye tests unitarios para todos los métodos nuevos.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## Restricciones/Consideraciones

### Base de Datos
- PostgreSQL 16 para datos estructurados (materiales, usuarios, progreso)
- MongoDB 7 para datos semi-estructurados (evaluaciones, respuestas)
- Ya existe índice en `materials.updated_at` (creado en tarea anterior)

### Dependencias
- Usar `edugo-shared` para error handling y logging
- RabbitMQ ya está configurado (PASO 2.1)
- S3 ya está configurado (PASO 2.2)

### Performance
- Evitar queries N+1
- Usar eager loading cuando sea apropiado
- Considerar paginación para queries grandes

### Testing
- Usar mocks para bases de datos en tests unitarios
- Testcontainers para tests de integración (opcional para este sprint)

## Criterios de Aceptación

- [x] ~~Optimización de índice PostgreSQL (materials.updated_at)~~ ✅ COMPLETADO
- [ ] Queries de materiales con versiones implementadas y testeadas
- [ ] Cálculo de puntajes funcionando correctamente
- [ ] Feedback detallado generándose para todas las evaluaciones
- [ ] UPSERT de progreso funcionando sin duplicados
- [ ] Query de estadísticas retornando métricas correctas
- [ ] Todos los tests pasando
- [ ] Código compilando sin errores
- [ ] Cobertura de tests ≥80% en código nuevo

## Estimación de Esfuerzo

**Total**: 1-1.5 días (~6-8 horas)

**Desglose**:
- Queries de materiales: 1-2 horas
- Cálculo de puntajes: 2-3 horas
- Feedback detallado: 1 hora
- UPSERT progreso: 1 hora
- Query estadísticas: 1-2 horas
- Tests y validación: 1 hora

## Referencias

- Plan Maestro: `sprint/docs/MASTER_PLAN_VISUAL.md` (FASE 2, PASO 2.3)
- Documentación anterior: `sprint/archived/sprint-2025-11-05-2038/`
- Código existente de servicios: `internal/application/service/`
- Repositorios: `internal/infrastructure/persistence/`

## Próximos Pasos Después de Este Sprint

Una vez completado este sprint (FASE 2.3), continuar con:
- **FASE 3**: Limpieza y Consolidación (eliminar código duplicado)
- **FASE 4**: Testing de Integración (tests con testcontainers)
