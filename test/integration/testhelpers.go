//go:build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/EduGoGroup/edugo-api-mobile/internal/container"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/storage/s3"
	"github.com/EduGoGroup/edugo-shared/logger"
	_ "github.com/lib/pq" // PostgreSQL driver
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// TestApp encapsula todo lo necesario para tests de integración
type TestApp struct {
	Container *container.Container
	DB        *sql.DB
	MongoDB   *mongo.Database
	Cleanup   func()
}

// testLogger es un logger mock para tests que no imprime nada
type testLogger struct{}

func (l *testLogger) Debug(msg string, fields ...interface{}) {}
func (l *testLogger) Info(msg string, fields ...interface{})  {}
func (l *testLogger) Warn(msg string, fields ...interface{})  {}
func (l *testLogger) Error(msg string, fields ...interface{}) {}
func (l *testLogger) Fatal(msg string, fields ...interface{}) {}
func (l *testLogger) Sync() error { return nil }
func (l *testLogger) With(fields ...interface{}) logger.Logger { return l }
func (l *testLogger) WithContext(ctx context.Context) logger.Logger { return l }

// SetupTestApp inicializa una aplicación completa para testing
// Levanta testcontainers y crea el Container DI con todas las dependencias
func SetupTestApp(t *testing.T) *TestApp {
	t.Helper()
	
	// Skip si tests están deshabilitados
	SkipIfIntegrationTestsDisabled(t)
	
	// Levantar testcontainers
	containers, cleanup := SetupContainers(t)
	
	ctx := context.Background()
	
	// Obtener connection strings
	pgConnStr, err := containers.Postgres.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		cleanup()
		t.Fatalf("Failed to get Postgres connection string: %v", err)
	}
	
	mongoConnStr, err := containers.MongoDB.ConnectionString(ctx)
	if err != nil {
		cleanup()
		t.Fatalf("Failed to get MongoDB connection string: %v", err)
	}
	
	rabbitConnStr, err := containers.RabbitMQ.AmqpURL(ctx)
	if err != nil {
		cleanup()
		t.Fatalf("Failed to get RabbitMQ connection string: %v", err)
	}
	
	// Conectar a PostgreSQL
	db, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		cleanup()
		t.Fatalf("Failed to open Postgres connection: %v", err)
	}
	
	// Verificar conexión
	if err := db.Ping(); err != nil {
		db.Close()
		cleanup()
		t.Fatalf("Failed to ping Postgres: %v", err)
	}
	
	t.Log("✅ PostgreSQL connected")
	
	// Crear schema básico para tests
	if err := initTestSchema(db); err != nil {
		db.Close()
		cleanup()
		t.Fatalf("Failed to init test schema: %v", err)
	}
	
	t.Log("✅ PostgreSQL schema initialized")
	
	// Conectar a MongoDB
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnStr))
	if err != nil {
		db.Close()
		cleanup()
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	
	// Verificar conexión
	if err := mongoClient.Ping(ctx, nil); err != nil {
		mongoClient.Disconnect(ctx)
		db.Close()
		cleanup()
		t.Fatalf("Failed to ping MongoDB: %v", err)
	}
	
	mongodb := mongoClient.Database("edugo")
	t.Log("✅ MongoDB connected")
	
	// Crear logger para tests (mock silencioso)
	testLogger := &testLogger{}
	
	// Crear RabbitMQ Publisher (opcional - puede fallar sin romper tests)
	publisher, err := createTestRabbitMQPublisher(rabbitConnStr, testLogger)
	if err != nil {
		t.Logf("⚠️  Warning: RabbitMQ publisher failed (non-critical): %v", err)
		// Usar mock publisher en lugar de fallar
		publisher = &mockPublisher{}
	}
	
	// Crear S3 Client (mock para tests)
	s3Client := createTestS3Client()
	
	// JWT Secret para tests
	jwtSecret := "test-jwt-secret-key-very-secure-for-testing-only"
	
	// Crear Container DI
	c := container.NewContainer(
		db,
		mongodb,
		publisher,
		s3Client,
		jwtSecret,
		testLogger,
	)
	
	t.Log("✅ Container DI initialized")
	
	// Cleanup extendido
	appCleanup := func() {
		t.Log("🧹 Cleaning up test app...")
		if c != nil {
			c.Close()
		}
		if mongodb != nil {
			mongoClient.Disconnect(ctx)
		}
		if db != nil {
			db.Close()
		}
		cleanup() // Terminar testcontainers
		t.Log("✅ Test app cleaned up")
	}
	
	return &TestApp{
		Container: c,
		DB:        db,
		MongoDB:   mongodb,
		Cleanup:   appCleanup,
	}
}

