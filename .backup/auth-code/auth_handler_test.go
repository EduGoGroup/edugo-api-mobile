package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// TestNewAuthHandler verifica el constructor del handler
func TestNewAuthHandler(t *testing.T) {
	// Arrange
	mockService := &MockAuthService{}
	logger := NewTestLogger()

	// Act
	handler := NewAuthHandler(mockService, logger)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.authService)
	assert.Equal(t, logger, handler.logger)
}

// TestAuthHandler_Login_Success verifica login exitoso con credenciales válidas
func TestAuthHandler_Login_Success(t *testing.T) {
	// Arrange
	expectedEmail := "test@example.com"
	expectedUserID := "user-123"
	expectedAccessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	expectedRefreshToken := "refresh-token-abc123"

	mockService := &MockAuthService{
		LoginFunc: func(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
			assert.Equal(t, expectedEmail, req.Email)
			assert.Equal(t, "password123", req.Password)

			return &dto.LoginResponse{
				AccessToken:  expectedAccessToken,
				RefreshToken: expectedRefreshToken,
				User: dto.UserInfo{
					ID:        expectedUserID,
					Email:     req.Email,
					FirstName: "Test",
					LastName:  "User",
					FullName:  "Test User",
				},
			}, nil
		},
	}

	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/login", handler.Login)

	reqBody := fmt.Sprintf(`{"email":"%s","password":"password123"}`, expectedEmail)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedAccessToken, response.AccessToken)
	assert.Equal(t, expectedRefreshToken, response.RefreshToken)
	assert.Equal(t, expectedUserID, response.User.ID)
	assert.Equal(t, expectedEmail, response.User.Email)
	assert.Equal(t, "Test User", response.User.FullName)
}

// TestAuthHandler_Login_InvalidCredentials verifica rechazo de credenciales inválidas
func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	mockService := &MockAuthService{
		LoginFunc: func(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
			return nil, errors.NewUnauthorizedError("invalid credentials")
		},
	}

	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/login", handler.Login)

	testCases := []struct {
		name     string
		email    string
		password string
	}{
		{
			name:     "password incorrecto",
			email:    "test@example.com",
			password: "wrong-password",
		},
		{
			name:     "email no existente",
			email:    "nonexistent@example.com",
			password: "password123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			reqBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, tc.email, tc.password)

			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, http.StatusUnauthorized, w.Code)
			assert.Contains(t, w.Body.String(), "invalid credentials")
			assert.Contains(t, w.Body.String(), "UNAUTHORIZED")
		})
	}
}

// TestAuthHandler_Login_InvalidRequest verifica validación de request body
func TestAuthHandler_Login_InvalidRequest(t *testing.T) {
	// Arrange
	mockService := &MockAuthService{}
	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/login", handler.Login)

	testCases := []struct {
		name        string
		body        string
		wantCode    int
		wantMessage string
	}{
		{
			name:        "JSON malformado",
			body:        `{"email": invalid}`,
			wantCode:    http.StatusBadRequest,
			wantMessage: "invalid request body",
		},
		{
			name:        "email vacío",
			body:        `{"email":"","password":"test123"}`,
			wantCode:    http.StatusBadRequest,
			wantMessage: "invalid request body",
		},
		{
			name:        "password vacío",
			body:        `{"email":"test@example.com","password":""}`,
			wantCode:    http.StatusBadRequest,
			wantMessage: "invalid request body",
		},
		{
			name:        "ambos campos vacíos",
			body:        `{}`,
			wantCode:    http.StatusBadRequest,
			wantMessage: "invalid request body",
		},
		{
			name:        "solo email sin password",
			body:        `{"email":"test@example.com"}`,
			wantCode:    http.StatusBadRequest,
			wantMessage: "invalid request body",
		},
		{
			name:        "solo password sin email",
			body:        `{"password":"test123"}`,
			wantCode:    http.StatusBadRequest,
			wantMessage: "invalid request body",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.wantCode, w.Code, "Status code incorrecto para: %s", tc.name)
			assert.Contains(t, w.Body.String(), tc.wantMessage)
			assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
		})
	}
}

// TestAuthHandler_Login_ServiceError verifica manejo de errores internos del servicio
func TestAuthHandler_Login_ServiceError(t *testing.T) {
	// Arrange
	mockService := &MockAuthService{
		LoginFunc: func(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
			return nil, fmt.Errorf("database connection failed")
		},
	}

	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/login", handler.Login)

	reqBody := `{"email":"test@example.com","password":"password123"}`

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "internal server error")
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}

// TestAuthHandler_Refresh_Success verifica refresh exitoso de access token
func TestAuthHandler_Refresh_Success(t *testing.T) {
	// Arrange
	validRefreshToken := "valid-refresh-token-xyz"
	newAccessToken := "new-access-token-abc"

	mockService := &MockAuthService{
		RefreshAccessTokenFunc: func(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
			assert.Equal(t, validRefreshToken, refreshToken)
			return &dto.RefreshResponse{
				AccessToken: newAccessToken,
			}, nil
		},
	}

	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/refresh", handler.Refresh)

	reqBody := fmt.Sprintf(`{"refresh_token":"%s"}`, validRefreshToken)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/refresh", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.RefreshResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, newAccessToken, response.AccessToken)
}

