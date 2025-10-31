package entity

import (
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/pkg/errors"
	"github.com/EduGoGroup/edugo-shared/pkg/types/enum"
)

// Material representa un material educativo (PDF)
type Material struct {
	id               valueobject.MaterialID
	title            string
	description      string
	authorID         valueobject.UserID
	subjectID        string // Puede venir de API Admin
	s3Key            string
	s3URL            string
	status           enum.MaterialStatus
	processingStatus enum.ProcessingStatus
	createdAt        time.Time
	updatedAt        time.Time
}

// NewMaterial crea un nuevo material
func NewMaterial(
	title string,
	description string,
	authorID valueobject.UserID,
	subjectID string,
) (*Material, error) {
	// Validaciones de negocio
	if title == "" {
		return nil, errors.NewValidationError("title is required")
	}

	if authorID.IsZero() {
		return nil, errors.NewValidationError("author_id is required")
	}

	now := time.Now()

	return &Material{
		id:               valueobject.NewMaterialID(),
		title:            title,
		description:      description,
		authorID:         authorID,
		subjectID:        subjectID,
		status:           enum.MaterialStatusDraft,
		processingStatus: enum.ProcessingStatusPending,
		createdAt:        now,
		updatedAt:        now,
	}, nil
}

// ReconstructMaterial reconstruye desde DB
func ReconstructMaterial(
	id valueobject.MaterialID,
	title, description string,
	authorID valueobject.UserID,
	subjectID, s3Key, s3URL string,
	status enum.MaterialStatus,
	processingStatus enum.ProcessingStatus,
	createdAt, updatedAt time.Time,
) *Material {
	return &Material{
		id:               id,
		title:            title,
		description:      description,
		authorID:         authorID,
		subjectID:        subjectID,
		s3Key:            s3Key,
		s3URL:            s3URL,
		status:           status,
		processingStatus: processingStatus,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
}

// Getters

func (m *Material) ID() valueobject.MaterialID       { return m.id }
func (m *Material) Title() string                    { return m.title }
func (m *Material) Description() string              { return m.description }
func (m *Material) AuthorID() valueobject.UserID     { return m.authorID }
func (m *Material) SubjectID() string                { return m.subjectID }
func (m *Material) S3Key() string                    { return m.s3Key }
func (m *Material) S3URL() string                    { return m.s3URL }
func (m *Material) Status() enum.MaterialStatus      { return m.status }
func (m *Material) ProcessingStatus() enum.ProcessingStatus { return m.processingStatus }
func (m *Material) CreatedAt() time.Time             { return m.createdAt }
func (m *Material) UpdatedAt() time.Time             { return m.updatedAt }

// Business Logic

// SetS3Info establece la información de S3 cuando se sube el archivo
func (m *Material) SetS3Info(s3Key, s3URL string) error {
	if s3Key == "" || s3URL == "" {
		return errors.NewValidationError("s3_key and s3_url are required")
	}

	m.s3Key = s3Key
	m.s3URL = s3URL
	m.processingStatus = enum.ProcessingStatusProcessing
	m.updatedAt = time.Now()
	return nil
}

// MarkProcessingComplete marca el procesamiento como completado
func (m *Material) MarkProcessingComplete() error {
	if m.processingStatus == enum.ProcessingStatusCompleted {
		return errors.NewBusinessRuleError("material already processed")
	}

	m.processingStatus = enum.ProcessingStatusCompleted
	m.updatedAt = time.Now()
	return nil
}

// Publish publica el material
func (m *Material) Publish() error {
	if m.status == enum.MaterialStatusPublished {
		return errors.NewBusinessRuleError("material is already published")
	}

	if m.processingStatus != enum.ProcessingStatusCompleted {
		return errors.NewBusinessRuleError("material must be processed before publishing")
	}

	m.status = enum.MaterialStatusPublished
	m.updatedAt = time.Now()
	return nil
}

// Archive archiva el material
func (m *Material) Archive() error {
	if m.status == enum.MaterialStatusArchived {
		return errors.NewBusinessRuleError("material is already archived")
	}

	m.status = enum.MaterialStatusArchived
	m.updatedAt = time.Now()
	return nil
}

// IsDraft verifica si está en draft
func (m *Material) IsDraft() bool {
	return m.status == enum.MaterialStatusDraft
}

// IsPublished verifica si está publicado
func (m *Material) IsPublished() bool {
	return m.status == enum.MaterialStatusPublished
}

// IsProcessed verifica si fue procesado
func (m *Material) IsProcessed() bool {
	return m.processingStatus == enum.ProcessingStatusCompleted
}
