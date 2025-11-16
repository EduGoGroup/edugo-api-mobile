# Tareas del Sprint 01 - Schema de Base de Datos

## Objetivo
Crear schema PostgreSQL completo para el Sistema de Evaluaciones con 4 tablas, migraciones idempotentes, 15+ índices optimizados, constraints de integridad referencial y datos de prueba.

---

## Tareas

### TASK-01-001: Crear Migración Principal de Assessments
**Tipo:** feature  
**Prioridad:** HIGH  
**Estimación:** 3h  
**Asignado a:** @ai-executor

#### Descripción
Crear archivo de migración SQL que defina las 4 tablas del sistema de evaluaciones con sus constraints, checks, y relaciones FK.

#### Pasos de Implementación

1. Crear archivo de migración en ruta:
   ```
   /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments.sql
   ```

2. Implementar en este orden exacto:
   
   **A. Header y función UUIDv7:**
   ```sql
   -- Migration: 06_assessments.sql
   -- Description: Schema para Sistema de Evaluaciones
   -- Dependencies: 01_base_tables.sql (materials, users)
   -- Date: 2025-11-14
   
   BEGIN;
   
   -- Verificar que función gen_uuid_v7() existe
   DO $$
   BEGIN
       IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'gen_uuid_v7') THEN
           RAISE EXCEPTION 'Función gen_uuid_v7() no encontrada. Ejecutar migración 01_base_tables.sql primero.';
       END IF;
   END $$;
   ```

   **B. Tabla assessment:**
   ```sql
   CREATE TABLE IF NOT EXISTS assessment (
       id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
       material_id UUID NOT NULL,
       mongo_document_id VARCHAR(24) NOT NULL,
       title VARCHAR(255) NOT NULL,
       total_questions INTEGER NOT NULL CHECK (total_questions > 0 AND total_questions <= 100),
       pass_threshold INTEGER NOT NULL DEFAULT 70 CHECK (pass_threshold >= 0 AND pass_threshold <= 100),
       max_attempts INTEGER DEFAULT NULL,
       time_limit_minutes INTEGER DEFAULT NULL,
       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
       
       CONSTRAINT fk_assessment_material 
           FOREIGN KEY (material_id) 
           REFERENCES materials(id) 
           ON DELETE CASCADE,
       
       CONSTRAINT unique_material_assessment 
           UNIQUE (material_id)
   );
   
   COMMENT ON TABLE assessment IS 'Metadatos de evaluaciones asociadas a materiales educativos';
   COMMENT ON COLUMN assessment.mongo_document_id IS 'ObjectId del documento en MongoDB collection material_assessment';
   COMMENT ON COLUMN assessment.pass_threshold IS 'Porcentaje mínimo para aprobar (0-100)';
   COMMENT ON COLUMN assessment.max_attempts IS 'Máximo de intentos permitidos (NULL = ilimitado)';
   ```

   **C. Tabla assessment_attempt:**
   ```sql
   CREATE TABLE IF NOT EXISTS assessment_attempt (
       id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
       assessment_id UUID NOT NULL,
       student_id UUID NOT NULL,
       score INTEGER NOT NULL CHECK (score >= 0 AND score <= 100),
       max_score INTEGER NOT NULL DEFAULT 100,
       time_spent_seconds INTEGER NOT NULL CHECK (time_spent_seconds > 0 AND time_spent_seconds <= 7200),
       started_at TIMESTAMP NOT NULL,
       completed_at TIMESTAMP NOT NULL,
       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
       idempotency_key VARCHAR(64) DEFAULT NULL,
       
       CONSTRAINT fk_attempt_assessment 
           FOREIGN KEY (assessment_id) 
           REFERENCES assessment(id) 
           ON DELETE CASCADE,
       
       CONSTRAINT fk_attempt_student 
           FOREIGN KEY (student_id) 
           REFERENCES users(id) 
           ON DELETE CASCADE,
       
       CONSTRAINT check_attempt_time_logical 
           CHECK (completed_at > started_at),
       
       CONSTRAINT check_attempt_duration 
           CHECK (EXTRACT(EPOCH FROM (completed_at - started_at)) = time_spent_seconds),
       
       CONSTRAINT unique_idempotency_key 
           UNIQUE (idempotency_key)
   );
   
   COMMENT ON TABLE assessment_attempt IS 'Intentos de estudiantes en evaluaciones (INMUTABLE)';
   COMMENT ON COLUMN assessment_attempt.time_spent_seconds IS 'Tiempo total del intento en segundos (max 2 horas)';
   COMMENT ON COLUMN assessment_attempt.idempotency_key IS 'Clave para prevenir intentos duplicados (Post-MVP)';
   ```

   **D. Tabla assessment_attempt_answer:**
   ```sql
   CREATE TABLE IF NOT EXISTS assessment_attempt_answer (
       id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
       attempt_id UUID NOT NULL,
       question_id VARCHAR(50) NOT NULL,
       selected_answer_id VARCHAR(50) NOT NULL,
       is_correct BOOLEAN NOT NULL,
       time_spent_seconds INTEGER NOT NULL CHECK (time_spent_seconds >= 0),
       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
       
       CONSTRAINT fk_answer_attempt 
           FOREIGN KEY (attempt_id) 
           REFERENCES assessment_attempt(id) 
           ON DELETE CASCADE,
       
       CONSTRAINT unique_attempt_question 
           UNIQUE (attempt_id, question_id)
   );
   
   COMMENT ON TABLE assessment_attempt_answer IS 'Respuestas individuales de cada pregunta en un intento';
   COMMENT ON COLUMN assessment_attempt_answer.question_id IS 'ID de la pregunta en MongoDB';
   COMMENT ON COLUMN assessment_attempt_answer.selected_answer_id IS 'ID de la opción seleccionada';
   COMMENT ON COLUMN assessment_attempt_answer.is_correct IS 'Si la respuesta fue correcta (calculado en servidor)';
   ```

   **E. Tabla material_summary_link (Opcional):**
   ```sql
   CREATE TABLE IF NOT EXISTS material_summary_link (
       id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
       material_id UUID NOT NULL,
       mongo_summary_id VARCHAR(24) NOT NULL,
       mongo_assessment_id VARCHAR(24) DEFAULT NULL,
       link_type VARCHAR(20) NOT NULL CHECK (link_type IN ('summary', 'assessment', 'both')),
       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
       
       CONSTRAINT fk_link_material 
           FOREIGN KEY (material_id) 
           REFERENCES materials(id) 
           ON DELETE CASCADE,
       
       CONSTRAINT unique_material_link 
           UNIQUE (material_id, link_type)
   );
   
   COMMENT ON TABLE material_summary_link IS 'Enlaces entre materiales PostgreSQL y documentos MongoDB';
   COMMENT ON COLUMN material_summary_link.link_type IS 'Tipo de enlace: summary, assessment, o both';
   ```

   **F. Commit:**
   ```sql
   COMMIT;
   ```

