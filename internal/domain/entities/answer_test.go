package entities_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entities"
	domainErrors "github.com/EduGoGroup/edugo-api-mobile/internal/domain/errors"
)

func TestNewAnswer_Success(t *testing.T) {
	attemptID := uuid.New()
	questionID := "q1"
	selectedAnswerID := "a2"
	isCorrect := true
	timeSpent := 45

	answer, err := entities.NewAnswer(attemptID, questionID, selectedAnswerID, isCorrect, timeSpent)

	require.NoError(t, err)
	require.NotNil(t, answer)

	assert.NotEqual(t, uuid.Nil, answer.ID)
	assert.Equal(t, attemptID, answer.AttemptID)
	assert.Equal(t, questionID, answer.QuestionID)
	assert.Equal(t, selectedAnswerID, answer.SelectedAnswerID)
	assert.Equal(t, isCorrect, answer.IsCorrect)
	assert.Equal(t, timeSpent, answer.TimeSpentSeconds)
	assert.False(t, answer.CreatedAt.IsZero())
}

func TestNewAnswer_InvalidAttemptID(t *testing.T) {
	_, err := entities.NewAnswer(uuid.Nil, "q1", "a1", true, 30)
	assert.ErrorIs(t, err, domainErrors.ErrInvalidAttemptID)
}

func TestNewAnswer_InvalidQuestionID(t *testing.T) {
	_, err := entities.NewAnswer(uuid.New(), "", "a1", true, 30)
	assert.ErrorIs(t, err, domainErrors.ErrInvalidQuestionID)
}

func TestNewAnswer_InvalidSelectedAnswerID(t *testing.T) {
	_, err := entities.NewAnswer(uuid.New(), "q1", "", true, 30)
	assert.ErrorIs(t, err, domainErrors.ErrInvalidSelectedAnswerID)
}

func TestNewAnswer_NegativeTimeSpent(t *testing.T) {
	_, err := entities.NewAnswer(uuid.New(), "q1", "a1", true, -10)
	assert.ErrorIs(t, err, domainErrors.ErrInvalidTimeSpent)
}

func TestNewAnswer_ZeroTimeSpent(t *testing.T) {
	// Zero time spent is valid (answered instantly)
	answer, err := entities.NewAnswer(uuid.New(), "q1", "a1", true, 0)
	require.NoError(t, err)
	assert.Equal(t, 0, answer.TimeSpentSeconds)
}

func TestAnswer_Validate_Success(t *testing.T) {
	answer, _ := entities.NewAnswer(uuid.New(), "q1", "a1", true, 30)
	err := answer.Validate()
	assert.NoError(t, err)
}

func TestAnswer_Validate_InvalidID(t *testing.T) {
	answer, _ := entities.NewAnswer(uuid.New(), "q1", "a1", true, 30)
	answer.ID = uuid.Nil
	err := answer.Validate()
	assert.ErrorIs(t, err, domainErrors.ErrInvalidAnswerID)
}

func TestAnswer_Validate_InvalidAttemptID(t *testing.T) {
	answer, _ := entities.NewAnswer(uuid.New(), "q1", "a1", true, 30)
	answer.AttemptID = uuid.Nil
	err := answer.Validate()
	assert.ErrorIs(t, err, domainErrors.ErrInvalidAttemptID)
}

func TestAnswer_Validate_EmptyQuestionID(t *testing.T) {
	answer, _ := entities.NewAnswer(uuid.New(), "q1", "a1", true, 30)
	answer.QuestionID = ""
	err := answer.Validate()
	assert.ErrorIs(t, err, domainErrors.ErrInvalidQuestionID)
}

func TestAnswer_Validate_EmptySelectedAnswerID(t *testing.T) {
	answer, _ := entities.NewAnswer(uuid.New(), "q1", "a1", true, 30)
	answer.SelectedAnswerID = ""
	err := answer.Validate()
	assert.ErrorIs(t, err, domainErrors.ErrInvalidSelectedAnswerID)
}

func TestAnswer_Validate_NegativeTimeSpent(t *testing.T) {
	answer, _ := entities.NewAnswer(uuid.New(), "q1", "a1", true, 30)
	answer.TimeSpentSeconds = -5
	err := answer.Validate()
	assert.ErrorIs(t, err, domainErrors.ErrInvalidTimeSpent)
}

func TestAnswer_IsCorrect_Variations(t *testing.T) {
	testCases := []struct {
		name      string
		isCorrect bool
	}{
		{"correct answer", true},
		{"incorrect answer", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			answer, err := entities.NewAnswer(uuid.New(), "q1", "a1", tc.isCorrect, 30)
			require.NoError(t, err)
			assert.Equal(t, tc.isCorrect, answer.IsCorrect)
		})
	}
}
