package service

import (
	"context"
	"errors"
	"testing"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProgressRepository es un mock del repositorio de progreso
type MockProgressRepository struct {
	mock.Mock
}

func (m *MockProgressRepository) Save(ctx context.Context, progress *entity.Progress) error {
	args := m.Called(ctx, progress)
	return args.Error(0)
}

func (m *MockProgressRepository) FindByMaterialAndUser(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) (*entity.Progress, error) {
	args := m.Called(ctx, materialID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Progress), args.Error(1)
}

func (m *MockProgressRepository) Update(ctx context.Context, progress *entity.Progress) error {
	args := m.Called(ctx, progress)
	return args.Error(0)
}

func (m *MockProgressRepository) Upsert(ctx context.Context, progress *entity.Progress) (*entity.Progress, error) {
	args := m.Called(ctx, progress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Progress), args.Error(1)
}

func (m *MockProgressRepository) CountActiveUsers(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProgressRepository) CalculateAverageProgress(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

// MockLogger es un mock del logger
type MockProgressLogger struct {
	mock.Mock
}

func (m *MockProgressLogger) Info(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockProgressLogger) Warn(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockProgressLogger) Error(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockProgressLogger) Debug(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockProgressLogger) Fatal(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockProgressLogger) Sync() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockProgressLogger) With(keysAndValues ...interface{}) logger.Logger {
	m.Called(keysAndValues)
	return m
}

// TestUpdateProgress_Success_ValidProgress prueba actualización exitosa con progreso válido
func TestUpdateProgress_Success_ValidProgress(t *testing.T) {
	// Arrange
	mockRepo := new(MockProgressRepository)
	mockLogger := new(MockProgressLogger)
	service := NewProgressService(mockRepo, mockLogger)

	ctx := context.Background()
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	userID := "660e8400-e29b-41d4-a716-446655440001"
	percentage := 75
	lastPage := 10

	// Crear progress esperado
	matID, _ := valueobject.MaterialIDFromString(materialID)
	uID, _ := valueobject.UserIDFromString(userID)
	expectedProgress := entity.NewProgress(matID, uID)
	_ = expectedProgress.UpdateProgress(percentage, lastPage)

	// Mock expectations
	mockLogger.On("Info", "updating progress", mock.Anything).Return()
	mockRepo.On("Upsert", ctx, mock.AnythingOfType("*entity.Progress")).
		Return(expectedProgress, nil)
	mockLogger.On("Info", "progress updated successfully", mock.Anything).Return()

	// Act
	err := service.UpdateProgress(ctx, materialID, userID, percentage, lastPage)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// TestUpdateProgress_Success_CompletedMaterial prueba completar material (percentage = 100)
func TestUpdateProgress_Success_CompletedMaterial(t *testing.T) {
	// Arrange
	mockRepo := new(MockProgressRepository)
	mockLogger := new(MockProgressLogger)
	service := NewProgressService(mockRepo, mockLogger)

	ctx := context.Background()
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	userID := "660e8400-e29b-41d4-a716-446655440001"
	percentage := 100
	lastPage := 50

	// Crear progress completado
	matID, _ := valueobject.MaterialIDFromString(materialID)
	uID, _ := valueobject.UserIDFromString(userID)
	completedProgress := entity.NewProgress(matID, uID)
	_ = completedProgress.UpdateProgress(percentage, lastPage)

	// Mock expectations
	mockLogger.On("Info", "updating progress", mock.Anything).Return()
	mockRepo.On("Upsert", ctx, mock.AnythingOfType("*entity.Progress")).
		Return(completedProgress, nil)
	mockLogger.On("Info", "material completed by user", mock.Anything).Return()
	mockLogger.On("Info", "progress updated successfully", mock.Anything).Return()

	// Act
	err := service.UpdateProgress(ctx, materialID, userID, percentage, lastPage)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, enum.ProgressStatusCompleted, completedProgress.Status())
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// TestUpdateProgress_Error_InvalidPercentageNegative prueba error con percentage negativo
func TestUpdateProgress_Error_InvalidPercentageNegative(t *testing.T) {
	// Arrange
	mockRepo := new(MockProgressRepository)
	mockLogger := new(MockProgressLogger)
	service := NewProgressService(mockRepo, mockLogger)

	ctx := context.Background()
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	userID := "660e8400-e29b-41d4-a716-446655440001"
	percentage := -10
	lastPage := 5

	// Mock expectations
	mockLogger.On("Info", "updating progress", mock.Anything).Return()
	mockLogger.On("Warn", "invalid percentage value", mock.Anything).Return()

	// Act
	err := service.UpdateProgress(ctx, materialID, userID, percentage, lastPage)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "percentage must be between 0 and 100")
	mockLogger.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}

// TestUpdateProgress_Error_InvalidPercentageOver100 prueba error con percentage > 100
func TestUpdateProgress_Error_InvalidPercentageOver100(t *testing.T) {
	// Arrange
	mockRepo := new(MockProgressRepository)
	mockLogger := new(MockProgressLogger)
	service := NewProgressService(mockRepo, mockLogger)

	ctx := context.Background()
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	userID := "660e8400-e29b-41d4-a716-446655440001"
	percentage := 150
	lastPage := 5

	// Mock expectations
	mockLogger.On("Info", "updating progress", mock.Anything).Return()
	mockLogger.On("Warn", "invalid percentage value", mock.Anything).Return()

	// Act
	err := service.UpdateProgress(ctx, materialID, userID, percentage, lastPage)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "percentage must be between 0 and 100")
	mockLogger.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}