3. Agregar tests de integridad al final del archivo:
   ```sql
   -- Tests de Integridad (comentados, ejecutar manualmente)
   -- DO $$
   -- BEGIN
   --     ASSERT (SELECT COUNT(*) FROM pg_tables WHERE tablename = 'assessment') = 1, 'Tabla assessment no creada';
   --     ASSERT (SELECT COUNT(*) FROM pg_tables WHERE tablename = 'assessment_attempt') = 1, 'Tabla assessment_attempt no creada';
   --     ASSERT (SELECT COUNT(*) FROM pg_tables WHERE tablename = 'assessment_attempt_answer') = 1, 'Tabla assessment_attempt_answer no creada';
   --     RAISE NOTICE 'Todas las tablas creadas exitosamente';
   -- END $$;
   ```

#### Criterios de Aceptación
- [ ] Archivo creado en `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments.sql`
- [ ] 4 tablas definidas con comentarios
- [ ] Constraints FK con ON DELETE CASCADE
- [ ] Checks de validación en campos numéricos
- [ ] Unique constraints correctos
- [ ] Función gen_uuid_v7() verificada antes de CREATE TABLE
- [ ] Migración idempotente (IF NOT EXISTS)

#### Comandos de Validación
```bash
# Verificar sintaxis SQL
psql -U postgres -d postgres -c "\i /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments.sql" --dry-run

# Ejecutar migración en BD de prueba
psql -U postgres -d edugo_test < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments.sql

# Verificar tablas creadas
psql -U postgres -d edugo_test -c "\dt assessment*"

# Verificar constraints
psql -U postgres -d edugo_test -c "SELECT conname, contype FROM pg_constraint WHERE conrelid = 'assessment'::regclass;"
```

#### Dependencias
- Requiere: Tabla `materials` existente
- Requiere: Tabla `users` existente
- Requiere: Función `gen_uuid_v7()` existente
- Usa: PostgreSQL 15+

#### Tiempo Estimado
3 horas

---

### TASK-01-002: Crear Índices Optimizados
**Tipo:** feature  
**Prioridad:** HIGH  
**Estimación:** 2h  
**Asignado a:** @ai-executor

#### Descripción
Crear índices B-tree y compuestos para optimizar queries frecuentes del sistema de evaluaciones.

#### Pasos de Implementación

