package repository

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoAssessmentRepository struct {
	assessments *mongo.Collection
	attempts    *mongo.Collection
}

func NewMongoAssessmentRepository(db *mongo.Database) repository.AssessmentRepository {
	return &mongoAssessmentRepository{
		assessments: db.Collection("material_assessments"),
		attempts:    db.Collection("assessment_attempts"),
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

	return assessment, nil
}

func (r *mongoAssessmentRepository) SaveAttempt(ctx context.Context, attempt *repository.AssessmentAttempt) error {
	doc := bson.M{
		"id":          attempt.ID,
		"material_id": attempt.MaterialID.String(),
		"user_id":     attempt.UserID.String(),
		"answers":     attempt.Answers,
		"score":       attempt.Score,
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
			ID:         doc["id"].(string),
			MaterialID: materialID,
			UserID:     userID,
			Score:      doc["score"].(float64),
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
		ID:         doc["id"].(string),
		MaterialID: materialID,
		UserID:     userID,
		Score:      doc["score"].(float64),
		AttemptedAt: doc["attempted_at"].(string),
	}, nil
}
