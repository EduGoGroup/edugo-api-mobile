package request

// CreateMaterialRequest representa la petición para crear un nuevo material
type CreateMaterialRequest struct {
	Title       string `json:"title" binding:"required" example:"Introducción a Pascal"`
	Description string `json:"description" binding:"required" example:"Material base sobre programación en Pascal"`
	SubjectID   string `json:"subject_id" binding:"required" example:"subject-uuid-123"`
	UnitID      string `json:"unit_id" binding:"required" example:"unit-uuid-456"`
	FileName    string `json:"file_name" binding:"required" example:"introduccion-pascal.pdf"`
	FileSize    int64  `json:"file_size" binding:"required" example:"2048576"`
} // @name CreateMaterialRequest

// UploadCompleteRequest representa la petición para notificar que el upload se completó
type UploadCompleteRequest struct {
	VersionKey string `json:"version_key" binding:"required" example:"materials/uuid-1/v1.pdf"`
	FileSize   int64  `json:"file_size" binding:"required" example:"2048576"`
	Checksum   string `json:"checksum" binding:"required" example:"sha256:abc123..."`
} // @name UploadCompleteRequest

// RecordAttemptRequest representa la petición para registrar un intento de quiz
type RecordAttemptRequest struct {
	Answers []Answer `json:"answers" binding:"required"`
} // @name RecordAttemptRequest

// Answer representa una respuesta individual del quiz
type Answer struct {
	QuestionID string `json:"question_id" binding:"required" example:"q-uuid-1"`
	Answer     string `json:"answer" binding:"required" example:"option_a"`
} // @name Answer

// UpdateProgressRequest representa la petición para actualizar el progreso
type UpdateProgressRequest struct {
	Progress float64 `json:"progress" binding:"required,min=0,max=100" example:"75.5"`
} // @name UpdateProgressRequest
