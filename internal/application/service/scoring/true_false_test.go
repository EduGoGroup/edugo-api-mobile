package scoring

import (
	"testing"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/stretchr/testify/assert"
)

func TestTrueFalseStrategy_CalculateScore(t *testing.T) {
	strategy := NewTrueFalseStrategy()

	tests := []struct {
		name              string
		question          repository.AssessmentQuestion
		userAnswer        interface{}
		expectedScore     float64
		expectedCorrect   bool
		expectExplanation bool
	}{
		{
			name: "respuesta_correcta_true_string",
			question: repository.AssessmentQuestion{
				ID:            "q1",
				QuestionText:  "¿El sol es una estrella?",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: "true",
				Explanation:   "El sol es una estrella de tipo G.",
			},
			userAnswer:        "true",
			expectedScore:     1.0,
			expectedCorrect:   true,
			expectExplanation: true,
		},
		{
			name: "respuesta_correcta_true_bool",
			question: repository.AssessmentQuestion{
				ID:            "q2",
				QuestionText:  "¿El agua hierve a 100°C?",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: true,
				Explanation:   "A presión atmosférica normal.",
			},
			userAnswer:        true,
			expectedScore:     1.0,
			expectedCorrect:   true,
			expectExplanation: true,
		},
		{
			name: "respuesta_correcta_formato_1",
			question: repository.AssessmentQuestion{
				ID:            "q3",
				QuestionText:  "Pregunta",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: true,
				Explanation:   "",
			},
			userAnswer:        "1",
			expectedScore:     1.0,
			expectedCorrect:   true,
			expectExplanation: true,
		},
		{
			name: "respuesta_correcta_formato_verdadero",
			question: repository.AssessmentQuestion{
				ID:            "q4",
				QuestionText:  "Pregunta",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: "true",
				Explanation:   "",
			},
			userAnswer:        "verdadero",
			expectedScore:     1.0,
			expectedCorrect:   true,
			expectExplanation: true,
		},
		{
			name: "respuesta_correcta_false_string",
			question: repository.AssessmentQuestion{
				ID:            "q5",
				QuestionText:  "¿La tierra es plana?",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: "false",
				Explanation:   "La tierra es esférica.",
			},
			userAnswer:        "false",
			expectedScore:     1.0,
			expectedCorrect:   true,
			expectExplanation: true,
		},
		{
			name: "respuesta_correcta_false_bool",
			question: repository.AssessmentQuestion{
				ID:            "q6",
				QuestionText:  "Pregunta",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: false,
				Explanation:   "",
			},
			userAnswer:        false,
			expectedScore:     1.0,
			expectedCorrect:   true,
			expectExplanation: true,
		},
		{
			name: "respuesta_correcta_formato_0",
			question: repository.AssessmentQuestion{
				ID:            "q7",
				QuestionText:  "Pregunta",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: false,
				Explanation:   "",
			},
			userAnswer:        "0",
			expectedScore:     1.0,
			expectedCorrect:   true,
			expectExplanation: true,
		},
		{
			name: "respuesta_correcta_formato_falso",
			question: repository.AssessmentQuestion{
				ID:            "q8",
				QuestionText:  "Pregunta",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: false,
				Explanation:   "",
			},
			userAnswer:        "falso",
			expectedScore:     1.0,
			expectedCorrect:   true,
			expectExplanation: true,
		},
		{
			name: "respuesta_incorrecta",
			question: repository.AssessmentQuestion{
				ID:            "q9",
				QuestionText:  "Pregunta",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: true,
				Explanation:   "Explicación de por qué es verdadero.",
			},
			userAnswer:        "false",
			expectedScore:     0.0,
			expectedCorrect:   false,
			expectExplanation: true,
		},
		{
			name: "respuesta_formato_invalido",
			question: repository.AssessmentQuestion{
				ID:            "q10",
				QuestionText:  "Pregunta",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: true,
				Explanation:   "",
			},
			userAnswer:        "maybe",
			expectedScore:     0.0,
			expectedCorrect:   false,
			expectExplanation: true,
		},
		{
			name: "respuesta_tipo_invalido",
			question: repository.AssessmentQuestion{
				ID:            "q11",
				QuestionText:  "Pregunta",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: true,
				Explanation:   "",
			},
			userAnswer:        123,
			expectedScore:     0.0,
			expectedCorrect:   false,
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

func TestTrueFalseStrategy_CorrectAnswer_MalConfigurado(t *testing.T) {
	strategy := NewTrueFalseStrategy()

	question := repository.AssessmentQuestion{
		ID:            "q_bad",
		QuestionText:  "Pregunta",
		QuestionType:  enum.AssessmentTypeTrueFalse,
		CorrectAnswer: []int{1, 2, 3}, // Tipo completamente incorrecto
		Explanation:   "",
	}

	score, isCorrect, explanation := strategy.CalculateScore(question, "true")

	assert.Equal(t, 0.0, score)
	assert.False(t, isCorrect)
	assert.Contains(t, explanation, "Error interno")
}

func TestNormalizeToBool(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expected  bool
		expectErr bool
	}{
		{"bool_true", true, true, false},
		{"bool_false", false, false, false},
		{"string_true", "true", true, false},
		{"string_True", "True", true, false},
		{"string_1", "1", true, false},
		{"string_verdadero", "verdadero", true, false},
		{"string_false", "false", false, false},
		{"string_False", "False", false, false},
		{"string_0", "0", false, false},
		{"string_falso", "falso", false, false},
		{"string_invalid", "maybe", false, true},
		{"int_invalid", 42, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizeToBool(tt.input)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
