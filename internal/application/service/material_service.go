package service

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/google/uuid"
)

// MaterialService define las operaciones de negocio para materiales
type MaterialService interface {
	CreateMaterial(ctx context.Context, req dto.CreateMaterialRequest, authorID string, schoolID string) (*dto.MaterialResponse, error)
	GetMaterial(ctx context.Context, id string) (*dto.MaterialResponse, error)
	GetMaterialWithVersions(ctx context.Context, id string) (*dto.MaterialWithVersionsResponse, error)
	NotifyUploadComplete(ctx context.Context, materialID string, req dto.UploadCompleteRequest) error
	ListMaterials(ctx context.Context, filters repository.ListFilters) ([]*dto.MaterialResponse, error)
}

type materialService struct {
	materialRepo     repository.MaterialRepository
	messagePublisher rabbitmq.Publisher
	logger           logger.Logger
}

func NewMaterialService(
	materialRepo repository.MaterialRepository,
	messagePublisher rabbitmq.Publisher,
	logger logger.Logger,
) MaterialService {
	return &materialService{
		materialRepo:     materialRepo,
		messagePublisher: messagePublisher,
		logger:           logger,
	}
}

func (s *materialService) CreateMaterial(
	ctx context.Context,
	req dto.CreateMaterialRequest,
	authorIDStr string,
	schoolIDStr string,
) (*dto.MaterialResponse, error) {
	// Validar request
	if err := req.Validate(); err != nil {
		s.logger.Warn("validation failed", "error", err)
		return nil, err
	}

	// Parsear author ID
	authorID, err := valueobject.UserIDFromString(authorIDStr)
	if err != nil {
		return nil, errors.NewValidationError("invalid author_id format")
	}

	// Parsear school ID del contexto de autenticación (JWT)
	schoolID, err := uuid.Parse(schoolIDStr)
	if err != nil {
		s.logger.Warn("invalid school_id format", "school_id", schoolIDStr, "error", err)
		return nil, errors.NewValidationError("invalid school_id format in authentication context")
	}

	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	var subject *string
	if req.Subject != "" {
		subject = &req.Subject
	}

	var grade *string
	if req.Grade != "" {
		grade = &req.Grade
	}

	material := &pgentities.Material{
		ID:                  valueobject.NewMaterialID().UUID().UUID,
		SchoolID:            schoolID,
		UploadedByTeacherID: authorID.UUID().UUID,
		AcademicUnitID:      nil,
		Title:               req.Title,
		Description:         description,
		Subject:             subject,
		Grade:               grade,
		FileURL:             "",
		FileType:            "",
		FileSizeBytes:       0,
		Status:              "uploaded",
		IsPublic:            false,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Persistir
	if err := s.materialRepo.Create(ctx, material); err != nil {
		s.logger.Error("failed to save material", "error", err)
		return nil, errors.NewDatabaseError("create material", err)
	}

	s.logger.Info("material created",
		"material_id", material.ID.String(),
		"author_id", authorID.String(),
		"title", material.Title,
	)

	// NOTA: El evento MaterialUploaded se publica en NotifyUploadComplete
	// cuando el frontend confirma que el archivo fue subido a S3 con datos reales

	return dto.ToMaterialResponse(material), nil
}

func (s *materialService) GetMaterial(ctx context.Context, id string) (*dto.MaterialResponse, error) {
	materialID, err := valueobject.MaterialIDFromString(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid material_id format")
	}

	material, err := s.materialRepo.FindByID(ctx, materialID)
	if err != nil || material == nil {
		return nil, errors.NewNotFoundError("material")
	}

	return dto.ToMaterialResponse(material), nil
}

func (s *materialService) NotifyUploadComplete(
	ctx context.Context,
	materialIDStr string,
	req dto.UploadCompleteRequest,
) error {
	// Validar
	if err := req.Validate(); err != nil {
		return err
	}

	materialID, err := valueobject.MaterialIDFromString(materialIDStr)
	if err != nil {
		return errors.NewValidationError("invalid material_id format")
	}

	// Buscar material
	material, err := s.materialRepo.FindByID(ctx, materialID)
	if err != nil || material == nil {
		return errors.NewNotFoundError("material")
	}

	// Actualizar con info de archivo
	material.FileURL = req.FileURL
	material.FileType = req.FileType
	material.FileSizeBytes = req.FileSizeBytes
	material.Status = "uploaded"
	material.UpdatedAt = time.Now()

	// Persistir
	if err := s.materialRepo.Update(ctx, material); err != nil {
		s.logger.Error("failed to update material", "error", err)
		return errors.NewDatabaseError("update material", err)
	}

	s.logger.Info("upload complete notified",
		"material_id", materialID.String(),
		"file_url", req.FileURL,
	)

	// Publicar evento MaterialUploaded con datos reales de S3
	payload := rabbitmq.MaterialUploadedPayload{
		MaterialID:    material.ID.String(),
		SchoolID:      material.SchoolID.String(),
		TeacherID:     material.UploadedByTeacherID.String(),
		FileURL:       req.FileURL,
		FileSizeBytes: req.FileSizeBytes,
		FileType:      req.FileType,
		Metadata: map[string]interface{}{
			"title":       material.Title,
			"description": material.Description,
		},
	}

	event := rabbitmq.NewMaterialUploadedEvent(payload)
	eventJSON, err := event.ToJSON()
	if err != nil {
		s.logger.Warn("failed to serialize material uploaded event",
			"material_id", material.ID.String(),
			"error", err,
		)
	} else {
		if err := s.messagePublisher.Publish(ctx, "edugo.materials", "material.uploaded", eventJSON); err != nil {
			s.logger.Warn("failed to publish material uploaded event",
				"material_id", material.ID.String(),
				"error", err,
			)
		} else {
			s.logger.Info("material uploaded event published",
				"material_id", material.ID.String(),
				"event_id", event.EventID,
				"file_url", req.FileURL,
				"file_size", req.FileSizeBytes,
			)
		}
	}

	return nil
}

func (s *materialService) ListMaterials(ctx context.Context, filters repository.ListFilters) ([]*dto.MaterialResponse, error) {
	materials, err := s.materialRepo.List(ctx, filters)
	if err != nil {
		s.logger.Error("failed to list materials", "error", err)
		return nil, errors.NewDatabaseError("list materials", err)
	}

	responses := make([]*dto.MaterialResponse, len(materials))
	for i, material := range materials {
		responses[i] = dto.ToMaterialResponse(material)
	}

	return responses, nil
}

// GetMaterialWithVersions obtiene un material incluyendo su historial completo de versiones
// Este método consulta el material junto con todas sus versiones en una sola operación de BD
func (s *materialService) GetMaterialWithVersions(ctx context.Context, id string) (*dto.MaterialWithVersionsResponse, error) {
	// Registrar inicio de operación para medir tiempo de ejecución
	startTime := time.Now()

	// Parsear y validar materialID
	materialID, err := valueobject.MaterialIDFromString(id)
	if err != nil {
		s.logger.Warn("invalid material_id format",
			"material_id", id,
			"error", err,
		)
		return nil, errors.NewValidationError("invalid material_id format")
	}

	// Invocar repository para obtener material con versiones
	material, versions, err := s.materialRepo.FindByIDWithVersions(ctx, materialID)
	if err != nil {
		s.logger.Error("failed to fetch material with versions",
			"material_id", materialID.String(),
			"error", err,
		)
		return nil, errors.NewDatabaseError("fetch material with versions", err)
	}

	// Validar que el material existe
	if material == nil {
		s.logger.Warn("material not found",
			"material_id", materialID.String(),
		)
		return nil, errors.NewNotFoundError("material")
	}

	// Transformar entidades de domain a DTOs
	response := dto.ToMaterialWithVersionsResponse(material, versions)

	// Logging contextual con métricas relevantes
	executionTime := time.Since(startTime)
	s.logger.Info("material with versions fetched successfully",
		"material_id", materialID.String(),
		"version_count", len(versions),
		"execution_time", executionTime,
	)

	return response, nil
}
