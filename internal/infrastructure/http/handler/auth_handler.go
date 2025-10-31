package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/edugo/api-mobile/internal/application/dto"
	"github.com/edugo/api-mobile/internal/application/service"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/logger"
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
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
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
