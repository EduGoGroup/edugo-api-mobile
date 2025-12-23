# Guía de Uso del Sistema de Bootstrap

Esta guía explica cómo usar el sistema de bootstrap de infraestructura para inicializar y gestionar recursos en la aplicación EduGo API Mobile.

## Tabla de Contenidos

- [Introducción](#introducción)
- [Conceptos Básicos](#conceptos-básicos)
- [Uso Básico](#uso-básico)
- [Recursos Opcionales](#recursos-opcionales)
- [Inyección de Mocks](#inyección-de-mocks)
- [Gestión del Ciclo de Vida](#gestión-del-ciclo-de-vida)
- [Configuración Avanzada](#configuración-avanzada)
- [Ejemplos Prácticos](#ejemplos-prácticos)
- [Best Practices](#best-practices)

## Introducción

El sistema de bootstrap proporciona una forma consistente y modular de inicializar todos los recursos de infraestructura de la aplicación:

- **Logger**: Sistema de logging estructurado (Zap)
- **PostgreSQL**: Base de datos principal
- **MongoDB**: Base de datos para evaluaciones y resúmenes
- **RabbitMQ**: Sistema de mensajería para eventos
- **S3**: Almacenamiento de archivos

### Beneficios

✅ **Simplicidad**: Una sola llamada inicializa todos los recursos  
✅ **Consistencia**: Todos los recursos siguen el mismo patrón  
✅ **Testabilidad**: Inyección fácil de mocks  
✅ **Flexibilidad**: Recursos opcionales para desarrollo  
✅ **Seguridad**: Cleanup automático de recursos  

## Conceptos Básicos

### Resources

La estructura `Resources` encapsula todos los recursos de infraestructura:

```go
type Resources struct {
    Logger           logger.Logger
    PostgreSQL       *sql.DB
    MongoDB          *mongo.Database
    RabbitMQPublisher rabbitmq.Publisher
    S3Client         s3.S3Storage
    JWTSecret        string
}
```

### Bootstrapper

El `Bootstrapper` orquesta la inicialización de recursos:

```go
type Bootstrapper struct {
    config    *config.Config
    options   *BootstrapOptions
    resources *Resources
    logger    logger.Logger
}
```

### BootstrapOptions

Opciones para configurar el comportamiento del bootstrap:

```go
type BootstrapOptions struct {
    // Recursos pre-construidos (para testing)
    Logger           logger.Logger
    PostgreSQL       *sql.DB
    MongoDB          *mongo.Database
    RabbitMQPublisher rabbitmq.Publisher
    S3Client         s3.S3Storage

    // Configuración de recursos opcionales
    OptionalResources map[string]bool
}
```

## Uso Básico

### Inicialización Simple

```go
package main

import (
    "context"
    "log"

    "github.com/EduGoGroup/edugo-api-mobile/internal/bootstrap"
    "github.com/EduGoGroup/edugo-api-mobile/internal/config"
)

func main() {
    ctx := context.Background()

    // 1. Cargar configuración
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 2. Crear bootstrapper
    b := bootstrap.New(cfg)

    // 3. Inicializar infraestructura
    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    if err != nil {
        log.Fatalf("Failed to initialize infrastructure: %v", err)
    }
    defer cleanup()

    // 4. Usar recursos
    resources.Logger.Info("Application started")

    // Tu código aquí...
}
```

### Acceso a Recursos

```go
// Logger
resources.Logger.Info("Message", zap.String("key", "value"))

// PostgreSQL
rows, err := resources.PostgreSQL.QueryContext(ctx, "SELECT * FROM users")

// MongoDB
collection := resources.MongoDB.Collection("assessments")

// RabbitMQ
err := resources.RabbitMQPublisher.Publish(ctx, "events", "user.created", data)

// S3
url, err := resources.S3Client.GeneratePresignedUploadURL(ctx, key, contentType, expires)

// JWT Secret
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, err := token.SignedString([]byte(resources.JWTSecret))
```

## Recursos Opcionales

### ¿Qué son los Recursos Opcionales?

Recursos que pueden no estar disponibles sin impedir el funcionamiento de la aplicación. Cuando un recurso opcional falla:

1. Se registra una advertencia
2. Se inyecta una implementación noop
3. La aplicación continúa funcionando

### Configuración

**Opción 1: Archivo YAML**

Edita `config/config-local.yaml`:

```yaml
infrastructure:
  optional_resources:
    - rabbitmq
    - s3
```

**Opción 2: Variables de Entorno**

```bash
# En .env
INFRASTRUCTURE_OPTIONAL_RESOURCES=rabbitmq,s3
```

**Opción 3: Código**

```go
b := bootstrap.New(cfg,
    bootstrap.WithOptionalResource("rabbitmq"),
    bootstrap.WithOptionalResource("s3"),
)
```

### Recursos Disponibles

| Recurso | ¿Opcional? | Implementación Noop |
|---------|-----------|---------------------|
| Logger | ❌ No | N/A |
| PostgreSQL | ❌ No | N/A |
| MongoDB | ❌ No | N/A |
| RabbitMQ | ✅ Sí | `noop.NoopPublisher` |
| S3 | ✅ Sí | `noop.NoopS3Storage` |

### Ejemplo: Desarrollo sin RabbitMQ

```go
func main() {
    ctx := context.Background()
    cfg, _ := config.Load()

    // Marcar RabbitMQ como opcional
    b := bootstrap.New(cfg,
        bootstrap.WithOptionalResource("rabbitmq"),
    )

    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer cleanup()

    // resources.RabbitMQPublisher será un NoopPublisher
    // Las llamadas a Publish() no harán nada pero no fallarán
    resources.RabbitMQPublisher.Publish(ctx, "events", "test", []byte("data"))
    // ^ Esto registra un mensaje de debug pero no falla
}
```

### Implementaciones Noop

#### NoopPublisher

```go
type NoopPublisher struct {
    logger logger.Logger
}

func (p *NoopPublisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
    p.logger.Debug("noop publisher: message not published (RabbitMQ not available)",
        zap.String("exchange", exchange),
        zap.String("routing_key", routingKey),
    )
    return nil
}

func (p *NoopPublisher) Close() error {
    return nil
}
```

#### NoopS3Storage

```go
type NoopS3Storage struct {
    logger logger.Logger
}

func (s *NoopS3Storage) GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expires time.Duration) (string, error) {
    s.logger.Debug("noop storage: presigned URL not generated (S3 not available)",
        zap.String("key", key),
    )
    return "", fmt.Errorf("S3 not available")
}
```

## Inyección de Mocks

### ¿Por Qué Inyectar Mocks?

- Testing sin dependencias externas
- Tests más rápidos y confiables
- Control total sobre el comportamiento de recursos
- Aislamiento de lógica de negocio

### Inyección de Logger

```go
mockLogger := logger.NewNoopLogger()

b := bootstrap.New(cfg,
    bootstrap.WithLogger(mockLogger),
)
```

### Inyección de PostgreSQL

```go
mockDB := setupMockDB(t)

b := bootstrap.New(cfg,
    bootstrap.WithPostgreSQL(mockDB),
)
```

### Inyección de MongoDB

```go
mockMongoDB := setupMockMongoDB(t)

b := bootstrap.New(cfg,
    bootstrap.WithMongoDB(mockMongoDB),
)
```

### Inyección de RabbitMQ

```go
type MockPublisher struct {
    PublishedMessages []Message
}

func (m *MockPublisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
    m.PublishedMessages = append(m.PublishedMessages, Message{
        Exchange:   exchange,
        RoutingKey: routingKey,
        Body:       body,
    })
    return nil
}

func (m *MockPublisher) Close() error {
    return nil
}

// Uso
mockPublisher := &MockPublisher{}
b := bootstrap.New(cfg,
    bootstrap.WithRabbitMQ(mockPublisher),
)
```

### Inyección de S3

```go
type MockS3Storage struct {
    GeneratedURLs map[string]string
}

func (m *MockS3Storage) GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expires time.Duration) (string, error) {
    url := fmt.Sprintf("https://mock-s3.example.com/%s", key)
    m.GeneratedURLs[key] = url
    return url, nil
}

// Uso
mockS3 := &MockS3Storage{GeneratedURLs: make(map[string]string)}
b := bootstrap.New(cfg,
    bootstrap.WithS3Client(mockS3),
)
```

### Ejemplo Completo de Testing

```go
func TestMaterialService(t *testing.T) {
    ctx := context.Background()
    cfg := testConfig()

    // Setup mocks
    mockLogger := logger.NewNoopLogger()
    mockDB := setupMockDB(t)
    mockMongoDB := setupMockMongoDB(t)
    mockPublisher := &MockPublisher{}
    mockS3 := &MockS3Storage{GeneratedURLs: make(map[string]string)}

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

    // Crear servicio
    service := service.NewMaterialService(resources)

    // Test
    material, err := service.CreateMaterial(ctx, &dto.CreateMaterialRequest{
        Title: "Test Material",
    })

    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, material)
    assert.Len(t, mockPublisher.PublishedMessages, 1)
    assert.Equal(t, "material.created", mockPublisher.PublishedMessages[0].RoutingKey)
}
```

## Gestión del Ciclo de Vida

### Cleanup Automático

El bootstrap retorna una función `cleanup()` que cierra todos los recursos:

```go
resources, cleanup, err := b.InitializeInfrastructure(ctx)
if err != nil {
    log.Fatal(err)
}
defer cleanup()
```

### Orden de Cleanup

Los recursos se cierran en orden inverso a su inicialización (LIFO):

1. S3 Client
2. RabbitMQ Publisher
3. MongoDB
4. PostgreSQL
5. Logger (si es necesario)

### Manejo de Errores en Cleanup

Si un recurso falla al cerrarse:

1. Se registra el error
2. Se continúa con el cleanup de otros recursos
3. Se retorna un error agregado al final

```go
// Ejemplo de error en cleanup
defer func() {
    if err := cleanup(); err != nil {
        log.Printf("Cleanup errors: %v", err)
        // La aplicación puede continuar o terminar según el contexto
    }
}()
```

### Cleanup Manual

Si necesitas control manual sobre el cleanup:

```go
resources, cleanup, err := b.InitializeInfrastructure(ctx)
if err != nil {
    log.Fatal(err)
}

// No usar defer, cleanup manual
// ... tu código ...

// Cleanup explícito
if err := cleanup(); err != nil {
    log.Printf("Cleanup failed: %v", err)
}
```

## Configuración Avanzada

### Deshabilitar Recursos

```go
b := bootstrap.New(cfg,
    bootstrap.WithDisabledResource("rabbitmq"),
)
```

### Combinar Opciones

```go
b := bootstrap.New(cfg,
    // Inyectar logger custom
    bootstrap.WithLogger(customLogger),

    // Marcar recursos como opcionales
    bootstrap.WithOptionalResource("rabbitmq"),
    bootstrap.WithOptionalResource("s3"),

    // Inyectar mock de PostgreSQL
    bootstrap.WithPostgreSQL(mockDB),
)
```

### Configuración Condicional

```go
func createBootstrapper(cfg *config.Config, env string) *bootstrap.Bootstrapper {
    var opts []bootstrap.BootstrapOption

    if env == "development" {
        // En desarrollo, RabbitMQ y S3 son opcionales
        opts = append(opts,
            bootstrap.WithOptionalResource("rabbitmq"),
            bootstrap.WithOptionalResource("s3"),
        )
    }

    if env == "test" {
        // En tests, usar mocks
        opts = append(opts,
            bootstrap.WithLogger(logger.NewNoopLogger()),
            bootstrap.WithPostgreSQL(testDB),
            bootstrap.WithMongoDB(testMongoDB),
        )
    }

    return bootstrap.New(cfg, opts...)
}
```

## Ejemplos Prácticos

### Ejemplo 1: Aplicación de Producción

```go
func main() {
    ctx := context.Background()
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Bootstrap con todos los recursos requeridos
    b := bootstrap.New(cfg)
    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    if err != nil {
        log.Fatalf("Failed to initialize infrastructure: %v", err)
    }
    defer cleanup()

    // Crear container
    container := container.NewContainer(resources)

    // Setup router
    router := router.SetupRouter(container)

    // Iniciar servidor
    resources.Logger.Info("Starting server", zap.Int("port", cfg.Server.Port))
    if err := router.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
        resources.Logger.Fatal("Failed to start server", zap.Error(err))
    }
}
```

### Ejemplo 2: Desarrollo Local

```go
func main() {
    ctx := context.Background()
    cfg, _ := config.Load()

    // Bootstrap con recursos opcionales para desarrollo
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

### Ejemplo 3: Tests Unitarios

```go
func TestAuthService(t *testing.T) {
    ctx := context.Background()
    cfg := &config.Config{
        Auth: config.AuthConfig{
            JWTSecret: "test-secret",
        },
    }

    // Mocks
    mockLogger := logger.NewNoopLogger()
    mockDB := setupMockDB(t)
    mockMongoDB := setupMockMongoDB(t)

    // Bootstrap con mocks
    b := bootstrap.New(cfg,
        bootstrap.WithLogger(mockLogger),
        bootstrap.WithPostgreSQL(mockDB),
        bootstrap.WithMongoDB(mockMongoDB),
        bootstrap.WithRabbitMQ(noop.NewNoopPublisher(mockLogger)),
        bootstrap.WithS3Client(noop.NewNoopS3Storage(mockLogger)),
    )

    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    require.NoError(t, err)
    defer cleanup()

    // Test service
    authService := service.NewAuthService(resources)
    token, err := authService.Login(ctx, "user@example.com", "password")

    assert.NoError(t, err)
    assert.NotEmpty(t, token)
}
```

### Ejemplo 4: Tests de Integración

```go
func TestIntegrationMaterialFlow(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    ctx := context.Background()
    cfg := loadIntegrationTestConfig()

    // Bootstrap con recursos reales (testcontainers)
    b := bootstrap.New(cfg)
    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    require.NoError(t, err)
    defer cleanup()

    // Test flujo completo
    container := container.NewContainer(resources)

    // 1. Crear material
    material, err := container.MaterialService.CreateMaterial(ctx, &dto.CreateMaterialRequest{
        Title: "Integration Test Material",
    })
    require.NoError(t, err)

    // 2. Verificar en base de datos
    var count int
    err = resources.PostgreSQL.QueryRowContext(ctx,
        "SELECT COUNT(*) FROM materials WHERE id = $1",
        material.ID,
    ).Scan(&count)
    require.NoError(t, err)
    assert.Equal(t, 1, count)
}
```

## Best Practices

### 1. Siempre Usar defer cleanup()

```go
// ✅ Correcto
resources, cleanup, err := b.InitializeInfrastructure(ctx)
if err != nil {
    log.Fatal(err)
}
defer cleanup()

// ❌ Incorrecto
resources, cleanup, err := b.InitializeInfrastructure(ctx)
// Olvidar defer cleanup()
```

### 2. Manejar Errores de Inicialización

```go
// ✅ Correcto
resources, cleanup, err := b.InitializeInfrastructure(ctx)
if err != nil {
    log.Fatalf("Failed to initialize: %v", err)
}
defer cleanup()

// ❌ Incorrecto
resources, cleanup, _ := b.InitializeInfrastructure(ctx)
// Ignorar errores
```

### 3. Usar Recursos Opcionales en Desarrollo

```go
// ✅ Correcto para desarrollo
if cfg.Environment == "development" {
    b = bootstrap.New(cfg,
        bootstrap.WithOptionalResource("rabbitmq"),
        bootstrap.WithOptionalResource("s3"),
    )
}

// ❌ Incorrecto para producción
// No marcar recursos críticos como opcionales en producción
```

### 4. Inyectar Todos los Mocks en Tests

```go
// ✅ Correcto
b := bootstrap.New(cfg,
    bootstrap.WithLogger(mockLogger),
    bootstrap.WithPostgreSQL(mockDB),
    bootstrap.WithMongoDB(mockMongoDB),
    bootstrap.WithRabbitMQ(mockPublisher),
    bootstrap.WithS3Client(mockS3),
)

// ❌ Incorrecto
b := bootstrap.New(cfg,
    bootstrap.WithLogger(mockLogger),
    // Olvidar otros mocks - puede intentar conectar a recursos reales
)
```

### 5. Crear Helpers para Tests

```go
// ✅ Correcto
func setupTestResources(t *testing.T) (*bootstrap.Resources, func()) {
    t.Helper()
    cfg := testConfig()
    b := bootstrap.New(cfg, testBootstrapOptions()...)
    resources, cleanup, err := b.InitializeInfrastructure(context.Background())
    require.NoError(t, err)
    return resources, cleanup
}

// Uso
func TestSomething(t *testing.T) {
    resources, cleanup := setupTestResources(t)
    defer cleanup()
    // ...
}
```

### 6. Logging Apropiado

```go
// ✅ Correcto
resources.Logger.Info("Operation completed",
    zap.String("operation", "create_material"),
    zap.String("material_id", id),
)

// ❌ Incorrecto
fmt.Printf("Operation completed: %s\n", id)
// No usar fmt.Printf, usar el logger estructurado
```

### 7. Verificar Recursos Opcionales

```go
// ✅ Correcto
if _, ok := resources.RabbitMQPublisher.(*noop.NoopPublisher); ok {
    resources.Logger.Warn("RabbitMQ not available, events will not be published")
}

// Alternativa: verificar en la lógica de negocio
err := resources.RabbitMQPublisher.Publish(ctx, "events", "key", data)
if err != nil {
    resources.Logger.Warn("Failed to publish event", zap.Error(err))
    // Continuar sin fallar
}
```

## Recursos Adicionales

- **[BOOTSTRAP_MIGRATION_GUIDE.md](BOOTSTRAP_MIGRATION_GUIDE.md)** - Guía de migración desde código legacy
- **[../internal/bootstrap/INTEGRATION_TESTS.md](../internal/bootstrap/INTEGRATION_TESTS.md)** - Guía de testing con bootstrap
- **[../config/OPTIONAL_RESOURCES.md](../config/OPTIONAL_RESOURCES.md)** - Configuración de recursos opcionales
- **[../.kiro/specs/infrastructure-bootstrap-refactor/](../.kiro/specs/infrastructure-bootstrap-refactor/)** - Especificación completa del diseño
