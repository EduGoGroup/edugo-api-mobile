# Documento de Diseño - Mejora de Estrategia de Testing

## Resumen Ejecutivo

Este documento presenta el diseño detallado para mejorar la estrategia de testing del proyecto edugo-api-mobile. El diseño se divide en tres fases principales: Análisis y Evaluación, Refactorización y Mejoras, e Implementación de Nuevos Tests. El objetivo es establecer una estrategia de testing robusta, escalable y mantenible que siga las mejores prácticas de la industria.

## Arquitectura General

### Estructura de Testing Propuesta

```
edugo-api-mobile/
├── internal/
│   ├── application/
│   │   ├── service/
│   │   │   ├── assessment_service.go
│   │   │   ├── assessment_service_test.go      # Tests unitarios junto al código
│   │   │   ├── material_service.go
│   │   │   └── material_service_test.go
│   │   └── dto/                                 # Sin tests (DTOs simples)
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── material.go
│   │   │   └── material_test.go                # Tests de lógica de dominio
│   │   └── valueobject/
│   │       ├── email.go
│   │       └── email_test.go
│   └── infrastructure/
│       ├── http/
│       │   └── handler/
│       │       ├── auth_handler.go
│       │       └── auth_handler_test.go        # Tests unitarios con mocks
│       └── persistence/
│           └── postgres/
│               └── repository/
│                   ├── user_repository_impl.go
│                   └── user_repository_impl_test.go
├── test/
│   ├── integration/                            # Tests de integración E2E
│   │   ├── config.go                           # Control de ejecución
│   │   ├── setup.go                            # Testcontainers setup
│   │   ├── testhelpers.go                      # Helpers y seeds
│   │   ├── auth_flow_test.go
│   │   ├── material_flow_test.go
│   │   └── assessment_flow_test.go
│   ├── fixtures/                               # Datos de prueba reutilizables
│   │   ├── users.json
│   │   ├── materials.json
│   │   └── assessments.json
│   └── scripts/                                # Scripts de setup para desarrollo
│       ├── setup_dev_env.sh
│       ├── seed_test_data.sh
│       └── teardown_dev_env.sh
├── .coverignore                                # Exclusiones de cobertura
├── .golangci.yml                               # Configuración de linter
└── Makefile                                    # Comandos de testing
```


### Principios de Diseño

1. **Tests Unitarios Junto al Código**: Los tests unitarios se ubican en el mismo paquete que el código que prueban, facilitando el mantenimiento y descubrimiento
2. **Tests de Integración Separados**: Los tests de integración se mantienen en `test/integration/` con build tags para ejecución controlada
3. **Reutilización de Infraestructura**: La infraestructura de testcontainers se puede reutilizar para desarrollo local
4. **Exclusiones Inteligentes**: Configuración clara de qué código excluir de cobertura (generado, DTOs, mocks)
5. **Helpers Centralizados**: Funciones helper reutilizables para setup, cleanup y seed de datos

## Componentes y Interfaces

### 1. Sistema de Análisis de Tests

#### Componente: TestAnalyzer

**Responsabilidad**: Analizar la estructura actual de tests y generar reportes

**Interfaz**:
```go
type TestAnalyzer interface {
    // AnalyzeStructure analiza la estructura de archivos de test
    AnalyzeStructure() (*TestStructureReport, error)
    
    // CalculateCoverage calcula cobertura por paquete
    CalculateCoverage() (*CoverageReport, error)
    
    // FindMissingTests identifica módulos sin tests
    FindMissingTests() ([]string, error)
    
    // ValidateIntegrationTests verifica que tests de integración funcionen
    ValidateIntegrationTests() (*ValidationReport, error)
}
```

**Modelo de Datos**:
```go
type TestStructureReport struct {
    TotalTestFiles      int
    UnitTestFiles       []string
    IntegrationTestFiles []string
    EmptyDirectories    []string
    TestsByPackage      map[string]int
}

type CoverageReport struct {
    OverallCoverage     float64
    PackageCoverage     map[string]float64
    UncoveredPackages   []string
    CriticalPackages    map[string]float64  // Servicios, dominio
}
```

### 2. Sistema de Configuración de Cobertura

#### Componente: CoverageConfig

**Responsabilidad**: Gestionar exclusiones y configuración de cobertura

**Archivo**: `.coverignore`
```
# Archivos generados
docs/docs.go
docs/swagger.json
docs/swagger.yaml

# DTOs y estructuras simples
internal/application/dto/
internal/domain/entity/
internal/infrastructure/http/request/
internal/infrastructure/http/response/

# Main y comandos
cmd/

# Mocks y helpers de testing
*_mock.go
*/mocks_test.go
*/testing_helpers.go
test/integration/testhelpers.go

# Configuración
internal/config/
tools/configctl/

# Noop implementations
internal/bootstrap/noop/
```

