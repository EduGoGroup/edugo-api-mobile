//go:build integration

package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMaterialTable(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS materials (
			id UUID PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			author_id UUID NOT NULL,
			subject_id VARCHAR(100),
			s3_key VARCHAR(500),
			s3_url TEXT,
			status VARCHAR(50) NOT NULL DEFAULT 'draft',
			processing_status VARCHAR(50) NOT NULL DEFAULT 'pending',
			is_deleted BOOLEAN DEFAULT false,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	require.NoError(t, err, "Failed to create materials table")
}

func TestMaterialRepository_Create(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	setupMaterialTable(t, db)

	repo := NewPostgresMaterialRepository(db)

	authorID := valueobject.NewUserID()
	material, err := entity.NewMaterial("Test Material", "Description", authorID, "subject-1")
	require.NoError(t, err)

	// Act
	err = repo.Create(context.Background(), material)

	// Assert
	require.NoError(t, err, "Create should not return error")

	// Verificar que se creó en DB
	var count int
	db.QueryRow("SELECT COUNT(*) FROM materials WHERE id = $1", material.ID().String()).Scan(&count)
	assert.Equal(t, 1, count, "Material should be in database")
}

func TestMaterialRepository_FindByID_MaterialExists(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	setupMaterialTable(t, db)

	repo := NewPostgresMaterialRepository(db)

	// Crear material directamente en DB
	materialID := valueobject.NewMaterialID()
	authorID := valueobject.NewUserID()

	_, err := db.Exec(`
		INSERT INTO materials (id, title, description, author_id, subject_id, status, processing_status)
		VALUES ($1, $2, $3, $4, $5, 'published', 'completed')
	`, materialID.String(), "Test Material", "Test Description", authorID.String(), "subject-1")
	require.NoError(t, err)

	// Act
	material, err := repo.FindByID(context.Background(), materialID)

	// Assert
	require.NoError(t, err, "FindByID should not return error when material exists")
	assert.NotNil(t, material)
	assert.Equal(t, "Test Material", material.Title())
	assert.Equal(t, "Test Description", material.Description())
}

func TestMaterialRepository_FindByID_MaterialNotFound(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	setupMaterialTable(t, db)

	repo := NewPostgresMaterialRepository(db)
	nonExistentID := valueobject.NewMaterialID()

	// Act
	material, err := repo.FindByID(context.Background(), nonExistentID)

	// Assert
	assert.NoError(t, err, "FindByID should not error on not found")
	assert.Nil(t, material, "Material should be nil when not found")
}

func TestMaterialRepository_FindByAuthorID(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	setupMaterialTable(t, db)

	repo := NewPostgresMaterialRepository(db)
	authorID := valueobject.NewUserID()

	// Crear 2 materiales del mismo autor
	for i := 1; i <= 2; i++ {
		materialID := valueobject.NewMaterialID()
		_, err := db.Exec(`
			INSERT INTO materials (id, title, description, author_id, status, processing_status)
			VALUES ($1, $2, $3, $4, 'published', 'completed')
		`, materialID.String(), "Material "+string(rune(i)), "Description", authorID.String())
		require.NoError(t, err)
	}

	// Act
	materials, err := repo.FindByAuthorID(context.Background(), authorID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, materials, 2, "Should find 2 materials")
}
