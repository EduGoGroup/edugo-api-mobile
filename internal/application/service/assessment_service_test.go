package service

import (
	"context"
	"errors"
	"testing"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAssessmentRepository es un mock del repositorio de assessments
type MockAssessmentRepository struct {
	mock.Mock
}

func (m *MockAssessmentRepository) SaveAssessment(ctx context.Context, assessment *repository.MaterialAssessment) error {
	args := m.Called(ctx, assessment)
	return args.Error(0)
}

func (m *MockAssessmentRepository) FindAssessmentByMaterialID(ctx context.Context, materialID valueobject.MaterialID) (*repository.MaterialAssessment, error) {
	args := m.Called(ctx, materialID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.MaterialAssessment), args.Error(1)
}

func (m *MockAssessmentRepository) SaveAttempt(ctx context.Context, attempt *repository.AssessmentAttempt) error {
	args := m.Called(ctx, attempt)
	return args.Error(0)
}

func (m *MockAssessmentRepository) FindAttemptsByUser(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) ([]*repository.AssessmentAttempt, error) {
	args := m.Called(ctx, materialID, userID)
	return args.Get(0).([]*repository.AssessmentAttempt), args.Error(1)
}

func (m *MockAssessmentRepository) GetBestAttempt(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) (*repository.AssessmentAttempt, error) {
	args := m.Called(ctx, materialID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.AssessmentAttempt), args.Error(1)
}

func (m *MockAssessmentRepository) SaveResult(ctx context.Context, result *repository.AssessmentResult) error {
	args := m.Called(ctx, result)
	return args.Error(0)
}

func (m *MockAssessmentRepository) CountCompletedAssessments(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAssessmentRepository) CalculateAverageScore(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

// Nota: MockPublisher y MockLogger ya están definidos en material_service_test.go

func TestAssessmentService_CalculateScore_TodasRespuestasCorrectas(t *testing.T) {
	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	assessmentID := "550e8400-e29b-41d4-a716-446655440001"
	userID := "550e8400-e29b-41d4-a716-446655440002"
	matID, _ := valueobject.MaterialIDFromString(assessmentID)

	assessment := &repository.MaterialAssessment{
		MaterialID: matID,
		Questions: []repository.AssessmentQuestion{
			{
				ID:            "q1",
				QuestionText:  "¿Cuál es la capital de Francia?",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				Options:       []string{"A. Madrid", "B. París", "C. Londres"},
				CorrectAnswer: "B",
				Explanation:   "París es la capital de Francia.",
			},
			{
				ID:            "q2",
				QuestionText:  "¿El sol es una estrella?",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: "true",
				Explanation:   "El sol es una estrella.",
			},
		},
		CreatedAt: "2025-11-05T00:00:00Z",
	}

	userResponses := map[string]interface{}{
		"q1": "B",
		"q2": "true",
	}

	// Mock expectations
	mockRepo.On("FindAssessmentByMaterialID", mock.Anything, matID).Return(assessment, nil)
	mockRepo.On("SaveResult", mock.Anything, mock.AnythingOfType("*repository.AssessmentResult")).Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	// TODO(sprint-01): Restaurar mock cuando se implemente evento assessment.completed
	// mockPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Act
	result, err := service.CalculateScore(context.Background(), assessmentID, userID, userResponses)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 100.0, result.Score)
	assert.Equal(t, 2, result.TotalQuestions)
	assert.Equal(t, 2, result.CorrectAnswers)
	assert.Len(t, result.Feedback, 2)
	assert.True(t, result.Feedback[0].IsCorrect)
	assert.True(t, result.Feedback[1].IsCorrect)

	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestAssessmentService_CalculateScore_RespuestasParciales(t *testing.T) {
	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	assessmentID := "550e8400-e29b-41d4-a716-446655440001"
	userID := "550e8400-e29b-41d4-a716-446655440002"
	matID, _ := valueobject.MaterialIDFromString(assessmentID)

	assessment := &repository.MaterialAssessment{
		MaterialID: matID,
		Questions: []repository.AssessmentQuestion{
			{
				ID:            "q1",
				QuestionText:  "Pregunta 1",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				CorrectAnswer: "A",
				Explanation:   "",
			},
			{
				ID:            "q2",
				QuestionText:  "Pregunta 2",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: true,
				Explanation:   "",
			},
			{
				ID:            "q3",
				QuestionText:  "Pregunta 3",
				QuestionType:  enum.AssessmentTypeShortAnswer,
				CorrectAnswer: "París",
				Explanation:   "",
			},
			{
				ID:            "q4",
				QuestionText:  "Pregunta 4",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				CorrectAnswer: "D",
				Explanation:   "",
			},
		},
		CreatedAt: "2025-11-05T00:00:00Z",
	}

	userResponses := map[string]interface{}{
		"q1": "A",     // Correcta
		"q2": "false", // Incorrecta
		"q3": "París", // Correcta
		"q4": "B",     // Incorrecta
	}

	// Mock expectations
	mockRepo.On("FindAssessmentByMaterialID", mock.Anything, matID).Return(assessment, nil)
	mockRepo.On("SaveResult", mock.Anything, mock.AnythingOfType("*repository.AssessmentResult")).Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	// TODO(sprint-01): Restaurar mock cuando se implemente evento assessment.completed
	// mockPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Act
	result, err := service.CalculateScore(context.Background(), assessmentID, userID, userResponses)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 50.0, result.Score) // 2 de 4 correctas = 50%
	assert.Equal(t, 4, result.TotalQuestions)
	assert.Equal(t, 2, result.CorrectAnswers)
	assert.Len(t, result.Feedback, 4)
	assert.True(t, result.Feedback[0].IsCorrect)  // q1
	assert.False(t, result.Feedback[1].IsCorrect) // q2
	assert.True(t, result.Feedback[2].IsCorrect)  // q3
	assert.False(t, result.Feedback[3].IsCorrect) // q4

	mockRepo.AssertExpectations(t)
}

func TestAssessmentService_CalculateScore_NingunaCorrecta(t *testing.T) {
	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	assessmentID := "550e8400-e29b-41d4-a716-446655440001"
	userID := "550e8400-e29b-41d4-a716-446655440002"
	matID, _ := valueobject.MaterialIDFromString(assessmentID)

	assessment := &repository.MaterialAssessment{
		MaterialID: matID,
		Questions: []repository.AssessmentQuestion{
			{
				ID:            "q1",
				QuestionText:  "Pregunta 1",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				CorrectAnswer: "A",
				Explanation:   "",
			},
			{
				ID:            "q2",
				QuestionText:  "Pregunta 2",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: true,
				Explanation:   "",
			},
		},
		CreatedAt: "2025-11-05T00:00:00Z",
	}

	userResponses := map[string]interface{}{
		"q1": "B",     // Incorrecta
		"q2": "false", // Incorrecta
	}

	// Mock expectations
	mockRepo.On("FindAssessmentByMaterialID", mock.Anything, matID).Return(assessment, nil)
	mockRepo.On("SaveResult", mock.Anything, mock.AnythingOfType("*repository.AssessmentResult")).Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	// TODO(sprint-01): Restaurar mock cuando se implemente evento assessment.completed
	// mockPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Act
	result, err := service.CalculateScore(context.Background(), assessmentID, userID, userResponses)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0.0, result.Score)
	assert.Equal(t, 2, result.TotalQuestions)
	assert.Equal(t, 0, result.CorrectAnswers)
	assert.Len(t, result.Feedback, 2)
	assert.False(t, result.Feedback[0].IsCorrect)
	assert.False(t, result.Feedback[1].IsCorrect)

	mockRepo.AssertExpectations(t)
}

