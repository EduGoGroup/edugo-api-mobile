package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// TestNewSummaryHandler verifica el constructor del handler
func TestNewSummaryHandler(t *testing.T) {
	// Arrange
	mockService := &MockSummaryService{}
	logger := NewTestLogger()

	// Act
	handler := NewSummaryHandler(mockService, logger)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.summaryService)
	assert.Equal(t, logger, handler.logger)
}

// TestSummaryHandler_GetSummary_Success verifica obtención exitosa de resumen de material
func TestSummaryHandler_GetSummary_Success(t *testing.T) {
	// Arrange
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	matID, _ := valueobject.MaterialIDFromString(materialID)

	expectedSummary := &repository.MaterialSummary{
		MaterialID: matID,
		MainIdeas: []string{
			"El material cubre conceptos fundamentales de programación",
			"Se enfoca en buenas prácticas y patrones de diseño",
			"Incluye ejemplos prácticos y ejercicios",
		},
		KeyConcepts: map[string]string{
			"SOLID":      "Principios de diseño orientado a objetos",
			"DRY":        "Don't Repeat Yourself - No te repitas",
			"Clean Code": "Código limpio y mantenible",
		},
		Sections: []repository.SummarySection{
			{
				Title:   "Introducción",
				Content: "Conceptos básicos de programación orientada a objetos",
				Page:    1,
			},
			{
				Title:   "Principios SOLID",
				Content: "Explicación detallada de cada principio SOLID",
				Page:    10,
			},
			{
				Title:   "Patrones de Diseño",
				Content: "Patrones comunes y sus aplicaciones",
				Page:    25,
			},
		},
		Glossary: map[string]string{
			"Abstracción":   "Ocultar detalles de implementación",
			"Encapsulación": "Agrupar datos y métodos relacionados",
			"Polimorfismo":  "Capacidad de tomar múltiples formas",
		},
		CreatedAt: "2024-01-15T10:30:00Z",
	}

	mockService := &MockSummaryService{
		GetSummaryFunc: func(ctx context.Context, matID string) (*repository.MaterialSummary, error) {
			assert.Equal(t, materialID, matID)
			return expectedSummary, nil
		},
	}

	logger := NewTestLogger()
	handler := NewSummaryHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id/summary", handler.GetSummary)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/materials/%s/summary", materialID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response repository.MaterialSummary
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, len(expectedSummary.MainIdeas), len(response.MainIdeas))
	assert.Equal(t, len(expectedSummary.KeyConcepts), len(response.KeyConcepts))
	assert.Equal(t, len(expectedSummary.Sections), len(response.Sections))
	assert.Equal(t, len(expectedSummary.Glossary), len(response.Glossary))
	assert.Equal(t, expectedSummary.CreatedAt, response.CreatedAt)
}

// TestSummaryHandler_GetSummary_MaterialNotFound verifica manejo de material sin resumen
func TestSummaryHandler_GetSummary_MaterialNotFound(t *testing.T) {
	// Arrange
	materialID := "550e8400-e29b-41d4-a716-446655440000"

	mockService := &MockSummaryService{
		GetSummaryFunc: func(ctx context.Context, matID string) (*repository.MaterialSummary, error) {
			return nil, errors.NewNotFoundError("summary")
		},
	}

	logger := NewTestLogger()
	handler := NewSummaryHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id/summary", handler.GetSummary)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/materials/%s/summary", materialID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "summary")
	assert.Contains(t, w.Body.String(), "NOT_FOUND")
}

// TestSummaryHandler_GetSummary_InvalidMaterialID verifica manejo de UUID inválido
func TestSummaryHandler_GetSummary_InvalidMaterialID(t *testing.T) {
	// Arrange
	invalidID := "not-a-valid-uuid"

	mockService := &MockSummaryService{
		GetSummaryFunc: func(ctx context.Context, matID string) (*repository.MaterialSummary, error) {
			return nil, errors.NewValidationError("invalid material_id")
		},
	}

	logger := NewTestLogger()
	handler := NewSummaryHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id/summary", handler.GetSummary)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/materials/%s/summary", invalidID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid material_id")
	assert.Contains(t, w.Body.String(), "VALIDATION_ERROR")
}

