package services

import (
	"errors"

	"github.com/google/uuid"

	pgentities "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure_stubs/postgres/entities"
	domainErrors "github.com/EduGoGroup/edugo-api-mobile/internal/domain/errors"
)

// AttemptDomainService contiene reglas de negocio de AssessmentAttempt y AssessmentAnswer
// Extrae la lógica que antes estaba embebida en entities.Attempt y entities.Answer
type AttemptDomainService struct{}

// NewAttemptDomainService crea una nueva instancia del servicio
func NewAttemptDomainService() *AttemptDomainService {
	return &AttemptDomainService{}
}

// IsPassed indica si un intento aprobó la evaluación
// Compara el score del intento con el threshold del assessment
func (s *AttemptDomainService) IsPassed(attempt *pgentities.AssessmentAttempt, passThreshold int) bool {
	return attempt.Score >= passThreshold
}

// GetCorrectAnswersCount cuenta respuestas correctas
func (s *AttemptDomainService) GetCorrectAnswersCount(answers []*pgentities.AssessmentAnswer) int {
	count := 0
	for _, answer := range answers {
		if answer.IsCorrect {
			count++
		}
	}
	return count
}

// GetIncorrectAnswersCount cuenta respuestas incorrectas
func (s *AttemptDomainService) GetIncorrectAnswersCount(answers []*pgentities.AssessmentAnswer) int {
	return len(answers) - s.GetCorrectAnswersCount(answers)
}

// GetTotalQuestions devuelve el total de preguntas respondidas
func (s *AttemptDomainService) GetTotalQuestions(answers []*pgentities.AssessmentAnswer) int {
	return len(answers)
}

// GetAccuracyPercentage calcula el porcentaje de acierto
func (s *AttemptDomainService) GetAccuracyPercentage(answers []*pgentities.AssessmentAnswer) int {
	total := len(answers)
	if total == 0 {
		return 0
	}

	correct := s.GetCorrectAnswersCount(answers)
	return (correct * 100) / total
}

// GetAverageTimePerQuestion calcula tiempo promedio por pregunta en segundos
func (s *AttemptDomainService) GetAverageTimePerQuestion(attempt *pgentities.AssessmentAttempt, totalQuestions int) int {
	if totalQuestions == 0 {
		return 0
	}

	return attempt.TimeSpentSeconds / totalQuestions
}

// ValidateAttempt valida un intento completo con sus respuestas
// Aplica todas las reglas de validación del dominio
func (s *AttemptDomainService) ValidateAttempt(attempt *pgentities.AssessmentAttempt, answers []*pgentities.AssessmentAnswer) error {
	// Validar IDs
	if attempt.ID == uuid.Nil {
		return domainErrors.ErrInvalidAttemptID
	}

	if attempt.AssessmentID == uuid.Nil {
		return domainErrors.ErrInvalidAssessmentID
	}

	if attempt.StudentID == uuid.Nil {
		return domainErrors.ErrInvalidStudentID
	}

	// Validar score (0-100)
	if attempt.Score < 0 || attempt.Score > 100 {
		return domainErrors.ErrInvalidScore
	}

	// Validar time spent (>0 y <=2 horas)
	if attempt.TimeSpentSeconds <= 0 || attempt.TimeSpentSeconds > 7200 {
		return domainErrors.ErrInvalidTimeSpent
	}

	// Validar tiempos
	if attempt.StartedAt.IsZero() {
		return domainErrors.ErrInvalidStartTime
	}

	if attempt.CompletedAt.IsZero() || !attempt.CompletedAt.After(attempt.StartedAt) {
		return domainErrors.ErrInvalidEndTime
	}

	// Validar que tenga respuestas
	if len(answers) == 0 {
		return domainErrors.ErrNoAnswersProvided
	}

	// Validar que el score calculado coincide con las respuestas
	correctCount := s.GetCorrectAnswersCount(answers)
	expectedScore := (correctCount * 100) / len(answers)

	if attempt.Score != expectedScore {
		return errors.New("domain: score mismatch with answers")
	}

	return nil
}

// ValidateAnswer valida una respuesta individual
func (s *AttemptDomainService) ValidateAnswer(answer *pgentities.AssessmentAnswer) error {
	// Validar ID
	if answer.ID == uuid.Nil {
		return domainErrors.ErrInvalidAnswerID
	}

	// Validar attempt ID
	if answer.AttemptID == uuid.Nil {
		return domainErrors.ErrInvalidAttemptID
	}

	// Validar question ID (ObjectId de MongoDB = 24 caracteres)
	if answer.QuestionID == "" {
		return domainErrors.ErrInvalidQuestionID
	}

	// Validar selected answer ID (ObjectId de MongoDB = 24 caracteres)
	if answer.SelectedAnswerID == "" {
		return domainErrors.ErrInvalidSelectedAnswerID
	}

	// Validar tiempo (debe ser positivo)
	if answer.TimeSpentSeconds < 0 {
		return domainErrors.ErrInvalidTimeSpent
	}

	return nil
}

// CalculateScore calcula el score basado en respuestas correctas
// Retorna porcentaje 0-100
func (s *AttemptDomainService) CalculateScore(answers []*pgentities.AssessmentAnswer) int {
	return s.GetAccuracyPercentage(answers)
}
