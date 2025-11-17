package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entities"
	domainErrors "github.com/EduGoGroup/edugo-api-mobile/internal/domain/errors"
)

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
	answer1, _ := entities.NewAnswer(attemptID, "q1", "a", true, 30)
	answer2, _ := entities.NewAnswer(attemptID, "q2", "b", false, 45)
	answer3, _ := entities.NewAnswer(attemptID, "q3", "c", true, 20)

	answers := []*entities.Answer{answer1, answer2, answer3}

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
	err := repo.Save(ctx, []*entities.Answer{})

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
	invalidAnswer := &entities.Answer{
		ID:               uuid.Nil, // Inválido
		AttemptID:        uuid.New(),
		QuestionID:       "q1",
		SelectedAnswerID: "a",
		IsCorrect:        true,
		TimeSpentSeconds: 30,
	}

	// Act
	err := repo.Save(ctx, []*entities.Answer{invalidAnswer})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid answer")
}

func TestPostgresAnswerRepository_Save_SingleAnswer(t *testing.T) {
	// Arrange
	var mockDB *sql.DB
	repo := NewPostgresAnswerRepository(mockDB)

	ctx := context.Background()

	answer, err := entities.NewAnswer(uuid.New(), "q1", "a", true, 30)
	require.NoError(t, err)

	// Act
	err = repo.Save(ctx, []*entities.Answer{answer})

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

// TODO: Claude Local - Tests de integración con testcontainers
// func TestPostgresAnswerRepository_Integration_BatchInsert(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping integration test")
// 	}
//
// 	// Test que verifica:
// 	// 1. Batch insert de múltiples respuestas
// 	// 2. Recuperación por attempt_id
// 	// 3. Recuperación por question_id con paginación
// 	// 4. Estadísticas de dificultad de preguntas
// }
