package postgres

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
)

// mockScreenRepository es un stub para desarrollo sin BD
type mockScreenRepository struct{}

// NewMockScreenRepository crea una nueva instancia del mock
func NewMockScreenRepository() repository.ScreenRepository {
	return &mockScreenRepository{}
}

func (r *mockScreenRepository) GetCombinedScreen(ctx context.Context, screenKey string, userID uuid.UUID) (*repository.CombinedScreen, error) {
	return nil, nil
}

func (r *mockScreenRepository) GetScreensForResource(ctx context.Context, resourceKey string) ([]*repository.ResourceScreenInfo, error) {
	return []*repository.ResourceScreenInfo{}, nil
}

func (r *mockScreenRepository) GetUserPreferences(ctx context.Context, screenKey string, userID uuid.UUID) (json.RawMessage, error) {
	return json.RawMessage("{}"), nil
}

func (r *mockScreenRepository) SaveUserPreferences(ctx context.Context, screenKey string, userID uuid.UUID, prefs json.RawMessage) error {
	return nil
}
