package mongodb

import (
	"context"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	mongoRepo "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/mongodb/repository"
)

// Este archivo contiene stubs de repositorios MongoDB para desarrollo y testing

// Summary Repository Stub
type mockSummaryRepository struct{}

func NewMockSummaryRepository() repository.SummaryRepository { return &mockSummaryRepository{} }
func (r *mockSummaryRepository) FindByMaterialID(ctx context.Context, materialID valueobject.MaterialID) (*repository.MaterialSummary, error) {
	return nil, nil
}
func (r *mockSummaryRepository) Exists(ctx context.Context, materialID valueobject.MaterialID) (bool, error) {
	return false, nil
}
func (r *mockSummaryRepository) List(ctx context.Context, limit, offset int) ([]*repository.MaterialSummary, error) {
	return []*repository.MaterialSummary{}, nil
}
func (r *mockSummaryRepository) Save(ctx context.Context, summary *repository.MaterialSummary) error {
	return nil
}
func (r *mockSummaryRepository) Delete(ctx context.Context, materialID valueobject.MaterialID) error {
	return nil
}

// Assessment Document Repository Stub
type mockAssessmentDocumentRepository struct{}

func NewMockAssessmentDocumentRepository() mongoRepo.AssessmentDocumentRepository {
	return &mockAssessmentDocumentRepository{}
}
func (r *mockAssessmentDocumentRepository) FindByMaterialID(ctx context.Context, materialID string) (*mongoRepo.AssessmentDocument, error) {
	return nil, nil
}
func (r *mockAssessmentDocumentRepository) FindByID(ctx context.Context, objectID string) (*mongoRepo.AssessmentDocument, error) {
	return nil, nil
}
func (r *mockAssessmentDocumentRepository) Save(ctx context.Context, doc *mongoRepo.AssessmentDocument) error {
	return nil
}
func (r *mockAssessmentDocumentRepository) Delete(ctx context.Context, id string) error { return nil }
