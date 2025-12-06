package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
)

func TestNewTimeSpent_Success(t *testing.T) {
	testCases := []struct {
		seconds         int
		expectedMinutes int
	}{
		{0, 0},
		{60, 1},
		{120, 2},
		{300, 5},
		{3600, 60},  // 1 hour
		{7200, 120}, // 2 hours (maximum)
	}

	for _, tc := range testCases {
		timeSpent, err := valueobject.NewTimeSpent(tc.seconds)
		require.NoError(t, err)
		assert.Equal(t, tc.seconds, timeSpent.Seconds())
		assert.Equal(t, tc.expectedMinutes, timeSpent.Minutes())
	}
}

func TestNewTimeSpent_Negative(t *testing.T) {
	testCases := []int{-1, -10, -100}

	for _, seconds := range testCases {
		_, err := valueobject.NewTimeSpent(seconds)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be negative")
	}
}

func TestNewTimeSpent_ExceedsMaximum(t *testing.T) {
	testCases := []int{7201, 8000, 10000}

	for _, seconds := range testCases {
		_, err := valueobject.NewTimeSpent(seconds)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot exceed 7200 seconds")
	}
}

func TestTimeSpent_String(t *testing.T) {
	testCases := []struct {
		seconds  int
		expected string
	}{
		{0, "0m0s"},
		{30, "0m30s"},
		{60, "1m0s"},
		{90, "1m30s"},
		{300, "5m0s"},
		{365, "6m5s"},
	}

	for _, tc := range testCases {
		timeSpent, _ := valueobject.NewTimeSpent(tc.seconds)
		assert.Equal(t, tc.expected, timeSpent.String())
	}
}

func TestTimeSpent_Equals(t *testing.T) {
	time1, _ := valueobject.NewTimeSpent(300)
	time2, _ := valueobject.NewTimeSpent(300)
	time3, _ := valueobject.NewTimeSpent(400)

	assert.True(t, time1.Equals(time2))
	assert.False(t, time1.Equals(time3))
}

func TestTimeSpent_Minutes_Conversion(t *testing.T) {
	testCases := []struct {
		seconds int
		minutes int
	}{
		{0, 0},
		{59, 0},
		{60, 1},
		{61, 1},
		{119, 1},
		{120, 2},
	}

	for _, tc := range testCases {
		timeSpent, _ := valueobject.NewTimeSpent(tc.seconds)
		assert.Equal(t, tc.minutes, timeSpent.Minutes())
	}
}