1. Agregar al final de `06_assessments.sql` (antes del COMMIT):

   **Índices para tabla assessment:**
   ```sql
   -- Índices: assessment
   CREATE INDEX IF NOT EXISTS idx_assessment_material_id 
       ON assessment(material_id);
   
   CREATE INDEX IF NOT EXISTS idx_assessment_mongo_document_id 
       ON assessment(mongo_document_id);
   
   CREATE INDEX IF NOT EXISTS idx_assessment_created_at 
       ON assessment(created_at DESC);
   
   COMMENT ON INDEX idx_assessment_material_id IS 'Query más frecuente: obtener assessment por material';
   ```

   **Índices para tabla assessment_attempt:**
   ```sql
   -- Índices: assessment_attempt
   CREATE INDEX IF NOT EXISTS idx_attempt_assessment_id 
       ON assessment_attempt(assessment_id);
   
   CREATE INDEX IF NOT EXISTS idx_attempt_student_id 
       ON assessment_attempt(student_id);
   
   CREATE INDEX IF NOT EXISTS idx_attempt_student_assessment 
       ON assessment_attempt(student_id, assessment_id, created_at DESC);
   
   CREATE INDEX IF NOT EXISTS idx_attempt_created_at 
       ON assessment_attempt(created_at DESC);
   
   CREATE INDEX IF NOT EXISTS idx_attempt_score 
       ON assessment_attempt(score);
   
   CREATE INDEX IF NOT EXISTS idx_attempt_idempotency_key 
       ON assessment_attempt(idempotency_key) 
       WHERE idempotency_key IS NOT NULL;
   
   COMMENT ON INDEX idx_attempt_student_assessment IS 'Historial de intentos de un estudiante en un assessment';
   COMMENT ON INDEX idx_attempt_idempotency_key IS 'Índice parcial para prevenir duplicados (Post-MVP)';
   ```

   **Índices para tabla assessment_attempt_answer:**
   ```sql
   -- Índices: assessment_attempt_answer
   CREATE INDEX IF NOT EXISTS idx_answer_attempt_id 
       ON assessment_attempt_answer(attempt_id);
   
   CREATE INDEX IF NOT EXISTS idx_answer_question_id 
       ON assessment_attempt_answer(question_id);
   
   CREATE INDEX IF NOT EXISTS idx_answer_is_correct 
       ON assessment_attempt_answer(is_correct);
   
   CREATE INDEX IF NOT EXISTS idx_answer_attempt_question 
       ON assessment_attempt_answer(attempt_id, question_id);
   
   COMMENT ON INDEX idx_answer_attempt_question IS 'Unique constraint enforcement + query optimization';
   ```

   **Índices para tabla material_summary_link:**
   ```sql
   -- Índices: material_summary_link
   CREATE INDEX IF NOT EXISTS idx_link_material_id 
       ON material_summary_link(material_id);
   
   CREATE INDEX IF NOT EXISTS idx_link_mongo_summary_id 
       ON material_summary_link(mongo_summary_id);
   
   CREATE INDEX IF NOT EXISTS idx_link_type 
       ON material_summary_link(link_type);
   ```

2. Crear análisis de tamaño estimado de índices:
   ```sql
   -- Análisis de Índices (comentado)
   -- SELECT 
   --     schemaname,
   --     tablename,
   --     indexname,
   --     pg_size_pretty(pg_relation_size(indexname::regclass)) as index_size
   -- FROM pg_indexes
   -- WHERE schemaname = 'public' AND tablename LIKE 'assessment%'
   -- ORDER BY pg_relation_size(indexname::regclass) DESC;
   ```

#### Criterios de Aceptación
- [ ] Mínimo 15 índices creados
- [ ] Índices con IF NOT EXISTS (idempotente)
- [ ] Índice compuesto para historial (student_id, assessment_id, created_at)
- [ ] Índice parcial para idempotency_key (WHERE NOT NULL)
- [ ] Comentarios en índices críticos
- [ ] Índices en orden descendente para created_at

#### Comandos de Validación
```bash
# Verificar índices creados
psql -U postgres -d edugo_test -c "SELECT indexname, indexdef FROM pg_indexes WHERE tablename LIKE 'assessment%' ORDER BY tablename, indexname;"

# Contar índices por tabla
psql -U postgres -d edugo_test -c "SELECT tablename, COUNT(*) as num_indexes FROM pg_indexes WHERE tablename LIKE 'assessment%' GROUP BY tablename;"

# Analizar uso de índices (después de queries)
psql -U postgres -d edugo_test -c "SELECT schemaname, tablename, indexname, idx_scan FROM pg_stat_user_indexes WHERE tablename LIKE 'assessment%' ORDER BY idx_scan DESC;"
```

#### Dependencias
- Requiere: TASK-01-001 completada (tablas creadas)
- Usa: PostgreSQL 15+ (soporte para índices parciales)

#### Tiempo Estimado
2 horas

---

### TASK-01-003: Crear Script de Seeds
**Tipo:** feature  
**Prioridad:** MEDIUM  
**Estimación:** 2h  
**Asignado a:** @ai-executor

#### Descripción
Crear script SQL con datos de prueba realistas para desarrollo y testing.

#### Pasos de Implementación

1. Crear archivo:
   ```
   /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/seeds/assessment_seeds.sql
   ```