// TestSummaryHandler_GetSummary_ServiceError verifica manejo de errores internos del servicio
func TestSummaryHandler_GetSummary_ServiceError(t *testing.T) {
	// Arrange
	materialID := "550e8400-e29b-41d4-a716-446655440000"

	mockService := &MockSummaryService{
		GetSummaryFunc: func(ctx context.Context, matID string) (*repository.MaterialSummary, error) {
			return nil, fmt.Errorf("database connection failed")
		},
	}

	logger := NewTestLogger()
	handler := NewSummaryHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id/summary", handler.GetSummary)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/materials/%s/summary", materialID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "internal server error")
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}

// TestSummaryHandler_GetSummary_DatabaseError verifica manejo de errores de base de datos
func TestSummaryHandler_GetSummary_DatabaseError(t *testing.T) {
	// Arrange
	materialID := "550e8400-e29b-41d4-a716-446655440000"

	mockService := &MockSummaryService{
		GetSummaryFunc: func(ctx context.Context, matID string) (*repository.MaterialSummary, error) {
			return nil, errors.NewDatabaseError("get summary", fmt.Errorf("connection timeout"))
		},
	}

	logger := NewTestLogger()
	handler := NewSummaryHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id/summary", handler.GetSummary)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/materials/%s/summary", materialID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "DATABASE_ERROR")
}

// TestSummaryHandler_GetSummary_EmptySummary verifica manejo de resumen vacío
func TestSummaryHandler_GetSummary_EmptySummary(t *testing.T) {
	// Arrange
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	matID, _ := valueobject.MaterialIDFromString(materialID)

	emptySummary := &repository.MaterialSummary{
		MaterialID:  matID,
		MainIdeas:   []string{},
		KeyConcepts: map[string]string{},
		Sections:    []repository.SummarySection{},
		Glossary:    map[string]string{},
		CreatedAt:   "2024-01-15T10:30:00Z",
	}

	mockService := &MockSummaryService{
		GetSummaryFunc: func(ctx context.Context, matID string) (*repository.MaterialSummary, error) {
			return emptySummary, nil
		},
	}

	logger := NewTestLogger()
	handler := NewSummaryHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id/summary", handler.GetSummary)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/materials/%s/summary", materialID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response repository.MaterialSummary
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Empty(t, response.MainIdeas)
	assert.Empty(t, response.KeyConcepts)
	assert.Empty(t, response.Sections)
	assert.Empty(t, response.Glossary)
}

// TestSummaryHandler_GetSummary_WithMultipleSections verifica resumen con múltiples secciones
func TestSummaryHandler_GetSummary_WithMultipleSections(t *testing.T) {
	// Arrange
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	matID, _ := valueobject.MaterialIDFromString(materialID)

	summaryWithSections := &repository.MaterialSummary{
		MaterialID: matID,
		MainIdeas: []string{
			"Idea 1",
			"Idea 2",
			"Idea 3",
		},
		KeyConcepts: map[string]string{
			"concepto1": "definición1",
			"concepto2": "definición2",
		},
		Sections: []repository.SummarySection{
			{Title: "Sección 1", Content: "Contenido 1", Page: 1},
			{Title: "Sección 2", Content: "Contenido 2", Page: 5},
			{Title: "Sección 3", Content: "Contenido 3", Page: 10},
			{Title: "Sección 4", Content: "Contenido 4", Page: 15},
			{Title: "Sección 5", Content: "Contenido 5", Page: 20},
		},
		Glossary: map[string]string{
			"término1": "definición1",
			"término2": "definición2",
			"término3": "definición3",
		},
		CreatedAt: "2024-01-15T10:30:00Z",
	}

	mockService := &MockSummaryService{
		GetSummaryFunc: func(ctx context.Context, matID string) (*repository.MaterialSummary, error) {
			return summaryWithSections, nil
		},
	}

	logger := NewTestLogger()
	handler := NewSummaryHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id/summary", handler.GetSummary)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/materials/%s/summary", materialID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response repository.MaterialSummary
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 3, len(response.MainIdeas))
	assert.Equal(t, 2, len(response.KeyConcepts))
	assert.Equal(t, 5, len(response.Sections))
	assert.Equal(t, 3, len(response.Glossary))

	// Verificar que las secciones mantienen el orden
	for i, section := range response.Sections {
		assert.Equal(t, summaryWithSections.Sections[i].Title, section.Title)
		assert.Equal(t, summaryWithSections.Sections[i].Page, section.Page)
	}
}

