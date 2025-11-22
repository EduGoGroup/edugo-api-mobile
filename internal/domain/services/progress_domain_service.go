package services

import (
	"time"

	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
)

// ProgressDomainService contiene reglas de negocio de Progress
// Extrae la lógica que antes estaba embebida en entity.Progress
type ProgressDomainService struct{}

// NewProgressDomainService crea una nueva instancia del servicio
func NewProgressDomainService() *ProgressDomainService {
	return &ProgressDomainService{}
}

// UpdateProgress actualiza el progreso con validaciones
// Reglas de negocio:
// - Percentage debe estar entre 0-100
// - Status se determina automáticamente según el percentage
//   - 0% = not_started
//   - 1-99% = in_progress
//   - 100% = completed
func (s *ProgressDomainService) UpdateProgress(progress *pgentities.Progress, percentage, lastPage int) error {
	// Validar rango de percentage
	if percentage < 0 || percentage > 100 {
		return errors.NewValidationError("percentage must be between 0 and 100")
	}

	// Actualizar campos
	progress.Percentage = percentage
	progress.LastPage = lastPage
	progress.LastAccessedAt = time.Now()
	progress.UpdatedAt = time.Now()

	// Determinar status según percentage (regla de negocio)
	// Valores según infrastructure: not_started, in_progress, completed
	if percentage == 0 {
		progress.Status = "not_started"
	} else if percentage >= 100 {
		progress.Status = "completed"
	} else {
		progress.Status = "in_progress"
	}

	return nil
}

// IsCompleted indica si el progreso está completado
func (s *ProgressDomainService) IsCompleted(progress *pgentities.Progress) bool {
	return progress.Status == "completed" || progress.Percentage >= 100
}

// IsStarted indica si el progreso fue iniciado
func (s *ProgressDomainService) IsStarted(progress *pgentities.Progress) bool {
	return progress.Status != "not_started" && progress.Percentage > 0
}

// IsInProgress indica si el progreso está en curso
func (s *ProgressDomainService) IsInProgress(progress *pgentities.Progress) bool {
	return progress.Status == "in_progress"
}