2. Implementar seeds con datos realistas:
   ```sql
   -- Seeds: Sistema de Evaluaciones
   -- Description: Datos de prueba para desarrollo y testing
   -- Dependencies: 06_assessments.sql, materials y users existentes
   -- Date: 2025-11-14
   
   BEGIN;
   
   -- Verificar que tablas existen
   DO $$
   BEGIN
       IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = 'assessment') THEN
           RAISE EXCEPTION 'Tabla assessment no existe. Ejecutar 06_assessments.sql primero.';
       END IF;
   END $$;
   
   -- Limpiar datos existentes (solo en dev/test)
   TRUNCATE TABLE assessment_attempt_answer CASCADE;
   TRUNCATE TABLE assessment_attempt CASCADE;
   TRUNCATE TABLE assessment CASCADE;
   TRUNCATE TABLE material_summary_link CASCADE;
   
   -- Seed 1: Assessment de Material de Pascal
   INSERT INTO assessment (
       id,
       material_id,
       mongo_document_id,
       title,
       total_questions,
       pass_threshold,
       max_attempts,
       time_limit_minutes,
       created_at,
       updated_at
   ) VALUES (
       '01936d9a-0000-7000-a000-000000000001',
       (SELECT id FROM materials WHERE title ILIKE '%pascal%' LIMIT 1),
       '507f1f77bcf86cd799439011',
       'Cuestionario: Introducción a Pascal',
       5,
       70,
       NULL, -- Ilimitado
       15,
       NOW() - INTERVAL '30 days',
       NOW() - INTERVAL '30 days'
   );
   
   -- Seed 2: Assessment de Material de Python
   INSERT INTO assessment (
       id,
       material_id,
       mongo_document_id,
       title,
       total_questions,
       pass_threshold,
       max_attempts,
       time_limit_minutes,
       created_at,
       updated_at
   ) VALUES (
       '01936d9a-0000-7000-a000-000000000002',
       (SELECT id FROM materials WHERE title ILIKE '%python%' LIMIT 1),
       '507f1f77bcf86cd799439012',
       'Cuestionario: Fundamentos de Python',
       10,
       75,
       3, -- Máximo 3 intentos
       30,
       NOW() - INTERVAL '15 days',
       NOW() - INTERVAL '15 days'
   );
   
   -- Seed 3: Assessment de Material de Algoritmos
   INSERT INTO assessment (
       id,
       material_id,
       mongo_document_id,
       title,
       total_questions,
       pass_threshold,
       max_attempts,
       time_limit_minutes,
       created_at,
       updated_at
   ) VALUES (
       '01936d9a-0000-7000-a000-000000000003',
       (SELECT id FROM materials WHERE title ILIKE '%algoritmo%' LIMIT 1),
       '507f1f77bcf86cd799439013',
       'Cuestionario: Algoritmos de Ordenamiento',
       8,
       80,
       NULL,
       20,
       NOW() - INTERVAL '7 days',
       NOW() - INTERVAL '7 days'
   );
   
   -- Attempts: Estudiante 1 - Pascal (3 intentos)
   INSERT INTO assessment_attempt (
       id,
       assessment_id,
       student_id,
       score,
       max_score,
       time_spent_seconds,
       started_at,
       completed_at,
       created_at
   ) VALUES 
   -- Intento 1: Reprobado
   (
       '01936d9a-0000-7000-b000-000000000001',
       '01936d9a-0000-7000-a000-000000000001',
       (SELECT id FROM users WHERE role = 'student' LIMIT 1),
       60,
       100,
       420, -- 7 minutos
       NOW() - INTERVAL '29 days',
       NOW() - INTERVAL '29 days' + INTERVAL '420 seconds',
       NOW() - INTERVAL '29 days'
   ),
   -- Intento 2: Aprobado
   (
       '01936d9a-0000-7000-b000-000000000002',
       '01936d9a-0000-7000-a000-000000000001',
       (SELECT id FROM users WHERE role = 'student' LIMIT 1),
       80,
       100,
       540, -- 9 minutos
       NOW() - INTERVAL '28 days',
       NOW() - INTERVAL '28 days' + INTERVAL '540 seconds',
       NOW() - INTERVAL '28 days'
   ),
   -- Intento 3: Excelente
   (
       '01936d9a-0000-7000-b000-000000000003',
       '01936d9a-0000-7000-a000-000000000001',
       (SELECT id FROM users WHERE role = 'student' LIMIT 1),
       100,
       100,
       360, -- 6 minutos
       NOW() - INTERVAL '27 days',
       NOW() - INTERVAL '27 days' + INTERVAL '360 seconds',
       NOW() - INTERVAL '27 days'
   );
   
   -- Answers del Intento 1 (60% - 3 de 5 correctas)
   INSERT INTO assessment_attempt_answer (
       id,
       attempt_id,
       question_id,
       selected_answer_id,
       is_correct,
       time_spent_seconds,
       created_at
   ) VALUES
   ('01936d9a-0000-7000-c000-000000000001', '01936d9a-0000-7000-b000-000000000001', 'q1', 'q1_opt_a', TRUE, 60, NOW() - INTERVAL '29 days'),
   ('01936d9a-0000-7000-c000-000000000002', '01936d9a-0000-7000-b000-000000000001', 'q2', 'q2_opt_b', FALSE, 90, NOW() - INTERVAL '29 days'),
   ('01936d9a-0000-7000-c000-000000000003', '01936d9a-0000-7000-b000-000000000001', 'q3', 'q3_opt_c', TRUE, 80, NOW() - INTERVAL '29 days'),
   ('01936d9a-0000-7000-c000-000000000004', '01936d9a-0000-7000-b000-000000000001', 'q4', 'q4_opt_a', FALSE, 100, NOW() - INTERVAL '29 days'),
   ('01936d9a-0000-7000-c000-000000000005', '01936d9a-0000-7000-b000-000000000001', 'q5', 'q5_opt_d', TRUE, 90, NOW() - INTERVAL '29 days');
   
   -- Answers del Intento 2 (80% - 4 de 5 correctas)
   INSERT INTO assessment_attempt_answer (
       id,
       attempt_id,
       question_id,
       selected_answer_id,
       is_correct,
       time_spent_seconds,
       created_at
   ) VALUES
   ('01936d9a-0000-7000-c000-000000000006', '01936d9a-0000-7000-b000-000000000002', 'q1', 'q1_opt_a', TRUE, 80, NOW() - INTERVAL '28 days'),
   ('01936d9a-0000-7000-c000-000000000007', '01936d9a-0000-7000-b000-000000000002', 'q2', 'q2_opt_c', TRUE, 120, NOW() - INTERVAL '28 days'),
   ('01936d9a-0000-7000-c000-000000000008', '01936d9a-0000-7000-b000-000000000002', 'q3', 'q3_opt_c', TRUE, 100, NOW() - INTERVAL '28 days'),
   ('01936d9a-0000-7000-c000-000000000009', '01936d9a-0000-7000-b000-000000000002', 'q4', 'q4_opt_b', FALSE, 120, NOW() - INTERVAL '28 days'),
   ('01936d9a-0000-7000-c000-000000000010', '01936d9a-0000-7000-b000-000000000002', 'q5', 'q5_opt_d', TRUE, 120, NOW() - INTERVAL '28 days');
   
   -- Answers del Intento 3 (100% - 5 de 5 correctas)
   INSERT INTO assessment_attempt_answer (
       id,
       attempt_id,
       question_id,
       selected_answer_id,
       is_correct,
       time_spent_seconds,
       created_at
   ) VALUES
   ('01936d9a-0000-7000-c000-000000000011', '01936d9a-0000-7000-b000-000000000003', 'q1', 'q1_opt_a', TRUE, 60, NOW() - INTERVAL '27 days'),
   ('01936d9a-0000-7000-c000-000000000012', '01936d9a-0000-7000-b000-000000000003', 'q2', 'q2_opt_c', TRUE, 70, NOW() - INTERVAL '27 days'),
   ('01936d9a-0000-7000-c000-000000000013', '01936d9a-0000-7000-b000-000000000003', 'q3', 'q3_opt_c', TRUE, 65, NOW() - INTERVAL '27 days'),
   ('01936d9a-0000-7000-c000-000000000014', '01936d9a-0000-7000-b000-000000000003', 'q4', 'q4_opt_c', TRUE, 80, NOW() - INTERVAL '27 days'),
   ('01936d9a-0000-7000-c000-000000000015', '01936d9a-0000-7000-b000-000000000003', 'q5', 'q5_opt_d', TRUE, 85, NOW() - INTERVAL '27 days');
   
   -- Material Summary Links
   INSERT INTO material_summary_link (
       id,
       material_id,
       mongo_summary_id,
       mongo_assessment_id,
       link_type,
       created_at
   ) VALUES
   (
       '01936d9a-0000-7000-d000-000000000001',
       (SELECT id FROM materials WHERE title ILIKE '%pascal%' LIMIT 1),
       '507f1f77bcf86cd799439010',
       '507f1f77bcf86cd799439011',
       'both',
       NOW() - INTERVAL '30 days'
   );
   
   COMMIT;
   
   -- Verificación de seeds
   DO $$
   DECLARE
       assessment_count INTEGER;
       attempt_count INTEGER;
       answer_count INTEGER;
   BEGIN
       SELECT COUNT(*) INTO assessment_count FROM assessment;
       SELECT COUNT(*) INTO attempt_count FROM assessment_attempt;
       SELECT COUNT(*) INTO answer_count FROM assessment_attempt_answer;
       
       RAISE NOTICE 'Seeds insertados exitosamente:';
       RAISE NOTICE '  - Assessments: %', assessment_count;
       RAISE NOTICE '  - Attempts: %', attempt_count;
       RAISE NOTICE '  - Answers: %', answer_count;
       
       ASSERT assessment_count >= 3, 'Se esperaban al menos 3 assessments';
       ASSERT attempt_count >= 3, 'Se esperaban al menos 3 attempts';
       ASSERT answer_count >= 15, 'Se esperaban al menos 15 answers';
   END $$;
   ```

