package postgres

import (
	"context"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repositories"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/google/uuid"
)

// Material Repository Stub
type mockMaterialRepository struct{}

func NewMockMaterialRepository() repository.MaterialRepository { return &mockMaterialRepository{} }
func (r *mockMaterialRepository) FindByID(ctx context.Context, id valueobject.MaterialID) (*pgentities.Material, error) {
	return nil, nil
}
func (r *mockMaterialRepository) FindByIDWithVersions(ctx context.Context, id valueobject.MaterialID) (*pgentities.Material, []*pgentities.MaterialVersion, error) {
	return nil, nil, nil
}
func (r *mockMaterialRepository) List(ctx context.Context, filters repository.ListFilters) ([]*pgentities.Material, error) {
	return []*pgentities.Material{}, nil
}
func (r *mockMaterialRepository) FindByAuthor(ctx context.Context, authorID valueobject.UserID) ([]*pgentities.Material, error) {
	return []*pgentities.Material{}, nil
}
func (r *mockMaterialRepository) Create(ctx context.Context, material *pgentities.Material) error {
	return nil
}
func (r *mockMaterialRepository) Update(ctx context.Context, material *pgentities.Material) error {
	return nil
}
func (r *mockMaterialRepository) UpdateStatus(ctx context.Context, id valueobject.MaterialID, status enum.MaterialStatus) error {
	return nil
}
func (r *mockMaterialRepository) UpdateProcessingStatus(ctx context.Context, id valueobject.MaterialID, status enum.ProcessingStatus) error {
	return nil
}
func (r *mockMaterialRepository) CountPublishedMaterials(ctx context.Context) (int64, error) {
	return 0, nil
}

// Progress Repository Stub
type mockProgressRepository struct{}

func NewMockProgressRepository() repository.ProgressRepository { return &mockProgressRepository{} }
func (r *mockProgressRepository) FindByUserAndMaterial(ctx context.Context, userID valueobject.UserID, materialID valueobject.MaterialID) (*pgentities.Progress, error) {
	return nil, nil
}
func (r *mockProgressRepository) FindByMaterialAndUser(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) (*pgentities.Progress, error) {
	return nil, nil
}
func (r *mockProgressRepository) FindByUser(ctx context.Context, userID valueobject.UserID) ([]*pgentities.Progress, error) {
	return []*pgentities.Progress{}, nil
}
func (r *mockProgressRepository) Create(ctx context.Context, progress *pgentities.Progress) error {
	return nil
}
func (r *mockProgressRepository) Save(ctx context.Context, progress *pgentities.Progress) error {
	return nil
}
func (r *mockProgressRepository) Update(ctx context.Context, progress *pgentities.Progress) error {
	return nil
}
func (r *mockProgressRepository) Upsert(ctx context.Context, progress *pgentities.Progress) (*pgentities.Progress, error) {
	return progress, nil
}
func (r *mockProgressRepository) CalculateAverageProgress(ctx context.Context) (float64, error) {
	return 0.0, nil
}
func (r *mockProgressRepository) CountActiveUsers(ctx context.Context) (int64, error) { return 0, nil }

// RefreshToken Repository Stub
type mockRefreshTokenRepository struct{}

func NewMockRefreshTokenRepository() repository.RefreshTokenRepository {
	return &mockRefreshTokenRepository{}
}
func (r *mockRefreshTokenRepository) Store(ctx context.Context, token repository.RefreshTokenData) error {
	return nil
}
func (r *mockRefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*repository.RefreshTokenData, error) {
	return nil, nil
}
func (r *mockRefreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error { return nil }
func (r *mockRefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	return nil
}
func (r *mockRefreshTokenRepository) DeleteExpired(ctx context.Context) (int64, error) { return 0, nil }

// LoginAttempt Repository Stub
type mockLoginAttemptRepository struct{}

func NewMockLoginAttemptRepository() repository.LoginAttemptRepository {
	return &mockLoginAttemptRepository{}
}
func (r *mockLoginAttemptRepository) CountFailedAttempts(ctx context.Context, identifier string, windowMinutes int) (int, error) {
	return 0, nil
}
func (r *mockLoginAttemptRepository) IsRateLimited(ctx context.Context, identifier string, maxAttempts int, windowMinutes int) (bool, error) {
	return false, nil
}
func (r *mockLoginAttemptRepository) RecordAttempt(ctx context.Context, attempt repository.LoginAttemptData) error {
	return nil
}

// Assessment Repository Stub (Sprint-03)
type mockAssessmentRepository struct{}

func NewMockAssessmentRepository() repositories.AssessmentRepository {
	return &mockAssessmentRepository{}
}
func (r *mockAssessmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*pgentities.Assessment, error) {
	return nil, nil
}
func (r *mockAssessmentRepository) FindByMaterialID(ctx context.Context, materialID uuid.UUID) (*pgentities.Assessment, error) {
	return nil, nil
}
func (r *mockAssessmentRepository) Save(ctx context.Context, assessment *pgentities.Assessment) error {
	return nil
}
func (r *mockAssessmentRepository) Delete(ctx context.Context, id uuid.UUID) error { return nil }

// Attempt Repository Stub
type mockAttemptRepository struct{}

func NewMockAttemptRepository() repositories.AttemptRepository { return &mockAttemptRepository{} }
func (r *mockAttemptRepository) FindByID(ctx context.Context, id uuid.UUID) (*pgentities.AssessmentAttempt, error) {
	return nil, nil
}
func (r *mockAttemptRepository) FindByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) ([]*pgentities.AssessmentAttempt, error) {
	return []*pgentities.AssessmentAttempt{}, nil
}
func (r *mockAttemptRepository) FindByStudent(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]*pgentities.AssessmentAttempt, error) {
	return []*pgentities.AssessmentAttempt{}, nil
}
func (r *mockAttemptRepository) Save(ctx context.Context, attempt *pgentities.AssessmentAttempt) error {
	return nil
}
func (r *mockAttemptRepository) CountByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) (int, error) {
	return 0, nil
}

// Answer Repository Stub
type mockAnswerRepository struct{}

func NewMockAnswerRepository() repositories.AnswerRepository { return &mockAnswerRepository{} }
func (r *mockAnswerRepository) FindByAttemptID(ctx context.Context, attemptID uuid.UUID) ([]*pgentities.AssessmentAttemptAnswer, error) {
	return []*pgentities.AssessmentAttemptAnswer{}, nil
}
func (r *mockAnswerRepository) Save(ctx context.Context, answers []*pgentities.AssessmentAttemptAnswer) error {
	return nil
}
func (r *mockAnswerRepository) FindByQuestionID(ctx context.Context, questionID string, limit, offset int) ([]*pgentities.AssessmentAttemptAnswer, error) {
	return []*pgentities.AssessmentAttemptAnswer{}, nil
}
