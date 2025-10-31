package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/edugo/api-mobile/internal/application/service"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/logger"
)

type StatsHandler struct {
	statsService service.StatsService
	logger       logger.Logger
}

func NewStatsHandler(statsService service.StatsService, logger logger.Logger) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
		logger:       logger,
	}
}

// GetMaterialStats godoc
// @Summary Get material statistics
// @Tags materials
// @Produce json
// @Param id path string true "Material ID"
// @Success 200 {object} service.MaterialStats
// @Router /materials/{id}/stats [get]
// @Security BearerAuth
func (h *StatsHandler) GetMaterialStats(c *gin.Context) {
	id := c.Param("id")

	stats, err := h.statsService.GetMaterialStats(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal error", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
