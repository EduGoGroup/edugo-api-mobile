# Análisis del Sprint - Completar Queries Complejas (FASE 2.3)

## Resumen Ejecutivo

Este sprint representa la culminación de la FASE 2 (TODOs de Servicios) del plan maestro, enfocándose en completar el 80% restante de las queries complejas pendientes. Se trata de una evolución arquitectónica que optimiza las capacidades de consulta del sistema, mejorando significativamente la eficiencia en la recuperación de datos relacionados con materiales educativos, evaluaciones, progreso de usuarios y estadísticas globales.

El sprint aborda cinco áreas funcionales críticas: consultas de materiales con versionado, cálculo automático de puntajes en evaluaciones, generación de feedback detallado, actualización idempotente de progreso mediante UPSERT, y agregación de estadísticas globales. Cada área presenta desafíos técnicos únicos que requieren soluciones específicas respetando los principios de Clean Architecture.

La complejidad técnica es moderada-alta, ya que involucra operaciones en dos bases de datos diferentes (PostgreSQL y MongoDB), requiere optimización de queries para prevenir problemas de performance (N+1), y demanda implementación de lógica de negocio no trivial como el cálculo de puntajes y generación de feedback contextual.

El impacto arquitectónico es principalmente a nivel de la capa de aplicación (servicios) y la capa de infraestructura (repositorios), manteniendo intacta la capa de dominio, lo cual es consistente con los principios de diseño establecidos en el proyecto.

## Objetivo del Sprint

Completar la implementación de consultas complejas en los servicios de aplicación del sistema EduGo API Mobile, específicamente:

1. Habilitar consultas de materiales educativos que incluyan información completa de sus versiones históricas
2. Implementar la lógica de evaluación automática con cálculo de puntajes para diferentes tipos de preguntas
3. Generar feedback educativo detallado por pregunta en las evaluaciones
4. Implementar operaciones UPSERT para actualización de progreso sin duplicados
5. Crear queries agregadas para estadísticas globales del sistema

Este sprint cierra el 80% restante del PASO 2.3 del plan maestro, completando así íntegramente la FASE 2.

## Arquitectura Propuesta

### Tipo de Arquitectura

**Clean Architecture (Arquitectura Hexagonal)** - Mantenimiento y evolución de arquitectura existente.

El proyecto ya implementa Clean Architecture con tres capas claramente diferenciadas, y este sprint respeta esa estructura sin introducir cambios arquitectónicos disruptivos.

### Descripción de Arquitectura

La arquitectura del sistema sigue el patrón de capas de Clean Architecture, donde las dependencias fluyen hacia adentro (desde infrastructure hacia domain):

**Capa de Dominio** (`internal/domain/`):
- Contiene entidades de negocio puras (Material, Assessment, Progress, User)
- Define interfaces de repositorio (contratos) que serán implementados por infrastructure
- Define Value Objects y reglas de negocio core
- **No se modifica en este sprint** - las interfaces existentes son suficientes

**Capa de Aplicación** (`internal/application/`):
- Servicios que orquestan casos de uso de negocio
- DTOs para transferencia de datos entre capas
- **Principal área de cambio**: MaterialService, AssessmentService, ProgressService, StatsService
- Coordina llamadas a repositorios y aplica lógica de negocio

**Capa de Infraestructura** (`internal/infrastructure/`):
- Implementaciones concretas de repositorios
- Persistencia en PostgreSQL (`persistence/postgres/repository/`)
- Persistencia en MongoDB (`persistence/mongodb/repository/`)
- HTTP handlers (Gin) para exponer servicios vía REST
- **Área de cambio secundaria**: Implementar queries SQL/NoSQL optimizadas

**Capa de Container** (`internal/container/`):
- Inyección de dependencias con Wire
- Conecta todas las capas en tiempo de compilación
- **No requiere cambios** - configuración actual es suficiente

### Componentes Principales

#### 1. MaterialService (Capa de Aplicación)
- **Responsabilidad**: Gestionar operaciones de consulta de materiales educativos incluyendo versionado histórico
- **Tecnologías**: Go, PostgreSQL (lib/pq)
- **Interacciones**:
  - Lee de MaterialRepository (PostgreSQL) usando joins para incluir versiones
  - Puede notificar eventos vía RabbitMQ (ya configurado en PASO 2.1)
  - Transforma entidades de dominio a DTOs para exposición
- **Cambios específicos**:
  - Implementar método `GetMaterialWithVersions(materialID string)` que retorne material con historial completo de versiones
  - Implementar método `GetMaterialByVersion(materialID string, version int)` para consulta de versión específica
  - Optimizar queries con LEFT JOIN eficiente hacia tabla `material_versions`

