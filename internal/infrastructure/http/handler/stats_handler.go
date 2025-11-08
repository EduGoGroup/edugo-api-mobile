package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
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
// @Description Retrieves statistics for a specific material including views, completion rate, and average score
// @Tags stats
// @Produce json
// @Param id path string true "Material ID (UUID format)"
// @Success 200 {object} service.MaterialStats "Material statistics retrieved successfully"
// @Failure 400 {object} ErrorResponse "Invalid material ID format"
// @Failure 404 {object} ErrorResponse "Material not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
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

// GetGlobalStats godoc
// @Summary Get global system statistics
// @Description Obtiene estadísticas globales del sistema (solo admins)
// @Tags stats
// @Produce json
// @Success 200 {object} dto.GlobalStatsDTO
// @Failure 403 {object} ErrorResponse "Forbidden - solo admins"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/stats/global [get]
// @Security BearerAuth
func (h *StatsHandler) GetGlobalStats(c *gin.Context) {
	// Obtener estadísticas globales del servicio
	stats, err := h.statsService.GetGlobalStats(c.Request.Context())
	if err != nil {
		h.logger.Error("error al obtener estadísticas globales")
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "error interno del servidor", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
