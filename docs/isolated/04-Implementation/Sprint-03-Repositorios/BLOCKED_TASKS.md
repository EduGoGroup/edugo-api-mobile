# 🚧 Tareas Bloqueadas - Sprint 03 Fase 1

**Autor:** Claude Code Web
**Fecha:** 2025-11-17
**Razón:** No tengo acceso a conexiones de base de datos reales (PostgreSQL y MongoDB)

---

## Resumen

He implementado **4 repositorios con stubs** (código que compila pero retorna datos mockeados).

**Tu trabajo:** Reemplazar los stubs por queries SQL/MongoDB reales.

---

## BLOCKED-001: PostgreSQL - AssessmentRepository

**Archivo:** `internal/infrastructure/persistence/postgres/repository/assessment_repository.go`

| Método | Línea | Stub Actual | Qué Hacer |
|--------|-------|-------------|-----------|
| `FindByID` | 29-51 | Retorna mock Assessment | Descomentar query SQL (líneas 53-99) |
| `FindByMaterialID` | 107-131 | Retorna mock Assessment | Descomentar query SQL (líneas 133-149) |
| `Save` | 152-202 | Simula guardado | Descomentar UPSERT SQL (líneas 204-217) |
| `Delete` | 219-238 | Simula eliminación | Descomentar DELETE SQL (líneas 240-254) |

**Complejidad:** Baja (queries simples con prepared statements)

**Tiempo estimado:** 1 hora

---

## BLOCKED-002: PostgreSQL - AttemptRepository ⚠️ CRÍTICO

**Archivo:** `internal/infrastructure/persistence/postgres/repository/attempt_repository.go`

| Método | Línea | Stub Actual | Qué Hacer |
|--------|-------|-------------|-----------|
| `FindByID` | 29-56 | Retorna mock Attempt + Answers | JOIN con `assessment_attempt_answer` (líneas 58-157) |
| `FindByStudentAndAssessment` | 150-178 | Retorna array mock | Query con ORDER BY (líneas 180-202) |
| `Save` | 199-279 | Simula guardado | **TRANSACCIÓN atómica** (líneas 281-327) |
| `CountByStudentAndAssessment` | 294-312 | Retorna count mock | Query COUNT (líneas 314-331) |
| `FindByStudent` | 327-361 | Retorna array mock | Query con LIMIT/OFFSET (líneas 363-395) |

**⚠️ CRÍTICO:** El método `Save` DEBE usar transacción:
```sql
BEGIN;
INSERT INTO assessment_attempt (...) VALUES (...);
INSERT INTO assessment_attempt_answer (...) VALUES (...); -- Batch
COMMIT; -- O ROLLBACK si falla
```

**Complejidad:** Alta (transacciones + JOIN)

**Tiempo estimado:** 2-3 horas

---

## BLOCKED-003: PostgreSQL - AnswerRepository

**Archivo:** `internal/infrastructure/persistence/postgres/repository/answer_repository.go`

| Método | Línea | Stub Actual | Qué Hacer |
|--------|-------|-------------|-----------|
| `FindByAttemptID` | 29-48 | Retorna array mock | Query con ORDER BY (líneas 50-99) |
| `Save` | 105-149 | Simula batch insert | Transacción + prepared statement (líneas 151-192) |
| `FindByQuestionID` | 165-191 | Retorna array mock | Query con LIMIT/OFFSET (líneas 193-225) |
| `GetQuestionDifficultyStats` | 202-233 | Retorna stats mock | Agregación SQL (líneas 235-258) |

**Complejidad:** Media (batch insert + agregación)

**Tiempo estimado:** 1-2 horas

---

## BLOCKED-004: MongoDB - AssessmentDocumentRepository

**Archivo:** `internal/infrastructure/persistence/mongodb/repository/assessment_document_repository.go`

