# Plan de Trabajo - Sprint: Optimización de Queries con Índice en Materials

## Resumen Ejecutivo

Este plan de trabajo implementa una optimización atómica de base de datos: la creación de un índice descendente en la columna `updated_at` de la tabla `materials` en PostgreSQL. La estrategia es quirúrgica y de bajo riesgo, diseñada para mejorar la performance de listados de materiales ordenados cronológicamente sin modificar código de aplicación.

**Objetivo**: Reducir latencia de queries con `ORDER BY updated_at DESC` de 50-200ms a 5-20ms (mejora de 5-10x).

**Alcance**: Solo capa de persistencia (PostgreSQL), sin cambios en código Go.

**Tiempo estimado**: 10-15 minutos de implementación + validación.

---

## Stack Tecnológico

- **Base de Datos**: PostgreSQL 14+
- **Feature**: Índices descendentes (`CREATE INDEX ... DESC`)
- **Herramientas**: psql CLI, EXPLAIN ANALYZE
- **Control de versiones**: Git
- **Driver Go**: lib/pq (sin modificaciones)

---

## 📋 Plan de Ejecución

### Fase 1: Preparación y Validación de Estado Actual

**Objetivo**: Establecer baseline de performance y preparar el ambiente local para la implementación del índice.

**Tareas**:

- [ ] **1.1** - Verificar conexión a base de datos local
  - **Descripción**: Confirmar que PostgreSQL está corriendo y accesible con las credenciales correctas
  - **Comando**: `psql -d edugo_db_local -c "SELECT current_database(), version();"`
  - **Criterio de aceptación**: Conexión exitosa, muestra nombre de BD y versión de PostgreSQL (debe ser 9.5+)

- [ ] **1.2** - Verificar existencia de tabla materials
  - **Descripción**: Confirmar que la tabla `materials` existe y contiene registros
  - **Comando**: `psql -d edugo_db_local -c "SELECT COUNT(*) as total_materials FROM materials;"`
  - **Criterio de aceptación**: Query ejecuta sin error y retorna cantidad de registros (puede ser 0 en BD local limpia)
  - 🔗 **Depende de**: Tarea 1.1

- [ ] **1.3** - Verificar índices existentes en tabla materials
  - **Descripción**: Listar todos los índices actuales de la tabla para confirmar que `idx_materials_updated_at` NO existe aún
  - **Comando**: `psql -d edugo_db_local -c "SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'materials';"`
  - **Criterio de aceptación**: Query retorna lista de índices (al menos PRIMARY KEY), `idx_materials_updated_at` NO está presente
  - 🔗 **Depende de**: Tarea 1.2

- [ ] **1.4** - Medir performance baseline (ANTES del índice)
  - **Descripción**: Ejecutar EXPLAIN ANALYZE de query de listado para documentar performance actual
  - **Comando**:
    ```sql
    psql -d edugo_db_local -c "EXPLAIN ANALYZE SELECT * FROM materials ORDER BY updated_at DESC LIMIT 20;"
    ```
  - **Criterio de aceptación**:
    - Query ejecuta exitosamente
    - Plan de ejecución muestra "Seq Scan" o "Sort" (no usa índice de updated_at)
    - Tiempo de ejecución documentado (puede ser muy rápido en BD vacía, eso es esperado)
  - 🔗 **Depende de**: Tarea 1.3
  - **Nota**: Si la BD local está vacía o con pocos registros, el tiempo será bajo. Esto es aceptable; la validación real será en QA/producción.

**Completitud de Fase**: 0/4 tareas completadas

---

### Fase 2: Creación del Script de Migración

**Objetivo**: Crear el archivo SQL de migración con contenido idempotente, documentado y siguiendo convenciones del proyecto.

**Tareas**:

- [ ] **2.1** - Verificar carpeta de scripts SQL
  - **Descripción**: Confirmar que existe la carpeta `scripts/postgresql/` en el proyecto
  - **Comando**: `ls -la scripts/postgresql/`
  - **Criterio de aceptación**: Carpeta existe, muestra scripts existentes numerados (`01_*.sql`, `02_*.sql`, etc.)

