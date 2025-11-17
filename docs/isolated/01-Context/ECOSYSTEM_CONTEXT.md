# ECOSYSTEM CONTEXT - API Mobile

## Posición en EduGo

**Rol:** Microservicio Frontend-facing (cara visible para app móvil)  
**Interacción:** Hub central de evaluaciones, nexo entre móvil y backend

---

## Mapa de Ecosistema

```
┌────────────────────────────────────────────────────────────┐
│                     APLICACIÓN MÓVIL                       │
│          (Android/iOS - Usuarios finales)                  │
└────────────────────┬─────────────────────────────────────┘
                     │ HTTP REST
                     ▼
        ┌────────────────────────────┐
        │   API MOBILE (Puerto 8080) │  ◄─── ESTE PROYECTO
        │  - Evaluaciones            │
        │  - Preguntas               │
        │  - Asignaciones            │
        │  - Resultados              │
        └────┬───────────┬──────┬────┘
             │           │      │
             │ TCP       │      │ HTTP
             ▼           ▼      ▼
        ┌─────────┐  ┌──────────────────┐  ┌──────────────┐
        │PostgreSQL   │    SHARED        │  │    WORKER    │
        │ (Relacional)│  v1.3.0+         │  │  (Async IA)  │
        │             │  - Logger        │  │              │
        │             │  - Database      │  │  - PDF → Txt │
        │             │  - Auth          │  │  - Resúmenes │
        │             │  - Messaging     │  │  - Quizzes   │
        └─────────────┴──────────────────┘  └──────────────┘
                          │
                          │ (consume/produce)
                          ▼
                    ┌──────────────┐
                    │   RabbitMQ   │
                    │   (Mensajes) │
                    └──────┬───────┘
                           │
                           ├──────────┬──────────┐
                           ▼          ▼          ▼
                        ┌─────┐  ┌────────┐  ┌──────┐
                        │ API │  │MongoDB │  │Worker│
                        │Admin│  │ (Docs) │  │      │
                        └─────┘  └────────┘  └──────┘

┌─────────────────────────────────────────────────────────────┐
│              DEV ENVIRONMENT (Docker Compose)               │
│  - Orquesta todos los servicios en desarrollo              │
│  - Gestiona volúmenes, redes, perfiles                     │
└─────────────────────────────────────────────────────────────┘
```

---

## Interacciones con Otros Servicios

### 1. Integración con SHARED (v1.3.0+)

**Módulos utilizados:**

#### Logger
```go
import "github.com/EduGoGroup/edugo-shared/logger"

// En API Mobile:
logger.Info("Evaluación creada", "evaluation_id", eval.ID)
logger.Error("Error al calcular puntuación", "error", err)
```

#### Database (PostgreSQL)
```go
import "github.com/EduGoGroup/edugo-shared/database"

// Conexión centralizada
db := database.GetPostgres()
var evaluation Evaluation
db.First(&evaluation, id)
```

#### Auth (JWT)
```go
import "github.com/EduGoGroup/edugo-shared/auth"

// Validación de tokens
middleware := auth.NewJWTMiddleware()
// En Gin handler
c.Use(middleware.Validate())
```

#### Messaging (RabbitMQ)
```go
import "github.com/EduGoGroup/edugo-shared/messaging"

// Publicar evento: solicitar quiz automático
publisher := messaging.NewPublisher()
publisher.Publish("evaluation.quiz.generate", payload)

// Consumir evento: quizzes generados
consumer := messaging.NewConsumer()
consumer.Subscribe("evaluation.quiz.generated", handler)
```

**Flujo de versionamiento:**
1. SHARED lanza versión v1.3.0 desde rama dev
2. API Mobile actualiza `go.mod` con nueva versión
3. API Mobile debe testear compatibilidad
4. Commit en rama feature de API Mobile

---

### 2. Integración con WORKER

**Comunicación vía RabbitMQ:**

