package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/edugo/api-mobile/internal/application/service"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/logger"
)

// AssessmentHandler maneja peticiones de assessments
type AssessmentHandler struct {
	assessmentService service.AssessmentService
	logger            logger.Logger
}

func NewAssessmentHandler(assessmentService service.AssessmentService, logger logger.Logger) *AssessmentHandler {
	return &AssessmentHandler{
		assessmentService: assessmentService,
		logger:            logger,
	}
}

// GetAssessment godoc
// @Summary Get material assessment/quiz
// @Tags materials
// @Produce json
// @Param id path string true "Material ID"
// @Success 200 {object} repository.MaterialAssessment
// @Router /materials/{id}/assessment [get]
// @Security BearerAuth
func (h *AssessmentHandler) GetAssessment(c *gin.Context) {
	id := c.Param("id")

	assessment, err := h.assessmentService.GetAssessment(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, assessment)
}

// RecordAttempt godoc
// @Summary Record assessment attempt
// @Tags materials
// @Accept json
// @Produce json
// @Param id path string true "Material ID"
// @Param request body map[string]interface{} true "Answers"
// @Success 200 {object} repository.AssessmentAttempt
// @Router /materials/{id}/assessment/attempts [post]
// @Security BearerAuth
func (h *AssessmentHandler) RecordAttempt(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	var answers map[string]interface{}
	if err := c.ShouldBindJSON(&answers); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request", Code: "INVALID_REQUEST"})
		return
	}

	attempt, err := h.assessmentService.RecordAttempt(c.Request.Context(), id, userID.(string), answers)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("attempt recorded", "material_id", id, "score", attempt.Score)
	c.JSON(http.StatusOK, attempt)
}