- [ ] **2.2** - Identificar número secuencial para el nuevo script
  - **Descripción**: Determinar el próximo número disponible para nombrar el script (debe ser `06_indexes_materials.sql`)
  - **Comando**: `ls scripts/postgresql/ | grep -E '^[0-9]+_' | sort -V | tail -1`
  - **Criterio de aceptación**: Último script identificado (ej: `05_*.sql`), próximo número es `06`
  - 🔗 **Depende de**: Tarea 2.1

- [ ] **2.3** - Crear archivo de script SQL con contenido completo
  - **Descripción**: Crear el archivo `scripts/postgresql/06_indexes_materials.sql` con el comando de creación de índice idempotente y documentado
  - **Archivos a crear**: `scripts/postgresql/06_indexes_materials.sql`
  - **Contenido del archivo**:
    ```sql
    -- ============================================================
    -- Migration: 06_indexes_materials.sql
    -- Description: Agregar índice descendente en materials.updated_at
    --              para optimizar queries de listado cronológico
    -- Author: Claude Code / EduGo Team
    -- Date: 2025-11-04
    -- ============================================================

    -- Objetivo:
    -- Mejorar performance de queries que ordenan materiales por fecha
    -- de actualización más reciente (patrón común en la aplicación).
    --
    -- Queries beneficiadas:
    -- 1. SELECT * FROM materials ORDER BY updated_at DESC LIMIT N;
    -- 2. SELECT * FROM materials WHERE course_id = X ORDER BY updated_at DESC;
    -- 3. SELECT * FROM materials WHERE type = 'Y' ORDER BY updated_at DESC;
    --
    -- Mejora esperada: 5-10x más rápido (de 50-200ms a 5-20ms)

    -- Crear índice descendente de forma idempotente
    CREATE INDEX IF NOT EXISTS idx_materials_updated_at
    ON materials(updated_at DESC);

    -- Verificación:
    -- Después de ejecutar este script, verificar con:
    -- SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'materials';
    --
    -- Validar uso del índice con:
    -- EXPLAIN ANALYZE SELECT * FROM materials ORDER BY updated_at DESC LIMIT 20;
    -- Debe mostrar: "Index Scan using idx_materials_updated_at"

    -- Rollback (si es necesario):
    -- DROP INDEX IF EXISTS idx_materials_updated_at;
    ```
  - **Criterio de aceptación**:
    - Archivo creado en ruta correcta
    - Contiene comentarios explicativos completos
    - Usa `CREATE INDEX IF NOT EXISTS` (idempotente)
    - Especifica dirección `DESC` en el índice
    - Incluye instrucciones de verificación y rollback
  - 🔗 **Depende de**: Tarea 2.2

- [ ] **2.4** - Validar sintaxis SQL del script
  - **Descripción**: Verificar que el script SQL no tiene errores de sintaxis antes de ejecutarlo
  - **Comando**: `psql -d edugo_db_local --dry-run -f scripts/postgresql/06_indexes_materials.sql` (o validar con linter SQL si está disponible)
  - **Alternativa**: Ejecutar en transacción y hacer rollback:
    ```bash
    psql -d edugo_db_local -c "BEGIN; \i scripts/postgresql/06_indexes_materials.sql; ROLLBACK;"
    ```
  - **Criterio de aceptación**: No hay errores de sintaxis, comando es válido
  - 🔗 **Depende de**: Tarea 2.3
  - **Nota**: Si `--dry-run` no está disponible en tu versión de psql, usar la alternativa con BEGIN/ROLLBACK.

**Completitud de Fase**: 0/4 tareas completadas

---

### Fase 3: Ejecución Local del Script

**Objetivo**: Aplicar el script de migración en la base de datos local y verificar que el índice se crea correctamente.

**Tareas**:

- [ ] **3.1** - Ejecutar script de migración en BD local
  - **Descripción**: Aplicar el script SQL para crear el índice en la base de datos local
  - **Comando**: `psql -d edugo_db_local -f scripts/postgresql/06_indexes_materials.sql`
  - **Criterio de aceptación**:
    - Script ejecuta sin errores
    - Mensaje de salida: `CREATE INDEX` o `NOTICE: relation "idx_materials_updated_at" already exists, skipping`
    - No hay mensajes de ERROR
  - 🔗 **Depende de**: Fase 2 completada