// GetPostgresConnString obtiene el connection string de PostgreSQL testcontainer
func GetPostgresConnString(t *testing.T, containers *TestContainers) string {
	t.Helper()
	ctx := context.Background()
	connStr, err := containers.Postgres.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to get Postgres connection string: %v", err)
	}
	return connStr
}

// GetMongoConnString obtiene el connection string de MongoDB testcontainer
func GetMongoConnString(t *testing.T, containers *TestContainers) string {
	t.Helper()
	ctx := context.Background()
	connStr, err := containers.MongoDB.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get MongoDB connection string: %v", err)
	}
	return connStr
}

// GetRabbitMQConnString obtiene el connection string de RabbitMQ testcontainer
func GetRabbitMQConnString(t *testing.T, containers *TestContainers) string {
	t.Helper()
	ctx := context.Background()
	connStr, err := containers.RabbitMQ.AmqpURL(ctx)
	if err != nil {
		t.Fatalf("Failed to get RabbitMQ connection string: %v", err)
	}
	return connStr
}

// CleanDatabase limpia todas las tablas de PostgreSQL para tests aislados
func CleanDatabase(t *testing.T, db *sql.DB) {
	t.Helper()
	
	// Limpiar en orden inverso de dependencias
	tables := []string{
		"refresh_tokens",
		"login_attempts",
		"assessment_results",
		"assessment_attempts",
		"assessments",
		"progress",
		"materials",
		"users",
	}
	
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Logf("⚠️  Warning: Failed to truncate %s: %v", table, err)
		}
	}
	
	t.Log("🧹 Database cleaned")
}

// CleanMongoCollections limpia las colecciones de MongoDB y crea índices necesarios
func CleanMongoCollections(t *testing.T, mongodb *mongo.Database) {
	t.Helper()
	
	// Lista de colecciones a limpiar
	collections := []string{
		"material_assessments",
		"assessment_attempts",
		"assessment_results",
	}
	
	ctx := context.Background()
	for _, collName := range collections {
		coll := mongodb.Collection(collName)
		if err := coll.Drop(ctx); err != nil {
			t.Logf("Warning: Failed to drop collection %s: %v", collName, err)
		}
	}
	
	// Crear índice UNIQUE en assessment_results (assessment_id, user_id) para prevenir duplicados
	resultsCollection := mongodb.Collection("assessment_results")
	indexModel := mongo.IndexModel{
		Keys: map[string]interface{}{
			"assessment_id": 1,
			"user_id":       1,
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := resultsCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		t.Logf("Warning: Failed to create unique index on assessment_results: %v", err)
	}
	
	t.Log("🧹 MongoDB collections cleaned")
}

// SeedTestUser crea un usuario de prueba en PostgreSQL
func SeedTestUser(t *testing.T, db *sql.DB) (userID string, email string) {
	t.Helper()
	
	email = "test@edugo.com"
	password := "password123"
	
	// Generar bcrypt hash dinámicamente
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate password hash: %v", err)
	}
	passwordHash := string(passwordHashBytes)
	
	err = db.QueryRow(`
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, 'Test', 'User', 'student', true)
		RETURNING id
	`, email, passwordHash).Scan(&userID)
	
	if err != nil {
		t.Fatalf("Failed to seed test user: %v", err)
	}
	
	t.Logf("👤 Test user created: %s (%s) with password: %s", email, userID, password)
	return userID, email
}

// SeedTestMaterial crea un material de prueba en PostgreSQL
func SeedTestMaterial(t *testing.T, db *sql.DB, authorID string) (materialID string) {
	t.Helper()
	return SeedTestMaterialWithTitle(t, db, authorID, "Test Material")
}

// SeedTestMaterialWithTitle crea un material de prueba con un título específico
func SeedTestMaterialWithTitle(t *testing.T, db *sql.DB, authorID, title string) (materialID string) {
	t.Helper()
	
	err := db.QueryRow(`
		INSERT INTO materials (author_id, title, description, status, processing_status)
		VALUES ($1, $2, 'Test material description', 'published', 'completed')
		RETURNING id
	`, authorID, title).Scan(&materialID)
	
	if err != nil {
		t.Fatalf("Failed to seed test material: %v", err)
	}
	
	t.Logf("📚 Test material created: %s (%s)", title, materialID)
	return materialID
}

