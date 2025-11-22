package services

import (
	"errors"

	"github.com/google/uuid"

	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	sharedErrors "github.com/EduGoGroup/edugo-shared/common/errors"
)

// AttemptDomainService contiene reglas de negocio de AssessmentAttempt y AssessmentAttemptAnswer
// Adaptado a estructura REAL de infrastructure (muchos campos nullable)
type AttemptDomainService struct{}

// NewAttemptDomainService crea una nueva instancia del servicio
func NewAttemptDomainService() *AttemptDomainService {
	return &AttemptDomainService{}
}

// IsPassed indica si un intento aprobó la evaluación
// Compara el score del intento con el threshold del assessment
func (s *AttemptDomainService) IsPassed(attempt *pgentities.AssessmentAttempt, passThreshold int) bool {
	if attempt.Score == nil {
		return false
	}
	return *attempt.Score >= float64(passThreshold)
}

// GetCorrectAnswersCount cuenta respuestas correctas
func (s *AttemptDomainService) GetCorrectAnswersCount(answers []*pgentities.AssessmentAttemptAnswer) int {
	count := 0
	for _, answer := range answers {
		if answer.IsCorrect != nil && *answer.IsCorrect {
			count++
		}
	}
	return count
}

// GetIncorrectAnswersCount cuenta respuestas incorrectas
func (s *AttemptDomainService) GetIncorrectAnswersCount(answers []*pgentities.AssessmentAttemptAnswer) int {
	total := 0
	for _, answer := range answers {
		if answer.IsCorrect != nil && !*answer.IsCorrect {
			total++
		}
	}
	return total
}

// GetTotalQuestions devuelve el total de preguntas respondidas
func (s *AttemptDomainService) GetTotalQuestions(answers []*pgentities.AssessmentAttemptAnswer) int {
	return len(answers)
}

// GetAccuracyPercentage calcula el porcentaje de acierto
func (s *AttemptDomainService) GetAccuracyPercentage(answers []*pgentities.AssessmentAttemptAnswer) float64 {
	total := len(answers)
	if total == 0 {
		return 0
	}

	correct := s.GetCorrectAnswersCount(answers)
	return (float64(correct) * 100.0) / float64(total)
}

// GetAverageTimePerQuestion calcula tiempo promedio por pregunta en segundos
func (s *AttemptDomainService) GetAverageTimePerQuestion(attempt *pgentities.AssessmentAttempt, totalQuestions int) int {
	if totalQuestions == 0 || attempt.TimeSpentSeconds == nil {
		return 0
	}

	return *attempt.TimeSpentSeconds / totalQuestions
}

// ValidateAttempt valida una entity AssessmentAttempt con sus respuestas
// Adaptado a campos nullable de infrastructure
func (s *AttemptDomainService) ValidateAttempt(attempt *pgentities.AssessmentAttempt, answers []*pgentities.AssessmentAttemptAnswer) error {
	// Validar IDs
	if attempt.ID == uuid.Nil {
		return sharedErrors.NewValidationError("invalid attempt id")
	}

	if attempt.AssessmentID == uuid.Nil {
		return sharedErrors.NewValidationError("invalid assessment id")
	}

	if attempt.StudentID == uuid.Nil {
		return sharedErrors.NewValidationError("invalid student id")
	}

	// Validar score (0-100) - nullable
	if attempt.Score != nil && (*attempt.Score < 0 || *attempt.Score > 100) {
		return sharedErrors.NewValidationError("score must be between 0 and 100")
	}

	// Validar time spent - nullable
	if attempt.TimeSpentSeconds != nil && (*attempt.TimeSpentSeconds <= 0 || *attempt.TimeSpentSeconds > 7200) {
		return sharedErrors.NewValidationError("time spent must be between 1 and 7200 seconds")
	}

	// Validar timestamps
	if attempt.StartedAt.IsZero() {
		return sharedErrors.NewValidationError("started_at is required")
	}

	// CompletedAt es nullable - solo validar si está presente
	if attempt.CompletedAt != nil && !attempt.CompletedAt.After(attempt.StartedAt) {
		return sharedErrors.NewValidationError("completed_at must be after started_at")
	}

	// Validar status (según migration: in_progress, completed, abandoned)
	validStatuses := map[string]bool{
		"in_progress": true,
		"completed":   true,
		"abandoned":   true,
	}
	if !validStatuses[attempt.Status] {
		return sharedErrors.NewValidationError("invalid status - must be in_progress, completed, or abandoned")
	}

	// Validar que existan respuestas si el intento está completado
	if attempt.Status == "completed" && len(answers) == 0 {
		return sharedErrors.NewValidationError("completed attempt must have at least one answer")
	}

	// Verificar que el score calculado coincide con el score almacenado (si ambos existen)
	if attempt.Score != nil && len(answers) > 0 {
		correctCount := s.GetCorrectAnswersCount(answers)
		expectedScore := (float64(correctCount) * 100.0) / float64(len(answers))
		// Permitir diferencia de 0.01 por redondeo
		if *attempt.Score < expectedScore-0.01 || *attempt.Score > expectedScore+0.01 {
			return errors.New("score mismatch with answers")
		}
	}

	return nil
}

// ValidateAnswer valida una entity AssessmentAttemptAnswer
// Adaptado a campos nullable de infrastructure
func (s *AttemptDomainService) ValidateAnswer(answer *pgentities.AssessmentAttemptAnswer) error {
	// Validar IDs
	if answer.ID == uuid.Nil {
		return sharedErrors.NewValidationError("invalid answer id")
	}

	if answer.AttemptID == uuid.Nil {
		return sharedErrors.NewValidationError("invalid attempt id")
	}

	// Validar question index (debe ser >= 0)
	if answer.QuestionIndex < 0 {
		return sharedErrors.NewValidationError("question index must be >= 0")
	}

	// StudentAnswer es nullable - solo validar si está presente
	if answer.StudentAnswer != nil && *answer.StudentAnswer == "" {
		return sharedErrors.NewValidationError("student answer cannot be empty string if provided")
	}

	// TimeSpentSeconds es nullable - solo validar si está presente
	if answer.TimeSpentSeconds != nil && *answer.TimeSpentSeconds < 0 {
		return sharedErrors.NewValidationError("time spent cannot be negative")
	}

	// PointsEarned y MaxPoints deben ser consistentes si ambos existen
	if answer.PointsEarned != nil && answer.MaxPoints != nil && *answer.PointsEarned > *answer.MaxPoints {
		return sharedErrors.NewValidationError("points earned cannot exceed max points")
	}

	return nil
}

// CalculateScore calcula el score basado en respuestas
// Útil para calcular el score antes de guardar
func (s *AttemptDomainService) CalculateScore(answers []*pgentities.AssessmentAttemptAnswer) float64 {
	return s.GetAccuracyPercentage(answers)
}

// IsCompleted indica si el intento está completado
func (s *AttemptDomainService) IsCompleted(attempt *pgentities.AssessmentAttempt) bool {
	return attempt.Status == "completed"
}

// IsInProgress indica si el intento está en progreso
func (s *AttemptDomainService) IsInProgress(attempt *pgentities.AssessmentAttempt) bool {
	return attempt.Status == "in_progress"
}

// IsAbandoned indica si el intento fue abandonado
func (s *AttemptDomainService) IsAbandoned(attempt *pgentities.AssessmentAttempt) bool {
	return attempt.Status == "abandoned"
}
