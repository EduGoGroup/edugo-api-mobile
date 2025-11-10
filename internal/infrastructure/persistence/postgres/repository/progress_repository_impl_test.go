//go:build integration

package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// setupProgressTestDB crea un testcontainer de PostgreSQL con todas las tablas necesarias
func setupProgressTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	ctx := context.Background()

	// Levantar PostgreSQL testcontainer
	pgContainer, err := postgres.Run(ctx, "postgres:15-alpine",
		postgres.WithDatabase("edugo"),
		postgres.WithUsername("edugo_user"),
		postgres.WithPassword("edugo_pass"),
	)
	require.NoError(t, err, "Failed to start PostgreSQL testcontainer")

	// Obtener connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Conectar a la base de datos
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	// Retry ping to ensure database is ready
	var pingErr error
	for i := 0; i < 10; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	require.NoError(t, pingErr, "Failed to ping database after retries")

	// Crear todas las tablas necesarias
	_, err = db.Exec(`
		-- Users table
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			role VARCHAR(50) NOT NULL,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		-- Materials table
		CREATE TABLE IF NOT EXISTS materials (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title VARCHAR(255) NOT NULL,
			description TEXT,
			author_id UUID NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'draft',
			processing_status VARCHAR(50) NOT NULL DEFAULT 'pending',
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		-- Material Progress table
		CREATE TABLE IF NOT EXISTS material_progress (
			material_id UUID NOT NULL REFERENCES materials(id),
			user_id UUID NOT NULL REFERENCES users(id),
			percentage INT DEFAULT 0 CHECK (percentage >= 0 AND percentage <= 100),
			last_page INT DEFAULT 0,
			status VARCHAR(50) DEFAULT 'not_started',
			last_accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			completed_at TIMESTAMP NULL,
			PRIMARY KEY (material_id, user_id)
		);
	`)
	require.NoError(t, err, "Failed to create tables")

	cleanup := func() {
		db.Close()
		pgContainer.Terminate(ctx)
	}

	return db, cleanup
}

// seedTestUserAndMaterial crea un usuario y material de prueba
func seedTestUserAndMaterial(t *testing.T, db *sql.DB) (valueobject.UserID, valueobject.MaterialID) {
	t.Helper()

	ctx := context.Background()

	// Crear usuario
	var userIDStr string
	err := db.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, 'Test', 'User', 'student', true)
		RETURNING id
	`, "test@example.com", "hashedpassword").Scan(&userIDStr)
	require.NoError(t, err, "Failed to create test user")

	userID, err := valueobject.UserIDFromString(userIDStr)
	require.NoError(t, err, "Failed to parse user ID")

	// Crear material
	var materialIDStr string
	err = db.QueryRowContext(ctx, `
		INSERT INTO materials (title, description, author_id, status, processing_status)
		VALUES ($1, $2, $3, 'published', 'completed')
		RETURNING id
	`, "Test Material", "Test Description", userIDStr).Scan(&materialIDStr)
	require.NoError(t, err, "Failed to create test material")

	materialID, err := valueobject.MaterialIDFromString(materialIDStr)
	require.NoError(t, err, "Failed to parse material ID")

	return userID, materialID
}

func TestProgressRepository_Upsert_CreateNewProgress(t *testing.T) {
	// Arrange
	db, cleanup := setupProgressTestDB(t)
	defer cleanup()

	userID, materialID := seedTestUserAndMaterial(t, db)
	repo := NewPostgresProgressRepository(db)

	// Crear nuevo progreso
	progress := entity.NewProgress(materialID, userID)
	err := progress.UpdateProgress(25, 5)
	require.NoError(t, err)

	// Act
	result, err := repo.Upsert(context.Background(), progress)

	// Assert
	require.NoError(t, err, "Upsert should not return error when creating new progress")
	assert.NotNil(t, result)
	assert.Equal(t, materialID.String(), result.MaterialID().String())
	assert.Equal(t, userID.String(), result.UserID().String())
	assert.Equal(t, 25, result.Percentage())
	assert.Equal(t, 5, result.LastPage())
	assert.Equal(t, enum.ProgressStatusInProgress, result.Status())

	// Verificar que se creó en DB
	var count int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM material_progress 
		WHERE material_id = $1 AND user_id = $2
	`, materialID.String(), userID.String()).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "Progress should be in database")
}

