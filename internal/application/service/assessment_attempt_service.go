package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/entities"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repositories"
	mongoRepo "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/mongodb/repository"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// AssessmentAttemptService maneja la lógica de negocio de intentos de evaluación
// Orquesta repositorios PostgreSQL y MongoDB
type AssessmentAttemptService interface {
	// GetAssessmentByMaterialID obtiene un assessment SIN respuestas correctas (sanitizado)
	GetAssessmentByMaterialID(ctx context.Context, materialID uuid.UUID) (*dto.AssessmentResponse, error)

	// CreateAttempt crea un intento, valida respuestas y calcula score en servidor
	CreateAttempt(ctx context.Context, studentID, materialID uuid.UUID, req dto.CreateAttemptRequest) (*dto.AttemptResultResponse, error)

	// GetAttemptResult obtiene los resultados de un intento específico
	GetAttemptResult(ctx context.Context, attemptID, studentID uuid.UUID) (*dto.AttemptResultResponse, error)

	// GetAttemptHistory obtiene el historial de intentos de un estudiante
	GetAttemptHistory(ctx context.Context, studentID uuid.UUID, limit, offset int) (*dto.AttemptHistoryResponse, error)
}

type assessmentAttemptService struct {
	assessmentRepo repositories.AssessmentRepository
	attemptRepo    repositories.AttemptRepository
	answerRepo     repositories.AnswerRepository
	mongoRepo      mongoRepo.AssessmentDocumentRepository
	logger         logger.Logger
}

// NewAssessmentAttemptService crea una nueva instancia del servicio
func NewAssessmentAttemptService(
	assessmentRepo repositories.AssessmentRepository,
	attemptRepo repositories.AttemptRepository,
	answerRepo repositories.AnswerRepository,
	mongoRepo mongoRepo.AssessmentDocumentRepository,
	logger logger.Logger,
) AssessmentAttemptService {
	return &assessmentAttemptService{
		assessmentRepo: assessmentRepo,
		attemptRepo:    attemptRepo,
		answerRepo:     answerRepo,
		mongoRepo:      mongoRepo,
		logger:         logger,
	}
}

// GetAssessmentByMaterialID obtiene assessment SIN respuestas correctas
func (s *assessmentAttemptService) GetAssessmentByMaterialID(ctx context.Context, materialID uuid.UUID) (*dto.AssessmentResponse, error) {
	// 1. Buscar assessment en PostgreSQL
	assessment, err := s.assessmentRepo.FindByMaterialID(ctx, materialID)
	if err != nil {
		s.logger.Error("failed to find assessment", zap.Error(err))
		return nil, errors.NewDatabaseError("find assessment", err)
	}
	if assessment == nil {
		return nil, errors.NewNotFoundError("assessment")
	}

	// 2. Buscar preguntas en MongoDB usando mongo_document_id
	mongoDoc, err := s.mongoRepo.FindByID(ctx, assessment.MongoDocumentID)
	if err != nil {
		s.logger.Error("failed to find mongo document", zap.Error(err))
		return nil, errors.NewDatabaseError("find mongo document", err)
	}
	if mongoDoc == nil {
		return nil, errors.NewNotFoundError("assessment questions")
	}

	// 3. Sanitizar preguntas (remover correct_answer y feedback)
	sanitizedQuestions := sanitizeQuestions(mongoDoc.Questions)

	// 4. Construir response DTO
	return &dto.AssessmentResponse{
		AssessmentID:         assessment.ID,
		MaterialID:           assessment.MaterialID,
		Title:                assessment.Title,
		TotalQuestions:       assessment.TotalQuestions,
		PassThreshold:        assessment.PassThreshold,
		MaxAttempts:          assessment.MaxAttempts,
		TimeLimitMinutes:     assessment.TimeLimitMinutes,
		EstimatedTimeMinutes: mongoDoc.Metadata.EstimatedTimeMinutes,
		Questions:            sanitizedQuestions,
	}, nil
}

