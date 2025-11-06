package service

import (
	"context"
	"errors"
	"testing"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestGetGlobalStats_Success valida que se obtienen estadísticas correctamente cuando todas las queries son exitosas
func TestGetGlobalStats_Success(t *testing.T) {
	// Arrange
	mockMaterialRepo := new(MockMaterialRepository)
	mockAssessmentRepo := new(MockAssessmentRepository)
	mockProgressRepo := new(MockProgressRepository)
	mockLogger := new(MockLogger)

	// Configurar mocks para retornar valores exitosos
	mockMaterialRepo.On("CountPublishedMaterials", mock.Anything).Return(int64(150), nil)
	mockAssessmentRepo.On("CountCompletedAssessments", mock.Anything).Return(int64(320), nil)
	mockAssessmentRepo.On("CalculateAverageScore", mock.Anything).Return(float64(78.5), nil)
	mockProgressRepo.On("CountActiveUsers", mock.Anything).Return(int64(85), nil)
	mockProgressRepo.On("CalculateAverageProgress", mock.Anything).Return(float64(62.3), nil)

	// Configurar logger para aceptar cualquier llamada
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	service := &statsService{
		logger:         mockLogger,
		materialStats:   mockMaterialRepo,
		assessmentStats: mockAssessmentRepo,
		progressStats:   mockProgressRepo,
	}

	// Act
	ctx := context.Background()
	stats, err := service.GetGlobalStats(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(150), stats.TotalPublishedMaterials)
	assert.Equal(t, int64(320), stats.TotalCompletedAssessments)
	assert.Equal(t, float64(78.5), stats.AverageAssessmentScore)
	assert.Equal(t, int64(85), stats.ActiveUsersLast30Days)
	assert.Equal(t, float64(62.3), stats.AverageProgress)
	assert.NotZero(t, stats.GeneratedAt)

	mockMaterialRepo.AssertExpectations(t)
	mockAssessmentRepo.AssertExpectations(t)
	mockProgressRepo.AssertExpectations(t)
}

// TestGetGlobalStats_MaterialRepoError valida manejo de error en query de materiales
func TestGetGlobalStats_MaterialRepoError(t *testing.T) {
	// Arrange
	mockMaterialRepo := new(MockMaterialRepository)
	mockAssessmentRepo := new(MockAssessmentRepository)
	mockProgressRepo := new(MockProgressRepository)
	mockLogger := new(MockLogger)

	// Configurar mock de materiales para retornar error
	mockMaterialRepo.On("CountPublishedMaterials", mock.Anything).Return(int64(0), errors.New("database connection error"))

	// Los demás repos retornan valores exitosos (pero no deberían usarse debido al error)
	mockAssessmentRepo.On("CountCompletedAssessments", mock.Anything).Return(int64(320), nil)
	mockAssessmentRepo.On("CalculateAverageScore", mock.Anything).Return(float64(78.5), nil)
	mockProgressRepo.On("CountActiveUsers", mock.Anything).Return(int64(85), nil)
	mockProgressRepo.On("CalculateAverageProgress", mock.Anything).Return(float64(62.3), nil)

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	service := &statsService{
		logger:         mockLogger,
		materialStats:   mockMaterialRepo,
		assessmentStats: mockAssessmentRepo,
		progressStats:   mockProgressRepo,
	}

	// Act
	ctx := context.Background()
	stats, err := service.GetGlobalStats(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, stats)
	mockMaterialRepo.AssertExpectations(t)
}

// TestGetGlobalStats_AssessmentRepoError valida manejo de error en query de evaluaciones
func TestGetGlobalStats_AssessmentRepoError(t *testing.T) {
	// Arrange
	mockMaterialRepo := new(MockMaterialRepository)
	mockAssessmentRepo := new(MockAssessmentRepository)
	mockProgressRepo := new(MockProgressRepository)
	mockLogger := new(MockLogger)

	mockMaterialRepo.On("CountPublishedMaterials", mock.Anything).Return(int64(150), nil)
	// Configurar assessment repo para retornar error en countDocuments
	mockAssessmentRepo.On("CountCompletedAssessments", mock.Anything).Return(int64(0), errors.New("mongo connection error"))
	mockAssessmentRepo.On("CalculateAverageScore", mock.Anything).Return(float64(0), nil)
	mockProgressRepo.On("CountActiveUsers", mock.Anything).Return(int64(85), nil)
	mockProgressRepo.On("CalculateAverageProgress", mock.Anything).Return(float64(62.3), nil)

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	service := &statsService{
		logger:         mockLogger,
		materialStats:   mockMaterialRepo,
		assessmentStats: mockAssessmentRepo,
		progressStats:   mockProgressRepo,
	}

	// Act
	ctx := context.Background()
	stats, err := service.GetGlobalStats(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, stats)
}

// TestGetGlobalStats_ProgressRepoError valida manejo de error en query de progreso
func TestGetGlobalStats_ProgressRepoError(t *testing.T) {
	// Arrange
	mockMaterialRepo := new(MockMaterialRepository)
	mockAssessmentRepo := new(MockAssessmentRepository)
	mockProgressRepo := new(MockProgressRepository)
	mockLogger := new(MockLogger)

	mockMaterialRepo.On("CountPublishedMaterials", mock.Anything).Return(int64(150), nil)
	mockAssessmentRepo.On("CountCompletedAssessments", mock.Anything).Return(int64(320), nil)
	mockAssessmentRepo.On("CalculateAverageScore", mock.Anything).Return(float64(78.5), nil)
	// Configurar progress repo para retornar error
	mockProgressRepo.On("CountActiveUsers", mock.Anything).Return(int64(0), errors.New("postgres connection error"))
	mockProgressRepo.On("CalculateAverageProgress", mock.Anything).Return(float64(0), nil)

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	service := &statsService{
		logger:         mockLogger,
		materialStats:   mockMaterialRepo,
		assessmentStats: mockAssessmentRepo,
		progressStats:   mockProgressRepo,
	}

	// Act
	ctx := context.Background()
	stats, err := service.GetGlobalStats(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, stats)
}

// TestGetGlobalStats_AllZeros valida que funciona con sistema vacío (sin datos)
func TestGetGlobalStats_AllZeros(t *testing.T) {
	// Arrange
	mockMaterialRepo := new(MockMaterialRepository)
	mockAssessmentRepo := new(MockAssessmentRepository)
	mockProgressRepo := new(MockProgressRepository)
	mockLogger := new(MockLogger)

	// Simular sistema vacío - todos los counts son 0
	mockMaterialRepo.On("CountPublishedMaterials", mock.Anything).Return(int64(0), nil)
	mockAssessmentRepo.On("CountCompletedAssessments", mock.Anything).Return(int64(0), nil)
	mockAssessmentRepo.On("CalculateAverageScore", mock.Anything).Return(float64(0), nil)
	mockProgressRepo.On("CountActiveUsers", mock.Anything).Return(int64(0), nil)
	mockProgressRepo.On("CalculateAverageProgress", mock.Anything).Return(float64(0), nil)

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	service := &statsService{
		logger:         mockLogger,
		materialStats:   mockMaterialRepo,
		assessmentStats: mockAssessmentRepo,
		progressStats:   mockProgressRepo,
	}

	// Act
	ctx := context.Background()
	stats, err := service.GetGlobalStats(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(0), stats.TotalPublishedMaterials)
	assert.Equal(t, int64(0), stats.TotalCompletedAssessments)
	assert.Equal(t, float64(0), stats.AverageAssessmentScore)
	assert.Equal(t, int64(0), stats.ActiveUsersLast30Days)
	assert.Equal(t, float64(0), stats.AverageProgress)

	mockMaterialRepo.AssertExpectations(t)
	mockAssessmentRepo.AssertExpectations(t)
	mockProgressRepo.AssertExpectations(t)
}

// TestGetGlobalStats_DTOStructure valida la estructura del DTO retornado
func TestGetGlobalStats_DTOStructure(t *testing.T) {
	// Arrange
	mockMaterialRepo := new(MockMaterialRepository)
	mockAssessmentRepo := new(MockAssessmentRepository)
	mockProgressRepo := new(MockProgressRepository)
	mockLogger := new(MockLogger)

	mockMaterialRepo.On("CountPublishedMaterials", mock.Anything).Return(int64(100), nil)
	mockAssessmentRepo.On("CountCompletedAssessments", mock.Anything).Return(int64(200), nil)
	mockAssessmentRepo.On("CalculateAverageScore", mock.Anything).Return(float64(85.0), nil)
	mockProgressRepo.On("CountActiveUsers", mock.Anything).Return(int64(50), nil)
	mockProgressRepo.On("CalculateAverageProgress", mock.Anything).Return(float64(70.0), nil)

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	service := &statsService{
		logger:         mockLogger,
		materialStats:   mockMaterialRepo,
		assessmentStats: mockAssessmentRepo,
		progressStats:   mockProgressRepo,
	}

	// Act
	ctx := context.Background()
	stats, err := service.GetGlobalStats(ctx)

	// Assert
	assert.NoError(t, err)
	assert.IsType(t, &dto.GlobalStatsDTO{}, stats)
	assert.NotZero(t, stats.GeneratedAt, "GeneratedAt debe estar establecido")
}