| Método | Línea | Stub Actual | Qué Hacer |
|--------|-------|-------------|-----------|
| `FindByMaterialID` | 82-109 | Retorna mock document | FindOne con filter (líneas 111-130) |
| `FindByID` | 132-158 | Retorna mock document | FindOne por ObjectID (líneas 160-181) |
| `Save` | 181-211 | Simula upsert | ReplaceOne con upsert=true (líneas 213-233) |
| `Delete` | 233-252 | Simula eliminación | DeleteOne (líneas 254-272) |
| `GetQuestionByID` | 271-311 | Retorna mock question | Aggregation con $elemMatch (líneas 313-347) |

**Complejidad:** Media (upsert + projection)

**Tiempo estimado:** 1-2 horas

---

## Checklist de Validación

Cuando termines de conectar con BD real, verifica:

### PostgreSQL
- [ ] Prepared statements usan `$1, $2, ...` (no concatenación)
- [ ] NULL values se manejan con `sql.NullInt32`, `sql.NullString`
- [ ] Transacciones hacen ROLLBACK en caso de error
- [ ] Queries retornan `nil` (no error) cuando no hay resultados
- [ ] Paginación funciona con LIMIT y OFFSET

### MongoDB
- [ ] Upsert usa `material_id` como filter (no `_id`)
- [ ] Se maneja `mongo.ErrNoDocuments` correctamente
- [ ] ObjectID se valida antes de queries
- [ ] Projection con `$elemMatch` funciona para preguntas

### Tests
- [ ] Tests unitarios siguen pasando (49 tests)
- [ ] Tests de integración ejecutados con testcontainers
- [ ] Coverage >= 80%

---

## Estructura de Archivos con Stubs

```
internal/infrastructure/persistence/
├── postgres/repository/
│   ├── assessment_repository.go       (268 líneas, 4 stubs)
│   ├── attempt_repository.go          (363 líneas, 5 stubs, 1 CRÍTICO)
│   └── answer_repository.go           (234 líneas, 4 stubs)
│
└── mongodb/repository/
    └── assessment_document_repository.go (313 líneas, 5 stubs)
```

---

## Comandos de Verificación

```bash
# Compilar (debe pasar)
go build ./internal/infrastructure/persistence/...

# Tests con stubs (debe pasar - 49 tests)
go test ./internal/infrastructure/persistence/... -short -v

# Tests de integración (después de conectar BD real)
go test ./internal/infrastructure/persistence/... -v -run Integration

# Ver qué archivos tienen stubs (buscar "STUB")
grep -r "STUB" internal/infrastructure/persistence/
```

---

## Decisiones Pendientes

### 1. Schema de `assessment_attempt_answer`

**Conflicto:**
- **Entity `Answer`:** Tiene campo `ID uuid.UUID` (línea 14 de `internal/domain/entities/answer.go`)
- **Schema SQL en docs:** Usa composite PK `(attempt_id, question_id)` (línea 264 de `DATA_MODEL.md`)

**Implementación actual:** He usado `ID uuid.UUID` según la entity.

**Acción requerida:**
- Si prefieres composite PK: Modificar entity `Answer` y queries SQL
- Si prefieres ID uuid: Actualizar schema en `DATA_MODEL.md`

### 2. Índices de Base de Datos

**Verifica que existan estos índices:**

PostgreSQL:
```sql
CREATE INDEX idx_assessment_material_id ON assessment(material_id);
CREATE INDEX idx_attempt_student_assessment ON assessment_attempt(student_id, assessment_id);
CREATE INDEX idx_answer_attempt_id ON assessment_attempt_answer(attempt_id);
```

MongoDB:
```javascript
db.material_assessment.createIndex({"material_id": 1}, {unique: true});
```

---

## Próximos Pasos

1. Leer `HANDOFF_TO_LOCAL.md` para detalles completos
2. Verificar conexión a PostgreSQL y MongoDB
3. Ejecutar migraciones (si existen)
4. Reemplazar stubs por queries reales (4-6 horas)
5. Ejecutar tests de integración
6. Crear PR a `dev`

---

**Estado:** ✅ Listo para Claude Code Local

**Riesgo:** 🟢 Bajo (código ya compila y tests pasan con stubs)

**Prioridad:** 🔴 Alta (bloqueante para Sprint 04 - Services y API)
