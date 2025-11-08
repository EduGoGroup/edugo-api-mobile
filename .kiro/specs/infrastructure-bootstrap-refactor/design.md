# Design Document

## Overview

Este diseño refactoriza el sistema de inicialización de infraestructura de la API Mobile, moviendo la responsabilidad de creación de recursos desde `cmd/main.go` hacia un nuevo paquete `internal/bootstrap`. El diseño implementa el patrón Factory con soporte para recursos opcionales, inyección de mocks, y gestión del ciclo de vida.

### Objetivos del Diseño

1. **Simplificar main.go**: Reducir de ~180 líneas a <50 líneas
2. **Consistencia**: Todos los recursos siguen el mismo patrón de inicialización
3. **Flexibilidad**: Soportar recursos opcionales y mocks para testing
4. **Mantenibilidad**: Código organizado en módulos cohesivos con responsabilidades claras

## Architecture

### Estructura de Paquetes

```
internal/
├── bootstrap/
│   ├── bootstrap.go          # Orquestador principal del bootstrap
│   ├── interfaces.go          # Interfaces de recursos de infraestructura
│   ├── config.go              # Configuración de recursos opcionales
│   ├── factories.go           # Factories para crear recursos
│   ├── lifecycle.go           # Gestión del ciclo de vida (cleanup)
│   └── noop/                  # Implementaciones noop para recursos opcionales
│       ├── logger.go
│       ├── publisher.go
│       └── storage.go
└── container/
    └── ...                    # Sin cambios significativos
```

### Flujo de Inicialización

```
main.go
   │
   ├─> config.Load()
   │
   ├─> bootstrap.New(config, options...)
   │      │
   │      ├─> CreateLogger()
   │      ├─> CreatePostgreSQL()
   │      ├─> CreateMongoDB()
   │      ├─> CreateRabbitMQ()
   │      └─> CreateS3Client()
   │
   ├─> bootstrap.InitializeInfrastructure()
   │      │
   │      └─> Returns: Resources + Cleanup function
   │
   ├─> container.NewContainer(resources...)
   │
   ├─> router.SetupRouter(container)
   │
   └─> server.Run()
        │
        └─> defer cleanup()
```

## Components and Interfaces

### 1. Bootstrap Orchestrator

**Archivo**: `internal/bootstrap/bootstrap.go`

```go
// Bootstrapper orquesta la inicialización de todos los recursos de infraestructura
type Bootstrapper struct {
    config    *config.Config
    options   *BootstrapOptions
    resources *Resources
    logger    logger.Logger
}

// Resources encapsula todos los recursos de infraestructura inicializados
type Resources struct {
    Logger           logger.Logger
    PostgreSQL       *sql.DB
    MongoDB          *mongo.Database
    RabbitMQPublisher rabbitmq.Publisher
    S3Client         S3Storage
    JWTSecret        string
}

// BootstrapOptions permite configurar el comportamiento del bootstrap
type BootstrapOptions struct {
    // Recursos pre-construidos (para testing)
    Logger           logger.Logger
    PostgreSQL       *sql.DB
    MongoDB          *mongo.Database
    RabbitMQPublisher rabbitmq.Publisher
    S3Client         S3Storage
    
    // Configuración de recursos opcionales
    OptionalResources map[string]bool
}

// New crea un nuevo Bootstrapper
func New(cfg *config.Config, opts ...BootstrapOption) *Bootstrapper

// InitializeInfrastructure inicializa todos los recursos y retorna cleanup function
func (b *Bootstrapper) InitializeInfrastructure(ctx context.Context) (*Resources, func() error, error)
```

**Responsabilidades**:
- Orquestar la inicialización de recursos en el orden correcto
- Aplicar opciones de configuración (mocks, recursos opcionales)
- Manejar errores y degradación graciosa
- Retornar función de cleanup para gestión del ciclo de vida

### 2. Resource Interfaces

**Archivo**: `internal/bootstrap/interfaces.go`

