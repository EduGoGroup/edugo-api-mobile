# 🗄️ Esquema de Base de Datos

## Visión General

EduGo API Mobile utiliza un enfoque **polyglot persistence** con dos bases de datos:

| Base de Datos | Versión | Propósito |
|---------------|---------|-----------|
| **PostgreSQL** | 16 | Datos relacionales estructurados |
| **MongoDB** | 7.0 | Documentos no estructurados (assessments, IA) |

---

## 📊 PostgreSQL - Diagrama ER

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                           POSTGRESQL SCHEMA                                          │
└─────────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────┐       ┌─────────────────────┐       ┌─────────────────────┐
│       schools       │       │        users        │       │   academic_units    │
├─────────────────────┤       ├─────────────────────┤       ├─────────────────────┤
│ id          UUID PK │◄──────│ school_id   UUID FK │       │ id          UUID PK │
│ name        VARCHAR │       │ id          UUID PK │       │ school_id   UUID FK │
│ code        VARCHAR │       │ email       VARCHAR │       │ name        VARCHAR │
│ created_at  TIMESTAMP│      │ role        VARCHAR │       │ grade       VARCHAR │
│ updated_at  TIMESTAMP│      │ first_name  VARCHAR │       │ created_at  TIMESTAMP│
└─────────────────────┘       │ last_name   VARCHAR │       └─────────────────────┘
                              │ created_at  TIMESTAMP│                │
                              │ updated_at  TIMESTAMP│                │
                              └──────────┬──────────┘                │
                                         │                            │
              ┌──────────────────────────┼────────────────────────────┘
              │                          │
              ▼                          ▼
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                                    materials                                         │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ id                      UUID        PRIMARY KEY                                      │
│ school_id               UUID        NOT NULL  FK → schools(id)                       │
│ uploaded_by_teacher_id  UUID        NOT NULL  FK → users(id)                         │
│ academic_unit_id        UUID        NULLABLE  FK → academic_units(id)                │
│ title                   VARCHAR(200) NOT NULL                                        │
│ description             TEXT        NULLABLE                                         │
│ subject                 VARCHAR(100) NULLABLE                                        │
│ grade                   VARCHAR(50)  NULLABLE                                        │
│ file_url                VARCHAR(500) NOT NULL                                        │
│ file_type               VARCHAR(100) NOT NULL                                        │
│ file_size_bytes         BIGINT      NOT NULL                                         │
│ status                  VARCHAR(50)  NOT NULL  DEFAULT 'uploaded'                    │
│ processing_status       VARCHAR(50)  NULLABLE                                        │
│ processing_started_at   TIMESTAMP   NULLABLE                                         │
│ processing_completed_at TIMESTAMP   NULLABLE                                         │
│ is_public               BOOLEAN     NOT NULL  DEFAULT false                          │
│ created_at              TIMESTAMP   NOT NULL  DEFAULT NOW()                          │
│ updated_at              TIMESTAMP   NOT NULL  DEFAULT NOW()                          │
│ deleted_at              TIMESTAMP   NULLABLE  (soft delete)                          │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ INDEXES:                                                                             │
│  • idx_materials_school_id                                                           │
│  • idx_materials_teacher_id                                                          │
│  • idx_materials_status                                                              │
│  • idx_materials_created_at                                                          │
└─────────────────────────────────────────────────────────────────────────────────────┘
              │
              │ 1:N
              ▼
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              material_versions                                       │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ id              UUID        PRIMARY KEY                                              │
│ material_id     UUID        NOT NULL  FK → materials(id)                             │
│ version_number  INTEGER     NOT NULL                                                 │
│ title           VARCHAR(200) NOT NULL                                                │
│ content_url     VARCHAR(500) NOT NULL                                                │
│ changed_by      UUID        NOT NULL  FK → users(id)                                 │
│ created_at      TIMESTAMP   NOT NULL  DEFAULT NOW()                                  │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ UNIQUE: (material_id, version_number)                                                │
└─────────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────────┐
│                                    progress                                          │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ id                    UUID        PRIMARY KEY                                        │
│ user_id               UUID        NOT NULL  FK → users(id)                           │
│ material_id           UUID        NOT NULL  FK → materials(id)                       │
│ progress_percentage   INTEGER     NOT NULL  CHECK (0-100)                            │
│ last_page             INTEGER     NULLABLE                                           │
│ time_spent_seconds    INTEGER     NOT NULL  DEFAULT 0                                │
│ started_at            TIMESTAMP   NOT NULL  DEFAULT NOW()                            │
│ last_accessed_at      TIMESTAMP   NOT NULL  DEFAULT NOW()                            │
│ completed_at          TIMESTAMP   NULLABLE                                           │
│ created_at            TIMESTAMP   NOT NULL  DEFAULT NOW()                            │
│ updated_at            TIMESTAMP   NOT NULL  DEFAULT NOW()                            │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ UNIQUE: (user_id, material_id)  -- Para UPSERT idempotente                           │
│ INDEXES:                                                                             │
│  • idx_progress_user_id                                                              │
│  • idx_progress_material_id                                                          │
│  • idx_progress_last_accessed_at  -- Para estadísticas de usuarios activos           │
└─────────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────────┐
│                                  assessments                                         │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ id                    UUID        PRIMARY KEY                                        │
│ material_id           UUID        NOT NULL  FK → materials(id)                       │
│ title                 VARCHAR(200) NOT NULL                                          │
│ pass_threshold        INTEGER     NOT NULL  DEFAULT 60                               │
│ max_attempts          INTEGER     NULLABLE                                           │
│ time_limit_minutes    INTEGER     NULLABLE                                           │
│ is_active             BOOLEAN     NOT NULL  DEFAULT true                             │
│ created_at            TIMESTAMP   NOT NULL  DEFAULT NOW()                            │
│ updated_at            TIMESTAMP   NOT NULL  DEFAULT NOW()                            │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ UNIQUE: (material_id)  -- Un assessment por material                                 │
└─────────────────────────────────────────────────────────────────────────────────────┘
              │
              │ 1:N
              ▼
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                                   questions                                          │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ id                UUID        PRIMARY KEY                                            │
│ assessment_id     UUID        NOT NULL  FK → assessments(id)                         │
│ question_text     TEXT        NOT NULL                                               │
│ question_type     VARCHAR(50) NOT NULL  -- multiple_choice, true_false, short_answer │
│ correct_answer    TEXT        NOT NULL                                               │
│ explanation       TEXT        NULLABLE                                               │
│ difficulty_level  VARCHAR(20) NOT NULL  DEFAULT 'medium'                             │
│ order_index       INTEGER     NOT NULL                                               │
│ created_at        TIMESTAMP   NOT NULL  DEFAULT NOW()                                │
└─────────────────────────────────────────────────────────────────────────────────────┘
              │
              │ 1:N
              ▼
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                                answer_options                                        │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ id              UUID        PRIMARY KEY                                              │
│ question_id     UUID        NOT NULL  FK → questions(id)                             │
│ option_text     TEXT        NOT NULL                                                 │
│ is_correct      BOOLEAN     NOT NULL  DEFAULT false                                  │
│ order_index     INTEGER     NOT NULL                                                 │
└─────────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              assessment_attempts                                     │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ id                    UUID        PRIMARY KEY                                        │
│ assessment_id         UUID        NOT NULL  FK → assessments(id)                     │
│ student_id            UUID        NOT NULL  FK → users(id)                           │
│ score                 INTEGER     NOT NULL                                           │
│ max_score             INTEGER     NOT NULL                                           │
│ passed                BOOLEAN     NOT NULL                                           │
│ time_spent_seconds    INTEGER     NOT NULL                                           │
│ started_at            TIMESTAMP   NOT NULL                                           │
│ completed_at          TIMESTAMP   NOT NULL                                           │
│ created_at            TIMESTAMP   NOT NULL  DEFAULT NOW()                            │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ INDEXES:                                                                             │
│  • idx_attempts_assessment_id                                                        │
│  • idx_attempts_student_id                                                           │
│  • idx_attempts_completed_at                                                         │
└─────────────────────────────────────────────────────────────────────────────────────┘
              │
              │ 1:N
              ▼
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                               attempt_answers                                        │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ id                    UUID        PRIMARY KEY                                        │
│ attempt_id            UUID        NOT NULL  FK → assessment_attempts(id)             │
│ question_id           UUID        NOT NULL  FK → questions(id)                       │
│ selected_answer_id    UUID        NULLABLE  FK → answer_options(id)                  │
│ answer_text           TEXT        NULLABLE  -- Para respuestas abiertas              │
│ is_correct            BOOLEAN     NOT NULL                                           │
│ time_spent_seconds    INTEGER     NOT NULL  DEFAULT 0                                │
└─────────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────────┐
│                               refresh_tokens                                         │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ id              UUID        PRIMARY KEY                                              │
│ user_id         UUID        NOT NULL  FK → users(id)                                 │
│ token           VARCHAR(500) NOT NULL  UNIQUE                                        │
│ expires_at      TIMESTAMP   NOT NULL                                                 │
│ revoked         BOOLEAN     NOT NULL  DEFAULT false                                  │
│ created_at      TIMESTAMP   NOT NULL  DEFAULT NOW()                                  │
└─────────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────────┐
│                                login_attempts                                        │
├─────────────────────────────────────────────────────────────────────────────────────┤
│ id              UUID        PRIMARY KEY                                              │
│ user_id         UUID        NULLABLE  FK → users(id)                                 │
│ email           VARCHAR(255) NOT NULL                                                │
│ ip_address      VARCHAR(45)  NOT NULL                                                │
│ user_agent      TEXT        NULLABLE                                                 │
│ success         BOOLEAN     NOT NULL                                                 │
│ failure_reason  VARCHAR(100) NULLABLE                                                │
│ created_at      TIMESTAMP   NOT NULL  DEFAULT NOW()                                  │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 📄 MongoDB - Colecciones

