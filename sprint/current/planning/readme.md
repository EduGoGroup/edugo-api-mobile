# Plan de Trabajo - Fase 2: Completar TODOs de Servicios

## Resumen del Proyecto

Completar la implementación de tres servicios fundamentales pendientes en la arquitectura Clean Architecture (Hexagonal) existente: RabbitMQ Messaging para eventos de dominio, AWS S3 Storage para URLs firmadas, y Queries Complejas optimizadas en PostgreSQL y MongoDB.

## Stack Tecnológico

- **Backend**: Go 1.21+ con framework Gin
- **Base de Datos Relacional**: PostgreSQL (driver lib/pq)
- **Base de Datos NoSQL**: MongoDB (mongo-driver)
- **Messaging**: RabbitMQ (cliente amqp091-go)
- **Storage**: AWS S3 (aws-sdk-go-v2)
- **Logger**: Zap (edugo-shared/logger)
- **Testing**: Testcontainers

---

## 📋 Plan de Ejecución

### Fase 1: Implementar RabbitMQ Messaging

**Objetivo**: Crear sistema de publicación de eventos de dominio usando RabbitMQ para notificar acciones críticas (material_uploaded, assessment_attempt_recorded).

**Tareas**:

- [ ] **1.1** - Agregar dependencias de RabbitMQ al proyecto
  - **Descripción**: Ejecutar `go get github.com/rabbitmq/amqp091-go@v1.9.0` y actualizar go.mod
  - **Archivos a crear/modificar**: `go.mod`, `go.sum`
  - **Criterio de aceptación**: Dependencia aparece en go.mod y compila sin errores

- [ ] **1.2** - Crear estructura de eventos de dominio
  - **Descripción**: Definir tipos de eventos (MaterialUploadedEvent, AssessmentAttemptRecordedEvent) con campos y métodos de serialización JSON
  - **Archivos a crear/modificar**:
    - `internal/infrastructure/messaging/events.go` (NUEVO)
  - **Criterio de aceptación**: Structs con tags JSON, métodos ToJSON() que retornan []byte

- [ ] **1.3** - Implementar RabbitMQ Publisher
  - **Descripción**: Crear cliente RabbitMQ con conexión persistente, declaración de exchanges, publicación de mensajes con publisher confirms
  - **Archivos a crear/modificar**:
    - `internal/infrastructure/messaging/rabbitmq/publisher.go` (NUEVO)
  - **Criterio de aceptación**:
    - Struct `RabbitMQPublisher` con métodos Connect(), Publish(), Close()
    - Manejo de reconexión automática
    - Logging con Zap de eventos publicados
  - 🔗 **Depende de**: Tarea 1.2

- [ ] **1.4** - Agregar configuración de RabbitMQ
  - **Descripción**: Extender config.go con struct RabbitMQConfig (URL, Exchange, QueueNames) y agregar sección en config.yaml
  - **Archivos a crear/modificar**:
    - `internal/config/config.go`
    - `config/config.yaml`
  - **Criterio de aceptación**:
    - Struct RabbitMQConfig con campos URL, Exchange string
    - config.yaml tiene sección rabbitmq con valores de ejemplo

- [ ] **1.5** - Inicializar RabbitMQ en main.go
  - **Descripción**: Instanciar RabbitMQPublisher en cmd/main.go con config de Viper, conectar al servidor, agregar defer Close()
  - **Archivos a crear/modificar**:
    - `cmd/main.go`
  - **Criterio de aceptación**:
    - Publisher inicializado después de DBs y antes de Container
    - Log de conexión exitosa
    - Graceful shutdown con Close()
  - 🔗 **Depende de**: Tareas 1.3, 1.4

- [ ] **1.6** - Agregar Publisher al Container DI
  - **Descripción**: Extender Container struct con campo MessagePublisher, inyectar en constructor NewContainer()
  - **Archivos a crear/modificar**:
    - `internal/container/container.go`
  - **Criterio de aceptación**:
    - Campo `MessagePublisher` en struct Container
    - Servicios pueden acceder a publisher via container
  - 🔗 **Depende de**: Tarea 1.5

