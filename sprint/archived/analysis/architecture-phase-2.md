# Arquitectura del Sistema - Fase 2: Completar TODOs de Servicios

## Resumen Ejecutivo

Esta fase se enfoca en completar la implementación de tres servicios fundamentales que quedaron pendientes en la arquitectura existente: **Messaging con RabbitMQ**, **Storage con AWS S3**, y **Queries Complejas** en los repositorios. La arquitectura mantiene el patrón **Clean Architecture (Hexagonal)** establecido en las fases anteriores.

**Nota**: Este análisis se enfoca en la Fase 2 del sprint.

## Tipo de Arquitectura

**Clean Architecture (Hexagonal)** - Continuación de la arquitectura establecida

La Fase 2 extiende la arquitectura existente con implementaciones de infraestructura pendientes, manteniendo la separación clara de capas:
- **Domain**: Entidades y contratos (sin cambios)
- **Application**: Servicios con lógica de negocio (modificaciones menores)
- **Infrastructure**: Nuevas implementaciones (RabbitMQ, S3, queries SQL/NoSQL)

## Diagrama de Arquitectura

```mermaid
flowchart TB
    subgraph "HTTP Layer"
        H1[Material Handler]
        H2[Assessment Handler]
    end
    
    subgraph "Application Layer"
        S1[Material Service]
        S2[Assessment Service]
        S3[Progress Service]
        S4[Stats Service]
    end
    
    subgraph "Infrastructure Layer - Nuevos Componentes"
        MSG[RabbitMQ Publisher]
        S3C[AWS S3 Client]
        
        subgraph "Persistence"
            PG[PostgreSQL]
            MG[MongoDB]
        end
    end
    
    subgraph "External Systems"
        RMQ[RabbitMQ Server]
        AWS[AWS S3 Bucket]
    end
    
    H1 --> S1
    H2 --> S2
    S1 --> MSG
    S2 --> MSG
    S1 --> S3C
    S1 --> PG
    S2 --> MG
    S3 --> PG
    S4 --> PG
    S4 --> MG
    
    MSG --> RMQ
    S3C --> AWS
    
    style MSG fill:#ffeb3b
    style S3C fill:#ffeb3b
    style PG fill:#4caf50
    style MG fill:#4caf50
```

## Componentes Principales

### 1. RabbitMQ Publisher (NUEVO)
- **Responsabilidad**: Publicar eventos de dominio a colas de mensajes
- **Tecnologías**: RabbitMQ client Go (amqp091-go)
- **Interacciones**: 
  - Material Service → publish `material_uploaded`
  - Assessment Service → publish `assessment_attempt_recorded`
- **Ubicación**: `internal/infrastructure/messaging/rabbitmq/publisher.go`

**Funcionalidades**:
- Conexión persistente a RabbitMQ
- Publicación de eventos con payload JSON
- Declaración automática de exchanges y queues
- Manejo de reconexión en caso de fallo

### 2. AWS S3 Client (NUEVO)
- **Responsabilidad**: Generar URLs firmadas para subida directa de archivos
- **Tecnologías**: AWS SDK for Go v2
- **Interacciones**: 
  - Material Service → solicita presigned URL
  - Frontend (indirecto) → usa URL para upload
- **Ubicación**: `internal/infrastructure/storage/s3/client.go`

**Funcionalidades**:
- Generación de presigned URLs con tiempo de expiración
- Configuración de bucket y región desde config
- Manejo de credenciales AWS (IAM roles o env vars)

### 3. Queries Complejas en Repositorios (ACTUALIZACIÓN)
- **Responsabilidad**: Implementar consultas avanzadas en PostgreSQL y MongoDB
- **Tecnologías**: SQL avanzado (JOINs, CTEs, UPSERT) + MongoDB aggregation pipeline
- **Interacciones**: 
  - Servicios → repositorios → bases de datos
- **Ubicación**: Repositorios existentes en `internal/infrastructure/persistence/`

**Queries a implementar**:
- **Material Repository**: Consulta de materiales con todas sus versiones
- **Assessment Repository**: Cálculo de puntajes y generación de feedback
- **Progress Repository**: UPSERT de progreso con actualización inteligente
- **Stats Repository**: Agregaciones complejas de estadísticas

### 4. Container DI (ACTUALIZACIÓN)
- **Responsabilidad**: Integrar nuevos componentes al contenedor de dependencias
- **Tecnologías**: Inyección de dependencias manual
- **Interacciones**: Orquestación de todos los componentes
- **Ubicación**: `internal/container/container.go`

