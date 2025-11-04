package service

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/common/types"
	"github.com/EduGoGroup/edugo-shared/logger"
	"go.uber.org/zap"
)

// AssessmentService define operaciones para assessments
type AssessmentService interface {
	GetAssessment(ctx context.Context, materialID string) (*repository.MaterialAssessment, error)
	RecordAttempt(ctx context.Context, materialID string, userID string, answers map[string]interface{}) (*repository.AssessmentAttempt, error)
}

type assessmentService struct {
	assessmentRepo   repository.AssessmentRepository
	messagePublisher rabbitmq.Publisher
	logger           logger.Logger
}

func NewAssessmentService(
	assessmentRepo repository.AssessmentRepository,
	messagePublisher rabbitmq.Publisher,
	logger logger.Logger,
) AssessmentService {
	return &assessmentService{
		assessmentRepo:   assessmentRepo,
		messagePublisher: messagePublisher,
		logger:           logger,
	}
}

func (s *assessmentService) GetAssessment(ctx context.Context, materialID string) (*repository.MaterialAssessment, error) {
	matID, err := valueobject.MaterialIDFromString(materialID)
	if err != nil {
		return nil, errors.NewValidationError("invalid material_id")
	}

	assessment, err := s.assessmentRepo.FindAssessmentByMaterialID(ctx, matID)
	if err != nil {
		s.logger.Error("failed to get assessment", "error", err)
		return nil, errors.NewDatabaseError("get assessment", err)
	}

	if assessment == nil {
		return nil, errors.NewNotFoundError("assessment")
	}

	return assessment, nil
}

func (s *assessmentService) RecordAttempt(ctx context.Context, materialID string, userIDStr string, answers map[string]interface{}) (*repository.AssessmentAttempt, error) {
	matID, err := valueobject.MaterialIDFromString(materialID)
	if err != nil {
		return nil, errors.NewValidationError("invalid material_id")
	}

	userID, err := valueobject.UserIDFromString(userIDStr)
	if err != nil {
		return nil, errors.NewValidationError("invalid user_id")
	}

	// Obtener assessment para calificar
	assessment, err := s.assessmentRepo.FindAssessmentByMaterialID(ctx, matID)
	if err != nil || assessment == nil {
		return nil, errors.NewNotFoundError("assessment")
	}

	// Calcular score (simplificado - en prod validar respuestas correctas)
	score := 75.0 // Mock score

	// Guardar intento
	attempt := &repository.AssessmentAttempt{
		ID:          types.NewUUID().String(),
		MaterialID:  matID,
		UserID:      userID,
		Answers:     answers,
		Score:       score,
		AttemptedAt: "",
	}

	if err := s.assessmentRepo.SaveAttempt(ctx, attempt); err != nil {
		s.logger.Error("failed to save attempt", "error", err)
		return nil, errors.NewDatabaseError("save attempt", err)
	}

	s.logger.Info("attempt recorded", "material_id", materialID, "user_id", userIDStr, "score", score)

	// Publicar evento de intento registrado
	event := messaging.AssessmentAttemptRecordedEvent{
		AttemptID:    attempt.ID,
		UserID:       userID.String(),
		AssessmentID: assessment.MaterialID.String(),
		Score:        score,
		SubmittedAt:  time.Now(),
	}

	eventJSON, err := event.ToJSON()
	if err != nil {
		s.logger.Warn("failed to serialize assessment attempt recorded event",
			zap.String("attempt_id", attempt.ID),
			zap.Error(err),
		)
	} else {
		// Publicar evento de forma asíncrona (no bloqueante)
		if err := s.messagePublisher.Publish(ctx, "edugo.materials", "assessment.attempt.recorded", eventJSON); err != nil {
			s.logger.Warn("failed to publish assessment attempt recorded event",
				zap.String("attempt_id", attempt.ID),
				zap.Error(err),
			)
		} else {
			s.logger.Info("assessment attempt recorded event published",
				zap.String("attempt_id", attempt.ID),
			)
		}
	}

	return attempt, nil
}
