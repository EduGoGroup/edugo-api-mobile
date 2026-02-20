package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/screenconfig"
)

// ============================================
// Mock: ScreenService
// ============================================

type MockScreenService struct {
	GetScreenFunc             func(ctx context.Context, screenKey string, userID uuid.UUID, platform string) (*dto.CombinedScreenDTO, error)
	GetNavigationConfigFunc   func(ctx context.Context, userID uuid.UUID, platform string, permissions []string) (*service.NavigationConfigDTO, error)
	SaveUserPreferencesFunc   func(ctx context.Context, screenKey string, userID uuid.UUID, prefs json.RawMessage) error
	GetScreensForResourceFunc func(ctx context.Context, resourceKey string) ([]*dto.ResourceScreenDTO, error)
}

func (m *MockScreenService) GetScreen(ctx context.Context, screenKey string, userID uuid.UUID, platform string) (*dto.CombinedScreenDTO, error) {
	if m.GetScreenFunc != nil {
		return m.GetScreenFunc(ctx, screenKey, userID, platform)
	}
	return &dto.CombinedScreenDTO{ScreenKey: screenKey}, nil
}

func (m *MockScreenService) GetNavigationConfig(ctx context.Context, userID uuid.UUID, platform string, permissions []string) (*service.NavigationConfigDTO, error) {
	if m.GetNavigationConfigFunc != nil {
		return m.GetNavigationConfigFunc(ctx, userID, platform, permissions)
	}
	return &service.NavigationConfigDTO{
		BottomNav:   []service.NavItemDTO{},
		DrawerItems: []service.NavItemDTO{},
		Version:     1,
	}, nil
}

func (m *MockScreenService) SaveUserPreferences(ctx context.Context, screenKey string, userID uuid.UUID, prefs json.RawMessage) error {
	if m.SaveUserPreferencesFunc != nil {
		return m.SaveUserPreferencesFunc(ctx, screenKey, userID, prefs)
	}
	return nil
}

func (m *MockScreenService) GetScreensForResource(ctx context.Context, resourceKey string) ([]*dto.ResourceScreenDTO, error) {
	if m.GetScreensForResourceFunc != nil {
		return m.GetScreensForResourceFunc(ctx, resourceKey)
	}
	return []*dto.ResourceScreenDTO{}, nil
}

// ============================================
// Helper: crear CombinedScreenDTO de prueba
// ============================================

func newTestScreenDTO(screenKey string) *dto.CombinedScreenDTO {
	return &dto.CombinedScreenDTO{
		ScreenID:        "si-" + screenKey,
		ScreenKey:       screenKey,
		ScreenName:      "Test Screen",
		Pattern:         "list",
		Version:         1,
		Template:        json.RawMessage(`{"zones": []}`),
		DataEndpoint:    "/v1/materials",
		DataConfig:      json.RawMessage(`{"method": "GET"}`),
		Actions:         json.RawMessage(`[]`),
		UserPreferences: json.RawMessage(`{}`),
		UpdatedAt:       time.Date(2026, 2, 14, 10, 0, 0, 0, time.UTC),
	}
}

// ============================================
// Tests: NewScreenHandler
// ============================================

func TestNewScreenHandler(t *testing.T) {
	mockService := &MockScreenService{}
	logger := NewTestLogger()

	handler := NewScreenHandler(mockService, logger)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.screenService)
	assert.Equal(t, logger, handler.logger)
}

// ============================================
// Tests: GetScreen
// ============================================

func TestScreenHandler_GetScreen_Success(t *testing.T) {
	// Arrange
	testUserID := uuid.New().String()
	screenKey := "materials-list"
	expectedScreen := newTestScreenDTO(screenKey)

	mockService := &MockScreenService{
		GetScreenFunc: func(ctx context.Context, sk string, uid uuid.UUID, platform string) (*dto.CombinedScreenDTO, error) {
			assert.Equal(t, screenKey, sk)
			assert.Equal(t, testUserID, uid.String())
			assert.Empty(t, platform)
			return expectedScreen, nil
		},
	}

	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/v1/screens/:screenKey", MockAuthMiddleware(testUserID, "school-1"), handler.GetScreen)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/screens/"+screenKey, nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.CombinedScreenDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, screenKey, response.ScreenKey)
	assert.Equal(t, "Test Screen", response.ScreenName)
	assert.Equal(t, screenconfig.Pattern("list"), response.Pattern)

	// Verificar headers de cache
	assert.NotEmpty(t, w.Header().Get("ETag"))
	assert.NotEmpty(t, w.Header().Get("Last-Modified"))
	assert.Equal(t, "max-age=3600", w.Header().Get("Cache-Control"))
}

func TestScreenHandler_GetScreen_NotFound(t *testing.T) {
	// Arrange
	testUserID := uuid.New().String()

	mockService := &MockScreenService{
		GetScreenFunc: func(ctx context.Context, sk string, uid uuid.UUID, platform string) (*dto.CombinedScreenDTO, error) {
			return nil, errors.NewNotFoundError("screen")
		},
	}

	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/v1/screens/:screenKey", MockAuthMiddleware(testUserID, "school-1"), handler.GetScreen)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/screens/nonexistent", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "NOT_FOUND")
}

