//go:build integration

package integration

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPostgresTablesExist(t *testing.T) {
	// Setup PostgreSQL testcontainer
	pgContainer, cleanup := SetupPostgres(t)
	defer cleanup()

	ctx := context.Background()

	// Obtener connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	// Conectar
	db, err := sql.Open("postgres", connStr)
	assert.NoError(t, err)
	defer db.Close()

	// Verificar conexión
	err = db.Ping()
	assert.NoError(t, err)
	t.Log("✅ Conectado a PostgreSQL")

	// Verificar que tablas principales existan
	tables := []string{
		"user",
		"school",
		"unit",
		"subject",
		"material",
		"student_material_progress",
		"quiz_attempt",
		"guardian_student_relation",
	}

	for _, table := range tables {
		var exists bool
		query := `SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)`
		err := db.QueryRow(query, table).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists, "Table %s should exist", table)
		t.Logf("✓ Table '%s' exists", table)
	}

	// Verificar que hay 17 tablas en total
	var tableCount int
	err = db.QueryRow(`
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
	`).Scan(&tableCount)
	assert.NoError(t, err)
	assert.Equal(t, 17, tableCount, "Should have 17 tables")
	t.Log("✅ All 17 tables exist")
}