#### Request (API Mobile → Worker)
```json
{
  "request_id": "req-12345",
  "type": "generate_assessment",
  "material_id": 42,
  "material_path": "s3://bucket/material-42.pdf",
  "config": {
    "num_questions": 10,
    "difficulty": "medium",
    "language": "es"
  },
  "timestamp": "2025-11-15T10:30:00Z"
}
```

**Routing:** Exchange `assessment.requests` → Queue `worker.assessment.requests`

#### Response (Worker → API Mobile)
```json
{
  "request_id": "req-12345",
  "status": "success",
  "evaluation_id": 1001,
  "questions_generated": 10,
  "timestamp": "2025-11-15T10:45:00Z",
  "generated_data": {
    "questions": [...],
    "metadata": {...}
  }
}
```

**Routing:** Exchange `assessment.responses` → Queue `api-mobile.assessment.responses`

**Casos de uso:**
1. Docente sube material (PDF)
2. API Mobile solicita generación de quiz a Worker
3. Worker procesa con OpenAI, genera preguntas
4. Worker envía respuesta con preguntas
5. API Mobile persiste preguntas en PostgreSQL
6. API Mobile notifica a cliente (webhook o polling)

---

### 3. Integración con API ADMIN

**Comunicación:** HTTP REST + Shared Context

#### Queries hacia API Admin
```go
// API Mobile necesita info de escuela/profesor
// que está en API Admin (jerarquía académica)
client := http.NewClient()
resp := client.Get("http://api-admin:8081/api/v1/schools/{school_id}")
```

#### Datos compartidos vía PostgreSQL
```sql
-- Tabla teachers (ambos APIs acceden)
SELECT * FROM teachers WHERE school_id = 1;

-- Tabla students
SELECT * FROM students WHERE school_id = 1;

-- Tabla academic_units (jerarquía)
SELECT * FROM academic_units WHERE school_id = 1;
```

#### Contexto compartido
```go
// Ambos APIs usan mismo contexto de autenticación
// desde SHARED
import "github.com/EduGoGroup/edugo-shared/auth"

// Extracto de headers con user_id, school_id
ctx := context.WithValue(r.Context(), "user_id", 42)
```

---

### 4. Bases de Datos Compartidas

#### PostgreSQL (datos relacionales persistentes)

**Tablas de API Mobile:**
```
├─ evaluations           (evaluaciones creadas)
├─ questions             (preguntas de evaluaciones)
├─ question_options      (opciones de respuesta)
├─ evaluation_assignments (asignación a estudiantes)
└─ answer_drafts         (respuestas guardadas localmente)
```

**Tablas compartidas con otros servicios:**
```
├─ users                 (usuarios del sistema)
├─ schools               (escuelas)
├─ academic_units        (jerarquía: facultades, departamentos)
├─ teachers              (docentes - referencia)
├─ students              (estudiantes - referencia)
├─ materials             (materiales educativos)
└─ enrollments           (inscripción de estudiantes)
```

#### MongoDB (almacenamiento documentos)

**Colecciones de API Mobile:**
```
├─ evaluation_results    (resultados de evaluaciones, scoring)
├─ evaluation_answers    (respuestas detalladas)
└─ evaluation_audit      (auditoría de cambios)
```

**Colecciones compartidas:**
```
├─ material_summary      (resúmenes generados por Worker)
├─ material_assessment   (quizzes generados por Worker)
└─ material_event        (eventos de procesamiento)
```

---

## Flujos Inter-microservicios

### Flujo A: Generar Quiz Automático (API Mobile + Worker)

```
Paso 1: Docente en app móvil
├─ Selecciona material educativo
├─ Hace click en "Generar Quiz"
└─ API Mobile POST /api/v1/evaluations/material/{id}/generate-quiz

Paso 2: API Mobile
├─ Valida que material existe (PostgreSQL)
├─ Obtiene path de S3 del material
├─ Publica mensaje a RabbitMQ:
│  └─ Exchange: "assessment.requests"
│     Queue: "worker.assessment.requests"
│     Payload: {request_id, material_id, material_path, config}
├─ Guarda request en caché (Redis o PostgeSQL)
└─ Retorna 202 Accepted con request_id

Paso 3: Worker consume mensaje
├─ Descarga material de S3
├─ Procesa con OpenAI GPT-4
├─ Genera preguntas y opciones
├─ Guarda en MongoDB (evaluation_assessment)
└─ Publica respuesta a RabbitMQ:
   └─ Exchange: "assessment.responses"
      Queue: "api-mobile.assessment.responses"
      Payload: {request_id, questions, evaluation_id}

Paso 4: API Mobile consume respuesta
├─ Recibe mensaje de Worker
├─ Persiste preguntas en PostgreSQL
├─ Crea evaluation object
├─ Marca request como completado
└─ Notifica a app móvil (websocket o polling)
```

