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

// TestProgressFlow_UpsertProgress prueba crear progreso inicial
func TestProgressFlow_UpsertProgress(t *testing.T) {
	// Setup
	app := SetupTestApp(t)
	defer app.Cleanup()
	
	// Limpiar bases de datos
	CleanDatabase(t, app.DB)
	CleanMongoCollections(t, app.MongoDB)
	
	// Seed usuario y material
	userID, _ := SeedTestUser(t, app.DB)
	materialID := SeedTestMaterial(t, app.DB, userID)
	
	// Crear router Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Registrar endpoint con inyección de user_id en contexto
	router.PUT("/api/v1/progress", func(c *gin.Context) {
		c.Set("user_id", userID)
		app.Container.Handlers.ProgressHandler.UpsertProgress(c)
	})
	
	// Preparar request de progreso
	progressReq := map[string]interface{}{
		"user_id":             userID,
		"material_id":         materialID,
		"progress_percentage": 50,
		"last_page":           25,
	}
	
	reqBody, err := json.Marshal(progressReq)
	require.NoError(t, err)
	
	// Ejecutar request
	req := httptest.NewRequest(http.MethodPut, "/api/v1/progress", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Validar response
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())
	
	assert.Equal(t, http.StatusOK, w.Code, "UpsertProgress should succeed")
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verificar estructura
	assert.Equal(t, userID, response["user_id"])
	assert.Equal(t, materialID, response["material_id"])
	assert.Equal(t, float64(50), response["progress_percentage"])
	assert.Equal(t, float64(25), response["last_page"])
	assert.Contains(t, response, "message")
	
	t.Logf("✅ Progress created: %d%%, page %v", int(response["progress_percentage"].(float64)), response["last_page"])
}

