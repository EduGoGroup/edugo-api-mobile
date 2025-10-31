package service

import (
	"context"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/pkg/errors"
	"github.com/EduGoGroup/edugo-shared/pkg/logger"
	"github.com/EduGoGroup/edugo-shared/pkg/types"
)

// AssessmentService define operaciones para assessments
type AssessmentService interface {
	GetAssessment(ctx context.Context, materialID string) (*repository.MaterialAssessment, error)
	RecordAttempt(ctx context.Context, materialID string, userID string, answers map[string]interface{}) (*repository.AssessmentAttempt, error)
}

type assessmentService struct {
	assessmentRepo repository.AssessmentRepository
	logger         logger.Logger
}

func NewAssessmentService(assessmentRepo repository.AssessmentRepository, logger logger.Logger) AssessmentService {
	return &assessmentService{
		assessmentRepo: assessmentRepo,
		logger:         logger,
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

	return attempt, nil
}
