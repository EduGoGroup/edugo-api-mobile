package response

import "time"

// MaterialListResponse representa la lista de materiales
type MaterialListResponse struct {
	Materials []MaterialSummary `json:"materials"`
	Total     int               `json:"total"`
	Page      int               `json:"page"`
} // @name MaterialListResponse

// MaterialSummary representa un resumen de material en la lista
type MaterialSummary struct {
	ID          string    `json:"id" example:"uuid-1"`
	Title       string    `json:"title" example:"Introducción a Pascal"`
	SubjectName string    `json:"subject_name" example:"Programación"`
	UnitName    string    `json:"unit_name" example:"5.º A - Programación"`
	Status      string    `json:"status" example:"new" enums:"new,in_progress,completed"`
	Progress    float64   `json:"progress" example:"0"`
	HasSummary  bool      `json:"has_summary" example:"true"`
	HasQuiz     bool      `json:"has_quiz" example:"true"`
	PublishedAt time.Time `json:"published_at" example:"2025-01-15T12:00:00Z"`
} // @name MaterialSummary

// MaterialDetailResponse representa el detalle completo de un material
type MaterialDetailResponse struct {
	Material        MaterialDetail `json:"material"`
	PDFURL          string         `json:"pdf_url" example:"https://s3.../presigned-url"`
	PDFURLExpiresAt time.Time      `json:"pdf_url_expires_at" example:"2025-01-29T11:15:00Z"`
	HasSummary      bool           `json:"has_summary" example:"true"`
	HasQuiz         bool           `json:"has_quiz" example:"true"`
} // @name MaterialDetailResponse

// MaterialDetail contiene información detallada del material
type MaterialDetail struct {
	ID          string    `json:"id" example:"uuid-1"`
	Title       string    `json:"title" example:"Introducción a Pascal"`
	Description string    `json:"description" example:"Material base sobre Pascal"`
	AuthorName  string    `json:"author_name" example:"Prof. García"`
	SubjectName string    `json:"subject_name" example:"Programación"`
	FileSize    int64     `json:"file_size" example:"2048576"`
	MyProgress  float64   `json:"my_progress" example:"0"`
	CreatedAt   time.Time `json:"created_at" example:"2025-01-15T12:00:00Z"`
} // @name MaterialDetail

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Message string `json:"message,omitempty" example:"Field 'title' is required"`
} // @name ErrorResponse

// CreateMaterialResponse representa la respuesta al crear un material
type CreateMaterialResponse struct {
	Status             string    `json:"status" example:"created"`
	MaterialID         string    `json:"material_id" example:"uuid-123"`
	UploadURL          string    `json:"upload_url" example:"https://s3.amazonaws.com/edugo-materials-prod/upload-url-mock"`
	UploadURLExpiresAt time.Time `json:"upload_url_expires_at" example:"2025-01-29T12:15:00Z"`
	MaxFileSizeBytes   int64     `json:"max_file_size_bytes" example:"104857600"`
} // @name CreateMaterialResponse

// UploadCompleteResponse representa la respuesta al completar un upload
type UploadCompleteResponse struct {
	Status     string `json:"status" example:"processing"`
	MaterialID string `json:"material_id" example:"uuid-123"`
	Message    string `json:"message" example:"Material queued for processing"`
} // @name UploadCompleteResponse

// SummarySection representa una sección del resumen
type SummarySection struct {
	Title                string `json:"title" example:"Contexto Histórico"`
	Content              string `json:"content" example:"Pascal fue desarrollado..."`
	Difficulty           string `json:"difficulty" example:"basic" enums:"basic,medium,advanced"`
	EstimatedTimeMinutes int    `json:"estimated_time_minutes,omitempty" example:"5"`
	Order                int    `json:"order" example:"1"`
} // @name SummarySection

// GlossaryTerm representa un término del glosario
type GlossaryTerm struct {
	Term       string `json:"term" example:"Compilador"`
	Definition string `json:"definition" example:"Programa que traduce código fuente a código máquina"`
	Order      int    `json:"order" example:"1"`
} // @name GlossaryTerm

