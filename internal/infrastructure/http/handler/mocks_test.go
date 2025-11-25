package handler

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/application/service"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/google/uuid"
)

// MockMaterialService para tests de material_handler
type MockMaterialService struct {
	CreateMaterialFunc          func(ctx context.Context, req dto.CreateMaterialRequest, authorID string) (*dto.MaterialResponse, error)
	GetMaterialFunc             func(ctx context.Context, id string) (*dto.MaterialResponse, error)
	GetMaterialWithVersionsFunc func(ctx context.Context, id string) (*dto.MaterialWithVersionsResponse, error)
	ListMaterialsFunc           func(ctx context.Context, filters repository.ListFilters) ([]*dto.MaterialResponse, error)
	NotifyUploadCompleteFunc    func(ctx context.Context, id string, req dto.UploadCompleteRequest) error
}

func (m *MockMaterialService) CreateMaterial(ctx context.Context, req dto.CreateMaterialRequest, authorID string) (*dto.MaterialResponse, error) {
	if m.CreateMaterialFunc != nil {
		return m.CreateMaterialFunc(ctx, req, authorID)
	}
	return &dto.MaterialResponse{ID: "test-id"}, nil
}

func (m *MockMaterialService) GetMaterial(ctx context.Context, id string) (*dto.MaterialResponse, error) {
	if m.GetMaterialFunc != nil {
		return m.GetMaterialFunc(ctx, id)
	}
	return &dto.MaterialResponse{ID: id}, nil
}

func (m *MockMaterialService) ListMaterials(ctx context.Context, filters repository.ListFilters) ([]*dto.MaterialResponse, error) {
	if m.ListMaterialsFunc != nil {
		return m.ListMaterialsFunc(ctx, filters)
	}
	return []*dto.MaterialResponse{}, nil
}

func (m *MockMaterialService) GetMaterialWithVersions(ctx context.Context, id string) (*dto.MaterialWithVersionsResponse, error) {
	if m.GetMaterialWithVersionsFunc != nil {
		return m.GetMaterialWithVersionsFunc(ctx, id)
	}
	return &dto.MaterialWithVersionsResponse{
		Material: &dto.MaterialResponse{ID: id},
		Versions: []*dto.MaterialVersionResponse{},
	}, nil
}

func (m *MockMaterialService) NotifyUploadComplete(ctx context.Context, id string, req dto.UploadCompleteRequest) error {
	if m.NotifyUploadCompleteFunc != nil {
		return m.NotifyUploadCompleteFunc(ctx, id, req)
	}
	return nil
}

// MockS3Storage para tests de S3 (implementa s3.S3Storage interface)
type MockS3Storage struct {
	GeneratePresignedUploadURLFunc   func(ctx context.Context, key, contentType string, expires time.Duration) (string, error)
	GeneratePresignedDownloadURLFunc func(ctx context.Context, key string, expires time.Duration) (string, error)
}

func (m *MockS3Storage) GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expires time.Duration) (string, error) {
	if m.GeneratePresignedUploadURLFunc != nil {
		return m.GeneratePresignedUploadURLFunc(ctx, key, contentType, expires)
	}
	return "https://mock-s3-url.com/presigned-upload", nil
}

func (m *MockS3Storage) GeneratePresignedDownloadURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	if m.GeneratePresignedDownloadURLFunc != nil {
		return m.GeneratePresignedDownloadURLFunc(ctx, key, expires)
	}
	return "https://mock-s3-url.com/presigned-download", nil
}

// MockAssessmentService para tests de assessment_handler
type MockAssessmentService struct {
	GetAssessmentFunc  func(ctx context.Context, materialID string) (*repository.MaterialAssessment, error)
	RecordAttemptFunc  func(ctx context.Context, materialID string, userID string, answers map[string]interface{}) (*repository.AssessmentAttempt, error)
	CalculateScoreFunc func(ctx context.Context, assessmentID string, userID string, userResponses map[string]interface{}) (*repository.AssessmentResult, error)
}

func (m *MockAssessmentService) GetAssessment(ctx context.Context, materialID string) (*repository.MaterialAssessment, error) {
	if m.GetAssessmentFunc != nil {
		return m.GetAssessmentFunc(ctx, materialID)
	}
	return &repository.MaterialAssessment{}, nil
}

func (m *MockAssessmentService) RecordAttempt(ctx context.Context, materialID string, userID string, answers map[string]interface{}) (*repository.AssessmentAttempt, error) {
	if m.RecordAttemptFunc != nil {
		return m.RecordAttemptFunc(ctx, materialID, userID, answers)
	}
	return &repository.AssessmentAttempt{}, nil
}

func (m *MockAssessmentService) CalculateScore(ctx context.Context, assessmentID string, userID string, userResponses map[string]interface{}) (*repository.AssessmentResult, error) {
	if m.CalculateScoreFunc != nil {
		return m.CalculateScoreFunc(ctx, assessmentID, userID, userResponses)
	}
	return &repository.AssessmentResult{
		ID:             "result-123",
		AssessmentID:   assessmentID,
		Score:          75.0,
		TotalQuestions: 4,
		CorrectAnswers: 3,
		Feedback:       []repository.FeedbackItem{},
	}, nil
}

