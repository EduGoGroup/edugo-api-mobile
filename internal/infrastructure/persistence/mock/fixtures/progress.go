package fixtures

import (
	"time"

	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

var (
	progressBaseTime = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
)

// ProgressKey es la clave compuesta para buscar progress (MaterialID + UserID)
type ProgressKey struct {
	MaterialID uuid.UUID
	UserID     uuid.UUID
}

// GetDefaultProgress retorna los registros de progreso de prueba
func GetDefaultProgress() map[ProgressKey]*pgentities.Progress {
	progress := make(map[ProgressKey]*pgentities.Progress)

	// Progreso del estudiante en Guía de Sumas (75% completado)
	progress[ProgressKey{MaterialID: MaterialGuidaSumasID, UserID: StudentID}] = &pgentities.Progress{
		MaterialID:     MaterialGuidaSumasID,
		UserID:         StudentID,
		Percentage:     75,
		LastPage:       8,
		Status:         "in_progress",
		LastAccessedAt: progressBaseTime.Add(24 * time.Hour),
		CreatedAt:      progressBaseTime,
		UpdatedAt:      progressBaseTime.Add(24 * time.Hour),
	}

	// Progreso del estudiante en Guía de Restas (30% completado)
	progress[ProgressKey{MaterialID: MaterialGuiaRestasID, UserID: StudentID}] = &pgentities.Progress{
		MaterialID:     MaterialGuiaRestasID,
		UserID:         StudentID,
		Percentage:     30,
		LastPage:       3,
		Status:         "in_progress",
		LastAccessedAt: progressBaseTime.Add(48 * time.Hour),
		CreatedAt:      progressBaseTime.Add(24 * time.Hour),
		UpdatedAt:      progressBaseTime.Add(48 * time.Hour),
	}

	return progress
}
