package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// MaterialReader define operaciones de lectura para Material
// Principio ISP: Separar lectura de escritura y estadísticas
type MaterialReader interface {
	// FindByID busca un material por ID
	FindByID(ctx context.Context, id valueobject.MaterialID) (*entity.Material, error)

	// FindByIDWithVersions busca un material por ID incluyendo su historial de versiones
	FindByIDWithVersions(ctx context.Context, id valueobject.MaterialID) (*entity.Material, []*entity.MaterialVersion, error)

	// List lista materiales con filtros
	List(ctx context.Context, filters ListFilters) ([]*entity.Material, error)

	// FindByAuthor busca materiales de un autor
	FindByAuthor(ctx context.Context, authorID valueobject.UserID) ([]*entity.Material, error)
}

// MaterialWriter define operaciones de escritura para Material
type MaterialWriter interface {
	// Create crea un nuevo material
	Create(ctx context.Context, material *entity.Material) error

	// Update actualiza un material
	Update(ctx context.Context, material *entity.Material) error

	// UpdateStatus actualiza el status del material
	UpdateStatus(ctx context.Context, id valueobject.MaterialID, status enum.MaterialStatus) error

	// UpdateProcessingStatus actualiza el estado de procesamiento
	UpdateProcessingStatus(ctx context.Context, id valueobject.MaterialID, status enum.ProcessingStatus) error
}

// MaterialStats define operaciones de estadísticas para Material
type MaterialStats interface {
	// CountPublishedMaterials cuenta total de materiales publicados (para estadísticas)
	CountPublishedMaterials(ctx context.Context) (int64, error)
}

// MaterialRepository agrega todas las capacidades de Material (PostgreSQL)
// Las implementaciones deben cumplir con todas las interfaces segregadas
type MaterialRepository interface {
	MaterialReader
	MaterialWriter
	MaterialStats
}

// ListFilters filtros para listar materiales
type ListFilters struct {
	Status    *enum.MaterialStatus
	AuthorID  *valueobject.UserID
	SubjectID *string
	Limit     int
	Offset    int
}