func TestAssessmentService_CalculateScore_PreguntaSinResponder(t *testing.T) {
	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	assessmentID := "550e8400-e29b-41d4-a716-446655440001"
	userID := "550e8400-e29b-41d4-a716-446655440002"
	matID, _ := valueobject.MaterialIDFromString(assessmentID)

	assessment := &repository.MaterialAssessment{
		MaterialID: matID,
		Questions: []repository.AssessmentQuestion{
			{
				ID:            "q1",
				QuestionText:  "Pregunta 1",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				CorrectAnswer: "A",
				Explanation:   "",
			},
			{
				ID:            "q2",
				QuestionText:  "Pregunta 2",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: true,
				Explanation:   "",
			},
		},
		CreatedAt: "2025-11-05T00:00:00Z",
	}

	userResponses := map[string]interface{}{
		"q1": "A", // Correcta
		// q2 no respondida
	}

	// Mock expectations
	mockRepo.On("FindAssessmentByMaterialID", mock.Anything, matID).Return(assessment, nil)
	mockRepo.On("SaveResult", mock.Anything, mock.AnythingOfType("*repository.AssessmentResult")).Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()
	// TODO(sprint-01): Restaurar mock cuando se implemente evento assessment.completed
	// mockPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Act
	result, err := service.CalculateScore(context.Background(), assessmentID, userID, userResponses)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 50.0, result.Score) // 1 de 2 correctas = 50%
	assert.Equal(t, 2, result.TotalQuestions)
	assert.Equal(t, 1, result.CorrectAnswers)
	assert.Len(t, result.Feedback, 2)
	assert.True(t, result.Feedback[0].IsCorrect)
	assert.False(t, result.Feedback[1].IsCorrect) // sin respuesta = incorrecta
	assert.Contains(t, result.Feedback[1].UserAnswer, "sin respuesta")

	mockRepo.AssertExpectations(t)
}

