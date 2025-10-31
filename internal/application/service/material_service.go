package service

import (
	"context"

	"github.com/edugo/api-mobile/internal/application/dto"
	"github.com/edugo/api-mobile/internal/domain/entity"
	"github.com/edugo/api-mobile/internal/domain/repository"
	"github.com/edugo/api-mobile/internal/domain/valueobject"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/logger"
)

// MaterialService define las operaciones de negocio para materiales
type MaterialService interface {
	CreateMaterial(ctx context.Context, req dto.CreateMaterialRequest, authorID string) (*dto.MaterialResponse, error)
	GetMaterial(ctx context.Context, id string) (*dto.MaterialResponse, error)
	NotifyUploadComplete(ctx context.Context, materialID string, req dto.UploadCompleteRequest) error
	ListMaterials(ctx context.Context, filters repository.ListFilters) ([]*dto.MaterialResponse, error)
}

type materialService struct {
	materialRepo repository.MaterialRepository
	logger       logger.Logger
}

func NewMaterialService(
	materialRepo repository.MaterialRepository,
	logger logger.Logger,
) MaterialService {
	return &materialService{
		materialRepo: materialRepo,
		logger:       logger,
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
