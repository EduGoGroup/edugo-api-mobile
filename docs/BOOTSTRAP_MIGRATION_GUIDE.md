# Guía de Migración al Sistema de Bootstrap

Esta guía te ayudará a migrar código existente que inicializa recursos de infraestructura directamente al nuevo sistema de bootstrap modular.

## Tabla de Contenidos

- [¿Por Qué Migrar?](#por-qué-migrar)
- [Cambios Principales](#cambios-principales)
- [Migración Paso a Paso](#migración-paso-a-paso)
- [Ejemplos de Código](#ejemplos-de-código)
- [Testing con Bootstrap](#testing-con-bootstrap)
- [Recursos Opcionales](#recursos-opcionales)
- [Troubleshooting](#troubleshooting)

## ¿Por Qué Migrar?

El nuevo sistema de bootstrap proporciona:

1. **Código más limpio**: `main.go` reducido de ~180 líneas a <50 líneas
2. **Mejor testabilidad**: Inyección fácil de mocks para testing
3. **Recursos opcionales**: Desarrollo sin infraestructura completa
4. **Gestión de ciclo de vida**: Cleanup automático de recursos
5. **Consistencia**: Todos los recursos siguen el mismo patrón
6. **Degradación graciosa**: La app continúa funcionando con funcionalidad reducida

## Cambios Principales

### 1. Inicialización de Recursos

**Antes:**
```go
// Inicialización manual de cada recurso
log := logger.NewZapLogger(cfg.Logger.Level, cfg.Logger.Format)
pgDB, err := database.InitPostgreSQL(ctx, cfg, log)
mongoDB, err := database.InitMongoDB(ctx, cfg, log)
publisher, err := rabbitmq.NewRabbitMQPublisher(cfg.Messaging.RabbitMQ.URL, "events", log)
s3Client, err := s3.NewS3Client(ctx, cfg.Storage.S3, log)
```

**Después:**
```go
// Bootstrap inicializa todos los recursos
b := bootstrap.New(cfg)
resources, cleanup, err := b.InitializeInfrastructure(ctx)
defer cleanup()
```

### 2. Firma del Container

**Antes:**
```go
func NewContainer(
    logger logger.Logger,
    pgDB *sql.DB,
    mongoDB *mongo.Database,
    publisher rabbitmq.Publisher,
    s3Client s3.S3Storage,
    jwtSecret string,
) *Container
```

**Después:**
```go
func NewContainer(resources *bootstrap.Resources) *Container
```

### 3. Cleanup de Recursos

**Antes:**
```go
defer pgDB.Close()
defer mongoDB.Client().Disconnect(ctx)
defer publisher.Close()
// ... más defer statements
```

**Después:**
```go
defer cleanup() // Una sola función limpia todo
```

## Migración Paso a Paso

### Paso 1: Actualizar main.go

**Código Original:**
```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/EduGoGroup/edugo-api-mobile/internal/config"
    "github.com/EduGoGroup/edugo-api-mobile/internal/container"
    "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/database"
    "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/router"
    "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
    "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/storage/s3"
    "github.com/EduGoGroup/edugo-api-mobile/pkg/logger"
    "go.uber.org/zap"
)

func main() {
    ctx := context.Background()
    
    // Cargar configuración
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Inicializar logger
    appLogger := logger.NewZapLogger(cfg.Logger.Level, cfg.Logger.Format)
    
    // Inicializar PostgreSQL
    pgDB, err := database.InitPostgreSQL(ctx, cfg, appLogger)
    if err != nil {
        appLogger.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
    }
    defer pgDB.Close()
    
    // Inicializar MongoDB
    mongoDB, err := database.InitMongoDB(ctx, cfg, appLogger)
    if err != nil {
        appLogger.Fatal("Failed to connect to MongoDB", zap.Error(err))
    }
    defer mongoDB.Client().Disconnect(ctx)
    
    // Inicializar RabbitMQ
    var publisher rabbitmq.Publisher
    publisher, err = rabbitmq.NewRabbitMQPublisher(
        cfg.Messaging.RabbitMQ.URL,
        cfg.Messaging.RabbitMQ.Exchange,
        appLogger,
    )
    if err != nil {
        appLogger.Warn("Failed to connect to RabbitMQ, using noop publisher", zap.Error(err))
        publisher = rabbitmq.NewNoopPublisher(appLogger)
    }
    defer publisher.Close()
    
    // Inicializar S3
    var s3Client s3.S3Storage
    s3Client, err = s3.NewS3Client(ctx, cfg.Storage.S3, appLogger)
    if err != nil {
        appLogger.Warn("Failed to initialize S3, using noop storage", zap.Error(err))
        s3Client = s3.NewNoopS3Storage(appLogger)
    }
    
    // Crear container
    c := container.NewContainer(
        appLogger,
        pgDB,
        mongoDB,
        publisher,
        s3Client,
        cfg.Auth.JWTSecret,
    )
    
    // Setup router
    r := router.SetupRouter(c)
    
    // Iniciar servidor
    appLogger.Info("Starting server", zap.Int("port", cfg.Server.Port))
    if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
        appLogger.Fatal("Failed to start server", zap.Error(err))
    }
}
```

**Código Migrado:**
```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/EduGoGroup/edugo-api-mobile/internal/bootstrap"
    "github.com/EduGoGroup/edugo-api-mobile/internal/config"
    "github.com/EduGoGroup/edugo-api-mobile/internal/container"
    "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/router"
    "go.uber.org/zap"
)

func main() {
    ctx := context.Background()
    
    // Cargar configuración
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Bootstrap de infraestructura
    b := bootstrap.New(cfg)
    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    if err != nil {
        log.Fatalf("Failed to initialize infrastructure: %v", err)
    }
    defer cleanup()
    
    // Crear container
    c := container.NewContainer(resources)
    
    // Setup router
    r := router.SetupRouter(c)
    
    // Iniciar servidor
    resources.Logger.Info("Starting server", zap.Int("port", cfg.Server.Port))
    if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
        resources.Logger.Fatal("Failed to start server", zap.Error(err))
    }
}
```

### Paso 2: Actualizar Container

**Código Original:**
```go
type Container struct {
    Logger           logger.Logger
    PostgreSQL       *sql.DB
    MongoDB          *mongo.Database
    RabbitMQPublisher rabbitmq.Publisher
    S3Client         s3.S3Storage
    JWTSecret        string
    
    // Services
    AuthService     *service.AuthService
    MaterialService *service.MaterialService
    // ...
}

func NewContainer(
    log logger.Logger,
    pgDB *sql.DB,
    mongoDB *mongo.Database,
    publisher rabbitmq.Publisher,
    s3Client s3.S3Storage,
    jwtSecret string,
) *Container {
    c := &Container{
        Logger:           log,
        PostgreSQL:       pgDB,
        MongoDB:          mongoDB,
        RabbitMQPublisher: publisher,
        S3Client:         s3Client,
        JWTSecret:        jwtSecret,
    }
    
    // Inicializar servicios
    c.initializeServices()
    
    return c
}
```

**Código Migrado:**
```go
type Container struct {
    Resources *bootstrap.Resources
    
    // Services
    AuthService     *service.AuthService
    MaterialService *service.MaterialService
    // ...
}

func NewContainer(resources *bootstrap.Resources) *Container {
    c := &Container{
        Resources: resources,
    }
    
    // Inicializar servicios
    c.initializeServices()
    
    return c
}

// Métodos helper para acceso a recursos
func (c *Container) Logger() logger.Logger {
    return c.Resources.Logger
}

func (c *Container) PostgreSQL() *sql.DB {
    return c.Resources.PostgreSQL
}

func (c *Container) MongoDB() *mongo.Database {
    return c.Resources.MongoDB
}

func (c *Container) Publisher() rabbitmq.Publisher {
    return c.Resources.RabbitMQPublisher
}

func (c *Container) S3Client() s3.S3Storage {
    return c.Resources.S3Client
}

func (c *Container) JWTSecret() string {
    return c.Resources.JWTSecret
}
```

### Paso 3: Actualizar Referencias en Servicios

Si tus servicios acceden directamente a campos del container, actualiza las referencias:

**Antes:**
```go
func (s *MaterialService) GetMaterial(ctx context.Context, id string) (*Material, error) {
    s.container.Logger.Info("Getting material", zap.String("id", id))
    // ...
}
```

**Después:**
```go
func (s *MaterialService) GetMaterial(ctx context.Context, id string) (*Material, error) {
    s.container.Logger().Info("Getting material", zap.String("id", id))
    // O si prefieres acceso directo:
    s.container.Resources.Logger.Info("Getting material", zap.String("id", id))
    // ...
}
```

## Ejemplos de Código

### Ejemplo 1: Aplicación Básica

```go
func main() {
    ctx := context.Background()
    cfg, _ := config.Load()
    
    // Bootstrap simple
    b := bootstrap.New(cfg)
    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer cleanup()
    
    // Usar recursos
    container := container.NewContainer(resources)
    router := router.SetupRouter(container)
    router.Run(":8080")
}
```

### Ejemplo 2: Con Recursos Opcionales

```go
func main() {
    ctx := context.Background()
    cfg, _ := config.Load()
    
    // Marcar RabbitMQ y S3 como opcionales
    b := bootstrap.New(cfg,
        bootstrap.WithOptionalResource("rabbitmq"),
        bootstrap.WithOptionalResource("s3"),
    )
    
    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer cleanup()
    
    // La app funciona sin RabbitMQ ni S3
    container := container.NewContainer(resources)
    router := router.SetupRouter(container)
    router.Run(":8080")
}
```

### Ejemplo 3: Testing con Mocks

```go
func TestMaterialHandler(t *testing.T) {
    ctx := context.Background()
    cfg := testConfig()
    
    // Crear mocks
    mockLogger := logger.NewNoopLogger()
    mockDB := setupMockDB(t)
    mockMongoDB := setupMockMongoDB(t)
    mockPublisher := &MockPublisher{}
    mockS3 := &MockS3Storage{}
    
    // Bootstrap con mocks
    b := bootstrap.New(cfg,
        bootstrap.WithLogger(mockLogger),
        bootstrap.WithPostgreSQL(mockDB),
        bootstrap.WithMongoDB(mockMongoDB),
        bootstrap.WithRabbitMQ(mockPublisher),
        bootstrap.WithS3Client(mockS3),
    )
    
    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    require.NoError(t, err)
    defer cleanup()
    
    // Crear container y ejecutar tests
    container := container.NewContainer(resources)
    handler := handler.NewMaterialHandler(container)
    
    // Test handler
    req := httptest.NewRequest("GET", "/materials/123", nil)
    w := httptest.NewRecorder()
    handler.GetMaterial(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## Testing con Bootstrap

### Tests Unitarios

Para tests unitarios, inyecta mocks de todos los recursos:

```go
func TestService(t *testing.T) {
    cfg := &config.Config{}
    
    b := bootstrap.New(cfg,
        bootstrap.WithLogger(mockLogger),
        bootstrap.WithPostgreSQL(mockDB),
        bootstrap.WithMongoDB(mockMongoDB),
        bootstrap.WithRabbitMQ(mockPublisher),
        bootstrap.WithS3Client(mockS3),
    )
    
    resources, cleanup, _ := b.InitializeInfrastructure(context.Background())
    defer cleanup()
    
    // Usar resources en tests
}
```

### Tests de Integración

Para tests de integración, puedes usar recursos reales o testcontainers:

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    ctx := context.Background()
    cfg := loadTestConfig()
    
    // Bootstrap con recursos reales
    b := bootstrap.New(cfg)
    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    require.NoError(t, err)
    defer cleanup()
    
    // Tests con recursos reales
}
```

### Helper para Tests

Crea un helper para simplificar setup de tests:

```go
// test/helpers/bootstrap.go
package helpers

func SetupTestBootstrap(t *testing.T) (*bootstrap.Resources, func()) {
    t.Helper()
    
    cfg := testConfig()
    b := bootstrap.New(cfg,
        bootstrap.WithLogger(logger.NewNoopLogger()),
        bootstrap.WithPostgreSQL(setupTestDB(t)),
        bootstrap.WithMongoDB(setupTestMongoDB(t)),
        bootstrap.WithRabbitMQ(noop.NewNoopPublisher(logger.NewNoopLogger())),
        bootstrap.WithS3Client(noop.NewNoopS3Storage(logger.NewNoopLogger())),
    )
    
    resources, cleanup, err := b.InitializeInfrastructure(context.Background())
    require.NoError(t, err)
    
    return resources, cleanup
}

// Uso en tests
func TestMyFeature(t *testing.T) {
    resources, cleanup := helpers.SetupTestBootstrap(t)
    defer cleanup()
    
    // Tests...
}
```

## Recursos Opcionales

### Configuración

Marca recursos como opcionales en `config/config-local.yaml`:

```yaml
infrastructure:
  optional_resources:
    - rabbitmq
    - s3
```

O mediante variables de entorno:

```bash
INFRASTRUCTURE_OPTIONAL_RESOURCES=rabbitmq,s3
```

### Comportamiento

Cuando un recurso opcional falla:

1. Se registra una advertencia en los logs
2. Se inyecta una implementación noop
3. La aplicación continúa funcionando
4. Las operaciones que usan ese recurso registran mensajes de debug

### Recursos Disponibles

| Recurso | ¿Puede ser opcional? | Implementación Noop |
|---------|---------------------|---------------------|
| Logger | ❌ No | N/A |
| PostgreSQL | ❌ No | N/A |
| MongoDB | ❌ No | N/A |
| RabbitMQ | ✅ Sí | `noop.NoopPublisher` |
| S3 | ✅ Sí | `noop.NoopS3Storage` |

### Ejemplo de Uso

```go
// Desarrollo local sin RabbitMQ
b := bootstrap.New(cfg,
    bootstrap.WithOptionalResource("rabbitmq"),
)

resources, cleanup, err := b.InitializeInfrastructure(ctx)
// err será nil incluso si RabbitMQ no está disponible
// resources.RabbitMQPublisher será un NoopPublisher
```

## Troubleshooting

### Error: "required resource failed to initialize"

**Causa**: Un recurso requerido (PostgreSQL o MongoDB) no está disponible.

**Solución**:
1. Verifica que el servicio esté corriendo
2. Verifica la configuración de conexión
3. Revisa los logs para detalles del error

### Error: "cleanup failed"

**Causa**: Error al cerrar uno o más recursos.

**Solución**:
- Revisa los logs para identificar qué recurso falló
- El cleanup continúa con otros recursos incluso si uno falla
- Generalmente no es crítico, pero investiga la causa

### Advertencia: "optional resource failed, using noop"

**Causa**: Un recurso opcional no está disponible.

**Solución**:
- Esto es esperado si no tienes el recurso configurado
- La aplicación funciona con funcionalidad reducida
- Si necesitas el recurso, verifica su configuración

### Tests fallan después de migración

**Causa**: Tests antiguos esperan la firma anterior del container.

**Solución**:
1. Actualiza los tests para usar `bootstrap.Resources`
2. Usa helpers de testing para simplificar setup
3. Inyecta mocks mediante opciones de bootstrap

### Performance degradado

**Causa**: Inicialización secuencial de recursos.

**Solución**:
- Actualmente los recursos se inicializan secuencialmente
- Futura optimización: inicialización paralela de recursos independientes
- Para desarrollo, usa recursos opcionales para reducir tiempo de inicio

## Recursos Adicionales

- **[internal/bootstrap/INTEGRATION_TESTS.md](../internal/bootstrap/INTEGRATION_TESTS.md)** - Guía de testing con bootstrap
- **[config/OPTIONAL_RESOURCES.md](../config/OPTIONAL_RESOURCES.md)** - Configuración de recursos opcionales
- **[.kiro/specs/infrastructure-bootstrap-refactor/](../.kiro/specs/infrastructure-bootstrap-refactor/)** - Especificación completa del diseño

## Preguntas Frecuentes

### ¿Necesito migrar todo mi código de una vez?

No. El sistema de bootstrap es compatible con código existente. Puedes migrar gradualmente:
1. Primero migra `main.go`
2. Luego actualiza el container
3. Finalmente actualiza tests

### ¿Puedo seguir usando las funciones de inicialización directamente?

Sí. Las funciones en `internal/infrastructure` se mantienen sin cambios. El bootstrap las usa internamente.

### ¿Cómo afecta esto a mis tests existentes?

Tests que no usan el container directamente no se ven afectados. Tests que crean el container necesitan actualizarse para usar `bootstrap.Resources`.

### ¿Puedo agregar nuevos recursos al bootstrap?

Sí. Sigue el patrón existente:
1. Define la interfaz en `interfaces.go`
2. Agrega factory en `factories.go`
3. Agrega campo en `Resources`
4. Agrega opción `With...()` en `config.go`
5. Actualiza `InitializeInfrastructure()` en `bootstrap.go`

### ¿El bootstrap afecta el performance?

No significativamente. El overhead es mínimo (validación de opciones, registro de cleanup). El tiempo de inicialización es el mismo que antes.
