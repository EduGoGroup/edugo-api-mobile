package entities_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"edugo-api-mobile/internal/domain/entities"
	domainErrors "edugo-api-mobile/internal/domain/errors"
)

func TestNewAssessment_Success(t *testing.T) {
	materialID := uuid.New()
	mongoDocID := "507f1f77bcf86cd799439011"
	title := "Cuestionario: Introducción a Pascal"
	totalQuestions := 5
	passThreshold := 70

	assessment, err := entities.NewAssessment(
		materialID,
		mongoDocID,
		title,
		totalQuestions,
		passThreshold,
	)

	require.NoError(t, err)
	require.NotNil(t, assessment)

	assert.NotEqual(t, uuid.Nil, assessment.ID)
	assert.Equal(t, materialID, assessment.MaterialID)
	assert.Equal(t, mongoDocID, assessment.MongoDocumentID)
	assert.Equal(t, title, assessment.Title)
	assert.Equal(t, totalQuestions, assessment.TotalQuestions)
	assert.Equal(t, passThreshold, assessment.PassThreshold)
	assert.Nil(t, assessment.MaxAttempts, "MaxAttempts should be nil by default")
	assert.Nil(t, assessment.TimeLimitMinutes, "TimeLimitMinutes should be nil by default")
	assert.False(t, assessment.CreatedAt.IsZero())
	assert.False(t, assessment.UpdatedAt.IsZero())
}

func TestNewAssessment_InvalidMaterialID(t *testing.T) {
	_, err := entities.NewAssessment(
		uuid.Nil,
		"507f1f77bcf86cd799439011",
		"Title",
		5,
		70,
	)

	assert.ErrorIs(t, err, domainErrors.ErrInvalidMaterialID)
}

func TestNewAssessment_InvalidMongoDocumentID(t *testing.T) {
	testCases := []struct {
		name       string
		mongoDocID string
	}{
		{"empty", ""},
		{"too short", "123"},
		{"too long", "507f1f77bcf86cd799439011EXTRA"},
		{"wrong length", "507f1f77bcf86cd79943901"},
		{"23 characters", "507f1f77bcf86cd7994390"},
		{"25 characters", "507f1f77bcf86cd799439011X"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := entities.NewAssessment(
				uuid.New(),
				tc.mongoDocID,
				"Title",
				5,
				70,
			)

			assert.ErrorIs(t, err, domainErrors.ErrInvalidMongoDocumentID)
		})
	}
}

func TestNewAssessment_EmptyTitle(t *testing.T) {
	_, err := entities.NewAssessment(
		uuid.New(),
		"507f1f77bcf86cd799439011",
		"",
		5,
		70,
	)

	assert.ErrorIs(t, err, domainErrors.ErrEmptyTitle)
}

func TestNewAssessment_InvalidTotalQuestions(t *testing.T) {
	testCases := []struct {
		name           string
		totalQuestions int
	}{
		{"zero questions", 0},
		{"negative questions", -1},
		{"too many questions", 101},
		{"way too many", 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := entities.NewAssessment(
				uuid.New(),
				"507f1f77bcf86cd799439011",
				"Title",
				tc.totalQuestions,
				70,
			)

			assert.ErrorIs(t, err, domainErrors.ErrInvalidTotalQuestions)
		})
	}
}

func TestNewAssessment_ValidTotalQuestions(t *testing.T) {
	testCases := []int{1, 5, 10, 50, 100}

	for _, totalQuestions := range testCases {
		t.Run(string(rune(totalQuestions)), func(t *testing.T) {
			assessment, err := entities.NewAssessment(
				uuid.New(),
				"507f1f77bcf86cd799439011",
				"Title",
				totalQuestions,
				70,
			)

			require.NoError(t, err)
			assert.Equal(t, totalQuestions, assessment.TotalQuestions)
		})
	}
}

func TestNewAssessment_InvalidPassThreshold(t *testing.T) {
	testCases := []struct {
		name          string
		passThreshold int
	}{
		{"negative threshold", -1},
		{"above 100", 101},
		{"way above 100", 150},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := entities.NewAssessment(
				uuid.New(),
				"507f1f77bcf86cd799439011",
				"Title",
				5,
				tc.passThreshold,
			)

			assert.ErrorIs(t, err, domainErrors.ErrInvalidPassThreshold)
		})
	}
}