- [ ] **1.7** - Integrar eventos en MaterialService
  - **Descripción**: Inyectar MessagePublisher en MaterialService, publicar MaterialUploadedEvent después de crear material en método CreateMaterial()
  - **Archivos a crear/modificar**:
    - `internal/application/service/material_service.go`
  - **Criterio de aceptación**:
    - Evento publicado con payload: material_id, title, content_type, uploaded_at
    - Si falla publicación, log warning pero NO falla el request HTTP
  - 🔗 **Depende de**: Tarea 1.6

- [ ] **1.8** - Integrar eventos en AssessmentService
  - **Descripción**: Inyectar MessagePublisher en AssessmentService, publicar AssessmentAttemptRecordedEvent después de registrar intento en método RecordAttempt()
  - **Archivos a crear/modificar**:
    - `internal/application/service/assessment_service.go`
  - **Criterio de aceptación**:
    - Evento publicado con payload: attempt_id, user_id, assessment_id, score, submitted_at
    - Si falla publicación, log warning pero NO falla el request HTTP
  - 🔗 **Depende de**: Tarea 1.6

- [ ] **1.9** - Crear tests unitarios de RabbitMQ Publisher
  - **Descripción**: Tests con mock de amqp.Connection para validar lógica de publicación, serialización de eventos, manejo de errores
  - **Archivos a crear/modificar**:
    - `internal/infrastructure/messaging/rabbitmq/publisher_test.go` (NUEVO)
  - **Criterio de aceptación**:
    - Al menos 5 tests: conexión exitosa, publicación exitosa, error de conexión, error de publicación, serialización JSON
    - Tests pasan con `go test ./internal/infrastructure/messaging/...`
  - 🔗 **Depende de**: Tarea 1.3

- [ ] **1.10** - Commit atómico de RabbitMQ Messaging
  - **Descripción**: Crear commit con mensaje "feat: implementar messaging RabbitMQ para eventos"
  - **Archivos incluidos**: Todos los archivos modificados en Fase 1
  - **Criterio de aceptación**:
    - `go build ./...` compila sin errores
    - Tests pasan: `go test ./...`
    - Commit creado en branch feature/fase2-servicios
  - 🔗 **Depende de**: Tareas 1.1 a 1.9

**Completitud de Fase 1**: 0/10 tareas completadas

---

### Fase 2: Implementar AWS S3 Presigned URLs

**Objetivo**: Crear cliente S3 que genere URLs firmadas para permitir upload directo de archivos desde el cliente a AWS S3, reduciendo carga del backend.

**Tareas**:

- [ ] **2.1** - Agregar dependencias de AWS SDK
  - **Descripción**: Ejecutar `go get` para aws-sdk-go-v2/config, aws-sdk-go-v2/service/s3, aws-sdk-go-v2/credentials
  - **Archivos a crear/modificar**: `go.mod`, `go.sum`
  - **Criterio de aceptación**: Dependencias aparecen en go.mod, versión >= v2.24.0

- [ ] **2.2** - Crear cliente AWS S3
  - **Descripción**: Implementar S3Client con método GeneratePresignedURL(fileName, contentType, expiration) usando S3 Presign API
  - **Archivos a crear/modificar**:
    - `internal/infrastructure/storage/s3/client.go` (NUEVO)
  - **Criterio de aceptación**:
    - Struct S3Client con campos: s3Client *s3.Client, bucket string, region string
    - Método GeneratePresignedURL retorna URL string y error
    - URL válida por tiempo configurado (default 15 min)
    - Logging de URLs generadas (sin exponer URL completa)
  - 🔗 **Depende de**: Tarea 2.1

- [ ] **2.3** - Agregar configuración de S3
  - **Descripción**: Extender config.go con S3Config (Region, Bucket, PresignedURLExpiration) y agregar sección en config.yaml
  - **Archivos a crear/modificar**:
    - `internal/config/config.go`
    - `config/config.yaml`
  - **Criterio de aceptación**:
    - Struct S3Config con campos Region, Bucket string, PresignedURLExpiration time.Duration
    - config.yaml tiene sección s3 con valores: region=us-east-1, bucket=edugo-materials, expiration=15m

- [ ] **2.4** - Inicializar S3Client en main.go
  - **Descripción**: Instanciar S3Client en cmd/main.go usando config de Viper, inicializar AWS config con credenciales (env vars o IAM roles)
  - **Archivos a crear/modificar**:
    - `cmd/main.go`
  - **Criterio de aceptación**:
    - S3Client inicializado después de RabbitMQ
    - Log de inicialización exitosa con región y bucket
    - Manejo de error si falla carga de credenciales AWS
  - 🔗 **Depende de**: Tareas 2.2, 2.3

