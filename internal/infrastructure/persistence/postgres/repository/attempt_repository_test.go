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

	domainErrors "github.com/EduGoGroup/edugo-api-mobile/internal/domain/errors"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// ptrStr crea un puntero a string
func ptrStr(s string) *string {
	return &s
}

// ptrBool crea un puntero a bool
func ptrBool(b bool) *bool {
	return &b
}

// ptrInt crea un puntero a int
func ptrInt(i int) *int {
	return &i
}

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
	now := time.Now().UTC()
	answer1 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionID:       "q1",
		StudentAnswer:    ptrStr("a"),
		IsCorrect:        ptrBool(true),
		TimeSpentSeconds: ptrInt(30),
		AnsweredAt:       now,
		CreatedAt:        now,
	}
	answer2 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionID:       "q2",
		StudentAnswer:    ptrStr("b"),
		IsCorrect:        ptrBool(true),
		TimeSpentSeconds: ptrInt(25),
		AnsweredAt:       now,
		CreatedAt:        now,
	}

	startedAt := now.Add(-5 * time.Minute)
	completedAt := now

	attempt := &pgentities.AssessmentAttempt{
		ID:               attemptID,
		AssessmentID:     uuid.New(),
		StudentID:        uuid.New(),
		Score:            100, // 2 correctas de 2 = (2*100)/2 = 100
		MaxScore:         100,
		TimeSpentSeconds: 55,
		StartedAt:        startedAt,
		CompletedAt:      completedAt,
		CreatedAt:        now,
		Answers:          []*pgentities.AssessmentAttemptAnswer{answer1, answer2},
		IdempotencyKey:   nil,
	}
	err := require.NoError(t, nil)

	// Act
	err = repo.Save(ctx, *attempt)

	// Assert
	assert.NoError(t, err)
}

func TestPostgresAttemptRepository_Save_NilAttempt(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAttemptRepository(mockDB)

	ctx := context.Background()

	// Act
	err := repo.Save(ctx, pgentities.AssessmentAttempt{})

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
	attempt := pgentities.AssessmentAttempt{
		ID:               uuid.New(),
		AssessmentID:     uuid.New(),
		StudentID:        uuid.New(),
		Score:            0,
		MaxScore:         100,
		TimeSpentSeconds: 60,
		StartedAt:        time.Now().UTC().Add(-1 * time.Minute),
		CompletedAt:      time.Now().UTC(),
		CreatedAt:        time.Now().UTC(),
		Answers:          []*pgentities.AssessmentAttemptAnswer{}, // Vacío!
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
