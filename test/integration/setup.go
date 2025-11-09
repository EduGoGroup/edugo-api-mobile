//go:build integration

package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/streadway/amqp"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

// TestContainers contiene todos los contenedores de prueba
type TestContainers struct {
	Postgres *postgres.PostgresContainer
	MongoDB  *mongodb.MongoDBContainer
	RabbitMQ *rabbitmq.RabbitMQContainer
}

// SetupContainers inicia todos los contenedores necesarios para testing
func SetupContainers(t *testing.T) (*TestContainers, func()) {
	ctx := context.Background()

	// PostgreSQL (schema se crea manualmente después)
	t.Log("🐘 Iniciando PostgreSQL testcontainer...")
	pgContainer, err := postgres.Run(ctx, "postgres:15-alpine",
		postgres.WithDatabase("edugo"),
		postgres.WithUsername("edugo_user"),
		postgres.WithPassword("edugo_pass"),
		// NO usar WithInitScripts porque la ruta es relativa y falla en tests
		// El schema se ejecutará manualmente en SetupTestApp
	)
	if err != nil {
		t.Fatalf("Failed to start Postgres: %v", err)
	}
	t.Log("✅ PostgreSQL ready")

	// MongoDB
	t.Log("🍃 Iniciando MongoDB testcontainer...")
	mongoContainer, err := mongodb.Run(ctx, "mongo:7.0",
		mongodb.WithUsername("edugo_admin"),
		mongodb.WithPassword("edugo_pass"),
	)
	if err != nil {
		pgContainer.Terminate(ctx)
		t.Fatalf("Failed to start MongoDB: %v", err)
	}
	t.Log("✅ MongoDB ready")

	// RabbitMQ
	t.Log("🐰 Iniciando RabbitMQ testcontainer...")
	rabbitContainer, err := rabbitmq.Run(ctx, "rabbitmq:3.12-management-alpine",
		rabbitmq.WithAdminUsername("edugo_user"),
		rabbitmq.WithAdminPassword("edugo_pass"),
	)
	if err != nil {
		pgContainer.Terminate(ctx)
		mongoContainer.Terminate(ctx)
		t.Fatalf("Failed to start RabbitMQ: %v", err)
	}
	t.Log("✅ RabbitMQ ready")

	// Configurar topología de RabbitMQ (exchanges, colas, bindings)
	if err := setupRabbitMQTopology(ctx, rabbitContainer); err != nil {
		t.Logf("⚠️  Warning: RabbitMQ topology setup failed (non-critical): %v", err)
		// No fallar el test, algunos tests pueden funcionar sin RabbitMQ
	} else {
		t.Log("✅ RabbitMQ topology configured")
	}

	containers := &TestContainers{
		Postgres: pgContainer,
		MongoDB:  mongoContainer,
		RabbitMQ: rabbitContainer,
	}

	// Cleanup function
	cleanup := func() {
		t.Log("🧹 Cleaning up testcontainers...")
		pgContainer.Terminate(ctx)
		mongoContainer.Terminate(ctx)
		rabbitContainer.Terminate(ctx)
		t.Log("✅ Testcontainers terminated")
	}

	return containers, cleanup
}

// setupRabbitMQTopology configura la topología de RabbitMQ para tests
// Crea exchanges, colas y bindings necesarios para el proyecto
func setupRabbitMQTopology(ctx context.Context, container *rabbitmq.RabbitMQContainer) error {
	// Obtener connection string de RabbitMQ
	connStr, err := container.AmqpURL(ctx)
	if err != nil {
		return fmt.Errorf("failed to get RabbitMQ connection string: %w", err)
	}

	// Conectar a RabbitMQ
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	// Crear canal
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()

	// 1. Crear exchange principal
	exchangeName := "edugo.events"
	err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange %s: %w", exchangeName, err)
	}

	// 2. Definir colas necesarias
	queues := []struct {
		name       string
		routingKey string
	}{
		{"material.created", "material.created"},
		{"material.updated", "material.updated"},
		{"material.deleted", "material.deleted"},
		{"assessment.completed", "assessment.completed"},
		{"progress.updated", "progress.updated"},
		{"user.registered", "user.registered"},
	}

	// 3. Crear colas y bindings
	for _, q := range queues {
		// Declarar cola
		_, err = ch.QueueDeclare(
			q.name, // name
			true,   // durable
			false,  // delete when unused
			false,  // exclusive
			false,  // no-wait
			nil,    // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", q.name, err)
		}

		// Crear binding entre exchange y cola
		err = ch.QueueBind(
			q.name,       // queue name
			q.routingKey, // routing key
			exchangeName, // exchange
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue %s: %w", q.name, err)
		}
	}

	return nil
}

// SetupPostgres inicia solo PostgreSQL
func SetupPostgres(t *testing.T) (*postgres.PostgresContainer, func()) {
	ctx := context.Background()

	t.Log("🐘 Iniciando PostgreSQL testcontainer...")
	pgContainer, err := postgres.Run(ctx, "postgres:15-alpine",
		postgres.WithDatabase("edugo"),
		postgres.WithUsername("edugo_user"),
		postgres.WithPassword("edugo_pass"),
		// Schema se ejecuta manualmente después de conectar
	)
	if err != nil {
		t.Fatalf("Failed to start Postgres: %v", err)
	}
	t.Log("✅ PostgreSQL ready")

	cleanup := func() {
		pgContainer.Terminate(ctx)
	}

	return pgContainer, cleanup
}

// SetupMongoDB inicia solo MongoDB
func SetupMongoDB(t *testing.T) (*mongodb.MongoDBContainer, func()) {
	ctx := context.Background()

	t.Log("🍃 Iniciando MongoDB testcontainer...")
	mongoContainer, err := mongodb.Run(ctx, "mongo:7.0",
		mongodb.WithUsername("edugo_admin"),
		mongodb.WithPassword("edugo_pass"),
	)
	if err != nil {
		t.Fatalf("Failed to start MongoDB: %v", err)
	}
	t.Log("✅ MongoDB ready")

	cleanup := func() {
		mongoContainer.Terminate(ctx)
	}

	return mongoContainer, cleanup
}

// SetupRabbitMQ inicia solo RabbitMQ
func SetupRabbitMQ(t *testing.T) (*rabbitmq.RabbitMQContainer, func()) {
	ctx := context.Background()

	t.Log("🐰 Iniciando RabbitMQ testcontainer...")
	rabbitContainer, err := rabbitmq.Run(ctx, "rabbitmq:3.12-management-alpine",
		rabbitmq.WithAdminUsername("edugo_user"),
		rabbitmq.WithAdminPassword("edugo_pass"),
	)
	if err != nil {
		t.Fatalf("Failed to start RabbitMQ: %v", err)
	}
	t.Log("✅ RabbitMQ ready")

	cleanup := func() {
		rabbitContainer.Terminate(ctx)
	}

	return rabbitContainer, cleanup
}
