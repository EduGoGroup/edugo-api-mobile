package entities

import (
	"time"

	"github.com/google/uuid"
)

// AssessmentAttempt representa un intento de evaluación por parte de un estudiante
// STUB TEMPORAL - Simula github.com/EduGoGroup/edugo-infrastructure/postgres/entities.AssessmentAttempt
// TODO FASE 2: Reemplazar por import real de infrastructure
type AssessmentAttempt struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	AssessmentID     uuid.UUID `gorm:"type:uuid;not null" json:"assessment_id"`
	StudentID        uuid.UUID `gorm:"type:uuid;not null" json:"student_id"`
	Score            int       `gorm:"not null" json:"score"`
	TimeSpentSeconds int       `gorm:"not null" json:"time_spent_seconds"`
	StartedAt        time.Time `gorm:"not null" json:"started_at"`
	CompletedAt      time.Time `gorm:"not null" json:"completed_at"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName especifica el nombre de la tabla en PostgreSQL
func (AssessmentAttempt) TableName() string {
	return "assessment_attempts"
}
