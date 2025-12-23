// Package repository contiene tipos de dominio para el sistema de assessments
package repository

import (
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// AssessmentQuestion representa una pregunta del quiz
// Usado por el sistema de scoring para evaluar respuestas
type AssessmentQuestion struct {
	ID              string
	QuestionText    string
	QuestionType    enum.AssessmentType
	Options         []string    // Para multiple choice
	CorrectAnswer   interface{} // String o int dependiendo del tipo
	Explanation     string
	DifficultyLevel string
}

// MaterialAssessment representa el quiz de un material (usado para conversiones)
type MaterialAssessment struct {
	MaterialID valueobject.MaterialID
	Questions  []AssessmentQuestion
	CreatedAt  string
}

// FeedbackItem representa el feedback detallado de una pregunta
type FeedbackItem struct {
	QuestionID    string
	IsCorrect     bool
	UserAnswer    string
	CorrectAnswer string
	Explanation   string
}

// AssessmentResult representa el resultado de una evaluación completada
type AssessmentResult struct {
	ID             string
	AssessmentID   string
	UserID         valueobject.UserID
	Score          float64
	TotalQuestions int
	CorrectAnswers int
	Feedback       []FeedbackItem
	SubmittedAt    string
}

// AssessmentAttempt representa un intento de resolver el quiz
type AssessmentAttempt struct {
	ID          string
	MaterialID  valueobject.MaterialID
	UserID      valueobject.UserID
	Answers     map[string]interface{} // question_id -> answer
	Score       float64
	AttemptedAt string
}
