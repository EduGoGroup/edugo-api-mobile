package scoring

import (
	"fmt"
	"strings"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
)

// TrueFalseStrategy implementa lógica de evaluación para preguntas de verdadero/falso
type TrueFalseStrategy struct{}

// NewTrueFalseStrategy crea una nueva estrategia para true/false
func NewTrueFalseStrategy() ScoringStrategy {
	return &TrueFalseStrategy{}
}

// CalculateScore evalúa una pregunta de verdadero/falso
// Acepta múltiples formatos: "true"/"false", "True"/"False", "1"/"0", "verdadero"/"falso"
func (s *TrueFalseStrategy) CalculateScore(question repository.AssessmentQuestion, userAnswer interface{}) (float64, bool, string) {
	// Convertir respuesta correcta a bool
	correctBool, err := normalizeToBool(question.CorrectAnswer)
	if err != nil {
		return 0.0, false, "Error interno: respuesta correcta mal configurada"
	}

	// Convertir respuesta de usuario a bool
	userBool, err := normalizeToBool(userAnswer)
	if err != nil {
		return 0.0, false, "Formato de respuesta inválido. Se esperaba: true/false, verdadero/falso, 1/0"
	}

	// Comparar respuestas
	isCorrect := userBool == correctBool

	// Preparar strings para explicación
	correctStr := boolToSpanish(correctBool)
	userStr := boolToSpanish(userBool)

	var explanation string
	if isCorrect {
		if question.Explanation != "" {
			explanation = fmt.Sprintf("¡Correcto! La respuesta es %s. %s", correctStr, question.Explanation)
		} else {
			explanation = fmt.Sprintf("¡Correcto! La respuesta es %s.", correctStr)
		}
		return 1.0, true, explanation
	}

	// Respuesta incorrecta
	if question.Explanation != "" {
		explanation = fmt.Sprintf("Incorrecto. Tu respuesta: %s. Respuesta correcta: %s. %s",
			userStr, correctStr, question.Explanation)
	} else {
		explanation = fmt.Sprintf("Incorrecto. Tu respuesta: %s. La respuesta correcta es: %s",
			userStr, correctStr)
	}

	return 0.0, false, explanation
}

// normalizeToBool convierte diferentes formatos de respuesta a bool
// Acepta: "true", "True", "1", "verdadero", "Verdadero" -> true
//
//	"false", "False", "0", "falso", "Falso" -> false
func normalizeToBool(value interface{}) (bool, error) {
	// Si ya es bool, retornar directamente
	if boolVal, ok := value.(bool); ok {
		return boolVal, nil
	}

	// Intentar convertir desde string
	strVal, ok := value.(string)
	if !ok {
		return false, fmt.Errorf("tipo de valor no soportado: %T", value)
	}

	// Normalizar: trim y lowercase
	normalized := strings.ToLower(strings.TrimSpace(strVal))

	switch normalized {
	case "true", "1", "verdadero", "sí", "si", "yes", "t", "v":
		return true, nil
	case "false", "0", "falso", "no", "f":
		return false, nil
	default:
		return false, fmt.Errorf("valor no reconocido: %s", strVal)
	}
}

// boolToSpanish convierte un bool a representación en español
func boolToSpanish(value bool) string {
	if value {
		return "Verdadero"
	}
	return "Falso"
}