#### Criterios de Aceptación
- [ ] Archivo creado en `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/seeds/assessment_seeds.sql`
- [ ] Mínimo 3 assessments de prueba
- [ ] Mínimo 3 attempts con diferentes scores (60%, 80%, 100%)
- [ ] Todas las answers con is_correct calculado
- [ ] TRUNCATE antes de INSERT (solo dev/test)
- [ ] Verificación automática al final
- [ ] UUIDs hardcodeados (reproducibles)

#### Comandos de Validación
```bash
# Ejecutar seeds
psql -U postgres -d edugo_test < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/seeds/assessment_seeds.sql

# Verificar datos insertados
psql -U postgres -d edugo_test -c "SELECT COUNT(*) as assessments FROM assessment;"
psql -U postgres -d edugo_test -c "SELECT COUNT(*) as attempts FROM assessment_attempt;"
psql -U postgres -d edugo_test -c "SELECT COUNT(*) as answers FROM assessment_attempt_answer;"

# Verificar scores calculados correctamente
psql -U postgres -d edugo_test -c "
    SELECT 
        aa.attempt_id,
        COUNT(*) as total_answers,
        SUM(CASE WHEN aa.is_correct THEN 1 ELSE 0 END) as correct_answers,
        (SUM(CASE WHEN aa.is_correct THEN 1 ELSE 0 END)::FLOAT / COUNT(*)::FLOAT * 100)::INTEGER as calculated_score,
        at.score as stored_score
    FROM assessment_attempt_answer aa
    JOIN assessment_attempt at ON aa.attempt_id = at.id
    GROUP BY aa.attempt_id, at.score;
"
```

