package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// ProgressReader define operaciones de lectura para Progress
// Principio ISP: Separar lectura de escritura y estadísticas
type ProgressReader interface {
	FindByMaterialAndUser(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) (*pgentities.Progress, error)
}

// ProgressWriter define operaciones de escritura para Progress
type ProgressWriter interface {
	Save(ctx context.Context, progress *pgentities.Progress) error
	Update(ctx context.Context, progress *pgentities.Progress) error
	// Upsert realiza INSERT o UPDATE idempotente usando ON CONFLICT de PostgreSQL
	Upsert(ctx context.Context, progress *pgentities.Progress) (*pgentities.Progress, error)
}

// ProgressStats define operaciones de estadísticas para Progress
type ProgressStats interface {
	// CountActiveUsers cuenta usuarios únicos con actividad en últimos 30 días (para estadísticas)
	CountActiveUsers(ctx context.Context) (int64, error)

	// CalculateAverageProgress calcula el promedio de progreso de todos los usuarios
	CalculateAverageProgress(ctx context.Context) (float64, error)
}

// ProgressRepository agrega todas las capacidades de Progress
// Las implementaciones deben cumplir con todas las interfaces segregadas
type ProgressRepository interface {
	ProgressReader
	ProgressWriter
	ProgressStats
}
