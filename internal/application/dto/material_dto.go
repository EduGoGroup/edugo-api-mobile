package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entity"
	"github.com/EduGoGroup/edugo-shared/common/validator"
)

// CreateMaterialRequest solicitud para crear material
type CreateMaterialRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=200" example:"Introduction to Calculus"`
	Description string `json:"description" binding:"max=1000" example:"A comprehensive guide to differential and integral calculus"`
	SubjectID   string `json:"subject_id" binding:"omitempty,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
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
	ID               string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title            string    `json:"title" example:"Introduction to Calculus"`
	Description      string    `json:"description" example:"A comprehensive guide to differential and integral calculus"`
	AuthorID         string    `json:"author_id" example:"660e8400-e29b-41d4-a716-446655440001"`
	SubjectID        string    `json:"subject_id,omitempty" example:"770e8400-e29b-41d4-a716-446655440002"`
	S3Key            string    `json:"s3_key,omitempty" example:"materials/550e8400-e29b-41d4-a716-446655440000/calculus.pdf"`
	S3URL            string    `json:"s3_url,omitempty" example:"https://s3.amazonaws.com/bucket/materials/550e8400-e29b-41d4-a716-446655440000/calculus.pdf"`
	Status           string    `json:"status" example:"published"`
	ProcessingStatus string    `json:"processing_status" example:"completed"`
	CreatedAt        time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt        time.Time `json:"updated_at" example:"2024-01-15T10:30:00Z"`
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
	S3Key string `json:"s3_key" example:"materials/550e8400-e29b-41d4-a716-446655440000/calculus.pdf"`
	S3URL string `json:"s3_url" example:"https://s3.amazonaws.com/bucket/materials/550e8400-e29b-41d4-a716-446655440000/calculus.pdf"`
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
	FileName    string `json:"file_name" binding:"required" example:"calculus.pdf"`
	ContentType string `json:"content_type" binding:"required" example:"application/pdf"`
}

// GenerateUploadURLResponse respuesta con URL presignada de subida
type GenerateUploadURLResponse struct {
	UploadURL string `json:"upload_url" example:"https://s3.amazonaws.com/bucket/materials/550e8400-e29b-41d4-a716-446655440000/calculus.pdf?X-Amz-Algorithm=..."`
	S3Key     string `json:"s3_key" example:"materials/550e8400-e29b-41d4-a716-446655440000/calculus.pdf"`
	ExpiresIn int    `json:"expires_in" example:"900"` // En segundos
}

// GenerateDownloadURLResponse respuesta con URL presignada de descarga
type GenerateDownloadURLResponse struct {
	DownloadURL string `json:"download_url" example:"https://s3.amazonaws.com/bucket/materials/550e8400-e29b-41d4-a716-446655440000/calculus.pdf?X-Amz-Algorithm=..."`
	ExpiresIn   int    `json:"expires_in" example:"3600"` // En segundos
}

// MaterialVersionResponse representa una versión de material
type MaterialVersionResponse struct {
	ID            string    `json:"id" example:"880e8400-e29b-41d4-a716-446655440003"`
	VersionNumber int       `json:"version_number" example:"2"`
	Title         string    `json:"title" example:"Introduction to Calculus - Updated"`
	ContentURL    string    `json:"content_url" example:"https://s3.amazonaws.com/bucket/materials/550e8400-e29b-41d4-a716-446655440000/calculus-v2.pdf"`
	ChangedBy     string    `json:"changed_by" example:"660e8400-e29b-41d4-a716-446655440001"`
	CreatedAt     time.Time `json:"created_at" example:"2024-01-20T14:30:00Z"`
}

// MaterialWithVersionsResponse respuesta de material con su historial de versiones
type MaterialWithVersionsResponse struct {
	Material *MaterialResponse          `json:"material"`
	Versions []*MaterialVersionResponse `json:"versions"`
}

// ToMaterialVersionResponse convierte entidad a DTO
func ToMaterialVersionResponse(version *entity.MaterialVersion) *MaterialVersionResponse {
	return &MaterialVersionResponse{
		ID:            version.ID().String(),
		VersionNumber: version.VersionNumber(),
		Title:         version.Title(),
		ContentURL:    version.ContentURL(),
		ChangedBy:     version.ChangedBy().String(),
		CreatedAt:     version.CreatedAt(),
	}
}

// ToMaterialWithVersionsResponse convierte material y versiones a DTO
func ToMaterialWithVersionsResponse(material *entity.Material, versions []*entity.MaterialVersion) *MaterialWithVersionsResponse {
	versionResponses := make([]*MaterialVersionResponse, len(versions))
	for i, version := range versions {
		versionResponses[i] = ToMaterialVersionResponse(version)
	}

	return &MaterialWithVersionsResponse{
		Material: ToMaterialResponse(material),
		Versions: versionResponses,
	}
}