#### 2. AssessmentService (Capa de Aplicación)
- **Responsabilidad**: Gestionar evaluaciones y cálculo automático de puntajes
- **Tecnologías**: Go, MongoDB (mongo-driver)
- **Interacciones**:
  - Lee evaluaciones y respuestas desde AssessmentRepository (MongoDB)
  - Calcula puntajes aplicando reglas de negocio complejas
  - Persiste resultados en colección `assessment_results`
  - Genera feedback detallado por pregunta
- **Cambios específicos**:
  - Implementar método `CalculateScore(assessmentID string, userResponses []Response)` con lógica multi-tipo
  - Soportar tipos de pregunta: multiple_choice, true_false, short_answer, fill_blank
  - Implementar estrategia de puntaje: correctas/totales * 100
  - Método `GenerateDetailedFeedback(assessmentID string, userResponses []Response)` que retorne array de feedback por pregunta

#### 3. ProgressService (Capa de Aplicación)
- **Responsabilidad**: Gestionar progreso de usuarios con actualización idempotente
- **Tecnologías**: Go, PostgreSQL (lib/pq)
- **Interacciones**:
  - Lee/escribe en ProgressRepository (PostgreSQL)
  - Usa operación UPSERT para evitar duplicados
  - Actualiza timestamp `last_updated_at`
- **Cambios específicos**:
  - Implementar método `UpdateProgress(userID string, materialID string, progress int)` con semántica UPSERT
  - Query SQL: `INSERT INTO user_progress ... ON CONFLICT (user_id, material_id) DO UPDATE SET ...`
  - Validar que progress esté en rango [0-100]

#### 4. StatsService (Capa de Aplicación)
- **Responsabilidad**: Generar estadísticas agregadas del sistema
- **Tecnologías**: Go, PostgreSQL + MongoDB
- **Interacciones**:
  - Consulta múltiples repositorios (MaterialRepository, AssessmentRepository, ProgressRepository)
  - Agrega métricas cross-database
  - Cacheo opcional (no implementado en este sprint)
- **Cambios específicos**:
  - Implementar método `GetGlobalStats()` que retorne:
    - Total de materiales publicados
    - Total de evaluaciones completadas
    - Promedio de puntajes globales
    - Total de usuarios activos
    - Métricas de progreso promedio
  - Optimizar con queries agregadas (COUNT, AVG, SUM)

#### 5. Repositorios de Infraestructura
- **MaterialRepositoryImpl** (PostgreSQL):
  - Implementar query con JOIN a `material_versions`
  - Mapeo eficiente de resultados con múltiples filas
  - Índice ya existe en `materials.updated_at` (creado en tarea anterior)

- **AssessmentRepositoryImpl** (MongoDB):
  - Implementar queries de lookup para evaluaciones + respuestas
  - Pipeline de agregación para cálculo de estadísticas
  - Optimizar con índices en campos frecuentemente consultados

- **ProgressRepositoryImpl** (PostgreSQL):
  - Implementar query UPSERT nativo de PostgreSQL
  - Constraint UNIQUE en (user_id, material_id) para prevenir duplicados

### Interacciones

**Flujo típico de consulta de material con versiones**:
1. Cliente HTTP → MaterialHandler (Gin)
2. Handler → MaterialService.GetMaterialWithVersions()
3. Service → MaterialRepository.FindByIDWithVersions()
4. Repository ejecuta query SQL con LEFT JOIN
5. Repository mapea filas a entidad Material con array de Versions
6. Service transforma a MaterialDTO
7. Handler serializa JSON y retorna al cliente

**Flujo típico de evaluación con cálculo de puntaje**:
1. Cliente HTTP → AssessmentHandler (Gin) con respuestas de usuario
2. Handler → AssessmentService.SubmitAssessment()
3. Service → AssessmentRepository.FindByID() para obtener respuestas correctas
4. Service ejecuta lógica de cálculo de puntaje (comparación respuesta correcta vs. enviada)
5. Service → AssessmentService.GenerateDetailedFeedback() genera feedback por pregunta
6. Service → AssessmentRepository.SaveResult() persiste resultado en MongoDB
7. Service → RabbitMQ publica evento "assessment_completed" (opcional)
8. Handler retorna resultado con puntaje y feedback al cliente

**Flujo típico de actualización de progreso**:
1. Cliente HTTP → ProgressHandler (Gin) con progreso actualizado
2. Handler → ProgressService.UpdateProgress()
3. Service valida progreso en rango [0-100]
4. Service → ProgressRepository.Upsert()
5. Repository ejecuta query UPSERT en PostgreSQL
6. Handler retorna confirmación al cliente

## Modelo de Datos

### Estrategia de Persistencia

**Híbrido: Relacional (PostgreSQL) + NoSQL (MongoDB)**

**Justificación**:
- **PostgreSQL** para datos estructurados con relaciones fuertes (materiales, usuarios, progreso)
- **MongoDB** para datos semi-estructurados flexibles (evaluaciones, respuestas variables)