### Colección: `assessments` (Quizzes generados por IA)

```javascript
// Collection: assessments
// Almacena los cuestionarios generados automáticamente por el Worker de IA
{
  "_id": ObjectId("..."),
  "material_id": "550e8400-e29b-41d4-a716-446655440000",  // UUID del material
  "questions": [
    {
      "id": "q1",
      "question_text": "¿Cuál es la derivada de x²?",
      "question_type": "multiple_choice",  // multiple_choice, true_false, short_answer
      "options": [
        { "id": "a", "text": "x" },
        { "id": "b", "text": "2x" },
        { "id": "c", "text": "x²" },
        { "id": "d", "text": "2" }
      ],
      "correct_answer": "b",  // ID de la opción correcta
      "explanation": "La derivada de x^n es n*x^(n-1), entonces d/dx(x²) = 2x",
      "difficulty_level": "medium"  // easy, medium, hard
    },
    {
      "id": "q2",
      "question_text": "La integral es la operación inversa de la derivada",
      "question_type": "true_false",
      "correct_answer": true,
      "explanation": "El teorema fundamental del cálculo establece esta relación"
    }
  ],
  "total_questions": 10,
  "estimated_time_minutes": 15,
  "pass_threshold": 70,
  "created_at": ISODate("2024-12-06T10:00:00Z"),
  "updated_at": ISODate("2024-12-06T10:00:00Z")
}

// Indexes
db.assessments.createIndex({ "material_id": 1 }, { unique: true })
```