- [ ] **2.5** - Agregar S3Client al Container DI
  - **Descripción**: Extender Container struct con campo S3Client, inyectar en NewContainer()
  - **Archivos a crear/modificar**:
    - `internal/container/container.go`
  - **Criterio de aceptación**:
    - Campo `S3Client` en struct Container
    - MaterialService puede acceder a S3Client via container
  - 🔗 **Depende de**: Tarea 2.4

- [ ] **2.6** - Integrar S3 en MaterialService
  - **Descripción**: Inyectar S3Client en MaterialService, generar presigned URL en CreateMaterial() ANTES de insertar en DB, retornar URL en response
  - **Archivos a crear/modificar**:
    - `internal/application/service/material_service.go`
  - **Criterio de aceptación**:
    - Método CreateMaterial llama a s3Client.GeneratePresignedURL()
    - URL incluida en MaterialResponse DTO
    - Si falla generación, retornar error 500 (es crítico)
  - 🔗 **Depende de**: Tarea 2.5

- [ ] **2.7** - Actualizar MaterialResponse DTO
  - **Descripción**: Agregar campo UploadURL string con tag JSON en MaterialResponse DTO
  - **Archivos a crear/modificar**:
    - `internal/application/dto/material_dto.go`
  - **Criterio de aceptación**:
    - Campo `UploadURL string` con tag `json:"upload_url,omitempty"`
    - Documentación en comentario explicando que es presigned URL temporal

- [ ] **2.8** - Crear tests unitarios de S3Client
  - **Descripción**: Tests con mock de AWS SDK para validar generación de URLs, formato de URLs, expiración
  - **Archivos a crear/modificar**:
    - `internal/infrastructure/storage/s3/client_test.go` (NUEVO)
  - **Criterio de aceptación**:
    - Al menos 4 tests: URL generada exitosamente, formato URL válido, error de AWS SDK, expiración correcta
    - Tests pasan con `go test ./internal/infrastructure/storage/...`
  - 🔗 **Depende de**: Tarea 2.2

- [ ] **2.9** - Commit atómico de S3 Presigned URLs
  - **Descripción**: Crear commit con mensaje "feat: implementar generación de URLs firmadas S3"
  - **Archivos incluidos**: Todos los archivos modificados en Fase 2
  - **Criterio de aceptación**:
    - `go build ./...` compila sin errores
    - Tests pasan: `go test ./...`
    - Commit creado en branch feature/fase2-servicios
  - 🔗 **Depende de**: Tareas 2.1 a 2.8

**Completitud de Fase 2**: 0/9 tareas completadas

---

### Fase 3: Implementar Queries Complejas

**Objetivo**: Implementar consultas optimizadas en PostgreSQL (JOINs, CTEs, UPSERT) y MongoDB (aggregation pipelines) para mejorar performance y funcionalidad de servicios.

**Tareas**:

#### 📌 Subtarea 3.A: Queries PostgreSQL

- [ ] **3.1** - Crear índice en material_versions
  - **Descripción**: Crear script de migración SQL para agregar índices en tabla material_versions (material_id, version_number)
  - **Archivos a crear/modificar**:
    - `scripts/postgresql/05_indexes_material_versions.sql` (NUEVO)
  - **Criterio de aceptación**:
    - Script crea índices: `idx_material_versions_material_id`, `idx_material_versions_material_id_version`
    - Índices incluyen columna version_number DESC para ordenamiento
    - Script es idempotente (IF NOT EXISTS)

- [ ] **3.2** - Crear índice en materials
  - **Descripción**: Crear script SQL para agregar índice en materials.updated_at
  - **Archivos a crear/modificar**:
    - `scripts/postgresql/06_indexes_materials.sql` (NUEVO)
  - **Criterio de aceptación**:
    - Script crea índice: `idx_materials_updated_at DESC`
    - Idempotente (IF NOT EXISTS)

