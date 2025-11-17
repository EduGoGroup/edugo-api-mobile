package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/application/service"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	ginmiddleware "github.com/EduGoGroup/edugo-shared/middleware/gin"
)

// AssessmentHandler maneja peticiones de assessments
type AssessmentHandler struct {
	assessmentService        service.AssessmentService
	assessmentAttemptService service.AssessmentAttemptService
	logger                   logger.Logger
}

func NewAssessmentHandler(
	assessmentService service.AssessmentService,
	assessmentAttemptService service.AssessmentAttemptService,
	logger logger.Logger,
) *AssessmentHandler {
	return &AssessmentHandler{
		assessmentService:        assessmentService,
		assessmentAttemptService: assessmentAttemptService,
		logger:                   logger,
	}
}

// GetAssessment godoc
// @Summary Get material assessment/quiz
// @Description Retrieves the assessment (quiz/test) associated with a specific material
// @Tags assessments
// @Produce json
// @Param id path string true "Material ID (UUID format)"
// @Success 200 {object} map[string]interface{} "Assessment retrieved successfully"
// @Failure 400 {object} ErrorResponse "Invalid material ID format"
// @Failure 404 {object} ErrorResponse "Assessment not found for this material"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/materials/{id}/assessment [get]
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
// @Description Records a user's attempt at completing an assessment with their answers
// @Tags assessments
// @Accept json
// @Produce json
// @Param id path string true "Material ID (UUID format)"
// @Param request body map[string]interface{} true "User answers (question_id -> answer mapping)"
// @Success 200 {object} map[string]interface{} "Attempt recorded successfully with score"
// @Failure 400 {object} ErrorResponse "Invalid request body or material ID"
// @Failure 401 {object} ErrorResponse "User not authenticated"
// @Failure 404 {object} ErrorResponse "Assessment not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/materials/{id}/assessment/attempts [post]
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
// @Success 200 {object} map[string]interface{} "Resultado con score y feedback detallado"
// @Failure 400 {object} ErrorResponse "Invalid request or assessment_id"
// @Failure 404 {object} ErrorResponse "Assessment not found"
// @Failure 409 {object} ErrorResponse "Assessment already completed by user"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/assessments/{id}/submit [post]
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
	if len(req.Responses) == 0 {
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
	Responses map[string]interface{} `json:"responses" binding:"required" swaggertype:"object" example:"{\"q1\":\"answer_a\",\"q2\":\"True\",\"q3\":\"Paris\"}"` // question_id -> answer
}

// ========== SPRINT-04: NUEVOS ENDPOINTS DEL SISTEMA DE EVALUACIONES ==========

// GetMaterialAssessment godoc
// @Summary Obtener cuestionario de un material (SIN respuestas correctas)
// @Description Retorna las preguntas de evaluación de un material sin exponer las respuestas correctas
// @Tags Evaluaciones
// @Security BearerAuth
// @Param id path string true "Material ID (UUID)"
// @Success 200 {object} dto.AssessmentResponse "Assessment obtenido exitosamente"
// @Failure 400 {object} ErrorResponse "Invalid material ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Assessment not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/materials/{id}/assessment [get]
func (h *AssessmentHandler) GetMaterialAssessment(c *gin.Context) {
	materialID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid material ID", Code: "INVALID_MATERIAL_ID"})
		return
	}

	assessment, err := h.assessmentAttemptService.GetAssessmentByMaterialID(c.Request.Context(), materialID)
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

// CreateMaterialAttempt godoc
// @Summary Crear intento de evaluación y obtener calificación
// @Description Crea un intento, valida respuestas en servidor, calcula score y retorna resultados con feedback
// @Tags Evaluaciones
// @Security BearerAuth
// @Param id path string true "Material ID (UUID)"
// @Param request body dto.CreateAttemptRequest true "Respuestas del estudiante"
// @Success 201 {object} dto.AttemptResultResponse "Attempt creado exitosamente"
// @Failure 400 {object} ErrorResponse "Invalid request or material ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Assessment not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/materials/{id}/assessment/attempts [post]
func (h *AssessmentHandler) CreateMaterialAttempt(c *gin.Context) {
	materialID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid material ID", Code: "INVALID_MATERIAL_ID"})
		return
	}

	// Obtener student ID del JWT
	studentIDStr := ginmiddleware.MustGetUserID(c)
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user ID", Code: "INVALID_USER_ID"})
		return
	}

	var req dto.CreateAttemptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	result, err := h.assessmentAttemptService.CreateAttempt(c.Request.Context(), studentID, materialID, req)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("attempt created",
		"material_id", materialID.String(),
		"student_id", studentID.String(),
		"score", result.Score,
		"passed", result.Passed,
	)

	c.JSON(http.StatusCreated, result)
}

// GetAttemptResults godoc
// @Summary Obtener resultados de un intento específico
// @Description Retorna los resultados detallados de un intento, incluyendo feedback por pregunta
// @Tags Evaluaciones
// @Security BearerAuth
// @Param id path string true "Attempt ID (UUID)"
// @Success 200 {object} dto.AttemptResultResponse "Resultados obtenidos exitosamente"
// @Failure 400 {object} ErrorResponse "Invalid attempt ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden - attempt does not belong to user"
// @Failure 404 {object} ErrorResponse "Attempt not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/attempts/{id}/results [get]
func (h *AssessmentHandler) GetAttemptResults(c *gin.Context) {
	attemptID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid attempt ID", Code: "INVALID_ATTEMPT_ID"})
		return
	}

	// Obtener student ID del JWT
	studentIDStr := ginmiddleware.MustGetUserID(c)
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user ID", Code: "INVALID_USER_ID"})
		return
	}

	result, err := h.assessmentAttemptService.GetAttemptResult(c.Request.Context(), attemptID, studentID)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUserAttemptHistory godoc
// @Summary Obtener historial de intentos del usuario autenticado
// @Description Retorna todos los intentos de evaluación del usuario, ordenados por fecha descendente
// @Tags Evaluaciones
// @Security BearerAuth
// @Param limit query int false "Número máximo de resultados" default(10)
// @Param offset query int false "Número de resultados a saltar" default(0)
// @Success 200 {object} dto.AttemptHistoryResponse "Historial obtenido exitosamente"
// @Failure 400 {object} ErrorResponse "Invalid query parameters"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/users/me/attempts [get]
func (h *AssessmentHandler) GetUserAttemptHistory(c *gin.Context) {
	// Obtener student ID del JWT
	studentIDStr := ginmiddleware.MustGetUserID(c)
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user ID", Code: "INVALID_USER_ID"})
		return
	}

	// Parsear query parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	history, err := h.assessmentAttemptService.GetAttemptHistory(c.Request.Context(), studentID, limit, offset)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, history)
}
