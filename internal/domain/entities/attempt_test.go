package entities_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entities"
	domainErrors "github.com/EduGoGroup/edugo-api-mobile/internal/domain/errors"
)

func TestNewAttempt_Success(t *testing.T) {
	assessmentID := uuid.New()
	studentID := uuid.New()

	// Crear 5 respuestas: 3 correctas, 2 incorrectas (60%)
	answers := []*entities.Answer{
		{ID: uuid.New(), QuestionID: "q1", SelectedAnswerID: "a1", IsCorrect: true, TimeSpentSeconds: 30},
		{ID: uuid.New(), QuestionID: "q2", SelectedAnswerID: "a2", IsCorrect: false, TimeSpentSeconds: 45},
		{ID: uuid.New(), QuestionID: "q3", SelectedAnswerID: "a3", IsCorrect: true, TimeSpentSeconds: 60},
		{ID: uuid.New(), QuestionID: "q4", SelectedAnswerID: "a4", IsCorrect: false, TimeSpentSeconds: 50},
		{ID: uuid.New(), QuestionID: "q5", SelectedAnswerID: "a5", IsCorrect: true, TimeSpentSeconds: 40},
	}

	startedAt := time.Now().Add(-5 * time.Minute)
	completedAt := time.Now()

	attempt, err := entities.NewAttempt(assessmentID, studentID, answers, startedAt, completedAt)

	require.NoError(t, err)
	require.NotNil(t, attempt)

	assert.NotEqual(t, uuid.Nil, attempt.ID)
	assert.Equal(t, assessmentID, attempt.AssessmentID)
	assert.Equal(t, studentID, attempt.StudentID)
	assert.Equal(t, 60, attempt.Score, "Score should be 60% (3/5 correct)")
	assert.Equal(t, 100, attempt.MaxScore)
	assert.Equal(t, 5, len(attempt.Answers))
	assert.True(t, attempt.TimeSpentSeconds > 0)
	assert.False(t, attempt.CreatedAt.IsZero())
}

func TestNewAttempt_ScoreCalculation(t *testing.T) {
	testCases := []struct {
		name           string
		correctCount   int
		totalQuestions int
		expectedScore  int
	}{
		{"0% - all wrong", 0, 5, 0},
		{"20% - 1 of 5", 1, 5, 20},
		{"40% - 2 of 5", 2, 5, 40},
		{"60% - 3 of 5", 3, 5, 60},
		{"80% - 4 of 5", 4, 5, 80},
		{"100% - all correct", 5, 5, 100},
		{"100% - single question", 1, 1, 100},
		{"75% - 3 of 4", 3, 4, 75},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			answers := make([]*entities.Answer, tc.totalQuestions)
			for i := 0; i < tc.totalQuestions; i++ {
				isCorrect := i < tc.correctCount
				answers[i] = &entities.Answer{
					ID:               uuid.New(),
					QuestionID:       "q",
					SelectedAnswerID: "a1",
					IsCorrect:        isCorrect,
					TimeSpentSeconds: 30,
				}
			}

			startedAt := time.Now().Add(-3 * time.Minute)
			completedAt := time.Now()

			attempt, err := entities.NewAttempt(uuid.New(), uuid.New(), answers, startedAt, completedAt)

			require.NoError(t, err)
			assert.Equal(t, tc.expectedScore, attempt.Score)
		})
	}
}

func TestNewAttempt_InvalidAssessmentID(t *testing.T) {
	answers := []*entities.Answer{{ID: uuid.New(), QuestionID: "q1", IsCorrect: true}}

	_, err := entities.NewAttempt(
		uuid.Nil,
		uuid.New(),
		answers,
		time.Now().Add(-1*time.Minute),
		time.Now(),
	)

	assert.ErrorIs(t, err, domainErrors.ErrInvalidAssessmentID)
}

func TestNewAttempt_InvalidStudentID(t *testing.T) {
	answers := []*entities.Answer{{ID: uuid.New(), QuestionID: "q1", IsCorrect: true}}

	_, err := entities.NewAttempt(
		uuid.New(),
		uuid.Nil,
		answers,
		time.Now().Add(-1*time.Minute),
		time.Now(),
	)

	assert.ErrorIs(t, err, domainErrors.ErrInvalidStudentID)
}

func TestNewAttempt_NoAnswers(t *testing.T) {
	_, err := entities.NewAttempt(
		uuid.New(),
		uuid.New(),
		[]*entities.Answer{},
		time.Now().Add(-1*time.Minute),
		time.Now(),
	)

	assert.ErrorIs(t, err, domainErrors.ErrNoAnswersProvided)
}

func TestNewAttempt_InvalidStartTime(t *testing.T) {
	answers := []*entities.Answer{{ID: uuid.New(), QuestionID: "q1", IsCorrect: true}}

	_, err := entities.NewAttempt(
		uuid.New(),
		uuid.New(),
		answers,
		time.Time{}, // Zero time
		time.Now(),
	)

	assert.ErrorIs(t, err, domainErrors.ErrInvalidStartTime)
}

