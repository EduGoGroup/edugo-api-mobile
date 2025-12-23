package repositories

import (
	"context"

	"github.com/google/uuid"

	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// AttemptRepository define el contrato para persistencia de intentos
// También implementa AssessmentStats para estadísticas globales
type AttemptRepository interface {
	// FindByID busca un intento por ID
	FindByID(ctx context.Context, id uuid.UUID) (*pgentities.AssessmentAttempt, error)

	// FindByStudentAndAssessment busca intentos de un estudiante en una evaluación
	FindByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) ([]*pgentities.AssessmentAttempt, error)

	// Save guarda un intento (solo INSERT, no UPDATE - inmutable)
	Save(ctx context.Context, attempt *pgentities.AssessmentAttempt) error

	// CountByStudentAndAssessment cuenta intentos de un estudiante
	CountByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) (int, error)

	// FindByStudent busca todos los intentos de un estudiante (historial)
	FindByStudent(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]*pgentities.AssessmentAttempt, error)

	// === Métodos de Estadísticas (implementa AssessmentStats) ===

	// CountCompletedAssessments cuenta el total de evaluaciones completadas
	CountCompletedAssessments(ctx context.Context) (int64, error)

	// CalculateAverageScore calcula el promedio de puntajes de evaluaciones completadas
	CalculateAverageScore(ctx context.Context) (float64, error)
}
