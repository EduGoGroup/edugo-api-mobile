# Validación del Sprint 01 - Schema de Base de Datos

## Pre-validación

### Verificar Estado del Proyecto
```bash
# Cambiar al directorio del proyecto
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Verificar estado de Git
git status

# Verificar rama actual
git branch --show-current
# Output esperado: develop o feature/assessments

# Verificar que PostgreSQL está corriendo
pg_isready -h localhost -p 5432
# Output esperado: localhost:5432 - accepting connections
```

---

## Checklist de Validación

### 1. Migración Principal (06_assessments.sql)

#### 1.1 Verificar Sintaxis SQL
```bash
# Verificar sintaxis sin ejecutar (dry-run simulado)
psql -U postgres -d postgres --single-transaction --set ON_ERROR_STOP=on -f /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments.sql -o /dev/null
```
**Criterio de éxito:** Comando retorna exit code 0 (sin errores de sintaxis)

#### 1.2 Ejecutar Migración en BD de Test
```bash
# Crear BD de test limpia
psql -U postgres -c "DROP DATABASE IF EXISTS edugo_test_sprint01;"
psql -U postgres -c "CREATE DATABASE edugo_test_sprint01;"

# Ejecutar migraciones previas (si existen)
psql -U postgres -d edugo_test_sprint01 < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/01_base_tables.sql

# Ejecutar migración de assessments
psql -U postgres -d edugo_test_sprint01 < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments.sql
```
**Criterio de éxito:** Migración completa sin errores, output muestra "COMMIT"

#### 1.3 Verificar Tablas Creadas
```bash
# Listar todas las tablas de assessments
psql -U postgres -d edugo_test_sprint01 -c "\dt assessment*"
```
**Criterio de éxito:** 
- ✅ 4 tablas creadas: `assessment`, `assessment_attempt`, `assessment_attempt_answer`, `material_summary_link`

```bash
# Verificar estructura de tabla assessment
psql -U postgres -d edugo_test_sprint01 -c "\d assessment"
```
**Criterio de éxito:**
- ✅ 10 columnas: id, material_id, mongo_document_id, title, total_questions, pass_threshold, max_attempts, time_limit_minutes, created_at, updated_at
- ✅ PK en `id`
- ✅ FK a `materials(id)`
- ✅ Unique constraint en `material_id`

```bash
# Verificar estructura de tabla assessment_attempt
psql -U postgres -d edugo_test_sprint01 -c "\d assessment_attempt"
```
**Criterio de éxito:**
- ✅ 9 columnas: id, assessment_id, student_id, score, max_score, time_spent_seconds, started_at, completed_at, created_at
- ✅ Check constraint `score >= 0 AND score <= 100`
- ✅ Check constraint `completed_at > started_at`
- ✅ Check constraint validando `time_spent_seconds`

#### 1.4 Verificar Constraints
```bash
# Listar todos los constraints de assessment
psql -U postgres -d edugo_test_sprint01 -c "
    SELECT conname, contype, confdeltype
    FROM pg_constraint 
    WHERE conrelid = 'assessment'::regclass
    ORDER BY conname;
"
```
**Criterio de éxito:**
- ✅ FK constraint: `fk_assessment_material` con `ON DELETE CASCADE`
- ✅ Unique constraint: `unique_material_assessment`
- ✅ Check constraints en `total_questions` (1-100)
- ✅ Check constraint en `pass_threshold` (0-100)

```bash
# Listar constraints de assessment_attempt
psql -U postgres -d edugo_test_sprint01 -c "
    SELECT conname, contype 
    FROM pg_constraint 
    WHERE conrelid = 'assessment_attempt'::regclass
    ORDER BY conname;
"
```
**Criterio de éxito:**
- ✅ 5+ constraints: FK a assessment, FK a users, checks de validación

---

### 2. Índices Optimizados

