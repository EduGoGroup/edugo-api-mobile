package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-mobile/internal/client"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// mockAuthServer crea un servidor HTTP mock para simular api-admin
func mockAuthServer(t *testing.T, response client.TokenInfo) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
}

func TestRemoteAuthMiddleware_ValidToken(t *testing.T) {
	// Crear servidor mock que retorna token válido
	server := mockAuthServer(t, client.TokenInfo{
		Valid:  true,
		UserID: "user-123",
		Email:  "test@example.com",
		Role:   "teacher",
	})
	defer server.Close()

	// Crear cliente de auth
	authClient := client.NewAuthClient(client.AuthClientConfig{
		BaseURL:      server.URL,
		CacheEnabled: false,
	})

	// Crear router con middleware
	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
	}))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id": GetUserID(c),
			"email":   GetUserEmail(c),
			"role":    GetUserRole(c),
		})
	})

	// Hacer request con token válido
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code esperado 200, obtenido %d", w.Code)
	}

	var response map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	if response["user_id"] != "user-123" {
		t.Errorf("user_id esperado 'user-123', obtenido '%s'", response["user_id"])
	}
	if response["email"] != "test@example.com" {
		t.Errorf("email esperado 'test@example.com', obtenido '%s'", response["email"])
	}
	if response["role"] != "teacher" {
		t.Errorf("role esperado 'teacher', obtenido '%s'", response["role"])
	}
}

func TestRemoteAuthMiddleware_InvalidToken(t *testing.T) {
	// Servidor que retorna token inválido
	server := mockAuthServer(t, client.TokenInfo{
		Valid: false,
		Error: "token expirado",
	})
	defer server.Close()

	authClient := client.NewAuthClient(client.AuthClientConfig{
		BaseURL:      server.URL,
		CacheEnabled: false,
	})

	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
	}))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code esperado 401, obtenido %d", w.Code)
	}
}

func TestRemoteAuthMiddleware_MissingToken(t *testing.T) {
	server := mockAuthServer(t, client.TokenInfo{Valid: true})
	defer server.Close()

	authClient := client.NewAuthClient(client.AuthClientConfig{
		BaseURL:      server.URL,
		CacheEnabled: false,
	})

	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
	}))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	// Sin header Authorization
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code esperado 401, obtenido %d", w.Code)
	}

	var response map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"] != "UNAUTHORIZED" {
		t.Errorf("code esperado 'UNAUTHORIZED', obtenido '%s'", response["code"])
	}
}

func TestRemoteAuthMiddleware_InvalidFormat(t *testing.T) {
	server := mockAuthServer(t, client.TokenInfo{Valid: true})
	defer server.Close()

	authClient := client.NewAuthClient(client.AuthClientConfig{
		BaseURL:      server.URL,
		CacheEnabled: false,
	})

	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
	}))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test con formato inválido (sin "Bearer")
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat token123")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code esperado 401, obtenido %d", w.Code)
	}
}

func TestRemoteAuthMiddleware_SkipPaths(t *testing.T) {
	server := mockAuthServer(t, client.TokenInfo{Valid: false})
	defer server.Close()

	authClient := client.NewAuthClient(client.AuthClientConfig{
		BaseURL:      server.URL,
		CacheEnabled: false,
	})

	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
		SkipPaths:  []string{"/health", "/public"},
	}))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "public"})
	})
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "protected"})
	})

	// /health debe pasar sin token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("/health: Status code esperado 200, obtenido %d", w.Code)
	}

	// /public debe pasar sin token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/public", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("/public: Status code esperado 200, obtenido %d", w.Code)
	}

	// /protected debe fallar sin token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("/protected: Status code esperado 401, obtenido %d", w.Code)
	}
}

func TestRemoteAuthMiddleware_CustomUnauthorizedHandler(t *testing.T) {
	server := mockAuthServer(t, client.TokenInfo{Valid: false})
	defer server.Close()

	authClient := client.NewAuthClient(client.AuthClientConfig{
		BaseURL:      server.URL,
		CacheEnabled: false,
	})

	customHandlerCalled := false
	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
		OnUnauthorized: func(c *gin.Context, message string) {
			customHandlerCalled = true
			c.JSON(http.StatusUnauthorized, gin.H{
				"custom_error": message,
			})
			c.Abort()
		},
	}))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	router.ServeHTTP(w, req)

	if !customHandlerCalled {
		t.Error("Custom unauthorized handler no fue llamado")
	}

	var response map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	if _, exists := response["custom_error"]; !exists {
		t.Error("Respuesta no contiene 'custom_error'")
	}
}

func TestRequireRole_Allowed(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Simular usuario autenticado
		c.Set("user_id", "user-123")
		c.Set("role", "admin")
		c.Next()
	})
	router.Use(RequireRole("admin", "superadmin"))
	router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin access"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code esperado 200, obtenido %d", w.Code)
	}
}

func TestRequireRole_Forbidden(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "user-123")
		c.Set("role", "student") // Role no permitido
		c.Next()
	})
	router.Use(RequireRole("admin", "teacher"))
	router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin access"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Status code esperado 403, obtenido %d", w.Code)
	}
}

func TestRequireRole_NoRole(t *testing.T) {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "user-123")
		// Sin role
		c.Next()
	})
	router.Use(RequireRole("admin"))
	router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin access"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code esperado 401, obtenido %d", w.Code)
	}
}

func TestGetUserHelpers(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		// Probar sin valores
		if GetUserID(c) != "" {
			t.Error("GetUserID debería retornar vacío sin valor")
		}
		if GetUserEmail(c) != "" {
			t.Error("GetUserEmail debería retornar vacío sin valor")
		}
		if GetUserRole(c) != "" {
			t.Error("GetUserRole debería retornar vacío sin valor")
		}

		// Setear valores
		c.Set("user_id", "test-user")
		c.Set("email", "test@test.com")
		c.Set("role", "tester")

		// Probar con valores
		if GetUserID(c) != "test-user" {
			t.Errorf("GetUserID esperado 'test-user', obtenido '%s'", GetUserID(c))
		}
		if GetUserEmail(c) != "test@test.com" {
			t.Errorf("GetUserEmail esperado 'test@test.com', obtenido '%s'", GetUserEmail(c))
		}
		if GetUserRole(c) != "tester" {
			t.Errorf("GetUserRole esperado 'tester', obtenido '%s'", GetUserRole(c))
		}

		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
}

func TestRemoteAuthMiddleware_BearerCaseInsensitive(t *testing.T) {
	server := mockAuthServer(t, client.TokenInfo{
		Valid:  true,
		UserID: "user-123",
	})
	defer server.Close()

	authClient := client.NewAuthClient(client.AuthClientConfig{
		BaseURL:      server.URL,
		CacheEnabled: false,
	})

	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
	}))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"user_id": GetUserID(c)})
	})

	// Test con "bearer" en minúsculas
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "bearer valid-token")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("bearer minúsculas: Status code esperado 200, obtenido %d", w.Code)
	}

	// Test con "BEARER" en mayúsculas
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "BEARER valid-token")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("BEARER mayúsculas: Status code esperado 200, obtenido %d", w.Code)
	}
}

// Suprimir variable no usada
var _ = context.Background
var _ = time.Now
