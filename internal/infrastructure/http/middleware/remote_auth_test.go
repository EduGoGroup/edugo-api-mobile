package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-mobile/internal/client"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

const (
	testJWTSecret = "test-secret-key-at-least-32-chars-long-for-security"
	testJWTIssuer = "edugo-central"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// generateTestToken genera un token JWT válido para pruebas
func generateTestToken(t *testing.T, userID, email string, role enum.SystemRole, expiresIn time.Duration) string {
	t.Helper()
	manager := auth.NewJWTManager(testJWTSecret, testJWTIssuer)
	token, err := manager.GenerateToken(userID, email, role, expiresIn)
	if err != nil {
		t.Fatalf("Error generando token de prueba: %v", err)
	}
	return token
}

// createAuthClient crea un cliente de auth configurado para validación local
func createAuthClient() *client.AuthClient {
	return client.NewAuthClient(client.AuthClientConfig{
		JWTSecret:    testJWTSecret,
		JWTIssuer:    testJWTIssuer,
		CacheEnabled: false,
	})
}

func TestRemoteAuthMiddleware_ValidToken(t *testing.T) {
	authClient := createAuthClient()

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

	// Generar token válido
	token := generateTestToken(t, "user-123", "test@example.com", enum.SystemRoleTeacher, 15*time.Minute)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code esperado 200, obtenido %d. Body: %s", w.Code, w.Body.String())
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
	authClient := createAuthClient()

	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
	}))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-that-is-not-a-jwt")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code esperado 401, obtenido %d", w.Code)
	}
}

func TestRemoteAuthMiddleware_ExpiredToken(t *testing.T) {
	authClient := createAuthClient()

	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
	}))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Generar token expirado
	token := generateTestToken(t, "user-123", "test@example.com", enum.SystemRoleTeacher, -1*time.Hour)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code esperado 401, obtenido %d", w.Code)
	}
}

func TestRemoteAuthMiddleware_MissingToken(t *testing.T) {
	authClient := createAuthClient()

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
	authClient := createAuthClient()

	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
	}))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Token sin "Bearer" prefix
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "invalid-format-token")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code esperado 401, obtenido %d", w.Code)
	}
}

func TestRemoteAuthMiddleware_SkipPaths(t *testing.T) {
	authClient := createAuthClient()

	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
		SkipPaths:  []string{"/health", "/public"},
	}))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "public endpoint"})
	})
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "protected"})
	})

	// /health sin token - debe pasar
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Errorf("/health sin token: esperado 200, obtenido %d", w1.Code)
	}

	// /public sin token - debe pasar
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/public", nil)
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("/public sin token: esperado 200, obtenido %d", w2.Code)
	}

	// /protected sin token - debe fallar
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w3, req3)
	if w3.Code != http.StatusUnauthorized {
		t.Errorf("/protected sin token: esperado 401, obtenido %d", w3.Code)
	}
}

func TestRemoteAuthMiddleware_CustomUnauthorizedHandler(t *testing.T) {
	authClient := createAuthClient()

	customHandlerCalled := false

	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
		OnUnauthorized: func(c *gin.Context, message string) {
			customHandlerCalled = true
			c.JSON(http.StatusUnauthorized, gin.H{
				"custom_error": true,
				"message":      message,
			})
			c.Abort()
		},
	}))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w, req)

	if !customHandlerCalled {
		t.Error("Custom handler debería haber sido llamado")
	}

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	if response["custom_error"] != true {
		t.Error("Respuesta debería contener custom_error: true")
	}
}

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Con user_id en contexto
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", "test-user-123")

	if GetUserID(c) != "test-user-123" {
		t.Errorf("GetUserID esperado 'test-user-123', obtenido '%s'", GetUserID(c))
	}

	// Sin user_id en contexto
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)

	if GetUserID(c2) != "" {
		t.Errorf("GetUserID sin contexto esperado '', obtenido '%s'", GetUserID(c2))
	}
}

func TestGetUserEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("email", "test@test.com")

	if GetUserEmail(c) != "test@test.com" {
		t.Errorf("GetUserEmail esperado 'test@test.com', obtenido '%s'", GetUserEmail(c))
	}

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)

	if GetUserEmail(c2) != "" {
		t.Errorf("GetUserEmail sin contexto esperado '', obtenido '%s'", GetUserEmail(c2))
	}
}

func TestGetUserRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "teacher")

	if GetUserRole(c) != "teacher" {
		t.Errorf("GetUserRole esperado 'teacher', obtenido '%s'", GetUserRole(c))
	}

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)

	if GetUserRole(c2) != "" {
		t.Errorf("GetUserRole sin contexto esperado '', obtenido '%s'", GetUserRole(c2))
	}
}

func TestRequireRole_Allowed(t *testing.T) {
	authClient := createAuthClient()

	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
	}))
	router.GET("/admin-only", RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
	})
	router.GET("/teacher-or-admin", RequireRole("admin", "teacher"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "access granted"})
	})

	// Admin accede a admin-only
	tokenAdmin := generateTestToken(t, "admin-1", "admin@test.com", enum.SystemRoleAdmin, 15*time.Minute)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/admin-only", nil)
	req1.Header.Set("Authorization", "Bearer "+tokenAdmin)
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Admin accediendo a admin-only: esperado 200, obtenido %d", w1.Code)
	}

	// Teacher accede a teacher-or-admin
	tokenTeacher := generateTestToken(t, "teacher-1", "teacher@test.com", enum.SystemRoleTeacher, 15*time.Minute)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/teacher-or-admin", nil)
	req2.Header.Set("Authorization", "Bearer "+tokenTeacher)
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Teacher accediendo a teacher-or-admin: esperado 200, obtenido %d", w2.Code)
	}
}

func TestRequireRole_Forbidden(t *testing.T) {
	authClient := createAuthClient()

	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
	}))
	router.GET("/admin-only", RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
	})

	// Student intenta acceder a admin-only
	tokenStudent := generateTestToken(t, "student-1", "student@test.com", enum.SystemRoleStudent, 15*time.Minute)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin-only", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudent)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Student accediendo a admin-only: esperado 403, obtenido %d", w.Code)
	}

	var response map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"] != "FORBIDDEN" {
		t.Errorf("code esperado 'FORBIDDEN', obtenido '%s'", response["code"])
	}
}

func TestRemoteAuthMiddleware_BearerCaseInsensitive(t *testing.T) {
	authClient := createAuthClient()

	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
	}))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	token := generateTestToken(t, "user-123", "test@example.com", enum.SystemRoleTeacher, 15*time.Minute)

	// bearer minúsculas
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/protected", nil)
	req1.Header.Set("Authorization", "bearer "+token)
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("bearer minúsculas: Status code esperado 200, obtenido %d", w1.Code)
	}

	// BEARER mayúsculas
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/protected", nil)
	req2.Header.Set("Authorization", "BEARER "+token)
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("BEARER mayúsculas: Status code esperado 200, obtenido %d", w2.Code)
	}
}

func TestRemoteAuthMiddleware_AllRoles(t *testing.T) {
	authClient := createAuthClient()

	router := gin.New()
	router.Use(RemoteAuthMiddleware(RemoteAuthConfig{
		AuthClient: authClient,
	}))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"role": GetUserRole(c),
		})
	})

	roles := []enum.SystemRole{
		enum.SystemRoleAdmin,
		enum.SystemRoleTeacher,
		enum.SystemRoleStudent,
		enum.SystemRoleGuardian,
	}

	for _, role := range roles {
		t.Run(string(role), func(t *testing.T) {
			token := generateTestToken(t, "user-"+string(role), "test@example.com", role, 15*time.Minute)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Role %s: esperado 200, obtenido %d", role, w.Code)
			}

			var response map[string]string
			_ = json.Unmarshal(w.Body.Bytes(), &response)

			if response["role"] != string(role) {
				t.Errorf("Role esperado '%s', obtenido '%s'", role, response["role"])
			}
		})
	}
}