#### 2.1 Verificar Cantidad de Índices
```bash
# Contar índices por tabla
psql -U postgres -d edugo_test_sprint01 -c "
    SELECT tablename, COUNT(*) as num_indexes 
    FROM pg_indexes 
    WHERE tablename LIKE 'assessment%' 
    GROUP BY tablename
    ORDER BY tablename;
"
```
**Criterio de éxito:**
- ✅ `assessment`: Mínimo 4 índices (PK + 3 secundarios)
- ✅ `assessment_attempt`: Mínimo 6 índices (PK + 5 secundarios)
- ✅ `assessment_attempt_answer`: Mínimo 4 índices
- ✅ **Total global: Mínimo 15 índices**

#### 2.2 Verificar Índices Específicos
```bash
# Listar índices de assessment_attempt
psql -U postgres -d edugo_test_sprint01 -c "
    SELECT indexname, indexdef 
    FROM pg_indexes 
    WHERE tablename = 'assessment_attempt'
    ORDER BY indexname;
"
```
**Criterio de éxito:**
- ✅ Índice compuesto: `idx_attempt_student_assessment` en (student_id, assessment_id, created_at DESC)
- ✅ Índice parcial: `idx_attempt_idempotency_key` con WHERE clause
- ✅ Todos con `IF NOT EXISTS` (idempotentes)

#### 2.3 Verificar Performance de Índices
```bash
# Ejecutar ANALYZE para actualizar estadísticas
psql -U postgres -d edugo_test_sprint01 -c "ANALYZE assessment;"
psql -U postgres -d edugo_test_sprint01 -c "ANALYZE assessment_attempt;"
psql -U postgres -d edugo_test_sprint01 -c "ANALYZE assessment_attempt_answer;"

# Verificar que índices están siendo usados (requiere datos de prueba)
psql -U postgres -d edugo_test_sprint01 < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/seeds/assessment_seeds.sql

# EXPLAIN de query de historial
psql -U postgres -d edugo_test_sprint01 -c "
    EXPLAIN (ANALYZE, BUFFERS) 
    SELECT * FROM assessment_attempt 
    WHERE student_id = (SELECT id FROM users WHERE role='student' LIMIT 1)
    ORDER BY created_at DESC 
    LIMIT 10;
"
```
**Criterio de éxito:**
- ✅ Plan usa `Index Scan` (NO `Seq Scan`)
- ✅ Tiempo de ejecución <10ms (con seeds)

---

### 3. Seeds de Datos de Prueba

#### 3.1 Ejecutar Seeds
```bash
# Ejecutar script de seeds
psql -U postgres -d edugo_test_sprint01 < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/seeds/assessment_seeds.sql
```
**Criterio de éxito:** Output muestra `COMMIT` y verificación exitosa

#### 3.2 Verificar Datos Insertados
```bash
# Contar registros por tabla
psql -U postgres -d edugo_test_sprint01 -c "
    SELECT 
        'assessment' as tabla, COUNT(*) as filas FROM assessment
    UNION ALL
    SELECT 
        'assessment_attempt', COUNT(*) FROM assessment_attempt
    UNION ALL
    SELECT 
        'assessment_attempt_answer', COUNT(*) FROM assessment_attempt_answer
    UNION ALL
    SELECT 
        'material_summary_link', COUNT(*) FROM material_summary_link;
"
```
**Criterio de éxito:**
- ✅ assessment: ≥3 filas
- ✅ assessment_attempt: ≥3 filas
- ✅ assessment_attempt_answer: ≥15 filas (5 preguntas × 3 intentos)
- ✅ material_summary_link: ≥1 fila

#### 3.3 Validar Integridad de Seeds
```bash
# Verificar que scores calculados coinciden con respuestas
psql -U postgres -d edugo_test_sprint01 -c "
    SELECT 
        aa.attempt_id,
        COUNT(*) as total_answers,
        SUM(CASE WHEN aa.is_correct THEN 1 ELSE 0 END) as correct_answers,
        (SUM(CASE WHEN aa.is_correct THEN 1 ELSE 0 END)::FLOAT / COUNT(*)::FLOAT * 100)::INTEGER as calculated_score,
        at.score as stored_score,
        CASE 
            WHEN (SUM(CASE WHEN aa.is_correct THEN 1 ELSE 0 END)::FLOAT / COUNT(*)::FLOAT * 100)::INTEGER = at.score 
            THEN '✅ OK' 
            ELSE '❌ MISMATCH' 
        END as validation
    FROM assessment_attempt_answer aa
    JOIN assessment_attempt at ON aa.attempt_id = at.id
    GROUP BY aa.attempt_id, at.score;
"
```
**Criterio de éxito:**
- ✅ Todos los intentos muestran `✅ OK` (score calculado = score almacenado)
- ✅ Intento 1: 60% (3/5 correctas)
- ✅ Intento 2: 80% (4/5 correctas)
- ✅ Intento 3: 100% (5/5 correctas)

