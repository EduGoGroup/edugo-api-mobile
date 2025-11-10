//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupMongoTestDB crea un testcontainer de MongoDB con las colecciones necesarias
func setupMongoTestDB(t *testing.T) (*mongo.Database, func()) {
	t.Helper()

	ctx := context.Background()

	// Levantar MongoDB testcontainer
	mongoContainer, err := mongodb.Run(ctx, "mongo:7.0",
		mongodb.WithUsername("edugo_admin"),
		mongodb.WithPassword("edugo_pass"),
	)
	require.NoError(t, err, "Failed to start MongoDB testcontainer")

	// Obtener connection string
	connStr, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	// Conectar a MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	require.NoError(t, err)

	// Retry ping to ensure database is ready
	var pingErr error
	for i := 0; i < 10; i++ {
		pingErr = client.Ping(ctx, nil)
		if pingErr == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	require.NoError(t, pingErr, "Failed to ping MongoDB after retries")

	db := client.Database("edugo")

	// Crear índice UNIQUE en assessment_results (assessment_id, user_id)
	resultsCollection := db.Collection("assessment_results")
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "assessment_id", Value: 1},
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err = resultsCollection.Indexes().CreateOne(ctx, indexModel)
	require.NoError(t, err, "Failed to create unique index on assessment_results")

	cleanup := func() {
		client.Disconnect(ctx)
		mongoContainer.Terminate(ctx)
	}

	return db, cleanup
}

// cleanMongoCollections limpia todas las colecciones de prueba
func cleanMongoCollections(t *testing.T, db *mongo.Database) {
	t.Helper()

	ctx := context.Background()
	collections := []string{
		"material_assessments",
		"assessment_attempts",
		"assessment_results",
	}

	for _, collName := range collections {
		coll := db.Collection(collName)
		_, err := coll.DeleteMany(ctx, bson.M{})
		if err != nil {
			t.Logf("Warning: Failed to clean collection %s: %v", collName, err)
		}
	}
}