#### Dependencias
- Requiere: TASK-01-001 completada (tablas creadas)
- Requiere: Tabla `materials` con datos
- Requiere: Tabla `users` con estudiantes

#### Tiempo Estimado
2 horas

---

### TASK-01-004: Crear Script de Rollback
**Tipo:** feature  
**Prioridad:** MEDIUM  
**Estimación:** 1h  
**Asignado a:** @ai-executor

#### Descripción
Crear script SQL para revertir la migración 06_assessments.sql de forma segura.

#### Pasos de Implementación

1. Crear archivo:
   ```
   /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments_rollback.sql
   ```

2. Implementar rollback en orden inverso a la creación:
   ```sql
   -- Rollback: 06_assessments.sql
   -- Description: Revertir schema de Sistema de Evaluaciones
   -- WARNING: Esta operación elimina TODOS los datos de evaluaciones
   -- Date: 2025-11-14
   
   BEGIN;
   
   -- Verificar que estamos en ambiente correcto
   DO $$
   BEGIN
       IF current_database() = 'edugo_prod' THEN
           RAISE EXCEPTION 'ROLLBACK PROHIBIDO EN PRODUCCIÓN. Use migraciones controladas.';
       END IF;
       
       RAISE NOTICE 'Ejecutando rollback en base de datos: %', current_database();
   END $$;
   
   -- Drop tablas en orden inverso (respetando FKs)
   DROP TABLE IF EXISTS material_summary_link CASCADE;
   DROP TABLE IF EXISTS assessment_attempt_answer CASCADE;
   DROP TABLE IF EXISTS assessment_attempt CASCADE;
   DROP TABLE IF EXISTS assessment CASCADE;
   
   -- Drop índices huérfanos (si existen)
   DROP INDEX IF EXISTS idx_assessment_material_id;
   DROP INDEX IF EXISTS idx_assessment_mongo_document_id;
   DROP INDEX IF EXISTS idx_assessment_created_at;
   DROP INDEX IF EXISTS idx_attempt_assessment_id;
   DROP INDEX IF EXISTS idx_attempt_student_id;
   DROP INDEX IF EXISTS idx_attempt_student_assessment;
   DROP INDEX IF EXISTS idx_attempt_created_at;
   DROP INDEX IF EXISTS idx_attempt_score;
   DROP INDEX IF EXISTS idx_attempt_idempotency_key;
   DROP INDEX IF EXISTS idx_answer_attempt_id;
   DROP INDEX IF EXISTS idx_answer_question_id;
   DROP INDEX IF EXISTS idx_answer_is_correct;
   DROP INDEX IF EXISTS idx_answer_attempt_question;
   DROP INDEX IF EXISTS idx_link_material_id;
   DROP INDEX IF EXISTS idx_link_mongo_summary_id;
   DROP INDEX IF EXISTS idx_link_type;
   
   COMMIT;
   
   -- Verificación de rollback
   DO $$
   BEGIN
       IF EXISTS (SELECT 1 FROM pg_tables WHERE tablename IN ('assessment', 'assessment_attempt', 'assessment_attempt_answer', 'material_summary_link')) THEN
           RAISE WARNING 'Algunas tablas no fueron eliminadas correctamente';
       ELSE
           RAISE NOTICE 'Rollback completado exitosamente. Todas las tablas eliminadas.';
       END IF;
   END $$;
   ```

#### Criterios de Aceptación
- [ ] Archivo creado en `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments_rollback.sql`
- [ ] Prevención de ejecución en producción
- [ ] Drop en orden correcto (respetando FKs)
- [ ] CASCADE para eliminar dependencias
- [ ] Verificación al final
- [ ] Idempotente (IF EXISTS)

#### Comandos de Validación
```bash
# Ejecutar rollback en BD de prueba
psql -U postgres -d edugo_test < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments_rollback.sql

# Verificar que tablas fueron eliminadas
psql -U postgres -d edugo_test -c "\dt assessment*"

# Re-ejecutar migración (debe funcionar)
psql -U postgres -d edugo_test < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments.sql

# Verificar que tablas fueron recreadas
psql -U postgres -d edugo_test -c "\dt assessment*"
```

#### Dependencias
- Requiere: Conocimiento de estructura de TASK-01-001

#### Tiempo Estimado
1 hora

---

### TASK-01-005: Tests de Integridad Referencial
**Tipo:** test  
**Prioridad:** HIGH  
**Estimación:** 2h  
**Asignado a:** @ai-executor

#### Descripción
Crear suite de tests SQL para validar integridad referencial, constraints y triggers.

#### Pasos de Implementación

1. Crear archivo:
   ```
   /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/tests/test_assessments_integrity.sql
   ```

