package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
)

// MockSummaryRepository mock del repositorio
type MockSummaryRepository struct {
	mock.Mock
}

func (m *MockSummaryRepository) FindByMaterialID(ctx context.Context, materialID valueobject.MaterialID) (*repository.MaterialSummary, error) {
	args := m.Called(ctx, materialID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.MaterialSummary), args.Error(1)
}

func (m *MockSummaryRepository) Exists(ctx context.Context, materialID valueobject.MaterialID) (bool, error) {
	args := m.Called(ctx, materialID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSummaryRepository) Save(ctx context.Context, summary *repository.MaterialSummary) error {
	args := m.Called(ctx, summary)
	return args.Error(0)
}

func (m *MockSummaryRepository) Delete(ctx context.Context, materialID valueobject.MaterialID) error {
	args := m.Called(ctx, materialID)
	return args.Error(0)
}

// TestNewSummaryService verifica que el constructor inicialice correctamente
func TestNewSummaryService(t *testing.T) {
	mockRepo := new(MockSummaryRepository)
	mockLogger := new(MockLogger)

	service := NewSummaryService(mockRepo, mockLogger)

	assert.NotNil(t, service)
}

// TestGetSummary_Success verifica obtener summary exitosamente
func TestGetSummary_Success(t *testing.T) {
	mockRepo := new(MockSummaryRepository)
	mockLogger := new(MockLogger)
	service := NewSummaryService(mockRepo, mockLogger)

	ctx := context.Background()
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	matID, _ := valueobject.MaterialIDFromString(materialID)

	expectedSummary := &repository.MaterialSummary{
		MaterialID:  matID,
		MainIdeas:   []string{"Idea 1", "Idea 2"},
		KeyConcepts: map[string]string{"concept": "definition"},
	}

	mockRepo.On("FindByMaterialID", ctx, matID).Return(expectedSummary, nil)

	result, err := service.GetSummary(ctx, materialID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedSummary, result)
	mockRepo.AssertExpectations(t)
}

// TestGetSummary_InvalidMaterialID verifica error con ID inválido
func TestGetSummary_InvalidMaterialID(t *testing.T) {
	mockRepo := new(MockSummaryRepository)
	mockLogger := new(MockLogger)
	service := NewSummaryService(mockRepo, mockLogger)

	ctx := context.Background()
	invalidID := "invalid-uuid"

	result, err := service.GetSummary(ctx, invalidID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid material_id")
}

// TestGetSummary_NotFound verifica error cuando no existe el summary
func TestGetSummary_NotFound(t *testing.T) {
	mockRepo := new(MockSummaryRepository)
	mockLogger := new(MockLogger)
	service := NewSummaryService(mockRepo, mockLogger)

	ctx := context.Background()
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	matID, _ := valueobject.MaterialIDFromString(materialID)

	mockRepo.On("FindByMaterialID", ctx, matID).Return(nil, nil)

	result, err := service.GetSummary(ctx, materialID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "summary")
	mockRepo.AssertExpectations(t)
}

// TestGetSummary_DatabaseError verifica error de base de datos
func TestGetSummary_DatabaseError(t *testing.T) {
	mockRepo := new(MockSummaryRepository)
	mockLogger := new(MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()
	service := NewSummaryService(mockRepo, mockLogger)

	ctx := context.Background()
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	matID, _ := valueobject.MaterialIDFromString(materialID)

	dbError := errors.New("database connection failed")
	mockRepo.On("FindByMaterialID", ctx, matID).Return(nil, dbError)

	result, err := service.GetSummary(ctx, materialID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
