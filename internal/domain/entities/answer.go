package entities

import (
	"time"

	"github.com/google/uuid"

	domainErrors "edugo-api-mobile/internal/domain/errors"
)

// Answer representa una respuesta individual a una pregunta en un intento
// Corresponde a la tabla `assessment_attempt_answer` en PostgreSQL
type Answer struct {
	ID               uuid.UUID
	AttemptID        uuid.UUID
	QuestionID       string // ID de la pregunta en MongoDB
	SelectedAnswerID string // ID de la opción seleccionada en MongoDB
	IsCorrect        bool
	TimeSpentSeconds int
	CreatedAt        time.Time
}

// NewAnswer crea una nueva respuesta
func NewAnswer(
	attemptID uuid.UUID,
	questionID string,
	selectedAnswerID string,
	isCorrect bool,
	timeSpent int,
) (*Answer, error) {
	if attemptID == uuid.Nil {
		return nil, domainErrors.ErrInvalidAttemptID
	}

	if questionID == "" {
		return nil, domainErrors.ErrInvalidQuestionID
	}

	if selectedAnswerID == "" {
		return nil, domainErrors.ErrInvalidSelectedAnswerID
	}

	if timeSpent < 0 {
		return nil, domainErrors.ErrInvalidTimeSpent
	}

	return &Answer{
		ID:               uuid.New(),
		AttemptID:        attemptID,
		QuestionID:       questionID,
		SelectedAnswerID: selectedAnswerID,
		IsCorrect:        isCorrect,
		TimeSpentSeconds: timeSpent,
		CreatedAt:        time.Now().UTC(),
	}, nil
}

// Validate verifica la validez de la respuesta
func (a *Answer) Validate() error {
	if a.ID == uuid.Nil {
		return domainErrors.ErrInvalidAnswerID
	}
	if a.AttemptID == uuid.Nil {
		return domainErrors.ErrInvalidAttemptID
	}
	if a.QuestionID == "" {
		return domainErrors.ErrInvalidQuestionID
	}
	if a.SelectedAnswerID == "" {
		return domainErrors.ErrInvalidSelectedAnswerID
	}
	if a.TimeSpentSeconds < 0 {
		return domainErrors.ErrInvalidTimeSpent
	}
	return nil
}