Esta estrategia ya está implementada en el proyecto y se mantiene en este sprint.

### Entidades Principales (sin modificar - ya definidas en domain)

**Entidad: Material** (PostgreSQL)
- Campos: id, title, description, type, content_url, author_id, published_at, updated_at, is_published
- Propósito: Representar recursos educativos del sistema
- **Relación nueva**: 1:N con MaterialVersion (historial de cambios)

**Entidad: MaterialVersion** (PostgreSQL)
- Campos: id, material_id, version_number, title, content_url, changed_by, created_at
- Propósito: Mantener historial de cambios de materiales
- **Relación**: N:1 con Material (FK: material_id)

**Entidad: Assessment** (MongoDB - colección)
- Campos: _id, title, description, questions (array), created_at, updated_at
- Propósito: Representar evaluaciones con preguntas y respuestas correctas
- Estructura: Document con subdocumentos anidados (flexible)

**Entidad: AssessmentResult** (MongoDB - colección)
- Campos: _id, assessment_id, user_id, score, total_questions, correct_answers, feedback (array), submitted_at
- Propósito: Almacenar resultados de evaluaciones completadas
- **Nueva colección** creada en este sprint

**Entidad: UserProgress** (PostgreSQL)
- Campos: id, user_id, material_id, progress_percentage, last_updated_at, completed_at
- Propósito: Rastrear progreso de usuarios en materiales
- **Constraint UNIQUE** en (user_id, material_id) para UPSERT

### Relaciones

**Material ↔ MaterialVersion**: 1:N
- Un material puede tener múltiples versiones históricas
- Cada versión apunta a un material específico
- Query JOIN: `SELECT m.*, v.* FROM materials m LEFT JOIN material_versions v ON m.id = v.material_id WHERE m.id = $1`

**User ↔ UserProgress ↔ Material**: N:M a través de tabla intermedia
- Un usuario puede tener progreso en múltiples materiales
- Un material puede tener progreso de múltiples usuarios
- Constraint UNIQUE previene duplicados

**User ↔ AssessmentResult ↔ Assessment**: N:M (MongoDB lookup)
- Un usuario puede completar múltiples evaluaciones
- Una evaluación puede ser completada por múltiples usuarios
- Resultados separados por documento

### Scripts de Creación (Borrador PostgreSQL)

```sql
-- Tabla material_versions (si no existe)
CREATE TABLE IF NOT EXISTS material_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    material_id UUID NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content_url TEXT NOT NULL,
    changed_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_material_version UNIQUE(material_id, version_number)
);

CREATE INDEX idx_material_versions_material_id ON material_versions(material_id);
CREATE INDEX idx_material_versions_created_at ON material_versions(created_at DESC);

-- Agregar constraint UNIQUE en user_progress si no existe
ALTER TABLE user_progress
ADD CONSTRAINT unique_user_material UNIQUE(user_id, material_id);

CREATE INDEX idx_user_progress_user_id ON user_progress(user_id);
CREATE INDEX idx_user_progress_material_id ON user_progress(material_id);
```

### Colecciones MongoDB (Borrador)

```javascript
// Colección assessment_results
db.createCollection("assessment_results");

db.assessment_results.createIndex({ "assessment_id": 1, "user_id": 1 }, { unique: true });
db.assessment_results.createIndex({ "submitted_at": -1 });
db.assessment_results.createIndex({ "user_id": 1, "submitted_at": -1 });

// Estructura de documento assessment_result
{
  "_id": ObjectId("..."),
  "assessment_id": "uuid-string",
  "user_id": "uuid-string",
  "score": 85.5,
  "total_questions": 10,
  "correct_answers": 8,
  "feedback": [
    {
      "question_id": "q1",
      "is_correct": true,
      "user_answer": "B",
      "correct_answer": "B",
      "explanation": "Correcto. La respuesta B es..."
    },
    {
      "question_id": "q2",
      "is_correct": false,
      "user_answer": "A",
      "correct_answer": "C",
      "explanation": "Incorrecto. La respuesta correcta es C porque..."
    }
  ],
  "submitted_at": ISODate("2025-11-05T...")
}
```

## Flujo de Procesos

### Proceso Principal: Consulta de Material con Versiones

**Descripción**: Este flujo permite a un usuario consultar un material educativo incluyendo todo su historial de versiones.

**Pasos detallados**:

1. **Recepción de solicitud HTTP**: Cliente envía `GET /api/v1/materials/{id}/versions` al handler
2. **Validación de entrada**: MaterialHandler valida UUID del material y permisos de usuario
3. **Invocación de servicio**: Handler llama a `MaterialService.GetMaterialWithVersions(materialID)`
4. **Consulta de datos**: Service invoca `MaterialRepository.FindByIDWithVersions(materialID)`
5. **Ejecución de query SQL**: Repository ejecuta LEFT JOIN entre `materials` y `material_versions`:
   ```sql
   SELECT
     m.id, m.title, m.description, m.type, m.content_url, m.published_at,
     v.id as version_id, v.version_number, v.title as version_title,
     v.content_url as version_url, v.created_at as version_created_at
   FROM materials m
   LEFT JOIN material_versions v ON m.id = v.material_id
   WHERE m.id = $1
   ORDER BY v.version_number DESC
   ```
6. **Mapeo de resultados**: Repository itera sobre filas retornadas y construye objeto Material con array de Versions
7. **Transformación a DTO**: Service transforma entidad de domain a `MaterialWithVersionsDTO` (incluye solo campos expuestos)
8. **Serialización y respuesta**: Handler serializa DTO a JSON y retorna con código 200
9. **Logging**: Logger registra operación exitosa con materialID, tiempo de ejecución
10. **Manejo de errores**: Si material no existe, retornar 404; si error de DB, retornar 500

**Variante - Consulta de versión específica**:
- Cliente envía `GET /api/v1/materials/{id}/versions/{versionNumber}`
- Query filtra por `v.version_number = $2`
- Retorna solo una versión específica

### Proceso Secundario: Evaluación con Cálculo de Puntaje

**Descripción**: Este flujo permite que un usuario complete una evaluación, el sistema calcule automáticamente el puntaje y genere feedback detallado.

**Pasos detallados**:

1. **Recepción de respuestas**: Cliente envía `POST /api/v1/assessments/{id}/submit` con body JSON:
   ```json
   {
     "user_id": "uuid",
     "responses": [
       {"question_id": "q1", "answer": "B"},
       {"question_id": "q2", "answer": "true"},
       {"question_id": "q3", "answer": "París"}
     ]
   }
   ```
2. **Validación de entrada**: AssessmentHandler valida que assessment existe y usuario está autenticado
3. **Obtención de evaluación**: Service llama a `AssessmentRepository.FindByID(assessmentID)` para obtener preguntas y respuestas correctas
4. **Inicialización de contadores**: Service inicializa variables: totalQuestions, correctAnswers, score
5. **Iteración y comparación**: Para cada respuesta de usuario:
   - Buscar pregunta correspondiente en assessment
   - Comparar respuesta enviada vs. respuesta correcta (case-insensitive, trim)
   - Incrementar correctAnswers si coincide
   - Almacenar resultado de pregunta para feedback
6. **Cálculo de puntaje**: `score = (correctAnswers / totalQuestions) * 100`
7. **Generación de feedback**: Service llama a método interno que construye array de feedback:
   ```go
   feedback := []FeedbackItem{
     {
       QuestionID: "q1",
       IsCorrect: true,
       UserAnswer: "B",
       CorrectAnswer: "B",
       Explanation: "Correcto. La opción B es la respuesta adecuada porque..."
     },
     // ...
   }
   ```
8. **Persistencia de resultado**: Service llama a `AssessmentRepository.SaveResult()` que inserta documento en MongoDB:
   ```javascript
   db.assessment_results.insertOne({
     assessment_id: "uuid",
     user_id: "uuid",
     score: 85.5,
     feedback: [...],
     submitted_at: new Date()
   })
   ```
9. **Publicación de evento (opcional)**: Service publica evento "assessment_completed" a RabbitMQ para notificaciones
10. **Respuesta al cliente**: Handler retorna resultado con puntaje, feedback y timestamp

**Casos especiales**:
- **Evaluación ya completada**: Verificar en `assessment_results` si existe documento con (assessment_id, user_id), retornar 409 Conflict
- **Respuesta inválida**: Si tipo de respuesta no coincide con tipo de pregunta, marcar como incorrecta y explicar en feedback
- **Pregunta sin responder**: Marcar como incorrecta, no incrementar totalQuestions

### Proceso Terciario: Actualización de Progreso (UPSERT)

**Descripción**: Este flujo permite actualizar el progreso de un usuario en un material de forma idempotente, sin generar duplicados.

**Pasos detallados**:

1. **Recepción de actualización**: Cliente envía `PUT /api/v1/progress` con body:
   ```json
   {
     "user_id": "uuid",
     "material_id": "uuid",
     "progress_percentage": 75
   }
   ```
2. **Validación de entrada**: ProgressHandler valida:
   - user_id y material_id son UUIDs válidos
   - progress_percentage está en rango [0-100]
   - Usuario autenticado coincide con user_id (o es admin)