- [ ] **3.2** - Verificar creación del índice
  - **Descripción**: Confirmar que el índice `idx_materials_updated_at` existe en el catálogo de PostgreSQL
  - **Comando**:
    ```bash
    psql -d edugo_db_local -c "SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'materials' AND indexname = 'idx_materials_updated_at';"
    ```
  - **Criterio de aceptación**:
    - Query retorna exactamente 1 registro
    - `indexname` = `idx_materials_updated_at`
    - `indexdef` contiene `CREATE INDEX idx_materials_updated_at ON materials USING btree (updated_at DESC)`
  - 🔗 **Depende de**: Tarea 3.1

- [ ] **3.3** - Validar que el índice es utilizado en queries
  - **Descripción**: Ejecutar EXPLAIN ANALYZE de query de prueba para confirmar que PostgreSQL usa el nuevo índice
  - **Comando**:
    ```bash
    psql -d edugo_db_local -c "EXPLAIN ANALYZE SELECT * FROM materials ORDER BY updated_at DESC LIMIT 20;"
    ```
  - **Criterio de aceptación**:
    - Plan de ejecución muestra `Index Scan using idx_materials_updated_at` o `Index Scan Backward using idx_materials_updated_at`
    - NO muestra "Seq Scan" ni "Sort" en el plan principal
    - Tiempo de ejecución documentado (comparar con baseline de Tarea 1.4)
  - 🔗 **Depende de**: Tarea 3.2
  - **Nota**: Si la BD local tiene muy pocos registros (<100), PostgreSQL puede optar por Seq Scan (es correcto según optimizador). En ese caso, documentar y validar en ambiente QA con más datos.

- [ ] **3.4** - Probar idempotencia del script
  - **Descripción**: Re-ejecutar el script para confirmar que es idempotente (no falla si el índice ya existe)
  - **Comando**: `psql -d edugo_db_local -f scripts/postgresql/06_indexes_materials.sql`
  - **Criterio de aceptación**:
    - Script ejecuta exitosamente
    - Mensaje: `NOTICE: relation "idx_materials_updated_at" already exists, skipping`
    - No hay ERROR
    - Índice sigue existiendo (no se duplica)
  - 🔗 **Depende de**: Tarea 3.3

**Completitud de Fase**: 0/4 tareas completadas

---

### Fase 4: Validación de Impacto en Aplicación

**Objetivo**: Confirmar que la aplicación Go sigue funcionando correctamente y que el índice no introduce regresiones.

**Tareas**:

- [ ] **4.1** - Verificar que la aplicación compila
  - **Descripción**: Ejecutar build del proyecto Go para asegurar que no hay errores de compilación (no debería haber, ya que no se modificó código)
  - **Comando**: `go build ./...`
  - **Criterio de aceptación**: Build exitoso sin errores, binario generado (si aplica)

- [ ] **4.2** - Ejecutar suite de tests unitarios
  - **Descripción**: Correr todos los tests unitarios del proyecto para confirmar que no hay regresiones
  - **Comando**: `go test ./... -v`
  - **Criterio de aceptación**:
    - Todos los tests pasan (resultado: PASS)
    - No hay tests fallidos relacionados con materiales o queries
  - 🔗 **Depende de**: Tarea 4.1
  - **Nota**: El índice es transparente para el código Go, por lo que los tests no deberían verse afectados.

- [ ] **4.3** - Ejecutar tests de integración (si existen)
  - **Descripción**: Si el proyecto tiene tests de integración que usan base de datos real, ejecutarlos para validar comportamiento end-to-end
  - **Comando**: `go test ./... -tags=integration -v` (ajustar según convención del proyecto)
  - **Criterio de aceptación**:
    - Tests de integración pasan exitosamente
    - Queries de listado de materiales se ejecutan sin errores
    - Performance igual o mejor que antes
  - 🔗 **Depende de**: Tarea 4.2
  - **Nota**: Si no hay tests de integración, marcar como "N/A - No aplica" y continuar.

- [ ] **4.4** - Probar manualmente endpoint de listado de materiales (opcional)
  - **Descripción**: Si es posible levantar el servidor localmente, hacer request manual al endpoint de materiales
  - **Comando (opcional)**:
    ```bash
    # Levantar servidor (ajustar según proyecto)
    go run cmd/main.go

    # En otra terminal, hacer request
    curl -X GET "http://localhost:8080/api/materials?sort=updated_at&order=desc&limit=20"
    ```
  - **Criterio de aceptación**:
    - Endpoint responde exitosamente (HTTP 200)
    - Retorna JSON con lista de materiales
    - No hay errores en logs del servidor
  - 🔗 **Depende de**: Tarea 4.3
  - **Nota**: Esta tarea es opcional. Si no es posible levantar servidor localmente, omitir y confiar en tests automatizados.