---

### 4. Rollback Script

#### 4.1 Verificar Sintaxis de Rollback
```bash
# Verificar sintaxis
psql -U postgres -d postgres --single-transaction --set ON_ERROR_STOP=on -f /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments_rollback.sql -o /dev/null
```
**Criterio de éxito:** Sin errores de sintaxis

#### 4.2 Ejecutar Rollback (en BD temporal)
```bash
# Crear BD temporal con migración
psql -U postgres -c "CREATE DATABASE edugo_rollback_test;"
psql -U postgres -d edugo_rollback_test < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/01_base_tables.sql
psql -U postgres -d edugo_rollback_test < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments.sql

# Ejecutar rollback
psql -U postgres -d edugo_rollback_test < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments_rollback.sql

# Verificar que tablas fueron eliminadas
psql -U postgres -d edugo_rollback_test -c "\dt assessment*"
```
**Criterio de éxito:**
- ✅ Output muestra "Did not find any relations" (0 tablas)
- ✅ Rollback completa con COMMIT

#### 4.3 Re-ejecutar Migración (Idempotencia)
```bash
# Re-ejecutar migración después de rollback
psql -U postgres -d edugo_rollback_test < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments.sql

# Verificar que tablas fueron recreadas
psql -U postgres -d edugo_rollback_test -c "\dt assessment*"
```
**Criterio de éxito:**
- ✅ 4 tablas recreadas exitosamente
- ✅ Sin errores (migración es idempotente)

#### 4.4 Limpiar BD Temporal
```bash
# Eliminar BD de test de rollback
psql -U postgres -c "DROP DATABASE edugo_rollback_test;"
```

---

### 5. Tests de Integridad Referencial

#### 5.1 Ejecutar Suite de Tests
```bash
# Ejecutar tests de integridad
psql -U postgres -d edugo_test_sprint01 < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/tests/test_assessments_integrity.sql
```
**Criterio de éxito:** Output muestra mínimo 7 mensajes "TEST PASSED"

#### 5.2 Validar Tests Individuales
```bash
# Contar tests pasados
psql -U postgres -d edugo_test_sprint01 < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/tests/test_assessments_integrity.sql 2>&1 | grep -c "TEST PASSED"
```
**Criterio de éxito:**
- ✅ Output: `7` (todos los tests pasaron)

```bash
# Verificar que NO hay tests fallidos
psql -U postgres -d edugo_test_sprint01 < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/tests/test_assessments_integrity.sql 2>&1 | grep "TEST FAILED"
```
**Criterio de éxito:**
- ✅ Sin output (no hay tests fallidos)

---

### 6. Validación de Comentarios y Documentación

#### 6.1 Verificar Comentarios en Tablas
```bash
# Listar comentarios de tablas
psql -U postgres -d edugo_test_sprint01 -c "
    SELECT 
        relname as table_name,
        obj_description(oid) as comment
    FROM pg_class 
    WHERE relname LIKE 'assessment%' 
      AND relkind = 'r'
    ORDER BY relname;
"
```
**Criterio de éxito:**
- ✅ 4 tablas tienen comentarios descriptivos
- ✅ `assessment_attempt` tiene comentario mencionando "INMUTABLE"

