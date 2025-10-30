package repository

import (
	"context"

	"github.com/edugo/api-mobile/internal/domain/entity"
	"github.com/edugo/api-mobile/internal/domain/valueobject"
)

// ProgressRepository define operaciones para Progress
type ProgressRepository interface {
	Save(ctx context.Context, progress *entity.Progress) error
	FindByMaterialAndUser(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) (*entity.Progress, error)
	Update(ctx context.Context, progress *entity.Progress) error
}
