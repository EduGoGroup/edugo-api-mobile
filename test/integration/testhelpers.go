//go:build integration

package integration

import (
	"context"
	"github.com/EduGoGroup/edugo-shared/testing/containers"
	"database/sql"
	"fmt"
	"testing"

	"github.com/EduGoGroup/edugo-api-mobile/internal/bootstrap"
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

func (l *testLogger) Debug(msg string, fields ...interface{})       {}
func (l *testLogger) Info(msg string, fields ...interface{})        {}
func (l *testLogger) Warn(msg string, fields ...interface{})        {}
func (l *testLogger) Error(msg string, fields ...interface{})       {}
func (l *testLogger) Fatal(msg string, fields ...interface{})       {}
func (l *testLogger) Sync() error                                   { return nil }
func (l *testLogger) With(fields ...interface{}) logger.Logger      { return l }
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

	// Crear Resources para Container DI
	resources := &bootstrap.Resources{
		Logger:            testLogger,
		PostgreSQL:        db,
		MongoDB:           mongodb,
		RabbitMQPublisher: publisher,
		S3Client:          s3Client,
		JWTSecret:         jwtSecret,
	}

	// Crear Container DI
	c := container.NewContainer(resources)

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
// SeedTestUser crea un usuario de prueba con credenciales conocidas
// Contraseña sin encriptar: Test1234!
// Email: test@edugo.com
// Role: student
func SeedTestUser(t *testing.T, db *sql.DB) (userID string, email string) {
	t.Helper()

	email = "test@edugo.com"
	password := "Test1234!" // Contraseña sin encriptar para uso en tests de login

	// Generar hash de contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Insertar usuario usando el ID generado dinámicamente
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

// SeedTestUserWithEmail crea un usuario de prueba con un email específico
// SeedTestUserWithEmail crea un usuario de prueba con email personalizado
// Contraseña sin encriptar: Test1234!
// Role: student (por defecto)
func SeedTestUserWithEmail(t *testing.T, db *sql.DB, email string) (userID string, emailOut string) {
	t.Helper()

	password := "Test1234!" // Contraseña sin encriptar para uso en tests de login

	// Generar hash de contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Insertar usuario usando el ID generado dinámicamente
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

		-- Material Progress table (nombre correcto según repositorio)
		CREATE TABLE IF NOT EXISTS material_progress (
			material_id UUID NOT NULL REFERENCES materials(id),
			user_id UUID NOT NULL REFERENCES users(id),
			percentage INT DEFAULT 0 CHECK (percentage >= 0 AND percentage <= 100),
			last_page INT DEFAULT 0,
			status VARCHAR(50) DEFAULT 'in_progress',
			last_accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			completed_at TIMESTAMP NULL,
			PRIMARY KEY (material_id, user_id)
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

// TestUser representa un usuario de prueba con sus credenciales sin encriptar
type TestUser struct {
	ID       string
	Email    string
	Password string // Contraseña sin encriptar para uso en tests
	Role     string
}

// SeedTestUsers crea múltiples usuarios de prueba con un rol específico
// Contraseña sin encriptar para todos: Test1234!
// Emails: test1@edugo.com, test2@edugo.com, etc.
func SeedTestUsers(t *testing.T, db *sql.DB, count int, role string) []TestUser {
	t.Helper()

	users := make([]TestUser, count)
	password := "Test1234!"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	for i := 0; i < count; i++ {
		email := fmt.Sprintf("test%d@edugo.com", i+1)
		firstName := fmt.Sprintf("Test%d", i+1)

		var userID string
		query := `
			INSERT INTO users (email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, 'User', $4, true, NOW(), NOW())
			RETURNING id
		`
		err := db.QueryRow(query, email, string(hashedPassword), firstName, role).Scan(&userID)
		if err != nil {
			t.Fatalf("Failed to seed test user %d: %v", i+1, err)
		}

		users[i] = TestUser{
			ID:       userID,
			Email:    email,
			Password: password,
			Role:     role,
		}
	}

	t.Logf("👥 Created %d test users with role: %s (password: %s)", count, role, password)
	return users
}

// TestScenario representa un escenario completo de prueba con múltiples entidades
type TestScenario struct {
	Teacher     TestUser
	Students    []TestUser
	Materials   []string // IDs de materiales
	Assessments []string // IDs de assessments (mismo que material ID)
}

// SeedCompleteTestScenario crea un escenario completo de prueba
// Incluye: 1 teacher, N students, materiales y assessments
// Contraseña para todos los usuarios: Test1234!
func SeedCompleteTestScenario(t *testing.T, db *sql.DB, mongodb *mongo.Database, numStudents int) *TestScenario {
	t.Helper()

	// 1. Crear teacher
	teacherID, teacherEmail := SeedTestUserWithEmail(t, db, "teacher@edugo.com")
	// Actualizar rol a teacher
	_, err := db.Exec("UPDATE users SET role = 'teacher' WHERE id = $1", teacherID)
	if err != nil {
		t.Fatalf("Failed to update teacher role: %v", err)
	}

	teacher := TestUser{
		ID:       teacherID,
		Email:    teacherEmail,
		Password: "Test1234!",
		Role:     "teacher",
	}

	// 2. Crear students
	students := SeedTestUsers(t, db, numStudents, "student")

	// 3. Crear materiales (2 por defecto)
	material1 := SeedTestMaterialWithTitle(t, db, teacher.ID, "Introducción a Go")
	material2 := SeedTestMaterialWithTitle(t, db, teacher.ID, "Testing en Go")
	materials := []string{material1, material2}

	// 4. Crear assessments para cada material
	assessment1 := SeedTestAssessment(t, mongodb, material1)
	assessment2 := SeedTestAssessment(t, mongodb, material2)
	assessments := []string{assessment1, assessment2}

	scenario := &TestScenario{
		Teacher:     teacher,
		Students:    students,
		Materials:   materials,
		Assessments: assessments,
	}

	t.Logf("🎬 Complete test scenario created:")
	t.Logf("   Teacher: %s (%s)", teacher.Email, teacher.ID)
	t.Logf("   Students: %d", len(students))
	t.Logf("   Materials: %d", len(materials))
	t.Logf("   Assessments: %d", len(assessments))
	t.Logf("   Password for all users: Test1234!")

	return scenario
}

// SeedTestUserWithRole crea un usuario de prueba con un rol específico
// Contraseña sin encriptar: Test1234!
func SeedTestUserWithRole(t *testing.T, db *sql.DB, email, role string) (userID string, emailOut string) {
	t.Helper()

	password := "Test1234!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	firstName := "Test"
	lastName := "User"
	if role == "teacher" {
		firstName = "Teacher"
	} else if role == "admin" {
		firstName = "Admin"
	}

	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, true, NOW(), NOW())
		RETURNING id
	`
	err = db.QueryRow(query, email, string(hashedPassword), firstName, lastName, role).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to seed test user with role %s: %v", role, err)
	}

	t.Logf("👤 Test %s created: %s (email: %s, password: %s)", role, userID, email, password)
	return userID, email
}

// SetupTestAppWithSharedContainers inicializa una aplicación completa para testing
// usando contenedores compartidos para mejor performance
func SetupTestAppWithSharedContainers(t *testing.T) *TestApp {
	t.Helper()

	// Skip si tests están deshabilitados
	SkipIfIntegrationTestsDisabled(t)

	// Obtener manager de shared/testing (singleton)
	manager, err := containers.GetManager(t, nil)
	if err != nil {
		t.Fatalf("Failed to get containers manager: %v", err)
	}

	ctx := context.Background()

	// Obtener conexiones directamente del manager
	db := manager.PostgreSQL().DB()
	t.Log("✅ PostgreSQL connected")

	// Crear schema básico para tests (si no existe)
	if err := initTestSchema(db); err != nil {
		t.Fatalf("Failed to init test schema: %v", err)
	}
	t.Log("✅ PostgreSQL schema initialized")

	// MongoDB
	mongodb := manager.MongoDB().Database()
	t.Log("✅ MongoDB connected")

	// Logger mock para tests
	testLogger := &testLogger{}

	// RabbitMQ connection string
	rabbitConnStr, err := manager.RabbitMQ().ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get RabbitMQ connection string: %v", err)
	}

	// Crear RabbitMQ Publisher (opcional)
	publisher, err := createTestRabbitMQPublisher(rabbitConnStr, testLogger)
	if err != nil {
		t.Logf("⚠️  Warning: RabbitMQ publisher failed (non-critical): %v", err)
		publisher = &mockPublisher{}
	}

	// S3 Client mock
	s3Client := createTestS3Client()

	// JWT Secret para tests
	jwtSecret := "test-jwt-secret-key-very-secure-for-testing-only"

	// Crear Resources para Container DI
	resources := &bootstrap.Resources{
		Logger:            testLogger,
		PostgreSQL:        db,
		MongoDB:           mongodb,
		RabbitMQPublisher: publisher,
		S3Client:          s3Client,
		JWTSecret:         jwtSecret,
	}

	// Crear Container DI
	c := container.NewContainer(resources)
	t.Log("✅ Container DI initialized")

	// Cleanup: NO cerrar db ni mongodb (managed por shared/testing)
	// NO cerrar Container porque cierra la DB
	appCleanup := func() {
		t.Log("✅ Test app cleaned up")
	}

	return &TestApp{
		Container: c,
		DB:        db,
		MongoDB:   mongodb,
		Cleanup:   appCleanup,
	}
}
