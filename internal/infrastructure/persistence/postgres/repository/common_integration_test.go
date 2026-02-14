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

	// Luego crear tablas aisladas para tests de repositorio.
	// Se evita acoplar estos tests al schema completo de infrastructure y sus FKs.
	schema := `
		CREATE TABLE assessment (
			id UUID PRIMARY KEY,
			material_id UUID NOT NULL,
			mongo_document_id VARCHAR(24) NOT NULL,
			questions_count INTEGER NOT NULL DEFAULT 0,
			total_questions INTEGER,
			title VARCHAR(255),
			pass_threshold INTEGER,
			max_attempts INTEGER DEFAULT NULL,
			time_limit_minutes INTEGER DEFAULT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'generated',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE INDEX idx_assessment_material_id ON assessment(material_id);

		CREATE TABLE assessment_attempt (
			id UUID PRIMARY KEY,
			assessment_id UUID NOT NULL,
			student_id UUID NOT NULL,
			score NUMERIC(5,2),
			max_score NUMERIC(5,2),
			percentage NUMERIC(5,2),
			status VARCHAR(50) NOT NULL DEFAULT 'completed',
			time_spent_seconds INTEGER,
			idempotency_key VARCHAR(64),
			started_at TIMESTAMP NOT NULL,
			completed_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL
		);

		CREATE INDEX idx_attempt_student_assessment ON assessment_attempt(student_id, assessment_id);

		CREATE TABLE assessment_attempt_answer (
			id UUID NOT NULL,
			attempt_id UUID NOT NULL,
			question_index INTEGER NOT NULL,
			student_answer TEXT,
			is_correct BOOLEAN,
			points_earned NUMERIC(5,2),
			max_points NUMERIC(5,2),
			time_spent_seconds INTEGER,
			answered_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			PRIMARY KEY (id)
		);

		CREATE INDEX idx_answer_attempt_id ON assessment_attempt_answer(attempt_id);
		CREATE INDEX idx_answer_question_index ON assessment_attempt_answer(question_index);
	`

	_, err = db.ExecContext(ctx, schema)
	return err
}