**Completitud de Fase**: 0/4 tareas completadas

---

### Fase 5: Control de Versiones y Documentación

**Objetivo**: Registrar el cambio en Git con un commit atómico bien documentado y actualizar el plan de sprint.

**Tareas**:

- [ ] **5.1** - Verificar estado de Git antes del commit
  - **Descripción**: Confirmar qué archivos han sido modificados/agregados antes de hacer commit
  - **Comando**: `git status`
  - **Criterio de aceptación**:
    - Muestra `scripts/postgresql/06_indexes_materials.sql` como archivo nuevo (untracked o en staging)
    - No hay otros archivos modificados no relacionados con este sprint
    - Branch actual es correcto (ej: `fix/debug-sprint-commands` o branch de trabajo)

- [ ] **5.2** - Agregar script SQL al staging area
  - **Descripción**: Agregar el archivo de script al área de staging de Git
  - **Comando**: `git add scripts/postgresql/06_indexes_materials.sql`
  - **Criterio de aceptación**:
    - `git status` muestra el archivo en "Changes to be committed"
    - Solo el archivo del script está en staging (no hay archivos adicionales no intencionados)
  - 🔗 **Depende de**: Tarea 5.1

- [ ] **5.3** - Crear commit con mensaje descriptivo
  - **Descripción**: Hacer commit del script con mensaje siguiendo convención del proyecto y footer de Claude Code
  - **Comando**:
    ```bash
    git commit -m "$(cat <<'EOF'
    perf(db): agregar índice en materials.updated_at para optimizar ordenamiento

    - Crear script 06_indexes_materials.sql
    - Índice descendente (DESC) para queries con ORDER BY updated_at DESC
    - Script idempotente con IF NOT EXISTS
    - Mejora esperada: 5-10x más rápido (50-200ms → 5-20ms)
    - Sin cambios en código Go (optimización transparente)

    Queries beneficiadas:
    - Listado de materiales recientes
    - Filtros por curso/tipo + ordenamiento cronológico

    Validado con EXPLAIN ANALYZE en ambiente local.

    🤖 Generated with [Claude Code](https://claude.com/claude-code)

    Co-Authored-By: Claude <noreply@anthropic.com>
    EOF
    )"
    ```
  - **Criterio de aceptación**:
    - Commit creado exitosamente
    - Mensaje de commit incluye:
      - Prefijo `perf(db):` (tipo de cambio)
      - Descripción concisa en primera línea
      - Bullet points de detalles en cuerpo
      - Footer de Claude Code
    - `git log -1` muestra el commit recién creado
  - 🔗 **Depende de**: Tarea 5.2

- [ ] **5.4** - Actualizar plan de sprint con estado completado
  - **Descripción**: Marcar todas las casillas de este plan como completadas y documentar resultado
  - **Archivos a modificar**: `sprint/current/planning/readme.md` y `sprint/current/readme.md`
  - **Cambios**:
    - Marcar todas las tareas de este plan con `[x]` en lugar de `[ ]`
    - Actualizar sección de "Completitud de Fase" con conteos correctos
    - Agregar nota al final del plan con resultado de validación (ej: "✅ Índice creado exitosamente. Validado con EXPLAIN ANALYZE. Tiempo de ejecución reducido.")
  - **Criterio de aceptación**:
    - Archivo `sprint/current/planning/readme.md` actualizado con checkboxes marcados
    - Archivo `sprint/current/readme.md` actualizado con progreso del sprint
    - Documentación refleja estado real del trabajo
  - 🔗 **Depende de**: Tarea 5.3

- [ ] **5.5** - Crear commit de actualización de documentación
  - **Descripción**: Hacer commit de los archivos de plan actualizados
  - **Comando**:
    ```bash
    git add sprint/current/planning/readme.md sprint/current/readme.md
    git commit -m "docs(sprint): marcar optimización de índice como completada

    - Actualizar checkboxes en planning/readme.md
    - Documentar resultado de validación
    - Sprint completado exitosamente

    🤖 Generated with [Claude Code](https://claude.com/claude-code)

    Co-Authored-By: Claude <noreply@anthropic.com>"
    ```
  - **Criterio de aceptación**:
    - Commit de documentación creado
    - `git log` muestra 2 commits nuevos (script + docs)
  - 🔗 **Depende de**: Tarea 5.4