```go
// LoggerFactory crea instancias de logger
type LoggerFactory interface {
    Create(level, format string) (logger.Logger, error)
}

// DatabaseFactory crea conexiones a bases de datos
type DatabaseFactory interface {
    CreatePostgreSQL(ctx context.Context, cfg *config.Config, log logger.Logger) (*sql.DB, error)
    CreateMongoDB(ctx context.Context, cfg *config.Config, log logger.Logger) (*mongo.Database, error)
}

// MessagingFactory crea clientes de mensajería
type MessagingFactory interface {
    CreatePublisher(url, exchange string, log logger.Logger) (rabbitmq.Publisher, error)
}

// StorageFactory crea clientes de almacenamiento
type StorageFactory interface {
    CreateS3Client(ctx context.Context, cfg s3.S3Config, log logger.Logger) (S3Storage, error)
}

// S3Storage define las operaciones de almacenamiento (ya existe en s3/interface.go)
type S3Storage interface {
    GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expires time.Duration) (string, error)
    GeneratePresignedDownloadURL(ctx context.Context, key string, expires time.Duration) (string, error)
}
```

**Responsabilidades**:
- Definir contratos para factories de recursos
- Permitir inyección de implementaciones alternativas
- Facilitar testing con mocks

### 3. Resource Factories

**Archivo**: `internal/bootstrap/factories.go`

```go
// DefaultFactories implementa todas las factories con implementaciones reales
type DefaultFactories struct{}

func (f *DefaultFactories) CreateLogger(level, format string) (logger.Logger, error) {
    return logger.NewZapLogger(level, format), nil
}

func (f *DefaultFactories) CreatePostgreSQL(ctx context.Context, cfg *config.Config, log logger.Logger) (*sql.DB, error) {
    return database.InitPostgreSQL(ctx, cfg, log)
}

func (f *DefaultFactories) CreateMongoDB(ctx context.Context, cfg *config.Config, log logger.Logger) (*mongo.Database, error) {
    return database.InitMongoDB(ctx, cfg, log)
}

func (f *DefaultFactories) CreatePublisher(url, exchange string, log logger.Logger) (rabbitmq.Publisher, error) {
    return rabbitmq.NewRabbitMQPublisher(url, exchange, log)
}

func (f *DefaultFactories) CreateS3Client(ctx context.Context, cfg s3.S3Config, log logger.Logger) (S3Storage, error) {
    return s3.NewS3Client(ctx, cfg, log)
}
```

**Responsabilidades**:
- Implementar la creación de recursos reales
- Encapsular la lógica de inicialización específica de cada recurso
- Delegar a las implementaciones existentes en `internal/infrastructure`

### 4. Bootstrap Configuration

**Archivo**: `internal/bootstrap/config.go`

```go
// ResourceConfig define la configuración de un recurso
type ResourceConfig struct {
    Name     string
    Optional bool
    Enabled  bool
}

// DefaultResourceConfig retorna la configuración por defecto
func DefaultResourceConfig() map[string]ResourceConfig {
    return map[string]ResourceConfig{
        "logger":     {Name: "logger", Optional: false, Enabled: true},
        "postgresql": {Name: "postgresql", Optional: false, Enabled: true},
        "mongodb":    {Name: "mongodb", Optional: false, Enabled: true},
        "rabbitmq":   {Name: "rabbitmq", Optional: true, Enabled: true},
        "s3":         {Name: "s3", Optional: true, Enabled: true},
    }
}

// BootstrapOption es una función que modifica BootstrapOptions
type BootstrapOption func(*BootstrapOptions)

// WithLogger inyecta un logger pre-construido
func WithLogger(log logger.Logger) BootstrapOption

// WithPostgreSQL inyecta una conexión PostgreSQL pre-construida
func WithPostgreSQL(db *sql.DB) BootstrapOption

// WithMongoDB inyecta una conexión MongoDB pre-construida
func WithMongoDB(db *mongo.Database) BootstrapOption

// WithRabbitMQ inyecta un publisher RabbitMQ pre-construido
func WithRabbitMQ(pub rabbitmq.Publisher) BootstrapOption

// WithS3Client inyecta un cliente S3 pre-construido
func WithS3Client(client S3Storage) BootstrapOption

// WithOptionalResource marca un recurso como opcional
func WithOptionalResource(resourceName string) BootstrapOption

// WithDisabledResource deshabilita un recurso
func WithDisabledResource(resourceName string) BootstrapOption
```

