package scoring

import (
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
)

// ScoringStrategy define la interfaz para estrategias de cálculo de puntajes
// Permite implementar diferentes lógicas de evaluación según el tipo de pregunta
type ScoringStrategy interface {
	// CalculateScore calcula el puntaje de una pregunta basándose en la respuesta del usuario
	// Retorna:
	// - score: 1.0 si es correcta, 0.0 si es incorrecta
	// - isCorrect: true si la respuesta coincide con la correcta
	// - explanation: mensaje contextual sobre la evaluación
	CalculateScore(question repository.AssessmentQuestion, userAnswer interface{}) (score float64, isCorrect bool, explanation string)
}

// ScoringResult encapsula el resultado de evaluar una pregunta
type ScoringResult struct {
	QuestionID     string
	IsCorrect      bool
	Score          float64
	UserAnswer     string
	CorrectAnswer  string
	Explanation    string
}
