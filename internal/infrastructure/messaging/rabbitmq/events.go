package rabbitmq

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Event es la estructura envelope estándar para todos los eventos
type Event struct {
	EventID      string      `json:"event_id"`
	EventType    string      `json:"event_type"`
	EventVersion string      `json:"event_version"`
	Timestamp    time.Time   `json:"timestamp"`
	Payload      interface{} `json:"payload"`
}

// MaterialUploadedPayload representa el payload del evento material.uploaded
type MaterialUploadedPayload struct {
	MaterialID    string                 `json:"material_id"`
	SchoolID      string                 `json:"school_id"`
	TeacherID     string                 `json:"teacher_id"`
	FileURL       string                 `json:"file_url"`
	FileSizeBytes int64                  `json:"file_size_bytes"`
	FileType      string                 `json:"file_type"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// NewMaterialUploadedEvent crea un nuevo evento material.uploaded con envelope estándar
func NewMaterialUploadedEvent(payload MaterialUploadedPayload) Event {
	return Event{
		EventID:      uuid.New().String(),
		EventType:    "material.uploaded",
		EventVersion: "1.0",
		Timestamp:    time.Now().UTC(),
		Payload:      payload,
	}
}

// MaterialCompletedPayload representa el payload del evento material.completed
// Se publica cuando un usuario completa un material (progress = 100%)
type MaterialCompletedPayload struct {
	MaterialID  string    `json:"material_id"`
	UserID      string    `json:"user_id"`
	CompletedAt time.Time `json:"completed_at"`
}

// NewMaterialCompletedEvent crea un nuevo evento material.completed con envelope estándar
func NewMaterialCompletedEvent(payload MaterialCompletedPayload) Event {
	return Event{
		EventID:      uuid.New().String(),
		EventType:    "material.completed",
		EventVersion: "1.0",
		Timestamp:    time.Now().UTC(),
		Payload:      payload,
	}
}

// AssessmentGeneratedPayload representa el payload del evento assessment.generated
type AssessmentGeneratedPayload struct {
	MaterialID       string `json:"material_id"`
	MongoDocumentID  string `json:"mongo_document_id"`
	QuestionsCount   int    `json:"questions_count"`
	ProcessingTimeMs int    `json:"processing_time_ms,omitempty"`
}

// NewAssessmentGeneratedEvent crea un nuevo evento assessment.generated con envelope estándar
func NewAssessmentGeneratedEvent(payload AssessmentGeneratedPayload) Event {
	return Event{
		EventID:      uuid.New().String(),
		EventType:    "assessment.generated",
		EventVersion: "1.0",
		Timestamp:    time.Now().UTC(),
		Payload:      payload,
	}
}

// ToJSON serializa el evento a JSON
func (e Event) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
