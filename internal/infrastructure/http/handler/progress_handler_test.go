package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// Tests de UpdateProgress (legacy) fueron eliminados
// El endpoint PATCH /materials/:id/progress fue removido
// Usar PUT /progress (UpsertProgress) en su lugar

// TestNewProgressHandler verifica el constructor del handler
func TestNewProgressHandler(t *testing.T) {
	// Arrange
	mockService := &MockProgressService{}
	logger := NewTestLogger()

	// Act
	handler := NewProgressHandler(mockService, logger)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.progressService)
	assert.Equal(t, logger, handler.logger)
}

// TestProgressHandler_UpsertProgress_Success verifica actualización exitosa de progreso con datos válidos
func TestProgressHandler_UpsertProgress_Success(t *testing.T) {
	// Arrange
	authenticatedUserID := "user-123"
	materialID := "material-456"
	progressPercentage := 75
	lastPage := 45

	mockService := &MockProgressService{
		UpdateProgressFunc: func(ctx context.Context, matID, usrID string, percentage, page int) error {
			assert.Equal(t, materialID, matID)
			assert.Equal(t, authenticatedUserID, usrID)
			assert.Equal(t, progressPercentage, percentage)
			assert.Equal(t, lastPage, page)
			return nil
		},
	}

	logger := NewTestLogger()
	handler := NewProgressHandler(mockService, logger)

	router := SetupTestRouter()
	router.PUT("/progress", MockUserIDMiddleware(authenticatedUserID), handler.UpsertProgress)

	reqBody := fmt.Sprintf(`{
		"user_id": "%s",
		"material_id": "%s",
		"progress_percentage": %d,
		"last_page": %d
	}`, authenticatedUserID, materialID, progressPercentage, lastPage)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/progress", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response ProgressResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, authenticatedUserID, response.UserID)
	assert.Equal(t, materialID, response.MaterialID)
	assert.Equal(t, progressPercentage, response.ProgressPercentage)
	assert.Equal(t, lastPage, response.LastPage)
	assert.Equal(t, "progress updated successfully", response.Message)
}

// TestProgressHandler_UpsertProgress_InvalidJSON verifica rechazo de JSON malformado
func TestProgressHandler_UpsertProgress_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := &MockProgressService{}
	logger := NewTestLogger()
	handler := NewProgressHandler(mockService, logger)

	router := SetupTestRouter()
	router.PUT("/progress", MockUserIDMiddleware("user-123"), handler.UpsertProgress)

	testCases := []struct {
		name string
		body string
	}{
		{
			name: "JSON malformado",
			body: `{"user_id": invalid}`,
		},
		{
			name: "JSON incompleto",
			body: `{"user_id": "user-123"`,
		},
		{
			name: "campo con tipo incorrecto",
			body: `{"user_id": "user-123", "material_id": "mat-456", "progress_percentage": "not-a-number"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/progress", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "invalid request body")
			assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
		})
	}
}

// TestProgressHandler_UpsertProgress_MissingRequiredFields verifica validación de campos requeridos
func TestProgressHandler_UpsertProgress_MissingRequiredFields(t *testing.T) {
	// Arrange
	mockService := &MockProgressService{}
	logger := NewTestLogger()
	handler := NewProgressHandler(mockService, logger)

	router := SetupTestRouter()
	router.PUT("/progress", MockUserIDMiddleware("user-123"), handler.UpsertProgress)

	testCases := []struct {
		name string
		body string
	}{
		{
			name: "sin user_id",
			body: `{"material_id": "mat-456", "progress_percentage": 75}`,
		},
		{
			name: "sin material_id",
			body: `{"user_id": "user-123", "progress_percentage": 75}`,
		},
		{
			name: "user_id vacío",
			body: `{"user_id": "", "material_id": "mat-456", "progress_percentage": 75}`,
		},
		{
			name: "material_id vacío",
			body: `{"user_id": "user-123", "material_id": "", "progress_percentage": 75}`,
		},
		{
			name: "ambos campos vacíos",
			body: `{"user_id": "", "material_id": "", "progress_percentage": 75}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/progress", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, http.StatusBadRequest, w.Code)
			// La validación puede ocurrir en dos lugares:
			// 1. En el binding de Gin (campos required)
			// 2. En la validación manual del handler (campos vacíos)
			assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
		})
	}
}

