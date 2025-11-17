package entities

import (
	"time"

	"github.com/google/uuid"

	domainErrors "github.com/EduGoGroup/edugo-api-mobile/internal/domain/errors"
)

// Assessment representa una evaluación de un material educativo
// Esta entity corresponde a la tabla `assessment` en PostgreSQL
type Assessment struct {
	ID               uuid.UUID
	MaterialID       uuid.UUID
	MongoDocumentID  string // ObjectId de MongoDB (24 caracteres hex)
	Title            string
	TotalQuestions   int
	PassThreshold    int  // Porcentaje 0-100 para aprobar
	MaxAttempts      *int // nil = intentos ilimitados
	TimeLimitMinutes *int // nil = sin límite de tiempo
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewAssessment crea una nueva evaluación con validaciones
// Este constructor aplica fail-fast: si alguna validación falla, retorna error inmediatamente
func NewAssessment(
	materialID uuid.UUID,
	mongoDocID string,
	title string,
	totalQuestions int,
	passThreshold int,
) (*Assessment, error) {
	// Validar material ID
	if materialID == uuid.Nil {
		return nil, domainErrors.ErrInvalidMaterialID
	}

	// Validar MongoDB document ID (24 caracteres hexadecimales)
	if len(mongoDocID) != 24 {
		return nil, domainErrors.ErrInvalidMongoDocumentID
	}

	// Validar que el título no esté vacío
	if title == "" {
		return nil, domainErrors.ErrEmptyTitle
	}

	// Validar total de preguntas (1-100 según schema PostgreSQL)
	if totalQuestions < 1 || totalQuestions > 100 {
		return nil, domainErrors.ErrInvalidTotalQuestions
	}

	// Validar umbral de aprobación (0-100)
	if passThreshold < 0 || passThreshold > 100 {
		return nil, domainErrors.ErrInvalidPassThreshold
	}

	now := time.Now().UTC()
	return &Assessment{
		ID:               uuid.New(),
		MaterialID:       materialID,
		MongoDocumentID:  mongoDocID,
		Title:            title,
		TotalQuestions:   totalQuestions,
		PassThreshold:    passThreshold,
		MaxAttempts:      nil, // Default: ilimitado
		TimeLimitMinutes: nil, // Default: sin límite
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// Validate verifica que la evaluación sea válida en su estado actual
func (a *Assessment) Validate() error {
	if a.ID == uuid.Nil {
		return domainErrors.ErrInvalidAssessmentID
	}
	if a.MaterialID == uuid.Nil {
		return domainErrors.ErrInvalidMaterialID
	}
	if len(a.MongoDocumentID) != 24 {
		return domainErrors.ErrInvalidMongoDocumentID
	}
	if a.Title == "" {
		return domainErrors.ErrEmptyTitle
	}
	if a.TotalQuestions < 1 || a.TotalQuestions > 100 {
		return domainErrors.ErrInvalidTotalQuestions
	}
	if a.PassThreshold < 0 || a.PassThreshold > 100 {
		return domainErrors.ErrInvalidPassThreshold
	}
	if a.MaxAttempts != nil && *a.MaxAttempts < 1 {
		return domainErrors.ErrInvalidMaxAttempts
	}
	if a.TimeLimitMinutes != nil && (*a.TimeLimitMinutes < 1 || *a.TimeLimitMinutes > 180) {
		return domainErrors.ErrInvalidTimeLimit
	}
	return nil
}

// CanAttempt verifica si un estudiante puede hacer otro intento
// Regla de negocio: si MaxAttempts es nil, intentos ilimitados
func (a *Assessment) CanAttempt(attemptCount int) bool {
	if a.MaxAttempts == nil {
		return true // Ilimitado
	}
	return attemptCount < *a.MaxAttempts
}

// IsTimeLimited indica si la evaluación tiene límite de tiempo
func (a *Assessment) IsTimeLimited() bool {
	return a.TimeLimitMinutes != nil && *a.TimeLimitMinutes > 0
}

// SetMaxAttempts establece el máximo de intentos permitidos
// Esta es una business rule: mínimo 1 intento debe permitirse
func (a *Assessment) SetMaxAttempts(max int) error {
	if max < 1 {
		return domainErrors.ErrInvalidMaxAttempts
	}
	a.MaxAttempts = &max
	a.UpdatedAt = time.Now().UTC()
	return nil
}

// SetTimeLimit establece el límite de tiempo en minutos
// Business rule: entre 1 y 180 minutos (3 horas)
func (a *Assessment) SetTimeLimit(minutes int) error {
	if minutes < 1 || minutes > 180 {
		return domainErrors.ErrInvalidTimeLimit
	}
	a.TimeLimitMinutes = &minutes
	a.UpdatedAt = time.Now().UTC()
	return nil
}

// RemoveMaxAttempts quita el límite de intentos (ilimitados)
func (a *Assessment) RemoveMaxAttempts() {
	a.MaxAttempts = nil
	a.UpdatedAt = time.Now().UTC()
}

// RemoveTimeLimit quita el límite de tiempo
func (a *Assessment) RemoveTimeLimit() {
	a.TimeLimitMinutes = nil
	a.UpdatedAt = time.Now().UTC()
}
