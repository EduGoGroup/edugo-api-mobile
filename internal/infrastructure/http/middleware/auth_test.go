package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/EduGoGroup/edugo-shared/logger"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Warn(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Error(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Debug(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Fatal(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Sync() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockLogger) With(keysAndValues ...interface{}) logger.Logger {
	args := m.Called(keysAndValues)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(logger.Logger)
}

// generateTestTokenWithContext genera un token JWT para tests con contexto RBAC
func generateTestTokenWithContext(t *testing.T, manager *auth.JWTManager, userID, email string, role enum.SystemRole, expiresIn time.Duration) string {
	t.Helper()

	activeContext := &auth.UserContext{
		RoleID:      "role-" + string(role),
		RoleName:    string(role),
		SchoolID:    "test-school-123",
		SchoolName:  "Test School",
		Permissions: []string{"read", "write"},
	}

	// Si expiresIn es negativo o menor a 1 minuto, generar token manualmente (para tests de expiración)
	if expiresIn < time.Minute {
		return generateExpiredTokenManually(t, userID, email, activeContext, expiresIn)
	}

	token, _, err := manager.GenerateTokenWithContext(userID, email, activeContext, expiresIn)
	if err != nil {
		t.Fatalf("Error generando token de prueba: %v", err)
	}
	return token
}

// generateExpiredTokenManually crea un token JWT manualmente para tests de expiración
func generateExpiredTokenManually(t *testing.T, userID, email string, activeContext *auth.UserContext, expiresIn time.Duration) string {
	t.Helper()

	now := time.Now()
	expiresAt := now.Add(expiresIn)

	claims := &auth.Claims{
		UserID:        userID,
		Email:         email,
		ActiveContext: activeContext,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "test-issuer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test-secret-key"))
	if err != nil {
		t.Fatalf("Error generando token manual: %v", err)
	}

	return tokenString
}

func TestAuthRequired_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	jwtManager := auth.NewJWTManager("test-secret-key", "test-issuer")
	mockLogger := new(MockLogger)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	// Generate valid token
	token := generateTestTokenWithContext(t, jwtManager, "user-123", "test@example.com", enum.SystemRoleStudent, 15*time.Minute)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c.Request = req

	middleware := AuthRequired(jwtManager, mockLogger)

	// Act
	middleware(c)

	// Assert
	assert.False(t, c.IsAborted())
	assert.Equal(t, "user-123", c.GetString("user_id"))
	assert.Equal(t, "test@example.com", c.GetString("email"))
	// Role is set in context
	_, exists := c.Get("role")
	assert.True(t, exists)
}

func TestAuthRequired_MissingAuthorizationHeader(t *testing.T) {
	t.Parallel()

	// Arrange
	jwtManager := auth.NewJWTManager("test-secret-key", "test-issuer")
	mockLogger := new(MockLogger)
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test", nil)
	// No Authorization header
	c.Request = req

	middleware := AuthRequired(jwtManager, mockLogger)

	// Act
	middleware(c)

	// Assert
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "authorization required")
}

func TestAuthRequired_InvalidHeaderFormat_NoBearer(t *testing.T) {
	t.Parallel()

	// Arrange
	jwtManager := auth.NewJWTManager("test-secret-key", "test-issuer")
	mockLogger := new(MockLogger)
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidToken")
	c.Request = req

	middleware := AuthRequired(jwtManager, mockLogger)

	// Act
	middleware(c)

	// Assert
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid authorization header")
}

func TestAuthRequired_InvalidHeaderFormat_OnlyBearer(t *testing.T) {
	t.Parallel()

	// Arrange
	jwtManager := auth.NewJWTManager("test-secret-key", "test-issuer")
	mockLogger := new(MockLogger)
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer")
	c.Request = req

	middleware := AuthRequired(jwtManager, mockLogger)

	// Act
	middleware(c)

	// Assert
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid authorization header")
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	t.Parallel()

	// Arrange
	jwtManager := auth.NewJWTManager("test-secret-key", "test-issuer")
	mockLogger := new(MockLogger)
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	c.Request = req

	middleware := AuthRequired(jwtManager, mockLogger)

	// Act
	middleware(c)

	// Assert
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid or expired token")
}

func TestAuthRequired_ExpiredToken(t *testing.T) {
	t.Parallel()

	// Arrange
	jwtManager := auth.NewJWTManager("test-secret-key", "test-issuer")
	mockLogger := new(MockLogger)
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()

	// Generate expired token (negative duration)
	token := generateTestTokenWithContext(t, jwtManager, "user-123", "test@example.com", enum.SystemRoleStudent, -1*time.Hour)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c.Request = req

	middleware := AuthRequired(jwtManager, mockLogger)

	// Act
	middleware(c)

	// Assert
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid or expired token")
}

func TestAuthRequired_WrongSecret(t *testing.T) {
	t.Parallel()

	// Arrange
	jwtManager1 := auth.NewJWTManager("secret-1", "test-issuer")
	jwtManager2 := auth.NewJWTManager("secret-2", "test-issuer")
	mockLogger := new(MockLogger)
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()

	// Generate token with one secret
	token := generateTestTokenWithContext(t, jwtManager1, "user-123", "test@example.com", enum.SystemRoleStudent, 15*time.Minute)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	c.Request = req

	// Try to validate with different secret
	middleware := AuthRequired(jwtManager2, mockLogger)

	// Act
	middleware(c)

	// Assert
	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid or expired token")
}
