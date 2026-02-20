package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
)

// ============================================
// Tests: GetNavigationConfig - Dynamic Navigation
// ============================================

func TestScreenService_GetNavigationConfig_FullPermissions(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockResourceReader := new(MockResourceReader)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockResourceReader, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	resources := []*repository.MenuResource{
		{ID: "r1", Key: "dashboard", DisplayName: "Home", Icon: strPtr("home"), SortOrder: 0, Scope: "system"},
		{ID: "r2", Key: "materials", DisplayName: "Materials", Icon: strPtr("folder"), SortOrder: 1, Scope: "school"},
		{ID: "r3", Key: "assessments", DisplayName: "Assessments", Icon: strPtr("quiz"), SortOrder: 2, Scope: "school"},
		{ID: "r4", Key: "settings", DisplayName: "Settings", Icon: strPtr("settings"), SortOrder: 4, Scope: "system"},
	}
	mappings := []*repository.ResourceScreenMapping{
		{ResourceKey: "dashboard", ScreenKey: "dashboard-teacher", IsDefault: true},
		{ResourceKey: "materials", ScreenKey: "materials-list", IsDefault: true},
		{ResourceKey: "assessments", ScreenKey: "assessments-list", IsDefault: true},
		{ResourceKey: "settings", ScreenKey: "app-settings", IsDefault: true},
	}

	mockResourceReader.On("GetMenuResources", ctx).Return(resources, nil)
	mockResourceReader.On("GetResourceScreenMappings", ctx, []string{"dashboard", "materials", "assessments", "settings"}).Return(mappings, nil)

	permissions := []string{"materials:read", "assessments:read"}

	// Act
	result, err := svc.GetNavigationConfig(ctx, userID, "mobile", permissions)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	// system scope + materials + assessments = 4 items
	totalItems := len(result.BottomNav) + len(result.DrawerItems)
	assert.Equal(t, 4, totalItems)
	assert.Equal(t, 1, result.Version)

	mockResourceReader.AssertExpectations(t)
}

func TestScreenService_GetNavigationConfig_PartialPermissions(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockResourceReader := new(MockResourceReader)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockResourceReader, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	resources := []*repository.MenuResource{
		{ID: "r1", Key: "dashboard", DisplayName: "Home", Icon: strPtr("home"), SortOrder: 0, Scope: "system"},
		{ID: "r2", Key: "materials", DisplayName: "Materials", Icon: strPtr("folder"), SortOrder: 1, Scope: "school"},
		{ID: "r3", Key: "assessments", DisplayName: "Assessments", Icon: strPtr("quiz"), SortOrder: 2, Scope: "school"},
	}
	mappings := []*repository.ResourceScreenMapping{
		{ResourceKey: "dashboard", ScreenKey: "dashboard-teacher", IsDefault: true},
		{ResourceKey: "materials", ScreenKey: "materials-list", IsDefault: true},
	}

	mockResourceReader.On("GetMenuResources", ctx).Return(resources, nil)
	// Solo dashboard (system) y materials (tiene permiso) - assessments filtrado
	mockResourceReader.On("GetResourceScreenMappings", ctx, []string{"dashboard", "materials"}).Return(mappings, nil)

	// Solo tiene permiso de materials, no de assessments
	permissions := []string{"materials:read"}

	// Act
	result, err := svc.GetNavigationConfig(ctx, userID, "mobile", permissions)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	// dashboard (system) + materials (permiso) = 2 items
	totalItems := len(result.BottomNav) + len(result.DrawerItems)
	assert.Equal(t, 2, totalItems)

	// Verificar que assessments fue filtrado
	for _, item := range result.BottomNav {
		assert.NotEqual(t, "assessments", item.Key, "assessments debe ser filtrado sin permiso")
	}

	mockResourceReader.AssertExpectations(t)
}

func TestScreenService_GetNavigationConfig_NoPermissions(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockResourceReader := new(MockResourceReader)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockResourceReader, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	resources := []*repository.MenuResource{
		{ID: "r1", Key: "materials", DisplayName: "Materials", Icon: strPtr("folder"), SortOrder: 1, Scope: "school"},
		{ID: "r2", Key: "assessments", DisplayName: "Assessments", Icon: strPtr("quiz"), SortOrder: 2, Scope: "school"},
	}

	mockResourceReader.On("GetMenuResources", ctx).Return(resources, nil)

	// Sin permisos - no hay recursos school scope
	permissions := []string{}

	// Act
	result, err := svc.GetNavigationConfig(ctx, userID, "mobile", permissions)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Empty(t, result.BottomNav, "sin permisos no debe haber items")
	assert.Empty(t, result.DrawerItems)

	mockResourceReader.AssertExpectations(t)
}

