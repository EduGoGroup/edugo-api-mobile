package enum

// SystemRole representa los roles del sistema
type SystemRole string

const (
	RoleAdmin    SystemRole = "admin"
	RoleTeacher  SystemRole = "teacher"
	RoleStudent  SystemRole = "student"
	RoleGuardian SystemRole = "guardian"
)

// MaterialStatus representa el estado de un material
type MaterialStatus string

const (
	MaterialDraft     MaterialStatus = "draft"
	MaterialPublished MaterialStatus = "published"
	MaterialArchived  MaterialStatus = "archived"
)

// AssessmentType representa tipos de preguntas
type AssessmentType string

const (
	AssessmentMultipleChoice AssessmentType = "multiple_choice"
	AssessmentTrueFalse      AssessmentType = "true_false"
	AssessmentShortAnswer    AssessmentType = "short_answer"
)

// ProgressStatus representa el estado de progreso de un estudiante
type ProgressStatus string

const (
	StatusNotStarted ProgressStatus = "not_started"
	StatusInProgress ProgressStatus = "in_progress"
	StatusCompleted  ProgressStatus = "completed"
)
