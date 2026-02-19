package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	ginmiddleware "github.com/EduGoGroup/edugo-shared/middleware/gin"
)

// ScreenHandler maneja las peticiones HTTP relacionadas con pantallas
type ScreenHandler struct {
	screenService service.ScreenService
	logger        logger.Logger
}

// NewScreenHandler crea una nueva instancia del handler de pantallas
func NewScreenHandler(screenService service.ScreenService, logger logger.Logger) *ScreenHandler {
	return &ScreenHandler{
		screenService: screenService,
		logger:        logger,
	}
}

// GetScreen godoc
// @Summary Get screen definition
// @Description Retrieves a combined screen definition by screen key
// @Tags screens
// @Produce json
// @Param screenKey path string true "Screen key identifier"
// @Param platform query string false "Platform (mobile, desktop, web)"
// @Success 200 {object} screenconfig.CombinedScreenDTO "Screen definition"
// @Success 304 "Not Modified"
// @Failure 404 {object} ErrorResponse "Screen not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/screens/{screenKey} [get]
// @Security BearerAuth
func (h *ScreenHandler) GetScreen(c *gin.Context) {
	screenKey := c.Param("screenKey")
	platform := c.Query("platform")
	userIDStr := ginmiddleware.MustGetUserID(c)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Warn("invalid user_id format", "user_id", userIDStr, "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user_id", Code: "INVALID_USER_ID"})
		return
	}

	screen, err := h.screenService.GetScreen(c.Request.Context(), screenKey, userID, platform)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		h.logger.Error("unexpected error getting screen", "screen_key", screenKey, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	// Generar ETag basado en el contenido
	screenJSON, _ := json.Marshal(screen)
	etag := fmt.Sprintf(`"%x"`, md5.Sum(screenJSON))

	// Soporte If-None-Match para cache condicional
	ifNoneMatch := c.GetHeader("If-None-Match")
	if ifNoneMatch != "" && ifNoneMatch == etag {
		c.Status(http.StatusNotModified)
		return
	}

	// Headers de cache
	c.Header("ETag", etag)
	c.Header("Last-Modified", screen.UpdatedAt.Format(http.TimeFormat))
	c.Header("Cache-Control", "max-age=3600")

	c.JSON(http.StatusOK, screen)
}

// GetScreensForResource godoc
// @Summary Get screens for a resource
// @Description Retrieves all screen configurations linked to a resource
// @Tags screens
// @Produce json
// @Param resourceKey path string true "Resource key identifier"
// @Success 200 {array} screenconfig.ResourceScreenDTO "List of screens"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/screens/resource/{resourceKey} [get]
// @Security BearerAuth
func (h *ScreenHandler) GetScreensForResource(c *gin.Context) {
	resourceKey := c.Param("resourceKey")

	screens, err := h.screenService.GetScreensForResource(c.Request.Context(), resourceKey)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		h.logger.Error("unexpected error getting screens for resource", "resource_key", resourceKey, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, screens)
}

// GetNavigation godoc
// @Summary Get navigation config
// @Description Retrieves the complete navigation structure for the authenticated user
// @Tags screens
// @Produce json
// @Param platform query string false "Platform (mobile, desktop, web)"
// @Success 200 {object} service.NavigationConfigDTO "Navigation configuration"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/screens/navigation [get]
// @Security BearerAuth
func (h *ScreenHandler) GetNavigation(c *gin.Context) {
	platform := c.DefaultQuery("platform", "mobile")
	userIDStr := ginmiddleware.MustGetUserID(c)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Warn("invalid user_id format", "user_id", userIDStr, "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user_id", Code: "INVALID_USER_ID"})
		return
	}

	// Extraer permisos del JWT context
	var permissions []string
	claims, claimsErr := ginmiddleware.GetClaims(c)
	if claimsErr == nil && claims != nil && claims.ActiveContext != nil {
		permissions = claims.ActiveContext.Permissions
	}

	nav, err := h.screenService.GetNavigationConfig(c.Request.Context(), userID, platform, permissions)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		h.logger.Error("unexpected error getting navigation", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, nav)
}

// SavePreferences godoc
// @Summary Save user preferences for a screen
// @Description Saves user-specific preferences for a screen
// @Tags screens
// @Accept json
// @Produce json
// @Param screenKey path string true "Screen key identifier"
// @Param request body json.RawMessage true "User preferences JSON"
// @Success 204 "Preferences saved"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/screens/{screenKey}/preferences [put]
// @Security BearerAuth
func (h *ScreenHandler) SavePreferences(c *gin.Context) {
	screenKey := c.Param("screenKey")
	userIDStr := ginmiddleware.MustGetUserID(c)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Warn("invalid user_id format", "user_id", userIDStr, "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user_id", Code: "INVALID_USER_ID"})
		return
	}

	var prefs json.RawMessage
	if err := c.ShouldBindJSON(&prefs); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	if err := h.screenService.SaveUserPreferences(c.Request.Context(), screenKey, userID, prefs); err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		h.logger.Error("unexpected error saving preferences", "screen_key", screenKey, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	c.Status(http.StatusNoContent)
}
