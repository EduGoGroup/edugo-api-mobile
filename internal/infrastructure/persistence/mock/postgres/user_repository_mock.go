package postgres

import (
	"context"
	"sync"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/mock/fixtures"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

type mockUserRepository struct {
	users map[string]*pgentities.User
	mu    sync.RWMutex
}

func NewMockUserRepository() repository.UserRepository {
	return &mockUserRepository{users: fixtures.GetDefaultUsers()}
}

func (r *mockUserRepository) FindByID(ctx context.Context, id valueobject.UserID) (*pgentities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if u, ok := r.users[id.UUID().String()]; ok {
		c := *u
		return &c, nil
	}
	return nil, nil
}

func (r *mockUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*pgentities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email == email.String() {
			c := *u
			return &c, nil
		}
	}
	return nil, nil
}

func (r *mockUserRepository) Update(ctx context.Context, user *pgentities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validar que el usuario existe antes de actualizar
	if _, exists := r.users[user.ID.String()]; !exists {
		return nil // En producción esto retornaría un error, pero para mock retornamos nil
	}

	// Crear copia para evitar mutaciones externas
	userCopy := *user
	r.users[user.ID.String()] = &userCopy
	return nil
}