### Flujo B: Enviar Evaluación (API Mobile + PostgreSQL)

```
Paso 1: Estudiante en app móvil
├─ Completa evaluación
├─ Hace click en "Enviar"
└─ API Mobile POST /api/v1/evaluations/{id}/submit

Paso 2: API Mobile (Validación)
├─ Obtiene evaluation de PostgreSQL
├─ Valida que no esté cerrada
├─ Obtiene assignment del estudiante
├─ Valida que no haya sido completada
└─ Procede a validación de respuestas

Paso 3: API Mobile (Cálculo)
├─ Itera cada respuesta del estudiante
├─ Compara con respuestas correctas
├─ Calcula puntos por pregunta
├─ Calcula total score y porcentaje
└─ Valida reglas de scoring

Paso 4: API Mobile (Persistencia)
├─ Guarda respuestas en PostgreSQL (answer_drafts→final)
├─ Guarda resultados en MongoDB (evaluation_results)
├─ Actualiza assignment status → "submitted"
├─ Publica evento a RabbitMQ (auditoría)
└─ Retorna resultado con feedback

Paso 5: Docente en API Admin
├─ Consulta GET /api/v1/evaluations/{id}/results
├─ Recibe resultados consolidados
└─ Decide si dar feedback adicional
```

### Flujo C: Sincronización de Contexto (SHARED)

```
Paso 1: Autenticación
├─ App móvil envía JWT en header Authorization
├─ API Mobile middleware (de SHARED) valida token
├─ Extrae user_id, school_id, roles
└─ Inyecta en contexto de request

Paso 2: Logging centralizado
├─ API Mobile usa logger de SHARED
├─ Todos los logs van a salida estándar/Datadog
├─ User_id y request_id se pasan en contexto
└─ Logs correlacionados entre servicios

Paso 3: Gestión de errores
├─ API Mobile captura panics
├─ Convierte a HTTP errors usando SHARED
├─ Retorna estructura estándar de error
└─ Logs incluyen stack trace

Paso 4: Timeout global
├─ API Mobile hereda timeout de SHARED (default 30s)
├─ Aplica a BD, RabbitMQ, HTTP calls
├─ Cancela contexto si se excede
└─ Retorna 504 Gateway Timeout
```

---

## Dependencias Directas

| Servicio | Versión | Tipo | Críticidad |
|----------|---------|------|-----------|
| SHARED | v1.3.0+ | Librería Go | CRÍTICA |
| PostgreSQL | 15+ | Base datos | CRÍTICA |
| MongoDB | 7.0+ | Base datos | ALTA |
| RabbitMQ | 3.12+ | Message broker | ALTA |
| WORKER | Latest | Microservicio | MEDIA |
| API ADMIN | Latest | Microservicio | MEDIA |

---

## Dependencias Transitivas

### A través de SHARED

API Mobile heredará estas dependencias de SHARED:
```
├─ PostgreSQL driver (pq)
├─ MongoDB driver (mongo-go-driver)
├─ RabbitMQ client (amqp)
├─ UUID generation (google/uuid)
├─ Date/Time utilities (carbon, chronos)
└─ Encryption utilities (crypto)
```

---

## Ciclo de Vida de Datos