**Responsabilidades**:
- Definir qué recursos son opcionales vs requeridos
- Proporcionar API funcional para configurar el bootstrap
- Permitir inyección de recursos pre-construidos

### 5. Lifecycle Management

**Archivo**: `internal/bootstrap/lifecycle.go`

```go
// CleanupFunc es una función que cierra un recurso
type CleanupFunc func() error

// LifecycleManager gestiona el ciclo de vida de recursos
type LifecycleManager struct {
    cleanupFuncs []CleanupFunc
    logger       logger.Logger
}

// Register registra una función de cleanup
func (lm *LifecycleManager) Register(name string, cleanup CleanupFunc)

// Cleanup ejecuta todas las funciones de cleanup en orden inverso
func (lm *LifecycleManager) Cleanup() error {
    var errors []error
    
    // Ejecutar en orden inverso (LIFO)
    for i := len(lm.cleanupFuncs) - 1; i >= 0; i-- {
        if err := lm.cleanupFuncs[i](); err != nil {
            errors = append(errors, err)
            // Continuar con los demás recursos
        }
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("errors during cleanup: %v", errors)
    }
    return nil
}
```

**Responsabilidades**:
- Registrar funciones de cleanup para cada recurso
- Ejecutar cleanup en orden inverso (LIFO)
- Manejar errores sin detener el cleanup de otros recursos

### 6. Noop Implementations

**Archivo**: `internal/bootstrap/noop/publisher.go`

```go
// NoopPublisher es una implementación noop de rabbitmq.Publisher
type NoopPublisher struct {
    logger logger.Logger
}

func NewNoopPublisher(log logger.Logger) *NoopPublisher {
    return &NoopPublisher{logger: log}
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

**Archivo**: `internal/bootstrap/noop/storage.go`

```go
// NoopS3Storage es una implementación noop de S3Storage
type NoopS3Storage struct {
    logger logger.Logger
}

func NewNoopS3Storage(log logger.Logger) *NoopS3Storage {
    return &NoopS3Storage{logger: log}
}

func (s *NoopS3Storage) GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expires time.Duration) (string, error) {
    s.logger.Debug("noop storage: presigned URL not generated (S3 not available)",
        zap.String("key", key),
    )
    return "", fmt.Errorf("S3 not available")
}

func (s *NoopS3Storage) GeneratePresignedDownloadURL(ctx context.Context, key string, expires time.Duration) (string, error) {
    s.logger.Debug("noop storage: presigned URL not generated (S3 not available)",
        zap.String("key", key),
    )
    return "", fmt.Errorf("S3 not available")
}
```

**Responsabilidades**:
- Proporcionar implementaciones seguras cuando recursos opcionales no están disponibles
- Registrar intentos de uso para debugging
- Evitar panics o errores fatales

## Data Models

### Bootstrap Result

```go
type BootstrapResult struct {
    Resources *Resources
    Cleanup   func() error
    Warnings  []string  // Advertencias sobre recursos opcionales no disponibles
}
```

### Resource Status

```go
type ResourceStatus struct {
    Name      string
    Available bool
    Error     error
    Duration  time.Duration
}
```

## Error Handling

### Estrategia de Errores

1. **Recursos Requeridos**: Error fatal, detener inicialización
2. **Recursos Opcionales**: Log warning, continuar con implementación noop
3. **Errores de Cleanup**: Log error, continuar con cleanup de otros recursos

### Tipos de Errores

```go
var (
    ErrRequiredResourceFailed = errors.New("required resource failed to initialize")
    ErrInvalidConfiguration   = errors.New("invalid bootstrap configuration")
)

