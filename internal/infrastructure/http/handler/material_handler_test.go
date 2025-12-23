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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// Helper para crear punteros a strings
func stringPtr(s string) *string {
	return &s
}

// TestNewMaterialHandler verifica el constructor del handler
func TestNewMaterialHandler(t *testing.T) {
	// Arrange
	mockService := &MockMaterialService{}
	mockS3 := &MockS3Storage{}
	logger := NewTestLogger()

	// Act
	handler := NewMaterialHandler(mockService, mockS3, logger)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.materialService)
	assert.Equal(t, logger, handler.logger)
}

// TestMaterialHandler_GenerateUploadURL_PathTraversalPrevention verifica prevención de path traversal
// Este es un test CRÍTICO de seguridad que previene acceso no autorizado a S3
func TestMaterialHandler_GenerateUploadURL_PathTraversalPrevention(t *testing.T) {
	// Arrange
	mockService := &MockMaterialService{
		GetMaterialFunc: func(ctx context.Context, id string) (*dto.MaterialResponse, error) {
			return &dto.MaterialResponse{ID: id, Title: "Test Material"}, nil
		},
	}

	logger := NewTestLogger()
	mockS3 := &MockS3Storage{}
	handler := NewMaterialHandler(mockService, mockS3, logger)

	router := SetupTestRouter()
	router.POST("/materials/:id/upload-url", handler.GenerateUploadURL)

	testCases := []struct {
		name        string
		fileName    string
		wantCode    int
		wantErrCode string
	}{
		{
			name:        "path traversal con .. debe rechazarse",
			fileName:    "../../../etc/passwd",
			wantCode:    http.StatusBadRequest,
			wantErrCode: "INVALID_FILENAME",
		},
		{
			name:        "path traversal con múltiples .. debe rechazarse",
			fileName:    "../../secret.txt",
			wantCode:    http.StatusBadRequest,
			wantErrCode: "INVALID_FILENAME",
		},
		{
			name:        "path traversal con / debe rechazarse",
			fileName:    "uploads/../../secret.txt",
			wantCode:    http.StatusBadRequest,
			wantErrCode: "INVALID_FILENAME",
		},
		{
			name:        "path traversal con \\ debe rechazarse",
			fileName:    "uploads\\\\..\\\\secret.txt", // Doble backslash para JSON
			wantCode:    http.StatusBadRequest,
			wantErrCode: "INVALID_FILENAME",
		},
		{
			name:        "nombre con / debe rechazarse",
			fileName:    "folder/file.pdf",
			wantCode:    http.StatusBadRequest,
			wantErrCode: "INVALID_FILENAME",
		},
		{
			name:        "nombre con \\ debe rechazarse",
			fileName:    "folder\\\\file.pdf", // Doble backslash para JSON
			wantCode:    http.StatusBadRequest,
			wantErrCode: "INVALID_FILENAME",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			reqBody := fmt.Sprintf(`{"file_name":"%s","content_type":"application/pdf"}`, tc.fileName)

			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/materials/test-id/upload-url", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.wantCode, w.Code, "Status code incorrecto para: %s", tc.name)
			if tc.wantErrCode != "" {
				assert.Contains(t, w.Body.String(), tc.wantErrCode, "Error code no encontrado en respuesta")
				assert.Contains(t, w.Body.String(), "invalid file name", "Mensaje de error no encontrado")
			}
		})
	}
}