func TestScreenService_GetNavigationConfig_SystemScopeAlwaysVisible(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockResourceReader := new(MockResourceReader)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockResourceReader, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	resources := []*repository.MenuResource{
		{ID: "r1", Key: "dashboard", DisplayName: "Home", Icon: strPtr("home"), SortOrder: 0, Scope: "system"},
		{ID: "r2", Key: "settings", DisplayName: "Settings", Icon: strPtr("settings"), SortOrder: 4, Scope: "system"},
		{ID: "r3", Key: "materials", DisplayName: "Materials", Icon: strPtr("folder"), SortOrder: 1, Scope: "school"},
	}
	mappings := []*repository.ResourceScreenMapping{
		{ResourceKey: "dashboard", ScreenKey: "dashboard-teacher", IsDefault: true},
		{ResourceKey: "settings", ScreenKey: "app-settings", IsDefault: true},
	}

	mockResourceReader.On("GetMenuResources", ctx).Return(resources, nil)
	mockResourceReader.On("GetResourceScreenMappings", ctx, []string{"dashboard", "settings"}).Return(mappings, nil)

	// Sin permisos de school scope, pero system scope siempre visible
	permissions := []string{}

	// Act
	result, err := svc.GetNavigationConfig(ctx, userID, "mobile", permissions)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	// Solo system scope items deben estar presentes
	totalItems := len(result.BottomNav) + len(result.DrawerItems)
	assert.Equal(t, 2, totalItems)

	mockResourceReader.AssertExpectations(t)
}

func TestScreenService_GetNavigationConfig_DesktopPlatform_AllInDrawer(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockResourceReader := new(MockResourceReader)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockResourceReader, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	resources := []*repository.MenuResource{
		{ID: "r1", Key: "dashboard", DisplayName: "Home", Icon: strPtr("home"), SortOrder: 0, Scope: "system"},
		{ID: "r2", Key: "materials", DisplayName: "Materials", Icon: strPtr("folder"), SortOrder: 1, Scope: "school"},
	}
	mappings := []*repository.ResourceScreenMapping{
		{ResourceKey: "dashboard", ScreenKey: "dashboard-teacher", IsDefault: true},
		{ResourceKey: "materials", ScreenKey: "materials-list", IsDefault: true},
	}

	mockResourceReader.On("GetMenuResources", ctx).Return(resources, nil)
	mockResourceReader.On("GetResourceScreenMappings", ctx, []string{"dashboard", "materials"}).Return(mappings, nil)

	permissions := []string{"materials:read"}

	// Act - platform desktop
	result, err := svc.GetNavigationConfig(ctx, userID, "desktop", permissions)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Empty(t, result.BottomNav, "desktop no debe tener bottomNav")
	assert.Len(t, result.DrawerItems, 2, "desktop debe tener todo en drawer")

	mockResourceReader.AssertExpectations(t)
}

func TestScreenService_GetNavigationConfig_EmptyResources(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockResourceReader := new(MockResourceReader)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockResourceReader, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	mockResourceReader.On("GetMenuResources", ctx).Return([]*repository.MenuResource{}, nil)

	permissions := []string{"materials:read"}

	// Act
	result, err := svc.GetNavigationConfig(ctx, userID, "mobile", permissions)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Empty(t, result.BottomNav)
	assert.Empty(t, result.DrawerItems)
	assert.Equal(t, 1, result.Version)

	mockResourceReader.AssertExpectations(t)
}

func TestScreenService_GetNavigationConfig_ResourceReaderError(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockResourceReader := new(MockResourceReader)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockResourceReader, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	dbErr := fmt.Errorf("database connection failed")
	mockResourceReader.On("GetMenuResources", ctx).Return(nil, dbErr)
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	permissions := []string{"materials:read"}

	// Act
	result, err := svc.GetNavigationConfig(ctx, userID, "mobile", permissions)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "get menu resources")

	mockResourceReader.AssertExpectations(t)
}

func TestScreenService_GetNavigationConfig_WithParentChild(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockResourceReader := new(MockResourceReader)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockResourceReader, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	parentID := "r1"

	resources := []*repository.MenuResource{
		{ID: "r1", Key: "academic", DisplayName: "Academic", Icon: strPtr("school"), SortOrder: 0, Scope: "system"},
		{ID: "r2", Key: "materials", DisplayName: "Materials", Icon: strPtr("folder"), ParentID: &parentID, SortOrder: 0, Scope: "system"},
		{ID: "r3", Key: "assessments", DisplayName: "Assessments", Icon: strPtr("quiz"), ParentID: &parentID, SortOrder: 1, Scope: "system"},
	}
	mappings := []*repository.ResourceScreenMapping{
		{ResourceKey: "materials", ScreenKey: "materials-list", IsDefault: true},
		{ResourceKey: "assessments", ScreenKey: "assessments-list", IsDefault: true},
	}

	mockResourceReader.On("GetMenuResources", ctx).Return(resources, nil)
	mockResourceReader.On("GetResourceScreenMappings", ctx, []string{"academic", "materials", "assessments"}).Return(mappings, nil)

	permissions := []string{}

	// Act
	result, err := svc.GetNavigationConfig(ctx, userID, "mobile", permissions)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	// Debe haber 1 item top-level con 2 hijos
	assert.Len(t, result.BottomNav, 1)
	assert.Equal(t, "academic", result.BottomNav[0].Key)
	assert.Len(t, result.BottomNav[0].Children, 2)
	assert.Equal(t, "materials", result.BottomNav[0].Children[0].Key)
	assert.Equal(t, "assessments", result.BottomNav[0].Children[1].Key)

	mockResourceReader.AssertExpectations(t)
}

