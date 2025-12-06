package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
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
			questionID, err := valueobject.NewQuestionID(id)
			require.NoError(t, err)
			assert.Equal(t, id, questionID.Value())
			assert.Equal(t, id, questionID.String())
		})
	}
}

func TestNewQuestionID_Empty(t *testing.T) {
	_, err := valueobject.NewQuestionID("")
	assert.ErrorIs(t, err, valueobject.ErrInvalidQuestionID)
}

func TestQuestionID_Equals(t *testing.T) {
	id1, _ := valueobject.NewQuestionID("q1")
	id2, _ := valueobject.NewQuestionID("q1")
	id3, _ := valueobject.NewQuestionID("q2")

	assert.True(t, id1.Equals(id2))
	assert.False(t, id1.Equals(id3))
}