func TestNewAssessment_ValidPassThreshold(t *testing.T) {
	testCases := []int{0, 50, 70, 100}

	for _, passThreshold := range testCases {
		t.Run(string(rune(passThreshold)), func(t *testing.T) {
			assessment, err := entities.NewAssessment(
				uuid.New(),
				"507f1f77bcf86cd799439011",
				"Title",
				5,
				passThreshold,
			)

			require.NoError(t, err)
			assert.Equal(t, passThreshold, assessment.PassThreshold)
		})
	}
}

func TestAssessment_Validate_Success(t *testing.T) {
	assessment, _ := entities.NewAssessment(
		uuid.New(),
		"507f1f77bcf86cd799439011",
		"Title",
		5,
		70,
	)

	err := assessment.Validate()
	assert.NoError(t, err)
}

func TestAssessment_Validate_WithLimits(t *testing.T) {
	assessment, _ := entities.NewAssessment(
		uuid.New(),
		"507f1f77bcf86cd799439011",
		"Title",
		5,
		70,
	)

	assessment.SetMaxAttempts(3)
	assessment.SetTimeLimit(30)

	err := assessment.Validate()
	assert.NoError(t, err)
}

func TestAssessment_CanAttempt_Unlimited(t *testing.T) {
	assessment, _ := entities.NewAssessment(
		uuid.New(),
		"507f1f77bcf86cd799439011",
		"Title",
		5,
		70,
	)

	// Sin límite de intentos (MaxAttempts = nil)
	assert.True(t, assessment.CanAttempt(0))
	assert.True(t, assessment.CanAttempt(10))
	assert.True(t, assessment.CanAttempt(100))
	assert.True(t, assessment.CanAttempt(1000))
}

func TestAssessment_CanAttempt_WithLimit(t *testing.T) {
	assessment, _ := entities.NewAssessment(
		uuid.New(),
		"507f1f77bcf86cd799439011",
		"Title",
		5,
		70,
	)

	// Establecer límite de 3 intentos
	err := assessment.SetMaxAttempts(3)
	require.NoError(t, err)

	assert.True(t, assessment.CanAttempt(0), "should allow attempt 1")
	assert.True(t, assessment.CanAttempt(1), "should allow attempt 2")
	assert.True(t, assessment.CanAttempt(2), "should allow attempt 3")
	assert.False(t, assessment.CanAttempt(3), "should NOT allow attempt 4")
	assert.False(t, assessment.CanAttempt(4), "should NOT allow attempt 5")
}

func TestAssessment_SetMaxAttempts_Success(t *testing.T) {
	assessment, _ := entities.NewAssessment(
		uuid.New(),
		"507f1f77bcf86cd799439011",
		"Title",
		5,
		70,
	)

	oldUpdatedAt := assessment.UpdatedAt
	time.Sleep(10 * time.Millisecond)

	err := assessment.SetMaxAttempts(5)

	assert.NoError(t, err)
	assert.NotNil(t, assessment.MaxAttempts)
	assert.Equal(t, 5, *assessment.MaxAttempts)
	assert.True(t, assessment.UpdatedAt.After(oldUpdatedAt), "UpdatedAt should be updated")
}

func TestAssessment_SetMaxAttempts_Invalid(t *testing.T) {
	assessment, _ := entities.NewAssessment(
		uuid.New(),
		"507f1f77bcf86cd799439011",
		"Title",
		5,
		70,
	)

	testCases := []int{0, -1, -10}

	for _, maxAttempts := range testCases {
		t.Run(string(rune(maxAttempts)), func(t *testing.T) {
			err := assessment.SetMaxAttempts(maxAttempts)
			assert.ErrorIs(t, err, domainErrors.ErrInvalidMaxAttempts)
		})
	}
}

func TestAssessment_RemoveMaxAttempts(t *testing.T) {
	assessment, _ := entities.NewAssessment(
		uuid.New(),
		"507f1f77bcf86cd799439011",
		"Title",
		5,
		70,
	)

	// Primero establecer límite
	assessment.SetMaxAttempts(3)
	assert.NotNil(t, assessment.MaxAttempts)

	oldUpdatedAt := assessment.UpdatedAt
	time.Sleep(10 * time.Millisecond)

	// Luego quitar límite
	assessment.RemoveMaxAttempts()
	assert.Nil(t, assessment.MaxAttempts)
	assert.True(t, assessment.UpdatedAt.After(oldUpdatedAt), "UpdatedAt should be updated")

	// Ahora debería permitir intentos ilimitados
	assert.True(t, assessment.CanAttempt(100))
}

