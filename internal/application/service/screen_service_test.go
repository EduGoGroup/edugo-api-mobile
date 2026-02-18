package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
)

// ============================================
// Mock: ScreenRepository
// ============================================

type MockScreenRepository struct {
	mock.Mock
}

func (m *MockScreenRepository) GetCombinedScreen(ctx context.Context, screenKey string, userID uuid.UUID) (*repository.CombinedScreen, error) {
	args := m.Called(ctx, screenKey, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.CombinedScreen), args.Error(1)
}

func (m *MockScreenRepository) GetScreensForResource(ctx context.Context, resourceKey string) ([]*repository.ResourceScreenInfo, error) {
	args := m.Called(ctx, resourceKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repository.ResourceScreenInfo), args.Error(1)
}

func (m *MockScreenRepository) GetUserPreferences(ctx context.Context, screenKey string, userID uuid.UUID) (json.RawMessage, error) {
	args := m.Called(ctx, screenKey, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(json.RawMessage), args.Error(1)
}

func (m *MockScreenRepository) SaveUserPreferences(ctx context.Context, screenKey string, userID uuid.UUID, prefs json.RawMessage) error {
	args := m.Called(ctx, screenKey, userID, prefs)
	return args.Error(0)
}

// ============================================
// Helper: crear CombinedScreen de prueba
// ============================================

func newTestCombinedScreen(screenKey string) *repository.CombinedScreen {
	definition := json.RawMessage(`{
		"navigation": {"topBar": {"title": "slot:page_title"}},
		"zones": [
			{
				"id": "list_content",
				"type": "simple-list",
				"distribution": "stacked"
			}
		],
		"platformOverrides": {
			"desktop": {
				"zones": {
					"list_content": {"distribution": "grid", "columns": 2}
				}
			}
		}
	}`)

	slotData := json.RawMessage(`{"page_title": "Materials List"}`)
	actions := json.RawMessage(`[{"id": "item-click", "trigger": "item_click", "type": "NAVIGATE"}]`)
	dataConfig := json.RawMessage(`{"method": "GET", "pagination": {"type": "offset"}}`)

	return &repository.CombinedScreen{
		ID:              "si-materials-list",
		ScreenKey:       screenKey,
		Name:            "Educational Materials",
		Pattern:         "list",
		Version:         1,
		Definition:      definition,
		SlotData:        slotData,
		Actions:         actions,
		DataEndpoint:    "/v1/materials",
		DataConfig:      dataConfig,
		UserPreferences: json.RawMessage(`{}`),
		LastUpdated:     time.Date(2026, 2, 14, 10, 0, 0, 0, time.UTC),
	}
}

// ============================================
// Tests: GetScreen
// ============================================

func TestScreenService_GetScreen_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	screenKey := "materials-list"

	combined := newTestCombinedScreen(screenKey)

	mockRepo.On("GetCombinedScreen", ctx, screenKey, userID).Return(combined, nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	// Act
	result, err := svc.GetScreen(ctx, screenKey, userID, "")

	// Assert
	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "si-materials-list", result.ScreenID)
	assert.Equal(t, screenKey, result.ScreenKey)
	assert.Equal(t, "Educational Materials", result.ScreenName)
	assert.Equal(t, "list", result.Pattern)
	assert.Equal(t, 1, result.Version)
	assert.Equal(t, "/v1/materials", result.DataEndpoint)

	mockRepo.AssertExpectations(t)
}

func TestScreenService_GetScreen_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	screenKey := "nonexistent-screen"

	mockRepo.On("GetCombinedScreen", ctx, screenKey, userID).Return(nil, nil)
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	// Act
	result, err := svc.GetScreen(ctx, screenKey, userID, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")

	mockRepo.AssertExpectations(t)
}