// CreateAttempt crea un intento, valida respuestas y calcula score en servidor
func (s *assessmentAttemptService) CreateAttempt(ctx context.Context, studentID, materialID uuid.UUID, req dto.CreateAttemptRequest) (*dto.AttemptResultResponse, error) {
	startTime := time.Now()

	// 1. Buscar assessment
	assessment, err := s.assessmentRepo.FindByMaterialID(ctx, materialID)
	if err != nil || assessment == nil {
		return nil, errors.NewNotFoundError("assessment")
	}

	// 2. Verificar si puede hacer otro intento (max_attempts)
	attemptCount, err := s.attemptRepo.CountByStudentAndAssessment(ctx, studentID, assessment.ID)
	if err != nil {
		s.logger.Error("failed to count attempts", zap.Error(err))
		return nil, errors.NewDatabaseError("count attempts", err)
	}

	if !assessment.CanAttempt(attemptCount) {
		return nil, errors.NewValidationError("max attempts reached")
	}

	// 3. Obtener preguntas con respuestas correctas de MongoDB
	mongoDoc, err := s.mongoRepo.FindByID(ctx, assessment.MongoDocumentID)
	if err != nil || mongoDoc == nil {
		return nil, errors.NewNotFoundError("assessment questions")
	}

	// 4. Validar que todas las preguntas fueron respondidas
	if len(req.Answers) != len(mongoDoc.Questions) {
		return nil, errors.NewValidationError(fmt.Sprintf("incomplete answers: expected %d, got %d", len(mongoDoc.Questions), len(req.Answers)))
	}

	// 5. VALIDAR RESPUESTAS Y CALCULAR SCORE EN SERVIDOR
	answers, correctCount, feedback := s.validateAndScoreAnswers(mongoDoc.Questions, req.Answers)

	// 6. Calcular timestamps
	completedAt := startTime.Add(time.Duration(req.TimeSpentSeconds) * time.Second)

	// 7. Crear entity Attempt con score calculado en servidor
	attempt, err := entities.NewAttempt(
		assessment.ID,
		studentID,
		answers,
		startTime,
		completedAt,
	)
	if err != nil {
		s.logger.Error("failed to create attempt entity", zap.Error(err))
		return nil, errors.NewValidationError("invalid attempt data")
	}

	// 8. Persistir intento (PostgreSQL con transacción ACID)
	if err := s.attemptRepo.Save(ctx, attempt); err != nil {
		s.logger.Error("failed to save attempt", zap.Error(err))
		return nil, errors.NewDatabaseError("save attempt", err)
	}

	// 9. Verificar si puede hacer más intentos
	canRetake := assessment.CanAttempt(attemptCount + 1)

	// 10. Calcular previous best score (opcional)
	var previousBestScore *int
	previousAttempts, _ := s.attemptRepo.FindByStudentAndAssessment(ctx, studentID, assessment.ID)
	if len(previousAttempts) > 1 { // Más de 1 porque ya guardamos el actual
		best := 0
		for _, prev := range previousAttempts {
			if prev.ID != attempt.ID && prev.Score > best {
				best = prev.Score
			}
		}
		if best > 0 {
			previousBestScore = &best
		}
	}

	s.logger.Info("attempt created successfully",
		zap.String("attempt_id", attempt.ID.String()),
		zap.String("student_id", studentID.String()),
		zap.Int("score", attempt.Score),
		zap.Int("correct_answers", correctCount),
	)

	// 11. Retornar resultado con feedback
	return &dto.AttemptResultResponse{
		AttemptID:         attempt.ID,
		AssessmentID:      assessment.ID,
		Score:             attempt.Score,
		MaxScore:          100,
		CorrectAnswers:    correctCount,
		TotalQuestions:    len(mongoDoc.Questions),
		PassThreshold:     assessment.PassThreshold,
		Passed:            attempt.IsPassed(assessment.PassThreshold),
		TimeSpentSeconds:  attempt.TimeSpentSeconds,
		StartedAt:         attempt.StartedAt,
		CompletedAt:       attempt.CompletedAt,
		Feedback:          feedback,
		CanRetake:         canRetake,
		PreviousBestScore: previousBestScore,
	}, nil
}

