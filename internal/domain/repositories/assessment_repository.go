package repositories

import (
	"context"

	"github.com/google/uuid"

	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// AssessmentRepository define el contrato para persistencia de evaluaciones
type AssessmentRepository interface {
	// FindByID busca una evaluación por ID
	FindByID(ctx context.Context, id uuid.UUID) (*pgentities.Assessment, error)

	// FindByMaterialID busca una evaluación por material ID
	FindByMaterialID(ctx context.Context, materialID uuid.UUID) (*pgentities.Assessment, error)

	// Save guarda una evaluación (INSERT o UPDATE)
	Save(ctx context.Context, assessment *pgentities.Assessment) error

	// Delete elimina una evaluación
	Delete(ctx context.Context, id uuid.UUID) error
}
