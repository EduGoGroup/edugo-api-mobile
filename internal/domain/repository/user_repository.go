package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

// UserReader define operaciones de lectura para User
// Principio ISP: Separar lectura de escritura
type UserReader interface {
	FindByID(ctx context.Context, id valueobject.UserID) (*pgentities.User, error)
	FindByEmail(ctx context.Context, email valueobject.Email) (*pgentities.User, error)
}

// UserWriter define operaciones de escritura para User
type UserWriter interface {
	Update(ctx context.Context, user *pgentities.User) error
}

// UserRepository agrega todas las capacidades de User
// Las implementaciones deben cumplir con todas las interfaces segregadas
type UserRepository interface {
	UserReader
	UserWriter
}
