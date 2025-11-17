package service

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"go.uber.org/zap"
)

// MaterialService define las operaciones de negocio para materiales
type MaterialService interface {
	CreateMaterial(ctx context.Context, req dto.CreateMaterialRequest, authorID string) (*dto.MaterialResponse, error)
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

	// Crear entidad de dominio
	material, err := entity.NewMaterial(
		req.Title,
		req.Description,
		authorID,
		req.SubjectID,
	)
	if err != nil {
		return nil, err
	}

	// Persistir
	if err := s.materialRepo.Create(ctx, material); err != nil {
		s.logger.Error("failed to save material", "error", err)
		return nil, errors.NewDatabaseError("create material", err)
	}

	s.logger.Info("material created",
		"material_id", material.ID().String(),
		"author_id", authorID.String(),
		"title", material.Title(),
	)

	// Publicar evento de material creado (nuevo formato con envelope)
	payload := messaging.MaterialUploadedPayload{
		MaterialID:    material.ID().String(),
		SchoolID:      "00000000-0000-0000-0000-000000000000", // TODO: obtener school_id del contexto
		TeacherID:     authorID.String(),
		FileURL:       "s3://edugo/materials/" + material.ID().String() + ".pdf", // TODO: URL real de S3
		FileSizeBytes: 0, // TODO: obtener tamaño real del archivo
		FileType:      "application/pdf",
		Metadata: map[string]interface{}{
			"title":       material.Title(),
			"description": material.Description(),
		},
	}

	event := messaging.NewMaterialUploadedEvent(payload)
	eventJSON, err := event.ToJSON()
	if err != nil {
		s.logger.Warn("failed to serialize material uploaded event",
			zap.String("material_id", material.ID().String()),
			zap.Error(err),
		)
	} else {
		// Publicar evento de forma asíncrona (no bloqueante)
		if err := s.messagePublisher.Publish(ctx, "edugo.materials", "material.uploaded", eventJSON); err != nil {
			s.logger.Warn("failed to publish material uploaded event",
				zap.String("material_id", material.ID().String()),
				zap.Error(err),
			)
		} else {
			s.logger.Info("material uploaded event published",
				zap.String("material_id", material.ID().String()),
				zap.String("event_id", event.EventID),
			)
		}
	}

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

	// Actualizar con info de S3
	if err := material.SetS3Info(req.S3Key, req.S3URL); err != nil {
		return err
	}

	// Persistir
	if err := s.materialRepo.Update(ctx, material); err != nil {
		s.logger.Error("failed to update material", "error", err)
		return errors.NewDatabaseError("update material", err)
	}

	s.logger.Info("upload complete notified",
		"material_id", materialID.String(),
		"s3_key", req.S3Key,
	)

	// TODO: Aquí se debería publicar evento a RabbitMQ
	// usando shared/messaging para que el worker procese el PDF

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
			zap.String("material_id", id),
			zap.Error(err),
		)
		return nil, errors.NewValidationError("invalid material_id format")
	}

	// Invocar repository para obtener material con versiones
	material, versions, err := s.materialRepo.FindByIDWithVersions(ctx, materialID)
	if err != nil {
		s.logger.Error("failed to fetch material with versions",
			zap.String("material_id", materialID.String()),
			zap.Error(err),
		)
		return nil, errors.NewDatabaseError("fetch material with versions", err)
	}

	// Validar que el material existe
	if material == nil {
		s.logger.Warn("material not found",
			zap.String("material_id", materialID.String()),
		)
		return nil, errors.NewNotFoundError("material")
	}

	// Transformar entidades de domain a DTOs
	response := dto.ToMaterialWithVersionsResponse(material, versions)

	// Logging contextual con métricas relevantes
	executionTime := time.Since(startTime)
	s.logger.Info("material with versions fetched successfully",
		zap.String("material_id", materialID.String()),
		zap.Int("version_count", len(versions)),
		zap.Duration("execution_time", executionTime),
	)

	return response, nil
}
