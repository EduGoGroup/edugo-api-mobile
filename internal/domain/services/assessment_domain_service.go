package services

import (
	"time"

	"github.com/google/uuid"

	pgentities "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure_stubs/postgres/entities"
	domainErrors "github.com/EduGoGroup/edugo-api-mobile/internal/domain/errors"
)

// AssessmentDomainService contiene reglas de negocio de Assessment
// Extrae la lógica que antes estaba embebida en entities.Assessment
type AssessmentDomainService struct{}

// NewAssessmentDomainService crea una nueva instancia del servicio
func NewAssessmentDomainService() *AssessmentDomainService {
	return &AssessmentDomainService{}
}

// ValidateAssessment valida una entity Assessment
// Aplica todas las reglas de validación del dominio
func (s *AssessmentDomainService) ValidateAssessment(a *pgentities.Assessment) error {
	// Validar IDs
	if a.ID == uuid.Nil {
		return domainErrors.ErrInvalidAssessmentID
	}

	if a.MaterialID == uuid.Nil {
		return domainErrors.ErrInvalidMaterialID
	}

	// Validar MongoDB document ID (24 caracteres hexadecimales)
	if len(a.MongoDocumentID) != 24 {
		return domainErrors.ErrInvalidMongoDocumentID
	}

	// Validar título
	if a.Title == "" {
		return domainErrors.ErrEmptyTitle
	}

	// Validar total de preguntas (1-100 según schema PostgreSQL)
	if a.TotalQuestions < 1 || a.TotalQuestions > 100 {
		return domainErrors.ErrInvalidTotalQuestions
	}

	// Validar pass threshold (0-100)
	if a.PassThreshold < 0 || a.PassThreshold > 100 {
		return domainErrors.ErrInvalidPassThreshold
	}

	// Validar max attempts (si está definido, debe ser >= 1)
	if a.MaxAttempts != nil && *a.MaxAttempts < 1 {
		return domainErrors.ErrInvalidMaxAttempts
	}

	// Validar time limit (si está definido, debe estar entre 1-180 minutos)
	if a.TimeLimitMinutes != nil && (*a.TimeLimitMinutes < 1 || *a.TimeLimitMinutes > 180) {
		return domainErrors.ErrInvalidTimeLimit
	}

	return nil
}

// CanAttempt verifica si un usuario puede hacer otro intento
// Regla de negocio: si MaxAttempts es nil, intentos ilimitados
func (s *AssessmentDomainService) CanAttempt(assessment *pgentities.Assessment, attemptCount int) bool {
	if assessment.MaxAttempts == nil {
		return true // Intentos ilimitados
	}

	return attemptCount < *assessment.MaxAttempts
}

// IsTimeLimited indica si la evaluación tiene límite de tiempo
func (s *AssessmentDomainService) IsTimeLimited(assessment *pgentities.Assessment) bool {
	return assessment.TimeLimitMinutes != nil && *assessment.TimeLimitMinutes > 0
}

// SetMaxAttempts establece máximo de intentos
// Valida que sea >= 1
func (s *AssessmentDomainService) SetMaxAttempts(assessment *pgentities.Assessment, max int) error {
	if max < 1 {
		return domainErrors.ErrInvalidMaxAttempts
	}

	assessment.MaxAttempts = &max
	assessment.UpdatedAt = time.Now().UTC()

	return nil
}

// SetTimeLimit establece límite de tiempo en minutos
// Valida que esté entre 1-180 minutos
func (s *AssessmentDomainService) SetTimeLimit(assessment *pgentities.Assessment, minutes int) error {
	if minutes < 1 || minutes > 180 {
		return domainErrors.ErrInvalidTimeLimit
	}

	assessment.TimeLimitMinutes = &minutes
	assessment.UpdatedAt = time.Now().UTC()

	return nil
}

// RemoveMaxAttempts quita el límite de intentos (pone en nil = ilimitado)
func (s *AssessmentDomainService) RemoveMaxAttempts(assessment *pgentities.Assessment) {
	assessment.MaxAttempts = nil
	assessment.UpdatedAt = time.Now().UTC()
}

// RemoveTimeLimit quita el límite de tiempo (pone en nil = sin límite)
func (s *AssessmentDomainService) RemoveTimeLimit(assessment *pgentities.Assessment) {
	assessment.TimeLimitMinutes = nil
	assessment.UpdatedAt = time.Now().UTC()
}
