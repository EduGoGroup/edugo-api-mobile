package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/application/service"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/storage/s3"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	ginmiddleware "github.com/EduGoGroup/edugo-shared/middleware/gin"
)

// MaterialHandler maneja las peticiones HTTP relacionadas con materiales
type MaterialHandler struct {
	materialService service.MaterialService
	s3Storage       s3.S3Storage // Usar interfaz en lugar de tipo concreto
	logger          logger.Logger
}

func NewMaterialHandler(materialService service.MaterialService, s3Storage s3.S3Storage, logger logger.Logger) *MaterialHandler {
	return &MaterialHandler{
		materialService: materialService,
		s3Storage:       s3Storage,
		logger:          logger,
	}
}

// CreateMaterial godoc
// @Summary Create a new material
// @Description Creates a new educational material with title, description, and optional subject. Requires authentication.
// @Tags materials
// @Accept json
// @Produce json
// @Param request body dto.CreateMaterialRequest true "Material data (title, description, subject_id)"
// @Success 201 {object} dto.MaterialResponse "Material created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body or validation error"
// @Failure 401 {object} ErrorResponse "User not authenticated"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/materials [post]
// @Security BearerAuth
func (h *MaterialHandler) CreateMaterial(c *gin.Context) {
	var req dto.CreateMaterialRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	// Obtener user_id del contexto (middleware de autenticación)
	authorID := ginmiddleware.MustGetUserID(c)

	material, err := h.materialService.CreateMaterial(c.Request.Context(), req, authorID)
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
// @Description Retrieves a specific educational material by its unique identifier
// @Tags materials
// @Produce json
// @Param id path string true "Material ID (UUID format)"
// @Success 200 {object} dto.MaterialResponse "Material found successfully"
// @Failure 400 {object} ErrorResponse "Invalid material ID format"
// @Failure 404 {object} ErrorResponse "Material not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/materials/{id} [get]
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

// GetMaterialWithVersions godoc
// @Summary Get material with version history
// @Description Get a material including its complete version history
// @Tags materials
// @Produce json
// @Param id path string true "Material ID (UUID format)"
// @Success 200 {object} dto.MaterialWithVersionsResponse
// @Failure 400 {object} ErrorResponse "Invalid UUID format"
// @Failure 404 {object} ErrorResponse "Material not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/materials/{id}/versions [get]
// @Security BearerAuth
func (h *MaterialHandler) GetMaterialWithVersions(c *gin.Context) {
	// Obtener materialID del path parameter
	id := c.Param("id")

	// Invocar servicio para obtener material con versiones
	result, err := h.materialService.GetMaterialWithVersions(c.Request.Context(), id)
	if err != nil {
		// Convertir error de aplicación a respuesta HTTP apropiada
		if appErr, ok := errors.GetAppError(err); ok {
			h.logger.Warn("get material with versions failed",
				"material_id", id,
				"error", appErr.Message,
				"code", appErr.Code,
			)
			c.JSON(appErr.StatusCode, ErrorResponse{
				Error: appErr.Message,
				Code:  string(appErr.Code),
			})
			return
		}

		// Error inesperado (no es error de aplicación)
		h.logger.Error("unexpected error getting material with versions",
			"material_id", id,
			"error", err,
		)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "internal server error",
			Code:  "INTERNAL_ERROR",
		})
		return
	}

	// Retornar respuesta exitosa con material y versiones
	h.logger.Info("material with versions retrieved successfully",
		"material_id", id,
		"version_count", len(result.Versions),
	)
	c.JSON(http.StatusOK, result)
}

// NotifyUploadComplete godoc
// @Summary Notify upload complete
// @Description Notifies the system that a PDF file has been successfully uploaded to S3
// @Tags materials
// @Accept json
// @Produce json
// @Param id path string true "Material ID (UUID format)"
// @Param request body dto.UploadCompleteRequest true "S3 key and URL information"
// @Success 204 "Upload notification processed successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body or material ID"
// @Failure 404 {object} ErrorResponse "Material not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/materials/{id}/upload-complete [post]
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
// @Summary List all materials
// @Description Retrieves a list of all educational materials available in the system
// @Tags materials
// @Produce json
// @Success 200 {array} dto.MaterialResponse "List of materials retrieved successfully"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /v1/materials [get]
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

