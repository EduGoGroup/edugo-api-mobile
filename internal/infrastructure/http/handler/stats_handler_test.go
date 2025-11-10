package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// TestNewStatsHandler verifica el constructor del handler
func TestNewStatsHandler(t *testing.T) {
	// Arrange
	mockService := &MockStatsService{}
	logger := NewTestLogger()

	// Act
	handler := NewStatsHandler(mockService, logger)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.statsService)
	assert.Equal(t, logger, handler.logger)
}

// TestStatsHandler_GetMaterialStats_Success verifica obtención exitosa de estadísticas de material
func TestStatsHandler_GetMaterialStats_Success(t *testing.T) {
	// Arrange
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	expectedStats := &service.MaterialStats{
		TotalViews:    150,
		AvgProgress:   67.5,
		TotalAttempts: 45,
		AvgScore:      78.3,
	}

	mockService := &MockStatsService{
		GetMaterialStatsFunc: func(ctx context.Context, matID string) (*service.MaterialStats, error) {
			assert.Equal(t, materialID, matID)
			return expectedStats, nil
		},
	}

	logger := NewTestLogger()
	handler := NewStatsHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id/stats", handler.GetMaterialStats)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/materials/%s/stats", materialID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response service.MaterialStats
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedStats.TotalViews, response.TotalViews)
	assert.Equal(t, expectedStats.AvgProgress, response.AvgProgress)
	assert.Equal(t, expectedStats.TotalAttempts, response.TotalAttempts)
	assert.Equal(t, expectedStats.AvgScore, response.AvgScore)
}

// TestStatsHandler_GetMaterialStats_MaterialNotFound verifica manejo de material inexistente
func TestStatsHandler_GetMaterialStats_MaterialNotFound(t *testing.T) {
	// Arrange
	materialID := "550e8400-e29b-41d4-a716-446655440000"

	mockService := &MockStatsService{
		GetMaterialStatsFunc: func(ctx context.Context, matID string) (*service.MaterialStats, error) {
			return nil, errors.NewNotFoundError("material not found")
		},
	}

	logger := NewTestLogger()
	handler := NewStatsHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id/stats", handler.GetMaterialStats)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/materials/%s/stats", materialID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "material not found")
	assert.Contains(t, w.Body.String(), "NOT_FOUND")
}

// TestStatsHandler_GetMaterialStats_InvalidMaterialID verifica manejo de UUID inválido
func TestStatsHandler_GetMaterialStats_InvalidMaterialID(t *testing.T) {
	// Arrange
	invalidID := "not-a-valid-uuid"

	mockService := &MockStatsService{
		GetMaterialStatsFunc: func(ctx context.Context, matID string) (*service.MaterialStats, error) {
			return nil, errors.NewValidationError("invalid material_id")
		},
	}

	logger := NewTestLogger()
	handler := NewStatsHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id/stats", handler.GetMaterialStats)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/materials/%s/stats", invalidID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid material_id")
	assert.Contains(t, w.Body.String(), "VALIDATION_ERROR")
}

// TestStatsHandler_GetMaterialStats_ServiceError verifica manejo de errores internos del servicio
func TestStatsHandler_GetMaterialStats_ServiceError(t *testing.T) {
	// Arrange
	materialID := "550e8400-e29b-41d4-a716-446655440000"

	mockService := &MockStatsService{
		GetMaterialStatsFunc: func(ctx context.Context, matID string) (*service.MaterialStats, error) {
			return nil, fmt.Errorf("database connection failed")
		},
	}

	logger := NewTestLogger()
	handler := NewStatsHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id/stats", handler.GetMaterialStats)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/materials/%s/stats", materialID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "internal error")
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}

// TestStatsHandler_GetMaterialStats_WithZeroValues verifica manejo de estadísticas con valores en cero
func TestStatsHandler_GetMaterialStats_WithZeroValues(t *testing.T) {
	// Arrange
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	expectedStats := &service.MaterialStats{
		TotalViews:    0,
		AvgProgress:   0.0,
		TotalAttempts: 0,
		AvgScore:      0.0,
	}

	mockService := &MockStatsService{
		GetMaterialStatsFunc: func(ctx context.Context, matID string) (*service.MaterialStats, error) {
			return expectedStats, nil
		},
	}

	logger := NewTestLogger()
	handler := NewStatsHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id/stats", handler.GetMaterialStats)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/materials/%s/stats", materialID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response service.MaterialStats
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 0, response.TotalViews)
	assert.Equal(t, 0.0, response.AvgProgress)
	assert.Equal(t, 0, response.TotalAttempts)
	assert.Equal(t, 0.0, response.AvgScore)
}

