package dto

import "time"

// GlobalStatsDTO representa las estadísticas globales del sistema
type GlobalStatsDTO struct {
	TotalPublishedMaterials   int64     `json:"total_published_materials"`
	TotalCompletedAssessments int64     `json:"total_completed_assessments"`
	AverageAssessmentScore    float64   `json:"average_assessment_score"`
	ActiveUsersLast30Days     int64     `json:"active_users_last_30_days"`
	AverageProgress           float64   `json:"average_progress"`
	GeneratedAt               time.Time `json:"generated_at"`
}
