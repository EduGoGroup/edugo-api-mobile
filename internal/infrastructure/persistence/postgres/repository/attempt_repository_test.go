// TODO: Estos tests unitarios requieren actualización para usar mocks reales (sqlmock)
// Los tests de integración en attempt_repository_integration_test.go
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

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entities"
	domainErrors "github.com/EduGoGroup/edugo-api-mobile/internal/domain/errors"
)

func TestNewPostgresAttemptRepository(t *testing.T) {
	var mockDB *sql.DB
	repo := NewPostgresAttemptRepository(mockDB)

	assert.NotNil(t, repo)
	assert.IsType(t, &PostgresAttemptRepository{}, repo)
}

func TestPostgresAttemptRepository_FindByID_Success(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAttemptRepository(mockDB)

	ctx := context.Background()
	attemptID := uuid.New()

	// Act
	result, err := repo.FindByID(ctx, attemptID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, attemptID, result.ID)
	assert.NotEmpty(t, result.Answers)
	assert.Greater(t, result.Score, 0)
	assert.Equal(t, 100, result.MaxScore)
}

func TestPostgresAttemptRepository_FindByID_InvalidID(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAttemptRepository(mockDB)

	ctx := context.Background()

	// Act
	result, err := repo.FindByID(ctx, uuid.Nil)

	// Assert
	assert.ErrorIs(t, err, domainErrors.ErrInvalidAttemptID)
	assert.Nil(t, result)
}

func TestPostgresAttemptRepository_FindByStudentAndAssessment_Success(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAttemptRepository(mockDB)

	ctx := context.Background()
	studentID := uuid.New()
	assessmentID := uuid.New()

	// Act
	results, err := repo.FindByStudentAndAssessment(ctx, studentID, assessmentID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, results)
	assert.NotEmpty(t, results)
	assert.Equal(t, assessmentID, results[0].AssessmentID)
	assert.Equal(t, studentID, results[0].StudentID)
}

func TestPostgresAttemptRepository_FindByStudentAndAssessment_InvalidIDs(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAttemptRepository(mockDB)

	ctx := context.Background()

	// Act
	results, err := repo.FindByStudentAndAssessment(ctx, uuid.Nil, uuid.Nil)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestPostgresAttemptRepository_Save_Success(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAttemptRepository(mockDB)

	ctx := context.Background()

	// Crear attempt válido con answers
	attemptID := uuid.New()
	answer1, _ := entities.NewAnswer(attemptID, "q1", "a", true, 30)
	answer2, _ := entities.NewAnswer(attemptID, "q2", "b", true, 25)

	startedAt := time.Now().UTC().Add(-5 * time.Minute)
	completedAt := time.Now().UTC()

	attempt, err := entities.NewAttempt(
		uuid.New(), // assessmentID
		uuid.New(), // studentID
		[]*entities.Answer{answer1, answer2},
		startedAt,
		completedAt,
	)
	require.NoError(t, err)

	// Act
	err = repo.Save(ctx, attempt)

	// Assert
	assert.NoError(t, err)
}

func TestPostgresAttemptRepository_Save_NilAttempt(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAttemptRepository(mockDB)

	ctx := context.Background()

	// Act
	err := repo.Save(ctx, nil)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestPostgresAttemptRepository_Save_NoAnswers(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAttemptRepository(mockDB)

	ctx := context.Background()

	// Attempt sin respuestas (inválido)
	attempt := &entities.Attempt{
		ID:               uuid.New(),
		AssessmentID:     uuid.New(),
		StudentID:        uuid.New(),
		Score:            0,
		MaxScore:         100,
		TimeSpentSeconds: 60,
		StartedAt:        time.Now().UTC().Add(-1 * time.Minute),
		CompletedAt:      time.Now().UTC(),
		CreatedAt:        time.Now().UTC(),
		Answers:          []*entities.Answer{}, // Vacío!
	}

	// Act
	err := repo.Save(ctx, attempt)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, domainErrors.ErrNoAnswersProvided)
}

func TestPostgresAttemptRepository_CountByStudentAndAssessment_Success(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAttemptRepository(mockDB)

	ctx := context.Background()
	studentID := uuid.New()
	assessmentID := uuid.New()

	// Act
	count, err := repo.CountByStudentAndAssessment(ctx, studentID, assessmentID)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 0)
}

func TestPostgresAttemptRepository_CountByStudentAndAssessment_InvalidIDs(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAttemptRepository(mockDB)

	ctx := context.Background()

	// Act
	count, err := repo.CountByStudentAndAssessment(ctx, uuid.Nil, uuid.Nil)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestPostgresAttemptRepository_FindByStudent_Success(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAttemptRepository(mockDB)

	ctx := context.Background()
	studentID := uuid.New()

	// Act
	results, err := repo.FindByStudent(ctx, studentID, 10, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, results)
	// Puede estar vacío o con resultados mock
}

func TestPostgresAttemptRepository_FindByStudent_Pagination(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAttemptRepository(mockDB)

	ctx := context.Background()
	studentID := uuid.New()

	// Act - Primera página
	page1, err := repo.FindByStudent(ctx, studentID, 5, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, page1)

	// Act - Segunda página
	page2, err := repo.FindByStudent(ctx, studentID, 5, 5)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, page2)
}

func TestPostgresAttemptRepository_FindByStudent_InvalidID(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAttemptRepository(mockDB)

	ctx := context.Background()

	// Act
	results, err := repo.FindByStudent(ctx, uuid.Nil, 10, 0)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, results)
}

// TODO: Claude Local - Tests de integración con testcontainers
// func TestPostgresAttemptRepository_Integration_WithTransaction(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping integration test")
// 	}
//
// 	// Test que verifica que:
// 	// 1. Se guarda attempt + answers en transacción
// 	// 2. Si falla una answer, se hace rollback del attempt
// 	// 3. Se pueden recuperar attempts con sus answers
// }
