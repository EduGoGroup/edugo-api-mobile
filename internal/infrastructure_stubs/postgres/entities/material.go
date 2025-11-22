package entities

import (
	"time"

	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/google/uuid"
)

// Material representa un material educativo (estructura de BD)
// STUB TEMPORAL - Simula github.com/EduGoGroup/edugo-infrastructure/postgres/entities.Material
// TODO FASE 2: Reemplazar por import real de infrastructure
type Material struct {
	ID               uuid.UUID              `gorm:"type:uuid;primary_key" json:"id"`
	Title            string                 `gorm:"type:varchar(255);not null" json:"title"`
	Description      string                 `gorm:"type:text" json:"description"`
	AuthorID         uuid.UUID              `gorm:"type:uuid;not null" json:"author_id"`
	SubjectID        string                 `gorm:"type:varchar(100)" json:"subject_id"`
	S3Key            string                 `gorm:"type:varchar(500)" json:"s3_key"`
	S3URL            string                 `gorm:"type:varchar(1000)" json:"s3_url"`
	Status           enum.MaterialStatus    `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
	ProcessingStatus enum.ProcessingStatus  `gorm:"type:varchar(20);not null;default:'pending'" json:"processing_status"`
	CreatedAt        time.Time              `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time              `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName especifica el nombre de la tabla en PostgreSQL
func (Material) TableName() string {
	return "materials"
}
