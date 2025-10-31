package service

import (
	"context"

	"github.com/edugo/api-mobile/internal/domain/valueobject"
	"github.com/edugo/shared/pkg/errors"
	"github.com/edugo/shared/pkg/logger"
)

type MaterialStats struct {
	TotalViews    int     `json:"total_views"`
	AvgProgress   float64 `json:"avg_progress"`
	TotalAttempts int     `json:"total_attempts"`
	AvgScore      float64 `json:"avg_score"`
}

type StatsService interface {
	GetMaterialStats(ctx context.Context, materialID string) (*MaterialStats, error)
}

type statsService struct {
	logger logger.Logger
}

func NewStatsService(logger logger.Logger) StatsService {
	return &statsService{logger: logger}
}

func (s *statsService) GetMaterialStats(ctx context.Context, materialID string) (*MaterialStats, error) {
	_, err := valueobject.MaterialIDFromString(materialID)
	if err != nil {
		return nil, errors.NewValidationError("invalid material_id")
	}

	// Mock stats por ahora (en prod hacer queries a DB)
	stats := &MaterialStats{
		TotalViews:    150,
		AvgProgress:   67.5,
		TotalAttempts: 45,
		AvgScore:      78.3,
	}

	return stats, nil
}