// TestProgressHandler_UpsertProgress_InvalidPercentage verifica validación de porcentaje
func TestProgressHandler_UpsertProgress_InvalidPercentage(t *testing.T) {
	// Arrange
	mockService := &MockProgressService{}
	logger := NewTestLogger()
	handler := NewProgressHandler(mockService, logger)

	router := SetupTestRouter()
	router.PUT("/progress", MockUserIDMiddleware("user-123"), handler.UpsertProgress)

	testCases := []struct {
		name       string
		percentage int
	}{
		{
			name:       "porcentaje negativo",
			percentage: -10,
		},
		{
			name:       "porcentaje mayor a 100",
			percentage: 150,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			reqBody := fmt.Sprintf(`{
				"user_id": "user-123",
				"material_id": "mat-456",
				"progress_percentage": %d
			}`, tc.percentage)

			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/progress", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "invalid request body")
			assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
		})
	}
}

// TestProgressHandler_UpsertProgress_Unauthorized verifica rechazo sin autenticación
func TestProgressHandler_UpsertProgress_Unauthorized(t *testing.T) {
	// Arrange
	mockService := &MockProgressService{}
	logger := NewTestLogger()
	handler := NewProgressHandler(mockService, logger)

	router := SetupTestRouter()
	// Sin MockUserIDMiddleware - simula usuario no autenticado
	// Agregamos un middleware que capture el panic de MustGetUserID
	router.Use(func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(http.StatusUnauthorized, ErrorResponse{
					Error: "user not authenticated",
					Code:  "UNAUTHORIZED",
				})
				c.Abort()
			}
		}()
		c.Next()
	})
	router.PUT("/progress", handler.UpsertProgress)

	reqBody := `{
		"user_id": "user-123",
		"material_id": "mat-456",
		"progress_percentage": 75
	}`

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/progress", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "user not authenticated")
	assert.Contains(t, w.Body.String(), "UNAUTHORIZED")
}

// TestProgressHandler_UpsertProgress_Forbidden verifica que usuario solo puede actualizar su propio progreso
func TestProgressHandler_UpsertProgress_Forbidden(t *testing.T) {
	// Arrange
	authenticatedUserID := "user-123"
	differentUserID := "user-456"

	mockService := &MockProgressService{}
	logger := NewTestLogger()
	handler := NewProgressHandler(mockService, logger)

	router := SetupTestRouter()
	router.PUT("/progress", MockUserIDMiddleware(authenticatedUserID), handler.UpsertProgress)

	reqBody := fmt.Sprintf(`{
		"user_id": "%s",
		"material_id": "mat-789",
		"progress_percentage": 75
	}`, differentUserID)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/progress", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "you can only update your own progress")
	assert.Contains(t, w.Body.String(), "FORBIDDEN")
}

// TestProgressHandler_UpsertProgress_MaterialNotFound verifica manejo de material inexistente
func TestProgressHandler_UpsertProgress_MaterialNotFound(t *testing.T) {
	// Arrange
	authenticatedUserID := "user-123"

	mockService := &MockProgressService{
		UpdateProgressFunc: func(ctx context.Context, materialID, userID string, percentage, lastPage int) error {
			return errors.NewNotFoundError("material not found")
		},
	}

	logger := NewTestLogger()
	handler := NewProgressHandler(mockService, logger)

	router := SetupTestRouter()
	router.PUT("/progress", MockUserIDMiddleware(authenticatedUserID), handler.UpsertProgress)

	reqBody := fmt.Sprintf(`{
		"user_id": "%s",
		"material_id": "nonexistent-material",
		"progress_percentage": 75
	}`, authenticatedUserID)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/progress", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "material not found")
	assert.Contains(t, w.Body.String(), "NOT_FOUND")
}

// TestProgressHandler_UpsertProgress_InvalidMaterialID verifica manejo de UUID inválido
func TestProgressHandler_UpsertProgress_InvalidMaterialID(t *testing.T) {
	// Arrange
	authenticatedUserID := "user-123"

	mockService := &MockProgressService{
		UpdateProgressFunc: func(ctx context.Context, materialID, userID string, percentage, lastPage int) error {
			return errors.NewValidationError("invalid material_id")
		},
	}

	logger := NewTestLogger()
	handler := NewProgressHandler(mockService, logger)

	router := SetupTestRouter()
	router.PUT("/progress", MockUserIDMiddleware(authenticatedUserID), handler.UpsertProgress)

	reqBody := fmt.Sprintf(`{
		"user_id": "%s",
		"material_id": "not-a-valid-uuid",
		"progress_percentage": 75
	}`, authenticatedUserID)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/progress", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid material_id")
	assert.Contains(t, w.Body.String(), "VALIDATION_ERROR")
}

