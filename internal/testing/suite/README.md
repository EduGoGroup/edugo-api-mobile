# Suite de Tests de Integración con Testcontainers

Suite reutilizable para tests de integración que comparte contenedores Docker entre tests, integrada con scripts de `edugo-infrastructure`.

## 🎯 Características

- ✅ **Contenedores Compartidos**: PostgreSQL, MongoDB, RabbitMQ se levantan UNA vez
- ✅ **Migraciones Automáticas**: Ejecuta migraciones de `edugo-infrastructure` en SetupSuite
- ✅ **Seeds Automáticos**: Carga datos de prueba de `edugo-infrastructure`
- ✅ **Cleanup Automático**: Limpia datos entre tests sin reiniciar contenedores
- ✅ **Performance**: ~75% más rápido que levantar contenedores por test

## 📦 Instalación

```bash
# Dependencias ya están en go.mod
go mod tidy
```

## 🚀 Uso Básico

### 1. Crear Suite de Test

```go
package mypackage_test

import (
	"testing"

	"github.com/EduGoGroup/edugo-api-mobile/internal/testing/suite"
	testifySuite "github.com/stretchr/testify/suite"
)

type MySuite struct {
	suite.IntegrationTestSuite
}

func TestMySuite(t *testing.T) {
	testifySuite.Run(t, new(MySuite))
}

func (s *MySuite) TestWithPostgres() {
	// s.PostgresDB está listo con migraciones y seeds
	var count int
	err := s.PostgresDB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	s.NoError(err)
	s.Greater(count, 0, "Debe haber usuarios del seed")
}
```

### 2. Ejecutar Tests

```bash
# Ejecutar todos los tests de integración
RUN_INTEGRATION_TESTS=true go test -tags=integration -v ./...

# Ejecutar suite específica
RUN_INTEGRATION_TESTS=true go test -tags=integration -v ./internal/testing/suite -run TestMySuite
```

## 📋 Recursos Disponibles en la Suite

La suite `IntegrationTestSuite` expone:

```go
type IntegrationTestSuite struct {
    suite.Suite

    // Contenedores Docker
    PostgresContainer *postgres.PostgresContainer
    MongoContainer    *mongodb.MongoDBContainer
    RabbitContainer   *rabbitmq.RabbitMQContainer

    // Conexiones listas para usar
    PostgresDB *sql.DB          // PostgreSQL con migraciones y seeds
    MongoDB    *mongo.Database  // MongoDB listo
    Logger     logger.Logger    // Logger Zap configurado
}
```

## 🔄 Ciclo de Vida de los Tests

```
SetupSuite() - UNA VEZ antes de todos los tests
├── Inicializar logger
├── Calcular paths a infrastructure
├── Levantar contenedores (PostgreSQL, MongoDB, RabbitMQ)
├── Aplicar migraciones desde edugo-infrastructure
└── Aplicar seeds desde edugo-infrastructure

SetupTest() - ANTES de cada test individual
├── Limpiar datos de PostgreSQL (TRUNCATE)
└── Re-aplicar seeds

Test() - El test individual

TearDownSuite() - UNA VEZ después de todos los tests
└── Detener y eliminar contenedores
```

## 📊 Ejemplos de Tests

### Test con PostgreSQL

```go
func (s *MySuite) TestPostgresQuery() {
    ctx := context.Background()

    // Los datos de seeds ya están cargados
    var email string
    err := s.PostgresDB.QueryRowContext(ctx,
        "SELECT email FROM users WHERE role = $1 LIMIT 1",
        "student",
    ).Scan(&email)

    s.NoError(err)
    s.NotEmpty(email)
}
```

### Test con RabbitMQ

