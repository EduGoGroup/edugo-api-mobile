//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMaterialFlow_CreateMaterial prueba la creación de un material
func TestMaterialFlow_CreateMaterial(t *testing.T) {
	// Setup
	app := SetupTestApp(t)
	defer app.Cleanup()
	
	// Limpiar base de datos
	CleanDatabase(t, app.DB)
	
	// Seed usuario de prueba (autor del material)
	userID, email := SeedTestUser(t, app.DB)
	t.Logf("✅ Test user created: %s (%s)", email, userID)
	
	// Crear router Gin con el handler de material
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Registrar endpoint de creación de material
	router.POST("/api/v1/materials", func(c *gin.Context) {
		// Simular middleware de autenticación inyectando user_id en contexto
		c.Set("user_id", userID)
		app.Container.Handlers.MaterialHandler.CreateMaterial(c)
	})
	
	// Preparar request de creación de material
	createReq := map[string]interface{}{
		"title":       "Introducción a Go",
		"description": "Material sobre programación en Go para principiantes",
	}
	
	reqBody, err := json.Marshal(createReq)
	require.NoError(t, err)
	
	// Ejecutar request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/materials", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Validar response
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())
	
	assert.Equal(t, http.StatusCreated, w.Code, "Material creation should succeed")
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verificar que tiene los campos esperados
	assert.Contains(t, response, "id", "Response should contain material ID")
	assert.Equal(t, "Introducción a Go", response["title"])
	assert.Equal(t, "Material sobre programación en Go para principiantes", response["description"])
	assert.Equal(t, userID, response["author_id"])
	assert.Contains(t, response, "status")
	assert.Contains(t, response, "processing_status")
	
	t.Logf("✅ Material created: %v", response["id"])
}

// TestMaterialFlow_GetMaterial prueba obtener un material por ID
func TestMaterialFlow_GetMaterial(t *testing.T) {
	// Setup
	app := SetupTestApp(t)
	defer app.Cleanup()
	
	// Limpiar base de datos
	CleanDatabase(t, app.DB)
	
	// Seed usuario y material de prueba
	userID, _ := SeedTestUser(t, app.DB)
	materialID := SeedTestMaterial(t, app.DB, userID)
	t.Logf("✅ Test material created: %s", materialID)
	
	// Verificar que el material existe en la BD
	var count int
	if err := app.DB.QueryRow("SELECT COUNT(*) FROM materials WHERE id = $1", materialID).Scan(&count); err != nil {
		t.Fatalf("Failed to verify material: %v", err)
	}
	t.Logf("✅ Material count in DB: %d", count)
	
	if count == 0 {
		t.Fatal("Material was not seeded properly")
	}
	
	// Ver el UUID real que está en la BD
	var dbID string
	if err := app.DB.QueryRow("SELECT id::text FROM materials WHERE id = $1", materialID).Scan(&dbID); err != nil {
		t.Fatalf("Failed to get material ID from DB: %v", err)
	}
	t.Logf("✅ Material ID from DB: %s", dbID)
	t.Logf("✅ Material ID we're searching: %s", materialID)
	
	// Crear router Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Registrar endpoint de obtener material
	router.GET("/api/v1/materials/:id", app.Container.Handlers.MaterialHandler.GetMaterial)
	
	// Ejecutar request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/materials/"+materialID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Validar response
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())
	
	assert.Equal(t, http.StatusOK, w.Code, "Get material should succeed")
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verificar estructura
	assert.Equal(t, materialID, response["id"])
	assert.Equal(t, "Test Material", response["title"])
	assert.Equal(t, userID, response["author_id"])
	
	t.Logf("✅ Material retrieved successfully")
}

// TestMaterialFlow_GetMaterialNotFound prueba obtener un material inexistente
func TestMaterialFlow_GetMaterialNotFound(t *testing.T) {
	// Setup
	app := SetupTestApp(t)
	defer app.Cleanup()
	
	// Limpiar base de datos
	CleanDatabase(t, app.DB)
	
	// Crear router Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Registrar endpoint
	router.GET("/api/v1/materials/:id", app.Container.Handlers.MaterialHandler.GetMaterial)
	
	// Ejecutar request con UUID inexistente
	fakeID := "00000000-0000-0000-0000-000000000000"
	req := httptest.NewRequest(http.MethodGet, "/api/v1/materials/"+fakeID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Validar response
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())
	
	assert.Equal(t, http.StatusNotFound, w.Code, "Should return 404 for non-existent material")
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Contains(t, response, "error")
	assert.Contains(t, response, "code")
	
	t.Logf("✅ 404 returned correctly")
}

// TestMaterialFlow_ListMaterials prueba listar materiales
func TestMaterialFlow_ListMaterials(t *testing.T) {
	// Setup
	app := SetupTestApp(t)
	defer app.Cleanup()
	
	// Limpiar base de datos
	CleanDatabase(t, app.DB)
	
	// Seed usuario y varios materiales
	userID, _ := SeedTestUser(t, app.DB)
	material1ID := SeedTestMaterial(t, app.DB, userID)
	material2ID := SeedTestMaterialWithTitle(t, app.DB, userID, "Advanced Go Patterns")
	
	t.Logf("✅ Test materials created: %s, %s", material1ID, material2ID)
	
	// Crear router Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Registrar endpoint de listar materiales
	router.GET("/api/v1/materials", app.Container.Handlers.MaterialHandler.ListMaterials)
	
	// Ejecutar request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/materials", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Validar response
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())
	
	assert.Equal(t, http.StatusOK, w.Code, "List materials should succeed")
	
	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verificar que hay al menos 2 materiales
	assert.GreaterOrEqual(t, len(response), 2, "Should return at least 2 materials")
	
	// Verificar estructura del primer material
	assert.Contains(t, response[0], "id")
	assert.Contains(t, response[0], "title")
	assert.Contains(t, response[0], "author_id")
	
	t.Logf("✅ Listed %d materials successfully", len(response))
}