**Completitud de Fase**: 0/5 tareas completadas

---

### Fase 6: Preparación para Deployment (Opcional)

**Objetivo**: Documentar pasos necesarios para ejecutar el script en ambientes QA y producción (no se ejecuta en este sprint, solo se prepara documentación).

**Tareas**:

- [ ] **6.1** - Documentar instrucciones de deployment para QA
  - **Descripción**: Crear/actualizar documento con pasos para que DevOps ejecute el script en ambiente QA
  - **Archivos a crear/modificar**: `docs/deployment/database-migrations.md` o similar
  - **Contenido mínimo**:
    - Comando para ejecutar script en QA
    - Comando para verificar índice creado
    - Query EXPLAIN ANALYZE para validar uso del índice
    - Rollback en caso de problemas
  - **Criterio de aceptación**: Documento existe con instrucciones claras y completas

- [ ] **6.2** - Documentar consideraciones para deployment en producción
  - **Descripción**: Agregar sección en documentación sobre precauciones para producción
  - **Archivos a modificar**: `docs/deployment/database-migrations.md` o similar
  - **Contenido mínimo**:
    - Verificar espacio en disco antes de crear índice
    - Ejecutar en ventana de bajo tráfico (ej: 2 AM)
    - Tiempo estimado de creación según tamaño de tabla
    - Plan de monitoreo post-deployment (latencia, uso de índice)
    - Comando de rollback en caso de problemas
  - **Criterio de aceptación**: Documentación completa, DevOps puede seguir pasos sin ambigüedad
  - 🔗 **Depende de**: Tarea 6.1

- [ ] **6.3** - Notificar al equipo sobre cambio pendiente (opcional)
  - **Descripción**: Informar a DevOps/QA que hay un script de migración listo para deployment
  - **Acción**: Enviar mensaje en canal de comunicación del equipo (Slack, email, etc.)
  - **Contenido sugerido**:
    - "Script de optimización de BD listo: `06_indexes_materials.sql`"
    - "Mejora esperada: 5-10x en queries de listado de materiales"
    - "Próximo paso: ejecutar en QA para validación"
    - Link al commit o PR
  - **Criterio de aceptación**: Equipo notificado, DevOps está al tanto del cambio pendiente
  - 🔗 **Depende de**: Tarea 6.2
  - **Nota**: Esta tarea es opcional y depende de la cultura del equipo. Si no aplica, omitir.

**Completitud de Fase**: 0/3 tareas completadas

**Nota**: Esta fase es opcional y puede ejecutarse después del sprint principal. Las tareas 5.3 ya incluye el commit necesario para que el cambio esté listo para deployment.

---

## 📊 Resumen de Dependencias

### Ruta Crítica

El sprint sigue una ruta secuencial obligatoria (no hay tareas independientes que puedan paralelizarse):

```
Fase 1 (Preparación) → Fase 2 (Script) → Fase 3 (Ejecución) → Fase 4 (Validación) → Fase 5 (Commit)
```

**Secuencia crítica de tareas**:
1. Tarea 1.1 → 1.2 → 1.3 → 1.4 (establecer baseline)
2. Tarea 2.1 → 2.2 → 2.3 → 2.4 (crear script validado)
3. Tarea 3.1 → 3.2 → 3.3 → 3.4 (aplicar índice)
4. Tarea 4.1 → 4.2 → 4.3 → 4.4 (validar app)
5. Tarea 5.1 → 5.2 → 5.3 → 5.4 → 5.5 (commit cambios)
6. Fase 6 es opcional y puede ejecutarse después

### Dependencias Entre Fases

- **Fase 2 depende de Fase 1**: Necesitamos verificar estado actual antes de crear script
- **Fase 3 depende de Fase 2**: No podemos ejecutar script hasta que esté creado y validado
- **Fase 4 depende de Fase 3**: Validación de app requiere que índice esté creado
- **Fase 5 depende de Fase 4**: Solo hacer commit si todas las validaciones pasan
- **Fase 6 es independiente**: Puede hacerse en paralelo o después de Fase 5