// TestAuthHandler_Refresh_InvalidToken verifica rechazo de refresh token inválido
func TestAuthHandler_Refresh_InvalidToken(t *testing.T) {
	// Arrange
	mockService := &MockAuthService{
		RefreshAccessTokenFunc: func(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
			return nil, errors.NewUnauthorizedError("invalid refresh token")
		},
	}

	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/refresh", handler.Refresh)

	reqBody := `{"refresh_token":"invalid-token"}`

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/refresh", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid refresh token")
}

// TestAuthHandler_Refresh_ExpiredToken verifica rechazo de refresh token expirado
func TestAuthHandler_Refresh_ExpiredToken(t *testing.T) {
	// Arrange
	mockService := &MockAuthService{
		RefreshAccessTokenFunc: func(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
			return nil, errors.NewUnauthorizedError("refresh token expired")
		},
	}

	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/refresh", handler.Refresh)

	reqBody := `{"refresh_token":"expired-token"}`

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/refresh", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "refresh token expired")
}

// TestAuthHandler_Refresh_InvalidRequest verifica validación de request body
func TestAuthHandler_Refresh_InvalidRequest(t *testing.T) {
	// Arrange
	mockService := &MockAuthService{}
	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/refresh", handler.Refresh)

	testCases := []struct {
		name string
		body string
	}{
		{
			name: "JSON malformado",
			body: `{"refresh_token": invalid}`,
		},
		{
			name: "refresh_token vacío",
			body: `{"refresh_token":""}`,
		},
		{
			name: "sin refresh_token",
			body: `{}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/refresh", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "invalid request body")
		})
	}
}

// TestAuthHandler_Logout_Success verifica logout exitoso
func TestAuthHandler_Logout_Success(t *testing.T) {
	// Arrange
	expectedUserID := "user-123"
	expectedRefreshToken := "refresh-token-to-revoke"

	mockService := &MockAuthService{
		LogoutFunc: func(ctx context.Context, userID, refreshToken string) error {
			assert.Equal(t, expectedUserID, userID)
			assert.Equal(t, expectedRefreshToken, refreshToken)
			return nil
		},
	}

	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/logout", MockUserIDMiddleware(expectedUserID), handler.Logout)

	reqBody := fmt.Sprintf(`{"refresh_token":"%s"}`, expectedRefreshToken)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/logout", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

// TestAuthHandler_Logout_Unauthenticated verifica error cuando no hay user_id en contexto
func TestAuthHandler_Logout_Unauthenticated(t *testing.T) {
	// Arrange
	mockService := &MockAuthService{}
	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/logout", handler.Logout) // Sin middleware de autenticación

	reqBody := `{"refresh_token":"some-token"}`

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/logout", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "user not authenticated")
	assert.Contains(t, w.Body.String(), "UNAUTHORIZED")
}

// TestAuthHandler_Logout_InvalidRequest verifica validación de request body
func TestAuthHandler_Logout_InvalidRequest(t *testing.T) {
	// Arrange
	mockService := &MockAuthService{}
	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/logout", MockUserIDMiddleware("user-123"), handler.Logout)

	testCases := []struct {
		name string
		body string
	}{
		{
			name: "JSON malformado",
			body: `{"refresh_token": invalid}`,
		},
		{
			name: "refresh_token vacío",
			body: `{"refresh_token":""}`,
		},
		{
			name: "sin refresh_token",
			body: `{}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/logout", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "invalid request body")
		})
	}
}

// TestAuthHandler_Logout_ServiceError verifica manejo de errores del servicio
func TestAuthHandler_Logout_ServiceError(t *testing.T) {
	// Arrange
	mockService := &MockAuthService{
		LogoutFunc: func(ctx context.Context, userID, refreshToken string) error {
			return errors.NewInternalError("database error", fmt.Errorf("connection failed"))
		},
	}

	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/logout", MockUserIDMiddleware("user-123"), handler.Logout)

	reqBody := `{"refresh_token":"valid-token"}`

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/logout", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestAuthHandler_RevokeAll_Success verifica revocación exitosa de todas las sesiones
func TestAuthHandler_RevokeAll_Success(t *testing.T) {
	// Arrange
	expectedUserID := "user-123"

	mockService := &MockAuthService{
		RevokeAllSessionsFunc: func(ctx context.Context, userID string) error {
			assert.Equal(t, expectedUserID, userID)
			return nil
		},
	}

	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/revoke-all", MockUserIDMiddleware(expectedUserID), handler.RevokeAll)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/revoke-all", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

// TestAuthHandler_RevokeAll_Unauthenticated verifica error cuando no hay user_id
func TestAuthHandler_RevokeAll_Unauthenticated(t *testing.T) {
	// Arrange
	mockService := &MockAuthService{}
	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/revoke-all", handler.RevokeAll) // Sin middleware

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/revoke-all", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "user not authenticated")
}

// TestAuthHandler_RevokeAll_ServiceError verifica manejo de errores del servicio
func TestAuthHandler_RevokeAll_ServiceError(t *testing.T) {
	// Arrange
	mockService := &MockAuthService{
		RevokeAllSessionsFunc: func(ctx context.Context, userID string) error {
			return errors.NewInternalError("database connection failed", fmt.Errorf("connection error"))
		},
	}

	logger := NewTestLogger()
	handler := NewAuthHandler(mockService, logger)

	router := SetupTestRouter()
	router.POST("/auth/revoke-all", MockUserIDMiddleware("user-123"), handler.RevokeAll)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/revoke-all", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "database connection failed")
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}
