package entity

import (
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMaterial_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	authorID := valueobject.NewUserID()
	title := "Introducción a Go"
	description := "Material de aprendizaje sobre Go"
	subjectID := "subject-123"

	// Act
	material, err := NewMaterial(title, description, authorID, subjectID)

	// Assert
	require.NoError(t, err, "NewMaterial should not return error for valid data")
	assert.NotNil(t, material)
	assert.False(t, material.ID().IsZero(), "Material should have a valid ID")
	assert.Equal(t, title, material.Title())
	assert.Equal(t, description, material.Description())
	assert.Equal(t, authorID, material.AuthorID())
	assert.Equal(t, subjectID, material.SubjectID())
	assert.Equal(t, enum.MaterialStatusDraft, material.Status())
	assert.Equal(t, enum.ProcessingStatusPending, material.ProcessingStatus())
	assert.False(t, material.CreatedAt().IsZero())
	assert.False(t, material.UpdatedAt().IsZero())
}

func TestNewMaterial_ValidationErrors(t *testing.T) {
	t.Parallel()

	authorID := valueobject.NewUserID()

	tests := []struct {
		name        string
		title       string
		description string
		authorID    valueobject.UserID
		subjectID   string
		wantErr     string
	}{
		{
			name:        "título vacío debe fallar",
			title:       "",
			description: "Description",
			authorID:    authorID,
			subjectID:   "subject-1",
			wantErr:     "title is required",
		},
		{
			name:        "author_id zero debe fallar",
			title:       "Title",
			description: "Description",
			authorID:    valueobject.UserID{},
			subjectID:   "subject-1",
			wantErr:     "author_id is required",
		},
		{
			name:        "descripción vacía es permitida",
			title:       "Title",
			description: "",
			authorID:    authorID,
			subjectID:   "subject-1",
			wantErr:     "", // No debe dar error
		},
		{
			name:        "subjectID vacío es permitido",
			title:       "Title",
			description: "Description",
			authorID:    authorID,
			subjectID:   "",
			wantErr:     "", // No debe dar error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			material, err := NewMaterial(tt.title, tt.description, tt.authorID, tt.subjectID)

			// Assert
			if tt.wantErr != "" {
				require.Error(t, err, "NewMaterial should return error")
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, material)
			} else {
				require.NoError(t, err, "NewMaterial should not return error")
				assert.NotNil(t, material)
			}
		})
	}
}

func TestReconstructMaterial(t *testing.T) {
	t.Parallel()

	// Arrange
	materialID := valueobject.NewMaterialID()
	authorID := valueobject.NewUserID()
	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()

	// Act
	material := ReconstructMaterial(
		materialID,
		"Test Material",
		"Test Description",
		authorID,
		"subject-123",
		"s3/path/file.pdf",
		"https://s3.amazonaws.com/file.pdf",
		enum.MaterialStatusPublished,
		enum.ProcessingStatusCompleted,
		createdAt,
		updatedAt,
	)

	// Assert
	assert.NotNil(t, material)
	assert.Equal(t, materialID, material.ID())
	assert.Equal(t, "Test Material", material.Title())
	assert.Equal(t, "Test Description", material.Description())
	assert.Equal(t, authorID, material.AuthorID())
	assert.Equal(t, "subject-123", material.SubjectID())
	assert.Equal(t, "s3/path/file.pdf", material.S3Key())
	assert.Equal(t, "https://s3.amazonaws.com/file.pdf", material.S3URL())
	assert.Equal(t, enum.MaterialStatusPublished, material.Status())
	assert.Equal(t, enum.ProcessingStatusCompleted, material.ProcessingStatus())
	assert.Equal(t, createdAt, material.CreatedAt())
	assert.Equal(t, updatedAt, material.UpdatedAt())
}

func TestMaterial_InitialState(t *testing.T) {
	t.Parallel()

	// Arrange
	authorID := valueobject.NewUserID()

	// Act
	material, err := NewMaterial("Title", "Description", authorID, "subject-1")
	require.NoError(t, err)

	// Assert
	assert.Equal(t, enum.MaterialStatusDraft, material.Status(), "New material should be in draft status")
	assert.Equal(t, enum.ProcessingStatusPending, material.ProcessingStatus(), "New material should have pending processing status")
	assert.Empty(t, material.S3Key(), "New material should not have S3 key")
	assert.Empty(t, material.S3URL(), "New material should not have S3 URL")
}