2. Implementar tests de integridad:
   ```sql
   -- Tests: Integridad Referencial del Sistema de Evaluaciones
   -- Description: Suite de tests para validar constraints y FKs
   -- Date: 2025-11-14
   
   BEGIN;
   
   -- Setup: Crear datos de prueba
   DO $$
   DECLARE
       test_material_id UUID;
       test_user_id UUID;
       test_assessment_id UUID;
   BEGIN
       -- Crear material de prueba
       INSERT INTO materials (id, title, content, processing_status, created_at)
       VALUES (gen_uuid_v7(), 'Test Material', 'Test content', 'completed', NOW())
       RETURNING id INTO test_material_id;
       
       -- Crear usuario de prueba
       INSERT INTO users (id, email, name, role, created_at)
       VALUES (gen_uuid_v7(), 'test@test.com', 'Test User', 'student', NOW())
       RETURNING id INTO test_user_id;
       
       RAISE NOTICE 'Setup completado: material_id=%, user_id=%', test_material_id, test_user_id;
   END $$;
   
   -- TEST 1: FK Constraint assessment -> materials
   DO $$
   BEGIN
       BEGIN
           INSERT INTO assessment (material_id, mongo_document_id, title, total_questions)
           VALUES ('00000000-0000-0000-0000-000000000000', '507f1f77bcf86cd799439011', 'Test', 5);
           
           RAISE EXCEPTION 'TEST FAILED: FK constraint assessment->materials no funcionó';
       EXCEPTION
           WHEN foreign_key_violation THEN
               RAISE NOTICE 'TEST PASSED: FK constraint assessment->materials OK';
       END;
   END $$;
   
   -- TEST 2: Unique Constraint material_id en assessment
   DO $$
   DECLARE
       test_material_id UUID;
   BEGIN
       SELECT id INTO test_material_id FROM materials WHERE title = 'Test Material' LIMIT 1;
       
       -- Primer insert OK
       INSERT INTO assessment (material_id, mongo_document_id, title, total_questions)
       VALUES (test_material_id, '507f1f77bcf86cd799439011', 'Test 1', 5);
       
       -- Segundo insert debe fallar
       BEGIN
           INSERT INTO assessment (material_id, mongo_document_id, title, total_questions)
           VALUES (test_material_id, '507f1f77bcf86cd799439012', 'Test 2', 5);
           
           RAISE EXCEPTION 'TEST FAILED: Unique constraint material_id no funcionó';
       EXCEPTION
           WHEN unique_violation THEN
               RAISE NOTICE 'TEST PASSED: Unique constraint material_id OK';
       END;
   END $$;
   
   -- TEST 3: Check Constraint total_questions (1-100)
   DO $$
   DECLARE
       test_material_id UUID;
   BEGIN
       INSERT INTO materials (id, title, content, processing_status, created_at)
       VALUES (gen_uuid_v7(), 'Test Material 2', 'Test', 'completed', NOW())
       RETURNING id INTO test_material_id;
       
       -- Intentar total_questions = 0 (debe fallar)
       BEGIN
           INSERT INTO assessment (material_id, mongo_document_id, title, total_questions)
           VALUES (test_material_id, '507f1f77bcf86cd799439013', 'Test', 0);
           
           RAISE EXCEPTION 'TEST FAILED: Check constraint total_questions no funcionó';
       EXCEPTION
           WHEN check_violation THEN
               RAISE NOTICE 'TEST PASSED: Check constraint total_questions OK';
       END;
   END $$;
   
   -- TEST 4: Check Constraint pass_threshold (0-100)
   DO $$
   DECLARE
       test_material_id UUID;
   BEGIN
       INSERT INTO materials (id, title, content, processing_status, created_at)
       VALUES (gen_uuid_v7(), 'Test Material 3', 'Test', 'completed', NOW())
       RETURNING id INTO test_material_id;
       
       -- Intentar pass_threshold = 150 (debe fallar)
       BEGIN
           INSERT INTO assessment (material_id, mongo_document_id, title, total_questions, pass_threshold)
           VALUES (test_material_id, '507f1f77bcf86cd799439014', 'Test', 5, 150);
           
           RAISE EXCEPTION 'TEST FAILED: Check constraint pass_threshold no funcionó';
       EXCEPTION
           WHEN check_violation THEN
               RAISE NOTICE 'TEST PASSED: Check constraint pass_threshold OK';
       END;
   END $$;
   
   -- TEST 5: FK Constraint attempt -> assessment (ON DELETE CASCADE)
   DO $$
   DECLARE
       test_material_id UUID;
       test_user_id UUID;
       test_assessment_id UUID;
       attempt_count INTEGER;
   BEGIN
       SELECT id INTO test_user_id FROM users WHERE email = 'test@test.com' LIMIT 1;
       
       INSERT INTO materials (id, title, content, processing_status, created_at)
       VALUES (gen_uuid_v7(), 'Test Material 4', 'Test', 'completed', NOW())
       RETURNING id INTO test_material_id;
       
       INSERT INTO assessment (id, material_id, mongo_document_id, title, total_questions)
       VALUES (gen_uuid_v7(), test_material_id, '507f1f77bcf86cd799439015', 'Test', 5)
       RETURNING id INTO test_assessment_id;
       
       INSERT INTO assessment_attempt (assessment_id, student_id, score, time_spent_seconds, started_at, completed_at)
       VALUES (test_assessment_id, test_user_id, 80, 300, NOW(), NOW() + INTERVAL '300 seconds');
       
       -- Eliminar assessment debe eliminar attempt (CASCADE)
       DELETE FROM assessment WHERE id = test_assessment_id;
       
       SELECT COUNT(*) INTO attempt_count FROM assessment_attempt WHERE assessment_id = test_assessment_id;
       
       IF attempt_count > 0 THEN
           RAISE EXCEPTION 'TEST FAILED: ON DELETE CASCADE no funcionó';
       ELSE
           RAISE NOTICE 'TEST PASSED: ON DELETE CASCADE OK';
       END IF;
   END $$;
   
   -- TEST 6: Check Constraint completed_at > started_at
   DO $$
   DECLARE
       test_assessment_id UUID;
       test_user_id UUID;
   BEGIN
       SELECT id INTO test_assessment_id FROM assessment WHERE title = 'Test 1' LIMIT 1;
       SELECT id INTO test_user_id FROM users WHERE email = 'test@test.com' LIMIT 1;
       
       BEGIN
           INSERT INTO assessment_attempt (assessment_id, student_id, score, time_spent_seconds, started_at, completed_at)
           VALUES (test_assessment_id, test_user_id, 80, 300, NOW(), NOW() - INTERVAL '1 hour');
           
           RAISE EXCEPTION 'TEST FAILED: Check constraint completed_at > started_at no funcionó';
       EXCEPTION
           WHEN check_violation THEN
               RAISE NOTICE 'TEST PASSED: Check constraint completed_at > started_at OK';
       END;
   END $$;
   
   -- TEST 7: Unique Constraint attempt_id + question_id
   DO $$
   DECLARE
       test_assessment_id UUID;
       test_user_id UUID;
       test_attempt_id UUID;
   BEGIN
       SELECT id INTO test_assessment_id FROM assessment WHERE title = 'Test 1' LIMIT 1;
       SELECT id INTO test_user_id FROM users WHERE email = 'test@test.com' LIMIT 1;
       
       INSERT INTO assessment_attempt (id, assessment_id, student_id, score, time_spent_seconds, started_at, completed_at)
       VALUES (gen_uuid_v7(), test_assessment_id, test_user_id, 80, 300, NOW(), NOW() + INTERVAL '300 seconds')
       RETURNING id INTO test_attempt_id;
       
       -- Primer insert OK
       INSERT INTO assessment_attempt_answer (attempt_id, question_id, selected_answer_id, is_correct, time_spent_seconds)
       VALUES (test_attempt_id, 'q1', 'q1_opt_a', TRUE, 60);
       
       -- Segundo insert de misma pregunta debe fallar
       BEGIN
           INSERT INTO assessment_attempt_answer (attempt_id, question_id, selected_answer_id, is_correct, time_spent_seconds)
           VALUES (test_attempt_id, 'q1', 'q1_opt_b', FALSE, 60);
           
           RAISE EXCEPTION 'TEST FAILED: Unique constraint attempt_id+question_id no funcionó';
       EXCEPTION
           WHEN unique_violation THEN
               RAISE NOTICE 'TEST PASSED: Unique constraint attempt_id+question_id OK';
       END;
   END $$;
   
   ROLLBACK;
   
   -- Resumen
   RAISE NOTICE '=================================';
   RAISE NOTICE 'TESTS DE INTEGRIDAD COMPLETADOS';
   RAISE NOTICE '=================================';
   ```

