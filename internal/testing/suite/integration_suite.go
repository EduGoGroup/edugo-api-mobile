//go:build integration
// +build integration

package suite

import (
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"time"

	infrastructureTesting "github.com/EduGoGroup/edugo-infrastructure/migrations"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"go.mongodb.org/mongo-driver/mongo"
)

// IntegrationTestSuite es la suite base para tests de integración
// Comparte contenedores Docker entre todos los tests para mejor performance
type IntegrationTestSuite struct {
	suite.Suite
	
	// Contenedores Docker (compartidos entre tests)
	PostgresContainer *postgres.PostgresContainer
	MongoContainer    *mongodb.MongoDBContainer
	RabbitContainer   *rabbitmq.RabbitMQContainer
	
	// Conexiones a recursos (compartidas entre tests)
	PostgresDB *sql.DB
	MongoDB    *mongo.Database
	Logger     logger.Logger
	
	// Contexto de la suite
	ctx context.Context
	
	// Paths a scripts de infrastructure
	migrationsPath string
	seedsPath      string
}

// SetupSuite se ejecuta UNA VEZ antes de todos los tests
// Levanta contenedores y aplica migraciones/seeds
func (s *IntegrationTestSuite) SetupSuite() {
	s.ctx = context.Background()
	
	// Inicializar logger
	s.Logger = logger.NewZapLogger("info", "json")
	s.Logger.Info("🚀 Iniciando suite de integración...")
	
	// Calcular paths a scripts de infrastructure
	s.calculateInfrastructurePaths()
	
	// Levantar contenedores compartidos
	s.startContainers()
	
	// Aplicar migraciones de infrastructure
	s.applyMigrations()
	
	// Aplicar seeds de infrastructure
	s.applySeeds()
	
	s.Logger.Info("✅ Suite de integración lista")
}

// TearDownSuite se ejecuta UNA VEZ después de todos los tests
// Limpia contenedores
func (s *IntegrationTestSuite) TearDownSuite() {
	s.Logger.Info("🧹 Limpiando suite de integración...")
	
	if s.PostgresDB != nil {
		_ = s.PostgresDB.Close()
	}
	
	if s.PostgresContainer != nil {
		_ = s.PostgresContainer.Terminate(s.ctx)
	}
	
	if s.MongoContainer != nil {
		_ = s.MongoContainer.Terminate(s.ctx)
	}
	
	if s.RabbitContainer != nil {
		_ = s.RabbitContainer.Terminate(s.ctx)
	}
	
	s.Logger.Info("✅ Suite limpiada")
}

// SetupTest se ejecuta ANTES de cada test individual
// Limpia datos pero NO reinicia contenedores
func (s *IntegrationTestSuite) SetupTest() {
	s.Logger.Info("🧪 Preparando test individual...")
	
	// Limpiar datos de PostgreSQL
	if err := infrastructureTesting.CleanDatabase(s.PostgresDB); err != nil {
		s.T().Fatalf("Error limpiando base de datos: %v", err)
	}
	
	// Re-aplicar seeds para tener datos frescos
	if err := infrastructureTesting.ApplySeeds(s.PostgresDB, s.seedsPath); err != nil {
		s.T().Fatalf("Error aplicando seeds: %v", err)
	}
	
	s.Logger.Info("✅ Test individual listo con datos frescos")
}

// startContainers levanta todos los contenedores de infraestructura
func (s *IntegrationTestSuite) startContainers() {
	s.Logger.Info("📦 Iniciando contenedores Docker...")
	
	// PostgreSQL
	s.Logger.Info("  - PostgreSQL...")
	pgContainer, err := postgres.Run(s.ctx, "postgres:16-alpine",
		postgres.WithDatabase("edugo_test"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_pass"),
	)
	s.Require().NoError(err, "PostgreSQL container debe iniciar")
	s.PostgresContainer = pgContainer
	
	// Conectar a PostgreSQL
	connStr, err := pgContainer.ConnectionString(s.ctx, "sslmode=disable")
	s.Require().NoError(err)
	
	db, err := sql.Open("postgres", connStr)
	s.Require().NoError(err)
	s.PostgresDB = db
	
	// MongoDB
	s.Logger.Info("  - MongoDB...")
	mongoContainer, err := mongodb.Run(s.ctx, "mongo:7.0",
		mongodb.WithUsername("test_admin"),
		mongodb.WithPassword("test_pass"),
	)
	s.Require().NoError(err, "MongoDB container debe iniciar")
	s.MongoContainer = mongoContainer
	
	// RabbitMQ
	s.Logger.Info("  - RabbitMQ...")
	rabbitContainer, err := rabbitmq.Run(s.ctx, "rabbitmq:3.12-management-alpine",
		rabbitmq.WithAdminUsername("test_user"),
		rabbitmq.WithAdminPassword("test_pass"),
	)
	s.Require().NoError(err, "RabbitMQ container debe iniciar")
	s.RabbitContainer = rabbitContainer
	
	// Esperar a que contenedores estén completamente listos
	time.Sleep(2 * time.Second)
	
	s.Logger.Info("✅ Contenedores iniciados")
}

// applyMigrations ejecuta migraciones de edugo-infrastructure
func (s *IntegrationTestSuite) applyMigrations() {
	s.Logger.Info("📝 Aplicando migraciones de infrastructure...")
	
	err := infrastructureTesting.ApplyMigrations(s.PostgresDB, s.migrationsPath)
	s.Require().NoError(err, "Migraciones deben aplicarse correctamente")
	
	s.Logger.Info("✅ Migraciones aplicadas")
}

// applySeeds ejecuta seeds de edugo-infrastructure
func (s *IntegrationTestSuite) applySeeds() {
	s.Logger.Info("🌱 Aplicando seeds de infrastructure...")
	
	err := infrastructureTesting.ApplySeeds(s.PostgresDB, s.seedsPath)
	s.Require().NoError(err, "Seeds deben aplicarse correctamente")
	
	s.Logger.Info("✅ Seeds aplicados")
}

// calculateInfrastructurePaths calcula paths relativos a edugo-infrastructure
// Asume que edugo-infrastructure está en el mismo nivel que edugo-api-mobile
func (s *IntegrationTestSuite) calculateInfrastructurePaths() {
	// Obtener directorio del archivo actual
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	
	// Subir desde internal/testing/suite hasta raíz del proyecto
	projectRoot := filepath.Join(currentDir, "..", "..", "..")
	
	// Calcular path a infrastructure (directorio hermano)
	infrastructureRoot := filepath.Join(projectRoot, "..", "edugo-infrastructure")
	
	s.migrationsPath = filepath.Join(infrastructureRoot, "database", "migrations", "postgres")
	s.seedsPath = filepath.Join(infrastructureRoot, "seeds", "postgres")
	
	s.Logger.Info("📁 Paths de infrastructure calculados",
		"migrations", s.migrationsPath,
		"seeds", s.seedsPath,
	)
}
