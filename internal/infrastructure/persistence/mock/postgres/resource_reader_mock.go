package postgres

import (
	"context"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
)

// mockResourceReader es un stub para desarrollo sin BD
type mockResourceReader struct{}

// NewMockResourceReader crea una nueva instancia del mock
func NewMockResourceReader() repository.ResourceReader {
	return &mockResourceReader{}
}

func (r *mockResourceReader) GetMenuResources(ctx context.Context) ([]*repository.MenuResource, error) {
	return []*repository.MenuResource{}, nil
}

func (r *mockResourceReader) GetResourceScreenMappings(ctx context.Context, resourceKeys []string) ([]*repository.ResourceScreenMapping, error) {
	return []*repository.ResourceScreenMapping{}, nil
}