// TestStatsHandler_GetGlobalStats_Success verifica obtención exitosa de estadísticas globales
func TestStatsHandler_GetGlobalStats_Success(t *testing.T) {
	// Arrange
	now := time.Now()
	expectedStats := &dto.GlobalStatsDTO{
		TotalPublishedMaterials:   150,
		TotalCompletedAssessments: 1250,
		AverageAssessmentScore:    78.5,
		ActiveUsersLast30Days:     320,
		AverageProgress:           65.3,
		GeneratedAt:               now,
	}

	mockService := &MockStatsService{
		GetGlobalStatsFunc: func(ctx context.Context) (*dto.GlobalStatsDTO, error) {
			return expectedStats, nil
		},
	}

	logger := NewTestLogger()
	handler := NewStatsHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/stats/global", handler.GetGlobalStats)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stats/global", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.GlobalStatsDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedStats.TotalPublishedMaterials, response.TotalPublishedMaterials)
	assert.Equal(t, expectedStats.TotalCompletedAssessments, response.TotalCompletedAssessments)
	assert.Equal(t, expectedStats.AverageAssessmentScore, response.AverageAssessmentScore)
	assert.Equal(t, expectedStats.ActiveUsersLast30Days, response.ActiveUsersLast30Days)
	assert.Equal(t, expectedStats.AverageProgress, response.AverageProgress)
}

// TestStatsHandler_GetGlobalStats_ServiceError verifica manejo de errores del servicio
func TestStatsHandler_GetGlobalStats_ServiceError(t *testing.T) {
	// Arrange
	mockService := &MockStatsService{
		GetGlobalStatsFunc: func(ctx context.Context) (*dto.GlobalStatsDTO, error) {
			return nil, errors.NewInternalError("database connection failed", fmt.Errorf("connection error"))
		},
	}

	logger := NewTestLogger()
	handler := NewStatsHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/stats/global", handler.GetGlobalStats)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stats/global", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "database connection failed")
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}

// TestStatsHandler_GetGlobalStats_WithZeroValues verifica estadísticas globales con valores en cero
func TestStatsHandler_GetGlobalStats_WithZeroValues(t *testing.T) {
	// Arrange
	now := time.Now()
	expectedStats := &dto.GlobalStatsDTO{
		TotalPublishedMaterials:   0,
		TotalCompletedAssessments: 0,
		AverageAssessmentScore:    0.0,
		ActiveUsersLast30Days:     0,
		AverageProgress:           0.0,
		GeneratedAt:               now,
	}

	mockService := &MockStatsService{
		GetGlobalStatsFunc: func(ctx context.Context) (*dto.GlobalStatsDTO, error) {
			return expectedStats, nil
		},
	}

	logger := NewTestLogger()
	handler := NewStatsHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/stats/global", handler.GetGlobalStats)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stats/global", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.GlobalStatsDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, int64(0), response.TotalPublishedMaterials)
	assert.Equal(t, int64(0), response.TotalCompletedAssessments)
	assert.Equal(t, 0.0, response.AverageAssessmentScore)
	assert.Equal(t, int64(0), response.ActiveUsersLast30Days)
	assert.Equal(t, 0.0, response.AverageProgress)
}

// TestStatsHandler_GetGlobalStats_GenericError verifica manejo de errores genéricos
func TestStatsHandler_GetGlobalStats_GenericError(t *testing.T) {
	// Arrange
	mockService := &MockStatsService{
		GetGlobalStatsFunc: func(ctx context.Context) (*dto.GlobalStatsDTO, error) {
			return nil, fmt.Errorf("unexpected error")
		},
	}

	logger := NewTestLogger()
	handler := NewStatsHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/stats/global", handler.GetGlobalStats)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stats/global", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error interno del servidor")
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}

// TestStatsHandler_GetMaterialStats_DifferentMaterialIDs verifica múltiples materiales
func TestStatsHandler_GetMaterialStats_DifferentMaterialIDs(t *testing.T) {
	// Arrange
	testCases := []struct {
		name       string
		materialID string
		stats      *service.MaterialStats
	}{
		{
			name:       "material con alta actividad",
			materialID: "550e8400-e29b-41d4-a716-446655440001",
			stats: &service.MaterialStats{
				TotalViews:    500,
				AvgProgress:   85.0,
				TotalAttempts: 200,
				AvgScore:      90.5,
			},
		},
		{
			name:       "material con baja actividad",
			materialID: "550e8400-e29b-41d4-a716-446655440002",
			stats: &service.MaterialStats{
				TotalViews:    10,
				AvgProgress:   25.0,
				TotalAttempts: 3,
				AvgScore:      55.0,
			},
		},
		{
			name:       "material nuevo sin actividad",
			materialID: "550e8400-e29b-41d4-a716-446655440003",
			stats: &service.MaterialStats{
				TotalViews:    0,
				AvgProgress:   0.0,
				TotalAttempts: 0,
				AvgScore:      0.0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockService := &MockStatsService{
				GetMaterialStatsFunc: func(ctx context.Context, matID string) (*service.MaterialStats, error) {
					assert.Equal(t, tc.materialID, matID)
					return tc.stats, nil
				},
			}

			logger := NewTestLogger()
			handler := NewStatsHandler(mockService, logger)

			router := SetupTestRouter()
			router.GET("/materials/:id/stats", handler.GetMaterialStats)

			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/materials/%s/stats", tc.materialID), nil)
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, http.StatusOK, w.Code)

			var response service.MaterialStats
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tc.stats.TotalViews, response.TotalViews)
			assert.Equal(t, tc.stats.AvgProgress, response.AvgProgress)
			assert.Equal(t, tc.stats.TotalAttempts, response.TotalAttempts)
			assert.Equal(t, tc.stats.AvgScore, response.AvgScore)
		})
	}
}
