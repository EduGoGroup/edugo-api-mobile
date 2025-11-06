package handler

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
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

// MockAuthService para tests de auth_handler
type MockAuthService struct {
	LoginFunc              func(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshAccessTokenFunc func(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
	LogoutFunc             func(ctx context.Context, userID, refreshToken string) error
	RevokeAllSessionsFunc  func(ctx context.Context, userID string) error
}

func (m *MockAuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, req)
	}
	return &dto.LoginResponse{
		AccessToken:  "mock-access-token",
		RefreshToken: "mock-refresh-token",
		User: dto.UserInfo{
			ID:        "user-123",
			Email:     req.Email,
			FirstName: "Test",
			LastName:  "User",
			FullName:  "Test User",
		},
	}, nil
}

func (m *MockAuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	if m.RefreshAccessTokenFunc != nil {
		return m.RefreshAccessTokenFunc(ctx, refreshToken)
	}
	return &dto.RefreshResponse{
		AccessToken: "new-access-token",
	}, nil
}

func (m *MockAuthService) Logout(ctx context.Context, userID, refreshToken string) error {
	if m.LogoutFunc != nil {
		return m.LogoutFunc(ctx, userID, refreshToken)
	}
	return nil
}

func (m *MockAuthService) RevokeAllSessions(ctx context.Context, userID string) error {
	if m.RevokeAllSessionsFunc != nil {
		return m.RevokeAllSessionsFunc(ctx, userID)
	}
	return nil
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
