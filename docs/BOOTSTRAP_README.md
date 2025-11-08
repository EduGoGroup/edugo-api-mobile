# Sistema de Bootstrap de Infraestructura

## Resumen Ejecutivo

El sistema de bootstrap de infraestructura es un módulo que centraliza y simplifica la inicialización de todos los recursos de infraestructura de la aplicación EduGo API Mobile. Proporciona una forma consistente, testeable y flexible de gestionar dependencias externas.

## ¿Qué Problema Resuelve?

### Antes del Bootstrap

```go
// main.go - ~180 líneas
func main() {
    // Inicialización manual de cada recurso
    log := logger.NewZapLogger(...)
    pgDB, err := database.InitPostgreSQL(...)
    if err != nil { log.Fatal(...) }
    defer pgDB.Close()
    
    mongoDB, err := database.InitMongoDB(...)
    if err != nil { log.Fatal(...) }
    defer mongoDB.Client().Disconnect(ctx)
    
    publisher, err := rabbitmq.NewRabbitMQPublisher(...)
    if err != nil {
        log.Warn(...)
        publisher = rabbitmq.NewNoopPublisher(...)
    }
    defer publisher.Close()
    
    s3Client, err := s3.NewS3Client(...)
    if err != nil {
        log.Warn(...)
        s3Client = s3.NewNoopS3Storage(...)
    }
    
    // Crear container con 6+ parámetros
    c := container.NewContainer(log, pgDB, mongoDB, publisher, s3Client, jwtSecret)
    
    // Setup router y servidor
    r := router.SetupRouter(c)
    r.Run(":8080")
}
```

**Problemas:**
- ❌ Código repetitivo y verboso
- ❌ Difícil de testear (no hay inyección de dependencias)
- ❌ Manejo inconsistente de recursos opcionales
- ❌ Múltiples defer statements
- ❌ Difícil de mantener

### Después del Bootstrap

```go
// main.go - <50 líneas
func main() {
    ctx := context.Background()
    cfg, _ := config.Load()
    
    // Bootstrap inicializa todo
    b := bootstrap.New(cfg)
    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer cleanup()
    
    // Crear container y servidor
    c := container.NewContainer(resources)
    r := router.SetupRouter(c)
    r.Run(":8080")
}
```

**Beneficios:**
- ✅ Código limpio y conciso
- ✅ Fácil de testear (inyección de mocks)
- ✅ Manejo consistente de recursos opcionales
- ✅ Un solo defer cleanup()
- ✅ Fácil de mantener

## Características Principales

### 1. Inicialización Unificada

Una sola llamada inicializa todos los recursos:

```go
resources, cleanup, err := b.InitializeInfrastructure(ctx)
```

### 2. Recursos Opcionales

Marca recursos como opcionales para desarrollo sin infraestructura completa:

```go
b := bootstrap.New(cfg,
    bootstrap.WithOptionalResource("rabbitmq"),
    bootstrap.WithOptionalResource("s3"),
)
```

### 3. Inyección de Mocks

Facilita testing con mocks:

```go
b := bootstrap.New(cfg,
    bootstrap.WithLogger(mockLogger),
    bootstrap.WithPostgreSQL(mockDB),
    bootstrap.WithMongoDB(mockMongoDB),
    bootstrap.WithRabbitMQ(mockPublisher),
    bootstrap.WithS3Client(mockS3),
)
```

### 4. Gestión de Ciclo de Vida

Cleanup automático de todos los recursos:

```go
defer cleanup() // Cierra todo en orden correcto
```

### 5. Degradación Graciosa

Cuando un recurso opcional falla, se inyecta una implementación noop:

```go
// Si RabbitMQ no está disponible
resources.RabbitMQPublisher.Publish(...) // No falla, solo registra debug
```

## Recursos Gestionados

| Recurso | Tipo | ¿Opcional? | Implementación Noop |
|---------|------|-----------|---------------------|
| Logger | Zap Logger | ❌ No | N/A |
| PostgreSQL | SQL Database | ❌ No | N/A |
| MongoDB | NoSQL Database | ❌ No | N/A |
| RabbitMQ | Message Broker | ✅ Sí | NoopPublisher |
| S3 | Object Storage | ✅ Sí | NoopS3Storage |

## Casos de Uso

### Desarrollo Local

Ejecuta la aplicación sin RabbitMQ o S3:

```go
b := bootstrap.New(cfg,
    bootstrap.WithOptionalResource("rabbitmq"),
    bootstrap.WithOptionalResource("s3"),
)
```

### Testing Unitario

Inyecta mocks de todos los recursos:

```go
b := bootstrap.New(cfg,
    bootstrap.WithLogger(mockLogger),
    bootstrap.WithPostgreSQL(mockDB),
    bootstrap.WithMongoDB(mockMongoDB),
    bootstrap.WithRabbitMQ(mockPublisher),
    bootstrap.WithS3Client(mockS3),
)
```

### Testing de Integración

Usa recursos reales o testcontainers:

```go
b := bootstrap.New(cfg)
resources, cleanup, err := b.InitializeInfrastructure(ctx)
require.NoError(t, err)
defer cleanup()
```

### Producción

Inicializa todos los recursos:

```go
b := bootstrap.New(cfg)
resources, cleanup, err := b.InitializeInfrastructure(ctx)
if err != nil {
    log.Fatal(err)
}
defer cleanup()
```

## Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│                         main.go                              │
│                                                              │
│  cfg := config.Load()                                        │
│  b := bootstrap.New(cfg, options...)                         │
│  resources, cleanup, err := b.InitializeInfrastructure(ctx)  │
│  defer cleanup()                                             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Bootstrap System                          │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Factories   │  │   Lifecycle  │  │    Config    │      │
│  │              │  │   Manager    │  │   Options    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       Resources                              │
│                                                              │
│  • Logger (Zap)                                              │
│  • PostgreSQL (*sql.DB)                                      │
│  • MongoDB (*mongo.Database)                                 │
│  • RabbitMQ Publisher (rabbitmq.Publisher)                   │
│  • S3 Client (s3.S3Storage)                                  │
│  • JWT Secret (string)                                       │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       Container                              │
│                                                              │
│  • Services                                                  │
│  • Repositories                                              │
│  • Handlers                                                  │
└─────────────────────────────────────────────────────────────┘
```

## Estructura de Archivos

```
internal/bootstrap/
├── bootstrap.go              # Orquestador principal
├── interfaces.go             # Interfaces de recursos
├── config.go                 # Configuración y opciones
├── factories.go              # Factories de recursos
├── lifecycle.go              # Gestión de ciclo de vida
├── bootstrap_test.go         # Tests unitarios
├── lifecycle_test.go         # Tests de lifecycle
├── bootstrap_integration_test.go  # Tests de integración
├── INTEGRATION_TESTS.md      # Documentación de tests
└── noop/
    ├── publisher.go          # Implementación noop de RabbitMQ
    └── storage.go            # Implementación noop de S3
```

## Documentación

### Para Empezar

1. **[BOOTSTRAP_INDEX.md](BOOTSTRAP_INDEX.md)** - Índice completo de documentación
2. **[BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md)** - Guía completa de uso
3. **[../QUICKSTART.md](../QUICKSTART.md)** - Quick start guide

### Para Migrar

4. **[BOOTSTRAP_MIGRATION_GUIDE.md](BOOTSTRAP_MIGRATION_GUIDE.md)** - Guía de migración paso a paso

### Para Testing

5. **[../internal/bootstrap/INTEGRATION_TESTS.md](../internal/bootstrap/INTEGRATION_TESTS.md)** - Guía de testing

### Para Configuración

6. **[../config/OPTIONAL_RESOURCES.md](../config/OPTIONAL_RESOURCES.md)** - Recursos opcionales

## Quick Start

### 1. Instalación

El sistema de bootstrap ya está incluido en el proyecto. No requiere instalación adicional.

### 2. Uso Básico

```go
package main

import (
    "context"
    "log"
    
    "github.com/EduGoGroup/edugo-api-mobile/internal/bootstrap"
    "github.com/EduGoGroup/edugo-api-mobile/internal/config"
    "github.com/EduGoGroup/edugo-api-mobile/internal/container"
)

func main() {
    ctx := context.Background()
    cfg, _ := config.Load()
    
    b := bootstrap.New(cfg)
    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer cleanup()
    
    c := container.NewContainer(resources)
    // ... resto de tu aplicación
}
```

### 3. Desarrollo sin Infraestructura Completa

```yaml
# config/config-local.yaml
infrastructure:
  optional_resources:
    - rabbitmq
    - s3
```

### 4. Testing

```go
func TestMyFeature(t *testing.T) {
    cfg := testConfig()
    b := bootstrap.New(cfg,
        bootstrap.WithLogger(mockLogger),
        bootstrap.WithPostgreSQL(mockDB),
        bootstrap.WithMongoDB(mockMongoDB),
        bootstrap.WithRabbitMQ(mockPublisher),
        bootstrap.WithS3Client(mockS3),
    )
    
    resources, cleanup, _ := b.InitializeInfrastructure(context.Background())
    defer cleanup()
    
    // Tests...
}
```

## Métricas de Impacto

### Reducción de Código

- **main.go**: De ~180 líneas a <50 líneas (72% reducción)
- **Container**: De 6+ parámetros a 1 parámetro (83% reducción)
- **Cleanup**: De 5+ defer statements a 1 defer (80% reducción)

### Mejora en Testabilidad

- **Antes**: Difícil inyectar mocks, requiere modificar código
- **Después**: Inyección simple con opciones funcionales

### Flexibilidad

- **Antes**: Recursos opcionales manejados manualmente con if/else
- **Después**: Configuración declarativa de recursos opcionales

## Compatibilidad

- ✅ **Sin breaking changes**: Las funciones existentes se mantienen
- ✅ **Backward compatible**: Migración gradual posible
- ✅ **Tests existentes**: Continúan funcionando

## Próximos Pasos

1. **Lee la documentación**: Empieza con [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md)
2. **Prueba en desarrollo**: Usa recursos opcionales para desarrollo local
3. **Escribe tests**: Usa inyección de mocks para testing
4. **Migra código existente**: Sigue [BOOTSTRAP_MIGRATION_GUIDE.md](BOOTSTRAP_MIGRATION_GUIDE.md)

## Soporte

- **Documentación completa**: [BOOTSTRAP_INDEX.md](BOOTSTRAP_INDEX.md)
- **Troubleshooting**: [BOOTSTRAP_MIGRATION_GUIDE.md - Troubleshooting](BOOTSTRAP_MIGRATION_GUIDE.md#troubleshooting)
- **Ejemplos**: [BOOTSTRAP_USAGE.md - Ejemplos Prácticos](BOOTSTRAP_USAGE.md#ejemplos-prácticos)

## Contribuir

Para contribuir al sistema de bootstrap:

1. Lee la especificación técnica en [../.kiro/specs/infrastructure-bootstrap-refactor/](../.kiro/specs/infrastructure-bootstrap-refactor/)
2. Revisa el diseño en [design.md](../.kiro/specs/infrastructure-bootstrap-refactor/design.md)
3. Sigue los patrones existentes en el código
4. Escribe tests para nuevas funcionalidades
5. Actualiza la documentación

## Licencia

MIT
