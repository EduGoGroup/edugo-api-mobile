package repositories

import (
	"context"

	"github.com/google/uuid"

	"edugo-api-mobile/internal/domain/entities"
)

// AssessmentRepository define el contrato para persistencia de evaluaciones
type AssessmentRepository interface {
	// FindByID busca una evaluación por ID
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Assessment, error)

	// FindByMaterialID busca una evaluación por material ID
	FindByMaterialID(ctx context.Context, materialID uuid.UUID) (*entities.Assessment, error)

	// Save guarda una evaluación (INSERT o UPDATE)
	Save(ctx context.Context, assessment *entities.Assessment) error

	// Delete elimina una evaluación
	Delete(ctx context.Context, id uuid.UUID) error
}
