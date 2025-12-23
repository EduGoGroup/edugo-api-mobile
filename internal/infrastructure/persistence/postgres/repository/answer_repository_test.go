// TODO: Estos tests unitarios requieren actualización para usar mocks reales (sqlmock)
// Los tests de integración en answer_repository_integration_test.go
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

// ptr es una función auxiliar para crear punteros
func ptr[T any](v T) *T {
	return &v
}

func TestNewPostgresAnswerRepository(t *testing.T) {
	var mockDB *sql.DB
	repo := NewPostgresAnswerRepository(mockDB)

	assert.NotNil(t, repo)
	assert.IsType(t, &PostgresAnswerRepository{}, repo)
}

func TestPostgresAnswerRepository_FindByAttemptID_Success(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAnswerRepository(mockDB)

	ctx := context.Background()
	attemptID := uuid.New()

	// Act
	results, err := repo.FindByAttemptID(ctx, attemptID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, results)
	assert.NotEmpty(t, results)

	// Verificar que todas las respuestas pertenecen al mismo intento
	for _, answer := range results {
		assert.Equal(t, attemptID, answer.AttemptID)
	}
}

func TestPostgresAnswerRepository_FindByAttemptID_InvalidID(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAnswerRepository(mockDB)

	ctx := context.Background()

	// Act
	results, err := repo.FindByAttemptID(ctx, uuid.Nil)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestPostgresAnswerRepository_Save_Success(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAnswerRepository(mockDB)

	ctx := context.Background()
	attemptID := uuid.New()

	// Crear múltiples respuestas
	now := time.Now()
	answer1 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionID:       "q1",
		StudentAnswer:    ptr("a"),
		IsCorrect:        ptr(true),
		TimeSpentSeconds: ptr(30),
		AnsweredAt:       now,
		CreatedAt:        now,
	}
	answer2 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionID:       "q2",
		StudentAnswer:    ptr("b"),
		IsCorrect:        ptr(false),
		TimeSpentSeconds: ptr(45),
		AnsweredAt:       now,
		CreatedAt:        now,
	}
	answer3 := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionID:       "q3",
		StudentAnswer:    ptr("c"),
		IsCorrect:        ptr(true),
		TimeSpentSeconds: ptr(20),
		AnsweredAt:       now,
		CreatedAt:        now,
	}

	answers := []*pgentities.AssessmentAttemptAnswer{answer1, answer2, answer3}

	// Act
	err := repo.Save(ctx, answers)

	// Assert
	assert.NoError(t, err)
}

func TestPostgresAnswerRepository_Save_EmptyArray(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAnswerRepository(mockDB)

	ctx := context.Background()

	// Act
	err := repo.Save(ctx, []*pgentities.AssessmentAttemptAnswer{})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no answers to save")
}

func TestPostgresAnswerRepository_Save_InvalidAnswer(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAnswerRepository(mockDB)

	ctx := context.Background()

	// Answer inválido (ID nil)
	invalidAnswer := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.Nil, // Inválido
		AttemptID:        uuid.New(),
		QuestionID:       "q1",
		SelectedAnswerID: "a",
		IsCorrect:        true,
		TimeSpentSeconds: 30,
	}

	// Act
	err := repo.Save(ctx, []*pgentities.AssessmentAttemptAnswer{invalidAnswer})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid answer")
}

func TestPostgresAnswerRepository_Save_SingleAnswer(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAnswerRepository(mockDB)

	ctx := context.Background()

	now := time.Now()
	answer := &pgentities.AssessmentAttemptAnswer{
		ID:               uuid.New(),
		AttemptID:        uuid.New(),
		QuestionID:       "q1",
		StudentAnswer:    ptr("a"),
		IsCorrect:        ptr(true),
		TimeSpentSeconds: ptr(30),
		AnsweredAt:       now,
		CreatedAt:        now,
	}

	// Act
	err := repo.Save(ctx, []*pgentities.AssessmentAttemptAnswer{answer})

	// Assert
	assert.NoError(t, err)
}

func TestPostgresAnswerRepository_FindByQuestionID_Success(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAnswerRepository(mockDB)

	ctx := context.Background()
	questionID := "q1"

	// Act
	results, err := repo.FindByQuestionID(ctx, questionID, 10, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, results)

	// Verificar que todas las respuestas son de la misma pregunta
	for _, answer := range results {
		assert.Equal(t, questionID, answer.QuestionID)
	}
}

func TestPostgresAnswerRepository_FindByQuestionID_InvalidID(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAnswerRepository(mockDB)

	ctx := context.Background()

	// Act
	results, err := repo.FindByQuestionID(ctx, "", 10, 0)

	// Assert
	assert.ErrorIs(t, err, domainErrors.ErrInvalidQuestionID)
	assert.Empty(t, results)
}

func TestPostgresAnswerRepository_FindByQuestionID_Pagination(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAnswerRepository(mockDB)

	ctx := context.Background()
	questionID := "q1"

	// Act - Primera página
	page1, err := repo.FindByQuestionID(ctx, questionID, 5, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, page1)

	// Act - Segunda página
	page2, err := repo.FindByQuestionID(ctx, questionID, 5, 5)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, page2)
}

func TestPostgresAnswerRepository_GetQuestionDifficultyStats_Success(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAnswerRepository(mockDB).(*PostgresAnswerRepository)

	ctx := context.Background()
	questionID := "q1"

	// Act
	total, correct, errorRate, err := repo.GetQuestionDifficultyStats(ctx, questionID)

	// Assert
	require.NoError(t, err)
	assert.Greater(t, total, 0)
	assert.GreaterOrEqual(t, correct, 0)
	assert.LessOrEqual(t, correct, total)
	assert.GreaterOrEqual(t, errorRate, 0.0)
	assert.LessOrEqual(t, errorRate, 1.0)
}

func TestPostgresAnswerRepository_GetQuestionDifficultyStats_InvalidID(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAnswerRepository(mockDB).(*PostgresAnswerRepository)

	ctx := context.Background()

	// Act
	total, correct, errorRate, err := repo.GetQuestionDifficultyStats(ctx, "")

	// Assert
	assert.ErrorIs(t, err, domainErrors.ErrInvalidQuestionID)
	assert.Equal(t, 0, total)
	assert.Equal(t, 0, correct)
	assert.Equal(t, 0.0, errorRate)
}

func TestPostgresAnswerRepository_GetQuestionDifficultyStats_ErrorRateCalculation(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAnswerRepository(mockDB).(*PostgresAnswerRepository)

	ctx := context.Background()
	questionID := "q1"

	// Act
	total, correct, errorRate, err := repo.GetQuestionDifficultyStats(ctx, questionID)

	// Assert
	require.NoError(t, err)

	// Verificar que el error rate es consistente
	expectedErrorRate := float64(total-correct) / float64(total)
	assert.InDelta(t, expectedErrorRate, errorRate, 0.01)
}
