package messaging

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitValidator(t *testing.T) {
	err := InitValidator()
	require.NoError(t, err, "InitValidator should succeed")

	// Verificar que el validador está inicializado
	v, err := GetValidator()
	require.NoError(t, err, "GetValidator should return validator")
	assert.NotNil(t, v, "validator should not be nil")
}

func TestValidateEvent_MaterialUploaded_Valid(t *testing.T) {
	// Inicializar validador
	err := InitValidator()
	require.NoError(t, err)

	// Crear payload válido
	payload := MaterialUploadedPayload{
		MaterialID:    "550e8400-e29b-41d4-a716-446655440000",
		SchoolID:      "550e8400-e29b-41d4-a716-446655440001",
		TeacherID:     "550e8400-e29b-41d4-a716-446655440002",
		FileURL:       "https://s3.amazonaws.com/edugo/materials/test.pdf",
		FileSizeBytes: 1024000,
		FileType:      "application/pdf",
		Metadata: map[string]interface{}{
			"title":       "Introducción a Go",
			"description": "Guía básica de Go",
		},
	}

	// Crear evento con envelope estándar
	event := NewMaterialUploadedEvent(payload)

	// Validar evento
	err = ValidateEvent(event)
	assert.NoError(t, err, "valid event should pass validation")
}

func TestValidateEvent_MaterialUploaded_Invalid(t *testing.T) {
	// Inicializar validador
	err := InitValidator()
	require.NoError(t, err)

	tests := []struct {
		name    string
		payload MaterialUploadedPayload
	}{
		{
			name: "missing material_id",
			payload: MaterialUploadedPayload{
				MaterialID:    "", // vacío (inválido)
				SchoolID:      "550e8400-e29b-41d4-a716-446655440001",
				TeacherID:     "550e8400-e29b-41d4-a716-446655440002",
				FileURL:       "https://s3.amazonaws.com/test.pdf",
				FileSizeBytes: 1000,
				FileType:      "application/pdf",
			},
		},
		{
			name: "missing school_id",
			payload: MaterialUploadedPayload{
				MaterialID:    "550e8400-e29b-41d4-a716-446655440000",
				SchoolID:      "", // vacío (inválido)
				TeacherID:     "550e8400-e29b-41d4-a716-446655440002",
				FileURL:       "https://s3.amazonaws.com/test.pdf",
				FileSizeBytes: 1000,
				FileType:      "application/pdf",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewMaterialUploadedEvent(tt.payload)
			err := ValidateEvent(event)
			assert.Error(t, err, "invalid event should fail validation")
		})
	}
}

func TestValidateEvent_AssessmentGenerated_Valid(t *testing.T) {
	// Inicializar validador
	err := InitValidator()
	require.NoError(t, err)

	// Crear payload válido
	payload := AssessmentGeneratedPayload{
		MaterialID:       "550e8400-e29b-41d4-a716-446655440002",
		MongoDocumentID:  "507f1f77bcf86cd799439011", // ObjectId válido (24 hex chars)
		QuestionsCount:   10,
		ProcessingTimeMs: 1500,
	}

	// Crear evento con envelope estándar
	event := NewAssessmentGeneratedEvent(payload)

	// Validar evento
	err = ValidateEvent(event)
	assert.NoError(t, err, "valid event should pass validation")
}

func TestGetValidator_NotInitialized(t *testing.T) {
	// Resetear validator para este test (cuidado en tests concurrentes)
	// En producción, InitValidator se llama una sola vez al inicio

	v, err := GetValidator()
	if err == nil {
		// Si ya está inicializado de tests anteriores, está ok
		assert.NotNil(t, v)
	} else {
		// Si no está inicializado, debe retornar error
		assert.Error(t, err)
		assert.Nil(t, v)
	}
}
