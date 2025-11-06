package scoring

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
)

// ShortAnswerStrategy implementa lógica de evaluación para preguntas de respuesta corta
// Usa comparación flexible: normalización de texto, eliminación de puntuación,
// y soporte para múltiples respuestas válidas separadas por "|"
type ShortAnswerStrategy struct{}

// NewShortAnswerStrategy crea una nueva estrategia para short answer
func NewShortAnswerStrategy() ScoringStrategy {
	return &ShortAnswerStrategy{}
}

// CalculateScore evalúa una pregunta de respuesta corta
// Características:
// - Normalización: lowercase, trim, eliminación de puntuación
// - Múltiples respuestas válidas: "París|Paris" acepta ambas
// - Comparación flexible para admitir variaciones ortográficas
func (s *ShortAnswerStrategy) CalculateScore(question repository.AssessmentQuestion, userAnswer interface{}) (float64, bool, string) {
	// Convertir respuesta correcta a string
	correctAnswer, ok := question.CorrectAnswer.(string)
	if !ok {
		return 0.0, false, "Error interno: respuesta correcta mal configurada"
	}

	// Convertir respuesta de usuario a string
	userAnswerStr, ok := userAnswer.(string)
	if !ok {
		return 0.0, false, "Formato de respuesta inválido. Se esperaba texto"
	}

	// Verificar que la respuesta no esté vacía
	if strings.TrimSpace(userAnswerStr) == "" {
		return 0.0, false, "No se proporcionó una respuesta"
	}

	// Normalizar respuesta de usuario
	normalizedUser := normalizeText(userAnswerStr)

	// Si la respuesta correcta contiene "|", significa que hay múltiples opciones válidas
	validAnswers := strings.Split(correctAnswer, "|")

	isCorrect := false
	var matchedAnswer string

	for _, validAnswer := range validAnswers {
		normalizedValid := normalizeText(validAnswer)
		if normalizedUser == normalizedValid {
			isCorrect = true
			matchedAnswer = strings.TrimSpace(validAnswer)
			break
		}
	}

	var explanation string
	if isCorrect {
		if question.Explanation != "" {
			explanation = fmt.Sprintf("¡Correcto! %s", question.Explanation)
		} else {
			explanation = fmt.Sprintf("¡Correcto! Tu respuesta \"%s\" es válida.", matchedAnswer)
		}
		return 1.0, true, explanation
	}

	// Respuesta incorrecta
	// Mostrar la primera opción válida (la más común)
	primaryAnswer := strings.TrimSpace(validAnswers[0])

	if question.Explanation != "" {
		explanation = fmt.Sprintf("Incorrecto. Tu respuesta: \"%s\". Respuesta correcta: \"%s\". %s",
			userAnswerStr, primaryAnswer, question.Explanation)
	} else {
		if len(validAnswers) > 1 {
			explanation = fmt.Sprintf("Incorrecto. Tu respuesta: \"%s\". Respuestas válidas: \"%s\"",
				userAnswerStr, correctAnswer)
		} else {
			explanation = fmt.Sprintf("Incorrecto. Tu respuesta: \"%s\". Respuesta correcta: \"%s\"",
				userAnswerStr, primaryAnswer)
		}
	}

	return 0.0, false, explanation
}

// normalizeText normaliza texto para comparación:
// - Convierte a lowercase
// - Elimina whitespace extra (trim y espacios múltiples)
// - Elimina puntuación común (.,;:!?¿¡)
// - Conserva caracteres con tilde y ñ
func normalizeText(text string) string {
	// Trim whitespace inicial y final
	normalized := strings.TrimSpace(text)

	// Convertir a lowercase
	normalized = strings.ToLower(normalized)

	// Eliminar puntuación común usando regex
	// Mantiene letras (incluyendo acentuadas), números y espacios
	reg := regexp.MustCompile(`[.,;:!?¿¡'"\-_(){}[\]]+`)
	normalized = reg.ReplaceAllString(normalized, "")

	// Reducir múltiples espacios a uno solo
	spaceReg := regexp.MustCompile(`\s+`)
	normalized = spaceReg.ReplaceAllString(normalized, " ")

	// Trim final después de procesamiento
	normalized = strings.TrimSpace(normalized)

	return normalized
}