### Tareas que Pueden Omitirse sin Bloquear el Sprint

- **Tarea 1.4** (Medir baseline): Recomendada pero no bloqueante. Si la BD local está vacía, el resultado no será representativo.
- **Tarea 4.3** (Tests de integración): Solo si el proyecto tiene tests de integración configurados.
- **Tarea 4.4** (Prueba manual): Opcional, útil para validación adicional pero no bloqueante.
- **Fase 6 completa** (Documentación de deployment): Puede hacerse después del sprint principal.

---

## 📈 Métricas del Plan

- **Total de fases**: 6 (5 obligatorias + 1 opcional)
- **Total de tareas**: 24 tareas (19 obligatorias + 5 opcionales)
- **Tareas con dependencias explícitas**: 18 tareas
- **Tareas independientes (inicio de fase)**: 6 tareas (1.1, 2.1, 4.1, 5.1, 6.1)
- **Estimación de tiempo**:
  - Fase 1: 3-5 minutos
  - Fase 2: 3-5 minutos
  - Fase 3: 2-3 minutos
  - Fase 4: 5-10 minutos (depende de suite de tests)
  - Fase 5: 3-5 minutos
  - Fase 6: 10-15 minutos (opcional)
  - **Total**: 16-28 minutos (sin Fase 6), 26-43 minutos (con Fase 6)

---

## 🎯 Estrategia de Ejecución Recomendada

### Enfoque: Ejecución Lineal Secuencial

Dado que este es un sprint atómico y las tareas tienen dependencias secuenciales fuertes, la estrategia recomendada es:

1. **Primera sesión (10-15 min)**: Ejecutar Fases 1-3
   - Preparar ambiente
   - Crear script
   - Aplicar índice localmente
   - **Hito**: Índice creado y verificado

2. **Segunda sesión (5-10 min)**: Ejecutar Fases 4-5
   - Validar aplicación
   - Crear commits
   - **Hito**: Cambio registrado en Git, listo para push/PR

3. **Tercera sesión (opcional, 10-15 min)**: Ejecutar Fase 6
   - Documentar deployment
   - Notificar equipo
   - **Hito**: Cambio listo para deployment a QA/producción

### Uso del Comando `/03-execution`

Para ejecutar fases específicas:

```bash
# Ejecutar todas las fases
/03-execution

# Ejecutar solo Fase 1 (preparación)
/03-execution phase-1

# Ejecutar solo Fase 3 (ejecución del script)
/03-execution phase-3

# Ejecutar tarea específica (si el comando lo soporta)
/03-execution task-3.2
```

### Puntos de Verificación (Checkpoints)

**Checkpoint 1**: Al completar Fase 1
- ✅ Baseline de performance documentado
- ✅ Índice NO existe aún
- **Decisión**: Proceder a Fase 2

**Checkpoint 2**: Al completar Fase 3
- ✅ Índice creado exitosamente
- ✅ PostgreSQL usa el índice en queries
- ✅ Script es idempotente
- **Decisión**: Proceder a validación de app (Fase 4)

**Checkpoint 3**: Al completar Fase 4
- ✅ Aplicación compila sin errores
- ✅ Tests pasan exitosamente
- **Decisión**: Proceder a commit (Fase 5)

**Checkpoint Final**: Al completar Fase 5
- ✅ Commit(s) creado(s) con mensaje apropiado
- ✅ Documentación de sprint actualizada
- **Decisión**: Sprint completado. Opcional: documentar deployment (Fase 6)

---

## 🚨 Manejo de Errores por Fase

### Fase 1: Problemas de Conexión

**Error**: No se puede conectar a PostgreSQL local
- **Causa posible**: PostgreSQL no está corriendo, credenciales incorrectas, puerto bloqueado
- **Solución**:
  1. Verificar que PostgreSQL está corriendo: `pg_ctl status` o `brew services list` (macOS)
  2. Iniciar PostgreSQL si está detenido: `brew services start postgresql` (macOS) o `sudo systemctl start postgresql` (Linux)
  3. Verificar credenciales en `config/config-local.yaml` o variables de entorno
  4. Verificar puerto: por defecto 5432
- **Impacto**: Bloquea todo el sprint hasta resolver

