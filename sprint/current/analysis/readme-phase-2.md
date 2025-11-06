# Resumen del Análisis - Fase 2: Completar TODOs de Servicios

**Alcance**: Análisis de la Fase 2

## Objetivo de la Fase

Completar la implementación de tres servicios fundamentales que quedaron pendientes en la arquitectura existente:

1. **RabbitMQ Messaging**: Sistema de publicación de eventos de dominio
2. **AWS S3 Storage**: Generación de URLs firmadas para upload directo de archivos
3. **Queries Complejas**: Consultas optimizadas en PostgreSQL y MongoDB con agregaciones avanzadas

**Esfuerzo estimado**: 3-4 días
**Commits esperados**: 3 commits atómicos (uno por subtarea)

---

## Arquitectura Propuesta

La Fase 2 extiende la **Clean Architecture (Hexagonal)** existente con nuevas implementaciones de infraestructura, manteniendo la separación clara de capas:

### Componentes Nuevos

1. **RabbitMQ Publisher** (`internal/infrastructure/messaging/rabbitmq/`)
   - Publicación de eventos `material_uploaded` y `assessment_attempt_recorded`
   - Conexión persistente con manejo de reconexión
   - Publicación asíncrona (no bloquea HTTP response)

2. **AWS S3 Client** (`internal/infrastructure/storage/s3/`)
   - Generación de presigned URLs con expiración de 15 minutos
   - Upload directo desde cliente a S3 (reduce carga del backend)
   - Configuración de bucket y credenciales desde config.yaml

3. **Queries Complejas en Repositorios** (actualización de repositorios existentes)
   - PostgreSQL: JOINs, CTEs, UPSERT con lógica condicional
   - MongoDB: Aggregation pipelines con múltiples stages

### Integración con Arquitectura Existente

```
HTTP Handlers → Services → [NUEVOS] Publisher + S3Client + Queries Complejas → External Systems
```

No se modifica:
- Contratos de API (endpoints mantienen misma firma)
- Entidades de dominio
- DTOs de aplicación

Solo se extiende:
- Servicios (inyectar nuevas dependencias)
- Repositorios (agregar métodos con queries complejas)
- Container DI (agregar RabbitMQ y S3)

---

## Componentes Principales

### 1. RabbitMQ Publisher
**Responsabilidad**: Publicar eventos de dominio a colas de mensajes

**Eventos a implementar**:
- `material_uploaded`: Cuando se crea un material
  - Payload: `{ material_id, title, content_type, uploaded_at }`
- `assessment_attempt_recorded`: Cuando se registra un intento de evaluación
  - Payload: `{ attempt_id, user_id, assessment_id, score, submitted_at }`

**Características**:
- Declaración automática de exchanges y queues
- Publisher confirms para garantizar entrega
- Logging de eventos publicados
- **No crítico**: Si falla publicación, request HTTP continúa (solo log warning)

---

### 2. AWS S3 Client
**Responsabilidad**: Generar URLs firmadas para upload directo de archivos

**Flujo de uso**:
1. Cliente solicita crear material → POST /materials
2. Backend genera presigned URL y retorna en response
3. Cliente hace PUT del archivo directamente a S3 usando la URL
4. Archivo nunca pasa por el backend (reduce carga y latencia)

**Configuración**:
```yaml
s3:
  region: "us-east-1"
  bucket: "edugo-materials"
  presigned_url_expiration: "15m"
```

**Seguridad**:
- URLs con tiempo de expiración corto (15 minutos)
- IAM roles para permisos granulares
- Bucket policies para restringir acceso

---

### 3. Queries Complejas

#### PostgreSQL Queries

**a) Materiales con Versiones** (`GetMaterialsWithVersions`)
- CTE para agregar metadatos (count de versiones, última versión)
- JSON aggregation para construir array de versiones
- Ordenamiento por fecha de actualización
- **Índices requeridos**: material_id, version_number, updated_at

**b) UPSERT de Progreso** (`UpdateProgress`)
- INSERT ... ON CONFLICT ... DO UPDATE
- Lógica condicional: solo actualiza si nuevo progreso > progreso actual
- Cálculo automático de `status` según porcentaje
- Establece `completed_at` solo la primera vez que llega a 100%
- **Índice UNIQUE requerido**: (user_id, material_id)

#### MongoDB Aggregations

**a) Cálculo de Score con Feedback** (`CalculateScoreWithFeedback`)
- Lookup de assessment original
- Cálculo de porcentaje
- Generación dinámica de feedback según rango:
  - >= 90%: "Excelente trabajo!"
  - >= 70%: "Buen trabajo!"
  - >= 50%: "Aprobado, revisa temas fallidos"
  - < 50%: "Necesitas repasar el material"

**b) Estadísticas de Usuario** (`GetUserStatistics`)
- Aggregation con 5 stages
- Cálculos: total intentos, promedio, máximo/mínimo, overall percentage
- Array de últimos 10 intentos ordenados por fecha
- **Índices requeridos**: (user_id, submitted_at), (user_id, assessment_id)

---

## Modelo de Datos

**SIN CAMBIOS EN EL ESQUEMA** - Solo optimización de queries sobre tablas/colecciones existentes.

### Índices Nuevos a Crear

**PostgreSQL**:
```sql
CREATE INDEX idx_material_versions_material_id ON material_versions(material_id);
CREATE INDEX idx_materials_updated_at ON materials(updated_at DESC);
CREATE UNIQUE INDEX idx_user_progress_user_material ON user_progress(user_id, material_id);
```

**MongoDB**:
```javascript
db.assessment_attempts.createIndex({ assessment_id: 1 });
db.assessment_attempts.createIndex({ user_id: 1, submitted_at: -1 });
```

