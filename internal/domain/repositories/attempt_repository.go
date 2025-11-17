package repositories

import (
	"context"

	"github.com/google/uuid"

	"edugo-api-mobile/internal/domain/entities"
)

// AttemptRepository define el contrato para persistencia de intentos
type AttemptRepository interface {
	// FindByID busca un intento por ID
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Attempt, error)

	// FindByStudentAndAssessment busca intentos de un estudiante en una evaluación
	FindByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) ([]*entities.Attempt, error)

	// Save guarda un intento (solo INSERT, no UPDATE - inmutable)
	Save(ctx context.Context, attempt *entities.Attempt) error

	// CountByStudentAndAssessment cuenta intentos de un estudiante
	CountByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) (int, error)

	// FindByStudent busca todos los intentos de un estudiante (historial)
	FindByStudent(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]*entities.Attempt, error)
}
