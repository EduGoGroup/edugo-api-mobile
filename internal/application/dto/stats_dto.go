package dto

import "time"

// GlobalStatsDTO representa las estadísticas globales del sistema
type GlobalStatsDTO struct {
	TotalPublishedMaterials   int64     `json:"total_published_materials" example:"150"`
	TotalCompletedAssessments int64     `json:"total_completed_assessments" example:"1250"`
	AverageAssessmentScore    float64   `json:"average_assessment_score" example:"78.5"`
	ActiveUsersLast30Days     int64     `json:"active_users_last_30_days" example:"320"`
	AverageProgress           float64   `json:"average_progress" example:"65.3"`
	GeneratedAt               time.Time `json:"generated_at" example:"2024-01-15T10:30:00Z"`
}
