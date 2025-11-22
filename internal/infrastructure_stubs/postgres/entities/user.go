package entities

import (
	"time"

	"github.com/google/uuid"
)

// User representa un usuario del sistema
// STUB TEMPORAL - Simula github.com/EduGoGroup/edugo-infrastructure/postgres/entities.User
// TODO FASE 2: Reemplazar por import real de infrastructure
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Role      string    `gorm:"type:varchar(50);not null" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName especifica el nombre de la tabla en PostgreSQL
func (User) TableName() string {
	return "users"
}