func TestScreenHandler_GetScreen_WithPlatform(t *testing.T) {
	// Arrange
	testUserID := uuid.New().String()
	screenKey := "materials-list"

	var capturedPlatform string
	mockService := &MockScreenService{
		GetScreenFunc: func(ctx context.Context, sk string, uid uuid.UUID, platform string) (*dto.CombinedScreenDTO, error) {
			capturedPlatform = platform
			return newTestScreenDTO(sk), nil
		},
	}

	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/v1/screens/:screenKey", MockAuthMiddleware(testUserID, "school-1"), handler.GetScreen)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/screens/"+screenKey+"?platform=desktop", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "desktop", capturedPlatform)
}

func TestScreenHandler_GetScreen_ETag_ConditionalRequest_304(t *testing.T) {
	// Arrange
	testUserID := uuid.New().String()
	screenKey := "materials-list"
	screenDTO := newTestScreenDTO(screenKey)

	mockService := &MockScreenService{
		GetScreenFunc: func(ctx context.Context, sk string, uid uuid.UUID, platform string) (*dto.CombinedScreenDTO, error) {
			return screenDTO, nil
		},
	}

	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/v1/screens/:screenKey", MockAuthMiddleware(testUserID, "school-1"), handler.GetScreen)

	// Act - Primera peticion para obtener el ETag
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/v1/screens/"+screenKey, nil)
	router.ServeHTTP(w1, req1)

	require.Equal(t, http.StatusOK, w1.Code)
	etag := w1.Header().Get("ETag")
	require.NotEmpty(t, etag, "primera respuesta debe incluir ETag")

	// Act - Segunda peticion con If-None-Match
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/v1/screens/"+screenKey, nil)
	req2.Header.Set("If-None-Match", etag)
	router.ServeHTTP(w2, req2)

	// Assert
	assert.Equal(t, http.StatusNotModified, w2.Code)
	assert.Empty(t, w2.Body.String(), "304 no debe tener body")
}

func TestScreenHandler_GetScreen_ETag_DifferentContent_200(t *testing.T) {
	// Arrange
	testUserID := uuid.New().String()

	mockService := &MockScreenService{
		GetScreenFunc: func(ctx context.Context, sk string, uid uuid.UUID, platform string) (*dto.CombinedScreenDTO, error) {
			return newTestScreenDTO(sk), nil
		},
	}

	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/v1/screens/:screenKey", MockAuthMiddleware(testUserID, "school-1"), handler.GetScreen)

	// Act - peticion con ETag diferente
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/screens/materials-list", nil)
	req.Header.Set("If-None-Match", `"different-etag"`)
	router.ServeHTTP(w, req)

	// Assert - debe retornar 200 porque el ETag no coincide
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestScreenHandler_GetScreen_InvalidUserID(t *testing.T) {
	// Arrange - user_id no es UUID valido
	mockService := &MockScreenService{}
	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/v1/screens/:screenKey", MockAuthMiddleware("not-a-uuid", "school-1"), handler.GetScreen)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/screens/materials-list", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "INVALID_USER_ID")
}

func TestScreenHandler_GetScreen_InternalError(t *testing.T) {
	// Arrange
	testUserID := uuid.New().String()

	mockService := &MockScreenService{
		GetScreenFunc: func(ctx context.Context, sk string, uid uuid.UUID, platform string) (*dto.CombinedScreenDTO, error) {
			return nil, fmt.Errorf("unexpected database error")
		},
	}

	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/v1/screens/:screenKey", MockAuthMiddleware(testUserID, "school-1"), handler.GetScreen)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/screens/materials-list", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}

// ============================================
// Tests: GetScreensForResource
// ============================================

func TestScreenHandler_GetScreensForResource_Success(t *testing.T) {
	// Arrange
	resourceKey := "materials"
	expectedScreens := []*dto.ResourceScreenDTO{
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
			ScreenKey:   "materials-detail",
			ScreenType:  "detail",
			IsDefault:   false,
		},
	}

	mockService := &MockScreenService{
		GetScreensForResourceFunc: func(ctx context.Context, rk string) ([]*dto.ResourceScreenDTO, error) {
			assert.Equal(t, resourceKey, rk)
			return expectedScreens, nil
		},
	}

	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/v1/screens/resource/:resourceKey", handler.GetScreensForResource)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/screens/resource/"+resourceKey, nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []*dto.ResourceScreenDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response, 2)
	assert.Equal(t, "materials-list", response[0].ScreenKey)
	assert.True(t, response[0].IsDefault)
}

