package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// MaterialAssessment representa el quiz de un material (almacenado en MongoDB)
type MaterialAssessment struct {
	MaterialID valueobject.MaterialID
	Questions  []AssessmentQuestion
	CreatedAt  string
}

// AssessmentQuestion representa una pregunta del quiz
type AssessmentQuestion struct {
	ID              string
	QuestionText    string
	QuestionType    enum.AssessmentType
	Options         []string    // Para multiple choice
	CorrectAnswer   interface{} // String o int dependiendo del tipo
	Explanation     string
	DifficultyLevel string
}

// AssessmentAttempt representa un intento de resolver el quiz
type AssessmentAttempt struct {
	ID          string
	MaterialID  valueobject.MaterialID
	UserID      valueobject.UserID
	Answers     map[string]interface{} // question_id -> answer
	Score       float64
	AttemptedAt string
}

// FeedbackItem representa el feedback detallado de una pregunta
type FeedbackItem struct {
	QuestionID    string
	IsCorrect     bool
	UserAnswer    string
	CorrectAnswer string
	Explanation   string
}

// AssessmentResult representa el resultado de una evaluación completada
// Se almacena en la colección assessment_results (diferente de assessment_attempts)
type AssessmentResult struct {
	ID             string
	AssessmentID   string
	UserID         valueobject.UserID
	Score          float64
	TotalQuestions int
	CorrectAnswers int
	Feedback       []FeedbackItem
	SubmittedAt    string
}

// AssessmentReader define operaciones de lectura para assessments (MongoDB)
// Principio ISP: Separar lectura de escritura y estadísticas
type AssessmentReader interface {
	// FindAssessmentByMaterialID busca el assessment de un material
	FindAssessmentByMaterialID(ctx context.Context, materialID valueobject.MaterialID) (*MaterialAssessment, error)

	// FindAttemptsByUser busca los intentos de un usuario para un material
	FindAttemptsByUser(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) ([]*AssessmentAttempt, error)

	// GetBestAttempt obtiene el mejor intento de un usuario
	GetBestAttempt(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) (*AssessmentAttempt, error)
}

// AssessmentWriter define operaciones de escritura para assessments
type AssessmentWriter interface {
	// SaveAssessment guarda o actualiza un assessment
	SaveAssessment(ctx context.Context, assessment *MaterialAssessment) error

	// SaveAttempt guarda un intento de assessment
	SaveAttempt(ctx context.Context, attempt *AssessmentAttempt) error

	// SaveResult guarda el resultado de una evaluación completada en assessment_results
	// Retorna error si la evaluación ya fue completada por el usuario (índice UNIQUE)
	SaveResult(ctx context.Context, result *AssessmentResult) error
}

// AssessmentStats define operaciones de estadísticas para assessments
type AssessmentStats interface {
	// CountCompletedAssessments cuenta el total de evaluaciones completadas (para estadísticas)
	CountCompletedAssessments(ctx context.Context) (int64, error)

	// CalculateAverageScore calcula el promedio de puntajes de todas las evaluaciones completadas
	CalculateAverageScore(ctx context.Context) (float64, error)
}

// AssessmentRepository agrega todas las capacidades de assessments (MongoDB)
// Las implementaciones deben cumplir con todas las interfaces segregadas
type AssessmentRepository interface {
	AssessmentReader
	AssessmentWriter
	AssessmentStats
}
