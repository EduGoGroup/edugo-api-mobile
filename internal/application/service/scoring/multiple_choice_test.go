package scoring

import (
	"testing"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/stretchr/testify/assert"
)

func TestMultipleChoiceStrategy_CalculateScore(t *testing.T) {
	strategy := NewMultipleChoiceStrategy()

	tests := []struct {
		name              string
		question          repository.AssessmentQuestion
		userAnswer        interface{}
		expectedScore     float64
		expectedCorrect   bool
		expectExplanation bool
	}{
		{
			name: "respuesta_correcta_exacta",
			question: repository.AssessmentQuestion{
				ID:            "q1",
				QuestionText:  "¿Cuál es la capital de Francia?",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				Options:       []string{"A. Madrid", "B. París", "C. Londres", "D. Berlín"},
				CorrectAnswer: "B",
				Explanation:   "París es la capital de Francia.",
			},
			userAnswer:        "B",
			expectedScore:     1.0,
			expectedCorrect:   true,
			expectExplanation: true,
		},
		{
			name: "respuesta_correcta_case_insensitive",
			question: repository.AssessmentQuestion{
				ID:            "q2",
				QuestionText:  "¿Cuál es la respuesta correcta?",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				Options:       []string{"A", "B", "C", "D"},
				CorrectAnswer: "C",
				Explanation:   "",
			},
			userAnswer:        "c",
			expectedScore:     1.0,
			expectedCorrect:   true,
			expectExplanation: true,
		},
		{
			name: "respuesta_correcta_con_whitespace",
			question: repository.AssessmentQuestion{
				ID:            "q3",
				QuestionText:  "Pregunta",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				Options:       []string{"A", "B"},
				CorrectAnswer: "A",
				Explanation:   "",
			},
			userAnswer:        " A ",
			expectedScore:     1.0,
			expectedCorrect:   true,
			expectExplanation: true,
		},
		{
			name: "respuesta_incorrecta",
			question: repository.AssessmentQuestion{
				ID:            "q4",
				QuestionText:  "Pregunta",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				Options:       []string{"A", "B", "C"},
				CorrectAnswer: "B",
				Explanation:   "La opción B es correcta porque...",
			},
			userAnswer:        "A",
			expectedScore:     0.0,
			expectedCorrect:   false,
			expectExplanation: true,
		},
		{
			name: "respuesta_invalida_tipo_incorrecto",
			question: repository.AssessmentQuestion{
				ID:            "q5",
				QuestionText:  "Pregunta",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				Options:       []string{"A", "B"},
				CorrectAnswer: "A",
				Explanation:   "",
			},
			userAnswer:        123,
			expectedScore:     0.0,
			expectedCorrect:   false,
			expectExplanation: true,
		},
		{
			name: "respuesta_vacia",
			question: repository.AssessmentQuestion{
				ID:            "q6",
				QuestionText:  "Pregunta",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				Options:       []string{"A", "B"},
				CorrectAnswer: "A",
				Explanation:   "",
			},
			userAnswer:        "",
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

func TestMultipleChoiceStrategy_CorrectAnswer_MalConfigurado(t *testing.T) {
	strategy := NewMultipleChoiceStrategy()

	question := repository.AssessmentQuestion{
		ID:            "q_bad",
		QuestionText:  "Pregunta",
		QuestionType:  enum.AssessmentTypeMultipleChoice,
		Options:       []string{"A", "B"},
		CorrectAnswer: 123, // Tipo incorrecto
		Explanation:   "",
	}

	score, isCorrect, explanation := strategy.CalculateScore(question, "A")

	assert.Equal(t, 0.0, score)
	assert.False(t, isCorrect)
	assert.Contains(t, explanation, "Error interno")
}
