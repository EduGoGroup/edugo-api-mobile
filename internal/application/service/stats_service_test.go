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
		logger:          mockLogger,
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
		logger:          mockLogger,
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
		logger:          mockLogger,
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
		logger:          mockLogger,
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
		logger:          mockLogger,
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
		logger:          mockLogger,
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

// TestGetGlobalStats_MultipleErrors valida manejo cuando múltiples queries fallan
func TestGetGlobalStats_MultipleErrors(t *testing.T) {
	// Arrange
	mockMaterialRepo := new(MockMaterialRepository)
	mockAssessmentRepo := new(MockAssessmentRepository)
	mockProgressRepo := new(MockProgressRepository)
	mockLogger := new(MockLogger)

	// Configurar múltiples repos para retornar errores
	mockMaterialRepo.On("CountPublishedMaterials", mock.Anything).Return(int64(0), errors.New("postgres error"))
	mockAssessmentRepo.On("CountCompletedAssessments", mock.Anything).Return(int64(0), errors.New("mongo error"))
	mockAssessmentRepo.On("CalculateAverageScore", mock.Anything).Return(float64(0), errors.New("mongo error"))
	mockProgressRepo.On("CountActiveUsers", mock.Anything).Return(int64(0), errors.New("postgres error"))
	mockProgressRepo.On("CalculateAverageProgress", mock.Anything).Return(float64(0), errors.New("postgres error"))

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	service := &statsService{
		logger:          mockLogger,
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

// TestGetGlobalStats_AverageScoreError valida error en cálculo de promedio de puntajes
func TestGetGlobalStats_AverageScoreError(t *testing.T) {
	// Arrange
	mockMaterialRepo := new(MockMaterialRepository)
	mockAssessmentRepo := new(MockAssessmentRepository)
	mockProgressRepo := new(MockProgressRepository)
	mockLogger := new(MockLogger)

	mockMaterialRepo.On("CountPublishedMaterials", mock.Anything).Return(int64(150), nil)
	mockAssessmentRepo.On("CountCompletedAssessments", mock.Anything).Return(int64(320), nil)
	// Error en cálculo de promedio
	mockAssessmentRepo.On("CalculateAverageScore", mock.Anything).Return(float64(0), errors.New("calculation error"))
	mockProgressRepo.On("CountActiveUsers", mock.Anything).Return(int64(85), nil)
	mockProgressRepo.On("CalculateAverageProgress", mock.Anything).Return(float64(62.3), nil)

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	service := &statsService{
		logger:          mockLogger,
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

// TestGetGlobalStats_AverageProgressError valida error en cálculo de promedio de progreso
func TestGetGlobalStats_AverageProgressError(t *testing.T) {
	// Arrange
	mockMaterialRepo := new(MockMaterialRepository)
	mockAssessmentRepo := new(MockAssessmentRepository)
	mockProgressRepo := new(MockProgressRepository)
	mockLogger := new(MockLogger)

	mockMaterialRepo.On("CountPublishedMaterials", mock.Anything).Return(int64(150), nil)
	mockAssessmentRepo.On("CountCompletedAssessments", mock.Anything).Return(int64(320), nil)
	mockAssessmentRepo.On("CalculateAverageScore", mock.Anything).Return(float64(78.5), nil)
	mockProgressRepo.On("CountActiveUsers", mock.Anything).Return(int64(85), nil)
	// Error en cálculo de promedio de progreso
	mockProgressRepo.On("CalculateAverageProgress", mock.Anything).Return(float64(0), errors.New("calculation error"))

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	service := &statsService{
		logger:          mockLogger,
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

// Tests para GetMaterialStats

// TestGetMaterialStats_Success valida obtención exitosa de estadísticas de material
func TestGetMaterialStats_Success(t *testing.T) {
	// Arrange
	mockMaterialRepo := new(MockMaterialRepository)
	mockAssessmentRepo := new(MockAssessmentRepository)
	mockProgressRepo := new(MockProgressRepository)
	mockLogger := new(MockLogger)

	service := NewStatsService(mockLogger, mockMaterialRepo, mockAssessmentRepo, mockProgressRepo)

	ctx := context.Background()
	materialID := "550e8400-e29b-41d4-a716-446655440000"

	// Act
	stats, err := service.GetMaterialStats(ctx, materialID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 150, stats.TotalViews)
	assert.Equal(t, 67.5, stats.AvgProgress)
	assert.Equal(t, 45, stats.TotalAttempts)
	assert.Equal(t, 78.3, stats.AvgScore)
}

// TestGetMaterialStats_InvalidMaterialID valida error con materialID inválido
func TestGetMaterialStats_InvalidMaterialID(t *testing.T) {
	// Arrange
	mockMaterialRepo := new(MockMaterialRepository)
	mockAssessmentRepo := new(MockAssessmentRepository)
	mockProgressRepo := new(MockProgressRepository)
	mockLogger := new(MockLogger)

	service := NewStatsService(mockLogger, mockMaterialRepo, mockAssessmentRepo, mockProgressRepo)

	ctx := context.Background()
	materialID := "invalid-uuid"

	// Act
	stats, err := service.GetMaterialStats(ctx, materialID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "invalid material_id")
}

// TestGetMaterialStats_EmptyMaterialID valida error con materialID vacío
func TestGetMaterialStats_EmptyMaterialID(t *testing.T) {
	// Arrange
	mockMaterialRepo := new(MockMaterialRepository)
	mockAssessmentRepo := new(MockAssessmentRepository)
	mockProgressRepo := new(MockProgressRepository)
	mockLogger := new(MockLogger)

	service := NewStatsService(mockLogger, mockMaterialRepo, mockAssessmentRepo, mockProgressRepo)

	ctx := context.Background()
	materialID := ""

	// Act
	stats, err := service.GetMaterialStats(ctx, materialID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "invalid material_id")
}

// TestGetMaterialStats_ValidUUID valida que acepta cualquier UUID válido
func TestGetMaterialStats_ValidUUID(t *testing.T) {
	// Arrange
	mockMaterialRepo := new(MockMaterialRepository)
	mockAssessmentRepo := new(MockAssessmentRepository)
	mockProgressRepo := new(MockProgressRepository)
	mockLogger := new(MockLogger)

	service := NewStatsService(mockLogger, mockMaterialRepo, mockAssessmentRepo, mockProgressRepo)

	ctx := context.Background()
	materialIDs := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"660e8400-e29b-41d4-a716-446655440001",
		"770e8400-e29b-41d4-a716-446655440002",
	}

	// Act & Assert
	for _, materialID := range materialIDs {
		stats, err := service.GetMaterialStats(ctx, materialID)
		assert.NoError(t, err)
		assert.NotNil(t, stats)
	}
}

// TestGetMaterialStats_MockedValues valida que retorna valores mockeados correctos
func TestGetMaterialStats_MockedValues(t *testing.T) {
	// Arrange
	mockMaterialRepo := new(MockMaterialRepository)
	mockAssessmentRepo := new(MockAssessmentRepository)
	mockProgressRepo := new(MockProgressRepository)
	mockLogger := new(MockLogger)

	service := NewStatsService(mockLogger, mockMaterialRepo, mockAssessmentRepo, mockProgressRepo)

	ctx := context.Background()
	materialID := "550e8400-e29b-41d4-a716-446655440000"

	// Act
	stats, err := service.GetMaterialStats(ctx, materialID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	// Verificar valores específicos mockeados
	assert.Greater(t, stats.TotalViews, 0, "TotalViews debe ser mayor a 0")
	assert.Greater(t, stats.AvgProgress, 0.0, "AvgProgress debe ser mayor a 0")
	assert.Greater(t, stats.TotalAttempts, 0, "TotalAttempts debe ser mayor a 0")
	assert.Greater(t, stats.AvgScore, 0.0, "AvgScore debe ser mayor a 0")
}

// TestNewStatsService valida la creación del servicio
func TestNewStatsService(t *testing.T) {
	// Arrange
	mockMaterialRepo := new(MockMaterialRepository)
	mockAssessmentRepo := new(MockAssessmentRepository)
	mockProgressRepo := new(MockProgressRepository)
	mockLogger := new(MockLogger)

	// Act
	service := NewStatsService(mockLogger, mockMaterialRepo, mockAssessmentRepo, mockProgressRepo)

	// Assert
	assert.NotNil(t, service)
	assert.Implements(t, (*StatsService)(nil), service)
}