// TestProgressFlow_UpsertProgressUpdate prueba actualizar progreso existente (idempotencia)
func TestProgressFlow_UpsertProgressUpdate(t *testing.T) {
	// Setup
	app := SetupTestApp(t)
	defer app.Cleanup()
	
	// Limpiar bases de datos
	CleanDatabase(t, app.DB)
	CleanMongoCollections(t, app.MongoDB)
	
	// Seed usuario y material
	userID, _ := SeedTestUser(t, app.DB)
	materialID := SeedTestMaterial(t, app.DB, userID)
	
	// Crear router Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.PUT("/api/v1/progress", func(c *gin.Context) {
		c.Set("user_id", userID)
		app.Container.Handlers.ProgressHandler.UpsertProgress(c)
	})
	
	// PRIMER upsert: 30%
	progressReq1 := map[string]interface{}{
		"user_id":             userID,
		"material_id":         materialID,
		"progress_percentage": 30,
		"last_page":           15,
	}
	reqBody1, _ := json.Marshal(progressReq1)
	
	req1 := httptest.NewRequest(http.MethodPut, "/api/v1/progress", bytes.NewBuffer(reqBody1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	
	t.Logf("First upsert - Status: %d", w1.Code)
	assert.Equal(t, http.StatusOK, w1.Code, "First upsert should succeed")
	
	// SEGUNDO upsert: 75% (actualizar progreso)
	progressReq2 := map[string]interface{}{
		"user_id":             userID,
		"material_id":         materialID,
		"progress_percentage": 75,
		"last_page":           38,
	}
	reqBody2, _ := json.Marshal(progressReq2)
	
	req2 := httptest.NewRequest(http.MethodPut, "/api/v1/progress", bytes.NewBuffer(reqBody2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	t.Logf("Second upsert - Status: %d", w2.Code)
	t.Logf("Second upsert - Body: %s", w2.Body.String())
	
	assert.Equal(t, http.StatusOK, w2.Code, "Second upsert should succeed (idempotent)")
	
	var response map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &response)
	
	// Verificar que el progreso se actualizó
	assert.Equal(t, float64(75), response["progress_percentage"], "Progress should be updated to 75%")
	assert.Equal(t, float64(38), response["last_page"], "Last page should be updated to 38")
	
	t.Logf("✅ Progress updated successfully: 30%% -> 75%%")
}

// TestProgressFlow_UpsertProgressUnauthorized prueba actualizar progreso de otro usuario (403)
func TestProgressFlow_UpsertProgressUnauthorized(t *testing.T) {
	// Setup
	app := SetupTestApp(t)
	defer app.Cleanup()
	
	// Limpiar bases de datos
	CleanDatabase(t, app.DB)
	CleanMongoCollections(t, app.MongoDB)
	
	// Seed dos usuarios y un material
	userID1, _ := SeedTestUser(t, app.DB)
	userID2, _ := SeedTestUserWithEmail(t, app.DB, "other@edugo.com")
	materialID := SeedTestMaterial(t, app.DB, userID1)
	
	// Crear router Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Usuario 1 intenta actualizar progreso de usuario 2
	router.PUT("/api/v1/progress", func(c *gin.Context) {
		c.Set("user_id", userID1) // Autenticado como user1
		app.Container.Handlers.ProgressHandler.UpsertProgress(c)
	})
	
	// Intentar actualizar progreso de otro usuario
	progressReq := map[string]interface{}{
		"user_id":             userID2, // Pero intenta actualizar progreso de user2
		"material_id":         materialID,
		"progress_percentage": 50,
		"last_page":           25,
	}
	
	reqBody, _ := json.Marshal(progressReq)
	
	req := httptest.NewRequest(http.MethodPut, "/api/v1/progress", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Validar response
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())
	
	assert.Equal(t, http.StatusForbidden, w.Code, "Should return 403 Forbidden")
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "error")
	assert.Contains(t, response["error"], "own progress")
	
	t.Logf("✅ Unauthorized update correctly rejected")
}

// TestProgressFlow_UpsertProgressInvalidData prueba datos inválidos (400)
func TestProgressFlow_UpsertProgressInvalidData(t *testing.T) {
	// Setup
	app := SetupTestApp(t)
	defer app.Cleanup()
	
	// Limpiar bases de datos
	CleanDatabase(t, app.DB)
	CleanMongoCollections(t, app.MongoDB)
	
	// Seed usuario
	userID, _ := SeedTestUser(t, app.DB)
	
	// Crear router Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.PUT("/api/v1/progress", func(c *gin.Context) {
		c.Set("user_id", userID)
		app.Container.Handlers.ProgressHandler.UpsertProgress(c)
	})
	
	// Test 1: Porcentaje inválido (> 100)
	progressReq1 := map[string]interface{}{
		"user_id":             userID,
		"material_id":         "some-material-id",
		"progress_percentage": 150, // Inválido
		"last_page":           25,
	}
	reqBody1, _ := json.Marshal(progressReq1)
	
	req1 := httptest.NewRequest(http.MethodPut, "/api/v1/progress", bytes.NewBuffer(reqBody1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	
	t.Logf("Invalid percentage - Status: %d", w1.Code)
	assert.Equal(t, http.StatusBadRequest, w1.Code, "Should return 400 for invalid percentage")
	
	// Test 2: user_id faltante
	progressReq2 := map[string]interface{}{
		"material_id":         "some-material-id",
		"progress_percentage": 50,
		"last_page":           25,
	}
	reqBody2, _ := json.Marshal(progressReq2)
	
	req2 := httptest.NewRequest(http.MethodPut, "/api/v1/progress", bytes.NewBuffer(reqBody2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	t.Logf("Missing user_id - Status: %d", w2.Code)
	assert.Equal(t, http.StatusBadRequest, w2.Code, "Should return 400 for missing user_id")
	
	t.Logf("✅ Invalid data correctly rejected")
}

// TestStatsFlow_GetMaterialStats prueba obtener estadísticas de un material
func TestStatsFlow_GetMaterialStats(t *testing.T) {
	// Setup
	app := SetupTestApp(t)
	defer app.Cleanup()
	
	// Limpiar bases de datos
	CleanDatabase(t, app.DB)
	CleanMongoCollections(t, app.MongoDB)
	
	// Seed usuario y material
	userID, _ := SeedTestUser(t, app.DB)
	materialID := SeedTestMaterial(t, app.DB, userID)
	
	// Crear router Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Registrar endpoint
	router.GET("/api/v1/materials/:id/stats", app.Container.Handlers.StatsHandler.GetMaterialStats)
	
	// Ejecutar request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/materials/"+materialID+"/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Validar response
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())
	
	assert.Equal(t, http.StatusOK, w.Code, "GetMaterialStats should succeed")
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verificar estructura básica de stats
	// Las stats pueden variar según la implementación, validamos que sea un objeto válido
	assert.NotNil(t, response, "Stats should not be nil")
	
	t.Logf("✅ Material stats retrieved successfully")
}

// TestStatsFlow_GetGlobalStats prueba obtener estadísticas globales
func TestStatsFlow_GetGlobalStats(t *testing.T) {
	// Setup
	app := SetupTestApp(t)
	defer app.Cleanup()
	
	// Limpiar bases de datos
	CleanDatabase(t, app.DB)
	CleanMongoCollections(t, app.MongoDB)
	
	// Seed algunos datos para que las stats tengan contenido
	userID, _ := SeedTestUser(t, app.DB)
	SeedTestMaterial(t, app.DB, userID)
	SeedTestMaterial(t, app.DB, userID)
	
	// Crear router Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Registrar endpoint
	router.GET("/api/v1/stats/global", app.Container.Handlers.StatsHandler.GetGlobalStats)
	
	// Ejecutar request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats/global", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Validar response
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())
	
	assert.Equal(t, http.StatusOK, w.Code, "GetGlobalStats should succeed")
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	// Verificar estructura básica de global stats
	assert.NotNil(t, response, "Global stats should not be nil")
	
	t.Logf("✅ Global stats retrieved successfully")
}
