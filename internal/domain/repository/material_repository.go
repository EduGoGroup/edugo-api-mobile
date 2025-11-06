package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// MaterialRepository define las operaciones de persistencia para Material (PostgreSQL)
type MaterialRepository interface {
	// Create crea un nuevo material
	Create(ctx context.Context, material *entity.Material) error

	// FindByID busca un material por ID
	FindByID(ctx context.Context, id valueobject.MaterialID) (*entity.Material, error)

	// FindByIDWithVersions busca un material por ID incluyendo su historial de versiones
	FindByIDWithVersions(ctx context.Context, id valueobject.MaterialID) (*entity.Material, []*entity.MaterialVersion, error)

	// Update actualiza un material
	Update(ctx context.Context, material *entity.Material) error

	// List lista materiales con filtros
	List(ctx context.Context, filters ListFilters) ([]*entity.Material, error)

	// FindByAuthor busca materiales de un autor
	FindByAuthor(ctx context.Context, authorID valueobject.UserID) ([]*entity.Material, error)

	// UpdateStatus actualiza el status del material
	UpdateStatus(ctx context.Context, id valueobject.MaterialID, status enum.MaterialStatus) error

	// UpdateProcessingStatus actualiza el estado de procesamiento
	UpdateProcessingStatus(ctx context.Context, id valueobject.MaterialID, status enum.ProcessingStatus) error

	// CountPublishedMaterials cuenta total de materiales publicados (para estadísticas)
	CountPublishedMaterials(ctx context.Context) (int64, error)
}

// ListFilters filtros para listar materiales
type ListFilters struct {
	Status    *enum.MaterialStatus
	AuthorID  *valueobject.UserID
	SubjectID *string
	Limit     int
	Offset    int
}
