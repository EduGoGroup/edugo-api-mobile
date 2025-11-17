package valueobjects_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobjects"
)

func TestNewQuestionID_Success(t *testing.T) {
	validIDs := []string{
		"q1",
		"q2",
		"question_123",
		"abc-def",
		"1",
	}

	for _, id := range validIDs {
		t.Run(id, func(t *testing.T) {
			questionID, err := valueobjects.NewQuestionID(id)
			require.NoError(t, err)
			assert.Equal(t, id, questionID.Value())
			assert.Equal(t, id, questionID.String())
		})
	}
}

func TestNewQuestionID_Empty(t *testing.T) {
	_, err := valueobjects.NewQuestionID("")
	assert.ErrorIs(t, err, valueobjects.ErrInvalidQuestionID)
}

func TestQuestionID_Equals(t *testing.T) {
	id1, _ := valueobjects.NewQuestionID("q1")
	id2, _ := valueobjects.NewQuestionID("q1")
	id3, _ := valueobjects.NewQuestionID("q2")

	assert.True(t, id1.Equals(id2))
	assert.False(t, id1.Equals(id3))
}
