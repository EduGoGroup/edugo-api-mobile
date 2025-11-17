//go:build integration
// +build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	testifySuite "github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/persistence/mongodb/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/testing/suite"
)

// AssessmentDocumentRepositoryIntegrationSuite tests de integración para MongoDB
type AssessmentDocumentRepositoryIntegrationSuite struct {
	suite.IntegrationTestSuite
	repo     repository.AssessmentDocumentRepository
	repoImpl *repository.MongoAssessmentDocumentRepository // Para métodos helper no públicos
}

// SetupSuite se ejecuta UNA VEZ antes de todos los tests
func (s *AssessmentDocumentRepositoryIntegrationSuite) SetupSuite() {
	s.IntegrationTestSuite.SetupSuite()

	// Conectar a MongoDB (si la suite padre no lo hizo)
	if s.MongoDB == nil {
		ctx := context.Background()
		connStr, err := s.MongoContainer.ConnectionString(ctx)
		s.Require().NoError(err)

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
		s.Require().NoError(err)

		s.MongoDB = client.Database("edugo_test")
	}
}

// SetupTest prepara cada test individual
func (s *AssessmentDocumentRepositoryIntegrationSuite) SetupTest() {
	s.IntegrationTestSuite.SetupTest()
	concreteRepo := repository.NewMongoAssessmentDocumentRepository(s.MongoDB)
	s.repo = concreteRepo
	s.repoImpl = concreteRepo.(*repository.MongoAssessmentDocumentRepository)

	// Limpiar colección MongoDB antes de cada test
	ctx := context.Background()
	_ = s.MongoDB.Collection("material_assessment").Drop(ctx)
}

// TestAssessmentDocumentRepositoryIntegration ejecuta la suite
func TestAssessmentDocumentRepositoryIntegration(t *testing.T) {
	testifySuite.Run(t, new(AssessmentDocumentRepositoryIntegrationSuite))
}

