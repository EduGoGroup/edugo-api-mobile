package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AssessmentDocumentRepository maneja documentos de evaluación en MongoDB
// Estos documentos contienen las preguntas, opciones y respuestas correctas
type AssessmentDocumentRepository interface {
	// FindByMaterialID busca un documento de assessment por material_id
	FindByMaterialID(ctx context.Context, materialID string) (*AssessmentDocument, error)

	// FindByID busca un documento por ObjectID de MongoDB
	FindByID(ctx context.Context, objectID string) (*AssessmentDocument, error)

	// Save guarda o actualiza un documento de assessment (upsert)
	Save(ctx context.Context, doc *AssessmentDocument) error

	// Delete elimina un documento de assessment
	Delete(ctx context.Context, objectID string) error
}

// AssessmentDocument representa el schema MongoDB de material_assessment_worker
type AssessmentDocument struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	MaterialID       string             `bson:"material_id"`
	Title            string             `bson:"title"`
	Questions        []Question         `bson:"questions"`
	TotalQuestions   int                `bson:"total_questions"`
	TotalPoints      int                `bson:"total_points"`
	AIModel          string             `bson:"ai_model"`
	ProcessingTimeMs int                `bson:"processing_time_ms"`
	Metadata         Metadata           `bson:"metadata"`
	Version          int                `bson:"version"`
	CreatedAt        time.Time          `bson:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at"`
}

// Question representa una pregunta del assessment
type Question struct {
	ID            string   `bson:"id"`
	Text          string   `bson:"text"`
	Type          string   `bson:"type"` // "multiple_choice", "true_false", "short_answer"
	Options       []Option `bson:"options"`
	CorrectAnswer string   `bson:"correct_answer"`
	Feedback      Feedback `bson:"feedback"`
	Difficulty    string   `bson:"difficulty,omitempty"` // Post-MVP
	Tags          []string `bson:"tags,omitempty"`       // Post-MVP
}

// Option representa una opción de respuesta
type Option struct {
	ID   string `bson:"id"`
	Text string `bson:"text"`
}

// Feedback contiene el feedback educativo
type Feedback struct {
	Correct   string `bson:"correct"`
	Incorrect string `bson:"incorrect"`
}

// Metadata contiene información de generación
type Metadata struct {
	GeneratedBy          string    `bson:"generated_by"`
	GenerationDate       time.Time `bson:"generation_date"`
	PromptVersion        string    `bson:"prompt_version"`
	TotalQuestions       int       `bson:"total_questions"`
	EstimatedTimeMinutes int       `bson:"estimated_time_minutes"`
}

// MongoAssessmentDocumentRepository implementa AssessmentDocumentRepository
type MongoAssessmentDocumentRepository struct {
	collection *mongo.Collection
}

// NewMongoAssessmentDocumentRepository crea una nueva instancia del repositorio
func NewMongoAssessmentDocumentRepository(db *mongo.Database) AssessmentDocumentRepository {
	return &MongoAssessmentDocumentRepository{
		collection: db.Collection("material_assessment_worker"),
	}
}

// FindByMaterialID busca un documento de assessment por material_id
func (r *MongoAssessmentDocumentRepository) FindByMaterialID(ctx context.Context, materialID string) (*AssessmentDocument, error) {
	if materialID == "" {
		return nil, fmt.Errorf("mongo: material_id cannot be empty")
	}

	filter := bson.M{"material_id": materialID}

	var doc AssessmentDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil // No encontrado
	}
	if err != nil {
		return nil, fmt.Errorf("mongo: error finding assessment document: %w", err)
	}

	return &doc, nil
}

// FindByID busca un documento por ObjectID de MongoDB
func (r *MongoAssessmentDocumentRepository) FindByID(ctx context.Context, objectID string) (*AssessmentDocument, error) {
	if objectID == "" {
		return nil, fmt.Errorf("mongo: objectID cannot be empty")
	}

	// Validar que sea un ObjectID válido
	oid, err := primitive.ObjectIDFromHex(objectID)
	if err != nil {
		return nil, fmt.Errorf("mongo: invalid ObjectID: %w", err)
	}

	filter := bson.M{"_id": oid}

	var doc AssessmentDocument
	err = r.collection.FindOne(ctx, filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("mongo: error finding assessment document: %w", err)
	}

	return &doc, nil
}

// Save guarda o actualiza un documento de assessment (upsert)
func (r *MongoAssessmentDocumentRepository) Save(ctx context.Context, doc *AssessmentDocument) error {
	if doc == nil {
		return fmt.Errorf("mongo: document cannot be nil")
	}

	if doc.MaterialID == "" {
		return fmt.Errorf("mongo: material_id is required")
	}

	if doc.Title == "" {
		return fmt.Errorf("mongo: title is required")
	}

	if len(doc.Questions) == 0 {
		return fmt.Errorf("mongo: at least one question is required")
	}

	// Actualizar timestamp
	doc.UpdatedAt = time.Now().UTC()

	// Buscar documento existente por material_id para preservar el _id
	filter := bson.M{"material_id": doc.MaterialID}
	existing, err := r.FindByMaterialID(ctx, doc.MaterialID)

	if err == nil && existing != nil {
		// Documento existe: preservar _id y created_at
		doc.ID = existing.ID
		doc.CreatedAt = existing.CreatedAt
	} else {
		// Documento nuevo: generar _id y created_at
		if doc.ID.IsZero() {
			doc.ID = primitive.NewObjectID()
		}
		if doc.CreatedAt.IsZero() {
			doc.CreatedAt = time.Now().UTC()
		}
	}

	// Upsert por material_id
	opts := options.Replace().SetUpsert(true)

	_, err = r.collection.ReplaceOne(ctx, filter, doc, opts)
	if err != nil {
		return fmt.Errorf("mongo: error saving assessment document: %w", err)
	}

	return nil
}

// Delete elimina un documento de assessment
func (r *MongoAssessmentDocumentRepository) Delete(ctx context.Context, objectID string) error {
	if objectID == "" {
		return fmt.Errorf("mongo: objectID cannot be empty")
	}

	// Validar ObjectID
	oid, err := primitive.ObjectIDFromHex(objectID)
	if err != nil {
		return fmt.Errorf("mongo: invalid ObjectID: %w", err)
	}

	filter := bson.M{"_id": oid}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("mongo: error deleting assessment document: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("mongo: document not found")
	}

	return nil
}

// GetQuestionByID busca una pregunta específica dentro de un documento
// Helper method útil para recuperar feedback durante corrección
func (r *MongoAssessmentDocumentRepository) GetQuestionByID(ctx context.Context, materialID, questionID string) (*Question, error) {
	if materialID == "" || questionID == "" {
		return nil, fmt.Errorf("mongo: materialID and questionID are required")
	}

	filter := bson.M{
		"material_id":  materialID,
		"questions.id": questionID,
	}

	projection := bson.M{
		"questions": bson.M{
			"$elemMatch": bson.M{"id": questionID},
		},
	}

	var result struct {
		Questions []Question `bson:"questions"`
	}

	err := r.collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err == mongo.ErrNoDocuments || len(result.Questions) == 0 {
		return nil, fmt.Errorf("mongo: question not found")
	}
	if err != nil {
		return nil, fmt.Errorf("mongo: error finding question: %w", err)
	}

	return &result.Questions[0], nil
}