// GetAttemptResult obtiene los resultados de un intento específico
func (s *assessmentAttemptService) GetAttemptResult(ctx context.Context, attemptID, studentID uuid.UUID) (*dto.AttemptResultResponse, error) {
	// 1. Buscar intento
	attempt, err := s.attemptRepo.FindByID(ctx, attemptID)
	if err != nil {
		s.logger.Error("failed to find attempt", zap.Error(err))
		return nil, errors.NewDatabaseError("find attempt", err)
	}
	if attempt == nil {
		return nil, errors.NewNotFoundError("attempt")
	}

	// 2. Verificar que el intento pertenece al estudiante (autorización)
	if attempt.StudentID != studentID {
		return nil, errors.NewForbiddenError("attempt does not belong to user")
	}

	// 3. Buscar assessment para obtener metadata
	assessment, err := s.assessmentRepo.FindByID(ctx, attempt.AssessmentID)
	if err != nil || assessment == nil {
		return nil, errors.NewNotFoundError("assessment")
	}

	// 4. Buscar preguntas en MongoDB para generar feedback
	mongoDoc, err := s.mongoRepo.FindByID(ctx, assessment.MongoDocumentID)
	if err != nil || mongoDoc == nil {
		return nil, errors.NewNotFoundError("assessment questions")
	}

	// 5. Generar feedback desde answers
	feedback := s.generateFeedback(mongoDoc.Questions, attempt.Answers)

	// 6. Verificar si puede hacer más intentos
	attemptCount, _ := s.attemptRepo.CountByStudentAndAssessment(ctx, studentID, assessment.ID)
	canRetake := assessment.CanAttempt(attemptCount)

	// 7. Calcular previous best score
	var previousBestScore *int
	previousAttempts, _ := s.attemptRepo.FindByStudentAndAssessment(ctx, studentID, assessment.ID)
	if len(previousAttempts) > 1 {
		best := 0
		for _, prev := range previousAttempts {
			if prev.ID != attempt.ID && prev.Score > best {
				best = prev.Score
			}
		}
		if best > 0 {
			previousBestScore = &best
		}
	}

	return &dto.AttemptResultResponse{
		AttemptID:         attempt.ID,
		AssessmentID:      assessment.ID,
		Score:             attempt.Score,
		MaxScore:          100,
		CorrectAnswers:    attempt.GetCorrectAnswersCount(),
		TotalQuestions:    attempt.GetTotalQuestions(),
		PassThreshold:     assessment.PassThreshold,
		Passed:            attempt.IsPassed(assessment.PassThreshold),
		TimeSpentSeconds:  attempt.TimeSpentSeconds,
		StartedAt:         attempt.StartedAt,
		CompletedAt:       attempt.CompletedAt,
		Feedback:          feedback,
		CanRetake:         canRetake,
		PreviousBestScore: previousBestScore,
	}, nil
}

// GetAttemptHistory obtiene el historial de intentos de un estudiante
func (s *assessmentAttemptService) GetAttemptHistory(ctx context.Context, studentID uuid.UUID, limit, offset int) (*dto.AttemptHistoryResponse, error) {
	// 1. Buscar intentos del estudiante
	attempts, err := s.attemptRepo.FindByStudent(ctx, studentID, limit, offset)
	if err != nil {
		s.logger.Error("failed to find attempts", zap.Error(err))
		return nil, errors.NewDatabaseError("find attempts", err)
	}

	// 2. Construir summaries
	summaries := make([]dto.AttemptSummaryDTO, 0, len(attempts))
	for _, attempt := range attempts {
		// Buscar assessment para obtener material_id y title
		assessment, err := s.assessmentRepo.FindByID(ctx, attempt.AssessmentID)
		if err != nil || assessment == nil {
			continue // Skip si no se encuentra assessment
		}

		summaries = append(summaries, dto.AttemptSummaryDTO{
			AttemptID:        attempt.ID,
			AssessmentID:     assessment.ID,
			MaterialID:       assessment.MaterialID,
			MaterialTitle:    assessment.Title,
			Score:            attempt.Score,
			MaxScore:         100,
			Passed:           attempt.IsPassed(assessment.PassThreshold),
			TimeSpentSeconds: attempt.TimeSpentSeconds,
			CompletedAt:      attempt.CompletedAt,
		})
	}

	// 3. Calcular page
	page := (offset / limit) + 1
	if limit == 0 {
		page = 1
	}

	return &dto.AttemptHistoryResponse{
		Attempts:   summaries,
		TotalCount: len(summaries),
		Page:       page,
		Limit:      limit,
	}, nil
}