func TestAssessmentRepository_SaveAssessment_ValidData(t *testing.T) {
	// Arrange
	db, cleanup := setupMongoTestDB(t)
	defer cleanup()

	repo := NewMongoAssessmentRepository(db)

	materialID, err := valueobject.MaterialIDFromString("550e8400-e29b-41d4-a716-446655440000")
	require.NoError(t, err)

	assessment := &repository.MaterialAssessment{
		MaterialID: materialID,
		Questions: []repository.AssessmentQuestion{
			{
				ID:            "q1",
				QuestionText:  "¿Qué es Go?",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				Options:       []string{"A) Un lenguaje", "B) Una base de datos", "C) Un framework"},
				CorrectAnswer: "A",
			},
			{
				ID:            "q2",
				QuestionText:  "¿Go es compilado?",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: true,
			},
		},
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	// Act
	err = repo.SaveAssessment(context.Background(), assessment)

	// Assert
	require.NoError(t, err, "SaveAssessment should not return error with valid data")

	// Verificar que se guardó en la colección
	collection := db.Collection("material_assessments")
	var doc bson.M
	err = collection.FindOne(context.Background(), bson.M{"material_id": materialID.String()}).Decode(&doc)
	require.NoError(t, err, "Assessment should be in database")
	assert.Equal(t, materialID.String(), doc["material_id"])
	assert.NotNil(t, doc["questions"])
	assert.NotNil(t, doc["created_at"])
}

func TestAssessmentRepository_SaveAssessment_UpsertBehavior(t *testing.T) {
	// Arrange
	db, cleanup := setupMongoTestDB(t)
	defer cleanup()

	repo := NewMongoAssessmentRepository(db)

	materialID, err := valueobject.MaterialIDFromString("550e8400-e29b-41d4-a716-446655440001")
	require.NoError(t, err)

	// Primer assessment con 2 preguntas
	assessment1 := &repository.MaterialAssessment{
		MaterialID: materialID,
		Questions: []repository.AssessmentQuestion{
			{
				ID:            "q1",
				QuestionText:  "Pregunta 1",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				Options:       []string{"A", "B", "C"},
				CorrectAnswer: "A",
			},
			{
				ID:            "q2",
				QuestionText:  "Pregunta 2",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: true,
			},
		},
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	// Guardar primer assessment
	err = repo.SaveAssessment(context.Background(), assessment1)
	require.NoError(t, err)

	// Segundo assessment con 3 preguntas (actualización)
	assessment2 := &repository.MaterialAssessment{
		MaterialID: materialID,
		Questions: []repository.AssessmentQuestion{
			{
				ID:            "q1",
				QuestionText:  "Pregunta 1 actualizada",
				QuestionType:  enum.AssessmentTypeMultipleChoice,
				Options:       []string{"A", "B", "C", "D"},
				CorrectAnswer: "B",
			},
			{
				ID:            "q2",
				QuestionText:  "Pregunta 2 actualizada",
				QuestionType:  enum.AssessmentTypeTrueFalse,
				CorrectAnswer: false,
			},
			{
				ID:            "q3",
				QuestionText:  "Pregunta 3 nueva",
				QuestionType:  enum.AssessmentTypeShortAnswer,
				CorrectAnswer: "respuesta correcta",
			},
		},
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	// Act - Guardar segundo assessment (debería reemplazar el primero)
	err = repo.SaveAssessment(context.Background(), assessment2)

	// Assert
	require.NoError(t, err, "SaveAssessment should not return error on upsert")

	// Verificar que solo hay un documento
	collection := db.Collection("material_assessments")
	count, err := collection.CountDocuments(context.Background(), bson.M{"material_id": materialID.String()})
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "Should have only one assessment document")

	// Verificar que tiene 3 preguntas (actualizado)
	var doc bson.M
	err = collection.FindOne(context.Background(), bson.M{"material_id": materialID.String()}).Decode(&doc)
	require.NoError(t, err)
	questions := doc["questions"].(bson.A)
	assert.Equal(t, 3, len(questions), "Should have 3 questions after update")
}

func TestAssessmentRepository_FindAssessmentByMaterialID_AssessmentExists(t *testing.T) {
	// Arrange
	db, cleanup := setupMongoTestDB(t)
	defer cleanup()

	repo := NewMongoAssessmentRepository(db)

	materialID, err := valueobject.MaterialIDFromString("550e8400-e29b-41d4-a716-446655440002")
	require.NoError(t, err)

	// Insertar assessment directamente en la colección
	collection := db.Collection("material_assessments")
	doc := bson.M{
		"material_id": materialID.String(),
		"questions": bson.A{
			bson.M{
				"id":            "q1",
				"text":          "¿Qué es Go?",
				"question_type": "multiple_choice",
				"options":       bson.A{"A) Un lenguaje", "B) Una base de datos"},
				"answer":        "A",
			},
			bson.M{
				"id":            "q2",
				"text":          "¿Go es compilado?",
				"question_type": "true_false",
				"answer":        true,
			},
		},
		"created_at": "2024-01-01T00:00:00Z",
	}
	_, err = collection.InsertOne(context.Background(), doc)
	require.NoError(t, err)

	// Act
	assessment, err := repo.FindAssessmentByMaterialID(context.Background(), materialID)

	// Assert
	require.NoError(t, err, "FindAssessmentByMaterialID should not return error when assessment exists")
	assert.NotNil(t, assessment)
	assert.Equal(t, materialID.String(), assessment.MaterialID.String())
	assert.Equal(t, 2, len(assessment.Questions))
	assert.Equal(t, "q1", assessment.Questions[0].ID)
	assert.Equal(t, "¿Qué es Go?", assessment.Questions[0].QuestionText)
	assert.Equal(t, enum.AssessmentTypeMultipleChoice, assessment.Questions[0].QuestionType)
	assert.Equal(t, 2, len(assessment.Questions[0].Options))
	assert.Equal(t, "A", assessment.Questions[0].CorrectAnswer)
	assert.Equal(t, "q2", assessment.Questions[1].ID)
	assert.Equal(t, enum.AssessmentTypeTrueFalse, assessment.Questions[1].QuestionType)
	assert.Equal(t, "2024-01-01T00:00:00Z", assessment.CreatedAt)
}

func TestAssessmentRepository_FindAssessmentByMaterialID_AssessmentNotFound(t *testing.T) {
	// Arrange
	db, cleanup := setupMongoTestDB(t)
	defer cleanup()

	repo := NewMongoAssessmentRepository(db)

	materialID, err := valueobject.MaterialIDFromString("550e8400-e29b-41d4-a716-446655440003")
	require.NoError(t, err)

	// No insertar ningún assessment

	// Act
	assessment, err := repo.FindAssessmentByMaterialID(context.Background(), materialID)

	// Assert
	require.NoError(t, err, "FindAssessmentByMaterialID should not return error when assessment not found")
	assert.Nil(t, assessment, "Assessment should be nil when not found")
}

func TestAssessmentRepository_SaveResult_ValidData(t *testing.T) {
	// Arrange
	db, cleanup := setupMongoTestDB(t)
	defer cleanup()

	repo := NewMongoAssessmentRepository(db)

	userID, err := valueobject.UserIDFromString("660e8400-e29b-41d4-a716-446655440000")
	require.NoError(t, err)

	result := &repository.AssessmentResult{
		AssessmentID:   "assessment-123",
		UserID:         userID,
		Score:          85.5,
		TotalQuestions: 10,
		CorrectAnswers: 8,
		Feedback: []repository.FeedbackItem{
			{
				QuestionID:    "q1",
				IsCorrect:     true,
				UserAnswer:    "A",
				CorrectAnswer: "A",
				Explanation:   "Correcto",
			},
			{
				QuestionID:    "q2",
				IsCorrect:     false,
				UserAnswer:    "B",
				CorrectAnswer: "C",
				Explanation:   "La respuesta correcta es C",
			},
		},
		SubmittedAt: time.Now().Format(time.RFC3339),
	}

	// Act
	err = repo.SaveResult(context.Background(), result)

	// Assert
	require.NoError(t, err, "SaveResult should not return error with valid data")

	// Verificar que se guardó en la colección
	collection := db.Collection("assessment_results")
	var doc bson.M
	err = collection.FindOne(context.Background(), bson.M{
		"assessment_id": "assessment-123",
		"user_id":       userID.String(),
	}).Decode(&doc)
	require.NoError(t, err, "Result should be in database")
	assert.Equal(t, "assessment-123", doc["assessment_id"])
	assert.Equal(t, userID.String(), doc["user_id"])
	assert.Equal(t, 85.5, doc["score"])
	assert.Equal(t, int32(10), doc["total_questions"])
	assert.Equal(t, int32(8), doc["correct_answers"])
	assert.NotNil(t, doc["feedback"])
	assert.NotNil(t, doc["submitted_at"])
}

func TestAssessmentRepository_SaveResult_DuplicateKey(t *testing.T) {
	// Arrange
	db, cleanup := setupMongoTestDB(t)
	defer cleanup()

	repo := NewMongoAssessmentRepository(db)

	userID, err := valueobject.UserIDFromString("660e8400-e29b-41d4-a716-446655440001")
	require.NoError(t, err)

	result1 := &repository.AssessmentResult{
		AssessmentID:   "assessment-456",
		UserID:         userID,
		Score:          75.0,
		TotalQuestions: 10,
		CorrectAnswers: 7,
		Feedback:       []repository.FeedbackItem{},
		SubmittedAt:    time.Now().Format(time.RFC3339),
	}

	// Guardar primer resultado
	err = repo.SaveResult(context.Background(), result1)
	require.NoError(t, err)

	// Intentar guardar segundo resultado para el mismo assessment y usuario
	result2 := &repository.AssessmentResult{
		AssessmentID:   "assessment-456",
		UserID:         userID,
		Score:          90.0,
		TotalQuestions: 10,
		CorrectAnswers: 9,
		Feedback:       []repository.FeedbackItem{},
		SubmittedAt:    time.Now().Format(time.RFC3339),
	}

	// Act
	err = repo.SaveResult(context.Background(), result2)

	// Assert
	require.Error(t, err, "SaveResult should return error on duplicate key")
	assert.Contains(t, err.Error(), "DuplicateKey", "Error should indicate duplicate key")

	// Verificar que solo hay un documento
	collection := db.Collection("assessment_results")
	count, err := collection.CountDocuments(context.Background(), bson.M{
		"assessment_id": "assessment-456",
		"user_id":       userID.String(),
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "Should have only one result document")

	// Verificar que el score sigue siendo el primero (75.0)
	var doc bson.M
	err = collection.FindOne(context.Background(), bson.M{
		"assessment_id": "assessment-456",
		"user_id":       userID.String(),
	}).Decode(&doc)
	require.NoError(t, err)
	assert.Equal(t, 75.0, doc["score"], "Score should remain the first one saved")
}

func TestAssessmentRepository_CountCompletedAssessments(t *testing.T) {
	// Arrange
	db, cleanup := setupMongoTestDB(t)
	defer cleanup()

	repo := NewMongoAssessmentRepository(db)

	// Insertar varios resultados
	collection := db.Collection("assessment_results")
	results := []interface{}{
		bson.M{
			"assessment_id":   "assessment-1",
			"user_id":         "user-1",
			"score":           80.0,
			"total_questions": 10,
			"correct_answers": 8,
			"feedback":        bson.A{},
			"submitted_at":    time.Now(),
		},
		bson.M{
			"assessment_id":   "assessment-2",
			"user_id":         "user-1",
			"score":           90.0,
			"total_questions": 10,
			"correct_answers": 9,
			"feedback":        bson.A{},
			"submitted_at":    time.Now(),
		},
		bson.M{
			"assessment_id":   "assessment-1",
			"user_id":         "user-2",
			"score":           70.0,
			"total_questions": 10,
			"correct_answers": 7,
			"feedback":        bson.A{},
			"submitted_at":    time.Now(),
		},
	}
	_, err := collection.InsertMany(context.Background(), results)
	require.NoError(t, err)

	// Act
	count, err := repo.CountCompletedAssessments(context.Background())

	// Assert
	require.NoError(t, err, "CountCompletedAssessments should not return error")
	assert.Equal(t, int64(3), count, "Should count all completed assessments")
}

func TestAssessmentRepository_CountCompletedAssessments_EmptyCollection(t *testing.T) {
	// Arrange
	db, cleanup := setupMongoTestDB(t)
	defer cleanup()

	repo := NewMongoAssessmentRepository(db)

	// No insertar ningún resultado

	// Act
	count, err := repo.CountCompletedAssessments(context.Background())

	// Assert
	require.NoError(t, err, "CountCompletedAssessments should not return error on empty collection")
	assert.Equal(t, int64(0), count, "Should return 0 for empty collection")
}

func TestAssessmentRepository_CalculateAverageScore(t *testing.T) {
	// Arrange
	db, cleanup := setupMongoTestDB(t)
	defer cleanup()

	repo := NewMongoAssessmentRepository(db)

	// Insertar varios resultados con diferentes scores
	collection := db.Collection("assessment_results")
	results := []interface{}{
		bson.M{
			"assessment_id":   "assessment-1",
			"user_id":         "user-1",
			"score":           80.0,
			"total_questions": 10,
			"correct_answers": 8,
			"feedback":        bson.A{},
			"submitted_at":    time.Now(),
		},
		bson.M{
			"assessment_id":   "assessment-2",
			"user_id":         "user-1",
			"score":           90.0,
			"total_questions": 10,
			"correct_answers": 9,
			"feedback":        bson.A{},
			"submitted_at":    time.Now(),
		},
		bson.M{
			"assessment_id":   "assessment-1",
			"user_id":         "user-2",
			"score":           70.0,
			"total_questions": 10,
			"correct_answers": 7,
			"feedback":        bson.A{},
			"submitted_at":    time.Now(),
		},
	}
	_, err := collection.InsertMany(context.Background(), results)
	require.NoError(t, err)

	// Act
	avgScore, err := repo.CalculateAverageScore(context.Background())

	// Assert
	require.NoError(t, err, "CalculateAverageScore should not return error")
	expectedAvg := (80.0 + 90.0 + 70.0) / 3.0
	assert.InDelta(t, expectedAvg, avgScore, 0.01, "Should calculate correct average score")
}

func TestAssessmentRepository_CalculateAverageScore_EmptyCollection(t *testing.T) {
	// Arrange
	db, cleanup := setupMongoTestDB(t)
	defer cleanup()

	repo := NewMongoAssessmentRepository(db)

	// No insertar ningún resultado

	// Act
	avgScore, err := repo.CalculateAverageScore(context.Background())

	// Assert
	require.NoError(t, err, "CalculateAverageScore should not return error on empty collection")
	assert.Equal(t, 0.0, avgScore, "Should return 0.0 for empty collection")
}
