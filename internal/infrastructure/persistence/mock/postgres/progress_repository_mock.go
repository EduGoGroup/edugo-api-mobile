package postgres

import (
	"context"
	"sync"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/mock/fixtures"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

type progressRepositoryMock struct {
	progress map[fixtures.ProgressKey]*pgentities.Progress
	mu       sync.RWMutex
}

// NewMockProgressRepository crea un nuevo repositorio mock de progreso con datos predeterminados
func NewMockProgressRepository() repository.ProgressRepository {
	return &progressRepositoryMock{progress: fixtures.GetDefaultProgress()}
}

func (r *progressRepositoryMock) FindByMaterialAndUser(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) (*pgentities.Progress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := fixtures.ProgressKey{
		MaterialID: materialID.UUID().UUID,
		UserID:     userID.UUID().UUID,
	}

	if p, ok := r.progress[key]; ok {
		copy := *p
		return &copy, nil
	}
	return nil, nil
}

func (r *progressRepositoryMock) Save(ctx context.Context, progress *pgentities.Progress) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fixtures.ProgressKey{
		MaterialID: progress.MaterialID,
		UserID:     progress.UserID,
	}

	copy := *progress
	now := time.Now()
	copy.CreatedAt = now
	copy.UpdatedAt = now
	r.progress[key] = &copy
	return nil
}

func (r *progressRepositoryMock) Update(ctx context.Context, progress *pgentities.Progress) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fixtures.ProgressKey{
		MaterialID: progress.MaterialID,
		UserID:     progress.UserID,
	}

	if _, exists := r.progress[key]; exists {
		copy := *progress
		copy.UpdatedAt = time.Now()
		r.progress[key] = &copy
	}
	return nil
}

func (r *progressRepositoryMock) Upsert(ctx context.Context, progress *pgentities.Progress) (*pgentities.Progress, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fixtures.ProgressKey{
		MaterialID: progress.MaterialID,
		UserID:     progress.UserID,
	}

	copy := *progress
	now := time.Now()

	if existing, exists := r.progress[key]; exists {
		// Update
		copy.CreatedAt = existing.CreatedAt
		copy.UpdatedAt = now
	} else {
		// Insert
		copy.CreatedAt = now
		copy.UpdatedAt = now
	}

	r.progress[key] = &copy
	return &copy, nil
}

func (r *progressRepositoryMock) CountActiveUsers(ctx context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Contar usuarios únicos con actividad en últimos 30 días
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	activeUsers := make(map[uuid.UUID]bool)

	for _, p := range r.progress {
		if p.LastAccessedAt.After(thirtyDaysAgo) {
			activeUsers[p.UserID] = true
		}
	}

	return int64(len(activeUsers)), nil
}

func (r *progressRepositoryMock) CalculateAverageProgress(ctx context.Context) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.progress) == 0 {
		return 0.0, nil
	}

	var total float64
	for _, p := range r.progress {
		total += float64(p.Percentage)
	}

	return total / float64(len(r.progress)), nil
}
