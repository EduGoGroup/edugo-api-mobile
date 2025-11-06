package entity

import (
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// MaterialVersion representa una versión histórica de un material educativo
// Cada vez que se realiza un cambio significativo al material, se crea una nueva versión
type MaterialVersion struct {
	id            valueobject.MaterialVersionID
	materialID    valueobject.MaterialID
	versionNumber int
	title         string
	contentURL    string
	changedBy     valueobject.UserID
	createdAt     time.Time
}

// NewMaterialVersion crea una nueva versión de un material
func NewMaterialVersion(
	materialID valueobject.MaterialID,
	versionNumber int,
	title string,
	contentURL string,
	changedBy valueobject.UserID,
) (*MaterialVersion, error) {
	// Validaciones de negocio
	if materialID.IsZero() {
		return nil, errors.NewValidationError("material_id is required")
	}

	if versionNumber <= 0 {
		return nil, errors.NewValidationError("version_number must be positive")
	}

	if title == "" {
		return nil, errors.NewValidationError("title is required")
	}

	if contentURL == "" {
		return nil, errors.NewValidationError("content_url is required")
	}

	if changedBy.IsZero() {
		return nil, errors.NewValidationError("changed_by is required")
	}

	return &MaterialVersion{
		id:            valueobject.NewMaterialVersionID(),
		materialID:    materialID,
		versionNumber: versionNumber,
		title:         title,
		contentURL:    contentURL,
		changedBy:     changedBy,
		createdAt:     time.Now(),
	}, nil
}

// ReconstructMaterialVersion reconstruye una versión desde la base de datos
func ReconstructMaterialVersion(
	id valueobject.MaterialVersionID,
	materialID valueobject.MaterialID,
	versionNumber int,
	title, contentURL string,
	changedBy valueobject.UserID,
	createdAt time.Time,
) *MaterialVersion {
	return &MaterialVersion{
		id:            id,
		materialID:    materialID,
		versionNumber: versionNumber,
		title:         title,
		contentURL:    contentURL,
		changedBy:     changedBy,
		createdAt:     createdAt,
	}
}

// Getters

func (mv *MaterialVersion) ID() valueobject.MaterialVersionID  { return mv.id }
func (mv *MaterialVersion) MaterialID() valueobject.MaterialID { return mv.materialID }
func (mv *MaterialVersion) VersionNumber() int                 { return mv.versionNumber }
func (mv *MaterialVersion) Title() string                      { return mv.title }
func (mv *MaterialVersion) ContentURL() string                 { return mv.contentURL }
func (mv *MaterialVersion) ChangedBy() valueobject.UserID      { return mv.changedBy }
func (mv *MaterialVersion) CreatedAt() time.Time               { return mv.createdAt }
