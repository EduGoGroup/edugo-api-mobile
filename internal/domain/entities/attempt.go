package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"

	domainErrors "edugo-api-mobile/internal/domain/errors"
)

// Attempt representa un intento de un estudiante en una evaluación
// Esta entity es INMUTABLE: una vez creada, no se puede modificar
// Corresponde a la tabla `assessment_attempt` en PostgreSQL
type Attempt struct {
	ID               uuid.UUID
	AssessmentID     uuid.UUID
	StudentID        uuid.UUID
	Score            int // 0-100
	MaxScore         int // Siempre 100
	TimeSpentSeconds int // Tiempo total en segundos
	StartedAt        time.Time
	CompletedAt      time.Time
	CreatedAt        time.Time
	Answers          []*Answer // Respuestas del intento
	IdempotencyKey   *string   // Para prevenir duplicados (Post-MVP)
}

// NewAttempt crea un nuevo intento COMPLETO
// Business rule: un intento se crea YA COMPLETADO, no hay estado "en progreso"
func NewAttempt(
	assessmentID uuid.UUID,
	studentID uuid.UUID,
	answers []*Answer,
	startedAt time.Time,
	completedAt time.Time,
) (*Attempt, error) {
	// Validaciones básicas
	if assessmentID == uuid.Nil {
		return nil, domainErrors.ErrInvalidAssessmentID
	}

	if studentID == uuid.Nil {
		return nil, domainErrors.ErrInvalidStudentID
	}

	if len(answers) == 0 {
		return nil, domainErrors.ErrNoAnswersProvided
	}

	if startedAt.IsZero() {
		return nil, domainErrors.ErrInvalidStartTime
	}

	if completedAt.IsZero() || !completedAt.After(startedAt) {
		return nil, domainErrors.ErrInvalidEndTime
	}

	// Calcular tiempo gastado
	timeSpent := int(completedAt.Sub(startedAt).Seconds())
	if timeSpent <= 0 || timeSpent > 7200 { // Máximo 2 horas
		return nil, domainErrors.ErrInvalidTimeSpent
	}

	// Calcular score basándose en respuestas correctas
	correctAnswers := 0
	for _, answer := range answers {
		if answer.IsCorrect {
			correctAnswers++
		}
	}

	totalQuestions := len(answers)
	score := (correctAnswers * 100) / totalQuestions

	return &Attempt{
		ID:               uuid.New(),
		AssessmentID:     assessmentID,
		StudentID:        studentID,
		Score:            score,
		MaxScore:         100,
		TimeSpentSeconds: timeSpent,
		StartedAt:        startedAt.UTC(),
		CompletedAt:      completedAt.UTC(),
		CreatedAt:        time.Now().UTC(),
		Answers:          answers,
		IdempotencyKey:   nil,
	}, nil
}

// NewAttemptWithIdempotency crea intento con clave de idempotencia (Post-MVP)
func NewAttemptWithIdempotency(
	assessmentID uuid.UUID,
	studentID uuid.UUID,
	answers []*Answer,
	startedAt time.Time,
	completedAt time.Time,
	idempotencyKey string,
) (*Attempt, error) {
	attempt, err := NewAttempt(assessmentID, studentID, answers, startedAt, completedAt)
	if err != nil {
		return nil, err
	}

	attempt.IdempotencyKey = &idempotencyKey
	return attempt, nil
}

// IsPassed indica si el intento aprobó la evaluación
func (a *Attempt) IsPassed(passThreshold int) bool {
	return a.Score >= passThreshold
}

// GetCorrectAnswersCount retorna la cantidad de respuestas correctas
func (a *Attempt) GetCorrectAnswersCount() int {
	count := 0
	for _, answer := range a.Answers {
		if answer.IsCorrect {
			count++
		}
	}
	return count
}

// GetIncorrectAnswersCount retorna la cantidad de respuestas incorrectas
func (a *Attempt) GetIncorrectAnswersCount() int {
	return len(a.Answers) - a.GetCorrectAnswersCount()
}

// GetTotalQuestions retorna el total de preguntas respondidas
func (a *Attempt) GetTotalQuestions() int {
	return len(a.Answers)
}

// GetAccuracyPercentage retorna el porcentaje de precisión (alias de Score)
func (a *Attempt) GetAccuracyPercentage() int {
	return a.Score
}

// GetAverageTimePerQuestion retorna el tiempo promedio por pregunta en segundos
func (a *Attempt) GetAverageTimePerQuestion() int {
	if len(a.Answers) == 0 {
		return 0
	}
	return a.TimeSpentSeconds / len(a.Answers)
}

// Validate verifica la integridad del intento
func (a *Attempt) Validate() error {
	if a.ID == uuid.Nil {
		return domainErrors.ErrInvalidAttemptID
	}
	if a.AssessmentID == uuid.Nil {
		return domainErrors.ErrInvalidAssessmentID
	}
	if a.StudentID == uuid.Nil {
		return domainErrors.ErrInvalidStudentID
	}
	if a.Score < 0 || a.Score > 100 {
		return domainErrors.ErrInvalidScore
	}
	if a.TimeSpentSeconds <= 0 || a.TimeSpentSeconds > 7200 {
		return domainErrors.ErrInvalidTimeSpent
	}
	if a.StartedAt.IsZero() {
		return domainErrors.ErrInvalidStartTime
	}
	if a.CompletedAt.IsZero() || !a.CompletedAt.After(a.StartedAt) {
		return domainErrors.ErrInvalidEndTime
	}
	if len(a.Answers) == 0 {
		return domainErrors.ErrNoAnswersProvided
	}

	// Verificar que el score calculado coincide
	correctCount := a.GetCorrectAnswersCount()
	expectedScore := (correctCount * 100) / len(a.Answers)
	if a.Score != expectedScore {
		return errors.New("domain: score mismatch with answers")
	}

	return nil
}
