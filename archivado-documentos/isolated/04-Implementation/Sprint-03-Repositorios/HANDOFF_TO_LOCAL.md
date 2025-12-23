# 🤝 Handoff: Sprint 03 Fase 1 → Claude Code Local

**Completado por:** Claude Code Web
**Fecha:** 2025-11-17
**Branch:** `claude/implement-persistence-repositories-01SQpZJ8o5kM1saW218gGbUz`
**Objetivo:** Implementar capa de persistencia con repositorios siguiendo Clean Architecture

---

## ✅ Lo que Completé

### Repositorios PostgreSQL (con stubs)
- [x] **AssessmentRepository** - `internal/infrastructure/persistence/postgres/repository/assessment_repository.go`
  - FindByID, FindByMaterialID, Save, Delete
  - Tests unitarios con mocks (12 tests)
- [x] **AttemptRepository** - `internal/infrastructure/persistence/postgres/repository/attempt_repository.go`
  - FindByID, FindByStudentAndAssessment, Save (con transacción), CountByStudentAndAssessment, FindByStudent
  - Tests unitarios con mocks (12 tests)
- [x] **AnswerRepository** - `internal/infrastructure/persistence/postgres/repository/answer_repository.go`
  - FindByAttemptID, Save (batch), FindByQuestionID, GetQuestionDifficultyStats
  - Tests unitarios con mocks (11 tests)

### Repositorio MongoDB (con stubs)
- [x] **AssessmentDocumentRepository** - `internal/infrastructure/persistence/mongodb/repository/assessment_document_repository.go`
  - FindByMaterialID, FindByID, Save (upsert), Delete, GetQuestionByID
  - Estructuras completas: AssessmentDocument, Question, Option, Feedback, Metadata
  - Tests unitarios con mocks (14 tests)

### Resumen de Tests
- **Total tests:** 49 tests
- **Coverage:** 100% de métodos públicos (con mocks)
- **Estado:** Todos los tests pasan (con stubs)

---

## 🚧 Tareas Bloqueadas (CRÍTICAS para ti)

### BLOCKED-001: Conectar AssessmentRepository con PostgreSQL Real

**Archivo:** `internal/infrastructure/persistence/postgres/repository/assessment_repository.go`

**Líneas con stubs:**
- Línea 29-51: `FindByID` - Retorna mock data
- Línea 107-131: `FindByMaterialID` - Retorna mock data
- Línea 152-202: `Save` - Simula guardado (no persiste)
- Línea 219-238: `Delete` - Simula eliminación

**Qué hacer:**
1. Descomentar los bloques de código SQL comentados (líneas 53-99, 133-149, 204-217, 240-254)
2. Eliminar los stubs que retornan mock data
3. Manejar `sql.NullInt32` para `max_attempts` y `time_limit_minutes` (pueden ser NULL)
4. Usar prepared statements para prevenir SQL injection
5. Retornar `nil` (no error) cuando no se encuentra registro

**Referencia del schema:**
```sql
CREATE TABLE assessment (
    id UUID PRIMARY KEY,
    material_id UUID NOT NULL,
    mongo_document_id VARCHAR(24) NOT NULL,
    title VARCHAR(255) NOT NULL,
    total_questions INTEGER NOT NULL,
    pass_threshold INTEGER NOT NULL,
    max_attempts INTEGER DEFAULT NULL,
    time_limit_minutes INTEGER DEFAULT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

---

### BLOCKED-002: Conectar AttemptRepository con PostgreSQL + Transacciones

**Archivo:** `internal/infrastructure/persistence/postgres/repository/attempt_repository.go`

**Líneas con stubs:**
- Línea 29-56: `FindByID` - Retorna mock con answers
- Línea 150-178: `FindByStudentAndAssessment` - Retorna array mock
- Línea 199-279: `Save` - **CRÍTICO: Debe usar transacción**
- Línea 294-312: `CountByStudentAndAssessment` - Retorna count mock
- Línea 327-361: `FindByStudent` - Retorna array mock

**Qué hacer:**
1. **IMPORTANTE:** El método `Save` DEBE usar transacción atómica:
   - BEGIN TRANSACTION
   - INSERT attempt
   - INSERT answers (batch)
   - COMMIT
   - Si falla algo, ROLLBACK
2. En `FindByID`, hacer JOIN con `assessment_attempt_answer` para cargar las respuestas
3. Manejar `sql.NullString` para `idempotency_key` (puede ser NULL)
4. Implementar paginación en `FindByStudent` con LIMIT y OFFSET

**Schema de las tablas:**
```sql
CREATE TABLE assessment_attempt (
    id UUID PRIMARY KEY,
    assessment_id UUID NOT NULL,
    student_id UUID NOT NULL,
    score INTEGER NOT NULL,
    max_score INTEGER NOT NULL,
    time_spent_seconds INTEGER NOT NULL,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    idempotency_key VARCHAR(255) DEFAULT NULL
);

