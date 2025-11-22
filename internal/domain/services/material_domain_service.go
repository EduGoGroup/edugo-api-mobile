package services

import (
	"time"

	pgentities "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure_stubs/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// MaterialDomainService contiene reglas de negocio de Material
// Extrae la lógica que antes estaba embebida en entity.Material
type MaterialDomainService struct{}

// NewMaterialDomainService crea una nueva instancia del servicio
func NewMaterialDomainService() *MaterialDomainService {
	return &MaterialDomainService{}
}

// SetS3Info establece información de S3 (extraído de entity.Material)
// Valida que los parámetros no estén vacíos y actualiza el material
func (s *MaterialDomainService) SetS3Info(material *pgentities.Material, s3Key, s3URL string) error {
	if s3Key == "" || s3URL == "" {
		return errors.NewValidationError("s3_key and s3_url are required")
	}

	material.S3Key = s3Key
	material.S3URL = s3URL
	material.ProcessingStatus = enum.ProcessingStatusProcessing
	material.UpdatedAt = time.Now()

	return nil
}

// MarkProcessingComplete marca procesamiento completado
// Valida que no esté ya procesado antes de cambiar el estado
func (s *MaterialDomainService) MarkProcessingComplete(material *pgentities.Material) error {
	if material.ProcessingStatus == enum.ProcessingStatusCompleted {
		return errors.NewBusinessRuleError("material already processed")
	}

	material.ProcessingStatus = enum.ProcessingStatusCompleted
	material.UpdatedAt = time.Now()

	return nil
}

// Publish publica el material
// Regla de negocio: un material debe estar procesado antes de publicarse
func (s *MaterialDomainService) Publish(material *pgentities.Material) error {
	if material.Status == enum.MaterialStatusPublished {
		return errors.NewBusinessRuleError("material is already published")
	}

	if material.ProcessingStatus != enum.ProcessingStatusCompleted {
		return errors.NewBusinessRuleError("material must be processed before publishing")
	}

	material.Status = enum.MaterialStatusPublished
	material.UpdatedAt = time.Now()

	return nil
}

// Archive archiva el material
// Valida que no esté ya archivado
func (s *MaterialDomainService) Archive(material *pgentities.Material) error {
	if material.Status == enum.MaterialStatusArchived {
		return errors.NewBusinessRuleError("material is already archived")
	}

	material.Status = enum.MaterialStatusArchived
	material.UpdatedAt = time.Now()

	return nil
}

// Query helpers - Métodos de consulta que no modifican el material

// IsDraft indica si el material está en estado draft
func (s *MaterialDomainService) IsDraft(material *pgentities.Material) bool {
	return material.Status == enum.MaterialStatusDraft
}

// IsPublished indica si el material está publicado
func (s *MaterialDomainService) IsPublished(material *pgentities.Material) bool {
	return material.Status == enum.MaterialStatusPublished
}

// IsProcessed indica si el material ya fue procesado
func (s *MaterialDomainService) IsProcessed(material *pgentities.Material) bool {
	return material.ProcessingStatus == enum.ProcessingStatusCompleted
}
