package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
)

// MaterialSummary representa el resumen de un material (almacenado en MongoDB)
type MaterialSummary struct {
	MaterialID  valueobject.MaterialID
	MainIdeas   []string
	KeyConcepts map[string]string
	Sections    []SummarySection
	Glossary    map[string]string
	CreatedAt   string
}

// SummarySection representa una sección del resumen
type SummarySection struct {
	Title   string
	Content string
	Page    int
}

// SummaryReader define operaciones de lectura para summaries (MongoDB)
// Principio ISP: Separar lectura de escritura
type SummaryReader interface {
	// FindByMaterialID busca el summary de un material
	FindByMaterialID(ctx context.Context, materialID valueobject.MaterialID) (*MaterialSummary, error)

	// Exists verifica si existe un summary
	Exists(ctx context.Context, materialID valueobject.MaterialID) (bool, error)
}

// SummaryWriter define operaciones de escritura para summaries
type SummaryWriter interface {
	// Save guarda o actualiza un summary
	Save(ctx context.Context, summary *MaterialSummary) error

	// Delete elimina un summary
	Delete(ctx context.Context, materialID valueobject.MaterialID) error
}

// SummaryRepository agrega todas las capacidades de summaries (MongoDB)
// Las implementaciones deben cumplir con todas las interfaces segregadas
type SummaryRepository interface {
	SummaryReader
	SummaryWriter
}