// TestProgressHandler_UpsertProgress_ServiceError verifica manejo de errores internos del servicio
func TestProgressHandler_UpsertProgress_ServiceError(t *testing.T) {
	// Arrange
	authenticatedUserID := "user-123"

	mockService := &MockProgressService{
		UpdateProgressFunc: func(ctx context.Context, materialID, userID string, percentage, lastPage int) error {
			return fmt.Errorf("database connection failed")
		},
	}

	logger := NewTestLogger()
	handler := NewProgressHandler(mockService, logger)

	router := SetupTestRouter()
	router.PUT("/progress", MockUserIDMiddleware(authenticatedUserID), handler.UpsertProgress)

	reqBody := fmt.Sprintf(`{
		"user_id": "%s",
		"material_id": "mat-456",
		"progress_percentage": 75
	}`, authenticatedUserID)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/progress", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "internal server error")
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}

// TestProgressHandler_UpsertProgress_ValidPercentageRange verifica valores válidos de porcentaje
func TestProgressHandler_UpsertProgress_ValidPercentageRange(t *testing.T) {
	// Arrange
	authenticatedUserID := "user-123"

	testCases := []struct {
		name       string
		percentage int
	}{
		{
			name:       "porcentaje 1",
			percentage: 1,
		},
		{
			name:       "porcentaje 50",
			percentage: 50,
		},
		{
			name:       "porcentaje 100",
			percentage: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockService := &MockProgressService{
				UpdateProgressFunc: func(ctx context.Context, materialID, userID string, percentage, lastPage int) error {
					assert.Equal(t, tc.percentage, percentage)
					return nil
				},
			}

			logger := NewTestLogger()
			handler := NewProgressHandler(mockService, logger)

			router := SetupTestRouter()
			router.PUT("/progress", MockUserIDMiddleware(authenticatedUserID), handler.UpsertProgress)

			reqBody := fmt.Sprintf(`{
				"user_id": "%s",
				"material_id": "mat-456",
				"progress_percentage": %d
			}`, authenticatedUserID, tc.percentage)

			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/progress", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, http.StatusOK, w.Code)

			var response ProgressResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, tc.percentage, response.ProgressPercentage)
		})
	}
}

// TestProgressHandler_UpsertProgress_WithoutLastPage verifica que last_page es opcional
func TestProgressHandler_UpsertProgress_WithoutLastPage(t *testing.T) {
	// Arrange
	authenticatedUserID := "user-123"
	materialID := "mat-456"
	progressPercentage := 75

	mockService := &MockProgressService{
		UpdateProgressFunc: func(ctx context.Context, matID, usrID string, percentage, page int) error {
			assert.Equal(t, materialID, matID)
			assert.Equal(t, authenticatedUserID, usrID)
			assert.Equal(t, progressPercentage, percentage)
			assert.Equal(t, 0, page) // last_page debería ser 0 si no se proporciona
			return nil
		},
	}

	logger := NewTestLogger()
	handler := NewProgressHandler(mockService, logger)

	router := SetupTestRouter()
	router.PUT("/progress", MockUserIDMiddleware(authenticatedUserID), handler.UpsertProgress)

	reqBody := fmt.Sprintf(`{
		"user_id": "%s",
		"material_id": "%s",
		"progress_percentage": %d
	}`, authenticatedUserID, materialID, progressPercentage)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/progress", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response ProgressResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, 0, response.LastPage)
}

// TestProgressHandler_UpsertProgress_Idempotency verifica que múltiples llamadas con mismos datos son seguras
func TestProgressHandler_UpsertProgress_Idempotency(t *testing.T) {
	// Arrange
	authenticatedUserID := "user-123"
	materialID := "mat-456"
	progressPercentage := 75
	lastPage := 45

	callCount := 0
	mockService := &MockProgressService{
		UpdateProgressFunc: func(ctx context.Context, matID, usrID string, percentage, page int) error {
			callCount++
			return nil
		},
	}

	logger := NewTestLogger()
	handler := NewProgressHandler(mockService, logger)

	router := SetupTestRouter()
	router.PUT("/progress", MockUserIDMiddleware(authenticatedUserID), handler.UpsertProgress)

	reqBody := fmt.Sprintf(`{
		"user_id": "%s",
		"material_id": "%s",
		"progress_percentage": %d,
		"last_page": %d
	}`, authenticatedUserID, materialID, progressPercentage, lastPage)

	// Act - Llamar 3 veces con los mismos datos
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/progress", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code, "Llamada %d debería ser exitosa", i+1)
	}

	// Verificar que el servicio fue llamado 3 veces
	assert.Equal(t, 3, callCount)
}
