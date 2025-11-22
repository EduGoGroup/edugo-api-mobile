package entities

import (
	"time"

	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/google/uuid"
)

// Progress representa el progreso de un usuario en un material
// STUB TEMPORAL - Simula github.com/EduGoGroup/edugo-infrastructure/postgres/entities.Progress
// TODO FASE 2: Reemplazar por import real de infrastructure
type Progress struct {
	ID             uuid.UUID           `gorm:"type:uuid;primary_key" json:"id"`
	UserID         uuid.UUID           `gorm:"type:uuid;not null" json:"user_id"`
	MaterialID     uuid.UUID           `gorm:"type:uuid;not null" json:"material_id"`
	Percentage     int                 `gorm:"not null;default:0" json:"percentage"`
	LastPage       int                 `gorm:"not null;default:0" json:"last_page"`
	Status         enum.ProgressStatus `gorm:"type:varchar(20);not null;default:'not_started'" json:"status"`
	LastAccessedAt time.Time           `gorm:"" json:"last_accessed_at"`
	CreatedAt      time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName especifica el nombre de la tabla en PostgreSQL
func (Progress) TableName() string {
	return "progress"
}