func TestScreenHandler_GetScreensForResource_ServiceError(t *testing.T) {
	// Arrange
	mockService := &MockScreenService{
		GetScreensForResourceFunc: func(ctx context.Context, rk string) ([]*dto.ResourceScreenDTO, error) {
			return nil, errors.NewDatabaseError("get screens", assert.AnError)
		},
	}

	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/v1/screens/resource/:resourceKey", handler.GetScreensForResource)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/screens/resource/materials", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "DATABASE_ERROR")
}

// ============================================
// Tests: GetNavigation
// ============================================

func TestScreenHandler_GetNavigation_Success(t *testing.T) {
	// Arrange
	testUserID := uuid.New().String()

	expectedNav := &service.NavigationConfigDTO{
		BottomNav: []service.NavItemDTO{
			{Key: "dashboard", Label: "Home", Icon: "home", ScreenKey: "dashboard-teacher", SortOrder: 0},
			{Key: "materials", Label: "Materials", Icon: "folder", ScreenKey: "materials-list", SortOrder: 1},
		},
		DrawerItems: []service.NavItemDTO{},
		Version:     3,
	}

	mockService := &MockScreenService{
		GetNavigationConfigFunc: func(ctx context.Context, uid uuid.UUID, platform string, permissions []string) (*service.NavigationConfigDTO, error) {
			assert.Equal(t, testUserID, uid.String())
			return expectedNav, nil
		},
	}

	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/v1/screens/navigation", MockAuthMiddleware(testUserID, "school-1"), handler.GetNavigation)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/screens/navigation", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response service.NavigationConfigDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response.BottomNav, 2)
	assert.Equal(t, "dashboard", response.BottomNav[0].Key)
	assert.Equal(t, "dashboard-teacher", response.BottomNav[0].ScreenKey)
	assert.Equal(t, 3, response.Version)
}

func TestScreenHandler_GetNavigation_InvalidUserID(t *testing.T) {
	// Arrange
	mockService := &MockScreenService{}
	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/v1/screens/navigation", MockAuthMiddleware("invalid-uuid", "school-1"), handler.GetNavigation)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/screens/navigation", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "INVALID_USER_ID")
}

func TestScreenHandler_GetNavigation_ServiceError(t *testing.T) {
	// Arrange
	testUserID := uuid.New().String()

	mockService := &MockScreenService{
		GetNavigationConfigFunc: func(ctx context.Context, uid uuid.UUID, platform string, permissions []string) (*service.NavigationConfigDTO, error) {
			return nil, fmt.Errorf("unexpected error")
		},
	}

	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/v1/screens/navigation", MockAuthMiddleware(testUserID, "school-1"), handler.GetNavigation)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/screens/navigation", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}

// ============================================
// Tests: SavePreferences
// ============================================

func TestScreenHandler_SavePreferences_Success(t *testing.T) {
	// Arrange
	testUserID := uuid.New().String()
	screenKey := "materials-list"

	mockService := &MockScreenService{
		SaveUserPreferencesFunc: func(ctx context.Context, sk string, uid uuid.UUID, prefs json.RawMessage) error {
			assert.Equal(t, screenKey, sk)
			assert.Equal(t, testUserID, uid.String())

			var p map[string]interface{}
			err := json.Unmarshal(prefs, &p)
			require.NoError(t, err)
			assert.Equal(t, "title", p["sortBy"])
			return nil
		},
	}

	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.PUT("/v1/screens/:screenKey/preferences", MockAuthMiddleware(testUserID, "school-1"), handler.SavePreferences)

	body := `{"sortBy": "title", "viewMode": "grid"}`

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/screens/"+screenKey+"/preferences", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestScreenHandler_SavePreferences_InvalidBody(t *testing.T) {
	// Arrange
	testUserID := uuid.New().String()

	mockService := &MockScreenService{}
	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.PUT("/v1/screens/:screenKey/preferences", MockAuthMiddleware(testUserID, "school-1"), handler.SavePreferences)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/screens/materials-list/preferences", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
}

func TestScreenHandler_SavePreferences_InvalidUserID(t *testing.T) {
	// Arrange
	mockService := &MockScreenService{}
	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.PUT("/v1/screens/:screenKey/preferences", MockAuthMiddleware("not-a-uuid", "school-1"), handler.SavePreferences)

	body := `{"sortBy": "title"}`

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/screens/materials-list/preferences", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "INVALID_USER_ID")
}

func TestScreenHandler_SavePreferences_ServiceError(t *testing.T) {
	// Arrange
	testUserID := uuid.New().String()

	mockService := &MockScreenService{
		SaveUserPreferencesFunc: func(ctx context.Context, sk string, uid uuid.UUID, prefs json.RawMessage) error {
			return errors.NewDatabaseError("save preferences", assert.AnError)
		},
	}

	logger := NewTestLogger()
	handler := NewScreenHandler(mockService, logger)

	router := SetupTestRouter()
	router.PUT("/v1/screens/:screenKey/preferences", MockAuthMiddleware(testUserID, "school-1"), handler.SavePreferences)

	body := `{"sortBy": "title"}`

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/screens/materials-list/preferences", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "DATABASE_ERROR")
}