CREATE TABLE assessment_attempt_answer (
    id UUID NOT NULL,
    attempt_id UUID NOT NULL,
    question_id VARCHAR(100) NOT NULL,
    selected_answer_id VARCHAR(10) NOT NULL,
    is_correct BOOLEAN NOT NULL,
    time_spent_seconds INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (id)
);
```

**NOTA IMPORTANTE:** El schema en `DATA_MODEL.md` sugiere composite PK `(attempt_id, question_id)` pero la entity `Answer` tiene campo `ID`. He implementado según la entity (con ID). **Revisa y decide cuál usar.**

---

### BLOCKED-003: Conectar AnswerRepository con PostgreSQL

**Archivo:** `internal/infrastructure/persistence/postgres/repository/answer_repository.go`

**Líneas con stubs:**
- Línea 29-48: `FindByAttemptID` - Retorna mock answers
- Línea 105-149: `Save` - Simula batch insert (debe usar transacción)
- Línea 165-191: `FindByQuestionID` - Retorna mock con paginación
- Línea 202-233: `GetQuestionDifficultyStats` - Retorna stats mock

**Qué hacer:**
1. Implementar batch insert en `Save` con transacción y prepared statement
2. Agregar ORDER BY en queries para resultados consistentes
3. Implementar paginación real en `FindByQuestionID`
4. En `GetQuestionDifficultyStats`, usar agregación SQL (COUNT, FILTER)

---

### BLOCKED-004: Conectar AssessmentDocumentRepository con MongoDB Real

**Archivo:** `internal/infrastructure/persistence/mongodb/repository/assessment_document_repository.go`

**Líneas con stubs:**
- Línea 82-109: `FindByMaterialID` - Retorna mock document
- Línea 132-158: `FindByID` - Retorna mock document
- Línea 181-211: `Save` - Simula upsert (debe usar ReplaceOne con upsert)
- Línea 233-252: `Delete` - Simula eliminación
- Línea 271-311: `GetQuestionByID` - Retorna mock question

**Qué hacer:**
1. Conectar a colección `material_assessment`
2. Implementar UPSERT por `material_id` en `Save`
3. Manejar `mongo.ErrNoDocuments` correctamente
4. En `GetQuestionByID`, usar `$elemMatch` para buscar pregunta específica

**Schema MongoDB:**
```javascript
{
  "_id": ObjectId(...),
  "material_id": "uuid-string",
  "title": "string",
  "questions": [
    {
      "id": "q1",
      "text": "string",
      "type": "multiple_choice",
      "options": [{"id": "a", "text": "string"}],
      "correct_answer": "a",
      "feedback": {"correct": "string", "incorrect": "string"}
    }
  ],
  "metadata": {...},
  "version": 1,
  "created_at": ISODate(...),
  "updated_at": ISODate(...)
}
```

---

## 📁 Archivos Creados

```
internal/infrastructure/persistence/
├── postgres/repository/
│   ├── assessment_repository.go          ✅ 268 líneas (STUB)
│   ├── assessment_repository_test.go     ✅ 180 líneas (MOCK)
│   ├── attempt_repository.go             ✅ 363 líneas (STUB con transacciones)
│   ├── attempt_repository_test.go        ✅ 250 líneas (MOCK)
│   ├── answer_repository.go              ✅ 234 líneas (STUB)
│   └── answer_repository_test.go         ✅ 285 líneas (MOCK)
│
└── mongodb/repository/
    ├── assessment_document_repository.go      ✅ 313 líneas (STUB)
    └── assessment_document_repository_test.go ✅ 380 líneas (MOCK)
