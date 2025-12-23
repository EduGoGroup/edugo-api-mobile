package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repositories"
	domainServices "github.com/EduGoGroup/edugo-api-mobile/internal/domain/services"
	mongoRepo "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/mongodb/repository"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
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
	assessmentRepo      repositories.AssessmentRepository
	attemptRepo         repositories.AttemptRepository
	answerRepo          repositories.AnswerRepository
	mongoRepo           mongoRepo.AssessmentDocumentRepository
	assessmentDomainSvc *domainServices.AssessmentDomainService
	attemptDomainSvc    *domainServices.AttemptDomainService
	logger              logger.Logger
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
		assessmentRepo:      assessmentRepo,
		attemptRepo:         attemptRepo,
		answerRepo:          answerRepo,
		mongoRepo:           mongoRepo,
		assessmentDomainSvc: domainServices.NewAssessmentDomainService(),
		attemptDomainSvc:    domainServices.NewAttemptDomainService(),
		logger:              logger,
	}
}

// GetAssessmentByMaterialID obtiene assessment SIN respuestas correctas
func (s *assessmentAttemptService) GetAssessmentByMaterialID(ctx context.Context, materialID uuid.UUID) (*dto.AssessmentResponse, error) {
	// 1. Buscar assessment en PostgreSQL
	assessment, err := s.assessmentRepo.FindByMaterialID(ctx, materialID)
	if err != nil {
		s.logger.Error("failed to find assessment", "error", err)
		return nil, errors.NewDatabaseError("find assessment", err)
	}
	if assessment == nil {
		return nil, errors.NewNotFoundError("assessment")
	}

	// 2. Buscar preguntas en MongoDB usando mongo_document_id
	mongoDoc, err := s.mongoRepo.FindByID(ctx, assessment.MongoDocumentID)
	if err != nil {
		s.logger.Error("failed to find mongo document", "error", err)
		return nil, errors.NewDatabaseError("find mongo document", err)
	}
	if mongoDoc == nil {
		return nil, errors.NewNotFoundError("assessment questions")
	}

	// 3. Sanitizar preguntas (remover correct_answer y feedback)
	sanitizedQuestions := sanitizeQuestions(mongoDoc.Questions)

	// 4. Convertir campos nullable a valores concretos para DTO
	title := ""
	if assessment.Title != nil {
		title = *assessment.Title
	}

	totalQuestions := assessment.QuestionsCount // Usar QuestionsCount no-nullable
	if assessment.TotalQuestions != nil {
		totalQuestions = *assessment.TotalQuestions
	}

	passThreshold := 60 // Default
	if assessment.PassThreshold != nil {
		passThreshold = *assessment.PassThreshold
	}

	// 5. Construir response DTO
	return &dto.AssessmentResponse{
		AssessmentID:         assessment.ID,
		MaterialID:           assessment.MaterialID,
		Title:                title,
		TotalQuestions:       totalQuestions,
		PassThreshold:        passThreshold,
		MaxAttempts:          assessment.MaxAttempts,      // Nullable OK en DTO
		TimeLimitMinutes:     assessment.TimeLimitMinutes, // Nullable OK en DTO
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
		s.logger.Error("failed to count attempts", "error", err)
		return nil, errors.NewDatabaseError("count attempts", err)
	}

	if !s.assessmentDomainSvc.CanAttempt(assessment, attemptCount) {
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

	// 6. Calcular timestamps y score
	completedAt := startTime.Add(time.Duration(req.TimeSpentSeconds) * time.Second)
	score := s.attemptDomainSvc.CalculateScore(answers)
	percentage := score // Score ya es porcentaje (0-100)
	maxScore := 100.0
	timeSpent := req.TimeSpentSeconds

	// 7. Crear entity Attempt manualmente (no existe constructor NewAttempt)
	attemptID := uuid.New()
	attempt := &pgentities.AssessmentAttempt{
		ID:               attemptID,
		AssessmentID:     assessment.ID,
		StudentID:        studentID,
		StartedAt:        startTime,
		CompletedAt:      &completedAt,
		Score:            &score,
		MaxScore:         &maxScore,
		Percentage:       &percentage,
		TimeSpentSeconds: &timeSpent,
		Status:           "completed",
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	// Asignar attemptID a todas las respuestas
	for _, answer := range answers {
		answer.AttemptID = attemptID
	}

	// 8. Persistir intento (PostgreSQL con transacción ACID)
	if err := s.attemptRepo.Save(ctx, attempt); err != nil {
		s.logger.Error("failed to save attempt", "error", err)
		return nil, errors.NewDatabaseError("save attempt", err)
	}

	// 9. Persistir respuestas
	if err := s.answerRepo.Save(ctx, answers); err != nil {
		s.logger.Error("failed to save answers", "error", err)
		return nil, errors.NewDatabaseError("save answers", err)
	}

	// 10. Verificar si puede hacer más intentos
	canRetake := s.assessmentDomainSvc.CanAttempt(assessment, attemptCount+1)

	// 11. Calcular previous best score (opcional)
	var previousBestScore *int
	previousAttempts, _ := s.attemptRepo.FindByStudentAndAssessment(ctx, studentID, assessment.ID)
	if len(previousAttempts) > 1 { // Más de 1 porque ya guardamos el actual
		best := 0.0
		for _, prev := range previousAttempts {
			if prev.ID != attempt.ID && prev.Score != nil && *prev.Score > best {
				best = *prev.Score
			}
		}
		if best > 0 {
			bestInt := int(best)
			previousBestScore = &bestInt
		}
	}

	s.logger.Info("attempt created successfully",
		"attempt_id", attempt.ID.String(),
		"student_id", studentID.String(),
		"score", score,
		"correct_answers", correctCount,
	)

	// 12. Obtener pass threshold (nullable, default 60)
	passThreshold := 60
	if assessment.PassThreshold != nil {
		passThreshold = *assessment.PassThreshold
	}

	// 13. Retornar resultado con feedback
	return &dto.AttemptResultResponse{
		AttemptID:         attempt.ID,
		AssessmentID:      assessment.ID,
		Score:             int(score),
		MaxScore:          100,
		CorrectAnswers:    correctCount,
		TotalQuestions:    len(mongoDoc.Questions),
		PassThreshold:     passThreshold, // DTO espera int, no *int
		Passed:            s.attemptDomainSvc.IsPassed(attempt, passThreshold),
		TimeSpentSeconds:  timeSpent,
		StartedAt:         attempt.StartedAt,
		CompletedAt:       *attempt.CompletedAt, // Desreferenciar *time.Time
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
		s.logger.Error("failed to find attempt", "error", err)
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

	// 5. Cargar respuestas del intento (no están en la entity)
	answers, err := s.answerRepo.FindByAttemptID(ctx, attemptID)
	if err != nil {
		s.logger.Error("failed to find answers", "error", err)
		return nil, errors.NewDatabaseError("find answers", err)
	}

	// 6. Generar feedback desde answers
	feedback := s.generateFeedback(mongoDoc.Questions, answers)

	// 7. Verificar si puede hacer más intentos
	attemptCount, _ := s.attemptRepo.CountByStudentAndAssessment(ctx, studentID, assessment.ID)
	canRetake := s.assessmentDomainSvc.CanAttempt(assessment, attemptCount)

	// 8. Calcular previous best score
	var previousBestScore *int
	previousAttempts, _ := s.attemptRepo.FindByStudentAndAssessment(ctx, studentID, assessment.ID)
	if len(previousAttempts) > 1 {
		best := 0.0
		for _, prev := range previousAttempts {
			if prev.ID != attempt.ID && prev.Score != nil && *prev.Score > best {
				best = *prev.Score
			}
		}
		if best > 0 {
			bestInt := int(best)
			previousBestScore = &bestInt
		}
	}

	// 9. Obtener pass threshold (nullable, default 60)
	passThreshold := 60
	if assessment.PassThreshold != nil {
		passThreshold = *assessment.PassThreshold
	}

	// 10. Obtener score del attempt (nullable)
	score := 0
	if attempt.Score != nil {
		score = int(*attempt.Score)
	}

	// 11. Obtener time spent (nullable)
	timeSpent := 0
	if attempt.TimeSpentSeconds != nil {
		timeSpent = *attempt.TimeSpentSeconds
	}

	return &dto.AttemptResultResponse{
		AttemptID:         attempt.ID,
		AssessmentID:      assessment.ID,
		Score:             score,
		MaxScore:          100,
		CorrectAnswers:    s.attemptDomainSvc.GetCorrectAnswersCount(answers),
		TotalQuestions:    s.attemptDomainSvc.GetTotalQuestions(answers),
		PassThreshold:     passThreshold, // DTO espera int, no *int
		Passed:            s.attemptDomainSvc.IsPassed(attempt, passThreshold),
		TimeSpentSeconds:  timeSpent,
		StartedAt:         attempt.StartedAt,
		CompletedAt:       *attempt.CompletedAt, // Desreferenciar *time.Time
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
		s.logger.Error("failed to find attempts", "error", err)
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

		// Obtener pass threshold (nullable, default 60)
		passThreshold := 60
		if assessment.PassThreshold != nil {
			passThreshold = *assessment.PassThreshold
		}

		// Obtener score (nullable)
		score := 0
		if attempt.Score != nil {
			score = int(*attempt.Score)
		}

		// Obtener time spent (nullable)
		timeSpent := 0
		if attempt.TimeSpentSeconds != nil {
			timeSpent = *attempt.TimeSpentSeconds
		}

		// Obtener material title (nullable)
		materialTitle := ""
		if assessment.Title != nil {
			materialTitle = *assessment.Title
		}

		summaries = append(summaries, dto.AttemptSummaryDTO{
			AttemptID:        attempt.ID,
			AssessmentID:     assessment.ID,
			MaterialID:       assessment.MaterialID,
			MaterialTitle:    materialTitle, // DTO espera string, no *string
			Score:            score,
			MaxScore:         100,
			Passed:           s.attemptDomainSvc.IsPassed(attempt, passThreshold),
			TimeSpentSeconds: timeSpent,
			CompletedAt:      *attempt.CompletedAt, // Desreferenciar *time.Time
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
) ([]*pgentities.AssessmentAttemptAnswer, int, []dto.AnswerFeedbackDTO) {
	answers := make([]*pgentities.AssessmentAttemptAnswer, 0, len(userAnswers))
	feedback := make([]dto.AnswerFeedbackDTO, 0, len(userAnswers))
	correctCount := 0

	// Crear mapa de preguntas para lookup rápido (por índice)
	questionMap := make(map[string]mongoRepo.Question)
	for i, q := range questions {
		questionMap[q.ID] = q
		// También mapear por índice para facilitar lookup
		questionMap[fmt.Sprintf("%d", i)] = q
	}

	for i, userAnswer := range userAnswers {
		// Buscar pregunta correspondiente
		question, exists := questionMap[userAnswer.QuestionID]
		if !exists {
			s.logger.Warn("invalid question_id", "question_id", userAnswer.QuestionID)
			continue
		}

		// Comparar respuesta del usuario con respuesta correcta (servidor-side)
		isCorrect := question.CorrectAnswer == userAnswer.SelectedAnswerID

		if isCorrect {
			correctCount++
		}

		// Calcular puntos (simple: 100/total_questions por respuesta correcta)
		pointsEarned := 0.0
		if isCorrect {
			pointsEarned = 100.0 / float64(len(questions))
		}
		maxPoints := 100.0 / float64(len(questions))

		// Crear entity Answer manualmente (no existe constructor NewAnswer)
		answer := &pgentities.AssessmentAttemptAnswer{
			ID:               uuid.New(),
			AttemptID:        uuid.Nil, // Se asigna después
			QuestionIndex:    i,        // Usar índice basado en posición
			StudentAnswer:    &userAnswer.SelectedAnswerID,
			IsCorrect:        &isCorrect,
			PointsEarned:     &pointsEarned,
			MaxPoints:        &maxPoints,
			TimeSpentSeconds: &userAnswer.TimeSpentSeconds,
			AnsweredAt:       time.Now().UTC(),
			CreatedAt:        time.Now().UTC(),
			UpdatedAt:        time.Now().UTC(),
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
func (s *assessmentAttemptService) generateFeedback(questions []mongoRepo.Question, answers []*pgentities.AssessmentAttemptAnswer) []dto.AnswerFeedbackDTO {
	feedback := make([]dto.AnswerFeedbackDTO, 0, len(answers))

	// Crear array indexado de preguntas para lookup por QuestionIndex
	if len(questions) == 0 {
		return feedback
	}

	for _, answer := range answers {
		// Validar que el índice esté dentro del rango
		if answer.QuestionIndex < 0 || answer.QuestionIndex >= len(questions) {
			s.logger.Warn("invalid question_index", "index", answer.QuestionIndex)
			continue
		}

		question := questions[answer.QuestionIndex]

		// Obtener selectedOption (nullable)
		selectedOption := ""
		if answer.StudentAnswer != nil {
			selectedOption = *answer.StudentAnswer
		}

		// Obtener isCorrect (nullable, default false)
		isCorrect := false
		if answer.IsCorrect != nil {
			isCorrect = *answer.IsCorrect
		}

		var message string
		if isCorrect {
			message = question.Feedback.Correct
		} else {
			message = question.Feedback.Incorrect
		}

		feedback = append(feedback, dto.AnswerFeedbackDTO{
			QuestionID:     question.ID,
			QuestionText:   question.Text,
			SelectedOption: selectedOption,
			CorrectAnswer:  question.CorrectAnswer,
			IsCorrect:      isCorrect,
			Message:        message,
		})
	}

	return feedback
}
