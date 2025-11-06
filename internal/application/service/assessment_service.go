package service

import (
	"context"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/service/scoring"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/common/types"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/EduGoGroup/edugo-shared/logger"
	"go.uber.org/zap"
)

// AssessmentService define operaciones para assessments
type AssessmentService interface {
	GetAssessment(ctx context.Context, materialID string) (*repository.MaterialAssessment, error)
	RecordAttempt(ctx context.Context, materialID string, userID string, answers map[string]interface{}) (*repository.AssessmentAttempt, error)
	CalculateScore(ctx context.Context, assessmentID string, userID string, userResponses map[string]interface{}) (*repository.AssessmentResult, error)
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

// CalculateScore calcula el puntaje de una evaluación y genera feedback detallado
// Implementa Strategy Pattern para evaluar diferentes tipos de preguntas
func (s *assessmentService) CalculateScore(ctx context.Context, assessmentID string, userIDStr string, userResponses map[string]interface{}) (*repository.AssessmentResult, error) {
	startTime := time.Now()

	// Validar IDs
	matID, err := valueobject.MaterialIDFromString(assessmentID)
	if err != nil {
		return nil, errors.NewValidationError("invalid assessment_id")
	}

	userID, err := valueobject.UserIDFromString(userIDStr)
	if err != nil {
		return nil, errors.NewValidationError("invalid user_id")
	}

	// 1. Obtener assessment con todas las preguntas
	assessment, err := s.assessmentRepo.FindAssessmentByMaterialID(ctx, matID)
	if err != nil {
		s.logger.Error("failed to fetch assessment", zap.Error(err))
		return nil, errors.NewDatabaseError("fetch assessment", err)
	}

	if assessment == nil {
		return nil, errors.NewNotFoundError("assessment")
	}

	// Validar que haya preguntas
	if len(assessment.Questions) == 0 {
		return nil, errors.NewValidationError("assessment has no questions")
	}

	// 2. Inicializar contadores y feedback
	totalQuestions := len(assessment.Questions)
	correctAnswers := 0
	feedback := make([]repository.FeedbackItem, 0, totalQuestions)

	s.logger.Info("calculating assessment score",
		zap.String("assessment_id", assessmentID),
		zap.String("user_id", userIDStr),
		zap.Int("total_questions", totalQuestions),
	)

	// 3. Crear estrategias de scoring
	strategies := map[enum.AssessmentType]scoring.ScoringStrategy{
		enum.AssessmentTypeMultipleChoice: scoring.NewMultipleChoiceStrategy(),
		enum.AssessmentTypeTrueFalse:      scoring.NewTrueFalseStrategy(),
		enum.AssessmentTypeShortAnswer:    scoring.NewShortAnswerStrategy(),
	}

	// 4. Iterar sobre cada pregunta y evaluar
	for _, question := range assessment.Questions {
		// Obtener respuesta del usuario para esta pregunta
		userAnswer, hasAnswer := userResponses[question.ID]
		if !hasAnswer {
			// Pregunta no respondida - marcar como incorrecta
			feedback = append(feedback, repository.FeedbackItem{
				QuestionID:    question.ID,
				IsCorrect:     false,
				UserAnswer:    "(sin respuesta)",
				CorrectAnswer: fmt.Sprintf("%v", question.CorrectAnswer),
				Explanation:   "No se proporcionó una respuesta para esta pregunta",
			})
			continue
		}

		// Seleccionar estrategia apropiada según tipo de pregunta
		strategy, exists := strategies[question.QuestionType]
		if !exists {
			s.logger.Warn("unsupported question type",
				zap.String("question_id", question.ID),
				zap.String("question_type", string(question.QuestionType)),
			)
			feedback = append(feedback, repository.FeedbackItem{
				QuestionID:    question.ID,
				IsCorrect:     false,
				UserAnswer:    fmt.Sprintf("%v", userAnswer),
				CorrectAnswer: fmt.Sprintf("%v", question.CorrectAnswer),
				Explanation:   "Tipo de pregunta no soportado para evaluación automática",
			})
			continue
		}

		// Calcular score para esta pregunta usando estrategia
		_, isCorrect, explanation := strategy.CalculateScore(question, userAnswer)

		// Incrementar contador si es correcta
		if isCorrect {
			correctAnswers++
		}

		// Agregar feedback de la pregunta
		feedback = append(feedback, repository.FeedbackItem{
			QuestionID:    question.ID,
			IsCorrect:     isCorrect,
			UserAnswer:    fmt.Sprintf("%v", userAnswer),
			CorrectAnswer: fmt.Sprintf("%v", question.CorrectAnswer),
			Explanation:   explanation,
		})
	}

	// 5. Calcular puntaje final: (correctas / totales) * 100
	finalScore := (float64(correctAnswers) / float64(totalQuestions)) * 100.0

	// 6. Crear resultado
	result := &repository.AssessmentResult{
		ID:             types.NewUUID().String(),
		AssessmentID:   assessmentID,
		UserID:         userID,
		Score:          finalScore,
		TotalQuestions: totalQuestions,
		CorrectAnswers: correctAnswers,
		Feedback:       feedback,
		SubmittedAt:    time.Now().Format(time.RFC3339),
	}

	// 7. Persistir resultado
	if err := s.assessmentRepo.SaveResult(ctx, result); err != nil {
		s.logger.Error("failed to save assessment result", zap.Error(err))
		return nil, errors.NewDatabaseError("save assessment result", err)
	}

	// Logging de éxito con métricas
	duration := time.Since(startTime)
	s.logger.Info("assessment score calculated successfully",
		zap.String("assessment_id", assessmentID),
		zap.String("user_id", userIDStr),
		zap.Float64("score", finalScore),
		zap.Int("correct_answers", correctAnswers),
		zap.Int("total_questions", totalQuestions),
		zap.Duration("duration", duration),
	)

	// 8. Publicar evento de evaluación completada (asíncrono, no bloqueante)
	event := messaging.AssessmentAttemptRecordedEvent{
		AttemptID:    result.ID,
		UserID:       userIDStr,
		AssessmentID: assessmentID,
		Score:        finalScore,
		SubmittedAt:  time.Now(),
	}

	eventJSON, err := event.ToJSON()
	if err != nil {
		s.logger.Warn("failed to serialize assessment completed event",
			zap.String("result_id", result.ID),
			zap.Error(err),
		)
	} else {
		if err := s.messagePublisher.Publish(ctx, "edugo.materials", "assessment.completed", eventJSON); err != nil {
			s.logger.Warn("failed to publish assessment completed event",
				zap.String("result_id", result.ID),
				zap.Error(err),
			)
		} else {
			s.logger.Info("assessment completed event published",
				zap.String("result_id", result.ID),
			)
		}
	}

	return result, nil
}