func TestProgressRepository_Upsert_UpdateExistingProgress(t *testing.T) {
	// Arrange
	db, cleanup := setupProgressTestDB(t)
	defer cleanup()

	userID, materialID := seedTestUserAndMaterial(t, db)
	repo := NewPostgresProgressRepository(db)

	// Crear progreso inicial
	initialProgress := entity.NewProgress(materialID, userID)
	err := initialProgress.UpdateProgress(25, 5)
	require.NoError(t, err)

	_, err = repo.Upsert(context.Background(), initialProgress)
	require.NoError(t, err)

	// Esperar un momento para asegurar que updated_at sea diferente
	time.Sleep(10 * time.Millisecond)

	// Actualizar progreso
	updatedProgress := entity.ReconstructProgress(
		materialID,
		userID,
		50,
		10,
		enum.ProgressStatusInProgress,
		time.Now(),
		initialProgress.CreatedAt(),
		time.Now(),
	)

	// Act
	result, err := repo.Upsert(context.Background(), updatedProgress)

	// Assert
	require.NoError(t, err, "Upsert should not return error when updating existing progress")
	assert.NotNil(t, result)
	assert.Equal(t, 50, result.Percentage())
	assert.Equal(t, 10, result.LastPage())
	assert.Equal(t, enum.ProgressStatusInProgress, result.Status())

	// Verificar que solo hay un registro en DB
	var count int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM material_progress 
		WHERE material_id = $1 AND user_id = $2
	`, materialID.String(), userID.String()).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "Should still have only one progress record")

	// Verificar que el porcentaje se actualizó
	var percentage int
	err = db.QueryRow(`
		SELECT percentage FROM material_progress 
		WHERE material_id = $1 AND user_id = $2
	`, materialID.String(), userID.String()).Scan(&percentage)
	require.NoError(t, err)
	assert.Equal(t, 50, percentage)
}

func TestProgressRepository_Upsert_CompleteProgress(t *testing.T) {
	// Arrange
	db, cleanup := setupProgressTestDB(t)
	defer cleanup()

	userID, materialID := seedTestUserAndMaterial(t, db)
	repo := NewPostgresProgressRepository(db)

	// Crear progreso completado
	progress := entity.NewProgress(materialID, userID)
	err := progress.UpdateProgress(100, 20)
	require.NoError(t, err)

	// Act
	result, err := repo.Upsert(context.Background(), progress)

	// Assert
	require.NoError(t, err, "Upsert should not return error when completing progress")
	assert.NotNil(t, result)
	assert.Equal(t, 100, result.Percentage())
	assert.Equal(t, enum.ProgressStatusCompleted, result.Status())

	// Verificar que completed_at se estableció
	var completedAt sql.NullTime
	err = db.QueryRow(`
		SELECT completed_at FROM material_progress 
		WHERE material_id = $1 AND user_id = $2
	`, materialID.String(), userID.String()).Scan(&completedAt)
	require.NoError(t, err)
	assert.True(t, completedAt.Valid, "completed_at should be set when percentage is 100")
}

func TestProgressRepository_FindByMaterialAndUser_ProgressExists(t *testing.T) {
	// Arrange
	db, cleanup := setupProgressTestDB(t)
	defer cleanup()

	userID, materialID := seedTestUserAndMaterial(t, db)
	repo := NewPostgresProgressRepository(db)

	// Crear progreso directamente en DB
	now := time.Now()
	_, err := db.Exec(`
		INSERT INTO material_progress (material_id, user_id, percentage, last_page, status, last_accessed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, materialID.String(), userID.String(), 75, 15, "in_progress", now, now, now)
	require.NoError(t, err)

	// Act
	progress, err := repo.FindByMaterialAndUser(context.Background(), materialID, userID)

	// Assert
	require.NoError(t, err, "FindByMaterialAndUser should not return error when progress exists")
	assert.NotNil(t, progress)
	assert.Equal(t, materialID.String(), progress.MaterialID().String())
	assert.Equal(t, userID.String(), progress.UserID().String())
	assert.Equal(t, 75, progress.Percentage())
	assert.Equal(t, 15, progress.LastPage())
	assert.Equal(t, enum.ProgressStatusInProgress, progress.Status())
}

func TestProgressRepository_FindByMaterialAndUser_ProgressNotFound(t *testing.T) {
	// Arrange
	db, cleanup := setupProgressTestDB(t)
	defer cleanup()

	userID, materialID := seedTestUserAndMaterial(t, db)
	repo := NewPostgresProgressRepository(db)

	// No crear ningún progreso

	// Act
	progress, err := repo.FindByMaterialAndUser(context.Background(), materialID, userID)

	// Assert
	require.NoError(t, err, "FindByMaterialAndUser should not return error when progress not found")
	assert.Nil(t, progress, "Progress should be nil when not found")
}

func TestProgressRepository_FindByMaterialAndUser_DifferentUser(t *testing.T) {
	// Arrange
	db, cleanup := setupProgressTestDB(t)
	defer cleanup()

	userID1, materialID := seedTestUserAndMaterial(t, db)
	repo := NewPostgresProgressRepository(db)

	// Crear segundo usuario
	var userID2Str string
	err := db.QueryRow(`
		INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
		VALUES ($1, $2, 'Test2', 'User2', 'student', true)
		RETURNING id
	`, "test2@example.com", "hashedpassword").Scan(&userID2Str)
	require.NoError(t, err)

	userID2, err := valueobject.UserIDFromString(userID2Str)
	require.NoError(t, err)

	// Crear progreso para userID1
	now := time.Now()
	_, err = db.Exec(`
		INSERT INTO material_progress (material_id, user_id, percentage, last_page, status, last_accessed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, materialID.String(), userID1.String(), 50, 10, "in_progress", now, now, now)
	require.NoError(t, err)

	// Act - Buscar progreso para userID2 (no existe)
	progress, err := repo.FindByMaterialAndUser(context.Background(), materialID, userID2)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, progress, "Should not find progress for different user")
}