// GenerateUploadURL godoc
// @Summary Generate presigned upload URL
// @Description Generate a presigned URL for uploading a material file to S3
// @Tags materials
// @Accept json
// @Produce json
// @Param id path string true "Material ID"
// @Param request body dto.GenerateUploadURLRequest true "Upload URL request"
// @Success 200 {object} dto.GenerateUploadURLResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/materials/{id}/upload-url [post]
// @Security BearerAuth
func (h *MaterialHandler) GenerateUploadURL(c *gin.Context) {
	materialID := c.Param("id")
	var req dto.GenerateUploadURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body", Code: "INVALID_REQUEST"})
		return
	}

	// Verificar que el material existe
	_, err := h.materialService.GetMaterial(c.Request.Context(), materialID)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	// Validar y sanitizar el nombre del archivo para prevenir path traversal
	if strings.Contains(req.FileName, "..") || strings.Contains(req.FileName, "/") || strings.Contains(req.FileName, "\\") {
		h.logger.Warn("invalid file name with path traversal attempt",
			"material_id", materialID,
			"file_name", req.FileName,
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid file name: must not contain path separators",
			Code:  "INVALID_FILENAME",
		})
		return
	}

	// Construir key de S3: materials/{material_id}/{filename}
	s3Key := "materials/" + materialID + "/" + req.FileName

	// Generar URL presignada (válida por 15 minutos)
	uploadURL, err := h.s3Storage.GeneratePresignedUploadURL(
		c.Request.Context(),
		s3Key,
		req.ContentType,
		15*time.Minute,
	)
	if err != nil {
		h.logger.Error("error generating presigned upload URL",
			"material_id", materialID,
			"error", err,
		)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to generate upload URL",
			Code:  "S3_ERROR",
		})
		return
	}

	h.logger.Info("presigned upload URL generated",
		"material_id", materialID,
		"s3_key", s3Key,
	)

	c.JSON(http.StatusOK, dto.GenerateUploadURLResponse{
		UploadURL: uploadURL,
		FileURL:   s3Key,
		ExpiresIn: 900, // 15 minutos en segundos
	})
}

// GenerateDownloadURL godoc
// @Summary Generate presigned download URL
// @Description Generate a presigned URL for downloading a material file from S3
// @Tags materials
// @Produce json
// @Param id path string true "Material ID"
// @Success 200 {object} dto.GenerateDownloadURLResponse
// @Failure 404 {object} ErrorResponse
// @Router /v1/materials/{id}/download-url [get]
// @Security BearerAuth
func (h *MaterialHandler) GenerateDownloadURL(c *gin.Context) {
	materialID := c.Param("id")

	// Verificar que el material existe y obtener la S3 key
	material, err := h.materialService.GetMaterial(c.Request.Context(), materialID)
	if err != nil {
		if appErr, ok := errors.GetAppError(err); ok {
			c.JSON(appErr.StatusCode, ErrorResponse{Error: appErr.Message, Code: string(appErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error", Code: "INTERNAL_ERROR"})
		return
	}

	// Verificar que el material tiene una FileURL (fue subido)
	if material.FileURL == "" {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "material file not uploaded yet",
			Code:  "FILE_NOT_FOUND",
		})
		return
	}

	// Generar URL presignada para descarga (válida por 1 hora)
	downloadURL, err := h.s3Storage.GeneratePresignedDownloadURL(
		c.Request.Context(),
		material.FileURL,
		1*time.Hour,
	)
	if err != nil {
		h.logger.Error("error generating presigned download URL",
			"material_id", materialID,
			"file_url", material.FileURL,
			"error", err,
		)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to generate download URL",
			Code:  "S3_ERROR",
		})
		return
	}

	h.logger.Info("presigned download URL generated",
		"material_id", materialID,
		"file_url", material.FileURL,
	)

	c.JSON(http.StatusOK, dto.GenerateDownloadURLResponse{
		DownloadURL: downloadURL,
		ExpiresIn:   3600, // 1 hora en segundos
	})
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid request body"`
	Code  string `json:"code" example:"INVALID_REQUEST"`
}