**Nuevas dependencias a inyectar**:
- RabbitMQ Publisher
- AWS S3 Client
- Configuración extendida (RabbitMQ URL, S3 credentials)

## Patrones de Diseño Aplicados

### 1. Repository Pattern (existente, extendido)
**Descripción**: Los repositorios encapsulan la lógica de acceso a datos
**Justificación**: Permite agregar queries complejas sin modificar contratos del dominio
**Implementación en Fase 2**:
- Queries SQL avanzadas en PostgreSQL repositories
- Aggregation pipelines en MongoDB repositories

### 2. Publisher Pattern (nuevo)
**Descripción**: Publicación de eventos de dominio a sistemas externos
**Justificación**: Desacoplamiento entre servicios y sistemas de mensajería
**Implementación en Fase 2**:
- Interface `MessagePublisher` en dominio (opcional)
- Implementación concreta `RabbitMQPublisher` en infraestructura
- Servicios publican eventos sin conocer RabbitMQ

### 3. Dependency Injection (existente, extendido)
**Descripción**: Inyección de dependencias a través del Container
**Justificación**: Facilita testing y desacoplamiento
**Implementación en Fase 2**:
- Agregar Publisher y S3 Client al Container
- Inyectar en servicios que los requieren

### 4. Presigned URL Pattern (nuevo)
**Descripción**: Generación de URLs temporales para acceso directo a S3
**Justificación**: Evita que archivos pasen por el backend (menor carga)
**Implementación en Fase 2**:
- Cliente S3 genera URLs con expiración
- Frontend sube archivos directamente a S3
- Backend solo genera la URL y guarda metadata

## Stack Tecnológico Recomendado

### Backend (sin cambios)
- **Framework**: Gin (HTTP server)
- **Lenguaje**: Go 1.21+
- **Logger**: Zap (edugo-shared/logger)
- **Autenticación**: JWT (edugo-shared/auth)

### Base de datos (sin cambios)
- **Relacional**: PostgreSQL (lib/pq driver)
- **NoSQL**: MongoDB (mongo-driver)

### Infraestructura (nuevos servicios)
- **Messaging**: RabbitMQ
  - Cliente: `github.com/rabbitmq/amqp091-go`
  - Exchanges: topic exchange para routing flexible
  - Queues: durable para persistencia de mensajes
- **Storage**: AWS S3
  - SDK: `github.com/aws/aws-sdk-go-v2`
  - Presigned URLs con expiración de 15 minutos
  - Bucket: configurado en `config/config.yaml`

### Configuración (actualización)
- **Config Management**: Viper (existente)
- **Nuevas variables**:
  ```yaml
  rabbitmq:
    url: "amqp://user:password@host:5672/"
    exchange: "edugo_events"
  
  s3:
    region: "us-east-1"
    bucket: "edugo-materials"
    presigned_url_expiration: "15m"
  ```

## Consideraciones No Funcionales

### Escalabilidad
- **RabbitMQ**: 
  - Mensajes persistentes (survive broker restart)
  - Publisher confirms para garantizar entrega
  - Conexión pool para múltiples go-routines
- **S3**: 
  - Upload directo reduce carga del backend
  - Presigned URLs permiten escalado horizontal sin estado compartido
- **Queries**: 
  - Índices optimizados en PostgreSQL
  - Projection en MongoDB para reducir transferencia de datos

### Seguridad
- **RabbitMQ**: 
  - Credenciales en variables de entorno
  - Conexión TLS en producción
- **S3**: 
  - Presigned URLs con tiempo de expiración corto (15 min)
  - IAM roles para permisos granulares
  - Bucket policies para restringir acceso
- **Queries**: 
  - Prepared statements para prevenir SQL injection
  - Sanitización de inputs en servicios

### Performance
- **RabbitMQ**: 
  - Publicación asíncrona (no bloquea request HTTP)
  - Batch publishing para múltiples eventos
- **S3**: 
  - URLs generadas on-demand (no pre-generadas)
  - Caché de cliente S3 para reutilización de conexiones
- **Queries**: 
  - CTEs para queries legibles y optimizadas
  - Aggregation pipeline en MongoDB con límites

### Mantenibilidad
- **RabbitMQ**: 
  - Centralizador de eventos en `events.go`
  - Tipos de eventos bien definidos
  - Logging de eventos publicados
- **S3**: 
  - Cliente S3 encapsulado
  - Configuración centralizada
  - Fácil cambio a otro storage provider (interfaz)
- **Queries**: 
  - Queries en métodos de repositorio
  - Tests de integración para validar queries
  - Comentarios en SQL/aggregations complejas

## Diagrama de Flujo de Eventos