func TestScreenService_GetScreen_DatabaseError(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	screenKey := "materials-list"

	dbErr := fmt.Errorf("database connection failed")
	mockRepo.On("GetCombinedScreen", ctx, screenKey, userID).Return(nil, dbErr)
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	// Act
	result, err := svc.GetScreen(ctx, screenKey, userID, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestScreenService_GetScreen_ResolvesSlots(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	screenKey := "materials-list"

	combined := newTestCombinedScreen(screenKey)
	// Definition contiene "slot:page_title" que debe resolverse a "Materials List"

	mockRepo.On("GetCombinedScreen", ctx, screenKey, userID).Return(combined, nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	// Act
	result, err := svc.GetScreen(ctx, screenKey, userID, "")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verificar que slot:page_title fue reemplazado por "Materials List"
	templateStr := string(result.Template)
	assert.NotContains(t, templateStr, "slot:page_title")
	assert.Contains(t, templateStr, "Materials List")

	mockRepo.AssertExpectations(t)
}

func TestScreenService_GetScreen_AppliesPlatformOverrides(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	screenKey := "materials-list"

	combined := newTestCombinedScreen(screenKey)

	mockRepo.On("GetCombinedScreen", ctx, screenKey, userID).Return(combined, nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	// Act
	result, err := svc.GetScreen(ctx, screenKey, userID, "desktop")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verificar que el override de desktop se aplico (distribution -> grid)
	var templateMap map[string]interface{}
	err = json.Unmarshal(result.Template, &templateMap)
	require.NoError(t, err)

	zones, ok := templateMap["zones"].([]interface{})
	require.True(t, ok)

	// Buscar zone "list_content"
	found := false
	for _, zone := range zones {
		zoneMap, ok := zone.(map[string]interface{})
		if !ok {
			continue
		}
		if zoneMap["id"] == "list_content" {
			assert.Equal(t, "grid", zoneMap["distribution"])
			assert.Equal(t, float64(2), zoneMap["columns"])
			found = true
			break
		}
	}
	assert.True(t, found, "zone list_content debe existir con overrides aplicados")

	// platformOverrides debe haberse eliminado del resultado
	_, hasPlatformOverrides := templateMap["platformOverrides"]
	assert.False(t, hasPlatformOverrides, "platformOverrides debe eliminarse del resultado")

	mockRepo.AssertExpectations(t)
}

func TestScreenService_GetScreen_NoPlatformOverrides_WhenNotRequested(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	screenKey := "materials-list"

	combined := newTestCombinedScreen(screenKey)

	mockRepo.On("GetCombinedScreen", ctx, screenKey, userID).Return(combined, nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	// Act - sin platform
	result, err := svc.GetScreen(ctx, screenKey, userID, "")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	// Sin platform, el template mantiene platformOverrides intacto
	var templateMap map[string]interface{}
	err = json.Unmarshal(result.Template, &templateMap)
	require.NoError(t, err)

	zones, ok := templateMap["zones"].([]interface{})
	require.True(t, ok)

	// zone list_content debe mantener distribution: stacked
	for _, zone := range zones {
		zoneMap, ok := zone.(map[string]interface{})
		if !ok {
			continue
		}
		if zoneMap["id"] == "list_content" {
			assert.Equal(t, "stacked", zoneMap["distribution"])
			break
		}
	}

	mockRepo.AssertExpectations(t)
}

func TestScreenService_GetScreen_CacheHit(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	screenKey := "materials-list"

	combined := newTestCombinedScreen(screenKey)
	userPrefs := json.RawMessage(`{"theme": "dark"}`)

	// Primera llamada: carga de BD
	mockRepo.On("GetCombinedScreen", ctx, screenKey, userID).Return(combined, nil).Once()
	// Segunda llamada: solo busca preferencias de usuario (cache hit)
	mockRepo.On("GetUserPreferences", ctx, screenKey, userID).Return(userPrefs, nil).Once()
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	// Act - primera llamada (cache miss -> BD)
	result1, err1 := svc.GetScreen(ctx, screenKey, userID, "")
	require.NoError(t, err1)
	require.NotNil(t, result1)

	// Act - segunda llamada (cache hit)
	result2, err2 := svc.GetScreen(ctx, screenKey, userID, "")
	require.NoError(t, err2)
	require.NotNil(t, result2)

	// Assert
	assert.Equal(t, result1.ScreenKey, result2.ScreenKey)
	// El repo.GetCombinedScreen solo debio llamarse 1 vez
	mockRepo.AssertNumberOfCalls(t, "GetCombinedScreen", 1)
	// GetUserPreferences se llama en la segunda llamada (cache hit con fusion de prefs)
	mockRepo.AssertCalled(t, "GetUserPreferences", ctx, screenKey, userID)
}

func TestScreenService_GetScreen_MergesUserPreferences(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	screenKey := "materials-list"

	combined := newTestCombinedScreen(screenKey)
	combined.UserPreferences = json.RawMessage(`{"sortBy": "date"}`)

	mockRepo.On("GetCombinedScreen", ctx, screenKey, userID).Return(combined, nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	// Act
	result, err := svc.GetScreen(ctx, screenKey, userID, "")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	var prefs map[string]interface{}
	err = json.Unmarshal(result.UserPreferences, &prefs)
	require.NoError(t, err)
	assert.Equal(t, "date", prefs["sortBy"])

	mockRepo.AssertExpectations(t)
}

// ============================================
// Tests: GetNavigationConfig
// ============================================

func TestScreenService_GetNavigationConfig_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	// Act
	result, err := svc.GetNavigationConfig(ctx, userID, "mobile")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.BottomNav)
	assert.NotNil(t, result.DrawerItems)
	assert.Greater(t, result.Version, 0)

	// Verificar que los items de navegacion tienen la estructura esperada
	for _, item := range result.BottomNav {
		assert.NotEmpty(t, item.Key)
		assert.NotEmpty(t, item.Label)
		assert.NotEmpty(t, item.Icon)
		assert.NotEmpty(t, item.ScreenKey)
	}
}

// ============================================
// Tests: SaveUserPreferences
// ============================================

func TestScreenService_SaveUserPreferences_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	screenKey := "materials-list"
	prefs := json.RawMessage(`{"sortBy": "title", "viewMode": "grid"}`)

	mockRepo.On("SaveUserPreferences", ctx, screenKey, userID, prefs).Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	// Act
	err := svc.SaveUserPreferences(ctx, screenKey, userID, prefs)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestScreenService_SaveUserPreferences_DatabaseError(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	screenKey := "materials-list"
	prefs := json.RawMessage(`{"sortBy": "title"}`)

	dbErr := fmt.Errorf("database error")
	mockRepo.On("SaveUserPreferences", ctx, screenKey, userID, prefs).Return(dbErr)
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	// Act
	err := svc.SaveUserPreferences(ctx, screenKey, userID, prefs)

	// Assert
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

// ============================================
// Tests: GetScreensForResource
// ============================================

func TestScreenService_GetScreensForResource_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger)

	ctx := context.Background()
	resourceKey := "materials"

	infos := []*repository.ResourceScreenInfo{
		{
			ResourceID:  "res-1",
			ResourceKey: "materials",
			ScreenKey:   "materials-list",
			ScreenType:  "list",
			IsDefault:   true,
		},
		{
			ResourceID:  "res-2",
			ResourceKey: "materials",
			ScreenKey:   "materials-grid",
			ScreenType:  "grid",
			IsDefault:   false,
		},
	}

	mockRepo.On("GetScreensForResource", ctx, resourceKey).Return(infos, nil)

	// Act
	result, err := svc.GetScreensForResource(ctx, resourceKey)

	// Assert
	require.NoError(t, err)
	require.Len(t, result, 2)

	assert.Equal(t, "materials-list", result[0].ScreenKey)
	assert.True(t, result[0].IsDefault)
	assert.Equal(t, "materials-grid", result[1].ScreenKey)
	assert.False(t, result[1].IsDefault)

	mockRepo.AssertExpectations(t)
}

func TestScreenService_GetScreensForResource_Empty(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger)

	ctx := context.Background()
	resourceKey := "nonexistent"

	mockRepo.On("GetScreensForResource", ctx, resourceKey).Return([]*repository.ResourceScreenInfo{}, nil)

	// Act
	result, err := svc.GetScreensForResource(ctx, resourceKey)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 0)

	mockRepo.AssertExpectations(t)
}

func TestScreenService_GetScreensForResource_DatabaseError(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger)

	ctx := context.Background()
	resourceKey := "materials"

	dbErr := fmt.Errorf("database error")
	mockRepo.On("GetScreensForResource", ctx, resourceKey).Return(nil, dbErr)
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	// Act
	result, err := svc.GetScreensForResource(ctx, resourceKey)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

// ============================================
// Tests: resolveSlots (funcion interna)
// ============================================

func TestResolveSlots_ReplacesSlotReferences(t *testing.T) {
	definition := json.RawMessage(`{"title": "slot:page_title", "subtitle": "static text"}`)
	slotData := json.RawMessage(`{"page_title": "My Page Title"}`)

	result := resolveSlots(definition, slotData)

	var resultMap map[string]interface{}
	err := json.Unmarshal(result, &resultMap)
	require.NoError(t, err)

	assert.Equal(t, "My Page Title", resultMap["title"])
	assert.Equal(t, "static text", resultMap["subtitle"])
}

func TestResolveSlots_EmptySlotData(t *testing.T) {
	definition := json.RawMessage(`{"title": "slot:page_title"}`)

	// Con slotData vacio, no deberia cambiar nada
	result := resolveSlots(definition, json.RawMessage(`{}`))

	assert.Equal(t, string(definition), string(result))
}

func TestResolveSlots_NullSlotData(t *testing.T) {
	definition := json.RawMessage(`{"title": "slot:page_title"}`)

	result := resolveSlots(definition, json.RawMessage(`null`))

	assert.Equal(t, string(definition), string(result))
}

func TestResolveSlots_NestedStructure(t *testing.T) {
	definition := json.RawMessage(`{
		"nav": {
			"topBar": {
				"title": "slot:header_title"
			}
		},
		"items": [
			{"label": "slot:item_label"},
			{"label": "fixed"}
		]
	}`)
	slotData := json.RawMessage(`{
		"header_title": "Dashboard",
		"item_label": "Home"
	}`)

	result := resolveSlots(definition, slotData)

	resultStr := string(result)
	assert.NotContains(t, resultStr, "slot:header_title")
	assert.NotContains(t, resultStr, "slot:item_label")
	assert.Contains(t, resultStr, "Dashboard")
	assert.Contains(t, resultStr, "Home")
	assert.Contains(t, resultStr, "fixed")
}

func TestResolveSlots_UnknownSlotKey_KeepsOriginal(t *testing.T) {
	definition := json.RawMessage(`{"title": "slot:unknown_key"}`)
	slotData := json.RawMessage(`{"other_key": "value"}`)

	result := resolveSlots(definition, slotData)

	var resultMap map[string]interface{}
	err := json.Unmarshal(result, &resultMap)
	require.NoError(t, err)

	// Un slot:xxx que no tiene match en slotData debe mantenerse como string
	assert.Equal(t, "slot:unknown_key", resultMap["title"])
}

// ============================================
// Tests: applyPlatformOverrides (funcion interna)
// ============================================

func TestApplyPlatformOverrides_AppliesDesktop(t *testing.T) {
	definition := json.RawMessage(`{
		"zones": [
			{"id": "list_content", "distribution": "stacked"}
		],
		"platformOverrides": {
			"desktop": {
				"zones": {
					"list_content": {"distribution": "grid", "columns": 3}
				}
			}
		}
	}`)

	result := applyPlatformOverrides(definition, "desktop")

	var resultMap map[string]interface{}
	err := json.Unmarshal(result, &resultMap)
	require.NoError(t, err)

	// platformOverrides debe eliminarse
	_, exists := resultMap["platformOverrides"]
	assert.False(t, exists)

	// zone debe tener los overrides aplicados
	zones := resultMap["zones"].([]interface{})
	zone := zones[0].(map[string]interface{})
	assert.Equal(t, "grid", zone["distribution"])
	assert.Equal(t, float64(3), zone["columns"])
}

func TestApplyPlatformOverrides_NoPlatformMatch(t *testing.T) {
	definition := json.RawMessage(`{
		"zones": [
			{"id": "list_content", "distribution": "stacked"}
		],
		"platformOverrides": {
			"desktop": {
				"zones": {
					"list_content": {"distribution": "grid"}
				}
			}
		}
	}`)

	// Platform "mobile" no tiene overrides definidos
	result := applyPlatformOverrides(definition, "mobile")

	var resultMap map[string]interface{}
	err := json.Unmarshal(result, &resultMap)
	require.NoError(t, err)

	// zones no deben cambiar
	zones := resultMap["zones"].([]interface{})
	zone := zones[0].(map[string]interface{})
	assert.Equal(t, "stacked", zone["distribution"])
}

func TestApplyPlatformOverrides_NoOverridesKey(t *testing.T) {
	definition := json.RawMessage(`{
		"zones": [
			{"id": "list_content", "distribution": "stacked"}
		]
	}`)

	result := applyPlatformOverrides(definition, "desktop")

	// Sin platformOverrides, el template debe ser identico
	assert.JSONEq(t, string(definition), string(result))
}

// ============================================
// Tests: withUserPreferences (cache hit)
// ============================================

func TestScreenService_WithUserPreferences_MergesOnCacheHit(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger).(*screenService)

	ctx := context.Background()
	userID := uuid.New()
	screenKey := "materials-list"

	cached := &dto.CombinedScreenDTO{
		ScreenID:        "si-1",
		ScreenKey:       screenKey,
		UserPreferences: json.RawMessage(`{}`),
	}

	userPrefs := json.RawMessage(`{"viewMode": "compact"}`)
	mockRepo.On("GetUserPreferences", ctx, screenKey, userID).Return(userPrefs, nil)

	// Act
	result, err := svc.withUserPreferences(ctx, cached, screenKey, userID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	var prefs map[string]interface{}
	err = json.Unmarshal(result.UserPreferences, &prefs)
	require.NoError(t, err)
	assert.Equal(t, "compact", prefs["viewMode"])

	// El cached original no debe modificarse
	assert.Equal(t, `{}`, string(cached.UserPreferences))
}

func TestScreenService_WithUserPreferences_ErrorFallsBackToCached(t *testing.T) {
	// Arrange
	mockRepo := new(MockScreenRepository)
	mockLogger := new(MockLogger)

	svc := NewScreenService(mockRepo, mockLogger).(*screenService)

	ctx := context.Background()
	userID := uuid.New()
	screenKey := "materials-list"

	cached := &dto.CombinedScreenDTO{
		ScreenID:        "si-1",
		ScreenKey:       screenKey,
		UserPreferences: json.RawMessage(`{"default": true}`),
	}

	mockRepo.On("GetUserPreferences", ctx, screenKey, userID).Return(nil, fmt.Errorf("db error"))
	mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

	// Act
	result, err := svc.withUserPreferences(ctx, cached, screenKey, userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, cached, result) // Retorna el cached original sin cambios
}