**Script de Cobertura**: `scripts/coverage.sh`
```bash
#!/bin/bash
# Genera reporte de cobertura excluyendo archivos configurados

EXCLUDE_PATTERN=$(cat .coverignore | grep -v '^#' | grep -v '^$' | tr '\n' '|' | sed 's/|$//')

go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep -vE "$EXCLUDE_PATTERN" > coverage-filtered.txt
go tool cover -html=coverage.out -o coverage.html
```


### 3. Infraestructura de Tests de Integración

#### Componente: TestInfrastructure

**Responsabilidad**: Gestionar testcontainers y recursos para tests de integración

**Estructura Actual (Mantener y Mejorar)**:
```go
// test/integration/setup.go
type TestContainers struct {
    Postgres *postgres.PostgresContainer
    MongoDB  *mongodb.MongoDBContainer
    RabbitMQ *rabbitmq.RabbitMQContainer
}

// Mejorar con configuración de RabbitMQ
func SetupContainers(t *testing.T) (*TestContainers, func()) {
    // ... código existente ...
    
    // NUEVO: Configurar RabbitMQ automáticamente
    if err := setupRabbitMQTopology(rabbitContainer); err != nil {
        t.Logf("Warning: RabbitMQ topology setup failed: %v", err)
    }
    
    return containers, cleanup
}

// NUEVO: Configurar exchanges y colas
func setupRabbitMQTopology(container *rabbitmq.RabbitMQContainer) error {
    conn, err := amqp.Dial(container.AmqpURL())
    if err != nil {
        return err
    }
    defer conn.Close()
    
    ch, err := conn.Channel()
    if err != nil {
        return err
    }
    defer ch.Close()
    
    // Crear exchange
    err = ch.ExchangeDeclare(
        "edugo.events",  // name
        "topic",         // type
        true,            // durable
        false,           // auto-deleted
        false,           // internal
        false,           // no-wait
        nil,             // arguments
    )
    if err != nil {
        return err
    }
    
    // Crear colas necesarias
    queues := []string{
        "material.created",
        "assessment.completed",
        "progress.updated",
    }
    
    for _, queueName := range queues {
        _, err = ch.QueueDeclare(
            queueName,
            true,  // durable
            false, // delete when unused
            false, // exclusive
            false, // no-wait
            nil,   // arguments
        )
        if err != nil {
            return err
        }
        
        // Bind queue to exchange
        err = ch.QueueBind(
            queueName,
            queueName,
            "edugo.events",
            false,
            nil,
        )
        if err != nil {
            return err
        }
    }
    
    return nil
}
```


### 4. Sistema de Gestión de Datos de Prueba

#### Componente: TestDataManager

**Responsabilidad**: Proporcionar funciones helper para crear y gestionar datos de prueba

**Mejoras a testhelpers.go**:

