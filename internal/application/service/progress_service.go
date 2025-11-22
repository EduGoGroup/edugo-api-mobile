package service

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"go.uber.org/zap"
)

type ProgressService interface {
	UpdateProgress(ctx context.Context, materialID string, userID string, percentage int, lastPage int) error
}

type progressService struct {
	progressRepo repository.ProgressRepository
	logger       logger.Logger
}

func NewProgressService(progressRepo repository.ProgressRepository, logger logger.Logger) ProgressService {
	return &progressService{
		progressRepo: progressRepo,
		logger:       logger,
	}
}

// UpdateProgress actualiza el progreso de un usuario en un material de forma idempotente.
// Usa operación UPSERT para evitar duplicados y simplificar lógica de cliente.
// Si progress=100, se publica evento "material_completed" a RabbitMQ (futuro).
func (s *progressService) UpdateProgress(ctx context.Context, materialID string, userIDStr string, percentage int, lastPage int) error {
	startTime := time.Now()

	// Logging de entrada con contexto
	s.logger.Info("updating progress",
		zap.String("material_id", materialID),
		zap.String("user_id", userIDStr),
		zap.Int("percentage", percentage),
		zap.Int("last_page", lastPage),
	)

	// Validar que percentage está en rango [0-100]
	if percentage < 0 || percentage > 100 {
		s.logger.Warn("invalid percentage value",
			zap.Int("percentage", percentage),
			zap.String("user_id", userIDStr),
		)
		return errors.NewValidationError("percentage must be between 0 and 100")
	}

	// Validar materialID
	matID, err := valueobject.MaterialIDFromString(materialID)
	if err != nil {
		s.logger.Error("invalid material_id", zap.Error(err))
		return errors.NewValidationError("invalid material_id")
	}

	// Validar userID
	userID, err := valueobject.UserIDFromString(userIDStr)
	if err != nil {
		s.logger.Error("invalid user_id", zap.Error(err))
		return errors.NewValidationError("invalid user_id")
	}

	// Determinar status basado en porcentaje
	var status string
	switch percentage {
	case 0:
		status = "not_started"
	case 100:
		status = "completed"
	default:
		status = "in_progress"
	}

	// Crear nueva entidad Progress con valores actualizados
	progress := &pgentities.Progress{
		MaterialID:     matID.UUID().UUID,
		UserID:         userID.UUID().UUID,
		Percentage:     percentage,
		LastPage:       lastPage,
		Status:         status,
		LastAccessedAt: time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Ejecutar operación UPSERT idempotente
	updatedProgress, err := s.progressRepo.Upsert(ctx, progress)
	if err != nil {
		s.logger.Error("failed to upsert progress",
			zap.Error(err),
			zap.String("material_id", materialID),
			zap.String("user_id", userIDStr),
		)
		return errors.NewDatabaseError("upsert progress", err)
	}

	// Verificar si material fue completado (progress = 100)
	isCompleted := updatedProgress.Percentage == 100
	if isCompleted {
		s.logger.Info("material completed by user",
			zap.String("material_id", materialID),
			zap.String("user_id", userIDStr),
			zap.Time("completed_at", updatedProgress.UpdatedAt),
		)

		// TODO (Fase futura): Publicar evento "material_completed" a RabbitMQ
		// Example:
		// event := events.MaterialCompleted{
		//     MaterialID: materialID,
		//     UserID: userIDStr,
		//     CompletedAt: updatedProgress.UpdatedAt(),
		// }
		// s.eventPublisher.Publish(ctx, "material.completed", event)
	}

	// Logging de éxito con métricas de performance
	elapsed := time.Since(startTime)
	s.logger.Info("progress updated successfully",
		zap.String("material_id", materialID),
		zap.String("user_id", userIDStr),
		zap.Int("percentage", percentage),
		zap.Bool("is_completed", isCompleted),
		zap.Duration("elapsed_ms", elapsed),
	)

	return nil
}