- [ ] **3.3** - Crear índice UNIQUE en user_progress
  - **Descripción**: Crear script SQL para índice UNIQUE compuesto (user_id, material_id) en user_progress, requerido para UPSERT
  - **Archivos a crear/modificar**:
    - `scripts/postgresql/07_indexes_user_progress.sql` (NUEVO)
  - **Criterio de aceptación**:
    - Índice UNIQUE: `idx_user_progress_user_material ON user_progress(user_id, material_id)`
    - Idempotente (IF NOT EXISTS)

- [ ] **3.4** - Implementar GetMaterialsWithVersions en MaterialRepository
  - **Descripción**: Agregar método GetMaterialsWithVersions() que ejecuta CTE con JSON aggregation para retornar materiales con array de versiones
  - **Archivos a crear/modificar**:
    - `internal/domain/repository/material_repository.go` (agregar método a interfaz)
    - `internal/infrastructure/persistence/postgres/repository/material_repository_impl.go`
  - **Criterio de aceptación**:
    - Query usa CTE `material_stats` con COUNT de versiones y MAX de version_number
    - json_agg construye array de versiones ordenado por version_number DESC
    - Retorna []Material con campo Versions []MaterialVersion poblado
    - Manejo de errores SQL con wrapping
  - 🔗 **Depende de**: Tareas 3.1, 3.2

- [ ] **3.5** - Implementar UpdateProgress con UPSERT en ProgressRepository
  - **Descripción**: Modificar UpdateProgress() para usar INSERT ... ON CONFLICT ... DO UPDATE con lógica condicional de GREATEST()
  - **Archivos a crear/modificar**:
    - `internal/infrastructure/persistence/postgres/repository/progress_repository_impl.go`
  - **Criterio de aceptación**:
    - Query usa ON CONFLICT (user_id, material_id) DO UPDATE
    - progress_percentage usa GREATEST para solo actualizar si nuevo > actual
    - status calculado automáticamente según porcentaje (CASE WHEN)
    - completed_at establecido solo primera vez que llega a 100%
    - RETURNING * para devolver registro actualizado
  - 🔗 **Depende de**: Tarea 3.3

- [ ] **3.6** - Actualizar MaterialService para usar GetMaterialsWithVersions
  - **Descripción**: Modificar método GetMaterials() en MaterialService para llamar al nuevo método del repositorio
  - **Archivos a crear/modificar**:
    - `internal/application/service/material_service.go`
  - **Criterio de aceptación**:
    - Servicio llama a repository.GetMaterialsWithVersions()
    - Response DTO incluye array de versiones
  - 🔗 **Depende de**: Tarea 3.4

- [ ] **3.7** - Actualizar ProgressService para usar UPSERT
  - **Descripción**: El método UpdateProgress() ya debería estar usando el repositorio actualizado, solo verificar integración
  - **Archivos a crear/modificar**:
    - `internal/application/service/progress_service.go` (verificar, puede no requerir cambios)
  - **Criterio de aceptación**:
    - Servicio llama a repository.UpdateProgress()
    - Lógica de negocio correcta (validaciones de porcentaje 0-100)
  - 🔗 **Depende de**: Tarea 3.5

#### 📌 Subtarea 3.B: Queries MongoDB

- [ ] **3.8** - Crear índices en assessment_attempts
  - **Descripción**: Crear función de migración Go que crea índices en colección assessment_attempts (assessment_id, user_id+submitted_at)
  - **Archivos a crear/modificar**:
    - `scripts/mongodb/indexes.go` (NUEVO) o ejecutar en main.go durante init
  - **Criterio de aceptación**:
    - Índice simple: `assessment_id` (ascending)
    - Índice compuesto: `user_id` (asc) + `submitted_at` (desc)
    - Logging de índices creados

- [ ] **3.9** - Implementar CalculateScoreWithFeedback en AssessmentRepository
  - **Descripción**: Agregar método que ejecuta aggregation pipeline con $lookup, $project y $switch para generar feedback dinámico
  - **Archivos a crear/modificar**:
    - `internal/domain/repository/assessment_repository.go` (agregar método a interfaz)
    - `internal/infrastructure/persistence/mongodb/repository/assessment_repository_impl.go`
  - **Criterio de aceptación**:
    - Pipeline con stages: $match (attempt_id), $lookup (assessments), $unwind, $project
    - Cálculo de percentage con $divide y $multiply
    - Feedback generado con $switch según rangos (>=90%, >=70%, >=50%, <50%)
    - Retorna struct AttemptWithFeedback
  - 🔗 **Depende de**: Tarea 3.8