// ProcessingMetadata contiene metadata del procesamiento NLP
type ProcessingMetadata struct {
	NLPProvider           string `json:"nlp_provider,omitempty" example:"openai"`
	Model                 string `json:"model,omitempty" example:"gpt-4"`
	TokensUsed            int    `json:"tokens_used,omitempty" example:"3500"`
	ProcessingTimeSeconds int    `json:"processing_time_seconds,omitempty" example:"45"`
	Language              string `json:"language,omitempty" example:"es"`
	PromptVersion         string `json:"prompt_version,omitempty" example:"v1.2"`
} // @name ProcessingMetadata

// MaterialSummaryResponse representa la respuesta con el resumen del material (estructura completa)
type MaterialSummaryResponse struct {
	Sections            []SummarySection   `json:"sections"`
	Glossary            []GlossaryTerm     `json:"glossary,omitempty"`
	ReflectionQuestions []string           `json:"reflection_questions,omitempty" example:"¿Cómo aplicarías este concepto?,¿Qué diferencias encuentras?"`
	ProcessingMetadata  ProcessingMetadata `json:"processing_metadata,omitempty"`
} // @name MaterialSummaryResponse

// QuestionOption representa una opción de respuesta
type QuestionOption struct {
	ID   string `json:"id" example:"a"`
	Text string `json:"text" example:"Un programa que traduce código"`
} // @name QuestionOption

// AssessmentResponse representa la respuesta con el quiz del material (estructura completa)
type AssessmentResponse struct {
	Title            string     `json:"title" example:"Cuestionario: Introducción a Pascal"`
	Description      string     `json:"description,omitempty" example:"Evaluación de conceptos básicos"`
	Questions        []Question `json:"questions"`
	TotalQuestions   int        `json:"total_questions" example:"5"`
	TotalPoints      int        `json:"total_points" example:"100"`
	PassingScore     int        `json:"passing_score" example:"70"`
	TimeLimitMinutes int        `json:"time_limit_minutes,omitempty" example:"15"`
} // @name AssessmentResponse

// Question representa una pregunta del quiz (estructura completa)
type Question struct {
	ID         string           `json:"id" example:"q1"`
	Text       string           `json:"text" example:"¿Qué es un compilador?"`
	Type       string           `json:"type" example:"multiple_choice" enums:"multiple_choice,true_false,short_answer"`
	Difficulty string           `json:"difficulty,omitempty" example:"basic" enums:"basic,medium,advanced"`
	Points     int              `json:"points" example:"20"`
	Order      int              `json:"order" example:"1"`
	Options    []QuestionOption `json:"options"`
} // @name Question

// QuestionFeedback representa el feedback detallado de una pregunta
type QuestionFeedback struct {
	QuestionID      string `json:"question_id" example:"q1"`
	IsCorrect       bool   `json:"is_correct" example:"true"`
	YourAnswer      string `json:"your_answer" example:"a"`
	FeedbackMessage string `json:"feedback_message" example:"¡Correcto! Un compilador traduce código fuente a código máquina."`
} // @name QuestionFeedback

// AttemptResultResponse representa la respuesta al registrar un intento (con feedback detallado)
type AttemptResultResponse struct {
	Score            float64            `json:"score" example:"85.5"`
	TotalPoints      float64            `json:"total_points" example:"100"`
	Passed           bool               `json:"passed" example:"true"`
	DetailedFeedback []QuestionFeedback `json:"detailed_feedback"`
} // @name AttemptResultResponse

// SuccessResponse representa una respuesta genérica de éxito
type SuccessResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message,omitempty" example:"Operation completed successfully"`
} // @name SuccessResponse

// MaterialStatsResponse representa las estadísticas de un material
type MaterialStatsResponse struct {
	TotalViews      int     `json:"total_views" example:"150"`
	TotalCompleted  int     `json:"total_completed" example:"45"`
	AverageProgress float64 `json:"average_progress" example:"65.5"`
	AverageScore    float64 `json:"average_score" example:"82.3"`
} // @name MaterialStatsResponse