// TestUpdateProgress_Error_InvalidMaterialID prueba error con materialID inválido
func TestUpdateProgress_Error_InvalidMaterialID(t *testing.T) {
	// Arrange
	mockRepo := new(MockProgressRepository)
	mockLogger := new(MockProgressLogger)
	service := NewProgressService(mockRepo, mockLogger)

	ctx := context.Background()
	materialID := "invalid-uuid"
	userID := "660e8400-e29b-41d4-a716-446655440001"
	percentage := 50
	lastPage := 5

	// Mock expectations
	mockLogger.On("Info", "updating progress", mock.Anything).Return()
	mockLogger.On("Error", "invalid material_id", mock.Anything).Return()

	// Act
	err := service.UpdateProgress(ctx, materialID, userID, percentage, lastPage)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid material_id")
	mockLogger.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}

// TestUpdateProgress_Error_InvalidUserID prueba error con userID inválido
func TestUpdateProgress_Error_InvalidUserID(t *testing.T) {
	// Arrange
	mockRepo := new(MockProgressRepository)
	mockLogger := new(MockProgressLogger)
	service := NewProgressService(mockRepo, mockLogger)

	ctx := context.Background()
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	userID := "invalid-uuid"
	percentage := 50
	lastPage := 5

	// Mock expectations
	mockLogger.On("Info", "updating progress", mock.Anything).Return()
	mockLogger.On("Error", "invalid user_id", mock.Anything).Return()

	// Act
	err := service.UpdateProgress(ctx, materialID, userID, percentage, lastPage)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user_id")
	mockLogger.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}

// TestUpdateProgress_Error_DatabaseError prueba error de base de datos durante UPSERT
func TestUpdateProgress_Error_DatabaseError(t *testing.T) {
	// Arrange
	mockRepo := new(MockProgressRepository)
	mockLogger := new(MockProgressLogger)
	service := NewProgressService(mockRepo, mockLogger)

	ctx := context.Background()
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	userID := "660e8400-e29b-41d4-a716-446655440001"
	percentage := 75
	lastPage := 10

	// Mock expectations
	mockLogger.On("Info", "updating progress", mock.Anything).Return()
	mockRepo.On("Upsert", ctx, mock.AnythingOfType("*entity.Progress")).
		Return(nil, errors.New("database connection error"))
	mockLogger.On("Error", "failed to upsert progress", mock.Anything).Return()

	// Act
	err := service.UpdateProgress(ctx, materialID, userID, percentage, lastPage)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

// TestUpdateProgress_Idempotency_MultipleCallsSameProgress prueba idempotencia con múltiples llamadas
func TestUpdateProgress_Idempotency_MultipleCallsSameProgress(t *testing.T) {
	// Arrange
	mockRepo := new(MockProgressRepository)
	mockLogger := new(MockProgressLogger)
	service := NewProgressService(mockRepo, mockLogger)

	ctx := context.Background()
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	userID := "660e8400-e29b-41d4-a716-446655440001"
	percentage := 50
	lastPage := 5

	// Crear progress esperado
	matID, _ := valueobject.MaterialIDFromString(materialID)
	uID, _ := valueobject.UserIDFromString(userID)
	expectedProgress := entity.NewProgress(matID, uID)
	_ = expectedProgress.UpdateProgress(percentage, lastPage)

	// Mock expectations (se llamará 3 veces con mismos parámetros)
	mockLogger.On("Info", "updating progress", mock.Anything).Return().Times(3)
	mockRepo.On("Upsert", ctx, mock.AnythingOfType("*entity.Progress")).
		Return(expectedProgress, nil).Times(3)
	mockLogger.On("Info", "progress updated successfully", mock.Anything).Return().Times(3)

	// Act - Llamar UpdateProgress 3 veces con mismos parámetros
	err1 := service.UpdateProgress(ctx, materialID, userID, percentage, lastPage)
	err2 := service.UpdateProgress(ctx, materialID, userID, percentage, lastPage)
	err3 := service.UpdateProgress(ctx, materialID, userID, percentage, lastPage)

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	// Verificar que Upsert fue llamado exactamente 3 veces (operación idempotente)
	mockRepo.AssertNumberOfCalls(t, "Upsert", 3)
}

// TestUpdateProgress_Idempotency_DifferentPercentages prueba múltiples actualizaciones con diferentes valores
func TestUpdateProgress_Idempotency_DifferentPercentages(t *testing.T) {
	// Arrange
	mockRepo := new(MockProgressRepository)
	mockLogger := new(MockProgressLogger)
	service := NewProgressService(mockRepo, mockLogger)

	ctx := context.Background()
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	userID := "660e8400-e29b-41d4-a716-446655440001"

	percentages := []int{25, 50, 75, 100}
	matID, _ := valueobject.MaterialIDFromString(materialID)
	uID, _ := valueobject.UserIDFromString(userID)

	// Mock expectations para cada llamada
	for _, p := range percentages {
		expectedProgress := entity.NewProgress(matID, uID)
		_ = expectedProgress.UpdateProgress(p, p/5) // lastPage proporcional

		mockLogger.On("Info", "updating progress", mock.Anything).Return().Once()
		mockRepo.On("Upsert", ctx, mock.AnythingOfType("*entity.Progress")).
			Return(expectedProgress, nil).Once()

		if p == 100 {
			mockLogger.On("Info", "material completed by user", mock.Anything).Return().Once()
		}

		mockLogger.On("Info", "progress updated successfully", mock.Anything).Return().Once()
	}

	// Act - Actualizar progreso incrementalmente
	for _, p := range percentages {
		err := service.UpdateProgress(ctx, materialID, userID, p, p/5)
		assert.NoError(t, err)
	}

	// Assert
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	// Verificar que Upsert fue llamado exactamente 4 veces
	mockRepo.AssertNumberOfCalls(t, "Upsert", 4)
}