3. **Invocación de servicio**: Handler llama a `ProgressService.UpdateProgress(userID, materialID, progressPercentage)`
4. **Preparación de query UPSERT**: Service prepara llamada a repository con datos validados
5. **Ejecución de UPSERT**: Repository ejecuta query PostgreSQL:
   ```sql
   INSERT INTO user_progress (user_id, material_id, progress_percentage, last_updated_at, completed_at)
   VALUES ($1, $2, $3, NOW(), CASE WHEN $3 = 100 THEN NOW() ELSE NULL END)
   ON CONFLICT (user_id, material_id)
   DO UPDATE SET
     progress_percentage = EXCLUDED.progress_percentage,
     last_updated_at = NOW(),
     completed_at = CASE
       WHEN EXCLUDED.progress_percentage = 100 THEN NOW()
       WHEN user_progress.completed_at IS NOT NULL THEN user_progress.completed_at
       ELSE NULL
     END
   RETURNING *;
   ```
6. **Mapeo de resultado**: Repository mapea fila retornada a entidad Progress
7. **Verificación de completitud**: Si progress = 100 y completed_at se actualizó:
   - Service publica evento "material_completed" a RabbitMQ
   - Opcional: Desbloquear siguiente material en secuencia
8. **Transformación a DTO**: Service transforma a ProgressDTO
9. **Respuesta al cliente**: Handler retorna progreso actualizado con código 200

**Ventajas del UPSERT**:
- **Idempotencia**: Múltiples llamadas con mismo progreso no generan errores
- **Simplicidad**: Cliente no necesita verificar si registro existe antes de actualizar
- **Atomicidad**: INSERT o UPDATE ocurren en una sola transacción

### Proceso Cuaternario: Generación de Estadísticas Globales

**Descripción**: Este flujo permite a administradores obtener métricas agregadas del sistema completo.

**Pasos detallados**:

1. **Recepción de solicitud**: Cliente admin envía `GET /api/v1/stats/global`
2. **Verificación de permisos**: StatsHandler valida que usuario tiene rol "admin"
3. **Invocación de servicio**: Handler llama a `StatsService.GetGlobalStats()`
4. **Consultas paralelas** (usando goroutines con sync.WaitGroup):
   - **Goroutine 1**: Consultar total de materiales publicados en PostgreSQL:
     ```sql
     SELECT COUNT(*) FROM materials WHERE is_published = true;
     ```
   - **Goroutine 2**: Consultar total de evaluaciones completadas en MongoDB:
     ```javascript
     db.assessment_results.countDocuments({})
     ```
   - **Goroutine 3**: Calcular promedio de puntajes en MongoDB:
     ```javascript
     db.assessment_results.aggregate([
       { $group: { _id: null, avgScore: { $avg: "$score" } } }
     ])
     ```
   - **Goroutine 4**: Consultar total de usuarios activos (con progreso reciente) en PostgreSQL:
     ```sql
     SELECT COUNT(DISTINCT user_id) FROM user_progress
     WHERE last_updated_at >= NOW() - INTERVAL '30 days';
     ```
   - **Goroutine 5**: Calcular progreso promedio global en PostgreSQL:
     ```sql
     SELECT AVG(progress_percentage) FROM user_progress;
     ```
5. **Agregación de resultados**: Service espera a que todas las goroutines terminen y construye objeto GlobalStats:
   ```go
   stats := GlobalStats{
     TotalMaterials: materialsCount,
     TotalAssessmentsCompleted: assessmentsCount,
     AverageScore: avgScore,
     ActiveUsers30Days: activeUsers,
     AverageProgressGlobal: avgProgress,
     GeneratedAt: time.Now(),
   }
   ```
6. **Transformación a DTO**: Service transforma a StatsDTO
7. **Respuesta al cliente**: Handler retorna estadísticas con código 200

**Optimizaciones**:
- Usar índices en campos de fecha (`last_updated_at`, `submitted_at`)
- Cachear resultado por 5-10 minutos (implementar en futuro)
- Limitar a usuarios con permiso de admin para prevenir abuso

## Patrones de Diseño Recomendados

### 1. Repository Pattern (Ya implementado)
**Por qué es apropiado**: Abstrae la lógica de acceso a datos, permitiendo cambiar implementaciones de persistencia sin afectar servicios. En este sprint se extiende para soportar queries más complejas manteniendo la misma interfaz.

### 2. Service Layer Pattern (Ya implementado)
**Por qué es apropiado**: Encapsula lógica de negocio en servicios reutilizables. En este sprint se enriquece con cálculo de puntajes y generación de feedback, que son casos de uso de negocio complejos.

### 3. DTO Pattern (Ya implementado)
**Por qué es apropiado**: Separa representación interna de entidades de su exposición externa. Permite evolucionar el dominio sin romper contratos de API.

### 4. Strategy Pattern (Nuevo - para cálculo de puntajes)
**Por qué es apropiado**: Diferentes tipos de preguntas pueden requerir diferentes estrategias de evaluación. Implementar interfaz `ScoringStrategy` permite extender fácilmente nuevos tipos de evaluación:
```go
type ScoringStrategy interface {
    CalculateScore(question Question, userAnswer string) (score float64, isCorrect bool)
}

type MultipleChoiceStrategy struct {}
type TrueFalseStrategy struct {}
type ShortAnswerStrategy struct {}
```

