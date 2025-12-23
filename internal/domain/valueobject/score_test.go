package valueobject_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
)

func TestNewScore_Success(t *testing.T) {
	testCases := []int{0, 10, 50, 70, 100}

	for _, value := range testCases {
		score, err := valueobject.NewScore(value)
		require.NoError(t, err)
		assert.Equal(t, value, score.Value())
	}
}

func TestNewScore_Invalid(t *testing.T) {
	testCases := []int{-1, -10, 101, 150, 1000}

	for _, value := range testCases {
		t.Run(fmt.Sprintf("value_%d", value), func(t *testing.T) {
			_, err := valueobject.NewScore(value)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "must be between 0 and 100")
		})
	}
}

func TestScore_IsPassing(t *testing.T) {
	testCases := []struct {
		score     int
		threshold int
		expected  bool
	}{
		{60, 70, false},
		{70, 70, true},
		{80, 70, true},
		{100, 70, true},
		{0, 70, false},
		{100, 100, true},
		{99, 100, false},
	}

	for _, tc := range testCases {
		score, _ := valueobject.NewScore(tc.score)
		result := score.IsPassing(tc.threshold)
		assert.Equal(t, tc.expected, result, "Score %d with threshold %d", tc.score, tc.threshold)
	}
}

func TestScore_IsFailing(t *testing.T) {
	score, _ := valueobject.NewScore(60)
	assert.True(t, score.IsFailing(70))
	assert.False(t, score.IsFailing(60))
	assert.False(t, score.IsFailing(50))
}

func TestScore_String(t *testing.T) {
	testCases := []struct {
		value    int
		expected string
	}{
		{0, "0%"},
		{50, "50%"},
		{100, "100%"},
	}

	for _, tc := range testCases {
		score, _ := valueobject.NewScore(tc.value)
		assert.Equal(t, tc.expected, score.String())
	}
}

func TestScore_Equals(t *testing.T) {
	score1, _ := valueobject.NewScore(80)
	score2, _ := valueobject.NewScore(80)
	score3, _ := valueobject.NewScore(70)

	assert.True(t, score1.Equals(score2))
	assert.False(t, score1.Equals(score3))
}