// SeedTestAssessment crea un assessment de prueba en MongoDB
func SeedTestAssessment(t *testing.T, mongodb *mongo.Database, materialID string) (assessmentID string) {
	t.Helper()
	
	// Assessment ID es el mismo que el material ID
	assessmentID = materialID
	
	// Crear assessment con 2 preguntas de prueba
	assessment := map[string]interface{}{
		"material_id": materialID,
		"questions": []map[string]interface{}{
			{
				"id":            "q1",
				"text":          "¿Qué es Go?",
				"question_type": "multiple_choice",
				"options":       []string{"A) Un lenguaje de programación", "B) Una base de datos", "C) Un framework", "D) Un editor"},
				"answer":        "A",
				"points":        1,
			},
			{
				"id":            "q2",
				"text":          "¿Go es compilado o interpretado?",
				"question_type": "multiple_choice",
				"options":       []string{"A) Interpretado", "B) Compilado", "C) Ambos", "D) Ninguno"},
				"answer":        "B",
				"points":        1,
			},
		},
		"created_at": "2024-01-01T00:00:00Z",
	}
	
	// Insertar en la colección material_assessments
	collection := mongodb.Collection("material_assessments")
	_, err := collection.InsertOne(context.Background(), assessment)
	if err != nil {
		t.Fatalf("Failed to seed test assessment: %v", err)
	}
	
	t.Logf("📝 Test assessment created for material: %s", materialID)
	return assessmentID
}

// createTestRabbitMQPublisher crea un publisher de RabbitMQ para tests
func createTestRabbitMQPublisher(url string, log logger.Logger) (rabbitmq.Publisher, error) {
	// Crear publisher real con RabbitMQ testcontainer
	exchange := "edugo.events" // Exchange de testing
	publisher, err := rabbitmq.NewRabbitMQPublisher(url, exchange, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ publisher: %w", err)
	}
	return publisher, nil
}

// createTestS3Client crea un cliente S3 para tests
func createTestS3Client() *s3.S3Client {
	// Para tests, usamos configuración mock/local
	// En producción real, se usaría LocalStack o MinIO para tests
	config := s3.S3Config{
		Region:          "us-east-1",
		BucketName:      "test-edugo-materials",
		AccessKeyID:     "test-access-key",
		SecretAccessKey: "test-secret-key",
		Endpoint:        "", // Vacío = AWS real (o LocalStack si está configurado en env)
	}
	
	ctx := context.Background()
	// Crear logger simple para S3 (silencioso)
	testLogger := &testLogger{}
	
	// Si falla (esperado en tests sin AWS), retornar nil
	// Los tests que usen S3 deberán mockear o skipear
	client, _ := s3.NewS3Client(ctx, config, testLogger)
	return client
}

// initTestSchema crea el schema mínimo necesario para tests
func initTestSchema(db *sql.DB) error {
	// Schema SQL mínimo para tests
	// Solo las tablas esenciales para los tests básicos
	schema := `
		-- Users table
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'student',
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Refresh tokens table
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			token_hash VARCHAR(255) UNIQUE NOT NULL,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			client_info JSONB,
			expires_at TIMESTAMP NOT NULL,
			revoked BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Login attempts table
		CREATE TABLE IF NOT EXISTS login_attempts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			identifier VARCHAR(255) NOT NULL,
			success BOOLEAN NOT NULL,
			attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Materials table
		CREATE TABLE IF NOT EXISTS materials (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title VARCHAR(255) NOT NULL,
			description TEXT,
			author_id UUID NOT NULL REFERENCES users(id),
			subject_id VARCHAR(255),
			status VARCHAR(50) DEFAULT 'draft',
			processing_status VARCHAR(50) DEFAULT 'pending',
			s3_key VARCHAR(500),
			s3_url VARCHAR(1000),
			is_deleted BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Progress table
		CREATE TABLE IF NOT EXISTS progress (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id),
			material_id UUID NOT NULL REFERENCES materials(id),
			progress_percentage NUMERIC(5,2) DEFAULT 0.00,
			last_accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, material_id)
		);

		-- Assessments table (mock/simplified)
		CREATE TABLE IF NOT EXISTS assessments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			material_id UUID NOT NULL REFERENCES materials(id),
			title VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Assessment attempts table
		CREATE TABLE IF NOT EXISTS assessment_attempts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			assessment_id UUID NOT NULL REFERENCES assessments(id),
			user_id UUID NOT NULL REFERENCES users(id),
			attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Assessment results table
		CREATE TABLE IF NOT EXISTS assessment_results (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			assessment_id UUID NOT NULL REFERENCES assessments(id),
			user_id UUID NOT NULL REFERENCES users(id),
			score NUMERIC(5,2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.Exec(schema)
	return err
}

// mockPublisher es un publisher mock para tests que no hace nada
type mockPublisher struct{}

func (m *mockPublisher) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	// No-op: tests no requieren RabbitMQ real
	return nil
}

func (m *mockPublisher) Close() error {
	return nil
}
