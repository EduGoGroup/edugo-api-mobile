package repository

import (
	"context"

	"github.com/edugo/api-mobile/internal/domain/valueobject"
	"github.com/edugo/shared/pkg/types/enum"
)

// MaterialAssessment representa el quiz de un material (almacenado en MongoDB)
type MaterialAssessment struct {
	MaterialID valueobject.MaterialID
	Questions  []AssessmentQuestion
	CreatedAt  string
}

// AssessmentQuestion representa una pregunta del quiz
type AssessmentQuestion struct {
	ID             string
	QuestionText   string
	QuestionType   enum.AssessmentType
	Options        []string // Para multiple choice
	CorrectAnswer  interface{} // String o int dependiendo del tipo
	Explanation    string
	DifficultyLevel string
}

// AssessmentAttempt representa un intento de resolver el quiz
type AssessmentAttempt struct {
	ID         string
	MaterialID valueobject.MaterialID
	UserID     valueobject.UserID
	Answers    map[string]interface{} // question_id -> answer
	Score      float64
	AttemptedAt string
}

// AssessmentRepository define las operaciones para assessments (MongoDB)
type AssessmentRepository interface {
	// SaveAssessment guarda o actualiza un assessment
	SaveAssessment(ctx context.Context, assessment *MaterialAssessment) error

	// FindAssessmentByMaterialID busca el assessment de un material
	FindAssessmentByMaterialID(ctx context.Context, materialID valueobject.MaterialID) (*MaterialAssessment, error)

	// SaveAttempt guarda un intento de assessment
	SaveAttempt(ctx context.Context, attempt *AssessmentAttempt) error

	// FindAttemptsByUser busca los intentos de un usuario para un material
	FindAttemptsByUser(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) ([]*AssessmentAttempt, error)

	// GetBestAttempt obtiene el mejor intento de un usuario
	GetBestAttempt(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) (*AssessmentAttempt, error)
}
