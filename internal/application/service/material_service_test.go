package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	apperrors "github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// MockMaterialRepository es un mock del repositorio de materiales
type MockMaterialRepository struct {
	mock.Mock
}

func (m *MockMaterialRepository) Create(ctx context.Context, material *entity.Material) error {
	args := m.Called(ctx, material)
	return args.Error(0)
}

func (m *MockMaterialRepository) FindByID(ctx context.Context, id valueobject.MaterialID) (*entity.Material, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Material), args.Error(1)
}

func (m *MockMaterialRepository) FindByIDWithVersions(ctx context.Context, id valueobject.MaterialID) (*entity.Material, []*entity.MaterialVersion, error) {
	args := m.Called(ctx, id)

	var material *entity.Material
	if args.Get(0) != nil {
		material = args.Get(0).(*entity.Material)
	}

	var versions []*entity.MaterialVersion
	if args.Get(1) != nil {
		versions = args.Get(1).([]*entity.MaterialVersion)
	}

	return material, versions, args.Error(2)
}

func (m *MockMaterialRepository) Update(ctx context.Context, material *entity.Material) error {
	args := m.Called(ctx, material)
	return args.Error(0)
}

func (m *MockMaterialRepository) List(ctx context.Context, filters repository.ListFilters) ([]*entity.Material, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Material), args.Error(1)
}

func (m *MockMaterialRepository) FindByAuthor(ctx context.Context, authorID valueobject.UserID) ([]*entity.Material, error) {
	args := m.Called(ctx, authorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Material), args.Error(1)
}

func (m *MockMaterialRepository) UpdateStatus(ctx context.Context, id valueobject.MaterialID, status enum.MaterialStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockMaterialRepository) UpdateProcessingStatus(ctx context.Context, id valueobject.MaterialID, status enum.ProcessingStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockMaterialRepository) CountPublishedMaterials(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// MockPublisher es un mock del publisher de RabbitMQ
type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(ctx context.Context, exchange, routingKey string, message []byte) error {
	args := m.Called(ctx, exchange, routingKey, message)
	return args.Error(0)
}

func (m *MockPublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockLogger es un mock del logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Warn(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Error(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Debug(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Fatal(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Sync() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockLogger) With(keysAndValues ...interface{}) logger.Logger {
	args := m.Called(keysAndValues)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(logger.Logger)
}

// Tests para GetMaterialWithVersions

func TestMaterialService_GetMaterialWithVersions_Success_WithVersions(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	materialID := valueobject.NewMaterialID()
	authorID := valueobject.NewUserID()
	changedByID := valueobject.NewUserID()

	// Material de prueba
	now := time.Now()
	material := entity.ReconstructMaterial(
		materialID,
		"Test Material",
		"Description",
		authorID,
		"",
		"s3://key",
		"https://s3.url",
		enum.MaterialStatusPublished,
		enum.ProcessingStatusCompleted,
		now,
		now,
	)

	// Versiones de prueba (ordenadas por version_number DESC)
	version1 := entity.ReconstructMaterialVersion(
		valueobject.NewMaterialVersionID(),
		materialID,
		2,
		"Version 2 Title",
		"https://s3.url/v2",
		changedByID,
		now,
	)

	version2 := entity.ReconstructMaterialVersion(
		valueobject.NewMaterialVersionID(),
		materialID,
		1,
		"Version 1 Title",
		"https://s3.url/v1",
		changedByID,
		now,
	)

	versions := []*entity.MaterialVersion{version1, version2}

	// Configurar mock
	mockRepo.On("FindByIDWithVersions", ctx, materialID).Return(material, versions, nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	// Act
	result, err := service.GetMaterialWithVersions(ctx, materialID.String())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Material)
	assert.Equal(t, materialID.String(), result.Material.ID)
	assert.Len(t, result.Versions, 2)
	assert.Equal(t, 2, result.Versions[0].VersionNumber) // Debe estar ordenado DESC
	assert.Equal(t, 1, result.Versions[1].VersionNumber)

	mockRepo.AssertExpectations(t)
}

func TestMaterialService_GetMaterialWithVersions_Success_WithoutVersions(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	materialID := valueobject.NewMaterialID()
	authorID := valueobject.NewUserID()

	// Material sin versiones
	now := time.Now()
	material := entity.ReconstructMaterial(
		materialID,
		"Test Material",
		"Description",
		authorID,
		"",
		"s3://key",
		"https://s3.url",
		enum.MaterialStatusPublished,
		enum.ProcessingStatusCompleted,
		now,
		now,
	)

	versions := []*entity.MaterialVersion{} // Array vacío

	// Configurar mock
	mockRepo.On("FindByIDWithVersions", ctx, materialID).Return(material, versions, nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	// Act
	result, err := service.GetMaterialWithVersions(ctx, materialID.String())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Material)
	assert.Len(t, result.Versions, 0) // Debe ser array vacío, no nil

	mockRepo.AssertExpectations(t)
}

func TestMaterialService_GetMaterialWithVersions_MaterialNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	materialID := valueobject.NewMaterialID()

	// Configurar mock para retornar material nil (no encontrado)
	mockRepo.On("FindByIDWithVersions", ctx, materialID).Return(nil, nil, nil)
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	// Act
	result, err := service.GetMaterialWithVersions(ctx, materialID.String())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	// Verificar que es un error de tipo NotFound
	appErr, ok := apperrors.GetAppError(err)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrorCodeNotFound, appErr.Code)

	mockRepo.AssertExpectations(t)
}

func TestMaterialService_GetMaterialWithVersions_InvalidMaterialID(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	invalidID := "not-a-valid-uuid"

	// Configurar mock logger
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	// Act
	result, err := service.GetMaterialWithVersions(ctx, invalidID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	// Verificar que es un error de validación
	appErr, ok := apperrors.GetAppError(err)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrorCodeValidation, appErr.Code)

	// NO debe llamar al repository porque la validación falló antes
	mockRepo.AssertNotCalled(t, "FindByIDWithVersions")
}

func TestMaterialService_GetMaterialWithVersions_DatabaseError(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	materialID := valueobject.NewMaterialID()

	// Configurar mock para retornar error de base de datos
	dbError := errors.New("database connection failed")
	mockRepo.On("FindByIDWithVersions", ctx, materialID).Return(nil, nil, dbError)
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	// Act
	result, err := service.GetMaterialWithVersions(ctx, materialID.String())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	// Verificar que es un error de base de datos
	appErr, ok := apperrors.GetAppError(err)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrorCodeDatabaseError, appErr.Code)

	mockRepo.AssertExpectations(t)
}
