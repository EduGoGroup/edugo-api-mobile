package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	ginmiddleware "github.com/EduGoGroup/edugo-shared/middleware/gin"
)

type ProgressHandler struct {
	progressService service.ProgressService
	logger          logger.Logger
}

func NewProgressHandler(progressService service.ProgressService, logger logger.Logger) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
		logger:          logger,
	}
}

// UpdateProgress godoc
// @Summary Update reading progress (legacy endpoint)
// @Tags materials
// @Accept json
// @Produce json
// @Param id path string true "Material ID"
// @Param request body map[string]int true "Progress data"
// @Success 204 "No Content"
// @Router /materials/{id}/progress [patch]
// @Security BearerAuth
func (h *ProgressHandler) UpdateProgress(c *gin.Context) {
	id := c.Param("id")
	userID := ginmiddleware.MustGetUserID(c)

	var req struct {
		Percentage int `json:"percentage"`
		LastPage   int `json:"last_page"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request", Code: "INVALID_REQUEST"})
		return
	}

	err := h.progressService.UpdateProgress(c.Request.Context(), id, userID, req.Percentage, req.LastPage)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal error", Code: "INTERNAL_ERROR"})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpsertProgress godoc
// @Summary Upsert progress idempotently (new UPSERT endpoint)
// @Description Updates user progress in a material using idempotent UPSERT operation. Multiple calls with same data are safe.
// @Tags progress
// @Accept json
// @Produce json
// @Param request body UpsertProgressRequest true "Progress data"
// @Success 200 {object} ProgressResponse "Progress updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request (bad UUID, percentage out of range)"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden (user can only update own progress)"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/progress [put]
// @Security BearerAuth
func (h *ProgressHandler) UpsertProgress(c *gin.Context) {
	// Obtener userID del contexto (usuario autenticado)
	authenticatedUserID := ginmiddleware.MustGetUserID(c)

	// Estructura de request
	var req UpsertProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	// Validar que user_id y material_id estén presentes
	if req.UserID == "" || req.MaterialID == "" {
		h.logger.Warn("missing required fields",
			"user_id", req.UserID,
			"material_id", req.MaterialID,
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "user_id and material_id are required",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	// Autorización: Usuario solo puede actualizar su propio progreso (a menos que sea admin)
	// TODO: Agregar verificación de rol admin cuando exista
	if req.UserID != authenticatedUserID {
		h.logger.Warn("user attempting to update progress of another user",
			"authenticated_user_id", authenticatedUserID,
			"target_user_id", req.UserID,
		)
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error: "you can only update your own progress",
			Code:  "FORBIDDEN",
		})
		return
	}

	// Invocar servicio para actualizar progreso
	err := h.progressService.UpdateProgress(
		c.Request.Context(),
		req.MaterialID,
		req.UserID,
		req.ProgressPercentage,
		req.LastPage,
	)

	if err != nil {
		// Manejar errores de aplicación
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("application error during progress update",
				"code", appErr.Code,
				"message", appErr.Message,
			)
			c.JSON(appErr.StatusCode, ErrorResponse{
				Error: appErr.Message,
				Code:  string(appErr.Code),
			})
			return
		}

		// Error genérico
		h.logger.Error("unexpected error during progress update", "error", err.Error())
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
			Code:  "INTERNAL_ERROR",
		})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, ProgressResponse{
		UserID:             req.UserID,
		MaterialID:         req.MaterialID,
		ProgressPercentage: req.ProgressPercentage,
		LastPage:           req.LastPage,
		Message:            "progress updated successfully",
	})
}

// UpsertProgressRequest representa la solicitud de actualización de progreso
type UpsertProgressRequest struct {
	UserID             string `json:"user_id" binding:"required"`
	MaterialID         string `json:"material_id" binding:"required"`
	ProgressPercentage int    `json:"progress_percentage" binding:"required,min=0,max=100"`
	LastPage           int    `json:"last_page"`
}

// ProgressResponse representa la respuesta de actualización de progreso
type ProgressResponse struct {
	UserID             string `json:"user_id"`
	MaterialID         string `json:"material_id"`
	ProgressPercentage int    `json:"progress_percentage"`
	LastPage           int    `json:"last_page"`
	Message            string `json:"message"`
}
