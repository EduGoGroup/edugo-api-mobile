package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestNewAssessmentHandler verifica la creación correcta del handler
func TestNewAssessmentHandler(t *testing.T) {
	// Arrange
	mockService := &MockAssessmentService{}
	logger := NewTestLogger()

	// Act
	handler := NewAssessmentHandler(mockService, logger)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.assessmentService)
	assert.Equal(t, logger, handler.logger)
}

// TestAssessmentHandler_SubmitAssessment_Success verifica envío exitoso con todas las respuestas correctas
func TestAssessmentHandler_SubmitAssessment_Success(t *testing.T) {
	// Arrange
	assessmentID := "assessment-123"
	userID := "user-456"

	mockService := &MockAssessmentService{
		CalculateScoreFunc: func(ctx context.Context, assID string, usrID string, responses map[string]interface{}) (*repository.AssessmentResult, error) {
			assert.Equal(t, assessmentID, assID)
			assert.Equal(t, userID, usrID)
			assert.NotNil(t, responses)

			return &repository.AssessmentResult{
				ID:             "result-789",
				AssessmentID:   assessmentID,
				Score:          100.0,
				TotalQuestions: 3,
				CorrectAnswers: 3,
				Feedback: []repository.FeedbackItem{
					{
						QuestionID:    "q1",
						IsCorrect:     true,
						UserAnswer:    "B",
						CorrectAnswer: "B",
						Explanation:   "Correcto",
					},
					{
						QuestionID:    "q2",
						IsCorrect:     true,
						UserAnswer:    "true",
						CorrectAnswer: "true",
						Explanation:   "Correcto",
					},
					{
						QuestionID:    "q3",
						IsCorrect:     true,
						UserAnswer:    "París",
						CorrectAnswer: "París",
						Explanation:   "Correcto",
					},
				},
				SubmittedAt: "2025-11-05T20:00:00Z",
			}, nil
		},
	}

	logger := NewTestLogger()
	handler := NewAssessmentHandler(mockService, logger)

	// Preparar request body
	requestBody := SubmitAssessmentRequest{
		Responses: map[string]interface{}{
			"q1": "B",
			"q2": "true",
			"q3": "París",
		},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	// Crear request HTTP
	req, _ := http.NewRequest("POST", "/v1/assessments/"+assessmentID+"/submit", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Crear response recorder y contexto Gin
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: assessmentID}}

	// Simular contexto de JWT con userID (normalmente lo agrega el middleware)
	c.Set("user_id", userID)

	// Act
	handler.SubmitAssessment(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response repository.AssessmentResult
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "result-789", response.ID)
	assert.Equal(t, assessmentID, response.AssessmentID)
	assert.Equal(t, 100.0, response.Score)
	assert.Equal(t, 3, response.TotalQuestions)
	assert.Equal(t, 3, response.CorrectAnswers)
	assert.Len(t, response.Feedback, 3)
}

// TestAssessmentHandler_SubmitAssessment_PartialCorrect verifica envío con respuestas parcialmente correctas
func TestAssessmentHandler_SubmitAssessment_PartialCorrect(t *testing.T) {
	// Arrange
	assessmentID := "assessment-123"
	userID := "user-456"

	mockService := &MockAssessmentService{
		CalculateScoreFunc: func(ctx context.Context, assID string, usrID string, responses map[string]interface{}) (*repository.AssessmentResult, error) {
			return &repository.AssessmentResult{
				ID:             "result-789",
				AssessmentID:   assessmentID,
				Score:          50.0,
				TotalQuestions: 4,
				CorrectAnswers: 2,
				Feedback: []repository.FeedbackItem{
					{
						QuestionID:    "q1",
						IsCorrect:     true,
						UserAnswer:    "A",
						CorrectAnswer: "A",
						Explanation:   "Correcto",
					},
					{
						QuestionID:    "q2",
						IsCorrect:     false,
						UserAnswer:    "false",
						CorrectAnswer: "true",
						Explanation:   "Incorrecto. La respuesta correcta es true",
					},
					{
						QuestionID:    "q3",
						IsCorrect:     true,
						UserAnswer:    "París",
						CorrectAnswer: "París",
						Explanation:   "Correcto",
					},
					{
						QuestionID:    "q4",
						IsCorrect:     false,
						UserAnswer:    "Londres",
						CorrectAnswer: "Madrid",
						Explanation:   "Incorrecto. La respuesta correcta es Madrid",
					},
				},
				SubmittedAt: "2025-11-05T20:00:00Z",
			}, nil
		},
	}

	logger := NewTestLogger()
	handler := NewAssessmentHandler(mockService, logger)

	// Preparar request body
	requestBody := SubmitAssessmentRequest{
		Responses: map[string]interface{}{
			"q1": "A",
			"q2": "false",
			"q3": "París",
			"q4": "Londres",
		},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	// Crear request HTTP
	req, _ := http.NewRequest("POST", "/v1/assessments/"+assessmentID+"/submit", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Crear response recorder y contexto Gin
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: assessmentID}}
	c.Set("user_id", userID)

	// Act
	handler.SubmitAssessment(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response repository.AssessmentResult
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 50.0, response.Score)
	assert.Equal(t, 2, response.CorrectAnswers)
	assert.Equal(t, 4, response.TotalQuestions)
	assert.Len(t, response.Feedback, 4)
}

// TestAssessmentHandler_SubmitAssessment_InvalidRequest verifica error con body inválido
func TestAssessmentHandler_SubmitAssessment_InvalidRequest(t *testing.T) {
	// Arrange
	assessmentID := "assessment-123"
	userID := "user-456"

	mockService := &MockAssessmentService{}
	logger := NewTestLogger()
	handler := NewAssessmentHandler(mockService, logger)

	// Body JSON inválido
	invalidJSON := []byte(`{"responses": "invalid"}`)

	req, _ := http.NewRequest("POST", "/v1/assessments/"+assessmentID+"/submit", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: assessmentID}}
	c.Set("user_id", userID)

	// Act
	handler.SubmitAssessment(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "INVALID_REQUEST", errorResponse.Code)
}