```go
// MEJORAR: Agregar comentarios con valores sin encriptar
func SeedTestUser(t *testing.T, db *sql.DB) (userID string, email string) {
    t.Helper()
    
    email = "test@edugo.com"
    password := "Test1234!"  // Contraseña sin encriptar: Test1234!
    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        t.Fatalf("Failed to hash password: %v", err)
    }
    
    query := `
        INSERT INTO users (email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
        VALUES ($1, $2, 'Test', 'User', 'student', true, NOW(), NOW())
        RETURNING id
    `
    err = db.QueryRow(query, email, string(hashedPassword)).Scan(&userID)
    if err != nil {
        t.Fatalf("Failed to seed test user: %v", err)
    }
    
    t.Logf("👤 Test user created: %s (email: %s, password: %s)", userID, email, password)
    return userID, email
}

// NUEVO: Seed múltiples usuarios con roles diferentes
func SeedTestUsers(t *testing.T, db *sql.DB, count int, role string) []TestUser {
    t.Helper()
    
    users := make([]TestUser, count)
    for i := 0; i < count; i++ {
        email := fmt.Sprintf("test%d@edugo.com", i+1)
        password := "Test1234!"
        
        hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        
        var userID string
        query := `
            INSERT INTO users (email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, true, NOW(), NOW())
            RETURNING id
        `
        err := db.QueryRow(query, email, string(hashedPassword), 
            fmt.Sprintf("Test%d", i+1), "User", role).Scan(&userID)
        if err != nil {
            t.Fatalf("Failed to seed test user %d: %v", i+1, err)
        }
        
        users[i] = TestUser{
            ID:       userID,
            Email:    email,
            Password: password,  // Guardar sin encriptar para tests
            Role:     role,
        }
    }
    
    t.Logf("👥 Created %d test users with role: %s", count, role)
    return users
}

type TestUser struct {
    ID       string
    Email    string
    Password string  // Sin encriptar para uso en tests
    Role     string
}

// NUEVO: Seed completo de escenario de prueba
func SeedCompleteTestScenario(t *testing.T, db *sql.DB, mongodb *mongo.Database) *TestScenario {
    t.Helper()
    
    // Crear usuarios
    teacher, _ := SeedTestUserWithEmail(t, db, "teacher@edugo.com")
    student1, _ := SeedTestUserWithEmail(t, db, "student1@edugo.com")
    student2, _ := SeedTestUserWithEmail(t, db, "student2@edugo.com")
    
    // Crear materiales
    material1 := SeedTestMaterialWithTitle(t, db, teacher, "Introducción a Go")
    material2 := SeedTestMaterialWithTitle(t, db, teacher, "Testing en Go")
    
    // Crear assessments
    assessment1 := SeedTestAssessment(t, mongodb, material1)
    assessment2 := SeedTestAssessment(t, mongodb, material2)
    
    return &TestScenario{
        Teacher:     teacher,
        Students:    []string{student1, student2},
        Materials:   []string{material1, material2},
        Assessments: []string{assessment1, assessment2},
    }
}

type TestScenario struct {
    Teacher     string
    Students    []string
    Materials   []string
    Assessments []string
}
```


### 5. Scripts de Setup para Desarrollo

#### Componente: DevEnvironmentScripts

**Responsabilidad**: Proporcionar scripts para configurar ambiente de desarrollo usando infraestructura de tests

**Script**: `test/scripts/setup_dev_env.sh`
```bash
#!/bin/bash
set -e

echo "🚀 Setting up development environment..."

# Verificar Docker
if ! docker ps > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker Desktop."
    exit 1
fi

# Levantar contenedores
echo "🐳 Starting containers..."
docker-compose -f docker-compose-dev.yml up -d

# Esperar a que estén listos
echo "⏳ Waiting for services to be ready..."
sleep 5

# Ejecutar schema SQL
echo "🗄️  Creating PostgreSQL schema..."
docker exec edugo-postgres-dev psql -U edugo_user -d edugo -f /scripts/schema.sql

# Cargar datos de prueba
echo "📊 Loading test data..."
docker exec edugo-postgres-dev psql -U edugo_user -d edugo -f /scripts/seed_data.sql

# Configurar MongoDB
echo "🍃 Setting up MongoDB..."
docker exec edugo-mongo-dev mongosh edugo --eval "
    db.createCollection('material_assessments');
    db.material_assessments.createIndex({ material_id: 1 }, { unique: true });
    db.createCollection('assessment_results');
    db.assessment_results.createIndex({ assessment_id: 1, user_id: 1 }, { unique: true });
"

# Configurar RabbitMQ
echo "🐰 Setting up RabbitMQ..."
docker exec edugo-rabbitmq-dev rabbitmqadmin declare exchange name=edugo.events type=topic durable=true
docker exec edugo-rabbitmq-dev rabbitmqadmin declare queue name=material.created durable=true
docker exec edugo-rabbitmq-dev rabbitmqadmin declare queue name=assessment.completed durable=true
docker exec edugo-rabbitmq-dev rabbitmqadmin declare binding source=edugo.events destination=material.created routing_key=material.created

echo "✅ Development environment ready!"
echo ""
echo "📝 Connection strings:"
echo "  PostgreSQL: postgresql://edugo_user:edugo_pass@localhost:5432/edugo"
echo "  MongoDB:    mongodb://edugo_admin:edugo_pass@localhost:27017/edugo"
echo "  RabbitMQ:   amqp://edugo_user:edugo_pass@localhost:5672/"
echo ""
echo "🌐 Web interfaces:"
echo "  RabbitMQ Management: http://localhost:15672 (user: edugo_user, pass: edugo_pass)"
```

**Script**: `test/scripts/teardown_dev_env.sh`
```bash
#!/bin/bash
set -e

echo "🧹 Tearing down development environment..."

docker-compose -f docker-compose-dev.yml down -v

echo "✅ Development environment cleaned up!"
```

**Docker Compose**: `docker-compose-dev.yml`
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: edugo-postgres-dev
    environment:
      POSTGRES_DB: edugo
      POSTGRES_USER: edugo_user
      POSTGRES_PASSWORD: edugo_pass
    ports:
      - "5432:5432"
    volumes:
      - ./scripts/postgresql:/scripts
      - postgres-dev-data:/var/lib/postgresql/data

  mongodb:
    image: mongo:7.0
    container_name: edugo-mongo-dev
    environment:
      MONGO_INITDB_ROOT_USERNAME: edugo_admin
      MONGO_INITDB_ROOT_PASSWORD: edugo_pass
      MONGO_INITDB_DATABASE: edugo
    ports:
      - "27017:27017"
    volumes:
      - mongo-dev-data:/data/db

  rabbitmq:
    image: rabbitmq:3.12-management-alpine
    container_name: edugo-rabbitmq-dev
    environment:
      RABBITMQ_DEFAULT_USER: edugo_user
      RABBITMQ_DEFAULT_PASS: edugo_pass
    ports:
      - "5672:5672"
      - "15672:15672"
    volumes:
      - rabbitmq-dev-data:/var/lib/rabbitmq

volumes:
  postgres-dev-data:
  mongo-dev-data:
  rabbitmq-dev-data:
```


### 6. Mejoras al Makefile

**Nuevos comandos para testing**:

```makefile
# ============================================
# Testing Avanzado
# ============================================

test-unit: ## Solo tests unitarios (rápido)
	@echo "$(YELLOW)🧪 Ejecutando tests unitarios...$(RESET)"
	@go test -v -short -race ./internal/... -timeout 2m
	@echo "$(GREEN)✓ Tests unitarios completados$(RESET)"

test-unit-coverage: ## Tests unitarios con cobertura
	@echo "$(YELLOW)📊 Tests unitarios con cobertura...$(RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@go test -v -short -race -coverprofile=$(COVERAGE_DIR)/unit-coverage.out ./internal/... -timeout 2m
	@./scripts/filter-coverage.sh $(COVERAGE_DIR)/unit-coverage.out
	@go tool cover -html=$(COVERAGE_DIR)/unit-coverage-filtered.out -o $(COVERAGE_DIR)/unit-coverage.html
	@echo "$(GREEN)✓ Reporte: $(COVERAGE_DIR)/unit-coverage.html$(RESET)"

test-integration-verbose: ## Tests de integración con logs detallados
	@echo "$(YELLOW)🐳 Tests de integración (verbose)...$(RESET)"
	@RUN_INTEGRATION_TESTS=true go test -v -tags=integration ./test/integration/... -timeout 10m

test-all: test-unit test-integration ## Ejecutar todos los tests

test-watch: ## Watch mode para tests (requiere entr)
	@echo "$(YELLOW)👀 Watching tests...$(RESET)"
	@find . -name "*.go" | entr -c make test-unit

coverage-report: ## Generar reporte de cobertura completo
	@echo "$(YELLOW)📊 Generando reporte de cobertura completo...$(RESET)"
	@mkdir -p $(COVERAGE_DIR)
	@go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./... -timeout 5m
	@./scripts/filter-coverage.sh $(COVERAGE_DIR)/coverage.out
	@go tool cover -html=$(COVERAGE_DIR)/coverage-filtered.out -o $(COVERAGE_DIR)/coverage.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage-filtered.out | tail -20
	@echo "$(GREEN)✓ Reporte: $(COVERAGE_DIR)/coverage.html$(RESET)"

coverage-check: ## Verificar que cobertura cumple umbral mínimo
	@echo "$(YELLOW)🎯 Verificando cobertura mínima...$(RESET)"
	@go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./... -timeout 5m
	@./scripts/check-coverage.sh $(COVERAGE_DIR)/coverage.out 60
	@echo "$(GREEN)✓ Cobertura cumple umbral mínimo$(RESET)"

# ============================================
# Desarrollo Local
# ============================================

dev-setup: ## Configurar ambiente de desarrollo con Docker
	@echo "$(YELLOW)🚀 Configurando ambiente de desarrollo...$(RESET)"
	@./test/scripts/setup_dev_env.sh

dev-teardown: ## Limpiar ambiente de desarrollo
	@echo "$(YELLOW)🧹 Limpiando ambiente de desarrollo...$(RESET)"
	@./test/scripts/teardown_dev_env.sh

dev-reset: dev-teardown dev-setup ## Resetear ambiente de desarrollo

dev-logs: ## Ver logs de contenedores de desarrollo
	@docker-compose -f docker-compose-dev.yml logs -f

# ============================================
# Análisis de Tests
# ============================================

test-analyze: ## Analizar estructura de tests
	@echo "$(YELLOW)🔍 Analizando estructura de tests...$(RESET)"
	@go run ./tools/test-analyzer/main.go

test-missing: ## Identificar módulos sin tests
	@echo "$(YELLOW)🔍 Buscando módulos sin tests...$(RESET)"
	@./scripts/find-missing-tests.sh

test-validate: ## Validar que todos los tests pasan
	@echo "$(YELLOW)✅ Validando tests...$(RESET)"
	@make test-unit
	@make test-integration
	@echo "$(GREEN)✓ Todos los tests pasan$(RESET)"
```


## Modelos de Datos

### Estructura de Reporte de Análisis

```go
// TestAnalysisReport representa el resultado del análisis de tests
type TestAnalysisReport struct {
    Timestamp           time.Time
    ProjectPath         string
    Summary             TestSummary
    CoverageAnalysis    CoverageAnalysis
    QualityMetrics      QualityMetrics
    Recommendations     []Recommendation
}

type TestSummary struct {
    TotalTestFiles      int
    UnitTestFiles       int
    IntegrationTestFiles int
    TotalTestCases      int
    PassingTests        int
    FailingTests        int
    SkippedTests        int
}

type CoverageAnalysis struct {
    OverallCoverage     float64
    PackageCoverage     map[string]PackageCoverage
    UncoveredFiles      []string
    CriticalGaps        []CoverageGap
}

type PackageCoverage struct {
    PackageName         string
    Coverage            float64
    TotalStatements     int
    CoveredStatements   int
    UncoveredFunctions  []string
}

type CoverageGap struct {
    PackageName         string
    CurrentCoverage     float64
    TargetCoverage      float64
    Priority            string  // "critical", "high", "medium", "low"
    Reason              string
}

type QualityMetrics struct {
    TestsWithMocks      int
    TestsWithoutMocks   int
    TestsFollowingAAA   int
    TestsWithCleanup    int
    AverageTestDuration time.Duration
}

type Recommendation struct {
    Type        string  // "coverage", "quality", "structure"
    Priority    string  // "critical", "high", "medium", "low"
    Description string
    Action      string
    Effort      string  // "low", "medium", "high"
}
```

### Configuración de Cobertura

```go
// CoverageConfig representa la configuración de exclusiones
type CoverageConfig struct {
    ExcludedPatterns    []string
    ExcludedPackages    []string
    MinimumCoverage     float64
    CriticalPackages    map[string]float64  // Paquete -> Cobertura mínima
}

// Ejemplo de configuración
var DefaultCoverageConfig = CoverageConfig{
    ExcludedPatterns: []string{
        "docs/",
        "*_mock.go",
        "*/mocks_test.go",
        "cmd/",
        "internal/application/dto/",
    },
    ExcludedPackages: []string{
        "github.com/EduGoGroup/edugo-api-mobile/docs",
        "github.com/EduGoGroup/edugo-api-mobile/cmd",
        "github.com/EduGoGroup/edugo-api-mobile/internal/bootstrap/noop",
    },
    MinimumCoverage: 60.0,
    CriticalPackages: map[string]float64{
        "internal/application/service":     70.0,
        "internal/domain/entity":           80.0,
        "internal/domain/valueobject":      80.0,
        "internal/infrastructure/http/handler": 60.0,
    },
}
```


## Manejo de Errores

### Estrategia de Manejo de Errores en Tests

1. **Tests Unitarios**:
   - Usar `t.Fatal()` para errores de setup que impiden continuar
   - Usar `t.Error()` para fallos de assertions que permiten continuar
   - Usar `assert.NoError()` y `require.NoError()` de testify apropiadamente

2. **Tests de Integración**:
   - Fallar rápido si testcontainers no se pueden levantar
   - Usar `t.Skip()` si recursos opcionales no están disponibles (ej: RabbitMQ)
   - Limpiar recursos incluso si el test falla (defer cleanup)

3. **Análisis de Tests**:
   - Capturar y reportar errores sin detener el análisis completo
   - Generar reportes parciales si algunos análisis fallan
   - Logging detallado de errores para debugging

### Casos de Error Comunes

```go
// Error: Docker no disponible
func TestIntegration_DockerNotAvailable(t *testing.T) {
    SkipIfIntegrationTestsDisabled(t)
    
    if !isDockerAvailable() {
        t.Skip("Docker is not available, skipping integration test")
    }
    
    // ... resto del test
}

// Error: Testcontainer falla al iniciar
func SetupContainers(t *testing.T) (*TestContainers, func()) {
    ctx := context.Background()
    
    pgContainer, err := postgres.Run(ctx, "postgres:15-alpine", ...)
    if err != nil {
        t.Fatalf("Failed to start PostgreSQL container: %v\nMake sure Docker is running", err)
    }
    
    // ... resto del setup
}

// Error: Seed de datos falla
func SeedTestUser(t *testing.T, db *sql.DB) (string, string) {
    t.Helper()
    
    var userID string
    err := db.QueryRow(query, ...).Scan(&userID)
    if err != nil {
        t.Fatalf("Failed to seed test user: %v\nQuery: %s", err, query)
    }
    
    return userID, email
}

// Error: Cobertura por debajo del umbral
func CheckCoverageThreshold(coverageFile string, threshold float64) error {
    coverage, err := parseCoverage(coverageFile)
    if err != nil {
        return fmt.Errorf("failed to parse coverage: %w", err)
    }
    
    if coverage < threshold {
        return fmt.Errorf("coverage %.2f%% is below threshold %.2f%%", coverage, threshold)
    }
    
    return nil
}
```

## Estrategia de Testing

### Pirámide de Testing

```
                    /\
                   /  \
                  / E2E \          <- 10% (Tests de integración completos)
                 /--------\
                /          \
               /  Integration \    <- 20% (Tests de integración de componentes)
              /--------------\
             /                \
            /   Unit Tests     \  <- 70% (Tests unitarios)
           /--------------------\
```

### Tipos de Tests y Cuándo Usarlos

1. **Tests Unitarios** (70% de los tests):
   - **Qué**: Funciones puras, lógica de negocio, validaciones
   - **Dónde**: Junto al código fuente (`*_test.go`)
   - **Dependencias**: Mocks para todas las dependencias externas
   - **Velocidad**: Muy rápido (< 100ms por test)
   - **Ejemplos**:
     - Validación de value objects (Email, MaterialID)
     - Lógica de scoring en servicios
     - Transformaciones de DTOs
     - Estrategias de scoring (multiple choice, true/false)

2. **Tests de Integración de Componentes** (20% de los tests):
   - **Qué**: Interacción entre capas (handler -> service -> repository)
   - **Dónde**: `test/integration/`
   - **Dependencias**: Testcontainers para BD, mocks para servicios externos
   - **Velocidad**: Medio (1-5s por test)
   - **Ejemplos**:
     - Repository con PostgreSQL real
     - Handler con service real pero repository mockeado
     - Service con repository real

3. **Tests E2E** (10% de los tests):
   - **Qué**: Flujos completos de usuario
   - **Dónde**: `test/integration/`
   - **Dependencias**: Todos los servicios reales (testcontainers)
   - **Velocidad**: Lento (5-20s por test)
   - **Ejemplos**:
     - Flujo completo de autenticación
     - Crear material -> crear assessment -> completar assessment
     - Progreso de usuario en múltiples materiales


## Plan de Implementación

### Fase 1: Análisis y Evaluación (Requisitos 1, 2, 11)

**Objetivo**: Entender el estado actual y validar que los tests existentes funcionan

**Tareas**:
1. Crear herramienta de análisis de estructura de tests
2. Ejecutar análisis de cobertura actual
3. Validar que todos los tests existentes pasan
4. Identificar módulos sin tests
5. Generar reporte de estado actual

**Entregables**:
- Reporte de análisis de tests (`docs/TEST_ANALYSIS_REPORT.md`)
- Lista de módulos sin cobertura
- Validación de tests existentes

**Criterios de Éxito**:
- Todos los tests unitarios existentes pasan (100%)
- Todos los tests de integración existentes pasan (100%)
- Reporte de cobertura generado con métricas por paquete
- Identificados todos los módulos críticos sin tests

### Fase 2: Configuración y Refactorización (Requisitos 3, 4, 5, 6, 7)

**Objetivo**: Establecer la infraestructura y configuración base para testing

**Tareas**:
1. Crear archivo `.coverignore` con exclusiones
2. Crear scripts de filtrado de cobertura
3. Eliminar carpetas vacías en `test/unit/`
4. Mejorar helpers de testcontainers con configuración de RabbitMQ
5. Agregar funciones helper para seed de datos complejos
6. Documentar contraseñas sin encriptar en comentarios
7. Crear scripts de setup para desarrollo local

**Entregables**:
- `.coverignore` configurado
- `scripts/filter-coverage.sh`
- `scripts/check-coverage.sh`
- `test/integration/testhelpers.go` mejorado
- `test/scripts/setup_dev_env.sh`
- `test/scripts/teardown_dev_env.sh`
- `docker-compose-dev.yml`

**Criterios de Éxito**:
- Cobertura se calcula excluyendo archivos configurados
- Helpers de seed documentan valores sin encriptar
- RabbitMQ se configura automáticamente en tests
- Scripts de desarrollo funcionan correctamente

### Fase 3: Mejora de Cobertura (Requisitos 9, 10)

**Objetivo**: Incrementar cobertura de tests en módulos críticos

**Tareas**:
1. Crear tests para value objects sin cobertura
2. Crear tests para entities de dominio
3. Crear tests para repositories
4. Crear tests para servicios con baja cobertura
5. Crear tests para handlers con baja cobertura
6. Documentar guías de testing

**Entregables**:
- Tests para `internal/domain/valueobject/`
- Tests para `internal/domain/entity/`
- Tests para `internal/infrastructure/persistence/*/repository/`
- Guía de testing (`docs/TESTING_GUIDE.md`)
- Plantillas de tests

**Criterios de Éxito**:
- Cobertura de servicios >= 70%
- Cobertura de dominio >= 80%
- Cobertura de handlers >= 60%
- Cobertura general >= 60%

### Fase 4: Automatización y CI/CD (Requisito 12)

**Objetivo**: Integrar testing en el pipeline de CI/CD

**Tareas**:
1. Actualizar Makefile con nuevos comandos
2. Configurar GitHub Actions para ejecutar tests
3. Configurar reporte de cobertura en CI
4. Configurar umbral mínimo de cobertura
5. Configurar badges de cobertura

**Entregables**:
- Makefile actualizado
- `.github/workflows/test.yml`
- `.github/workflows/coverage.yml`
- Badges en README

**Criterios de Éxito**:
- Tests se ejecutan automáticamente en cada PR
- Build falla si cobertura cae por debajo del umbral
- Reportes de cobertura se publican automáticamente
- Badges muestran estado actual de tests y cobertura


## Decisiones de Diseño y Justificaciones

### 1. Tests Unitarios Junto al Código vs Carpeta Separada

**Decisión**: Mantener tests unitarios junto al código fuente (`*_test.go` en el mismo paquete)

**Justificación**:
- **Descubrimiento**: Es más fácil encontrar los tests relacionados con un archivo
- **Mantenimiento**: Al modificar código, los tests están inmediatamente visibles
- **Convención Go**: Es la práctica estándar en la comunidad Go
- **Tooling**: Las herramientas de Go esperan esta estructura
- **Cobertura**: `go test` automáticamente encuentra y ejecuta estos tests

**Alternativa Rechazada**: Carpeta `test/unit/` separada
- Requiere duplicar estructura de carpetas
- Dificulta el descubrimiento de tests
- No es idiomático en Go
- La carpeta actual solo tiene `.gitkeep` (vacía)

**Acción**: Eliminar `test/unit/` y sus subcarpetas vacías

### 2. Tests de Integración en Carpeta Separada

**Decisión**: Mantener tests de integración en `test/integration/` con build tags

**Justificación**:
- **Separación de Concerns**: Tests de integración son diferentes (lentos, requieren Docker)
- **Control de Ejecución**: Build tags permiten ejecutarlos selectivamente
- **Infraestructura Compartida**: Helpers y setup se reutilizan entre tests
- **Desarrollo Local**: Infraestructura se puede reutilizar para desarrollo

**Implementación Actual**: Ya está bien implementado, solo necesita mejoras menores

### 3. Exclusiones de Cobertura

**Decisión**: Usar archivo `.coverignore` y script de filtrado

**Justificación**:
- **Claridad**: Archivo dedicado es más claro que flags en comandos
- **Mantenibilidad**: Fácil agregar/quitar exclusiones
- **Documentación**: El archivo sirve como documentación de qué se excluye y por qué
- **Flexibilidad**: Script permite lógica compleja de filtrado

**Alternativa Rechazada**: Flags en `go test`
- Comandos muy largos y difíciles de mantener
- No hay forma nativa de excluir patrones en Go

### 4. RabbitMQ Opcional en Tests

**Decisión**: RabbitMQ es opcional, usar mock si falla

**Justificación**:
- **Robustez**: Tests no fallan si RabbitMQ tiene problemas
- **Velocidad**: Mock es más rápido para tests que no necesitan mensajería real
- **Flexibilidad**: Tests críticos pueden requerir RabbitMQ real, otros no

**Implementación**:
```go
publisher, err := createTestRabbitMQPublisher(rabbitConnStr, testLogger)
if err != nil {
    t.Logf("⚠️  Warning: RabbitMQ publisher failed (non-critical): %v", err)
    publisher = &mockPublisher{}  // Fallback a mock
}
```

### 5. Helpers Centralizados vs Duplicados

**Decisión**: Helpers centralizados en `test/integration/testhelpers.go`

**Justificación**:
- **DRY**: No duplicar código de setup entre tests
- **Consistencia**: Todos los tests usan los mismos helpers
- **Mantenibilidad**: Cambios en helpers se propagan a todos los tests
- **Documentación**: Un solo lugar para documentar cómo crear datos de prueba

### 6. Documentar Contraseñas Sin Encriptar

**Decisión**: Agregar comentarios con valores sin encriptar en helpers de seed

**Justificación**:
- **Usabilidad**: Desarrolladores necesitan saber las contraseñas para tests manuales
- **Debugging**: Facilita debugging de problemas de autenticación
- **Documentación**: Sirve como documentación de datos de prueba
- **Seguridad**: No es un problema porque son datos de prueba, no producción

**Ejemplo**:
```go
password := "Test1234!"  // Contraseña sin encriptar: Test1234!
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

### 7. Reutilización de Infraestructura para Desarrollo

**Decisión**: Proporcionar scripts que reutilizan testcontainers para desarrollo local

**Justificación**:
- **Consistencia**: Mismo ambiente en tests y desarrollo
- **Eficiencia**: No duplicar configuración de Docker
- **Facilidad**: Un comando para levantar todo el ambiente
- **Datos de Prueba**: Mismo seed data en tests y desarrollo

**Implementación**: Scripts `setup_dev_env.sh` y `docker-compose-dev.yml`


## Consideraciones de Rendimiento

### Optimización de Tests Unitarios

1. **Paralelización**:
   ```go
   func TestSomething(t *testing.T) {
       t.Parallel()  // Ejecutar en paralelo con otros tests
       // ... test code
   }
   ```

2. **Mocks Eficientes**:
   - Usar mocks en lugar de dependencias reales
   - Reutilizar mocks entre tests cuando sea posible
   - Evitar setup complejo en cada test

3. **Tabla de Tests**:
   ```go
   func TestValidation(t *testing.T) {
       tests := []struct {
           name    string
           input   string
           wantErr bool
       }{
           {"valid email", "test@example.com", false},
           {"invalid email", "invalid", true},
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               t.Parallel()
               // ... test logic
           })
       }
   }
   ```

### Optimización de Tests de Integración

1. **Reutilización de Contenedores**:
   - Levantar contenedores una vez por suite de tests
   - Limpiar datos entre tests en lugar de recrear contenedores

2. **Cleanup Eficiente**:
   ```go
   func CleanDatabase(t *testing.T, db *sql.DB) {
       t.Helper()
       // TRUNCATE es más rápido que DELETE
       tables := []string{"refresh_tokens", "materials", "users"}
       for _, table := range tables {
           db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
       }
   }
   ```

3. **Timeouts Apropiados**:
   ```go
   // Tests de integración pueden tardar más
   go test -timeout 10m -tags=integration ./test/integration/...
   ```

### Métricas de Rendimiento Esperadas

| Tipo de Test | Tiempo Promedio | Timeout Recomendado |
|--------------|-----------------|---------------------|
| Test Unitario | < 100ms | 2m |
| Test de Integración (componente) | 1-5s | 5m |
| Test E2E | 5-20s | 10m |
| Suite Completa Unitaria | < 5s | 5m |
| Suite Completa Integración | 3-6 min | 15m |

## Seguridad en Testing

### Datos Sensibles en Tests

1. **Contraseñas**:
   - Usar contraseñas de prueba simples y documentadas
   - Nunca usar contraseñas reales de producción
   - Documentar valores sin encriptar en comentarios

2. **Tokens y Secrets**:
   - Usar valores hardcodeados para tests
   - Nunca usar secrets reales de producción
   - Ejemplo: `jwtSecret := "test-jwt-secret-key-very-secure-for-testing-only"`

3. **Datos de Usuarios**:
   - Usar emails de prueba claramente identificables (`test@edugo.com`)
   - Usar nombres genéricos (`Test User`, `Student 1`)
   - No usar datos personales reales

### Aislamiento de Tests

1. **Bases de Datos**:
   - Cada test debe limpiar sus datos
   - Usar transacciones cuando sea posible
   - Testcontainers proporciona aislamiento completo

2. **Estado Global**:
   - Evitar variables globales en tests
   - Cada test debe ser independiente
   - Usar `t.Cleanup()` para garantizar limpieza

3. **Concurrencia**:
   - Tests paralelos no deben compartir estado
   - Usar `t.Parallel()` solo cuando sea seguro
   - Cuidado con race conditions

## Documentación y Guías

### Estructura de Documentación

```
docs/
├── TESTING_GUIDE.md              # Guía principal de testing
├── TESTING_UNIT_GUIDE.md         # Guía específica de tests unitarios
├── TESTING_INTEGRATION_GUIDE.md  # Guía específica de tests de integración
├── TEST_ANALYSIS_REPORT.md       # Reporte de análisis actual
└── TEST_COVERAGE_PLAN.md         # Plan para mejorar cobertura
```

### Contenido de Guías

**TESTING_GUIDE.md**:
- Filosofía de testing del proyecto
- Tipos de tests y cuándo usarlos
- Estructura de carpetas
- Comandos make disponibles
- Mejores prácticas

**TESTING_UNIT_GUIDE.md**:
- Cómo escribir tests unitarios
- Uso de mocks
- Patrón AAA (Arrange-Act-Assert)
- Ejemplos de tests por tipo de componente
- Plantillas de tests

**TESTING_INTEGRATION_GUIDE.md**:
- Cómo escribir tests de integración
- Uso de testcontainers
- Helpers disponibles
- Seed de datos
- Troubleshooting

**TEST_ANALYSIS_REPORT.md**:
- Estado actual de tests
- Cobertura por módulo
- Módulos sin tests
- Recomendaciones

**TEST_COVERAGE_PLAN.md**:
- Metas de cobertura por módulo
- Priorización de tests faltantes
- Timeline de implementación
- Responsables

## Conclusión

Este diseño proporciona una estrategia completa y robusta para mejorar el sistema de testing del proyecto edugo-api-mobile. La implementación se divide en 4 fases claras, cada una con objetivos, tareas y criterios de éxito bien definidos.

**Puntos Clave**:
1. Tests unitarios junto al código (idiomático Go)
2. Tests de integración separados con testcontainers
3. Exclusiones de cobertura configurables
4. Infraestructura reutilizable para desarrollo
5. Helpers centralizados y bien documentados
6. Automatización completa en CI/CD

**Próximos Pasos**:
1. Revisar y aprobar este diseño
2. Proceder a crear el plan de implementación detallado (tasks.md)
3. Ejecutar Fase 1: Análisis y Evaluación