#### Criterios de Aceptación
- [ ] Archivo creado en `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/tests/test_assessments_integrity.sql`
- [ ] Mínimo 7 tests de integridad
- [ ] Tests de FK constraints
- [ ] Tests de unique constraints
- [ ] Tests de check constraints
- [ ] Tests de ON DELETE CASCADE
- [ ] Todos los tests pasan
- [ ] Uso de ROLLBACK (no afecta BD)

#### Comandos de Validación
```bash
# Ejecutar suite de tests
psql -U postgres -d edugo_test < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/tests/test_assessments_integrity.sql

# Verificar que todos los tests pasaron (buscar "TEST PASSED" en output)
psql -U postgres -d edugo_test < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/tests/test_assessments_integrity.sql | grep "TEST PASSED"

# Contar tests pasados vs fallidos
psql -U postgres -d edugo_test < /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/tests/test_assessments_integrity.sql | grep -c "TEST PASSED"
```

#### Dependencias
- Requiere: TASK-01-001 completada (tablas y constraints creados)
- Requiere: TASK-01-002 completada (índices creados)

#### Tiempo Estimado
2 horas

---

## Resumen del Sprint

**Total de Tareas:** 5  
**Estimación Total:** 10 horas  
**Archivos a Crear:** 4 archivos SQL + 1 carpeta tests

**Entregables:**
1. `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments.sql`
2. `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments_rollback.sql`
3. `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/seeds/assessment_seeds.sql`
4. `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/tests/test_assessments_integrity.sql`

**Criterios de Éxito:**
- [ ] 4 tablas PostgreSQL creadas
- [ ] 15+ índices optimizados
- [ ] Migraciones idempotentes
- [ ] Seeds de datos de prueba
- [ ] Tests de integridad pasando
- [ ] Rollback funcional

---

**Generado con:** Claude Code  
**Sprint:** 01/06  
**Última actualización:** 2025-11-14
