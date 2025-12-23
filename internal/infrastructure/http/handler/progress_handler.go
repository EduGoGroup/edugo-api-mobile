package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/service"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/middleware"
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
// @Router /v1/progress [put]
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

	// Autorización: Usuario solo puede actualizar su propio progreso
	// Excepción: admin y super_admin pueden actualizar el progreso de cualquier usuario
	if req.UserID != authenticatedUserID {
		if !middleware.IsAdminRole(c) {
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
		// Admin actualizando progreso de otro usuario
		h.logger.Info("admin updating progress of another user",
			"admin_user_id", authenticatedUserID,
			"target_user_id", req.UserID,
		)
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
	UserID             string `json:"user_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	MaterialID         string `json:"material_id" binding:"required" example:"660e8400-e29b-41d4-a716-446655440001"`
	ProgressPercentage int    `json:"progress_percentage" binding:"required,min=0,max=100" example:"75"`
	LastPage           int    `json:"last_page" example:"45"`
}

// ProgressResponse representa la respuesta de actualización de progreso
type ProgressResponse struct {
	UserID             string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	MaterialID         string `json:"material_id" example:"660e8400-e29b-41d4-a716-446655440001"`
	ProgressPercentage int    `json:"progress_percentage" example:"75"`
	LastPage           int    `json:"last_page" example:"45"`
	Message            string `json:"message" example:"progress updated successfully"`
}