// TestMaterialHandler_GenerateUploadURL_ValidFileNames verifica nombres de archivo válidos
func TestMaterialHandler_GenerateUploadURL_ValidFileNames(t *testing.T) {
	testCases := []struct {
		name            string
		fileName        string
		contentType     string
		expectedFileURL string
	}{
		{
			name:            "nombre simple válido",
			fileName:        "document.pdf",
			contentType:     "application/pdf",
			expectedFileURL: "materials/test-id/document.pdf",
		},
		{
			name:            "nombre con guiones",
			fileName:        "my-document-2024.pdf",
			contentType:     "application/pdf",
			expectedFileURL: "materials/test-id/my-document-2024.pdf",
		},
		{
			name:            "nombre con guiones bajos",
			fileName:        "my_document_final.pdf",
			contentType:     "application/pdf",
			expectedFileURL: "materials/test-id/my_document_final.pdf",
		},
		{
			name:            "nombre con espacios",
			fileName:        "my document.pdf",
			contentType:     "application/pdf",
			expectedFileURL: "materials/test-id/my document.pdf",
		},
		{
			name:            "imagen PNG",
			fileName:        "diagram.png",
			contentType:     "image/png",
			expectedFileURL: "materials/test-id/diagram.png",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			expectedURL := "https://s3.amazonaws.com/bucket/" + tc.expectedFileURL + "?presigned-params"

			mockService := &MockMaterialService{
				GetMaterialFunc: func(ctx context.Context, id string) (*dto.MaterialResponse, error) {
					return &dto.MaterialResponse{ID: id, Title: "Test Material"}, nil
				},
			}

			mockS3 := &MockS3Storage{
				GeneratePresignedUploadURLFunc: func(ctx context.Context, key, contentType string, expires time.Duration) (string, error) {
					// Verificar que el S3 key es el esperado
					assert.Equal(t, tc.expectedFileURL, key)
					assert.Equal(t, tc.contentType, contentType)
					return expectedURL, nil
				},
			}

			logger := NewTestLogger()
			handler := NewMaterialHandler(mockService, mockS3, logger)

			router := SetupTestRouter()
			router.POST("/materials/:id/upload-url", handler.GenerateUploadURL)

			reqBody := fmt.Sprintf(`{"file_name":"%s","content_type":"%s"}`, tc.fileName, tc.contentType)

			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/materials/test-id/upload-url", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, http.StatusOK, w.Code)

			var response dto.GenerateUploadURLResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, expectedURL, response.UploadURL)
			assert.Equal(t, tc.expectedFileURL, response.FileURL)
			assert.Equal(t, 900, response.ExpiresIn) // 15 minutos = 900 segundos
		})
	}
}

// TestMaterialHandler_GenerateUploadURL_MaterialNotFound verifica error cuando material no existe
func TestMaterialHandler_GenerateUploadURL_MaterialNotFound(t *testing.T) {
	// Arrange
	mockService := &MockMaterialService{
		GetMaterialFunc: func(ctx context.Context, id string) (*dto.MaterialResponse, error) {
			// Simular error de material no encontrado
			return nil, fmt.Errorf("material not found")
		},
	}

	logger := NewTestLogger()
	mockS3 := &MockS3Storage{}
	handler := NewMaterialHandler(mockService, mockS3, logger)

	router := SetupTestRouter()
	router.POST("/materials/:id/upload-url", handler.GenerateUploadURL)

	reqBody := `{"file_name":"document.pdf","content_type":"application/pdf"}`

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/materials/nonexistent-id/upload-url", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

// TestMaterialHandler_GenerateUploadURL_InvalidRequest verifica validación de request
func TestMaterialHandler_GenerateUploadURL_InvalidRequest(t *testing.T) {
	// Arrange
	mockService := &MockMaterialService{}
	logger := NewTestLogger()
	mockS3 := &MockS3Storage{}
	handler := NewMaterialHandler(mockService, mockS3, logger)

	router := SetupTestRouter()
	router.POST("/materials/:id/upload-url", handler.GenerateUploadURL)

	testCases := []struct {
		name     string
		body     string
		wantCode int
	}{
		{
			name:     "JSON malformado",
			body:     `{"file_name": invalid}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "campos vacíos",
			body:     `{}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "file_name vacío",
			body:     `{"file_name":"","content_type":"application/pdf"}`,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/materials/test-id/upload-url", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.wantCode, w.Code)
		})
	}
}

// TestMaterialHandler_CreateMaterial_Success verifica creación exitosa de material
func TestMaterialHandler_CreateMaterial_Success(t *testing.T) {
	// Arrange
	expectedID := "test-material-123"
	expectedTitle := "Test Material"
	testUserID := "user-123"
	testSchoolID := "550e8400-e29b-41d4-a716-446655440000" // UUID válido

	mockService := &MockMaterialService{
		CreateMaterialFunc: func(ctx context.Context, req dto.CreateMaterialRequest, authorID string, schoolID string) (*dto.MaterialResponse, error) {
			assert.Equal(t, expectedTitle, req.Title)
			assert.Equal(t, testUserID, authorID)
			assert.Equal(t, testSchoolID, schoolID)

			return &dto.MaterialResponse{
				ID:    expectedID,
				Title: req.Title,
			}, nil
		},
	}

	logger := NewTestLogger()
	mockS3 := &MockS3Storage{}
	handler := NewMaterialHandler(mockService, mockS3, logger)

	router := SetupTestRouter()
	router.POST("/materials", MockAuthMiddleware(testUserID, testSchoolID), handler.CreateMaterial)

	reqBody := fmt.Sprintf(`{"title":"%s","description":"Test description"}`, expectedTitle)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/materials", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.MaterialResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedID, response.ID)
	assert.Equal(t, expectedTitle, response.Title)
}