// TestAssessmentHandler_SubmitAssessment_EmptyResponses verifica error con responses vacías
func TestAssessmentHandler_SubmitAssessment_EmptyResponses(t *testing.T) {
	// Arrange
	assessmentID := "assessment-123"
	userID := "user-456"

	mockService := &MockAssessmentService{}
	logger := NewTestLogger()
	handler := NewAssessmentHandler(mockService, logger)

	// Body con responses vacío
	requestBody := SubmitAssessmentRequest{
		Responses: map[string]interface{}{},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/v1/assessments/"+assessmentID+"/submit", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: assessmentID}}
	c.Set("user_id", userID)

	// Act
	handler.SubmitAssessment(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "INVALID_REQUEST", errorResponse.Code)
	assert.Contains(t, errorResponse.Error, "responses map is required")
}

// TestAssessmentHandler_SubmitAssessment_AssessmentNotFound verifica error 404 cuando assessment no existe
func TestAssessmentHandler_SubmitAssessment_AssessmentNotFound(t *testing.T) {
	// Arrange
	assessmentID := "non-existent-assessment"
	userID := "user-456"

	mockService := &MockAssessmentService{
		CalculateScoreFunc: func(ctx context.Context, assID string, usrID string, responses map[string]interface{}) (*repository.AssessmentResult, error) {
			return nil, errors.NewNotFoundError("assessment")
		},
	}

	logger := NewTestLogger()
	handler := NewAssessmentHandler(mockService, logger)

	requestBody := SubmitAssessmentRequest{
		Responses: map[string]interface{}{
			"q1": "A",
		},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/v1/assessments/"+assessmentID+"/submit", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: assessmentID}}
	c.Set("user_id", userID)

	// Act
	handler.SubmitAssessment(c)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var errorResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "NOT_FOUND", errorResponse.Code)
}

// TestAssessmentHandler_SubmitAssessment_InvalidAssessmentID verifica error 400 con UUID inválido
func TestAssessmentHandler_SubmitAssessment_InvalidAssessmentID(t *testing.T) {
	// Arrange
	assessmentID := "invalid-uuid"
	userID := "user-456"

	mockService := &MockAssessmentService{
		CalculateScoreFunc: func(ctx context.Context, assID string, usrID string, responses map[string]interface{}) (*repository.AssessmentResult, error) {
			// Simular error de validación de MaterialID/UUID
			_, err := valueobject.MaterialIDFromString(assID)
			if err != nil {
				return nil, errors.NewValidationError("invalid assessment_id")
			}
			return nil, nil
		},
	}

	logger := NewTestLogger()
	handler := NewAssessmentHandler(mockService, logger)

	requestBody := SubmitAssessmentRequest{
		Responses: map[string]interface{}{
			"q1": "A",
		},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/v1/assessments/"+assessmentID+"/submit", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: assessmentID}}
	c.Set("user_id", userID)

	// Act
	handler.SubmitAssessment(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "VALIDATION_ERROR", errorResponse.Code)
}

// TestAssessmentHandler_SubmitAssessment_DatabaseError verifica error 500 con problema de BD
func TestAssessmentHandler_SubmitAssessment_DatabaseError(t *testing.T) {
	// Arrange
	assessmentID := "assessment-123"
	userID := "user-456"

	mockService := &MockAssessmentService{
		CalculateScoreFunc: func(ctx context.Context, assID string, usrID string, responses map[string]interface{}) (*repository.AssessmentResult, error) {
			// Retornar error de BD genérico (no relacionado con save)
			return nil, errors.NewDatabaseError("fetch assessment", assert.AnError)
		},
	}

	logger := NewTestLogger()
	handler := NewAssessmentHandler(mockService, logger)

	requestBody := SubmitAssessmentRequest{
		Responses: map[string]interface{}{
			"q1": "A",
		},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/v1/assessments/"+assessmentID+"/submit", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: assessmentID}}
	c.Set("user_id", userID)

	// Act
	handler.SubmitAssessment(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "DATABASE_ERROR", errorResponse.Code)
}

// TestAssessmentHandler_SubmitAssessment_AssessmentAlreadyCompleted verifica error 409 cuando evaluación ya fue completada
func TestAssessmentHandler_SubmitAssessment_AssessmentAlreadyCompleted(t *testing.T) {
	// Arrange
	assessmentID := "assessment-123"
	userID := "user-456"

	mockService := &MockAssessmentService{
		CalculateScoreFunc: func(ctx context.Context, assID string, usrID string, responses map[string]interface{}) (*repository.AssessmentResult, error) {
			// Simular error de duplicado al guardar resultado
			return nil, errors.NewDatabaseError("save assessment result", assert.AnError)
		},
	}

	logger := NewTestLogger()
	handler := NewAssessmentHandler(mockService, logger)

	requestBody := SubmitAssessmentRequest{
		Responses: map[string]interface{}{
			"q1": "A",
		},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/v1/assessments/"+assessmentID+"/submit", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: assessmentID}}
	c.Set("user_id", userID)

	// Act
	handler.SubmitAssessment(c)

	// Assert
	// Nota: El handler detecta "save assessment result" en DATABASE_ERROR y asume que es duplicado (409)
	assert.Equal(t, http.StatusConflict, w.Code)

	var errorResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "ASSESSMENT_ALREADY_COMPLETED", errorResponse.Code)
	assert.Contains(t, errorResponse.Error, "assessment already completed")
}
