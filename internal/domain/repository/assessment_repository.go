// Package repository contiene interfaces de repositorio para MongoDB (LEGACY)
//
// IMPORTANTE: Este paquete contiene el sistema LEGACY de assessments basado en MongoDB.
// El nuevo sistema de assessments usa PostgreSQL y está en internal/domain/repositories/.
//
// Estado de migración:
//   - AssessmentStats: ACTIVO - Usado por StatsService para estadísticas globales
//   - AssessmentReader: LEGACY - No se usa activamente, pendiente migración
//   - AssessmentWriter: LEGACY - No se usa activamente, pendiente migración
//
// Plan de consolidación: Ver docs/technical/ASSESSMENT_CONSOLIDATION.md
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
//
// DEPRECATED: Esta interfaz es parte del sistema legacy de MongoDB.
// El nuevo sistema usa PostgreSQL (ver internal/domain/repositories/assessment_repository.go).
// Mantener hasta completar migración de datos históricos.
type AssessmentReader interface {
	// FindAssessmentByMaterialID busca el assessment de un material
	FindAssessmentByMaterialID(ctx context.Context, materialID valueobject.MaterialID) (*MaterialAssessment, error)

	// FindAttemptsByUser busca los intentos de un usuario para un material
	FindAttemptsByUser(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) ([]*AssessmentAttempt, error)

	// GetBestAttempt obtiene el mejor intento de un usuario
	GetBestAttempt(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) (*AssessmentAttempt, error)
}

// AssessmentWriter define operaciones de escritura para assessments
//
// DEPRECATED: Esta interfaz es parte del sistema legacy de MongoDB.
// El nuevo sistema usa PostgreSQL (ver internal/domain/repositories/assessment_repository.go).
// Mantener hasta completar migración de datos históricos.
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
//
// ACTIVO: Esta interfaz SÍ se usa activamente por StatsService.
// Consulta datos de MongoDB para estadísticas globales del sistema.
// Pendiente: Migrar a PostgreSQL cuando se complete consolidación.
type AssessmentStats interface {
	// CountCompletedAssessments cuenta el total de evaluaciones completadas (para estadísticas)
	CountCompletedAssessments(ctx context.Context) (int64, error)

	// CalculateAverageScore calcula el promedio de puntajes de todas las evaluaciones completadas
	CalculateAverageScore(ctx context.Context) (float64, error)
}

// AssessmentRepository agrega todas las capacidades de assessments (MongoDB)
// Las implementaciones deben cumplir con todas las interfaces segregadas
//
// LEGACY: Este repositorio compuesto es parte del sistema legacy.
// Solo AssessmentStats se usa activamente. Reader y Writer están deprecated.
// Ver docs/technical/ASSESSMENT_CONSOLIDATION.md para plan de migración.
type AssessmentRepository interface {
	AssessmentReader // DEPRECATED
	AssessmentWriter // DEPRECATED
	AssessmentStats  // ACTIVO
}