// TestMaterialHandler_CreateMaterial_InvalidRequest verifica validación de request
func TestMaterialHandler_CreateMaterial_InvalidRequest(t *testing.T) {
	// Arrange
	mockService := &MockMaterialService{}
	logger := NewTestLogger()
	mockS3 := &MockS3Storage{}
	handler := NewMaterialHandler(mockService, mockS3, logger)

	router := SetupTestRouter()
	router.POST("/materials", MockAuthMiddleware("user-123", "550e8400-e29b-41d4-a716-446655440000"), handler.CreateMaterial)

	testCases := []struct {
		name string
		body string
	}{
		{
			name: "JSON malformado",
			body: `{"title": invalid}`,
		},
		{
			name: "campos vacíos",
			body: `{}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/materials", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "invalid request body")
		})
	}
}

// TestMaterialHandler_GetMaterial_Success verifica obtención exitosa de material
func TestMaterialHandler_GetMaterial_Success(t *testing.T) {
	// Arrange
	expectedID := "material-123"
	expectedTitle := "Test Material"

	mockService := &MockMaterialService{
		GetMaterialFunc: func(ctx context.Context, id string) (*dto.MaterialResponse, error) {
			assert.Equal(t, expectedID, id)
			return &dto.MaterialResponse{
				ID:    id,
				Title: expectedTitle,
			}, nil
		},
	}

	logger := NewTestLogger()
	mockS3 := &MockS3Storage{}
	handler := NewMaterialHandler(mockService, mockS3, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id", handler.GetMaterial)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/materials/"+expectedID, nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.MaterialResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedID, response.ID)
	assert.Equal(t, expectedTitle, response.Title)
}

// TestMaterialHandler_GenerateDownloadURL_FileNotUploaded verifica error cuando archivo no existe
func TestMaterialHandler_GenerateDownloadURL_FileNotUploaded(t *testing.T) {
	// Arrange
	mockService := &MockMaterialService{
		GetMaterialFunc: func(ctx context.Context, id string) (*dto.MaterialResponse, error) {
			return &dto.MaterialResponse{
				ID:      id,
				Title:   "Material sin archivo",
				FileURL: "", // Sin S3 key (archivo no subido)
			}, nil
		},
	}

	logger := NewTestLogger()
	mockS3 := &MockS3Storage{}
	handler := NewMaterialHandler(mockService, mockS3, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id/download-url", handler.GenerateDownloadURL)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/materials/test-id/download-url", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "FILE_NOT_FOUND")
	assert.Contains(t, w.Body.String(), "material file not uploaded yet")
}

// ============================================
// Tests: ListMaterials
// ============================================

func TestMaterialHandler_ListMaterials_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	mockService := &MockMaterialService{
		ListMaterialsFunc: func(ctx context.Context, filters repository.ListFilters) ([]*dto.MaterialResponse, error) {
			return []*dto.MaterialResponse{
				{
					ID:          "material-1",
					Title:       "Material 1",
					Description: stringPtr("Description 1"),
				},
				{
					ID:          "material-2",
					Title:       "Material 2",
					Description: stringPtr("Description 2"),
				},
			}, nil
		},
	}

	logger := NewTestLogger()
	mockS3 := &MockS3Storage{}
	handler := NewMaterialHandler(mockService, mockS3, logger)

	req, _ := http.NewRequest("GET", "/v1/materials", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler.ListMaterials(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
}

func TestMaterialHandler_ListMaterials_EmptyList(t *testing.T) {
	t.Parallel()

	// Arrange
	mockService := &MockMaterialService{
		ListMaterialsFunc: func(ctx context.Context, filters repository.ListFilters) ([]*dto.MaterialResponse, error) {
			return []*dto.MaterialResponse{}, nil
		},
	}

	logger := NewTestLogger()
	mockS3 := &MockS3Storage{}
	handler := NewMaterialHandler(mockService, mockS3, logger)

	req, _ := http.NewRequest("GET", "/v1/materials", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler.ListMaterials(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 0)
}

func TestMaterialHandler_ListMaterials_DatabaseError(t *testing.T) {
	t.Parallel()

	// Arrange
	mockService := &MockMaterialService{
		ListMaterialsFunc: func(ctx context.Context, filters repository.ListFilters) ([]*dto.MaterialResponse, error) {
			return nil, errors.NewDatabaseError("list materials", assert.AnError)
		},
	}

	logger := NewTestLogger()
	mockS3 := &MockS3Storage{}
	handler := NewMaterialHandler(mockService, mockS3, logger)

	req, _ := http.NewRequest("GET", "/v1/materials", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler.ListMaterials(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "DATABASE_ERROR", errorResponse.Code)
}
