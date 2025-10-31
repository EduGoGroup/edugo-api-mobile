package service

import (
	"context"

	"github.com/edugo/api-mobile/internal/domain/entity"
	"github.com/edugo/api-mobile/internal/domain/repository"
	"github.com/edugo/api-mobile/internal/domain/valueobject"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/logger"
)

type ProgressService interface {
	UpdateProgress(ctx context.Context, materialID string, userID string, percentage int, lastPage int) error
}

type progressService struct {
	progressRepo repository.ProgressRepository
	logger       logger.Logger
}

func NewProgressService(progressRepo repository.ProgressRepository, logger logger.Logger) ProgressService {
	return &progressService{
		progressRepo: progressRepo,
		logger:       logger,
	}
}

func (s *progressService) UpdateProgress(ctx context.Context, materialID string, userIDStr string, percentage int, lastPage int) error {
	matID, err := valueobject.MaterialIDFromString(materialID)
	if err != nil {
		return errors.NewValidationError("invalid material_id")
	}

	userID, err := valueobject.UserIDFromString(userIDStr)
	if err != nil {
		return errors.NewValidationError("invalid user_id")
	}

	// Buscar o crear progress
	progress, err := s.progressRepo.FindByMaterialAndUser(ctx, matID, userID)
	if err != nil {
		return errors.NewDatabaseError("find progress", err)
	}

	if progress == nil {
		progress = entity.NewProgress(matID, userID)
		if err := progress.UpdateProgress(percentage, lastPage); err != nil {
			return err
		}
		if err := s.progressRepo.Save(ctx, progress); err != nil {
			return errors.NewDatabaseError("save progress", err)
		}
	} else {
		if err := progress.UpdateProgress(percentage, lastPage); err != nil {
			return err
		}
		if err := s.progressRepo.Update(ctx, progress); err != nil {
			return errors.NewDatabaseError("update progress", err)
		}
	}

	s.logger.Info("progress updated", "material_id", materialID, "percentage", percentage)
	return nil
}
