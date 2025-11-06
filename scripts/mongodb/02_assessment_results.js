// =====================================================
// Script: 02_assessment_results.js
// Descripción: Colección de resultados de evaluaciones completadas
// Autor: EduGo Team
// Fecha: 2025-11-05
// Versión: 1.0.0
// =====================================================

// Conectar a la base de datos (ajustar según ambiente)
// db = db.getSiblingDB('edugo_mobile');

print("Creando colección assessment_results...");

// =====================================================
// Crear Colección con Validación de Schema
// =====================================================

db.createCollection("assessment_results", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["assessment_id", "user_id", "score", "total_questions", "correct_answers", "submitted_at"],
            properties: {
                assessment_id: {
                    bsonType: "string",
                    description: "ID del assessment (UUID). Required."
                },
                user_id: {
                    bsonType: "string",
                    description: "ID del usuario que completó el assessment (UUID). Required."
                },
                score: {
                    bsonType: "double",
                    minimum: 0,
                    maximum: 100,
                    description: "Puntaje obtenido (0-100). Required."
                },
                total_questions: {
                    bsonType: "int",
                    minimum: 1,
                    description: "Total de preguntas en el assessment. Required."
                },
                correct_answers: {
                    bsonType: "int",
                    minimum: 0,
                    description: "Total de respuestas correctas. Required."
                },
                feedback: {
                    bsonType: "array",
                    description: "Array de feedback detallado por pregunta. Optional.",
                    items: {
                        bsonType: "object",
                        required: ["question_id", "is_correct"],
                        properties: {
                            question_id: {
                                bsonType: "string",
                                description: "ID de la pregunta"
                            },
                            is_correct: {
                                bsonType: "bool",
                                description: "Si la respuesta fue correcta"
                            },
                            user_answer: {
                                bsonType: "string",
                                description: "Respuesta enviada por el usuario"
                            },
                            correct_answer: {
                                bsonType: "string",
                                description: "Respuesta correcta"
                            },
                            explanation: {
                                bsonType: "string",
                                description: "Explicación contextual del resultado"
                            }
                        }
                    }
                },
                submitted_at: {
                    bsonType: "date",
                    description: "Fecha y hora de envío del assessment. Required."
                }
            }
        }
    }
});

print("Colección assessment_results creada exitosamente.");

// =====================================================
// Crear Índices de Performance
// =====================================================

print("Creando índices de performance...");

// Índice UNIQUE compuesto en (assessment_id, user_id)
// Previene que un usuario complete el mismo assessment múltiples veces
db.assessment_results.createIndex(
    { "assessment_id": 1, "user_id": 1 },
    {
        unique: true,
        name: "idx_assessment_user_unique",
        background: true
    }
);

print("✓ Índice UNIQUE creado: idx_assessment_user_unique");

// Índice en submitted_at para ordenar resultados por fecha
// Usado en queries de "evaluaciones recientes"
db.assessment_results.createIndex(
    { "submitted_at": -1 },
    {
        name: "idx_submitted_at_desc",
        background: true
    }
);

print("✓ Índice creado: idx_submitted_at_desc");

// Índice compuesto en (user_id, submitted_at)
// Optimiza queries de historial de evaluaciones por usuario
db.assessment_results.createIndex(
    { "user_id": 1, "submitted_at": -1 },
    {
        name: "idx_user_submitted",
        background: true
    }
);

print("✓ Índice creado: idx_user_submitted");

// Índice en score para estadísticas y filtros por puntaje
db.assessment_results.createIndex(
    { "score": 1 },
    {
        name: "idx_score",
        background: true
    }
);

print("✓ Índice creado: idx_score");

// =====================================================
// Verificar Índices Creados
// =====================================================

print("\nÍndices creados en assessment_results:");
db.assessment_results.getIndexes().forEach(function(index) {
    print("  - " + index.name + ": " + JSON.stringify(index.key));
});

// =====================================================
// Insertar Documento de Ejemplo (Opcional - solo para testing)
// =====================================================

// Descomentar para insertar un documento de ejemplo:
/*
db.assessment_results.insertOne({
    assessment_id: "550e8400-e29b-41d4-a716-446655440001",
    user_id: "550e8400-e29b-41d4-a716-446655440002",
    score: 85.5,
    total_questions: 10,
    correct_answers: 8,
    feedback: [
        {
            question_id: "q1",
            is_correct: true,
            user_answer: "B",
            correct_answer: "B",
            explanation: "Correcto. La opción B es la respuesta adecuada porque..."
        },
        {
            question_id: "q2",
            is_correct: false,
            user_answer: "A",
            correct_answer: "C",
            explanation: "Incorrecto. La respuesta correcta es C porque..."
        }
    ],
    submitted_at: new Date()
});

print("\nDocumento de ejemplo insertado.");
*/

// =====================================================
// Estadísticas de la Colección
// =====================================================

print("\nEstadísticas de la colección assessment_results:");
const stats = db.assessment_results.stats();
print("  - Total documentos: " + stats.count);
print("  - Tamaño total: " + (stats.size / 1024).toFixed(2) + " KB");
print("  - Índices: " + stats.nindexes);

print("\n✅ Script completado exitosamente.");
