//go:build integration
// +build integration

package repository_test

import (
	"context"
	"database/sql"
)

// createAssessmentTables crea las tablas necesarias para los tests de assessment
// Función compartida entre todos los tests de integración de PostgreSQL
func createAssessmentTables(db *sql.DB) error {
	ctx := context.Background()

	// Primero eliminar tablas existentes para evitar conflictos de schema
	dropSQL := `
		DROP TABLE IF EXISTS assessment_attempt_answer CASCADE;
		DROP TABLE IF EXISTS assessment_attempt CASCADE;
		DROP TABLE IF EXISTS assessment CASCADE;
	`

	_, err := db.ExecContext(ctx, dropSQL)
	if err != nil {
		return err
	}

	// Luego crear las tablas con el schema correcto
	schema := `
		CREATE TABLE assessment (
			id UUID PRIMARY KEY,
			material_id UUID NOT NULL,
			mongo_document_id VARCHAR(24) NOT NULL,
			title VARCHAR(255) NOT NULL,
			total_questions INTEGER NOT NULL,
			pass_threshold INTEGER NOT NULL,
			max_attempts INTEGER DEFAULT NULL,
			time_limit_minutes INTEGER DEFAULT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE INDEX idx_assessment_material_id ON assessment(material_id);

		CREATE TABLE assessment_attempt (
			id UUID PRIMARY KEY,
			assessment_id UUID NOT NULL,
			student_id UUID NOT NULL,
			score INTEGER NOT NULL,
			max_score INTEGER NOT NULL,
			time_spent_seconds INTEGER NOT NULL,
			started_at TIMESTAMP NOT NULL,
			completed_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL,
			idempotency_key VARCHAR(255) DEFAULT NULL
		);

		CREATE INDEX idx_attempt_student_assessment ON assessment_attempt(student_id, assessment_id);

		CREATE TABLE assessment_attempt_answer (
			id UUID NOT NULL,
			attempt_id UUID NOT NULL,
			question_id VARCHAR(100) NOT NULL,
			selected_answer_id VARCHAR(10) NOT NULL,
			is_correct BOOLEAN NOT NULL,
			time_spent_seconds INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL,
			PRIMARY KEY (id)
		);

		CREATE INDEX idx_answer_attempt_id ON assessment_attempt_answer(attempt_id);
		CREATE INDEX idx_answer_question_id ON assessment_attempt_answer(question_id);
	`

	_, err = db.ExecContext(ctx, schema)
	return err
}
