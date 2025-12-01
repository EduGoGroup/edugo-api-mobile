package postgres

import (
	"context"
	"sync"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/mock/fixtures"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/google/uuid"
)

type materialRepositoryMock struct {
	materials map[uuid.UUID]*pgentities.Material
	mu        sync.RWMutex
}

// NewMockMaterialRepository crea un nuevo repositorio mock de materiales con datos predeterminados
func NewMockMaterialRepository() repository.MaterialRepository {
	return &materialRepositoryMock{materials: fixtures.GetDefaultMaterials()}
}

func (r *materialRepositoryMock) FindByID(ctx context.Context, id valueobject.MaterialID) (*pgentities.Material, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if m, ok := r.materials[id.UUID().UUID]; ok {
		copy := *m
		return &copy, nil
	}
	return nil, nil
}

func (r *materialRepositoryMock) FindByIDWithVersions(ctx context.Context, id valueobject.MaterialID) (*pgentities.Material, []*pgentities.MaterialVersion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if m, ok := r.materials[id.UUID().UUID]; ok {
		copy := *m
		// Mock: retornar material sin versiones históricas
		return &copy, []*pgentities.MaterialVersion{}, nil
	}
	return nil, nil, nil
}

func (r *materialRepositoryMock) List(ctx context.Context, filters repository.ListFilters) ([]*pgentities.Material, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*pgentities.Material

	for _, m := range r.materials {
		// Aplicar filtros
		if filters.Status != nil && m.Status != string(*filters.Status) {
			continue
		}
		if filters.AuthorID != nil && m.UploadedByTeacherID != filters.AuthorID.UUID().UUID {
			continue
		}

		copy := *m
		result = append(result, &copy)
	}

	// Aplicar paginación
	start := filters.Offset
	if start > len(result) {
		return []*pgentities.Material{}, nil
	}

	end := start + filters.Limit
	if end > len(result) || filters.Limit == 0 {
		end = len(result)
	}

	return result[start:end], nil
}

func (r *materialRepositoryMock) FindByAuthor(ctx context.Context, authorID valueobject.UserID) ([]*pgentities.Material, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*pgentities.Material
	for _, m := range r.materials {
		if m.UploadedByTeacherID == authorID.UUID().UUID {
			copy := *m
			result = append(result, &copy)
		}
	}
	return result, nil
}

func (r *materialRepositoryMock) Create(ctx context.Context, material *pgentities.Material) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copy := *material
	copy.CreatedAt = time.Now()
	copy.UpdatedAt = time.Now()
	r.materials[material.ID] = &copy
	return nil
}

func (r *materialRepositoryMock) Update(ctx context.Context, material *pgentities.Material) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.materials[material.ID]; !exists {
		return nil
	}

	copy := *material
	copy.UpdatedAt = time.Now()
	r.materials[material.ID] = &copy
	return nil
}

func (r *materialRepositoryMock) UpdateStatus(ctx context.Context, id valueobject.MaterialID, status enum.MaterialStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if m, exists := r.materials[id.UUID().UUID]; exists {
		m.Status = string(status)
		m.UpdatedAt = time.Now()
	}
	return nil
}

func (r *materialRepositoryMock) UpdateProcessingStatus(ctx context.Context, id valueobject.MaterialID, status enum.ProcessingStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if m, exists := r.materials[id.UUID().UUID]; exists {
		now := time.Now()
		switch status {
		case enum.ProcessingStatusProcessing:
			m.ProcessingStartedAt = &now
		case enum.ProcessingStatusCompleted:
			m.ProcessingCompletedAt = &now
		}
		m.UpdatedAt = now
	}
	return nil
}

func (r *materialRepositoryMock) CountPublishedMaterials(ctx context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var count int64
	for _, m := range r.materials {
		if m.Status == "ready" && m.IsPublic {
			count++
		}
	}
	return count, nil
}
