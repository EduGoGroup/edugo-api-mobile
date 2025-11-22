package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	apperrors "github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// Helper function para crear pointer a string
func stringPtr(s string) *string {
	return &s
}

// MockMaterialRepository es un mock del repositorio de materiales
type MockMaterialRepository struct {
	mock.Mock
}

func (m *MockMaterialRepository) Create(ctx context.Context, material *pgentities.Material) error {
	args := m.Called(ctx, material)
	return args.Error(0)
}

func (m *MockMaterialRepository) FindByID(ctx context.Context, id valueobject.MaterialID) (*pgentities.Material, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pgentities.Material), args.Error(1)
}

func (m *MockMaterialRepository) FindByIDWithVersions(ctx context.Context, id valueobject.MaterialID) (*pgentities.Material, []*pgentities.MaterialVersion, error) {
	args := m.Called(ctx, id)

	var material *pgentities.Material
	if args.Get(0) != nil {
		material = args.Get(0).(*pgentities.Material)
	}

	var versions []*pgentities.MaterialVersion
	if args.Get(1) != nil {
		versions = args.Get(1).([]*pgentities.MaterialVersion)
	}

	return material, versions, args.Error(2)
}

func (m *MockMaterialRepository) Update(ctx context.Context, material *pgentities.Material) error {
	args := m.Called(ctx, material)
	return args.Error(0)
}

func (m *MockMaterialRepository) List(ctx context.Context, filters repository.ListFilters) ([]*pgentities.Material, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*pgentities.Material), args.Error(1)
}

func (m *MockMaterialRepository) FindByAuthor(ctx context.Context, authorID valueobject.UserID) ([]*pgentities.Material, error) {
	args := m.Called(ctx, authorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*pgentities.Material), args.Error(1)
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

// Tests para CreateMaterial

func TestMaterialService_CreateMaterial_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	authorID := valueobject.NewUserID()
	req := dto.CreateMaterialRequest{
		Title:       "Test Material",
		Description: "Test Description",
		Subject:     valueobject.NewUserID().String(),
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*pgentities.Material")).Return(nil)
	mockPublisher.On("Publish", ctx, "edugo.materials", "material.uploaded", mock.Anything).Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	// Act
	result, err := service.CreateMaterial(ctx, req, authorID.String())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Title, result.Title)
	assert.Equal(t, req.Description, result.Description)
	assert.Equal(t, authorID.String(), result.UploadedByTeacherID)

	mockRepo.AssertExpectations(t)
}

func TestMaterialService_CreateMaterial_ValidationError_EmptyTitle(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	authorID := valueobject.NewUserID()
	req := dto.CreateMaterialRequest{
		Title:       "", // Empty title
		Description: "Test Description",
	}

	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	// Act
	result, err := service.CreateMaterial(ctx, req, authorID.String())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	appErr, ok := apperrors.GetAppError(err)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrorCodeValidation, appErr.Code)

	mockRepo.AssertNotCalled(t, "Create")
}

func TestMaterialService_CreateMaterial_ValidationError_TitleTooShort(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	authorID := valueobject.NewUserID()
	req := dto.CreateMaterialRequest{
		Title:       "AB", // Too short (min 3)
		Description: "Test Description",
	}

	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	// Act
	result, err := service.CreateMaterial(ctx, req, authorID.String())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertNotCalled(t, "Create")
}

func TestMaterialService_CreateMaterial_InvalidAuthorID(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	req := dto.CreateMaterialRequest{
		Title:       "Test Material",
		Description: "Test Description",
	}

	// Act
	result, err := service.CreateMaterial(ctx, req, "invalid-uuid")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	appErr, ok := apperrors.GetAppError(err)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrorCodeValidation, appErr.Code)

	mockRepo.AssertNotCalled(t, "Create")
}

func TestMaterialService_CreateMaterial_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	authorID := valueobject.NewUserID()
	req := dto.CreateMaterialRequest{
		Title:       "Test Material",
		Description: "Test Description",
	}

	dbError := errors.New("database connection failed")
	mockRepo.On("Create", ctx, mock.AnythingOfType("*pgentities.Material")).Return(dbError)
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	// Act
	result, err := service.CreateMaterial(ctx, req, authorID.String())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	appErr, ok := apperrors.GetAppError(err)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrorCodeDatabaseError, appErr.Code)

	mockRepo.AssertExpectations(t)
}