### 5. Template Method Pattern (Nuevo - para feedback)
**Por qué es apropiado**: El proceso de generación de feedback tiene pasos comunes (validar respuesta, comparar, formatear) pero detalles específicos por tipo de pregunta. Implementar método plantilla permite reutilizar estructura:
```go
func (s *AssessmentService) GenerateFeedback(question Question, userAnswer string) FeedbackItem {
    // Pasos comunes
    isCorrect := s.compareAnswer(question, userAnswer)
    explanation := s.getExplanation(question, isCorrect)

    // Detalle específico según tipo
    return s.formatFeedback(question, userAnswer, isCorrect, explanation)
}
```

### 6. Builder Pattern (Opcional - para queries complejas)
**Por qué es apropiado**: Construir queries SQL/MongoDB complejas con múltiples condiciones opcionales de forma fluida y legible:
```go
query := NewMaterialQueryBuilder().
    WithVersions().
    FilterByPublished(true).
    OrderByUpdatedAt(DESC).
    Build()
```

## Stack Tecnológico Recomendado

### Backend (Ya establecido - no cambia)
- **Lenguaje**: Go 1.21+
- **Framework Web**: Gin 1.9+
- **Justificación**: Alto rendimiento, routing eficiente, middleware robusto

### Base de Datos (Ya establecido - no cambia)
- **Relacional**: PostgreSQL 16
  - **Justificación**: Soporte UPSERT nativo, índices eficientes, JOINs optimizados, constraints fuertes
  - **Driver**: lib/pq (driver nativo de PostgreSQL)
- **NoSQL**: MongoDB 7
  - **Justificación**: Flexibilidad para esquemas de evaluación variables, agregación pipeline potente
  - **Driver**: mongo-driver oficial

### Messaging (Ya configurado en PASO 2.1)
- **RabbitMQ**: Para eventos asíncronos (assessment_completed, material_completed)
- **Justificación**: Ya configurado, permite desacoplar notificaciones

### Storage (Ya configurado en PASO 2.2)
- **AWS S3**: Para almacenamiento de contenido de materiales
- **Justificación**: URLs firmadas ya implementadas, escalable

### Logging y Error Handling
- **Logger**: Zap Logger de `edugo-shared`
- **Errors**: Error types de `edugo-shared/common/errors`
- **Justificación**: Logging estructurado de alto rendimiento, manejo de errores estandarizado

### Testing
- **Framework**: Testing nativo de Go (`testing` package)
- **Mocks**: Interfaces para mocks de repositorios
- **Justificación**: Tests unitarios rápidos sin dependencias externas

## Consideraciones No Funcionales

### Escalabilidad

**Estrategia**:
- **Queries optimizadas**: Uso de índices apropiados en PostgreSQL y MongoDB para evitar table scans
- **Paginación**: Aunque no se implementa en este sprint, los métodos deben diseñarse para soportar paginación futura
- **Conexiones a BD**: Connection pooling ya configurado en container
- **Caché**: No se implementa en este sprint, pero queries de estadísticas son candidatos para cacheo futuro (Redis)

**Bottlenecks identificados**:
- Query de estadísticas globales puede ser costoso con millones de registros → Solución: índices + cacheo
- Query de materiales con muchas versiones (100+) → Solución: limitar versiones retornadas a últimas 50

### Seguridad

**Medidas**:
- **Autenticación**: JWT ya configurado en middleware (edugo-shared/auth)
- **Autorización**: Validar que usuario solo accede a su propio progreso (excepto admins)
- **Sanitización**: Validar inputs en handlers antes de pasar a servicios
- **SQL Injection**: Usar prepared statements con placeholders ($1, $2) - nunca string concatenation
- **NoSQL Injection**: Usar struct binding en MongoDB, evitar queries dinámicas construidas con strings

**Permisos**:
- `/materials/{id}/versions`: Autenticado (todos los usuarios)
- `/assessments/{id}/submit`: Autenticado (usuario propio)
- `/progress`: Autenticado (usuario propio o admin)
- `/stats/global`: Autenticado (solo admins)

### Performance

**Optimizaciones propuestas**:

1. **Índices en PostgreSQL** (crear si no existen):
   ```sql
   CREATE INDEX idx_material_versions_material_id ON material_versions(material_id);
   CREATE INDEX idx_user_progress_user_material ON user_progress(user_id, material_id);
   CREATE INDEX idx_user_progress_updated_at ON user_progress(last_updated_at);
   ```

