package entity

import (
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProgress(t *testing.T) {
	t.Parallel()

	// Arrange
	materialID := valueobject.NewMaterialID()
	userID := valueobject.NewUserID()

	// Act
	progress := NewProgress(materialID, userID)

	// Assert
	assert.NotNil(t, progress)
	assert.Equal(t, materialID, progress.MaterialID())
	assert.Equal(t, userID, progress.UserID())
	assert.Equal(t, 0, progress.Percentage(), "New progress should start at 0%")
	assert.Equal(t, 0, progress.LastPage(), "New progress should start at page 0")
	assert.Equal(t, enum.ProgressStatusNotStarted, progress.Status())
	assert.False(t, progress.LastAccessedAt().IsZero())
	assert.False(t, progress.CreatedAt().IsZero())
	assert.False(t, progress.UpdatedAt().IsZero())
}

func TestReconstructProgress(t *testing.T) {
	t.Parallel()

	// Arrange
	materialID := valueobject.NewMaterialID()
	userID := valueobject.NewUserID()
	lastAccessed := time.Now().Add(-1 * time.Hour)
	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()

	// Act
	progress := ReconstructProgress(
		materialID,
		userID,
		75,
		10,
		enum.ProgressStatusInProgress,
		lastAccessed,
		createdAt,
		updatedAt,
	)

	// Assert
	assert.NotNil(t, progress)
	assert.Equal(t, materialID, progress.MaterialID())
	assert.Equal(t, userID, progress.UserID())
	assert.Equal(t, 75, progress.Percentage())
	assert.Equal(t, 10, progress.LastPage())
	assert.Equal(t, enum.ProgressStatusInProgress, progress.Status())
	assert.Equal(t, lastAccessed, progress.LastAccessedAt())
	assert.Equal(t, createdAt, progress.CreatedAt())
	assert.Equal(t, updatedAt, progress.UpdatedAt())
}

func TestProgress_UpdateProgress_ValidPercentages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		percentage     int
		expectedStatus enum.ProgressStatus
		shouldBeValid  bool
	}{
		{
			name:           "0% es válido (not started)",
			percentage:     0,
			expectedStatus: enum.ProgressStatusNotStarted,
			shouldBeValid:  true,
		},
		{
			name:           "1% es válido (in progress)",
			percentage:     1,
			expectedStatus: enum.ProgressStatusInProgress,
			shouldBeValid:  true,
		},
		{
			name:           "50% es válido (in progress)",
			percentage:     50,
			expectedStatus: enum.ProgressStatusInProgress,
			shouldBeValid:  true,
		},
		{
			name:           "99% es válido (in progress)",
			percentage:     99,
			expectedStatus: enum.ProgressStatusInProgress,
			shouldBeValid:  true,
		},
		{
			name:           "100% es válido (completed)",
			percentage:     100,
			expectedStatus: enum.ProgressStatusCompleted,
			shouldBeValid:  true,
		},
		{
			name:           "-1% es inválido",
			percentage:     -1,
			expectedStatus: enum.ProgressStatusNotStarted,
			shouldBeValid:  false,
		},
		{
			name:           "101% es inválido",
			percentage:     101,
			expectedStatus: enum.ProgressStatusNotStarted,
			shouldBeValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			materialID := valueobject.NewMaterialID()
			userID := valueobject.NewUserID()
			progress := NewProgress(materialID, userID)

			// Act
			err := progress.UpdateProgress(tt.percentage, 1)

			// Assert
			if tt.shouldBeValid {
				require.NoError(t, err, "UpdateProgress should not return error for valid percentage")
				assert.Equal(t, tt.percentage, progress.Percentage())
				assert.Equal(t, tt.expectedStatus, progress.Status())
			} else {
				require.Error(t, err, "UpdateProgress should return error for invalid percentage")
				assert.Contains(t, err.Error(), "percentage must be between 0 and 100")
			}
		})
	}
}

func TestProgress_Getters(t *testing.T) {
	t.Parallel()

	// Arrange
	materialID := valueobject.NewMaterialID()
	userID := valueobject.NewUserID()
	lastAccessed := time.Now().Add(-2 * time.Hour)
	createdAt := time.Now().Add(-48 * time.Hour)
	updatedAt := time.Now().Add(-1 * time.Hour)

	progress := ReconstructProgress(
		materialID,
		userID,
		80,
		15,
		enum.ProgressStatusInProgress,
		lastAccessed,
		createdAt,
		updatedAt,
	)

	// Act & Assert
	assert.Equal(t, materialID, progress.MaterialID())
	assert.Equal(t, userID, progress.UserID())
	assert.Equal(t, 80, progress.Percentage())
	assert.Equal(t, 15, progress.LastPage())
	assert.Equal(t, enum.ProgressStatusInProgress, progress.Status())
	assert.Equal(t, lastAccessed, progress.LastAccessedAt())
	assert.Equal(t, createdAt, progress.CreatedAt())
	assert.Equal(t, updatedAt, progress.UpdatedAt())
}