**Error**: Tabla `materials` no existe
- **Causa posible**: Migraciones anteriores no se han ejecutado en BD local
- **Solución**:
  1. Ejecutar scripts de migración previos: `01_*.sql`, `02_*.sql`, etc.
  2. Verificar que tabla se creó: `\dt materials` en psql
  3. Si es necesario, ejecutar seed de datos de prueba
- **Impacto**: Bloquea Fase 2 y siguientes

### Fase 2: Problemas con Script SQL

**Error**: Sintaxis SQL incorrecta en el script
- **Causa posible**: Error tipográfico, palabra clave mal escrita, sintaxis no válida
- **Solución**:
  1. Revisar script carácter por carácter
  2. Comparar con ejemplo en este plan (Tarea 2.3)
  3. Ejecutar en transacción de prueba: `BEGIN; \i script.sql; ROLLBACK;`
  4. Buscar documentación de PostgreSQL para sintaxis de CREATE INDEX
- **Impacto**: Bloquea Fase 3 hasta corregir

### Fase 3: Problemas al Crear Índice

**Error**: `ERROR: could not create unique index "idx_materials_updated_at"`
- **Causa posible**: Valores duplicados o NULL en `updated_at` (poco probable, pero posible)
- **Solución**:
  1. Verificar si hay valores NULL: `SELECT COUNT(*) FROM materials WHERE updated_at IS NULL;`
  2. Si hay NULLs, actualizar: `UPDATE materials SET updated_at = created_at WHERE updated_at IS NULL;`
  3. Re-ejecutar script
- **Impacto**: Bloquea Fase 4 hasta resolver

**Error**: `ERROR: out of memory` durante creación de índice
- **Causa posible**: Tabla demasiado grande, memoria insuficiente
- **Solución**:
  1. Cerrar aplicaciones que consumen memoria
  2. Aumentar `maintenance_work_mem` en PostgreSQL temporalmente
  3. Considerar `CREATE INDEX CONCURRENTLY` (tarda más pero consume menos memoria)
- **Impacto**: Bloquea sprint hasta resolver problema de recursos

**Error**: Índice no se usa en EXPLAIN ANALYZE (Tarea 3.3)
- **Causa posible**: Tabla muy pequeña, estadísticas desactualizadas, configuración de PostgreSQL
- **Solución**:
  1. Ejecutar `ANALYZE materials;` para actualizar estadísticas
  2. Revisar configuración de `random_page_cost` (debería ser ~1.1 para SSD)
  3. Si tabla tiene <100 registros, es esperado que no use índice (Seq Scan es más rápido)
  4. **Decisión**: Si tabla es pequeña, documentar y continuar. El índice se usará al crecer la tabla.
- **Impacto**: No bloquea sprint (es comportamiento aceptable del optimizador)

### Fase 4: Problemas con Tests

**Error**: Tests fallan después de crear índice
- **Causa posible**: Tests asumen un orden específico sin `ORDER BY`, comportamiento de BD cambió
- **Solución**:
  1. Revisar tests fallidos para entender qué esperaban
  2. Verificar que tests no dependan de orden implícito de resultados
  3. Si el problema es real del código (poco probable), investigar
  4. Si el problema es de los tests (asumen orden no determinístico), corregir tests
- **Impacto**: Puede bloquear commit si los tests son parte de CI/CD obligatorio

**Error**: Aplicación no compila
- **Causa posible**: Problema no relacionado con el índice (índice es transparente para Go)
- **Solución**:
  1. Verificar que no se modificó código Go accidentalmente
  2. Ejecutar `go mod tidy` por si hay problema de dependencias
  3. Revisar errores de compilación
  4. **Si el error existía antes del sprint**: documentar y notificar al usuario
- **Impacto**: Bloquea commit hasta resolver

### Fase 5: Problemas con Git

**Error**: `git commit` falla por hooks de pre-commit
- **Causa posible**: Linter detecta problema, formato incorrecto, tests fallan en hook
- **Solución**:
  1. Revisar salida del hook para entender qué falló
  2. Corregir el problema reportado
  3. Re-intentar commit
  4. Si el hook es problemático: discutir con usuario antes de bypasear con `--no-verify`
- **Impacto**: Bloquea finalización del sprint hasta resolver

---

## 📝 Notas Adicionales

