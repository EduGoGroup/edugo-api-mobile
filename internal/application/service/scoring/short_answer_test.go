package scoring

import (
	"testing"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/stretchr/testify/assert"
)

func TestShortAnswerStrategy_CalculateScore(t *testing.T) {
	strategy := NewShortAnswerStrategy()

	tests := []struct {
		name             string
		question         repository.AssessmentQuestion
		userAnswer       interface{}
		expectedScore    float64
		expectedCorrect  bool
		expectExplanation bool
	}{
		{
			name: "respuesta_correcta_exacta",
			question: repository.AssessmentQuestion{
				ID:              "q1",
				QuestionText:    "¿Cuál es la capital de Francia?",
				QuestionType:    enum.AssessmentTypeShortAnswer,
				CorrectAnswer:   "París",
				Explanation:     "París es la capital de Francia.",
			},
			userAnswer:       "París",
			expectedScore:    1.0,
			expectedCorrect:  true,
			expectExplanation: true,
		},
		{
			name: "respuesta_correcta_case_insensitive",
			question: repository.AssessmentQuestion{
				ID:              "q2",
				QuestionText:    "¿Cuál es la capital de Francia?",
				QuestionType:    enum.AssessmentTypeShortAnswer,
				CorrectAnswer:   "París",
				Explanation:     "",
			},
			userAnswer:       "parís",
			expectedScore:    1.0,
			expectedCorrect:  true,
			expectExplanation: true,
		},
		{
			name: "respuesta_correcta_con_whitespace",
			question: repository.AssessmentQuestion{
				ID:              "q3",
				QuestionText:    "Pregunta",
				QuestionType:    enum.AssessmentTypeShortAnswer,
				CorrectAnswer:   "respuesta",
				Explanation:     "",
			},
			userAnswer:       "  respuesta  ",
			expectedScore:    1.0,
			expectedCorrect:  true,
			expectExplanation: true,
		},
		{
			name: "respuesta_correcta_con_puntuacion",
			question: repository.AssessmentQuestion{
				ID:              "q4",
				QuestionText:    "Pregunta",
				QuestionType:    enum.AssessmentTypeShortAnswer,
				CorrectAnswer:   "respuesta",
				Explanation:     "",
			},
			userAnswer:       "respuesta.",
			expectedScore:    1.0,
			expectedCorrect:  true,
			expectExplanation: true,
		},
		{
			name: "respuesta_correcta_multiples_opciones_primera",
			question: repository.AssessmentQuestion{
				ID:              "q5",
				QuestionText:    "¿Cuál es la capital de Francia?",
				QuestionType:    enum.AssessmentTypeShortAnswer,
				CorrectAnswer:   "París|Paris",
				Explanation:     "Aceptamos ambas ortografías.",
			},
			userAnswer:       "París",
			expectedScore:    1.0,
			expectedCorrect:  true,
			expectExplanation: true,
		},
		{
			name: "respuesta_correcta_multiples_opciones_segunda",
			question: repository.AssessmentQuestion{
				ID:              "q6",
				QuestionText:    "¿Cuál es la capital de Francia?",
				QuestionType:    enum.AssessmentTypeShortAnswer,
				CorrectAnswer:   "París|Paris",
				Explanation:     "",
			},
			userAnswer:       "Paris",
			expectedScore:    1.0,
			expectedCorrect:  true,
			expectExplanation: true,
		},
		{
			name: "respuesta_incorrecta",
			question: repository.AssessmentQuestion{
				ID:              "q7",
				QuestionText:    "¿Cuál es la capital de Francia?",
				QuestionType:    enum.AssessmentTypeShortAnswer,
				CorrectAnswer:   "París",
				Explanation:     "La capital de Francia es París.",
			},
			userAnswer:       "Madrid",
			expectedScore:    0.0,
			expectedCorrect:  false,
			expectExplanation: true,
		},
		{
			name: "respuesta_vacia",
			question: repository.AssessmentQuestion{
				ID:              "q8",
				QuestionText:    "Pregunta",
				QuestionType:    enum.AssessmentTypeShortAnswer,
				CorrectAnswer:   "respuesta",
				Explanation:     "",
			},
			userAnswer:       "",
			expectedScore:    0.0,
			expectedCorrect:  false,
			expectExplanation: true,
		},
		{
			name: "respuesta_solo_whitespace",
			question: repository.AssessmentQuestion{
				ID:              "q9",
				QuestionText:    "Pregunta",
				QuestionType:    enum.AssessmentTypeShortAnswer,
				CorrectAnswer:   "respuesta",
				Explanation:     "",
			},
			userAnswer:       "   ",
			expectedScore:    0.0,
			expectedCorrect:  false,
			expectExplanation: true,
		},
		{
			name: "respuesta_tipo_invalido",
			question: repository.AssessmentQuestion{
				ID:              "q10",
				QuestionText:    "Pregunta",
				QuestionType:    enum.AssessmentTypeShortAnswer,
				CorrectAnswer:   "respuesta",
				Explanation:     "",
			},
			userAnswer:       123,
			expectedScore:    0.0,
			expectedCorrect:  false,
			expectExplanation: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, isCorrect, explanation := strategy.CalculateScore(tt.question, tt.userAnswer)

			assert.Equal(t, tt.expectedScore, score, "Score should match")
			assert.Equal(t, tt.expectedCorrect, isCorrect, "IsCorrect should match")
			if tt.expectExplanation {
				assert.NotEmpty(t, explanation, "Explanation should not be empty")
			}
		})
	}
}

func TestShortAnswerStrategy_CorrectAnswer_MalConfigurado(t *testing.T) {
	strategy := NewShortAnswerStrategy()

	question := repository.AssessmentQuestion{
		ID:            "q_bad",
		QuestionText:  "Pregunta",
		QuestionType:  enum.AssessmentTypeShortAnswer,
		CorrectAnswer: 123, // Tipo incorrecto
		Explanation:   "",
	}

	score, isCorrect, explanation := strategy.CalculateScore(question, "respuesta")

	assert.Equal(t, 0.0, score)
	assert.False(t, isCorrect)
	assert.Contains(t, explanation, "Error interno")
}

func TestNormalizeText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"texto_simple", "hola", "hola"},
		{"texto_uppercase", "HOLA", "hola"},
		{"texto_mixto", "HoLa", "hola"},
		{"con_whitespace", "  hola  ", "hola"},
		{"con_puntuacion", "hola.", "hola"},
		{"con_multiples_puntuaciones", "¿Hola, cómo estás?", "hola cómo estás"},
		{"espacios_multiples", "hola    mundo", "hola mundo"},
		{"con_tilde", "París", "parís"},
		{"con_enie", "España", "españa"},
		{"combinado", "  ¡HOLA, Mundo!  ", "hola mundo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeText(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
