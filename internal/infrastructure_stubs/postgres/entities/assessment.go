package entities

import (
	"time"

	"github.com/google/uuid"
)

// Assessment representa una evaluación de un material educativo
// STUB TEMPORAL - Simula github.com/EduGoGroup/edugo-infrastructure/postgres/entities.Assessment
// TODO FASE 2: Reemplazar por import real de infrastructure
type Assessment struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	MaterialID       uuid.UUID `gorm:"type:uuid;not null" json:"material_id"`
	MongoDocumentID  string    `gorm:"type:varchar(24);not null" json:"mongo_document_id"`
	Title            string    `gorm:"type:varchar(255);not null" json:"title"`
	TotalQuestions   int       `gorm:"not null" json:"total_questions"`
	PassThreshold    int       `gorm:"not null" json:"pass_threshold"`
	MaxAttempts      *int      `gorm:"" json:"max_attempts"`
	TimeLimitMinutes *int      `gorm:"" json:"time_limit_minutes"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName especifica el nombre de la tabla en PostgreSQL
func (Assessment) TableName() string {
	return "assessment"
}
