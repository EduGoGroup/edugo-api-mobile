# Preguntas y Decisiones del Sprint 01 - Schema de Base de Datos

## Q001: ¿Qué tipo de ID usar para las tablas?
**Contexto:** PostgreSQL soporta múltiples tipos de IDs (SERIAL, UUID v4, UUID v7, ULID). Necesitamos decidir qué usar para las tablas de evaluaciones.

**Opciones:**

### 1. **UUID v7 (Opción Recomendada)**
- **Pros:**
  - Ordenamiento cronológico natural (timestamp en primeros bytes)
  - Compatible con UUID estándar (16 bytes)
  - Mejor performance en índices B-tree que UUID v4
  - No revela información secuencial como SERIAL
  - Permite generación distribuida sin colisiones
  - Soportado por función `gen_uuid_v7()` ya existente en EduGo
  
- **Contras:**
  - Requiere función custom en PostgreSQL (ya implementada)
  - 16 bytes vs 8 bytes de BIGSERIAL
  - No estándar en PostgreSQL (requiere extension o función)

### 2. **UUID v4 (Random)**
- **Pros:**
  - Nativo en PostgreSQL con `gen_random_uuid()`
  - Completamente random, sin ordenamiento
  - No requiere funciones custom
  
- **Contras:**
  - Performance degradada en índices B-tree por aleatoriedad
  - Sin ordenamiento cronológico
  - Más fragmentación en disco

### 3. **BIGSERIAL**
- **Pros:**
  - Más pequeño (8 bytes)
  - Ordenamiento natural
  - Performance óptima
  
- **Contras:**
  - Revela información de volumen de negocio
  - Problemas en entornos distribuidos
  - Predecible (riesgo de seguridad menor)

**Decisión por Defecto:** **UUID v7**

**Justificación:**
- EduGo ya usa UUID v7 en otras tablas (materials, users)
- Mejor balance entre performance y distribución
- Ordenamiento cronológico útil para queries de historial
- Consistencia con arquitectura existente

**Implementación:**
```sql
-- Usar en todas las tablas de evaluaciones
CREATE TABLE assessment (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    -- ...
);

CREATE TABLE assessment_attempt (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    -- ...
);

CREATE TABLE assessment_attempt_answer (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    -- ...
);
```

---

## Q002: ¿Los intentos (attempts) deben ser mutables o inmutables?
**Contexto:** Necesitamos decidir si permitir UPDATE en la tabla `assessment_attempt` o tratarla como append-only (solo INSERT).

**Opciones:**

### 1. **Inmutable (Append-Only) - RECOMENDADO**
- **Pros:**
  - Auditoría completa (nunca se pierde información)
  - Previene modificación de calificaciones históricas
  - Más simple en lógica de negocio (no hay estados intermedios)
  - Mejor para analytics (datos nunca cambian)
  - Previene race conditions en concurrencia
  
- **Contras:**
  - No se puede "editar" un intento si hay error
  - Correcciones requieren nuevo registro
  - Más espacio en disco (pero marginal)

### 2. **Mutable (Permite UPDATE)**
- **Pros:**
  - Flexibilidad para correcciones
  - Menos filas en tabla (solo 1 fila por intento)
  
- **Contras:**
  - Riesgo de pérdida de información histórica
  - Complejidad en auditoría
  - Posibles race conditions
  - Necesita triggers para audit trail

**Decisión por Defecto:** **Inmutable (Append-Only)**

**Justificación:**
- Cumple con requisitos de auditoría educativa
- Previene modificación fraudulenta de calificaciones
- Alineado con mejores prácticas de event sourcing
- Más simple de implementar y mantener

**Implementación:**
```sql
-- No crear UPDATE triggers
-- En aplicación Go, solo permitir INSERT
CREATE TABLE assessment_attempt (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    -- Todos los campos NOT NULL (no hay "draft" state)
    assessment_id UUID NOT NULL,
    student_id UUID NOT NULL,
    score INTEGER NOT NULL,
    -- ...
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
    -- NO HAY updated_at (inmutable)
);

COMMENT ON TABLE assessment_attempt IS 'Intentos de estudiantes (INMUTABLE - solo INSERT permitido)';

-- Prevenir UPDATE a nivel aplicación (no a nivel BD por performance)
-- Implementar en Repository Go:
// func (r *AttemptRepository) Update() error {
//     return errors.New("assessment_attempt is immutable")
// }
```