```go
func (s *MySuite) TestRabbitMQPublish() {
    ctx := context.Background()

    // Obtener URL de RabbitMQ
    rabbitURL, err := s.RabbitContainer.AmqpURL(ctx)
    s.NoError(err)

    // Conectar y usar
    conn, err := amqp.Dial(rabbitURL)
    s.NoError(err)
    defer conn.Close()

    ch, err := conn.Channel()
    s.NoError(err)

    // Publicar mensaje
    err = ch.Publish("exchange", "routing.key", false, false,
        amqp.Publishing{
            Body: []byte(`{"event": "test"}`),
        },
    )
    s.NoError(err)
}
```

### Test con Cleanup de Datos

```go
func (s *MySuite) TestDataIsolation_First() {
    // Insertar dato de prueba
    _, err := s.PostgresDB.Exec(
        "INSERT INTO users (id, email, ...) VALUES ($1, $2, ...)",
        uuid.New(), "test@example.com", ...,
    )
    s.NoError(err)
}

func (s *MySuite) TestDataIsolation_Second() {
    // El dato del test anterior NO existe
    // porque SetupTest() limpia entre tests
    var count int
    err := s.PostgresDB.QueryRow(
        "SELECT COUNT(*) FROM users WHERE email = $1",
        "test@example.com",
    ).Scan(&count)
    s.NoError(err)
    s.Equal(0, count, "Datos del test anterior deben estar limpios")
}
```

## ⚙️ Configuración

### Variables de Entorno

```bash
# Requerida para ejecutar tests de integración
export RUN_INTEGRATION_TESTS=true
```

### Paths de Infrastructure

La suite automáticamente calcula los paths relativos a `edugo-infrastructure`:

```
edugo-api-mobile/
└── internal/testing/suite/

edugo-infrastructure/
├── database/migrations/postgres/  ← Migraciones
└── seeds/postgres/                ← Datos de prueba
```

## 🧪 Tests de Ejemplo Incluidos

### `integration_suite_test.go`
- ✅ TestPostgresConnection - Verifica tablas creadas
- ✅ TestSeedsApplied - Verifica seeds cargados
- ✅ TestDataCleanupBetweenTests - Verifica limpieza
- ✅ TestContainersAreShared - Verifica compartición

### `rabbitmq_test.go`
- ✅ TestRabbitMQConnection - Verifica conexión
- ✅ TestPublishMessage - Publica mensaje a exchange
- ✅ TestConsumeMessage - Consume mensaje de queue
- ✅ TestMultiplePublishers - Múltiples publishers concurrentes

## 📈 Performance

| Método | Tiempo | Descripción |
|--------|--------|-------------|
| **Contenedores por test** | ~60s | Cada test levanta/detiene contenedores |
| **Contenedores compartidos** | ~15s | Suite levanta UNA vez |
| **Ganancia** | **75%** | Reducción de tiempo |

## 🔧 Troubleshooting

### Error: "Access denied - path outside allowed directories"

Los scripts de infrastructure deben estar en el directorio hermano:
```
/path/to/edugo-api-mobile/
/path/to/edugo-infrastructure/  ← Debe existir aquí
```

### Error: "FK constraint violation"

Los seeds se ejecutan en orden correcto automáticamente.
Si agregaste nuevos seeds, verifica las dependencias en:
`edugo-infrastructure/testing/postgres.go:orderSeedsByDependencies()`

### Tests lentos

- Verifica que Docker Desktop tenga recursos suficientes (7GB+ RAM)
- Los contenedores se comparten, no deberían reiniciarse entre tests

## 📚 Referencias

- **Infrastructure Testing Package**: `github.com/EduGoGroup/edugo-infrastructure/testing`
- **Testcontainers**: https://golang.testcontainers.org/
- **Testify Suite**: https://pkg.go.dev/github.com/stretchr/testify/suite

## 🤝 Contribuir

Para agregar nuevos tests:

1. Extender `IntegrationTestSuite`
2. Usar recursos compartidos (`PostgresDB`, `MongoDB`, `RabbitContainer`)
3. Confiar en `SetupTest()` para datos limpios
4. Ejecutar con `RUN_INTEGRATION_TESTS=true`

**No crear nuevos contenedores** - usa los compartidos para mejor performance.