- [ ] **3.10** - Implementar GetUserStatistics en StatsRepository
  - **Descripción**: Agregar método que ejecuta aggregation pipeline de 5 stages para calcular estadísticas del usuario
  - **Archivos a crear/modificar**:
    - `internal/domain/repository/stats_repository.go` (agregar método a interfaz)
    - `internal/infrastructure/persistence/mongodb/repository/stats_repository_impl.go`
  - **Criterio de aceptación**:
    - Pipeline: $match (user_id), $lookup (assessments), $unwind, $group (cálculos), $project (formato)
    - Cálculos: total_attempts, average_score, highest/lowest, overall_percentage
    - Array recent_attempts con $slice (últimos 10) y $sortArray
    - Retorna struct UserStatistics
  - 🔗 **Depende de**: Tarea 3.8

- [ ] **3.11** - Actualizar AssessmentService para usar CalculateScoreWithFeedback
  - **Descripción**: Modificar RecordAttempt() para llamar al nuevo método del repositorio después de insertar intento
  - **Archivos a crear/modificar**:
    - `internal/application/service/assessment_service.go`
  - **Criterio de aceptación**:
    - Después de insertar attempt, llamar a repository.CalculateScoreWithFeedback()
    - Feedback incluido en response DTO
  - 🔗 **Depende de**: Tarea 3.9

- [ ] **3.12** - Actualizar StatsService para usar GetUserStatistics
  - **Descripción**: Modificar método GetStatistics() para usar aggregation pipeline en lugar de queries simples
  - **Archivos a crear/modificar**:
    - `internal/application/service/stats_service.go`
  - **Criterio de aceptación**:
    - Servicio llama a repository.GetUserStatistics()
    - Response DTO con todas las estadísticas calculadas
  - 🔗 **Depende de**: Tarea 3.10

#### 📌 Subtarea 3.C: Testing

- [ ] **3.13** - Crear tests de integración para queries PostgreSQL
  - **Descripción**: Tests con Testcontainers PostgreSQL para validar GetMaterialsWithVersions y UpdateProgress
  - **Archivos a crear/modificar**:
    - `test/integration/postgres_queries_test.go` (NUEVO)
  - **Criterio de aceptación**:
    - Test GetMaterialsWithVersions: seed 2 materiales con 3 versiones cada uno, validar array de versiones
    - Test UpdateProgress INSERT: usuario sin progreso → crea registro
    - Test UpdateProgress UPDATE mayor: 50% → 75% → actualiza
    - Test UpdateProgress UPDATE menor: 75% → 50% → NO actualiza (GREATEST)
    - Test UpdateProgress completed: 100% → status=completed, completed_at NOT NULL
    - Todos los tests pasan
  - 🔗 **Depende de**: Tareas 3.4, 3.5

- [ ] **3.14** - Crear tests de integración para aggregations MongoDB
  - **Descripción**: Tests con Testcontainers MongoDB para validar CalculateScoreWithFeedback y GetUserStatistics
  - **Archivos a crear/modificar**:
    - `test/integration/mongodb_aggregations_test.go` (NUEVO)
  - **Criterio de aceptación**:
    - Test CalculateScoreWithFeedback: seed assessment + attempt con 80% → feedback="Buen trabajo!"
    - Test GetUserStatistics: seed 5 intentos de 2 assessments → validar totales, promedio, recent_attempts
    - Todos los tests pasan
  - 🔗 **Depende de**: Tareas 3.9, 3.10

- [ ] **3.15** - Commit atómico de Queries Complejas
  - **Descripción**: Crear commit con mensaje "feat: implementar consultas complejas en servicios"
  - **Archivos incluidos**: Todos los archivos modificados en Fase 3
  - **Criterio de aceptación**:
    - `go build ./...` compila sin errores
    - Tests unitarios pasan: `go test ./internal/...`
    - Tests de integración pasan: `go test ./test/integration/...`
    - Commit creado en branch feature/fase2-servicios
  - 🔗 **Depende de**: Tareas 3.1 a 3.14

**Completitud de Fase 3**: 0/15 tareas completadas

---

## 📊 Resumen de Dependencias

### Grafo de Dependencias