func TestNewAttempt_InvalidEndTime(t *testing.T) {
	answers := []*entities.Answer{{ID: uuid.New(), QuestionID: "q1", IsCorrect: true}}
	now := time.Now()

	testCases := []struct {
		name        string
		startedAt   time.Time
		completedAt time.Time
	}{
		{"zero time", now, time.Time{}},
		{"completed before started", now, now.Add(-1 * time.Minute)},
		{"completed same as started", now, now},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := entities.NewAttempt(
				uuid.New(),
				uuid.New(),
				answers,
				tc.startedAt,
				tc.completedAt,
			)

			assert.ErrorIs(t, err, domainErrors.ErrInvalidEndTime)
		})
	}
}

func TestNewAttempt_InvalidTimeSpent(t *testing.T) {
	answers := []*entities.Answer{{ID: uuid.New(), QuestionID: "q1", IsCorrect: true}}
	now := time.Now()

	testCases := []struct {
		name        string
		startedAt   time.Time
		completedAt time.Time
	}{
		{"over 2 hours", now.Add(-3 * time.Hour), now},
		{"exactly 2 hours + 1 second", now.Add(-2*time.Hour - 1*time.Second), now},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := entities.NewAttempt(
				uuid.New(),
				uuid.New(),
				answers,
				tc.startedAt,
				tc.completedAt,
			)

			assert.ErrorIs(t, err, domainErrors.ErrInvalidTimeSpent)
		})
	}
}

func TestNewAttemptWithIdempotency_Success(t *testing.T) {
	answers := []*entities.Answer{{ID: uuid.New(), QuestionID: "q1", IsCorrect: true}}
	startedAt := time.Now().Add(-1 * time.Minute)
	completedAt := time.Now()
	idempotencyKey := "unique-key-12345"

	attempt, err := entities.NewAttemptWithIdempotency(
		uuid.New(),
		uuid.New(),
		answers,
		startedAt,
		completedAt,
		idempotencyKey,
	)

	require.NoError(t, err)
	require.NotNil(t, attempt.IdempotencyKey)
	assert.Equal(t, idempotencyKey, *attempt.IdempotencyKey)
}

func TestAttempt_IsPassed(t *testing.T) {
	testCases := []struct {
		name          string
		score         int
		passThreshold int
		shouldPass    bool
	}{
		{"60% with 70% threshold - fail", 60, 70, false},
		{"70% with 70% threshold - pass", 70, 70, true},
		{"80% with 70% threshold - pass", 80, 70, true},
		{"100% with 70% threshold - pass", 100, 70, true},
		{"0% - fail", 0, 70, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Crear intento con score específico
			totalQuestions := 10
			correctCount := (tc.score * totalQuestions) / 100

			answers := make([]*entities.Answer, totalQuestions)
			for i := 0; i < totalQuestions; i++ {
				answers[i] = &entities.Answer{
					ID:               uuid.New(),
					QuestionID:       "q",
					IsCorrect:        i < correctCount,
					TimeSpentSeconds: 30,
				}
			}

			attempt, _ := entities.NewAttempt(
				uuid.New(),
				uuid.New(),
				answers,
				time.Now().Add(-5*time.Minute),
				time.Now(),
			)

			assert.Equal(t, tc.shouldPass, attempt.IsPassed(tc.passThreshold))
		})
	}
}

func TestAttempt_GetCorrectAnswersCount(t *testing.T) {
	answers := []*entities.Answer{
		{ID: uuid.New(), QuestionID: "q1", IsCorrect: true},
		{ID: uuid.New(), QuestionID: "q2", IsCorrect: false},
		{ID: uuid.New(), QuestionID: "q3", IsCorrect: true},
		{ID: uuid.New(), QuestionID: "q4", IsCorrect: true},
	}

	attempt, _ := entities.NewAttempt(
		uuid.New(),
		uuid.New(),
		answers,
		time.Now().Add(-2*time.Minute),
		time.Now(),
	)

	assert.Equal(t, 3, attempt.GetCorrectAnswersCount())
	assert.Equal(t, 1, attempt.GetIncorrectAnswersCount())
	assert.Equal(t, 4, attempt.GetTotalQuestions())
}

func TestAttempt_GetAccuracyPercentage(t *testing.T) {
	// 3 correct out of 5 = 60%
	answers := make([]*entities.Answer, 5)
	for i := 0; i < 5; i++ {
		answers[i] = &entities.Answer{
			ID:               uuid.New(),
			QuestionID:       "q",
			IsCorrect:        i < 3,
			TimeSpentSeconds: 30,
		}
	}

	attempt, _ := entities.NewAttempt(
		uuid.New(),
		uuid.New(),
		answers,
		time.Now().Add(-3*time.Minute),
		time.Now(),
	)

	assert.Equal(t, 60, attempt.GetAccuracyPercentage())
}

