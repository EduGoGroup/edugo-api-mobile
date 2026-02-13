package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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

func TestAuthRequired_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	jwtManager := auth.NewJWTManager("test-secret-key", "test-issuer")
	mockLogger := new(MockLogger)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	// Generate valid token
	token, err := jwtManager.GenerateToken("user-123", "test@example.com", enum.SystemRoleStudent, 15*time.Minute) //nolint:staticcheck // test legacy con tokens sin ActiveContext
	assert.NoError(t, err)

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
	token, err := jwtManager.GenerateToken("user-123", "test@example.com", enum.SystemRoleStudent, -1*time.Hour) //nolint:staticcheck // test legacy con tokens sin ActiveContext
	assert.NoError(t, err)

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
	token, err := jwtManager1.GenerateToken("user-123", "test@example.com", enum.SystemRoleStudent, 15*time.Minute) //nolint:staticcheck // test legacy con tokens sin ActiveContext
	assert.NoError(t, err)

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
