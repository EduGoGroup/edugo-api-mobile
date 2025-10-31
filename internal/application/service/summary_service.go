package service

import (
	"context"

	"github.com/edugo/api-mobile/internal/domain/repository"
	"github.com/edugo/api-mobile/internal/domain/valueobject"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/logger"
)

// SummaryService define operaciones para summaries
type SummaryService interface {
	GetSummary(ctx context.Context, materialID string) (*repository.MaterialSummary, error)
}

type summaryService struct {
	summaryRepo repository.SummaryRepository
	logger      logger.Logger
}

func NewSummaryService(summaryRepo repository.SummaryRepository, logger logger.Logger) SummaryService {
	return &summaryService{
		summaryRepo: summaryRepo,
		logger:      logger,
	}
}

func (s *summaryService) GetSummary(ctx context.Context, materialID string) (*repository.MaterialSummary, error) {
	matID, err := valueobject.MaterialIDFromString(materialID)
	if err != nil {
		return nil, errors.NewValidationError("invalid material_id")
	}

	summary, err := s.summaryRepo.FindByMaterialID(ctx, matID)
	if err != nil {
		s.logger.Error("failed to get summary", "error", err)
		return nil, errors.NewDatabaseError("get summary", err)
	}

	if summary == nil {
		return nil, errors.NewNotFoundError("summary")
	}

	return summary, nil
}
