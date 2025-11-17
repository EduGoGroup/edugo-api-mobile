package repositories

import (
	"context"

	"github.com/google/uuid"

	"edugo-api-mobile/internal/domain/entities"
)

// AnswerRepository define el contrato para persistencia de respuestas
type AnswerRepository interface {
	// FindByAttemptID busca todas las respuestas de un intento
	FindByAttemptID(ctx context.Context, attemptID uuid.UUID) ([]*entities.Answer, error)

	// Save guarda una o múltiples respuestas (batch insert)
	Save(ctx context.Context, answers []*entities.Answer) error

	// FindByQuestionID busca todas las respuestas para una pregunta específica
	// Útil para analytics: identificar preguntas difíciles
	FindByQuestionID(ctx context.Context, questionID string, limit, offset int) ([]*entities.Answer, error)
}