func TestAssessment_IsTimeLimited(t *testing.T) {
	assessment, _ := entities.NewAssessment(
		uuid.New(),
		"507f1f77bcf86cd799439011",
		"Title",
		5,
		70,
	)

	// Sin límite de tiempo
	assert.False(t, assessment.IsTimeLimited())

	// Con límite de tiempo
	assessment.SetTimeLimit(30)
	assert.True(t, assessment.IsTimeLimited())

	// Quitar límite
	assessment.RemoveTimeLimit()
	assert.False(t, assessment.IsTimeLimited())
}

func TestAssessment_SetTimeLimit_Success(t *testing.T) {
	assessment, _ := entities.NewAssessment(
		uuid.New(),
		"507f1f77bcf86cd799439011",
		"Title",
		5,
		70,
	)

	testCases := []int{1, 15, 30, 60, 120, 180}

	for _, minutes := range testCases {
		t.Run(string(rune(minutes)), func(t *testing.T) {
			oldUpdatedAt := assessment.UpdatedAt
			time.Sleep(10 * time.Millisecond)

			err := assessment.SetTimeLimit(minutes)
			assert.NoError(t, err)
			assert.NotNil(t, assessment.TimeLimitMinutes)
			assert.Equal(t, minutes, *assessment.TimeLimitMinutes)
			assert.True(t, assessment.UpdatedAt.After(oldUpdatedAt), "UpdatedAt should be updated")
		})
	}
}

func TestAssessment_SetTimeLimit_Invalid(t *testing.T) {
	assessment, _ := entities.NewAssessment(
		uuid.New(),
		"507f1f77bcf86cd799439011",
		"Title",
		5,
		70,
	)

	testCases := []int{0, -1, 181, 200, 1000}

	for _, minutes := range testCases {
		t.Run(string(rune(minutes)), func(t *testing.T) {
			err := assessment.SetTimeLimit(minutes)
			assert.ErrorIs(t, err, domainErrors.ErrInvalidTimeLimit)
		})
	}
}

func TestAssessment_RemoveTimeLimit(t *testing.T) {
	assessment, _ := entities.NewAssessment(
		uuid.New(),
		"507f1f77bcf86cd799439011",
		"Title",
		5,
		70,
	)

	// Primero establecer límite
	assessment.SetTimeLimit(30)
	assert.NotNil(t, assessment.TimeLimitMinutes)
	assert.True(t, assessment.IsTimeLimited())

	oldUpdatedAt := assessment.UpdatedAt
	time.Sleep(10 * time.Millisecond)

	// Luego quitar límite
	assessment.RemoveTimeLimit()
	assert.Nil(t, assessment.TimeLimitMinutes)
	assert.False(t, assessment.IsTimeLimited())
	assert.True(t, assessment.UpdatedAt.After(oldUpdatedAt), "UpdatedAt should be updated")
}

func TestAssessment_Validate_InvalidMaxAttempts(t *testing.T) {
	assessment, _ := entities.NewAssessment(
		uuid.New(),
		"507f1f77bcf86cd799439011",
		"Title",
		5,
		70,
	)

	// Forzar un valor inválido (no se puede hacer via SetMaxAttempts)
	invalidValue := 0
	assessment.MaxAttempts = &invalidValue

	err := assessment.Validate()
	assert.ErrorIs(t, err, domainErrors.ErrInvalidMaxAttempts)
}

func TestAssessment_Validate_InvalidTimeLimit(t *testing.T) {
	assessment, _ := entities.NewAssessment(
		uuid.New(),
		"507f1f77bcf86cd799439011",
		"Title",
		5,
		70,
	)

	testCases := []struct {
		name    string
		minutes int
	}{
		{"zero minutes", 0},
		{"negative minutes", -1},
		{"over 180 minutes", 181},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Forzar un valor inválido directamente
			assessment.TimeLimitMinutes = &tc.minutes

			err := assessment.Validate()
			assert.ErrorIs(t, err, domainErrors.ErrInvalidTimeLimit)
		})
	}
}
