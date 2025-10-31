package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/edugo/api-mobile/internal/application/service"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/logger"
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
// @Summary Update reading progress
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
	userID, _ := c.Get("user_id")

	var req struct {
		Percentage int `json:"percentage"`
		LastPage   int `json:"last_page"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request", Code: "INVALID_REQUEST"})
		return
	}

	err := h.progressService.UpdateProgress(c.Request.Context(), id, userID.(string), req.Percentage, req.LastPage)
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
