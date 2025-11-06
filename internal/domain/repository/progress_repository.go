package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
)

// ProgressRepository define operaciones para Progress
type ProgressRepository interface {
	Save(ctx context.Context, progress *entity.Progress) error
	FindByMaterialAndUser(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) (*entity.Progress, error)
	Update(ctx context.Context, progress *entity.Progress) error
	// Upsert realiza INSERT o UPDATE idempotente usando ON CONFLICT de PostgreSQL
	Upsert(ctx context.Context, progress *entity.Progress) (*entity.Progress, error)

	// CountActiveUsers cuenta usuarios únicos con actividad en últimos 30 días (para estadísticas)
	CountActiveUsers(ctx context.Context) (int64, error)

	// CalculateAverageProgress calcula el promedio de progreso de todos los usuarios
	CalculateAverageProgress(ctx context.Context) (float64, error)
}
