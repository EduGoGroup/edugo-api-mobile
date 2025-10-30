package mongodb

import "time"

// MaterialSummaryDocument representa un documento de resumen en MongoDB
type MaterialSummaryDocument struct {
	MaterialID          string             `bson:"material_id"`
	Version             int                `bson:"version"`
	Status              string             `bson:"status"` // pending | processing | completed | failed
	Sections            []SummarySection   `bson:"sections"`
	Glossary            []GlossaryTerm     `bson:"glossary,omitempty"`
	ReflectionQuestions []string           `bson:"reflection_questions,omitempty"`
	ProcessingMetadata  ProcessingMetadata `bson:"processing_metadata,omitempty"`
	CreatedAt           time.Time          `bson:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at,omitempty"`
}

// SummarySection representa una sección del resumen en MongoDB
type SummarySection struct {
	Title                string `bson:"title"`
	Content              string `bson:"content"`
	Difficulty           string `bson:"difficulty"` // basic | medium | advanced
	EstimatedTimeMinutes int    `bson:"estimated_time_minutes,omitempty"`
	Order                int    `bson:"order"`
}

// GlossaryTerm representa un término del glosario en MongoDB
type GlossaryTerm struct {
	Term       string `bson:"term"`
	Definition string `bson:"definition"`
	Order      int    `bson:"order"`
}

// ProcessingMetadata contiene metadata del procesamiento NLP
type ProcessingMetadata struct {
	NLPProvider           string `bson:"nlp_provider,omitempty"`
	Model                 string `bson:"model,omitempty"`
	TokensUsed            int    `bson:"tokens_used,omitempty"`
	ProcessingTimeSeconds int    `bson:"processing_time_seconds,omitempty"`
	Language              string `bson:"language,omitempty"`
	PromptVersion         string `bson:"prompt_version,omitempty"`
}

// MaterialAssessmentDocument representa un documento de quiz en MongoDB
type MaterialAssessmentDocument struct {
	MaterialID         string             `bson:"material_id"`
	Title              string             `bson:"title"`
	Description        string             `bson:"description,omitempty"`
	Questions          []QuestionDocument `bson:"questions"`
	TotalQuestions     int                `bson:"total_questions"`
	TotalPoints        int                `bson:"total_points"`
	PassingScore       int                `bson:"passing_score"`
	TimeLimitMinutes   int                `bson:"time_limit_minutes,omitempty"`
	Version            int                `bson:"version"`
	ProcessingMetadata ProcessingMetadata `bson:"processing_metadata,omitempty"`
	CreatedAt          time.Time          `bson:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at,omitempty"`
}

// QuestionDocument representa una pregunta en MongoDB
type QuestionDocument struct {
	ID            string              `bson:"id"`
	Text          string              `bson:"text"`
	Type          string              `bson:"type"`                 // multiple_choice | true_false | short_answer
	Difficulty    string              `bson:"difficulty,omitempty"` // basic | medium | advanced
	Points        int                 `bson:"points"`
	Order         int                 `bson:"order"`
	Options       []QuestionOption    `bson:"options"`
	CorrectAnswer string              `bson:"correct_answer"`
	Feedback      QuestionFeedbackDoc `bson:"feedback,omitempty"`
}

// QuestionOption representa una opción de respuesta en MongoDB
type QuestionOption struct {
	ID   string `bson:"id"`
	Text string `bson:"text"`
}

// QuestionFeedbackDoc contiene feedback para respuestas correctas e incorrectas
type QuestionFeedbackDoc struct {
	Correct   string `bson:"correct"`
	Incorrect string `bson:"incorrect"`
}

// MaterialEventDocument representa un evento de consumo en MongoDB
type MaterialEventDocument struct {
	EventID    string                 `bson:"event_id"`
	MaterialID string                 `bson:"material_id"`
	StudentID  string                 `bson:"student_id"`
	EventType  string                 `bson:"event_type"` // summary_viewed | quiz_completed | material_downloaded
	EventData  map[string]interface{} `bson:"event_data"`
	Timestamp  time.Time              `bson:"timestamp"`
}
