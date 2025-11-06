package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	ginmiddleware "github.com/EduGoGroup/edugo-shared/middleware/gin"
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
	userID := ginmiddleware.MustGetUserID(c)

	var answers map[string]interface{}
	if err := c.ShouldBindJSON(&answers); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request", Code: "INVALID_REQUEST"})
		return
	}

	attempt, err := h.assessmentService.RecordAttempt(c.Request.Context(), id, userID, answers)
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

// SubmitAssessment godoc
// @Summary Submit assessment with automatic scoring and detailed feedback
// @Description Calcula automáticamente el puntaje de una evaluación y genera feedback detallado por pregunta
// @Tags assessments
// @Accept json
// @Produce json
// @Param id path string true "Assessment ID (Material ID)"
// @Param request body SubmitAssessmentRequest true "User responses"
// @Success 200 {object} repository.AssessmentResult "Resultado con score y feedback detallado"
// @Failure 400 {object} ErrorResponse "Invalid request or assessment_id"
// @Failure 404 {object} ErrorResponse "Assessment not found"
// @Failure 409 {object} ErrorResponse "Assessment already completed by user"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /assessments/{id}/submit [post]
// @Security BearerAuth
func (h *AssessmentHandler) SubmitAssessment(c *gin.Context) {
	assessmentID := c.Param("id")
	userID := ginmiddleware.MustGetUserID(c)

	// Parsear body JSON con respuestas del usuario
	var req SubmitAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	// Validar que el body tenga respuestas
	if req.Responses == nil || len(req.Responses) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "responses map is required",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	// Invocar servicio para calcular score y generar feedback
	result, err := h.assessmentService.CalculateScore(
		c.Request.Context(),
		assessmentID,
		userID,
		req.Responses,
	)

	if err != nil {
		// Manejo de errores específicos según tipo
		if appErr, ok := errors.GetAppError(err); ok {
			// Determinar código HTTP apropiado
			statusCode := appErr.StatusCode

			// Si el error es de duplicado (evaluación ya completada), retornar 409 Conflict
			if appErr.Code == "DATABASE_ERROR" && appErr.Message == "database error during save assessment result" {
				// TODO: Mejorar detección de error de duplicado chequeando el error subyacente
				// Por ahora, asumimos que error de save es duplicado (índice UNIQUE en MongoDB)
				statusCode = http.StatusConflict
				c.JSON(statusCode, ErrorResponse{
					Error: "assessment already completed by this user",
					Code:  "ASSESSMENT_ALREADY_COMPLETED",
				})
				return
			}

			c.JSON(statusCode, ErrorResponse{
				Error: appErr.Message,
				Code:  string(appErr.Code),
			})
			return
		}

		h.logger.Error("failed to calculate score", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
			Code:  "INTERNAL_ERROR",
		})
		return
	}

	// Retornar resultado con score y feedback detallado
	h.logger.Info("assessment submitted successfully",
		"assessment_id", assessmentID,
		"user_id", userID,
		"score", result.Score,
		"correct_answers", result.CorrectAnswers,
		"total_questions", result.TotalQuestions,
	)

	c.JSON(http.StatusOK, result)
}

// SubmitAssessmentRequest representa el body del request para submit
type SubmitAssessmentRequest struct {
	Responses map[string]interface{} `json:"responses" binding:"required"` // question_id -> answer
}