#### 6.2 Verificar Comentarios en Columnas Clave
```bash
# Comentarios de columnas críticas
psql -U postgres -d edugo_test_sprint01 -c "
    SELECT 
        a.attname as column_name,
        col_description(a.attrelid, a.attnum) as comment
    FROM pg_attribute a
    WHERE a.attrelid = 'assessment'::regclass
      AND a.attnum > 0
      AND NOT a.attisdropped
      AND col_description(a.attrelid, a.attnum) IS NOT NULL
    ORDER BY a.attnum;
"
```
**Criterio de éxito:**
- ✅ Mínimo 3 columnas con comentarios (mongo_document_id, pass_threshold, max_attempts)

#### 6.3 Verificar Comentarios en Índices
```bash
# Comentarios de índices
psql -U postgres -d edugo_test_sprint01 -c "
    SELECT 
        indexname,
        obj_description(indexrelid) as comment
    FROM pg_stat_user_indexes
    WHERE schemaname = 'public' 
      AND indexrelname LIKE 'idx_attempt_%'
      AND obj_description(indexrelid) IS NOT NULL;
"
```
**Criterio de éxito:**
- ✅ Mínimo 2 índices con comentarios explicativos

---

### 7. Performance y Tamaño

#### 7.1 Analizar Tamaño de Tablas e Índices
```bash
# Tamaño de tablas
psql -U postgres -d edugo_test_sprint01 -c "
    SELECT 
        schemaname,
        tablename,
        pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
    FROM pg_tables
    WHERE tablename LIKE 'assessment%'
    ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
"
```
**Criterio de éxito:**
- ✅ Todas las tablas <1 MB (sin datos masivos)

#### 7.2 Analizar Tamaño de Índices Individuales
```bash
# Tamaño de índices
psql -U postgres -d edugo_test_sprint01 -c "
    SELECT 
        schemaname,
        tablename,
        indexname,
        pg_size_pretty(pg_relation_size(indexname::regclass)) as index_size
    FROM pg_indexes
    WHERE tablename LIKE 'assessment%'
    ORDER BY pg_relation_size(indexname::regclass) DESC;
"
```
**Criterio de éxito:**
- ✅ Índices más grandes no superan 100 KB (con seeds)

---

### 8. Verificación de Dependencias

#### 8.1 Ejecutar Script de Verificación de Dependencias
```bash
# Crear y ejecutar script de verificación
cat > /tmp/verify_deps.sql << 'EOF'
DO $$
DECLARE
    materials_count INTEGER;
    users_count INTEGER;
    has_gen_uuid_v7 BOOLEAN;
BEGIN
    RAISE NOTICE '=== VERIFICANDO DEPENDENCIAS SPRINT 01 ===';
    
    SELECT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'gen_uuid_v7') INTO has_gen_uuid_v7;
    IF NOT has_gen_uuid_v7 THEN
        RAISE EXCEPTION '❌ FALTA: Función gen_uuid_v7()';
    ELSE
        RAISE NOTICE '✅ OK: gen_uuid_v7()';
    END IF;
    
    SELECT COUNT(*) INTO materials_count FROM materials;
    RAISE NOTICE '✅ OK: materials (% filas)', materials_count;
    
    SELECT COUNT(*) INTO users_count FROM users WHERE role = 'student';
    RAISE NOTICE '✅ OK: users (% estudiantes)', users_count;
    
    IF (SELECT current_setting('server_version_num')::INTEGER) < 150000 THEN
        RAISE WARNING '⚠️  PostgreSQL < 15';
    ELSE
        RAISE NOTICE '✅ OK: PostgreSQL 15+';
    END IF;
    
    RAISE NOTICE '=== DEPENDENCIAS OK ===';
END $$;
EOF

psql -U postgres -d edugo_test_sprint01 < /tmp/verify_deps.sql
```
**Criterio de éxito:**
- ✅ Output muestra "DEPENDENCIAS OK"
- ✅ Sin excepciones lanzadas

---

## Criterios de Éxito Globales del Sprint

### Checklist Final

- [ ] **Migración Principal**
  - [ ] 06_assessments.sql ejecuta sin errores
  - [ ] 4 tablas creadas con estructura correcta
  - [ ] Todos los constraints funcionando
  - [ ] Comentarios en tablas y columnas clave