type ResourceError struct {
    ResourceName string
    Operation    string
    Err          error
}

func (e *ResourceError) Error() string {
    return fmt.Sprintf("resource '%s' failed during %s: %v", e.ResourceName, e.Operation, e.Err)
}
```

## Testing Strategy

### Unit Tests

1. **bootstrap_test.go**: Probar orquestación con mocks
2. **factories_test.go**: Probar cada factory individualmente
3. **lifecycle_test.go**: Probar cleanup en diferentes escenarios
4. **config_test.go**: Probar opciones de configuración

### Integration Tests

1. **Recursos opcionales**: Verificar degradación graciosa
2. **Inyección de mocks**: Verificar que mocks se usan correctamente
3. **Cleanup**: Verificar que recursos se cierran en orden correcto

### Test Helpers

```go
// TestBootstrapOptions retorna opciones para testing
func TestBootstrapOptions() *BootstrapOptions {
    return &BootstrapOptions{
        Logger:           logger.NewNoopLogger(),
        PostgreSQL:       mockDB,
        MongoDB:          mockMongoDB,
        RabbitMQPublisher: noop.NewNoopPublisher(logger.NewNoopLogger()),
        S3Client:         noop.NewNoopS3Storage(logger.NewNoopLogger()),
    }
}
```

## Migration Strategy

### Fase 1: Crear Infraestructura Bootstrap

1. Crear paquete `internal/bootstrap`
2. Implementar interfaces y factories
3. Implementar implementaciones noop
4. Escribir tests unitarios

### Fase 2: Refactorizar Main

1. Actualizar `cmd/main.go` para usar bootstrap
2. Simplificar lógica de inicialización
3. Mantener comportamiento existente

### Fase 3: Actualizar Container

1. Adaptar `internal/container` para recibir `Resources`
2. Actualizar tests de integración
3. Verificar que todo funciona correctamente

### Compatibilidad

- **Sin breaking changes**: Las funciones existentes en `internal/infrastructure` se mantienen
- **Backward compatible**: El container sigue recibiendo los mismos tipos
- **Incremental**: Se puede migrar gradualmente

## Performance Considerations

1. **Inicialización Paralela**: Recursos independientes pueden inicializarse en paralelo (futura optimización)
2. **Lazy Loading**: Recursos opcionales pueden cargarse bajo demanda (futura optimización)
3. **Connection Pooling**: Mantener configuración existente de pools

## Security Considerations

1. **Secrets**: JWT secret y credenciales se manejan igual que antes
2. **Logging**: No registrar información sensible durante bootstrap
3. **Error Messages**: No exponer detalles de infraestructura en errores públicos

## Example Usage

### Uso Normal (Producción)

```go
func main() {
    ctx := context.Background()
    cfg, _ := config.Load()
    
    // Bootstrap con configuración por defecto
    b := bootstrap.New(cfg)
    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer cleanup()
    
    // Crear container y arrancar servidor
    c := container.NewContainer(resources)
    r := router.SetupRouter(c)
    r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
```

### Uso con Recursos Opcionales (Desarrollo)

```go
func main() {
    ctx := context.Background()
    cfg, _ := config.Load()
    
    // Bootstrap con RabbitMQ y S3 opcionales
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
    c := container.NewContainer(resources)
    r := router.SetupRouter(c)
    r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
```

### Uso con Mocks (Testing)

```go
func TestAPI(t *testing.T) {
    ctx := context.Background()
    cfg := testConfig()
    
    // Bootstrap con todos los recursos mockeados
    b := bootstrap.New(cfg,
        bootstrap.WithLogger(mockLogger),
        bootstrap.WithPostgreSQL(mockDB),
        bootstrap.WithMongoDB(mockMongoDB),
        bootstrap.WithRabbitMQ(mockPublisher),
        bootstrap.WithS3Client(mockS3),
    )
    resources, cleanup, _ := b.InitializeInfrastructure(ctx)
    defer cleanup()
    
    // Tests con mocks
    c := container.NewContainer(resources)
    // ... tests
}
```