```mermaid
graph TD
    %% Fase 1: RabbitMQ
    T1_1[1.1: Deps RabbitMQ] --> T1_2[1.2: Events structs]
    T1_2 --> T1_3[1.3: Publisher impl]
    T1_1 --> T1_4[1.4: Config RabbitMQ]
    T1_3 --> T1_5[1.5: Init main.go]
    T1_4 --> T1_5
    T1_5 --> T1_6[1.6: Container DI]
    T1_6 --> T1_7[1.7: MaterialService]
    T1_6 --> T1_8[1.8: AssessmentService]
    T1_3 --> T1_9[1.9: Tests Publisher]
    T1_7 --> T1_10[1.10: Commit RabbitMQ]
    T1_8 --> T1_10
    T1_9 --> T1_10

    %% Fase 2: S3
    T2_1[2.1: Deps AWS] --> T2_2[2.2: S3Client impl]
    T2_1 --> T2_3[2.3: Config S3]
    T2_2 --> T2_4[2.4: Init S3 main.go]
    T2_3 --> T2_4
    T2_4 --> T2_5[2.5: Container S3]
    T2_5 --> T2_6[2.6: S3 in MaterialService]
    T2_6 --> T2_7[2.7: Update DTO]
    T2_2 --> T2_8[2.8: Tests S3]
    T2_7 --> T2_9[2.9: Commit S3]
    T2_8 --> T2_9

    %% Fase 3A: PostgreSQL
    T3_1[3.1: Índices material_versions] --> T3_4[3.4: GetMaterialsWithVersions]
    T3_2[3.2: Índices materials] --> T3_4
    T3_3[3.3: Índice user_progress] --> T3_5[3.5: UPSERT UpdateProgress]
    T3_4 --> T3_6[3.6: MaterialService update]
    T3_5 --> T3_7[3.7: ProgressService update]
    T3_4 --> T3_13[3.13: Tests PostgreSQL]
    T3_5 --> T3_13

    %% Fase 3B: MongoDB
    T3_8[3.8: Índices MongoDB] --> T3_9[3.9: CalculateScoreWithFeedback]
    T3_8 --> T3_10[3.10: GetUserStatistics]
    T3_9 --> T3_11[3.11: AssessmentService update]
    T3_10 --> T3_12[3.12: StatsService update]
    T3_9 --> T3_14[3.14: Tests MongoDB]
    T3_10 --> T3_14

    %% Commit final Fase 3
    T3_6 --> T3_15[3.15: Commit Queries]
    T3_7 --> T3_15
    T3_11 --> T3_15
    T3_12 --> T3_15
    T3_13 --> T3_15
    T3_14 --> T3_15
```

### Ruta Crítica

Las siguientes tareas están en la ruta crítica y deben completarse en orden:

**Fase 1**:
1. 1.1 → 1.2 → 1.3 → 1.5 → 1.6 → 1.7 → 1.10

**Fase 2**:
1. 2.1 → 2.2 → 2.4 → 2.5 → 2.6 → 2.7 → 2.9

**Fase 3 (PostgreSQL)**:
1. 3.1, 3.2 → 3.4 → 3.6 → 3.13 → 3.15
2. 3.3 → 3.5 → 3.7 → 3.13 → 3.15

**Fase 3 (MongoDB)**:
1. 3.8 → 3.9 → 3.11 → 3.14 → 3.15
2. 3.8 → 3.10 → 3.12 → 3.14 → 3.15

### Tareas Independientes (pueden ejecutarse en paralelo)

**Fase 1**:
- Tarea 1.4 (Config) puede hacerse en paralelo con 1.2-1.3
- Tarea 1.9 (Tests) puede hacerse después de 1.3, independiente de 1.7-1.8

**Fase 2**:
- Tarea 2.3 (Config) puede hacerse en paralelo con 2.2
- Tarea 2.8 (Tests) puede hacerse después de 2.2, independiente de 2.6-2.7

**Fase 3**:
- Tareas 3.1, 3.2, 3.3 (Índices PostgreSQL) son independientes entre sí
- Tarea 3.8 (Índices MongoDB) es independiente de tareas PostgreSQL
- Tareas 3.4 y 3.5 pueden ejecutarse en paralelo (diferentes repositorios)
- Tareas 3.9 y 3.10 pueden ejecutarse en paralelo (diferentes repositorios)