### Colección: `assessment_attempts` (Intentos de usuarios)

```javascript
// Collection: assessment_attempts
// Almacena cada intento de un usuario en un cuestionario
{
  "_id": ObjectId("..."),
  "material_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "660e8400-e29b-41d4-a716-446655440001",
  "answers": {
    "q1": "b",      // question_id -> respuesta
    "q2": true,
    "q3": "Paris"
  },
  "score": 85.5,
  "attempted_at": ISODate("2024-12-06T14:30:00Z")
}

// Indexes
db.assessment_attempts.createIndex({ "material_id": 1, "user_id": 1 })
db.assessment_attempts.createIndex({ "user_id": 1, "attempted_at": -1 })
```

### Colección: `assessment_results` (Resultados finales)

```javascript
// Collection: assessment_results
// Almacena resultados detallados con feedback
{
  "_id": ObjectId("..."),
  "assessment_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "660e8400-e29b-41d4-a716-446655440001",
  "score": 80.0,
  "total_questions": 10,
  "correct_answers": 8,
  "feedback": [
    {
      "question_id": "q1",
      "is_correct": true,
      "user_answer": "b",
      "correct_answer": "b",
      "explanation": "Correcto! La derivada de x² es 2x"
    },
    {
      "question_id": "q2",
      "is_correct": false,
      "user_answer": "a",
      "correct_answer": "c",
      "explanation": "Incorrecto. La respuesta correcta es..."
    }
  ],
  "submitted_at": ISODate("2024-12-06T14:35:00Z")
}

// Indexes - UNIQUE para evitar duplicados
db.assessment_results.createIndex(
  { "assessment_id": 1, "user_id": 1 },
  { unique: true }
)
```

