// TODO: Estos tests unitarios requieren actualización para usar mocks reales (mongomock)
// Los tests de integración en assessment_document_repository_integration_test.go
// validan el funcionamiento real con testcontainers

//go:build skip_for_now
// +build skip_for_now

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewMongoAssessmentDocumentRepository(t *testing.T) {
	// Mock MongoDB Database
	var mockDB *mongo.Database

	repo := NewMongoAssessmentDocumentRepository(mockDB)

	assert.NotNil(t, repo)
	assert.IsType(t, &MongoAssessmentDocumentRepository{}, repo)
}

func TestMongoAssessmentDocumentRepository_FindByMaterialID_Success(t *testing.T) {
	// Arrange
	var mockDB *mongo.Database
	repo := NewMongoAssessmentDocumentRepository(mockDB)

	ctx := context.Background()
	materialID := "01936d9a-7f8e-7e4c-9d3f-987654321cba"

	// Act
	result, err := repo.FindByMaterialID(ctx, materialID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, materialID, result.MaterialID)
	assert.NotEmpty(t, result.Title)
	assert.NotEmpty(t, result.Questions)
	assert.Greater(t, len(result.Questions), 0)
}

func TestMongoAssessmentDocumentRepository_FindByMaterialID_EmptyID(t *testing.T) {
	// Arrange
	var mockDB *mongo.Database
	repo := NewMongoAssessmentDocumentRepository(mockDB)

	ctx := context.Background()

	// Act
	result, err := repo.FindByMaterialID(ctx, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestMongoAssessmentDocumentRepository_FindByID_Success(t *testing.T) {
	// Arrange
	var mockDB *mongo.Database
	repo := NewMongoAssessmentDocumentRepository(mockDB)

	ctx := context.Background()
	objectID := "507f1f77bcf86cd799439011" // ObjectID válido

	// Act
	result, err := repo.FindByID(ctx, objectID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.MaterialID)
	assert.NotEmpty(t, result.Questions)
}

func TestMongoAssessmentDocumentRepository_FindByID_InvalidObjectID(t *testing.T) {
	// Arrange
	var mockDB *mongo.Database
	repo := NewMongoAssessmentDocumentRepository(mockDB)

	ctx := context.Background()

	// Act
	result, err := repo.FindByID(ctx, "invalid-object-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid ObjectID")
}

func TestMongoAssessmentDocumentRepository_FindByID_EmptyID(t *testing.T) {
	// Arrange
	var mockDB *mongo.Database
	repo := NewMongoAssessmentDocumentRepository(mockDB)

	ctx := context.Background()

	// Act
	result, err := repo.FindByID(ctx, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestMongoAssessmentDocumentRepository_Save_Success(t *testing.T) {
	// Arrange
	var mockDB *mongo.Database
	repo := NewMongoAssessmentDocumentRepository(mockDB)

	ctx := context.Background()

	doc := &AssessmentDocument{
		MaterialID: "01936d9a-7f8e-7e4c-9d3f-987654321cba",
		Title:      "Test Assessment",
		Questions: []Question{
			{
				ID:   "q1",
				Text: "Test question?",
				Type: "multiple_choice",
				Options: []Option{
					{ID: "a", Text: "Option A"},
					{ID: "b", Text: "Option B"},
				},
				CorrectAnswer: "a",
				Feedback: Feedback{
					Correct:   "Correct!",
					Incorrect: "Try again.",
				},
			},
		},
		Metadata: Metadata{
			GeneratedBy:          "test",
			TotalQuestions:       1,
			EstimatedTimeMinutes: 5,
		},
		Version:   1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Act
	err := repo.Save(ctx, doc)

	// Assert
	assert.NoError(t, err)
}

func TestMongoAssessmentDocumentRepository_Save_NilDocument(t *testing.T) {
	// Arrange
	var mockDB *mongo.Database
	repo := NewMongoAssessmentDocumentRepository(mockDB)

	ctx := context.Background()

	// Act
	err := repo.Save(ctx, nil)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestMongoAssessmentDocumentRepository_Save_EmptyMaterialID(t *testing.T) {
	// Arrange
	var mockDB *mongo.Database
	repo := NewMongoAssessmentDocumentRepository(mockDB)

	ctx := context.Background()

	doc := &AssessmentDocument{
		MaterialID: "", // Vacío!
		Title:      "Test",
		Questions:  []Question{{ID: "q1"}},
	}

	// Act
	err := repo.Save(ctx, doc)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "material_id is required")
}

func TestMongoAssessmentDocumentRepository_Save_EmptyTitle(t *testing.T) {
	// Arrange
	var mockDB *mongo.Database
	repo := NewMongoAssessmentDocumentRepository(mockDB)

	ctx := context.Background()

	doc := &AssessmentDocument{
		MaterialID: "01936d9a-7f8e-7e4c-9d3f-987654321cba",
		Title:      "", // Vacío!
		Questions:  []Question{{ID: "q1"}},
	}

	// Act
	err := repo.Save(ctx, doc)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")
}

func TestMongoAssessmentDocumentRepository_Save_NoQuestions(t *testing.T) {
	// Arrange
	var mockDB *mongo.Database
	repo := NewMongoAssessmentDocumentRepository(mockDB)

	ctx := context.Background()

	doc := &AssessmentDocument{
		MaterialID: "01936d9a-7f8e-7e4c-9d3f-987654321cba",
		Title:      "Test",
		Questions:  []Question{}, // Vacío!
	}

	// Act
	err := repo.Save(ctx, doc)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one question is required")
}

func TestMongoAssessmentDocumentRepository_Delete_Success(t *testing.T) {
	// Arrange
	var mockDB *mongo.Database
	repo := NewMongoAssessmentDocumentRepository(mockDB)

	ctx := context.Background()
	objectID := "507f1f77bcf86cd799439011"

	// Act
	err := repo.Delete(ctx, objectID)

	// Assert
	assert.NoError(t, err)
}

func TestMongoAssessmentDocumentRepository_Delete_InvalidObjectID(t *testing.T) {
	// Arrange
	var mockDB *mongo.Database
	repo := NewMongoAssessmentDocumentRepository(mockDB)

	ctx := context.Background()

	// Act
	err := repo.Delete(ctx, "invalid-id")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid ObjectID")
}

func TestMongoAssessmentDocumentRepository_Delete_EmptyID(t *testing.T) {
	// Arrange
	var mockDB *mongo.Database
	repo := NewMongoAssessmentDocumentRepository(mockDB)

	ctx := context.Background()

	// Act
	err := repo.Delete(ctx, "")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestMongoAssessmentDocumentRepository_GetQuestionByID_Success(t *testing.T) {
	// Arrange
	var mockDB *mongo.Database
	repo := NewMongoAssessmentDocumentRepository(mockDB).(*MongoAssessmentDocumentRepository)

	ctx := context.Background()
	materialID := "01936d9a-7f8e-7e4c-9d3f-987654321cba"
	questionID := "q1"

	// Act
	result, err := repo.GetQuestionByID(ctx, materialID, questionID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, questionID, result.ID)
	assert.NotEmpty(t, result.Text)
	assert.NotEmpty(t, result.Options)
}

func TestMongoAssessmentDocumentRepository_GetQuestionByID_EmptyIDs(t *testing.T) {
	// Arrange
	var mockDB *mongo.Database
	repo := NewMongoAssessmentDocumentRepository(mockDB).(*MongoAssessmentDocumentRepository)

	ctx := context.Background()

	// Act
	result, err := repo.GetQuestionByID(ctx, "", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "are required")
}

func TestAssessmentDocument_Validation(t *testing.T) {
	// Test que verifica la estructura del documento
	doc := &AssessmentDocument{
		ID:         primitive.NewObjectID(),
		MaterialID: "01936d9a-7f8e-7e4c-9d3f-987654321cba",
		Title:      "Test Assessment",
		Questions: []Question{
			{
				ID:            "q1",
				Text:          "What is Go?",
				Type:          "multiple_choice",
				Options:       []Option{{ID: "a", Text: "Language"}},
				CorrectAnswer: "a",
				Feedback:      Feedback{Correct: "Yes", Incorrect: "No"},
			},
		},
		Metadata: Metadata{
			GeneratedBy:          "openai-gpt4",
			TotalQuestions:       1,
			EstimatedTimeMinutes: 5,
		},
		Version:   1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Verificar campos requeridos
	assert.NotEmpty(t, doc.MaterialID)
	assert.NotEmpty(t, doc.Title)
	assert.NotEmpty(t, doc.Questions)
	assert.Greater(t, doc.Version, 0)
}

func TestQuestion_Validation(t *testing.T) {
	// Test que verifica la estructura de una pregunta
	question := Question{
		ID:   "q1",
		Text: "Sample question?",
		Type: "multiple_choice",
		Options: []Option{
			{ID: "a", Text: "Option A"},
			{ID: "b", Text: "Option B"},
		},
		CorrectAnswer: "a",
		Feedback: Feedback{
			Correct:   "Correct!",
			Incorrect: "Try again.",
		},
	}

	// Verificar campos requeridos
	assert.NotEmpty(t, question.ID)
	assert.NotEmpty(t, question.Text)
	assert.NotEmpty(t, question.Type)
	assert.NotEmpty(t, question.Options)
	assert.NotEmpty(t, question.CorrectAnswer)
	assert.NotEmpty(t, question.Feedback.Correct)
	assert.NotEmpty(t, question.Feedback.Incorrect)
}
