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

type mongoSummaryRepository struct {
	collection *mongo.Collection
}

func NewMongoSummaryRepository(db *mongo.Database) repository.SummaryRepository {
	return &mongoSummaryRepository{
		collection: db.Collection("material_summaries"),
	}
}

func (r *mongoSummaryRepository) Save(ctx context.Context, summary *repository.MaterialSummary) error {
	doc := bson.M{
		"material_id":  summary.MaterialID.String(),
		"main_ideas":   summary.MainIdeas,
		"key_concepts": summary.KeyConcepts,
		"sections":     summary.Sections,
		"glossary":     summary.Glossary,
		"created_at":   time.Now(),
	}

	filter := bson.M{"material_id": summary.MaterialID.String()}
	opts := options.Replace().SetUpsert(true)
	_, err := r.collection.ReplaceOne(ctx, filter, doc, opts)
	return err
}

func (r *mongoSummaryRepository) FindByMaterialID(ctx context.Context, materialID valueobject.MaterialID) (*repository.MaterialSummary, error) {
	filter := bson.M{"material_id": materialID.String()}

	var doc bson.M
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	summary := &repository.MaterialSummary{
		MaterialID: materialID,
		CreatedAt:  doc["created_at"].(string),
	}

	if mainIdeas, ok := doc["main_ideas"].([]interface{}); ok {
		for _, idea := range mainIdeas {
			summary.MainIdeas = append(summary.MainIdeas, idea.(string))
		}
	}

	if keyConcepts, ok := doc["key_concepts"].(map[string]interface{}); ok {
		summary.KeyConcepts = make(map[string]string)
		for k, v := range keyConcepts {
			summary.KeyConcepts[k] = v.(string)
		}
	}

	return summary, nil
}

func (r *mongoSummaryRepository) Delete(ctx context.Context, materialID valueobject.MaterialID) error {
	filter := bson.M{"material_id": materialID.String()}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *mongoSummaryRepository) Exists(ctx context.Context, materialID valueobject.MaterialID) (bool, error) {
	filter := bson.M{"material_id": materialID.String()}
	count, err := r.collection.CountDocuments(ctx, filter)
	return count > 0, err
}
