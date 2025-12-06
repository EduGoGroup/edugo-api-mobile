package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests de SubmitAssessment (legacy) fueron eliminados
// El endpoint POST /assessments/:id/submit fue removido
// Usar POST /materials/:id/assessment/attempts en su lugar

// TestNewAssessmentHandler verifica la creación correcta del handler
func TestNewAssessmentHandler(t *testing.T) {
	// Arrange
	mockService := &MockAssessmentService{}
	logger := NewTestLogger()

	// Act
	mockAttemptService := &MockAssessmentAttemptService{}
	handler := NewAssessmentHandler(mockService, mockAttemptService, logger)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.assessmentService)
	assert.Equal(t, logger, handler.logger)
}
