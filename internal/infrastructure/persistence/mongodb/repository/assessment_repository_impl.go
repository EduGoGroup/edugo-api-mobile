package repository

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoAssessmentRepository struct {
	assessments *mongo.Collection
	attempts    *mongo.Collection
	results     *mongo.Collection
}

func NewMongoAssessmentRepository(db *mongo.Database) repository.AssessmentRepository {
	return &mongoAssessmentRepository{
		assessments: db.Collection("material_assessments"),
		attempts:    db.Collection("assessment_attempts"),
		results:     db.Collection("assessment_results"),
	}
}

func (r *mongoAssessmentRepository) SaveAssessment(ctx context.Context, assessment *repository.MaterialAssessment) error {
	doc := bson.M{
		"material_id": assessment.MaterialID.String(),
		"questions":   assessment.Questions,
		"created_at":  time.Now(),
	}

	filter := bson.M{"material_id": assessment.MaterialID.String()}
	opts := options.Replace().SetUpsert(true)
	_, err := r.assessments.ReplaceOne(ctx, filter, doc, opts)
	return err
}

func (r *mongoAssessmentRepository) FindAssessmentByMaterialID(ctx context.Context, materialID valueobject.MaterialID) (*repository.MaterialAssessment, error) {
	filter := bson.M{"material_id": materialID.String()}

	var doc bson.M
	err := r.assessments.FindOne(ctx, filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	assessment := &repository.MaterialAssessment{
		MaterialID: materialID,
		CreatedAt:  doc["created_at"].(string),
	}

	// Parsear questions si existen
	if questionsData, ok := doc["questions"].(bson.A); ok {
		questions := make([]repository.AssessmentQuestion, 0, len(questionsData))
		for _, qData := range questionsData {
			qMap := qData.(bson.M)
			
			// Parsear options si existen
			var options []string
			if opts, ok := qMap["options"].(bson.A); ok {
				for _, opt := range opts {
					options = append(options, opt.(string))
				}
			}
			
			question := repository.AssessmentQuestion{
				ID:            getString(qMap, "id"),
				QuestionText:  getString(qMap, "text"),
				QuestionType:  enum.AssessmentType(getString(qMap, "question_type")),
				Options:       options,
				CorrectAnswer: qMap["answer"],
			}
			questions = append(questions, question)
		}
		assessment.Questions = questions
	}

	return assessment, nil
}

// getString es un helper para extraer strings de bson.M de forma segura
func getString(m bson.M, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (r *mongoAssessmentRepository) SaveAttempt(ctx context.Context, attempt *repository.AssessmentAttempt) error {
	doc := bson.M{
		"id":           attempt.ID,
		"material_id":  attempt.MaterialID.String(),
		"user_id":      attempt.UserID.String(),
		"answers":      attempt.Answers,
		"score":        attempt.Score,
		"attempted_at": time.Now(),
	}

	_, err := r.attempts.InsertOne(ctx, doc)
	return err
}

func (r *mongoAssessmentRepository) FindAttemptsByUser(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) ([]*repository.AssessmentAttempt, error) {
	filter := bson.M{
		"material_id": materialID.String(),
		"user_id":     userID.String(),
	}

	cursor, err := r.attempts.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "attempted_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var attempts []*repository.AssessmentAttempt
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}

		attempt := &repository.AssessmentAttempt{
			ID:          doc["id"].(string),
			MaterialID:  materialID,
			UserID:      userID,
			Score:       doc["score"].(float64),
			AttemptedAt: doc["attempted_at"].(string),
		}
		attempts = append(attempts, attempt)
	}

	return attempts, nil
}

func (r *mongoAssessmentRepository) GetBestAttempt(ctx context.Context, materialID valueobject.MaterialID, userID valueobject.UserID) (*repository.AssessmentAttempt, error) {
	filter := bson.M{
		"material_id": materialID.String(),
		"user_id":     userID.String(),
	}

	var doc bson.M
	err := r.attempts.FindOne(ctx, filter, options.FindOne().SetSort(bson.D{{Key: "score", Value: -1}})).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &repository.AssessmentAttempt{
		ID:          doc["id"].(string),
		MaterialID:  materialID,
		UserID:      userID,
		Score:       doc["score"].(float64),
		AttemptedAt: doc["attempted_at"].(string),
	}, nil
}

func (r *mongoAssessmentRepository) SaveResult(ctx context.Context, result *repository.AssessmentResult) error {
	// Convertir feedback items a formato BSON
	feedbackDocs := make([]bson.M, 0, len(result.Feedback))
	for _, item := range result.Feedback {
		feedbackDocs = append(feedbackDocs, bson.M{
			"question_id":    item.QuestionID,
			"is_correct":     item.IsCorrect,
			"user_answer":    item.UserAnswer,
			"correct_answer": item.CorrectAnswer,
			"explanation":    item.Explanation,
		})
	}

	// Crear documento para insertar
	doc := bson.M{
		"assessment_id":   result.AssessmentID,
		"user_id":         result.UserID.String(),
		"score":           result.Score,
		"total_questions": result.TotalQuestions,
		"correct_answers": result.CorrectAnswers,
		"feedback":        feedbackDocs,
		"submitted_at":    time.Now(),
	}

	// InsertOne - el índice UNIQUE en (assessment_id, user_id) previene duplicados
	_, err := r.results.InsertOne(ctx, doc)
	if err != nil {
		// Verificar si es error de duplicado (código 11000)
		if mongo.IsDuplicateKeyError(err) {
			return mongo.CommandError{Code: 11000, Name: "DuplicateKey", Message: "assessment already completed by user"}
		}
		return err
	}

	return nil
}

// CountCompletedAssessments cuenta el total de evaluaciones completadas
// Usado para estadísticas globales del sistema
func (r *mongoAssessmentRepository) CountCompletedAssessments(ctx context.Context) (int64, error) {
	count, err := r.results.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CalculateAverageScore calcula el promedio de puntajes de todas las evaluaciones completadas
// Usa pipeline de agregación de MongoDB para calcular el promedio
func (r *mongoAssessmentRepository) CalculateAverageScore(ctx context.Context) (float64, error) {
	// Pipeline: { $group: { _id: null, avgScore: { $avg: "$score" } } }
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "avgScore", Value: bson.D{{Key: "$avg", Value: "$score"}}},
		}}},
	}

	cursor, err := r.results.Aggregate(ctx, pipeline)
	if err != nil {
		return 0.0, err
	}
	defer cursor.Close(ctx)

	// Leer resultado
	if cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return 0.0, err
		}

		// Extraer avgScore (puede ser nil si no hay documentos)
		if avgScore, ok := result["avgScore"].(float64); ok {
			return avgScore, nil
		}
	}

	// Si no hay resultados, retornar 0.0 (no hay evaluaciones completadas)
	return 0.0, nil
}