```
┌─ Creación en API Mobile
│  ├─ Docente crea evaluación (POST)
│  └─ Guarda en PostgreSQL
│
├─ Enriquecimiento
│  ├─ Solicita quiz automático a Worker
│  └─ Recibe preguntas generadas
│
├─ Distribución
│  ├─ Asigna a estudiantes
│  └─ App móvil descarga evaluación
│
├─ Transformación
│  ├─ Estudiante responde
│  └─ API Mobile valida y calcula
│
├─ Persistencia Final
│  ├─ Respuestas → PostgreSQL (final_answers)
│  └─ Resultados → MongoDB (evaluation_results)
│
└─ Auditoría
   ├─ Eventos a RabbitMQ
   ├─ Logs centralizados via SHARED
   └─ Acceso desde API Admin para reportes
```

---

## Compatibilidad Entre Versiones

### API Mobile ↔ SHARED

| API Mobile | SHARED | Compatibilidad |
|-----------|--------|----------------|
| v1.0.0 | v1.3.0 | ✅ Compatible |
| v1.0.0 | v1.4.0 | ✅ Compatible (forward) |
| v1.0.0 | v2.0.0 | ❌ Breaking changes |

**Política:** API Mobile siempre debe mantener compatibilidad con última versión MINOR de SHARED.

### API Mobile ↔ PostgreSQL

| Versión | PostgreSQL | Compatibilidad |
|---------|-----------|----------------|
| v1.0+ | 14+ | ✅ Funciona |
| v1.0+ | 15 | ✅ Recomendado |
| v1.0+ | 13 | ⚠️ Legacy (deprecado) |

**Política:** API Mobile usa GORM, que soporta PostgreSQL 13+, pero se recomienda 15+.

---

## Cambios en Otros Servicios que Afectan a API Mobile

### Si SHARED cambia
- API Mobile debe actualizar go.mod
- Realizar testing de integración
- Potencial redeployment necesario

### Si PostgreSQL schema cambia
- Migraciones deben ser backward compatible
- API Mobile debe testar antes de upgrade
- Coordinar con API Admin (comparten BD)

### Si RabbitMQ cambia
- Actualizar configuración de conexión
- Cambios en formato de mensajes → impacta Worker
- Testing end-to-end requerido

### Si Worker cambia formato de mensajes
- API Mobile debe adaptar parsing
- Cambios en estructura de preguntas → impacta schema
- Versionar formato de payload

---

## Impacto de Cambios en API Mobile

### Si API Mobile cambia endpoint
- **Impacto:** App móvil debe actualizar
- **Versioning:** Mantener compatibilidad con versión N-1
- **Migración:** Soportar versiones múltiples vía URL prefix (/api/v1 vs /api/v2)

### Si API Mobile cambia schema de evaluación
- **Impacto:** Migraciones de BD
- **Validación:** API Mobile debe ser backward compatible
- **Testing:** Incluir tests de migración

### Si API Mobile publica nuevos eventos RabbitMQ
- **Impacto:** Cualquier consumidor (audit, reporting)
- **Política:** Versionar eventos con versionNumber o topic pattern

---

## Puntos de Fallo Críticos

| Fallo | Impacto | Mitigation |
|------|--------|-----------|
| PostgreSQL caída | API Mobile no funciona | Replicación + failover |
| MongoDB caída | Resultados no persisten | Caché en Redis + retry |
| RabbitMQ caída | Quizzes no se generan | Dead letter queues + retry |
| WORKER caída | Quizzes manuales OK | Usar evaluaciones manuales |
| SHARED versión incompatible | API Mobile no inicia | CI/CD validación de versiones |

---

## Checklist de Integración

- [ ] SHARED v1.3.0+ importado en go.mod
- [ ] Conexión PostgreSQL funcionando
- [ ] Conexión MongoDB funcionando
- [ ] Conexión RabbitMQ funcionando
- [ ] Logger centralizado de SHARED activo
- [ ] Auth middleware de SHARED validando tokens
- [ ] Health checks incluyendo todas dependencias
- [ ] Timeout global de SHARED aplicado
- [ ] Tests de integración con todos servicios
- [ ] Documentación de endpoints en Swagger
- [ ] Monitoreo de queues RabbitMQ
- [ ] Alertas de desconexión de BD
