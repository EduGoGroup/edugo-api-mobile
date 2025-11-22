package entities

import (
	"time"

	"github.com/google/uuid"
)

// AssessmentAnswer representa una respuesta de un intento de evaluación
// STUB TEMPORAL - Simula github.com/EduGoGroup/edugo-infrastructure/postgres/entities.AssessmentAnswer
// TODO FASE 2: Reemplazar por import real de infrastructure
type AssessmentAnswer struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	AttemptID        uuid.UUID `gorm:"type:uuid;not null" json:"attempt_id"`
	QuestionID       string    `gorm:"type:varchar(24);not null" json:"question_id"`
	SelectedAnswerID string    `gorm:"type:varchar(24);not null" json:"selected_answer_id"`
	IsCorrect        bool      `gorm:"not null" json:"is_correct"`
	TimeSpentSeconds int       `gorm:"not null" json:"time_spent_seconds"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName especifica el nombre de la tabla en PostgreSQL
func (AssessmentAnswer) TableName() string {
	return "assessment_answers"
}