2. **Índices en MongoDB** (crear):
   ```javascript
   db.assessment_results.createIndex({ "assessment_id": 1, "user_id": 1 }, { unique: true });
   db.assessment_results.createIndex({ "submitted_at": -1 });
   db.assessments.createIndex({ "_id": 1, "questions.id": 1 });
   ```

3. **Evitar N+1 queries**:
   - En MaterialService: Usar LEFT JOIN en lugar de queries separadas por versión
   - En AssessmentService: Fetch assessment con todas las preguntas en una sola query

4. **Timeouts**:
   - Configurar context timeout de 5 segundos para queries de BD
   - Si query excede timeout, retornar error 504 Gateway Timeout

5. **Query selectiva**:
   - Solo SELECT campos necesarios, evitar `SELECT *` cuando sea posible
   - Proyectar solo campos necesarios en MongoDB

**Métricas objetivo**:
- **Latencia p95**: < 200ms para queries simples (material con versiones)
- **Latencia p95**: < 500ms para queries complejas (estadísticas globales)
- **Throughput**: 500 req/s por instancia

### Mantenibilidad

**Prácticas recomendadas**:

1. **Comentarios en código**:
   - Explicar lógica compleja de cálculo de puntajes
   - Documentar queries SQL/MongoDB complejas con comentarios inline
   - Explicar decisiones de diseño no obvias

2. **Tests unitarios**:
   - Crear tests para cada método nuevo en servicios
   - Usar table-driven tests para diferentes tipos de preguntas
   - Alcanzar mínimo 80% de cobertura en código nuevo
   ```go
   func TestAssessmentService_CalculateScore(t *testing.T) {
       tests := []struct {
           name           string
           assessmentID   string
           userResponses  []Response
           expectedScore  float64
           expectedError  error
       }{
           {
               name: "all_correct_answers",
               // ...
           },
           {
               name: "partial_correct_answers",
               // ...
           },
           {
               name: "no_correct_answers",
               // ...
           },
       }
       // ...
   }
   ```

3. **Separación de responsabilidades**:
   - Service: Lógica de negocio (cálculo, validación, orquestación)
   - Repository: Solo acceso a datos (queries, mapeo)
   - Handler: Solo validación de entrada y serialización

4. **Logging contextual**:
   ```go
   logger.Info("calculating assessment score",
       zap.String("assessment_id", assessmentID),
       zap.String("user_id", userID),
       zap.Int("total_questions", len(responses)),
       zap.Float64("score", score),
   )
   ```

5. **Error handling consistente**:
   ```go
   if err != nil {
       logger.Error("failed to fetch assessment", zap.Error(err))
       return nil, errors.NewInternalError("error al obtener evaluación")
   }
   ```

## Riesgos Identificados

### 1. Performance degradation con datasets grandes
**Descripción**: Queries de materiales con muchas versiones o estadísticas con millones de registros pueden causar timeouts y degradación de performance.

**Impacto**: Alto - Afecta experiencia de usuario y disponibilidad del sistema

**Probabilidad**: Media - Depende del crecimiento del sistema

**Mitigación**:
- Implementar índices apropiados en todas las tablas/colecciones afectadas
- Limitar número de versiones retornadas (últimas 50)
- Configurar query timeout de 5 segundos
- Implementar monitoring de query time en logs
- Plan futuro: Implementar caché para estadísticas (Redis)

### 2. Inconsistencia en cálculo de puntajes entre tipos de preguntas
**Descripción**: Diferentes tipos de preguntas (multiple choice, true/false, short answer) pueden requerir lógica de comparación diferente, causando puntajes incorrectos.

**Impacto**: Alto - Afecta integridad de resultados educativos

**Probabilidad**: Media - Complejidad de implementación

**Mitigación**:
- Implementar Strategy Pattern para aislar lógica de cada tipo
- Tests exhaustivos con table-driven tests para todos los tipos
- Revisión de código enfocada en lógica de comparación
- Documentar claramente criterios de evaluación en comentarios

### 3. Race conditions en actualización de progreso
**Descripción**: Múltiples requests simultáneos de actualización de progreso del mismo usuario en mismo material podrían causar condiciones de carrera.

**Impacto**: Bajo - UPSERT maneja esto, pero podría haber inconsistencias temporales

**Probabilidad**: Baja - Frontend debería prevenir esto, pero backend debe ser robusto

**Mitigación**:
- Usar UPSERT con ON CONFLICT para garantizar atomicidad
- Constraint UNIQUE previene duplicados
- Transaction isolation level SERIALIZABLE si es necesario
- Tests de concurrencia simulando requests paralelos

### 4. Queries N+1 en evaluación con muchas preguntas
**Descripción**: Si se fetch cada pregunta individualmente para comparar respuestas, se genera query N+1 con degradación de performance.

**Impacto**: Medio - Performance degradada en evaluaciones largas

**Probabilidad**: Media - Depende de implementación

