package rabbitmq

import (
	"testing"

	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMaterialUploadedEvent_ToJSON valida la serialización del evento MaterialUploadedEvent
func TestMaterialUploadedEvent_ToJSON(t *testing.T) {
	// Arrange
	event := messaging.MaterialUploadedEvent{
		MaterialID:  "mat_123",
		Title:       "Introducción a Go",
		ContentType: "application/pdf",
		UploadedAt:  time.Date(2024, 11, 4, 10, 0, 0, 0, time.UTC),
	}

	// Act
	jsonData, err := event.ToJSON()

	// Assert
	require.NoError(t, err, "ToJSON should not return error")
	assert.NotEmpty(t, jsonData, "JSON data should not be empty")
	assert.Contains(t, string(jsonData), "mat_123", "JSON should contain material_id")
	assert.Contains(t, string(jsonData), "Introducción a Go", "JSON should contain title")
	assert.Contains(t, string(jsonData), "application/pdf", "JSON should contain content_type")
}

// TestMaterialUploadedEvent_ToJSON_EmptyFields valida serialización con campos vacíos
func TestMaterialUploadedEvent_ToJSON_EmptyFields(t *testing.T) {
	// Arrange
	event := messaging.MaterialUploadedEvent{}

	// Act
	jsonData, err := event.ToJSON()

	// Assert
	require.NoError(t, err, "ToJSON should handle empty fields")
	assert.NotEmpty(t, jsonData, "JSON data should not be empty")
	assert.Contains(t, string(jsonData), `"material_id":""`, "JSON should contain empty material_id")
}

// TestAssessmentAttemptRecordedEvent_ToJSON valida la serialización del evento AssessmentAttemptRecordedEvent
func TestAssessmentAttemptRecordedEvent_ToJSON(t *testing.T) {
	// Arrange
	event := messaging.AssessmentAttemptRecordedEvent{
		AttemptID:    "attempt_456",
		UserID:       "user_789",
		AssessmentID: "assess_101",
		Score:        85.5,
		SubmittedAt:  time.Date(2024, 11, 4, 11, 30, 0, 0, time.UTC),
	}

	// Act
	jsonData, err := event.ToJSON()

	// Assert
	require.NoError(t, err, "ToJSON should not return error")
	assert.NotEmpty(t, jsonData, "JSON data should not be empty")
	assert.Contains(t, string(jsonData), "attempt_456", "JSON should contain attempt_id")
	assert.Contains(t, string(jsonData), "user_789", "JSON should contain user_id")
	assert.Contains(t, string(jsonData), "assess_101", "JSON should contain assessment_id")
	assert.Contains(t, string(jsonData), "85.5", "JSON should contain score")
}

// TestAssessmentAttemptRecordedEvent_ToJSON_ZeroScore valida serialización con score cero
func TestAssessmentAttemptRecordedEvent_ToJSON_ZeroScore(t *testing.T) {
	// Arrange
	event := messaging.AssessmentAttemptRecordedEvent{
		AttemptID:    "attempt_999",
		UserID:       "user_888",
		AssessmentID: "assess_777",
		Score:        0.0,
		SubmittedAt:  time.Now(),
	}

	// Act
	jsonData, err := event.ToJSON()

	// Assert
	require.NoError(t, err, "ToJSON should handle zero score")
	assert.NotEmpty(t, jsonData, "JSON data should not be empty")
	assert.Contains(t, string(jsonData), "attempt_999", "JSON should contain attempt_id")
	assert.Contains(t, string(jsonData), `"score":0`, "JSON should contain zero score")
}

// TestPublisher_Interface valida que RabbitMQPublisher implementa la interfaz Publisher
func TestPublisher_Interface(t *testing.T) {
	// Arrange & Act
	var _ Publisher = (*RabbitMQPublisher)(nil)

	// Assert - si compila, significa que implementa la interfaz correctamente
	assert.True(t, true, "RabbitMQPublisher implements Publisher interface")
}
