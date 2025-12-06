// TODO: Estos tests unitarios requieren actualización para usar mocks reales (sqlmock)
// Los tests de integración en assessment_repository_integration_test.go
// validan el funcionamiento real con testcontainers

//go:build skip_for_now
// +build skip_for_now

package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainErrors "github.com/EduGoGroup/edugo-api-mobile/internal/domain/errors"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

func TestNewPostgresAssessmentRepository(t *testing.T) {
	// Mock DB
	var mockDB *sql.DB

	repo := NewPostgresAssessmentRepository(mockDB)

	assert.NotNil(t, repo)
	assert.IsType(t, &PostgresAssessmentRepository{}, repo)
}

func TestPostgresAssessmentRepository_FindByID_Success(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAssessmentRepository(mockDB)

	ctx := context.Background()
	assessmentID := uuid.New()

	// Act
	result, err := repo.FindByID(ctx, assessmentID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, assessmentID, result.ID)
	assert.NotEmpty(t, result.Title)
	assert.Equal(t, 5, result.TotalQuestions)
	assert.Equal(t, 70, result.PassThreshold)
	assert.NotNil(t, result.MaxAttempts)
	assert.Equal(t, 3, *result.MaxAttempts)
}

func TestPostgresAssessmentRepository_FindByID_InvalidID(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAssessmentRepository(mockDB)

	ctx := context.Background()

	// Act
	result, err := repo.FindByID(ctx, uuid.Nil)

	// Assert
	assert.ErrorIs(t, err, domainErrors.ErrInvalidAssessmentID)
	assert.Nil(t, result)
}

func TestPostgresAssessmentRepository_FindByMaterialID_Success(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAssessmentRepository(mockDB)

	ctx := context.Background()
	materialID := uuid.New()

	// Act
	result, err := repo.FindByMaterialID(ctx, materialID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, materialID, result.MaterialID)
	assert.NotEmpty(t, result.Title)
	assert.Equal(t, 24, len(result.MongoDocumentID))
}

func TestPostgresAssessmentRepository_FindByMaterialID_InvalidID(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAssessmentRepository(mockDB)

	ctx := context.Background()

	// Act
	result, err := repo.FindByMaterialID(ctx, uuid.Nil)

	// Assert
	assert.ErrorIs(t, err, domainErrors.ErrInvalidMaterialID)
	assert.Nil(t, result)
}

func TestPostgresAssessmentRepository_Save_Success(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAssessmentRepository(mockDB)

	ctx := context.Background()

	maxAttempts := 5
	assessment := &pgentities.Assessment{
		ID:               uuid.New(),
		MaterialID:       uuid.New(),
		MongoDocumentID:  "507f1f77bcf86cd799439011",
		Title:            "Test Assessment",
		TotalQuestions:   10,
		PassThreshold:    75,
		MaxAttempts:      &maxAttempts,
		TimeLimitMinutes: nil,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	// Act
	err := repo.Save(ctx, assessment)

	// Assert
	assert.NoError(t, err)
}

func TestPostgresAssessmentRepository_Save_NilAssessment(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAssessmentRepository(mockDB)

	ctx := context.Background()

	// Act
	err := repo.Save(ctx, nil)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestPostgresAssessmentRepository_Save_InvalidAssessment(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAssessmentRepository(mockDB)

	ctx := context.Background()

	// Assessment inválido (ID nil)
	assessment := &pgentities.Assessment{
		ID:             uuid.Nil, // Inválido
		MaterialID:     uuid.New(),
		Title:          "Test",
		TotalQuestions: 5,
		PassThreshold:  70,
	}

	// Act
	err := repo.Save(ctx, assessment)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid assessment")
}

func TestPostgresAssessmentRepository_Delete_Success(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAssessmentRepository(mockDB)

	ctx := context.Background()
	assessmentID := uuid.New()

	// Act
	err := repo.Delete(ctx, assessmentID)

	// Assert
	assert.NoError(t, err)
}

func TestPostgresAssessmentRepository_Delete_InvalidID(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAssessmentRepository(mockDB)

	ctx := context.Background()

	// Act
	err := repo.Delete(ctx, uuid.Nil)

	// Assert
	assert.ErrorIs(t, err, domainErrors.ErrInvalidAssessmentID)
}
