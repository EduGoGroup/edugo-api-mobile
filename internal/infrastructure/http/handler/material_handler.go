package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/edugo/api-mobile/internal/application/dto"
	"github.com/edugo/api-mobile/internal/application/service"
	"github.com/edugo/api-mobile/internal/domain/repository"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/logger"
)

// MaterialHandler maneja las peticiones HTTP relacionadas con materiales
type MaterialHandler struct {
	materialService service.MaterialService
	logger          logger.Logger
}

func NewMaterialHandler(materialService service.MaterialService, logger logger.Logger) *MaterialHandler {
	return &MaterialHandler{
		materialService: materialService,
		logger:          logger,
	}
}

// CreateMaterial godoc
// @Summary Create a new material
// @Description Creates a new educational material
// @Tags materials
// @Accept json
// @Produce json
// @Param request body dto.CreateMaterialRequest true "Material data"
// @Success 201 {object} dto.MaterialResponse
// @Failure 400 {object} ErrorResponse
// @Router /materials [post]
// @Security BearerAuth
func (h *MaterialHandler) CreateMaterial(c *gin.Context) {
	var req dto.CreateMaterialRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	// Obtener user_id del contexto (middleware de autenticación)
	authorID, _ := c.Get("user_id")

	material, err := h.materialService.CreateMaterial(c.Request.Context(), req, authorID.(string))
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Error("create material failed", "error", appErr.Message)
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}

		h.logger.Error("unexpected error", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("material created", "material_id", material.ID)
	c.JSON(http.StatusCreated, material)
}

// GetMaterial godoc
// @Summary Get material by ID
// @Tags materials
// @Produce json
// @Param id path string true "Material ID"
// @Success 200 {object} dto.MaterialResponse
// @Failure 404 {object} ErrorResponse
// @Router /materials/{id} [get]
// @Security BearerAuth
func (h *MaterialHandler) GetMaterial(c *gin.Context) {
	id := c.Param("id")

	material, err := h.materialService.GetMaterial(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, material)
}

// NotifyUploadComplete godoc
// @Summary Notify upload complete
// @Description Notify that PDF upload to S3 is complete
// @Tags materials
// @Accept json
// @Produce json
// @Param id path string true "Material ID"
// @Param request body dto.UploadCompleteRequest true "S3 info"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Router /materials/{id}/upload-complete [post]
// @Security BearerAuth
func (h *MaterialHandler) NotifyUploadComplete(c *gin.Context) {
	id := c.Param("id")
	var req dto.UploadCompleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	err := h.materialService.NotifyUploadComplete(c.Request.Context(), id, req)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	h.logger.Info("upload complete", "material_id", id)
	c.Status(http.StatusNoContent)
}

// ListMaterials godoc
// @Summary List materials
// @Tags materials
// @Produce json
// @Success 200 {array} dto.MaterialResponse
// @Router /materials [get]
// @Security BearerAuth
func (h *MaterialHandler) ListMaterials(c *gin.Context) {
	// Por ahora sin filtros (se pueden agregar después)
	materials, err := h.materialService.ListMaterials(c.Request.Context(), repository.ListFilters{})
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal error", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, materials)
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}