---

## 📈 Métricas del Plan

- **Total de fases**: 3
- **Total de tareas**: 34
- **Tareas con dependencias**: 26
- **Tareas independientes**: 8
- **Estimación**: 3-4 días de trabajo (según análisis arquitectónico)

**Distribución por fase**:
- Fase 1 (RabbitMQ): 10 tareas (~1-2 días)
- Fase 2 (S3): 9 tareas (~1 día)
- Fase 3 (Queries): 15 tareas (~1-2 días)

---

## 🎯 Estrategia de Ejecución Recomendada

1. **Primera iteración** (Día 1-2): Completar Fase 1 (RabbitMQ Messaging)
   - Fundamento para eventos de dominio
   - Permite probar integración con servicios
   - Commit atómico al finalizar

2. **Segunda iteración** (Día 2-3): Completar Fase 2 (AWS S3)
   - Funcionalidad independiente de Fase 1
   - Mejora inmediata en performance (upload directo)
   - Commit atómico al finalizar

3. **Tercera iteración** (Día 3-4): Completar Fase 3 (Queries Complejas)
   - Optimización de funcionalidades existentes
   - Subtareas PostgreSQL y MongoDB pueden ejecutarse en paralelo
   - Testing exhaustivo antes de commit
   - Commit atómico al finalizar

**Nota**: Usa `/03-execution phase-1` para ejecutar fases específicas, o `/03-execution` para ejecutar todo el plan.

---

## 📝 Notas Adicionales

### Consideraciones Importantes

1. **Variables de Entorno Requeridas**:
   ```bash
   # RabbitMQ
   RABBITMQ_URL=amqp://user:password@localhost:5672/

   # AWS S3
   AWS_REGION=us-east-1
   AWS_ACCESS_KEY_ID=your-key-id
   AWS_SECRET_ACCESS_KEY=your-secret-key
   ```

2. **Servicios Externos Necesarios**:
   - RabbitMQ server running (puerto 5672)
   - AWS S3 bucket creado y configurado
   - Credenciales AWS configuradas (env vars o IAM roles)

3. **Testing**:
   - Tests unitarios no requieren servicios externos (usan mocks)
   - Tests de integración usan Testcontainers (requiere Docker)
   - Ejecutar tests de integración solo en entorno local o CI

4. **Manejo de Errores**:
   - Eventos RabbitMQ: NO CRÍTICOS (log warning si fallan, request continúa)
   - S3 URLs: CRÍTICAS (retornar 500 si falla generación)
   - Queries: CRÍTICAS (retornar error apropiado al cliente)

5. **Performance**:
   - RabbitMQ: Publicación asíncrona (considerar go-routine si bloquea)
   - S3: Caché de cliente AWS (reutilizar conexiones)
   - Queries: EXPLAIN ANALYZE en PostgreSQL para validar índices
   - MongoDB: Verificar aggregation con explain() si es lenta

6. **Seguridad**:
   - NO hardcodear credenciales en código
   - Usar secrets manager en producción
   - TLS para RabbitMQ y MongoDB en producción
   - IAM roles para AWS en lugar de access keys cuando sea posible

---

## ✅ Checklist de Validación Final

Antes de considerar la Fase 2 completada, verificar:

- [ ] Los 3 commits atómicos creados (RabbitMQ, S3, Queries)
- [ ] `go build ./...` compila sin errores ni warnings
- [ ] `go test ./...` pasa todos los tests unitarios
- [ ] `go test ./test/integration/...` pasa tests de integración
- [ ] RabbitMQ publisher conecta y publica eventos correctamente
- [ ] S3 client genera URLs válidas que permiten upload
- [ ] Queries PostgreSQL retornan datos correctos (validar con query manual)
- [ ] Aggregations MongoDB retornan datos correctos (validar con mongo shell)
- [ ] Índices creados en PostgreSQL y MongoDB (verificar con EXPLAIN)
- [ ] No hay archivos duplicados o código muerto
- [ ] Logging apropiado en todos los componentes nuevos
- [ ] Configuración documentada en config.yaml con ejemplos

---

**Plan generado**: 2025-11-04
**Basado en**: Análisis arquitectónico de Fase 2
**Listo para**: Ejecución con `/03-execution`