// TestSummaryHandler_GetSummary_DifferentMaterials verifica múltiples materiales
func TestSummaryHandler_GetSummary_DifferentMaterials(t *testing.T) {
	// Arrange
	testCases := []struct {
		name       string
		materialID string
		summary    *repository.MaterialSummary
	}{
		{
			name:       "material con resumen completo",
			materialID: "550e8400-e29b-41d4-a716-446655440001",
			summary: func() *repository.MaterialSummary {
				matID, _ := valueobject.MaterialIDFromString("550e8400-e29b-41d4-a716-446655440001")
				return &repository.MaterialSummary{
					MaterialID:  matID,
					MainIdeas:   []string{"Idea 1", "Idea 2", "Idea 3"},
					KeyConcepts: map[string]string{"concepto1": "def1", "concepto2": "def2"},
					Sections: []repository.SummarySection{
						{Title: "Intro", Content: "Contenido intro", Page: 1},
						{Title: "Desarrollo", Content: "Contenido desarrollo", Page: 5},
					},
					Glossary:  map[string]string{"término1": "def1"},
					CreatedAt: "2024-01-15T10:30:00Z",
				}
			}(),
		},
		{
			name:       "material con resumen básico",
			materialID: "550e8400-e29b-41d4-a716-446655440002",
			summary: func() *repository.MaterialSummary {
				matID, _ := valueobject.MaterialIDFromString("550e8400-e29b-41d4-a716-446655440002")
				return &repository.MaterialSummary{
					MaterialID:  matID,
					MainIdeas:   []string{"Idea única"},
					KeyConcepts: map[string]string{"concepto": "definición"},
					Sections: []repository.SummarySection{
						{Title: "Única sección", Content: "Contenido", Page: 1},
					},
					Glossary:  map[string]string{},
					CreatedAt: "2024-01-16T10:30:00Z",
				}
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockService := &MockSummaryService{
				GetSummaryFunc: func(ctx context.Context, matID string) (*repository.MaterialSummary, error) {
					assert.Equal(t, tc.materialID, matID)
					return tc.summary, nil
				},
			}

			logger := NewTestLogger()
			handler := NewSummaryHandler(mockService, logger)

			router := SetupTestRouter()
			router.GET("/materials/:id/summary", handler.GetSummary)

			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/materials/%s/summary", tc.materialID), nil)
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, http.StatusOK, w.Code)

			var response repository.MaterialSummary
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, len(tc.summary.MainIdeas), len(response.MainIdeas))
			assert.Equal(t, len(tc.summary.KeyConcepts), len(response.KeyConcepts))
			assert.Equal(t, len(tc.summary.Sections), len(response.Sections))
			assert.Equal(t, tc.summary.CreatedAt, response.CreatedAt)
		})
	}
}

// TestSummaryHandler_GetSummary_WithSpecialCharacters verifica manejo de caracteres especiales
func TestSummaryHandler_GetSummary_WithSpecialCharacters(t *testing.T) {
	// Arrange
	materialID := "550e8400-e29b-41d4-a716-446655440000"
	matID, _ := valueobject.MaterialIDFromString(materialID)

	summaryWithSpecialChars := &repository.MaterialSummary{
		MaterialID: matID,
		MainIdeas: []string{
			"Idea con acentos: áéíóú",
			"Idea con ñ y símbolos: ¿? ¡!",
			"Idea con comillas: \"texto entre comillas\"",
		},
		KeyConcepts: map[string]string{
			"concepto_con_guión": "definición con & símbolos",
			"concepto/slash":     "definición con <> caracteres",
		},
		Sections: []repository.SummarySection{
			{
				Title:   "Sección con símbolos: @#$%",
				Content: "Contenido con múltiples\nlíneas\ny tabulaciones\t",
				Page:    1,
			},
		},
		Glossary: map[string]string{
			"término_especial": "definición con € y otros símbolos",
		},
		CreatedAt: "2024-01-15T10:30:00Z",
	}

	mockService := &MockSummaryService{
		GetSummaryFunc: func(ctx context.Context, matID string) (*repository.MaterialSummary, error) {
			return summaryWithSpecialChars, nil
		},
	}

	logger := NewTestLogger()
	handler := NewSummaryHandler(mockService, logger)

	router := SetupTestRouter()
	router.GET("/materials/:id/summary", handler.GetSummary)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/materials/%s/summary", materialID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response repository.MaterialSummary
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verificar que los caracteres especiales se mantienen correctamente
	assert.Contains(t, response.MainIdeas[0], "áéíóú")
	assert.Contains(t, response.MainIdeas[1], "ñ")
	assert.NotEmpty(t, response.KeyConcepts)
	assert.NotEmpty(t, response.Sections)
}