// MockProgressService para tests de progress_handler
type MockProgressService struct {
	UpdateProgressFunc func(ctx context.Context, materialID, userID string, percentage, lastPage int) error
}

func (m *MockProgressService) UpdateProgress(ctx context.Context, materialID, userID string, percentage, lastPage int) error {
	if m.UpdateProgressFunc != nil {
		return m.UpdateProgressFunc(ctx, materialID, userID, percentage, lastPage)
	}
	return nil
}

// MockStatsService para tests de stats_handler
type MockStatsService struct {
	GetMaterialStatsFunc func(ctx context.Context, materialID string) (*service.MaterialStats, error)
	GetGlobalStatsFunc   func(ctx context.Context) (*dto.GlobalStatsDTO, error)
}

func (m *MockStatsService) GetMaterialStats(ctx context.Context, materialID string) (*service.MaterialStats, error) {
	if m.GetMaterialStatsFunc != nil {
		return m.GetMaterialStatsFunc(ctx, materialID)
	}
	return &service.MaterialStats{
		TotalViews:    150,
		AvgProgress:   67.5,
		TotalAttempts: 45,
		AvgScore:      78.3,
	}, nil
}

func (m *MockStatsService) GetGlobalStats(ctx context.Context) (*dto.GlobalStatsDTO, error) {
	if m.GetGlobalStatsFunc != nil {
		return m.GetGlobalStatsFunc(ctx)
	}
	return &dto.GlobalStatsDTO{
		TotalPublishedMaterials:   150,
		TotalCompletedAssessments: 1250,
		AverageAssessmentScore:    78.5,
		ActiveUsersLast30Days:     320,
		AverageProgress:           65.3,
	}, nil
}

// MockSummaryService para tests de summary_handler
type MockSummaryService struct {
	GetSummaryFunc func(ctx context.Context, materialID string) (*repository.MaterialSummary, error)
}

func (m *MockSummaryService) GetSummary(ctx context.Context, materialID string) (*repository.MaterialSummary, error) {
	if m.GetSummaryFunc != nil {
		return m.GetSummaryFunc(ctx, materialID)
	}
	return &repository.MaterialSummary{
		MainIdeas:   []string{"Idea principal 1", "Idea principal 2"},
		KeyConcepts: map[string]string{"concepto1": "definición1"},
		Sections: []repository.SummarySection{
			{Title: "Introducción", Content: "Contenido de introducción", Page: 1},
		},
		Glossary:  map[string]string{"término1": "definición1"},
		CreatedAt: "2024-01-15T10:30:00Z",
	}, nil
}

// MockAssessmentAttemptService para tests de assessment_handler (SPRINT-04)
type MockAssessmentAttemptService struct {
	GetAssessmentByMaterialIDFunc func(ctx context.Context, materialID uuid.UUID) (*dto.AssessmentResponse, error)
	CreateAttemptFunc             func(ctx context.Context, studentID, materialID uuid.UUID, req dto.CreateAttemptRequest) (*dto.AttemptResultResponse, error)
	GetAttemptResultFunc          func(ctx context.Context, attemptID, studentID uuid.UUID) (*dto.AttemptResultResponse, error)
	GetAttemptHistoryFunc         func(ctx context.Context, studentID uuid.UUID, limit, offset int) (*dto.AttemptHistoryResponse, error)
}

func (m *MockAssessmentAttemptService) GetAssessmentByMaterialID(ctx context.Context, materialID uuid.UUID) (*dto.AssessmentResponse, error) {
	if m.GetAssessmentByMaterialIDFunc != nil {
		return m.GetAssessmentByMaterialIDFunc(ctx, materialID)
	}
	return &dto.AssessmentResponse{}, nil
}

func (m *MockAssessmentAttemptService) CreateAttempt(ctx context.Context, studentID, materialID uuid.UUID, req dto.CreateAttemptRequest) (*dto.AttemptResultResponse, error) {
	if m.CreateAttemptFunc != nil {
		return m.CreateAttemptFunc(ctx, studentID, materialID, req)
	}
	return &dto.AttemptResultResponse{}, nil
}

func (m *MockAssessmentAttemptService) GetAttemptResult(ctx context.Context, attemptID, studentID uuid.UUID) (*dto.AttemptResultResponse, error) {
	if m.GetAttemptResultFunc != nil {
		return m.GetAttemptResultFunc(ctx, attemptID, studentID)
	}
	return &dto.AttemptResultResponse{}, nil
}

func (m *MockAssessmentAttemptService) GetAttemptHistory(ctx context.Context, studentID uuid.UUID, limit, offset int) (*dto.AttemptHistoryResponse, error) {
	if m.GetAttemptHistoryFunc != nil {
		return m.GetAttemptHistoryFunc(ctx, studentID, limit, offset)
	}
	return &dto.AttemptHistoryResponse{}, nil
}