---

## Q003: ¿Usar particionamiento en assessment_attempt?
**Contexto:** La tabla `assessment_attempt` puede crecer significativamente (estimado: 100K filas/año). ¿Particionar por fecha?

**Opciones:**

### 1. **No Particionar (MVP - RECOMENDADO)**
- **Pros:**
  - Más simple de implementar
  - Menos overhead de mantenimiento
  - PostgreSQL 15+ maneja 100K-1M filas sin problemas
  - Índices B-tree suficientes para performance
  
- **Contras:**
  - Puede necesitar particionamiento en el futuro
  - Queries full-table scan más lentas (pero mitigado por índices)

### 2. **Particionar por RANGE (created_at)**
- **Pros:**
  - Queries por rango de fecha más rápidas
  - Archivado más fácil (DROP partition antigua)
  - Mejor para tablas >10M filas
  
- **Contras:**
  - Complejidad operacional (crear partitions mensuales/anuales)
  - Más difícil de gestionar en migraciones
  - Overkill para volumen actual

**Decisión por Defecto:** **No Particionar (Post-MVP)**

**Justificación:**
- Volumen estimado (100K/año) no requiere particionamiento
- Índices optimizados suficientes para performance
- YAGNI (You Aren't Gonna Need It) - implementar cuando sea necesario
- Más fácil migrar a particionamiento después que deshacer particionamiento

**Implementación (MVP):**
```sql
-- Crear tabla sin particionamiento
CREATE TABLE assessment_attempt (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    -- ...
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Índice en created_at para queries de rango
CREATE INDEX idx_attempt_created_at ON assessment_attempt(created_at DESC);
```

**Implementación Futura (si >1M filas):**
```sql
-- Post-MVP: Convertir a tabla particionada
CREATE TABLE assessment_attempt_partitioned (
    -- mismo schema
) PARTITION BY RANGE (created_at);

-- Crear particiones por año
CREATE TABLE assessment_attempt_2025 PARTITION OF assessment_attempt_partitioned
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

CREATE TABLE assessment_attempt_2026 PARTITION OF assessment_attempt_partitioned
    FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');
```

---

## Q004: ¿Validar time_spent_seconds con EXTRACT vs. columna separada?
**Contexto:** Necesitamos verificar que `time_spent_seconds` coincida con `completed_at - started_at`.

**Opciones:**

### 1. **Check Constraint con EXTRACT - RECOMENDADO**
- **Pros:**
  - Validación a nivel BD (imposible insertar datos incorrectos)
  - Cero código de aplicación para validar
  - Garantía de integridad de datos
  
- **Contras:**
  - Constraint puede fallar por redondeo de milisegundos
  - Más estricto (menos flexible)

### 2. **Validación Solo en Aplicación**
- **Pros:**
  - Más flexible (puede aceptar ±1 segundo de diferencia)
  - Sin overhead de constraint checking
  
- **Contras:**
  - Posible inconsistencia si se inserta directamente por SQL
  - Requiere tests exhaustivos en aplicación

**Decisión por Defecto:** **Check Constraint con EXTRACT**

**Justificación:**
- Integridad de datos crítica para auditoría
- Previene bugs en aplicación
- PostgreSQL optimiza constraints eficientemente

**Implementación:**
```sql
CREATE TABLE assessment_attempt (
    -- ...
    time_spent_seconds INTEGER NOT NULL CHECK (time_spent_seconds > 0 AND time_spent_seconds <= 7200),
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP NOT NULL,
    
    -- Validar que completed_at > started_at
    CONSTRAINT check_attempt_time_logical 
        CHECK (completed_at > started_at),
    
    -- Validar que time_spent_seconds = completed_at - started_at
    CONSTRAINT check_attempt_duration 
        CHECK (EXTRACT(EPOCH FROM (completed_at - started_at)) = time_spent_seconds)
);
```

**Manejo de Edge Cases:**
```go
// En aplicación Go, calcular time_spent_seconds exactamente
completedAt := time.Now()
startedAt := attempt.StartedAt
timeSpent := int(completedAt.Sub(startedAt).Seconds())

// Guardar
attempt.CompletedAt = completedAt
attempt.TimeSpentSeconds = timeSpent
```

---

## Q005: ¿Índice compuesto vs índices separados para historial de estudiante?
**Contexto:** Query más frecuente: "Obtener historial de intentos de un estudiante en un assessment".

**Query típico:**
```sql
SELECT * FROM assessment_attempt 
WHERE student_id = $1 
  AND assessment_id = $2 
ORDER BY created_at DESC 
LIMIT 10;
```

**Opciones:**

### 1. **Índice Compuesto (student_id, assessment_id, created_at) - RECOMENDADO**
- **Pros:**
  - Covering index para query de historial
  - Una sola búsqueda en índice (sin lookup adicional)
  - Ordenamiento ya incluido
  
- **Contras:**
  - Más espacio en disco (pero marginal)
  - Solo útil si query usa student_id primero

### 2. **Índices Separados (student_id, assessment_id, created_at)**
- **Pros:**
  - Más flexible para queries variados
  - Menos espacio total
  
- **Contras:**
  - PostgreSQL debe combinar índices (bitmap scan)
  - Performance inferior para query específico de historial

**Decisión por Defecto:** **Índice Compuesto + Índices Separados**

**Justificación:**
- Query de historial es el caso de uso más frecuente (>60%)
- Espacio adicional es despreciable (<5% del tamaño de tabla)
- También mantener índices separados para queries alternativos

**Implementación:**
```sql
-- Índice compuesto para historial (query más frecuente)
CREATE INDEX idx_attempt_student_assessment 
    ON assessment_attempt(student_id, assessment_id, created_at DESC);

-- Índices separados para queries alternativos
CREATE INDEX idx_attempt_student_id 
    ON assessment_attempt(student_id);

CREATE INDEX idx_attempt_assessment_id 
    ON assessment_attempt(assessment_id);

CREATE INDEX idx_attempt_created_at 
    ON assessment_attempt(created_at DESC);

-- Comentar para documentar uso
COMMENT ON INDEX idx_attempt_student_assessment IS 'Historial de intentos de un estudiante en un assessment';
```

**Query optimizado usará:**
```sql
EXPLAIN ANALYZE
SELECT * FROM assessment_attempt 
WHERE student_id = '01936d9a-7f8e-7000-a000-123456789abc'
  AND assessment_id = '01936d9a-7f8e-7000-a000-987654321cba'
ORDER BY created_at DESC 
LIMIT 10;

-- Output esperado:
-- Index Scan using idx_attempt_student_assessment on assessment_attempt
-- (sin Seq Scan)
```

---

## Q006: ¿Campo mongo_document_id como VARCHAR(24) o TEXT?
**Contexto:** MongoDB ObjectId tiene longitud fija de 24 caracteres hexadecimales.

**Opciones:**

### 1. **VARCHAR(24) - RECOMENDADO**
- **Pros:**
  - Longitud fija, más eficiente que TEXT
  - Valida longitud automáticamente
  - Menos espacio en disco
  - Índices más pequeños
  
- **Contras:**
  - Si MongoDB cambia formato de ObjectId, requiere ALTER TABLE
  - Menos flexible

### 2. **TEXT**
- **Pros:**
  - Más flexible (acepta cualquier longitud)
  - No requiere cambios si formato cambia
  
- **Contras:**
  - Más espacio en disco
  - Sin validación de longitud

**Decisión por Defecto:** **VARCHAR(24)**

**Justificación:**
- MongoDB ObjectId es estándar desde 2009 (muy estable)
- Ahorro de espacio significativo en índices
- Validación automática de formato

**Implementación:**
```sql
CREATE TABLE assessment (
    -- ...
    mongo_document_id VARCHAR(24) NOT NULL,
    -- ...
);

-- Índice para joins con MongoDB
CREATE INDEX idx_assessment_mongo_document_id 
    ON assessment(mongo_document_id);

-- Validación adicional (opcional)
ALTER TABLE assessment 
    ADD CONSTRAINT check_mongo_id_format 
    CHECK (mongo_document_id ~ '^[0-9a-f]{24}$');
```

---

## Q007: ¿Agregar campo idempotency_key a assessment_attempt?
**Contexto:** Prevenir intentos duplicados por reintentos de HTTP requests.

**Opciones:**

### 1. **Agregar idempotency_key (Post-MVP) - RECOMENDADO**
- **Pros:**
  - Previene duplicados por network errors
  - Alineado con mejores prácticas de APIs REST
  - Útil para operaciones no-idempotentes (POST)
  
- **Contras:**
  - Campo adicional (32-64 bytes por fila)
  - Índice adicional
  - Complejidad en cliente (debe generar key)

### 2. **No Usar idempotency_key**
- **Pros:**
  - Más simple
  - Sin overhead de espacio
  
- **Contras:**
  - Posibles intentos duplicados si cliente reintenta
  - Sin protección contra double-submit

**Decisión por Defecto:** **Agregar campo (NULLABLE en MVP, NOT NULL en Post-MVP)**

**Justificación:**
- Protección contra edge cases de conectividad
- Campo nullable permite implementación gradual
- Estándar en APIs modernas (Stripe, PayPal, etc.)

**Implementación (MVP):**
```sql
CREATE TABLE assessment_attempt (
    -- ...
    idempotency_key VARCHAR(64) DEFAULT NULL,
    -- ...
    
    CONSTRAINT unique_idempotency_key 
        UNIQUE (idempotency_key)
);

-- Índice parcial (solo índices no-NULL para ahorrar espacio)
CREATE INDEX idx_attempt_idempotency_key 
    ON assessment_attempt(idempotency_key) 
    WHERE idempotency_key IS NOT NULL;
```

**Uso en API:**
```go
// Cliente genera key única
idempotencyKey := fmt.Sprintf("%s-%s-%d", 
    studentID, assessmentID, time.Now().Unix())

// POST con header
headers := map[string]string{
    "Idempotency-Key": idempotencyKey,
}

// Servidor valida antes de INSERT
existing, err := repo.FindByIdempotencyKey(idempotencyKey)
if existing != nil {
    // Retornar resultado existente (no crear duplicado)
    return existing, nil
}
```

---

## Q008: ¿Tabla material_summary_link es necesaria en MVP?
**Contexto:** Enlace entre PostgreSQL (materials) y MongoDB (summary/assessment). ¿Es realmente necesario?

**Opciones:**

### 1. **No Incluir (Simplificar MVP)**
- **Pros:**
  - Más simple, una tabla menos
  - Campo `mongo_document_id` en `assessment` es suficiente
  - Sin overhead de joins
  
- **Contras:**
  - No hay enlace centralizado a `material_summary`
  - Si en el futuro se necesita link a múltiples documentos MongoDB, requiere refactor

### 2. **Incluir (RECOMENDADO para escalabilidad)**
- **Pros:**
  - Tabla de enlace centralizada
  - Facilita queries "obtener todos los documentos MongoDB de un material"
  - Extensible para futuros tipos de documentos (events, analytics)
  
- **Contras:**
  - Tabla adicional
  - JOIN adicional en queries

**Decisión por Defecto:** **Incluir tabla (OPCIONAL en Sprint 01, REQUERIDA en Post-MVP)**

**Justificación:**
- Bajo costo (tabla pequeña, ~10K filas)
- Facilita escalabilidad futura
- Alineado con patrón de integración PostgreSQL-MongoDB

**Implementación:**
```sql
-- Crear tabla (opcional en Sprint 01)
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

-- Uso alternativo sin tabla (usar solo assessment.mongo_document_id)
-- SELECT mongo_document_id FROM assessment WHERE material_id = $1;
```

**Decisión Final:** **Incluir en migración pero marcar como OPCIONAL**

---

## Resumen de Decisiones

| Pregunta | Decisión por Defecto | Sprint de Implementación |
|----------|----------------------|--------------------------|
| Q001: Tipo de ID | UUID v7 | Sprint 01 |
| Q002: Mutabilidad de attempts | Inmutable (Append-Only) | Sprint 01 |
| Q003: Particionamiento | No (Post-MVP si >1M filas) | Post-MVP |
| Q004: Validación time_spent | Check Constraint con EXTRACT | Sprint 01 |
| Q005: Índices de historial | Compuesto + Separados | Sprint 01 |
| Q006: mongo_document_id tipo | VARCHAR(24) | Sprint 01 |
| Q007: idempotency_key | Agregar (NULLABLE en MVP) | Sprint 01 |
| Q008: Tabla material_summary_link | Incluir (OPCIONAL) | Sprint 01 |

---

**Generado con:** Claude Code  
**Sprint:** 01/06  
**Última actualización:** 2025-11-14