### Colección: `summaries` (Resúmenes generados por IA)

```javascript
// Collection: summaries
// Almacena resúmenes de materiales generados por el Worker de IA
{
  "_id": ObjectId("..."),
  "material_id": "550e8400-e29b-41d4-a716-446655440000",
  "summary_text": "Este documento trata sobre los fundamentos del cálculo diferencial...",
  "key_points": [
    "Concepto de límite",
    "Definición de derivada",
    "Reglas de derivación",
    "Aplicaciones prácticas"
  ],
  "word_count": 250,
  "reading_time_minutes": 2,
  "language": "es",
  "generated_by": "gpt-4",
  "created_at": ISODate("2024-12-06T10:05:00Z")
}

// Indexes
db.summaries.createIndex({ "material_id": 1 }, { unique: true })
```

---

## 🔗 Diagrama de Relaciones Completo

```
                                    PostgreSQL
    ┌─────────────────────────────────────────────────────────────────────┐
    │                                                                     │
    │   schools ◄─────────────────┬──────────────────► academic_units    │
    │      │                      │                          │           │
    │      │                      │                          │           │
    │      ▼                      ▼                          ▼           │
    │   users ◄────────────── materials ─────────► material_versions     │
    │      │                      │                                      │
    │      │         ┌────────────┼────────────┐                         │
    │      │         │            │            │                         │
    │      │         ▼            ▼            ▼                         │
    │      │    progress    assessments   (to MongoDB)                   │
    │      │                     │                                       │
    │      │                     ▼                                       │
    │      │               questions                                     │
    │      │                     │                                       │
    │      │                     ▼                                       │
    │      │             answer_options                                  │
    │      │                                                             │
    │      ├─────► assessment_attempts                                   │
    │      │              │                                              │
    │      │              ▼                                              │
    │      │       attempt_answers                                       │
    │      │                                                             │
    │      ├─────► refresh_tokens                                        │
    │      │                                                             │
    │      └─────► login_attempts                                        │
    │                                                                     │
    └─────────────────────────────────────────────────────────────────────┘

                                    MongoDB
    ┌─────────────────────────────────────────────────────────────────────┐
    │                                                                     │
    │   materials (PostgreSQL) ──────────► assessments (IA generated)    │
    │          │                                  │                       │
    │          │                                  ▼                       │
    │          │                       assessment_attempts               │
    │          │                                  │                       │
    │          │                                  ▼                       │
    │          │                        assessment_results               │
    │          │                                                          │
    │          └────────────────────────► summaries (IA generated)       │
    │                                                                     │
    └─────────────────────────────────────────────────────────────────────┘
```

