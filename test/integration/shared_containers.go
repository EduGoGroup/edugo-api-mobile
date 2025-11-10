//go:build integration

package integration

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SharedContainers mantiene contenedores compartidos entre tests
type SharedContainers struct {
	Postgres *postgres.PostgresContainer
	MongoDB  *mongodb.MongoDBContainer
	RabbitMQ *rabbitmq.RabbitMQContainer
	mu       sync.Mutex
}

var (
	sharedContainers *SharedContainers
	setupOnce        sync.Once
	setupError       error
)

// GetSharedContainers obtiene o crea los contenedores compartidos
// Los contenedores se crean UNA SOLA VEZ y se reutilizan entre todos los tests
func GetSharedContainers(t *testing.T) (*SharedContainers, error) {
	setupOnce.Do(func() {
		ctx := context.Background()

		t.Log("🚀 Iniciando contenedores compartidos (UNA SOLA VEZ)...")

		// PostgreSQL
		t.Log("🐘 Creando PostgreSQL compartido...")
		pgContainer, err := postgres.Run(ctx, "postgres:15-alpine",
			postgres.WithDatabase("edugo_test"),
			postgres.WithUsername("edugo_user"),
			postgres.WithPassword("edugo_pass"),
		)
		if err != nil {
			setupError = fmt.Errorf("failed to start shared PostgreSQL: %w", err)
			return
		}
		t.Log("✅ PostgreSQL compartido listo")

		// MongoDB
		t.Log("🍃 Creando MongoDB compartido...")
		mongoContainer, err := mongodb.Run(ctx, "mongo:7.0")
		if err != nil {
			pgContainer.Terminate(ctx)
			setupError = fmt.Errorf("failed to start shared MongoDB: %w", err)
			return
		}
		t.Log("✅ MongoDB compartido listo")

		// RabbitMQ
		t.Log("🐰 Creando RabbitMQ compartido...")
		rabbitContainer, err := rabbitmq.Run(ctx, "rabbitmq:3.12-alpine",
			rabbitmq.WithAdminUsername("edugo_user"),
			rabbitmq.WithAdminPassword("edugo_pass"),
		)
		if err != nil {
			pgContainer.Terminate(ctx)
			mongoContainer.Terminate(ctx)
			setupError = fmt.Errorf("failed to start shared RabbitMQ: %w", err)
			return
		}
		t.Log("✅ RabbitMQ compartido listo")

		sharedContainers = &SharedContainers{
			Postgres: pgContainer,
			MongoDB:  mongoContainer,
			RabbitMQ: rabbitContainer,
		}

		t.Log("🎉 Todos los contenedores compartidos están listos")
	})

	return sharedContainers, setupError
}

// CleanSharedDatabases limpia los datos de las bases de datos compartidas
// entre tests (TRUNCATE, no DROP) para permitir reutilización
func CleanSharedDatabases(t *testing.T, containers *SharedContainers) error {
	ctx := context.Background()

	// Limpiar PostgreSQL
	connStr, err := containers.Postgres.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to get postgres connection string: %w", err)
	}

	db, err := ConnectPostgresWithRetry(connStr, 3)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer db.Close()

	// Truncar todas las tablas (en orden correcto para respetar foreign keys)
	tables := []string{
		"refresh_tokens",
		"login_attempts",
		"student_material_progress",
		"quiz_attempt",
		"guardian_student_relation",
		"material_versions",
		"materials",
		"users",
		"subjects",
		"units",
		"schools",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			// Ignorar errores de tablas que no existen
			t.Logf("⚠️  Warning: Failed to truncate %s: %v", table, err)
		}
	}

	// Limpiar MongoDB
	mongoConnStr, err := containers.MongoDB.ConnectionString(ctx)
	if err != nil {
		return fmt.Errorf("failed to get mongodb connection string: %w", err)
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnStr))
	if err != nil {
		return fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	defer mongoClient.Disconnect(ctx)

	// Limpiar colecciones de MongoDB
	mongoDB := mongoClient.Database("edugo_test")
	collections := []string{
		"assessment_results",
		"material_summaries",
	}

	for _, collection := range collections {
		err := mongoDB.Collection(collection).Drop(ctx)
		if err != nil {
			t.Logf("⚠️  Warning: Failed to drop collection %s: %v", collection, err)
		}
	}

	t.Log("🧹 Bases de datos compartidas limpiadas")
	return nil
}

// TerminateSharedContainers destruye los contenedores compartidos
// Debe llamarse al final de TODOS los tests (en TestMain)
func TerminateSharedContainers() error {
	if sharedContainers == nil {
		return nil
	}

	ctx := context.Background()

	var errors []error

	if sharedContainers.Postgres != nil {
		if err := sharedContainers.Postgres.Terminate(ctx); err != nil {
			errors = append(errors, fmt.Errorf("failed to terminate postgres: %w", err))
		}
	}

	if sharedContainers.MongoDB != nil {
		if err := sharedContainers.MongoDB.Terminate(ctx); err != nil {
			errors = append(errors, fmt.Errorf("failed to terminate mongodb: %w", err))
		}
	}

	if sharedContainers.RabbitMQ != nil {
		if err := sharedContainers.RabbitMQ.Terminate(ctx); err != nil {
			errors = append(errors, fmt.Errorf("failed to terminate rabbitmq: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors terminating containers: %v", errors)
	}

	return nil
}
