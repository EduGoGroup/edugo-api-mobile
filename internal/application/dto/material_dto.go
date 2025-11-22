package dto

import (
	"time"

	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/validator"
)

// CreateMaterialRequest solicitud para crear material
type CreateMaterialRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=200" example:"Introduction to Calculus"`
	Description string `json:"description" binding:"max=1000" example:"A comprehensive guide to differential and integral calculus"`
	Subject     string `json:"subject" binding:"omitempty" example:"Mathematics"`
	Grade       string `json:"grade" binding:"omitempty" example:"12th Grade"`
}

func (r *CreateMaterialRequest) Validate() error {
	v := validator.New()

	v.Required(r.Title, "title")
	v.MinLength(r.Title, 3, "title")
	v.MaxLength(r.Title, 200, "title")

	v.MaxLength(r.Description, 1000, "description")

	return v.GetError()
}

// MaterialResponse respuesta de material
// Adaptado a estructura REAL de infrastructure
type MaterialResponse struct {
	ID                    string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	SchoolID              string     `json:"school_id" example:"660e8400-e29b-41d4-a716-446655440001"`
	UploadedByTeacherID   string     `json:"uploaded_by_teacher_id" example:"770e8400-e29b-41d4-a716-446655440002"`
	AcademicUnitID        *string    `json:"academic_unit_id,omitempty" example:"880e8400-e29b-41d4-a716-446655440003"`
	Title                 string     `json:"title" example:"Introduction to Calculus"`
	Description           *string    `json:"description,omitempty" example:"A comprehensive guide to differential and integral calculus"`
	Subject               *string    `json:"subject,omitempty" example:"Mathematics"`
	Grade                 *string    `json:"grade,omitempty" example:"12th Grade"`
	FileURL               string     `json:"file_url" example:"https://s3.amazonaws.com/bucket/materials/550e8400.pdf"`
	FileType              string     `json:"file_type" example:"application/pdf"`
	FileSizeBytes         int64      `json:"file_size_bytes" example:"1048576"`
	Status                string     `json:"status" example:"uploaded"` // uploaded, processing, ready, failed
	ProcessingStartedAt   *time.Time `json:"processing_started_at,omitempty" example:"2024-01-15T10:30:00Z"`
	ProcessingCompletedAt *time.Time `json:"processing_completed_at,omitempty" example:"2024-01-15T10:35:00Z"`
	IsPublic              bool       `json:"is_public" example:"false"`
	CreatedAt             time.Time  `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt             time.Time  `json:"updated_at" example:"2024-01-15T10:30:00Z"`
	DeletedAt             *time.Time `json:"deleted_at,omitempty"`
}

// ToMaterialResponse convierte entity infrastructure a DTO
func ToMaterialResponse(material *pgentities.Material) *MaterialResponse {
	resp := &MaterialResponse{
		ID:                  material.ID.String(),
		SchoolID:            material.SchoolID.String(),
		UploadedByTeacherID: material.UploadedByTeacherID.String(),
		Title:               material.Title,
		FileURL:             material.FileURL,
		FileType:            material.FileType,
		FileSizeBytes:       material.FileSizeBytes,
		Status:              material.Status,
		IsPublic:            material.IsPublic,
		CreatedAt:           material.CreatedAt,
		UpdatedAt:           material.UpdatedAt,
	}

	// Campos nullable
	if material.AcademicUnitID != nil {
		academicUnitStr := material.AcademicUnitID.String()
		resp.AcademicUnitID = &academicUnitStr
	}

	resp.Description = material.Description
	resp.Subject = material.Subject
	resp.Grade = material.Grade
	resp.ProcessingStartedAt = material.ProcessingStartedAt
	resp.ProcessingCompletedAt = material.ProcessingCompletedAt
	resp.DeletedAt = material.DeletedAt

	return resp
}

// UploadCompleteRequest notificación de subida completa
type UploadCompleteRequest struct {
	FileURL       string `json:"file_url" example:"https://s3.amazonaws.com/bucket/materials/file.pdf"`
	FileType      string `json:"file_type" example:"application/pdf"`
	FileSizeBytes int64  `json:"file_size_bytes" example:"1048576"`
}

func (r *UploadCompleteRequest) Validate() error {
	v := validator.New()
	v.Required(r.FileURL, "file_url")
	v.URL(r.FileURL, "file_url")
	v.Required(r.FileType, "file_type")
	// FileSizeBytes validación simple
	return v.GetError()
}

// GenerateUploadURLRequest solicitud para generar URL de subida presignada
type GenerateUploadURLRequest struct {
	FileName    string `json:"file_name" binding:"required" example:"calculus.pdf"`
	ContentType string `json:"content_type" binding:"required" example:"application/pdf"`
}

// GenerateUploadURLResponse respuesta con URL presignada de subida
type GenerateUploadURLResponse struct {
	UploadURL string `json:"upload_url" example:"https://s3.amazonaws.com/bucket/materials/550e8400.pdf?X-Amz-Algorithm=..."`
	FileURL   string `json:"file_url" example:"https://s3.amazonaws.com/bucket/materials/550e8400.pdf"`
	ExpiresIn int    `json:"expires_in" example:"900"` // En segundos
}

// GenerateDownloadURLResponse respuesta con URL presignada de descarga
type GenerateDownloadURLResponse struct {
	DownloadURL string `json:"download_url" example:"https://s3.amazonaws.com/bucket/materials/550e8400.pdf?X-Amz-Algorithm=..."`
	ExpiresIn   int    `json:"expires_in" example:"3600"` // En segundos
}

// MaterialVersionResponse representa una versión de material
// Adaptado a estructura REAL de infrastructure
type MaterialVersionResponse struct {
	ID            string    `json:"id" example:"880e8400-e29b-41d4-a716-446655440003"`
	MaterialID    string    `json:"material_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	VersionNumber int       `json:"version_number" example:"2"`
	Title         string    `json:"title" example:"Introduction to Calculus - Updated"`
	ContentURL    string    `json:"content_url" example:"https://s3.amazonaws.com/bucket/materials/file-v2.pdf"`
	ChangedBy     string    `json:"changed_by" example:"660e8400-e29b-41d4-a716-446655440001"`
	CreatedAt     time.Time `json:"created_at" example:"2024-01-20T14:30:00Z"`
}

// MaterialWithVersionsResponse respuesta de material con su historial de versiones
type MaterialWithVersionsResponse struct {
	Material *MaterialResponse          `json:"material"`
	Versions []*MaterialVersionResponse `json:"versions"`
}

// ToMaterialVersionResponse convierte entidad a DTO
func ToMaterialVersionResponse(version *pgentities.MaterialVersion) *MaterialVersionResponse {
	return &MaterialVersionResponse{
		ID:            version.ID.String(),
		MaterialID:    version.MaterialID.String(),
		VersionNumber: version.VersionNumber,
		Title:         version.Title,
		ContentURL:    version.ContentURL,
		ChangedBy:     version.ChangedBy.String(),
		CreatedAt:     version.CreatedAt,
	}
}

// ToMaterialWithVersionsResponse convierte material y versiones a DTO
func ToMaterialWithVersionsResponse(material *pgentities.Material, versions []*pgentities.MaterialVersion) *MaterialWithVersionsResponse {
	versionResponses := make([]*MaterialVersionResponse, len(versions))
	for i, version := range versions {
		versionResponses[i] = ToMaterialVersionResponse(version)
	}

	return &MaterialWithVersionsResponse{
		Material: ToMaterialResponse(material),
		Versions: versionResponses,
	}
}
