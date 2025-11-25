//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/client"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestRemoteAuth_Integration verifica el flujo completo de autenticación remota
// Este test simula api-admin y verifica que api-mobile valida tokens correctamente
func TestRemoteAuth_Integration(t *testing.T) {
	// Crear mock de api-admin
	apiAdminMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/verify" {
			t.Errorf("Path inesperado: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Decodificar request
		var req struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Simular validación de token
		var response client.TokenInfo
		switch req.Token {
		case "valid-teacher-token":
			response = client.TokenInfo{
				Valid:     true,
				UserID:    "teacher-123",
				Email:     "teacher@example.com",
				Role:      "teacher",
				ExpiresAt: time.Now().Add(1 * time.Hour),
			}
		case "valid-student-token":
			response = client.TokenInfo{
				Valid:     true,
				UserID:    "student-456",
				Email:     "student@example.com",
				Role:      "student",
				ExpiresAt: time.Now().Add(1 * time.Hour),
			}
		case "expired-token":
			response = client.TokenInfo{
				Valid: false,
				Error: "token expirado",
			}
		default:
			response = client.TokenInfo{
				Valid: false,
				Error: "token inválido",
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer apiAdminMock.Close()

	// Crear AuthClient apuntando al mock
	authClient := client.NewAuthClient(client.AuthClientConfig{
		BaseURL:      apiAdminMock.URL,
		Timeout:      5 * time.Second,
		CacheEnabled: true,
		CacheTTL:     60 * time.Second,
	})

	// Crear router con RemoteAuthMiddleware
	router := gin.New()
	router.Use(middleware.RemoteAuthMiddleware(middleware.RemoteAuthConfig{
		AuthClient: authClient,
	}))

	// Endpoint protegido de prueba
	router.GET("/v1/materials", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id": middleware.GetUserID(c),
			"email":   middleware.GetUserEmail(c),
			"role":    middleware.GetUserRole(c),
			"message": "materials list",
		})
	})

	// Test 1: Token válido de teacher
	t.Run("ValidTeacherToken", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/materials", nil)
		req.Header.Set("Authorization", "Bearer valid-teacher-token")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Status esperado 200, obtenido %d: %s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)

		if response["user_id"] != "teacher-123" {
			t.Errorf("user_id esperado 'teacher-123', obtenido '%v'", response["user_id"])
		}
		if response["role"] != "teacher" {
			t.Errorf("role esperado 'teacher', obtenido '%v'", response["role"])
		}
	})

	// Test 2: Token válido de student
	t.Run("ValidStudentToken", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/materials", nil)
		req.Header.Set("Authorization", "Bearer valid-student-token")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Status esperado 200, obtenido %d", w.Code)
		}

		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)

		if response["user_id"] != "student-456" {
			t.Errorf("user_id esperado 'student-456', obtenido '%v'", response["user_id"])
		}
		if response["role"] != "student" {
			t.Errorf("role esperado 'student', obtenido '%v'", response["role"])
		}
	})

	// Test 3: Token expirado
	t.Run("ExpiredToken", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/materials", nil)
		req.Header.Set("Authorization", "Bearer expired-token")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Status esperado 401, obtenido %d", w.Code)
		}
	})

	// Test 4: Token inválido
	t.Run("InvalidToken", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/materials", nil)
		req.Header.Set("Authorization", "Bearer invalid-random-token")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Status esperado 401, obtenido %d", w.Code)
		}
	})

	// Test 5: Sin token
	t.Run("NoToken", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/materials", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Status esperado 401, obtenido %d", w.Code)
		}
	})

	// Test 6: Formato de token inválido
	t.Run("InvalidTokenFormat", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/materials", nil)
		req.Header.Set("Authorization", "InvalidFormat token123")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Status esperado 401, obtenido %d", w.Code)
		}
	})
}

// TestRemoteAuth_CacheIntegration verifica que el cache funciona correctamente
func TestRemoteAuth_CacheIntegration(t *testing.T) {
	callCount := 0

	// Mock que cuenta llamadas
	apiAdminMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		response := client.TokenInfo{
			Valid:  true,
			UserID: "user-123",
			Role:   "teacher",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer apiAdminMock.Close()

	authClient := client.NewAuthClient(client.AuthClientConfig{
		BaseURL:      apiAdminMock.URL,
		CacheEnabled: true,
		CacheTTL:     5 * time.Second,
	})

	router := gin.New()
	router.Use(middleware.RemoteAuthMiddleware(middleware.RemoteAuthConfig{
		AuthClient: authClient,
	}))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// Hacer 5 requests con el mismo token
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer same-token")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: status esperado 200, obtenido %d", i+1, w.Code)
		}
	}

	// Con cache habilitado, solo debería haber 1 llamada a api-admin
	// (o pocas debido a race conditions)
	if callCount > 2 {
		t.Errorf("Cache no está funcionando: esperadas <=2 llamadas, obtenidas %d", callCount)
	}
	t.Logf("Llamadas a api-admin con 5 requests: %d (cache habilitado)", callCount)
}

// TestRemoteAuth_RequireRoleIntegration verifica control de acceso por roles
func TestRemoteAuth_RequireRoleIntegration(t *testing.T) {
	apiAdminMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Token string `json:"token"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)

		var response client.TokenInfo
		if req.Token == "admin-token" {
			response = client.TokenInfo{Valid: true, UserID: "admin-1", Role: "admin"}
		} else if req.Token == "teacher-token" {
			response = client.TokenInfo{Valid: true, UserID: "teacher-1", Role: "teacher"}
		} else {
			response = client.TokenInfo{Valid: true, UserID: "student-1", Role: "student"}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer apiAdminMock.Close()

	authClient := client.NewAuthClient(client.AuthClientConfig{
		BaseURL:      apiAdminMock.URL,
		CacheEnabled: false,
	})

	router := gin.New()
	router.Use(middleware.RemoteAuthMiddleware(middleware.RemoteAuthConfig{
		AuthClient: authClient,
	}))

	// Endpoint solo para admin
	router.GET("/admin", middleware.RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"access": "admin only"})
	})

	// Endpoint para admin y teacher
	router.GET("/materials/create", middleware.RequireRole("admin", "teacher"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"access": "admin or teacher"})
	})

	// Test: Admin accede a /admin
	t.Run("AdminAccessAdmin", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer admin-token")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Status esperado 200, obtenido %d", w.Code)
		}
	})

	// Test: Teacher NO accede a /admin
	t.Run("TeacherDeniedAdmin", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer teacher-token")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Status esperado 403, obtenido %d", w.Code)
		}
	})

	// Test: Teacher accede a /materials/create
	t.Run("TeacherAccessMaterials", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/materials/create", nil)
		req.Header.Set("Authorization", "Bearer teacher-token")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Status esperado 200, obtenido %d", w.Code)
		}
	})

	// Test: Student NO accede a /materials/create
	t.Run("StudentDeniedMaterials", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/materials/create", nil)
		req.Header.Set("Authorization", "Bearer student-token")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Status esperado 403, obtenido %d", w.Code)
		}
	})
}