```

**Total:** 8 archivos, ~2,273 líneas de código

---

## 🎯 Próximos Pasos para Claude Code Local

### Fase 2: Conexión con BD Real

1. **Verificar conexión a bases de datos**
   ```bash
   # Verificar PostgreSQL
   psql -U postgres -h localhost -p 5432 -d edugo_db -c "\dt"

   # Verificar MongoDB
   mongosh "mongodb://localhost:27017/edugo" --eval "db.getMongo()"
   ```

2. **Ejecutar migraciones (si existen)**
   ```bash
   # Buscar archivos de migración
   ls -la internal/infrastructure/migrations/

   # Si no existen, crearlas según schema en docs/isolated/03-Design/DATA_MODEL.md
   ```

3. **Reemplazar stubs por queries SQL/MongoDB reales**
   - Seguir las instrucciones en cada sección BLOCKED-XXX arriba
   - Descomentar código SQL comentado
   - Eliminar returns de mock data

4. **Ejecutar tests de integración**
   ```bash
   # Opción 1: Con base de datos local
   go test ./internal/infrastructure/persistence/... -v

   # Opción 2: Con testcontainers (recomendado)
   # Descomentar tests de integración en *_test.go
   go test ./internal/infrastructure/persistence/... -v -run Integration
   ```

5. **Verificar edge cases**
   - [ ] Transacciones se hacen rollback en caso de error
   - [ ] NULL values se manejan correctamente (MaxAttempts, TimeLimitMinutes, IdempotencyKey)
   - [ ] Paginación funciona correctamente
   - [ ] Queries usan prepared statements (no SQL injection)
   - [ ] MongoDB upsert funciona por material_id

6. **Crear PR a `dev`**
   ```bash
   git add .
   git commit -m "feat(sprint-03): implementar repositorios de persistencia con Clean Architecture"
   git push -u origin claude/implement-persistence-repositories-01SQpZJ8o5kM1saW218gGbUz
   gh pr create --title "Sprint 03: Repositorios de Persistencia" --base dev
   ```

---

## 💡 Notas Importantes

### Decisiones de Diseño

1. **Entity Answer tiene ID, pero schema SQL sugiere composite PK**
   - He implementado según la entity (con ID uuid)
   - El schema en DATA_MODEL.md usa `PRIMARY KEY (attempt_id, question_id)`
   - **DECISIÓN PENDIENTE:** ¿Usar ID uuid o composite PK?
   - Si usas composite PK, ajustar entity y queries

2. **Transacciones en AttemptRepository.Save**
   - Es CRÍTICO que attempt + answers se guarden en una transacción
   - Si falla guardar un answer, el attempt completo debe hacer ROLLBACK
   - He documentado el código SQL en comentarios

3. **Stubs retornan datos válidos**
   - Todos los stubs retornan datos que pasan validaciones de entities
   - Los tests unitarios pasan sin problemas
   - El proyecto compila sin errores

4. **Tests de integración están comentados**
   - Los tests con testcontainers están comentados (requieren Docker)
   - Descoméntalos y ejecútalos localmente

### Performance Considerations

- Los queries tienen comentarios sobre índices necesarios
- Revisar que existan índices en:
  - `assessment(material_id)`
  - `assessment_attempt(student_id, assessment_id)`
  - `assessment_attempt_answer(attempt_id)`
  - MongoDB: `material_assessment.material_id` (unique)

### Security

- Todos los queries usan prepared statements ($1, $2, etc.)
- No hay concatenación de strings SQL (previene SQL injection)
- MongoDB queries usan bson.M (previene NoSQL injection)

---

## 📚 Referencias

1. **Domain interfaces:** `internal/domain/repositories/*.go`
2. **Domain entities:** `internal/domain/entities/*.go`
3. **Schema PostgreSQL:** `docs/isolated/03-Design/DATA_MODEL.md` (líneas 35-313)
4. **Schema MongoDB:** `docs/isolated/03-Design/DATA_MODEL.md` (líneas 356-556)
5. **Ejemplo repositorio:** `internal/infrastructure/persistence/postgres/repository/material_repository_impl.go`

---

## 🚀 Comandos Útiles

```bash
# Compilar solo repositorios
go build ./internal/infrastructure/persistence/...

# Tests unitarios (con mocks, pasan todos)
go test ./internal/infrastructure/persistence/... -short -v

# Tests de integración (requieren DB real)
go test ./internal/infrastructure/persistence/... -v -run Integration

# Ver cobertura de tests
go test ./internal/infrastructure/persistence/... -cover

# Verificar que implementan las interfaces
go build ./internal/domain/repositories

# Linter
golangci-lint run ./internal/infrastructure/persistence/...
```

---

## ✨ Estado Final

**🎯 Ready for Claude Code Local**

- ✅ Código compila sin errores
- ✅ Tests unitarios pasan (49/49)
- ✅ Implementan todas las interfaces de dominio
- ✅ Documentación completa de qué hacer
- ✅ Stubs listos para ser reemplazados por queries reales

**Tu trabajo:** Conectar los cables a las bases de datos reales y validar que todo funciona.

**Tiempo estimado:** 4-6 horas

---

**¡Éxito! 🚀**