**Mitigación**:
- Fetch completo de assessment con todas las preguntas en una sola query
- Usar pipeline de agregación en MongoDB para obtener todo de una vez
- No iterar haciendo queries individuales

### 5. Falta de validación de tipos de respuesta
**Descripción**: Si cliente envía respuesta de tipo incorrecto (ej: texto en pregunta booleana), podría causar panics o resultados incorrectos.

**Impacto**: Medio - Errores de runtime y resultados incorrectos

**Probabilidad**: Media - Clientes pueden enviar datos malformados

**Mitigación**:
- Validación exhaustiva en handler antes de pasar a service
- Type assertions con manejo de panic recovery
- Retornar error 400 Bad Request con mensaje descriptivo
- Tests con inputs malformados

### 6. Falta de tests de integración
**Descripción**: Tests unitarios con mocks pueden no detectar problemas de integración con bases de datos reales (sintaxis SQL/MongoDB incorrecta, tipos incompatibles).

**Impacto**: Alto - Bugs en producción

**Probabilidad**: Media - Complejidad de queries

**Mitigación**:
- Validar manualmente queries en consola de PostgreSQL/MongoDB
- Ejecutar aplicación localmente con bases de datos reales
- Plan futuro: Implementar tests de integración con testcontainers (FASE 4)

## Siguientes Pasos Recomendados

### Fase de Implementación (Orden sugerido)

1. **[2 horas] Implementar queries de materiales con versiones**
   - Crear método `FindByIDWithVersions` en MaterialRepositoryImpl
   - Crear método `GetMaterialWithVersions` en MaterialService
   - Agregar handler endpoint `GET /materials/{id}/versions`
   - Tests unitarios

2. **[3 horas] Implementar cálculo de puntajes en AssessmentService**
   - Crear método `CalculateScore` con lógica multi-tipo
   - Implementar Strategy Pattern para tipos de pregunta
   - Crear método `SaveResult` en AssessmentRepositoryImpl
   - Tests exhaustivos (table-driven tests)

3. **[1 hora] Implementar generación de feedback detallado**
   - Crear método `GenerateDetailedFeedback` en AssessmentService
   - Integrar con cálculo de puntajes
   - Tests unitarios

4. **[1 hora] Implementar UPSERT de progreso**
   - Crear método `Upsert` en ProgressRepositoryImpl con query ON CONFLICT
   - Crear método `UpdateProgress` en ProgressService
   - Tests unitarios (incluir test de idempotencia)

5. **[2 horas] Implementar query de estadísticas globales**
   - Crear métodos de consulta en cada repositorio (MaterialRepository, AssessmentRepository, ProgressRepository)
   - Crear método `GetGlobalStats` en StatsService con goroutines paralelas
   - Agregar handler endpoint `GET /stats/global`
   - Tests unitarios

6. **[1 hora] Validación final y refinamiento**
   - Ejecutar `go build ./...` y verificar compilación
   - Ejecutar `go test ./...` y verificar que todos los tests pasan
   - Verificar cobertura de tests: `go test -cover ./...` (objetivo: ≥80%)
   - Prueba manual con Postman/curl de endpoints nuevos
   - Revisar logs y error handling

7. **[0.5 horas] Documentación**
   - Agregar comentarios en código complejo
   - Actualizar README si es necesario
   - Documentar ejemplos de uso en comentarios

8. **[0.5 horas] Commit atómico**
   - Agregar archivos modificados a staging
   - Crear commit con mensaje descriptivo siguiendo formato del proyecto
   - Actualizar `sprint/current/readme.md` marcando todas las casillas como completadas

### Post-Sprint

1. **Crear PR** con el commit del sprint
2. **Solicitar code review** enfocado en:
   - Correctitud de queries SQL/MongoDB
   - Lógica de cálculo de puntajes
   - Manejo de errores y edge cases
3. **Ejecutar linters y formatters**: `gofmt`, `golangci-lint`
4. **Merge a main** después de aprobación
5. **Continuar con FASE 3** del plan maestro: Limpieza y Consolidación (eliminar código duplicado)

### Mejoras Futuras (Fuera del alcance de este sprint)

- Implementar caché de estadísticas con Redis
- Agregar paginación a consulta de versiones
- Implementar tests de integración con testcontainers (FASE 4)
- Agregar monitoreo de query time con Prometheus
- Implementar soft delete para materiales y evaluaciones
- Soportar tipos de pregunta adicionales (essay, file upload)

---

💡 **Nota**: Este es un análisis rápido sin diagramas. Para análisis completo con diagramas visuales de arquitectura, modelo de datos y flujos de proceso, ejecuta: `/01-analysis --mode=full`

---

**Generado**: 2025-11-05
**Modo**: quick (análisis ejecutivo sin diagramas)
**Alcance**: Completo (todas las fases del sprint)
**Fuente**: sprint/current/readme.md
