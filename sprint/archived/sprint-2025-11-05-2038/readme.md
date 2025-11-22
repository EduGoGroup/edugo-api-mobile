# Sprint: Optimización de Queries - Índice en Materials

## Descripción

Implementar índice en la tabla `materials` de PostgreSQL para optimizar queries de ordenamiento por fecha de actualización. Esta es una tarea atómica pequeña pero funcional que mejora el performance de listados de materiales ordenados cronológicamente.

## Contexto

Esta tarea es parte de la Fase 3 del plan general (Implementar Queries Complejas), específicamente la subtarea 3.2. Se eligió como tarea de prueba para validar el sistema de comandos/agentes porque:

- ✅ Es **atómica** y autocontenida (1 archivo SQL)
- ✅ Es **funcional** (mejora performance real)
- ✅ No tiene **dependencias** de otras tareas
- ✅ Es **verificable** con EXPLAIN en PostgreSQL
- ✅ Es **segura** (no rompe funcionalidad existente)

## Objetivo

Crear script de migración SQL que agregue índice en la columna `updated_at` de la tabla `materials` para optimizar queries con ORDER BY updated_at DESC.

## Requisitos

### Requisito Funcional

- [ ] **RF-1**: Crear índice descendente en `materials.updated_at`
  - El índice debe ser descendente (DESC) porque las queries ordenan por más reciente primero
  - El índice debe mejorar performance de queries tipo: `SELECT * FROM materials ORDER BY updated_at DESC LIMIT 10`

### Requisitos Técnicos

- [ ] **RT-1**: Script SQL debe ser idempotente
  - Usar `CREATE INDEX IF NOT EXISTS` para evitar errores si ya existe
  - El script debe poder ejecutarse múltiples veces sin error

- [ ] **RT-2**: Seguir convención de nombres
  - Nombre del índice: `idx_materials_updated_at`
  - Nombre del archivo: `06_indexes_materials.sql`
  - Ubicación: `scripts/postgresql/`

- [ ] **RT-3**: Incluir comentarios en SQL
  - Explicar propósito del índice
  - Documentar queries que se benefician

### Requisitos de Validación

- [ ] **RV-1**: Verificar índice con EXPLAIN
  - Ejecutar EXPLAIN ANALYZE de query antes y después
  - Confirmar que el plan de ejecución usa el índice
  - Documentar mejora de performance (si es medible)

- [ ] **RV-2**: Proyecto compila sin errores
  - `go build ./...` debe pasar
  - No hay errores de sintaxis SQL

## Entregables Esperados

1. **Script SQL**: `scripts/postgresql/06_indexes_materials.sql`
   - Índice creado con IF NOT EXISTS
   - Comentarios explicativos
   - Sintaxis PostgreSQL válida

2. **Documentación** (opcional pero recomendado):
   - Resultado de EXPLAIN ANALYZE antes/después
   - Mejora de performance observada

3. **Commit atómico**:
   - Mensaje: `perf(db): agregar índice en materials.updated_at para optimizar ordenamiento`
   - Solo incluye el archivo SQL creado
   - Proyecto compila sin errores

## Ejemplo de Implementación Esperada

```sql
-- scripts/postgresql/06_indexes_materials.sql

-- Propósito: Optimizar queries que ordenan materiales por fecha de actualización
-- Beneficia queries tipo: SELECT * FROM materials ORDER BY updated_at DESC LIMIT N

-- Crear índice descendente en updated_at
-- DESC porque las queries más comunes ordenan de más reciente a más antiguo
CREATE INDEX IF NOT EXISTS idx_materials_updated_at
ON materials(updated_at DESC);

-- Verificar índice creado:
-- SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'materials';
```

## Queries que se Benefician

El índice optimizará estas queries comunes:

1. **Listar materiales recientes**:
   ```sql
   SELECT * FROM materials
   ORDER BY updated_at DESC
   LIMIT 20;
   ```

2. **Materiales actualizados en rango de fechas**:
   ```sql
   SELECT * FROM materials
   WHERE updated_at >= '2025-01-01'
   ORDER BY updated_at DESC;
   ```

## Restricciones/Consideraciones

### ✅ Hacer:
- Usar `IF NOT EXISTS` para idempotencia
- Crear índice como `DESC` (matches ORDER BY DESC)
- Ubicar en carpeta `scripts/postgresql/`
- Seguir numeración secuencial (06_)

### ❌ No Hacer:
- No modificar código Go (solo SQL)
- No modificar estructura de tabla (solo índice)
- No crear índices adicionales no solicitados
- No hacer commit si el proyecto no compila

### 🔍 Validación Manual

Para verificar que el índice funciona:

```bash
# 1. Ejecutar script SQL
psql -d edugo_db -f scripts/postgresql/06_indexes_materials.sql

# 2. Verificar índice creado
psql -d edugo_db -c "SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'materials';"

# 3. Ver plan de ejecución (debe usar idx_materials_updated_at)
psql -d edugo_db -c "EXPLAIN ANALYZE SELECT * FROM materials ORDER BY updated_at DESC LIMIT 10;"
```

## Criterios de Éxito

- [x] Script SQL creado en ubicación correcta
- [x] Índice usa IF NOT EXISTS (idempotente)
- [x] Índice es descendente (DESC)
- [x] Nombre de índice sigue convención: `idx_materials_updated_at`
- [x] Comentarios explican propósito
- [x] Proyecto compila sin errores
- [x] Commit atómico creado con mensaje apropiado

## Estimación

- **Complejidad**: Baja
- **Tiempo estimado**: 10-15 minutos
- **Archivos a crear**: 1 (script SQL)
- **Archivos a modificar**: 0

---

**Sprint para**: Validación de sistema de comandos/agentes  
**Tarea**: Fase 3, Subtarea 3.2  
**Fecha**: 2025-11-04  
**Branch**: `fix/debug-sprint-commands`