// ========== HELPERS ==========

// sanitizeQuestions remueve correct_answer y feedback de las preguntas
// CRÍTICO: Nunca exponer respuestas correctas al cliente
func sanitizeQuestions(questions []mongoRepo.Question) []dto.QuestionDTO {
	sanitized := make([]dto.QuestionDTO, len(questions))
	for i, q := range questions {
		options := make([]dto.OptionDTO, len(q.Options))
		for j, opt := range q.Options {
			options[j] = dto.OptionDTO{
				ID:   opt.ID,
				Text: opt.Text,
			}
		}

		sanitized[i] = dto.QuestionDTO{
			ID:      q.ID,
			Text:    q.Text,
			Type:    q.Type,
			Options: options,
			// ❌ NO incluir: CorrectAnswer, Feedback
		}
	}
	return sanitized
}

// validateAndScoreAnswers valida respuestas contra MongoDB y calcula score en servidor
// CRÍTICO: Score SIEMPRE calculado en servidor, NUNCA confiar en cliente
func (s *assessmentAttemptService) validateAndScoreAnswers(
	questions []mongoRepo.Question,
	userAnswers []dto.UserAnswerDTO,
) ([]*entities.Answer, int, []dto.AnswerFeedbackDTO) {
	answers := make([]*entities.Answer, 0, len(userAnswers))
	feedback := make([]dto.AnswerFeedbackDTO, 0, len(userAnswers))
	correctCount := 0

	// Crear mapa de preguntas para lookup rápido
	questionMap := make(map[string]mongoRepo.Question)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	for _, userAnswer := range userAnswers {
		// Buscar pregunta correspondiente
		question, exists := questionMap[userAnswer.QuestionID]
		if !exists {
			s.logger.Warn("invalid question_id", zap.String("question_id", userAnswer.QuestionID))
			continue
		}

		// Comparar respuesta del usuario con respuesta correcta (servidor-side)
		isCorrect := question.CorrectAnswer == userAnswer.SelectedAnswerID

		if isCorrect {
			correctCount++
		}

		// Crear entity Answer
		answer, err := entities.NewAnswer(
			uuid.Nil, // Attempt ID se asigna después
			userAnswer.QuestionID,
			userAnswer.SelectedAnswerID,
			isCorrect,
			userAnswer.TimeSpentSeconds,
		)
		if err != nil {
			s.logger.Error("failed to create answer entity", zap.Error(err))
			continue
		}

		answers = append(answers, answer)

		// Generar feedback educativo
		var message string
		if isCorrect {
			message = question.Feedback.Correct
		} else {
			message = question.Feedback.Incorrect
		}

		feedback = append(feedback, dto.AnswerFeedbackDTO{
			QuestionID:     question.ID,
			QuestionText:   question.Text,
			SelectedOption: userAnswer.SelectedAnswerID,
			CorrectAnswer:  question.CorrectAnswer,
			IsCorrect:      isCorrect,
			Message:        message,
		})
	}

	return answers, correctCount, feedback
}

// generateFeedback genera feedback desde answers persistidas
func (s *assessmentAttemptService) generateFeedback(questions []mongoRepo.Question, answers []*entities.Answer) []dto.AnswerFeedbackDTO {
	feedback := make([]dto.AnswerFeedbackDTO, 0, len(answers))

	// Crear mapa de preguntas
	questionMap := make(map[string]mongoRepo.Question)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	for _, answer := range answers {
		question, exists := questionMap[answer.QuestionID]
		if !exists {
			continue
		}

		var message string
		if answer.IsCorrect {
			message = question.Feedback.Correct
		} else {
			message = question.Feedback.Incorrect
		}

		feedback = append(feedback, dto.AnswerFeedbackDTO{
			QuestionID:     answer.QuestionID,
			QuestionText:   question.Text,
			SelectedOption: answer.SelectedAnswerID,
			CorrectAnswer:  question.CorrectAnswer,
			IsCorrect:      answer.IsCorrect,
			Message:        message,
		})
	}

	return feedback
}
