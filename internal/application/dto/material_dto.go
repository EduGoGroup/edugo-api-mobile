package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-shared/common/validator"
)

// CreateMaterialRequest solicitud para crear material
type CreateMaterialRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=200"`
	Description string `json:"description" binding:"max=1000"`
	SubjectID   string `json:"subject_id" binding:"omitempty,uuid"`
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
	S3Key            string    `json:"s3_key,omitempty"`
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
		S3Key:            material.S3Key(),
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

// GenerateUploadURLRequest solicitud para generar URL de subida presignada
type GenerateUploadURLRequest struct {
	FileName    string `json:"file_name" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

// GenerateUploadURLResponse respuesta con URL presignada de subida
type GenerateUploadURLResponse struct {
	UploadURL string `json:"upload_url"`
	S3Key     string `json:"s3_key"`
	ExpiresIn int    `json:"expires_in"` // En segundos
}

// GenerateDownloadURLResponse respuesta con URL presignada de descarga
type GenerateDownloadURLResponse struct {
	DownloadURL string `json:"download_url"`
	ExpiresIn   int    `json:"expires_in"` // En segundos
}
