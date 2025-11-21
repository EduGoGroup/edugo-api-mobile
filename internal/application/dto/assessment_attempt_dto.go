package dto

import (
	"time"

	"github.com/google/uuid"
)

// AssessmentResponse representa la respuesta al obtener un cuestionario
// IMPORTANTE: No incluye correct_answer ni feedback (sanitizado)
type AssessmentResponse struct {
	AssessmentID         uuid.UUID     `json:"assessment_id"`
	MaterialID           uuid.UUID     `json:"material_id"`
	Title                string        `json:"title"`
	TotalQuestions       int           `json:"total_questions"`
	PassThreshold        int           `json:"pass_threshold"`
	MaxAttempts          *int          `json:"max_attempts,omitempty"`
	TimeLimitMinutes     *int          `json:"time_limit_minutes,omitempty"`
	EstimatedTimeMinutes int           `json:"estimated_time_minutes"`
	Questions            []QuestionDTO `json:"questions"`
}

// QuestionDTO representa una pregunta sin respuesta correcta (sanitizado)
type QuestionDTO struct {
	ID      string      `json:"id"`
	Text    string      `json:"text"`
	Type    string      `json:"type"`
	Options []OptionDTO `json:"options"`
	// ❌ NO incluir: CorrectAnswer, Feedback
}

// OptionDTO representa una opción de respuesta
type OptionDTO struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// CreateAttemptRequest representa el body para crear un intento
type CreateAttemptRequest struct {
	Answers          []UserAnswerDTO `json:"answers" binding:"required,min=1,dive"`
	TimeSpentSeconds int             `json:"time_spent_seconds" binding:"required,min=1,max=7200"`
}

// UserAnswerDTO representa una respuesta del usuario
type UserAnswerDTO struct {
	QuestionID       string `json:"question_id" binding:"required"`
	SelectedAnswerID string `json:"selected_answer_id" binding:"required"`
	TimeSpentSeconds int    `json:"time_spent_seconds" binding:"required,min=0"`
}

// AttemptResultResponse representa el resultado de un intento
type AttemptResultResponse struct {
	AttemptID         uuid.UUID           `json:"attempt_id"`
	AssessmentID      uuid.UUID           `json:"assessment_id"`
	Score             int                 `json:"score"`
	MaxScore          int                 `json:"max_score"`
	CorrectAnswers    int                 `json:"correct_answers"`
	TotalQuestions    int                 `json:"total_questions"`
	PassThreshold     int                 `json:"pass_threshold"`
	Passed            bool                `json:"passed"`
	TimeSpentSeconds  int                 `json:"time_spent_seconds"`
	StartedAt         time.Time           `json:"started_at"`
	CompletedAt       time.Time           `json:"completed_at"`
	Feedback          []AnswerFeedbackDTO `json:"feedback"`
	CanRetake         bool                `json:"can_retake"`
	PreviousBestScore *int                `json:"previous_best_score,omitempty"`
}

// AnswerFeedbackDTO representa el feedback de una respuesta
type AnswerFeedbackDTO struct {
	QuestionID     string `json:"question_id"`
	QuestionText   string `json:"question_text"`
	SelectedOption string `json:"selected_option"`
	CorrectAnswer  string `json:"correct_answer"`
	IsCorrect      bool   `json:"is_correct"`
	Message        string `json:"message"`
}

// AttemptSummaryDTO representa un resumen de intento para historial
type AttemptSummaryDTO struct {
	AttemptID        uuid.UUID `json:"attempt_id"`
	AssessmentID     uuid.UUID `json:"assessment_id"`
	MaterialID       uuid.UUID `json:"material_id"`
	MaterialTitle    string    `json:"material_title"`
	Score            int       `json:"score"`
	MaxScore         int       `json:"max_score"`
	Passed           bool      `json:"passed"`
	TimeSpentSeconds int       `json:"time_spent_seconds"`
	CompletedAt      time.Time `json:"completed_at"`
}

// AttemptHistoryResponse representa el historial de intentos
type AttemptHistoryResponse struct {
	Attempts   []AttemptSummaryDTO `json:"attempts"`
	TotalCount int                 `json:"total_count"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
}
