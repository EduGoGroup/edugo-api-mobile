package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// AuthHandler maneja las peticiones de autenticación
type AuthHandler struct {
	authService service.AuthService
	logger      logger.Logger
}

func NewAuthHandler(authService service.AuthService, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password, returns JWT access token and refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials (email and password)"
// @Success 200 {object} dto.LoginResponse "Successfully authenticated, returns tokens and user info"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Warn("login failed", "error", appErr.Message, "email", req.Email)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("login successful", "email", req.Email, "user_id", response.User.ID)
	c.JSON(http.StatusOK, response)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Obtiene un nuevo access token usando un refresh token válido
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh token to exchange for new access token"
// @Success 200 {object} dto.RefreshResponse "New access token generated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 401 {object} ErrorResponse "Invalid or expired refresh token"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	response, err := h.authService.RefreshAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Warn("refresh failed", "error", appErr.Message)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("token refreshed successfully")
	c.JSON(http.StatusOK, response)
}

// Logout godoc
// @Summary User logout
// @Description Revoca el refresh token del usuario (cierra sesión). Requiere autenticación con Bearer token.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh token a revocar"
// @Success 204 "No content - Logout exitoso"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 401 {object} ErrorResponse "User not authenticated or invalid token"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/auth/logout [post]
// @Security BearerAuth
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated", Code: "UNAUTHORIZED"})
		return
	}

	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	err := h.authService.Logout(c.Request.Context(), userID.(string), req.RefreshToken)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Warn("logout failed", "error", appErr.Message)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("logout successful", "user_id", userID)
	c.Status(http.StatusNoContent)
}

// RevokeAll godoc
// @Summary Revoke all sessions
// @Description Revoca todos los refresh tokens del usuario (cierra todas las sesiones activas). Requiere autenticación con Bearer token.
// @Tags auth
// @Produce json
// @Success 204 "No content - Todas las sesiones revocadas exitosamente"
// @Failure 401 {object} ErrorResponse "User not authenticated"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/auth/revoke-all [post]
// @Security BearerAuth
func (h *AuthHandler) RevokeAll(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated", Code: "UNAUTHORIZED"})
		return
	}

	err := h.authService.RevokeAllSessions(c.Request.Context(), userID.(string))
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("all sessions revoked", "user_id", userID)
	c.Status(http.StatusNoContent)
}

// GetCurrentUser godoc
// @Summary Get current authenticated user
// @Description Obtiene la información del usuario autenticado actual desde el JWT token
// @Tags auth
// @Produce json
// @Success 200 {object} dto.UserInfo "User information"
// @Failure 401 {object} ErrorResponse "User not authenticated"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/auth/me [get]
// @Security BearerAuth
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Extraer user_id del contexto (seteado por middleware de auth)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated", Code: "UNAUTHORIZED"})
		return
	}

	userInfo, err := h.authService.GetCurrentUser(c.Request.Context(), userID.(string))
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Warn("get current user failed", "error", appErr.Message)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("current user retrieved", "user_id", userID)
	c.JSON(http.StatusOK, userInfo)
}
