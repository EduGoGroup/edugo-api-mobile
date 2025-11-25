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

// TestAuthFlow_LoginSuccess prueba el flujo completo de login exitoso
func TestAuthFlow_LoginSuccess(t *testing.T) {
	// Setup
	app := SetupTestAppWithSharedContainers(t)
	defer app.Cleanup()

	// Limpiar base de datos
	CleanDatabase(t, app.DB)

	// Seed usuario de prueba
	userID, email := SeedTestUser(t, app.DB)
	t.Logf("✅ Test user created: %s (%s)", email, userID)

	// Verificar que el usuario existe en la BD
	var count int
	err := app.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to verify user: %v", err)
	}
	t.Logf("✅ User count in DB: %d", count)

	if count == 0 {
		t.Fatal("User was not seeded properly")
	}

	// Crear router Gin con el handler de auth
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Registrar endpoint de login
	router.POST("/api/v1/auth/login", app.Container.Handlers.AuthHandler.Login)

	// Preparar request de login
	loginReq := map[string]string{
		"email":    email,
		"password": "Test1234!", // Password del SeedTestUser helper
	}
	reqBody, err := json.Marshal(loginReq)
	require.NoError(t, err)

	// Ejecutar request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validar response
	t.Logf("Response status: %d", w.Code)
	t.Logf("Response body: %s", w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code, "Login should succeed")

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verificar que hay tokens
	assert.Contains(t, response, "access_token", "Response should contain access_token")
	assert.Contains(t, response, "refresh_token", "Response should contain refresh_token")

	accessToken := response["access_token"].(string)
	refreshToken := response["refresh_token"].(string)

	assert.NotEmpty(t, accessToken, "Access token should not be empty")
	assert.NotEmpty(t, refreshToken, "Refresh token should not be empty")

	t.Logf("✅ Login successful - Access Token: %s...", accessToken[:20])
	t.Logf("✅ Refresh Token: %s...", refreshToken[:20])
}

// TestAuthFlow_LoginInvalidCredentials prueba login con credenciales inválidas
func TestAuthFlow_LoginInvalidCredentials(t *testing.T) {
	// Setup
	app := SetupTestAppWithSharedContainers(t)
	defer app.Cleanup()

	// Limpiar base de datos
	CleanDatabase(t, app.DB)

	// Seed usuario de prueba
	_, email := SeedTestUser(t, app.DB)

	// Crear router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/auth/login", app.Container.Handlers.AuthHandler.Login)

	// Request con password incorrecta
	loginReq := map[string]string{
		"email":    email,
		"password": "wrong-password",
	}
	reqBody, _ := json.Marshal(loginReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validar que falla
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Login with wrong password should fail")

	t.Log("✅ Invalid credentials correctly rejected")
}

// TestAuthFlow_LoginNonexistentUser prueba login con usuario que no existe
func TestAuthFlow_LoginNonexistentUser(t *testing.T) {
	// Setup
	app := SetupTestAppWithSharedContainers(t)
	defer app.Cleanup()

	// Limpiar base de datos (sin seed)
	CleanDatabase(t, app.DB)

	// Crear router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/auth/login", app.Container.Handlers.AuthHandler.Login)

	// Request con usuario inexistente
	loginReq := map[string]string{
		"email":    "nonexistent@edugo.com",
		"password": "any-password",
	}
	reqBody, _ := json.Marshal(loginReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validar que falla
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Login with nonexistent user should fail")

	t.Log("✅ Nonexistent user correctly rejected")
}

// TODO: Agregar más tests de Auth Flow:
// - TestAuthFlow_RefreshToken
// - TestAuthFlow_Logout
// - TestAuthFlow_RateLimiting
// - TestAuthFlow_AccessProtectedResource
