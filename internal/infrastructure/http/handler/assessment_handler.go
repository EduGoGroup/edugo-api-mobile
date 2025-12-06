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
	assessmentAttemptService service.AssessmentAttemptService
	logger                   logger.Logger
}

func NewAssessmentHandler(
	assessmentAttemptService service.AssessmentAttemptService,
	logger logger.Logger,
) *AssessmentHandler {
	return &AssessmentHandler{
		assessmentAttemptService: assessmentAttemptService,
		logger:                   logger,
	}
}

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
