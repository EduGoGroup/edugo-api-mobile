package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-shared/pkg/validator"
)

// CreateMaterialRequest solicitud para crear material
type CreateMaterialRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	SubjectID   string `json:"subject_id"`
}

func (r *CreateMaterialRequest) Validate() error {
	v := validator.New()

	v.Required(r.Title, "title")
	v.MinLength(r.Title, 3, "title")
	v.MaxLength(r.Title, 200, "title")

	v.MaxLength(r.Description, 1000, "description")

	if r.SubjectID != "" {
		v.UUID(r.SubjectID, "subject_id")
	}

	return v.GetError()
}

// MaterialResponse respuesta de material
type MaterialResponse struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	AuthorID         string    `json:"author_id"`
	SubjectID        string    `json:"subject_id,omitempty"`
	S3URL            string    `json:"s3_url,omitempty"`
	Status           string    `json:"status"`
	ProcessingStatus string    `json:"processing_status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func ToMaterialResponse(material *entity.Material) *MaterialResponse {
	return &MaterialResponse{
		ID:               material.ID().String(),
		Title:            material.Title(),
		Description:      material.Description(),
		AuthorID:         material.AuthorID().String(),
		SubjectID:        material.SubjectID(),
		S3URL:            material.S3URL(),
		Status:           material.Status().String(),
		ProcessingStatus: material.ProcessingStatus().String(),
		CreatedAt:        material.CreatedAt(),
		UpdatedAt:        material.UpdatedAt(),
	}
}

// UploadCompleteRequest notificación de subida completa
type UploadCompleteRequest struct {
	S3Key string `json:"s3_key"`
	S3URL string `json:"s3_url"`
}

func (r *UploadCompleteRequest) Validate() error {
	v := validator.New()
	v.Required(r.S3Key, "s3_key")
	v.Required(r.S3URL, "s3_url")
	v.URL(r.S3URL, "s3_url")
	return v.GetError()
}
