package repository

import (
	"context"

	"github.com/edugo/api-mobile/internal/domain/entity"
	"github.com/edugo/api-mobile/internal/domain/valueobject"
)

// UserRepository define las operaciones de persistencia para User
type UserRepository interface {
	FindByID(ctx context.Context, id valueobject.UserID) (*entity.User, error)
	FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
}