---

## 📋 Valores de Enums

### Material Status
```
uploaded     → Material creado, archivo subido a S3
processing   → Worker procesando (generando summary/quiz)
ready        → Listo para uso
failed       → Error en procesamiento
```

### Processing Status
```
pending      → En cola de procesamiento
in_progress  → Siendo procesado por Worker
completed    → Procesamiento exitoso
failed       → Error en procesamiento
```

### Question Type
```
multiple_choice  → Selección múltiple (una respuesta correcta)
true_false       → Verdadero/Falso
short_answer     → Respuesta corta (texto libre)
```

### Difficulty Level
```
easy    → Pregunta fácil
medium  → Pregunta de dificultad media
hard    → Pregunta difícil
```

### User Role
```
admin      → Administrador del sistema
teacher    → Docente
student    → Estudiante
```

---

## 🔧 Scripts de Migración

### Crear tablas PostgreSQL (ejemplo)

```sql
-- Materials table
CREATE TABLE IF NOT EXISTS materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id),
    uploaded_by_teacher_id UUID NOT NULL REFERENCES users(id),
    academic_unit_id UUID REFERENCES academic_units(id),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    subject VARCHAR(100),
    grade VARCHAR(50),
    file_url VARCHAR(500) NOT NULL,
    file_type VARCHAR(100) NOT NULL,
    file_size_bytes BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'uploaded',
    processing_status VARCHAR(50),
    processing_started_at TIMESTAMP,
    processing_completed_at TIMESTAMP,
    is_public BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_materials_school_id ON materials(school_id);
CREATE INDEX idx_materials_teacher_id ON materials(uploaded_by_teacher_id);
CREATE INDEX idx_materials_status ON materials(status);

-- Progress table with UPSERT support
CREATE TABLE IF NOT EXISTS progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    material_id UUID NOT NULL REFERENCES materials(id),
    progress_percentage INTEGER NOT NULL CHECK (progress_percentage >= 0 AND progress_percentage <= 100),
    last_page INTEGER,
    time_spent_seconds INTEGER NOT NULL DEFAULT 0,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_accessed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, material_id)
);

CREATE INDEX idx_progress_user_id ON progress(user_id);
CREATE INDEX idx_progress_material_id ON progress(material_id);
CREATE INDEX idx_progress_last_accessed ON progress(last_accessed_at);
```

### Crear índices MongoDB

```javascript
// assessments collection
db.assessments.createIndex({ "material_id": 1 }, { unique: true });

// assessment_attempts collection
db.assessment_attempts.createIndex({ "material_id": 1, "user_id": 1 });
db.assessment_attempts.createIndex({ "user_id": 1, "attempted_at": -1 });

// assessment_results collection
db.assessment_results.createIndex(
  { "assessment_id": 1, "user_id": 1 },
  { unique: true }
);

// summaries collection
db.summaries.createIndex({ "material_id": 1 }, { unique: true });
```

---

## 📊 Queries Frecuentes

### Obtener progreso de un usuario en todos sus materiales
```sql
SELECT m.id, m.title, p.progress_percentage, p.last_accessed_at
FROM materials m
LEFT JOIN progress p ON p.material_id = m.id AND p.user_id = :user_id
WHERE m.deleted_at IS NULL
ORDER BY p.last_accessed_at DESC NULLS LAST;
```

### Contar usuarios activos (últimos 30 días)
```sql
SELECT COUNT(DISTINCT user_id) as active_users
FROM progress
WHERE last_accessed_at >= NOW() - INTERVAL '30 days';
```

### Promedio de progreso global
```sql
SELECT AVG(progress_percentage) as avg_progress
FROM progress;
```

### Obtener assessment con mejor score de un usuario (MongoDB)
```javascript
db.assessment_attempts.find({
  material_id: "...",
  user_id: "..."
}).sort({ score: -1 }).limit(1);
```