### Validación de Éxito del Sprint

Al completar todas las tareas obligatorias (Fases 1-5), el sprint es exitoso si:

✅ Archivo `scripts/postgresql/06_indexes_materials.sql` existe y está bajo control de versiones
✅ Índice `idx_materials_updated_at` creado en BD local y verificado con query de catálogo
✅ EXPLAIN ANALYZE muestra que el índice se usa (o se documenta por qué no se usa si tabla es pequeña)
✅ Script es idempotente (puede ejecutarse múltiples veces sin error)
✅ Aplicación compila sin errores
✅ Tests pasan exitosamente (unitarios + integración si aplica)
✅ Commit creado con mensaje apropiado y footer de Claude Code
✅ Documentación de sprint actualizada con estado completado

### Consideraciones para BD Local Vacía

Si la base de datos local está vacía o tiene muy pocos registros (<100):

- **Tarea 1.4** (Baseline): El tiempo será muy bajo (1-5ms). Documentar que la BD está vacía.
- **Tarea 3.3** (Validar uso del índice): PostgreSQL puede elegir Seq Scan en lugar del índice (es correcto según optimizador).
- **Solución**: Documentar en el commit que la validación de performance real se hará en ambiente QA con datos reales.
- **Alternativa**: Si es posible, hacer seed de datos de prueba (1000-10000 registros) para validar uso del índice localmente.

### Próximos Pasos Después del Sprint

Una vez completado este sprint:

1. **Push del branch**: Ejecutar `git push origin fix/debug-sprint-commands` (o branch correspondiente)
2. **Crear Pull Request**: Seguir proceso estándar del proyecto para PR
3. **Deployment a QA**: DevOps ejecuta script en ambiente QA según documentación de Fase 6
4. **Validación en QA**: Confirmar mejora de performance con datos reales
5. **Deployment a Producción**: Ejecutar script en producción en ventana de bajo tráfico
6. **Monitoreo Post-Deployment**: Observar métricas de latencia por 24-48 horas
7. **Actualizar Sprint**: Archivar sprint actual y preparar próximo sprint

### Rollback en Caso de Problemas

Si después de crear el índice se detecta un problema crítico (muy poco probable):

```sql
-- Rollback: remover índice
DROP INDEX IF EXISTS idx_materials_updated_at;

-- Verificar que se eliminó
SELECT indexname FROM pg_indexes WHERE tablename = 'materials';
```

El rollback es instantáneo y seguro. La aplicación seguirá funcionando (con performance degradada, pero funcionando).

### Recursos Útiles

- **Documentación de PostgreSQL sobre índices**: https://www.postgresql.org/docs/current/indexes.html
- **EXPLAIN ANALYZE tutorial**: https://www.postgresql.org/docs/current/using-explain.html
- **Guía de optimización de queries**: https://www.postgresql.org/docs/current/performance-tips.html

---

## ✅ Criterios de Completitud del Sprint

Antes de marcar el sprint como completado, verificar:

- [ ] Todas las tareas obligatorias de Fases 1-5 están marcadas como completadas `[x]`
- [ ] Script SQL existe en `scripts/postgresql/06_indexes_materials.sql`
- [ ] Índice creado y verificado con `SELECT * FROM pg_indexes WHERE tablename = 'materials';`
- [ ] EXPLAIN ANALYZE ejecutado y documentado (resultado aceptable)
- [ ] Aplicación compila: `go build ./...` exitoso
- [ ] Tests pasan: `go test ./...` exitoso
- [ ] Al menos 1 commit creado con mensaje apropiado y footer de Claude Code
- [ ] Documentación de sprint actualizada en `sprint/current/readme.md`
- [ ] No hay errores pendientes ni tareas bloqueadas

**Estado al completar**: Sprint exitoso, optimización implementada, listo para deployment a QA.

---

**Plan generado el**: 2025-11-04
**Responsable**: Claude Code (Agente de Planificación)
**Basado en**: `sprint/current/analysis/readme.md`
**Próximo comando**: `/03-execution` para comenzar ejecución del plan

---

💡 **Tip**: Este plan fue diseñado para ser ejecutable tanto manualmente (siguiendo tarea por tarea) como automáticamente (usando `/03-execution`). Cada tarea tiene criterios de aceptación claros para facilitar validación.