func TestMaterialService_CreateMaterial_PublishEventFailure(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	authorID := valueobject.NewUserID()
	req := dto.CreateMaterialRequest{
		Title:       "Test Material",
		Description: "Test Description",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*pgentities.Material")).Return(nil)
	mockPublisher.On("Publish", ctx, "edugo.materials", "material.uploaded", mock.Anything).Return(errors.New("rabbitmq connection failed"))
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	// Act - Should succeed even if event publishing fails
	result, err := service.CreateMaterial(ctx, req, authorID.String())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

// Tests para GetMaterial

func TestMaterialService_GetMaterial_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	materialID := valueobject.NewMaterialID()
	authorID := valueobject.NewUserID()

	now := time.Now()
	material := &pgentities.Material{
		ID:                  materialID.UUID().UUID,
		Title:               "Test Material",
		Description:         stringPtr("Description"),
		UploadedByTeacherID: authorID.UUID().UUID,
		Subject:             stringPtr(""),
		FileURL:             "https://s3.url",
		Status:              string(enum.MaterialStatusPublished),
		IsPublic:            false,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	mockRepo.On("FindByID", ctx, materialID).Return(material, nil)

	// Act
	result, err := service.GetMaterial(ctx, materialID.String())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, materialID.String(), result.ID)
	assert.Equal(t, "Test Material", result.Title)

	mockRepo.AssertExpectations(t)
}

func TestMaterialService_GetMaterial_InvalidID(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()

	// Act
	result, err := service.GetMaterial(ctx, "invalid-uuid")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	appErr, ok := apperrors.GetAppError(err)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrorCodeValidation, appErr.Code)

	mockRepo.AssertNotCalled(t, "FindByID")
}

func TestMaterialService_GetMaterial_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	materialID := valueobject.NewMaterialID()

	mockRepo.On("FindByID", ctx, materialID).Return(nil, nil)

	// Act
	result, err := service.GetMaterial(ctx, materialID.String())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	appErr, ok := apperrors.GetAppError(err)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrorCodeNotFound, appErr.Code)

	mockRepo.AssertExpectations(t)
}

func TestMaterialService_GetMaterial_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	materialID := valueobject.NewMaterialID()

	dbError := errors.New("database error")
	mockRepo.On("FindByID", ctx, materialID).Return(nil, dbError)

	// Act
	result, err := service.GetMaterial(ctx, materialID.String())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	appErr, ok := apperrors.GetAppError(err)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrorCodeNotFound, appErr.Code)

	mockRepo.AssertExpectations(t)
}

// Tests para NotifyUploadComplete

