package entities

import (
	"time"

	"github.com/google/uuid"
)

// MaterialVersion representa una versión de un material
// STUB TEMPORAL - Simula github.com/EduGoGroup/edugo-infrastructure/postgres/entities.MaterialVersion
// TODO FASE 2: Reemplazar por import real de infrastructure
type MaterialVersion struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	MaterialID uuid.UUID `gorm:"type:uuid;not null" json:"material_id"`
	VersionNum int       `gorm:"not null" json:"version_num"`
	S3Key      string    `gorm:"type:varchar(500)" json:"s3_key"`
	S3URL      string    `gorm:"type:varchar(1000)" json:"s3_url"`
	Changes    string    `gorm:"type:text" json:"changes"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName especifica el nombre de la tabla en PostgreSQL
func (MaterialVersion) TableName() string {
	return "material_versions"
}
