package scoring

import (
	"fmt"
	"strings"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
)

// MultipleChoiceStrategy implementa lógica de evaluación para preguntas de selección múltiple
type MultipleChoiceStrategy struct{}

// NewMultipleChoiceStrategy crea una nueva estrategia para multiple choice
func NewMultipleChoiceStrategy() ScoringStrategy {
	return &MultipleChoiceStrategy{}
}

// CalculateScore evalúa una pregunta de selección múltiple
// Compara la respuesta del usuario (case-insensitive, trimmed) con la respuesta correcta
func (s *MultipleChoiceStrategy) CalculateScore(question repository.AssessmentQuestion, userAnswer interface{}) (float64, bool, string) {
	// Convertir respuesta correcta a string
	correctAnswer, ok := question.CorrectAnswer.(string)
	if !ok {
		return 0.0, false, "Error interno: respuesta correcta mal configurada"
	}

	// Convertir respuesta de usuario a string
	userAnswerStr, ok := userAnswer.(string)
	if !ok {
		return 0.0, false, "Formato de respuesta inválido. Se esperaba una opción (A, B, C, D)"
	}

	// Normalizar respuestas: trim whitespace y convertir a lowercase
	normalizedUser := strings.ToLower(strings.TrimSpace(userAnswerStr))
	normalizedCorrect := strings.ToLower(strings.TrimSpace(correctAnswer))

	// Comparar respuestas
	isCorrect := normalizedUser == normalizedCorrect

	var explanation string
	if isCorrect {
		if question.Explanation != "" {
			explanation = fmt.Sprintf("¡Correcto! %s", question.Explanation)
		} else {
			explanation = fmt.Sprintf("¡Correcto! La opción %s es la respuesta adecuada.", strings.ToUpper(correctAnswer))
		}
		return 1.0, true, explanation
	}

	// Respuesta incorrecta
	if question.Explanation != "" {
		explanation = fmt.Sprintf("Incorrecto. Tu respuesta: %s. Respuesta correcta: %s. %s",
			strings.ToUpper(userAnswerStr),
			strings.ToUpper(correctAnswer),
			question.Explanation)
	} else {
		explanation = fmt.Sprintf("Incorrecto. Tu respuesta: %s. La respuesta correcta es: %s",
			strings.ToUpper(userAnswerStr),
			strings.ToUpper(correctAnswer))
	}

	return 0.0, false, explanation
}