func TestAttempt_GetAverageTimePerQuestion(t *testing.T) {
	answers := []*entities.Answer{
		{ID: uuid.New(), QuestionID: "q1", IsCorrect: true, TimeSpentSeconds: 30},
		{ID: uuid.New(), QuestionID: "q2", IsCorrect: true, TimeSpentSeconds: 60},
		{ID: uuid.New(), QuestionID: "q3", IsCorrect: true, TimeSpentSeconds: 90},
	}

	startedAt := time.Now().Add(-3 * time.Minute)
	completedAt := time.Now()

	attempt, _ := entities.NewAttempt(
		uuid.New(),
		uuid.New(),
		answers,
		startedAt,
		completedAt,
	)

	avgTime := attempt.GetAverageTimePerQuestion()
	assert.True(t, avgTime > 0)
	assert.True(t, avgTime <= 180) // ~180 segundos / 3 preguntas = ~60 seg promedio
}

func TestAttempt_GetAverageTimePerQuestion_ZeroAnswers(t *testing.T) {
	// Esta situación no debería ocurrir debido a validaciones en NewAttempt,
	// pero si ocurre, debe retornar 0
	attempt := &entities.Attempt{
		Answers:          []*entities.Answer{},
		TimeSpentSeconds: 100,
	}

	assert.Equal(t, 0, attempt.GetAverageTimePerQuestion())
}

func TestAttempt_Validate_Success(t *testing.T) {
	answers := []*entities.Answer{
		{ID: uuid.New(), QuestionID: "q1", IsCorrect: true, TimeSpentSeconds: 60},
	}

	attempt, _ := entities.NewAttempt(
		uuid.New(),
		uuid.New(),
		answers,
		time.Now().Add(-2*time.Minute),
		time.Now(),
	)

	err := attempt.Validate()
	assert.NoError(t, err)
}

func TestAttempt_Validate_InvalidFields(t *testing.T) {
	// Crear un intento válido base
	createValidAttempt := func() *entities.Attempt {
		answers := []*entities.Answer{
			{ID: uuid.New(), QuestionID: "q1", IsCorrect: true, TimeSpentSeconds: 60},
		}
		attempt, _ := entities.NewAttempt(
			uuid.New(),
			uuid.New(),
			answers,
			time.Now().Add(-2*time.Minute),
			time.Now(),
		)
		return attempt
	}

	testCases := []struct {
		name          string
		modifier      func(*entities.Attempt)
		expectedError error
	}{
		{
			"invalid attempt ID",
			func(a *entities.Attempt) { a.ID = uuid.Nil },
			domainErrors.ErrInvalidAttemptID,
		},
		{
			"invalid assessment ID",
			func(a *entities.Attempt) { a.AssessmentID = uuid.Nil },
			domainErrors.ErrInvalidAssessmentID,
		},
		{
			"invalid student ID",
			func(a *entities.Attempt) { a.StudentID = uuid.Nil },
			domainErrors.ErrInvalidStudentID,
		},
		{
			"negative score",
			func(a *entities.Attempt) { a.Score = -1 },
			domainErrors.ErrInvalidScore,
		},
		{
			"score over 100",
			func(a *entities.Attempt) { a.Score = 101 },
			domainErrors.ErrInvalidScore,
		},
		{
			"zero time spent",
			func(a *entities.Attempt) { a.TimeSpentSeconds = 0 },
			domainErrors.ErrInvalidTimeSpent,
		},
		{
			"negative time spent",
			func(a *entities.Attempt) { a.TimeSpentSeconds = -1 },
			domainErrors.ErrInvalidTimeSpent,
		},
		{
			"time spent over 2 hours",
			func(a *entities.Attempt) { a.TimeSpentSeconds = 7201 },
			domainErrors.ErrInvalidTimeSpent,
		},
		{
			"zero start time",
			func(a *entities.Attempt) { a.StartedAt = time.Time{} },
			domainErrors.ErrInvalidStartTime,
		},
		{
			"zero completed time",
			func(a *entities.Attempt) { a.CompletedAt = time.Time{} },
			domainErrors.ErrInvalidEndTime,
		},
		{
			"completed before started",
			func(a *entities.Attempt) {
				a.CompletedAt = a.StartedAt.Add(-1 * time.Minute)
			},
			domainErrors.ErrInvalidEndTime,
		},
		{
			"no answers",
			func(a *entities.Attempt) { a.Answers = []*entities.Answer{} },
			domainErrors.ErrNoAnswersProvided,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attempt := createValidAttempt()
			tc.modifier(attempt)
			err := attempt.Validate()
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestAttempt_Validate_ScoreMismatch(t *testing.T) {
	answers := []*entities.Answer{
		{ID: uuid.New(), QuestionID: "q1", IsCorrect: true}, // 1/1 = 100%
	}

	attempt, _ := entities.NewAttempt(
		uuid.New(),
		uuid.New(),
		answers,
		time.Now().Add(-1*time.Minute),
		time.Now(),
	)

	// Manualmente cambiar el score para que no coincida
	attempt.Score = 50 // Debería ser 100%

	err := attempt.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "score mismatch")
}