- [ ] **Índices**
  - [ ] Mínimo 15 índices creados
  - [ ] Índice compuesto para historial de estudiante
  - [ ] Índice parcial para idempotency_key
  - [ ] EXPLAIN muestra uso de índices (no Seq Scan)

- [ ] **Seeds**
  - [ ] ≥3 assessments insertados
  - [ ] ≥3 attempts con scores variados (60%, 80%, 100%)
  - [ ] ≥15 answers con is_correct correcto
  - [ ] Scores calculados coinciden con stored

- [ ] **Rollback**
  - [ ] 06_assessments_rollback.sql elimina todas las tablas
  - [ ] Rollback + Re-migración funciona (idempotencia)

- [ ] **Tests de Integridad**
  - [ ] 7/7 tests de integridad pasando
  - [ ] FK constraints validados
  - [ ] Check constraints validados
  - [ ] ON DELETE CASCADE funciona

- [ ] **Documentación**
  - [ ] Comentarios en tablas críticas
  - [ ] Comentarios en columnas clave
  - [ ] README.md de Sprint 01 actualizado

---

## Comandos de Rollback en Caso de Error

### Si Migración Falla a Mitad de Ejecución
```bash
# Conectar a BD
psql -U postgres -d edugo_test_sprint01

# Rollback manual si transacción quedó abierta
ROLLBACK;

# Eliminar tablas parcialmente creadas
DROP TABLE IF EXISTS assessment_attempt_answer CASCADE;
DROP TABLE IF EXISTS assessment_attempt CASCADE;
DROP TABLE IF EXISTS assessment CASCADE;
DROP TABLE IF EXISTS material_summary_link CASCADE;
```

### Si Necesitas Limpiar BD Completamente
```bash
# Eliminar BD de test
psql -U postgres -c "DROP DATABASE IF EXISTS edugo_test_sprint01;"

# Recrear desde cero
psql -U postgres -c "CREATE DATABASE edugo_test_sprint01;"
psql -U postgres -d edugo_test_sprint01 < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/01_base_tables.sql
psql -U postgres -d edugo_test_sprint01 < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments.sql
```

---

## Reporte de Validación

Al completar todas las validaciones, generar reporte:

```bash
# Crear reporte de validación
cat > /Users/jhoanmedina/source/EduGo/Analisys/AnalisisEstandarizado/spec-01-evaluaciones/03-Sprints/Sprint-01-Schema-BD/VALIDATION_REPORT.md << EOF
# Reporte de Validación - Sprint 01

**Fecha:** $(date +%Y-%m-%d)  
**Ejecutado por:** $(whoami)

## Resultados

- ✅ Migración ejecutada exitosamente
- ✅ 4 tablas creadas
- ✅ 15+ índices creados
- ✅ Seeds insertados (3 assessments, 3 attempts, 15 answers)
- ✅ Rollback funcional
- ✅ 7/7 tests de integridad pasando

## Métricas

- **Tablas creadas:** 4
- **Índices creados:** $(psql -U postgres -d edugo_test_sprint01 -t -c "SELECT COUNT(*) FROM pg_indexes WHERE tablename LIKE 'assessment%';")
- **Constraints:** $(psql -U postgres -d edugo_test_sprint01 -t -c "SELECT COUNT(*) FROM pg_constraint WHERE conrelid::regclass::text LIKE 'assessment%';")
- **Tamaño total:** $(psql -U postgres -d edugo_test_sprint01 -t -c "SELECT pg_size_pretty(SUM(pg_total_relation_size(schemaname||'.'||tablename))) FROM pg_tables WHERE tablename LIKE 'assessment%';")

## Estado

**SPRINT 01: COMPLETADO ✅**
EOF

# Mostrar reporte
cat /Users/jhoanmedina/source/EduGo/Analisys/AnalisisEstandarizado/spec-01-evaluaciones/03-Sprints/Sprint-01-Schema-BD/VALIDATION_REPORT.md
```

---

**Generado con:** Claude Code  
**Sprint:** 01/06  
**Última actualización:** 2025-11-14
