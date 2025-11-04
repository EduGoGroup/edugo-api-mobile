package database

import (
	"context"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestInitPostgreSQL_Success verifica que la inicialización de PostgreSQL funcione correctamente
func TestInitPostgreSQL_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Levantar contenedor de PostgreSQL para testing
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err)
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("Error terminando contenedor: %v", err)
		}
	}()

	// Obtener información de conexión del contenedor
	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)

	mappedPort, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Configurar config de prueba con datos reales del contenedor
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{
				Host:           host,
				Port:           mappedPort.Int(),
				Database:       "testdb",
				User:           "testuser",
				Password:       "testpass",
				MaxConnections: 10,
				SSLMode:        "disable",
			},
		},
	}

	// Crear logger de prueba
	testLogger := logger.NewZapLogger("debug", "json")

	// Ejecutar la función a testear
	db, err := InitPostgreSQL(ctx, cfg, testLogger)

	// Verificaciones
	require.NoError(t, err, "No debería haber error al inicializar PostgreSQL")
	require.NotNil(t, db, "La conexión DB no debería ser nil")

	// Verificar que la conexión está activa
	err = db.PingContext(ctx)
	assert.NoError(t, err, "Debería poder hacer ping a la base de datos")

	// Verificar configuración del pool de conexiones
	stats := db.Stats()
	assert.Equal(t, 10, stats.MaxOpenConnections, "MaxOpenConnections debería estar configurado")

	// Limpiar
	err = db.Close()
	assert.NoError(t, err, "No debería haber error al cerrar la conexión")
}

// TestInitPostgreSQL_WithDefaultMaxConnections verifica que funcione sin MaxConnections configurado
func TestInitPostgreSQL_WithDefaultMaxConnections(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Levantar contenedor de PostgreSQL
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err)
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("Error terminando contenedor: %v", err)
		}
	}()

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)

	mappedPort, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Config sin MaxConnections (0 = default)
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{
				Host:           host,
				Port:           mappedPort.Int(),
				Database:       "testdb",
				User:           "testuser",
				Password:       "testpass",
				MaxConnections: 0, // Sin límite explícito
				SSLMode:        "disable",
			},
		},
	}

	testLogger := logger.NewZapLogger("debug", "json")

	db, err := InitPostgreSQL(ctx, cfg, testLogger)

	require.NoError(t, err)
	require.NotNil(t, db)

	// Verificar conexión
	err = db.PingContext(ctx)
	assert.NoError(t, err)

	// Verificar que usa valores por defecto de Go (sin límite)
	stats := db.Stats()
	assert.Equal(t, 0, stats.MaxOpenConnections, "Debería usar valor por defecto sin límite")

	err = db.Close()
	assert.NoError(t, err)
}

// TestInitPostgreSQL_InvalidConnectionString verifica manejo de errores con conexión inválida
func TestInitPostgreSQL_InvalidConnectionString(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Config con configuración inválida
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{
				Host:     "invalid-host-that-does-not-exist",
				Port:     9999,
				Database: "testdb",
				User:     "testuser",
				Password: "testpass",
				SSLMode:  "disable",
			},
		},
	}

	testLogger := logger.NewZapLogger("debug", "json")

	// Crear contexto con timeout corto para no esperar mucho
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	db, err := InitPostgreSQL(ctxTimeout, cfg, testLogger)

	// Debería fallar al hacer ping
	assert.Error(t, err, "Debería haber error con host inválido")
	assert.Contains(t, err.Error(), "error pinging postgres", "El error debería mencionar el ping")

	if db != nil {
		db.Close()
	}
}

// TestInitPostgreSQL_ContextCancellation verifica que respete la cancelación del contexto
func TestInitPostgreSQL_ContextCancellation(t *testing.T) {
	t.Parallel()

	// Crear contexto que ya está cancelado
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancelar inmediatamente

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{
				Host:     "localhost",
				Port:     5432,
				Database: "testdb",
				User:     "testuser",
				Password: "testpass",
				SSLMode:  "disable",
			},
		},
	}

	testLogger := logger.NewZapLogger("debug", "json")

	db, err := InitPostgreSQL(ctx, cfg, testLogger)

	// Debería fallar porque el contexto está cancelado
	assert.Error(t, err, "Debería haber error con contexto cancelado")

	if db != nil {
		db.Close()
	}
}

// TestInitPostgreSQL_MultipleConnections verifica que el pool de conexiones funcione correctamente
func TestInitPostgreSQL_MultipleConnections(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Levantar contenedor de PostgreSQL
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err)
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("Error terminando contenedor: %v", err)
		}
	}()

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)

	mappedPort, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Postgres: config.PostgresConfig{
				Host:           host,
				Port:           mappedPort.Int(),
				Database:       "testdb",
				User:           "testuser",
				Password:       "testpass",
				MaxConnections: 5,
				SSLMode:        "disable",
			},
		},
	}

	testLogger := logger.NewZapLogger("debug", "json")

	db, err := InitPostgreSQL(ctx, cfg, testLogger)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	// Ejecutar varias queries concurrentes para verificar el pool
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			var result int
			err := db.QueryRowContext(ctx, "SELECT $1", id).Scan(&result)
			assert.NoError(t, err)
			assert.Equal(t, id, result)
			done <- true
		}(i)
	}

	// Esperar a que todas las goroutines terminen
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verificar estadísticas del pool
	stats := db.Stats()
	assert.Equal(t, 5, stats.MaxOpenConnections, "MaxOpenConnections debería ser 5")
	assert.LessOrEqual(t, stats.OpenConnections, 5, "No debería haber más de 5 conexiones abiertas")
}
