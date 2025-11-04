package messaging

import (
	"encoding/json"
	"time"
)

// MaterialUploadedEvent representa el evento cuando un material es creado/subido
type MaterialUploadedEvent struct {
	MaterialID  string    `json:"material_id"`
	Title       string    `json:"title"`
	ContentType string    `json:"content_type"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

// ToJSON serializa el evento a JSON
func (e MaterialUploadedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// AssessmentAttemptRecordedEvent representa el evento cuando se registra un intento de evaluación
type AssessmentAttemptRecordedEvent struct {
	AttemptID    string    `json:"attempt_id"`
	UserID       string    `json:"user_id"`
	AssessmentID string    `json:"assessment_id"`
	Score        float64   `json:"score"`
	SubmittedAt  time.Time `json:"submitted_at"`
}

// ToJSON serializa el evento a JSON
func (e AssessmentAttemptRecordedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