---

## Flujo Principal

### Creación de Material con S3 Upload

```
1. POST /materials (metadata) → Handler
2. Service genera ID de material
3. S3Client genera presigned URL ← AWS SDK
4. Repository inserta metadata en PostgreSQL
5. Publisher publica evento material_uploaded → RabbitMQ
6. Handler retorna 201 + material + presignedURL
7. Cliente hace PUT a presignedURL → S3 (upload directo)
```

**Tiempo estimado**: ~300ms (sin considerar upload a S3 que es directo del cliente)

**Puntos críticos**:
- Generación de presigned URL (latencia variable de AWS API)
- Publicación a RabbitMQ (asíncrona, no bloquea response)

---

## Stack Tecnológico

### Backend (sin cambios)
- **Framework**: Gin
- **Lenguaje**: Go 1.21+
- **Logger**: Zap (edugo-shared)

### Nuevas Dependencias
```go
require (
    github.com/rabbitmq/amqp091-go v1.9.0         // RabbitMQ client
    github.com/aws/aws-sdk-go-v2 v1.24.0          // AWS SDK
    github.com/aws/aws-sdk-go-v2/service/s3 v1.48.0
)
```

### Infraestructura Externa
- **RabbitMQ**: Servidor de mensajería (puerto 5672)
- **AWS S3**: Storage de archivos (región configurable)

---

## Consideraciones Importantes

### Escalabilidad
- **RabbitMQ**: Mensajes persistentes, publisher confirms
- **S3**: Upload directo reduce carga del backend, permite escalado horizontal
- **Queries**: Índices optimizados, paginación donde aplique

### Seguridad
- **RabbitMQ**: Credenciales en env vars, TLS en producción
- **S3**: Presigned URLs con expiración corta, IAM roles
- **Queries**: Prepared statements, sanitización de inputs

### Performance
- **RabbitMQ**: Publicación asíncrona (no bloquea)
- **S3**: Caché de cliente para reutilizar conexiones
- **Queries**: CTEs optimizadas, aggregation con projection

### Mantenibilidad
- **RabbitMQ**: Eventos centralizados en `events.go`
- **S3**: Cliente encapsulado, fácil cambio a otro provider
- **Queries**: Tests de integración con Testcontainers

---

## Riesgos Identificados

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|--------------|---------|------------|
| RabbitMQ no disponible | Media | Alto | Circuit breaker, logging de eventos fallidos, evento no crítico |
| S3 timeout en generación de URLs | Baja | Medio | Timeout configurado, retry logic, caché de cliente |
| Queries complejas lentas | Media | Alto | Índices optimizados, EXPLAIN ANALYZE, paginación |
| Credenciales AWS expuestas | Baja | Crítico | IAM roles en producción, secrets manager, no hardcodear |

---

## Archivos a Crear

```
internal/infrastructure/messaging/rabbitmq/publisher.go  (nuevo)
internal/infrastructure/messaging/events.go              (nuevo)
internal/infrastructure/storage/s3/client.go             (nuevo)
```

## Archivos a Modificar

```
cmd/main.go                                  (inicializar RabbitMQ y S3)
internal/container/container.go              (inyectar nuevas dependencias)
internal/config/config.go                    (structs RabbitMQ y S3)
config/config.yaml                           (configuración)
internal/application/service/*_service.go    (inyectar publisher y S3)
internal/infrastructure/persistence/*/repository/*_impl.go  (queries complejas)
```

---

## Commits Atómicos Esperados

### Commit 1: RabbitMQ Messaging
```
feat: implementar messaging RabbitMQ para eventos

- Crear RabbitMQPublisher con conexión persistente
- Definir eventos material_uploaded y assessment_attempt_recorded
- Integrar publisher en MaterialService y AssessmentService
- Configurar RabbitMQ en config.yaml y main.go
- Agregar publisher al Container DI
```

### Commit 2: S3 Presigned URLs
```
feat: implementar generación de URLs firmadas S3

- Crear S3Client con generación de presigned URLs
- Integrar S3Client en MaterialService.CreateMaterial
- Configurar AWS S3 en config.yaml
- Agregar S3Client al Container DI
- Retornar presignedURL en MaterialResponse
```

### Commit 3: Queries Complejas
```
feat: implementar consultas complejas en servicios

- Implementar GetMaterialsWithVersions con CTE en PostgreSQL
- Implementar UpdateProgress con UPSERT en PostgreSQL
- Implementar CalculateScoreWithFeedback con aggregation en MongoDB
- Implementar GetUserStatistics con pipeline en MongoDB
- Crear índices necesarios en PostgreSQL y MongoDB
- Actualizar servicios para usar nuevos métodos de repositorio
```

---

## Siguientes Pasos Recomendados

1. **Ejecutar `/02-planning`** para generar el plan de tareas detallado
2. **Ejecutar `/03-execution phase-2`** para implementar las tareas
3. **Ejecutar `/04-review`** para consolidar el estado y validar completitud

---

## Testing

### Tests Unitarios
- Mock de RabbitMQ connection
- Mock de AWS SDK
- Validación de payload de eventos

### Tests de Integración (Testcontainers)
- Flujo completo: Material → S3 URL → Event published
- UPSERT con datos existentes/nuevos
- Aggregations con datos reales

---

📁 **Documentación completa**: Ver archivos `architecture-phase-2.md`, `data-model-phase-2.md`, y `process-diagram-phase-2.md` en esta carpeta.

---

**Análisis generado**: 2025-11-04
**Fase**: 2 de 4
**Estado**: ✅ Listo para planning
