package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/pkg/errors"
	"github.com/EduGoGroup/edugo-shared/pkg/logger"
)

// SummaryHandler maneja peticiones de summaries
type SummaryHandler struct {
	summaryService service.SummaryService
	logger         logger.Logger
}

func NewSummaryHandler(summaryService service.SummaryService, logger logger.Logger) *SummaryHandler {
	return &SummaryHandler{
		summaryService: summaryService,
		logger:         logger,
	}
}

// GetSummary godoc
// @Summary Get material summary
// @Tags materials
// @Produce json
// @Param id path string true "Material ID"
// @Success 200 {object} repository.MaterialSummary
// @Router /materials/{id}/summary [get]
// @Security BearerAuth
func (h *SummaryHandler) GetSummary(c *gin.Context) {
	id := c.Param("id")

	summary, err := h.summaryService.GetSummary(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
