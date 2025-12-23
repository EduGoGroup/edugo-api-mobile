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
	mockAttemptService := &MockAssessmentAttemptService{}
	logger := NewTestLogger()

	// Act
	handler := NewAssessmentHandler(mockAttemptService, logger)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, mockAttemptService, handler.assessmentAttemptService)
	assert.Equal(t, logger, handler.logger)
}