func TestMaterialService_NotifyUploadComplete_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	materialID := valueobject.NewMaterialID()
	authorID := valueobject.NewUserID()

	now := time.Now()
	material := &pgentities.Material{
		ID:                  materialID.UUID().UUID,
		Title:               "Test Material",
		Description:         stringPtr("Description"),
		UploadedByTeacherID: authorID.UUID().UUID,
		Subject:             stringPtr(""),
		FileURL:             "",
		Status:              string(enum.MaterialStatusDraft),
		IsPublic:            false,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	req := dto.UploadCompleteRequest{
		FileURL: "https://s3.amazonaws.com/bucket/materials/test.pdf",
	}

	mockRepo.On("FindByID", ctx, materialID).Return(material, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*pgentities.Material")).Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	// Act
	err := service.NotifyUploadComplete(ctx, materialID.String(), req)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestMaterialService_NotifyUploadComplete_ValidationError_EmptyS3Key(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	materialID := valueobject.NewMaterialID()

	req := dto.UploadCompleteRequest{
		FileURL: "", // Empty
	}

	// Act
	err := service.NotifyUploadComplete(ctx, materialID.String(), req)

	// Assert
	assert.Error(t, err)

	mockRepo.AssertNotCalled(t, "FindByID")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestMaterialService_NotifyUploadComplete_ValidationError_InvalidURL(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	materialID := valueobject.NewMaterialID()

	req := dto.UploadCompleteRequest{
		FileURL: "not-a-valid-url",
	}

	// Act
	err := service.NotifyUploadComplete(ctx, materialID.String(), req)

	// Assert
	assert.Error(t, err)

	mockRepo.AssertNotCalled(t, "FindByID")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestMaterialService_NotifyUploadComplete_InvalidMaterialID(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()

	req := dto.UploadCompleteRequest{
		FileURL: "https://s3.amazonaws.com/bucket/materials/test.pdf",
	}

	// Act
	err := service.NotifyUploadComplete(ctx, "invalid-uuid", req)

	// Assert
	assert.Error(t, err)

	appErr, ok := apperrors.GetAppError(err)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrorCodeValidation, appErr.Code)

	mockRepo.AssertNotCalled(t, "FindByID")
}

func TestMaterialService_NotifyUploadComplete_MaterialNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	materialID := valueobject.NewMaterialID()

	req := dto.UploadCompleteRequest{
		FileURL: "https://s3.amazonaws.com/bucket/materials/test.pdf",
	}

	mockRepo.On("FindByID", ctx, materialID).Return(nil, nil)

	// Act
	err := service.NotifyUploadComplete(ctx, materialID.String(), req)

	// Assert
	assert.Error(t, err)

	appErr, ok := apperrors.GetAppError(err)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrorCodeNotFound, appErr.Code)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestMaterialService_NotifyUploadComplete_UpdateError(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	materialID := valueobject.NewMaterialID()
	authorID := valueobject.NewUserID()

	now := time.Now()
	material := &pgentities.Material{
		ID:                  materialID.UUID().UUID,
		Title:               "Test Material",
		Description:         stringPtr("Description"),
		UploadedByTeacherID: authorID.UUID().UUID,
		Subject:             stringPtr(""),
		FileURL:             "",
		Status:              string(enum.MaterialStatusDraft),
		IsPublic:            false,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	req := dto.UploadCompleteRequest{
		FileURL: "https://s3.amazonaws.com/bucket/materials/test.pdf",
	}

	dbError := errors.New("database error")
	mockRepo.On("FindByID", ctx, materialID).Return(material, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*pgentities.Material")).Return(dbError)
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	// Act
	err := service.NotifyUploadComplete(ctx, materialID.String(), req)

	// Assert
	assert.Error(t, err)

	appErr, ok := apperrors.GetAppError(err)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrorCodeDatabaseError, appErr.Code)

	mockRepo.AssertExpectations(t)
}

// Tests para ListMaterials

func TestMaterialService_ListMaterials_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	filters := repository.ListFilters{
		Limit:  10,
		Offset: 0,
	}

	materialID1 := valueobject.NewMaterialID()
	materialID2 := valueobject.NewMaterialID()
	authorID := valueobject.NewUserID()
	now := time.Now()

	materials := []*pgentities.Material{
		{
			ID:                  materialID1.UUID().UUID,
			Title:               "Material 1",
			Description:         stringPtr("Description 1"),
			UploadedByTeacherID: authorID.UUID().UUID,
			Subject:             stringPtr(""),
			FileURL:             "https://s3.url1",
			Status:              string(enum.MaterialStatusPublished),
			IsPublic:            false,
			CreatedAt:           now,
			UpdatedAt:           now,
		},
		{
			ID:                  materialID2.UUID().UUID,
			Title:               "Material 2",
			Description:         stringPtr("Description 2"),
			UploadedByTeacherID: authorID.UUID().UUID,
			Subject:             stringPtr(""),
			FileURL:             "https://s3.url2",
			Status:              string(enum.MaterialStatusPublished),
			IsPublic:            false,
			CreatedAt:           now,
			UpdatedAt:           now,
		},
	}

	mockRepo.On("List", ctx, filters).Return(materials, nil)

	// Act
	result, err := service.ListMaterials(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "Material 1", result[0].Title)
	assert.Equal(t, "Material 2", result[1].Title)

	mockRepo.AssertExpectations(t)
}

func TestMaterialService_ListMaterials_EmptyList(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	filters := repository.ListFilters{
		Limit:  10,
		Offset: 0,
	}

	mockRepo.On("List", ctx, filters).Return([]*pgentities.Material{}, nil)

	// Act
	result, err := service.ListMaterials(ctx, filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)

	mockRepo.AssertExpectations(t)
}

func TestMaterialService_ListMaterials_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockMaterialRepository)
	mockPublisher := new(MockPublisher)
	mockLogger := new(MockLogger)

	service := NewMaterialService(mockRepo, mockPublisher, mockLogger)

	ctx := context.Background()
	filters := repository.ListFilters{
		Limit:  10,
		Offset: 0,
	}

	dbError := errors.New("database error")
	mockRepo.On("List", ctx, filters).Return(nil, dbError)
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	// Act
	result, err := service.ListMaterials(ctx, filters)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	appErr, ok := apperrors.GetAppError(err)
	assert.True(t, ok)
	assert.Equal(t, apperrors.ErrorCodeDatabaseError, appErr.Code)

	mockRepo.AssertExpectations(t)
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
	material := &pgentities.Material{
		ID:                  materialID.UUID().UUID,
		Title:               "Test Material",
		Description:         stringPtr("Description"),
		UploadedByTeacherID: authorID.UUID().UUID,
		Subject:             stringPtr(""),
		FileURL:             "https://s3.url",
		Status:              string(enum.MaterialStatusPublished),
		IsPublic:            false,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	// Versiones de prueba (ordenadas por version_number DESC)
	version1 := &pgentities.MaterialVersion{
		ID:            uuid.New(),
		MaterialID:    materialID.UUID().UUID,
		VersionNumber: 2,
		Title:         "Version 2 Title",
		ContentURL:    "https://s3.url/v2",
		ChangedBy:     changedByID.UUID().UUID,
		CreatedAt:     now,
	}

	version2 := &pgentities.MaterialVersion{
		ID:            uuid.New(),
		MaterialID:    materialID.UUID().UUID,
		VersionNumber: 1,
		Title:         "Version 1 Title",
		ContentURL:    "https://s3.url/v1",
		ChangedBy:     changedByID.UUID().UUID,
		CreatedAt:     now,
	}

	versions := []*pgentities.MaterialVersion{version1, version2}

	// Configurar mock
	mockRepo.On("FindByIDWithVersions", ctx, materialID).Return(*material, versions, nil)
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
	material := &pgentities.Material{
		ID:                  materialID.UUID().UUID,
		Title:               "Test Material",
		Description:         stringPtr("Description"),
		UploadedByTeacherID: authorID.UUID().UUID,
		Subject:             stringPtr(""),
		FileURL:             "https://s3.url",
		Status:              string(enum.MaterialStatusPublished),
		IsPublic:            false,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	versions := []*pgentities.MaterialVersion{} // Array vacío

	// Configurar mock
	mockRepo.On("FindByIDWithVersions", ctx, materialID).Return(*material, versions, nil)
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