func TestAssessmentService_CalculateScore_AssessmentNoExiste(t *testing.T) {
	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	assessmentID := "550e8400-e29b-41d4-a716-446655440001"
	userID := "550e8400-e29b-41d4-a716-446655440002"
	matID, _ := valueobject.MaterialIDFromString(assessmentID)

	// Mock expectations - assessment no existe
	mockRepo.On("FindAssessmentByMaterialID", mock.Anything, matID).Return(nil, nil)
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	// Act
	result, err := service.CalculateScore(context.Background(), assessmentID, userID, map[string]interface{}{})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestAssessmentService_CalculateScore_ErrorBaseDatos(t *testing.T) {
	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	assessmentID := "550e8400-e29b-41d4-a716-446655440001"
	userID := "550e8400-e29b-41d4-a716-446655440002"
	matID, _ := valueobject.MaterialIDFromString(assessmentID)

	// Mock expectations - error de BD
	mockRepo.On("FindAssessmentByMaterialID", mock.Anything, matID).Return(nil, errors.New("database error"))
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	// Act
	result, err := service.CalculateScore(context.Background(), assessmentID, userID, map[string]interface{}{})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestAssessmentService_CalculateScore_IDInvalido(t *testing.T) {
	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	// Act - assessmentID inválido
	result, err := service.CalculateScore(context.Background(), "invalid-id", "valid-user-id", map[string]interface{}{})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============================================
// Tests: GetAssessment
// ============================================

func TestAssessmentService_GetAssessment_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	materialID := "550e8400-e29b-41d4-a716-446655440001"
	matID, _ := valueobject.MaterialIDFromString(materialID)

	expectedAssessment := &repository.MaterialAssessment{
		MaterialID: matID,
		Questions: []repository.AssessmentQuestion{
			{
				ID:            "q1",
				QuestionText:  "Test question",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				CorrectAnswer: "A",
			},
		},
		CreatedAt: "2025-11-05T00:00:00Z",
	}

	mockRepo.On("FindAssessmentByMaterialID", mock.Anything, matID).Return(expectedAssessment, nil)
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	// Act
	result, err := service.GetAssessment(context.Background(), materialID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, matID, result.MaterialID)
	assert.Len(t, result.Questions, 1)

	mockRepo.AssertExpectations(t)
}

func TestAssessmentService_GetAssessment_InvalidMaterialID(t *testing.T) {
	t.Parallel()

	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	// Act
	result, err := service.GetAssessment(context.Background(), "invalid-uuid")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid material_id")
}

func TestAssessmentService_GetAssessment_NotFound(t *testing.T) {
	t.Parallel()

	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	materialID := "550e8400-e29b-41d4-a716-446655440001"
	matID, _ := valueobject.MaterialIDFromString(materialID)

	mockRepo.On("FindAssessmentByMaterialID", mock.Anything, matID).Return(nil, nil)
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	// Act
	result, err := service.GetAssessment(context.Background(), materialID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")

	mockRepo.AssertExpectations(t)
}

func TestAssessmentService_GetAssessment_DatabaseError(t *testing.T) {
	t.Parallel()

	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	materialID := "550e8400-e29b-41d4-a716-446655440001"
	matID, _ := valueobject.MaterialIDFromString(materialID)

	mockRepo.On("FindAssessmentByMaterialID", mock.Anything, matID).Return(nil, errors.New("database error"))
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	// Act
	result, err := service.GetAssessment(context.Background(), materialID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

// ============================================
// Tests: RecordAttempt
// ============================================

func TestAssessmentService_RecordAttempt_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	materialID := "550e8400-e29b-41d4-a716-446655440001"
	userID := "550e8400-e29b-41d4-a716-446655440002"
	matID, _ := valueobject.MaterialIDFromString(materialID)

	assessment := &repository.MaterialAssessment{
		MaterialID: matID,
		Questions: []repository.AssessmentQuestion{
			{
				ID:            "q1",
				QuestionText:  "Test question",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				CorrectAnswer: "A",
			},
		},
		CreatedAt: "2025-11-05T00:00:00Z",
	}

	answers := map[string]interface{}{
		"q1": "A",
	}

	mockRepo.On("FindAssessmentByMaterialID", mock.Anything, matID).Return(assessment, nil)
	mockRepo.On("SaveAttempt", mock.Anything, mock.AnythingOfType("*repository.AssessmentAttempt")).Return(nil)
	mockPublisher.On("Publish", mock.Anything, "edugo.materials", "assessment.attempt.recorded", mock.Anything).Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()

	// Act
	result, err := service.RecordAttempt(context.Background(), materialID, userID, answers)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, matID, result.MaterialID)
	assert.Equal(t, 75.0, result.Score) // Mock score

	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestAssessmentService_RecordAttempt_InvalidMaterialID(t *testing.T) {
	t.Parallel()

	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	userID := "550e8400-e29b-41d4-a716-446655440002"
	answers := map[string]interface{}{"q1": "A"}

	// Act
	result, err := service.RecordAttempt(context.Background(), "invalid-uuid", userID, answers)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid material_id")
}

func TestAssessmentService_RecordAttempt_InvalidUserID(t *testing.T) {
	t.Parallel()

	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	materialID := "550e8400-e29b-41d4-a716-446655440001"
	answers := map[string]interface{}{"q1": "A"}

	// Act
	result, err := service.RecordAttempt(context.Background(), materialID, "invalid-uuid", answers)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid user_id")
}

func TestAssessmentService_RecordAttempt_AssessmentNotFound(t *testing.T) {
	t.Parallel()

	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	materialID := "550e8400-e29b-41d4-a716-446655440001"
	userID := "550e8400-e29b-41d4-a716-446655440002"
	matID, _ := valueobject.MaterialIDFromString(materialID)
	answers := map[string]interface{}{"q1": "A"}

	mockRepo.On("FindAssessmentByMaterialID", mock.Anything, matID).Return(nil, nil)

	// Act
	result, err := service.RecordAttempt(context.Background(), materialID, userID, answers)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")

	mockRepo.AssertExpectations(t)
}

func TestAssessmentService_RecordAttempt_SaveAttemptError(t *testing.T) {
	t.Parallel()

	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	materialID := "550e8400-e29b-41d4-a716-446655440001"
	userID := "550e8400-e29b-41d4-a716-446655440002"
	matID, _ := valueobject.MaterialIDFromString(materialID)

	assessment := &repository.MaterialAssessment{
		MaterialID: matID,
		Questions: []repository.AssessmentQuestion{
			{
				ID:            "q1",
				QuestionText:  "Test question",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				CorrectAnswer: "A",
			},
		},
		CreatedAt: "2025-11-05T00:00:00Z",
	}

	answers := map[string]interface{}{"q1": "A"}

	mockRepo.On("FindAssessmentByMaterialID", mock.Anything, matID).Return(assessment, nil)
	mockRepo.On("SaveAttempt", mock.Anything, mock.AnythingOfType("*repository.AssessmentAttempt")).Return(errors.New("database error"))
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	// Act
	result, err := service.RecordAttempt(context.Background(), materialID, userID, answers)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestAssessmentService_RecordAttempt_PublishEventFailure(t *testing.T) {
	t.Parallel()

	// Arrange
	mockRepo := new(MockAssessmentRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewAssessmentService(mockRepo, mockPublisher, mockLogger)

	materialID := "550e8400-e29b-41d4-a716-446655440001"
	userID := "550e8400-e29b-41d4-a716-446655440002"
	matID, _ := valueobject.MaterialIDFromString(materialID)

	assessment := &repository.MaterialAssessment{
		MaterialID: matID,
		Questions: []repository.AssessmentQuestion{
			{
				ID:            "q1",
				QuestionText:  "Test question",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				CorrectAnswer: "A",
			},
		},
		CreatedAt: "2025-11-05T00:00:00Z",
	}

	answers := map[string]interface{}{"q1": "A"}

	mockRepo.On("FindAssessmentByMaterialID", mock.Anything, matID).Return(assessment, nil)
	mockRepo.On("SaveAttempt", mock.Anything, mock.AnythingOfType("*repository.AssessmentAttempt")).Return(nil)
	mockPublisher.On("Publish", mock.Anything, "edugo.materials", "assessment.attempt.recorded", mock.Anything).Return(errors.New("publish error"))
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()

	// Act
	result, err := service.RecordAttempt(context.Background(), materialID, userID, answers)

	// Assert - Should succeed even if publish fails (non-blocking)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}
