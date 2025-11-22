package services

import (
	"time"

	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// MaterialDomainService contiene reglas de negocio de Material
// Extrae la lógica que antes estaba embebida en entity.Material
type MaterialDomainService struct{}

// NewMaterialDomainService crea una nueva instancia del servicio
func NewMaterialDomainService() *MaterialDomainService {
	return &MaterialDomainService{}
}

// SetFileInfo establece información del archivo subido
// Adaptado a la estructura REAL de infrastructure (FileURL, FileType, FileSizeBytes)
// Status: uploaded → processing
func (s *MaterialDomainService) SetFileInfo(material *pgentities.Material, fileURL string, fileType string, fileSizeBytes int64) error {
	if fileURL == "" {
		return errors.NewValidationError("file_url is required")
	}
	if fileType == "" {
		return errors.NewValidationError("file_type is required")
	}

	material.FileURL = fileURL
	material.FileType = fileType
	material.FileSizeBytes = fileSizeBytes
	material.Status = "processing" // Valores según migration: uploaded, processing, ready, failed
	now := time.Now()
	material.ProcessingStartedAt = &now
	material.UpdatedAt = now

	return nil
}

// MarkProcessingComplete marca procesamiento completado
// Status: processing → ready
func (s *MaterialDomainService) MarkProcessingComplete(material *pgentities.Material) error {
	if material.Status == "ready" {
		return errors.NewBusinessRuleError("material already processed")
	}

	material.Status = "ready" // Material listo para usar
	now := time.Now()
	material.ProcessingCompletedAt = &now
	material.UpdatedAt = now

	return nil
}

// Publish publica el material (lo hace público)
// Regla de negocio: un material debe estar procesado (ready) antes de publicarse
func (s *MaterialDomainService) Publish(material *pgentities.Material) error {
	if material.IsPublic {
		return errors.NewBusinessRuleError("material is already published")
	}

	if material.Status != "ready" {
		return errors.NewBusinessRuleError("material must be processed before publishing")
	}

	material.IsPublic = true
	material.UpdatedAt = time.Now()

	return nil
}

// Archive archiva el material (soft delete)
// Usa DeletedAt para marcar como archivado
func (s *MaterialDomainService) Archive(material *pgentities.Material) error {
	if material.DeletedAt != nil {
		return errors.NewBusinessRuleError("material is already archived")
	}

	now := time.Now()
	material.DeletedAt = &now
	material.UpdatedAt = now

	return nil
}

// Query helpers - Métodos de consulta que no modifican el material

// IsUploaded indica si el material fue subido (estado inicial)
func (s *MaterialDomainService) IsUploaded(material *pgentities.Material) bool {
	return material.Status == "uploaded"
}

// IsProcessing indica si el material está siendo procesado
func (s *MaterialDomainService) IsProcessing(material *pgentities.Material) bool {
	return material.Status == "processing"
}

// IsReady indica si el material está listo (procesado)
func (s *MaterialDomainService) IsReady(material *pgentities.Material) bool {
	return material.Status == "ready"
}

// IsFailed indica si el procesamiento falló
func (s *MaterialDomainService) IsFailed(material *pgentities.Material) bool {
	return material.Status == "failed"
}

// IsPublished indica si el material está publicado (público)
func (s *MaterialDomainService) IsPublished(material *pgentities.Material) bool {
	return material.IsPublic
}

// IsArchived indica si el material está archivado (soft deleted)
func (s *MaterialDomainService) IsArchived(material *pgentities.Material) bool {
	return material.DeletedAt != nil
}