```mermaid
sequenceDiagram
    participant FE as Frontend
    participant H as Material Handler
    participant S as Material Service
    participant S3 as S3 Client
    participant MQ as RabbitMQ Publisher
    participant DB as PostgreSQL
    
    FE->>H: POST /materials (metadata)
    H->>S: CreateMaterial(dto)
    S->>S3: GeneratePresignedURL(fileName)
    S3-->>S: presignedURL
    S->>DB: INSERT material metadata
    DB-->>S: materialID
    S->>MQ: Publish(material_uploaded)
    MQ-->>S: ACK
    S-->>H: MaterialResponse + presignedURL
    H-->>FE: 201 Created + presignedURL
    FE->>S3: PUT file to presignedURL
    S3-->>FE: 200 OK
```

## Cambios en Componentes Existentes

### Material Service
**Modificaciones**:
- Inyectar `S3Client` y `MessagePublisher`
- Método `CreateMaterial`: generar presigned URL antes de retornar
- Método `CreateMaterial`: publicar evento `material_uploaded` después de DB insert
- Método `GetMaterialsWithVersions`: usar query compleja del repositorio

### Assessment Service
**Modificaciones**:
- Inyectar `MessagePublisher`
- Método `RecordAttempt`: publicar evento `assessment_attempt_recorded`
- Método `CalculateScore`: usar query compleja para generar feedback

### Progress Service
**Modificaciones**:
- Método `UpdateProgress`: usar UPSERT en lugar de UPDATE/INSERT separados

### Stats Service
**Modificaciones**:
- Método `GetStatistics`: usar aggregation pipeline compleja

## Nuevos Archivos a Crear

```
internal/
├── infrastructure/
│   ├── messaging/
│   │   ├── rabbitmq/
│   │   │   └── publisher.go          ← NUEVO
│   │   └── events.go                 ← NUEVO (definiciones de eventos)
│   └── storage/
│       └── s3/
│           └── client.go             ← NUEVO
config/
└── config.yaml                       ← ACTUALIZAR (agregar RabbitMQ + S3)
```

## Archivos a Modificar

```
cmd/main.go                                           ← Inicializar RabbitMQ y S3
internal/container/container.go                       ← Agregar nuevas dependencias
internal/config/config.go                             ← Agregar structs RabbitMQ y S3
internal/application/service/material_service.go      ← S3 + eventos
internal/application/service/assessment_service.go    ← Eventos
internal/application/service/progress_service.go      ← UPSERT
internal/application/service/stats_service.go         ← Aggregations
internal/infrastructure/persistence/postgres/repository/
  └── material_repository_impl.go                     ← Query con JOIN
internal/infrastructure/persistence/mongodb/repository/
  └── assessment_repository_impl.go                   ← Aggregation pipeline
```

## Dependencias Nuevas (go.mod)

```go
require (
    github.com/rabbitmq/amqp091-go v1.9.0  // Cliente RabbitMQ
    github.com/aws/aws-sdk-go-v2 v1.24.0   // AWS SDK Core
    github.com/aws/aws-sdk-go-v2/service/s3 v1.48.0  // S3 Service
    github.com/aws/aws-sdk-go-v2/config v1.26.0      // AWS Config
)
```

## Plan de Testing

### Tests Unitarios (por componente)
- **RabbitMQ Publisher**: 
  - Mock de conexión RabbitMQ
  - Validar payload JSON
  - Verificar declaración de exchanges
- **S3 Client**: 
  - Mock de AWS SDK
  - Validar formato de presigned URLs
  - Verificar expiración de URLs
- **Queries Complejas**: 
  - Tests de integración con Testcontainers
  - Validar resultados de queries
  - Verificar performance de índices

### Tests de Integración
- Flujo completo: Material creation → S3 URL → Event published
- Flujo Assessment → Event → Score calculation
- UPSERT de progreso con datos existentes/nuevos

## Riesgos y Mitigaciones

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|--------------|---------|------------|
| RabbitMQ no disponible | Media | Alto | Implementar circuit breaker, logging de eventos fallidos |
| S3 timeout en generación de URLs | Baja | Medio | Timeout configurado, retry logic, caché de cliente |
| Queries complejas lentas | Media | Alto | Índices optimizados, EXPLAIN ANALYZE, paginación |
| Credenciales AWS expuestas | Baja | Crítico | Usar IAM roles en producción, secrets manager |
| Mensajes duplicados en RabbitMQ | Media | Medio | Idempotencia en consumidores (fuera de esta fase) |

## Siguiente Paso Recomendado

Ejecutar `/02-planning` para generar el plan de tareas detallado de esta fase.