// TestSave_Insert valida que Save inserta un nuevo documento
func (s *AssessmentDocumentRepositoryIntegrationSuite) TestSave_Insert() {
	ctx := context.Background()

	// Arrange
	materialID := uuid.New().String()
	doc := &repository.AssessmentDocument{
		MaterialID: materialID,
		Title:      "Test Assessment",
		Questions: []repository.Question{
			{
				ID:            "q1",
				Text:          "¿Qué es Go?",
				Type:          "multiple_choice",
				Options:       []repository.Option{{ID: "a", Text: "Un lenguaje"}},
				CorrectAnswer: "a",
				Feedback:      repository.Feedback{Correct: "Correcto!", Incorrect: "Incorrecto"},
			},
		},
		Metadata: repository.Metadata{
			GeneratedBy:          "test",
			GenerationDate:       time.Now().UTC(),
			PromptVersion:        "v1",
			TotalQuestions:       1,
			EstimatedTimeMinutes: 5,
		},
		Version:   1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Act
	err := s.repo.Save(ctx, doc)

	// Assert
	s.NoError(err, "Save debe insertar sin errores")
	s.False(doc.ID.IsZero(), "Debe asignar ObjectID")

	// Verificar que se insertó
	found, err := s.repo.FindByMaterialID(ctx, materialID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(materialID, found.MaterialID)
	s.Equal("Test Assessment", found.Title)
	s.Equal(1, len(found.Questions))
}

// TestSave_Upsert valida que Save actualiza documento existente (UPSERT por material_id)
func (s *AssessmentDocumentRepositoryIntegrationSuite) TestSave_Upsert() {
	ctx := context.Background()

	// Arrange - Insertar documento inicial
	materialID := uuid.New().String()
	doc := &repository.AssessmentDocument{
		MaterialID: materialID,
		Title:      "Original Title",
		Questions: []repository.Question{
			{ID: "q1", Text: "Original question", Type: "multiple_choice", Options: []repository.Option{{ID: "a", Text: "A"}}, CorrectAnswer: "a", Feedback: repository.Feedback{Correct: "OK", Incorrect: "NO"}},
		},
		Metadata:  repository.Metadata{GeneratedBy: "test", GenerationDate: time.Now().UTC(), PromptVersion: "v1", TotalQuestions: 1, EstimatedTimeMinutes: 5},
		Version:   1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := s.repo.Save(ctx, doc)
	s.Require().NoError(err)

	originalID := doc.ID

	// Act - Actualizar el mismo material_id (UPSERT)
	updatedDoc := &repository.AssessmentDocument{
		MaterialID: materialID,
		Title:      "Updated Title",
		Questions: []repository.Question{
			{ID: "q1", Text: "Updated question", Type: "multiple_choice", Options: []repository.Option{{ID: "a", Text: "A"}}, CorrectAnswer: "a", Feedback: repository.Feedback{Correct: "OK", Incorrect: "NO"}},
			{ID: "q2", Text: "New question", Type: "multiple_choice", Options: []repository.Option{{ID: "b", Text: "B"}}, CorrectAnswer: "b", Feedback: repository.Feedback{Correct: "OK", Incorrect: "NO"}},
		},
		Metadata:  repository.Metadata{GeneratedBy: "test", GenerationDate: time.Now().UTC(), PromptVersion: "v2", TotalQuestions: 2, EstimatedTimeMinutes: 10},
		Version:   2,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = s.repo.Save(ctx, updatedDoc)

	// Assert
	s.NoError(err, "Save debe actualizar sin errores")

	// Verificar que se actualizó (mismo ObjectID)
	found, err := s.repo.FindByMaterialID(ctx, materialID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal("Updated Title", found.Title)
	s.Equal(2, len(found.Questions), "Debe tener 2 preguntas después del update")
	s.Equal(2, found.Version)

	// Verificar que NO se creó documento duplicado
	foundByID, err := s.repo.FindByID(ctx, originalID.Hex())
	s.NoError(err)
	s.NotNil(foundByID)
	s.Equal("Updated Title", foundByID.Title, "El documento original debe haberse actualizado")
}

// TestFindByMaterialID_NotFound valida que retorna nil cuando no encuentra
func (s *AssessmentDocumentRepositoryIntegrationSuite) TestFindByMaterialID_NotFound() {
	ctx := context.Background()

	// Act
	found, err := s.repo.FindByMaterialID(ctx, uuid.New().String())

	// Assert
	s.NoError(err)
	s.Nil(found, "Debe retornar nil cuando no encuentra")
}

// TestFindByID_Success valida que FindByID encuentra por ObjectID
func (s *AssessmentDocumentRepositoryIntegrationSuite) TestFindByID_Success() {
	ctx := context.Background()

	// Arrange - Insertar documento
	materialID := uuid.New().String()
	doc := &repository.AssessmentDocument{
		MaterialID: materialID,
		Title:      "Test",
		Questions:  []repository.Question{{ID: "q1", Text: "Q", Type: "multiple_choice", Options: []repository.Option{{ID: "a", Text: "A"}}, CorrectAnswer: "a", Feedback: repository.Feedback{Correct: "OK", Incorrect: "NO"}}},
		Metadata:   repository.Metadata{GeneratedBy: "test", GenerationDate: time.Now().UTC(), PromptVersion: "v1", TotalQuestions: 1, EstimatedTimeMinutes: 5},
		Version:    1,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	err := s.repo.Save(ctx, doc)
	s.Require().NoError(err)

	objectID := doc.ID.Hex()

	// Act
	found, err := s.repo.FindByID(ctx, objectID)

	// Assert
	s.NoError(err)
	s.NotNil(found)
	s.Equal(doc.ID, found.ID)
	s.Equal(materialID, found.MaterialID)
}

// TestFindByID_InvalidObjectID valida que maneja ObjectID inválido
func (s *AssessmentDocumentRepositoryIntegrationSuite) TestFindByID_InvalidObjectID() {
	ctx := context.Background()

	// Act
	found, err := s.repo.FindByID(ctx, "invalid-object-id")

	// Assert
	s.Error(err, "Debe fallar con ObjectID inválido")
	s.Nil(found)
}

// TestDelete_Success valida que Delete elimina correctamente
func (s *AssessmentDocumentRepositoryIntegrationSuite) TestDelete_Success() {
	ctx := context.Background()

	// Arrange - Insertar documento
	materialID := uuid.New().String()
	doc := &repository.AssessmentDocument{
		MaterialID: materialID,
		Title:      "To Delete",
		Questions:  []repository.Question{{ID: "q1", Text: "Q", Type: "multiple_choice", Options: []repository.Option{{ID: "a", Text: "A"}}, CorrectAnswer: "a", Feedback: repository.Feedback{Correct: "OK", Incorrect: "NO"}}},
		Metadata:   repository.Metadata{GeneratedBy: "test", GenerationDate: time.Now().UTC(), PromptVersion: "v1", TotalQuestions: 1, EstimatedTimeMinutes: 5},
		Version:    1,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	err := s.repo.Save(ctx, doc)
	s.Require().NoError(err)

	objectID := doc.ID.Hex()

	// Act
	err = s.repo.Delete(ctx, objectID)

	// Assert
	s.NoError(err)

	// Verificar que se eliminó
	found, err := s.repo.FindByID(ctx, objectID)
	s.NoError(err)
	s.Nil(found, "Documento debe estar eliminado")
}

// TestDelete_NotFound valida que Delete falla cuando no encuentra
func (s *AssessmentDocumentRepositoryIntegrationSuite) TestDelete_NotFound() {
	ctx := context.Background()

	// Act - Intentar eliminar documento que no existe
	newObjectID := primitive.NewObjectID()
	err := s.repo.Delete(ctx, newObjectID.Hex())

	// Assert
	s.Error(err, "Delete debe fallar cuando no encuentra el documento")
}

// TestGetQuestionByID_Success valida que encuentra pregunta con $elemMatch
func (s *AssessmentDocumentRepositoryIntegrationSuite) TestGetQuestionByID_Success() {
	ctx := context.Background()

	// Arrange - Insertar documento con múltiples preguntas
	materialID := uuid.New().String()
	doc := &repository.AssessmentDocument{
		MaterialID: materialID,
		Title:      "Multi-question Assessment",
		Questions: []repository.Question{
			{ID: "q1", Text: "Question 1", Type: "multiple_choice", Options: []repository.Option{{ID: "a", Text: "A"}}, CorrectAnswer: "a", Feedback: repository.Feedback{Correct: "OK1", Incorrect: "NO1"}},
			{ID: "q2", Text: "Question 2", Type: "multiple_choice", Options: []repository.Option{{ID: "b", Text: "B"}}, CorrectAnswer: "b", Feedback: repository.Feedback{Correct: "OK2", Incorrect: "NO2"}},
			{ID: "q3", Text: "Question 3", Type: "multiple_choice", Options: []repository.Option{{ID: "c", Text: "C"}}, CorrectAnswer: "c", Feedback: repository.Feedback{Correct: "OK3", Incorrect: "NO3"}},
		},
		Metadata:  repository.Metadata{GeneratedBy: "test", GenerationDate: time.Now().UTC(), PromptVersion: "v1", TotalQuestions: 3, EstimatedTimeMinutes: 15},
		Version:   1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := s.repo.Save(ctx, doc)
	s.Require().NoError(err)

	// Act
	question, err := s.repoImpl.GetQuestionByID(ctx, materialID, "q2")

	// Assert
	s.NoError(err)
	s.NotNil(question)
	s.Equal("q2", question.ID)
	s.Equal("Question 2", question.Text)
	s.Equal("OK2", question.Feedback.Correct)
}

// TestGetQuestionByID_NotFound valida que falla cuando no encuentra la pregunta
func (s *AssessmentDocumentRepositoryIntegrationSuite) TestGetQuestionByID_NotFound() {
	ctx := context.Background()

	// Arrange - Insertar documento
	materialID := uuid.New().String()
	doc := &repository.AssessmentDocument{
		MaterialID: materialID,
		Title:      "Test",
		Questions:  []repository.Question{{ID: "q1", Text: "Q", Type: "multiple_choice", Options: []repository.Option{{ID: "a", Text: "A"}}, CorrectAnswer: "a", Feedback: repository.Feedback{Correct: "OK", Incorrect: "NO"}}},
		Metadata:   repository.Metadata{GeneratedBy: "test", GenerationDate: time.Now().UTC(), PromptVersion: "v1", TotalQuestions: 1, EstimatedTimeMinutes: 5},
		Version:    1,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	err := s.repo.Save(ctx, doc)
	s.Require().NoError(err)

	// Act
	question, err := s.repoImpl.GetQuestionByID(ctx, materialID, "nonexistent")

	// Assert
	s.Error(err, "Debe fallar cuando no encuentra la pregunta")
	s.Nil(question)
}

// TestSave_ValidationErrors valida que Save valida campos requeridos
func (s *AssessmentDocumentRepositoryIntegrationSuite) TestSave_ValidationErrors() {
	ctx := context.Background()

	// Test 1: Sin material_id
	doc1 := &repository.AssessmentDocument{
		Title:     "Test",
		Questions: []repository.Question{{ID: "q1", Text: "Q", Type: "multiple_choice", Options: []repository.Option{{ID: "a", Text: "A"}}, CorrectAnswer: "a", Feedback: repository.Feedback{Correct: "OK", Incorrect: "NO"}}},
	}
	err := s.repo.Save(ctx, doc1)
	s.Error(err, "Debe fallar sin material_id")

	// Test 2: Sin título
	doc2 := &repository.AssessmentDocument{
		MaterialID: uuid.New().String(),
		Questions:  []repository.Question{{ID: "q1", Text: "Q", Type: "multiple_choice", Options: []repository.Option{{ID: "a", Text: "A"}}, CorrectAnswer: "a", Feedback: repository.Feedback{Correct: "OK", Incorrect: "NO"}}},
	}
	err = s.repo.Save(ctx, doc2)
	s.Error(err, "Debe fallar sin título")

	// Test 3: Sin preguntas
	doc3 := &repository.AssessmentDocument{
		MaterialID: uuid.New().String(),
		Title:      "Test",
		Questions:  []repository.Question{},
	}
	err = s.repo.Save(ctx, doc3)
	s.Error(err, "Debe fallar sin preguntas")
}
